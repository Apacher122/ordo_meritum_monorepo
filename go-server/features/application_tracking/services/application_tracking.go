package services

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/rs/zerolog/log"

	"github.com/ordo_meritum/database/jobs"
	db_models "github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/features/application_tracking/models/dto"
	request "github.com/ordo_meritum/features/application_tracking/models/requests"
	pretty "github.com/ordo_meritum/features/application_tracking/utils"

	app_schemas "github.com/ordo_meritum/features/application_tracking/models/schemas"
	"github.com/ordo_meritum/shared/libs/llm"
	"github.com/ordo_meritum/shared/templates/instructions"
	prompts "github.com/ordo_meritum/shared/templates/prompts"
	"github.com/ordo_meritum/shared/utils/formatters"
)

type AppTrackerService struct {
	jobRepo jobs.Repository
}

func NewAppTrackerService(jobRepo jobs.Repository) *AppTrackerService {
	return &AppTrackerService{
		jobRepo: jobRepo,
	}
}

func (s *AppTrackerService) QueueApplicationTracking(
	ctx context.Context,
	apiKey string,
	uid string,
	requestBody request.JobPostingRequest,
) (any, error) {
	l := log.With().
		Str("service", "application-tracking").
		Str("uid", uid).
		Logger()

	l.Info().Msg("Starting application tracking process")

	l.Info().Msg("Extracting structured job info from content...")
	parsedJob, err := s.parseJobDescriptionWithLLM(
		ctx,
		&requestBody,
		apiKey,
	)

	if err != nil {
		l.Error().Err(err).Msg("could not extract job info from LLM")
		return nil, fmt.Errorf("could not extract job info from LLM: %w", err)
	}

	l.Info().Msg("Persisting full job posting to database...")
	res, err := s.jobRepo.InsertFullJobPosting(ctx, requestBody.JobDescription, parsedJob, uid)
	if err != nil {
		l.Error().Err(err).Msg("failed to insert full job posting")
		return nil, fmt.Errorf("failed to insert full job posting")
	}

	l.Info().Msg("Successfully tracked new job.")
	return res.ID, nil
}

func (s *AppTrackerService) GetTrackedApplicationByID(ctx context.Context, firebaseUID string, roleID int) (*jobs.FullJobPosting, error) {
	return s.jobRepo.GetFullJobPosting(ctx, roleID)
}

func (s *AppTrackerService) UpdateApplicationStatus(ctx context.Context, firebaseUID string, roleID int, status db_models.AppStatus) error {
	return s.jobRepo.UpdateApplicationDetails(ctx, roleID, firebaseUID, &status, nil)
}

func (s *AppTrackerService) ListTrackedApplications(ctx context.Context, firebaseUID string) ([]*jobs.UserJobPosting, error) {
	return s.jobRepo.GetAllUserJobPostings(ctx, firebaseUID)
}

func LoadPrompt(promptFile string, data map[string]string) (string, error) {
	promptBytes, err := prompts.Prompts.ReadFile(promptFile)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("prompt").Parse(string(promptBytes))
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, data); err != nil {
		return "", err
	}

	return out.String(), nil
}

func (s *AppTrackerService) parseJobDescriptionWithLLM(
	ctx context.Context,
	r *request.JobPostingRequest,
	apiKey string,
) (*dto.JobDescription, error) {
	llmProvider, err := llm.GetProvider("cohere")
	if err != nil {
		return nil, err
	}

	jobPost := pretty.FormatJobPostingRequest(r)
	promptData := map[string]string{
		"JobPost": jobPost,
	}
	prompt, err := formatters.FormatTemplate(prompts.Prompts, "jobInfoExtraction.txt", promptData)

	if err != nil {
		return nil, fmt.Errorf("failed to format prompt template: %w", err)
	}

	instructionBytes, err := instructions.Instructions.ReadFile("jobInfoExtraction.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to read instructions file: %w", err)
	}

	rawResponse, err := llmProvider.Generate(
		ctx,
		string(instructionBytes),
		prompt,
		app_schemas.JobDescriptionResponseFormat,
		apiKey,
	)

	if err != nil {
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}

	cleanedJSON := llm.FormatLLMResponse(rawResponse)
	log.Info().Msgf("Cleaned JSON: %s", cleanedJSON)
	var llmResponse dto.JobDescription
	if err := json.Unmarshal([]byte(cleanedJSON), &llmResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal LLM response for job info: %w. Raw response: %s", err, rawResponse)
	}

	return &llmResponse, nil
}
