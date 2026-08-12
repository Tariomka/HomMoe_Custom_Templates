package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenContentItemsAreSorted_TheyAreOrderedByName(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	items := []models.SidMapping{{Name: "beta"}, {Name: "Alpha"}}

	// Act
	sorted := handler.SortContentItemsByName(items)

	// Assert
	assert.Equal(t, []models.SidMapping{{Name: "Alpha"}, {Name: "beta"}}, sorted)
}
