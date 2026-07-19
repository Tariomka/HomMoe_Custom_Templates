package store_test

import (
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common/snapshot"
	"github.com/stretchr/testify/assert"
)

func TestWhenPathRequested_UsesFailureExtension(t *testing.T) {
	t.Parallel()
	// Arrange
	root := t.TempDir()
	store := snapshot.NewStoreWithRoot(root)

	// Act
	failurePath := store.FailurePath("someFile", "SomeTest", 2)

	// Assert
	assert.Equal(t, filepath.Join(root, "someFile", "SomeTest_2.failure"), failurePath)
}
