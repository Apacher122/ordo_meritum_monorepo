package formatters

import (
	"bytes"
	"embed"
	"log"
	"text/template"
)

func FormatTemplate(fs embed.FS, filename string, data any) (string, error) {
	content, err := fs.ReadFile(filename)

	if err != nil {
		return "", nil
	}

	empl, err := template.New(filename).Parse(string(content))
	if err != nil {
		log.Printf("could not parse template %s: %s", filename, err)
		return "", nil
	}

	var buf bytes.Buffer
	if err := empl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
