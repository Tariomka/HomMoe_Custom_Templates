package guardStrength_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/stretchr/testify/assert"
)

func TestWhenCalled_ReturnsSilverGuardStrengthValues(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := common_connections.GuardStrength{
		Default:  20_000,
		Weakest:  18_000,
		Low:      21_000,
		Medium:   24_000,
		High:     27_000,
		VeryHigh: 30_000,
	}

	// Act
	strength := common_connections.GetSilverGuardStrength()

	// Assert
	assert.Equal(t, expected, strength)
}
