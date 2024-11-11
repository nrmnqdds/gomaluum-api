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

		profile.Level = strings.TrimSpace(e.ChildText(".row .col-md-12 .box.box-default .panel-body.row .col-md-4[style='text-align:center; padding:10px; floaf:left;'] h4:nth-of-type(2)"))

		profile.Kuliyyah = strings.TrimSpace(e.ChildText(".row .col-md-12 .box.box-default .panel-body.row .col-md-4[style='margin-top:3%;'] p"))

		_ic := strings.TrimSpace(e.ChildText(".row:nth-of-type(2) .col-md-12 .nav-tabs-custom .tab-content .tab-pane.active .row .col-md-3 p:nth-of-type(2)"))
		profile.IC = strings.TrimSpace(strings.Split(_ic, ":")[1])

		_gender := strings.TrimSpace(e.ChildText(".row:nth-of-type(2) .col-md-12 .nav-tabs-custom .tab-content .tab-pane.active .row .col-md-3 p:nth-of-type(3)"))
		profile.Gender = strings.TrimSpace(strings.Split(_gender, ":")[1])

		_birthday := strings.TrimSpace(e.ChildText(".row:nth-of-type(2) .col-md-12 .nav-tabs-custom .tab-content .tab-pane.active .row .col-md-3 p:nth-of-type(4)"))
		profile.Birthday = strings.TrimSpace(strings.Split(_birthday, ":")[1])

		_religion := strings.TrimSpace(e.ChildText(".row:nth-of-type(2) .col-md-12 .nav-tabs-custom .tab-content .tab-pane.active .row .col-md-3 p:nth-of-type(5)"))
		profile.Religion = strings.TrimSpace(strings.Split(_religion, ":")[1])

		_maritalStatus := strings.TrimSpace(e.ChildText(".row:nth-of-type(2) .col-md-12 .nav-tabs-custom .tab-content .tab-pane.active .row .col-md-9 p:nth-of-type(2)"))
		profile.MaritalStatus = strings.TrimSpace(strings.Split(_maritalStatus, ":")[1])

		_address := strings.TrimSpace(e.ChildText(".row:nth-of-type(2) .col-md-12 .nav-tabs-custom .tab-content .tab-pane.active .row .col-md-9 p:nth-of-type(3)"))
		// profile.Address = strings.TrimSpace(strings.Split(_address, ":")[1])
		addressParts := strings.Split(_address, ":")
		if len(addressParts) > 1 {
			profile.Address = formatAddress(addressParts[1])
		} else {
			profile.Address = addressParts[0]
		}
	})

	if err := c.Visit(helpers.ImaluumProfilePage); err != nil {
		return nil, dtos.ErrInternalServerError
	}

	profile.ImageURL = fmt.Sprintf("https://smartcard.iium.edu.my/packages/card/printing/camera/uploads/original/%s.jpeg", profile.MatricNo)

	if validationErr := helpers.Validator.Struct(&profile); validationErr != nil {
		return nil, dtos.ErrFailedToScrape
	}

	return &profile, nil
}

func formatAddress(address string) string {
	// Split the address into lines
	lines := strings.Split(address, "\n")

	// Clean each line
	var formattedLines []string
	for _, line := range lines {
		// Remove tabs, extra spaces, and trim whitespace
		cleaned := strings.TrimSpace(strings.ReplaceAll(line, "\t", ""))
		if cleaned != "" {
			formattedLines = append(formattedLines, cleaned)
		}
	}

	// Join the lines with commas
	return strings.Join(formattedLines, ", ")
}
