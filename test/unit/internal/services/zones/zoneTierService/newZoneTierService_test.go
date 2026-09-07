package zoneTierService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenServiceIsCreated_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	service := zones.NewZoneTierService()

	// Assert
	assert.NotNil(t, service)
}
