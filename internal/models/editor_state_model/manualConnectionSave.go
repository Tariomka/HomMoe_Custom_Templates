package editor_state_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// ManualConnectionSave adds behaviour to the behaviour-free
// editor_state.ManualConnectionSave entity.
type ManualConnectionSave struct {
	editor_state.ManualConnectionSave
}

// ToManualConnectionSaves converts live editor connections into their
// serializable form, preserving the IsUserAdded flag that the .rmg.json schema
// has nowhere to put.
func ToManualConnectionSaves(connections []template_model.Connection) []ManualConnectionSave {
	if len(connections) == 0 {
		return nil
	}

	return linq.FromSlice(connections).
		Select(func(connection template_model.Connection) ManualConnectionSave {
			return ManualConnectionSave{
				Connection:  template_model.ToConnectionEntity(connection),
				IsUserAdded: connection.IsUserAdded}
		}).ToSlice()
}

// FromManualConnectionSaves rebuilds live editor connections from their
// serialized form, restoring the IsUserAdded flag.
func FromManualConnectionSaves(saves []ManualConnectionSave) []template_model.Connection {
	if len(saves) == 0 {
		return nil
	}

	return linq.FromSlice(saves).
		Select(func(save ManualConnectionSave) template_model.Connection {
			connection := template_model.ToConnectionModel(save.Connection)
			connection.IsUserAdded = save.IsUserAdded
			return connection
		}).ToSlice()
}

func ToManualConnectionSaveModels(saves []editor_state.ManualConnectionSave) []ManualConnectionSave {
	if len(saves) == 0 {
		return nil
	}

	return linq.FromSlice(saves).
		Select(func(save editor_state.ManualConnectionSave) ManualConnectionSave {
			return ManualConnectionSave{ManualConnectionSave: save}
		}).ToSlice()
}

func ToManualConnectionSaveEntities(saves []ManualConnectionSave) []editor_state.ManualConnectionSave {
	if len(saves) == 0 {
		return nil
	}

	return linq.FromSlice(saves).
		Select(func(save ManualConnectionSave) editor_state.ManualConnectionSave {
			return save.ManualConnectionSave
		}).ToSlice()
}

func (this ManualConnectionSave) Clone() ManualConnectionSave {
	return ManualConnectionSave{
		Connection:  cloneConnection(this.Connection),
		IsUserAdded: this.IsUserAdded,
	}
}

func cloneConnection(source entities.Connection) entities.Connection {
	clone := source
	clone.Road = helpers.ClonePointer(source.Road)
	clone.PortalPlacementRulesFrom = clonePlacementRules(source.PortalPlacementRulesFrom)
	clone.PortalPlacementRulesTo = clonePlacementRules(source.PortalPlacementRulesTo)
	return clone
}

func clonePlacementRules(source []entities.PlacementRule) []entities.PlacementRule {
	clone := slices.Clone(source)
	for ruleIndex := range clone {
		// Args elements are opaque: JSON decoding only ever boxes immutable
		// scalars into them, so cloning the slice is enough.
		clone[ruleIndex].Args = slices.Clone(clone[ruleIndex].Args)
	}
	return clone
}
