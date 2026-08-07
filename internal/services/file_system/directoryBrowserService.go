package file_system

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	internal_constants "github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// DirectoryBrowserService lists and creates directories on the local machine.
// It is deliberately stateless: every call takes the directory it operates on,
// so the same instance serves every open dialog.
type DirectoryBrowserService struct{}

func NewDirectoryBrowserService() IDirectoryBrowserService {
	return &DirectoryBrowserService{}
}

// ListEntries reads directory and returns its contents with subdirectories
// first, each group sorted case-insensitively by name. Entries whose name
// starts with a dot - plus hidden/system entries on Windows - are omitted
// unless showHidden is true. When filterSuffixes is non-empty, files must end
// with one of them (case-insensitive) to be included; subdirectories are never
// filtered so navigation always stays possible.
func (this *DirectoryBrowserService) ListEntries(
	directory string,
	filterSuffixes []string,
	showHidden bool) ([]models.DirectoryEntry, error) {
	osEntries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	directories := make([]models.DirectoryEntry, 0, len(osEntries))
	files := make([]models.DirectoryEntry, 0, len(osEntries))
	for _, osEntry := range osEntries {
		name := osEntry.Name()
		if !showHidden && isHidden(name, osEntry) {
			continue
		}

		fullPath := filepath.Join(directory, name)
		if osEntry.IsDir() {
			directories = append(directories, models.DirectoryEntry{Name: name, Path: fullPath, IsDir: true})
			continue
		}

		if !matchesFilter(name, filterSuffixes) {
			continue
		}

		files = append(files, models.DirectoryEntry{Name: name, Path: fullPath, IsDir: false})
	}

	slices.SortStableFunc(directories, compareByLowercaseName)
	slices.SortStableFunc(files, compareByLowercaseName)

	return slices.Concat(directories, files), nil
}

// ListRoots enumerates the existing volume roots. On Windows that is every
// drive letter that answers a stat; other systems have a single root and get an
// empty result, which matches the caller's "there is nothing above / to ascend
// to" expectation. The platform is checked explicitly because `A:\` is a legal
// relative file name on Unix and would otherwise be mistaken for a volume.
func (this *DirectoryBrowserService) ListRoots() []models.DirectoryEntry {
	if runtime.GOOS != windowsOS {
		return nil
	}

	roots := make([]models.DirectoryEntry, 0, 26)
	for letter := 'A'; letter <= 'Z'; letter++ {
		root := string(letter) + `:\`
		if _, err := os.Stat(root); err == nil {
			roots = append(roots, models.DirectoryEntry{Name: root, Path: root, IsDir: true})
		}
	}

	return roots
}

// CreateDirectory creates name inside parent and returns the created path. The
// name is trimmed and rejected when empty, when it is a relative-path token or
// when it contains a separator or a volume colon, so a "new folder" prompt can
// never write outside parent. An empty parent is refused outright, because
// joining onto it would silently target the process working directory. Errors
// are sentinels from common_errors or the raw os error, leaving the wording of
// any message to the presentation layer.
func (this *DirectoryBrowserService) CreateDirectory(parent, name string) (string, error) {
	if parent == "" {
		return "", common_errors.ErrDirectoryParentEmpty
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return "", common_errors.ErrDirectoryNameEmpty
	}

	if name == "." || name == ".." ||
		strings.ContainsRune(name, '/') ||
		strings.ContainsRune(name, os.PathSeparator) ||
		strings.ContainsRune(name, ':') {
		return "", common_errors.ErrDirectoryNameInvalid
	}

	target := filepath.Join(parent, name)
	if err := os.Mkdir(target, internal_constants.FolderPermission); err != nil {
		return "", err
	}

	return target, nil
}

// isHidden reports whether an entry should be omitted while hidden entries are
// suppressed: dotfiles on every platform, plus hidden/system entries on
// Windows. An entry whose info cannot be read is kept, so an unreadable entry
// stays visible rather than silently disappearing.
func isHidden(name string, osEntry os.DirEntry) bool {
	if strings.HasPrefix(name, ".") {
		return true
	}

	info, err := osEntry.Info()
	if err != nil {
		return false
	}

	return hasHiddenAttribute(info)
}

// matchesFilter reports whether name ends with one of the suffixes. An empty
// suffix list matches everything.
func matchesFilter(name string, filterSuffixes []string) bool {
	if len(filterSuffixes) == 0 {
		return true
	}

	lowerName := strings.ToLower(name)
	for _, suffix := range filterSuffixes {
		if strings.HasSuffix(lowerName, strings.ToLower(suffix)) {
			return true
		}
	}

	return false
}

func compareByLowercaseName(first, second models.DirectoryEntry) int {
	return strings.Compare(strings.ToLower(first.Name), strings.ToLower(second.Name))
}
