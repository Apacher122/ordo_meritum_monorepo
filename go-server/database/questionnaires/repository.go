package questionnaires

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/ordo_meritum/database/models"
)

type Repository interface {
	CreateOrUpdate(ctx context.Context, questionnaire *models.CandidateQuestionnaire) (*models.CandidateQuestionnaire, error)
	GetByFirebaseUID(ctx context.Context, firebaseUID string) (*models.CandidateQuestionnaire, error)
	Delete(ctx context.Context, firebaseUID string) error
}

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateOrUpdate(ctx context.Context, questionnaire *models.CandidateQuestionnaire) (*models.CandidateQuestionnaire, error) {
	query := `
        INSERT INTO candidate_questionnaires (firebase_uid, title, brief_history, questions)
        VALUES (:firebase_uid, :title, :brief_history, :questions)
        ON CONFLICT (firebase_uid) DO UPDATE SET
            questions = EXCLUDED.questions,
            brief_history = EXCLUDED.brief_history,
            title = EXCLUDED.title,
            updated_at = NOW()
        RETURNING *
    `

	rows, err := r.db.NamedQueryContext(ctx, query, questionnaire)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var createdOrUpdated models.CandidateQuestionnaire
	if rows.Next() {
		if err := rows.StructScan(&createdOrUpdated); err != nil {
			return nil, err
		}
	}
	log.Printf("DATABASE: Upserted questionnaire for user %s", questionnaire.FirebaseUID)
	return &createdOrUpdated, nil
}

func (r *postgresRepository) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*models.CandidateQuestionnaire, error) {
	var questionnaire models.CandidateQuestionnaire
	query := "SELECT * FROM candidate_questionnaires WHERE firebase_uid = $1"

	err := r.db.GetContext(ctx, &questionnaire, query, firebaseUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	log.Printf("DATABASE: Fetched questionnaire for user %s", firebaseUID)
	return &questionnaire, nil
}

func (r *postgresRepository) Delete(ctx context.Context, firebaseUID string) error {
	query := "DELETE FROM candidate_questionnaires WHERE firebase_uid = $1"
	result, err := r.db.ExecContext(ctx, query, firebaseUID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("DATABASE: No questionnaire found to delete for user %s", firebaseUID)
	} else {
		log.Printf("DATABASE: Deleted questionnaire for user %s", firebaseUID)
	}
	return nil
}

var _ Repository = (*postgresRepository)(nil)
