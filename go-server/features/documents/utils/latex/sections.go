package latex

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ordo_meritum/features/documents/models/domain"
)

func GenerateResumeSections(resume *domain.Resume) string {
	var sections bytes.Buffer
	sections.WriteString(summarySection(resume.Summary))
	sections.WriteString(experienceSection(resume.Experiences))
	sections.WriteString(skillsSection(resume.Skills))
	sections.WriteString(projectsSection(resume.Projects))
	return sections.String()
}

func summarySection(summary []domain.SummaryBody) string {
	if len(summary) == 0 {
		return ""
	}

	var sentences []string
	for _, item := range summary {
		sentences = append(sentences, item.Sentence)
	}

	joined := strings.Join(sentences, " ")

	return fmt.Sprintf("\\cvsection{Summary}\n\\cvtext{%s}\n", EscapeChars(joined))
}

func experienceSection(experiences []domain.Experience) string {
	if len(experiences) == 0 {
		return ""
	}
	var b bytes.Buffer
	b.WriteString("\\cvsection{Experience}\n\\begin{cventries}\n")
	for _, exp := range experiences {
		var descItems string
		for _, item := range exp.BulletPoints {
			descItems += fmt.Sprintf("\n    \\item{%s}", EscapeChars(item.Text))
		}
		entry := fmt.Sprintf(`\cventry
  {%s}
  {%s}
  {%s}
  {%s -- %s}
  {\begin{cvitems}%s
    \end{cvitems}}`+"\n",
			EscapeChars(exp.Position), EscapeChars(exp.Company), "",
			EscapeChars(exp.Start), EscapeChars(exp.End), descItems)
		b.WriteString(entry)
	}
	b.WriteString("\\end{cventries}\n")
	return b.String()
}

func skillsSection(skills []domain.Skills) string {
	if len(skills) == 0 {
		return ""
	}
	skillString := ""

	skillString += "\\begin{cvskills}\n"

	for _, skill := range skills {
		skillString += "\\cvskill\n"
		skillString += fmt.Sprintf("\t{%s}\n", EscapeChars(skill.Category))

		for _, item := range skill.SkillItem {
			skillString += fmt.Sprintf("\t{%s}\n", EscapeChars(item))
		}
	}
	skillString += "\\end{cvskills}\n"
	return skillString
}

func projectsSection(projects []domain.Project) string {
	if len(projects) == 0 {
		return ""
	}
	var b bytes.Buffer
	b.WriteString("\\cvsection{Projects}\n\\begin{cventries}\n")
	for _, proj := range projects {
		var descItems string
		for _, item := range proj.BulletPoints {
			descItems += fmt.Sprintf("\n    \\item{%s}", EscapeChars(item.Text))
		}
		entry := fmt.Sprintf(`\cventry
  {%s}
  {}
  {}
  {}
  {\begin{cvitems}%s
    \end{cvitems}}`+"\n",
			EscapeChars(proj.Name), descItems)
		b.WriteString(entry)
	}
	b.WriteString("\\end{cventries}\n")
	return b.String()
}
