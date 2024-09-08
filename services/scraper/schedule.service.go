package scraper

import (
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/labstack/echo/v4"
	"github.com/lucsky/cuid"
	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/internal"
)

var (
	schedule []dtos.ScheduleResponse
	logger   = internal.NewLogger()
)

func ScheduleScraper(e echo.Context) ([]dtos.ScheduleResponse, *dtos.CustomError) {
	c := colly.NewCollector()

	var wg sync.WaitGroup
	scheduleChan := make(chan []dtos.ScheduleResponse, 100)

	cookie, err := e.Cookie("MOD_AUTH_CAS")
	if err != nil {
		return nil, dtos.ErrUnauthorized
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+cookie.Value)
		r.Headers.Set("User-Agent", internal.RandomString())
	})

	c.OnHTML(".box.box-primary .box-header.with-border .dropdown ul.dropdown-menu li[style*='font-size:16px']", func(e *colly.HTMLElement) {
		type Session struct {
			sessionName  string
			sessionQuery string
		}

		e.ForEach("a", func(i int, element *colly.HTMLElement) {
			wg.Add(1)
			go getScheduleFromSession(c, element.Attr("href"), element.Text, cookie.Value, &wg, scheduleChan)
		})
	})

	if err := c.Visit(internal.IMALUUM_SCHEDULE_PAGE); err != nil {
		return nil, dtos.ErrFailedToGoToURL
	}

	go func() {
		wg.Wait()
		close(scheduleChan)

		logger.Info("Closed schedule channel")
	}()

	logger.Info("Returned schedule", schedule)

	return schedule, nil
}

func getScheduleFromSession(c *colly.Collector, sessionQuery string, sessionName string, cookieValue string, wg *sync.WaitGroup, ch chan<- []dtos.ScheduleResponse) *dtos.CustomError {
	defer wg.Done()

	url := internal.IMALUUM_SCHEDULE_PAGE + sessionQuery

	subjects := []dtos.Subject{}

	c.OnHTML(".box-body table.table.table-hover", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, element *colly.HTMLElement) {
			tds := element.ChildTexts("td")

			days := []uint8{}
			weekTime := []dtos.WeekTime{}

			if len(tds) == 0 {
				// Skip the first row
				return
			}

			if len(tds) == 9 {

				courseCode := strings.TrimSpace(tds[0])
				courseName := strings.TrimSpace(tds[1])
				section, err := strconv.Atoi(strings.TrimSpace(tds[2]))
				if err != nil {
					return
				}

				chr, err := strconv.Atoi(strings.TrimSpace(tds[3]))
				if err != nil {
					return
				}

				_days := strings.Split(strings.Replace(strings.TrimSpace(tds[5]), " ", "", -1), "-")

				for _, day := range _days {
					dayNum := internal.GetScheduleDays(day)
					days = append(days, dayNum)
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

				for _, day := range days {

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
						Day:   day,
					})
				}

				venue := strings.TrimSpace(tds[7])
				lecturer := strings.TrimSpace(tds[8])

				subjects = append(subjects, dtos.Subject{
					SessionName: sessionName,
					Id:          cuid.New(),
					CourseCode:  courseCode,
					CourseName:  courseName,
					Section:     uint8(section),
					Chr:         uint8(chr),
					Timestamps:  weekTime,
					Venue:       venue,
					Lecturer:    lecturer,
				})

			}

			if len(tds) == 4 {
				courseCode := subjects[len(subjects)-1].CourseCode
				courseName := subjects[len(subjects)-1].CourseName
				section := subjects[len(subjects)-1].Section
				chr := subjects[len(subjects)-1].Chr

				_days := strings.Split(strings.Replace(strings.TrimSpace(tds[5]), " ", "", -1), "-")

				for _, day := range _days {
					dayNum := internal.GetScheduleDays(day)
					days = append(days, dayNum)
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
					SessionName: sessionName,
					Id:          cuid.Slug(),
					CourseCode:  courseCode,
					CourseName:  courseName,
					Section:     section,
					Chr:         chr,
					Timestamps:  weekTime,
					Venue:       venue,
					Lecturer:    lecturer,
				})
			}
		})
	})

	if err := c.Visit(url); err != nil {
		return dtos.ErrFailedToGoToURL
	}

	schedule = append(schedule, dtos.ScheduleResponse{
		SessionName:  sessionName,
		SessionQuery: sessionQuery,
		Schedule:     subjects,
	})

	logger.Info("Returned schedule", schedule)

	ch <- schedule

	return nil
}
