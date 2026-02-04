package llmErrors

import "fmt"

type LLMError struct {
	LLMProvider     string
	ErrorCode       int
	Err             error
	ProviderMessage string
}

func (e *LLMError) Error() string {
	if e.ProviderMessage != "" {
		return fmt.Sprintf("llm provider '%s' error: %v (provider message: %s)", e.LLMProvider, e.Err, e.ProviderMessage)
	}
	return fmt.Sprintf("llm provider '%s' error: %v", e.LLMProvider, e.Err)
}

func (e *LLMError) Unwrap() error {
	return e.Err
}

var (
	/* -- Config/Init --*/
	ErrInvalidAPIKey = fmt.Errorf("invalid_API_key")
	ErrFailedToInit  = fmt.Errorf("failed_to_init_llm_provider")

	/* -- API Requests and Input --*/
	ErrUnsupportedSchema = fmt.Errorf("unsupported_schema_type")
	ErrUnsupportedModel  = fmt.Errorf("unsupported_model")
	ErrInvalidProvider   = fmt.Errorf("invalid_llm_provider")

	/* -- API Responses and Network --*/
	ErrAuthenticationFailed = fmt.Errorf("authentication_failed")
	ErrRequestTimeout       = fmt.Errorf("request_timed_out")
	ErrServiceUnavailable   = fmt.Errorf("service_unavailable")
	ErrQuotaExceeded        = fmt.Errorf("quota_exceeded")
	ErrModelOverload        = fmt.Errorf("model_overloaded")

	/* -- Output/Content --*/
	ErrNoContent         = fmt.Errorf("no_content")
	ErrContentBlocked    = fmt.Errorf("content_blocked")
	ErrMalformedResponse = fmt.Errorf("malformed_response")
	ErrResponseNotText   = fmt.Errorf("response_part_not_text")
)
