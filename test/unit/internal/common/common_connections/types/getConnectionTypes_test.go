package types_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/stretchr/testify/assert"
)

func TestWhenCalled_ReturnsDirectAndPortalOnly(t *testing.T) {
	t.Parallel()
	// Arrange
	connectionTypeValues := registry.GetConnectionTypeValues()
	expected := []string{connectionTypeValues.Direct, connectionTypeValues.Portal}

	// Act
	connectionTypes := common_connections.GetConnectionTypes()

	// Assert
	assert.Equal(t, expected, connectionTypes)
}
