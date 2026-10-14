package manualEditDeltaService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenRichEditsAreExpandedOnTheSameGeneration_TheZonesAreRebuiltExactly(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	request := captureRequest(configuration, richEdit(configuration))
	edits := newService().Capture(request)

	// Act
	expansion := newService().Expand(expandRequest(configuration, edits))

	// Assert
	assert.Equal(t, request.Zones, expansion.Zones)
}

func TestWhenRichEditsAreExpandedOnTheSameGeneration_TheConnectionsAreRebuiltExactly(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	request := captureRequest(configuration, richEdit(configuration))
	edits := newService().Capture(request)

	// Act
	expansion := newService().Expand(expandRequest(configuration, edits))

	// Assert
	assert.Equal(t, request.Connections, expansion.Connections)
}

func TestWhenNonDeterministicEditsAreExpanded_TheConnectionsAreRebuiltExactly(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	request := captureRequest(configuration, richEdit(configuration))
	request.NonDeterministic = true
	edits := newService().Capture(request)
	expand := expandRequest(configuration, edits)
	expand.NonDeterministic = true

	// Act
	expansion := newService().Expand(expand)

	// Assert
	assert.Equal(t, request.Connections, expansion.Connections)
}

func TestWhenNothingChanged_ReportsTheEditsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	edits := newService().Capture(captureRequest(configuration, richEdit(configuration)))

	// Act
	expansion := newService().Expand(expandRequest(configuration, edits))

	// Assert
	assert.False(t, expansion.Changed)
}

func TestWhenADerivedValueChanged_TheUntouchedFieldFollowsTheFreshZone(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	edits := newService().Capture(captureRequest(configuration, richEdit(configuration)))
	request := expandRequest(configuration, edits)
	request.FreshZones[1].GuardMultiplier = 9.5

	// Act
	expansion := newService().Expand(request)

	// Assert
	assert.InDelta(t, 9.5, zoneNamed(expansion.Zones, goldZoneName).GuardMultiplier, 0)
}

func TestWhenADerivedValueChanged_TheOverriddenFieldKeepsTheEdit(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	edits := newService().Capture(captureRequest(configuration, richEdit(configuration)))
	request := expandRequest(configuration, edits)
	request.FreshZones[1].Size = 3

	// Act
	expansion := newService().Expand(request)

	// Assert
	assert.InDelta(t, 2.5, zoneNamed(expansion.Zones, goldZoneName).Size, 0)
}

func TestWhenAGeneratedZoneIsNoLongerGenerated_WarnsAboutIt(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	edits := newService().Capture(captureRequest(configuration, unchanged))
	request := expandRequest(configuration, edits)
	request.FreshZones = request.FreshZones[:2]

	// Act
	expansion := newService().Expand(request)

	// Assert
	assert.Equal(t, []string{`zone "Neutral-H" is no longer generated; kept as an editor-made zone`},
		expansion.Warnings)
}

func TestWhenAGeneratedZoneIsNoLongerGenerated_RecordsItAsAdded(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	edits := newService().Capture(captureRequest(configuration, unchanged))
	request := expandRequest(configuration, edits)
	request.FreshZones = request.FreshZones[:2]

	// Act
	expansion := newService().Expand(request)

	// Assert
	assert.True(t, zoneEdit(expansion.ManualEdits, silverZoneName).IsAdded())
}

func TestWhenAGeneratedZoneIsNoLongerGenerated_RebuildsItWithOneCastle(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	edits := newService().Capture(captureRequest(configuration, unchanged))
	request := expandRequest(configuration, edits)
	request.FreshZones = request.FreshZones[:2]

	// Act
	expansion := newService().Expand(request)

	// Assert
	assert.Equal(t, 1, test_helpers.NewZoneEditorService().CountZoneCastles(
		zoneNamed(expansion.Zones, silverZoneName)))
}

func TestWhenAnAddedZoneHasNoRecordedCastles_DerivesItWithoutCastles(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	edits := newService().Capture(captureRequest(configuration, withAddedZone(configuration)))
	edits.ManualZones[3].Castles = nil

	// Act
	expansion := newService().Expand(expandRequest(configuration, edits))

	// Assert
	assert.Equal(t, 0, test_helpers.NewZoneEditorService().CountZoneCastles(
		zoneNamed(expansion.Zones, addedZoneName)))
}

func TestWhenNeutralCastlesChange_ADerivedZoneFollowsTheFreshZone(t *testing.T) {
	t.Parallel()
	// Arrange
	request := castleChangeRequest(t)

	// Act
	expansion := newService().Expand(request)

	// Assert
	assert.Equal(t, request.FreshZones[1].MainObjects, zoneNamed(expansion.Zones, goldZoneName).MainObjects)
}

func TestWhenNeutralCastlesChange_AReprofiledZoneGetsTheNewCount(t *testing.T) {
	t.Parallel()
	// Arrange
	request := castleChangeRequest(t)

	// Act
	expansion := newService().Expand(request)

	// Assert
	assert.Equal(t, 3, test_helpers.NewZoneEditorService().CountZoneCastles(
		zoneNamed(expansion.Zones, silverZoneName)))
}

func TestWhenNeutralCastlesChange_AReprofiledZoneRecordsTheNewCount(t *testing.T) {
	t.Parallel()
	// Arrange
	request := castleChangeRequest(t)

	// Act
	expansion := newService().Expand(request)

	// Assert
	assert.Equal(t, new(3), zoneEdit(expansion.ManualEdits, silverZoneName).Castles)
}

func TestWhenNeutralCastlesChange_AnAddedZoneStaysAdded(t *testing.T) {
	t.Parallel()
	// Arrange
	request := castleChangeRequest(t)

	// Act
	expansion := newService().Expand(request)

	// Assert
	assert.True(t, zoneEdit(expansion.ManualEdits, addedZoneName).IsAdded())
}

func TestWhenNeutralCastlesChange_AnAddedZoneGetsTheNewCount(t *testing.T) {
	t.Parallel()
	// Arrange
	request := castleChangeRequest(t)

	// Act
	expansion := newService().Expand(request)

	// Assert
	assert.Equal(t, 3, test_helpers.NewZoneEditorService().CountZoneCastles(
		zoneNamed(expansion.Zones, addedZoneName)))
}

func TestWhenNeutralCastlesChange_ReportsTheEditsChanged(t *testing.T) {
	t.Parallel()
	// Arrange
	request := castleChangeRequest(t)

	// Act
	expansion := newService().Expand(request)

	// Assert
	assert.True(t, expansion.Changed)
}

func TestWhenNeutralCastlesChange_AGeneratedZoneWithEditedMainObjectsGetsTheNewCount(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	edits := newService().Capture(captureRequest(configuration,
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			zones[1].MainObjects = append(zones[1].MainObjects, template_model.MainObject{Type: "Mine"})
			return zones, connections
		}))
	configuration.ZoneConfiguration.NeutralZoneCastles = 3
	request := expandRequest(configuration, edits)
	request.CastleChanges = editor_state_model.CastleSettingChanges{NeutralSimple: true}

	// Act
	expansion := newService().Expand(request)

	// Assert
	assert.Equal(t, 3, test_helpers.NewZoneEditorService().CountZoneCastles(
		zoneNamed(expansion.Zones, goldZoneName)))
}

func TestWhenPlayerCastlesChange_ASpawnZoneWithEditedMainObjectsIsRebuilt(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	edits := newService().Capture(captureRequest(configuration,
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			zones[0].MainObjects = append(zones[0].MainObjects, template_model.MainObject{Type: "City"})
			return zones, connections
		}))
	configuration.ZoneConfiguration.PlayerZoneCastles = 2
	configuration.ZoneConfiguration.PlayerOwnedCastles = 0
	request := expandRequest(configuration, edits)
	request.CastleChanges = editor_state_model.CastleSettingChanges{PlayerCastles: true}

	// Act
	expansion := newService().Expand(request)

	// Assert
	assert.Equal(t, 2, test_helpers.NewZoneEditorService().CountZoneCastles(
		zoneNamed(expansion.Zones, spawnZoneName)))
}

func TestWhenParallelConnectionsAreExpanded_BothAreRebuiltInOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	request := captureRequest(configuration,
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			first := userAddedConnection(silverZoneName, zones, configuration)
			first.Name = "Manual-A-H"
			second := userAddedConnection(silverZoneName, zones, configuration)
			second.Name = "Manual-A-H-2"
			second.GuardValue = 777
			return zones, append(connections, first, second)
		})
	edits := newService().Capture(request)

	// Act
	expansion := newService().Expand(expandRequest(configuration, edits))

	// Assert
	assert.Equal(t, request.Connections, expansion.Connections)
}

// richEdit moves, overrides, clears and re-tiers zones, adds a zone, and
// edits, re-routes and adds connections.
func richEdit(configuration *config.GeneratorConfig) func(
	zones []template_model.Zone, connections []template_model.Connection) (
	[]template_model.Zone, []template_model.Connection) {
	return func(zones []template_model.Zone, connections []template_model.Connection) (
		[]template_model.Zone, []template_model.Connection) {
		zones, connections = withAddedZone(configuration)(zones, connections)
		zones[0].GuardedContentPool = nil
		zones[1].ManualPosition = new(data.NewVec2(0.9, 0.1))
		zones[1].Size = 2.5
		zones[1].GuardReactionDistribution = []int{}
		test_helpers.NewZoneEditorService().ApplyNeutralZoneQuality(
			&zones[2], neutral_zone.QualityLow, 2, newTuning(configuration, len(zones)))
		connections[0].GuardValue = 80000
		connections[1].From, connections[1].To = connections[1].To, connections[1].From
		portal := userAddedConnection(addedZoneName, zones, configuration)
		test_helpers.NewZoneEditorService().ChangeConnectionType(&portal, "Portal", configuration.GenerateRoads)
		portal.Name = "Manual-A-Q"
		return zones, append(connections, portal)
	}
}

func withAddedZone(configuration *config.GeneratorConfig) func(
	zones []template_model.Zone, connections []template_model.Connection) (
	[]template_model.Zone, []template_model.Connection) {
	return func(zones []template_model.Zone, connections []template_model.Connection) (
		[]template_model.Zone, []template_model.Connection) {
		added := factoryZone(addedZoneLabel, neutral_zone.QualityMedium, 1, configuration, len(zones)+1)
		added.ManualPosition = new(data.NewVec2(0.8, 0.2))
		return append(zones, added), connections
	}
}

// castleChangeRequest stores rich edits, then raises the simple-mode neutral
// castle count to 3 and flags that option as changed.
func castleChangeRequest(t *testing.T) models.ManualEditExpandRequest {
	t.Helper()
	configuration := newConfiguration()
	edits := newService().Capture(captureRequest(configuration, richEdit(configuration)))
	configuration.ZoneConfiguration.NeutralZoneCastles = 3
	request := expandRequest(configuration, edits)
	request.CastleChanges = editor_state_model.CastleSettingChanges{NeutralSimple: true}
	return request
}
