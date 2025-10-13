package requests

type CoverLetterRequest struct {
	CompanyProperName string          `json:"companyProperName"`
	JobTitle          string          `json:"jobTitle"`
	Body              CoverLetterBody `json:"body"`
}

type CoverLetterBody struct {
	About      string `json:"about"`
	Experience string `json:"experience"`
	WhatIBring string `json:"whatIBring"`
}
