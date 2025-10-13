package dto

import "github.com/ordo_meritum/database/models"

type MatchSummaryPayload struct {
	MatchSummary struct {
		ShouldApply          models.ShouldApply `json:"should_apply"`
		ShouldApplyReasoning string             `json:"should_apply_reasoning"`
		OverallMatchSummary  struct {
			OverallMatchScore int      `json:"overall_match_score"`
			Suggestions       []string `json:"suggestions"`
			Summary           []struct {
				SummaryText        string             `json:"summary_text"`
				SummaryTemperature models.Temperature `json:"summary_temperature"`
			} `json:"summary"`
		} `json:"overall_match_summary"`
		Metrics []struct {
			ScoreTitle    string  `json:"score_title"`
			RawScore      float64 `json:"raw_score"`
			WeightedScore float64 `json:"weighted_score"`
			ScoreWeight   float64 `json:"score_weight"`
			ScoreReason   string  `json:"score_reason"`
			IsCompatible  bool    `json:"isCompatible"`
			Strength      string  `json:"strength"`
			Weaknesses    string  `json:"weaknesses"`
		} `json:"metrics"`
	} `json:"match_summary"`
}
