package requests

import (
	document_requests "github.com/ordo_meritum/features/documents/models/requests"
	"github.com/ordo_meritum/shared/models/requests"
)

type JobGuideRequest = requests.RequestBody[JobGuidePayload, JobGuideOptions]

type JobGuidePayload struct {
	JobID         int                                    `json:"job_id"`
	EducationInfo document_requests.EducationInfoPayload `json:"education_info,omitempty"`
	CoverLetter   document_requests.CoverLetterPayload   `json:"cover_letter,omitzero"`
}

type JobGuideOptions struct {
	GuideType   string `json:"guide_type,omitempty"`
	LLMProvider string `json:"llm"`
	LlmModel    string `json:"llmModel,omitempty"`
	GetNew      bool   `json:"getNew,omitempty"`
}
