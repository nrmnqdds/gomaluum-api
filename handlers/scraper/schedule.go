package scraper

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/labstack/echo/v4"
	"github.com/lucsky/cuid"
	"github.com/nrmnqdds/gomaluum-api/internal"
	"log"
	"strings"
)

func ScheduleScraper(e echo.Context) error {
	fmt.Println("Running ScheduleScraperService")
	c := colly.NewCollector()

	tp := internal.NewTransport()

	c.WithTransport(tp)

	cookie, err := e.Cookie("MOD_AUTH_CAS")

	if err != nil {
		fmt.Println(err)
		return e.JSON(400, map[string]string{"error": err.Error()})
	}

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+cookie.Value)
		r.Headers.Set("User-Agent", internal.RandomString())
	})

	schedule := []Subject{}

	c.OnHTML(".box.box-primary .box-header.with-border .dropdown ul.dropdown-menu li[style*='font-size:16px']", func(e *colly.HTMLElement) {

		type Session struct {
			sessionName  string
			sessionQuery string
		}

		var sessionList []Session

		e.ForEach("a", func(i int, element *colly.HTMLElement) {
			sessionList = append(sessionList, Session{
				sessionName:  element.Text,
				sessionQuery: element.Attr("href"),
			})
		})

		for _, session := range sessionList {
			schedule := getScheduleFromSession(session.sessionQuery, session.sessionName, cookie.Value)

			schedule = append(schedule, schedule...)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		log.Println("====================================")
		log.Println("Duration:", tp.Duration())
		log.Println("Request duration:", tp.ReqDuration())
		log.Println("Connection duration:", tp.ConnDuration())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("====================================")
		log.Println("Error:", err)

		log.Println("Duration:", tp.Duration())
		log.Println("Request duration:", tp.ReqDuration())
		log.Println("Connection duration:", tp.ConnDuration())

		return
	})

	c.Visit(internal.IMALUUM_SCHEDULE_PAGE)

	return e.JSON(200, schedule)
}

type WeekTime struct {
	start string
	end   string
	day   int32
}

type Subject struct {
	sessionName string
	id          string
	courseCode  string
	courseName  string
	section     string
	chr         string
	timestamps  []WeekTime
	venue       string
	lecturer    string
}

type ScheduleResponse struct {
	SessionName  string    `json:"session_name"`
	SessionQuery string    `json:"session_query"`
	Schedule     []Subject `json:"schedule"`
}

func getScheduleFromSession(sessionQuery string, sessionName string, cookieValue string) []Subject {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+cookieValue)
		r.Headers.Set("User-Agent", internal.RandomString())
	})

	url := internal.IMALUUM_SCHEDULE_PAGE + sessionQuery

	var schedule []Subject

	c.OnHTML(".box-body table.table.table-hover", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, element *colly.HTMLElement) {
			tds := element.ChildTexts("td")

			if len(tds) == 0 {
				return
			}

			if len(tds) == 9 {

				courseCode := strings.TrimSpace(tds[0])
				courseName := strings.TrimSpace(tds[1])
				section := strings.TrimSpace(tds[2])
				chr := strings.TrimSpace(tds[3])

				_days := strings.Split(strings.Replace(strings.TrimSpace(tds[5]), " ", "", -1), "-")

				var days []int32

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

				var weekTime []WeekTime

				for _, day := range days {

					timeTemp := tds[6]
					time := strings.Split(strings.Replace(strings.TrimSpace(timeTemp), " ", "", -1), "-")

					if len(time) != 2 {
						continue
					}

					start := strings.TrimSpace(time[0])
					end := strings.TrimSpace(time[1])

					weekTime = append(weekTime, WeekTime{
						start: start,
						end:   end,
						day:   day,
					})
				}

				venue := strings.TrimSpace(tds[7])
				lecturer := strings.TrimSpace(tds[8])

				schedule = append(schedule, Subject{
					sessionName: sessionName,
					id:          cuid.New(),
					courseCode:  courseCode,
					courseName:  courseName,
					section:     section,
					chr:         chr,
					timestamps:  weekTime,
					venue:       venue,
					lecturer:    lecturer,
				})
			}

			if len(tds) == 4 {
				courseCode := schedule[len(schedule)-1].courseCode
				courseName := schedule[len(schedule)-1].courseName
				section := schedule[len(schedule)-1].section
				chr := schedule[len(schedule)-1].chr

				_days := strings.Split(strings.Replace(strings.TrimSpace(tds[5]), " ", "", -1), "-")

				var days []int32

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

				var weekTime []WeekTime

				for _, day := range days {

					timeTemp := tds[6]
					time := strings.Split(strings.Replace(strings.TrimSpace(timeTemp), " ", "", -1), "-")

					if len(time) != 2 {
						continue
					}

					start := strings.TrimSpace(time[0])
					end := strings.TrimSpace(time[1])

					weekTime = append(weekTime, WeekTime{
						start: start,
						end:   end,
						day:   day,
					})
				}

				venue := strings.TrimSpace(tds[7])
				lecturer := strings.TrimSpace(tds[8])

				schedule = append(schedule, Subject{
					sessionName: sessionName,
					id:          cuid.Slug(),
					courseCode:  courseCode,
					courseName:  courseName,
					section:     section,
					chr:         chr,
					timestamps:  weekTime,
					venue:       venue,
					lecturer:    lecturer,
				})

			}
		})
	})

	c.Visit(url)

	return schedule

}
