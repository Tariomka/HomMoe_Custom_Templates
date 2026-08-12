package construction_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const entitiesImportPath = "github.com/Tariomka/hommoe_custom_templates/internal/entities"

func TestWhenProductionEntityConstructionIsScanned_UsesBuildersForInvariantRichTypes(t *testing.T) {
	t.Parallel()
	// Arrange
	repositoryRoot := getRepositoryRoot(t)
	builderOwnedTypes := map[string]bool{
		"Border": true, "Connection": true, "MainObject": true, "MandatoryContentItem": true,
		"Orientation": true, "PlacementRule": true, "Road": true, "TypedRef": true,
		"Variant": true, "Zone": true,
	}

	// Act
	violations := findDirectEntityLiterals(t, repositoryRoot, builderOwnedTypes)

	// Assert
	assert.Empty(t, violations)
}

func TestWhenSliceElementOmitsEntityType_DetectsBuilderOwnedLiteral(t *testing.T) {
	t.Parallel()
	// Arrange
	fileSet := token.NewFileSet()
	expression, err := parser.ParseExprFrom(fileSet, "fixture.go", `[]entities.Zone{{Name: "A"}}`, 0)
	require.NoError(t, err)
	literal := expression.(*ast.CompositeLit)

	// Act
	violations := findElidedEntityLiterals(
		fileSet, "fixture.go", literal, "entities", map[string]bool{"Zone": true})

	// Assert
	assert.Equal(t, []string{"fixture.go:1 Zone"}, violations)
}

func TestWhenMapValueOmitsEntityType_DetectsBuilderOwnedLiteral(t *testing.T) {
	t.Parallel()
	// Arrange
	fileSet := token.NewFileSet()
	expression, err := parser.ParseExprFrom(
		fileSet, "fixture.go", `map[string]entities.Zone{"A": {Name: "A"}}`, 0)
	require.NoError(t, err)
	literal := expression.(*ast.CompositeLit)

	// Act
	violations := findElidedEntityLiterals(
		fileSet, "fixture.go", literal, "entities", map[string]bool{"Zone": true})

	// Assert
	assert.Equal(t, []string{"fixture.go:1 Zone"}, violations)
}

func TestWhenPointerElementOmitsEntityType_DetectsBuilderOwnedLiteral(t *testing.T) {
	t.Parallel()
	// Arrange
	fileSet := token.NewFileSet()
	expression, err := parser.ParseExprFrom(fileSet, "fixture.go", `[]*entities.Zone{{Name: "A"}}`, 0)
	require.NoError(t, err)
	literal := expression.(*ast.CompositeLit)

	// Act
	violations := findElidedEntityLiterals(
		fileSet, "fixture.go", literal, "entities", map[string]bool{"Zone": true})

	// Assert
	assert.Equal(t, []string{"fixture.go:1 Zone"}, violations)
}

func TestWhenNestedElementOmitsEntityType_DetectsBuilderOwnedLiteral(t *testing.T) {
	t.Parallel()
	// Arrange
	fileSet := token.NewFileSet()
	expression, err := parser.ParseExprFrom(fileSet, "fixture.go", `[][]entities.Zone{{{Name: "A"}}}`, 0)
	require.NoError(t, err)
	literal := expression.(*ast.CompositeLit)

	// Act
	violations := findElidedEntityLiterals(
		fileSet, "fixture.go", literal, "entities", map[string]bool{"Zone": true})

	// Assert
	assert.Equal(t, []string{"fixture.go:1 Zone"}, violations)
}

func getRepositoryRoot(t *testing.T) string {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "..", ".."))
}

func findDirectEntityLiterals(
	t *testing.T,
	repositoryRoot string,
	builderOwnedTypes map[string]bool,
) []string {
	t.Helper()
	var violations []string
	for _, sourceDirectory := range []string{"app", "internal"} {
		sourceRoot := filepath.Join(repositoryRoot, sourceDirectory)
		err := filepath.Walk(sourceRoot, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if shouldSkipPath(repositoryRoot, path, info) {
				return nil
			}
			fileViolations, parseErr := findFileEntityLiterals(repositoryRoot, path, builderOwnedTypes)
			if parseErr != nil {
				return parseErr
			}
			violations = append(violations, fileViolations...)
			return nil
		})
		require.NoError(t, err)
	}
	return violations
}

func shouldSkipPath(repositoryRoot, path string, info os.FileInfo) bool {
	if info.IsDir() || filepath.Ext(path) != ".go" {
		return true
	}
	fileName := filepath.Base(path)
	if strings.HasSuffix(fileName, "_test.go") || strings.HasSuffix(fileName, "_testexports.go") {
		return true
	}
	relativePath, err := filepath.Rel(repositoryRoot, path)
	if err != nil {
		return false
	}
	normalizedPath := filepath.ToSlash(relativePath)
	return strings.HasPrefix(normalizedPath, "internal/entities/") ||
		strings.HasPrefix(normalizedPath, "internal/services/builders/")
}

func findFileEntityLiterals(
	repositoryRoot string,
	path string,
	builderOwnedTypes map[string]bool,
) ([]string, error) {
	fileSet := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fileSet, path, nil, 0)
	if err != nil {
		return nil, err
	}
	entitiesAlias, err := findEntitiesAlias(parsedFile)
	if err != nil || entitiesAlias == "" {
		return nil, err
	}
	relativePath, err := filepath.Rel(repositoryRoot, path)
	if err != nil {
		return nil, err
	}

	var violations []string
	ast.Inspect(parsedFile, func(node ast.Node) bool {
		literal, ok := node.(*ast.CompositeLit)
		if !ok {
			return true
		}
		violations = append(violations, findElidedEntityLiterals(
			fileSet, filepath.ToSlash(relativePath), literal, entitiesAlias, builderOwnedTypes)...)
		selector, ok := literal.Type.(*ast.SelectorExpr)
		if !ok || !builderOwnedTypes[selector.Sel.Name] {
			return true
		}
		packageName, ok := selector.X.(*ast.Ident)
		if !ok || packageName.Name != entitiesAlias {
			return true
		}
		position := fileSet.Position(literal.Pos())
		violations = append(violations,
			filepath.ToSlash(relativePath)+":"+strconv.Itoa(position.Line)+" "+selector.Sel.Name)
		return true
	})
	return violations, nil
}

func findElidedEntityLiterals(
	fileSet *token.FileSet,
	relativePath string,
	literal *ast.CompositeLit,
	entitiesAlias string,
	builderOwnedTypes map[string]bool,
) []string {
	var elementType ast.Expr
	switch collectionType := literal.Type.(type) {
	case *ast.ArrayType:
		elementType = collectionType.Elt
	case *ast.MapType:
		elementType = collectionType.Value
	default:
		return nil
	}
	return findElidedElements(
		fileSet, relativePath, literal.Elts, elementType, entitiesAlias, builderOwnedTypes)
}

func findElidedElements(
	fileSet *token.FileSet,
	relativePath string,
	elements []ast.Expr,
	elementType ast.Expr,
	entitiesAlias string,
	builderOwnedTypes map[string]bool,
) []string {
	var violations []string
	for _, element := range elements {
		candidate := element
		if keyValue, isMapElement := element.(*ast.KeyValueExpr); isMapElement {
			candidate = keyValue.Value
		}
		elidedLiteral, ok := candidate.(*ast.CompositeLit)
		if !ok || elidedLiteral.Type != nil {
			continue
		}
		if entityType := builderOwnedEntityType(elementType, entitiesAlias, builderOwnedTypes); entityType != "" {
			position := fileSet.Position(elidedLiteral.Pos())
			violations = append(violations,
				relativePath+":"+strconv.Itoa(position.Line)+" "+entityType)
			continue
		}
		switch nestedType := elementType.(type) {
		case *ast.ArrayType:
			violations = append(violations, findElidedElements(
				fileSet, relativePath, elidedLiteral.Elts, nestedType.Elt, entitiesAlias, builderOwnedTypes)...)
		case *ast.MapType:
			violations = append(violations, findElidedElements(
				fileSet, relativePath, elidedLiteral.Elts, nestedType.Value, entitiesAlias, builderOwnedTypes)...)
		}
	}
	return violations
}

func builderOwnedEntityType(
	elementType ast.Expr,
	entitiesAlias string,
	builderOwnedTypes map[string]bool,
) string {
	if pointerType, ok := elementType.(*ast.StarExpr); ok {
		return builderOwnedEntityType(pointerType.X, entitiesAlias, builderOwnedTypes)
	}
	selector, ok := elementType.(*ast.SelectorExpr)
	if !ok || !builderOwnedTypes[selector.Sel.Name] {
		return ""
	}
	packageName, ok := selector.X.(*ast.Ident)
	if !ok || packageName.Name != entitiesAlias {
		return ""
	}
	return selector.Sel.Name
}

func findEntitiesAlias(parsedFile *ast.File) (string, error) {
	for _, importSpec := range parsedFile.Imports {
		importPath, err := strconv.Unquote(importSpec.Path.Value)
		if err != nil {
			return "", err
		}
		if importPath != entitiesImportPath {
			continue
		}
		if importSpec.Name != nil {
			return importSpec.Name.Name, nil
		}
		return "entities", nil
	}
	return "", nil
}
