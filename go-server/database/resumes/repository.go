package resumes

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/features/documents/models/domain"
)

const dateFormat = "Jan. 2006"

type Repository interface {
	UpsertResume(ctx context.Context, firebaseUID string, roleID int, resume *domain.Resume) error
	GetFullResume(ctx context.Context, firebaseUID string, roleID int) (*domain.Resume, error)
}

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) getResumeID(ctx context.Context, runner sqlx.ExtContext, firebaseUID string, roleID int) (int, error) {
	var resumeID int
	query := "SELECT id FROM resumes WHERE firebase_uid = $1 AND role_id = $2"
	var err error

	switch v := runner.(type) {
	case *sqlx.DB:
		err = v.GetContext(ctx, &resumeID, query, firebaseUID, roleID)
	case *sqlx.Tx:
		err = v.GetContext(ctx, &resumeID, query, firebaseUID, roleID)
	default:
		return 0, fmt.Errorf("unsupported runner type: %T", v)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no resume found for user %s and role %d", firebaseUID, roleID)
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

func (r *postgresRepository) UpsertResume(ctx context.Context, firebaseUID string, roleID int, resume *domain.Resume) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	resumeID, err := r.getResumeID(ctx, tx, firebaseUID, roleID)
	if err != nil {
		return err
	}

	if err := r.dropResume(ctx, tx, resumeID); err != nil {
		return fmt.Errorf("failed to drop existing resume data: %w", err)
	}

	for _, s := range resume.Skills {
		var skillID int
		skillQuery := "INSERT INTO skills (resume_id, category, justification_for_changes) VALUES ($1, $2, $3) RETURNING id"
		err := tx.GetContext(ctx, &skillID, skillQuery, resumeID, s.Category, s.JustificationForChanges)
		if err != nil {
			return err
		}

		var skillItems []models.SkillItem
		for _, item := range s.SkillItem {
			skillItems = append(skillItems, models.SkillItem{SkillID: skillID, Name: item})
		}
		if len(skillItems) > 0 {
			_, err = tx.NamedExecContext(ctx, "INSERT INTO skill_items (skill_id, name) VALUES (:skill_id, :name)", skillItems)
			if err != nil {
				return err
			}
		}
	}

	for _, e := range resume.Experiences {
		start, _ := time.Parse(dateFormat, e.Start)
		end, _ := time.Parse(dateFormat, e.End)
		var expID int
		expQuery := "INSERT INTO experiences (resume_id, position, company, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		err := tx.GetContext(ctx, &expID, expQuery, resumeID, e.Position, e.Company, start, end)
		if err != nil {
			return err
		}

		var expDescs []models.ExperienceDescription
		for _, desc := range e.BulletPoints {
			expDescs = append(expDescs, models.ExperienceDescription{ExpID: expID, Text: desc.Text, NewSuggestion: desc.IsNewSuggestion, JustificationForChange: &desc.JustificationForChange})
		}
		if len(expDescs) > 0 {
			_, err = tx.NamedExecContext(ctx, "INSERT INTO experience_descriptions (exp_id, text, new_suggestion, justification_for_change) VALUES (:exp_id, :text, :new_suggestion, :justification_for_change)", expDescs)
			if err != nil {
				return err
			}
		}
	}

	for _, p := range resume.Projects {
		var projectID int
		projQuery := `INSERT INTO projects (resume_id, name, role, status) VALUES ($1, $2, $3, 'ACTIVE') RETURNING id`
		err := tx.GetContext(ctx, &projectID, projQuery, resumeID, p.Name, p.Role)
		if err != nil {
			return err
		}

		var projDescs []models.ProjectDescription
		for _, desc := range p.BulletPoints {
			projDescs = append(projDescs, models.ProjectDescription{ProjectID: projectID, Text: desc.Text, NewSuggestion: desc.IsNewSuggestion, JustificationForChange: &desc.JustificationForChange})
		}
		if len(projDescs) > 0 {
			_, err = tx.NamedExecContext(ctx, "INSERT INTO project_descriptions (project_id, text, new_suggestion, justification_for_change) VALUES (:project_id, :text, :new_suggestion, :justification_for_change)", projDescs)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *postgresRepository) GetFullResume(ctx context.Context, firebaseUID string, roleID int) (*domain.Resume, error) {
	resumeID, err := r.getResumeID(ctx, r.db, firebaseUID, roleID)
	if err != nil {
		return nil, err
	}
	var resume models.Resume
	if err := r.db.GetContext(ctx, &resume, "SELECT * FROM resumes WHERE id = $1", resumeID); err != nil {
		return nil, err
	}

	var experiences []models.Experience
	if err := r.db.SelectContext(ctx, &experiences, "SELECT * FROM experiences WHERE resume_id = $1", resumeID); err != nil {
		return nil, err
	}
	var projects []models.Project
	if err := r.db.SelectContext(ctx, &projects, "SELECT * FROM projects WHERE resume_id = $1", resumeID); err != nil {
		return nil, err
	}

	var skills []models.Skill
	if err := r.db.SelectContext(ctx, &skills, "SELECT * FROM skills WHERE resume_id = $1", resumeID); err != nil {
		return nil, err
	}

	var expDescs []models.ExperienceDescription
	expIDs := getIDs(experiences, func(e models.Experience) int { return e.ID })
	if len(expIDs) > 0 {
		query, args, _ := sqlx.In("SELECT * FROM experience_descriptions WHERE exp_id IN (?)", expIDs)
		if err := r.db.SelectContext(ctx, &expDescs, r.db.Rebind(query), args...); err != nil {
			return nil, err
		}
	}

	var projDescs []models.ProjectDescription
	projIDs := getIDs(projects, func(p models.Project) int { return p.ID })
	if len(projIDs) > 0 {
		query, args, _ := sqlx.In("SELECT * FROM project_descriptions WHERE project_id IN (?)", projIDs)
		if err := r.db.SelectContext(ctx, &projDescs, r.db.Rebind(query), args...); err != nil {
			return nil, err
		}
	}

	var skillItems []models.SkillItem
	skillIDs := getIDs(skills, func(s models.Skill) int { return s.ID })
	if len(skillIDs) > 0 {
		query, args, _ := sqlx.In("SELECT * FROM skill_items WHERE skill_id IN (?)", skillIDs)
		if err := r.db.SelectContext(ctx, &skillItems, r.db.Rebind(query), args...); err != nil {
			return nil, err
		}
	}

	var matchSummary models.MatchSummary
	err = r.db.GetContext(ctx, &matchSummary, "SELECT * FROM match_summaries WHERE resume_id = $1 LIMIT 1", resumeID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	var overviewSummary string
	if err != sql.ErrNoRows {
		var overview models.MatchSummaryOverview
		err = r.db.GetContext(ctx, &overview, "SELECT summary FROM match_summary_overviews WHERE match_summary_id = $1 LIMIT 1", matchSummary.ID)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if err == nil {
			overviewSummary = overview.Summary
		}
	}

	resumePayload := &domain.Resume{}

	if overviewSummary != "" {
		resumePayload.Summary = []domain.SummaryBody{
			{
				Sentence:               overviewSummary,
				JustificationForChange: "",
				NewSuggestion:          false,
			},
		}
	}

	expMap := make(map[int][]models.ExperienceDescription)
	for _, d := range expDescs {
		expMap[d.ExpID] = append(expMap[d.ExpID], d)
	}
	for _, e := range experiences {
		startDateStr := e.StartDate.Format(dateFormat)

		var endDateStr string
		if e.EndDate != nil {
			endDateStr = e.EndDate.Format(dateFormat)
		} else {
			endDateStr = "Present"
		}

		payloadExp := domain.Experience{
			Position: e.Position,
			Company:  e.Company,
			Start:    startDateStr,
			End:      endDateStr,
		}

		resumePayload.Experiences = append(resumePayload.Experiences, payloadExp)
	}

	projMap := make(map[int][]models.ProjectDescription)
	for _, d := range projDescs {
		projMap[d.ProjectID] = append(projMap[d.ProjectID], d)
	}
	for _, p := range projects {
		payloadProj := domain.Project{
			Name: p.Name,
			Role: p.Role,
		}
		for _, desc := range projMap[p.ID] {
			payloadProj.BulletPoints = append(payloadProj.BulletPoints, struct {
				Text                   string `json:"text"`
				IsNewSuggestion        bool   `json:"is_new_suggestion"`
				JustificationForChange string `json:"justification_for_change"`
			}{
				Text:                   desc.Text,
				IsNewSuggestion:        false,
				JustificationForChange: "",
			})
		}
		resumePayload.Projects = append(resumePayload.Projects, payloadProj)
	}

	skillMap := make(map[int][]models.SkillItem)
	for _, i := range skillItems {
		skillMap[i.SkillID] = append(skillMap[i.SkillID], i)
	}
	for _, s := range skills {
		payloadSkill := domain.Skills{
			Category:                s.Category,
			JustificationForChanges: "",
		}
		for _, item := range skillMap[s.ID] {
			payloadSkill.SkillItem = append(payloadSkill.SkillItem, item.Name)
		}
		resumePayload.Skills = append(resumePayload.Skills, payloadSkill)
	}

	return resumePayload, nil
}

func getIDs[T any](slice []T, getID func(T) int) []int {
	ids := make([]int, len(slice))
	for i, item := range slice {
		ids[i] = getID(item)
	}
	return ids
}

var _ Repository = (*postgresRepository)(nil)
