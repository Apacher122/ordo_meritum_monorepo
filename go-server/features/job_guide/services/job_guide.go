package services

import (
	"context"

	"github.com/ordo_meritum/database/guides"
	"github.com/ordo_meritum/features/job_guide/models/dto"
)

type JobGuideService struct {
	guideRepo guides.Repository
}

func NewJobGuideService(guideRepo guides.Repository) *JobGuideService {
	return &JobGuideService{
		guideRepo: guideRepo,
	}
}

func (s *JobGuideService) GetCompanyInfo(ctx context.Context, companyName string, jobID int) (*dto.CompanyInfo, error) {
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

func (s *JobGuideService) GetMatchSummary(ctx context.Context, firebaseUID string, roleID int) error {
	// log.Printf("Generating match summary for user %s, role %d", firebaseUID, roleID)
	// resumeJSON, _ := json.Marshal(resume)
	// jobDescJSON, _ := json.Marshal(jobDesc)
	// prompt := fmt.Sprintf(constants.MatchSummaryPrompt, string(resumeJSON), string(jobDescJSON))

	// rawResponse, err := s.llmProvider.Generate(ctx, prompt)
	// if err != nil {
	// 	return fmt.Errorf("LLM generation for match summary failed: %w", err)
	// }

	// cleanedJSON := formatLLMResponse(rawResponse)

	// // Unmarshal the complex, nested LLM response into the DTO.
	// var summaryPayload dto.MatchSummaryPayload
	// if err := json.Unmarshal([]byte(cleanedJSON), &summaryPayload); err != nil {
	// 	return fmt.Errorf("failed to unmarshal match summary payload from LLM: %w", err)
	// }

	// // Pass the entire DTO to the repository, which handles the complex, multi-table insertion.
	// err = s.guideRepo.InsertMatchSummary(ctx, firebaseUID, roleID, &summaryPayload)
	// if err != nil {
	// 	return fmt.Errorf("failed to save match summary to database: %w", err)
	// }

	return nil
}

func (s *JobGuideService) GetGuidingAnswers(ctx context.Context, uid string, jobID int) (*dto.GuidingQuestions, error) {
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
