package manualConnectionSaveModel_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

// referenceFieldCase reads a value reachable only through a slice or pointer,
// and mutates that same value in place.
type referenceFieldCase struct {
	read   func(save editor_state_model.ManualConnectionSaveModel) any
	mutate func(save editor_state_model.ManualConnectionSaveModel)
}

func TestWhenSaveIsCloned_ScalarFieldsAreCopied(t *testing.T) {
	t.Parallel()
	// Arrange
	save := newPopulatedSave()

	// Act
	clone := save.Clone()

	// Assert
	assert.Equal(t, save, clone)
}

// TestWhenAReferenceFieldIsMutatedInPlaceOnTheClone_SourceIsUnchanged walks
// every slice and pointer reachable from entities.Connection. That entity lives
// in the protected template tree and cannot carry a Clone of its own, so this
// is the only place its copy semantics are pinned.
func TestWhenAReferenceFieldIsMutatedInPlaceOnTheClone_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	for caseName, fieldCase := range referenceFieldCases() {
		t.Run(caseName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			save := newPopulatedSave()
			clone := save.Clone()
			expected := fieldCase.read(save)

			// Act
			fieldCase.mutate(clone)

			// Assert
			assert.Equal(t, expected, fieldCase.read(save))
		})
	}
}

func referenceFieldCases() map[string]referenceFieldCase {
	return map[string]referenceFieldCase{
		"WhenRoadFlagIsMutated_SourceIsUnchanged": {
			read:   func(save editor_state_model.ManualConnectionSaveModel) any { return *save.Connection.Road },
			mutate: func(save editor_state_model.ManualConnectionSaveModel) { *save.Connection.Road = false },
		},
		"WhenPortalPlacementRuleFromIsMutated_SourceIsUnchanged": {
			read: func(save editor_state_model.ManualConnectionSaveModel) any {
				return save.Connection.PortalPlacementRulesFrom[0].Type
			},
			mutate: func(save editor_state_model.ManualConnectionSaveModel) {
				save.Connection.PortalPlacementRulesFrom[0].Type = "changed"
			},
		},
		"WhenPortalPlacementRuleFromArgsIsMutated_SourceIsUnchanged": {
			read: func(save editor_state_model.ManualConnectionSaveModel) any {
				return save.Connection.PortalPlacementRulesFrom[0].Args[0]
			},
			mutate: func(save editor_state_model.ManualConnectionSaveModel) {
				save.Connection.PortalPlacementRulesFrom[0].Args[0] = "changed"
			},
		},
		"WhenPortalPlacementRuleToIsMutated_SourceIsUnchanged": {
			read: func(save editor_state_model.ManualConnectionSaveModel) any {
				return save.Connection.PortalPlacementRulesTo[0].Weight
			},
			mutate: func(save editor_state_model.ManualConnectionSaveModel) {
				save.Connection.PortalPlacementRulesTo[0].Weight = 99
			},
		},
		"WhenPortalPlacementRuleToArgsIsMutated_SourceIsUnchanged": {
			read: func(save editor_state_model.ManualConnectionSaveModel) any {
				return save.Connection.PortalPlacementRulesTo[0].Args[0]
			},
			mutate: func(save editor_state_model.ManualConnectionSaveModel) {
				save.Connection.PortalPlacementRulesTo[0].Args[0] = "changed"
			},
		},
	}
}

// newPopulatedSave builds a save whose every reference-typed field carries data,
// so that a missed copy in cloneConnection shows up as shared storage.
func newPopulatedSave() editor_state_model.ManualConnectionSaveModel {
	connection := entities.Connection{
		Name:                     "connection",
		From:                     "a",
		To:                       "b",
		Road:                     new(true),
		PortalPlacementRulesFrom: []entities.PlacementRule{{Type: "Road", Args: []any{"fromArg"}, Weight: 1}},
		PortalPlacementRulesTo:   []entities.PlacementRule{{Type: "Crossroads", Args: []any{"toArg"}, Weight: 2}},
	}

	return editor_state_model.ManualConnectionSaveModel{
		Connection: connection, IsUserAdded: true,
	}
}
