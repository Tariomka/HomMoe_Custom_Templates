package zoneCandidate_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoCandidatesExist_ReturnsEmptyString(t *testing.T) {
	t.Parallel()
	// Arrange
	candidates := utils.CreateHubZoneCandidates(neutralZone.Plans{}, nil)

	// Act
	label := candidates.GetFirstCandidateLabel()

	// Assert
	assert.Empty(t, label)
}

func TestWhenCandidatesExist_ReturnsFirstCandidateLetter(t *testing.T) {
	t.Parallel()
	// Arrange
	plans := neutralZone.Plans{
		{Label: "C", Quality: neutralZone.QualityMedium, CastleCount: 0},
		{Label: "D", Quality: neutralZone.QualityMedium, CastleCount: 0},
	}
	candidates := utils.CreateHubZoneCandidates(plans, []map[string]int{{"C": 1, "D": 2}})

	// Act
	label := candidates.GetFirstCandidateLabel()

	// Assert
	assert.Equal(t, "C", label)
}
