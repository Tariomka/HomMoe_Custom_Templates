package previewLayoutCache_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenCacheIsCreated_ReturnsNonNilCache(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	cache := models.NewPreviewLayoutCache()

	// Assert
	assert.NotNil(t, cache)
}
