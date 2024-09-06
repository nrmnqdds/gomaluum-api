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

type WeekTime struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Day   uint8  `json:"day"`
}

type Subject struct {
	SessionName string     `json:"session_name"`
	Id          string     `json:"id"`
	CourseCode  string     `json:"course_code"`
	CourseName  string     `json:"course_name"`
	Section     uint8      `json:"section"`
	Chr         uint8      `json:"chr"`
	Timestamps  []WeekTime `json:"timestamps"`
	Venue       string     `json:"venue"`
	Lecturer    string     `json:"lecturer"`
}

type ScheduleResponse struct {
	SessionName  string     `json:"session_name"`
	SessionQuery string     `json:"session_query"`
	Schedule     []*Subject `json:"schedule"`
}

func ScheduleScraper(e echo.Context) ([]*ScheduleResponse, *dtos.CustomError) {
	c := colly.NewCollector()

	cookie, err := e.Cookie("MOD_AUTH_CAS")
	if err != nil {
		return nil, dtos.ErrUnauthorized
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+cookie.Value)
		r.Headers.Set("User-Agent", internal.RandomString())
	})

	var wg sync.WaitGroup
	scheduleChan := make(chan *ScheduleResponse)
	schedule := []*ScheduleResponse{}

	// Goroutine to collect the results from the channel
	go func() {
		for resp := range scheduleChan {
			schedule = append(schedule, resp)
		}
    wg.Done()
	}()

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

		wg.Add(len(sessionList)) // Add to the WaitGroup
		for _, session := range sessionList {
			go func(session Session) {
				_schedule, err := getScheduleFromSession(session.sessionQuery, session.sessionName, cookie.Value)
				if err != nil {
					return
				}

				scheduleChan <- &ScheduleResponse{
					SessionName:  session.sessionName,
					SessionQuery: session.sessionQuery,
					Schedule:     _schedule,
				}
			}(session)
		}
	})

	c.Visit(internal.IMALUUM_SCHEDULE_PAGE)

	wg.Wait()
	close(scheduleChan)

	return schedule, nil
}

func getScheduleFromSession(sessionQuery string, sessionName string, cookieValue string) ([]*Subject, *dtos.CustomError) {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+cookieValue)
		r.Headers.Set("User-Agent", internal.RandomString())
	})

	url := internal.IMALUUM_SCHEDULE_PAGE + sessionQuery

	schedule := []*Subject{}

	c.OnHTML(".box-body table.table.table-hover", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, element *colly.HTMLElement) {
			tds := element.ChildTexts("td")

			days := []uint8{}
			weekTime := []WeekTime{}

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
					if strings.Contains(day, "SUN") {
						days = append(days, 0)
					} else if strings.Contains(day, "MON") || day == "M" {
						days = append(days, 1)
					} else if strings.Contains(day, "TUE") || day == "T" {
						days = append(days, 2)
					} else if strings.Contains(day, "WED") || day == "W" {
						days = append(days, 3)
					} else if strings.Contains(day, "THUR") || day == "TH" {
						days = append(days, 4)
					} else if strings.Contains(day, "FRI") || day == "F" {
						days = append(days, 5)
					} else if strings.Contains(day, "SAT") {
						days = append(days, 6)
					}
				}

				for _, day := range days {

					timeTemp := tds[6]
					time := strings.Split(strings.Replace(strings.TrimSpace(timeTemp), " ", "", -1), "-")

					if len(time) != 2 {
						continue
					}

					start := strings.TrimSpace(time[0])
					end := strings.TrimSpace(time[1])

					weekTime = append(weekTime, WeekTime{
						Start: start,
						End:   end,
						Day:   day,
					})
				}

				venue := strings.TrimSpace(tds[7])
				lecturer := strings.TrimSpace(tds[8])

				schedule = append(schedule, &Subject{
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
					if strings.Contains(day, "SUN") {
						days = append(days, 0)
					} else if strings.Contains(day, "MON") || day == "M" {
						days = append(days, 1)
					} else if strings.Contains(day, "TUE") || day == "T" {
						days = append(days, 2)
					} else if strings.Contains(day, "WED") || day == "W" {
						days = append(days, 3)
					} else if strings.Contains(day, "THUR") || day == "TH" {
						days = append(days, 4)
					} else if strings.Contains(day, "FRI") || day == "F" {
						days = append(days, 5)
					} else if strings.Contains(day, "SAT") {
						days = append(days, 6)
					}
				}

				for _, day := range days {

					timeTemp := tds[6]
					time := strings.Split(strings.Replace(strings.TrimSpace(timeTemp), " ", "", -1), "-")

					if len(time) != 2 {
						continue
					}

					start := strings.TrimSpace(time[0])
					end := strings.TrimSpace(time[1])

					weekTime = append(weekTime, WeekTime{
						Start: start,
						End:   end,
						Day:   day,
					})
				}

				venue := strings.TrimSpace(tds[7])
				lecturer := strings.TrimSpace(tds[8])

				schedule = append(schedule, &Subject{
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

	c.Visit(url)

	return schedule, nil
}
