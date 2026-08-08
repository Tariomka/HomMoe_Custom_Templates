package zoneContentEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zone_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenItemsAreSorted_TheyAreOrderedCaseInsensitivelyByName(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	items := []models.SidMapping{{Name: "beta"}, {Name: "Alpha"}}

	// Act
	sorted := service.SortContentItemsByName(items)

	// Assert
	assert.Equal(t, []models.SidMapping{{Name: "Alpha"}, {Name: "beta"}}, sorted)
}

func TestWhenItemsAreSorted_TheSourceCatalogueKeepsItsOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()
	items := []models.SidMapping{{Name: "beta"}, {Name: "Alpha"}}

	// Act
	service.SortContentItemsByName(items)

	// Assert
	assert.Equal(t, []models.SidMapping{{Name: "beta"}, {Name: "Alpha"}}, items)
}
