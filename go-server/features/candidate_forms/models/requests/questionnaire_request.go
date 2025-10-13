package requests

type QuestionCategory string

const (
	ProblemSolving   QuestionCategory = "problem_solving"
	Communication    QuestionCategory = "communication"
	EmpathyTeamwork  QuestionCategory = "empathy_teamwork"
	Organization     QuestionCategory = "organization"
	Adaptability     QuestionCategory = "adaptability"
	Motivation       QuestionCategory = "motivation"
	StressManagement QuestionCategory = "stress_management"
	Creativity       QuestionCategory = "creativity"
)

type QuestionnaireRequest struct {
	BriefHistory        string                `json:"brief_history,omitempty"`
	QuestionsByCategory []QuestionsByCategory `json:"questions_form"`
}

type QuestionsByCategory struct {
	Category  QuestionCategory `json:"category"`
	Questions []QuestionAnswer `json:"questions"`
}

type QuestionAnswer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
