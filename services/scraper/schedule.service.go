package scraper

import (
	"slices"
	"sort"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/labstack/echo/v4"
	"github.com/lucsky/cuid"
	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/internal"
	"github.com/rung/go-safecast"
	"github.com/sourcegraph/conc/pool"
)

var logger = internal.NewLogger()

func ScheduleScraper(e echo.Context) ([]dtos.ScheduleResponse, *dtos.CustomError) {
	var (
		c        = colly.NewCollector()
		schedule []dtos.ScheduleResponse
		mu       sync.Mutex
		isLatest = e.QueryParam("latest")
		// semester       = e.QueryParam("semester")
		// year           = e.QueryParam("year")
		sessionQueries = []string{}
		p              = pool.New().WithMaxGoroutines(20)
	)

	cookie, err := e.Cookie("MOD_AUTH_CAS")
	if err != nil {
		return nil, dtos.ErrUnauthorized
	}

	if cookie.Value == "" {
		return nil, dtos.ErrUnauthorized
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+cookie.Value)
		r.Headers.Set("User-Agent", internal.RandomString())
	})

	c.OnHTML(".box.box-primary .box-header.with-border .dropdown ul.dropdown-menu li[style*='font-size:16px']", func(e *colly.HTMLElement) {
		if isLatest == "true" {
			if len(sessionQueries) > 0 {
				return
			}
		}

		latestSession := e.ChildAttr("a", "href")

		// Check if the session is already in the list
		if slices.Contains(sessionQueries, latestSession) {
			// If it is, return
			return
		}

		// If it's not, add it to the list
		sessionQueries = append(sessionQueries, latestSession)

		sessionName := e.ChildText("a")

		p.Go(func() {
			getScheduleFromSession(c, &latestSession, &sessionName, &schedule, &mu)
		})
	})

	if err := c.Visit(internal.IMALUUM_SCHEDULE_PAGE); err != nil {
		return nil, dtos.ErrFailedToGoToURL
	}

	p.Wait()
	c.Wait()

	if len(schedule) == 0 {
		return nil, dtos.ErrFailedToScrape
	}

	sort.Slice(schedule, func(i, j int) bool {
		return internal.CompareSessionNames(schedule[i].SessionName, schedule[j].SessionName)
	})

	return schedule, nil
}

func getScheduleFromSession(c *colly.Collector, sessionQuery *string, sessionName *string, schedule *[]dtos.ScheduleResponse, mu *sync.Mutex) {
	defer mu.Unlock()

	url := internal.IMALUUM_SCHEDULE_PAGE + *sessionQuery

	subjects := []dtos.Subject{}

	mu.Lock()

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

			section, err := safecast.Atoi8(strings.TrimSpace(tds[2]))
			if err != nil {
				return
			}

			chr, err := safecast.Atoi8(strings.TrimSpace(tds[3]))
			if err != nil {
				return
			}

			_days := strings.Split(strings.Replace(strings.TrimSpace(tds[5]), " ", "", -1), "-")

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
				dayNum := internal.GetScheduleDays(day)
				timeTemp := tds[6]
				time := strings.Split(strings.Replace(strings.TrimSpace(timeTemp), " ", "", -1), "-")

				if len(time) != 2 {
					continue
				}

				start := strings.TrimSpace(time[0])
				end := strings.TrimSpace(time[1])

				weekTime = append(weekTime, dtos.WeekTime{
					Start: start,
					End:   end,
					Day:   dayNum,
				})
			}

			venue := strings.TrimSpace(tds[7])
			lecturer := strings.TrimSpace(tds[8])

			subjects = append(subjects, dtos.Subject{
				Id:         cuid.New(),
				CourseCode: courseCode,
				CourseName: courseName,
				Section:    uint8(section),
				Chr:        uint8(chr),
				Timestamps: weekTime,
				Venue:      venue,
				Lecturer:   lecturer,
			})

		}

		// Handles for merged cell usually at time or day or venue
		if len(tds) == 4 {
			courseCode := subjects[len(subjects)-1].CourseCode
			courseName := subjects[len(subjects)-1].CourseName
			section := subjects[len(subjects)-1].Section
			chr := subjects[len(subjects)-1].Chr

			_days := strings.Split(strings.Replace(strings.TrimSpace(tds[0]), " ", "", -1), "-")

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
				dayNum := internal.GetScheduleDays(day)
				timeTemp := tds[1]
				time := strings.Split(strings.Replace(strings.TrimSpace(timeTemp), " ", "", -1), "-")

				if len(time) != 2 {
					continue
				}

				start := strings.TrimSpace(time[0])
				end := strings.TrimSpace(time[1])

				weekTime = append(weekTime, dtos.WeekTime{
					Start: start,
					End:   end,
					Day:   dayNum,
				})
			}

			venue := strings.TrimSpace(tds[2])
			lecturer := strings.TrimSpace(tds[3])

			subjects = append(subjects, dtos.Subject{
				Id:         cuid.Slug(),
				CourseCode: courseCode,
				CourseName: courseName,
				Section:    section,
				Chr:        chr,
				Timestamps: weekTime,
				Venue:      venue,
				Lecturer:   lecturer,
			})
		}
	})

	if err := c.Visit(url); err != nil {
		return
	}

	*schedule = append(*schedule, dtos.ScheduleResponse{
		Id:           cuid.Slug(),
		SessionName:  *sessionName,
		SessionQuery: *sessionQuery,
		Schedule:     subjects,
	})
}
