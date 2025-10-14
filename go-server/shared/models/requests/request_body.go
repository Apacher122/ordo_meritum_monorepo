package requests

type RequestBody[T any, U any] struct {
	Payload T `json:"payload"`
	Options U `json:"options,omitempty"`
}

// type Payload struct {
// 	UID            string                    `json:"uid"`
// 	JobID          int                       `json:"jobId,omitempty"`
// 	CompanyName    string                    `json:"companyName,omitempty"`
// 	Position       string                    `json:"position,omitempty"`
// 	JobDescription JobDescriptionRequest     `json:"jobDescription,omitempty"`
// 	Resume         ResumeRequest             `json:"resume,omitempty"`
// 	Coverletter    CoverLetterRequest        `json:"coverletter,omitempty"`
// 	UserInfo       UserInfoRequest           `json:"userInfo"`
// 	AdditionalInfo json.RawMessage           `json:"additionalInfo"`
// 	Corrections    []string                  `json:"corrections,omitempty"`
// 	WritingSamples []string                  `json:"writingSamples,omitempty"`
// 	Questionnaire  user.QuestionnaireRequest `json:"questionnaire,omitempty"`
// }

// type Options struct {
// 	DocType     string `json:"docType,omitempty"`
// 	LlmProvider string `json:"llm,omitempty"`
// 	LlmModel    string `json:"llmModel,omitempty"`
// 	GetNew      bool   `json:"getNew,omitempty"`
// }
