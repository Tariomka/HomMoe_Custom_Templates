package snapshotStore_test

import (
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common/snapshot"
	"github.com/stretchr/testify/assert"
)

func TestWhenConstructed_RootsUnderIntegrationCommonSnapshots(t *testing.T) {
	t.Parallel()
	// Arrange
	store := snapshot.NewSnapshotStore()

	// Act
	goldenPath := store.GoldenPath("someFile", "SomeTest", 1)

	// Assert
	assert.Contains(t, goldenPath, filepath.Join("test_helpers", "integration_common", "__snapshots__"))
}
