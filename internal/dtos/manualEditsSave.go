package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/entities"

// ManualZoneSave persists a zone edited in the manual zone editor. The zone's
// ManualPosition is captured separately because entities.Zone deliberately
// omits it from JSON (json:"-"), yet that normalized position is the essential
// piece of a hand-made layout and must survive a save/load round-trip.
type ManualZoneSave struct {
	Zone           entities.Zone `json:"zone"`
	ManualPosition *[2]float64   `json:"manualPosition,omitempty"`
}

// ManualConnectionSave persists a connection edited in the manual zone editor,
// capturing the runtime-only IsUserAdded flag that entities.Connection omits
// from JSON (json:"-").
type ManualConnectionSave struct {
	Connection  entities.Connection `json:"connection"`
	IsUserAdded bool                `json:"isUserAdded,omitempty"`
}

// ToManualZoneSaves converts live editor zones into their serializable form,
// preserving each zone's ManualPosition outside the entities.Zone JSON.
func ToManualZoneSaves(zones []entities.Zone) []ManualZoneSave {
	if len(zones) == 0 {
		return nil
	}
	saves := make([]ManualZoneSave, 0, len(zones))
	for _, zone := range zones {
		saves = append(saves, ManualZoneSave{Zone: zone, ManualPosition: zone.ManualPosition})
	}
	return saves
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

// FromManualZoneSaves rebuilds live editor zones from their serialized form,
// restoring each zone's ManualPosition.
func FromManualZoneSaves(saves []ManualZoneSave) []entities.Zone {
	if len(saves) == 0 {
		return nil
	}
	zones := make([]entities.Zone, 0, len(saves))
	for _, save := range saves {
		zone := save.Zone
		zone.ManualPosition = save.ManualPosition
		zones = append(zones, zone)
	}
	return zones
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
