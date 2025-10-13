package requests

type ResumeRequest struct {
	Skills      []SkillsRequest     `json:"skills"`
	Experiences []ExperienceRequest `json:"experiences"`
	Projects    []ProjectRequest    `json:"projects"`
}

type SkillsRequest struct {
	Skill string `json:"skill"`
}

type ExperienceRequest struct {
	BulletPoints []string `json:"bulletPoints"`
	Company      string   `json:"company"`
	ID           string   `json:"id"`
	Position     string   `json:"position"`
	Years        string   `json:"years"`
}

type ProjectRequest struct {
	BulletPoints []string `json:"bulletPoints"`
	Description  string   `json:"description"`
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Years        string   `json:"years"`
}
