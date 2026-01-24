package events

import (
	"github.com/ordo_meritum/features/documents/models/domain"
	"github.com/ordo_meritum/features/documents/models/requests"
)

// DocumentEvent represents the payload sent to Kafka for document compilation.
// It contains all the necessary information, including user details and the
// generated document content (resume or cover letter), to create a final PDF.
type DocumentEvent struct {
	JobID         int                           `json:"jobID"`
	UserId        string                        `json:"userID"`
	CompanyName   string                        `json:"companyName"`
	DocType       string                        `json:"docType"`
	UserInfo      requests.UserInfoPayload      `json:"userInfo"`
	EducationInfo requests.EducationInfoPayload `json:"educationInfo,omitzero"`
	Resume        domain.Resume                 `json:"resume,omitzero"`
	CoverLetter   domain.CoverLetter            `json:"coverLetter,omitzero"`
}
