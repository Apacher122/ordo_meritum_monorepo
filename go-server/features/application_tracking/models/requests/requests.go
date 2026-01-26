package request

import (
	"time"

	db_models "github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/shared/models/requests"
)

type JobPostingRequestV1 = requests.RequestBody[JobPostingPayload, JobPostingOptions]

type ApplicationUpdateRequest = requests.RequestBody[ApplicationUpdatePayload, struct{}]

// JobPostingPayload defines the raw data extracted or provided for a job posting.
type JobPostingPayload struct {
	CompanyName    string `json:"company"`
	JobTitle       string `json:"job_title"`
	Link           string `json:"website"`
	ApplicantCOunt string `json:"applicant_count"`
	TimeAgo        string `json:"time_ago"`
	JobDescription string `json:"job_description"`
}

// JobPostingOptions specifies parameters for processing the job posting payload, like the LLM to use.
type JobPostingOptions struct {
	LlmProvider string `json:"llm"`
	LlmModel    string `json:"llmModel"`
}

// ApplicationUpdatePayload defines the fields that can be updated for an existing job application.
// Fields are pointers to allow partial updates (only non-nil fields are updated).
type ApplicationUpdatePayload struct {
	JobID                  int                  `json:"job_id"`
	JobTitle               *string              `json:"job_title,omitempty"`
	Link                   *string              `json:"website,omitempty"`
	ApplicationStatus      *db_models.AppStatus `json:"application_status,omitempty"`
	InterviewCount         *int                 `json:"interview_count,omitempty"`
	InitialApplicationDate *time.Time           `json:"initial_application_date,omitempty"`
}

// JobPostingRequest defines the structure for submitting a job posting.
type JobPostingRequest struct {
	RawHtml     string `json:"raw_html"`
	CompanyName string `json:"company_name"`
	JobTitle    string `json:"job_title"`
	Link        string `json:"website"`
}

// JobPostingEvent represents the structured data extracted from a job posting,
// typically after processing by an LLM.
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

// // FormatJobPostingRequest converts a JobPostingRequest into a plain text format,
// // likley intended for use in LLM prompts.
// func FormatJobPostingRequest(jp *JobPostingRequest) string {
// 	return fmt.Sprintf(`
// Company: %s
// Position: %s
// URL: %s
// Number of Applicants: %s
// Post Age: %s

// Job Description:
// %s
// 	`,
// 		jp.CompanyName,
// 		jp.JobTitle,
// 		jp.Link,
// 	)
// }
