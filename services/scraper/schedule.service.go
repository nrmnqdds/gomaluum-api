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
	mu       sync.Mutex
)

func ScheduleScraper(e echo.Context) ([]dtos.ScheduleResponse, *dtos.CustomError) {
	c := colly.NewCollector()

	var wg sync.WaitGroup

	cookie, err := e.Cookie("MOD_AUTH_CAS")
	if err != nil {
		return nil, dtos.ErrUnauthorized
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+cookie.Value)
		r.Headers.Set("User-Agent", internal.RandomString())
	})

	// Target the session dropdown
	c.OnHTML(".box.box-primary .box-header.with-border .dropdown ul.dropdown-menu li[style*='font-size:16px']", func(e *colly.HTMLElement) {
		// Get the schedule from the dropdown
		e.ForEach("a", func(i int, element *colly.HTMLElement) {
			wg.Add(1)

			// Spawn a goroutine to get the schedule from the session
			// go getScheduleFromSession(c, element.Attr("href"), element.Text, cookie.Value, &wg, scheduleChan)
			go getScheduleFromSession(c, element.Attr("href"), element.Text, cookie.Value, &wg)
		})
	})

	if err := c.Visit(internal.IMALUUM_SCHEDULE_PAGE); err != nil {
		return nil, dtos.ErrFailedToGoToURL
	}

	// Wait for all waitgroup to finish
	wg.Wait()

	return schedule, nil
}

func getScheduleFromSession(c *colly.Collector, sessionQuery string, sessionName string, cookieValue string, wg *sync.WaitGroup) *dtos.CustomError {
	mu.Lock()
	defer mu.Unlock()
	defer wg.Done()

	url := internal.IMALUUM_SCHEDULE_PAGE + sessionQuery

	subjects := []dtos.Subject{}

	// Target the schedule table
	c.OnHTML(".box-body table.table.table-hover", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, element *colly.HTMLElement) {
			tds := element.ChildTexts("td")

			weekTime := []dtos.WeekTime{}

			// Perfect row
			if len(tds) == 9 {
				courseCode := strings.TrimSpace(tds[0])
				courseName := strings.TrimSpace(tds[1])

				// Convert section string to int
				section, err := strconv.Atoi(strings.TrimSpace(tds[2]))
				if err != nil {
					return
				}

				// Convert credit hours string to int
				chr, err := strconv.Atoi(strings.TrimSpace(tds[3]))
				if err != nil {
					return
				}

				// Split the days e.g. "Mon-Fri" -> ["Mon", "Fri"]
				_days := strings.Split(strings.Replace(strings.TrimSpace(tds[5]), " ", "", -1), "-")

				for _, day := range _days {

					// Get the day number from the day string
					// e.g. "Mon" -> 1
					dayNum := internal.GetScheduleDays(day)

					timeTemp := tds[6]

					// Split the time e.g. "0800 - 1000" -> ["0800", "1000"]
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

			// Row with merged cells
			if len(tds) == 4 {

				// Use the last subject to get the course code, course name, section, and credit hours
				courseCode := subjects[len(subjects)-1].CourseCode
				courseName := subjects[len(subjects)-1].CourseName
				section := subjects[len(subjects)-1].Section
				chr := subjects[len(subjects)-1].Chr

				// Split the days e.g. "Mon-Fri" -> ["Mon", "Fri"]
				_days := strings.Split(strings.Replace(strings.TrimSpace(tds[0]), " ", "", -1), "-")

				for _, day := range _days {

					// Get the day number from the day string
					// e.g. "Mon" -> 1
					dayNum := internal.GetScheduleDays(day)

					timeTemp := tds[1]

					// Split the time e.g. "0800 - 1000" -> ["0800", "1000"]
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

	return nil
}
