package store_test

import (
	"path/filepath"
	"strconv"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common/snapshot"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenPathRequested_JoinsRootFileNameAndNumberedGolden(t *testing.T) {
	t.Parallel()
	// Arrange
	root := t.TempDir()
	actionNumber := gofakeit.Number(1, 9)
	store := snapshot.NewStoreWithRoot(root)

	// Act
	goldenPath := store.GoldenPath("window_snapshot_integration_test", "TestTabs_Snapshot", actionNumber)

	// Assert
	assert.Equal(
		t,
		filepath.Join(
			root,
			"window_snapshot_integration_test",
			"TestTabs_Snapshot_"+strconv.Itoa(actionNumber)+".golden",
		),
		goldenPath,
	)
}

func TestWhenTestNameHasSubtestSlash_SanitizesToUnderscore(t *testing.T) {
	t.Parallel()
	// Arrange
	store := snapshot.NewStoreWithRoot(t.TempDir())

	// Act
	goldenPath := store.GoldenPath("someFile", "TestParent/sub case", 1)

	// Assert
	assert.Equal(t, "TestParent_sub_case_1.golden", filepath.Base(goldenPath))
}
