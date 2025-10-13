package guides

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/ordo_meritum/database/models"
	dto "github.com/ordo_meritum/features/job_guide/models/responses"
)

type Repository interface {
	InsertMatchSummary(ctx context.Context, firebaseUID string, roleID int, summary *dto.MatchSummaryPayload) error
}

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) getResumeID(ctx context.Context, tx *sqlx.Tx, firebaseUID string, roleID int) (int, error) {
	var resumeID int
	query := "SELECT id FROM resumes WHERE firebase_uid = $1 AND role_id = $2"
	err := tx.GetContext(ctx, &resumeID, query, firebaseUID, roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no resume found for user %s and role %d", firebaseUID, roleID)
		}
		return 0, err
	}
	return resumeID, nil
}

func (r *postgresRepository) InsertMatchSummary(ctx context.Context, firebaseUID string, roleID int, summaryPayload *dto.MatchSummaryPayload) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	resumeID, err := r.getResumeID(ctx, tx, firebaseUID, roleID)
	if err != nil {
		return err
	}

	var existingMatchIDs []int
	err = tx.SelectContext(ctx, &existingMatchIDs, "SELECT id FROM match_summaries WHERE resume_id = $1", resumeID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("could not query for existing match summaries: %w", err)
	}

	if len(existingMatchIDs) > 0 {
		query, args, err := sqlx.In("DELETE FROM match_summary_overviews WHERE match_summary_id IN (?)", existingMatchIDs)
		if err != nil {
			return fmt.Errorf("could not build query for deleting overviews: %w", err)
		}
		query = tx.Rebind(query)
		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to delete existing match_summary_overviews: %w", err)
		}

		query, args, err = sqlx.In("DELETE FROM metrics WHERE match_summary_id IN (?)", existingMatchIDs)
		if err != nil {
			return fmt.Errorf("could not build query for deleting metrics: %w", err)
		}
		query = tx.Rebind(query)
		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to delete existing metrics: %w", err)
		}

		query, args, err = sqlx.In("DELETE FROM match_summaries WHERE id IN (?)", existingMatchIDs)
		if err != nil {
			return fmt.Errorf("could not build query for deleting summaries: %w", err)
		}
		query = tx.Rebind(query)
		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to delete existing match_summaries: %w", err)
		}
	}

	overallSummary := summaryPayload.MatchSummary.OverallMatchSummary
	var matchID int
	summaryQuery := `
        INSERT INTO match_summaries (resume_id, should_apply, reasoning, overall_match_score, suggestions)
        VALUES ($1, $2, $3, $4, $5) RETURNING id
    `
	err = tx.GetContext(ctx, &matchID, summaryQuery,
		resumeID,
		summaryPayload.MatchSummary.ShouldApply,
		summaryPayload.MatchSummary.ShouldApplyReasoning,
		overallSummary.OverallMatchScore,
		pq.Array(overallSummary.Suggestions),
	)
	if err != nil {
		return fmt.Errorf("failed to insert into match_summaries: %w", err)
	}

	overviewQuery := `INSERT INTO match_summary_overviews (match_summary_id, summary, summary_temperature) VALUES (:match_summary_id, :summary, :summary_temperature)`
	var overviewsToInsert []models.MatchSummaryOverview
	for _, item := range overallSummary.Summary {
		overviewsToInsert = append(overviewsToInsert, models.MatchSummaryOverview{
			MatchSummaryID:     matchID,
			Summary:            item.SummaryText,
			SummaryTemperature: item.SummaryTemperature,
		})
	}
	if len(overviewsToInsert) > 0 {
		_, err = tx.NamedExecContext(ctx, overviewQuery, overviewsToInsert)
		if err != nil {
			return fmt.Errorf("failed to batch insert match_summary_overviews: %w", err)
		}
	}

	metricsQuery := `INSERT INTO metrics (match_summary_id, score_title, raw_score, weighted_score, score_weight, score_reason, is_compatible, strength, weaknesses) VALUES (:match_summary_id, :score_title, :raw_score, :weighted_score, :score_weight, :score_reason, :is_compatible, :strength, :weaknesses)`
	var metricsToInsert []models.Metrics
	for _, metric := range summaryPayload.MatchSummary.Metrics {
		metricsToInsert = append(metricsToInsert, models.Metrics{
			MatchSummaryID: matchID,
			ScoreTitle:     metric.ScoreTitle,
			RawScore:       metric.RawScore,
			WeightedScore:  metric.WeightedScore,
			ScoreWeight:    metric.ScoreWeight,
			ScoreReason:    &metric.ScoreReason,
			IsCompatible:   &metric.IsCompatible,
			Strength:       &metric.Strength,
			Weaknesses:     &metric.Weaknesses,
		})
	}
	if len(metricsToInsert) > 0 {
		_, err = tx.NamedExecContext(ctx, metricsQuery, metricsToInsert)
		if err != nil {
			return fmt.Errorf("failed to batch insert metrics: %w", err)
		}
	}

	return tx.Commit()
}

var _ Repository = (*postgresRepository)(nil)
