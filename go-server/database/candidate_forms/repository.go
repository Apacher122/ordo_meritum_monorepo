package candidate_forms

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	db_models "github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/features/candidate_forms/models/requests"
)

type Repository interface {
	UpsertPersonalityProfilee(ctx context.Context, firebaseUID string, oceanData db_models.OceanProfile, discData db_models.DiscProfile) error
	GetPersonalityProfile(ctx context.Context, firebaseUID string) (*db_models.OceanProfile, *db_models.DiscProfile, error)
	UpsertQuestionnaire(ctx context.Context, firebaseUID string, req requests.QuestionnaireRequest) (*db_models.CandidateQuestionnaire, error)
	GetQuestionnaireByFirebaseUID(ctx context.Context, firebaseUID string) (*db_models.CandidateQuestionnaire, error)
}

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) Repository {
	return &postgresRepository{db: db}
}

var ErrUserNotFound = errors.New("user not found")

func (r *postgresRepository) UpsertPersonalityProfilee(
	ctx context.Context,
	firebaseUID string,
	oceanData db_models.OceanProfile,
	discData db_models.DiscProfile,
) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	var profile db_models.PersonalityProfile
	err = tx.GetContext(ctx, &profile, "SELECT * FROM personality_profiles WHERE firebase_uid = $1", firebaseUID)

	if errors.Is(err, sql.ErrNoRows) {
		err = tx.GetContext(ctx, &profile, "INSERT INTO personality_profiles (firebase_uid) VALUES ($1) RETURNING *", firebaseUID)
		if err != nil {
			return fmt.Errorf("could not insert new personality profile: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO ocean_profiles (
				personality_profiles_id, openness_to_experience_score, openness_to_experience_reasoning,
				conscientiousness_score, conscientiousness_reasoning, extraversion_score, extraversion_reasoning,
				agreeableness_score, agreeableness_reasoning, neuroticism_score, neuroticism_reasoning, summary
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			profile.ID, oceanData.OpennessScore, oceanData.OpennessReasoning,
			oceanData.ConscientiousnessScore, oceanData.ConscientiousnessReasoning, oceanData.ExtraversionScore, oceanData.ExtraversionReasoning,
			oceanData.AgreeablenessScore, oceanData.AgreeablenessReasoning, oceanData.NeuroticismScore, oceanData.NeuroticismReasoning, oceanData.Summary)
		if err != nil {
			return fmt.Errorf("could not insert new ocean profile: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO disc_profiles (
				personality_profiles_id, dominance, influence, steadiness, conscientiousness, summary
			) VALUES ($1, $2, $3, $4, $5, $6)`,
			profile.ID, discData.Dominance, discData.Influence, discData.Steadiness, discData.Conscientiousness, discData.Summary)
		if err != nil {
			return fmt.Errorf("could not insert new disc profile: %w", err)
		}

	} else if err == nil {
		_, err = tx.ExecContext(ctx, `
			UPDATE ocean_profiles SET
				openness_to_experience_score = $2, openness_to_experience_reasoning = $3,
				conscientiousness_score = $4, conscientiousness_reasoning = $5, extraversion_score = $6,
				extraversion_reasoning = $7, agreeableness_score = $8, agreeableness_reasoning = $9,
				neuroticism_score = $10, neuroticism_reasoning = $11, summary = $12
			WHERE personality_profiles_id = $1`,
			profile.ID, oceanData.OpennessScore, oceanData.OpennessReasoning,
			oceanData.ConscientiousnessScore, oceanData.ConscientiousnessReasoning, oceanData.ExtraversionScore,
			oceanData.ExtraversionReasoning, oceanData.AgreeablenessScore, oceanData.AgreeablenessReasoning,
			oceanData.NeuroticismScore, oceanData.NeuroticismReasoning, oceanData.Summary)
		if err != nil {
			return fmt.Errorf("could not update ocean profile: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE disc_profiles SET
				dominance = $2, influence = $3, steadiness = $4, conscientiousness = $5, summary = $6
			WHERE personality_profiles_id = $1`,
			profile.ID, discData.Dominance, discData.Influence, discData.Steadiness, discData.Conscientiousness, discData.Summary)
		if err != nil {
			return fmt.Errorf("could not update disc profile: %w", err)
		}

	} else {
		return fmt.Errorf("could not check for existing profile: %w", err)
	}

	return tx.Commit()
}

func (r *postgresRepository) GetPersonalityProfile(
	ctx context.Context,
	firebaseUID string,
) (*db_models.OceanProfile, *db_models.DiscProfile, error) {
	type fullProfile struct {
		db_models.OceanProfile
		db_models.DiscProfile
	}

	query := `
		SELECT
			o.*, d.*
		FROM personality_profiles p
		JOIN ocean_profiles o ON p.id = o.personality_profiles_id
		JOIN disc_profiles  d ON p.id = d.personality_profiles_id
		WHERE p.firebase_uid = $1`

	var profile fullProfile
	err := r.db.GetContext(ctx, &profile, query, firebaseUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, ErrUserNotFound
		}
		return nil, nil, fmt.Errorf("could not get personality profile: %w", err)
	}

	return &profile.OceanProfile, &profile.DiscProfile, nil
}

func (r *postgresRepository) UpsertQuestionnaire(
	ctx context.Context,
	firebaseUID string,
	req requests.QuestionnaireRequest,
) (*db_models.CandidateQuestionnaire, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	var questionnaireID int

	upsertQuestionnaireQuery := `
		INSERT INTO candidate_questionnaires (firebase_uid, brief_history, title)
		VALUES ($1, $2, $3)
		ON CONFLICT (firebase_uid) DO UPDATE
		SET brief_history = EXCLUDED.brief_history, title = EXCLUDED.title, updated_at = NOW()
		RETURNING id`

	placeholderTitle := "Candidate Questionnaire"
	err = tx.QueryRowxContext(ctx, upsertQuestionnaireQuery, firebaseUID, req.BriefHistory, placeholderTitle).Scan(&questionnaireID)
	if err != nil {
		return nil, fmt.Errorf("could not upsert questionnaire: %w", err)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM questions WHERE questionnaire_id = $1", questionnaireID)
	if err != nil {
		return nil, fmt.Errorf("could not delete old questions: %w", err)
	}

	insertQuestionQuery := `
		INSERT INTO questions (questionnaire_id, category, question, answer)
		VALUES ($1, $2, $3, $4)`

	for _, category := range req.QuestionsByCategory {
		for _, qa := range category.Questions {
			_, err = tx.ExecContext(ctx, insertQuestionQuery, questionnaireID, category.Category, qa.Question, qa.Answer)
			if err != nil {
				return nil, fmt.Errorf("could not insert question: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	return r.getCompleteQuestionnaire(ctx, questionnaireID)
}

func (r *postgresRepository) getCompleteQuestionnaire(
	ctx context.Context,
	questionnaireID int,
) (*db_models.CandidateQuestionnaire, error) {
	var questionnaire db_models.CandidateQuestionnaire
	err := r.db.GetContext(ctx, &questionnaire, "SELECT * FROM candidate_questionnaires WHERE id = $1", questionnaireID)
	if err != nil {
		return nil, err
	}

	err = r.db.SelectContext(ctx, &questionnaire.Questions, "SELECT * FROM questions WHERE questionnaire_id = $1", questionnaireID)
	if err != nil {
		return nil, err
	}

	return &questionnaire, nil
}

func (r *postgresRepository) GetQuestionnaireByFirebaseUID(
	ctx context.Context,
	firebaseUID string,
) (*db_models.CandidateQuestionnaire, error) {
	var questionnaire db_models.CandidateQuestionnaire

	getQuestionnaireQuery := "SELECT * FROM candidate_questionnaires WHERE firebase_uid = $1"
	err := r.db.GetContext(ctx, &questionnaire, getQuestionnaireQuery, firebaseUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("could not get questionnaire: %w", err)
	}

	getQuestionsQuery := "SELECT * FROM questions WHERE questionnaire_id = $1"
	err = r.db.SelectContext(ctx, &questionnaire.Questions, getQuestionsQuery, questionnaire.ID)
	if err != nil {
		return nil, fmt.Errorf("could not get questions for questionnaire: %w", err)
	}

	return &questionnaire, nil
}

var _ Repository = (*postgresRepository)(nil)
