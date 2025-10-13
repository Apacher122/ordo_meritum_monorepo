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
	"github.com/ordo_meritum/features/documents/models/mocks"
	"github.com/ordo_meritum/features/documents/models/requests"
	"github.com/ordo_meritum/features/documents/models/schemas"
	"github.com/ordo_meritum/shared/libs/llm"
	"github.com/ordo_meritum/shared/templates/instructions"
	"github.com/ordo_meritum/shared/templates/prompts"
	"github.com/ordo_meritum/shared/utils/formatters"
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
	requestBody requests.RequestBody,
	uid string,
) (int, error) {
	l := log.With().
		Str("service", "documents-service").
		Str("documentType", "resume").
		Str("uid", uid).
		Int("jobId", requestBody.Options.JobID).
		Logger()

	l.Info().Msg("Starting resume generation process")

	kafkaRequest := mocks.GetMockDocumentEvent(uid, requestBody.Options.JobID, "resume")

	// kafkaRequest, err := s.updateResumeWithLLM(ctx, &requestBody, apiKey, uid)
	// if err != nil {
	// 	l.Error().Err(err).Msg("Failed to update resume with LLM")
	// 	return 0, fmt.Errorf("failed to update resume with LLM: %w", err)
	// }

	messageBytes, err := json.Marshal(kafkaRequest)
	if err != nil {
		l.Error().Err(err).Msg("Failed to marshal Kafka request")
		return 0, fmt.Errorf("failed to marshal Kafka request: %w", err)
	}

	kafkaCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = s.LatexWriter.WriteMessages(kafkaCtx, kafka.Message{
		Key:   []byte(strconv.Itoa(kafkaRequest.JobID)),
		Value: messageBytes,
	})
	if err != nil {
		l.Error().Err(err).Msg("Error writing to Kafka")
		return 0, fmt.Errorf("failed to write to kafka: %w", err)
	}

	l.Info().Msg("Successfully queued resume for compilation")
	return kafkaRequest.JobID, nil
}

func (s *DocumentService) QueueCoverLetterGeneration(
	ctx context.Context,
	apiKey string,
	requestBody requests.RequestBody,
	uid string,
) (int, error) {
	l := log.With().
		Str("service", "documents-service").
		Str("documentType", "cover-letter").
		Str("uid", uid).
		Int("jobId", requestBody.Options.JobID).
		Logger()

	l.Info().Msg("Starting cover letter generation process")

	// currentResume, err := s.resumeRepo.GetFullResume(ctx, uid, requestBody.Options.JobID)
	// if err != nil {
	// 	l.Error().Err(err).Msg("Failed to get latest resume")
	// 	return 0, fmt.Errorf("failed to get latest resume: %w", err)
	// }

	kafkaRequest := mocks.GetMockDocumentEvent(uid, requestBody.Options.JobID, "cover-letter")

	// kafkaRequest, err := s.updateCoverLetterWithLLM(
	// 	ctx,
	// 	&requestBody,
	// 	currentResume,
	// 	apiKey,
	// 	uid,
	// )

	// if err != nil {
	// 	l.Error().Err(err).Msg("Failed to update cover letter with LLM")
	// 	return 0, fmt.Errorf("failed to update cover letter with LLM: %w", err)
	// }

	messageBytes, err := json.Marshal(kafkaRequest)
	if err != nil {
		l.Error().Err(err).Msg("Failed to marshal Kafka request")
		return 0, fmt.Errorf("failed to marshal Kafka request: %w", err)
	}

	kafkaCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = s.LatexWriter.WriteMessages(kafkaCtx, kafka.Message{
		Key:   []byte(strconv.Itoa(kafkaRequest.JobID)),
		Value: messageBytes,
	})
	if err != nil {
		l.Error().Err(err).Msg("Error writing to Kafka")
		return 0, fmt.Errorf("failed to write to kafka: %w", err)
	}

	l.Info().Msg("Successfully queued cover letter for compilation")
	return kafkaRequest.JobID, nil
}

func (s *DocumentService) updateResumeWithLLM(
	ctx context.Context,
	r *requests.RequestBody,
	apiKey string,
	uid string,
) (*events.DocumentEvent, error) {
	jobID := r.Options.JobID
	j, err := s.jobRepo.GetFullJobPosting(ctx, jobID)
	if err != nil {
		return nil, err
	}

	llmProvider, err := llm.GetProvider(r.Options.LLM)
	if err != nil {
		return nil, err
	}

	jobPost := formatters.PrettyJobPost(*j)
	prettyresume := formatters.FormatResumeForLLMWithXML(r.Payload.Resume)
	additionalInfo, err := formatters.FormatAboutForLLMWithXML(r.Payload.AdditionalInfo)
	if err != nil {
		return nil, err
	}

	data := map[string]any{
		"JobPost":        jobPost,
		"Resume":         prettyresume,
		"AdditionalInfo": additionalInfo,
	}

	prompt, err := formatters.FormatTemplate(prompts.Prompts, "resume.txt", data)
	if err != nil {
		return nil, err
	}

	instructions, err := instructions.Instructions.ReadFile("resume.txt")
	if err != nil {
		return nil, err
	}

	rawResponse, err := llmProvider.Generate(
		ctx,
		string(instructions),
		prompt,
		schemas.GeminiResumeSchema,
		apiKey,
	)
	if err != nil {
		return nil, err
	}

	cleanedJSON := llm.FormatLLMResponse(rawResponse)

	var llmResume domain.Resume
	if err := json.Unmarshal([]byte(cleanedJSON), &llmResume); err != nil {
		return nil, fmt.Errorf("failed to unmarshal LLM resume response: %w. Raw response: %s", err, rawResponse)
	}
	s.resumeRepo.UpsertResume(ctx, uid, jobID, &llmResume)

	load := events.DocumentEvent{
		JobID:         jobID,
		UserId:        uid,
		CompanyName:   j.CompanyName,
		DocType:       "resume",
		UserInfo:      r.Payload.UserInfo,
		EducationInfo: r.Payload.EducationInfo,
		Resume:        llmResume,
	}
	return &load, nil
}

func (s *DocumentService) updateCoverLetterWithLLM(
	ctx context.Context,
	r *requests.RequestBody,
	currentResume *domain.Resume,
	apiKey string,
	uid string,
) (*events.DocumentEvent, error) {
	jobID := r.Options.JobID
	j, err := s.jobRepo.GetFullJobPosting(ctx, jobID)
	if err != nil {
		return nil, err
	}

	llmProvider, err := llm.GetProvider(r.Options.LLM)
	if err != nil {
		return nil, err
	}

	additionalInfo, err := formatters.FormatAboutForLLMWithXML(r.Payload.AdditionalInfo)
	if err != nil {
		return nil, err
	}

	promptData := map[string]any{
		"JobPost":        formatters.PrettyJobPost(*j),
		"Education":      formatters.PrettyEducation(r.Payload.EducationInfo),
		"Resume":         formatters.FormatResumePayloadForLLMWithXML(*currentResume),
		"AdditionalInfo": additionalInfo,
		"Corrections":    strings.Join(r.Options.Corrections, "\n- "),
		"WritingSamples": strings.Join(r.Options.WritingSamples, "\n- "),
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
