package resumes

import (
	"context"

	"github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/features/documents/models/domain"

	"github.com/jmoiron/sqlx"
)

func (r *postgresRepository) UpsertProjects(
	ctx context.Context,
	tx *sqlx.Tx,
	resumeID int,
	projects []domain.Project,
) error {
	for _, p := range projects {
		var projectID int
		projQuery := `INSERT INTO projects (resume_id, name, role, status) VALUES ($1, $2, $3, 'ACTIVE') RETURNING id`
		err := tx.GetContext(ctx, &projectID, projQuery, resumeID, p.Name, p.Role)
		if err != nil {
			return err
		}

		if len(p.BulletPoints) == 0 {
			continue
		}

		var projDescs []models.ProjectDescription
		for _, desc := range p.BulletPoints {
			projDescs = append(projDescs, models.ProjectDescription{ProjectID: projectID, Text: desc.Text, NewSuggestion: desc.IsNewSuggestion, JustificationForChange: &desc.JustificationForChange})
		}

		_, err = tx.NamedExecContext(ctx, "INSERT INTO project_descriptions (project_id, text, new_suggestion, justification_for_change) VALUES (:project_id, :text, :new_suggestion, :justification_for_change)", projDescs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *postgresRepository) GetResumeProjects(ctx context.Context, resumeID int) ([]domain.Project, error) {
	var projects []models.Project
	if err := r.db.SelectContext(ctx, &projects, "SELECT * FROM projects WHERE resume_id = $1", resumeID); err != nil {
		return nil, err
	}
	if len(projects) == 0 {
		return nil, nil
	}

	projIDs := getIDs(projects, func(p models.Project) int { return p.ID })
	var projDescs []models.ProjectDescription
	query, args, _ := sqlx.In("SELECT * FROM project_descriptions WHERE project_id IN (?)", projIDs)
	if err := r.db.SelectContext(ctx, &projDescs, r.db.Rebind(query), args...); err != nil {
		return nil, err
	}

	return MapProjectsToDomain(projects, projDescs), nil
}
