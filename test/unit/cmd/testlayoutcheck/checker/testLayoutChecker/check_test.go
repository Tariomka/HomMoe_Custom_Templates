package testLayoutChecker_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/cmd/testlayoutcheck/checker"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testExportsSource = "//go:build integration_test\n\npackage editor\n\nfunc (this *Window) TabCount() int { return 0 }\n"
	testExportsPath   = "app/gui/editor/window_testexports.go"
)

func writeTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for relativePath, content := range files {
		path := filepath.Join(root, filepath.FromSlash(relativePath))
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o750))
		require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	}

	return root
}

func reportedRules(violations []checker.Violation) []string {
	rules := make([]string, 0, len(violations))
	for _, violation := range violations {
		rules = append(rules, violation.Rule)
	}

	return rules
}

func TestWhenTreeFollowsTheRules_ReturnsNoViolations(t *testing.T) {
	t.Parallel()
	// Arrange
	root := writeTree(t, map[string]string{
		testExportsPath:            testExportsSource,
		"app/gui/editor/window.go": "package editor\n\ntype Window struct{}\n",
		"test/integration/window_integration_test.go": "//go:build integration_test\n\n" +
			"package integration_test\n\nfunc use(window *Window) int { return window.TabCount() }\n",
		"test/integration/gui/render_integration_test.go": "//go:build integration_test && gui\n\n" +
			"package gui_test\n",
		"test/unit/internal/services/service/do_test.go": "package service_test\n",
	})

	// Act
	violations, err := checker.NewTestLayoutChecker().Check(root)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, violations)
}

func TestWhenUnitTestCarriesIntegrationTag_ReportsTaggedUnitTest(t *testing.T) {
	t.Parallel()
	// Arrange
	root := writeTree(t, map[string]string{
		"test/unit/internal/services/service/do_test.go": "//go:build integration_test\n\npackage service_test\n",
	})

	// Act
	violations, err := checker.NewTestLayoutChecker().Check(root)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{checker.RuleTaggedUnitTest}, reportedRules(violations))
}

func TestWhenProductionFileCarriesIntegrationTag_ReportsTaggedProductionFile(t *testing.T) {
	t.Parallel()
	// Arrange
	root := writeTree(t, map[string]string{
		"internal/services/service.go": "//go:build integration_test\n\npackage services\n",
	})

	// Act
	violations, err := checker.NewTestLayoutChecker().Check(root)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{checker.RuleTaggedProductionFile}, reportedRules(violations))
}

func TestWhenTestExportsFileCarriesIntegrationTag_ReturnsNoViolations(t *testing.T) {
	t.Parallel()
	// Arrange
	root := writeTree(t, map[string]string{testExportsPath: testExportsSource})

	// Act
	violations, err := checker.NewTestLayoutChecker().Check(root)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, violations)
}

func TestWhenTestHelperOutsideProductionCarriesIntegrationTag_ReturnsNoViolations(t *testing.T) {
	t.Parallel()
	// Arrange
	root := writeTree(t, map[string]string{
		"test/test_helpers/integration_common/appRunner.go": "//go:build integration_test\n\n" +
			"package integration_common\n",
	})

	// Act
	violations, err := checker.NewTestLayoutChecker().Check(root)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, violations)
}

func TestWhenGuiIntegrationTestOmitsGuiTag_ReportsMissingGuiTag(t *testing.T) {
	t.Parallel()
	// Arrange
	root := writeTree(t, map[string]string{
		"test/integration/gui/render_integration_test.go": "//go:build integration_test\n\npackage gui_test\n",
	})

	// Act
	violations, err := checker.NewTestLayoutChecker().Check(root)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{checker.RuleMissingGuiTag}, reportedRules(violations))
}

func TestWhenTestUsesTestOnlyExportWithoutTag_ReportsMissingIntegrationTag(t *testing.T) {
	t.Parallel()
	// Arrange
	accessor := gofakeit.LetterN(8)
	root := writeTree(t, map[string]string{
		"app/gui/drivers/state_testexports.go": "//go:build integration_test\n\npackage drivers\n\n" +
			"func (this *State) " + accessor + "() {}\n",
		"test/integration/state_integration_test.go": "package integration_test\n\n" +
			"func use(state *State) { state." + accessor + "() }\n",
	})

	// Act
	violations, err := checker.NewTestLayoutChecker().Check(root)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{checker.RuleMissingIntegrationTag}, reportedRules(violations))
}

func TestWhenMissingIntegrationTagIsReported_NamesTheOffendingFile(t *testing.T) {
	t.Parallel()
	// Arrange
	root := writeTree(t, map[string]string{
		testExportsPath: testExportsSource,
		"test/integration/state_integration_test.go": "package integration_test\n\n" +
			"func use(window *Window) int { return window.TabCount() }\n",
	})

	// Act
	violations, err := checker.NewTestLayoutChecker().Check(root)

	// Assert
	require.Len(t, violations, 1)
	require.NoError(t, err)
	assert.Equal(t, "test/integration/state_integration_test.go", violations[0].Path)
}

func TestWhenFileIsGatedByAnUnrelatedNegatedTag_ReturnsNoViolations(t *testing.T) {
	t.Parallel()
	// Arrange
	root := writeTree(t, map[string]string{
		"internal/composition/wire_gen.go": "//go:build !wireinject\n\npackage composition\n",
		"internal/helpers/io_other.go":     "//go:build !windows\n\npackage helpers\n",
	})

	// Act
	violations, err := checker.NewTestLayoutChecker().Check(root)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, violations)
}

func TestWhenFileRequiresTheTagToBeAbsent_ReturnsNoViolations(t *testing.T) {
	t.Parallel()
	// Arrange
	root := writeTree(t, map[string]string{
		"internal/services/service.go": "//go:build !integration_test\n\npackage services\n",
	})

	// Act
	violations, err := checker.NewTestLayoutChecker().Check(root)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, violations)
}

func TestWhenDirectoryIsExcludedFromTheWalk_ReturnsNoViolations(t *testing.T) {
	t.Parallel()
	// Arrange
	root := writeTree(t, map[string]string{
		"data/generated.go":   "//go:build integration_test\n\npackage data\n",
		"output/generated.go": "//go:build integration_test\n\npackage output\n",
		"tmp/scratch.go":      "//go:build integration_test\n\npackage tmp\n",
		".agent/scratch.go":   "//go:build integration_test\n\npackage agent\n",
	})

	// Act
	violations, err := checker.NewTestLayoutChecker().Check(root)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, violations)
}

func TestWhenSourceCannotBeParsed_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	root := writeTree(t, map[string]string{"internal/services/broken.go": "package services\n\nfunc (\n"})

	// Act
	_, err := checker.NewTestLayoutChecker().Check(root)

	// Assert
	assert.Error(t, err)
}

func TestWhenRootDoesNotExist_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	root := filepath.Join(t.TempDir(), "missing")

	// Act
	_, err := checker.NewTestLayoutChecker().Check(root)

	// Assert
	assert.Error(t, err)
}
