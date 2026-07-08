package math_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/utils"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneAreaEqualsReferenceArea_ReturnsOne(t *testing.T) {
	// Arrange
	mapSize := 160
	totalZones := 4 // 160*160/4 equals the 80x80 reference zone area.

	// Act
	actual := utils.ComputeContentScale(mapSize, totalZones)

	// Assert
	assert.Equal(t, 1.0, actual)
}

func TestWhenZoneAreaIsFourTimesReferenceArea_ReturnsTwo(t *testing.T) {
	// Arrange
	mapSize := 320
	totalZones := 4

	// Act
	actual := utils.ComputeContentScale(mapSize, totalZones)

	// Assert
	assert.Equal(t, 2.0, actual)
}

func TestWhenZoneAreaIsTiny_ClampsScaleToMinimumHalf(t *testing.T) {
	// Arrange
	mapSize := 20
	totalZones := 30

	// Act
	actual := utils.ComputeContentScale(mapSize, totalZones)

	// Assert
	assert.Equal(t, 0.5, actual)
}

func TestWhenZoneAreaIsHuge_ClampsScaleToMaximumTwoAndHalf(t *testing.T) {
	// Arrange
	mapSize := 900
	totalZones := 1

	// Act
	actual := utils.ComputeContentScale(mapSize, totalZones)

	// Assert
	assert.Equal(t, 2.5, actual)
}

func TestWhenTotalZonesIsZero_TreatsZoneCountAsOne(t *testing.T) {
	// Arrange
	mapSize := 160
	totalZones := 0 // Guarded by math.Max(1, ...) so the zone area is the whole map.

	// Act
	actual := utils.ComputeContentScale(mapSize, totalZones)

	// Assert
	assert.Equal(t, 2.0, actual)
}
