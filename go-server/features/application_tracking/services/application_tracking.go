package services

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/ordo_meritum/database/jobs"
	db_models "github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/features/application_tracking/models/dto"
	request "github.com/ordo_meritum/features/application_tracking/models/requests"

	app_schemas "github.com/ordo_meritum/features/application_tracking/models/schemas"
	"github.com/ordo_meritum/shared/embeds"
	"github.com/ordo_meritum/shared/libs/llm"
	"github.com/ordo_meritum/shared/templates/instructions"
	prompts "github.com/ordo_meritum/shared/templates/prompts"
	formatters "github.com/ordo_meritum/shared/utils/formatters/pretty"
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

	parsedJob, err := s.parseJobDescriptionWithLLM(
		ctx,
		&requestBody,
		apiKey,
	)

	if err != nil {
		l.Error().Err(err).Msg("could not extract job info from LLM")
		return nil, err
	}

	l.Info().Msg("Persisting full job posting to database...")
	res, err := s.jobRepo.InsertFullJobPosting(ctx, requestBody.JobDescription, parsedJob, uid)
	if err != nil {
		l.Error().Err(err).Msg("failed to insert full job posting")
		return nil, err
	}

	l.Info().Msg("Successfully tracked new job.")
	return res.ID, nil
}

func (s *AppTrackerService) GetTrackedApplicationByID(
	ctx context.Context,
	roleID int,
) (*jobs.FullJobPosting, error) {
	return s.jobRepo.GetFullJobPosting(ctx, roleID)
}

func (s *AppTrackerService) UpdateApplicationStatus(
	ctx context.Context,
	firebaseUID string,
	roleID int,
	status db_models.AppStatus,
) error {
	return s.jobRepo.UpdateApplicationDetails(ctx, roleID, firebaseUID, &status, nil)
}

func (s *AppTrackerService) ListTrackedApplications(
	ctx context.Context,
	firebaseUID string,
) ([]*jobs.UserJobPosting, error) {
	return s.jobRepo.GetAllUserJobPostings(ctx, firebaseUID)
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

	jobPost := request.FormatJobPostingRequest(r)
	promptData := map[string]string{
		"JobPost": jobPost,
	}
	prompt, err := formatters.FormatTemplate(prompts.Prompts, "jobInfoExtraction.txt", promptData)

	if err != nil {
		return nil, fmt.Errorf("failed to format prompt template: %w", err)
	}

	instructions, err := embeds.ReadFile(instructions.Instructions, "jobInfoExtraction.txt")
	if err != nil {
		return nil, err
	}

	rawResponse, err := llmProvider.Generate(
		ctx,
		instructions,
		prompt,
		app_schemas.JobDescriptionResponseFormat,
		apiKey,
	)

	if err != nil {
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}

	cleanedJSON := llm.FormatLLMResponse(rawResponse)

	var llmResponse dto.JobDescription
	if err := json.Unmarshal([]byte(cleanedJSON), &llmResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal LLM response for job info: %w. Raw response: %s", err, rawResponse)
	}

	return &llmResponse, nil
}
