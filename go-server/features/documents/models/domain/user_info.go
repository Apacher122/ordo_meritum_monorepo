package domain

type UserInfo struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	CurrentLocation string `json:"current_location"`
	Email           string `json:"email"`
	Github          string `json:"github,omitempty"`
	Linkedin        string `json:"linkedin,omitempty"`
	Mobile          string `json:"mobile,omitempty"`
	Summary         string `json:"summary,omitempty"`
}

type EducationInfo struct {
	CourseWork *string  `json:"coursework,omitempty"`
	Degree     string   `json:"degree"`
	Location   string   `json:"location"`
	School     string   `json:"school"`
	StartEnd   string   `json:"start_end"`
	GPA        *float64 `json:"gpa,omitempty"`
	Honors     *string  `json:"honors,omitempty"`
}
