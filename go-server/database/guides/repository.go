package guides

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/features/job_guide/models/domain"
	"github.com/ordo_meritum/shared/contexts"
	error_response "github.com/ordo_meritum/shared/types/errors"
)

type Repository interface {
	UpsertMatchSummary(ctx context.Context, summaryPayload *domain.MatchSummary) error
	GetMatchSummary(ctx context.Context, roleID int) (*domain.MatchSummary, error)
}

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) getResumeID(ctx context.Context, tx *sqlx.Tx, roleID int) (int, error) {
	var resumeID int

	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		return 0, error_response.ErrNoUserContext
	}

	query := "SELECT id FROM resumes WHERE firebase_uid = $1 AND role_id = $2"
	err := tx.GetContext(ctx, &resumeID, query, userCtx.UID, roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no resume found for user %s and role %d", userCtx.UID, roleID)
		}
		return 0, err
	}
	return resumeID, nil
}

func (r *postgresRepository) UpsertMatchSummary(ctx context.Context, summaryPayload *domain.MatchSummary) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		return error_response.ErrNoUserContext
	}

	var resumeID int
	err = tx.GetContext(ctx, &resumeID, "SELECT id FROM resumes WHERE user_id = $1", userCtx.UID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no resume found for user ID %s", userCtx.UID)
		}
		return fmt.Errorf("could not get resume ID: %w", err)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM match_summaries WHERE resume_id = $1", resumeID)
	if err != nil {
		return fmt.Errorf("failed to delete existing match_summaries: %w", err)
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

	if len(overallSummary.Summary) > 0 {
		overviewQuery := `INSERT INTO match_summary_overviews (match_summary_id, summary, summary_temperature) VALUES (:match_summary_id, :summary, :summary_temperature)`
		var overviewsToInsert []models.MatchSummaryOverview
		for _, item := range overallSummary.Summary {
			overviewsToInsert = append(overviewsToInsert, models.MatchSummaryOverview{
				MatchSummaryID:     matchID,
				Summary:            item.SummaryText,
				SummaryTemperature: item.SummaryTemperature,
			})
		}
		_, err = tx.NamedExecContext(ctx, overviewQuery, overviewsToInsert)
		if err != nil {
			return fmt.Errorf("failed to batch insert match_summary_overviews: %w", err)
		}
	}

	if len(summaryPayload.MatchSummary.Metrics) > 0 {
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
		_, err = tx.NamedExecContext(ctx, metricsQuery, metricsToInsert)
		if err != nil {
			return fmt.Errorf("failed to batch insert metrics: %w", err)
		}
	}

	return tx.Commit()
}

func (r *postgresRepository) GetMatchSummary(ctx context.Context, roleID int) (*domain.MatchSummary, error) {
	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		return nil, error_response.ErrNoUserContext
	}

	var summaryDB models.MatchSummary
	querySummary := `
        SELECT ms.* FROM match_summaries ms
        JOIN resumes r ON ms.resume_id = r.id
        WHERE r.user_id = $1
        ORDER BY ms.created_at DESC
        LIMIT 1;`
	err := r.db.Get(&summaryDB, querySummary, userCtx.UID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no match summary found for user ID %s", userCtx.UID)
		}
		return nil, fmt.Errorf("error fetching match summary: %w", err)
	}

	var overviewsDB []models.MatchSummaryOverview
	queryOverviews := `SELECT * FROM match_summary_overviews WHERE match_summary_id = $1;`
	err = r.db.Select(&overviewsDB, queryOverviews, summaryDB.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching overviews: %w", err)
	}

	var metricsDB []models.Metrics
	queryMetrics := `SELECT * FROM metrics WHERE match_summary_id = $1;`
	err = r.db.Select(&metricsDB, queryMetrics, summaryDB.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching metrics: %w", err)
	}

	return transformToDomain(summaryDB, overviewsDB, metricsDB), nil
}

func transformToDomain(summary models.MatchSummary, overviews []models.MatchSummaryOverview, metrics []models.Metrics) *domain.MatchSummary {
	domainSummary := &domain.MatchSummary{}

	domainSummary.MatchSummary.ShouldApply = models.ShouldApply(summary.ShouldApply)
	domainSummary.MatchSummary.ShouldApplyReasoning = summary.Reasoning

	if summary.OverallMatchScore != nil {
		domainSummary.MatchSummary.OverallMatchSummary.OverallMatchScore = *summary.OverallMatchScore
	}
	domainSummary.MatchSummary.OverallMatchSummary.Suggestions = summary.Suggestions

	domainSummary.MatchSummary.OverallMatchSummary.Summary = make([]struct {
		SummaryText        string             `json:"summary_text"`
		SummaryTemperature models.Temperature `json:"summary_temperature"`
	}, len(overviews))

	for i, o := range overviews {
		domainSummary.MatchSummary.OverallMatchSummary.Summary[i] = struct {
			SummaryText        string             `json:"summary_text"`
			SummaryTemperature models.Temperature `json:"summary_temperature"`
		}{
			SummaryText:        o.Summary,
			SummaryTemperature: models.Temperature(o.SummaryTemperature),
		}
	}

	domainSummary.MatchSummary.Metrics = make([]struct {
		ScoreTitle    string  `json:"score_title"`
		RawScore      float64 `json:"raw_score"`
		WeightedScore float64 `json:"weighted_score"`
		ScoreWeight   float64 `json:"score_weight"`
		ScoreReason   string  `json:"score_reason"`
		IsCompatible  bool    `json:"isCompatible"`
		Strength      string  `json:"strength"`
		Weaknesses    string  `json:"weaknesses"`
	}, len(metrics))

	for i, m := range metrics {
		metric := struct {
			ScoreTitle    string  `json:"score_title"`
			RawScore      float64 `json:"raw_score"`
			WeightedScore float64 `json:"weighted_score"`
			ScoreWeight   float64 `json:"score_weight"`
			ScoreReason   string  `json:"score_reason"`
			IsCompatible  bool    `json:"isCompatible"`
			Strength      string  `json:"strength"`
			Weaknesses    string  `json:"weaknesses"`
		}{
			ScoreTitle:    m.ScoreTitle,
			RawScore:      m.RawScore,
			WeightedScore: m.WeightedScore,
			ScoreWeight:   m.ScoreWeight,
		}
		if m.ScoreReason != nil {
			metric.ScoreReason = *m.ScoreReason
		}
		if m.IsCompatible != nil {
			metric.IsCompatible = *m.IsCompatible
		}
		if m.Strength != nil {
			metric.Strength = *m.Strength
		}
		if m.Weaknesses != nil {
			metric.Weaknesses = *m.Weaknesses
		}
		domainSummary.MatchSummary.Metrics[i] = metric
	}

	return domainSummary
}

var _ Repository = (*postgresRepository)(nil)
