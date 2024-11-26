package scraper

import (
	"sort"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/lucsky/cuid"
	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/helpers"
)

func ResultScraper(d *dtos.ScheduleRequestProps) (*[]dtos.ResultResponse, *dtos.CustomError) {
	logger, _ := helpers.NewLogger()

	logger.Info("Result service called")

	var (
		e              = d.Echo
		c              = colly.NewCollector()
		_cookie        string
		wg             sync.WaitGroup
		result         []dtos.ResultResponse
		sessionQueries []string
		sessionNames   []string
	)

	if d.Token == "" {
		cookie, err := e.Cookie("MOD_AUTH_CAS")
		if err != nil {
			logger.Error("No cookie found!")
			return nil, dtos.ErrUnauthorized
		}

		logger.Debug("Found cookie")
		_cookie = cookie.Value
	} else {
		logger.Debug("Using token from login directly as cookie")
		_cookie = d.Token
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+_cookie)
		r.Headers.Set("User-Agent", helpers.RandomString())
	})

	resultChan := make(chan dtos.ResultResponse)

	c.OnHTML(".box.box-primary .box-header.with-border .dropdown ul.dropdown-menu", func(e *colly.HTMLElement) {
		sessionQueries = e.ChildAttrs("li[style*='font-size:16px'] a", "href")
		sessionNames = e.ChildTexts("li[style*='font-size:16px'] a")
	})

	if err := c.Visit(helpers.ImaluumResultPage); err != nil {
		logger.Error("Failed to go to URL")
		return nil, dtos.ErrFailedToGoToURL
	}

	// Filter out unwanted session
	filteredQueries := make([]string, 0)
	filteredNames := make([]string, 0)
	for i := range sessionQueries {
		if sessionQueries[i] != "?ses=1111/1111&sem=1" && sessionQueries[i] != "?ses=0000/0000&sem=0" {
			filteredQueries = append(filteredQueries, sessionQueries[i])
			filteredNames = append(filteredNames, sessionNames[i])
		}
	}
	sessionQueries = filteredQueries
	sessionNames = filteredNames

	for i := range sessionQueries {
		wg.Add(1)

		clone := c.Clone()

		go getResultFromSession(clone, &_cookie, &sessionQueries[i], &sessionNames[i], resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for s := range resultChan {
		result = append(result, s)
	}

	if len(result) == 0 {
		logger.Error("Schedule is empty")
		return nil, dtos.ErrFailedToScrape
	}

	sort.Slice(result, func(i, j int) bool {
		return helpers.CompareSessionNames(result[i].SessionName, result[j].SessionName)
	})

	return &result, nil
}

func getResultFromSession(c *colly.Collector, cookie *string, sessionQuery *string, sessionName *string, resultChan chan<- dtos.ResultResponse, wg *sync.WaitGroup) {
	defer wg.Done()

	logger, _ := helpers.NewLogger()
	logger.Debugf("Running scraper for session: %v", *sessionName)

	url := helpers.ImaluumResultPage + *sessionQuery

	var (
		subjects     []dtos.Result
		mu           sync.Mutex
		courseCode   string
		courseName   string
		courseGrade  string
		courseCredit string
		gpa          string
		cgpa         string
		chr          string
		status       string
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+*cookie)
		r.Headers.Set("User-Agent", helpers.RandomString())
	})

	c.OnHTML(".box-body table.table.table-hover tbody tr", func(e *colly.HTMLElement) {
		tds := e.ChildTexts("td")

		courseCode = strings.TrimSpace(tds[0])

		courseName = strings.TrimSpace(tds[1])

		courseGrade = strings.TrimSpace(tds[2])

		courseCredit = strings.TrimSpace(tds[3])

		words := strings.Fields(courseCode)
		if words[0] == "Total" {
			logger.Debug("cgpa usecase found:")
			logger.Debugf("tds[0]: %v", tds[0])
			logger.Debugf("tds[1]: %v", tds[1])
			logger.Debugf("tds[2]: %v", tds[2])
			logger.Debugf("tds[3]: %v", tds[3])

			gpaWord := strings.Fields(courseName)
			chr = strings.TrimSpace(gpaWord[1])
			gpa = strings.TrimSpace(gpaWord[2])
			status = strings.TrimSpace(gpaWord[3])

			cgpaWord := strings.Fields(courseCredit)
			cgpa = strings.TrimSpace(cgpaWord[2])
			return
		}

		mu.Lock()
		subjects = append(subjects, dtos.Result{
			CourseCode:   courseCode,
			CourseName:   courseName,
			CourseGrade:  courseGrade,
			CourseCredit: courseCredit,
		})
		mu.Unlock()
	})

	if err := c.Visit(url); err != nil {
		return
	}

	resultChan <- dtos.ResultResponse{
		ID:           cuid.Slug(),
		SessionName:  *sessionName,
		SessionQuery: *sessionQuery,
		GpaValue:     gpa,
		CgpaValue:    cgpa,
		CreditHours:  chr,
		Status:       status,
		Result:       subjects,
	}
}
