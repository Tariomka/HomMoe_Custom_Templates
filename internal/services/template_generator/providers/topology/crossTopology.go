package topology

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
)

// CrossTopologyService radiates one arm per player out of a central zone. Each
// arm is a chain of neutral zones running from the centre out to the player at
// its tip, so the zones and connections form a cross / star whose number of
// arms follows the player count.
type CrossTopologyService struct {
	RandomTopologyService
}

func NewCrossTopologyService() *CrossTopologyService {
	return &CrossTopologyService{
		RandomTopologyService: *NewRandomTopologyService(),
	}
}

func (this *CrossTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutralZone.Plans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) entities.Variant {
	return this.createVariantFromLayout(
		configuration, playerLabels, neutralZones, tuning, holdCityNeutralLabel, this.createCrossLayout)
}

// createCrossLayout places the central zone first, then walks each arm from the
// centre outward laying its neutral zones followed by the player at the tip.
func (this *CrossTopologyService) createCrossLayout(
	playerLabels []string,
	neutralZones neutralZone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
	const (
		centreX      = 0.5
		centreY      = 0.5
		playerRadius = 0.42
		armNear      = 0.14
		armFar       = 0.34
		startAngle   = -math.Pi / 2.0
	)
	playerCount := len(playerLabels)
	neutralCount := len(neutralZones)

	var allLabels []string
	var positions models.Positions

	// Central zone is the heart of the cross.
	centreIndex := -1
	if neutralCount >= 1 {
		centreIndex = len(allLabels)
		allLabels = append(allLabels, neutralZones[0].Label)
		positions.Add(data.NewVec2(centreX, centreY))
	}

	// Distribute the remaining neutral zones across the arms, round-robin.
	armNeutralCounts := make([]int, playerCount)
	for i := range neutralCount - 1 {
		armNeutralCounts[i%max(1, playerCount)]++
	}

	armIndices := make([][]int, playerCount)
	nextNeutral := 1
	for arm := range playerCount {
		angle := startAngle
		if playerCount > 0 {
			angle += 2.0 * math.Pi * float64(arm) / float64(playerCount)
		}

		// Arm neutral zones from the centre outward.
		count := armNeutralCounts[arm]
		for k := range count {
			radius := (armNear + armFar) / 2.0
			if count > 1 {
				radius = armNear + (armFar-armNear)*float64(k)/float64(count-1)
			}
			armIndices[arm] = append(armIndices[arm], len(allLabels))
			allLabels = append(allLabels, neutralZones[nextNeutral].Label)
			positions.Add(circlePoint(angle, centreX, centreY, radius))
			nextNeutral++
		}

		// Player zone at the arm tip.
		armIndices[arm] = append(armIndices[arm], len(allLabels))
		allLabels = append(allLabels, playerLabels[arm])
		positions.Add(circlePoint(angle, centreX, centreY, playerRadius))
	}

	pairs := this.createCrossPairs(centreIndex, armIndices, playerCount)
	return allLabels, positions, pairs
}

func (this *CrossTopologyService) createCrossPairs(
	centreIndex int,
	armIndices [][]int,
	playerCount int) []models.ConnectionIndexes {
	builder := newPairBuilder()

	for arm := range playerCount {
		indices := armIndices[arm]
		if len(indices) == 0 {
			continue
		}
		// The centre joins the innermost zone of each arm.
		if centreIndex >= 0 {
			builder.add(centreIndex, indices[0])
		}
		// Chain the arm zones from the centre outward to the player tip.
		for k := 0; k+1 < len(indices); k++ {
			builder.add(indices[k], indices[k+1])
		}
	}

	// Without a centre zone (no neutral zones) join the player tips in a ring so
	// they still form a closed, cross-like outline.
	if centreIndex < 0 && playerCount >= 2 {
		for i := range playerCount {
			tip := armIndices[i][len(armIndices[i])-1]
			nextTip := armIndices[(i+1)%playerCount][len(armIndices[(i+1)%playerCount])-1]
			builder.add(tip, nextTip)
		}
	}
	return builder.pairs
}
