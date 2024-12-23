package dtos

type ResultResponse struct {
	ID           string `json:"id"`
	SessionName  string `json:"session_name"`
	SessionQuery string `json:"session_query"`
	GpaValue     string `json:"gpa_value"`
	CgpaValue    string `json:"cgpa_value"`
	Status       string `json:"status"`
	// Remarks      string   `json:"remarks"`
	CreditHours string   `json:"credit_hours"`
	Result      []Result `json:"result"`
}

type Result struct {
	ID           string `json:"id"`
	CourseCode   string `json:"course_code"`
	CourseName   string `json:"course_name"`
	CourseGrade  string `json:"course_grade"`
	CourseCredit string `json:"course_credit"`
}
