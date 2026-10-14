package manualZoneEdit_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenAnAddedZoneIsMarkedReprofiled_OriginIsCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit
	edit = edit.AsAdded(1)

	// Act
	reprofiled := edit.AsReprofiled(2)

	// Assert
	assert.Equal(t, editor_state.ZoneOrigin(""), reprofiled.Origin)
}

func TestWhenAZoneIsMarkedReprofiled_CastlesAreRecorded(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit

	// Act
	reprofiled := edit.AsReprofiled(4)

	// Assert
	assert.Equal(t, new(4), reprofiled.Castles)
}

func TestWhenAZoneIsMarkedReprofiled_IsAddedIsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit
	edit = edit.AsAdded(1)

	// Act
	reprofiled := edit.AsReprofiled(2)

	// Assert
	assert.False(t, reprofiled.IsAdded())
}
