package wire_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/stretchr/testify/assert"
)

func TestWhenInjectorIsInvoked_ReturnsAssembledGuiHandler(t *testing.T) {
	t.Parallel()

	// Arrange

	// Act
	guiHandler := composition.InitializeGuiHandler()

	// Assert
	assert.NotNil(t, guiHandler)
}
