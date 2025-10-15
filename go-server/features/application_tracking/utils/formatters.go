package utils

import (
	"fmt"

	request "github.com/ordo_meritum/features/application_tracking/models/requests"
)

func FormatJobPostingRequest(jp *request.JobPostingRequest) string {
	return fmt.Sprintf(`
Company: %s
Position: %s
URL: %s
Number of Applicants: %s
Post Age: %s

Job Description:
%s
	`,
		jp.CompanyName,
		jp.JobTitle,
		jp.Link,
		jp.ApplicantCount,
		jp.TimeAgo,
		jp.JobDescription,
	)
}
