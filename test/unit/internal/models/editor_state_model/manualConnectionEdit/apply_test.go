package manualConnectionEdit_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenAppliedToTheDerivedConnection_EditedConnectionIsReproduced(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	edited := newEditedConnection()
	edit := editor_state_model.NewManualConnectionEdit(derived, edited)

	// Act
	result := edit.Apply(derived)

	// Assert
	assert.Equal(t, edited, result)
}

func TestWhenEditedRoadFlagIsNil_AppliedConnectionHasNilRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	edited := derived.Clone()
	edited.Road = nil
	edit := editor_state_model.NewManualConnectionEdit(derived, edited)

	// Act
	result := edit.Apply(derived)

	// Assert
	assert.Nil(t, result.Road)
}

func TestWhenEditedRulesAreNil_AppliedConnectionHasNilRules(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	edited := derived.Clone()
	edited.PortalPlacementRulesFrom = nil
	edit := editor_state_model.NewManualConnectionEdit(derived, edited)

	// Act
	result := edit.Apply(derived)

	// Assert
	assert.Nil(t, result.PortalPlacementRulesFrom)
}

func TestWhenAppliedToADifferentDerivedConnection_StructureComesFromTheEdit(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := editor_state_model.NewManualConnectionEdit(newDerivedConnection(), newEditedConnection())
	otherDerived := template_model.Connection{Name: "Other", From: "X", To: "Y", ConnectionType: "Direct"}

	// Act
	result := edit.Apply(otherDerived)

	// Assert
	assert.Equal(t, "Edited", result.Name)
}

func TestWhenAppliedToADifferentDerivedConnection_UntouchedFieldsComeFromDerived(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	edited := derived.Clone()
	edited.Length = 9
	edit := editor_state_model.NewManualConnectionEdit(derived, edited)
	otherDerived := template_model.Connection{GatePlacement: "Edge"}

	// Act
	result := edit.Apply(otherDerived)

	// Assert
	assert.Equal(t, "Edge", result.GatePlacement)
}

func TestWhenAppliedToAnEmptyConnection_PopulatedConnectionIsReproduced(t *testing.T) {
	t.Parallel()
	// Arrange
	edited := newEditedConnection()
	edit := editor_state_model.NewManualConnectionEdit(template_model.Connection{}, edited)

	// Act
	result := edit.Apply(template_model.Connection{})

	// Assert
	assert.Equal(t, edited, result)
}

func TestWhenTheResultRuleIsMutated_ReapplyingYieldsTheOriginalRule(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	edit := editor_state_model.NewManualConnectionEdit(derived, newEditedConnection())
	edit.Apply(derived).PortalPlacementRulesFrom[0].Type = "mutated"

	// Act
	result := edit.Apply(derived)

	// Assert
	assert.Equal(t, "edited", result.PortalPlacementRulesFrom[0].Type)
}

func TestWhenTheResultRoadFlagIsMutated_ReapplyingYieldsTheOriginalFlag(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	edit := editor_state_model.NewManualConnectionEdit(derived, newEditedConnection())
	*edit.Apply(derived).Road = true

	// Act
	result := edit.Apply(derived)

	// Assert
	assert.Equal(t, new(false), result.Road)
}

func TestWhenTheResultRuleIsMutated_DerivedConnectionIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedConnection()
	original := derived.Clone()
	edit := editor_state_model.NewManualConnectionEdit(derived, derived)

	// Act
	edit.Apply(derived).PortalPlacementRulesFrom[0].Type = "mutated"

	// Assert
	assert.Equal(t, original, derived)
}
