package guardStrength_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/stretchr/testify/assert"
)

func TestWhenCalled_ReturnsBronzeGuardStrengthValues(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := common_connections.GuardStrength{
		Default:  15_000,
		Weakest:  3_000,
		Low:      6_000,
		Medium:   9_000,
		High:     12_000,
		VeryHigh: 16_000,
	}

	// Act
	strength := common_connections.GetBronzeGuardStrength()

	// Assert
	assert.Equal(t, expected, strength)
}
