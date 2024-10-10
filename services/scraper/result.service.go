package scraper

import (
	"slices"
	"sort"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/lucsky/cuid"
	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/internal"
	"github.com/sourcegraph/conc/pool"
)

var logger = internal.NewLogger()

func ResultScraper(d *dtos.ScheduleRequestProps) (*[]dtos.ResultResponse, *dtos.CustomError) {
	e := d.Echo

	var (
		c        = colly.NewCollector()
		result   []dtos.ResultResponse
		mu       sync.Mutex
		isLatest = e.QueryParam("latest")
		// semester       = e.QueryParam("semester")
		// year           = e.QueryParam("year")
		sessionQueries = []string{}
		p              = pool.New().WithMaxGoroutines(20)
		_cookie        string
	)

	cookie, err := e.Cookie("MOD_AUTH_CAS")

	if err != nil {
		if d.Token == "" {
			return nil, dtos.ErrUnauthorized
		}

		_cookie = d.Token
	} else {
		_cookie = cookie.Value
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+_cookie)
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
			getResultFromSession(c, &latestSession, &sessionName, &result, &mu)
		})
	})

	if err := c.Visit(internal.IMALUUM_RESULT_PAGE); err != nil {
		return nil, dtos.ErrFailedToGoToURL
	}

	p.Wait()
	c.Wait()

	if len(result) == 0 {
		return nil, dtos.ErrFailedToScrape
	}

	sort.Slice(result, func(i, j int) bool {
		return internal.CompareSessionNames(result[i].SessionName, result[j].SessionName)
	})

	return &result, nil
}

func getResultFromSession(c *colly.Collector, sessionQuery *string, sessionName *string, result *[]dtos.ResultResponse, mu *sync.Mutex) {
	defer mu.Unlock()

	url := internal.IMALUUM_RESULT_PAGE + *sessionQuery

	mu.Lock()

	c.OnHTML(".box-body table.table.table-hover tr", func(e *colly.HTMLElement) {
		tds := e.ChildTexts("td")

		if len(tds) == 0 {
			// Skip the first row
			return
		}
		// grab last td
		// lastTD := tds[len(tds)-1]
		// logger.Infof("Last TD: %s", lastTD)

		logger.Infof("TDs: %v", tds)
	})

	if err := c.Visit(url); err != nil {
		return
	}

	*result = append(*result, dtos.ResultResponse{
		ID:           cuid.Slug(),
		SessionName:  *sessionName,
		SessionQuery: *sessionQuery,
		GpaValue:     "",
		CgpaValue:    "",
		Status:       "",
		Remarks:      "",
		Result:       []dtos.Result{},
	})
}
