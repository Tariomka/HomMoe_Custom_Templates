package zoneCandidate_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/stretchr/testify/assert"
)

func TestWhenMinDistancesDiffer_PutsFarthestZoneFirst(t *testing.T) {
	// Arrange
	plans := models.NeutralZonePlans{
		{Label: "X", Quality: models.QualityMedium, CastleCount: 0},
		{Label: "Y", Quality: models.QualityMedium, CastleCount: 0},
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
	// Arrange
	plans := models.NeutralZonePlans{
		{Label: "Y", Quality: models.QualityMedium, CastleCount: 0},
		{Label: "X", Quality: models.QualityMedium, CastleCount: 0},
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
	// Arrange
	plans := models.NeutralZonePlans{
		{Label: "Y", Quality: models.QualityLow, CastleCount: 0},
		{Label: "X", Quality: models.QualityHigh, CastleCount: 0},
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
	// Arrange
	plans := models.NeutralZonePlans{
		{Label: "Y", Quality: models.QualityMedium, CastleCount: 0},
		{Label: "X", Quality: models.QualityMedium, CastleCount: 2},
	}
	distancesByPlayer := []map[string]int{{"X": 2, "Y": 2}}

	// Act
	label := utils.CreateHubZoneCandidates(plans, distancesByPlayer).
		SortForHubCity().
		GetFirstCandidateLabel()

	// Assert
	assert.Equal(t, "X", label)
}
