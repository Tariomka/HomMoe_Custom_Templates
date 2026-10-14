package manualEditSettings_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheSettingsAreCloned_CloneEqualsTheSource(t *testing.T) {
	t.Parallel()
	// Arrange
	settings := newPopulatedManualEditSettings()

	// Act
	clone := settings.Clone()

	// Assert
	assert.Equal(t, settings, clone)
}

func TestWhenTheSettingsAreEmpty_CloneListsStayNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var settings editor_state_model.ManualEditSettings

	// Act
	clone := settings.Clone()

	// Assert
	assert.Equal(t, editor_state_model.ManualEditSettings{}, clone)
}

func TestWhenAZoneEditOfTheCloneIsMutated_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	settings := newPopulatedManualEditSettings()
	clone := settings.Clone()

	// Act
	clone.ManualZones[0].Name = "mutated"

	// Assert
	assert.Equal(t, "Zone A", settings.ManualZones[0].Name)
}

func TestWhenAZoneOverrideOfTheCloneIsMutated_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	settings := newPopulatedManualEditSettings()
	clone := settings.Clone()

	// Act
	clone.ManualZones[0].FieldOverrides.GuardedContentPool[0] = "mutated"

	// Assert
	assert.Equal(t, []string{"guarded"}, settings.ManualZones[0].FieldOverrides.GuardedContentPool)
}

func TestWhenAConnectionEditOfTheCloneIsMutated_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	settings := newPopulatedManualEditSettings()
	clone := settings.Clone()

	// Act
	clone.ManualConnections[0].Name = "mutated"

	// Assert
	assert.Equal(t, "A-B", settings.ManualConnections[0].Name)
}

func TestWhenAConnectionRuleOfTheCloneIsMutated_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	settings := newPopulatedManualEditSettings()
	clone := settings.Clone()

	// Act
	clone.ManualConnections[0].FieldOverrides.PortalPlacementRulesFrom[0].Type = "mutated"

	// Assert
	assert.Equal(t, "rule", settings.ManualConnections[0].FieldOverrides.PortalPlacementRulesFrom[0].Type)
}

func TestWhenAZoneEditIsAppendedToTheClone_SourceKeepsItsLength(t *testing.T) {
	t.Parallel()
	// Arrange
	settings := newPopulatedManualEditSettings()
	clone := settings.Clone()

	// Act
	clone.ManualZones = append(clone.ManualZones, editor_state_model.ManualZoneEdit{})

	// Assert
	assert.Len(t, settings.ManualZones, 1)
}
