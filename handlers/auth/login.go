package auth

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nrmnqdds/gomaluum-api/internal"
)

func LoginUser(c echo.Context) error {

	tp := internal.NewTransport()

	type LoginSchema struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	u := new(LoginSchema)
	if err := c.Bind(u); err != nil {
		return err
	}

	formVal := url.Values{
		"username":    {string(u.Username)},
		"password":    {string(u.Password)},
		"execution":   {"e1s1"},
		"_eventId":    {"submit"},
		"geolocation": {""},
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to initialize cookiejar!")
	}

	client := &http.Client{
		Transport: tp,
		Jar:       jar,
	}

	urlObj, _ := url.Parse("https://imaluum.iium.edu.my/")

	// First request
	req_first, _ := http.NewRequest("GET", "https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome", nil)
	req_first.Header.Set("Connection", "Keep-Alive")
	req_first.Header.Set("Accept-Language", "en-US")
	req_first.Header.Set("User-Agent", "Mozilla/5.0")

	resp_first, err := client.Do(req_first)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to send first request!")
	}
	defer resp_first.Body.Close()

	client.Jar.SetCookies(urlObj, resp_first.Cookies())

	// Second request
	req_second, _ := http.NewRequest("POST", "https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome", strings.NewReader(formVal.Encode()))
	req_second.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req_second.Header.Set("Connection", "Keep-Alive")
	req_second.Header.Set("Accept-Language", "en-US")
	req_second.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req_second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to send second request!")
	}
	defer resp.Body.Close()

	cookies := client.Jar.Cookies(urlObj)
	client.Jar.SetCookies(urlObj, cookies)

	newcookie := new(http.Cookie)

	for _, cookie := range cookies {
		if cookie.Name == "MOD_AUTH_CAS" {
			newcookie.Name = cookie.Name
			newcookie.Value = cookie.Value
			c.SetCookie(newcookie)

			log.Println("====================================")
			log.Println("Duration:", tp.Duration())
			log.Println("Request duration:", tp.ReqDuration())
			log.Println("Connection duration:", tp.ConnDuration())
			return c.JSON(http.StatusOK, "Success")
		}
	}

	return c.JSON(http.StatusInternalServerError, "Failed to login!")
}
