package formatters

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"text/template"

	"github.com/ordo_meritum/database/jobs"
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
<job_title> %s </job_title>
<company_name> %s </company_name>
<salary_range> %s </salary_range>
<applicant_count> %d </applicant_count>

<job_requirements_overview>
<skills_required> %s </skills_required>
<skills_nice_to_have> %s </skills_nice_to_have>
</job_requirements_overview>

<job_requirement_by_category>
<years_of_experience> %s </years_of_experience>
<education_level> %s </education_level>
<tools> %s </tools>
<programming_languages> %s </programming_languages>
<frameworks_and_libraries> %s </frameworks_and_libraries>
<databases> %s </databases>
<cloud_technologies> %s </cloud_technologies>
<soft_skills> %s </soft_skills>
<certifications> %s </certifications>
</job_requirement_by_category>

<full_job_description>
%s
</full_job_description>

<company_info>
<company_culture> %s </company_culture>
<company_values> %s </company_values>
<industry_keywords> %s </industry_keywords>
</company_info>
	`,
		job.JobTitle,
		job.CompanyName,
		PtrString(job.SalaryRange, "Not specified"),
		PtrInt(job.ApplicantCount, 0),

		FormatArray(job.Requirements),
		FormatArray(job.NiceToHaves),

		PtrString(job.YearsOfExp, "Not specified"),
		PtrString(job.EducationLevel, "Not specified"),
		FormatArray(job.Tools),
		FormatArray(job.ProgrammingLanguages),
		FormatArray(job.FrameworksAndLibraries),
		FormatArray(job.Databases),
		FormatArray(job.CloudTechnologies),
		FormatArray(job.SoftSkills),
		FormatArray(job.Certifications),

		PtrString(job.Description, "No description provided"),

		PtrString(job.CompanyCulture, "Not specified"),
		PtrString(job.CompanyValues, "Not specified"),
		FormatArray(job.IndustryKeywords),
	)
}

func JSONListToBulletPoints(list []string) string {
	return strings.Join(list, "\n- ")
}

func FormatTemplate(fs embed.FS, filename string, data any) (string, error) {
	content, err := fs.ReadFile(filename)

	if err != nil {
		return "", nil
	}

	empl, err := template.New(filename).Parse(string(content))
	if err != nil {
		log.Printf("could not parse template %s: %s", filename, err)
		return "", nil
	}

	var buf bytes.Buffer
	if err := empl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func cleanXMLTag(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	return nonAlphanumericRegex.ReplaceAllString(s, "")
}
