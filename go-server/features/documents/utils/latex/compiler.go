package latex

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CompileToPDF(texPath string) ([]byte, error) {
	workspaceDir := filepath.Dir(texPath)
	texFilename := filepath.Base(texPath)
	pdfFilename := strings.TrimSuffix(texFilename, ".tex") + ".pdf"
	pdfPath := filepath.Join(workspaceDir, pdfFilename)

	for i := 0; i < 2; i++ {
		cmd := exec.Command("pdflatex", "-interaction=nonstopmode", texFilename)
		cmd.Dir = workspaceDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("pdflatex failed on run %d. Output:\n%s\nError: %w", i+1, string(output), err)
		}
	}

	return os.ReadFile(pdfPath)
}

func CopyTemplateAssets(sourceDir, destDir string) error {
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		destPath := filepath.Join(destDir, strings.TrimPrefix(path, sourceDir))
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		if filepath.Ext(path) == ".tex" {
			return nil
		}
		sourceFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer sourceFile.Close()
		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()
		_, err = io.Copy(destFile, sourceFile)
		return err
	})
}
