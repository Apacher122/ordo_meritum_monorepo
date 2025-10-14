package schemas

import (
	cohere "github.com/cohere-ai/cohere-go/v2"
)

var JobDescriptionResponseFormat = cohere.JsonResponseFormatV2{
	JsonSchema: map[string]interface{}{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"type":    "object",
		"properties": map[string]interface{}{
			"job_title":    map[string]interface{}{"type": "string"},
			"company_name": map[string]interface{}{"type": "string"},
			"years_of_exp": map[string]interface{}{
				"type":    "string",
				"default": "Not Specified",
			},
			"education_level": map[string]interface{}{
				"type":    "string",
				"default": "Not Specified",
			},
			"website": map[string]interface{}{"type": "string"},
			"applicant_count": map[string]interface{}{
				"type":    "integer",
				"default": 0,
			},
			"post_age": map[string]interface{}{
				"type":    "string",
				"default": "Not Specified",
			},
			"skills_required": map[string]interface{}{
				"type":  "array",
				"items": map[string]interface{}{"type": "string"},
			},
			"skills_nice_to_haves": map[string]interface{}{
				"type":  "array",
				"items": map[string]interface{}{"type": "string"},
			},
			"tools_and_technologies": map[string]interface{}{
				"type":  "array",
				"items": map[string]interface{}{"type": "string"},
			},
			"company_culture": map[string]interface{}{
				"type":    "string",
				"default": "",
			},
			"salary_range": map[string]interface{}{
				"type":    "string",
				"default": "Not Specified",
			},
		},
		"required": []string{"job_title", "company_name", "website", "post_age", "salary_range"},
	},
}
