package zoneEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityChangesToHigh_ReprofilesZoneAsHigh(t *testing.T) {
	t.Parallel()
	// Arrange
	tuning := defaultTuning()
	zone := connection_editor.NewDefaultNeutralZone("Z", neutralZone.QualityLow, 0, false, tuning)

	// Act
	connection_editor.ApplyNeutralZoneQuality(&zone, neutralZone.QualityHigh, 2, tuning)

	// Assert
	assert.Equal(t, neutralZone.QualityHigh, neutralZone.GetQualityFrom(zone))
}

func TestWhenTwoCastlesAreRequested_RebuildsTwoCastles(t *testing.T) {
	t.Parallel()
	// Arrange
	tuning := defaultTuning()
	zone := connection_editor.NewDefaultNeutralZone("Z", neutralZone.QualityLow, 0, false, tuning)

	// Act
	connection_editor.ApplyNeutralZoneQuality(&zone, neutralZone.QualityHigh, 2, tuning)

	// Assert
	assert.Equal(t, 2, connection_editor.CountZoneCastles(zone))
}

// Adding castles to a former connector zone (0 -> 3 castles) must produce the
// stone castle<->castle roads that link the new castles to the primary one.
// Regression test for the missing third-castle road.
func TestWhenCastlesAreAddedToConnectorZone_RegeneratesCastleRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	tuning := defaultTuning()
	zone := connection_editor.NewDefaultNeutralZone("Z", neutralZone.QualityMedium, 0, true, tuning)

	// Act
	connection_editor.ApplyNeutralZoneQuality(&zone, neutralZone.QualityHigh, 3, tuning)

	// Assert
	assert.Equal(t, []string{"1", "2"}, castleRoadTargets(zone),
		"three castles must be linked by stone roads 0->1 and 0->2")
}

// Reducing the castle count must drop the now-dangling castle roads.
func TestWhenCastleCountShrinksToOne_RemovesStaleCastleRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	tuning := defaultTuning()
	zone := connection_editor.NewDefaultNeutralZone("Z", neutralZone.QualityHigh, 3, true, tuning)

	// Act
	connection_editor.ApplyNeutralZoneQuality(&zone, neutralZone.QualityHigh, 1, tuning)

	// Assert
	assert.Empty(t, castleRoadTargets(zone),
		"a single-castle zone must have no castle<->castle roads")
}
