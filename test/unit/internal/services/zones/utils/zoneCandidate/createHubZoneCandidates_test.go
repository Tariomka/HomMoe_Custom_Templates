package zoneCandidate_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlanLabelsHaveDistances_CreatesOneCandidatePerPlan(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{
		{Label: "C", Quality: neutral_zone.QualityMedium, CastleCount: 1},
		{Label: "D", Quality: neutral_zone.QualityLow, CastleCount: 0},
	}
	distancesByPlayer := []map[string]int{{"C": 1, "D": 2}}

	// Act
	candidates := utils.CreateHubZoneCandidates(plans, distancesByPlayer)

	// Assert
	assert.Len(t, *candidates, 2)
}

func TestWhenDistanceIsMissingForPlayer_TreatsZoneAsVeryFar(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{
		{Label: "X", Quality: neutral_zone.QualityMedium, CastleCount: 0},
		{Label: "Y", Quality: neutral_zone.QualityMedium, CastleCount: 0},
	}
	distancesByPlayer := []map[string]int{{"X": 1}, {"X": 1}}

	// Act
	label := utils.CreateHubZoneCandidates(plans, distancesByPlayer).
		SortForHubCity().
		GetFirstCandidateLabel()

	// Assert
	assert.Equal(t, "Y", label)
}
