package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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
	resp_first, _ := client.Get("https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome")
	defer resp_first.Body.Close()
	client.Jar.SetCookies(urlObj, resp_first.Cookies())
	cookies1 := resp_first.Cookies()
	resp, _ := client.PostForm("https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome", formVal)
	defer resp.Body.Close()
	newCook := append(cookies1, resp.Cookies()...)
	client.Jar.SetCookies(urlObj, newCook)

	cookies := client.Jar.Cookies(urlObj)

	for _, cookie := range cookies {
		if cookie.Name == "MOD_AUTH_CAS" {
			c.SetCookie(cookie.Name, cookie.Value, 3600, "/", "localhost", false, true)
			c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
			return
		}
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})

}
