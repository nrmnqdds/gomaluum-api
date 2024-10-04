package auth

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/internal"
)

var (
	client *http.Client
	urlObj *url.URL
)

func LoginUser(user *dtos.LoginDTO) (*dtos.LoginResponseDTO, *dtos.CustomError) {
	jar, _ := cookiejar.New(nil)

	logger := internal.NewLogger()
	client = &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}
	urlObj, err := url.Parse("https://imaluum.iium.edu.my/home")
	if err != nil {
		logger.Errorf("Failed to parse url: %v", err)
		return nil, dtos.ErrFailedToLogin
	}

	formVal := url.Values{
		"username":    {user.Username},
		"password":    {user.Password},
		"execution":   {"e1s1"},
		"_eventId":    {"submit"},
		"geolocation": {""},
	}

	// First request
	req_first, _ := http.NewRequest("GET", "https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome", nil)
	setHeaders(req_first)

	resp_first, err := client.Do(req_first)
	if err != nil {
		logger.Errorf("Failed to login first request: %v", err)
		return nil, dtos.ErrFailedToLogin
	}
	resp_first.Body.Close()

	client.Jar.SetCookies(urlObj, resp_first.Cookies())

	// Second request
	req_second, _ := http.NewRequest("POST", "https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome", strings.NewReader(formVal.Encode()))
	req_second.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	setHeaders(req_second)

	resp, err := client.Do(req_second)
	if err != nil {
		logger.Errorf("Failed to login second request: %v", err)
		return nil, dtos.ErrFailedToLogin
	}
	resp.Body.Close()

	cookies := client.Jar.Cookies(urlObj)

	for _, cookie := range cookies {
		if cookie.Name == "MOD_AUTH_CAS" {
			return &dtos.LoginResponseDTO{
				Username: user.Username,
				Token:    cookie.Value,
			}, nil
		}
	}

	return nil, dtos.ErrFailedToLogin
}

func setHeaders(req *http.Request) {
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("User-Agent", "Mozilla/5.0")
}
