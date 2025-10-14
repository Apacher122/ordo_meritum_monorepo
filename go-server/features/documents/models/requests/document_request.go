package requests

import (
	"encoding/json"

	"github.com/ordo_meritum/shared/models/requests"
)

type DocumentRequest = requests.RequestBody[DocumentPayload, DocumentOptions]

type DocumentPayload struct {
	Resume         ResumePayload        `json:"resume,omitempty"`
	UserInfo       UserInfoPayload      `json:"userInfo"`
	AdditionalInfo json.RawMessage      `json:"additionalInfo"`
	EducationInfo  EducationInfoPayload `json:"educationInfo"`
	Coverletter    CoverLetterPayload   `json:"coverletter,omitempty"`
}

type DocumentOptions struct {
	JobID          int      `json:"jobId"`
	DocType        string   `json:"docType"`
	LlmProvider    string   `json:"llm"`
	LlmModel       string   `json:"llmModel"`
	GetNew         bool     `json:"getNew,omitempty"`
	Corrections    []string `json:"corrections,omitempty"`
	WritingSamples []string `json:"writingSamples,omitempty"`
}

/*
--- Request Payloads ---

These structs define the shape of the `payload` field in the main request body
sent from the client, and serve as the core data used to generate the document.

The names are pretty self-explanatory.
*/

type ResumePayload struct {
	Skills      []SkillsPayload     `json:"skills"`
	Experiences []ExperiencePayload `json:"experiences"`
	Projects    []ProjectPayload    `json:"projects"`
}

type SkillsPayload struct {
	Skill string `json:"skill"`
}

type ExperiencePayload struct {
	BulletPoints []string `json:"bulletPoints"`
	Company      string   `json:"company"`
	ID           string   `json:"id"`
	Position     string   `json:"position"`
	Years        string   `json:"years"`
}

type ProjectPayload struct {
	BulletPoints []string `json:"bulletPoints"`
	Description  string   `json:"description"`
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Years        string   `json:"years"`
}

type CoverLetterPayload struct {
	CompanyProperName string                 `json:"companyProperName"`
	JobTitle          string                 `json:"jobTitle"`
	Body              CoverLetterPayloadBody `json:"body"`
}

type CoverLetterPayloadBody struct {
	About      string `json:"about"`
	Experience string `json:"experience"`
	WhatIBring string `json:"whatIBring"`
}

type EducationInfoPayload struct {
	CourseWork string `json:"coursework"`
	Degree     string `json:"degree"`
	Location   string `json:"location"`
	School     string `json:"school"`
	StartEnd   string `json:"start_end"`
}

type UserInfoPayload struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	CurrentLocation string `json:"current_location"`
	Email           string `json:"email"`
	Github          string `json:"github,omitempty"`
	Linkedin        string `json:"linkedin,omitempty"`
	Mobile          string `json:"mobile,omitempty"`
	Summary         string `json:"summary,omitempty"`
}
