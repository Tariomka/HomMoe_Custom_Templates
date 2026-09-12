//go:build integration_test

package integration_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The tier a zone carries has to survive a save/load round trip, because
// nothing on disk can reconstruct it: the .rmg.json schema has no field for it,
// and inferring it back out of the content pools is the guesswork this batch
// replaces.
//
// Plastic is the case worth pinning. Its ordinal is 0, so the persisted field
// has to be a POINTER - a plain int8 with omitempty would drop it from the file
// and the zone would load back as "tier never recorded".

func TestWhenAPlasticZoneIsSaved_ItsTierSurvivesTheLoad(t *testing.T) {
	// Arrange
	path := saveStateWithNeutralTier(t, new(neutral_zone.QualityLowest))

	// Act
	tier := reloadedNeutralTier(t, path)

	// Assert
	require.NotNil(t, tier, "the recorded Plastic tier was dropped on the round trip")
	assert.Equal(t, neutral_zone.QualityLowest, *tier)
}

func TestWhenAGoldZoneIsSaved_ItsTierSurvivesTheLoad(t *testing.T) {
	// Arrange
	path := saveStateWithNeutralTier(t, new(neutral_zone.QualityHigh))

	// Act
	tier := reloadedNeutralTier(t, path)

	// Assert
	require.NotNil(t, tier)
	assert.Equal(t, neutral_zone.QualityHigh, *tier)
}

// A .gen.json written before the tier was persisted carries no quality key, so
// an absent tier has to load as "not recorded" and fall back to inference
// rather than reading back as Plastic.
func TestWhenASavedZoneRecordsNoTier_TheFileCarriesNoQuality(t *testing.T) {
	// Arrange
	path := saveStateWithNeutralTier(t, nil)

	// Act
	raw, err := os.ReadFile(path)

	// Assert
	require.NoError(t, err)
	assert.NotContains(t, string(raw), `"quality"`)
}

func TestWhenASavedZoneRecordsNoTier_ItLoadsAsUnrecorded(t *testing.T) {
	// Arrange
	path := saveStateWithNeutralTier(t, nil)

	// Act
	tier := reloadedNeutralTier(t, path)

	// Assert
	assert.Nil(t, tier)
}

// A quality edit no longer moves the zone on its own: it carries every incident
// guard to the same named tier of the new table. Both halves have to survive the
// round trip, and neither gains a field to do it - the tier is the zone's
// existing quality key and the guard is the plain number it always was, with the
// preset name inferred back from it when the panel is reopened.

func TestWhenAQualityEditMovedAnIncidentGuard_TheGuardSurvivesTheLoad(t *testing.T) {
	// Arrange
	edit := saveStateWithRetieredNeutralZone(t)

	// Act
	reloaded := reloadedState(t, edit.path)

	// Assert
	assert.Equal(t, edit.guardValue,
		reloadedGuardValue(t, reloaded, edit.connectionName))
}

func TestWhenAQualityEditMovedAnIncidentGuard_TheZoneTierSurvivesTheLoad(t *testing.T) {
	// Arrange
	edit := saveStateWithRetieredNeutralZone(t)

	// Act
	reloaded := reloadedState(t, edit.path)

	// Assert
	assert.Equal(t, edit.quality, reloadedZoneTier(t, reloaded, edit.zoneName))
}

// retieredGuardEdit is what a quality edit through the real handler produced.
// The round trip is asserted against the numbers that edit actually chose rather
// than against a preset table copied into the test.
type retieredGuardEdit struct {
	path           string
	zoneName       string
	connectionName string
	quality        neutral_zone.Quality
	guardValue     int
}

// saveStateWithRetieredNeutralZone generates a template, re-tiers one neutral
// zone through the real quality handler - passing the whole working graph, the
// way the dialog does - applies the returned zones and connections as a manual
// edit, and writes the state to disk.
func saveStateWithRetieredNeutralZone(t *testing.T) retieredGuardEdit {
	t.Helper()
	directory := t.TempDir()

	state := newUIState()
	state.UpdateState(func(s *editor_state_model.EditorState) { s.NeutralZoneCount = 4 })
	state.AutoRegenerate(time.Now())
	template := state.GetLastTemplate()
	require.NotNil(t, template, "expected a generated template")
	require.NotEmpty(t, template.Variants)

	zones := template.Variants[0].Zones
	connections := template.Variants[0].Connections
	edge, zoneName, ok := neutralToPlayerEdge(connections)
	require.True(t, ok, "the generated template has no neutral-to-player connection to re-tier")
	require.NotEmpty(t, edge.Name, "a generated connection is expected to carry a name")
	target, ok := zoneByName(zones, zoneName)
	require.True(t, ok, "the connection references a zone the template does not hold")

	handler := composition.InitializeGuiHandler()
	quality := neutral_zone.QualityHigh
	if handler.GetZoneQuality(target) == quality {
		quality = neutral_zone.QualityLow
	}
	mutation := handler.ApplyZoneEditorQuality(dtos.ZoneEditorQualityRequestDto{
		Zone:            target,
		Quality:         quality,
		CastleCount:     1,
		Tuning:          handler.GetZoneEditorOptions(state.GetStateDto(), len(zones)).Tuning,
		Zones:           zones,
		Connections:     connections,
		PlayerZoneNames: playerZoneNames(zones),
	})
	edited, ok := connectionByName(mutation.Connections, edge.Name)
	require.True(t, ok, "the quality edit dropped the incident connection")
	require.NotEqual(t, edge.GuardValue, edited.GuardValue,
		"the quality edit left the incident guard where it was, so there is nothing to round trip")

	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{
		Zones:       mutation.Zones,
		Connections: mutation.Connections,
	})
	state.SaveStateToFile(filepath.Join(directory, "guard.gen.json"))
	_, isError := state.GetStatus()
	require.False(t, isError, "saving the state failed")

	return retieredGuardEdit{
		path:           filepath.Join(directory, state.GetStateData().TemplateName+".gen.json"),
		zoneName:       zoneName,
		connectionName: edge.Name,
		quality:        quality,
		guardValue:     edited.GuardValue,
	}
}

// reloadedState loads a saved editor state into a fresh session.
func reloadedState(t *testing.T, path string) editor_state_model.EditorState {
	t.Helper()
	reloaded := newUIState()
	reloaded.LoadStateFromFile(path)
	_, isError := reloaded.GetStatus()
	require.False(t, isError, "loading the state failed")

	return reloaded.GetStateData()
}

func reloadedGuardValue(
	t *testing.T, state editor_state_model.EditorState, connectionName string) int {
	t.Helper()
	for _, save := range state.ManualConnections {
		if save.Connection.Name == connectionName {
			return save.Connection.GuardValue
		}
	}
	t.Fatalf("the restored snapshot holds no connection called %q", connectionName)

	return 0
}

func reloadedZoneTier(
	t *testing.T, state editor_state_model.EditorState, zoneName string) neutral_zone.Quality {
	t.Helper()
	for _, zone := range state.ManualZones {
		if zone.Name == zoneName {
			require.NotNil(t, zone.Quality, "the re-tiered zone came back with no recorded tier")
			return *zone.Quality
		}
	}
	t.Fatalf("the restored snapshot holds no zone called %q", zoneName)

	return neutral_zone.QualityUnknown
}

// neutralToPlayerEdge finds a connection between a neutral zone and a player
// spawn, and names its neutral end. A spawn resolves as Unknown, the weakest
// tier there is, so the neutral end is unambiguously the one that decides the
// guard table.
func neutralToPlayerEdge(connections []template_model.Connection) (
	template_model.Connection, string, bool) {
	for _, connection := range connections {
		switch {
		case zone_helpers.IsZoneNameNeutral(connection.From) &&
			zone_helpers.IsZoneNamePlayer(connection.To):
			return connection, connection.From, true
		case zone_helpers.IsZoneNameNeutral(connection.To) &&
			zone_helpers.IsZoneNamePlayer(connection.From):
			return connection, connection.To, true
		}
	}

	return template_model.Connection{}, "", false
}

func zoneByName(zones []template_model.Zone, name string) (template_model.Zone, bool) {
	for _, zone := range zones {
		if zone.Name == name {
			return zone, true
		}
	}

	return template_model.Zone{}, false
}

func connectionByName(connections []template_model.Connection, name string) (
	template_model.Connection, bool) {
	for _, connection := range connections {
		if connection.Name == name {
			return connection, true
		}
	}

	return template_model.Connection{}, false
}

func playerZoneNames(zones []template_model.Zone) []string {
	names := make([]string, 0, len(zones))
	for _, zone := range zones {
		if zone_helpers.IsZoneNamePlayer(zone.Name) {
			names = append(names, zone.Name)
		}
	}

	return names
}

// saveStateWithNeutralTier generates a template, stamps the given tier on every
// neutral zone, applies the layout as a manual edit and writes the state to
// disk. It returns the path the state actually wrote.
func saveStateWithNeutralTier(t *testing.T, quality *neutral_zone.Quality) string {
	t.Helper()
	directory := t.TempDir()

	state := newUIState()
	state.UpdateState(func(s *editor_state_model.EditorState) { s.NeutralZoneCount = 4 })
	state.AutoRegenerate(time.Now())
	template := state.GetLastTemplate()
	require.NotNil(t, template, "expected a generated template")
	require.NotEmpty(t, template.Variants)

	zones := append([]template_model.Zone(nil), template.Variants[0].Zones...)
	require.NotZero(t, countNeutralZones(zones), "the generated template has no neutral zone to re-tier")
	for index := range zones {
		if zone_helpers.IsZoneNameNeutral(zones[index].Name) {
			zones[index].Quality = quality
		}
	}
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{
		Zones:       zones,
		Connections: template.Variants[0].Connections,
	})

	state.SaveStateToFile(filepath.Join(directory, "tier.gen.json"))
	_, isError := state.GetStatus()
	require.False(t, isError, "saving the state failed")

	return filepath.Join(directory, state.GetStateData().TemplateName+".gen.json")
}

// reloadedNeutralTier loads the file into a fresh session and reads the tier
// off the restored manual snapshot, which is what a later regeneration
// reapplies over the freshly generated template.
func reloadedNeutralTier(t *testing.T, path string) *neutral_zone.Quality {
	t.Helper()
	reloaded := newUIState()
	reloaded.LoadStateFromFile(path)
	_, isError := reloaded.GetStatus()
	require.False(t, isError, "loading the state failed")

	zones := reloaded.GetStateData().ManualZones
	require.NotEmpty(t, zones, "the manual snapshot was not persisted")
	for _, zone := range zones {
		if zone_helpers.IsZoneNameNeutral(zone.Name) {
			return zone.Quality
		}
	}
	t.Fatal("the restored snapshot holds no neutral zone")

	return nil
}

func countNeutralZones(zones []template_model.Zone) int {
	count := 0
	for _, zone := range zones {
		if zone_helpers.IsZoneNameNeutral(zone.Name) {
			count++
		}
	}
	return count
}
