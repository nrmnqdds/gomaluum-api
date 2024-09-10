package dtos

type WeekTime struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Day   uint8  `json:"day"`
}

type Subject struct {
	Id         string     `json:"id"`
	CourseCode string     `json:"course_code"`
	CourseName string     `json:"course_name"`
	Section    uint8      `json:"section"`
	Chr        uint8      `json:"chr"`
	Timestamps []WeekTime `json:"timestamps"`
	Venue      string     `json:"venue"`
	Lecturer   string     `json:"lecturer"`
}

type ScheduleResponse struct {
	Id           string    `json:"id"`
	SessionName  string    `json:"session_name"`
	SessionQuery string    `json:"session_query"`
	Schedule     []Subject `json:"schedule"`
}
