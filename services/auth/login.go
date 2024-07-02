package auth

import (
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
)

func LoginUser(c *gin.Context, username string, password string) {
	fmt.Println("Running LoginUser")

	jar, err := cookiejar.New(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client := &http.Client{
		Jar:       jar,
		Transport: transport,
		Timeout:   time.Second * 10,
	}

	urlObj, err := url.Parse("https://imaluum.iium.edu.my/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("URL parsed", urlObj)

	g, _ := errgroup.WithContext(c.Request.Context())

	var resp_first *http.Response
	g.Go(func() error {
		var err error
		resp_first, err = client.Get("https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome")
		return err
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp_first.Body.Close()

	fmt.Println("Get request done")

	client.Jar.SetCookies(urlObj, resp_first.Cookies())

	formVal := url.Values{
		"username":    {username},
		"password":    {password},
		"execution":   {"e1s1"},
		"_eventId":    {"submit"},
		"geolocation": {""},
	}

	var resp *http.Response
	g.Go(func() error {
		var err error
		resp, err = client.PostForm("https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome", formVal)
		return err
	})

	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

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
