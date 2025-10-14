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
	"github.com/ordo_meritum/features/documents/models/domain"
	"github.com/ordo_meritum/features/documents/models/events"
	"github.com/ordo_meritum/features/documents/models/requests"
	"github.com/ordo_meritum/features/documents/models/schemas"
	"github.com/ordo_meritum/shared/libs/llm"
	"github.com/ordo_meritum/shared/templates/instructions"
	"github.com/ordo_meritum/shared/templates/prompts"
	"github.com/ordo_meritum/shared/utils/formatters"
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

func (s *DocumentService) QueueResumeGeneration(
	ctx context.Context,
	apiKey string,
	requestBody requests.DocumentRequest,
	uid string,
) (int, error) {
	return s.queueDocumentGeneration(ctx, apiKey, requestBody, uid, "resume")
}

func (s *DocumentService) QueueCoverLetterGeneration(
	ctx context.Context,
	apiKey string,
	requestBody requests.DocumentRequest,
	uid string,
) (int, error) {
	return s.queueDocumentGeneration(ctx, apiKey, requestBody, uid, "cover-letter")
}

func (s *DocumentService) queueDocumentGeneration(
	ctx context.Context,
	apiKey string,
	requestBody requests.DocumentRequest,
	uid string,
	docType string,
) (int, error) {
	l := s.serviceLogger(uid, requestBody.Options.JobID, docType)
	l.Info().Msgf("Starting %s generation process", docType)

	// MOCK
	// kafkaRequest := mocks.GetMockDocumentEvent(uid, requestBody.Options.JobID, "cover-letter")
	var kafkaRequest *events.DocumentEvent
	var err error
	if docType == "resume" {
		kafkaRequest, err = s.updateResumeWithLLM(ctx, &requestBody, apiKey, uid)
	} else {
		currentResume, _ := s.resumeRepo.GetFullResume(ctx, uid, requestBody.Options.JobID)
		kafkaRequest, err = s.updateCoverLetterWithLLM(ctx, &requestBody, currentResume, apiKey, uid)
	}

	if err != nil {
		l.Error().Err(err).Msgf("Failed to update %s with LLM", docType)
		return 0, err
	}

	if err := s.sendKafkaMessage(ctx, kafkaRequest); err != nil {
		l.Error().Err(err).Msg("Error writing to Kafka")
		return 0, err
	}

	l.Info().Msgf("Successfully queued %s for compilation", docType)
	return kafkaRequest.JobID, nil
}

func (s *DocumentService) sendKafkaMessage(
	ctx context.Context,
	event *events.DocumentEvent,
) error {
	messageBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal Kafka request: %w", err)
	}

	kafkaCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = s.LatexWriter.WriteMessages(kafkaCtx, kafka.Message{
		Key:   []byte(strconv.Itoa(event.JobID)),
		Value: messageBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to write to kafka: %w", err)
	}
	return nil
}

func (s *DocumentService) updateResumeWithLLM(
	ctx context.Context,
	r *requests.DocumentRequest,
	apiKey string,
	uid string,
) (*events.DocumentEvent, error) {
	j, err := s.jobRepo.GetFullJobPosting(ctx, r.Options.JobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get job posting: %w", err)
	}

	promptData, err := buildResumePromptData(j, &r.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to build resume prompt data: %w", err)
	}

	var llmResume domain.Resume
	err = s.generateLLMContent(
		ctx,
		r.Options.LlmProvider,
		apiKey,
		"resume.txt",
		promptData,
		schemas.GeminiResumeSchema,
		&llmResume,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate resume from LLM: %w", err)
	}
	if err := s.resumeRepo.UpsertResume(ctx, uid, r.Options.JobID, &llmResume); err != nil {
		return nil, fmt.Errorf("failed to upsert resume: %w", err)
	}

	return &events.DocumentEvent{
		JobID:         r.Options.JobID,
		UserId:        uid,
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
	apiKey string,
	uid string,
) (*events.DocumentEvent, error) {
	jobID := r.Options.JobID
	j, err := s.jobRepo.GetFullJobPosting(ctx, jobID)
	if err != nil {
		return nil, err
	}

	llmProvider, err := llm.GetProvider(r.Options.LlmProvider)
	if err != nil {
		return nil, err
	}

	promptData, err := buildCoverLetterPromptData(j, &r.Payload, r.Options, currentResume)
	if err != nil {
		return nil, err
	}
	prompt, err := formatters.FormatTemplate(prompts.Prompts, "coverletter.txt", promptData)
	if err != nil {
		return nil, err
	}

	instructions, err := instructions.Instructions.ReadFile("coverletter.txt")
	if err != nil {
		return nil, err
	}

	rawResponse, err := llmProvider.Generate(
		ctx,
		string(instructions),
		prompt,
		schemas.GeminiCoverLetterSchema,
		apiKey,
	)
	if err != nil {
		return nil, err
	}

	cleanedJSON := llm.FormatLLMResponse(rawResponse)
	var llmCoverLetter domain.CoverLetterBody
	if err := json.Unmarshal([]byte(cleanedJSON), &llmCoverLetter); err != nil {
		return nil, fmt.Errorf("failed to unmarshal LLM resume response: %w. Raw response: %s", err, rawResponse)
	}

	coverLetterPayload := domain.CoverLetter{
		CompanyProperName: j.CompanyProperName,
		JobTitle:          j.JobTitle,
		Body:              llmCoverLetter,
	}

	load := events.DocumentEvent{
		JobID:         jobID,
		UserId:        uid,
		CompanyName:   j.CompanyName,
		DocType:       "resume",
		UserInfo:      r.Payload.UserInfo,
		EducationInfo: r.Payload.EducationInfo,
		CoverLetter:   coverLetterPayload,
	}
	return &load, nil
}

func (s *DocumentService) generateLLMContent(
	ctx context.Context,
	providerName, apiKey, instructionsFile string,
	promptData any,
	schema any,
	target interface{},
) error {
	llmProvider, err := llm.GetProvider(providerName)
	if err != nil {
		return err
	}

	prompt, err := formatters.FormatTemplate(prompts.Prompts, instructionsFile, promptData)
	if err != nil {
		return fmt.Errorf("failed to format prompt template: %w", err)
	}

	instructionBytes, err := instructions.Instructions.ReadFile(instructionsFile)
	if err != nil {
		return fmt.Errorf("failed to read instructions file: %w", err)
	}

	rawResponse, err := llmProvider.Generate(ctx, string(instructionBytes), prompt, schema, apiKey)
	if err != nil {
		return fmt.Errorf("LLM generation failed: %w", err)
	}

	cleanedJSON := llm.FormatLLMResponse(rawResponse)
	if err := json.Unmarshal([]byte(cleanedJSON), target); err != nil {
		return fmt.Errorf("failed to unmarshal LLM response: %w. Raw response: %s", err, rawResponse)
	}

	return nil
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

func buildResumePromptData(
	j *jobs.FullJobPosting,
	payload *requests.DocumentPayload,
) (map[string]any, error) {
	additionalInfo, err := formatters.FormatAboutForLLMWithXML(payload.AdditionalInfo)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"JobPost":        formatters.PrettyJobPost(*j),
		"Resume":         formatters.FormatResumeForLLMWithXML(payload.Resume),
		"AdditionalInfo": additionalInfo,
	}, nil
}

func buildCoverLetterPromptData(j *jobs.FullJobPosting, payload *requests.DocumentPayload, opts requests.DocumentOptions, resume *domain.Resume) (map[string]any, error) {
	additionalInfo, err := formatters.FormatAboutForLLMWithXML(payload.AdditionalInfo)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"JobPost":        formatters.PrettyJobPost(*j),
		"Education":      formatters.PrettyEducation(payload.EducationInfo),
		"Resume":         formatters.FormatResumePayloadForLLMWithXML(*resume),
		"AdditionalInfo": additionalInfo,
		"Corrections":    strings.Join(opts.Corrections, "\n- "),
		"WritingSamples": strings.Join(opts.WritingSamples, "\n- "),
	}, nil
}
