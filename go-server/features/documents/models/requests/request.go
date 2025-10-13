package requests

import (
	"encoding/json"
)

type RequestBody struct {
	Payload Payload `json:"payload"`
	Options Options `json:"options"`
}

type Payload struct {
	Resume         ResumeRequest        `json:"resume,omitempty"`
	Coverletter    CoverLetterRequest   `json:"coverletter,omitempty"`
	UserInfo       UserInfoRequest      `json:"userInfo"`
	AdditionalInfo json.RawMessage      `json:"additionalInfo"`
	EducationInfo  EducationInfoRequest `json:"educationInfo"`
}

type Options struct {
	JobID          int      `json:"jobId"`
	DocType        string   `json:"docType"`
	LLM            string   `json:"llm"`
	GetNew         bool     `json:"getNew,omitempty"`
	Corrections    []string `json:"corrections,omitempty"`
	WritingSamples []string `json:"writingSamples,omitempty"`
}
