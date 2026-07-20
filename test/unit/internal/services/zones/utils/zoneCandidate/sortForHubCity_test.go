package zoneCandidate_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/stretchr/testify/assert"
)

func TestWhenMinDistancesDiffer_PutsFarthestZoneFirst(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{
		{Label: "X", Quality: neutral_zone.QualityMedium, CastleCount: 0},
		{Label: "Y", Quality: neutral_zone.QualityMedium, CastleCount: 0},
	}
	distancesByPlayer := []map[string]int{{"X": 1, "Y": 3}}

	// Act
	label := utils.CreateHubZoneCandidates(plans, distancesByPlayer).
		SortForHubCity().
		GetFirstCandidateLabel()

	// Assert
	assert.Equal(t, "Y", label)
}

func TestWhenMinDistancesTie_PutsLowestVarianceFirst(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{
		{Label: "Y", Quality: neutral_zone.QualityMedium, CastleCount: 0},
		{Label: "X", Quality: neutral_zone.QualityMedium, CastleCount: 0},
	}
	distancesByPlayer := []map[string]int{
		{"X": 2, "Y": 2},
		{"X": 2, "Y": 4},
	}

	// Act
	label := utils.CreateHubZoneCandidates(plans, distancesByPlayer).
		SortForHubCity().
		GetFirstCandidateLabel()

	// Assert
	assert.Equal(t, "X", label)
}

func TestWhenVarianceAlsoTies_PutsHighestQualityFirst(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{
		{Label: "Y", Quality: neutral_zone.QualityLow, CastleCount: 0},
		{Label: "X", Quality: neutral_zone.QualityHigh, CastleCount: 0},
	}
	distancesByPlayer := []map[string]int{{"X": 2, "Y": 2}}

	// Act
	label := utils.CreateHubZoneCandidates(plans, distancesByPlayer).
		SortForHubCity().
		GetFirstCandidateLabel()

	// Assert
	assert.Equal(t, "X", label)
}

func TestWhenQualityAlsoTies_PutsCastleZoneFirst(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{
		{Label: "Y", Quality: neutral_zone.QualityMedium, CastleCount: 0},
		{Label: "X", Quality: neutral_zone.QualityMedium, CastleCount: 2},
	}
	distancesByPlayer := []map[string]int{{"X": 2, "Y": 2}}

	// Act
	label := utils.CreateHubZoneCandidates(plans, distancesByPlayer).
		SortForHubCity().
		GetFirstCandidateLabel()

	// Assert
	assert.Equal(t, "X", label)
}
