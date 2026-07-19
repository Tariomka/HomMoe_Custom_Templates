package store_test

import (
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common/snapshot"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenRootGiven_PathsStartWithRoot(t *testing.T) {
	t.Parallel()
	// Arrange
	root := filepath.Join(t.TempDir(), gofakeit.Word())
	store := snapshot.NewStoreWithRoot(root)

	// Act
	goldenPath := store.GoldenPath("someFile", "SomeTest", 1)

	// Assert
	assert.Equal(t, root, goldenPath[:len(root)])
}
