package requests

// UserInfoPayload holds the personal and contact details of a user as
// provided in a request payload.
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
