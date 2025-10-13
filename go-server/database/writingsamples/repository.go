package writingsamples

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/ordo_meritum/database/models"
)

type Repository interface {
	CreateOrUpdate(ctx context.Context, firebaseUID string, samples []models.CandidateWritingSample) error
	GetByFirebaseUID(ctx context.Context, firebaseUID string) ([]models.CandidateWritingSample, error)
}

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateOrUpdate(ctx context.Context, firebaseUID string, samples []models.CandidateWritingSample) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM candidate_writing_samples WHERE firebase_uid = $1", firebaseUID)
	if err != nil {
		return err
	}

	if len(samples) > 0 {
		query := "INSERT INTO candidate_writing_samples (firebase_uid, content) VALUES (:firebase_uid, :content)"
		_, err = tx.NamedExecContext(ctx, query, samples)
		if err != nil {
			return err
		}
	}

	log.Printf("DATABASE: Upserted %d writing samples for user %s", len(samples), firebaseUID)
	return tx.Commit()
}

func (r *postgresRepository) GetByFirebaseUID(ctx context.Context, firebaseUID string) ([]models.CandidateWritingSample, error) {
	var samples []models.CandidateWritingSample
	query := "SELECT * FROM candidate_writing_samples WHERE firebase_uid = $1"
	err := r.db.SelectContext(ctx, &samples, query, firebaseUID)
	if err != nil {
		return nil, err
	}
	log.Printf("DATABASE: Fetched %d writing samples for user %s", len(samples), firebaseUID)
	return samples, nil
}

var _ Repository = (*postgresRepository)(nil)
