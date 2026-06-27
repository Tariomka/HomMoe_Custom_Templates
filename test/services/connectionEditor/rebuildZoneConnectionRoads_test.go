package connectionEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
)

func roadTargets(zone entities.Zone, refType string) map[string]bool {
	out := map[string]bool{}
	for _, road := range zone.Roads {
		if road.To.Type == refType && len(road.To.Args) > 0 {
			out[road.To.Args[0]] = true
		}
		if road.From.Type == refType && len(road.From.Args) > 0 {
			out[road.From.Args[0]] = true
		}
	}
	return out
}

// A zone with a castle (MainObject 0) plus a remote-foothold road and a single
// connection road; after adding a second connection in the editor the rebuild
// must give the zone roads to BOTH connections while keeping the foothold road.
func TestRebuildZoneConnectionRoads_AddsRoadsToAllConnectionsAndKeepsFoothold(t *testing.T) {
	mainObjectZero := entities.TypedRef{Type: "MainObject", Args: []string{"0"}}
	zones := []entities.Zone{
		{
			Name:        "Spawn-A",
			MainObjects: []entities.MainObject{{Type: "Spawn"}, {Type: "City"}},
			Roads: []entities.Road{
				{Type: "Stone", From: mainObjectZero, To: entities.TypedRef{Type: "MainObject", Args: []string{"1"}}},
				{From: mainObjectZero, To: entities.TypedRef{Type: "MandatoryContent", Args: []string{"name_remote_foothold_1"}}},
				{From: mainObjectZero, To: entities.TypedRef{Type: "Connection", Args: []string{"Rnd-A-B"}}},
			},
		},
		{
			Name:        "Spawn-B",
			MainObjects: []entities.MainObject{{Type: "Spawn"}, {Type: "City"}},
			Roads: []entities.Road{
				{Type: "Stone", From: mainObjectZero, To: entities.TypedRef{Type: "MainObject", Args: []string{"1"}}},
				{From: mainObjectZero, To: entities.TypedRef{Type: "MandatoryContent", Args: []string{"name_remote_foothold_1"}}},
				{From: mainObjectZero, To: entities.TypedRef{Type: "Connection", Args: []string{"Rnd-A-B"}}},
			},
		},
	}
	// The editor added a second connection (Rnd-A-C) touching Spawn-A and a new
	// neutral zone, but no road was created for it yet.
	zones = append(zones, entities.Zone{
		Name:        "Neutral-C",
		MainObjects: []entities.MainObject{{Type: "City"}},
		Roads:       []entities.Road{{Type: "Stone", From: mainObjectZero, To: entities.TypedRef{Type: "MainObject", Args: []string{"1"}}}},
	})
	connections := []entities.Connection{
		{Name: "Rnd-A-B", From: "Spawn-A", To: "Spawn-B", ConnectionType: "Direct"},
		{Name: "Rnd-A-C", From: "Spawn-A", To: "Neutral-C", ConnectionType: "Direct"},
	}

	connection_editor.RebuildZoneConnectionRoads(zones, connections)

	spawnA := zones[0]
	connRoads := roadTargets(spawnA, "Connection")
	if !connRoads["Rnd-A-B"] {
		t.Error("Spawn-A lost its road to Rnd-A-B")
	}
	if !connRoads["Rnd-A-C"] {
		t.Error("Spawn-A has no road to the newly added connection Rnd-A-C")
	}
	footholdRoads := roadTargets(spawnA, "MandatoryContent")
	if !footholdRoads["name_remote_foothold_1"] {
		t.Error("Spawn-A lost its remote-foothold road")
	}

	neutralC := zones[2]
	if !roadTargets(neutralC, "Connection")["Rnd-A-C"] {
		t.Error("Neutral-C has no road to Rnd-A-C")
	}
}

// Connections added in the manual zone editor start nameless; the rebuild must
// give them a name and create roads on both endpoint zones.
func TestRebuildZoneConnectionRoads_NamesAndRoadsNamelessManualConnections(t *testing.T) {
	mainObjectZero := entities.TypedRef{Type: "MainObject", Args: []string{"0"}}
	zones := []entities.Zone{
		{
			Name:        "Spawn-E",
			MainObjects: []entities.MainObject{{Type: "Spawn"}, {Type: "City"}},
			Roads:       []entities.Road{{Type: "Stone", From: mainObjectZero, To: entities.TypedRef{Type: "MainObject", Args: []string{"1"}}}},
		},
		{
			Name:        "Neutral-M",
			MainObjects: []entities.MainObject{{Type: "City"}},
			Roads:       []entities.Road{{Type: "Stone", From: mainObjectZero, To: entities.TypedRef{Type: "MainObject", Args: []string{"1"}}}},
		},
	}
	// A nameless, user-added connection, exactly as produced by the editor.
	connections := []entities.Connection{
		{From: "Spawn-E", To: "Neutral-M", ConnectionType: "Direct", IsUserAdded: true},
	}

	connection_editor.RebuildZoneConnectionRoads(zones, connections)

	if connections[0].Name == "" {
		t.Fatal("nameless manual connection was not assigned a name")
	}
	name := connections[0].Name

	if !roadTargets(zones[0], "Connection")[name] {
		t.Errorf("Spawn-E has no road to manual connection %q", name)
	}
	if !roadTargets(zones[1], "Connection")[name] {
		t.Errorf("Neutral-M has no road to manual connection %q", name)
	}
}

// A zone whose castle count grew (e.g. re-tiered in the editor) but which never
// had castle<->castle roads must get fresh stone roads linking each new castle
// to the primary one. Regression test for the missing roads to extra castles.
func TestRebuildZoneConnectionRoads_RegeneratesCastleRoadsForAddedCastles(t *testing.T) {
	mainObjectZero := entities.TypedRef{Type: "MainObject", Args: []string{"0"}}
	zones := []entities.Zone{
		{
			Name:        "Neutral-G",
			MainObjects: []entities.MainObject{{Type: "City"}, {Type: "City"}, {Type: "City"}},
			// Only connection roads exist - no castle roads at all, mirroring a
			// connector zone that had castles added to it.
			Roads: []entities.Road{
				{From: mainObjectZero, To: entities.TypedRef{Type: "Connection", Args: []string{"Rnd-G-H"}}},
			},
		},
	}
	connections := []entities.Connection{
		{Name: "Rnd-G-H", From: "Neutral-G", To: "Neutral-H", ConnectionType: "Direct"},
	}

	connection_editor.RebuildZoneConnectionRoads(zones, connections)

	castleTargets := castleRoadTargets(zones[0])
	if len(castleTargets) != 2 || castleTargets[0] != "1" || castleTargets[1] != "2" {
		t.Errorf("expected stone roads 0->1 and 0->2, got %v", castleTargets)
	}
	if !roadTargets(zones[0], "Connection")["Rnd-G-H"] {
		t.Error("Neutral-G lost its connection road after rebuild")
	}
}

// A zone whose castle count shrank must not keep dangling castle roads that
// point at main objects that no longer exist.
func TestRebuildZoneConnectionRoads_DropsDanglingCastleRoads(t *testing.T) {
	mainObjectZero := entities.TypedRef{Type: "MainObject", Args: []string{"0"}}
	zones := []entities.Zone{
		{
			Name:        "Neutral-G",
			MainObjects: []entities.MainObject{{Type: "City"}},
			Roads: []entities.Road{
				{Type: "Stone", From: mainObjectZero, To: entities.TypedRef{Type: "MainObject", Args: []string{"1"}}},
				{Type: "Stone", From: mainObjectZero, To: entities.TypedRef{Type: "MainObject", Args: []string{"2"}}},
			},
		},
	}

	connection_editor.RebuildZoneConnectionRoads(zones, nil)

	if targets := castleRoadTargets(zones[0]); len(targets) != 0 {
		t.Errorf("single-castle zone should have no castle roads, got %v", targets)
	}
}

// castleRoadTargets returns the MainObject indices linked by the zone's stone
// castle<->castle roads from the primary main object.
func castleRoadTargets(zone entities.Zone) []string {
	var targets []string
	for _, road := range zone.Roads {
		if road.From.Type == "MainObject" && road.To.Type == "MainObject" && len(road.To.Args) > 0 {
			targets = append(targets, road.To.Args[0])
		}
	}
	return targets
}
