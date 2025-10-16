package domain

import (
	"fmt"
	"strconv"
	"strings"
)

type JobDescription struct {
	JobTitle               string   `json:"job_title"`
	CompanyName            string   `json:"company_name"`
	YearsOfExp             string   `json:"years_of_exp"`
	EducationLevel         string   `json:"education_level"`
	Website                string   `json:"website"`
	ApplicantCount         int      `json:"applicant_count"`
	PostAge                string   `json:"post_age"`
	SkillsRequired         []string `json:"skills_required"`
	SkillsNiceToHaves      []string `json:"skills_nice_to_haves"`
	ToolsAndTechnologies   []string `json:"tools_and_technologies"`
	ProgrammingLanguages   []string `json:"programming_languages"`
	FrameworksAndLibraries []string `json:"frameworks_and_libraries"`
	Databases              []string `json:"databases"`
	CloudTechnologies      []string `json:"cloud_technologies"`
	IndustryKeywords       []string `json:"industry_keywords"`
	SoftSkills             []string `json:"soft_skills"`
	Certifications         []string `json:"certifications"`
	CompanyCulture         string   `json:"company_culture"`
	CompanyValues          string   `json:"company_values"`
	SalaryRange            string   `json:"salary_range"`
}

func (jd *JobDescription) FormatForLLM() string {
	var builder strings.Builder

	appendString := func(key, value string) {
		if value != "" {
			builder.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		}
	}

	appendSlice := func(key string, values []string) {
		if len(values) > 0 {
			builder.WriteString(fmt.Sprintf("%s: %s\n", key, strings.Join(values, ", ")))
		}
	}

	builder.WriteString("--- Job Description ---\n")
	appendString("Job Title", jd.JobTitle)
	appendString("Company Name", jd.CompanyName)
	appendString("Years of Experience Required", jd.YearsOfExp)
	appendString("Education Level", jd.EducationLevel)
	appendString("Salary Range", jd.SalaryRange)
	appendString("Company Website", jd.Website)
	appendString("Company Culture", jd.CompanyCulture)
	appendString("Company Values", jd.CompanyValues)

	builder.WriteString("\n--- Key Requirements ---\n")
	appendSlice("Required Skills", jd.SkillsRequired)
	appendSlice("Nice-to-Have Skills", jd.SkillsNiceToHaves)
	appendSlice("Programming Languages", jd.ProgrammingLanguages)
	appendSlice("Frameworks and Libraries", jd.FrameworksAndLibraries)
	appendSlice("Tools and Technologies", jd.ToolsAndTechnologies)
	appendSlice("Databases", jd.Databases)
	appendSlice("Cloud Technologies", jd.CloudTechnologies)
	appendSlice("Industry Keywords", jd.IndustryKeywords)
	appendSlice("Soft Skills", jd.SoftSkills)
	appendSlice("Certifications", jd.Certifications)

	builder.WriteString("\n--- Posting Details ---\n")
	appendString("Post Age", jd.PostAge)
	if jd.ApplicantCount > 0 {
		appendString("Number of Applicants", strconv.Itoa(jd.ApplicantCount))
	}
	builder.WriteString("-----------------------\n")

	return builder.String()
}
