package schemas

import (
	"google.golang.org/genai"
)

var GeminiResumeSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"company_name": {
			Type: genai.TypeString,
		},
		"match_summary": {
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"should_apply": {
					Type: genai.TypeString,
					Enum: []string{"Strong Yes", "Yes", "No", "Strong No", "Maybe"},
				},
				"should_apply_reasoning": {
					Type: genai.TypeString,
				},
				"metrics": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"score_title": {
								Type: genai.TypeString,
								Enum: []string{
									"Keyword & Phrases",
									"Experience Alignment",
									"Education & Credentials",
									"Skills & Competencies",
									"Achievements & Quantifiable Results",
									"Job-Specific Filters",
									"Cultural & Organizational Fit (Emerging Factor)",
								},
							},
							"raw_score":      {Type: genai.TypeNumber},
							"weighted_score": {Type: genai.TypeNumber},
							"score_weight":   {Type: genai.TypeNumber},
							"score_reason":   {Type: genai.TypeString},
							"is_compatible":  {Type: genai.TypeBoolean},
							"strength":       {Type: genai.TypeString},
							"weaknesses":     {Type: genai.TypeString},
						},
					},
				},
				"overall_match_summary": {
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"overall_match_score": {
							Type: genai.TypeNumber,
						},
						"summary": {
							Type: genai.TypeArray,
							Items: &genai.Schema{
								Type: genai.TypeObject,
								Properties: map[string]*genai.Schema{
									"summary_text": {
										Type: genai.TypeString,
									},
									"summary_temperature": {
										Type: genai.TypeString,
										Enum: []string{"Good", "Neutral", "Bad"},
									},
								},
							},
						},
						"suggestions": {
							Type:  genai.TypeArray,
							Items: &genai.Schema{Type: genai.TypeString},
						},
					},
				},
			},
		},
	},
}
