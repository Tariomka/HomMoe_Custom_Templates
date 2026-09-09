package file_system

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type PathResolutionService struct{}

func NewPathResolutionService() IPathResolutionService {
	return &PathResolutionService{}
}

func (this *PathResolutionService) ResolveStartDirectory(preferred string) string {
	preferred = strings.TrimSpace(preferred)
	if preferred != "" {
		if absolute, err := filepath.Abs(preferred); err == nil {
			preferred = absolute
		}

		preferred = filepath.Clean(preferred)
		for candidate := preferred; ; {
			if info, err := os.Stat(candidate); err == nil && info.IsDir() {
				return candidate
			}

			parent := filepath.Dir(candidate)
			if parent == candidate {
				break
			}

			candidate = parent
		}
	}

	if home, err := os.UserHomeDir(); err == nil {
		if info, err := os.Stat(home); err == nil && info.IsDir() {
			return home
		}
	}

	if workingDirectory, err := os.Getwd(); err == nil {
		return workingDirectory
	}

	return "."
}

func (this *PathResolutionService) ParentDirectory(current string) string {
	if current == "" {
		return ""
	}

	parent := filepath.Dir(current)
	if parent == current {
		if runtime.GOOS == windowsOS {
			return ""
		}

		return current
	}

	return parent
}

func (this *PathResolutionService) ResolveSaveTarget(directory, filename, requiredSuffix string) (string, bool) {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return "", false
	}

	filename = filepath.Base(filename)
	if filename == "." || filename == ".." || filename == string(os.PathSeparator) {
		return "", false
	}

	if helpers.IsReservedFilename(filename) {
		return "", false
	}

	if requiredSuffix != "" && !strings.HasSuffix(strings.ToLower(filename), strings.ToLower(requiredSuffix)) {
		filename += requiredSuffix
	}

	return filepath.Join(directory, filename), true
}

func (this *PathResolutionService) PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (this *PathResolutionService) DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func (this *PathResolutionService) FindGameTemplateDirectory() (string, error) {
	return helpers.FindOldenEraTemplatesDir(false)
}
