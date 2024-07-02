package scraper

import (
	"fmt"
	"github.com/gocolly/colly/v2"
)

const IMALUUM_HOME_PAGE = "https://imaluum.iium.edu.my/home"

func ScheduleScraperService(cookie string) {
	fmt.Println("Running ScheduleScraperService")
	c := colly.NewCollector()

	// Find and visit all links
	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// 	e.Request.Visit(e.Attr("href"))
	// })
	//
	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// })
	//
	// c.Visit("http://go-colly.org/")

	c.Visit(IMALUUM_HOME_PAGE)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("Cookie", "MOD_AUTH_CAS="+cookie)
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
	})
}
