package zoneEditorGeometryService_test

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

// singleGuideZone parks one zone whose horizontal guides (312 / 350 / 388) are
// the only alignment lines a dragged zone can hold onto.
func singleGuideZone() map[string]image.Point {
	return map[string]image.Point{"A": image.Pt(350, 350)}
}

func TestWhenTheZoneRadiusIsUnknown_TheDraggedPositionIsUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	result := service.SnapPosition(image.Pt(200, 355), singleGuideZone(), 0, "B")

	// Assert
	assert.Equal(t, image.Pt(200, 355), result.Position)
}

func TestWhenTheZoneRadiusIsUnknown_NoGuideIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	result := service.SnapPosition(image.Pt(200, 355), singleGuideZone(), 0, "B")

	// Assert
	assert.Equal(t, []bool{false, false}, []bool{result.HasGuideX, result.HasGuideY})
}

func TestWhenADraggedZoneIsNearAnotherZonesGuide_ItHoldsOntoIt(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	result := service.SnapPosition(image.Pt(200, 355), singleGuideZone(), fixtureZoneRadius, "B")

	// Assert
	assert.Equal(t, image.Pt(201, 350), result.Position)
}

func TestWhenAZoneGuideIsHeld_ItsCoordinateIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	result := service.SnapPosition(image.Pt(200, 355), singleGuideZone(), fixtureZoneRadius, "B")

	// Assert
	assert.InDelta(t, 312.0, result.GuideY, 1e-9)
}

func TestWhenOnlyTheGridIsInReach_NoZoneGuideIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	result := service.SnapPosition(image.Pt(200, 355), singleGuideZone(), fixtureZoneRadius, "B")

	// Assert
	assert.False(t, result.HasGuideX)
}

func TestWhenTheDraggedZoneIsTheOnlyZone_ItDoesNotHoldOntoItself(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	result := service.SnapPosition(image.Pt(350, 355), singleGuideZone(), fixtureZoneRadius, "A")

	// Assert
	assert.False(t, result.HasGuideY)
}
