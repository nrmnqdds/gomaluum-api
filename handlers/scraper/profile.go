package scraper

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/labstack/echo/v4"
	"github.com/nrmnqdds/gomaluum-api/internal"
)

type Profile struct {
	imageURL string
	name     string
	matricNo string
}

func ProfileScraper(e echo.Context) error {
	c := colly.NewCollector()

	tp := internal.NewTransport()

	c.WithTransport(tp)

	cookie, err := e.Cookie("MOD_AUTH_CAS")

	profile := Profile{}

	if err != nil {
		log.Println(err)
		return e.JSON(400, map[string]string{"error": err.Error()})
	}

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+cookie.Value)
		r.Headers.Set("User-Agent", internal.RandomString())
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		_name := e.ChildText(".row .col-md-12 .box.box-default .panel-body.row .col-md-4[style='text-align:center; padding:10px; floaf:left;'] h4[style='margin-top:1%;']")
		profile.name = strings.TrimSpace(_name)

		_matricNo := e.ChildText(".row .col-md-12 .box.box-default .panel-body.row .col-md-4[style='margin-top:3%;'] h4")
		profile.matricNo = strings.TrimSpace(strings.Split(strings.TrimSpace(_matricNo), "|")[0])
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

	c.Visit(internal.IMALUUM_PROFILE_PAGE)

	profile.imageURL = fmt.Sprintf("https://smartcard.iium.edu.my/packages/card/printing/camera/uploads/original/%s.jpeg", profile.matricNo)

	return e.JSON(200, map[string]string{
		"imageURL": profile.imageURL,
		"name":     profile.name,
		"matricNo": profile.matricNo,
	})
}
