package manualConnectionEdit_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenConnectionsAreEqual_NoFieldIsOverridden(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := newDerivedConnection()

	// Act
	edit := editor_state_model.NewManualConnectionEdit(connection, connection)

	// Assert
	assert.Equal(t, editor_state.ConnectionFieldOverrides{}, edit.FieldOverrides)
}

func TestWhenConnectionsAreEqual_UnsetIsEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := newDerivedConnection()

	// Act
	edit := editor_state_model.NewManualConnectionEdit(connection, connection)

	// Assert
	assert.Empty(t, edit.Unset)
}

func TestWhenOnlyAScalarDiffers_ThatScalarIsStoredAsAPointerOverride(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	edited := derived.Clone()
	edited.Length = 7.5

	// Act
	edit := editor_state_model.NewManualConnectionEdit(derived, edited)

	// Assert
	assert.Equal(t, new(7.5), edit.FieldOverrides.Length)
}

func TestWhenOnlyAScalarDiffers_OtherFieldsStayUnset(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	edited := derived.Clone()
	edited.Length = 7.5

	// Act
	edit := editor_state_model.NewManualConnectionEdit(derived, edited)

	// Assert
	assert.Nil(t, edit.FieldOverrides.GuardValue)
}

func TestWhenEditedNumberIsZeroAndDerivedIsNot_ZeroOverrideIsStored(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	edited := derived.Clone()
	edited.GuardValue = 0

	// Act
	edit := editor_state_model.NewManualConnectionEdit(derived, edited)

	// Assert
	assert.Equal(t, new(0), edit.FieldOverrides.GuardValue)
}

func TestWhenEditedFlagIsFalseAndDerivedIsTrue_FalseOverrideIsStored(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	derived.SimTurnSquad = true
	edited := derived.Clone()
	edited.SimTurnSquad = false

	// Act
	edit := editor_state_model.NewManualConnectionEdit(derived, edited)

	// Assert
	assert.Equal(t, new(false), edit.FieldOverrides.SimTurnSquad)
}

func TestWhenTheRoadFlagDiffers_RoadOverrideHoldsTheEditedFlag(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	edited := derived.Clone()
	edited.Road = new(false)

	// Act
	edit := editor_state_model.NewManualConnectionEdit(derived, edited)

	// Assert
	assert.Equal(t, new(false), edit.FieldOverrides.Road)
}

func TestWhenTheRoadFlagIsNilAndDerivedHasOne_RoadIsListedInUnset(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	edited := derived.Clone()
	edited.Road = nil

	// Act
	edit := editor_state_model.NewManualConnectionEdit(derived, edited)

	// Assert
	assert.Equal(t, []string{"road"}, edit.Unset)
}

func TestWhenTheRoadFlagIsSetAndDerivedHasNone_RoadOverrideIsStored(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	derived.Road = nil
	edited := derived.Clone()
	edited.Road = new(true)

	// Act
	edit := editor_state_model.NewManualConnectionEdit(derived, edited)

	// Assert
	assert.Equal(t, new(true), edit.FieldOverrides.Road)
}

func TestWhenEditedRulesAreNilAndDerivedHasRules_RuleListIsListedInUnset(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	edited := derived.Clone()
	edited.PortalPlacementRulesFrom = nil

	// Act
	edit := editor_state_model.NewManualConnectionEdit(derived, edited)

	// Assert
	assert.Equal(t, []string{"portalPlacementRulesFrom"}, edit.Unset)
}

func TestWhenEditedRulesAreEmptyAndDerivedHasRules_EmptyOverrideIsStored(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	edited := derived.Clone()
	edited.PortalPlacementRulesFrom = []template_model.PlacementRule{}

	// Act
	edit := editor_state_model.NewManualConnectionEdit(derived, edited)

	// Assert
	assert.Equal(t, []template_entity.PlacementRule{}, edit.FieldOverrides.PortalPlacementRulesFrom)
}

func TestWhenStructureFieldsDiffer_StructureIsCopiedFromTheEditedConnection(t *testing.T) {
	t.Parallel()
	for caseName, testCase := range structureCases() {
		t.Run(caseName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			edited := newEditedConnection()

			// Act
			edit := editor_state_model.NewManualConnectionEdit(newDerivedConnection(), edited)

			// Assert
			assert.Equal(t, testCase.expected, testCase.read(edit))
		})
	}
}

func TestWhenOnlyStructureFieldsDiffer_NoFieldIsOverridden(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	edited := derived.Clone()
	edited.Name = "Renamed"
	edited.From = "Other"
	edited.To = "Another"
	edited.ConnectionType = "Direct"
	edited.IsUserAdded = true

	// Act
	edit := editor_state_model.NewManualConnectionEdit(derived, edited)

	// Assert
	assert.Equal(t, editor_state.ConnectionFieldOverrides{}, edit.FieldOverrides)
}

func TestWhenEditedRuleIsMutatedAfterwards_StoredOverrideIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	edited := newEditedConnection()
	edit := editor_state_model.NewManualConnectionEdit(derived, edited)

	// Act
	edited.PortalPlacementRulesFrom[0].Type = "mutated"

	// Assert
	assert.Equal(t, "edited", edit.FieldOverrides.PortalPlacementRulesFrom[0].Type)
}

func TestWhenEditedRoadFlagIsMutatedAfterwards_StoredOverrideIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	edited := newEditedConnection()
	edit := editor_state_model.NewManualConnectionEdit(derived, edited)

	// Act
	*edited.Road = true

	// Assert
	assert.Equal(t, new(false), edit.FieldOverrides.Road)
}

type structureCase struct {
	expected any
	read     func(edit editor_state_model.ManualConnectionEdit) any
}

func structureCases() map[string]structureCase {
	return map[string]structureCase{
		"WhenTheNameDiffers_NameIsCopied": {
			expected: "Edited",
			read:     func(edit editor_state_model.ManualConnectionEdit) any { return edit.Name },
		},
		"WhenTheSourceDiffers_SourceIsCopied": {
			expected: "Zone A",
			read:     func(edit editor_state_model.ManualConnectionEdit) any { return edit.From },
		},
		"WhenTheTargetDiffers_TargetIsCopied": {
			expected: "Zone C",
			read:     func(edit editor_state_model.ManualConnectionEdit) any { return edit.To },
		},
		"WhenTheTypeDiffers_TypeIsCopied": {
			expected: "Portal",
			read:     func(edit editor_state_model.ManualConnectionEdit) any { return edit.ConnectionType },
		},
		"WhenTheConnectionIsUserAdded_ProvenanceIsCopied": {
			expected: true,
			read:     func(edit editor_state_model.ManualConnectionEdit) any { return edit.IsUserAdded },
		},
	}
}
