package constants

import "path/filepath"

var (
	RootDir, _   = filepath.Abs(".")
	TemplatesDir = filepath.Join(RootDir, "shared", "templates")
	PromptsDir   = filepath.Join(TemplatesDir, "prompts")
	LatexDir     = filepath.Join(TemplatesDir, "latex")
)
