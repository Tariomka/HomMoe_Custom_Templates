package connectionEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenListingUserCreatableTypes_ReturnsDirectAndPortalOnly(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	connectionTypes := connection_editor.UserCreatableConnectionTypes()

	// Assert
	assert.Equal(t, []string{"Direct", "Portal"}, connectionTypes)
}
