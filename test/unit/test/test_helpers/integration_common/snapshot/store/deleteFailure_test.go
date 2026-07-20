package store_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common/snapshot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenFileExists_RemovesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	store := snapshot.NewStoreWithRoot(t.TempDir())
	failurePath := store.FailurePath("someFile", "SomeTest", 1)
	require.NoError(t, store.SaveFailure(failurePath, sampleScreenshot()))

	// Act
	err := store.DeleteFailure(failurePath)

	// Assert
	require.NoError(t, err)
	assert.NoFileExists(t, failurePath)
}

func TestWhenFileIsMissing_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	store := snapshot.NewStoreWithRoot(t.TempDir())

	// Act
	err := store.DeleteFailure(store.FailurePath("someFile", "NeverSaved", 1))

	// Assert
	assert.NoError(t, err)
}
