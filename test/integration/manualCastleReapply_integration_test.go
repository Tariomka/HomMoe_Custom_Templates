//go:build integration_test

package integration_test

import (
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// regenerateAfterDebounce drives AutoRegenerate through the debounce window a
// non-layout change goes through: one call arms the timer, a later call fires
// the regeneration.
func regenerateAfterDebounce(state *drivers.State, now time.Time) {
	state.AutoRegenerate(now)
	state.AutoRegenerate(now.Add(time.Second))
}

// connectionKeys captures the identity of every connection so tests can prove
// the manual connection graph survives a regeneration untouched.
func connectionKeys(connections []entities.Connection) map[[3]string]bool {
	keys := make(map[[3]string]bool, len(connections))
	for _, connection := range connections {
		keys[[3]string{connection.From, connection.To, connection.ConnectionType}] = true
	}
	return keys
}

func retierZone(state *drivers.State, zones []entities.Zone, index int, quality neutral_zone.Quality, castles int) {
	configuration := mappers.NewConfigMapper().FromEditorState(state.GetStateData())
	tuning := test_helpers.NewGenerationTuning(configuration, len(zones))
	test_helpers.NewZoneEditorService().ApplyNeutralZoneQuality(&zones[index], quality, castles, tuning)
}

func findNeutralOfQuality(t *testing.T, zones []entities.Zone, quality neutral_zone.Quality) int {
	t.Helper()
	for i, zone := range zones {
		if zone_helpers.IsZoneNameNeutral(zone.Name) && zone_services.NewZoneClassifier().GetQuality(zone) == quality {
			return i
		}
	}
	t.Fatalf("no neutral zone of quality %v found", quality)
	return -1
}

// TestCastleOptionChange_AfterManualEdits_UpdatesSnapshotCastles covers the
// core requirement: after manual zone editing, changing the simple-mode
// neutral castle count must update every neutral zone in the manual snapshot
// (including a manually re-tiered one) while the manual connection graph,
// positions and qualities stay exactly as edited.
func TestCastleOptionChange_AfterManualEdits_UpdatesSnapshotCastles(t *testing.T) {
	now := time.Now()
	state := newUIState()
	state.UpdateState(func(s *dtos.EditorStateDto) { s.NeutralZoneCount = 4 })
	state.AutoRegenerate(now)

	template := state.GetLastTemplate()
	require.NotNil(t, template)
	require.NotEmpty(t, template.Variants)

	// Manual session: re-tier one neutral zone to High, stamp positions and
	// add a user connection.
	zones := append([]entities.Zone(nil), template.Variants[0].Zones...)
	retiered := findNeutralOfQuality(t, zones, neutral_zone.QualityMedium)
	retierZone(state, zones, retiered, neutral_zone.QualityHigh, 1)
	retieredName := zones[retiered].Name
	for i := range zones {
		zones[i].ManualPosition = &[2]float64{0.1 * float64(i+1), 0.05 * float64(i+1)}
	}
	connections := append([]entities.Connection(nil), template.Variants[0].Connections...)
	connections = append(connections, entities.Connection{
		From: zones[0].Name, To: zones[1].Name, ConnectionType: "Portal", IsUserAdded: true,
	})
	state.ApplyEditedZones(zones, connections)
	expectedConnections := connectionKeys(connections)

	// The "last change after manual editing" is a castle count.
	state.UpdateState(func(s *dtos.EditorStateDto) { s.NeutralZoneCastles = 3 })
	regenerateAfterDebounce(state, now)

	got := state.GetLastTemplate()
	require.NotNil(t, got)
	require.NotEmpty(t, got.Variants)
	for _, zone := range got.Variants[0].Zones {
		if !zone_helpers.IsZoneNameNeutral(zone.Name) {
			continue
		}
				assert.Equalf(t, 3, test_helpers.NewZoneEditorService().CountZoneCastles(zone),
			"zone %s must follow the new simple-mode castle count", zone.Name)
		if zone.Name == retieredName {
			assert.Equal(t, neutral_zone.QualityHigh, zone_services.NewZoneClassifier().GetQuality(zone),
				"the manual quality change must survive the castle update")
		}
		assert.NotNilf(t, zone.ManualPosition, "zone %s lost its manual position", zone.Name)
	}

	gotConnections := connectionKeys(got.Variants[0].Connections)
	assert.Equal(t, expectedConnections, gotConnections,
		"the manual connection graph must survive a castle-count regeneration untouched")

	// The updated counts must be persisted back into the snapshot so saves and
	// later regenerations carry them.
	for _, save := range state.GetStateData().ManualZones {
		if zone_helpers.IsZoneNameNeutral(save.Zone.Name) {
						assert.Equal(t, 3, test_helpers.NewZoneEditorService().CountZoneCastles(save.Zone))
		}
	}
}

// TestAdvancedTierCastleChange_UpdatesByManualQuality covers the advanced-mode
// requirement: a per-tier castles-per-zone change must drive zones by their
// MANUAL quality, not the quality the generator originally planned.
func TestAdvancedTierCastleChange_UpdatesByManualQuality(t *testing.T) {
	now := time.Now()
	state := newUIState()
	state.UpdateState(func(s *dtos.EditorStateDto) {
		s.AdvancedMode = true
		s.NeutralLowCastleCount = 2  // two Low zones with one castle each
		s.NeutralHighCastleCount = 1 // one High zone with one castle
	})
	state.AutoRegenerate(now)

	template := state.GetLastTemplate()
	require.NotNil(t, template)
	require.NotEmpty(t, template.Variants)

	// Manually promote one of the Low zones to High.
	zones := append([]entities.Zone(nil), template.Variants[0].Zones...)
	promoted := findNeutralOfQuality(t, zones, neutral_zone.QualityLow)
	retierZone(state, zones, promoted, neutral_zone.QualityHigh, 1)
	promotedName := zones[promoted].Name
	state.ApplyEditedZones(zones, template.Variants[0].Connections)

	state.UpdateState(func(s *dtos.EditorStateDto) { s.NeutralHighCastlesPerZone = 3 })
	regenerateAfterDebounce(state, now)

	got := state.GetLastTemplate()
	require.NotNil(t, got)
	require.NotEmpty(t, got.Variants)
	for _, zone := range got.Variants[0].Zones {
		if !zone_helpers.IsZoneNameNeutral(zone.Name) {
			continue
		}
		expected := 1
		if zone_services.NewZoneClassifier().GetQuality(zone) == neutral_zone.QualityHigh {
			expected = 3
		}
				assert.Equalf(t, expected, test_helpers.NewZoneEditorService().CountZoneCastles(zone),
			"zone %s (quality %v)", zone.Name, zone_services.NewZoneClassifier().GetQuality(zone))
		if zone.Name == promotedName {
			assert.Equal(t, neutral_zone.QualityHigh, zone_services.NewZoneClassifier().GetQuality(zone))
		}
	}
}

// TestNonCastleChange_AfterManualEdits_KeepsSnapshotVerbatim confirms that a
// non-castle option change reapplies the manual snapshot without touching its
// castles or guard values.
func TestNonCastleChange_AfterManualEdits_KeepsSnapshotVerbatim(t *testing.T) {
	now := time.Now()
	state := newUIState()
	state.UpdateState(func(s *dtos.EditorStateDto) { s.NeutralZoneCount = 3 })
	state.AutoRegenerate(now)

	template := state.GetLastTemplate()
	require.NotNil(t, template)
	require.NotEmpty(t, template.Variants)

	zones := append([]entities.Zone(nil), template.Variants[0].Zones...)
	edited := findNeutralOfQuality(t, zones, neutral_zone.QualityMedium)
	retierZone(state, zones, edited, neutral_zone.QualityHigh, 2)
	zones[edited].GuardMultiplier = 7.5 // explicit manual guard edit
	editedName := zones[edited].Name
	state.ApplyEditedZones(zones, template.Variants[0].Connections)

	// Non-castle, non-layout change.
	state.UpdateState(func(s *dtos.EditorStateDto) { s.NeutralStackStrengthPercent = 150 })
	regenerateAfterDebounce(state, now)

	got := state.GetLastTemplate()
	require.NotNil(t, got)
	require.NotEmpty(t, got.Variants)
	found := false
	for _, zone := range got.Variants[0].Zones {
		if zone.Name != editedName {
			continue
		}
		found = true
				assert.Equal(t, 2, test_helpers.NewZoneEditorService().CountZoneCastles(zone),
			"a non-castle option change must not touch the manual castle count")
		assert.Equal(t, 7.5, zone.GuardMultiplier,
			"a non-castle option change must not touch manual guard values")
		assert.Equal(t, neutral_zone.QualityHigh, zone_services.NewZoneClassifier().GetQuality(zone))
	}
	assert.True(t, found, "manually edited zone disappeared from the regenerated template")
}
