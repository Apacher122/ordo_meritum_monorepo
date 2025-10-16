package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ordo_meritum/database/guides"
	"github.com/ordo_meritum/database/jobs"
	"github.com/ordo_meritum/database/resumes"
	job_mappers "github.com/ordo_meritum/features/application_tracking/utils/mappers"
	"github.com/ordo_meritum/features/job_guide/models/domain"
	"github.com/ordo_meritum/features/job_guide/models/requests"
	"github.com/ordo_meritum/features/job_guide/models/schemas"
	"github.com/ordo_meritum/shared/contexts"
	"github.com/ordo_meritum/shared/libs/llm"
	"github.com/ordo_meritum/shared/templates/instructions"
	"github.com/ordo_meritum/shared/templates/prompts"
	error_messages "github.com/ordo_meritum/shared/utils/errors"
	shared_formatters "github.com/ordo_meritum/shared/utils/formatters"
)

type JobGuideService struct {
	guideRepo  guides.Repository
	resumeRepo resumes.Repository
	jobsRepo   jobs.Repository
}

func NewJobGuideService(guideRepo guides.Repository, resumeRepo resumes.Repository, jobsRepo jobs.Repository) *JobGuideService {
	return &JobGuideService{
		guideRepo:  guideRepo,
		resumeRepo: resumeRepo,
		jobsRepo:   jobsRepo,
	}
}

func (s *JobGuideService) GetCompanyInfo(ctx context.Context, companyName string, jobID int) (*domain.CompanyInfo, error) {
	// log.Printf("Fetching company info for: %s", companyName)
	// prompt := fmt.Sprintf(constants.CompanyInfoPrompt, companyName, jobDescription)
	// rawResponse, err := s.llmProvider.Generate(ctx, prompt)
	// if err != nil {
	// 	return nil, fmt.Errorf("LLM generation for company info failed: %w", err)
	// }
	// cleanedJSON := formatLLMResponse(rawResponse)
	// var companyInfo domain.CompanyInfo
	// if err := json.Unmarshal([]byte(cleanedJSON), &companyInfo); err != nil {
	// 	return nil, fmt.Errorf("failed to unmarshal company info: %w", err)
	// }
	// return &companyInfo, nil
	return nil, nil
}

func (s *JobGuideService) GetMatchSummary(ctx context.Context, r *requests.JobGuideRequests) error {

	promptData, err := s.buildMatchSummaryPromptData(ctx, &r.Payload)
	if err != nil {
		return err
	}

	var matchSummary domain.MatchSummary
	err = s.generateLLMContent(
		ctx,
		r.Options.LLMProvider,
		"matchsummary.txt",
		promptData,
		schemas.GeminiResumeSchema,
		matchSummary,
	)

	if err != nil {
		return err
	}

	// Instructions & prompts

	// Call llm

	// Insert and Send data

	// log.Printf("Generating match summary for user %s, role %d", firebaseUID, roleID)
	// resumeJSON, _ := json.Marshal(resume)
	// jobDescJSON, _ := json.Marshal(jobDesc)
	// prompt := fmt.Sprintf(constants.MatchSummaryPrompt, string(resumeJSON), string(jobDescJSON))

	// rawResponse, err := s.llmProvider.Generate(ctx, prompt)
	// if err != nil {
	// 	return fmt.Errorf("LLM generation for match summary failed: %w", err)
	// }

	// cleanedJSON := formatLLMResponse(rawResponse)

	// // Unmarshal the complex, nested LLM response into the domain.
	// var summaryPayload domain.MatchSummaryPayload
	// if err := json.Unmarshal([]byte(cleanedJSON), &summaryPayload); err != nil {
	// 	return fmt.Errorf("failed to unmarshal match summary payload from LLM: %w", err)
	// }

	// // Pass the entire domain to the repository, which handles the complex, multi-table insertion.
	// err = s.guideRepo.InsertMatchSummary(ctx, firebaseUID, roleID, &summaryPayload)
	// if err != nil {
	// 	return fmt.Errorf("failed to save match summary to database: %w", err)
	// }

	return nil
}

func (s *JobGuideService) GetGuidingAnswers(ctx context.Context, uid string, jobID int) (*domain.GuidingQuestions, error) {
	// log.Printf("Generating guiding answers for %s at %s", jobDesc.Role, jobDesc.Company)
	// resumeJSON, _ := json.Marshal(resume)
	// jobDescJSON, _ := json.Marshal(jobDesc)
	// prompt := fmt.Sprintf(constants.GuidingQuestionsPrompt, string(resumeJSON), string(jobDescJSON))

	// rawResponse, err := s.llmProvider.Generate(ctx, prompt)
	// if err != nil {
	// 	return nil, fmt.Errorf("LLM generation for guiding answers failed: %w", err)
	// }
	// cleanedJSON := formatLLMResponse(rawResponse)
	// var questions domain.GuidingQuestions
	// if err := json.Unmarshal([]byte(cleanedJSON), &questions); err != nil {
	// 	return nil, fmt.Errorf("failed to unmarshal guiding questions: %w", err)
	// }
	// return &questions, nil
	return nil, nil
}

func (s *JobGuideService) generateLLMContent(
	ctx context.Context,
	providerName, instructionsFile string,
	promptData any,
	schema any,
	target interface{},
) error {
	llmProvider, err := llm.GetProvider(providerName)
	if err != nil {
		return err
	}

	prompt, err := shared_formatters.FormatTemplate(prompts.Prompts, instructionsFile, promptData)
	if err != nil {
		return fmt.Errorf("failed to format prompt template: %w", err)
	}

	instructionBytes, err := instructions.Instructions.ReadFile(instructionsFile)
	if err != nil {
		return fmt.Errorf("failed to read instructions file: %w", err)
	}

	rawResponse, err := llmProvider.Generate(ctx, string(instructionBytes), prompt, schema)
	if err != nil {
		return fmt.Errorf("LLM generation failed: %w", err)
	}

	cleanedJSON := llm.FormatLLMResponse(rawResponse)
	if err := json.Unmarshal([]byte(cleanedJSON), target); err != nil {
		return fmt.Errorf("failed to unmarshal LLM response: %w. Raw response: %s", err, rawResponse)
	}

	return nil
}

func (s *JobGuideService) buildMatchSummaryPromptData(ctx context.Context, payload *requests.JobGuidePayload) (map[string]any, error) {
	_, ok := contexts.FromContext(ctx)
	if !ok {
		return nil, error_messages.ErrorMessage(error_messages.ERR_USER_NO_CONTEXT)
	}

	r, err := s.resumeRepo.GetFullResume(ctx, payload.JobID)
	if err != nil {
		return nil, err
	}
	j, err := s.jobsRepo.GetFullJobPosting(ctx, payload.JobID)
	if err != nil {
		return nil, err
	}

	resume := r.FormatForLLM()
	jobPost := job_mappers.NewJobDescriptionFromPost(j)

	return map[string]any{
		"JobPost":     jobPost.FormatForLLM(),
		"Applicants":  jobPost.ApplicantCount,
		"Education":   payload.EducationInfo.FormatForLLM(),
		"Resume":      resume,
		"CoverLetter": payload.CoverLetter,
	}, nil
}
