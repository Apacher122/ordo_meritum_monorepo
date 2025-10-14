package resumes

import (
	"context"

	"github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/features/documents/models/domain"

	"github.com/jmoiron/sqlx"
)

func (r *postgresRepository) UpsertSkills(
	ctx context.Context,
	tx *sqlx.Tx,
	resumeID int,
	resume *domain.Resume,
) error {
	for _, s := range resume.Skills {
		var skillID int
		skillQuery := "INSERT INTO skills (resume_id, category, justification_for_changes) VALUES ($1, $2, $3) RETURNING id"
		err := tx.GetContext(ctx, &skillID, skillQuery, resumeID, s.Category, s.JustificationForChanges)
		if err != nil {
			return err
		}

		if len(s.SkillItem) == 0 {
			continue
		}

		var skillItems []models.SkillItem
		for _, item := range s.SkillItem {
			skillItems = append(skillItems, models.SkillItem{SkillID: skillID, Name: item})
		}

		_, err = tx.NamedExecContext(ctx, "INSERT INTO skill_items (skill_id, name) VALUES (:skill_id, :name)", skillItems)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *postgresRepository) GetResumeSkills(ctx context.Context, resumeID int) ([]domain.Skills, error) {
	var skills []models.Skill
	if err := r.db.SelectContext(ctx, &skills, "SELECT * FROM skills WHERE resume_id = $1", resumeID); err != nil {
		return nil, err
	}
	if len(skills) == 0 {
		return nil, nil
	}

	skillIDs := getIDs(skills, func(s models.Skill) int { return s.ID })
	var skillItems []models.SkillItem
	query, args, _ := sqlx.In("SELECT * FROM skill_items WHERE skill_id IN (?)", skillIDs)
	if err := r.db.SelectContext(ctx, &skillItems, r.db.Rebind(query), args...); err != nil {
		return nil, err
	}

	return MapSkillsToDomain(skills, skillItems), nil
}
