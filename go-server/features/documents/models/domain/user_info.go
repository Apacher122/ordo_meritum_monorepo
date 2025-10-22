package domain

import (
	"fmt"
	"strings"
)

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

func (e *EducationInfo) FormatForLLM() string {
	var sb strings.Builder

	sb.WriteString("<education_section>\n")

	sb.WriteString(fmt.Sprintf("\t<school>%s</school>\n", strings.TrimSpace(e.School)))
	sb.WriteString(fmt.Sprintf("\t<degree>%s</degree>\n", strings.TrimSpace(e.Degree)))
	sb.WriteString(fmt.Sprintf("\t<dates>%s</dates>\n", strings.TrimSpace(e.StartEnd)))
	sb.WriteString(fmt.Sprintf("\t<location>%s</location>\n", strings.TrimSpace(e.Location)))

	if e.GPA != nil {
		sb.WriteString(fmt.Sprintf("\t<gpa>%.2f</gpa>\n", *e.GPA))
	}

	if e.Honors != nil && *e.Honors != "" {
		sb.WriteString(fmt.Sprintf("\t<honors>%s</honors>\n", strings.TrimSpace(*e.Honors)))
	}

	if e.CourseWork != nil && *e.CourseWork != "" {
		sb.WriteString(fmt.Sprintf("\t<relevant_coursework>%s</relevant_coursework>\n", strings.TrimSpace(*e.CourseWork)))
	}

	sb.WriteString("</education_section>")

	return sb.String()
}

func FormatEducationHistoryForLLM(educations []*EducationInfo) string {
	if len(educations) == 0 {
		return "<education_section></education_section>"
	}

	var sb strings.Builder

	sb.WriteString("<education_section>\n")

	for _, edu := range educations {
		sb.WriteString(edu.FormatForLLM())
	}

	sb.WriteString("</education_section>")

	return sb.String()
}
