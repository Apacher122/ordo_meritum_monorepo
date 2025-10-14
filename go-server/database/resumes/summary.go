package resumes

import (
	"context"
	"database/sql"

	"github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/features/documents/models/domain"
)

func (r *postgresRepository) GetResumeSummary(ctx context.Context, resumeID int) ([]domain.SummaryBody, error) {
	var matchSummary models.MatchSummary
	err := r.db.GetContext(ctx, &matchSummary, "SELECT * FROM match_summaries WHERE resume_id = $1 LIMIT 1", resumeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No summary is not an error
		}
		return nil, err
	}

	var overview models.MatchSummaryOverview
	err = r.db.GetContext(ctx, &overview, "SELECT summary FROM match_summary_overviews WHERE match_summary_id = $1 LIMIT 1", matchSummary.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No overview is not an error
		}
		return nil, err
	}

	return MapSummaryToDomain(overview.Summary), nil
}
