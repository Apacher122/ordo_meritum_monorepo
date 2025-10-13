package schemas

import (
	"google.golang.org/genai"
)

var GeminiResumeSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"summary": {
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"sentence":                 {Type: genai.TypeString},
					"justification_for_change": {Type: genai.TypeString},
					"is_new_suggestion":        {Type: genai.TypeBoolean},
				},
			},
		},
		"skills": {
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"category":                  {Type: genai.TypeString},
					"justification_for_changes": {Type: genai.TypeString},
					"skill": {
						Type:  genai.TypeArray,
						Items: &genai.Schema{Type: genai.TypeString},
					},
				},
			},
		},
		"experiences": {
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"position": {Type: genai.TypeString},
					"company":  {Type: genai.TypeString},
					"start":    {Type: genai.TypeString},
					"end":      {Type: genai.TypeString},
					"bulletPoints": {
						Type: genai.TypeArray,
						Items: &genai.Schema{
							Type: genai.TypeObject,
							Properties: map[string]*genai.Schema{
								"text":                     {Type: genai.TypeString},
								"is_new_suggestion":        {Type: genai.TypeBoolean},
								"justification_for_change": {Type: genai.TypeString},
							},
						},
					},
				},
			},
		},
		"projects": {
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"name":   {Type: genai.TypeString},
					"role":   {Type: genai.TypeString},
					"status": {Type: genai.TypeString},
					"bulletPoints": {
						Type: genai.TypeArray,
						Items: &genai.Schema{
							Type: genai.TypeObject,
							Properties: map[string]*genai.Schema{
								"text":                     {Type: genai.TypeString},
								"is_new_suggestion":        {Type: genai.TypeBoolean},
								"justification_for_change": {Type: genai.TypeString},
							},
						},
					},
				},
			},
		},
	},
	Required: []string{"summary", "skills", "experiences", "projects"},
}

var GeminiCoverLetterSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"about": {
			Type: genai.TypeString,
		},
		"experience": {
			Type: genai.TypeString,
		},
		"whatIBring": {
			Type: genai.TypeString,
		},
		"revisionSummary": {
			Type: genai.TypeString,
		},
	},
	Required: []string{"about", "experience", "whatIBring", "revisionSummary"},
}
