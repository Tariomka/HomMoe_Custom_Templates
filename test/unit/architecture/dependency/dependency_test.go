package dependency_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const modulePath = "github.com/Tariomka/hommoe_custom_templates"

func TestWhenInternalImportsAreScanned_DoesNotDependOnApp(t *testing.T) {
	t.Parallel()
	// Arrange
	repositoryRoot := getRepositoryRoot(t)

	// Act
	violations := findImportsWithPrefix(t, repositoryRoot, "internal", modulePath+"/app/")

	// Assert
	assert.Empty(t, violations)
}

func TestWhenAppImportsAreScanned_DoesNotDependOnForbiddenInternalPackages(t *testing.T) {
	t.Parallel()
	// Arrange
	repositoryRoot := getRepositoryRoot(t)

	// Act
	violations := findForbiddenAppImports(t, repositoryRoot)

	// Assert
	assert.Empty(t, violations)
}

func getRepositoryRoot(t *testing.T) string {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "..", ".."))
}

func findImportsWithPrefix(t *testing.T, repositoryRoot, sourceDirectory, forbiddenPrefix string) map[string][]string {
	t.Helper()
	return findImports(t, repositoryRoot, sourceDirectory, func(importPath string) bool {
		return strings.HasPrefix(importPath, forbiddenPrefix)
	})
}

func findForbiddenAppImports(t *testing.T, repositoryRoot string) map[string][]string {
	t.Helper()
	allowedRoots := []string{
		modulePath + "/internal/common",
		modulePath + "/internal/dtos",
		modulePath + "/internal/entities",
		modulePath + "/internal/handlers",
		modulePath + "/internal/helpers",
		modulePath + "/internal/models",
		modulePath + "/internal/registry",
	}
	return findImports(t, repositoryRoot, "app", func(importPath string) bool {
		if !strings.HasPrefix(importPath, modulePath+"/internal/") {
			return false
		}
		return !slices.ContainsFunc(allowedRoots, func(allowedRoot string) bool {
			return importPath == allowedRoot || strings.HasPrefix(importPath, allowedRoot+"/")
		})
	})
}

func findImports(
	t *testing.T,
	repositoryRoot string,
	sourceDirectory string,
	isForbidden func(string) bool,
) map[string][]string {
	t.Helper()
	violations := map[string][]string{}
	sourceRoot := filepath.Join(repositoryRoot, sourceDirectory)
	err := filepath.Walk(sourceRoot, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}
		parsedFile, parseErr := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if parseErr != nil {
			return parseErr
		}
		for _, importSpec := range parsedFile.Imports {
			importPath, unquoteErr := strconv.Unquote(importSpec.Path.Value)
			if unquoteErr != nil {
				return unquoteErr
			}
			if isForbidden(importPath) {
				relativePath, relativeErr := filepath.Rel(repositoryRoot, path)
				if relativeErr != nil {
					return relativeErr
				}
				normalizedPath := filepath.ToSlash(relativePath)
				violations[normalizedPath] = append(violations[normalizedPath], importPath)
			}
		}
		return nil
	})
	require.NoError(t, err)
	return violations
}
