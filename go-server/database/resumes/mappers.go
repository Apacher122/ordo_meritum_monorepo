package resumes

import (
	"github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/features/documents/models/domain"
)

func MapExperiencesToDomain(experiences []models.Experience, descs []models.ExperienceDescription) []domain.Experience {
	expMap := make(map[int][]models.ExperienceDescription)
	for _, d := range descs {
		expMap[d.ExpID] = append(expMap[d.ExpID], d)
	}

	domainExperiences := make([]domain.Experience, 0, len(experiences))
	for _, e := range experiences {
		startDateStr := e.StartDate.Format(dateFormat)
		endDateStr := "Present"
		if e.EndDate != nil {
			endDateStr = e.EndDate.Format(dateFormat)
		}

		bulletPoints := make([]domain.BulletPoint, 0, len(expMap[e.ID]))
		for _, desc := range expMap[e.ID] {
			bulletPoints = append(bulletPoints, domain.BulletPoint{
				Text:                   desc.Text,
				IsNewSuggestion:        desc.NewSuggestion,
				JustificationForChange: *desc.JustificationForChange,
			})
		}

		domainExperiences = append(domainExperiences, domain.Experience{
			Position:     e.Position,
			Company:      e.Company,
			Start:        startDateStr,
			End:          endDateStr,
			BulletPoints: bulletPoints,
		})
	}
	return domainExperiences
}

func MapProjectsToDomain(projects []models.Project, descs []models.ProjectDescription) []domain.Project {
	projMap := make(map[int][]models.ProjectDescription)
	for _, d := range descs {
		projMap[d.ProjectID] = append(projMap[d.ProjectID], d)
	}

	domainProjects := make([]domain.Project, 0, len(projects))
	for _, p := range projects {
		bulletPoints := make([]domain.BulletPoint, 0, len(projMap[p.ID]))
		for _, desc := range projMap[p.ID] {
			bulletPoints = append(bulletPoints, domain.BulletPoint{
				Text:                   desc.Text,
				IsNewSuggestion:        desc.NewSuggestion,
				JustificationForChange: *desc.JustificationForChange,
			})
		}

		domainProjects = append(domainProjects, domain.Project{
			Name:         p.Name,
			Role:         p.Role,
			BulletPoints: bulletPoints,
		})
	}
	return domainProjects
}

func MapSkillsToDomain(skills []models.Skill, items []models.SkillItem) []domain.Skills {
	skillMap := make(map[int][]models.SkillItem)
	for _, i := range items {
		skillMap[i.SkillID] = append(skillMap[i.SkillID], i)
	}

	domainSkills := make([]domain.Skills, 0, len(skills))
	for _, s := range skills {
		skillItems := make([]string, 0, len(skillMap[s.ID]))
		for _, item := range skillMap[s.ID] {
			skillItems = append(skillItems, item.Name)
		}

		domainSkills = append(domainSkills, domain.Skills{
			Category:                s.Category,
			JustificationForChanges: *s.JustificationForChanges,
			SkillItem:               skillItems,
		})
	}
	return domainSkills
}

func MapSummaryToDomain(overviewSummary string) []domain.SummaryBody {
	if overviewSummary == "" {
		return nil
	}
	return []domain.SummaryBody{
		{
			Sentence:               overviewSummary,
			JustificationForChange: "",
			NewSuggestion:          false,
		},
	}
}
