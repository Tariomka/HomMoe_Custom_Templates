package wire_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/stretchr/testify/assert"
)

func TestWhenRegenerationInjectorIsInvoked_ReturnsAssembledRegenerationHandler(t *testing.T) {
	t.Parallel()

	// Arrange

	// Act
	regenerationHandler := composition.InitializeRegenerationHandler()

	// Assert
	assert.NotNil(t, regenerationHandler)
}
