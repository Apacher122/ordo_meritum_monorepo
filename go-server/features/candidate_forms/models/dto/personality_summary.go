package dto

type PersonalitySummary struct {
	OCEAN OCEANProfile `json:"ocean,omitempty"`
	DISC  DISCProfile  `json:"disc,omitempty"`
}

// --- OCEAN Model ---

type OCEANCategory string

const (
	Openness          OCEANCategory = "Openness"
	Conscientiousness OCEANCategory = "Conscientiousness"
	Extraversion      OCEANCategory = "Extraversion"
	Agreeableness     OCEANCategory = "Agreeableness"
	Neuroticism       OCEANCategory = "Neuroticism"
)

type OCEANProfile struct {
	Scores  []OCEANScore `json:"scores"`
	Summary string       `json:"summary"`
}

type OCEANScore struct {
	Category  OCEANCategory `json:"category"`
	Score     int           `json:"score"`
	Reasoning string        `json:"reasoning"`
}

// --- DISC Model -

type DISCCategory string

const (
	Dominance   DISCCategory = "Dominance"
	Influence   DISCCategory = "Influence"
	Steadiness  DISCCategory = "Steadiness"
	Consistency DISCCategory = "Conscientiousness"
)

type DISCProfile struct {
	Scores  []DISCScore `json:"scores"`
	Summary string      `json:"summary"`
}

type DISCScore struct {
	Category  DISCCategory `json:"category"`
	Reasoning string       `json:"reasoning"`
}
