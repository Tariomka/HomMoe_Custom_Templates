package roadPolicyService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenRoadFactoryIsProvided_ReturnsUsableService(t *testing.T) {
	t.Parallel()
	// Arrange
	roadFactory := zones.NewRoadFactory()

	// Act
	service := zones.NewRoadPolicyService(roadFactory)

	// Assert
	assert.NotNil(t, service)
}
