package requests

import (
	"fmt"

	"github.com/ordo_meritum/shared/utils/formatters"
)

// EducationInfoPayload represents a single educational entry in a request payload.
type EducationInfoPayload struct {
	CourseWork *string  `json:"coursework,omitempty"`
	Degree     string   `json:"degree"`
	Location   string   `json:"location"`
	School     string   `json:"school"`
	StartEnd   string   `json:"start_end"`
	GPA        *float64 `json:"gpa,omitempty"`
	Honors     *string  `json:"honors,omitempty"`
}

// FormatForLLM converts the EducationInfoPayload struct into a plain text format
// suitable for inclusion in a prompt for a large language model.
func (e *EducationInfoPayload) FormatForLLM() string {
	return fmt.Sprintf(`
School: %s
Degree: %s
Location: %s
Dates: %s

Coursework:
%s
	`,
		formatters.PtrString(&e.School, "Not specified"),
		formatters.PtrString(&e.Degree, "Not specified"),
		formatters.PtrString(&e.Location, "Not specified"),
		formatters.PtrString(&e.StartEnd, "Not specified"),
		formatters.PtrString(e.CourseWork, "No coursework listed"),
	)
}
