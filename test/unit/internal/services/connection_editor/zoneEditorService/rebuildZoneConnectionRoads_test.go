package zoneEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenZoneHasTwoConnections_KeepsRoadToExistingConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := buildFootholdScenario()

	// Act
	test_helpers.NewZoneEditorService().RebuildZoneConnectionRoads(zones, connections)

	// Assert
	assert.True(t, roadTargets(zones[0], "Connection")["Rnd-A-B"],
		"Spawn-A lost its road to Rnd-A-B")
}

func TestWhenConnectionWasAddedInEditor_CreatesRoadForIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := buildFootholdScenario()

	// Act
	test_helpers.NewZoneEditorService().RebuildZoneConnectionRoads(zones, connections)

	// Assert
	assert.True(t, roadTargets(zones[0], "Connection")["Rnd-A-C"],
		"Spawn-A has no road to the newly added connection Rnd-A-C")
}

func TestWhenZoneHasFootholdRoad_KeepsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := buildFootholdScenario()

	// Act
	test_helpers.NewZoneEditorService().RebuildZoneConnectionRoads(zones, connections)

	// Assert
	assert.True(t, roadTargets(zones[0], "MandatoryContent")["name_remote_foothold_1"],
		"Spawn-A lost its remote-foothold road")
}

func TestWhenNewConnectionTouchesAnotherZone_CreatesRoadOnThatZoneToo(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := buildFootholdScenario()

	// Act
	test_helpers.NewZoneEditorService().RebuildZoneConnectionRoads(zones, connections)

	// Assert
	assert.True(t, roadTargets(zones[2], "Connection")["Rnd-A-C"],
		"Neutral-C has no road to Rnd-A-C")
}

func TestWhenManualConnectionIsNameless_AssignsItAName(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := buildNamelessManualScenario()

	// Act
	test_helpers.NewZoneEditorService().RebuildZoneConnectionRoads(zones, connections)

	// Assert
	assert.NotEmpty(t, connections[0].Name,
		"nameless manual connection was not assigned a name")
}

func TestWhenManualConnectionIsNameless_CreatesRoadOnFromZone(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := buildNamelessManualScenario()

	// Act
	test_helpers.NewZoneEditorService().RebuildZoneConnectionRoads(zones, connections)

	// Assert
	require.NotEmpty(t, connections[0].Name)
	assert.True(t, roadTargets(zones[0], "Connection")[connections[0].Name],
		"Spawn-E has no road to the manual connection")
}

func TestWhenManualConnectionIsNameless_CreatesRoadOnToZone(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := buildNamelessManualScenario()

	// Act
	test_helpers.NewZoneEditorService().RebuildZoneConnectionRoads(zones, connections)

	// Assert
	require.NotEmpty(t, connections[0].Name)
	assert.True(t, roadTargets(zones[1], "Connection")[connections[0].Name],
		"Neutral-M has no road to the manual connection")
}

func TestWhenZoneGainedCastlesWithoutCastleRoads_CreatesStoneRoadsToEachCastle(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := buildCastleGrowthScenario()

	// Act
	test_helpers.NewZoneEditorService().RebuildZoneConnectionRoads(zones, connections)

	// Assert
	assert.Equal(t, []string{"1", "2"}, castleRoadTargets(zones[0]),
		"expected stone roads 0->1 and 0->2")
}

func TestWhenCastleRoadsAreRegenerated_KeepsConnectionRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := buildCastleGrowthScenario()

	// Act
	test_helpers.NewZoneEditorService().RebuildZoneConnectionRoads(zones, connections)

	// Assert
	assert.True(t, roadTargets(zones[0], "Connection")["Rnd-G-H"],
		"Neutral-G lost its connection road after rebuild")
}

func TestWhenCastleCountShrank_DropsDanglingCastleRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		{
			Name:        "Neutral-G",
			MainObjects: []template_model.MainObject{{Type: "City"}},
			Roads: []template_model.Road{
				{
					Type: "Stone",
					From: mainObjectZeroRef(),
					To:   template_model.TypedRef{Type: "MainObject", Args: []string{"1"}},
				},
				{
					Type: "Stone",
					From: mainObjectZeroRef(),
					To:   template_model.TypedRef{Type: "MainObject", Args: []string{"2"}},
				},
			},
		},
	}

	// Act
	test_helpers.NewZoneEditorService().RebuildZoneConnectionRoads(zones, nil)

	// Assert
	assert.Empty(t, castleRoadTargets(zones[0]),
		"single-castle zone should have no castle roads")
}

func TestWhenZoneHasNoMainObjects_CreatesConnectorRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		{Name: "Neutral-K"},
		{Name: "Neutral-L", MainObjects: []template_model.MainObject{{Type: "City"}}},
	}
	connections := []entities.Connection{
		{Name: "Rnd-K-L", From: "Neutral-K", To: "Neutral-L", ConnectionType: "Direct"},
	}

	// Act
	test_helpers.NewZoneEditorService().RebuildZoneConnectionRoads(zones, connections)

	// Assert
	assert.True(t, roadTargets(zones[0], "Connection")["Rnd-K-L"],
		"castle-less connector zone must still receive a road for its connection")
}

// buildFootholdScenario returns two spawn zones sharing connection Rnd-A-B plus
// a neutral zone; Spawn-A additionally has a remote-foothold road and the
// editor added a second connection (Rnd-A-C) that has no road yet.
func buildFootholdScenario() ([]template_model.Zone, []entities.Connection) {
	zones := []template_model.Zone{
		{
			Name:        "Spawn-A",
			MainObjects: []template_model.MainObject{{Type: "Spawn"}, {Type: "City"}},
			Roads: []template_model.Road{
				{
					Type: "Stone",
					From: mainObjectZeroRef(),
					To:   template_model.TypedRef{Type: "MainObject", Args: []string{"1"}},
				},
				{
					From: mainObjectZeroRef(),
					To:   template_model.TypedRef{Type: "MandatoryContent", Args: []string{"name_remote_foothold_1"}},
				},
				{From: mainObjectZeroRef(), To: template_model.TypedRef{Type: "Connection", Args: []string{"Rnd-A-B"}}},
			},
		},
		{
			Name:        "Spawn-B",
			MainObjects: []template_model.MainObject{{Type: "Spawn"}, {Type: "City"}},
			Roads: []template_model.Road{
				{
					Type: "Stone",
					From: mainObjectZeroRef(),
					To:   template_model.TypedRef{Type: "MainObject", Args: []string{"1"}},
				},
				{From: mainObjectZeroRef(), To: template_model.TypedRef{Type: "Connection", Args: []string{"Rnd-A-B"}}},
			},
		},
		{
			Name:        "Neutral-C",
			MainObjects: []template_model.MainObject{{Type: "City"}},
			Roads:       nil,
		},
	}
	connections := []entities.Connection{
		{Name: "Rnd-A-B", From: "Spawn-A", To: "Spawn-B", ConnectionType: "Direct"},
		{Name: "Rnd-A-C", From: "Spawn-A", To: "Neutral-C", ConnectionType: "Direct"},
	}
	return zones, connections
}

// buildCastleGrowthScenario returns a three-castle zone that has only a
// connection road (no castle roads at all), mirroring a connector zone that
// had castles added to it in the editor.
func buildCastleGrowthScenario() ([]template_model.Zone, []entities.Connection) {
	zones := []template_model.Zone{
		{
			Name:        "Neutral-G",
			MainObjects: []template_model.MainObject{{Type: "City"}, {Type: "City"}, {Type: "City"}},
			Roads: []template_model.Road{
				{From: mainObjectZeroRef(), To: template_model.TypedRef{Type: "Connection", Args: []string{"Rnd-G-H"}}},
			},
		},
	}
	connections := []entities.Connection{
		{Name: "Rnd-G-H", From: "Neutral-G", To: "Neutral-H", ConnectionType: "Direct"},
	}
	return zones, connections
}

// buildNamelessManualScenario returns two zones joined by a single nameless,
// user-added connection, exactly as produced by the manual zone editor.
func buildNamelessManualScenario() ([]template_model.Zone, []entities.Connection) {
	zones := []template_model.Zone{
		{
			Name:        "Spawn-E",
			MainObjects: []template_model.MainObject{{Type: "Spawn"}, {Type: "City"}},
		},
		{
			Name:        "Neutral-M",
			MainObjects: []template_model.MainObject{{Type: "City"}},
		},
	}
	connections := []entities.Connection{
		{From: "Spawn-E", To: "Neutral-M", ConnectionType: "Direct", IsUserAdded: true},
	}
	return zones, connections
}
