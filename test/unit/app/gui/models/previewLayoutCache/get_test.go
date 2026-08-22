package previewLayoutCache_test

import (
	"errors"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenKeyIsUnchanged_BuildsOnlyOnce(t *testing.T) {
	t.Parallel()
	// Arrange
	callCount := 0
	build := newCountingBuild(&callCount)
	canvasSide := float64(gofakeit.Number(100, 900))
	cache := models.NewPreviewLayoutCache()

	// Act
	for range 5 {
		_, err := cache.Get(1, config.TopologyRandom, canvasSide, build)
		require.NoError(t, err)
	}

	// Assert
	assert.Equal(t, 1, callCount)
}

func TestWhenKeyIsUnchanged_ReturnsTheCachedLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	callCount := 0
	build := newCountingBuild(&callCount)
	cache := models.NewPreviewLayoutCache()
	first, err := cache.Get(1, config.TopologyRandom, 600, build)
	require.NoError(t, err)

	// Act
	second, err := cache.Get(1, config.TopologyRandom, 600, build)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, first, second)
}

func TestWhenTemplateRevisionChanges_RebuildsTheLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	callCount := 0
	build := newCountingBuild(&callCount)
	cache := models.NewPreviewLayoutCache()
	_, err := cache.Get(1, config.TopologyRandom, 600, build)
	require.NoError(t, err)

	// Act
	_, err = cache.Get(2, config.TopologyRandom, 600, build)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 2, callCount)
}

func TestWhenTopologyChanges_RebuildsTheLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	callCount := 0
	build := newCountingBuild(&callCount)
	cache := models.NewPreviewLayoutCache()
	_, err := cache.Get(1, config.TopologyRandom, 600, build)
	require.NoError(t, err)

	// Act
	_, err = cache.Get(1, config.TopologyRing, 600, build)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 2, callCount)
}

func TestWhenCanvasSideChanges_RebuildsTheLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	callCount := 0
	build := newCountingBuild(&callCount)
	cache := models.NewPreviewLayoutCache()
	_, err := cache.Get(1, config.TopologyRandom, 600, build)
	require.NoError(t, err)

	// Act
	_, err = cache.Get(1, config.TopologyRandom, 601, build)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 2, callCount)
}

func TestWhenKeyChanges_ReturnsTheRebuiltLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	callCount := 0
	build := newCountingBuild(&callCount)
	cache := models.NewPreviewLayoutCache()
	_, err := cache.Get(1, config.TopologyRandom, 600, build)
	require.NoError(t, err)

	// Act
	actual, err := cache.Get(2, config.TopologyRandom, 600, build)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, preview.Layout{ZoneRadius: 2}, actual)
}

func TestWhenBuildFails_ReturnsTheError(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedError := errors.New(gofakeit.Sentence(3))
	cache := models.NewPreviewLayoutCache()

	// Act
	_, err := cache.Get(1, config.TopologyRandom, 600, func() (preview.Layout, error) {
		return preview.Layout{}, expectedError
	})

	// Assert
	assert.ErrorIs(t, err, expectedError)
}

func TestWhenBuildFails_IsNotCached(t *testing.T) {
	t.Parallel()
	// Arrange
	callCount := 0
	cache := models.NewPreviewLayoutCache()
	failingBuild := func() (preview.Layout, error) {
		callCount++
		return preview.Layout{}, errors.New(gofakeit.Sentence(3))
	}
	_, err := cache.Get(1, config.TopologyRandom, 600, failingBuild)
	require.Error(t, err)

	// Act
	_, err = cache.Get(1, config.TopologyRandom, 600, failingBuild)

	// Assert
	require.Error(t, err)
	assert.Equal(t, 2, callCount)
}

// newCountingBuild returns a build function that reports how often it ran and
// hands back a layout whose radius identifies the call that produced it.
func newCountingBuild(callCount *int) func() (preview.Layout, error) {
	return func() (preview.Layout, error) {
		*callCount++
		return preview.Layout{ZoneRadius: float64(*callCount)}, nil
	}
}
