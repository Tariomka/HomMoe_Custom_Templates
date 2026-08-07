package file_system

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

// PathResolutionService normalizes requested locations and typed-in file names.
// Like DirectoryBrowserService it is stateless: everything it needs arrives as
// arguments.
type PathResolutionService struct{}

func NewPathResolutionService() IPathResolutionService {
	return &PathResolutionService{}
}

// ResolveStartDirectory normalizes preferred and guarantees an existing
// directory is returned: it makes the path absolute, then climbs to the nearest
// existing ancestor (which also handles being handed a file path), then falls
// back to the user's home directory and finally the working directory.
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

// ParentDirectory returns the directory to ascend to from current, or current
// itself when it is already the top (so an Up control can be disabled). The
// empty string means the synthetic volume-root listing, which is also what a
// Windows volume root ascends to.
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

// ResolveSaveTarget builds an absolute save path inside directory from a
// user-typed name. Any directory component is stripped so the target can never
// escape directory, and requiredSuffix is appended when missing. ok is false
// when the name is empty, resolves to nothing usable, or names a Windows
// device.
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

// PathExists reports whether anything - file or directory - is currently
// present at path.
func (this *PathResolutionService) PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirectoryExists reports whether path is currently a directory, so callers can
// tell a file they may overwrite from a folder they may not.
func (this *PathResolutionService) DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
