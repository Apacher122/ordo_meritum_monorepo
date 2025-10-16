package domain

type CompanyInfo struct {
	CompanyName        string   `json:"companyName"`
	CompanyDescription string   `json:"companyDescription"`
	CompanyValues      []string `json:"companyValues"`
	CompanyMission     string   `json:"companyMission"`
	CompanyVision      string   `json:"companyVision"`
}
