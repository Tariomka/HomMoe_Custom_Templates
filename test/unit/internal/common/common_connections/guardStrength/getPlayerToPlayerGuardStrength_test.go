package guardStrength_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenCalled_ReturnsPlayerToPlayerGuardStrengthValues(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := models.GuardStrength{
		Default:  30_000,
		Weakest:  10_000,
		Low:      22_000,
		Medium:   34_000,
		High:     46_000,
		VeryHigh: 58_000,
	}

	// Act
	strength := common_connections.GetPlayerToPlayerGuardStrength()

	// Assert
	assert.Equal(t, expected, strength)
}
