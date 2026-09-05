package previewLayoutCache_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
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
		cache.Get(1, config.TopologyRandom, canvasSide, build)
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
	first := cache.Get(1, config.TopologyRandom, 600, build)

	// Act
	second := cache.Get(1, config.TopologyRandom, 600, build)

	// Assert
	assert.Equal(t, first, second)
}

func TestWhenTemplateRevisionChanges_RebuildsTheLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	callCount := 0
	build := newCountingBuild(&callCount)
	cache := models.NewPreviewLayoutCache()
	cache.Get(1, config.TopologyRandom, 600, build)

	// Act
	cache.Get(2, config.TopologyRandom, 600, build)

	// Assert
	assert.Equal(t, 2, callCount)
}

func TestWhenTopologyChanges_RebuildsTheLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	callCount := 0
	build := newCountingBuild(&callCount)
	cache := models.NewPreviewLayoutCache()
	cache.Get(1, config.TopologyRandom, 600, build)

	// Act
	cache.Get(1, config.TopologyRing, 600, build)

	// Assert
	assert.Equal(t, 2, callCount)
}

func TestWhenCanvasSideChanges_RebuildsTheLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	callCount := 0
	build := newCountingBuild(&callCount)
	cache := models.NewPreviewLayoutCache()
	cache.Get(1, config.TopologyRandom, 600, build)

	// Act
	cache.Get(1, config.TopologyRandom, 601, build)

	// Assert
	assert.Equal(t, 2, callCount)
}

func TestWhenKeyChanges_ReturnsTheRebuiltLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	callCount := 0
	build := newCountingBuild(&callCount)
	cache := models.NewPreviewLayoutCache()
	cache.Get(1, config.TopologyRandom, 600, build)

	// Act
	actual := cache.Get(2, config.TopologyRandom, 600, build)

	// Assert
	assert.Equal(t, preview.Layout{ZoneRadius: 2}, actual)
}

// newCountingBuild returns a build function that reports how often it ran and
// hands back a layout whose radius identifies the call that produced it.
func newCountingBuild(callCount *int) func() preview.Layout {
	return func() preview.Layout {
		*callCount++
		return preview.Layout{ZoneRadius: float64(*callCount)}
	}
}
