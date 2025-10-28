package schemaregistry

import (
	"fmt"

	app_schemas "github.com/ordo_meritum/features/application_tracking/models/schemas"
	doc_schemas "github.com/ordo_meritum/features/documents/models/schemas"
)

type SchemaBuilderFunc func(data any) (any, error)

// staticSchemaBuilder is a helper that wraps a static schema object into a
// SchemaBuilderFunc, ignoring any provided data.
func staticSchemaBuilder(schema any) SchemaBuilderFunc {
	return func(data any) (any, error) {
		return schema, nil
	}
}

var (
	Resume              = "resume"
	Coverletter         = "coverletter"
	MatchSummary        = "match_summary"
	ApplicationTracking = "application_tracking"
)

// ProviderSchemaRegistry maps LLM provider names to their supported schemas.
// Each schema is associated with a builder function that can construct it,
// allowing for both static and dynamic schema generation.
var ProviderSchemaRegistry = map[string]map[string]SchemaBuilderFunc{
	"gemini": {
		"resume":      doc_schemas.BuildResumeSchema,
		"coverletter": staticSchemaBuilder(doc_schemas.GeminiCoverLetterSchema),
	},
	"cohere": {
		"resume":               staticSchemaBuilder(doc_schemas.CohereResumeSchema),
		"coverletter":          staticSchemaBuilder(doc_schemas.CohereCoverLetterSchemaFormat),
		"application_tracking": staticSchemaBuilder(app_schemas.CohereJobDescriptionSchemaFormat),
	},
}

// GetSchema retrieves a schema from the registry for a specific provider and schema name.
// It uses the registered SchemaBuilderFunc to construct the schema, passing along any
// necessary data for dynamic schema creation. It returns an error if the provider
// or schema name is not found.
func GetSchema(provider, schemaName string, data any) (any, error) {
	providerSchemas, ok := ProviderSchemaRegistry[provider]
	if !ok {
		return nil, fmt.Errorf("provider '%s' not found in schema registry", provider)
	}

	builder, ok := providerSchemas[schemaName]
	if !ok {
		return nil, fmt.Errorf("schema '%s' not found for provider '%s'", schemaName, provider)
	}

	return builder(data)
}
