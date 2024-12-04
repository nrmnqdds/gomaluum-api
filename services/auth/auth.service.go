package auth

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/helpers"
)

func LoginUser(user *dtos.LoginDTO) (*dtos.LoginResponseDTO, *dtos.CustomError) {
	jar, _ := cookiejar.New(nil)

	logger, _ := helpers.NewLogger()
	client := &http.Client{
		Jar: jar,
	}
	urlObj, err := url.Parse("https://imaluum.iium.edu.my/")
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
	logger.Debug("Making first request")
	reqFirst, _ := http.NewRequest("GET", "https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome", nil)
	setHeaders(reqFirst)

	respFirst, err := client.Do(reqFirst)
	if err != nil {
		logger.Errorf("Failed to login first request: %v", err)
		return nil, dtos.ErrFailedToLogin
	}
	respFirst.Body.Close()

	client.Jar.SetCookies(urlObj, respFirst.Cookies())

	// Second request
	logger.Debug("Making second request")
	reqSecond, _ := http.NewRequest("POST", "https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome", strings.NewReader(formVal.Encode()))
	reqSecond.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	setHeaders(reqSecond)

	resp, err := client.Do(reqSecond)
	if err != nil {
		logger.Errorf("Failed to login second request: %v", err)
		return nil, dtos.ErrFailedToLogin
	}
	resp.Body.Close()

	cookies := client.Jar.Cookies(urlObj)

	for _, cookie := range cookies {
		if cookie.Name == "MOD_AUTH_CAS" {

			// Construct a new API object using a global API key
			// cloudflareClient, err := cloudflare.New(helpers.GetEnv("CLOUDFLARE_API_KEY"), helpers.GetEnv("CLOUDFLARE_API_EMAIL"))
			// alternatively, you can use a scoped API token
			cloudflareClient, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
			if err != nil {
				logger.Errorf("Error initiating cloudflare client: %s", err.Error())
			}

			// Most API calls require a Context
			ctx := context.Background()

			kvEntryParams := cloudflare.WriteWorkersKVEntryParams{
				NamespaceID: os.Getenv("KV_NAMESPACE_ID"),
				Key:         user.Username,
				Value:       []byte(user.Password),
			}

			kvResourceContainer := &cloudflare.ResourceContainer{
				Level:      "accounts",
				Identifier: os.Getenv("KV_USER_ID"),
				Type:       "account",
			}

			_, cerr := cloudflareClient.WriteWorkersKVEntry(ctx, kvResourceContainer, kvEntryParams)
			if cerr != nil {
				logger.Errorf("Error writing to KV: %s", cerr.Error())
			}
			logger.Info("Successfully wrote to KV")

			return &dtos.LoginResponseDTO{
				Username: user.Username,
				Token:    cookie.Value,
			}, nil
		}
	}

	return nil, dtos.ErrFailedToLogin
}

// Function to set headers for a request.
func setHeaders(req *http.Request) {
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("User-Agent", "Mozilla/5.0")
}
