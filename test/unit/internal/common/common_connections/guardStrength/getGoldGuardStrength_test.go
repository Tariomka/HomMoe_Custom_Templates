package guardStrength_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenCalled_ReturnsGoldGuardStrengthValues(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := models.GuardStrength{
		Default:  25_000,
		Weakest:  36_000,
		Low:      42_000,
		Medium:   48_000,
		High:     54_000,
		VeryHigh: 60_000,
	}

	// Act
	strength := common_connections.GetGoldGuardStrength()

	// Assert
	assert.Equal(t, expected, strength)
}
