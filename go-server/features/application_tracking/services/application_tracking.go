package services

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/ordo_meritum/database/jobs"
	db_models "github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/features/application_tracking/models/domain"
	request "github.com/ordo_meritum/features/application_tracking/models/requests"

	app_schemas "github.com/ordo_meritum/features/application_tracking/models/schemas"
	"github.com/ordo_meritum/shared/contexts"
	"github.com/ordo_meritum/shared/embeds"
	"github.com/ordo_meritum/shared/libs/llm"
	"github.com/ordo_meritum/shared/templates/instructions"
	prompts "github.com/ordo_meritum/shared/templates/prompts"
	error_messages "github.com/ordo_meritum/shared/utils/errors"
	formatters "github.com/ordo_meritum/shared/utils/formatters"
)

var serviceName = "application-tracking"

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
	requestBody request.JobPostingRequest,
) (any, error) {
	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		return nil, error_messages.ErrorMessage(error_messages.ERR_USER_NO_CONTEXT)
	}
	l := log.With().
		Str("service", serviceName).
		Str("uid", userCtx.UID).
		Logger()

	l.Info().Msg("Starting application tracking process")

	parsedJob, err := s.parseJobDescriptionWithLLM(
		ctx,
		&requestBody,
	)

	if err != nil {
		error_messages.ErrorLog(error_messages.ERR_LLM_NO_CONTENT, l.Error()).Msg("could not extract job info from LLM")
		return nil, err
	}

	l.Info().Msg("Persisting full job posting to database...")
	cn := formatters.ToSnakeCase(parsedJob.CompanyName)
	res, err := s.jobRepo.InsertFullJobPosting(ctx, requestBody.JobDescription, parsedJob, cn, parsedJob.CompanyName)
	if err != nil {
		error_messages.ErrorLog(error_messages.ERR_DB_FAILED_TO_INSERT, l.Error()).Msg("failed to insert full job posting")
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
	roleID int,
	status db_models.AppStatus,
) error {
	return s.jobRepo.UpdateApplicationDetails(ctx, roleID, &status, nil)
}

func (s *AppTrackerService) ListTrackedApplications(
	ctx context.Context,
) ([]*jobs.UserJobPosting, error) {
	return s.jobRepo.GetAllUserJobPostings(ctx)
}

func (s *AppTrackerService) parseJobDescriptionWithLLM(
	ctx context.Context,
	r *request.JobPostingRequest,
) (*domain.JobDescription, error) {
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
	)

	if err != nil {
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}

	cleanedJSON := llm.FormatLLMResponse(rawResponse)

	var llmResponse domain.JobDescription
	if err := json.Unmarshal([]byte(cleanedJSON), &llmResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal LLM response for job info: %w. Raw response: %s", err, rawResponse)
	}

	return &llmResponse, nil
}
