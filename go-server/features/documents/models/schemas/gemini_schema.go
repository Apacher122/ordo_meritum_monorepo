package schemas

import (
	"fmt"

	"github.com/ordo_meritum/features/documents/models/requests"
	"google.golang.org/genai"
)

func BuildResumeSchema(data any) (any, error) {
	experiences, ok := data.([]requests.ExperiencePayload)
	if !ok {
		return nil, fmt.Errorf("invalid data type for BuildResumeSchema; expected []requests.ExperiencePayload")
	}
	var allPositions, allCompanies []string
	for _, exp := range experiences {
		allPositions = append(allPositions, exp.Position)
		allCompanies = append(allCompanies, exp.Company)
	}

	return &genai.Schema{
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
						"position": {Type: genai.TypeString, Enum: allPositions},
						"company":  {Type: genai.TypeString, Enum: allCompanies},
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
								Required: []string{"text"},
							},
						},
					},
					Required: []string{"position", "company", "start", "end"},
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
	}, nil
}

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
					"position": {Type: genai.TypeString, Description: "The job title of the listed experience. This value must appear EXACTLY as the user enters it in the original resume."},
					"company":  {Type: genai.TypeString, Description: "The company name of the listed experience. This value must appear EXACTLY as the user enters it in the original resume."},
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
							Required: []string{"text"},
						},
					},
				},
				Required: []string{"position", "company", "start", "end"},
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
