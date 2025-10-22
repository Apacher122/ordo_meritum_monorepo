package resumes

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/features/documents/models/domain"
	"github.com/ordo_meritum/shared/contexts"
	error_response "github.com/ordo_meritum/shared/types/errors"
	error_messages "github.com/ordo_meritum/shared/utils/errors"
	"golang.org/x/sync/errgroup"
)

const dateFormat = "Jan. 2006"

type Repository interface {
	UpsertResume(ctx context.Context, roleID int, resume *domain.Resume, education *domain.EducationInfo) error
	GetFullResume(ctx context.Context, roleID int) (*domain.Resume, error)
	GetFullEducation(ctx context.Context, roleID int) ([]*domain.EducationInfo, error)
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

func (r *postgresRepository) GetFullEducation(
	ctx context.Context,
	jobID int,
) ([]*domain.EducationInfo, error) {
	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		return nil, error_messages.ErrorMessage(error_messages.ERR_USER_NO_CONTEXT)
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		SELECT e.*
		FROM education e
		JOIN resumes r ON e.resume_id = r.id
		WHERE r.firebase_uid = $1 AND r.role_id = $2;
	`

	var educationModels []*models.Education
	if err := tx.SelectContext(ctx, &educationModels, query, userCtx.UID, jobID); err != nil {
		return nil, fmt.Errorf("failed to get education for user %s and job %d: %w", userCtx.UID, jobID, err)
	}

	var domainEducations []*domain.EducationInfo
	for _, dbModel := range educationModels {
		domainModel, err := transformEducationToDomainModel(dbModel)
		if err != nil {
			fmt.Printf("Warning: failed to transform education model for school %s: %v\n", dbModel.School, err)
			continue
		}
		domainEducations = append(domainEducations, domainModel)
	}

	return domainEducations, nil
}

func (r *postgresRepository) UpsertResume(ctx context.Context, roleID int, resume *domain.Resume, education *domain.EducationInfo) error {
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
		return fmt.Errorf("failed to upsert skills %w", err)
	}

	if err := r.UpsertExperiences(ctx, tx, resumeID, resume.Experiences); err != nil {
		return fmt.Errorf("failed to upsert experiences %w", err)
	}

	if err := r.UpsertProjects(ctx, tx, resumeID, resume.Projects); err != nil {
		return fmt.Errorf("failed to upsert projects %w", err)
	}

	if err := r.UpsertEducation(ctx, tx, resumeID, education); err != nil {
		return fmt.Errorf("failed to upsert educations %w", err)
	}

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

func transformEducationToDomainModel(dbModel *models.Education) (*domain.EducationInfo, error) {
	if dbModel == nil {
		return nil, fmt.Errorf("database model cannot be nil")
	}

	domainModel := &domain.EducationInfo{
		School:     dbModel.School,
		Degree:     dbModel.Degree,
		Location:   dbModel.Location,
		StartEnd:   dbModel.Dates,
		CourseWork: dbModel.Courses,
		GPA:        dbModel.GPA,
		Honors:     dbModel.Honors,
	}

	return domainModel, nil
}

var _ Repository = (*postgresRepository)(nil)
