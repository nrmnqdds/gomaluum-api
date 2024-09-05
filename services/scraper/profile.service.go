package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/labstack/echo/v4"
	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/internal"
)

type Profile struct {
  ImageURL string `json:"image_url"`
  Name     string `json:"name"`
  MatricNo string `json:"matric_no"`
}

func ProfileScraper(e echo.Context) (*Profile, *dtos.CustomError) {
	c := colly.NewCollector()

	cookie, err := e.Cookie("MOD_AUTH_CAS")

	profile := Profile{}

	if err != nil {
		return nil, dtos.ErrUnauthorized
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+cookie.Value)
		r.Headers.Set("User-Agent", internal.RandomString())
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		_name := e.ChildText(".row .col-md-12 .box.box-default .panel-body.row .col-md-4[style='text-align:center; padding:10px; floaf:left;'] h4[style='margin-top:1%;']")
		profile.Name = strings.TrimSpace(_name)

		_matricNo := e.ChildText(".row .col-md-12 .box.box-default .panel-body.row .col-md-4[style='margin-top:3%;'] h4")
		profile.MatricNo = strings.TrimSpace(strings.Split(strings.TrimSpace(_matricNo), "|")[0])
	})

	c.Visit(internal.IMALUUM_PROFILE_PAGE)

	profile.ImageURL = fmt.Sprintf("https://smartcard.iium.edu.my/packages/card/printing/camera/uploads/original/%s.jpeg", profile.MatricNo)

	return &profile, nil
}
