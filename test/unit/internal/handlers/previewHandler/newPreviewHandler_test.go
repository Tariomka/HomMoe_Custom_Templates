package previewHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenHandlerIsCreated_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	handler := handlers.NewPreviewHandler(&test_helpers.PreviewLayoutServiceMock{})

	// Assert
	assert.NotNil(t, handler)
}
