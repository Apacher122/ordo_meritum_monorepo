package types

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type HttpErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
