package manualZoneEdit_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoCastleCountIsRecorded_GetCastlesReportsNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit

	// Act
	_, ok := edit.GetCastles()

	// Assert
	assert.False(t, ok)
}

func TestWhenNoCastleCountIsRecorded_GetCastlesReturnsZero(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit

	// Act
	castles, _ := edit.GetCastles()

	// Assert
	assert.Equal(t, 0, castles)
}

func TestWhenACastleCountIsRecorded_GetCastlesReportsFound(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit
	edit = edit.AsAdded(3)

	// Act
	_, ok := edit.GetCastles()

	// Assert
	assert.True(t, ok)
}

func TestWhenACastleCountIsRecorded_GetCastlesReturnsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit
	edit = edit.AsAdded(3)

	// Act
	castles, _ := edit.GetCastles()

	// Assert
	assert.Equal(t, 3, castles)
}

func TestWhenAZeroCastleCountIsRecorded_GetCastlesReportsFound(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit
	edit = edit.AsReprofiled(0)

	// Act
	_, ok := edit.GetCastles()

	// Assert
	assert.True(t, ok)
}
