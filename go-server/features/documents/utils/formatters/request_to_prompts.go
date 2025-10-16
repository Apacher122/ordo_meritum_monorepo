package formatters

import (
	"fmt"
	"strings"

	"github.com/ordo_meritum/features/documents/models/requests"
)

func FormatResumeRequestForLLMWithXML(request *requests.DocumentPayload) string {
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
