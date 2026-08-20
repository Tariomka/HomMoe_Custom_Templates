package zoneEditorGeometryService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

// singleGuideZone parks one zone whose horizontal guides (312 / 350 / 388) are
// the only alignment lines a dragged zone can hold onto.
func singleGuideZone() map[string]models.Position {
	return map[string]models.Position{"A": data.NewVec2(350.0, 350.0)}
}

func TestWhenTheZoneRadiusIsUnknown_TheDraggedPositionIsUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	result := service.SnapPosition(data.NewVec2(200.0, 355.0), singleGuideZone(), 0, "B")

	// Assert
	assert.Equal(t, data.NewVec2(200.0, 355.0), result.Position)
}

func TestWhenTheZoneRadiusIsUnknown_NoGuideIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	result := service.SnapPosition(data.NewVec2(200.0, 355.0), singleGuideZone(), 0, "B")

	// Assert
	assert.Equal(t, []bool{false, false}, []bool{result.HasGuideX, result.HasGuideY})
}

func TestWhenADraggedZoneIsNearAnotherZonesGuide_ItHoldsOntoIt(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	result := service.SnapPosition(data.NewVec2(200.0, 355.0), singleGuideZone(), fixtureZoneRadius, "B")

	// Assert
	assert.InDelta(t, 350.0, result.Position.Y, 1e-9)
}

// The grid step is 2*38/7 px, so the leading edge at x=162 holds onto the 15th
// grid line and carries the centre 6/7 px to the right - a fraction the editor
// now keeps instead of rounding away.
func TestWhenOnlyTheGridIsInReach_TheDraggedPositionKeepsTheFractionalCorrection(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	result := service.SnapPosition(data.NewVec2(200.0, 355.0), singleGuideZone(), fixtureZoneRadius, "B")

	// Assert
	assert.InDelta(t, 200.0+6.0/7.0, result.Position.X, 1e-9)
}

func TestWhenAZoneGuideIsHeld_ItsCoordinateIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	result := service.SnapPosition(data.NewVec2(200.0, 355.0), singleGuideZone(), fixtureZoneRadius, "B")

	// Assert
	assert.InDelta(t, 312.0, result.GuideY, 1e-9)
}

func TestWhenOnlyTheGridIsInReach_NoZoneGuideIsReported(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	result := service.SnapPosition(data.NewVec2(200.0, 355.0), singleGuideZone(), fixtureZoneRadius, "B")

	// Assert
	assert.False(t, result.HasGuideX)
}

func TestWhenTheDraggedZoneIsTheOnlyZone_ItDoesNotHoldOntoItself(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	result := service.SnapPosition(data.NewVec2(350.0, 355.0), singleGuideZone(), fixtureZoneRadius, "A")

	// Assert
	assert.False(t, result.HasGuideY)
}
