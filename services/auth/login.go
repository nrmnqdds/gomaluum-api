package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func LoginUser(c *gin.Context, username string, password string) {
	fmt.Println("Running LoginUser")

	formVal := url.Values{
		"username":    {string(username)},
		"password":    {string(password)},
		"execution":   {"e1s1"},
		"_eventId":    {"submit"},
		"geolocation": {""},
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		// error handling
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client := &http.Client{
		Jar: jar,
	}

	urlObj, _ := url.Parse("https://imaluum.iium.edu.my/")

	// First request
	req_first, _ := http.NewRequest("GET", "https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome", nil)
	req_first.Header.Set("Connection", "Keep-Alive")
	req_first.Header.Set("Accept-Language", "en-US")
	req_first.Header.Set("User-Agent", "Mozilla/5.0")

	resp_first, err := client.Do(req_first)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp_first.Body.Close()

	client.Jar.SetCookies(urlObj, resp_first.Cookies())
	// cookies1 := resp_first.Cookies()

	// Second request
	req_second, _ := http.NewRequest("POST", "https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome", strings.NewReader(formVal.Encode()))
	req_second.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req_second.Header.Set("Connection", "Keep-Alive")
	req_second.Header.Set("Accept-Language", "en-US")
	req_second.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req_second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	cookies := client.Jar.Cookies(urlObj)
	fmt.Println(cookies)
	client.Jar.SetCookies(urlObj, cookies)

	for _, cookie := range cookies {
		if cookie.Name == "MOD_AUTH_CAS" {
			c.SetCookie(cookie.Name, cookie.Value, 3600, "/", "localhost", false, true)
			c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
			return
		}
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})

}
