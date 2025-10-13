package schemas

import (
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
)

var schemaFS embed.FS

func LoadSchema(featureName, schemaName string) (map[string]any, error) {
	path := filepath.ToSlash(fmt.Sprintf("../../features/%s/models/schemas/%s", featureName, schemaName))
	bytes, err := schemaFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema %s for feature %s: %w", schemaName, featureName, err)
	}

	var schema map[string]interface{}
	if err := json.Unmarshal(bytes, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse schema %s for feature %s: %w", schemaName, featureName, err)
	}

	return schema, nil
}
