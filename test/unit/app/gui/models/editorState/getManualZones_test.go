package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenManualZonesWereStored_ZonesRoundTripWithManualPositions(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	zones := []entities.Zone{
		{
			Name:           "Zone A",
			Size:           gofakeit.Float64Range(0.5, 2),
			ManualPosition: &[2]float64{gofakeit.Float64Range(0, 1), gofakeit.Float64Range(0, 1)},
		},
		{Name: "Zone B"},
	}
	state.SetManualEdits(zones, nil)

	// Act
	restored := state.GetManualZones()

	// Assert
	assert.Equal(t, zones, restored)
}

func TestWhenNoManualZonesWereStored_NilZonesAreReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	restored := state.GetManualZones()

	// Assert
	assert.Nil(t, restored)
}
