package scraper

import (
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/lucsky/cuid"
	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/helpers"
	"github.com/rung/go-safecast"
)

func ScheduleScraper(d *dtos.ScheduleRequestProps) (*[]dtos.ScheduleResponse, *dtos.CustomError) {
	logger, _ := helpers.NewLogger()

	logger.Info("Schedule service called")

	var (
		e              = d.Echo
		c              = colly.NewCollector()
		_cookie        string
		wg             sync.WaitGroup
		schedule       []dtos.ScheduleResponse
		sessionQueries []string
		sessionNames   []string
	)

	// If token from POST request is empty, check existing cookie
	if d.Token == "" {
		cookie, err := e.Cookie("MOD_AUTH_CAS")
		if err != nil {
			logger.Error("No cookie found!")
			return nil, dtos.ErrUnauthorized
		}

		logger.Debug("Found cookie")
		_cookie = cookie.Value
	} else {
		// If token from POST request is not empty, use it as cookie
		logger.Debug("Using token from login directly as cookie")
		_cookie = d.Token
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+_cookie)
		r.Headers.Set("User-Agent", helpers.RandomString())
	})

	scheduleChan := make(chan dtos.ScheduleResponse)

	c.OnHTML(".box.box-primary .box-header.with-border .dropdown ul.dropdown-menu", func(e *colly.HTMLElement) {
		sessionQueries = e.ChildAttrs("li[style*='font-size:16px'] a", "href")
		sessionNames = e.ChildTexts("li[style*='font-size:16px'] a")
	})

	if err := c.Visit(helpers.ImaluumSchedulePage); err != nil {
		logger.Error("Failed to go to URL")
		return nil, dtos.ErrFailedToGoToURL
	}

	for i := range sessionQueries {
		wg.Add(1)

		clone := c.Clone()

		go getScheduleFromSession(clone, _cookie, sessionQueries[i], sessionNames[i], scheduleChan, &wg)
	}

	go func() {
		wg.Wait()
		close(scheduleChan)
	}()

	for s := range scheduleChan {
		schedule = append(schedule, s)
	}

	if len(schedule) == 0 {
		logger.Error("Schedule is empty")
		return nil, dtos.ErrFailedToScrape
	}

	sort.Slice(schedule, func(i, j int) bool {
		return helpers.CompareSessionNames(schedule[i].SessionName, schedule[j].SessionName)
	})

	return &schedule, nil
}

func getScheduleFromSession(c *colly.Collector, cookie string, sessionQuery string, sessionName string, scheduleChan chan<- dtos.ScheduleResponse, wg *sync.WaitGroup) {
	defer wg.Done()

	logger, _ := helpers.NewLogger()
	logger.Debugf("Running scraper for session: %v", sessionName)

	url := helpers.ImaluumSchedulePage + sessionQuery

	var mu sync.Mutex
	subjects := []dtos.Subject{}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+cookie)
		r.Headers.Set("User-Agent", helpers.RandomString())
	})

	c.OnHTML(".box-body table.table.table-hover tr", func(e *colly.HTMLElement) {
		tds := e.ChildTexts("td")

		weekTime := []dtos.WeekTime{}

		if len(tds) == 0 {
			// Skip the first row
			return
		}

		// Handles for perfect cell
		if len(tds) == 9 {
			courseCode := strings.TrimSpace(tds[0])
			courseName := strings.TrimSpace(tds[1])

			section, err := safecast.Atoi32(strings.TrimSpace(tds[2]))
			if err != nil {
				return
			}

			chr, err := strconv.ParseFloat(strings.TrimSpace(tds[3]), 32)
			if err != nil {
				return
			}

			// Split the days
			_days := strings.Split(strings.Replace(strings.TrimSpace(tds[5]), " ", "", -1), "-")

			// Handles weird ass day format
			switch _days[0] {
			case "MTW":
				_days = []string{"M", "T", "W"}
			case "TWTH":
				_days = []string{"T", "W", "TH"}
			case "MTWTH":
				_days = []string{"M", "T", "W", "TH"}
			case "MTWTHF":
				_days = []string{"M", "T", "W", "TH", "F"}
			}

			for _, day := range _days {
				dayNum := helpers.GetScheduleDays(day)
				timeTemp := tds[6]

				// `timeFullForm` refers to schedule time from iMaluum
				// e.g.: 800-920 or 1000-1120
				timeFullForm := strings.Replace(strings.TrimSpace(timeTemp), " ", "", -1)

				// in some cases, iMaluum will return "-" as time
				// if `timeFullForm` equals `TimeSeparator`, then we skip this row
				if timeFullForm == helpers.TimeSeparator {
					continue
				}

				// safely split time entry
				time := strings.Split(timeFullForm, helpers.TimeSeparator)

				start := strings.TrimSpace(time[0])
				end := strings.TrimSpace(time[1])

				if len(start) == 3 {
					start = "0" + start
				}

				if len(end) == 3 {
					end = "0" + end
				}

				weekTime = append(weekTime, dtos.WeekTime{
					Start: start,
					End:   end,
					Day:   dayNum,
				})
			}

			venue := strings.TrimSpace(tds[7])
			lecturer := strings.TrimSpace(tds[8])

			mu.Lock()
			subjects = append(subjects, dtos.Subject{
				ID:         cuid.New(),
				CourseCode: courseCode,
				CourseName: courseName,
				Section:    uint32(section),
				Chr:        chr,
				Timestamps: weekTime,
				Venue:      venue,
				Lecturer:   lecturer,
			})
			mu.Unlock()

		}

		// Handles for merged cell usually at time or day or venue
		if len(tds) == 4 {
			mu.Lock()
			lastSubject := subjects[len(subjects)-1]
			mu.Unlock()
			courseCode := lastSubject.CourseCode
			courseName := lastSubject.CourseName
			section := lastSubject.Section
			chr := lastSubject.Chr

			// Split the days
			_days := strings.Split(strings.Replace(strings.TrimSpace(tds[0]), " ", "", -1), "-")

			// Handles weird ass day format
			switch _days[0] {
			case "MTW":
				_days = []string{"M", "T", "W"}
			case "TWTH":
				_days = []string{"T", "W", "TH"}
			case "MTWTH":
				_days = []string{"M", "T", "W", "TH"}
			case "MTWTHF":
				_days = []string{"M", "T", "W", "TH", "F"}
			}

			for _, day := range _days {
				dayNum := helpers.GetScheduleDays(day)
				timeTemp := tds[1]
				time := strings.Split(strings.Replace(strings.TrimSpace(timeTemp), " ", "", -1), "-")

				if len(time) != 2 {
					continue
				}

				start := strings.TrimSpace(time[0])
				end := strings.TrimSpace(time[1])

				if len(start) == 3 {
					start = "0" + start
				}
				if len(end) == 3 {
					end = "0" + end
				}

				weekTime = append(weekTime, dtos.WeekTime{
					Start: start,
					End:   end,
					Day:   dayNum,
				})
			}

			venue := strings.TrimSpace(tds[2])
			lecturer := strings.TrimSpace(tds[3])

			mu.Lock()
			subjects = append(subjects, dtos.Subject{
				ID:         cuid.Slug(),
				CourseCode: courseCode,
				CourseName: courseName,
				Section:    section,
				Chr:        chr,
				Timestamps: weekTime,
				Venue:      venue,
				Lecturer:   lecturer,
			})
			mu.Unlock()
		}
	})

	if err := c.Visit(url); err != nil {
		return
	}

	scheduleChan <- dtos.ScheduleResponse{
		ID:           cuid.Slug(),
		SessionName:  sessionName,
		SessionQuery: sessionQuery,
		Schedule:     subjects,
	}
}
