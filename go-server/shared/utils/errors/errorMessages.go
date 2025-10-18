package error_messages

import (
	"fmt"
)

type ErrorBody struct {
	ErrCode string
	ErrMsg  error
}

var (
	ERR_LLM_INVALID_API_KEY        = "ERR_LLM_INVALID_API_KEY"
	ERR_LLM_FAILED_TO_INIT         = "ERR_LLM_FAILED_TO_INIT"
	ERR_LLM_UNSUPPORTED_SCHEMA     = "ERR_LLM_UNSUPPORTED_SCHEMA"
	ERR_LLM_UNSUPPORTED_MODEL      = "ERR_LLM_UNSUPPORTED_MODEL"
	ERR_LLM_INVALID_PROVIDER       = "ERR_LLM_INVALID_PROVIDER"
	ERR_LLM_AUTHENTICATION         = "ERR_LLM_AUTHENTICATION"
	ERR_LLM_REQUEST_TIMEOUT        = "ERR_LLM_REQUEST_TIMEOUT"
	ERR_LLM_SERVICE_UNAVAILABLE    = "ERR_LLM_SERVICE_UNAVAILABLE"
	ERR_LLM_QUOTA_EXCEEDED         = "ERR_LLM_QUOTA_EXCEEDED"
	ERR_LLM_MODEL_OVERLOADED       = "ERR_LLM_MODEL_OVERLOADED"
	ERR_LLM_NO_CONTENT             = "ERR_LLM_NO_CONTENT "
	ERR_LLM_CONTENT_BLOCKED        = "ERR_LLM_CONTENT_BLOCKED"
	ERR_LLM_MALFORMED_RESPONSE     = "ERR_LLM_MALFORMED_RESPONSE"
	ERR_LLM_RESPONSE_NOT_TEXT      = "ERR_LLM_RESPONSE_NOT_TEXT"
	ERR_LLM_PROMPT_FORMATTING      = "ERR_TEMPLATE_FORMATTING"
	ERR_LLM_INSTRUCTION_FORMATTING = "ERR_LLM_INSTRUCTION_FORMATTING"

	ERR_DB_FAILED_TO_INSERT  = "ERR_DB_FAILED_TO_INSERT"
	ERR_DB_FAILED_TO_UPSERT  = "ERR_DB_FAILED_TO_UPSERT"
	ERR_DB_FAILED_TO_GET     = "ERR_DB_FAILED_TO_GET"
	ERR_DB_FAILED_TO_UPDATE  = "ERR_DB_FAILED_TO_UPDATE"
	ERR_DB_FAILED_TO_DELETE  = "ERR_DB_FAILED_TO_DELETE"
	ERR_DB_FAILED_TO_CONNECT = "ERR_DB_FAILED_TO_CONNECT"

	ERR_USER_NOT_AUTHORIZED = "ERR_USER_NOT_AUTHORIZED"
	ERR_USER_NO_CONTEXT     = "ERR_USER_NO_CONTEXT"

	ERR_INVALID_REQUEST_FORMAT = "ERR_INVALID_REQUEST_FORMAT"
	ERR_INVALID_SCHEMA         = "ERR_INVALID_SCHEMA"

	ERR_KAFKA_FAILED_TO_WRITE = "ERR_KAFKA_FAILED_TO_WRITE"
	ERR_KAFKA_INVALID_REQUEST = "ERR_KAFKA_INVALID_REQUEST"
)

var (
	ErrLlmResponseFailed   = "llm provider failed to respond"
	ErrLlmOutputFail       = "llm provider failed to provide output"
	ErrPromptTemplate      = "failed to format prompt template"
	ErrInstructionTemplate = "failed to format instruction template"
)

func ErrorMessage(msg string) error {
	switch msg {
	case ERR_LLM_INVALID_API_KEY:
		return fmt.Errorf("invalid API key")
	case ERR_LLM_FAILED_TO_INIT:
		return fmt.Errorf("failed to initialize llm provider")
	case ERR_LLM_UNSUPPORTED_SCHEMA:
		return fmt.Errorf("unsupported schema type for llm generation")
	case ERR_LLM_UNSUPPORTED_MODEL:
		return fmt.Errorf("unsupported model for llm generation")
	case ERR_LLM_INVALID_PROVIDER:
		return fmt.Errorf("invalid llm provider. how did you do this???")
	case ERR_LLM_AUTHENTICATION:
		return fmt.Errorf("authentication failed with llm provider")
	case ERR_LLM_REQUEST_TIMEOUT:
		return fmt.Errorf("request timed out")
	case ERR_LLM_SERVICE_UNAVAILABLE:
		return fmt.Errorf("service is unavailable")
	case ERR_LLM_QUOTA_EXCEEDED:
		return fmt.Errorf("quota exceeded")
	case ERR_LLM_MODEL_OVERLOADED:
		return fmt.Errorf("model is overloaded")
	case ERR_LLM_NO_CONTENT:
		return fmt.Errorf("provider returned no content")
	case ERR_LLM_CONTENT_BLOCKED:
		return fmt.Errorf("content was blocked by content safety filters")
	case ERR_LLM_MALFORMED_RESPONSE:
		return fmt.Errorf("malformed response from llm provider")
	case ERR_LLM_RESPONSE_NOT_TEXT:
		return fmt.Errorf("response part was not of expected type TextPart")
	case ERR_LLM_PROMPT_FORMATTING:
		return fmt.Errorf("failed to format prompt template")
	case ERR_LLM_INSTRUCTION_FORMATTING:
		return fmt.Errorf("failed to format instruction template")

	case ERR_DB_FAILED_TO_INSERT:
		return fmt.Errorf("failed to insert information to db")
	case ERR_DB_FAILED_TO_DELETE:
		return fmt.Errorf("failed to delete information from db")
	case ERR_DB_FAILED_TO_GET:
		return fmt.Errorf("failed to get information from db")
	case ERR_DB_FAILED_TO_UPDATE:
		return fmt.Errorf("failed to update information in db")
	case ERR_DB_FAILED_TO_UPSERT:
		return fmt.Errorf("failed to upsert to db")
	case ERR_DB_FAILED_TO_CONNECT:
		return fmt.Errorf("failed to connect to db")

	case ERR_USER_NOT_AUTHORIZED:
		return fmt.Errorf("unauthorized access attempt")
	case ERR_USER_NO_CONTEXT:
		return fmt.Errorf("user does not exist in this context")

	case ERR_INVALID_REQUEST_FORMAT:
		return fmt.Errorf("invalid request format")
	case ERR_INVALID_SCHEMA:
		return fmt.Errorf("invalid schema")

	case ERR_KAFKA_FAILED_TO_WRITE:
		return fmt.Errorf("failed to write to kafka")
	case ERR_KAFKA_INVALID_REQUEST:
		return fmt.Errorf("kafka request is invalid")

	default:
		return fmt.Errorf("unknown error")
	}
}
