package pretty_strings

import (
	"fmt"

	"github.com/ordo_meritum/features/documents/models/requests"
	"github.com/ordo_meritum/shared/utils/formatters"
)

const defaultNotSpecified = "Not specified"

func PrettyEducation(e requests.EducationInfoPayload) string {
	return fmt.Sprintf(`
School: %s
Degree: %s
Location: %s
Dates: %s

Coursework:
%s
	`,
		formatters.PtrString(&e.School, defaultNotSpecified),
		formatters.PtrString(&e.Degree, defaultNotSpecified),
		formatters.PtrString(&e.Location, defaultNotSpecified),
		formatters.PtrString(&e.StartEnd, defaultNotSpecified),
		formatters.PtrString(&e.CourseWork, defaultNotSpecified),
	)
}
