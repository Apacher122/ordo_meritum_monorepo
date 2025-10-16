package schemaregistry

import (
	"fmt"

	app_schemas "github.com/ordo_meritum/features/application_tracking/models/schemas"
	doc_schemas "github.com/ordo_meritum/features/documents/models/schemas"
)

var (
	Resume              = "resume"
	Coverletter         = "coverletter"
	MatchSummary        = "match_summary"
	ApplicationTracking = "application_tracking"
)

var ProviderSchemaRegistry = map[string]map[string]any{
	"gemini": {
		"resume":      doc_schemas.GeminiResumeSchema,
		"coverletter": doc_schemas.GeminiCoverLetterSchema,
	},
	"cohere": {
		"resume":               doc_schemas.CohereResumeSchema,
		"coverletter":          doc_schemas.CohereCoverLetterSchemaFormat,
		"application_tracking": app_schemas.CohereJobDescriptionSchemaFormat,
	},
}

func GetSchema(provider, schemaName string) (any, error) {
	providerSchemas, ok := ProviderSchemaRegistry[provider]
	if !ok {
		return nil, fmt.Errorf("provider '%s' not found in schema registry", provider)
	}

	schema, ok := providerSchemas[schemaName]
	if !ok {
		return nil, fmt.Errorf("document type '%s' not found for provider '%s'", schemaName, provider)
	}

	return schema, nil
}
