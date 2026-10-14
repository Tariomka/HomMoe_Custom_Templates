package editorState_test

import (
	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type stateValidationFunc func(editor_state_model.EditorState, bool) editor_state_dto.EditorStateValidationDto

func (this stateValidationFunc) ValidateEditorState(
	stateDto editor_state_model.EditorState,
	fixIssues bool,
) editor_state_dto.EditorStateValidationDto {
	return this(stateDto, fixIssues)
}

func newEditorState() *models.EditorState {
	return newEditorStateWithValidation(
		func(stateDto editor_state_model.EditorState, _ bool) editor_state_dto.EditorStateValidationDto {
			return editor_state_dto.EditorStateValidationDto{State: stateDto}
		},
	)
}

func newEditorStateWithValidation(validation stateValidationFunc) *models.EditorState {
	return models.NewEditorState(validation)
}

// manualZoneFieldCase reads a value reachable only through a slice or pointer,
// and mutates that same value in place.
type manualZoneFieldCase struct {
	read   func(edit editor_state_model.ManualZoneEdit) any
	mutate func(edit editor_state_model.ManualZoneEdit)
}

// manualZoneFieldCases covers the nested storage a caller still holds a handle
// on after handing edits to the state, or receives back from it.
func manualZoneFieldCases() map[string]manualZoneFieldCase {
	return map[string]manualZoneFieldCase{
		"WhenManualPositionIsMutated_TheSnapshotIsUnchanged": {
			read:   func(edit editor_state_model.ManualZoneEdit) any { return *edit.ManualPosition },
			mutate: func(edit editor_state_model.ManualZoneEdit) { edit.ManualPosition.X = 9 },
		},
		"WhenGeneratorPositionIsMutated_TheSnapshotIsUnchanged": {
			read:   func(edit editor_state_model.ManualZoneEdit) any { return *edit.GeneratorPosition },
			mutate: func(edit editor_state_model.ManualZoneEdit) { edit.GeneratorPosition.Y = 9 },
		},
		"WhenGeneratorRingIsMutated_TheSnapshotIsUnchanged": {
			read:   func(edit editor_state_model.ManualZoneEdit) any { return *edit.GeneratorRing },
			mutate: func(edit editor_state_model.ManualZoneEdit) { *edit.GeneratorRing = 9 },
		},
		"WhenQualityIsMutated_TheSnapshotIsUnchanged": {
			read:   func(edit editor_state_model.ManualZoneEdit) any { return *edit.Quality },
			mutate: func(edit editor_state_model.ManualZoneEdit) { *edit.Quality = int8(neutral_zone.QualityHighest) },
		},
		"WhenEncounterHolesSettingsIsMutated_TheSnapshotIsUnchanged": {
			read: func(edit editor_state_model.ManualZoneEdit) any { return *edit.FieldOverrides.EncounterHolesSettings },
			mutate: func(edit editor_state_model.ManualZoneEdit) {
				edit.FieldOverrides.EncounterHolesSettings.TwoHoleEncounters = 9
			},
		},
		"WhenAContentPoolEntryIsMutated_TheSnapshotIsUnchanged": {
			read: func(edit editor_state_model.ManualZoneEdit) any { return edit.FieldOverrides.GuardedContentPool[0] },
			mutate: func(edit editor_state_model.ManualZoneEdit) {
				edit.FieldOverrides.GuardedContentPool[0] = "changed"
			},
		},
		"WhenAMainObjectFactionIsMutated_TheSnapshotIsUnchanged": {
			read: func(edit editor_state_model.ManualZoneEdit) any {
				return edit.FieldOverrides.MainObjects[0].Factions[0]
			},
			mutate: func(edit editor_state_model.ManualZoneEdit) {
				edit.FieldOverrides.MainObjects[0].Factions[0] = "changed"
			},
		},
		"WhenARoadIsMutated_TheSnapshotIsUnchanged": {
			read:   func(edit editor_state_model.ManualZoneEdit) any { return edit.FieldOverrides.Roads[0].Type },
			mutate: func(edit editor_state_model.ManualZoneEdit) { edit.FieldOverrides.Roads[0].Type = "changed" },
		},
		"WhenARoadFlagIsMutated_TheSnapshotIsUnchanged": {
			read:   func(edit editor_state_model.ManualZoneEdit) any { return *edit.FieldOverrides.Roads[0].Road },
			mutate: func(edit editor_state_model.ManualZoneEdit) { *edit.FieldOverrides.Roads[0].Road = false },
		},
	}
}

// newPopulatedManualZone builds a zone whose reference-typed fields all carry
// data, so a shallow copy at the snapshot boundary shows up as shared storage.
func newPopulatedManualZone() template_model.Zone {
	return template_model.Zone{
		Name:              "Zone A",
		Quality:           new(neutral_zone.QualityLowest),
		GeneratorPosition: new(data.NewVec2(1.0, 2.0)),
		GeneratorRing:     new(3),
		ManualPosition:    new(data.NewVec2(4.0, 5.0)),
		EncounterHolesSettings: &template_model.EncounterHolesSettings{
			AffectedEncounters: 1,
			TwoHoleEncounters:  2,
		},
		GuardedContentPool: []string{"guarded"},
		MainObjects: []template_model.MainObject{{
			Type:     "City",
			Factions: []string{"faction"},
		}},
		Roads: []template_model.Road{{
			Type: "road",
			From: template_model.TypedRef{Type: "from", Args: []string{"fromArg"}},
			To:   template_model.TypedRef{Type: "to", Args: []string{"toArg"}},
			Road: new(true),
		}},
	}
}

// newPopulatedManualEdits encodes newPopulatedManualZone as a manual-edit layer.
func newPopulatedManualEdits() editor_state_model.ManualEditSettings {
	return editor_state_model.ManualEditSettingsFromSnapshot([]template_model.Zone{newPopulatedManualZone()}, nil)
}
