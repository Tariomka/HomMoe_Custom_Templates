// Package connection_editor contains the model- and logic-layer behaviour of the
// Zone Connection Editor: guard-strength presets, zone tier classification,
// connection cloning, and graph-validation helpers. The visual canvas itself is
// rendered by the GUI layer; everything testable lives here.
package connection_editor

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
)

// NewDefaultConnection builds a user-added Direct connection between two zones,
// seeded with the tier's generator-default guard value and the standard (15%)
// weekly increment.
func NewDefaultConnection(from, to string, zones []entities.Zone, playerZoneNames map[string]bool) entities.Connection {
	quality := zone_helpers.GetZoneConnectionGuardQuality(
		from, to, zones, linq.FromMap(playerZoneNames).SelectKeys().ToSlice())
	return entities.Connection{
		From:                 from,
		To:                   to,
		ConnectionType:       "Direct",
		GuardValue:           common_connections.GetGuardStrengthForQuality(quality).Default,
		GuardZone:            from,
		GuardMatchGroup:      "rnd_guard_" + helpers.GetZoneLabel(from) + "_" + helpers.GetZoneLabel(to),
		GuardWeeklyIncrement: common_connections.GetGuardWeeklyIncrements().Standard,
		IsUserAdded:          true,
	}
}

// FindIsolatedZones returns the names of zones not referenced by any connection.
func FindIsolatedZones(zones []entities.Zone, connections []entities.Connection) []string {
	var isolated []string
	for _, zone := range zones {
		referenced := false
		for _, connection := range connections {
			if connection.From == zone.Name || connection.To == zone.Name {
				referenced = true
				break
			}
		}
		if !referenced {
			isolated = append(isolated, zone.Name)
		}
	}
	return isolated
}

// ComputeHasErrors reports whether any connection references a zone name that
// does not exist in the zone list.
func ComputeHasErrors(zones []entities.Zone, connections []entities.Connection) bool {
	zoneNames := make(map[string]bool, len(zones))
	for _, zone := range zones {
		zoneNames[zone.Name] = true
	}
	for _, connection := range connections {
		if !zoneNames[connection.From] || !zoneNames[connection.To] {
			return true
		}
	}
	return false
}

// HasDuplicateName reports whether a connection other than current shares its
// (case-insensitive) name. current must be an element of connections; identity is
// compared by pointer to mirror the C# ReferenceEquals check.
func HasDuplicateName(connections []entities.Connection, current *entities.Connection) bool {
	if current == nil || len(current.Name) == 0 {
		return false
	}
	for i := range connections {
		if &connections[i] == current {
			continue
		}
		if strings.EqualFold(connections[i].Name, current.Name) {
			return true
		}
	}
	return false
}
