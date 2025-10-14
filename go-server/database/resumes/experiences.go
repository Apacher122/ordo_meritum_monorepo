package resumes

import (
	"context"
	"time"

	"github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/features/documents/models/domain"

	"github.com/jmoiron/sqlx"
)

func (r *postgresRepository) UpsertExperiences(
	ctx context.Context,
	tx *sqlx.Tx,
	resumeID int,
	experiences []domain.Experience,
) error {
	for _, e := range experiences {
		start, _ := time.Parse(dateFormat, e.Start)
		end, _ := time.Parse(dateFormat, e.End)
		var expID int
		expQuery := "INSERT INTO experiences (resume_id, position, company, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		err := tx.GetContext(ctx, &expID, expQuery, resumeID, e.Position, e.Company, start, end)
		if err != nil {
			return err
		}

		if len(e.BulletPoints) == 0 {
			continue
		}

		var expDescs []models.ExperienceDescription
		for _, desc := range e.BulletPoints {
			expDescs = append(expDescs, models.ExperienceDescription{ExpID: expID, Text: desc.Text, NewSuggestion: desc.IsNewSuggestion, JustificationForChange: &desc.JustificationForChange})
		}

		_, err = tx.NamedExecContext(ctx, "INSERT INTO experience_descriptions (exp_id, text, new_suggestion, justification_for_change) VALUES (:exp_id, :text, :new_suggestion, :justification_for_change)", expDescs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *postgresRepository) GetResumeExperiences(ctx context.Context, resumeID int) ([]domain.Experience, error) {
	var experiences []models.Experience
	if err := r.db.SelectContext(ctx, &experiences, "SELECT * FROM experiences WHERE resume_id = $1", resumeID); err != nil {
		return nil, err
	}
	if len(experiences) == 0 {
		return nil, nil
	}

	expIDs := getIDs(experiences, func(e models.Experience) int { return e.ID })
	var expDescs []models.ExperienceDescription
	query, args, _ := sqlx.In("SELECT * FROM experience_descriptions WHERE exp_id IN (?)", expIDs)
	if err := r.db.SelectContext(ctx, &expDescs, r.db.Rebind(query), args...); err != nil {
		return nil, err
	}

	return MapExperiencesToDomain(experiences, expDescs), nil
}
