package position_layout

import (
	"math"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

type PositionLayoutService struct{}

func NewPositionLayoutService() *PositionLayoutService {
	return &PositionLayoutService{}
}

func (this *PositionLayoutService) CreatePositionsFromPlans(
	orderedLabels,
	playerLabels []string,
	neutralZonePlans neutral_zone.Plans,
) models.Positions {
	count := len(orderedLabels)
	if count == 0 {
		return nil
	}

	getTierRadius := func(tier int) float64 {
		switch tier {
		case 0:
			return 0.38
		case 1:
			return 0.27
		case 2:
			return 0.16
		default:
			return 0.06
		}
	}

	byTier := make(map[int][]int)
	for index, label := range orderedLabels {
		tier := 0
		if !slices.Contains(playerLabels, label) {
			tier = neutralZonePlans.GetTier(label)
		}
		byTier[tier] = append(byTier[tier], index)
	}

	positions := make(models.Positions, count)
	for tier, indices := range byTier {
		radius := getTierRadius(tier)
		positionCount := float64(len(indices))
		offset := float64(tier) * math.Pi / positionCount
		for index, positionIndex := range indices {
			angle := 2*math.Pi*float64(index)/positionCount + offset
			jitter := float64(index%3-1) * 0.008
			positions[positionIndex] = data.NewVec2(
				max(0.05, min(0.95, 0.5+math.Cos(angle+jitter)*radius)),
				max(0.05, min(0.95, 0.5+math.Sin(angle+jitter)*radius)),
			)
		}
	}
	return positions
}
