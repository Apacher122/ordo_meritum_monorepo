package requests

import (
	document_requests "github.com/ordo_meritum/features/documents/models/requests"
	"github.com/ordo_meritum/shared/models/requests"
)

// JobGuideRequest is a generic request body structure for job guide-related
// operations, containing a payload and a set of options.
type JobGuideRequest = requests.RequestBody[JobGuidePayload, JobGuideOptions]

// JobGuidePayload defines the core data used for job guide generation.
type JobGuidePayload struct {
	JobID         int                                    `json:"job_id"`
	EducationInfo document_requests.EducationInfoPayload `json:"education_info,omitempty"`
	CoverLetter   document_requests.CoverLetterPayload   `json:"cover_letter,omitzero"`
}

// JobGuideOptions specifies the parameters and settings for a job guide request.
type JobGuideOptions struct {
	GuideType   string `json:"guide_type,omitempty"`
	LLMProvider string `json:"llm"`
	LlmModel    string `json:"llmModel,omitempty"`
	GetNew      bool   `json:"getNew,omitempty"`
}
