package domain

// CoverLetter represents the main components of a cover letter.
type CoverLetter struct {
	CompanyProperName string          `json:"companyProperName"`
	JobTitle          string          `json:"jobTitle"`
	Body              CoverLetterBody `json:"body"`
}

// CoverLetterBody contains the distinct paragraphs that make up the
// main content of the cover letter.
type CoverLetterBody struct {
	About           string `json:"about"`
	Experience      string `json:"experience"`
	WhatIBring      string `json:"whatIBring"`
	RevisionSummary string `json:"revisionSummary"`
}
