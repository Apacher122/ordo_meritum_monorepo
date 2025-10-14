package domain

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

// -- Personality Archetypes --

type PersonalityArchetype string

const (
	Strategist PersonalityArchetype = "The Strategist"
	Innovator  PersonalityArchetype = " The Innovator"
	Diplomat   PersonalityArchetype = "The Diplomat"
	Anchor     PersonalityArchetype = "The Anchor"
	Visionary  PersonalityArchetype = "The Visionary"
	Executor   PersonalityArchetype = "The Executor"
	Analyst    PersonalityArchetype = "The Analyst"
	Builder    PersonalityArchetype = "The Builder"
	Connector  PersonalityArchetype = "The Connector"
)

type ArchetypeMetrics string

const (
	Creativity        ArchetypeMetrics = "Creativity vs. Structure"
	Collaboration     ArchetypeMetrics = "Collaboration vs. Independence"
	ActionOrientation ArchetypeMetrics = "Action Orientation"
	RiskTolerance     ArchetypeMetrics = "Risk Tolerance"
	Empathy           ArchetypeMetrics = "Empathy & Communication"
	Vision            ArchetypeMetrics = "Vision vs. Detail Focus"
	Adaptability      ArchetypeMetrics = "Adaptability"
)

type ArchetypeProfile struct {
	Archetype      ArchetypeSummary   `json:"archetype"`
	MetricInsights []KeyMetricInsight `json:"metric_insights"`
}

type ArchetypeSummary struct {
	Archetype PersonalityArchetype `json:"archetype"`
	Summary   string               `json:"summary"`
}

type KeyMetricInsight struct {
	Metric ArchetypeMetrics `json:"metric"`
	Score  int              `json:"score"`
}
