package zoneCandidate_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoCandidatesExist_ReturnsEmptyString(t *testing.T) {
	// Arrange
	candidates := utils.CreateHubZoneCandidates(models.NeutralZonePlans{}, nil)

	// Act
	label := candidates.GetFirstCandidateLabel()

	// Assert
	assert.Equal(t, "", label)
}

func TestWhenCandidatesExist_ReturnsFirstCandidateLetter(t *testing.T) {
	// Arrange
	plans := models.NeutralZonePlans{
		{Label: "C", Quality: models.QualityMedium, CastleCount: 0},
		{Label: "D", Quality: models.QualityMedium, CastleCount: 0},
	}
	candidates := utils.CreateHubZoneCandidates(plans, []map[string]int{{"C": 1, "D": 2}})

	// Act
	label := candidates.GetFirstCandidateLabel()

	// Assert
	assert.Equal(t, "C", label)
}
