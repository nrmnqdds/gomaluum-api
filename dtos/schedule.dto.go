package dtos

import "github.com/labstack/echo/v4"

type WeekTime struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Day   uint8  `json:"day"`
}

type Subject struct {
	Id         string     `json:"id"`
	CourseCode string     `json:"course_code"`
	CourseName string     `json:"course_name"`
	Venue      string     `json:"venue"`
	Lecturer   string     `json:"lecturer"`
	Timestamps []WeekTime `json:"timestamps"`
	Chr        float64    `json:"chr"`
	Section    uint8      `json:"section"`
}

type ScheduleResponse struct {
	Id           string    `json:"id"`
	SessionName  string    `json:"session_name"`
	SessionQuery string    `json:"session_query"`
	Schedule     []Subject `json:"schedule"`
}

type ScheduleRequestProps struct {
	Echo  echo.Context
	Token string
}
