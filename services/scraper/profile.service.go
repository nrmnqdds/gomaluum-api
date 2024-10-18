package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/labstack/echo/v4"
	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/helpers"
)

func ProfileScraper(e echo.Context) (*dtos.Profile, *dtos.CustomError) {
	c := colly.NewCollector()

	profile := dtos.Profile{}

	cookie, err := e.Cookie("MOD_AUTH_CAS")
	if err != nil {
		return nil, dtos.ErrUnauthorized
	}

	if cookie.Value == "" {
		return nil, dtos.ErrUnauthorized
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+cookie.Value)
		r.Headers.Set("User-Agent", helpers.RandomString())
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		profile.Name = strings.TrimSpace(e.ChildText(".row .col-md-12 .box.box-default .panel-body.row .col-md-4[style='text-align:center; padding:10px; floaf:left;'] h4[style='margin-top:1%;']"))

		_matricNo := strings.TrimSpace(e.ChildText(".row .col-md-12 .box.box-default .panel-body.row .col-md-4[style='margin-top:3%;'] h4"))
		profile.MatricNo = strings.TrimSpace(strings.Split(_matricNo, "|")[0])
	})

	if err := c.Visit(helpers.IMALUUM_PROFILE_PAGE); err != nil {
		return nil, dtos.ErrInternalServerError
	}

	profile.ImageURL = fmt.Sprintf("https://smartcard.iium.edu.my/packages/card/printing/camera/uploads/original/%s.jpeg", profile.MatricNo)

	if validationErr := helpers.Validator.Struct(&profile); validationErr != nil {
		return nil, dtos.ErrFailedToScrape
	}

	return &profile, nil
}
