package dependency_test

import (
	"maps"
	"path"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// productionRoots are the directories the layering rules police. Everything
// under test/ is exempt: a test legitimately names whatever it asserts on.
//
//nolint:gochecknoglobals // shared, read-only rule input for this file's tests.
var productionRoots = []string{"app", "internal", "cmd"}

// entityNamerPrefixes are the packages allowed to name an entity type
// (plan §0.5.4). internal/helpers is permitted only through its per-domain
// *_helpers subpackages, which is checked separately.
//
//nolint:gochecknoglobals // shared, read-only rule input for this file's tests.
var entityNamerPrefixes = []string{
	"internal/repositories/",
	"internal/models/",
	"internal/entities/",
	"internal/mappers/",
}

// dtoNamerPrefixes are the packages allowed to name a DTO: the API boundary
// itself, the DTO package, and the consumer that crosses it (plan §0.4).
//
//nolint:gochecknoglobals // shared, read-only rule input for this file's tests.
var dtoNamerPrefixes = []string{
	"internal/handlers/",
	"internal/dtos/",
	"app/",
}

// entityNamerAllowList records the packages that named an entity before the
// rule existed. Base internal/entities is the .rmg.json vocabulary the whole
// generator is built out of, so this list is large by design and shrinks one
// package at a time. **Only ever remove entries.** The residual breach is
// tracked in .agent/backlog/backlog-opus5.md.
//
//nolint:gochecknoglobals // shared, read-only rule input for this file's tests.
var entityNamerAllowList = []string{
	"app/gui/dialogs",
	"app/gui/drivers",
	"app/gui/editor",
	"app/gui/models",
	"internal/dtos",
	"internal/handlers",
	"internal/handlers/handler_interfaces",
	"internal/services/builders/mandatory_content",
	"internal/services/builders/placement_rule",
	"internal/services/builders/variant_content",
	"internal/services/connection_editor",
	"internal/services/content_rules",
	"internal/services/file_service",
	"internal/services/preview_service",
	"internal/services/template_generator",
	"internal/services/template_generator/providers",
	"internal/services/template_generator/providers/provider_interfaces",
	"internal/services/template_generator/providers/topology",
	"internal/services/template_generator/providers/topology/base",
	"internal/services/template_generator/providers/topology/topology_interfaces",
	"internal/services/template_generator/providers/topology/tournament_variant",
	"internal/services/zones",
	"internal/services/zones/zone_interfaces",
}

// dtoNamerAllowList records the services that consumed DTOs before the rule
// existed. Each needs the treatment internal/services/editor got in Phase 10:
// a model-side request/result pair with the handler mapping onto it.
// **Only ever remove entries.**
//
//nolint:gochecknoglobals // shared, read-only rule input for this file's tests.
var dtoNamerAllowList = []string{
	"internal/services/bonuses",
	"internal/services/pickers",
	"internal/services/zone_content",
}

func TestWhenEntityImportsAreScanned_DoesNotDependOnHigherLayers(t *testing.T) {
	t.Parallel()
	// Arrange
	repositoryRoot := getRepositoryRoot(t)

	// Act
	violations := findImports(t, repositoryRoot, "internal/entities", isForbiddenAboveTheEntityLayer)

	// Assert
	assert.Empty(t, violations)
}

func TestWhenEntityConsumersAreScanned_OnlyPermittedPackagesNameAnEntity(t *testing.T) {
	t.Parallel()
	// Arrange
	repositoryRoot := getRepositoryRoot(t)

	// Act
	violations := findUnlistedNamers(t, repositoryRoot, isEntityImport, isPermittedEntityNamer, entityNamerAllowList)

	// Assert
	assert.Empty(t, violations)
}

func TestWhenDtoConsumersAreScanned_OnlyTheApiBoundaryAndAppNameADto(t *testing.T) {
	t.Parallel()
	// Arrange
	repositoryRoot := getRepositoryRoot(t)

	// Act
	violations := findUnlistedNamers(t, repositoryRoot, isDtoImport, isPermittedDtoNamer, dtoNamerAllowList)

	// Assert
	assert.Empty(t, violations)
}

func TestWhenAnEntityWouldImportAModel_TheRuleReportsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	forbidden := modulePath + "/internal/models/editor_state_model"

	// Act
	reported := isForbiddenAboveTheEntityLayer(forbidden)

	// Assert
	assert.True(t, reported)
}

func TestWhenAnEntityImportsAnotherEntity_TheRuleStaysSilent(t *testing.T) {
	t.Parallel()
	// Arrange
	permitted := modulePath + "/internal/entities/topology"

	// Act
	reported := isForbiddenAboveTheEntityLayer(permitted)

	// Assert
	assert.False(t, reported)
}

func TestWhenAnEntityImportsTheGenericDataHelpers_TheRuleStaysSilent(t *testing.T) {
	t.Parallel()
	// Arrange
	carvedOut := modulePath + "/internal/helpers/data"

	// Act
	reported := isForbiddenAboveTheEntityLayer(carvedOut)

	// Assert
	assert.False(t, reported)
}

func TestWhenAnEntityWouldImportADomainHelper_TheRuleReportsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	forbidden := modulePath + "/internal/helpers/editor_state_helpers"

	// Act
	reported := isForbiddenAboveTheEntityLayer(forbidden)

	// Assert
	assert.True(t, reported)
}

func TestWhenAnUnlistedPackageWouldNameAnEntity_TheRuleReportsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	unlisted := "internal/services/brand_new/service.go"

	// Act
	permitted := isPermittedEntityNamer(unlisted) || isAllowListed(unlisted, entityNamerAllowList)

	// Assert
	assert.False(t, permitted)
}

func TestWhenADomainHelperNamesAnEntity_TheRuleStaysSilent(t *testing.T) {
	t.Parallel()
	// Arrange
	domainHelper := "internal/helpers/editor_state_helpers/zoneContentRow.go"

	// Act
	permitted := isPermittedEntityNamer(domainHelper)

	// Assert
	assert.True(t, permitted)
}

func TestWhenAHelperOutsideADomainPackageNamesAnEntity_TheRuleReportsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	looseHelper := "internal/helpers/io.go"

	// Act
	permitted := isPermittedEntityNamer(looseHelper)

	// Assert
	assert.False(t, permitted)
}

func TestWhenAnUnlistedServiceWouldNameADto_TheRuleReportsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	unlisted := "internal/services/brand_new/service.go"

	// Act
	permitted := isPermittedDtoNamer(unlisted) || isAllowListed(unlisted, dtoNamerAllowList)

	// Assert
	assert.False(t, permitted)
}

// findUnlistedNamers collects every production file that imports the guarded
// layer from a package that is neither permitted nor on the allow-list.
func findUnlistedNamers(
	t *testing.T,
	repositoryRoot string,
	isGuardedImport func(string) bool,
	isPermittedNamer func(string) bool,
	allowList []string,
) map[string][]string {
	t.Helper()
	violations := map[string][]string{}
	for _, root := range productionRoots {
		found := findImports(t, repositoryRoot, root, isGuardedImport)
		maps.DeleteFunc(found, func(filePath string, _ []string) bool {
			return isPermittedNamer(filePath) || isAllowListed(filePath, allowList)
		})
		maps.Copy(violations, found)
	}

	return violations
}

func isForbiddenAboveTheEntityLayer(importPath string) bool {
	// internal/helpers/data holds generic data structures (Vec2, Tuple,
	// Adjacency), so §4.4 carves it out of the helpers ban.
	if hasPackagePrefix(importPath, modulePath+"/internal/helpers/data") {
		return false
	}

	forbidden := []string{"models", "dtos", "services", "handlers", "helpers"}
	return slices.ContainsFunc(forbidden, func(layer string) bool {
		return hasPackagePrefix(importPath, modulePath+"/internal/"+layer)
	})
}

func isEntityImport(importPath string) bool {
	return hasPackagePrefix(importPath, modulePath+"/internal/entities")
}

func isDtoImport(importPath string) bool {
	return hasPackagePrefix(importPath, modulePath+"/internal/dtos")
}

func isPermittedEntityNamer(filePath string) bool {
	if isDomainHelper(filePath) {
		return true
	}

	return slices.ContainsFunc(entityNamerPrefixes, func(prefix string) bool {
		return strings.HasPrefix(filePath, prefix)
	})
}

func isPermittedDtoNamer(filePath string) bool {
	return slices.ContainsFunc(dtoNamerPrefixes, func(prefix string) bool {
		return strings.HasPrefix(filePath, prefix)
	})
}

// isDomainHelper reports whether the file sits in one of internal/helpers'
// per-domain *_helpers subpackages, the only part of helpers §0.5.4 permits to
// name an entity.
func isDomainHelper(filePath string) bool {
	segments := strings.Split(filePath, "/")
	return len(segments) > 3 &&
		segments[0] == "internal" &&
		segments[1] == "helpers" &&
		strings.HasSuffix(segments[2], "_helpers")
}

func isAllowListed(filePath string, allowList []string) bool {
	return slices.Contains(allowList, path.Dir(filePath))
}

// hasPackagePrefix matches a package path exactly or as a parent of it, so
// "internal/models" never matches a sibling such as "internal/modelsomething".
func hasPackagePrefix(importPath, packagePath string) bool {
	return importPath == packagePath || strings.HasPrefix(importPath, packagePath+"/")
}
