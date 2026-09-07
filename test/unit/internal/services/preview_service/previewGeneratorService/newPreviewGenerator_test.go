package previewGeneratorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenEmbeddedAssetsAreValid_ReturnsNoError(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	_, err := preview_service.NewPreviewGenerator(
		preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService()),
	)

	// Assert
	assert.NoError(t, err)
}

func TestWhenEmbeddedAssetsAreValid_ReturnsGenerator(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	generator, err := preview_service.NewPreviewGenerator(
		preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService()),
	)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, generator)
}
