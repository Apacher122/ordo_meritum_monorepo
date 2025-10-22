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

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func (s *DocumentService) QueueResumeGeneration(
	ctx context.Context,
	requestBody requests.DocumentRequest,
) (int, error) {
	return s.queueDocumentGeneration(ctx, requestBody, "resume")
}

func (s *DocumentService) QueueCoverLetterGeneration(
	ctx context.Context,
	requestBody requests.DocumentRequest,
) (int, error) {
	return s.queueDocumentGeneration(ctx, requestBody, "cover-letter")
}

func (s *DocumentService) queueDocumentGeneration(
	ctx context.Context,
	requestBody requests.DocumentRequest,
	docType string,
) (int, error) {
	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_USER_NO_CONTEXT}.ErrorLog()
		return 0, error_messages.ErrorMessage(error_messages.ERR_USER_NO_CONTEXT)
	}
	l := s.serviceLogger(userCtx.UID, requestBody.Options.JobID, docType)
	l.Info().Msgf("Starting %s generation process", docType)

	var kafkaRequest *events.DocumentEvent
	var err error
	if docType == "resume" {
		kafkaRequest, err = s.updateResumeWithLLM(ctx, &requestBody)
		if err != nil {
			return 0, err
		}
	} else {
		currentResume, err := s.resumeRepo.GetFullResume(ctx, requestBody.Options.JobID)
		if err != nil {
			return 0, err
		}

		kafkaRequest, err = s.updateCoverLetterWithLLM(ctx, &requestBody, currentResume)
		if err != nil {
			return 0, err
		}
	}

	if err := s.sendKafkaMessage(ctx, kafkaRequest); err != nil {
		return 0, err
	}

	lg.InfoLoggerType{Service: &service, JobID: &kafkaRequest.JobID, Message: fmt.Sprintf("Successfully queued %s for compilation", docType)}.InfoLog()
	return kafkaRequest.JobID, nil
}

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

func (s *DocumentService) updateResumeWithLLM(
	ctx context.Context,
	r *requests.DocumentRequest,
) (*events.DocumentEvent, error) {

	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_USER_NO_CONTEXT}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_USER_NO_CONTEXT)
	}

	j, err := s.jobRepo.GetFullJobPosting(ctx, r.Options.JobID)
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
	err = s.generateLLMContent(
		ctx,
		r,
		"resume.txt",
		promptData,
		schemaregistry.Resume,
		&llmResume,
	)
	if err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_LLM_NO_CONTENT, Error: err}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_LLM_NO_CONTENT)
	}

	if err := s.resumeRepo.UpsertResume(ctx, r.Options.JobID, &llmResume, education); err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_DB_FAILED_TO_UPSERT, Error: err}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_DB_FAILED_TO_UPSERT)
	}

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

func (s *DocumentService) updateCoverLetterWithLLM(
	ctx context.Context,
	r *requests.DocumentRequest,
	currentResume *domain.Resume,
) (*events.DocumentEvent, error) {
	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_USER_NO_CONTEXT}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_USER_NO_CONTEXT)
	}

	jobID := r.Options.JobID
	j, err := s.jobRepo.GetFullJobPosting(ctx, jobID)
	if err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_DB_FAILED_TO_GET, Error: err}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_DB_FAILED_TO_GET)
	}

	llmProvider, err := llm.GetProvider(r.Options.LlmProvider)
	if err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_LLM_NO_CONTENT, Error: err}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_LLM_NO_CONTENT)
	}

	schema, err := schemaregistry.GetSchema(r.Options.LlmProvider, schemaregistry.Coverletter, nil)
	if err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_INVALID_SCHEMA, Error: err}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_INVALID_SCHEMA)
	}

	promptData, err := buildCoverLetterPromptData(j, &r.Payload, r.Options, currentResume)
	if err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_LLM_PROMPT_FORMATTING, Error: err}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_LLM_PROMPT_FORMATTING)
	}
	prompt, err := shared_formatters.FormatTemplate(prompts.Prompts, "coverletter.txt", promptData)
	if err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_LLM_PROMPT_FORMATTING, Error: err}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_LLM_PROMPT_FORMATTING)
	}

	instructions, err := instructions.Instructions.ReadFile("coverletter.txt")
	if err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_LLM_INSTRUCTION_FORMATTING, Error: err}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_LLM_INSTRUCTION_FORMATTING)
	}

	rawResponse, err := llmProvider.Generate(
		ctx,
		string(instructions),
		prompt,
		schema,
	)
	if err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_LLM_NO_CONTENT, Error: err}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_LLM_NO_CONTENT)
	}

	cleanedJSON := llm.FormatLLMResponse(rawResponse)
	var llmCoverLetter domain.CoverLetterBody
	if err := json.Unmarshal([]byte(cleanedJSON), &llmCoverLetter); err != nil {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_LLM_MALFORMED_RESPONSE, Error: err}.ErrorLog()
		return nil, error_messages.ErrorMessage(error_messages.ERR_LLM_MALFORMED_RESPONSE)
	}

	coverLetterPayload := domain.CoverLetter{
		CompanyProperName: j.CompanyProperName,
		JobTitle:          j.JobTitle,
		Body:              llmCoverLetter,
	}

	load := events.DocumentEvent{
		JobID:         jobID,
		UserId:        userCtx.UID,
		CompanyName:   j.CompanyName,
		DocType:       "cover-letter",
		UserInfo:      r.Payload.UserInfo,
		EducationInfo: r.Payload.EducationInfo,
		CoverLetter:   coverLetterPayload,
	}
	return &load, nil
}

func (s *DocumentService) generateLLMContent(
	ctx context.Context,
	r *requests.DocumentRequest,
	instructionsFile string,
	promptData any,
	schemaType string,
	target interface{},
) error {
	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		lg.ErrorLoggerType{Service: &service, ErrorCode: &error_messages.ERR_USER_NO_CONTEXT}.ErrorLog()
		return error_messages.ErrorMessage(error_messages.ERR_USER_NO_CONTEXT)
	}

	llmProvider, err := llm.GetProvider(r.Options.LlmProvider)
	if err != nil {
		return err
	}
	schema, err := schemaregistry.GetSchema(r.Options.LlmProvider, schemaType, r.Payload.Resume.Experiences)
	if err != nil {
		return err
	}

	prompt, err := shared_formatters.FormatTemplate(prompts.Prompts, instructionsFile, promptData)
	if err != nil {
		lg.ErrorLoggerType{Uid: &userCtx.UID, Service: &service, ErrorCode: &error_messages.ERR_LLM_PROMPT_FORMATTING, Error: fmt.Errorf("failed to format prompt template: %w", err)}.ErrorLog()
		return fmt.Errorf("failed to format prompt template: %w", err)
	}

	instructionBytes, err := instructions.Instructions.ReadFile(instructionsFile)
	if err != nil {
		lg.ErrorLoggerType{Uid: &userCtx.UID, Service: &service, ErrorCode: &error_messages.ERR_LLM_INSTRUCTION_FORMATTING, Error: fmt.Errorf("failed to read instructions file: %w", err)}.ErrorLog()
		return fmt.Errorf("failed to read instructions file: %w", err)
	}

	rawResponse, err := llmProvider.Generate(ctx, string(instructionBytes), prompt, schema)
	if err != nil {
		lg.ErrorLoggerType{Uid: &userCtx.UID, Service: &service, ErrorCode: &error_messages.ERR_LLM_NO_CONTENT, Error: fmt.Errorf("LLM generation failed: %w", err)}.ErrorLog()
		return fmt.Errorf("LLM generation failed: %w", err)
	}

	log.Printf("LLM response: %s", rawResponse)

	cleanedJSON := llm.FormatLLMResponse(rawResponse)
	if err := json.Unmarshal([]byte(cleanedJSON), target); err != nil {
		lg.ErrorLoggerType{Uid: &userCtx.UID, Service: &service, ErrorCode: &error_messages.ERR_LLM_MALFORMED_RESPONSE, Error: fmt.Errorf("failed to unmarshal LLM response: %w. Raw response: %s", err, rawResponse)}.ErrorLog()
		return fmt.Errorf("failed to unmarshal LLM response: %w. Raw response: %s", err, rawResponse)
	}

	return nil
}

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

func buildCoverLetterPromptData(j *jobs.FullJobPosting, payload *requests.DocumentPayload, opts requests.DocumentOptions, resume *domain.Resume) (map[string]any, error) {
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

func (s *DocumentService) serviceLogger(
	uid string,
	jobID int,
	docType string,
) zerolog.Logger {
	return log.With().
		Str("service", "documents-service").
		Str("uid", uid).
		Int("jobID", jobID).
		Str("docType", docType).
		Logger()
}
