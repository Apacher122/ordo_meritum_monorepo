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
	"github.com/ordo_meritum/shared/utils/logger"
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
) (any, *error_messages.ErrorBody) {
	elg := &logger.ErrorLoggerType{Service: &serviceName}
	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		elg.ErrorCode = &error_messages.ERR_USER_NO_CONTEXT
		elg.ErrorLog()
		return nil, &error_messages.ErrorBody{ErrCode: error_messages.ERR_USER_NO_CONTEXT}
	}
	parsedJob, err := s.parseJobDescriptionWithLLM(
		ctx,
		&requestBody,
	)

	if err != nil {
		elg.Uid = &userCtx.UID
		elg.ErrorCode = &err.ErrCode
		elg.Error = err.ErrMsg
		elg.ErrorLog()
		return nil, err
	}

	logger.InfoLoggerType{Service: &serviceName, Uid: &userCtx.UID, Message: "Persisting full job posting to database..."}.InfoLog()
	cn := formatters.ToSnakeCase(parsedJob.CompanyName)
	res, err := s.jobRepo.InsertFullJobPosting(ctx, requestBody.JobDescription, parsedJob, cn, parsedJob.CompanyName)
	if err != nil {
		elg.ErrorCode = &err.ErrCode
		elg.Error = err.ErrMsg
		elg.ErrorLog()
		return nil, err
	}

	logger.InfoLoggerType{Service: &serviceName, Uid: &userCtx.UID, JobID: &res.RoleID, Message: "Successfully tracked new job."}.InfoLog()
	return res.ID, nil
}

func (s *AppTrackerService) GetTrackedApplicationByID(
	ctx context.Context,
	roleID int,
) (*jobs.FullJobPosting, *error_messages.ErrorBody) {
	return s.jobRepo.GetFullJobPosting(ctx, roleID)
}

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

func (s *AppTrackerService) ListTrackedApplications(
	ctx context.Context,
) ([]*jobs.UserJobPosting, *error_messages.ErrorBody) {
	return s.jobRepo.GetAllUserJobPostings(ctx)
}

func (s *AppTrackerService) parseJobDescriptionWithLLM(
	ctx context.Context,
	r *request.JobPostingRequest,
) (*domain.JobDescription, *error_messages.ErrorBody) {
	llmProvider, err := llm.GetProvider("cohere")
	if err != nil {
		return nil, &error_messages.ErrorBody{}
	}

	jobPost := request.FormatJobPostingRequest(r)
	promptData := map[string]string{
		"JobPost": jobPost,
	}
	prompt, err := formatters.FormatTemplate(prompts.Prompts, "jobInfoExtraction.txt", promptData)

	if err != nil {
		return nil, &error_messages.ErrorBody{ErrCode: error_messages.ERR_LLM_PROMPT_FORMATTING, ErrMsg: fmt.Errorf("failed to format prompt template: %w", err.ErrMsg)}
	}

	instructions, errMsg := embeds.ReadFile(instructions.Instructions, "jobInfoExtraction.txt")
	if err != nil {
		return nil, &error_messages.ErrorBody{ErrCode: error_messages.ERR_LLM_INSTRUCTION_FORMATTING, ErrMsg: errMsg}
	}

	sch, err := schemaregistry.GetSchema("cohere", schemaregistry.ApplicationTracking)
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
		return nil, err
	}

	cleanedJSON := llm.FormatLLMResponse(rawResponse)

	var llmResponse domain.JobDescription
	if err := json.Unmarshal([]byte(cleanedJSON), &llmResponse); err != nil {
		return nil, &error_messages.ErrorBody{ErrCode: error_messages.ERR_LLM_MALFORMED_RESPONSE, ErrMsg: fmt.Errorf("failed to unmarshal LLM response for job info: %w. Raw response: %s", err, rawResponse)}
	}

	return &llmResponse, nil
}
