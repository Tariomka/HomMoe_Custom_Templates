package roadPolicyService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

// ── Castle routes ────────────────────────────────────────────────────

func TestWhenZoneGainedCastles_AddsTheMissingCastleRoutes(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 3)}

	// Act
	newPolicy().RebuildZoneConnectionRoads(zones, nil)

	// Assert
	assert.Equal(t, []string{"1", "2"}, roadTargets(zones[0], "MainObject"))
}

func TestWhenCastleCountShrank_DropsTheDanglingCastleRoutes(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 1,
		template_model.Road{Type: "Stone", From: mainObjectRef("0"), To: mainObjectRef("1")},
		template_model.Road{Type: "Stone", From: mainObjectRef("0"), To: mainObjectRef("2")})}

	// Act
	newPolicy().RebuildZoneConnectionRoads(zones, nil)

	// Assert
	assert.Empty(t, zones[0].Roads)
}

// The editor reaches a castleless hub that hosts the arena through this entry
// point, so the marker must not become a road anchor here either.
func TestWhenTheMarkerOnlyZoneIsRebuilt_LinksItsConnectionsToEachOther(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{arenaZone()}
	connections := []template_model.Connection{
		{Name: "Rnd-Hub-A", From: "Hub", To: "A", ConnectionType: "Direct"},
		{Name: "Rnd-Hub-B", From: "Hub", To: "B", ConnectionType: "Direct"},
	}

	// Act
	newPolicy().RebuildZoneConnectionRoads(zones, connections)

	// Assert
	assert.Equal(t, []template_model.Road{
		road(connectionRef("Rnd-Hub-A"), connectionRef("Rnd-Hub-B")),
	}, zones[0].Roads)
}

// ── Unknown content ──────────────────────────────────────────────────

func TestWhenContentIsUnknown_KeepsTheRoadToANamedContentItem(t *testing.T) {
	t.Parallel()
	// Arrange
	foothold := road(mainObjectRef("0"), contentRef("name_remote_foothold_1"))
	zones := []template_model.Zone{castleZone("A", 1, foothold)}

	// Act
	newPolicy().RebuildZoneConnectionRoads(zones, nil)

	// Assert
	assert.Equal(t, []template_model.Road{foothold}, zones[0].Roads)
}

func TestWhenContentIsUnknown_KeepsTheRoadFromANamedContentItem(t *testing.T) {
	t.Parallel()
	// Arrange
	foothold := road(contentRef("name_remote_foothold_1"), mainObjectRef("0"))
	zones := []template_model.Zone{castleZone("A", 1, foothold)}

	// Act
	newPolicy().RebuildZoneConnectionRoads(zones, nil)

	// Assert
	assert.Equal(t, []template_model.Road{foothold}, zones[0].Roads)
}

// ── Connection scope ─────────────────────────────────────────────────

func TestWhenApproachIsMissing_AddsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 1)}
	connections := []template_model.Connection{
		{Name: "Rnd-A-B", From: "A", To: "B", ConnectionType: "Direct"},
	}

	// Act
	newPolicy().RebuildZoneConnectionRoads(zones, connections)

	// Assert
	assert.Equal(t, []string{"Rnd-A-B"}, roadTargets(zones[0], "Connection"))
}

func TestWhenRoadTargetsAConnectionOfAnotherZone_DropsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(mainObjectRef("0"), connectionRef("Rnd-B-C"))),
	}
	connections := []template_model.Connection{
		{Name: "Rnd-B-C", From: "B", To: "C", ConnectionType: "Direct"},
	}

	// Act
	newPolicy().RebuildZoneConnectionRoads(zones, connections)

	// Assert
	assert.Empty(t, zones[0].Roads)
}

func TestWhenRoadTargetsADeletedConnection_DropsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(mainObjectRef("0"), connectionRef("Deleted-A-B"))),
	}

	// Act
	newPolicy().RebuildZoneConnectionRoads(zones, nil)

	// Assert
	assert.Empty(t, zones[0].Roads)
}

func TestWhenTheSameConnectionIsListedTwice_CreatesASingleApproach(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 1)}
	connections := []template_model.Connection{
		{Name: "Rnd-A-B", From: "A", To: "B", ConnectionType: "Direct"},
		{Name: "Rnd-A-B", From: "A", To: "B", ConnectionType: "Direct"},
	}

	// Act
	newPolicy().RebuildZoneConnectionRoads(zones, connections)

	// Assert
	assert.Equal(t, []string{"Rnd-A-B"}, roadTargets(zones[0], "Connection"))
}

func TestWhenConnectionHasEmptyEndpoints_DropsTheRoadThatTargetsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(mainObjectRef("0"), connectionRef("Rnd-A-B"))),
	}
	connections := []template_model.Connection{{Name: "Rnd-A-B", ConnectionType: "Direct"}}

	// Act
	newPolicy().RebuildZoneConnectionRoads(zones, connections)

	// Assert
	assert.Empty(t, zones[0].Roads)
}

func TestWhenConnectionIsNameless_LeavesTheNamingToTheCaller(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 1)}
	connections := []template_model.Connection{{From: "A", To: "B", ConnectionType: "Direct"}}

	// Act
	newPolicy().RebuildZoneConnectionRoads(zones, connections)

	// Assert
	assert.Empty(t, connections[0].Name)
}

func TestWhenConnectionIsNameless_AddsNoApproachForIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 1)}
	connections := []template_model.Connection{{From: "A", To: "B", ConnectionType: "Direct"}}

	// Act
	newPolicy().RebuildZoneConnectionRoads(zones, connections)

	// Assert
	assert.Empty(t, zones[0].Roads)
}

// ── Road flags ───────────────────────────────────────────────────────

func TestWhenConnectionCarriesNoRoadFlag_StampsItTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []template_model.Connection{
		{Name: "Rnd-A-B", From: "A", To: "B", ConnectionType: "Direct"},
	}

	// Act
	newPolicy().RebuildZoneConnectionRoads(nil, connections)

	// Assert
	assert.Equal(t, new(true), connections[0].Road)
}

func TestWhenExplicitPortalSaysNoRoad_LeavesThatFlagAlone(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []template_model.Connection{
		{Name: "Portal-A-B", From: "A", To: "B", ConnectionType: "Portal", Road: new(false)},
	}

	// Act
	newPolicy().RebuildZoneConnectionRoads(nil, connections)

	// Assert
	assert.Equal(t, new(false), connections[0].Road)
}
