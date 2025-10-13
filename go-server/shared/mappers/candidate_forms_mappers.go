package mappers

import (
	db_models "github.com/ordo_meritum/database/models"
	"github.com/ordo_meritum/features/candidate_forms/models/dto"
)

func MapDBToDTO(dbOcean db_models.OceanProfile, dbDisc db_models.DiscProfile) dto.PersonalitySummary {
	oceanScores := []dto.OCEANScore{
		{Category: dto.Openness, Score: dbOcean.OpennessScore, Reasoning: dbOcean.OpennessReasoning},
		{Category: dto.Conscientiousness, Score: dbOcean.ConscientiousnessScore, Reasoning: dbOcean.ConscientiousnessReasoning}, {Category: dto.Extraversion, Score: dbOcean.ExtraversionScore, Reasoning: dbOcean.ExtraversionReasoning},
		{Category: dto.Agreeableness, Score: dbOcean.AgreeablenessScore, Reasoning: dbOcean.AgreeablenessReasoning},
		{Category: dto.Neuroticism, Score: dbOcean.NeuroticismScore, Reasoning: dbOcean.NeuroticismReasoning},
	}

	oceanProfile := &dto.OCEANProfile{
		Scores:  oceanScores,
		Summary: dbOcean.Summary,
	}

	discScores := []dto.DISCScore{
		{Category: dto.Dominance, Reasoning: dbDisc.Dominance},
		{Category: dto.Influence, Reasoning: dbDisc.Influence},
		{Category: dto.Steadiness, Reasoning: dbDisc.Steadiness},
		{Category: dto.Consistency, Reasoning: dbDisc.Conscientiousness},
	}

	discProfile := &dto.DISCProfile{
		Scores:  discScores,
		Summary: dbDisc.Summary,
	}

	return dto.PersonalitySummary{
		OCEAN: *oceanProfile,
		DISC:  *discProfile,
	}
}

func MapDTOToDB(summary dto.PersonalitySummary) (db_models.OceanProfile, db_models.DiscProfile) {
	var dbOcean db_models.OceanProfile
	var dbDisc db_models.DiscProfile

	dbOcean.Summary = summary.OCEAN.Summary
	for _, score := range summary.OCEAN.Scores {
		switch score.Category {
		case dto.Openness:
			dbOcean.OpennessScore = score.Score
			dbOcean.OpennessReasoning = score.Reasoning
		case dto.Conscientiousness:
			dbOcean.ConscientiousnessScore = score.Score
			dbOcean.ConscientiousnessReasoning = score.Reasoning
		case dto.Extraversion:
			dbOcean.ExtraversionScore = score.Score
			dbOcean.ExtraversionReasoning = score.Reasoning
		case dto.Agreeableness:
			dbOcean.AgreeablenessScore = score.Score
			dbOcean.AgreeablenessReasoning = score.Reasoning
		case dto.Neuroticism:
			dbOcean.NeuroticismScore = score.Score
			dbOcean.NeuroticismReasoning = score.Reasoning
		}
	}

	dbDisc.Summary = summary.DISC.Summary
	for _, score := range summary.DISC.Scores {
		switch score.Category {
		case dto.Dominance:
			dbDisc.Dominance = score.Reasoning
		case dto.Influence:
			dbDisc.Influence = score.Reasoning
		case dto.Steadiness:
			dbDisc.Steadiness = score.Reasoning
		case dto.Consistency:
			dbDisc.Conscientiousness = score.Reasoning
		}
	}

	return dbOcean, dbDisc
}
