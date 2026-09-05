package footholdDistancePresets_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_distances"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenRequested_ReturnsTheFootholdPlacementBounds(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := common_distances.FootholdDistancePresets{
		Crossroads:       models.DistancePreset{Name: "Foothold Crossroads", Min: 0.2, Max: 0.3},
		NearMainCastle:   models.DistancePreset{Name: "Foothold Main Castle", Min: 0.2, Max: 0.4},
		NearSecondCastle: models.DistancePreset{Name: "Foothold Second Castle", Min: 0.5, Max: 0.5},
	}

	// Act
	result := common_distances.GetFootholdDistancePresets()

	// Assert
	assert.Equal(t, expected, result)
}
