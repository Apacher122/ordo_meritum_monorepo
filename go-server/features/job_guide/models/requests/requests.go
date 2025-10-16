package requests

import (
	document_requests "github.com/ordo_meritum/features/documents/models/requests"
	"github.com/ordo_meritum/shared/models/requests"
)

type JobGuideRequests = requests.RequestBody[JobGuidePayload, JobGuideOptions]

type JobGuidePayload struct {
	JobID         int                                    `json:"job_id"`
	EducationInfo document_requests.EducationInfoPayload `json:"education_info"`
	CoverLetter   document_requests.CoverLetterPayload   `json:"cover_letter,omitzero"`
}

type JobGuideOptions struct {
	GuideType   string `json:"guide_type"`
	LLMProvider string `json:"llm"`
	LlmModel    string `json:"llmModel"`
	GetNew      bool   `json:"getNew,omitempty"`
}
