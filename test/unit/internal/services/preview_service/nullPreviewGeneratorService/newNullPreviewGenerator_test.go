package nullPreviewGeneratorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	"github.com/stretchr/testify/assert"
)

func TestWhenConstructed_ReturnsAGenerator(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	generator := preview_service.NewNullPreviewGenerator()

	// Assert
	assert.NotNil(t, generator)
}
