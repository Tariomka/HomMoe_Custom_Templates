package editor_state_dto

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

// ManualConnectionSave persists a connection edited in the manual zone editor,
// capturing the runtime-only IsUserAdded flag that entities.Connection omits
// from JSON (json:"-").
type ManualConnectionSave struct {
	Connection  entities.Connection `json:"connection"`
	IsUserAdded bool                `json:"isUserAdded,omitempty"`
}

// ToManualConnectionSaves converts live editor connections into their
// serializable form, preserving the IsUserAdded flag.
func ToManualConnectionSaves(connections []entities.Connection) []ManualConnectionSave {
	if len(connections) == 0 {
		return nil
	}
	saves := make([]ManualConnectionSave, 0, len(connections))
	for _, connection := range connections {
		saves = append(saves, ManualConnectionSave{Connection: connection, IsUserAdded: connection.IsUserAdded})
	}
	return saves
}

// FromManualConnectionSaves rebuilds live editor connections from their
// serialized form, restoring the IsUserAdded flag.
func FromManualConnectionSaves(saves []ManualConnectionSave) []entities.Connection {
	if len(saves) == 0 {
		return nil
	}
	connections := make([]entities.Connection, 0, len(saves))
	for _, save := range saves {
		connection := save.Connection
		connection.IsUserAdded = save.IsUserAdded
		connections = append(connections, connection)
	}
	return connections
}

// Clone returns a copy that shares no backing array or pointer with the
// receiver. entities.Connection lives in the protected template tree and
// therefore carries no Clone of its own, so every one of its reference-typed
// fields is copied here; a field added there must be added to cloneConnection
// as well.
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
