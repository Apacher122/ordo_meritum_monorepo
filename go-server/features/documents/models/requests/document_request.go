package requests

import (
	"encoding/json"

	"github.com/ordo_meritum/shared/models/requests"
)

// DocumentRequest is a generic request body structure for document-related
// operations, containing a payload and a set of options.
type DocumentRequest = requests.RequestBody[DocumentPayload, DocumentOptions]

// DocumentPayload defines the core data used for document generation,
// including user information and resume content.
type DocumentPayload struct {
	Resume         ResumePayload        `json:"resume,omitempty"`
	UserInfo       UserInfoPayload      `json:"userInfo"`
	AdditionalInfo json.RawMessage      `json:"additionalInfo"`
	EducationInfo  EducationInfoPayload `json:"educationInfo"`
	Coverletter    CoverLetterPayload   `json:"coverletter,omitzero"`
}

// DocumentOptions specifies the parameters and settings for a document
// generation request.
type DocumentOptions struct {
	JobID          int      `json:"jobId"`
	DocType        string   `json:"docType"`
	LlmProvider    string   `json:"llm"`
	LlmModel       string   `json:"llmModel"`
	GetNew         bool     `json:"getNew,omitempty"`
	Corrections    []string `json:"corrections,omitempty"`
	WritingSamples []string `json:"writingSamples,omitempty"`
}

// ResumePayload contains the detailed components of a resume provided in a request.
type ResumePayload struct {
	Skills      []SkillsPayload     `json:"skills"`
	Experiences []ExperiencePayload `json:"experiences"`
	Projects    []ProjectPayload    `json:"projects"`
}

// SkillsPayload represents a single skill item in a request.
type SkillsPayload struct {
	Skill string `json:"skill"`
}

// ExperiencePayload represents a single work experience entry in a request.
type ExperiencePayload struct {
	BulletPoints []string `json:"bulletPoints"`
	Company      string   `json:"company"`
	ID           string   `json:"id"`
	Position     string   `json:"jobTitle"`
	Years        string   `json:"years"`
}

// ProjectPayload represents a single project entry in a request.
type ProjectPayload struct {
	BulletPoints []string `json:"bulletPoints"`
	Description  string   `json:"description"`
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Years        string   `json:"years"`
}

// CoverLetterPayload contains the components of a cover letter provided in a request.
type CoverLetterPayload struct {
	CompanyProperName string                 `json:"companyProperName"`
	JobTitle          string                 `json:"jobTitle"`
	Body              CoverLetterPayloadBody `json:"body"`
}

// CoverLetterPayloadBody represents the paragraph-level content of a cover letter in a request.
type CoverLetterPayloadBody struct {
	About      string `json:"about,omitempty"`
	Experience string `json:"experience,omitempty"`
	WhatIBring string `json:"whatIBring,omitempty"`
	Paragraph  string `json:"paragraph,omitempty"`
}
