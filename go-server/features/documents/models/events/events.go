package events

import (
	"github.com/ordo_meritum/features/documents/models/domain"
	"github.com/ordo_meritum/features/documents/models/requests"
)

type DocumentEvent struct {
	JobID         int                           `json:"jobID"`
	UserId        string                        `json:"userID"`
	CompanyName   string                        `json:"companyName"`
	DocType       string                        `json:"docType"`
	UserInfo      requests.UserInfoRequest      `json:"userInfo"`
	EducationInfo requests.EducationInfoRequest `json:"educationInfo"`
	Resume        domain.Resume                 `json:"resume,omitzero"`
	CoverLetter   domain.CoverLetter            `json:"coverLetter,omitzero"`
}
