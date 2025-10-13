package domain

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
	SkillItem               []string `json:"skills"`
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
