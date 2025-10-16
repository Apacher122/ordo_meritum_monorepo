package mappers

import (
	"github.com/ordo_meritum/database/jobs"
	"github.com/ordo_meritum/features/application_tracking/models/domain"
)

func NewJobDescriptionFromPost(post *jobs.FullJobPosting) *domain.JobDescription {
	jd := &domain.JobDescription{
		JobTitle:               post.JobTitle,
		CompanyName:            post.CompanyName,
		SkillsRequired:         post.Requirements,
		SkillsNiceToHaves:      post.NiceToHaves,
		ToolsAndTechnologies:   post.Tools,
		ProgrammingLanguages:   post.ProgrammingLanguages,
		FrameworksAndLibraries: post.FrameworksAndLibraries,
		Databases:              post.Databases,
		CloudTechnologies:      post.CloudTechnologies,
		IndustryKeywords:       post.IndustryKeywords,
		SoftSkills:             post.SoftSkills,
		Certifications:         post.Certifications,
	}

	if post.YearsOfExp != nil {
		jd.YearsOfExp = *post.YearsOfExp
	}
	if post.EducationLevel != nil {
		jd.EducationLevel = *post.EducationLevel
	}
	if post.ApplicantCount != nil {
		jd.ApplicantCount = *post.ApplicantCount
	}
	if post.CompanyCulture != nil {
		jd.CompanyCulture = *post.CompanyCulture
	}
	if post.CompanyValues != nil {
		jd.CompanyValues = *post.CompanyValues
	}
	if post.SalaryRange != nil {
		jd.SalaryRange = *post.SalaryRange
	}

	return jd
}
