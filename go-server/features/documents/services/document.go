package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ordo_meritum/database/jobs"
	"github.com/ordo_meritum/database/resumes"
	apps_mappers "github.com/ordo_meritum/features/application_tracking/utils/mappers"
	"github.com/ordo_meritum/features/documents/models/domain"
	"github.com/ordo_meritum/features/documents/models/events"
	"github.com/ordo_meritum/features/documents/models/requests"
	"github.com/ordo_meritum/features/documents/utils/formatters"
	"github.com/ordo_meritum/shared/contexts"
	"github.com/ordo_meritum/shared/libs/llm"
	schemaregistry "github.com/ordo_meritum/shared/libs/llm/schema_registry"
	"github.com/ordo_meritum/shared/templates/instructions"
	"github.com/ordo_meritum/shared/templates/prompts"
	error_messages "github.com/ordo_meritum/shared/utils/errors"
	shared_formatters "github.com/ordo_meritum/shared/utils/formatters"
	lg "github.com/ordo_meritum/shared/utils/logger"

	"github.com/segmentio/kafka-go"
)

type DocumentService struct {
	jobRepo     jobs.Repository
	resumeRepo  resumes.Repository
	LatexWriter *kafka.Writer
}

func NewDocumentService(
	jobRepo jobs.Repository,
	resumeRepo resumes.Repository,
	latexWriter *kafka.Writer,
) *DocumentService {
	return &DocumentService{
		jobRepo:     jobRepo,
		resumeRepo:  resumeRepo,
		LatexWriter: latexWriter,
	}
}

var service = "documents-service"

// QueueDocumentGeneration orchestrates the asynchronous generation of a document.
//
// This function acts as the central router for creating different types of documents.
// It follows these steps:
//  1. Validates the user context from the request.
//  2. Determines the required document type (e.g., "resume", "cover-letter")
//     from the request body.
//  3. Dispatches the request to the appropriate internal method for content
//     generation with an LLM (e.g., updateResumeWithLLM).
//  4. Consolidates error handling from the generation process.
//  5. If generation is successful, it queues the resulting document event for
//     final compilation by sending it as a message to a Kafka topic.
//
// It returns the jobID for the queued document and an error if any part of the
// process fails.
func (s *DocumentService) QueueDocumentGeneration(
	ctx context.Context,
	requestBody *requests.DocumentRequest,
) (int, error) {
	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_USER_NO_CONTEXT}.ErrorLog()
		return 0, error_messages.ErrorMessage(error_messages.ERR_USER_NO_CONTEXT)
	}

	lg.InfoLoggerType{Service: &service, Uid: &userCtx.UID, Message: fmt.Sprintf("Starting %s generation process", requestBody.Options.DocType)}.InfoLog()

	var kafkaRequest *events.DocumentEvent
	var err error

	switch requestBody.Options.DocType {
	case "resume":
		kafkaRequest, err = s.updateResumeWithLLM(ctx, requestBody)
	case "cover-letter":
		var currentResume *domain.Resume
		currentResume, err = s.resumeRepo.GetFullResume(ctx, requestBody.Options.JobID)
		if err == nil {
			kafkaRequest, err = s.updateCoverLetterWithLLM(ctx, requestBody, currentResume)
		}
	default:
		err = fmt.Errorf("unsupported document type for generation: %s", requestBody.Options.DocType)
	}

	if err != nil {
		lg.ErrorLoggerType{
			Service:   &service,
			ErrorCode: &error_messages.ERR_PROCESS_FAILED,
			Error:     err,
			JobID:     &requestBody.Options.JobID,
			Uid:       &userCtx.UID,
		}.ErrorLog()
	}

	if err := s.sendKafkaMessage(ctx, kafkaRequest); err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_KAFKA_FAILED_TO_WRITE, Error: err}.ErrorLog()
		return 0, err
	}

	lg.InfoLoggerType{Service: &service, JobID: &kafkaRequest.JobID, Message: fmt.Sprintf("Successfully queued %s for compilation", requestBody.Options.DocType)}.InfoLog()
	return kafkaRequest.JobID, nil
}

// sendKafkaMessage marshals a DocumentEvent and sends it to the configured Kafka topic.
// It uses a context with a 10-second timeout for the write operation to prevent indefinite blocking.
// Returns an error if marshalling fails or the message cannot be written to Kafka.
func (s *DocumentService) sendKafkaMessage(
	ctx context.Context,
	event *events.DocumentEvent,
) error {
	messageBytes, err := json.Marshal(event)
	if err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_KAFKA_MALFORMED_RESPONSE, Error: fmt.Errorf("failed to marshal Kafka request: %w", err)}.ErrorLog()
		return fmt.Errorf("failed to marshal Kafka request: %w", err)
	}

	kafkaCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = s.LatexWriter.WriteMessages(kafkaCtx, kafka.Message{
		Key:   []byte(strconv.Itoa(event.JobID)),
		Value: messageBytes,
	})
	if err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_KAFKA_FAILED_TO_WRITE, Error: fmt.Errorf("failed to write to kafka: %w", err)}.ErrorLog()
		return fmt.Errorf("failed to write to kafka: %w", err)
	}
	return nil
}

// updateResumeWithLLM handles the process of generating new resume content using an LLM.
// It fetches job posting data, prepares a prompt, calls the LLM to get a structured response,
// and then persists the newly generated resume content to the database.
// It returns a DocumentEvent suitable for Kafka queuing or an error if any step fails.
func (s *DocumentService) updateResumeWithLLM(
	ctx context.Context,
	r *requests.DocumentRequest,
) (*events.DocumentEvent, error) {
	userCtx, j, err := s.prepareGenerationData(ctx, r.Options.JobID)
	if err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_DB_FAILED_TO_GET, Error: err}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_DB_FAILED_TO_GET)
	}

	promptData, err := buildResumePromptData(j, &r.Payload)
	if err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_LLM_PROMPT_FORMATTING, Error: err}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_LLM_PROMPT_FORMATTING)
	}

	e := r.Payload.EducationInfo
	education, err := formatters.NewEducationInfoFromPayload(&e)
	if err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_INVALID_REQUEST_FORMAT, Error: err}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_INVALID_REQUEST_FORMAT)
	}

	var llmResume domain.Resume
	errBody := s.generateLLMContent(
		ctx,
		r,
		"resume.txt",
		promptData,
		schemaregistry.Resume,
		&llmResume,
	)
	if errBody != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &errBody.ErrCode, Error: errBody.ErrMsg}.ErrorLog()
		return nil, errBody.ErrMsg
	}

	if err := s.resumeRepo.UpsertResume(ctx, r.Options.JobID, &llmResume, education); err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_DB_FAILED_TO_UPSERT, Error: err}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_DB_FAILED_TO_UPSERT)
	}

	lg.InfoLoggerType{Service: &service, Uid: &userCtx.UID, JobID: &r.Options.JobID, Message: llmResume.FormatForLLM()}.InfoLog()

	return &events.DocumentEvent{
		JobID:         r.Options.JobID,
		UserId:        userCtx.UID,
		CompanyName:   j.CompanyName,
		DocType:       "resume",
		UserInfo:      r.Payload.UserInfo,
		EducationInfo: r.Payload.EducationInfo,
		Resume:        llmResume,
	}, nil
}

// updateCoverLetterWithLLM handles the process of generating a cover letter using an LLM.
// It fetches job data and the user's current resume to build a comprehensive prompt,
// calls the LLM, and formats the response.
// It returns a DocumentEvent ready for Kafka queuing or an error if any step fails.
func (s *DocumentService) updateCoverLetterWithLLM(
	ctx context.Context,
	r *requests.DocumentRequest,
	currentResume *domain.Resume,
) (*events.DocumentEvent, error) {
	userCtx, j, err := s.prepareGenerationData(ctx, r.Options.JobID)
	if err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_DB_FAILED_TO_GET, Error: err}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_DB_FAILED_TO_GET)
	}

	promptData, err := buildCoverLetterPromptData(j, &r.Payload, &r.Options, currentResume)
	if err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_LLM_PROMPT_FORMATTING, Error: err}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_LLM_PROMPT_FORMATTING)
	}

	var llmCoverLetter domain.CoverLetterBody
	errBody := s.generateLLMContent(
		ctx,
		r,
		"coverletter.txt",
		promptData,
		schemaregistry.Coverletter,
		&llmCoverLetter,
	)
	if errBody != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &errBody.ErrCode, Error: errBody.ErrMsg}.ErrorLog()
		return nil, errBody.ErrMsg
	}

	coverLetterPayload := domain.CoverLetter{
		CompanyProperName: j.CompanyProperName,
		JobTitle:          j.JobTitle,
		Body:              llmCoverLetter,
	}

	load := events.DocumentEvent{
		JobID:       r.Options.JobID,
		UserId:      userCtx.UID,
		CompanyName: j.CompanyName,
		DocType:     "cover-letter",
		UserInfo:    r.Payload.UserInfo,
		CoverLetter: coverLetterPayload,
	}
	return &load, nil
}

// generateLLMContent provides a generic interface for interacting with an LLM.
// It retrieves the specified LLM provider, formats the prompt and instructions,
// generates content, and unmarshals the structured JSON response into the target interface.
//
// It returns a custom ErrorBody pointer which contains a specific error code and message,
// allowing the calling function to log and handle errors with more context.
// A nil return value indicates success.
func (s *DocumentService) generateLLMContent(
	ctx context.Context,
	r *requests.DocumentRequest,
	instructionsFile string,
	promptData any,
	schemaType string,
	target interface{},
) *error_messages.ErrorBody {
	_, ok := contexts.FromContext(ctx)
	if !ok {
		return &error_messages.ErrorBody{ErrCode: error_messages.ERR_USER_NO_CONTEXT, ErrMsg: error_messages.ErrorMessage(error_messages.ERR_USER_NO_CONTEXT)}
	}

	llmProvider, err := llm.GetProvider(r.Options.LlmProvider)
	if err != nil {
		return &error_messages.ErrorBody{ErrCode: error_messages.ERR_LLM_NO_CONTENT, ErrMsg: err}
	}

	schema, err := schemaregistry.GetSchema(r.Options.LlmProvider, schemaType, r.Payload.Resume.Experiences)
	if err != nil {
		return &error_messages.ErrorBody{ErrCode: error_messages.ERR_INVALID_SCHEMA, ErrMsg: err}
	}

	prompt, err := shared_formatters.FormatTemplate(prompts.Prompts, instructionsFile, promptData)
	if err != nil {
		return &error_messages.ErrorBody{ErrCode: error_messages.ERR_LLM_PROMPT_FORMATTING, ErrMsg: fmt.Errorf("failed to format prompt template: %w", err)}
	}

	instructionBytes, err := instructions.Instructions.ReadFile(instructionsFile)
	if err != nil {
		return &error_messages.ErrorBody{ErrCode: error_messages.ERR_LLM_INSTRUCTION_FORMATTING, ErrMsg: fmt.Errorf("failed to read instructions file: %w", err)}
	}

	rawResponse, err := llmProvider.Generate(ctx, string(instructionBytes), prompt, schema)
	if err != nil {
		return &error_messages.ErrorBody{ErrCode: error_messages.ERR_LLM_NO_CONTENT, ErrMsg: fmt.Errorf("LLM generation failed: %w", err)}
	}

	cleanedJSON := llm.FormatLLMResponse(rawResponse)
	if err := json.Unmarshal([]byte(cleanedJSON), target); err != nil {
		return &error_messages.ErrorBody{ErrCode: error_messages.ERR_LLM_MALFORMED_RESPONSE, ErrMsg: fmt.Errorf("failed to unmarshal LLM response: %w. Raw response: %s", err, rawResponse)}
	}

	return nil
}

// prepareGenerationData prepares the necessary data for document generation.
// It retrieves the user context and a full job posting from the database.
// If either the user context or job posting cannot be fetched, it returns an error.
// The returned user context and job posting are used to generate the document content using an LLM.
//
// It returns a pointer to a UserContext, a pointer to a FullJobPosting, and an error if any part of the process fails.
func (s *DocumentService) prepareGenerationData(ctx context.Context, jobID int) (*contexts.UserContext, *jobs.FullJobPosting, error) {
	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		return nil, nil, error_messages.ErrorMessage(error_messages.ERR_USER_NO_CONTEXT)
	}

	jobPosting, err := s.jobRepo.GetFullJobPosting(ctx, jobID)
	if err != nil {
		return nil, nil, error_messages.ErrorMessage(error_messages.ERR_DB_FAILED_TO_GET)
	}

	return userCtx, jobPosting, nil
}

// buildResumePromptData constructs the data map needed to populate the resume generation prompt.
// It combines formatted data from the job posting and the user's document request payload.
func buildResumePromptData(
	j *jobs.FullJobPosting,
	payload *requests.DocumentPayload,
) (map[string]any, error) {
	additionalInfo, err := shared_formatters.FormatAboutForLLMWithXML(payload.AdditionalInfo)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"JobPost":        shared_formatters.FormatJobPostForLLM(*j),
		"Resume":         formatters.FormatResumeRequestForLLMWithXML(payload),
		"AdditionalInfo": additionalInfo,
	}, nil
}

// buildCoverLetterPromptData constructs the data map for the cover letter generation prompt.
// It aggregates data from the job posting, user payload, existing resume, and other options
// to provide the LLM with comprehensive context.
func buildCoverLetterPromptData(j *jobs.FullJobPosting, payload *requests.DocumentPayload, opts *requests.DocumentOptions, resume *domain.Resume) (map[string]any, error) {
	additionalInfo := ""
	var err error
	if payload.AdditionalInfo != nil {
		additionalInfo, err = shared_formatters.FormatAboutForLLMWithXML(payload.AdditionalInfo)

	}
	if err != nil {
		return nil, err
	}
	jobPost := apps_mappers.NewJobDescriptionFromPost(j)
	return map[string]any{
		"JobPost":        jobPost.FormatForLLM(),
		"Education":      payload.EducationInfo.FormatForLLM(),
		"Resume":         resume.FormatForLLM(),
		"AdditionalInfo": additionalInfo,
		"Corrections":    strings.Join(opts.Corrections, "\n- "),
		"WritingSamples": strings.Join(opts.WritingSamples, "\n- "),
	}, nil
}
