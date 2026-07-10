package manualZoneSave_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestWhenSaveListIsEmpty_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var saves []editor_state_dto.ManualZoneSave

	// Act
	zones := editor_state_dto.FromManualZoneSaves(saves)

	// Assert
	assert.Nil(t, zones)
}

func TestWhenSavesCarryManualPositions_RestoresEachPositionOntoZone(t *testing.T) {
	t.Parallel()
	// Arrange
	firstPosition := &[2]float64{0.1, 0.9}
	secondPosition := &[2]float64{0.6, 0.4}
	saves := []editor_state_dto.ManualZoneSave{
		{Zone: entities.Zone{Name: "Zone A"}, ManualPosition: firstPosition},
		{Zone: entities.Zone{Name: "Zone B"}, ManualPosition: secondPosition},
	}
	expected := []entities.Zone{
		{Name: "Zone A", ManualPosition: firstPosition},
		{Name: "Zone B", ManualPosition: secondPosition},
	}

	// Act
	zones := editor_state_dto.FromManualZoneSaves(saves)

	// Assert
	assert.Equal(t, expected, zones)
}

func TestWhenSavePositionDiffersFromEmbeddedZonePosition_SavePositionWins(t *testing.T) {
	t.Parallel()
	// Arrange
	savedPosition := &[2]float64{0.2, 0.3}
	staleEmbeddedPosition := &[2]float64{0.8, 0.8}
	saves := []editor_state_dto.ManualZoneSave{
		{
			Zone:           entities.Zone{Name: "Zone A", ManualPosition: staleEmbeddedPosition},
			ManualPosition: savedPosition,
		},
	}

	// Act
	zones := editor_state_dto.FromManualZoneSaves(saves)

	// Assert
	assert.Same(t, savedPosition, zones[0].ManualPosition)
}
