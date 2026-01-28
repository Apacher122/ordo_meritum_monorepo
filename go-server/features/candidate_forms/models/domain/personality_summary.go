package domain

// PersonalitySummary aggregates results from different personality assessments,
// currently OCEAN and DISC profiles.
type PersonalitySummary struct {
	OCEAN OCEANProfile `json:"ocean,omitempty"`
	DISC  DISCProfile  `json:"disc,omitempty"`
}

// OCEANCategory represents the five factors in the OCEAN personality model.
type OCEANCategory string

const (
	Openness          OCEANCategory = "Openness"
	Conscientiousness OCEANCategory = "Conscientiousness"
	Extraversion      OCEANCategory = "Extraversion"
	Agreeableness     OCEANCategory = "Agreeableness"
	Neuroticism       OCEANCategory = "Neuroticism"
)

// OCEANProfile holds the scores and summary for an OCEAN assessment.
type OCEANProfile struct {
	Scores  []OCEANScore `json:"scores"`
	Summary string       `json:"summary"`
}

// OCEANScore represents a score and reasoning for a specific OCEAN category.
type OCEANScore struct {
	Category  OCEANCategory `json:"category"`
	Score     int           `json:"score"`
	Reasoning string        `json:"reasoning"`
}

// DISCCategory represents the four factors in the DISC assessment model.
type DISCCategory string

const (
	Dominance   DISCCategory = "Dominance"
	Influence   DISCCategory = "Influence"
	Steadiness  DISCCategory = "Steadiness"
	Consistency DISCCategory = "Conscientiousness"
)

// DISCProfile holds the reasoning scores and summary for a DISC assessment.
type DISCProfile struct {
	Scores  []DISCScore `json:"scores"`
	Summary string      `json:"summary"`
}

// DISCScore represents the reasoning for a specific DISC category.
type DISCScore struct {
	Category  DISCCategory `json:"category"`
	Reasoning string       `json:"reasoning"`
}

// PersonalityArchetype defines distinct personality profiles based on assessment metrics.
type PersonalityArchetype string

// Constants define the recognized personality archetypes.
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

// ArchetypeMetrics defines the key dimensions used to differentiate archetypes.
type ArchetypeMetrics string

// Constants define the valid metrics for archetype analysis.
const (
	Creativity        ArchetypeMetrics = "Creativity vs. Structure"
	Collaboration     ArchetypeMetrics = "Collaboration vs. Independence"
	ActionOrientation ArchetypeMetrics = "Action Orientation"
	RiskTolerance     ArchetypeMetrics = "Risk Tolerance"
	Empathy           ArchetypeMetrics = "Empathy & Communication"
	Vision            ArchetypeMetrics = "Vision vs. Detail Focus"
	Adaptability      ArchetypeMetrics = "Adaptability"
)

// ArchetypeProfile holds the determined archetype and related metric insights.
type ArchetypeProfile struct {
	Archetype      ArchetypeSummary   `json:"archetype"`
	MetricInsights []KeyMetricInsight `json:"metric_insights"`
}

// ArchetypeSummary contains the name and description of a personality archetype.
type ArchetypeSummary struct {
	Archetype PersonalityArchetype `json:"archetype"`
	Summary   string               `json:"summary"`
}

// KeyMetricInsight represents a score for a specific archetype metric.
type KeyMetricInsight struct {
	Metric ArchetypeMetrics `json:"metric"`
	Score  int              `json:"score"`
}
