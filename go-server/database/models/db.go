package models

import (
	"time"

	"github.com/lib/pq"
)

type AppStatus string
type ShouldApply string
type Temperature string
type ProjectStatus string

const (
	StatusReject       AppStatus = "REJECT"
	StatusOffered      AppStatus = "OFFERED"
	StatusOpen         AppStatus = "OPEN"
	StatusClosed       AppStatus = "CLOSED"
	StatusMoved        AppStatus = "MOVED"
	StatusNotApplied   AppStatus = "NOT_APPLIED"
	StatusGhosted      AppStatus = "GHOSTED"
	StatusInterviewing AppStatus = "INTERVIEWING"

	ApplyStrongYes ShouldApply = "Strong Yes"
	ApplyYes       ShouldApply = "Yes"
	ApplyNo        ShouldApply = "No"
	ApplyStrongNo  ShouldApply = "Strong No"
	ApplyMaybe     ShouldApply = "Maybe"

	TempGood    Temperature = "Good"
	TempNeutral Temperature = "Neutral"
	TempBad     Temperature = "Bad"

	ProjectPlanned   ProjectStatus = "PLANNED"
	ProjectActive    ProjectStatus = "ACTIVE"
	ProjectCompleted ProjectStatus = "COMPLETED"
	ProjectOnHold    ProjectStatus = "ON_HOLD"
)

// --- Table Structs ---

type User struct {
	FirebaseUID string    `db:"firebase_uid"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Company struct {
	ID             int     `db:"id"`
	CompanyName    string  `db:"company_name"`
	Description    *string `db:"description"`
	Website        *string `db:"website"`
	Industry       *string `db:"industry"`
	Size           *string `db:"size"`
	Location       *string `db:"location"`
	CompanyCulture *string `db:"company_culture"`
	CompanyValues  *string `db:"company_values"`
	Benefits       *string `db:"benefits"`
}

type Role struct {
	ID                   int        `db:"id"`
	CompanyID            int        `db:"company_id"`
	JobTitle             string     `db:"job_title"`
	Description          *string    `db:"description"`
	SalaryRange          *string    `db:"salary_range"`
	TypicalSalaryAsk     *string    `db:"typical_salary_ask"`
	TypicalSalaryReason  *string    `db:"typical_salary_reason"`
	AdvisedSalaryAsk     *string    `db:"advised_salary_ask"`
	AdvisedSalaryReason  *string    `db:"advised_salary_reason"`
	ApplicationProcess   *string    `db:"application_process"`
	ExpectedResponseTime *string    `db:"expected_response_time"`
	ApplicationStatus    *AppStatus `db:"application_status"`
	UserApplied          *bool      `db:"user_applied"`
	CreatedAt            time.Time  `db:"created_at"`
	UpdatedAt            time.Time  `db:"updated_at"`
}

type JobRequirements struct {
	ID                      int            `db:"id"`
	RoleID                  int            `db:"role_id"`
	OtherInfo               *string        `db:"other_info"`
	YearsOfExp              *string        `db:"years_of_exp"`
	EducationLevel          *string        `db:"education_level"`
	Tools                   pq.StringArray `db:"tools"`
	ProgrammingLanguages    pq.StringArray `db:"programming_languages"`
	FrameworksAndLibraries  pq.StringArray `db:"frameworks_and_libraries"`
	Databases               pq.StringArray `db:"databases"`
	CloudTechnologies       pq.StringArray `db:"cloud_technologies"`
	IndustryKeywords        pq.StringArray `db:"industry_keywords"`
	SoftSkills              pq.StringArray `db:"soft_skills"`
	Certifications          pq.StringArray `db:"certifications"`
	Requirements            pq.StringArray `db:"requirements"`
	NiceToHaves             pq.StringArray `db:"nice_to_haves"`
	ApplicantCount          *int           `db:"applicant_count"`
	CodeAssessmentCompleted *bool          `db:"code_assessment_completed"`
	InterviewCount          *int           `db:"interview_count"`
	InitialApplicationDate  *time.Time     `db:"initial_application_date"`
	CreatedAt               time.Time      `db:"created_at"`
	UpdatedAt               time.Time      `db:"updated_at"`
}

type Resume struct {
	ID          int        `db:"id"`
	FirebaseUID string     `db:"firebase_uid"`
	RoleID      int        `db:"role_id"`
	AppliedOn   *time.Time `db:"applied_on"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

type MatchSummary struct {
	ID                int            `db:"id"`
	ResumeID          int            `db:"resume_id"`
	ShouldApply       ShouldApply    `db:"should_apply"`
	Reasoning         string         `db:"reasoning"`
	OverallMatchScore *int           `db:"overall_match_score"`
	Suggestions       pq.StringArray `db:"suggestions"`
	CreatedAt         time.Time      `db:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at"`
}

type MatchSummaryOverview struct {
	ID                 int         `db:"id"`
	MatchSummaryID     int         `db:"match_summary_id"`
	Summary            string      `db:"summary"`
	SummaryTemperature Temperature `db:"summary_temperature"`
	CreatedAt          time.Time   `db:"created_at"`
	UpdatedAt          time.Time   `db:"updated_at"`
}

type Metrics struct {
	ID             int     `db:"id"`
	MatchSummaryID int     `db:"match_summary_id"`
	ScoreTitle     string  `db:"score_title"`
	RawScore       float64 `db:"raw_score"`
	WeightedScore  float64 `db:"weighted_score"`
	ScoreWeight    float64 `db:"score_weight"`
	ScoreReason    *string `db:"score_reason"`
	IsCompatible   *bool   `db:"is_compatible"`
	Strength       *string `db:"strength"`
	Weaknesses     *string `db:"weaknesses"`
}

type Experience struct {
	ID        int        `db:"id"`
	ResumeID  int        `db:"resume_id"`
	Position  string     `db:"position"`
	Company   string     `db:"company"`
	StartDate time.Time  `db:"start_date"`
	EndDate   *time.Time `db:"end_date"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}

type ExperienceDescription struct {
	ID                     int     `db:"id"`
	ExpID                  int     `db:"exp_id"`
	Text                   string  `db:"text"`
	JustificationForChange *string `db:"justification_for_change"`
	NewSuggestion          bool    `db:"new_suggestion"`
}

type Project struct {
	ID        int           `db:"id"`
	ResumeID  int           `db:"resume_id"`
	Name      string        `db:"name"`
	Role      string        `db:"role"`
	Status    ProjectStatus `db:"status"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
}

type ProjectDescription struct {
	ID                     int     `db:"id"`
	ProjectID              int     `db:"project_id"`
	Text                   string  `db:"text"`
	JustificationForChange *string `db:"justification_for_change"`
	NewSuggestion          bool    `db:"new_suggestion"`
}

type Skill struct {
	ID                      int       `db:"id"`
	ResumeID                int       `db:"resume_id"`
	Category                string    `db:"category"`
	JustificationForChanges *string   `db:"justification_for_changes"`
	CreatedAt               time.Time `db:"created_at"`
	UpdatedAt               time.Time `db:"updated_at"`
}

type SkillItem struct {
	ID      int    `db:"id"`
	SkillID int    `db:"skill_id"`
	Name    string `db:"name"`
}

type Education struct {
	ID       int      `db:"id"`
	ResumeID int      `db:"resume_id"`
	School   string   `db:"school"`
	Degree   string   `db:"degree"`
	Dates    string   `db:"dates"`
	Location string   `db:"location"`
	Courses  *string  `db:"courses"`
	GPA      *float64 `db:"gpa"`
	Honors   *string  `db:"honors"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CandidateQuestionnaire struct {
	ID           int         `db:"id"`
	FirebaseUID  string      `db:"firebase_uid"`
	Title        *string     `db:"title"`
	BriefHistory *string     `db:"brief_history"`
	Questions    []Questions `db:"questions"`
	CreatedAt    time.Time   `db:"created_at"`
	UpdatedAt    time.Time   `db:"updated_at"`
}

type Questions struct {
	ID              int    `db:"id"`
	QuestionnaireID int    `db:"questionnaire_id"`
	Category        string `db:"category"`
	Question        string `db:"question"`
	Answer          string `db:"answer"`
}

type CandidateWritingSample struct {
	ID          int       `db:"id"`
	FirebaseUID string    `db:"firebase_uid"`
	Content     string    `db:"content"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type QuestionnaireResponse struct {
	ID              int       `db:"id"`
	QuestionnaireID int       `db:"questionnaire_id"`
	Question        string    `db:"question"`
	Response        string    `db:"response"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

type PersonalityProfile struct {
	ID             int    `db:"id"`
	FirebaseUID    string `db:"firebase_uid"`
	OceanProfileID int    `db:"ocean_profile_id"`
	DiscProfileID  int    `db:"disc_profile_id"`
}

type OceanProfile struct {
	ID                         int    `db:"id"`
	PersonalityProfilesID      int    `db:"personality_profiles_id"`
	OpennessScore              int    `db:"openness_to_experience_score"`
	OpennessReasoning          string `db:"openness_to_experience_reasoning"`
	ConscientiousnessScore     int    `db:"conscientiousness_score"`
	ConscientiousnessReasoning string `db:"conscientiousness_reasoning"`
	ExtraversionScore          int    `db:"extraversion_score"`
	ExtraversionReasoning      string `db:"extraversion_reasoning"`
	AgreeablenessScore         int    `db:"agreeableness_score"`
	AgreeablenessReasoning     string `db:"agreeableness_reasoning"`
	NeuroticismScore           int    `db:"neuroticism_score"`
	NeuroticismReasoning       string `db:"neuroticism_reasoning"`
	Summary                    string `db:"summary"`
}

type DiscProfile struct {
	ID                    int    `db:"id"`
	PersonalityProfilesID int    `db:"personality_profiles_id"`
	Dominance             string `db:"dominance"`
	Influence             string `db:"influence"`
	Steadiness            string `db:"steadiness"`
	Conscientiousness     string `db:"conscientiousness"`
	Summary               string `db:"summary"`
}
