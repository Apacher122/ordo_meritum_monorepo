package services

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/ordo_meritum/database/candidate_forms"
	"github.com/ordo_meritum/features/candidate_forms/models/dto"
	"github.com/ordo_meritum/shared/mappers"
	"github.com/ordo_meritum/shared/models/requests"
)

type CandidateFormsService struct {
	candidateFormRepo candidate_forms.Repository
}

func NewCandidateFormService(
	candidateFormRepo candidate_forms.Repository,
) *CandidateFormsService {
	return &CandidateFormsService{
		candidateFormRepo: candidateFormRepo,
	}
}

func (s *CandidateFormsService) SaveCandidateQuestionnaire(
	ctx context.Context,
	firebaseUID string,
	apiKey string,
	requestBody *requests.RequestBody,
) error {
	l := log.With().
		Str("service", "candidate-forms").
		Str("uid", firebaseUID).
		Logger()

	form := requestBody.Payload.Questionnaire
	if form.QuestionsByCategory == nil {
		err := fmt.Errorf("questionsByCategory is nil")
		l.Error().Err(err).Msg("questionsByCategory is nil")
		return err
	}

	_, err := s.candidateFormRepo.UpsertQuestionnaire(ctx, firebaseUID, requestBody.Payload.Questionnaire)
	if err != nil {
		l.Error().Err(err).Msg("Error saving questionnaire")
		return err
	}

	return nil
}

func (s *CandidateFormsService) SavePersonalityProfile(
	ctx context.Context,
	firebaseUID string,
	summary dto.PersonalitySummary,
) (*dto.PersonalitySummary, error) {
	dbOcean, dbDisc := mappers.MapDTOToDB(summary)

	err := s.candidateFormRepo.UpsertPersonalityProfilee(ctx, firebaseUID, dbOcean, dbDisc)
	if err != nil {
		log.Error().Err(err).Msg("Error saving profile")
		return nil, err
	}

	return &summary, nil
}

func (s *CandidateFormsService) GetPersonalityProfile(ctx context.Context, firebaseUID string) (*dto.PersonalitySummary, error) {
	dbOcean, dbDisc, err := s.candidateFormRepo.GetPersonalityProfile(ctx, firebaseUID)
	if err != nil {
		log.Error().Err(err).Msg("Error getting profile")
		return nil, fmt.Errorf("service error getting profile: %w", err)
	}

	summaryDTO := mappers.MapDBToDTO(*dbOcean, *dbDisc)

	return &summaryDTO, nil
}

func (s *CandidateFormsService) CreatePersonalityProfile(ctx context.Context, firebaseUID string) error {
	return nil
}
