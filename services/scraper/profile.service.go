package scraper

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/labstack/echo/v4"
	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/internal"
)

var (
	collectorPool sync.Pool
	once          sync.Once
)

func initCollectorPool() {
	collectorPool = sync.Pool{
		New: func() interface{} {
			return colly.NewCollector(colly.Async(true))
		},
	}
}

func ProfileScraper(e echo.Context) (*dtos.Profile, *dtos.CustomError) {
	once.Do(initCollectorPool)

	cookie, err := e.Cookie("MOD_AUTH_CAS")
	if err != nil {
		return nil, dtos.ErrUnauthorized
	}

	c := collectorPool.Get().(*colly.Collector)
	defer collectorPool.Put(c)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+cookie.Value)
		r.Headers.Set("User-Agent", internal.RandomString())
	})

	var profile dtos.Profile
	var wg sync.WaitGroup
	wg.Add(1)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		defer wg.Done()
		profile.Name = strings.TrimSpace(e.ChildText(".row .col-md-12 .box.box-default .panel-body.row .col-md-4[style='text-align:center; padding:10px; floaf:left;'] h4[style='margin-top:1%;']"))

		_matricNo := strings.TrimSpace(e.ChildText(".row .col-md-12 .box.box-default .panel-body.row .col-md-4[style='margin-top:3%;'] h4"))
		profile.MatricNo = strings.TrimSpace(strings.Split(_matricNo, "|")[0])
	})

	if err := c.Visit(internal.IMALUUM_PROFILE_PAGE); err != nil {
		return nil, dtos.ErrInternalServerError
	}

	wg.Wait()

	profile.ImageURL = fmt.Sprintf("https://smartcard.iium.edu.my/packages/card/printing/camera/uploads/original/%s.jpeg", profile.MatricNo)

	c.Wait()

	return &profile, nil
}
