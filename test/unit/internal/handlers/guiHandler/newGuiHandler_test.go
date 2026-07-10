package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/stretchr/testify/assert"
)

func TestWhenHandlerIsConstructed_ReturnsNonNilHandler(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	handler := handlers.NewGuiHandler()

	// Assert
	assert.NotNil(t, handler)
}
