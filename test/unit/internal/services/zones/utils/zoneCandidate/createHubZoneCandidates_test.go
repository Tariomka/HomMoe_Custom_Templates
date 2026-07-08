package zoneCandidate_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlanLabelsHaveDistances_CreatesOneCandidatePerPlan(t *testing.T) {
	// Arrange
	plans := models.NeutralZonePlans{
		{Label: "C", Quality: models.QualityMedium, CastleCount: 1},
		{Label: "D", Quality: models.QualityLow, CastleCount: 0},
	}
	distancesByPlayer := []map[string]int{{"C": 1, "D": 2}}

	// Act
	candidates := utils.CreateHubZoneCandidates(plans, distancesByPlayer)

	// Assert
	assert.Len(t, *candidates, 2)
}

func TestWhenDistanceIsMissingForPlayer_TreatsZoneAsVeryFar(t *testing.T) {
	// Arrange
	plans := models.NeutralZonePlans{
		{Label: "X", Quality: models.QualityMedium, CastleCount: 0},
		{Label: "Y", Quality: models.QualityMedium, CastleCount: 0},
	}
	distancesByPlayer := []map[string]int{{"X": 1}, {"X": 1}}

	// Act
	label := utils.CreateHubZoneCandidates(plans, distancesByPlayer).
		SortForHubCity().
		GetFirstCandidateLabel()

	// Assert
	assert.Equal(t, "Y", label)
}
