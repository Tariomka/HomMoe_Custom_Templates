// Package roadPolicyService_test holds the shared arrangement helpers of the
// roadPolicyService.go unit tests.
package roadPolicyService_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

func newPolicy() zone_interfaces.IRoadPolicyService {
	return zones.NewRoadPolicyService(zones.NewRoadFactory())
}

func mainObjectRef(index string) template_model.TypedRef {
	return template_model.TypedRef{Type: "MainObject", Args: []string{index}}
}

func connectionRef(name string) template_model.TypedRef {
	return template_model.TypedRef{Type: "Connection", Args: []string{name}}
}

func contentRef(name string) template_model.TypedRef {
	return template_model.TypedRef{Type: "MandatoryContent", Args: []string{name}}
}

func road(from template_model.TypedRef, to template_model.TypedRef) template_model.Road {
	return template_model.Road{From: from, To: to}
}

// roadTargets returns the first argument of every road pointing at the given
// reference type, in road order.
func roadTargets(zone template_model.Zone, referenceType string) []string {
	var targets []string
	for _, item := range zone.Roads {
		if item.To.Type == referenceType && len(item.To.Args) > 0 {
			targets = append(targets, item.To.Args[0])
		}
	}
	return targets
}

// roadEndpoints returns the first argument of every road endpoint - From as
// well as To - pointing at the given reference type, in road order.
func roadEndpoints(zone template_model.Zone, referenceType string) []string {
	var endpoints []string
	for _, item := range zone.Roads {
		for _, reference := range []template_model.TypedRef{item.From, item.To} {
			if reference.Type == referenceType && len(reference.Args) > 0 {
				endpoints = append(endpoints, reference.Args[0])
			}
		}
	}
	return endpoints
}

// castleZone returns a zone with the requested number of main objects and the
// supplied roads, referencing one mandatory-content group.
func castleZone(name string, mainObjectCount int, roads ...template_model.Road) template_model.Zone {
	mainObjects := make([]template_model.MainObject, 0, mainObjectCount)
	for range mainObjectCount {
		mainObjects = append(mainObjects, template_model.MainObject{Type: "City"})
	}
	return template_model.Zone{
		Name:             name,
		MainObjects:      mainObjects,
		MandatoryContent: template_model.StringList{"mandatory_content_" + name},
		Roads:            roads,
	}
}

func contentGroup(zoneName string, itemNames ...string) template_model.MandatoryContent {
	items := make([]template_model.MandatoryContentItem, 0, len(itemNames))
	for _, itemName := range itemNames {
		items = append(items, template_model.MandatoryContentItem{Name: itemName})
	}
	return template_model.MandatoryContent{Name: "mandatory_content_" + zoneName, Content: items}
}

// arenaMarker returns the gladiator-arena main object the generator appends to
// the zone that hosts the win condition.
func arenaMarker() template_model.MainObject {
	return template_model.MainObject{Type: "GladiatorArena"}
}

// arenaZone returns the shape a castleless hub takes once the arena is placed:
// a zone whose only main object is the marker.
func arenaZone() template_model.Zone {
	return template_model.Zone{
		Name:             "Hub",
		MainObjects:      []template_model.MainObject{arenaMarker()},
		MandatoryContent: template_model.StringList{"mandatory_content_Hub"},
	}
}
