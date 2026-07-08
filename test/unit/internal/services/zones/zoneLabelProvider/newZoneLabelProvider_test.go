package zoneLabelProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenProviderIsConstructed_ReturnsNonNilProvider(t *testing.T) {
	// Arrange

	// Act
	provider := zones.NewZoneLabelProvider()

	// Assert
	assert.NotNil(t, provider)
}
