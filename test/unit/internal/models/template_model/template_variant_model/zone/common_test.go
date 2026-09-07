package zone_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
)

// newPopulatedZoneEntity mirrors newPopulatedZone on the wire-format side. Only
// the nested converters matter here: encounter holes hang off a pointer, and
// main objects and roads each reach a converter of their own.
func newPopulatedZoneEntity() template_entity.Zone {
	return template_entity.Zone{
		Name: "zone",
		EncounterHolesSettings: &template_entity.EncounterHolesSettings{
			AffectedEncounters: 1,
			TwoHoleEncounters:  2,
		},
		MainObjects: []template_entity.MainObject{{
			Type:          "City",
			Factions:      []string{"faction"},
			PlacementArgs: []string{"arg"},
			Faction:       &template_entity.TypedRef{Type: "FromList", Args: []string{"factionArg"}},
		}},
		Roads: []template_entity.Road{{
			Type: "road",
			From: template_entity.TypedRef{Type: "from", Args: []string{"fromArg"}},
			To:   template_entity.TypedRef{Type: "to", Args: []string{"toArg"}},
			Road: new(true),
		}},
	}
}
