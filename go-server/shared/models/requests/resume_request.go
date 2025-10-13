package requests

type ResumeRequest struct {
	Summary       string              `json:"summary.omitempty"`
	Skills        []SkillsRequest     `json:"skills"`
	EducationInfo EducationInfo       `json:"educationInfo"`
	Experiences   []ExperienceRequest `json:"experiences"`
	Projects      []ProjectRequest    `json:"projects,omitempty"`
}

type SkillsRequest struct {
	Skill string `json:"skill"`
}

type EducationInfo struct {
	CourseWork string `json:"coursework"`
	Degree     string `json:"degree"`
	Location   string `json:"location"`
	School     string `json:"school"`
	StartEnd   string `json:"start_end"`
}

type ExperienceRequest struct {
	BulletPoints []BulletPoint `json:"bulletPoints"`
	Company      string        `json:"company"`
	ID           string        `json:"id"`
	Position     string        `json:"position"`
	Years        string        `json:"years"`
}

type ProjectRequest struct {
	BulletPoints []BulletPoint `json:"bulletPoints"`
	Description  string        `json:"description"`
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Years        string        `json:"years"`
}

type BulletPoint struct {
	Text string `json:"text"`
}
