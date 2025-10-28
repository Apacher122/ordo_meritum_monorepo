package schemas

import (
	"fmt"

	"github.com/ordo_meritum/features/documents/models/requests"
	"google.golang.org/genai"
)

// BuildResumeSchema dynamically constructs a genai.Schema for resume generation.
// It takes a slice of experience payloads to populate the 'enum' fields for
// 'position' and 'company', ensuring the LLM's output aligns with the provided experiences.
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
					Required: []string{"sentence", "is_new_suggestion", "justification_for_change"},
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
					Required: []string{"category", "justification_for_changes", "skill"},
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
								Required: []string{"text", "is_new_suggestion", "justification_for_change"},
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
								Required: []string{"text", "is_new_suggestion", "justification_for_change"},
							},
						},
					},
				},
			},
		},
		Required: []string{"summary", "skills", "experiences", "projects"},
	}, nil
}

// GeminiResumeSchema defines a static genai.Schema for resume generation
// used with the Gemini API.
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

// GeminiCoverLetterSchema defines the genai.Schema for cover letter generation
// used with the Gemini API.
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

// GeminiSummarySchema defines the genai.Schema for generating only the summary
// section of a resume with the Gemini API.
var GeminiSummarySchema = &genai.Schema{
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
				Required: []string{"sentence", "is_new_suggestion", "justification_for_change"},
			},
		},
	},
}

// BuildExperienceSchema dynamically constructs a genai.Schema for generating
// only the experiences section of a resume. It populates 'enum' fields for
// 'position' and 'company' based on the input data.
func BuildExperienceSchema(data any) (any, error) {
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
								Required: []string{"text", "is_new_suggestion", "justification_for_change"},
							},
						},
						"relevantSkillsMentioned": {
							Type: genai.TypeArray,
							Items: &genai.Schema{
								Type: genai.TypeString,
							},
						},
					},
					Required: []string{"position", "company", "start", "end", "bulletPoints", "relevantSkillsMentioned"},
				},
			},
		},
	}, nil

}

// GeminiProjectSchema defines the genai.Schema for generating only the projects
// section of a resume with the Gemini API.
var GeminiProjectSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"projects": {
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
					Required: []string{"text", "is_new_suggestion", "justification_for_change"},
				},
				"relevantSkillsMentioned": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeString,
					},
				},
			},
			Required: []string{"name", "role", "status", "bulletPoints", "relevantSkillsMentioned"},
		},
	},
}

// GeminiSkillSchema defines the genai.Schema for generating only the skills
// section of a resume with the Gemini API.
var GeminiSkillSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"skills": {
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"category":          {Type: genai.TypeString},
					"is_new_suggestion": {Type: genai.TypeBoolean},
					"skill": {
						Type:  genai.TypeArray,
						Items: &genai.Schema{Type: genai.TypeString},
					},
				},
				Required: []string{"category", "is_new_suggestion", "skill"},
			},
		},
	},
}
