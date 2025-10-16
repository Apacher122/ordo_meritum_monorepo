package schemas

import (
	cohere "github.com/cohere-ai/cohere-go/v2"
)

var CohereResumeSchema = map[string]any{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type":    "object",
	"properties": map[string]any{
		"summary": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"sentence": map[string]any{
						"type": "string",
					},
					"justification_for_change": map[string]any{
						"type": "string",
					},
					"is_new_suggestion": map[string]any{
						"type": "boolean",
					},
				},
			},
			"skills": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"category": map[string]any{
							"type": "string",
						},
						"justification_for_changes": map[string]any{
							"type": "string",
						},
						"skill": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "string",
							},
						},
					},
				},
			},
			"experiences": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"position": map[string]any{
							"type": "string",
						},
						"company": map[string]any{
							"type": "string",
						},
						"start": map[string]any{
							"type": "string",
						},
						"end": map[string]any{
							"type": "string",
						},
						"bulletPoints": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"text": map[string]any{
										"type": "string",
									},
									"is_new_suggestion": map[string]any{
										"type": "boolean",
									},
									"justification_for_change": map[string]any{
										"type": "string",
									},
								},
							},
						},
					},
				},
				"projects": map[string]any{
					"type": "array",
					"items": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"name": map[string]any{
								"type": "string",
							},
							"role": map[string]any{
								"type": "string",
							},
							"status": map[string]any{
								"type": "string",
							},
							"bulletPoints": map[string]any{
								"type": "array",
								"items": map[string]any{
									"type": "object",
									"properties": map[string]any{
										"text": map[string]any{
											"type": "string",
										},
										"is_new_suggestion": map[string]any{
											"type": "boolean",
										},
										"justification_for_change": map[string]any{
											"type": "string",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
	"required": []string{"resume", "summary", "experiences", "projects", "skills"},
}

var CohereResumeSchemaFormat = cohere.JsonResponseFormatV2{
	JsonSchema: CohereResumeSchema,
}

var CoverLetterSchema = map[string]any{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type":    "object",
	"properties": map[string]any{
		"about":           map[string]any{"type": "string"},
		"experience":      map[string]any{"type": "string"},
		"whatIBring":      map[string]any{"type": "string"},
		"revisionSummary": map[string]any{"type": "string"},
	},
	"required": []string{"about", "experience", "whatIBring", "revisionSummary"},
}

var CohereCoverLetterSchemaFormat = cohere.JsonResponseFormatV2{
	JsonSchema: CoverLetterSchema,
}
