package scraper

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/labstack/echo/v4"
	"github.com/lucsky/cuid"
	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/internal"
)

func ScheduleScraper(e echo.Context) ([]dtos.ScheduleResponse, *dtos.CustomError) {
	c := colly.NewCollector(
		colly.Async(true),
	)

	// c.Limit(&colly.LimitRule{
	// 	DomainGlob: "*",
	// 	Parallelism: 2,
	// })

	var wg sync.WaitGroup
	scheduleChan := make(chan dtos.ScheduleResponse)
	schedule := []dtos.ScheduleResponse{}
	var mu sync.Mutex

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

		sessionList := []Session{}

		e.ForEach("a", func(i int, element *colly.HTMLElement) {
			sessionList = append(sessionList, Session{
				sessionName:  element.Text,
				sessionQuery: element.Attr("href"),
			})
		})

		for _, session := range sessionList {
			wg.Add(1)

			go func(session Session) {
				// defer wg.Done()

				mu.Lock()

				_schedule, err := getScheduleFromSession(c, session.sessionQuery, session.sessionName, cookie.Value)
				if err != nil {
					return
				}

				if len(_schedule) != 0 {
					unmarshaled, err := json.Marshal(_schedule)
					if err != nil {
						log.Println("failed to marshal schedule")
					}
					log.Println("_schedule: ", string(unmarshaled))
					scheduleChan <- dtos.ScheduleResponse{
						SessionName:  session.sessionName,
						SessionQuery: session.sessionQuery,
						Schedule:     _schedule,
					}

					mu.Unlock()
					wg.Done()
					return
				}
			}(session)
		}
	})

	if err := c.Visit(internal.IMALUUM_SCHEDULE_PAGE); err != nil {
		return nil, dtos.ErrFailedToGoToURL
	}

	// Goroutine to collect the results from the channel
	go func() {
		for resp := range scheduleChan {
			unmarshaled, err := json.Marshal(resp)
			if err != nil {
				log.Println("failed to marshal schedule")
			}
			log.Println("resp: ", string(unmarshaled))
			schedule = append(schedule, resp)
		}
	}()

	wg.Wait()
	close(scheduleChan)

	c.Wait()

	log.Println("returned schedule: ", schedule)

	return schedule, nil
}

func getScheduleFromSession(c *colly.Collector, sessionQuery string, sessionName string, cookieValue string) ([]dtos.Subject, *dtos.CustomError) {
	url := internal.IMALUUM_SCHEDULE_PAGE + sessionQuery

	schedule := []dtos.Subject{}

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

				schedule = append(schedule, dtos.Subject{
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
				courseCode := schedule[len(schedule)-1].CourseCode
				courseName := schedule[len(schedule)-1].CourseName
				section := schedule[len(schedule)-1].Section
				chr := schedule[len(schedule)-1].Chr

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

				schedule = append(schedule, dtos.Subject{
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
		return nil, dtos.ErrFailedToGoToURL
	}

	c.Wait()

	return schedule, nil
}
