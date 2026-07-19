package snapshotStore_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common/snapshot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenCalled_WritesFailureFile(t *testing.T) {
	t.Parallel()
	// Arrange
	store := snapshot.NewSnapshotStoreWithRoot(t.TempDir())
	failurePath := store.FailurePath("someFile", "SomeTest", 3)

	// Act
	err := store.SaveFailure(failurePath, sampleScreenshot())

	// Assert
	require.NoError(t, err)
	assert.FileExists(t, failurePath)
}
