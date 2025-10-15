package resumes

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ordo_meritum/features/documents/models/domain"
	"github.com/ordo_meritum/shared/contexts"
	error_response "github.com/ordo_meritum/shared/types/errors"
	"golang.org/x/sync/errgroup"
)

const dateFormat = "Jan. 2006"

type Repository interface {
	UpsertResume(ctx context.Context, roleID int, resume *domain.Resume) error
	GetFullResume(ctx context.Context, roleID int) (*domain.Resume, error)
}

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) getResumeID(ctx context.Context, runner sqlx.ExtContext, roleID int) (int, error) {
	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		return 0, error_response.ErrNoUserContext
	}

	var resumeID int
	query := "SELECT id FROM resumes WHERE firebase_uid = $1 AND role_id = $2"
	var err error

	switch v := runner.(type) {
	case *sqlx.DB:
		err = v.GetContext(ctx, &resumeID, query, userCtx.UID, roleID)
	case *sqlx.Tx:
		err = v.GetContext(ctx, &resumeID, query, userCtx.UID, roleID)
	default:
		return 0, fmt.Errorf("unsupported runner type: %T", v)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no resume found for user %s and role %d", userCtx.UID, roleID)
		}
		return 0, err
	}
	return resumeID, nil
}

func (r *postgresRepository) dropResume(ctx context.Context, tx *sqlx.Tx, resumeID int) error {
	subqueryExp := "SELECT id FROM experiences WHERE resume_id = $1"
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("DELETE FROM experience_descriptions WHERE exp_id IN (%s)", subqueryExp), resumeID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM experiences WHERE resume_id = $1", resumeID); err != nil {
		return err
	}

	subqueryProj := "SELECT id FROM projects WHERE resume_id = $1"
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("DELETE FROM project_descriptions WHERE project_id IN (%s)", subqueryProj), resumeID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM projects WHERE resume_id = $1", resumeID); err != nil {
		return err
	}

	subquerySkill := "SELECT id FROM skills WHERE resume_id = $1"
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("DELETE FROM skill_items WHERE skill_id IN (%s)", subquerySkill), resumeID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM skills WHERE resume_id = $1", resumeID); err != nil {
		return err
	}

	return nil
}

func (r *postgresRepository) UpsertResume(ctx context.Context, roleID int, resume *domain.Resume) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	resumeID, err := r.getResumeID(ctx, tx, roleID)
	if err != nil {
		return err
	}

	if err := r.dropResume(ctx, tx, resumeID); err != nil {
		return fmt.Errorf("failed to drop existing resume data: %w", err)
	}

	if err := r.UpsertSkills(ctx, tx, resumeID, resume); err != nil {
		return fmt.Errorf("failed to upsert skills")
	}

	// for _, s := range resume.Skills {
	// 	var skillID int
	// 	skillQuery := "INSERT INTO skills (resume_id, category, justification_for_changes) VALUES ($1, $2, $3) RETURNING id"
	// 	err := tx.GetContext(ctx, &skillID, skillQuery, resumeID, s.Category, s.JustificationForChanges)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	var skillItems []models.SkillItem
	// 	for _, item := range s.SkillItem {
	// 		skillItems = append(skillItems, models.SkillItem{SkillID: skillID, Name: item})
	// 	}
	// 	if len(skillItems) > 0 {
	// 		_, err = tx.NamedExecContext(ctx, "INSERT INTO skill_items (skill_id, name) VALUES (:skill_id, :name)", skillItems)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	if err := r.UpsertExperiences(ctx, tx, resumeID, resume.Experiences); err != nil {
		return fmt.Errorf("failed to upsert experiences")
	}

	// for _, e := range resume.Experiences {
	// 	start, _ := time.Parse(dateFormat, e.Start)
	// 	end, _ := time.Parse(dateFormat, e.End)
	// 	var expID int
	// 	expQuery := "INSERT INTO experiences (resume_id, position, company, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	// 	err := tx.GetContext(ctx, &expID, expQuery, resumeID, e.Position, e.Company, start, end)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	var expDescs []models.ExperienceDescription
	// 	for _, desc := range e.BulletPoints {
	// 		expDescs = append(expDescs, models.ExperienceDescription{ExpID: expID, Text: desc.Text, NewSuggestion: desc.IsNewSuggestion, JustificationForChange: &desc.JustificationForChange})
	// 	}
	// 	if len(expDescs) > 0 {
	// 		_, err = tx.NamedExecContext(ctx, "INSERT INTO experience_descriptions (exp_id, text, new_suggestion, justification_for_change) VALUES (:exp_id, :text, :new_suggestion, :justification_for_change)", expDescs)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	if err := r.UpsertProjects(ctx, tx, resumeID, resume.Projects); err != nil {
		return fmt.Errorf("failed to upsert projects")
	}

	// for _, p := range resume.Projects {
	// 	var projectID int
	// 	projQuery := `INSERT INTO projects (resume_id, name, role, status) VALUES ($1, $2, $3, 'ACTIVE') RETURNING id`
	// 	err := tx.GetContext(ctx, &projectID, projQuery, resumeID, p.Name, p.Role)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	var projDescs []models.ProjectDescription
	// 	for _, desc := range p.BulletPoints {
	// 		projDescs = append(projDescs, models.ProjectDescription{ProjectID: projectID, Text: desc.Text, NewSuggestion: desc.IsNewSuggestion, JustificationForChange: &desc.JustificationForChange})
	// 	}
	// 	if len(projDescs) > 0 {
	// 		_, err = tx.NamedExecContext(ctx, "INSERT INTO project_descriptions (project_id, text, new_suggestion, justification_for_change) VALUES (:project_id, :text, :new_suggestion, :justification_for_change)", projDescs)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return tx.Commit()
}

func (r *postgresRepository) GetFullResume(
	ctx context.Context,
	roleID int,
) (*domain.Resume, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	resumeID, err := r.getResumeID(ctx, tx, roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &domain.Resume{}, nil
		}
		return nil, err
	}

	var resumePayload domain.Resume
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		resumePayload.Experiences, err = r.GetResumeExperiences(gCtx, resumeID)
		return err
	})

	g.Go(func() error {
		var err error
		resumePayload.Projects, err = r.GetResumeProjects(gCtx, resumeID)
		return err
	})

	g.Go(func() error {
		var err error
		resumePayload.Skills, err = r.GetResumeSkills(gCtx, resumeID)
		return err
	})

	g.Go(func() error {
		var err error
		resumePayload.Summary, err = r.GetResumeSummary(gCtx, resumeID)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &resumePayload, nil
}

func getIDs[T any](slice []T, getID func(T) int) []int {
	ids := make([]int, len(slice))
	for i, item := range slice {
		ids[i] = getID(item)
	}
	return ids
}

var _ Repository = (*postgresRepository)(nil)
