package request

type JobPostingRequest struct {
	CompanyName    string `json:"company"`
	JobTitle       string `json:"job_title"`
	Link           string `json:"website"`
	ApplicantCount int    `json:"applicant_count"`
	TimeAgo        string `json:"time_ago"`
	JobDescription string `json:"job_description"`
}

type JobPostingEvent struct {
	JobTitle               string   `json:"job_title"`
	Company                string   `json:"company_name"`
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
