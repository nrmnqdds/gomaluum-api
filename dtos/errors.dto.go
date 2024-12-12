package dtos

import (
	"fmt"
	"net/http"
)

type CustomError struct {
	OriginalErr error
	Message     string `json:"message,omitempty"`
	StatusCode  int    `json:"status,omitempty"`
}

// Error returns the error message
func (e *CustomError) Error() string {
	if e.OriginalErr != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.OriginalErr)
	}
	return e.Message
}

func (e *CustomError) CustomError() *CustomError {
	return &CustomError{
		OriginalErr: e.OriginalErr,
		Message:     e.Message,
		StatusCode:  e.StatusCode,
	}
}

// Unwrap returns the original error
func (e *CustomError) Unwrap() error {
	return e.OriginalErr
}

// WrapError wraps an original error with a predefined CustomError
func WrapError(predefError *CustomError, originalErr error) *CustomError {
	return &CustomError{
		OriginalErr: originalErr,
		Message:     predefError.Message,
		StatusCode:  predefError.StatusCode,
	}
}

var (
	ErrInvalidRequestPayload = &CustomError{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid request payload!",
	}

	ErrInternalServerError = &CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal server error!",
	}

	ErrUnauthorized = &CustomError{
		StatusCode: http.StatusUnauthorized,
		Message:    "Please login first!",
	}

	ErrFailedToInitCookieJar = &CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Failed to initialize cookie jar!",
	}

	ErrFailedToScrape = &CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Failed to scrape data!",
	}

	ErrFailedToLogin = &CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Failed to login!",
	}

	ErrSelectorsNotFound = &CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Selectors not found!",
	}

	ErrFailedToGoToURL = &CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Failed to go to URL!",
	}

	ErrFailedToLimit = &CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Failed to set goroutine limit!",
	}

	ErrFailedToInitLogger = &CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Failed to initialize logger!",
	}

	ErrFailedToGetSchedule = &CustomError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Failed to get schedule!",
	}
)
