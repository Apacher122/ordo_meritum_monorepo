package jobs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/features/application_tracking/models/domain"
	request "github.com/ordo_meritum/features/application_tracking/models/requests"
	"github.com/ordo_meritum/shared/contexts"
	error_response "github.com/ordo_meritum/shared/types/errors"
	em "github.com/ordo_meritum/shared/utils/errors"
)

type FullJobPosting struct {
	JobTitle               string         `db:"job_title"`
	Description            *string        `db:"description"`
	CompanyName            string         `db:"company_name"`
	CompanyProperName      string         `db:"proper_name"`
	CompanyCulture         *string        `db:"company_culture"`
	CompanyValues          *string        `db:"company_values"`
	Requirements           pq.StringArray `db:"requirements"`
	NiceToHaves            pq.StringArray `db:"nice_to_haves"`
	EducationLevel         *string        `db:"education_level"`
	YearsOfExp             *string        `db:"years_of_exp"`
	Tools                  pq.StringArray `db:"tools"`
	ProgrammingLanguages   pq.StringArray `db:"programming_languages"`
	FrameworksAndLibraries pq.StringArray `db:"frameworks_and_libraries"`
	Databases              pq.StringArray `db:"databases"`
	CloudTechnologies      pq.StringArray `db:"cloud_technologies"`
	IndustryKeywords       pq.StringArray `db:"industry_keywords"`
	SoftSkills             pq.StringArray `db:"soft_skills"`
	Certifications         pq.StringArray `db:"certifications"`
	ApplicantCount         *int           `db:"applicant_count"`
	SalaryRange            *string        `db:"salary_range"`
}

type UserJobPosting struct {
	RoleID                 int               `db:"role_id"`
	JobTitle               string            `db:"job_title"`
	CompanyName            string            `db:"company_name"`
	CompanyProperName      string            `db:"proper_name"`
	Website                *string           `db:"website"`
	ApplicationStatus      *models.AppStatus `db:"application_status"`
	UserApplied            *bool             `db:"user_applied"`
	InterviewCount         *int              `db:"interview_count"`
	InitialApplicationDate *time.Time        `db:"initial_application_date"`
}

type Repository interface {
	GetFullJobPosting(ctx context.Context, roleID int) (*FullJobPosting, *em.ErrorBody)
	InsertFullJobPosting(ctx context.Context, jobRawText string, jobPost *domain.JobDescription, companyName string, properName string) (*models.JobRequirements, *em.ErrorBody)
	GetAllUserJobPostings(ctx context.Context) ([]*UserJobPosting, *em.ErrorBody)
	UpdateApplicationDetails(ctx context.Context, roleID int, req *request.ApplicationUpdateRequest) error
	DeleteJobPostByID(ctx context.Context, roleID int) error
}

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetFullJobPosting(ctx context.Context, roleID int) (*FullJobPosting, *em.ErrorBody) {
	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		return nil, &em.ErrorBody{ErrCode: em.ERR_USER_NO_CONTEXT, ErrMsg: em.ErrorMessage(em.ERR_USER_NO_CONTEXT)}
	}

	query := `
        SELECT
            r.job_title, r.description,
            c.company_name, c.proper_name, c.company_culture, c.company_values,
            j.requirements, j.nice_to_haves, j.education_level, j.years_of_exp,
            j.tools, j.programming_languages, j.frameworks_and_libraries, j.databases,
            j.cloud_technologies, j.industry_keywords, j.soft_skills, j.certifications,
            j.applicant_count,
            r.salary_range
        FROM roles r
        INNER JOIN companies c ON r.company_id = c.id
        INNER JOIN job_requirements j ON r.id = j.role_id
				INNER JOIN resumes res ON r.id = res.role_id
        WHERE r.id = $1 and res.firebase_uid = $2`
	var job FullJobPosting
	err := r.db.GetContext(ctx, &job, query, roleID, userCtx.UID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &em.ErrorBody{ErrCode: em.ERR_DB_FAILED_TO_GET, ErrMsg: fmt.Errorf("job with role ID %d not found", roleID)}
		}
		return nil, &em.ErrorBody{ErrCode: em.ERR_DB_FAILED_TO_GET, ErrMsg: em.ErrorMessage(em.ERR_DB_FAILED_TO_GET)}
	}
	return &job, nil
}

func (r *postgresRepository) InsertFullJobPosting(
	ctx context.Context,
	jobRawText string,
	jobPost *domain.JobDescription,
	companyName string,
	properName string,
) (*models.JobRequirements, *em.ErrorBody) {
	userCtx, _ := contexts.FromContext(ctx)
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, &em.ErrorBody{ErrCode: em.ERR_DB_FAILED_TO_INSERT, ErrMsg: err}
	}
	defer tx.Rollback()

	var companyID int
	companyQuery := `
        INSERT INTO companies (company_name, proper_name, company_culture, company_values)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (company_name) DO UPDATE
        SET company_culture = EXCLUDED.company_culture, company_values = EXCLUDED.company_values
        RETURNING id`
	err = tx.GetContext(ctx, &companyID, companyQuery, companyName, jobPost.CompanyName, jobPost.CompanyCulture, jobPost.CompanyValues)
	if err != nil {
		return nil, &em.ErrorBody{ErrCode: em.ERR_DB_FAILED_TO_UPSERT, ErrMsg: fmt.Errorf("failed to upsert company: %w", err)}
	}

	var roleID int
	roleQuery := `
        INSERT INTO roles (job_title, description, company_id, salary_range)
        VALUES ($1, $2, $3, $4) RETURNING id`
	err = tx.GetContext(ctx, &roleID, roleQuery, jobPost.JobTitle, jobRawText, companyID, jobPost.SalaryRange)
	if err != nil {
		return nil, &em.ErrorBody{ErrCode: em.ERR_DB_FAILED_TO_INSERT, ErrMsg: fmt.Errorf("failed to create role: %w", err)}
	}

	resumeQuery := `INSERT INTO resumes (role_id, firebase_uid) VALUES ($1, $2)`
	_, err = tx.ExecContext(ctx, resumeQuery, roleID, userCtx.UID)
	if err != nil {
		return nil, &em.ErrorBody{ErrCode: em.ERR_DB_FAILED_TO_INSERT, ErrMsg: fmt.Errorf("failed to create resume entry: %w", err)}
	}

	reqs := models.JobRequirements{
		RoleID:                 roleID,
		EducationLevel:         &jobPost.EducationLevel,
		ApplicantCount:         &jobPost.ApplicantCount,
		YearsOfExp:             &jobPost.YearsOfExp,
		Tools:                  jobPost.ToolsAndTechnologies,
		ProgrammingLanguages:   jobPost.ProgrammingLanguages,
		FrameworksAndLibraries: jobPost.FrameworksAndLibraries,
		Databases:              jobPost.Databases,
		CloudTechnologies:      jobPost.CloudTechnologies,
		IndustryKeywords:       jobPost.IndustryKeywords,
		Requirements:           jobPost.SkillsRequired,
		NiceToHaves:            jobPost.SkillsNiceToHaves,
		SoftSkills:             jobPost.SoftSkills,
		Certifications:         jobPost.Certifications,
	}
	reqQuery := `
        INSERT INTO job_requirements (role_id, education_level, applicant_count, years_of_exp, tools, programming_languages, frameworks_and_libraries, databases, cloud_technologies, industry_keywords, requirements, nice_to_haves, soft_skills, certifications)
        VALUES (:role_id, :education_level, :applicant_count, :years_of_exp, :tools, :programming_languages, :frameworks_and_libraries, :databases, :cloud_technologies, :industry_keywords, :requirements, :nice_to_haves, :soft_skills, :certifications)
        RETURNING *`
	rows, err := tx.NamedQuery(reqQuery, &reqs)
	if err != nil {
		return nil, &em.ErrorBody{ErrCode: em.ERR_DB_FAILED_TO_INSERT, ErrMsg: fmt.Errorf("failed to create job requirements: %w", err)}
	}
	defer rows.Close()

	var createdReqs models.JobRequirements
	if rows.Next() {
		if err := rows.StructScan(&createdReqs); err != nil {
			return nil, &em.ErrorBody{ErrCode: em.ERR_DB_FAILED_TO_INSERT, ErrMsg: err}
		}
	}
	err = tx.Commit()
	return &createdReqs, &em.ErrorBody{ErrCode: em.ERR_DB_FAILED_TO_INSERT, ErrMsg: err}
}

func (r *postgresRepository) GetAllUserJobPostings(ctx context.Context) ([]*UserJobPosting, *em.ErrorBody) {
	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		return nil, &em.ErrorBody{ErrCode: em.ERR_USER_NO_CONTEXT, ErrMsg: em.ErrorMessage(em.ERR_USER_NO_CONTEXT)}
	}
	query := `
        SELECT
            r.id as role_id,
            r.job_title as job_title,
            c.company_name as company_name,
						c.proper_name as proper_name,
            c.website as website,
            r.application_status as application_status,
            r.user_applied as user_applied,
            j.interview_count as interview_count,
            res.applied_on as initial_application_date
        FROM resumes res
        INNER JOIN roles r ON res.role_id = r.id
        INNER JOIN companies c ON r.company_id = c.id
        INNER JOIN job_requirements j ON r.id = j.role_id
        WHERE res.firebase_uid = $1`
	var jobs []*UserJobPosting
	err := r.db.SelectContext(ctx, &jobs, query, userCtx.UID)
	if err != nil {
		return jobs, &em.ErrorBody{ErrCode: em.ERR_DB_FAILED_TO_GET, ErrMsg: err}
	}
	return jobs, nil
}

func (r *postgresRepository) UpdateApplicationDetails(ctx context.Context, roleID int, req *request.ApplicationUpdateRequest) error {
	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		return error_response.ErrNoUserContext
	}

	payload := req.Payload

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM resumes WHERE role_id = $1 AND firebase_uid = $2)"
	err := r.db.QueryRow(query, roleID, userCtx.UID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking user authorization for role %d: %w", roleID, err)
	}
	if !exists {
		return fmt.Errorf("user not authorized to update application for role ID %d", roleID)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	if payload.Link != nil {
		var companyID int
		err := tx.QueryRow("SELECT company_id FROM roles WHERE id = $1", roleID).Scan(&companyID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("no role found with id %d", roleID)
			}
			return fmt.Errorf("could not query company for role %d: %w", roleID, err)
		}
		_, err = tx.Exec("UPDATE companies SET website = $1, updated_at = NOW() WHERE id = $2", payload.Link, companyID)
		if err != nil {
			return fmt.Errorf("failed to update company website: %w", err)
		}
	}

	roleUpdates := make(map[string]interface{})
	if payload.JobTitle != nil {
		roleUpdates["job_title"] = *payload.JobTitle
	}
	if payload.ApplicationStatus != nil {
		roleUpdates["application_status"] = *payload.ApplicationStatus
	}

	if len(roleUpdates) > 0 {
		if err := executeUpdate(tx, "roles", roleUpdates, "id", roleID); err != nil {
			return fmt.Errorf("failed to update role: %w", err)
		}
	}

	reqUpdates := make(map[string]interface{})
	if payload.InterviewCount != nil {
		reqUpdates["interview_count"] = *payload.InterviewCount
	}
	if payload.InitialApplicationDate != nil {
		reqUpdates["initial_application_date"] = *payload.InitialApplicationDate
	}

	if len(reqUpdates) > 0 {
		if err := executeUpdate(tx, "job_requirements", reqUpdates, "role_id", roleID); err != nil {
			return fmt.Errorf("failed to update job requirements: %w", err)
		}
	}

	return tx.Commit()
}

func (r *postgresRepository) DeleteJobPostByID(ctx context.Context, roleID int) error {
	userCtx, ok := contexts.FromContext(ctx)
	if !ok {
		return error_response.ErrNoUserContext
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	authQuery := "SELECT EXISTS(SELECT 1 FROM resumes WHERE role_id = $1 AND firebase_uid = $2)"
	if err := tx.GetContext(ctx, &exists, authQuery, roleID, userCtx.UID); err != nil {
		return fmt.Errorf("authorization check failed: %w", err)
	}
	if !exists {
		return errors.New("user is not authorized to delete this job posting or it does not exist")
	}

	if _, err := tx.ExecContext(ctx, "DELETE FROM job_requirements WHERE role_id = $1", roleID); err != nil {
		return fmt.Errorf("failed to delete from job_requirements: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM resumes WHERE role_id = $1 AND firebase_uid = $2", roleID, userCtx.UID); err != nil {
		return fmt.Errorf("failed to delete from resumes: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM roles WHERE id = $1", roleID); err != nil {
		return fmt.Errorf("failed to delete from roles: %w", err)
	}

	return tx.Commit()
}

func executeUpdate(tx *sql.Tx, table string, updates map[string]interface{}, whereCol string, whereVal int) error {
	var setClauses []string
	args := make([]interface{}, 0, len(updates)+1)
	i := 1

	for key, val := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, i))
		args = append(args, val)
		i++
	}

	query := fmt.Sprintf("UPDATE %s SET %s, updated_at = NOW() WHERE %s = $%d",
		table,
		strings.Join(setClauses, ", "),
		whereCol,
		i,
	)
	args = append(args, whereVal)

	result, err := tx.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows found to update in table %s for id %d", table, whereVal)
	}

	return nil
}

var _ Repository = (*postgresRepository)(nil)
