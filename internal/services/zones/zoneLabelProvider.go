package zones

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/utils"
)

type ZoneLabelProvider struct {
	zoneLabels []string
}

func NewZoneLabelProvider() *ZoneLabelProvider {
	return &ZoneLabelProvider{
		zoneLabels: []string{
			"A", "B", "C", "D", "E", "F", "G", "H",
			"I", "J", "K", "L", "M", "N", "O", "P",
			"Q", "R", "S", "T", "U", "V", "W", "X",
			"Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF",
		},
	}
}

func (this *ZoneLabelProvider) CreatePlayerLabels(playerCount int) []string {
	letters := make([]string, playerCount)
	copy(letters, this.zoneLabels[:playerCount])
	return letters
}

func (this *ZoneLabelProvider) CreateNeutralZonePlans(
	configuration config.GeneratorConfig) neutralZone.Plans {
	var plans neutralZone.Plans
	maxNeutral := max(0, len(this.zoneLabels)-configuration.PlayerCount)

	add := func(requested int, quality neutralZone.Quality, castleCount int) {
		count := helpers.Clamp(requested, 0, maxNeutral)
		// The plans-length guard enforces the label capacity cumulatively across add calls.
		for i := 0; i < count && len(plans) < maxNeutral; i++ {
			plans.AddPlan(this.zoneLabels[configuration.PlayerCount+len(plans)], quality, castleCount)
		}
	}

	advanced := configuration.ZoneConfiguration.Advanced
	advancedTotal := advanced.NeutralLowNoCastleCount + advanced.NeutralLowCastleCount +
		advanced.NeutralMediumNoCastleCount + advanced.NeutralMediumCastleCount +
		advanced.NeutralHighNoCastleCount + advanced.NeutralHighCastleCount
	if advancedTotal > 0 {
		lowCastlesPerZone := helpers.Clamp(advanced.NeutralLowCastlesPerZone, 0, 4)
		medCastlesPerZone := helpers.Clamp(advanced.NeutralMediumCastlesPerZone, 0, 4)
		highCastlesPerZone := helpers.Clamp(advanced.NeutralHighCastlesPerZone, 0, 4)
		add(advanced.NeutralLowNoCastleCount, neutralZone.QualityLow, 0)
		add(advanced.NeutralLowCastleCount, neutralZone.QualityLow, lowCastlesPerZone)
		add(advanced.NeutralMediumNoCastleCount, neutralZone.QualityMedium, 0)
		add(advanced.NeutralMediumCastleCount, neutralZone.QualityMedium, medCastlesPerZone)
		add(advanced.NeutralHighNoCastleCount, neutralZone.QualityHigh, 0)
		add(advanced.NeutralHighCastleCount, neutralZone.QualityHigh, highCastlesPerZone)
	} else {
		castleCount := helpers.Clamp(configuration.ZoneConfiguration.NeutralZoneCastles, 0, 4)
		add(configuration.ZoneConfiguration.NeutralZoneCount, neutralZone.QualityMedium, castleCount)
	}
	if configuration.Topology == config.TopologySharedWeb && len(plans) == 0 && maxNeutral > 0 {
		plans.AddMediumPlan(
			this.zoneLabels[configuration.PlayerCount],
			helpers.Clamp(configuration.ZoneConfiguration.NeutralZoneCastles, 0, 4))
	}
	return plans
}

func (this *ZoneLabelProvider) GetHoldCityLabel(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutralZone.Plans) string {
	if !neutralZones.Any() || !configuration.IsHubCityToHold() {
		return ""
	}

	adjacency := this.createTopologyAdjacency(configuration, playerLabels, neutralZones)
	var distancesByPlayer []map[string]int
	for _, label := range playerLabels {
		distancesByPlayer = append(distancesByPlayer, adjacency.GetDistancesFrom(label))
	}

	return utils.CreateHubZoneCandidates(neutralZones, distancesByPlayer).
		SortForHubCity().
		GetFirstCandidateLabel()
}

func (this *ZoneLabelProvider) CreateZoneName(label string, playerLabels []string) string {
	if slices.Contains(playerLabels, label) {
		return "Spawn-" + label
	}
	return "Neutral-" + label
}

func (this *ZoneLabelProvider) CreateOrderedZoneLabels(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutralZone.Plans,
	isRing bool) []string {
	if configuration.Topology == config.TopologyCircles {
		if isRing {
			return this.CreateBalancedRingZoneLabels(playerLabels, neutralZones)
		}
		return this.CreateBalancedChainZoneLabels(playerLabels, neutralZones)
	}

	neutralLabels := linq.FromSlice(neutralZones).
		SelectString(func(x neutralZone.Plan) string { return x.Label }).
		ToSlice()
	return append(playerLabels, neutralLabels...)
}

func (this *ZoneLabelProvider) CreateBalancedRingZoneLabels(
	playerLabels []string,
	neutralZones neutralZone.Plans) []string {
	if len(playerLabels) == 0 {
		return this.CreateBalancedNeutralRingZoneLabels(neutralZones, 1)
	}

	if len(neutralZones) == 0 {
		return playerLabels
	}

	caps := utils.GetEvenGapCapacities(len(playerLabels), len(neutralZones))
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, caps, false)
	var ordered []string
	for i, playerLabel := range playerLabels {
		ordered = append(ordered, playerLabel)
		for _, nz := range utils.OrderNeutralsWithinGap(gaps[i]) {
			ordered = append(ordered, nz.Label)
		}
	}
	return ordered
}

func (this *ZoneLabelProvider) CreateBalancedChainZoneLabels(
	playerLabels []string,
	neutralZones neutralZone.Plans) []string {
	if len(playerLabels) == 0 {
		return linq.FromSlice(neutralZones).
			SelectString(func(x neutralZone.Plan) string { return x.Label }).
			ToSlice()
	}

	gapCount := len(playerLabels) + 1
	capacities := make([]int, gapCount)
	remaining := len(neutralZones)
	// Distribute extra neutrals only into interior gaps so that the first
	// and last positions of the chain are always player zones. Degenerate cases (0 or 1
	// player) fall back to even distribution across every gap
	interiorGapCount := max(0, gapCount-2)
	if interiorGapCount > 0 {
		extras := utils.GetEvenGapCapacities(interiorGapCount, remaining)
		for i := 1; i < gapCount-1; i++ {
			capacities[i] += extras[i-1]
		}
	} else {
		extras := utils.GetEvenGapCapacities(gapCount, remaining)
		for i := range gapCount {
			capacities[i] += extras[i]
		}
	}
	neutralZoneGaps := utils.AssignNeutralZonesToGaps(neutralZones, capacities, true)
	orderedLabels := linq.FromSlice(utils.OrderEdgeGap(neutralZoneGaps[0], true)).
		SelectString(func(x neutralZone.Plan) string { return x.Label }).
		ToSlice()
	for index, playerLabel := range playerLabels {
		orderedLabels = append(orderedLabels, playerLabel)
		neutralZoneGap := neutralZoneGaps[index+1]
		trailing := index == len(playerLabels)-1
		var gap neutralZone.Plans
		if trailing {
			gap = utils.OrderEdgeGap(neutralZoneGap, false)
		} else {
			gap = utils.OrderNeutralsWithinGap(neutralZoneGap)
		}
		for _, zonePlan := range gap {
			orderedLabels = append(orderedLabels, zonePlan.Label)
		}
	}
	return orderedLabels
}

func (this *ZoneLabelProvider) CreateBalancedNeutralRingZoneLabels(
	neutralZones neutralZone.Plans,
	playerCount int) []string {
	if len(neutralZones) < 2 {
		labels := make([]string, len(neutralZones))
		for index, zonePlan := range neutralZones {
			labels[index] = zonePlan.Label
		}
		return labels
	}

	caps := utils.GetEvenGapCapacities(max(1, playerCount), len(neutralZones))
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, caps, false)
	var labels []string
	for _, gap := range gaps {
		for _, zonePlan := range utils.OrderNeutralsWithinGap(gap) {
			labels = append(labels, zonePlan.Label)
		}
	}
	return labels
}

func (this *ZoneLabelProvider) createTopologyAdjacency(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutralZone.Plans) models.ZoneAdjacency {
	adjacency := models.ZoneAdjacency{}

	isIsolated := configuration.NoDirectPlayerConnections && len(playerLabels) > 1

	switch configuration.Topology {
	case config.TopologyChain:
		orderedLabels := this.CreateOrderedZoneLabels(configuration, playerLabels, neutralZones, false)
		for current := range len(orderedLabels) - 1 {
			next := current + 1
			if isIsolated &&
				slices.Contains(playerLabels, orderedLabels[current]) &&
				slices.Contains(playerLabels, orderedLabels[next]) {
				continue
			}
			adjacency.Link(orderedLabels[current], orderedLabels[next])
		}
	case config.TopologyRing, config.TopologyCircles:
		orderedLabels := this.CreateOrderedZoneLabels(configuration, playerLabels, neutralZones, true)
		for current := range orderedLabels {
			next := (current + 1) % len(orderedLabels)
			if isIsolated &&
				slices.Contains(playerLabels, orderedLabels[current]) &&
				slices.Contains(playerLabels, orderedLabels[next]) {
				continue
			}
			adjacency.Link(orderedLabels[current], orderedLabels[next])
		}
	default:
		orderedLabels := this.CreateOrderedZoneLabels(configuration, playerLabels, neutralZones, true)
		for current := range orderedLabels {
			next := (current + 1) % len(orderedLabels)
			adjacency.Link(orderedLabels[current], orderedLabels[next])
		}
	}

	return adjacency
}
