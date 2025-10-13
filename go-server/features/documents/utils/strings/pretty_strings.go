package pretty_strings

import (
	"fmt"

	"github.com/ordo_meritum/features/documents/models/requests"
)

func PrettyEducation(e requests.EducationInfoRequest) string {
	return fmt.Sprintf(`
School: %s
Degree: %s
Location: %s
Dates: %s

Coursework:
%s
	`,
		PtrString(&e.School, "Not specified"),
		PtrString(&e.Degree, "Not specified"),
		PtrString(&e.Location, "Not specified"),
		PtrString(&e.StartEnd, "Not specified"),
		PtrString(&e.CourseWork, "No coursework listed"),
	)
}

func PtrString(s *string, fallback string) string {
	if s == nil || *s == "" {
		return fallback
	}
	return *s
}
