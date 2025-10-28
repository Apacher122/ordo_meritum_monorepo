package domain

// GuidingQuestions represents a set of questions and corresponding answers,
// often used for interview preparation or self-reflection.
type GuidingQuestions struct {
	Questions []string `json:"questions"`
	Answers   []string `json:"answers"`
}
