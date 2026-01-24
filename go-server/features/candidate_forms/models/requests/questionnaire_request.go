package requests

import "github.com/ordo_meritum/shared/models/requests"

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

// QuestionnaireRequest is a generic request body structure for questionnaire
// submissions, containing a payload and options.
type QuestionnaireRequest = requests.RequestBody[QuestionnarePayload, QuestionnaireOptions]

// QuestionnarePayload defines the core data submitted in a questionnaire request.
type QuestionnarePayload struct {
	BriefHistory        string                `json:"brief_history,omitempty"`
	QuestionsByCategory []QuestionsByCategory `json:"questions_form"`
}

// QuestionsByCategory groups questions and answers by their category.
type QuestionsByCategory struct {
	Category  QuestionCategory `json:"category"`
	Questions []QuestionAnswer `json:"questions"`
}

// QuestionAnswer represents a single question and its corresponding answer.
type QuestionAnswer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// QuestionnaireOptions specifies parameters for processing the questionnaire,
// such as the LLM provider and model to use.
type QuestionnaireOptions struct {
	LlmProvider string `json:"llm"`
	LlmModel    string `json:"llmModel"`
}
