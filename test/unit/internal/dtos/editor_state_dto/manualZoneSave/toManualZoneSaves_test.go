package manualZoneSave_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneListIsEmpty_ReturnsNil(t *testing.T) {
	// Arrange
	var zones []entities.Zone

	// Act
	saves := editor_state_dto.ToManualZoneSaves(zones)

	// Assert
	assert.Nil(t, saves)
}

func TestWhenZonesHaveManualPositions_PreservesEachPositionInSave(t *testing.T) {
	// Arrange
	firstPosition := &[2]float64{0.25, 0.75}
	secondPosition := &[2]float64{0.5, 0.5}
	zones := []entities.Zone{
		{Name: "Zone A", ManualPosition: firstPosition},
		{Name: "Zone B", ManualPosition: secondPosition},
	}
	expected := []editor_state_dto.ManualZoneSave{
		{Zone: zones[0], ManualPosition: firstPosition},
		{Zone: zones[1], ManualPosition: secondPosition},
	}

	// Act
	saves := editor_state_dto.ToManualZoneSaves(zones)

	// Assert
	assert.Equal(t, expected, saves)
}

func TestWhenZoneHasNoManualPosition_SavesNilPosition(t *testing.T) {
	// Arrange
	zones := []entities.Zone{{Name: "Zone A"}}

	// Act
	saves := editor_state_dto.ToManualZoneSaves(zones)

	// Assert
	assert.Nil(t, saves[0].ManualPosition)
}
