package guardStrength_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenCalled_ReturnsHubGuardStrengthValues(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := models.GuardStrength{
		Default:  35_000,
		Weakest:  45_000,
		Low:      52_000,
		Medium:   62_000,
		High:     70_000,
		VeryHigh: 75_000,
	}

	// Act
	strength := common_connections.GetHubGuardStrength()

	// Assert
	assert.Equal(t, expected, strength)
}
