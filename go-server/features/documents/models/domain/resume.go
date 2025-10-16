package domain

import (
	"fmt"
	"strings"
)

type Resume struct {
	Summary     []SummaryBody `json:"summary,omitempty"`
	Skills      []Skills      `json:"skills"`
	Experiences []Experience  `json:"experiences"`
	Projects    []Project     `json:"projects"`
}

type SummaryBody struct {
	Sentence               string `json:"sentence"`
	JustificationForChange string `json:"justification_for_change,omitempty"`
	NewSuggestion          bool   `json:"is_new_suggestion,omitempty"`
}

type Skills struct {
	Category                string   `json:"category,omitempty"`
	SkillItem               []string `json:"skill"`
	JustificationForChanges string   `json:"justification_for_changes,omitempty"`
}

type Experience struct {
	BulletPoints []BulletPoint `json:"bulletPoints"`
	Company      string        `json:"company"`
	ID           string        `json:"id"`
	Position     string        `json:"position"`
	Start        string        `json:"start"`
	End          string        `json:"end"`
}

type Project struct {
	BulletPoints []BulletPoint `json:"bulletPoints"`
	Role         string        `json:"role"`
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Status       string        `json:"status"`
}

type BulletPoint struct {
	Text                   string `json:"text"`
	IsNewSuggestion        bool   `json:"is_new_suggestion"`
	JustificationForChange string `json:"justification_for_change"`
}

func (r *Resume) FormatForLLM() string {
	var sb strings.Builder

	sb.WriteString("<resume_content>\n")

	if len(r.Summary) > 0 {
		sb.WriteString("\t<summary>\n")
		for _, s := range r.Summary {
			sb.WriteString(fmt.Sprintf("\t\t<sentence>%s</sentence>\n", strings.TrimSpace(s.Sentence)))
		}
		sb.WriteString("\t</summary>\n")
	}

	if len(r.Experiences) > 0 {
		sb.WriteString("\t<experiences>\n")
		for _, exp := range r.Experiences {
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

	if len(r.Projects) > 0 {
		sb.WriteString("\t<personal_projects>\n")
		for _, proj := range r.Projects {
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

	if len(r.Skills) > 0 {
		sb.WriteString("\t<skills_section>\n")
		for _, skillCat := range r.Skills {
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
