package connectionEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func defaultTuning() models.GenerationTuning {
	return models.GenerationTuning{
		ContentScale:                   1.0,
		ResourceDensityMultiplier:      0.5,
		StructureDensityMultiplier:     1.0,
		NeutralStackStrengthMultiplier: 1.0,
		BorderGuardStrengthMultiplier:  1.0,
	}
}

// ════════════════════════════════════════════════════════════════════════
// NextFreeZoneLabel
// ════════════════════════════════════════════════════════════════════════

func TestNextFreeZoneLabel_EmptyList_ReturnsA(t *testing.T) {
	assert.Equal(t, "A", connection_editor.NextFreeZoneLabel(nil))
}

func TestNextFreeZoneLabel_SkipsUsedLetters(t *testing.T) {
	zones := []entities.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-B"},
	}
	assert.Equal(t, "C", connection_editor.NextFreeZoneLabel(zones))
}

func TestNextFreeZoneLabel_SharedLetterAcrossPrefixes_CountsOnce(t *testing.T) {
	zones := []entities.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-A"},
	}
	assert.Equal(t, "B", connection_editor.NextFreeZoneLabel(zones))
}

// ════════════════════════════════════════════════════════════════════════
// NewDefaultNeutralZone
// ════════════════════════════════════════════════════════════════════════

func TestNewDefaultNeutralZone_BasicShape(t *testing.T) {
	zone := connection_editor.NewDefaultNeutralZone("Q", models.QualityMedium, 1, false, defaultTuning())

	assert.Equal(t, "Neutral-Q", zone.Name)
	assert.Nil(t, zone.MandatoryContent, "manual zones must not reference template-level mandatory content")
	assert.Equal(t, 1, connection_editor.CountZoneCastles(zone))
	assert.Equal(t, models.QualityMedium, connection_editor.QualityOfZone(zone))
}

func TestNewDefaultNeutralZone_NoCastles(t *testing.T) {
	zone := connection_editor.NewDefaultNeutralZone("R", models.QualityLow, 0, false, defaultTuning())
	assert.Equal(t, 0, connection_editor.CountZoneCastles(zone))
}

// ════════════════════════════════════════════════════════════════════════
// QualityOfZone / ApplyNeutralZoneQuality
// ════════════════════════════════════════════════════════════════════════

func TestQualityOfZone_DetectsAllPresets(t *testing.T) {
	tuning := defaultTuning()
	for _, quality := range []models.NeutralZoneQuality{models.QualityLow, models.QualityMedium, models.QualityHigh} {
		zone := connection_editor.NewDefaultNeutralZone("Z", quality, 1, false, tuning)
		assert.Equal(t, quality, connection_editor.QualityOfZone(zone))
	}
}

func TestApplyNeutralZoneQuality_ChangesPoolsAndCastles(t *testing.T) {
	tuning := defaultTuning()
	zone := connection_editor.NewDefaultNeutralZone("Z", models.QualityLow, 0, false, tuning)

	connection_editor.ApplyNeutralZoneQuality(&zone, models.QualityHigh, 2, tuning)

	assert.Equal(t, models.QualityHigh, connection_editor.QualityOfZone(zone))
	assert.Equal(t, 2, connection_editor.CountZoneCastles(zone))
}

// Adding castles to a former connector zone (0 -> 3 castles) must produce the
// stone castle<->castle roads that link the new castles to the primary one.
// Regression test for the missing third-castle road.
func TestApplyNeutralZoneQuality_RegeneratesCastleRoads(t *testing.T) {
	tuning := defaultTuning()
	// Start as a castle-less connector zone (no MainObject roads at all).
	zone := connection_editor.NewDefaultNeutralZone("Z", models.QualityMedium, 0, true, tuning)

	connection_editor.ApplyNeutralZoneQuality(&zone, models.QualityHigh, 3, tuning)

	castleRoads := castleRoadTargets(zone)
	assert.Equal(t, []string{"1", "2"}, castleRoads,
		"three castles must be linked by stone roads 0->1 and 0->2")
}

// Reducing the castle count must drop the now-dangling castle roads.
func TestApplyNeutralZoneQuality_RemovesStaleCastleRoads(t *testing.T) {
	tuning := defaultTuning()
	zone := connection_editor.NewDefaultNeutralZone("Z", models.QualityHigh, 3, true, tuning)

	connection_editor.ApplyNeutralZoneQuality(&zone, models.QualityHigh, 1, tuning)

	assert.Empty(t, castleRoadTargets(zone),
		"a single-castle zone must have no castle<->castle roads")
}

// ════════════════════════════════════════════════════════════════════════
// CanDeleteZone / RemoveZone
// ════════════════════════════════════════════════════════════════════════

func TestCanDeleteZone_ProtectsSpawnZones(t *testing.T) {
	players := map[string]bool{"Spawn-A": true, "Spawn-B": true}
	assert.False(t, connection_editor.CanDeleteZone("Spawn-A", players))
	assert.True(t, connection_editor.CanDeleteZone("Neutral-C", players))
	assert.True(t, connection_editor.CanDeleteZone("Hub", players))
}

func TestRemoveZone_RemovesZoneAndItsConnections(t *testing.T) {
	zones := []entities.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-C"},
		{Name: "Neutral-D"},
	}
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-C"},
		{From: "Neutral-C", To: "Neutral-D"},
		{From: "Spawn-A", To: "Neutral-D"},
	}

	keptZones, keptConnections := connection_editor.RemoveZone(zones, connections, "Neutral-C")

	assert.Len(t, keptZones, 2)
	for _, zone := range keptZones {
		assert.NotEqual(t, "Neutral-C", zone.Name)
	}
	assert.Len(t, keptConnections, 1)
	assert.Equal(t, "Spawn-A", keptConnections[0].From)
	assert.Equal(t, "Neutral-D", keptConnections[0].To)
}

// ════════════════════════════════════════════════════════════════════════
// FindOpenPosition
// ════════════════════════════════════════════════════════════════════════

func TestFindOpenPosition_EmptyBoard_NearCenterEdgeOfGrid(t *testing.T) {
	position := connection_editor.FindOpenPosition(nil)
	assert.InDelta(t, 0.5, position[0], 0.401)
	assert.InDelta(t, 0.5, position[1], 0.401)
}

func TestFindOpenPosition_AvoidsOccupiedCorner(t *testing.T) {
	occupied := [][2]float64{{0.1, 0.1}, {0.1, 0.2}, {0.2, 0.1}}
	position := connection_editor.FindOpenPosition(occupied)
	// The best spot should be far from the cluttered top-left corner.
	assert.Greater(t, position[0]+position[1], 1.0)
}
