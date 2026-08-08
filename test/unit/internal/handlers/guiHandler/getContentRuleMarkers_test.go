package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenARowHasNoRules_TheMarkerTextIsEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()

	// Act
	markers := handler.GetContentRuleMarkers(models.SidMapping{}, nil)

	// Assert
	assert.Empty(t, markers)
}
