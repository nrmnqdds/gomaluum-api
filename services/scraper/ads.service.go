package scraper

import (
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/lucsky/cuid"
	"github.com/nrmnqdds/gomaluum-api/dtos"
)

func AdsScraper() (*[]dtos.Ads, *dtos.CustomError) {
	c := colly.NewCollector()

	ads := []dtos.Ads{}

	c.OnHTML("div[style*='width:100%; clear:both;height:100px']", func(e *colly.HTMLElement) {
		ads = append(ads, dtos.Ads{
      Title:    strings.TrimSpace(e.ChildText("a")),
			ImageURL: strings.TrimSpace(e.ChildAttr("img", "src")),
			Link:     strings.TrimSpace(e.ChildAttr("a", "href")),
			ID:       cuid.New(),
		})
	})

	if err := c.Visit("https://souq.iium.edu.my/embeded"); err != nil {
		return nil, dtos.ErrInternalServerError
	}

	return &ads, nil
}
