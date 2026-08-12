package roadFactory_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenConstructed_ReturnsFactory(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	factory := zones.NewRoadFactory()

	// Assert
	assert.NotNil(t, factory)
}
