package auth

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/nrmnqdds/gomaluum-api/dtos"
	"github.com/nrmnqdds/gomaluum-api/helpers"
	"github.com/nrmnqdds/gomaluum-api/tasks"
)

func LoginUser(user *dtos.LoginDTO) (*dtos.LoginResponseDTO, *dtos.CustomError) {
	jar, _ := cookiejar.New(nil)

	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     helpers.GetEnv("REDIS_URL"),
		Username: helpers.GetEnv("REDIS_USERNAME"),
		Password: helpers.GetEnv("REDIS_PASSWORD"),
	})
	defer asynqClient.Close()

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

	logger.Debugf("Cookies: %v", cookies)
	for _, cookie := range cookies {
		if cookie.Name == "MOD_AUTH_CAS" {

			task, err := tasks.NewSaveToRedisTask(user.Username, user.Password)
			if err != nil {
				logger.Warnf("Failed to set user password to redis: %v", err)
			}

			info, err := asynqClient.Enqueue(task, asynq.ProcessIn(1*time.Minute))
			if err != nil {
				logger.Warnf("could not schedule task: %v", err)
			}
			logger.Infof("enqueued task: id=%s queue=%s", info.ID, info.Queue)

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
