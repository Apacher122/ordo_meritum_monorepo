package llmErrors

import "fmt"

type LLMError struct {
	LLMProvider     string
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
	ErrInvalidAPIKey = fmt.Errorf("invalid API key")
	ErrFailedToInit  = fmt.Errorf("failed to initialize llm provider")

	/* -- API Requests and Input --*/
	ErrUnsupportedSchema = fmt.Errorf("unsupported schema type for llm generation")
	ErrUnsupportedModel  = fmt.Errorf("unsupported model for llm generation")
	ErrInvalidProvider   = fmt.Errorf("invalid llm provider. how did you do this???")

	/* -- API Responses and Network --*/
	ErrAuthenticationFailed = fmt.Errorf("authentication failed with llm provider")
	ErrRequestTimeout       = fmt.Errorf("request timed out")
	ErrServiceUnavailable   = fmt.Errorf("service is unavailable")
	ErrQuotaExceeded        = fmt.Errorf("quota exceeded")
	ErrModelOverload        = fmt.Errorf("model is overloaded")

	/* -- Output/Content --*/
	ErrNoContent         = fmt.Errorf("provider returned no content")
	ErrContentBlocked    = fmt.Errorf("content was blocked by content safety filters")
	ErrMalformedResponse = fmt.Errorf("malformed response from llm provider")
	ErrResponseNotText   = fmt.Errorf("response part was not of expected type TextPart")
)
