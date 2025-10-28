package services

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/ordo_meritum/database/jobs"
	"github.com/ordo_meritum/features/application_tracking/models/domain"
	request "github.com/ordo_meritum/features/application_tracking/models/requests"

	"github.com/ordo_meritum/shared/contexts"
	"github.com/ordo_meritum/shared/embeds"
	"github.com/ordo_meritum/shared/libs/llm"
	schemaregistry "github.com/ordo_meritum/shared/libs/llm/schema_registry"
	"github.com/ordo_meritum/shared/templates/instructions"
	prompts "github.com/ordo_meritum/shared/templates/prompts"
	error_messages "github.com/ordo_meritum/shared/utils/errors"
	formatters "github.com/ordo_meritum/shared/utils/formatters"
	lg "github.com/ordo_meritum/shared/utils/logger"
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

// QueueApplicationTracking starts the application tracking process, given a job posting request.
// It will return the ID of the job posting in the database and an error if any errors occur.
// If the user context is not found, it will return an error with ErrorCode set to ERR_USER_NO_CONTEXT.
// If the job description cannot be parsed with LLM, it will return an error with ErrorCode set to ERR_LLM_NO_CONTENT.
// If the database fails to insert the full job posting, it will return an error with ErrorCode set to ERR_DB_FAILED_TO_INSERT.
func (s *AppTrackerService) QueueApplicationTracking(
	ctx context.Context,
	requestBody *request.JobPostingRequest,
) (int, error) {
	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		return 0, error_messages.ErrorMessage(error_messages.ERR_USER_NO_CONTEXT)
	}
	l := log.With().
		Str("service", serviceName).
		Str("uid", userCtx.UID).
		Logger()

	l.Info().Msg("Starting application tracking process")

	parsedJob, err := s.parseJobDescriptionWithLLM(
		ctx,
		requestBody,
	)

	if err != nil {
		lg.ErrorLoggerType{Service: &serviceName, ErrorCode: &error_messages.ERR_LLM_NO_CONTENT, Error: err}.ErrorLog()
		return 0, err
	}

	l.Info().Msg("Persisting full job posting to database...")
	cn := formatters.ToSnakeCase(parsedJob.CompanyName)
	res, err := s.jobRepo.InsertFullJobPosting(ctx, requestBody.JobDescription, parsedJob, cn, parsedJob.CompanyName)
	if err != nil {
		lg.ErrorLoggerType{Service: &serviceName, ErrorCode: &error_messages.ERR_DB_FAILED_TO_INSERT, Error: err}.ErrorLog()
		return 0, err
	}

	l.Info().Msg(fmt.Sprintf("Successfully created job posting with ID: %d and Company Name: %s", res.ID, parsedJob.CompanyName))
	return res.ID, nil
}

// GetTrackedApplicationByID retrieves a job posting by its role ID.
// If the role ID cannot be found, it will return an error with ErrorCode set to ERR_DB_FAILED_TO_GET.
// If the user context is not found, it will return an error with ErrorCode set to ERR_USER_NO_CONTEXT.
func (s *AppTrackerService) GetTrackedApplicationByID(
	ctx context.Context,
	roleID int,
) (*jobs.FullJobPosting, error) {
	return s.jobRepo.GetFullJobPosting(ctx, roleID)
}

// UpdateApplicationStatus updates the application status for a given job posting.
// It takes in a request which contains the role ID and the application status.
// It will return an error if any errors occur.
// If the role ID cannot be found, it will return an error with ErrorCode set to ERR_DB_FAILED_TO_GET.
// If the user context is not found, it will return an error with ErrorCode set to ERR_USER_NO_CONTEXT.
func (s *AppTrackerService) UpdateApplicationStatus(
	ctx context.Context,
	request *request.ApplicationUpdateRequest,
) error {
	l := log.With().
		Str("service", serviceName).
		Logger()
	err := s.jobRepo.UpdateApplicationDetails(ctx, request.Payload.JobID, request)
	if err != nil {
		l.Error().Err(err).Msg("Error updating application status")
		return err
	}
	return nil
}

// ListTrackedApplications retrieves all job postings for a given user.
// It will return an error if any errors occur.
// If the user context is not found, it will return an error with ErrorCode set to ERR_USER_NO_CONTEXT.
func (s *AppTrackerService) ListTrackedApplications(
	ctx context.Context,
) ([]*jobs.UserJobPosting, error) {
	return s.jobRepo.GetAllUserJobPostings(ctx)
}

// parseJobDescriptionWithLLM generates a job description from a job posting request using an LLM.
// It takes in a job posting request and returns a job description if successful.
// If the LLM provider cannot be found, it will return an error with ErrorCode set to ERR_LLM_NO_CONTENT.
// If the prompt cannot be formatted, it will return an error with ErrorCode set to ERR_LLM_PROMPT_FORMATTING.
// If the instructions file cannot be read, it will return an error with ErrorCode set to ERR_LLM_INSTRUCTION_FORMATTING.
// If the LLM response cannot be unmarshalled, it will return an error with ErrorCode set to ERR_LLM_MALFORMED_RESPONSE.
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

	sch, err := schemaregistry.GetSchema("cohere", schemaregistry.ApplicationTracking, nil)
	if err != nil {
		return nil, err
	}

	rawResponse, err := llmProvider.Generate(
		ctx,
		instructions,
		prompt,
		sch,
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
