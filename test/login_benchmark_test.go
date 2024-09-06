package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/internal"
	"github.com/nrmnqdds/gomaluum-api/services/auth"
	"github.com/nrmnqdds/gomaluum-api/services/scraper"
)

func BenchmarkProfileScraper(b *testing.B) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Create a new Echo context with a mock request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/api/v1/profile")

	c.SetCookie(&http.Cookie{
		Name:  "MOD_AUTH_CAS",
		Value: "ba163a937130a173ae89c1fffe6f2ef3",
	})

	// Check cookie
	cookie, err := c.Cookie("MOD_AUTH_CAS")
	if err != nil {

		// Login user
		user := dtos.LoginDTO{
			Username: internal.GetEnv("LOGIN_USERNAME"),
			Password: internal.GetEnv("LOGIN_PASSWORD"),
		}

		data, err := auth.LoginUser(&user)
		if err != nil {
			log.Println("error: ", err)
			b.Errorf("Login failed: %v", err)
		}

		c.SetCookie(&http.Cookie{
			Name:  "MOD_AUTH_CAS",
			Value: data.Token,
		})
	}
  log.Println(cookie)


	for i := 0; i < b.N; i++ {
		// Run the scraper
		profile, err := scraper.ProfileScraper(c)

		// Check for errors
		if err != nil && err != dtos.ErrUnauthorized {
			b.Errorf("Unexpected error: %v", err)
		}

		// Check the result to ensure the scraper ran
		if profile == nil {
			b.Errorf("ProfileScraper returned nil")
		}
	}
}
