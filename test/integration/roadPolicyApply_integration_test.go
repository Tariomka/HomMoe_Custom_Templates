package integration_test

// Road policy end to end: a real generation on the application's object graph,
// real manual Applies through the production handler, and the resulting
// template read back through the production mapper and JSON codec - the same
// bytes a save would write, without writing anything.

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/topology"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Generation ───────────────────────────────────────────────────────

func TestWhenGeneratingWithRoads_EveryRoadOfTheSavedTemplateResolves(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = true
	})

	// Act
	saved := serializedTemplate(t, state)

	// Assert
	assert.Empty(t, unresolvableRoads(saved))
}

func TestWhenGeneratingWithoutRoads_EveryRoadOfTheSavedTemplateResolves(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = false
	})

	// Act
	saved := serializedTemplate(t, state)

	// Assert
	assert.Empty(t, unresolvableRoads(saved))
}

func TestWhenGeneratingWithRoads_TheSavedConnectionsCarryRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = true
	})

	// Act
	saved := serializedTemplate(t, state)

	// Assert
	assert.Empty(t, savedConnectionsFailingFlag(saved, true))
}

func TestWhenGeneratingWithoutRoads_TheSavedConnectionsAreRoadless(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = false
	})

	// Act
	saved := serializedTemplate(t, state)

	// Assert
	assert.Empty(t, savedConnectionsFailingFlag(saved, false))
}

func TestWhenGeneratingWithoutRoads_TheSavedZonesKeepTheirInternalRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = false
		editorState.PlayerZoneCastles = 2
	})

	// Act
	saved := serializedTemplate(t, state)

	// Assert
	assert.NotEmpty(t, savedInternalRoads(saved))
}

// ── Apply, both directions ───────────────────────────────────────────

func TestWhenRoadsAreSwitchedOffAndApplied_TheNonPortalApproachRoadsAreGone(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = true
	})
	require.NotEmpty(t, savedApproachTargets(serializedTemplate(t, state)))

	// Act
	applyCurrentGraph(t, state, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = false
	})

	// Assert
	assert.Empty(t, savedApproachTargets(serializedTemplate(t, state)))
}

func TestWhenRoadsAreSwitchedOffAndApplied_TheInternalZoneRoadsSurvive(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = true
		editorState.PlayerZoneCastles = 2
	})
	expected := savedInternalRoads(serializedTemplate(t, state))

	// Act
	applyCurrentGraph(t, state, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = false
	})

	// Assert
	assert.Equal(t, expected, savedInternalRoads(serializedTemplate(t, state)))
}

func TestWhenRoadsAreSwitchedBackOnAndApplied_TheApproachRoadsReturn(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = true
	})
	expected := savedApproachTargetSets(serializedTemplate(t, state))
	applyCurrentGraph(t, state, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = false
	})
	require.Empty(t, savedApproachTargets(serializedTemplate(t, state)))

	// Act
	applyCurrentGraph(t, state, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = true
	})

	// Assert
	assert.Equal(t, expected, savedApproachTargetSets(serializedTemplate(t, state)))
}

func TestWhenRoadsAreSwitchedBackOnAndApplied_EveryRoadOfTheSavedTemplateResolves(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = true
	})
	applyCurrentGraph(t, state, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = false
	})

	// Act
	applyCurrentGraph(t, state, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = true
	})

	// Assert
	assert.Empty(t, unresolvableRoads(serializedTemplate(t, state)))
}

// An imported or hand-edited template may carry any connection type and any
// road flag; with roads on, every non-portal one comes back roaded.
func TestWhenRoadsAreOnAndMixedConnectionsAreApplied_EveryNonPortalIsRoaded(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = true
	})
	zones, connections := currentGraph(t, state)
	connections = append(connections, mixedNonPortalConnections(zones)...)

	// Act
	applyGraph(t, state, zones, connections, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = true
	})

	// Assert
	assert.Empty(t, savedConnectionsFailingFlag(serializedTemplate(t, state), true))
}

func TestWhenRoadsAreOffAndMixedConnectionsAreApplied_EveryNonPortalIsRoadless(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = true
	})
	zones, connections := currentGraph(t, state)
	connections = append(connections, mixedNonPortalConnections(zones)...)

	// Act
	applyGraph(t, state, zones, connections, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = false
	})

	// Assert
	assert.Empty(t, savedConnectionsFailingFlag(serializedTemplate(t, state), false))
}

func TestWhenRoadsAreOffAndPortalsAreApplied_ThePortalFlagsArePreserved(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = true
	})
	zones, connections := currentGraph(t, state)
	portals := portalConnections(zones)
	connections = append(connections, portals...)

	// Act
	applyGraph(t, state, zones, connections, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = false
	})

	// Assert
	assert.Equal(t, portalFlagsByName(portals), savedPortalFlags(t, serializedTemplate(t, state)))
}

func TestWhenRoadsAreOnAndPortalsAreApplied_ThePortalFlagsArePreserved(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = true
	})
	zones, connections := currentGraph(t, state)
	portals := portalConnections(zones)
	connections = append(connections, portals...)

	// Act
	applyGraph(t, state, zones, connections, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = true
	})

	// Assert
	assert.Equal(t, portalFlagsByName(portals), savedPortalFlags(t, serializedTemplate(t, state)))
}

// A roadless portal still holds its approaches: the zone-road checkbox never
// owned them.
func TestWhenRoadsAreOffAndARoadlessPortalIsApplied_ItStillGetsAnApproachRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = true
	})
	zones, connections := currentGraph(t, state)
	roadless := false
	connections = append(connections, template_model.Connection{
		Name:           "Portal-Roadless",
		From:           zones[0].Name,
		To:             zones[1].Name,
		ConnectionType: "portal",
		Road:           &roadless,
		IsUserAdded:    true,
	})

	// Act
	applyGraph(t, state, zones, connections, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = false
	})

	// Assert
	assert.Contains(t, savedApproachTargets(serializedTemplate(t, state)), "Portal-Roadless")
}

// ── Remote footholds ─────────────────────────────────────────────────

func TestWhenFootholdsAreDisabledAndApplied_TheFootholdRoadsAreGone(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.SpawnRemoteFootholds = true
		editorState.RemoteFootholdCount = 2
	})
	require.NotEmpty(t, savedFootholdTargets(serializedTemplate(t, state)))

	// Act
	applyCurrentGraph(t, state, func(editorState *editor_state_model.EditorState) {
		editorState.SpawnRemoteFootholds = false
	})

	// Assert
	assert.Empty(t, savedFootholdTargets(serializedTemplate(t, state)))
}

func TestWhenTheFootholdCountDropsAndApplied_TheStaleFootholdRoadIsGone(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.SpawnRemoteFootholds = true
		editorState.RemoteFootholdCount = 2
	})
	require.Contains(t, savedFootholdTargets(serializedTemplate(t, state)), "name_remote_foothold_2")

	// Act
	applyCurrentGraph(t, state, func(editorState *editor_state_model.EditorState) {
		editorState.RemoteFootholdCount = 1
	})

	// Assert
	assert.NotContains(t, savedFootholdTargets(serializedTemplate(t, state)), "name_remote_foothold_2")
}

func TestWhenTheFootholdCountRisesAndApplied_TheNewFootholdRoadAppears(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.SpawnRemoteFootholds = true
		editorState.RemoteFootholdCount = 1
	})

	// Act
	applyCurrentGraph(t, state, func(editorState *editor_state_model.EditorState) {
		editorState.RemoteFootholdCount = 2
	})

	// Assert
	assert.Contains(t, savedFootholdTargets(serializedTemplate(t, state)), "name_remote_foothold_2")
}

// ── Castle edits and revert ──────────────────────────────────────────

func TestWhenACastleIsRemovedAndApplied_TheStaleCastleRoadIsDropped(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.PlayerZoneCastles = 2
	})
	zones, connections := currentGraph(t, state)
	edited := indexOfZoneWithMainObjects(t, zones, 3)
	zones[edited].MainObjects = zones[edited].MainObjects[:len(zones[edited].MainObjects)-1]
	staleIndex := strconv.Itoa(len(zones[edited].MainObjects))

	// Act
	applyGraph(t, state, zones, connections, func(*editor_state_model.EditorState) {})

	// Assert
	assert.NotContains(t,
		savedMainObjectTargets(serializedTemplate(t, state), zones[edited].Name),
		staleIndex)
}

func TestWhenACastleIsRemovedAndApplied_TheRemainingCastleRoadsSurvive(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.PlayerZoneCastles = 2
	})
	zones, connections := currentGraph(t, state)
	edited := indexOfZoneWithMainObjects(t, zones, 3)
	zones[edited].MainObjects = zones[edited].MainObjects[:len(zones[edited].MainObjects)-1]

	// Act
	applyGraph(t, state, zones, connections, func(*editor_state_model.EditorState) {})

	// Assert
	assert.Contains(t,
		savedMainObjectTargets(serializedTemplate(t, state), zones[edited].Name),
		"1")
}

// Reverting to the previewed base must land on the roads a fresh generation
// with the same settings produces.
func TestWhenTheRoadlessBaseIsReapplied_TheSavedRoadsMatchTheBase(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.GenerateRoads = false
	})
	expected := savedRoadSignature(serializedTemplate(t, state))
	base, ok := state.PreviewBaseZones()
	require.True(t, ok)

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{
		Zones:        append([]template_model.Zone(nil), base.Zones...),
		Connections:  append([]template_model.Connection(nil), base.Connections...),
		RevertToBase: true,
	})

	// Assert
	assert.Equal(t, expected, savedRoadSignature(serializedTemplate(t, state)))
}

// ── Gladiator arena ──────────────────────────────────────────────────

// A hub with no castles still hosts the arena, so the marker ends up as the
// zone's only main object. It is a win-condition prop: no road may anchor on
// it, and the hub stays wired like the connector zone it is.

func TestWhenTheCastlelessHubHostsTheArena_ItsRoadsAnchorOnNoMainObject(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newArenaHubSession(t)

	// Act
	hub := savedArenaZone(t, serializedTemplate(t, state))

	// Assert
	assert.Empty(t, referenceArgs(hub.Roads, "MainObject"))
}

func TestWhenTheCastlelessHubHostsTheArena_ItStillLinksItsSpokes(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newArenaHubSession(t)

	// Act
	hub := savedArenaZone(t, serializedTemplate(t, state))

	// Assert
	assert.NotEmpty(t, referenceArgs(hub.Roads, "Connection"))
}

func TestWhenTheCastlelessHubHostsTheArenaAndIsApplied_ItsRoadsStillAnchorOnNoMainObject(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newArenaHubSession(t)

	// Act
	applyCurrentGraph(t, state, func(*editor_state_model.EditorState) {})

	// Assert
	assert.Empty(t, referenceArgs(savedArenaZone(t, serializedTemplate(t, state)).Roads, "MainObject"))
}

// ── Session helpers ──────────────────────────────────────────────────

// newRoadSession generates a template on the application's own object graph
// with the requested generation settings.
func newRoadSession(t *testing.T, settings func(*editor_state_model.EditorState)) *drivers.State {
	t.Helper()
	state := drivers.NewUIState(
		composition.InitializeGuiHandler(),
		composition.InitializeFileSystemHandler(),
		composition.InitializeRegenerationHandler(),
		false)
	state.UpdateState(func(editorState *editor_state_model.EditorState) {
		editorState.PlayerCount = 3
		editorState.NeutralZoneCount = 3
		editorState.RandomPortals = false
		settings(editorState)
	})
	state.Generate()
	require.NotNil(t, state.GetLastTemplate(), "the session must start from a generated template")
	return state
}

// currentGraph returns copies of the live zones and connections, the way the
// editor dialog hands its own slices to Apply.
func currentGraph(t *testing.T, state *drivers.State) ([]template_model.Zone, []template_model.Connection) {
	t.Helper()
	template := state.GetLastTemplate()
	require.NotNil(t, template)
	require.NotEmpty(t, template.Variants)
	return append([]template_model.Zone(nil), template.Variants[0].Zones...),
		append([]template_model.Connection(nil), template.Variants[0].Connections...)
}

func applyCurrentGraph(t *testing.T, state *drivers.State, settings func(*editor_state_model.EditorState)) {
	t.Helper()
	zones, connections := currentGraph(t, state)
	applyGraph(t, state, zones, connections, settings)
}

// newArenaHubSession generates a hub template whose hub holds no castles while
// the gladiator-arena win condition is on, so the generator's own placement
// leaves the hub with the arena marker as its only main object.
func newArenaHubSession(t *testing.T) *drivers.State {
	t.Helper()
	return newRoadSession(t, func(editorState *editor_state_model.EditorState) {
		editorState.Topology = topology.TopologyHubAndSpoke
		editorState.HubZoneCastles = 0
		editorState.GladiatorArena = true
		editorState.GenerateRoads = true
	})
}

func applyGraph(
	t *testing.T,
	state *drivers.State,
	zones []template_model.Zone,
	connections []template_model.Connection,
	settings func(*editor_state_model.EditorState)) {
	t.Helper()
	state.UpdateState(settings)
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: zones, Connections: connections})
}

// serializedTemplate writes the live template through the production mapper and
// JSON codec and reads it back, so the assertions run on the bytes a save would
// produce without writing a file.
func serializedTemplate(t *testing.T, state *drivers.State) template_entity.RmgTemplate {
	t.Helper()
	template := state.GetLastTemplate()
	require.NotNil(t, template)
	encoded, err := json.Marshal(mappers.NewTemplateMapper().ToEntity(*template))
	require.NoError(t, err)

	var decoded template_entity.RmgTemplate
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	return decoded
}

// ── Graph arrangement ────────────────────────────────────────────────

// mixedNonPortalConnections covers the types an imported or hand-edited
// template can carry, with every road-flag state the setting must overwrite.
func mixedNonPortalConnections(zones []template_model.Zone) []template_model.Connection {
	roaded := true
	roadless := false
	types := []string{"Direct", "Default", "", "GladiatorArena", "Proximity", "SomeImportedType"}
	flags := []*bool{&roaded, &roadless, nil}

	var connections []template_model.Connection
	for index, connectionType := range types {
		connections = append(connections, template_model.Connection{
			Name:           fmt.Sprintf("Imported-%d", index),
			From:           zones[index%len(zones)].Name,
			To:             zones[(index+1)%len(zones)].Name,
			ConnectionType: connectionType,
			Road:           flags[index%len(flags)],
			IsUserAdded:    true,
		})
	}
	return connections
}

// portalConnections declares explicit portals in both letter cases with each
// road-flag state, none of which the setting may touch.
func portalConnections(zones []template_model.Zone) []template_model.Connection {
	roaded := true
	roadless := false
	flags := []*bool{&roaded, &roadless, nil}
	types := []string{"Portal", "portal", "Portal"}

	var connections []template_model.Connection
	for index, flag := range flags {
		connections = append(connections, template_model.Connection{
			Name:           fmt.Sprintf("Imported-Portal-%d", index),
			From:           zones[index%len(zones)].Name,
			To:             zones[(index+1)%len(zones)].Name,
			ConnectionType: types[index],
			Road:           flag,
			IsUserAdded:    true,
		})
	}
	return connections
}

func portalFlagsByName(connections []template_model.Connection) map[string]*bool {
	flags := make(map[string]*bool, len(connections))
	for _, connection := range connections {
		flags[connection.Name] = connection.Road
	}
	return flags
}

func indexOfZoneWithMainObjects(t *testing.T, zones []template_model.Zone, minimum int) int {
	t.Helper()
	for index, zone := range zones {
		if len(zone.MainObjects) >= minimum {
			return index
		}
	}
	t.Fatalf("no zone with at least %d main objects was generated", minimum)
	return -1
}

// ── Saved-template inspection ────────────────────────────────────────

// savedConnectionsFailingFlag names every non-portal connection whose saved
// road flag is not the stamped setting.
func savedConnectionsFailingFlag(saved template_entity.RmgTemplate, expected bool) []string {
	var failing []string
	for _, connection := range saved.Variants[0].Connections {
		if strings.EqualFold(connection.ConnectionType, "Portal") {
			continue
		}
		if connection.Road == nil || *connection.Road != expected {
			failing = append(failing, connection.Name)
		}
	}
	return failing
}

// savedPortalFlags returns the road flag every applied portal ended up with.
func savedPortalFlags(t *testing.T, saved template_entity.RmgTemplate) map[string]*bool {
	t.Helper()
	flags := make(map[string]*bool)
	for _, connection := range saved.Variants[0].Connections {
		if strings.HasPrefix(connection.Name, "Imported-Portal-") {
			flags[connection.Name] = connection.Road
		}
	}
	return flags
}

// savedApproachTargets returns the connection every zone road leads out to.
func savedApproachTargets(saved template_entity.RmgTemplate) []string {
	var targets []string
	for _, zone := range saved.Variants[0].Zones {
		targets = append(targets, referenceArgs(zone.Roads, "Connection")...)
	}
	return targets
}

// savedApproachTargetSets groups the connections each zone reaches by road,
// deduplicated: generation can leave a zone several roads to the same
// connection, and restoring roads only owes the zone one of each.
func savedApproachTargetSets(saved template_entity.RmgTemplate) map[string][]string {
	targets := make(map[string][]string, len(saved.Variants[0].Zones))
	for _, zone := range saved.Variants[0].Zones {
		names := referenceArgs(zone.Roads, "Connection")
		if len(names) == 0 {
			continue
		}
		slices.Sort(names)
		targets[zone.Name] = slices.Compact(names)
	}
	return targets
}

func savedFootholdTargets(saved template_entity.RmgTemplate) []string {
	var targets []string
	for _, zone := range saved.Variants[0].Zones {
		targets = append(targets, referenceArgs(zone.Roads, "MandatoryContent")...)
	}
	return targets
}

func savedMainObjectTargets(saved template_entity.RmgTemplate, zoneName string) []string {
	for _, zone := range saved.Variants[0].Zones {
		if zone.Name == zoneName {
			return referenceArgs(zone.Roads, "MainObject")
		}
	}
	return nil
}

// savedArenaZone returns the zone the generator marked with the arena, failing
// the test when it placed none - the regression it guards is only meaningful
// with a real marker present.
func savedArenaZone(t *testing.T, saved template_entity.RmgTemplate) template_entity.Zone {
	t.Helper()
	for _, zone := range saved.Variants[0].Zones {
		for _, mainObject := range zone.MainObjects {
			if mainObject.Type == "GladiatorArena" {
				require.Len(t, zone.MainObjects, 1, "the hub was expected to hold the marker alone")
				return zone
			}
		}
	}
	t.Fatal("the generator placed no gladiator-arena main object")
	return template_entity.Zone{}
}

// savedInternalRoads lists the roads that stay inside their zone.
func savedInternalRoads(saved template_entity.RmgTemplate) []string {
	var internal []string
	for _, zone := range saved.Variants[0].Zones {
		for _, road := range zone.Roads {
			if road.From.Type != "Connection" && road.To.Type != "Connection" {
				internal = append(internal, zone.Name+":"+roadKey(road))
			}
		}
	}
	return internal
}

// savedRoadSignature captures every zone's roads so two templates can be
// compared road by road.
func savedRoadSignature(saved template_entity.RmgTemplate) map[string][]string {
	signature := make(map[string][]string, len(saved.Variants[0].Zones))
	for _, zone := range saved.Variants[0].Zones {
		keys := make([]string, 0, len(zone.Roads))
		for _, road := range zone.Roads {
			keys = append(keys, roadKey(road))
		}
		signature[zone.Name] = keys
	}
	return signature
}

func referenceArgs(roads []template_entity.Road, referenceType string) []string {
	var args []string
	for _, road := range roads {
		for _, reference := range []template_entity.TypedRef{road.From, road.To} {
			if reference.Type == referenceType && len(reference.Args) > 0 {
				args = append(args, reference.Args[0])
			}
		}
	}
	return args
}

func roadKey(road template_entity.Road) string {
	return fmt.Sprintf("%s|%s->%s", road.Type, referenceKey(road.From), referenceKey(road.To))
}

func referenceKey(reference template_entity.TypedRef) string {
	return reference.Type + "(" + strings.Join(reference.Args, ",") + ")"
}

// unresolvableRoads reports every road of the saved template whose endpoint the
// template itself cannot resolve: a main object that is not there, a connection
// that is missing, not incident or roadless, or a mandatory-content item absent
// from the groups the zone references.
func unresolvableRoads(saved template_entity.RmgTemplate) []string {
	variant := saved.Variants[0]
	contentGroups := make(map[string]map[string]bool, len(saved.MandatoryContent))
	for _, group := range saved.MandatoryContent {
		items := make(map[string]bool, len(group.Content))
		for _, item := range group.Content {
			items[item.Name] = true
		}
		contentGroups[group.Name] = items
	}

	var problems []string
	for _, zone := range variant.Zones {
		reachable := reachableConnections(variant, zone.Name)
		for _, road := range zone.Roads {
			for _, reference := range []template_entity.TypedRef{road.From, road.To} {
				problem := describeUnresolvedReference(zone, reference, reachable, contentGroups)
				if problem != "" {
					problems = append(problems, zone.Name+": "+problem)
				}
			}
		}
	}
	return problems
}

// reachableConnections returns the connections a zone may hold a road to: the
// incident ones that either carry roads or are explicit portals.
func reachableConnections(variant template_entity.Variant, zoneName string) map[string]bool {
	reachable := make(map[string]bool)
	for _, connection := range variant.Connections {
		if connection.From != zoneName && connection.To != zoneName {
			continue
		}
		if strings.EqualFold(connection.ConnectionType, "Portal") ||
			(connection.Road != nil && *connection.Road) {
			reachable[connection.Name] = true
		}
	}
	return reachable
}

func describeUnresolvedReference(
	zone template_entity.Zone,
	reference template_entity.TypedRef,
	reachable map[string]bool,
	contentGroups map[string]map[string]bool) string {
	if len(reference.Args) == 0 {
		return ""
	}
	name := reference.Args[0]

	switch reference.Type {
	case "MainObject":
		index, err := strconv.Atoi(name)
		if err == nil && (index < 0 || index >= len(zone.MainObjects)) {
			return fmt.Sprintf("road to main object %s of %d", name, len(zone.MainObjects))
		}
	case "Connection":
		if !reachable[name] {
			return "road to unreachable connection " + name
		}
	case "MandatoryContent":
		for _, groupName := range zone.MandatoryContent {
			if contentGroups[groupName][name] {
				return ""
			}
		}
		return "road to unreachable content item " + name
	}
	return ""
}
