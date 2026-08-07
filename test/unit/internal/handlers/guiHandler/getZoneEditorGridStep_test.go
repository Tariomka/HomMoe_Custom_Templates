package guiHandler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenTheGridStepRequested_HoldsSevenCellsPerZoneDiameter(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()

	// Act
	step := handler.GetZoneEditorGridStep(7)

	// Assert
	assert.InDelta(t, 2.0, step, 1e-9)
}
