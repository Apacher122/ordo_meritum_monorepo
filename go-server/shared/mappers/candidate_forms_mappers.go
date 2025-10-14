package mappers

import (
	db_models "github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/features/candidate_forms/models/domain"
)

func MapDBToDTO(dbOcean db_models.OceanProfile, dbDisc db_models.DiscProfile) domain.PersonalitySummary {
	oceanScores := []domain.OCEANScore{
		{Category: domain.Openness, Score: dbOcean.OpennessScore, Reasoning: dbOcean.OpennessReasoning},
		{Category: domain.Conscientiousness, Score: dbOcean.ConscientiousnessScore, Reasoning: dbOcean.ConscientiousnessReasoning}, {Category: domain.Extraversion, Score: dbOcean.ExtraversionScore, Reasoning: dbOcean.ExtraversionReasoning},
		{Category: domain.Agreeableness, Score: dbOcean.AgreeablenessScore, Reasoning: dbOcean.AgreeablenessReasoning},
		{Category: domain.Neuroticism, Score: dbOcean.NeuroticismScore, Reasoning: dbOcean.NeuroticismReasoning},
	}

	oceanProfile := &domain.OCEANProfile{
		Scores:  oceanScores,
		Summary: dbOcean.Summary,
	}

	discScores := []domain.DISCScore{
		{Category: domain.Dominance, Reasoning: dbDisc.Dominance},
		{Category: domain.Influence, Reasoning: dbDisc.Influence},
		{Category: domain.Steadiness, Reasoning: dbDisc.Steadiness},
		{Category: domain.Consistency, Reasoning: dbDisc.Conscientiousness},
	}

	discProfile := &domain.DISCProfile{
		Scores:  discScores,
		Summary: dbDisc.Summary,
	}

	return domain.PersonalitySummary{
		OCEAN: *oceanProfile,
		DISC:  *discProfile,
	}
}

func MapDTOToDB(summary domain.PersonalitySummary) (db_models.OceanProfile, db_models.DiscProfile) {
	var dbOcean db_models.OceanProfile
	var dbDisc db_models.DiscProfile

	dbOcean.Summary = summary.OCEAN.Summary
	for _, score := range summary.OCEAN.Scores {
		switch score.Category {
		case domain.Openness:
			dbOcean.OpennessScore = score.Score
			dbOcean.OpennessReasoning = score.Reasoning
		case domain.Conscientiousness:
			dbOcean.ConscientiousnessScore = score.Score
			dbOcean.ConscientiousnessReasoning = score.Reasoning
		case domain.Extraversion:
			dbOcean.ExtraversionScore = score.Score
			dbOcean.ExtraversionReasoning = score.Reasoning
		case domain.Agreeableness:
			dbOcean.AgreeablenessScore = score.Score
			dbOcean.AgreeablenessReasoning = score.Reasoning
		case domain.Neuroticism:
			dbOcean.NeuroticismScore = score.Score
			dbOcean.NeuroticismReasoning = score.Reasoning
		}
	}

	dbDisc.Summary = summary.DISC.Summary
	for _, score := range summary.DISC.Scores {
		switch score.Category {
		case domain.Dominance:
			dbDisc.Dominance = score.Reasoning
		case domain.Influence:
			dbDisc.Influence = score.Reasoning
		case domain.Steadiness:
			dbDisc.Steadiness = score.Reasoning
		case domain.Consistency:
			dbDisc.Conscientiousness = score.Reasoning
		}
	}

	return dbOcean, dbDisc
}
