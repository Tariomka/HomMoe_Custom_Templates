package editor_state_dto

import "github.com/Tariomka/hommoe_custom_templates/internal/entities"

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
