package zoneEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_common_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	neutralZoneA = "Neutral-A"
	neutralZoneB = "Neutral-B"
	playerZoneA  = "Spawn-A"
	hubZone      = "Hub"
)

func TestWhenNeutralZoneIsRetiered_MovesEachNamedPresetToTheNewTable(t *testing.T) {
	t.Parallel()
	// Bronze -> Gold, one case per named preset of the ordered table.
	cases := []struct {
		name     string
		bronze   int
		expected int
	}{
		{"Default", 15000, 25000},
		{"Weakest", 3000, 36000},
		{"Low", 6000, 42000},
		{"Medium", 9000, 48000},
		{"High", 12000, 54000},
		{"VeryHigh", 16000, 60000},
	}
	for _, testCase := range cases {
		t.Run("WhenGuardIs"+testCase.name+"_BecomesTheNewTiersSamePreset", func(t *testing.T) {
			t.Parallel()
			// Arrange
			zones, connections := playerEdgeGraph(neutral_zone.QualityLow, testCase.bronze)

			// Act
			_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)

			// Assert
			assert.Equal(t, testCase.expected, guardOf(t, result, "edge"))
		})
	}
}

// Default is a named tier, not a fixed number: it must follow the table just
// like the other presets do.
func TestWhenGuardIsBronzeDefault_BecomesGoldDefault(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := playerEdgeGraph(neutral_zone.QualityLow, 15000)

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(t, 25000, guardOf(t, result, "edge"))
}

func TestWhenNeutralZoneIsDowngraded_MovesTheGuardToTheWeakerTable(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := playerEdgeGraph(neutral_zone.QualityHigh, 48000)

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityLow, 0)

	// Assert
	assert.Equal(t, 9000, guardOf(t, result, "edge"))
}

// GuardZone decides placement only; the stronger endpoint decides the table.
func TestWhenTheWeakerEndpointRises_TakesTheRisenEndpointsTable(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		neutralZone("A", neutral_zone.QualityLow),
		neutralZone("B", neutral_zone.QualityMedium),
	}
	connections := []template_model.Connection{
		{Name: "edge", From: neutralZoneA, To: neutralZoneB, GuardValue: 24000, GuardZone: neutralZoneB},
	}

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(t, 48000, guardOf(t, result, "edge"))
}

func TestWhenTheOppositeEndpointStaysStronger_LeavesTheGuardUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		neutralZone("A", neutral_zone.QualityLow),
		neutralZone("B", neutral_zone.QualityHigh),
	}
	connections := []template_model.Connection{
		{Name: "edge", From: neutralZoneA, To: neutralZoneB, GuardValue: 48000},
	}

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityLowest, 0)

	// Assert
	assert.Equal(t, 48000, guardOf(t, result, "edge"))
}

func TestWhenBothEndpointsShareTheTier_StillFollowsTheEditedEndpoint(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		neutralZone("A", neutral_zone.QualityLow),
		neutralZone("B", neutral_zone.QualityLow),
	}
	connections := []template_model.Connection{
		{Name: "edge", From: neutralZoneA, To: neutralZoneB, GuardValue: 9000},
	}

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(t, 48000, guardOf(t, result, "edge"))
}

func TestWhenTheEditedZoneIsTheFromEndpoint_StillRemapsTheGuard(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{neutralZone("A", neutral_zone.QualityLow), {Name: playerZoneA}}
	connections := []template_model.Connection{
		{Name: "edge", From: neutralZoneA, To: playerZoneA, GuardValue: 15000},
	}

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(t, 25000, guardOf(t, result, "edge"))
}

// A hub endpoint resolves to the strongest tier, so a neutral edit below it
// cannot move the guard.
func TestWhenTheOppositeEndpointIsTheHub_LeavesTheGuardOnTheHubTable(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{neutralZone("A", neutral_zone.QualityLow), {Name: hubZone}}
	connections := []template_model.Connection{
		{Name: "edge", From: neutralZoneA, To: hubZone, GuardValue: 62000},
	}

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(t, 62000, guardOf(t, result, "edge"))
}

func TestWhenTheOppositeEndpointHasAnUnknownName_LetsTheNeutralEndpointDecide(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{neutralZone("A", neutral_zone.QualityLow), {Name: "Oddity"}}
	connections := []template_model.Connection{
		{Name: "edge", From: neutralZoneA, To: "Oddity", GuardValue: 15000},
	}

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(t, 25000, guardOf(t, result, "edge"))
}

// A name no zone carries falls back to the Bronze tier, which outranks a
// downgraded Plastic endpoint.
func TestWhenTheOppositeEndpointIsMissing_UsesTheBronzeFallbackTier(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{neutralZone("A", neutral_zone.QualityHigh)}
	connections := []template_model.Connection{
		{Name: "edge", From: neutralZoneA, To: "Nowhere", GuardValue: 48000},
	}

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityLowest, 0)

	// Assert
	assert.Equal(t, 9000, guardOf(t, result, "edge"))
}

// A zone loaded from a raw .rmg.json has no recorded tier; the inferred one
// decides both whether the tier changed and which table the guard came from.
func TestWhenTheEditedZoneHasNoRecordedTier_InfersItFromTheZoneContent(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := neutralZone("A", neutral_zone.QualityLow)
	zone.Quality = nil
	zones := []template_model.Zone{zone, {Name: playerZoneA}}
	connections := []template_model.Connection{
		{Name: "edge", From: playerZoneA, To: neutralZoneA, GuardValue: 15000},
	}

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(t, 25000, guardOf(t, result, "edge"))
}

// Plastic shares Bronze's preset table, so a Bronze Medium guard is recognised
// on a Plastic zone.
func TestWhenThePlasticZoneIsRetiered_ReadsItsGuardFromTheBronzeTable(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := playerEdgeGraph(neutral_zone.QualityLowest, 9000)

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityMedium, 0)

	// Assert
	assert.Equal(t, 24000, guardOf(t, result, "edge"))
}

// A generated Plastic guard of 10,000 is absent from the Bronze table it is
// matched against, so it stays a custom number. Batch E does not fix that
// table discrepancy.
func TestWhenThePlasticGuardIsTenThousand_KeepsItAsACustomValue(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := playerEdgeGraph(neutral_zone.QualityLowest, 10000)

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(t, 10000, guardOf(t, result, "edge"))
}

func TestWhenTheGuardMatchesNoPreset_KeepsTheCustomNumber(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		guard int
	}{
		{"Negative", -5000},
		{"Zero", 0},
		{"Arbitrary", 12345},
	}
	for _, testCase := range cases {
		t.Run("WhenGuardIs"+testCase.name+"_StaysUnchanged", func(t *testing.T) {
			t.Parallel()
			// Arrange
			zones, connections := playerEdgeGraph(neutral_zone.QualityLow, testCase.guard)

			// Act
			_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)

			// Assert
			assert.Equal(t, testCase.guard, guardOf(t, result, "edge"))
		})
	}
}

func TestWhenTheRequestedTierIsTheCurrentOne_LeavesIncidentGuardsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := playerEdgeGraph(neutral_zone.QualityLow, 15000)

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityLow, 3)

	// Assert
	assert.Equal(t, 15000, guardOf(t, result, "edge"))
}

func TestWhenOnlyTheCastleCountChanges_StillRebuildsTheZoneCastles(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := playerEdgeGraph(neutral_zone.QualityLow, 15000)

	// Act
	result, _ := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityLow, 3)

	// Assert
	assert.Equal(t, 3, test_helpers.NewZoneEditorService().CountZoneCastles(result[0]))
}

func TestWhenTheTierChanges_RecordsItOnTheEditedZone(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := playerEdgeGraph(neutral_zone.QualityLow, 15000)

	// Act
	result, _ := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(t, neutral_zone.QualityHigh, *result[0].Quality)
}

func TestWhenParallelEdgesShareEndpoints_RemapsEachOneIndependently(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{neutralZone("A", neutral_zone.QualityLow), {Name: playerZoneA}}
	connections := []template_model.Connection{
		{Name: "first", From: neutralZoneA, To: playerZoneA, GuardValue: 15000},
		{Name: "second", From: neutralZoneA, To: playerZoneA, GuardValue: 9000},
	}

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(t, []int{25000, 48000}, []int{result[0].GuardValue, result[1].GuardValue})
}

func TestWhenTheEdgeIsASelfLoop_RemapsItsGuardToo(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{neutralZone("A", neutral_zone.QualityLow)}
	connections := []template_model.Connection{
		{Name: "loop", From: neutralZoneA, To: neutralZoneA, GuardValue: 15000},
	}

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(t, 25000, guardOf(t, result, "loop"))
}

func TestWhenAnEdgeDoesNotTouchTheEditedZone_LeavesItsGuardUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		neutralZone("A", neutral_zone.QualityLow),
		neutralZone("B", neutral_zone.QualityLow),
		{Name: playerZoneA},
	}
	connections := []template_model.Connection{
		{Name: "incident", From: neutralZoneA, To: playerZoneA, GuardValue: 15000},
		{Name: "unrelated", From: neutralZoneB, To: playerZoneA, GuardValue: 15000},
	}

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(t, 15000, guardOf(t, result, "unrelated"))
}

func TestWhenTheGuardIsRemapped_ChangesNothingElseOnTheConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	road := true
	rules := placementRules()
	original := template_model.Connection{
		Name:                     "edge",
		From:                     playerZoneA,
		To:                       neutralZoneA,
		ConnectionType:           "Portal",
		Length:                   3.5,
		SimTurnSquad:             true,
		Road:                     &road,
		GuardZone:                neutralZoneA,
		GuardEscape:              true,
		GuardValue:               15000,
		GuardRandomization:       0.25,
		GuardWeeklyIncrement:     1.5,
		GuardMatchGroup:          "group",
		GatePlacement:            "gate",
		PortalPlacementRulesFrom: rules,
		PortalPlacementRulesTo:   rules,
		IsUserAdded:              true,
	}
	zones := []template_model.Zone{neutralZone("A", neutral_zone.QualityLow), {Name: playerZoneA}}
	expected := original.Clone()
	expected.GuardValue = 25000

	// Act
	_, result := editQuality(
		zones, []template_model.Connection{original}, neutralZoneA, neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(t, expected, result[0])
}

func TestWhenIncidentConnectionTypesVary_OnlyThePresetGuardChanges(t *testing.T) {
	t.Parallel()
	for _, connectionType := range []string{"Direct", "Portal", "Arena", "Proximity", "Default", ""} {
		t.Run("WhenTypeIs"+connectionType+"_PreservesOtherFields", func(t *testing.T) {
			t.Parallel()
			// Arrange
			roadTrue, roadFalse := true, false
			zones, _ := playerEdgeGraph(neutral_zone.QualityLow, 15000)
			connections := []template_model.Connection{
				{Name: "nil", From: neutralZoneA, To: playerZoneA, ConnectionType: connectionType, GuardValue: 15000},
				{
					Name: "true", From: neutralZoneA, To: playerZoneA,
					ConnectionType: connectionType, GuardValue: 15000, Road: &roadTrue,
				},
				{
					Name: "false", From: neutralZoneA, To: playerZoneA,
					ConnectionType: connectionType, GuardValue: 15000, Road: &roadFalse,
				},
			}
			expected := make([]template_model.Connection, len(connections))
			for index := range connections {
				expected[index] = connections[index].Clone()
				expected[index].GuardValue = 25000
			}

			// Act
			_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)

			// Assert
			assert.Equal(t, expected, result)
		})
	}
}

func TestWhenTheEditRuns_PreservesConnectionOrderAndCount(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{neutralZone("A", neutral_zone.QualityLow), {Name: playerZoneA}}
	connections := []template_model.Connection{
		{Name: "first", From: neutralZoneA, To: playerZoneA, GuardValue: 15000},
		{Name: "second", From: playerZoneA, To: neutralZoneA, GuardValue: 12345},
		{Name: "third", From: neutralZoneA, To: playerZoneA, GuardValue: 9000},
	}

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(t, []string{"first", "second", "third"}, connectionNames(result))
}

func TestWhenTheEditedZoneIsNotInTheGraph_ReturnsTheInputUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{neutralZone("A", neutral_zone.QualityLow), {Name: playerZoneA}}
	connections := []template_model.Connection{
		{Name: "edge", From: neutralZoneA, To: playerZoneA, GuardValue: 15000},
	}

	// Act
	resultZones, resultConnections := editQuality(
		zones, connections, "Neutral-Z", neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(
		t,
		[]any{zones, connections},
		[]any{resultZones, resultConnections})
}

func TestWhenTheEditedZoneIsNotInTheGraph_StillReturnsClonedZones(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{neutralZone("A", neutral_zone.QualityLow)}

	// Act
	resultZones, _ := editQuality(zones, nil, "Neutral-Z", neutral_zone.QualityHigh, 0)
	resultZones[0].Name = "mutated"

	// Assert
	assert.Equal(t, neutralZoneA, zones[0].Name)
}

func TestWhenTheGraphIsNil_ReturnsNilCollections(t *testing.T) {
	t.Parallel()
	// Act
	resultZones, resultConnections := editQuality(nil, nil, neutralZoneA, neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(t, []any{[]template_model.Zone(nil), []template_model.Connection(nil)},
		[]any{resultZones, resultConnections})
}

func TestWhenThereAreNoConnections_StillReprofilesTheZone(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{neutralZone("A", neutral_zone.QualityLow)}

	// Act
	result, _ := editQuality(zones, []template_model.Connection{}, neutralZoneA, neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(t, neutral_zone.QualityHigh, *result[0].Quality)
}

func TestWhenTheEditRuns_LeavesTheSourceZonesUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := playerEdgeGraph(neutral_zone.QualityLow, 15000)

	// Act
	_, _ = editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 3)

	// Assert
	assert.Equal(t, neutral_zone.QualityLow, *zones[0].Quality)
}

func TestWhenTheEditRuns_LeavesTheSourceGuardValuesUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	zones, connections := playerEdgeGraph(neutral_zone.QualityLow, 15000)

	// Act
	_, _ = editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)

	// Assert
	assert.Equal(t, 15000, connections[0].GuardValue)
}

// The nested slices of the result must not alias the request's.
func TestWhenTheResultZoneIsMutated_LeavesTheSourceContentPoolsUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		neutralZone("A", neutral_zone.QualityLow),
		neutralZone("B", neutral_zone.QualityLow),
	}
	original := append([]string(nil), zones[0].GuardedContentPool...)

	// Act
	result, _ := editQuality(zones, nil, neutralZoneB, neutral_zone.QualityHigh, 0)
	result[0].GuardedContentPool[0] = "overwritten"

	// Assert
	assert.Equal(t, original, zones[0].GuardedContentPool)
}

func TestWhenTheResultConnectionIsMutated_LeavesTheSourcePlacementRulesUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{neutralZone("A", neutral_zone.QualityLow), {Name: playerZoneA}}
	connections := []template_model.Connection{
		{
			Name:                     "edge",
			From:                     neutralZoneA,
			To:                       playerZoneA,
			GuardValue:               15000,
			PortalPlacementRulesFrom: placementRules(),
		},
	}

	// Act
	_, result := editQuality(zones, connections, neutralZoneA, neutral_zone.QualityHigh, 0)
	result[0].PortalPlacementRulesFrom[0].Args[0] = "overwritten"

	// Assert
	assert.Equal(t, []any{"placement"}, connections[0].PortalPlacementRulesFrom[0].Args)
}

// neutralZone builds a tiered neutral zone through the same factory the editor
// uses, so its recorded tier and content pools are the production ones.
func neutralZone(label string, quality neutral_zone.Quality) template_model.Zone {
	return test_helpers.NewZoneEditorService().
		NewDefaultNeutralZone(label, quality, 0, false, defaultTuning())
}

// guardOf returns the guard value of the named connection in the result.
func guardOf(t *testing.T, connections []template_model.Connection, name string) int {
	t.Helper()
	for _, connection := range connections {
		if connection.Name == name {
			return connection.GuardValue
		}
	}
	require.Failf(t, "connection not found", "no connection named %q", name)
	return 0
}

// editQuality runs the graph edit with the production collaborators.
func editQuality(
	zones []template_model.Zone,
	connections []template_model.Connection,
	zoneName string,
	quality neutral_zone.Quality,
	castleCount int) ([]template_model.Zone, []template_model.Connection) {
	return test_helpers.NewZoneEditorService().
		ApplyNeutralZoneQualityEdit(models.NeutralZoneQualityEditRequest{
			Zones:           zones,
			Connections:     connections,
			PlayerZoneNames: []string{playerZoneA},
			ZoneName:        zoneName,
			Quality:         quality,
			CastleCount:     castleCount,
			Tuning:          defaultTuning(),
		})
}

// playerEdgeGraph is a neutral zone of the given tier joined to a spawn zone,
// so the neutral endpoint alone decides the connection's guard table.
func playerEdgeGraph(
	quality neutral_zone.Quality,
	guardValue int) ([]template_model.Zone, []template_model.Connection) {
	zones := []template_model.Zone{neutralZone("A", quality), {Name: playerZoneA}}
	connections := []template_model.Connection{
		{Name: "edge", From: playerZoneA, To: neutralZoneA, GuardValue: guardValue},
	}
	return zones, connections
}

// connectionNames lists the connection names in their current order.
func connectionNames(connections []template_model.Connection) []string {
	names := make([]string, 0, len(connections))
	for _, connection := range connections {
		names = append(names, connection.Name)
	}
	return names
}

// placementRules builds a one-entry portal placement rule list with a nested
// argument slice, so result/source aliasing is observable.
func placementRules() []template_common_model.PlacementRule {
	return template_common_model.ToPlacementRuleModels(
		[]template_entity.PlacementRule{{Args: []any{"placement"}}})
}
