package errors

import (
	"errors"
)

var (
	ErrLoginFailed    = errors.New("username or password is incorrect")
	ErrURLParseFailed = errors.New("failed to parse URL")
)
