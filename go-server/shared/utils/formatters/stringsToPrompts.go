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
	"github.com/ordo_meritum/features/documents/models/domain"
	"github.com/ordo_meritum/features/documents/models/requests"
)

func FormatResumeForLLMWithXML(request *requests.DocumentPayload) string {
	var sb strings.Builder

	payload := request.Resume

	sb.WriteString("<resume_content>\n")

	if len(payload.Experiences) > 0 {
		sb.WriteString("\t<experiences>\n")
		for _, exp := range payload.Experiences {
			sb.WriteString("\t\t<job>\n")
			sb.WriteString(fmt.Sprintf("\t\t\t<position>%s</position>\n", exp.Position))
			sb.WriteString(fmt.Sprintf("\t\t\t<company>%s</company>\n", exp.Company))
			sb.WriteString(fmt.Sprintf("\t\t\t<dates>%s</dates>\n", exp.Years))
			sb.WriteString("\t\t\t<experience_bullet_points>\n")
			for _, point := range exp.BulletPoints {
				sb.WriteString(fmt.Sprintf("\t\t\t\t<experience_bullet>%s</experience_bullet>\n", strings.TrimSpace(point)))
			}
			sb.WriteString("\t\t\t</experience_bullet_points>\n")
			sb.WriteString("\t\t</job>\n")
		}
		sb.WriteString("\t</experiences>\n")
	}

	if len(payload.Projects) > 0 {
		sb.WriteString("\t<personal_projects>\n")
		for _, proj := range payload.Projects {
			sb.WriteString("\t\t<project>\n")
			sb.WriteString(fmt.Sprintf("\t\t\t<project_name>%s</project_name>\n", proj.Name))
			sb.WriteString(fmt.Sprintf("\t\t\t<candidate_role_in_project>%s</candidate_role_in_projec>\n", proj.Description))
			sb.WriteString("\t\t\t<project_bullet_points>\n")
			for _, point := range proj.BulletPoints {
				sb.WriteString(fmt.Sprintf("\t\t\t\t<project_bullet>%s</project_bullet>\n", strings.TrimSpace(point)))
			}
			sb.WriteString("\t\t\t</project_bullet_points>\n")
			sb.WriteString("\t\t</project>\n")
		}
		sb.WriteString("\t</personal_projects>\n")
	}

	if len(payload.Skills) > 0 {
		sb.WriteString("\t<skills_section>\n")
		sb.WriteString("\t\t<skill_list>\n")
		for _, s := range payload.Skills {
			sb.WriteString(fmt.Sprintf("\t\t\t<skill>%s</skill>\n", strings.TrimSpace(s.Skill)))
		}

		sb.WriteString("\t\t</skill_list>\n")
		sb.WriteString("\t</skills_section>\n")
	}
	sb.WriteString("</resume_content>")

	return sb.String()
}

func FormatResumePayloadForLLMWithXML(payload domain.Resume) string {
	var sb strings.Builder

	sb.WriteString("<resume_content>\n")

	if len(payload.Summary) > 0 {
		sb.WriteString("\t<summary>\n")
		for _, s := range payload.Summary {
			sb.WriteString(fmt.Sprintf("\t\t<sentence>%s</sentence>\n", strings.TrimSpace(s.Sentence)))
		}
		sb.WriteString("\t</summary>\n")
	}

	if len(payload.Experiences) > 0 {
		sb.WriteString("\t<experiences>\n")
		for _, exp := range payload.Experiences {
			sb.WriteString("\t\t<job>\n")
			sb.WriteString(fmt.Sprintf("\t\t\t<title>%s</title>\n", exp.Position))
			sb.WriteString(fmt.Sprintf("\t\t\t<company>%s</company>\n", exp.Company))
			sb.WriteString(fmt.Sprintf("\t\t\t<dates>%s - %s</dates>\n", exp.Start, exp.End))
			sb.WriteString("\t\t\t<bullet_points>\n")
			for _, point := range exp.BulletPoints {
				sb.WriteString(fmt.Sprintf("\t\t\t\t<bullet>%s</bullet>\n", strings.TrimSpace(point.Text)))
			}
			sb.WriteString("\t\t\t</bullet_points>\n")
			sb.WriteString("\t\t</job>\n")
		}
		sb.WriteString("\t</experiences>\n")
	}

	if len(payload.Projects) > 0 {
		sb.WriteString("\t<personal_projects>\n")
		for _, proj := range payload.Projects {
			sb.WriteString("\t\t<project>\n")
			sb.WriteString(fmt.Sprintf("\t\t\t<project_name>%s</project_name>\n", proj.Name))
			sb.WriteString(fmt.Sprintf("\t\t\t<candidate_role_in_project>%s</candidate_role_in_projec>\n", proj.Role))
			sb.WriteString("\t\t\t<project_bullet_points>\n")
			for _, point := range proj.BulletPoints {
				sb.WriteString(fmt.Sprintf("\t\t\t\t<project_bullet>%s</project_bullet>\n", strings.TrimSpace(point.Text)))
			}
			sb.WriteString("\t\t\t</project_bullet_points>\n")
			sb.WriteString("\t\t</project>\n")
		}
		sb.WriteString("\t</personal_projects>\n")
	}

	if len(payload.Skills) > 0 {
		sb.WriteString("\t<skills_section>\n")
		for _, skillCat := range payload.Skills {
			sb.WriteString(fmt.Sprintf("\t\t<skill_category name=\"%s\">\n", skillCat.Category))
			for _, skill := range skillCat.SkillItem {
				trimmedSkill := strings.TrimSpace(skill)
				if trimmedSkill != "" {
					sb.WriteString(fmt.Sprintf("\t\t\t<skill>%s</skill>\n", trimmedSkill))
				}
			}
			sb.WriteString("\t\t</skill_category>\n")
		}
		sb.WriteString("\t</skills_section>\n")
	}

	sb.WriteString("</resume_content>")

	return sb.String()
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9_]+`)

func cleanXMLTag(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	return nonAlphanumericRegex.ReplaceAllString(s, "")
}

/**
 * FormatEssayForLLMWithXML takes a JSON string containing various text sections
 * and formats them into a single, XML-tagged string for an LLM prompt.
 * @param jsonData The raw JSON data as a byte slice.
 * @returns A formatted XML string and an error if parsing fails.
 */
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

func FormatEducationForLLM(e requests.EducationInfoPayload) string {
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
