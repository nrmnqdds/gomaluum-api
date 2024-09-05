package auth

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/nrmnqdds/gomaluum-api/dtos"
)

func LoginUser(user *dtos.LoginDTO) (*dtos.LoginResponseDTO, *dtos.CustomError) {
	formVal := url.Values{
		"username":    {string(user.Username)},
		"password":    {string(user.Password)},
		"execution":   {"e1s1"},
		"_eventId":    {"submit"},
		"geolocation": {""},
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, dtos.ErrFailedToInitCookieJar
	}

	client := &http.Client{
		// Transport: tp,
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
		return nil, dtos.ErrFailedToLogin
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
		return nil, dtos.ErrFailedToLogin
	}
	defer resp.Body.Close()

	cookies := client.Jar.Cookies(urlObj)
	client.Jar.SetCookies(urlObj, cookies)

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
