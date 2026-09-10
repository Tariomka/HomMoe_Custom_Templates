package roadPolicyService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Connection road flags ────────────────────────────────────────────

func TestWhenRoadsAreOn_StampsRoadTrueOnEveryNonPortalType(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name           string
		connectionType string
	}{
		{"WhenTypeIsDirect_StampsTrue", "Direct"},
		{"WhenTypeIsDefault_StampsTrue", "Default"},
		{"WhenTypeIsEmpty_StampsTrue", ""},
		{"WhenTypeIsProximity_StampsTrue", "Proximity"},
		{"WhenTypeIsGladiatorArena_StampsTrue", "GladiatorArena"},
		{"WhenTypeIsCustom_StampsTrue", "SomeImportedType"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			connections := []template_model.Connection{
				{Name: "Rnd-A-B", From: "A", To: "B", ConnectionType: testCase.connectionType},
			}

			// Act
			newPolicy().Reconcile(models.RoadReconciliationRequest{
				Connections:   connections,
				GenerateRoads: true,
			})

			// Assert
			assert.Equal(t, new(true), connections[0].Road)
		})
	}
}

func TestWhenRoadsAreOff_StampsRoadFalseOnEveryNonPortalType(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name           string
		connectionType string
	}{
		{"WhenTypeIsDirect_StampsFalse", "Direct"},
		{"WhenTypeIsDefault_StampsFalse", "Default"},
		{"WhenTypeIsEmpty_StampsFalse", ""},
		{"WhenTypeIsProximity_StampsFalse", "Proximity"},
		{"WhenTypeIsGladiatorArena_StampsFalse", "GladiatorArena"},
		{"WhenTypeIsCustom_StampsFalse", "SomeImportedType"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			connections := []template_model.Connection{
				{Name: "Rnd-A-B", From: "A", To: "B", ConnectionType: testCase.connectionType},
			}

			// Act
			newPolicy().Reconcile(models.RoadReconciliationRequest{Connections: connections})

			// Assert
			assert.Equal(t, new(false), connections[0].Road)
		})
	}
}

func TestWhenImportedConnectionSaysNoRoadAndRoadsAreOn_OverwritesItWithTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []template_model.Connection{
		{Name: "Rnd-A-B", From: "A", To: "B", ConnectionType: "Direct", Road: new(false)},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Connections:   connections,
		GenerateRoads: true,
	})

	// Assert
	assert.Equal(t, new(true), connections[0].Road)
}

func TestWhenImportedConnectionSaysRoadAndRoadsAreOff_OverwritesItWithFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []template_model.Connection{
		{Name: "Rnd-A-B", From: "A", To: "B", ConnectionType: "Direct", Road: new(true)},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Connections: connections})

	// Assert
	assert.Equal(t, new(false), connections[0].Road)
}

func TestWhenConnectionIsAnExplicitPortal_LeavesItsOmittedRoadFlagAlone(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []template_model.Connection{
		{Name: "Portal-A-B", From: "A", To: "B", ConnectionType: "portal"},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Connections: connections})

	// Assert
	assert.Nil(t, connections[0].Road)
}

func TestWhenExplicitPortalSaysNoRoad_KeepsThatFlagWithRoadsOn(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []template_model.Connection{
		{Name: "Portal-A-B", From: "A", To: "B", ConnectionType: "Portal", Road: new(false)},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Connections:   connections,
		GenerateRoads: true,
	})

	// Assert
	assert.Equal(t, new(false), connections[0].Road)
}

func TestWhenNonPortalCarriesPortalPlacementRules_StillFollowsTheSetting(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []template_model.Connection{{
		Name:                     "Rnd-A-B",
		From:                     "A",
		To:                       "B",
		ConnectionType:           "Direct",
		PortalPlacementRulesFrom: []template_model.PlacementRule{{Type: "MainObject"}},
	}}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Connections: connections})

	// Assert
	assert.Equal(t, new(false), connections[0].Road)
}

// ── Approach roads ───────────────────────────────────────────────────

func TestWhenRoadsAreOnAndApproachIsMissing_AddsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 1)}
	connections := []template_model.Connection{
		{Name: "Rnd-A-B", From: "A", To: "B", ConnectionType: "Direct"},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:         zones,
		Connections:   connections,
		GenerateRoads: true,
	})

	// Assert
	assert.Equal(t, []string{"Rnd-A-B"}, roadTargets(zones[0], "Connection"))
}

func TestWhenRoadsAreOff_RemovesTheNonPortalApproach(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(mainObjectRef("0"), connectionRef("Rnd-A-B"))),
	}
	connections := []template_model.Connection{
		{Name: "Rnd-A-B", From: "A", To: "B", ConnectionType: "Direct"},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Zones: zones, Connections: connections})

	// Assert
	assert.Empty(t, roadTargets(zones[0], "Connection"))
}

func TestWhenRoadsAreTurnedBackOn_RestoresTheRemovedApproach(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 1)}
	connections := []template_model.Connection{
		{Name: "Rnd-A-B", From: "A", To: "B", ConnectionType: "Direct", Road: new(false)},
	}
	policy := newPolicy()
	policy.Reconcile(models.RoadReconciliationRequest{Zones: zones, Connections: connections})
	require.Empty(t, roadTargets(zones[0], "Connection"))

	// Act
	policy.Reconcile(models.RoadReconciliationRequest{
		Zones:         zones,
		Connections:   connections,
		GenerateRoads: true,
	})

	// Assert
	assert.Equal(t, []string{"Rnd-A-B"}, roadTargets(zones[0], "Connection"))
}

func TestWhenRoadsAreOffAndPortalSaysNoRoad_KeepsThePortalApproach(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(mainObjectRef("0"), connectionRef("Portal-A-B"))),
	}
	connections := []template_model.Connection{
		{Name: "Portal-A-B", From: "A", To: "B", ConnectionType: "Portal", Road: new(false)},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Zones: zones, Connections: connections})

	// Assert
	assert.Equal(t, []string{"Portal-A-B"}, roadTargets(zones[0], "Connection"))
}

func TestWhenRoadsAreOffAndPortalApproachIsMissing_AddsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 1)}
	connections := []template_model.Connection{
		{Name: "Portal-A-B", From: "A", To: "B", ConnectionType: "PORTAL"},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Zones: zones, Connections: connections})

	// Assert
	assert.Equal(t, []string{"Portal-A-B"}, roadTargets(zones[0], "Connection"))
}

func TestWhenConnectorSegmentMixesDirectAndPortal_DropsTheWholeSegment(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{{
		Name:  "Gate-Z",
		Roads: []template_model.Road{road(connectionRef("Rnd-A-Z"), connectionRef("Portal-Z-B"))},
	}}
	connections := []template_model.Connection{
		{Name: "Rnd-A-Z", From: "A", To: "Gate-Z", ConnectionType: "Direct"},
		{Name: "Portal-Z-B", From: "Gate-Z", To: "B", ConnectionType: "Portal"},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Zones: zones, Connections: connections})

	// Assert
	assert.NotContains(t, roadEndpoints(zones[0], "Connection"), "Rnd-A-Z")
}

func TestWhenConnectorSegmentMixesDirectAndPortal_RebuildsThePortalApproachAlone(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{{
		Name:  "Gate-Z",
		Roads: []template_model.Road{road(connectionRef("Rnd-A-Z"), connectionRef("Portal-Z-B"))},
	}}
	connections := []template_model.Connection{
		{Name: "Rnd-A-Z", From: "A", To: "Gate-Z", ConnectionType: "Direct"},
		{Name: "Portal-Z-B", From: "Gate-Z", To: "B", ConnectionType: "Portal"},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Zones: zones, Connections: connections})

	// Assert
	assert.Equal(t, []template_model.Road{
		road(connectionRef("Portal-Z-B"), connectionRef("Portal-Z-B")),
	}, zones[0].Roads)
}

func TestWhenApproachAlreadyExistsInReverseOrder_DoesNotDuplicateIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(connectionRef("Rnd-A-B"), mainObjectRef("0"))),
	}
	connections := []template_model.Connection{
		{Name: "Rnd-A-B", From: "A", To: "B", ConnectionType: "Direct"},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:         zones,
		Connections:   connections,
		GenerateRoads: true,
	})

	// Assert
	assert.Len(t, zones[0].Roads, 1)
}

func TestWhenRoadTargetsAConnectionOfAnotherZone_RemovesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(mainObjectRef("0"), connectionRef("Rnd-B-C"))),
	}
	connections := []template_model.Connection{
		{Name: "Rnd-B-C", From: "B", To: "C", ConnectionType: "Direct"},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:         zones,
		Connections:   connections,
		GenerateRoads: true,
	})

	// Assert
	assert.Empty(t, zones[0].Roads)
}

func TestWhenRoadTargetsAConnectionThatNoLongerExists_RemovesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(mainObjectRef("0"), connectionRef("Deleted-A-B"))),
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Zones: zones, GenerateRoads: true})

	// Assert
	assert.Empty(t, zones[0].Roads)
}

func TestWhenConnectionRefCarriesNoName_KeepsTheOpaqueRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	opaque := road(mainObjectRef("0"), template_model.TypedRef{Type: "Connection"})
	zones := []template_model.Zone{castleZone("A", 1, opaque)}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Zones: zones, GenerateRoads: true})

	// Assert
	assert.Equal(t, []template_model.Road{opaque}, zones[0].Roads)
}

func TestWhenConnectionIsNameless_CreatesNoApproachForIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 1)}
	connections := []template_model.Connection{{From: "A", To: "B", ConnectionType: "Direct"}}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:         zones,
		Connections:   connections,
		GenerateRoads: true,
	})

	// Assert
	assert.Empty(t, zones[0].Roads)
}

func TestWhenConnectionLoopsBackToItsOwnZone_CreatesASingleApproach(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 1)}
	connections := []template_model.Connection{
		{Name: "Rnd-A-A", From: "A", To: "A", ConnectionType: "Direct"},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:         zones,
		Connections:   connections,
		GenerateRoads: true,
	})

	// Assert
	assert.Equal(t, []string{"Rnd-A-A"}, roadTargets(zones[0], "Connection"))
}

func TestWhenZoneHasNoMainObjects_LinksItsConnectionsToEachOther(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{{Name: "Gate-Z"}}
	connections := []template_model.Connection{
		{Name: "Rnd-A-Z", From: "A", To: "Gate-Z", ConnectionType: "Direct"},
		{Name: "Rnd-Z-B", From: "Gate-Z", To: "B", ConnectionType: "Direct"},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:         zones,
		Connections:   connections,
		GenerateRoads: true,
	})

	// Assert
	assert.Equal(t, []template_model.Road{
		road(connectionRef("Rnd-A-Z"), connectionRef("Rnd-Z-B")),
	}, zones[0].Roads)
}

// ── Main-object anchors ──────────────────────────────────────────────

func TestWhenRoadIsAnchoredOnAMissingMainObject_RemovesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(mainObjectRef("2"), contentRef("custom_item"))),
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:            zones,
		MandatoryContent: []template_model.MandatoryContent{contentGroup("A", "custom_item")},
		GenerateRoads:    true,
	})

	// Assert
	assert.Empty(t, zones[0].Roads)
}

func TestWhenMainObjectAnchorIsNotAnIndex_KeepsTheRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	opaque := road(mainObjectRef("named_castle"), contentRef("custom_item"))
	zones := []template_model.Zone{castleZone("A", 1, opaque)}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:            zones,
		MandatoryContent: []template_model.MandatoryContent{contentGroup("A", "custom_item")},
		GenerateRoads:    true,
	})

	// Assert
	assert.Equal(t, []template_model.Road{opaque}, zones[0].Roads)
}

func TestWhenZoneHasSeveralMainObjects_AddsNoCastleRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 3)}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Zones: zones, GenerateRoads: true})

	// Assert
	assert.Empty(t, zones[0].Roads)
}

func TestWhenZoneHasAValidCustomCastleRoad_KeepsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	custom := template_model.Road{
		Type:       "Stone",
		From:       mainObjectRef("1"),
		To:         mainObjectRef("2"),
		GuardValue: 500,
	}
	zones := []template_model.Zone{castleZone("A", 3, custom)}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Zones: zones, GenerateRoads: true})

	// Assert
	assert.Equal(t, []template_model.Road{custom}, zones[0].Roads)
}

func TestWhenRoadUsesAnUnmodelledReferenceType_KeepsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	crossroads := road(template_model.TypedRef{Type: "Crossroads"}, contentRef("custom_item"))
	zones := []template_model.Zone{castleZone("A", 1, crossroads)}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:            zones,
		MandatoryContent: []template_model.MandatoryContent{contentGroup("A", "custom_item")},
		GenerateRoads:    true,
	})

	// Assert
	assert.Equal(t, []template_model.Road{crossroads}, zones[0].Roads)
}

// ── Arena marker ─────────────────────────────────────────────────────

// A castleless hub still hosts the arena, so the marker is the zone's only
// main object. It is a win-condition prop, not something a road may lead to.

func TestWhenZoneHoldsOnlyTheArenaMarker_LinksItsConnectionsToEachOther(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{arenaZone()}
	connections := []template_model.Connection{
		{Name: "Rnd-Hub-A", From: "Hub", To: "A", ConnectionType: "Direct"},
		{Name: "Rnd-Hub-B", From: "Hub", To: "B", ConnectionType: "Direct"},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:         zones,
		Connections:   connections,
		GenerateRoads: true,
	})

	// Assert
	assert.Equal(t, []template_model.Road{
		road(connectionRef("Rnd-Hub-A"), connectionRef("Rnd-Hub-B")),
	}, zones[0].Roads)
}

func TestWhenZoneHoldsOnlyTheArenaMarker_AnchorsNothingOnAMainObject(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{arenaZone()}
	connections := []template_model.Connection{
		{Name: "Rnd-Hub-A", From: "Hub", To: "A", ConnectionType: "Direct"},
		{Name: "Rnd-Hub-B", From: "Hub", To: "B", ConnectionType: "Direct"},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:         zones,
		Connections:   connections,
		GenerateRoads: true,
	})

	// Assert
	assert.Empty(t, roadEndpoints(zones[0], "MainObject"))
}

func TestWhenRoadsAreOffAndTheMarkerOnlyZoneHasAPortal_KeepsThePortalApproach(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{arenaZone()}
	connections := []template_model.Connection{
		{Name: "Rnd-Hub-A", From: "Hub", To: "A", ConnectionType: "Direct"},
		{Name: "Portal-Hub-B", From: "Hub", To: "B", ConnectionType: "Portal"},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Zones: zones, Connections: connections})

	// Assert
	assert.Equal(t, []template_model.Road{
		road(connectionRef("Portal-Hub-B"), connectionRef("Portal-Hub-B")),
	}, zones[0].Roads)
}

func TestWhenZoneHoldsOnlyTheArenaMarker_AddsNoFootholdRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{arenaZone()}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:                zones,
		MandatoryContent:     []template_model.MandatoryContent{contentGroup("Hub", "name_remote_foothold_1")},
		GenerateRoads:        true,
		SpawnRemoteFootholds: true,
		RemoteFootholdCount:  1,
	})

	// Assert
	assert.Empty(t, roadTargets(zones[0], "MandatoryContent"))
}

func TestWhenTheArenaMarkerFollowsTheCastle_AnchorsTheApproachOnTheCastle(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("Hub", 1)}
	zones[0].MainObjects = append(zones[0].MainObjects, arenaMarker())
	connections := []template_model.Connection{
		{Name: "Rnd-Hub-A", From: "Hub", To: "A", ConnectionType: "Direct"},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:         zones,
		Connections:   connections,
		GenerateRoads: true,
	})

	// Assert
	assert.Equal(t, []string{"0"}, roadEndpoints(zones[0], "MainObject"))
}

// An imported template is free to list the marker first, so the anchor is the
// first main object that is not the marker rather than index 0.
func TestWhenTheArenaMarkerPrecedesTheCastle_AnchorsTheApproachOnTheCastle(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("Hub", 1)}
	zones[0].MainObjects = append([]template_model.MainObject{arenaMarker()}, zones[0].MainObjects...)
	connections := []template_model.Connection{
		{Name: "Rnd-Hub-A", From: "Hub", To: "A", ConnectionType: "Direct"},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:         zones,
		Connections:   connections,
		GenerateRoads: true,
	})

	// Assert
	assert.Equal(t, []string{"1"}, roadEndpoints(zones[0], "MainObject"))
}

func TestWhenTheArenaMarkerPrecedesTheCastle_AnchorsTheFootholdOnTheCastle(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("Hub", 1)}
	zones[0].MainObjects = append([]template_model.MainObject{arenaMarker()}, zones[0].MainObjects...)

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:                zones,
		MandatoryContent:     []template_model.MandatoryContent{contentGroup("Hub", "name_remote_foothold_1")},
		GenerateRoads:        true,
		SpawnRemoteFootholds: true,
		RemoteFootholdCount:  1,
	})

	// Assert
	assert.Equal(t, []string{"1"}, roadEndpoints(zones[0], "MainObject"))
}

// ── Mandatory-content targets ────────────────────────────────────────

func TestWhenFinalContentIsNil_RemovesTheRoadToANamedContentItem(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(mainObjectRef("0"), contentRef("name_remote_foothold_1"))),
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Zones: zones, GenerateRoads: true})

	// Assert
	assert.Empty(t, zones[0].Roads)
}

func TestWhenFinalContentIsNil_RemovesTheRoadFromANamedContentItem(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(contentRef("name_remote_foothold_1"), mainObjectRef("0"))),
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Zones: zones, GenerateRoads: true})

	// Assert
	assert.Empty(t, zones[0].Roads)
}

func TestWhenFinalContentIsEmpty_RemovesTheRoadToANamedContentItem(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(mainObjectRef("0"), contentRef("name_remote_foothold_1"))),
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:            zones,
		MandatoryContent: []template_model.MandatoryContent{},
		GenerateRoads:    true,
	})

	// Assert
	assert.Empty(t, zones[0].Roads)
}

func TestWhenFinalContentIsEmpty_RemovesTheRoadFromANamedContentItem(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(contentRef("name_remote_foothold_1"), mainObjectRef("0"))),
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:            zones,
		MandatoryContent: []template_model.MandatoryContent{},
		GenerateRoads:    true,
	})

	// Assert
	assert.Empty(t, zones[0].Roads)
}

func TestWhenFinalContentIsNilAndContentRefIsUnnamed_KeepsTheOpaqueRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	opaque := road(mainObjectRef("0"), template_model.TypedRef{Type: "MandatoryContent"})
	zones := []template_model.Zone{castleZone("A", 1, opaque)}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Zones: zones, GenerateRoads: true})

	// Assert
	assert.Equal(t, []template_model.Road{opaque}, zones[0].Roads)
}

func TestWhenFootholdsAreDisabled_RemovesTheFootholdRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(mainObjectRef("0"), contentRef("name_remote_foothold_1"))),
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:            zones,
		MandatoryContent: []template_model.MandatoryContent{contentGroup("A", "some_other_item")},
		GenerateRoads:    true,
	})

	// Assert
	assert.Empty(t, zones[0].Roads)
}

func TestWhenFootholdCountShrank_RemovesTheRoadToTheDroppedFoothold(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 1,
		road(mainObjectRef("0"), contentRef("name_remote_foothold_1")),
		road(mainObjectRef("0"), contentRef("name_remote_foothold_2")))}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:                zones,
		MandatoryContent:     []template_model.MandatoryContent{contentGroup("A", "name_remote_foothold_1")},
		GenerateRoads:        true,
		SpawnRemoteFootholds: true,
		RemoteFootholdCount:  1,
	})

	// Assert
	assert.Equal(t, []string{"name_remote_foothold_1"}, roadTargets(zones[0], "MandatoryContent"))
}

func TestWhenFootholdCountGrew_AddsTheRoadToTheNewFoothold(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(mainObjectRef("0"), contentRef("name_remote_foothold_1"))),
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones: zones,
		MandatoryContent: []template_model.MandatoryContent{
			contentGroup("A", "name_remote_foothold_1", "name_remote_foothold_2"),
		},
		GenerateRoads:        true,
		SpawnRemoteFootholds: true,
		RemoteFootholdCount:  2,
	})

	// Assert
	assert.Equal(t,
		[]string{"name_remote_foothold_1", "name_remote_foothold_2"},
		roadTargets(zones[0], "MandatoryContent"))
}

func TestWhenFootholdsAreDisabledInSettings_AddsNoFootholdRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 1)}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:               zones,
		MandatoryContent:    []template_model.MandatoryContent{contentGroup("A", "name_remote_foothold_1")},
		GenerateRoads:       true,
		RemoteFootholdCount: 1,
	})

	// Assert
	assert.Empty(t, zones[0].Roads)
}

func TestWhenFootholdItemIsMissingFromTheZoneContent_AddsNoRoadToIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 1)}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:                zones,
		MandatoryContent:     []template_model.MandatoryContent{contentGroup("A", "name_remote_foothold_1")},
		GenerateRoads:        true,
		SpawnRemoteFootholds: true,
		RemoteFootholdCount:  2,
	})

	// Assert
	assert.Equal(t, []string{"name_remote_foothold_1"}, roadTargets(zones[0], "MandatoryContent"))
}

func TestWhenZoneHasNoMainObjects_AddsNoFootholdRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{{
		Name:             "Gate-Z",
		MandatoryContent: template_model.StringList{"mandatory_content_Gate-Z"},
	}}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:                zones,
		MandatoryContent:     []template_model.MandatoryContent{contentGroup("Gate-Z", "name_remote_foothold_1")},
		GenerateRoads:        true,
		SpawnRemoteFootholds: true,
		RemoteFootholdCount:  1,
	})

	// Assert
	assert.Empty(t, zones[0].Roads)
}

func TestWhenContentItemBelongsToAnotherZonesGroup_RemovesTheRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(mainObjectRef("0"), contentRef("name_remote_foothold_1"))),
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:            zones,
		MandatoryContent: []template_model.MandatoryContent{contentGroup("B", "name_remote_foothold_1")},
		GenerateRoads:    true,
	})

	// Assert
	assert.Empty(t, zones[0].Roads)
}

func TestWhenCustomRoadTargetsAnExistingContentItem_KeepsItWithItsAttributes(t *testing.T) {
	t.Parallel()
	// Arrange
	custom := template_model.Road{
		Type:         "Stone",
		From:         mainObjectRef("0"),
		To:           contentRef("name_remote_foothold_custom"),
		Road:         new(false),
		SimTurnSquad: true,
	}
	zones := []template_model.Zone{castleZone("A", 1, custom)}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:            zones,
		MandatoryContent: []template_model.MandatoryContent{contentGroup("A", "name_remote_foothold_custom")},
		GenerateRoads:    true,
	})

	// Assert
	assert.Equal(t, []template_model.Road{custom}, zones[0].Roads)
}

func TestWhenZoneAlreadyHasItsFootholdRoad_DoesNotDuplicateIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(contentRef("name_remote_foothold_1"), mainObjectRef("0"))),
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:                zones,
		MandatoryContent:     []template_model.MandatoryContent{contentGroup("A", "name_remote_foothold_1")},
		GenerateRoads:        true,
		SpawnRemoteFootholds: true,
		RemoteFootholdCount:  1,
	})

	// Assert
	assert.Len(t, zones[0].Roads, 1)
}

// ── Preservation ─────────────────────────────────────────────────────

func TestWhenZoneRoadsSurvive_KeepsTheirOriginalOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	first := road(mainObjectRef("0"), contentRef("custom_item"))
	second := template_model.Road{Type: "Stone", From: mainObjectRef("0"), To: mainObjectRef("1")}
	zones := []template_model.Zone{castleZone("A", 2, first, second)}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:            zones,
		MandatoryContent: []template_model.MandatoryContent{contentGroup("A", "custom_item")},
		GenerateRoads:    true,
	})

	// Assert
	assert.Equal(t, []template_model.Road{first, second}, zones[0].Roads)
}

func TestWhenZoneHasNoRoadsAndNothingToAdd_LeavesTheListNil(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 1)}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Zones: zones, GenerateRoads: true})

	// Assert
	assert.Nil(t, zones[0].Roads)
}

// A declared but empty list is the author's, so nothing to judge means nothing
// to rewrite.
func TestWhenZoneRoadsAreEmptyAndNothingToAdd_LeavesTheListEmptyNotNil(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 1)}
	zones[0].Roads = []template_model.Road{}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Zones: zones, GenerateRoads: true})

	// Assert
	assert.NotNil(t, zones[0].Roads)
}

// Dropping every road leaves the list nil rather than an empty one: the zone
// no longer declares roads at all, so the saved template omits the field.
func TestWhenEveryRoadIsDropped_LeavesTheListNil(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		castleZone("A", 1, road(mainObjectRef("0"), connectionRef("Deleted-A-B"))),
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{Zones: zones, GenerateRoads: true})

	// Assert
	assert.Nil(t, zones[0].Roads)
}

func TestWhenSeveralZonesShareAConnection_GivesEachOfThemAnApproach(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{castleZone("A", 1), castleZone("B", 1)}
	connections := []template_model.Connection{
		{Name: "Rnd-A-B", From: "A", To: "B", ConnectionType: "Direct"},
	}

	// Act
	newPolicy().Reconcile(models.RoadReconciliationRequest{
		Zones:         zones,
		Connections:   connections,
		GenerateRoads: true,
	})

	// Assert
	assert.Equal(t, []string{"Rnd-A-B"}, roadTargets(zones[1], "Connection"))
}
