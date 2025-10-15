package embeds

import (
	"embed"
	"fmt"
)

// ReadFile reads and returns the content of a file from the given embed.FS.
func ReadFile(fs embed.FS, fileName string) (string, error) {
	bytes, err := fs.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded file '%s': %w", fileName, err)
	}
	return string(bytes), nil
}
