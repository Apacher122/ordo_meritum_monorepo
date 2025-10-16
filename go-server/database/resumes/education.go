package resumes

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/features/documents/models/domain"
)

func (r *postgresRepository) UpsertEducation(
	ctx context.Context,
	tx *sqlx.Tx,
	resumeID int,
	education *domain.EducationInfo,
) error {
	educationUpsertQuery := `
		INSERT INTO education (
			resume_id, school, degree, dates, location, courses, gpa, honors
		) VALUES (
			:resume_id, :school, :degree, :dates, :location, :courses, :gpa, :honors
		)
		ON CONFLICT (resume_id, school, degree) DO UPDATE SET
			dates = EXCLUDED.dates,
			location = EXCLUDED.location,
			courses = EXCLUDED.courses,
			gpa = EXCLUDED.gpa,
			honors = EXCLUDED.honors,
			updated_at = NOW();`

	dbModel, err := transformEducationToDBModel(education, resumeID)
	if err != nil {
		return fmt.Errorf("failed to transform education domain model: %w", err)
	}

	_, err = tx.NamedExecContext(ctx, educationUpsertQuery, dbModel)
	if err != nil {
		return fmt.Errorf("failed to upsert education for school %s: %w", dbModel.School, err)
	}

	return nil
}

func transformEducationToDBModel(info *domain.EducationInfo, resumeID int) (*models.Education, error) {
	dbModel := &models.Education{
		ResumeID: resumeID,
		School:   info.School,
		Degree:   info.Degree,
		Location: info.Location,
		Dates:    info.StartEnd,
		Courses:  info.CourseWork,
		GPA:      info.GPA,
		Honors:   info.Honors,
	}

	return dbModel, nil
}
