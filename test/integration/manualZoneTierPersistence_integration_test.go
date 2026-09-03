//go:build integration_test

package integration_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

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
		Connections: template_model.ToConnectionEntities(template.Variants[0].Connections),
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

	saves := reloaded.GetStateData().ManualZones
	require.NotEmpty(t, saves, "the manual snapshot was not persisted")
	for _, zone := range editor_state_model.FromManualZoneSaves(saves) {
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
