package zoneCandidate_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoCandidatesExist_ReturnsEmptyString(t *testing.T) {
	t.Parallel()
	// Arrange
	candidates := utils.CreateHubZoneCandidates(neutral_zone.Plans{}, nil)

	// Act
	label := candidates.GetFirstCandidateLabel()

	// Assert
	assert.Empty(t, label)
}

func TestWhenCandidatesExist_ReturnsFirstCandidateLetter(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutral_zone.Plans{
		{Label: "C", Quality: neutral_zone.QualityMedium, CastleCount: 0},
		{Label: "D", Quality: neutral_zone.QualityMedium, CastleCount: 0},
	}
	candidates := utils.CreateHubZoneCandidates(plans, []map[string]int{{"C": 1, "D": 2}})

	// Act
	label := candidates.GetFirstCandidateLabel()

	// Assert
	assert.Equal(t, "C", label)
}
