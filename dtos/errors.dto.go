package dtos

type CustomError struct {
	OriginalErr error
	Message     string `json:"message,omitempty"`
	StatusCode  int    `json:"status,omitempty"`
}

var (
	ErrInvalidRequestPayload = &CustomError{
		StatusCode: 400,
		Message:    "Invalid request payload!",
	}

	ErrInternalServerError = &CustomError{
		StatusCode: 500,
		Message:    "Internal server error!",
	}

	ErrUnauthorized = &CustomError{
		StatusCode: 401,
		Message:    "Please login first!",
	}

	ErrFailedToInitCookieJar = &CustomError{
		StatusCode: 500,
		Message:    "Failed to initialize cookie jar!",
	}

	ErrFailedToScrape = &CustomError{
		StatusCode: 500,
		Message:    "Failed to scrape data!",
	}

	ErrFailedToLogin = &CustomError{
		StatusCode: 500,
		Message:    "Failed to login!",
	}

	ErrSelectorsNotFound = &CustomError{
		StatusCode: 500,
		Message:    "Selectors not found!",
	}

	ErrFailedToGoToURL = &CustomError{
		StatusCode: 500,
		Message:    "Failed to go to URL!",
	}

	ErrFailedToLimit = &CustomError{
		StatusCode: 500,
		Message:    "Failed to set goroutine limit!",
	}

	ErrFailedToInitLogger = &CustomError{
		StatusCode: 500,
		Message:    "Failed to initialize logger!",
	}

	ErrFailedToGetSchedule = &CustomError{
		StatusCode: 500,
		Message:    "Failed to get schedule!",
	}
)
