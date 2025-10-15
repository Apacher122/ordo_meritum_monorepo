package error_response

import "fmt"

type ErrorResponse[T any] struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
	Details   T      `json:"details,omitempty"`
}

type ValidationDetail struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
}

type ResourceErrorDetail struct {
	ResourceType string `json:"resource_type"`
	ResourceID   any    `json:"resource_id"`
}

var (
	ErrNoUserID            = fmt.Errorf("no authenticated user found in context")
	ErrFailedUserId        = fmt.Errorf("failed to get user id")
	ErrFailedRequestBody   = fmt.Errorf("failed to get request body")
	ErrInternalServerError = fmt.Errorf("internal server error")
	ErrResourceUnavailable = fmt.Errorf("resource unavailable")
	ErrNoUserContext       = fmt.Errorf("no user found in context")
)

var (
	UNAUTHORIZED_USER     = "UNAUTHORIZED_USER"
	BAD_REQUEST           = "BAD_REQUEST"
	INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"
	RESOURCE_UNAVAILABLE  = "RESOURCE_UNAVAILABLE"
)
