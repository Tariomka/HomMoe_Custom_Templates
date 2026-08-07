package zoneEditorGeometryService_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenTheZoneRadiusIsKnown_TheGridHoldsSevenCellsPerZoneDiameter(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	step := service.GridStep(7)

	// Assert
	assert.InDelta(t, 2.0, step, 1e-9)
}

func TestWhenTheZoneRadiusIsUnknown_TheGridCollapses(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	step := service.GridStep(0)

	// Assert
	assert.Zero(t, step)
}
