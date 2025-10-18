package formatters

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/ordo_meritum/database/jobs"
	error_messages "github.com/ordo_meritum/shared/utils/errors"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9_]+`)

func FormatAboutForLLMWithXML(jsonData []byte) (string, error) {
	var sections map[string]string
	if err := json.Unmarshal(jsonData, &sections); err != nil {
		return "", fmt.Errorf("failed to unmarshal json data: %w", err)
	}

	var sb strings.Builder
	sb.WriteString("<additional_info>\n")

	for key, value := range sections {
		if value == "" {
			continue
		}

		tagName := cleanXMLTag(key)

		sb.WriteString(fmt.Sprintf("\t<%s>\n", tagName))
		sb.WriteString(fmt.Sprintf("\t\t%s\n", strings.TrimSpace(value)))
		sb.WriteString(fmt.Sprintf("\t</%s>\n", tagName))
	}

	sb.WriteString("</additional_info>")

	return sb.String(), nil
}

func FormatJobPostForLLM(job jobs.FullJobPosting) string {
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

func JSONListToBulletPoints(list []string) string {
	return strings.Join(list, "\n- ")
}

func FormatTemplate(fs embed.FS, filename string, data any) (string, *error_messages.ErrorBody) {
	content, err := fs.ReadFile(filename)

	if err != nil {
		return "", &error_messages.ErrorBody{ErrMsg: err}
	}

	empl, err := template.New(filename).Parse(string(content))
	if err != nil {
		return "", &error_messages.ErrorBody{ErrMsg: fmt.Errorf("could not parse template %s: %s", filename, err)}
	}

	var buf bytes.Buffer
	if err := empl.Execute(&buf, data); err != nil {
		return "", &error_messages.ErrorBody{ErrMsg: err}
	}

	return buf.String(), nil
}

func cleanXMLTag(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	return nonAlphanumericRegex.ReplaceAllString(s, "")
}
