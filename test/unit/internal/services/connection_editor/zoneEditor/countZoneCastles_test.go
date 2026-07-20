package zoneEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneHasMixedMainObjects_CountsOnlyCities(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := entities.Zone{
		MainObjects: []entities.MainObject{
			{Type: "Spawn"},
			{Type: "City"},
			{Type: "AbandonedOutpost"},
			{Type: "City"},
		},
	}

	// Act
	count := connection_editor.CountZoneCastles(zone)

	// Assert
	assert.Equal(t, 2, count)
}

func TestWhenCityTypeDiffersInCase_DoesNotCountIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := entities.Zone{
		MainObjects: []entities.MainObject{{Type: "city"}},
	}

	// Act
	count := connection_editor.CountZoneCastles(zone)

	// Assert
	assert.Equal(t, 0, count)
}

func TestWhenZoneHasNoMainObjects_ReturnsZero(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := entities.Zone{}

	// Act
	count := connection_editor.CountZoneCastles(zone)

	// Assert
	assert.Equal(t, 0, count)
}
