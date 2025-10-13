package models

type JobDescription struct {
	ID               string   `json:"id,omitempty"`
	Company          string   `json:"company"`
	Role             string   `json:"role"`
	Description      string   `json:"description,omitempty"`
	Requirements     []string `json:"requirements"`
	Responsibilities []string `json:"responsibilities"`
	Location         string   `json:"location"`
	Salary           string   `json:"salary"`
	Status           string   `json:"status,omitempty"`
}
