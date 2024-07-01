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
	resp_first, _ := client.Get("https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome")
	defer resp_first.Body.Close()
	client.Jar.SetCookies(urlObj, resp_first.Cookies())
	cookies1 := resp_first.Cookies()
	resp, _ := client.PostForm("https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome", formVal)
	defer resp.Body.Close()
	newCook := append(cookies1, resp.Cookies()...)
	client.Jar.SetCookies(urlObj, newCook)
	// return client

	cookies := client.Jar.Cookies(urlObj)

	// filtered := filterMODAuthCAS(client.Jar.Cookies(urlObj)
	filtered := filterMODAuthCAS(cookies)
	// fmt.Println(filtered)fmt.Println(client.Jar.Cookies(urlObj))

	c.JSON(http.StatusOK, gin.H{"cookies": client.Jar.Cookies(urlObj)})
}

func filterMODAuthCAS(arr []string) []string {
	var result []string
	for _, s := range arr {
		parts := strings.Split(s, " ")
		for _, part := range parts {
			if strings.HasPrefix(part, "MOD_AUTH_CAS=") {
				result = append(result, part)
				break
			}
		}
	}
	return result
}
