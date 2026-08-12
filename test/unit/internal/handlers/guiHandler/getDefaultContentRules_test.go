package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenDefaultContentRulesAreRequested_TheContentIsGuardedByDefault(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()

	// Act
	rules := handler.GetDefaultContentRules(models.SidMapping{})

	// Assert
	assert.True(t, *rules[0].IsGuarded)
}
