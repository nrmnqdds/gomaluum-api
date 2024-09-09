package scraper

import (
	_"regexp"
  _"strings"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/labstack/echo/v4"
	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/internal"
)

var result []dtos.ResultResponse

func ResultScraper(e echo.Context) ([]dtos.ResultResponse, *dtos.CustomError) {
	c := colly.NewCollector()

	var wg sync.WaitGroup
	resultChan := make(chan []dtos.ResultResponse, 100)

	cookie, err := e.Cookie("MOD_AUTH_CAS")
	if err != nil {
		return nil, dtos.ErrUnauthorized
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+cookie.Value)
		r.Headers.Set("User-Agent", internal.RandomString())
	})

	c.OnHTML(".box.box-primary .box-header.with-border .dropdown ul.dropdown-menu li[style*='font-size:16px']", func(e *colly.HTMLElement) {
		e.ForEach("a", func(i int, element *colly.HTMLElement) {
			wg.Add(1)
			go getResultFromSession(c, element.Attr("href"), element.Text, cookie.Value, &wg, resultChan)
		})
	})

	if err := c.Visit(internal.IMALUUM_RESULT_PAGE); err != nil {
		return nil, dtos.ErrFailedToGoToURL
	}

	wg.Wait()
	close(resultChan)
	return result, nil
}

func getResultFromSession(c *colly.Collector, sessionQuery string, sessionName string, cookieValue string, wg *sync.WaitGroup, ch chan<- []dtos.ResultResponse) *dtos.CustomError {
	defer wg.Done()

	url := internal.IMALUUM_RESULT_PAGE + sessionQuery

	// subjects := []dtos.Result{}

	logger := internal.NewLogger()

	c.OnHTML(".box-body table.table.table-hover", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, element *colly.HTMLElement) {
			tds := element.ChildTexts("td")

			if len(tds) == 0 {
				// Skip the first row
				return
			}

      logger.Info("tds: ", tds)
      // logger.Info("tds[7]: ", tds[7])

			// TODO
			// Takbayar yuran
			// if len(tds) >= 4 {
			// 	courseCode := strings.TrimSpace(tds[0])
			// 	if strings.Split(courseCode, "  ")[0] == "Total Credit Points" {
			// 		return
			// 	}
			//
			// 	courseName := strings.TrimSpace(tds[1])
			// 	courseGrade := strings.TrimSpace(tds[2])
			// 	courseCredit := strings.TrimSpace(tds[3])
			//
			// 	subjects = append(subjects, dtos.Result{
			// 		CourseCode:   courseCode,
			// 		CourseName:   courseName,
			// 		CourseGrade:  courseGrade,
			// 		CourseCredit: courseCredit,
			// 	})
			//
			// }

			// re := regexp.MustCompile(`/\s{2,}/`)
			// neutralized := re.Split(strings.TrimSpace(tds[1]), -1)

			// neutralized := strings.TrimSpace(tds[1])

			// if len(neutralized) == 0 {
			// 	return
			// }

			// logger.Info("neutralized: ", neutralized)

			// gpaValue := neutralized[2]
			// status := neutralized[3]
			// remarks := neutralized[4]
			//
			// logger.Info("GPA Value: ", gpaValue)
			// logger.Info("Status: ", status)
			// logger.Info("Remarks: ", remarks)
			//
			//       const neutralized1 = tds[1].textContent.trim().split(/\s{2,}/) || [];
			// const gpaValue = neutralized1[2];
			// const status = neutralized1[3];
			// const remarks = neutralized1[4];
			//
			// const neutralized2 = tds[3].textContent.trim().split(/\s{2,}/) || [];
			// const cgpaValue = neutralized2[2];
			//
			// // Remove the last row
			// rows.pop();
		})
	})

	if err := c.Visit(url); err != nil {
		return dtos.ErrFailedToGoToURL
	}

	ch <- result

	return nil
}
