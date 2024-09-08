package auth

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/nrmnqdds/gomaluum-api/dtos"
)

var (
	client *http.Client
	urlObj *url.URL
)

func init() {
	jar, _ := cookiejar.New(nil)
	client = &http.Client{Jar: jar}
	urlObj, _ = url.Parse("https://imaluum.iium.edu.my/")
}

func LoginUser(user *dtos.LoginDTO) (*dtos.LoginResponseDTO, *dtos.CustomError) {
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
		return nil, dtos.ErrFailedToLogin
	}
	resp_first.Body.Close()

	client.Jar.SetCookies(urlObj, resp_first.Cookies())

	// Second request
	req_second, _ := http.NewRequest("POST", "https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome", strings.NewReader(formVal.Encode()))
	req_second.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	setHeaders(req_second)

	resp, err := client.Do(req_second)
	if err != nil {
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

	return nil, dtos.ErrInternalServerError
}

func setHeaders(req *http.Request) {
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("User-Agent", "Mozilla/5.0")
}
