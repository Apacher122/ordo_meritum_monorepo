package formatters

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/ordo_meritum/features/documents/models/domain"
	"github.com/ordo_meritum/features/documents/models/requests"
)

func FormatResumeForLLMWithXML(payload requests.ResumeRequest) string {
	var sb strings.Builder

	sb.WriteString("<resume_content>\n")

	if len(payload.Experiences) > 0 {
		sb.WriteString("\t<experiences>\n")
		for _, exp := range payload.Experiences {
			sb.WriteString("\t\t<job>\n")
			sb.WriteString(fmt.Sprintf("\t\t\t<title>%s</title>\n", exp.Position))
			sb.WriteString(fmt.Sprintf("\t\t\t<company>%s</company>\n", exp.Company))
			sb.WriteString(fmt.Sprintf("\t\t\t<dates>%s</dates>\n", exp.Years))
			sb.WriteString("\t\t\t<bullet_points>\n")
			for _, point := range exp.BulletPoints {
				sb.WriteString(fmt.Sprintf("\t\t\t\t<bullet>%s</bullet>\n", strings.TrimSpace(point)))
			}
			sb.WriteString("\t\t\t</bullet_points>\n")
			sb.WriteString("\t\t</job>\n")
		}
		sb.WriteString("\t</experiences>\n")
	}

	if len(payload.Projects) > 0 {
		sb.WriteString("\t<projects>\n")
		for _, proj := range payload.Projects {
			sb.WriteString("\t\t<project>\n")
			sb.WriteString(fmt.Sprintf("\t\t\t<name>%s</name>\n", proj.Name))
			sb.WriteString(fmt.Sprintf("\t\t\t<role>%s</role>\n", proj.Description))
			sb.WriteString("\t\t\t<bullet_points>\n")
			for _, point := range proj.BulletPoints {
				sb.WriteString(fmt.Sprintf("\t\t\t\t<bullet>%s</bullet>\n", strings.TrimSpace(point)))
			}
			sb.WriteString("\t\t\t</bullet_points>\n")
			sb.WriteString("\t\t</project>\n")
		}
		sb.WriteString("\t</projects>\n")
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
		sb.WriteString("\t<projects>\n")
		for _, proj := range payload.Projects {
			sb.WriteString("\t\t<project>\n")
			sb.WriteString(fmt.Sprintf("\t\t\t<name>%s</name>\n", proj.Name))
			sb.WriteString(fmt.Sprintf("\t\t\t<role>%s</role>\n", proj.Role))
			sb.WriteString("\t\t\t<bullet_points>\n")
			for _, point := range proj.BulletPoints {
				sb.WriteString(fmt.Sprintf("\t\t\t\t<bullet>%s</bullet>\n", strings.TrimSpace(point.Text)))
			}
			sb.WriteString("\t\t\t</bullet_points>\n")
			sb.WriteString("\t\t</project>\n")
		}
		sb.WriteString("\t</projects>\n")
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
