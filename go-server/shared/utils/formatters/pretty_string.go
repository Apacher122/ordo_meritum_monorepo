package formatters

import (
	"fmt"
	"strings"

	"github.com/ordo_meritum/database/jobs"
	"github.com/ordo_meritum/features/documents/models/requests"
)

func PrettyJobPost(job jobs.FullJobPosting) string {
	return fmt.Sprintf(`
Job Title: %s
Company: %s
Salary Range: %s
Years of Experience: %s
Education Level: %s

Description:
%s

Company Culture:
%s

Company Values:
%s

Required Tools: %s
Programming Languages: %s
Frameworks & Libraries: %s
Databases: %s
Cloud Technologies: %s
Industry Keywords: %s
Soft Skills: %s
Certifications: %s

Requirements:
%s

Nice to Have:
%s

Applicant Count: %d
	`,
		job.JobTitle,
		job.CompanyName,
		PtrString(job.SalaryRange, "Not specified"),
		PtrString(job.YearsOfExp, "Not specified"),
		PtrString(job.EducationLevel, "Not specified"),

		PtrString(job.Description, "No description provided"),
		PtrString(job.CompanyCulture, "Not specified"),
		PtrString(job.CompanyValues, "Not specified"),

		FormatArray(job.Tools),
		FormatArray(job.ProgrammingLanguages),
		FormatArray(job.FrameworksAndLibraries),
		FormatArray(job.Databases),
		FormatArray(job.CloudTechnologies),
		FormatArray(job.IndustryKeywords),
		FormatArray(job.SoftSkills),
		FormatArray(job.Certifications),

		FormatArray(job.Requirements),
		FormatArray(job.NiceToHaves),

		PtrInt(job.ApplicantCount, 0),
	)
}

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

func JSONListToBulletPoints(list []string) string {
	return strings.Join(list, "\n- ")
}
