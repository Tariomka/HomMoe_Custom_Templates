package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenARowHasNoRules_TheDisplayNameIsTheContentName(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	content := models.SidMapping{Name: "Crystal Mine"}

	// Act
	displayName := handler.GetContentRowDisplayName(content, nil)

	// Assert
	assert.Equal(t, "Crystal Mine", displayName)
}
