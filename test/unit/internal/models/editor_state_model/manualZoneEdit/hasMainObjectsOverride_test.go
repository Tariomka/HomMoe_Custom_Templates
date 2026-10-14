package manualZoneEdit_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenMainObjectsAreOverridden_HasMainObjectsOverrideIsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit
	edit.FieldOverrides.MainObjects = []template_entity.MainObject{{Type: "City"}}

	// Act
	has := edit.HasMainObjectsOverride()

	// Assert
	assert.True(t, has)
}

func TestWhenMainObjectsAreOverriddenWithAnEmptyList_HasMainObjectsOverrideIsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit
	edit.FieldOverrides.MainObjects = []template_entity.MainObject{}

	// Act
	has := edit.HasMainObjectsOverride()

	// Assert
	assert.True(t, has)
}

func TestWhenMainObjectsAreUnset_HasMainObjectsOverrideIsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit
	edit.Unset = []string{"layout", "mainObjects"}

	// Act
	has := edit.HasMainObjectsOverride()

	// Assert
	assert.True(t, has)
}

func TestWhenOnlyOtherFieldsAreOverriddenOrUnset_HasMainObjectsOverrideIsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit
	edit.FieldOverrides.Size = new(2.0)
	edit.Unset = []string{"roads"}

	// Act
	has := edit.HasMainObjectsOverride()

	// Assert
	assert.False(t, has)
}

func TestWhenNothingIsRecorded_HasMainObjectsOverrideIsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit

	// Act
	has := edit.HasMainObjectsOverride()

	// Assert
	assert.False(t, has)
}
