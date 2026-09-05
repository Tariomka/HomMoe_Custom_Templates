package topology

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

// CrossTopologyService radiates one arm per player out of a central zone. Each
// arm is a chain of neutral zones running from the center out to the player at
// its tip, so the zones and connections form a cross / star whose number of
// arms follows the player count.
type CrossTopologyService struct {
	PositionedTopologyBuilder
}

func NewCrossTopologyService(
	zoneFactory zone_interfaces.IZoneFactory,
	roadFactory zone_interfaces.IRoadFactory,
	zoneLabelProvider zone_interfaces.IZoneLabelProvider,
	connectionService base.ITopologyConnectionService) *CrossTopologyService {
	return &CrossTopologyService{
		PositionedTopologyBuilder: *NewPositionedTopologyBuilder(
			zoneFactory, roadFactory, zoneLabelProvider, connectionService),
	}
}

func (this *CrossTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) template_model.Variant {
	return this.BuildVariant(
		configuration, playerLabels, neutralZones, tuning, holdCityNeutralLabel, this.createCrossLayout, nil)
}

// createCrossLayout places the central zone first, then walks each arm from the
// center outward laying its neutral zones followed by the player at the tip.
func (this *CrossTopologyService) createCrossLayout(
	playerLabels []string,
	neutralZones neutral_zone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
	const (
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
	centerIndex := -1
	if neutralCount >= 1 {
		centerIndex = len(allLabels)
		allLabels = append(allLabels, neutralZones[0].Label)
		positions.Add(data.NewVec2(layoutCenter, layoutCenter))
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

		// Arm neutral zones from the center outward.
		count := armNeutralCounts[arm]
		for k := range count {
			radius := (armNear + armFar) / 2.0
			if count > 1 {
				radius = armNear + (armFar-armNear)*float64(k)/float64(count-1)
			}
			armIndices[arm] = append(armIndices[arm], len(allLabels))
			allLabels = append(allLabels, neutralZones[nextNeutral].Label)
			positions.Add(circlePoint(angle, radius))
			nextNeutral++
		}

		// Player zone at the arm tip.
		armIndices[arm] = append(armIndices[arm], len(allLabels))
		allLabels = append(allLabels, playerLabels[arm])
		positions.Add(circlePoint(angle, playerRadius))
	}

	pairs := this.createCrossPairs(centerIndex, armIndices, playerCount)
	return allLabels, positions, pairs
}

func (this *CrossTopologyService) createCrossPairs(
	centerIndex int,
	armIndices [][]int,
	playerCount int) []models.ConnectionIndexes {
	builder := newPairBuilder()

	for arm := range playerCount {
		indices := armIndices[arm]
		if len(indices) == 0 {
			continue
		}

		// The center joins the innermost zone of each arm.
		if centerIndex >= 0 {
			builder.add(centerIndex, indices[0])
		}
		// Chain the arm zones from the center outward to the player tip.
		for k := 0; k+1 < len(indices); k++ {
			builder.add(indices[k], indices[k+1])
		}
	}

	// Without a center zone (no neutral zones) join the player tips in a ring so
	// they still form a closed, cross-like outline.
	if centerIndex < 0 && playerCount >= 2 {
		for i := range playerCount {
			tip := armIndices[i][len(armIndices[i])-1]
			nextTip := armIndices[(i+1)%playerCount][len(armIndices[(i+1)%playerCount])-1]
			builder.add(tip, nextTip)
		}
	}
	return builder.pairs
}
