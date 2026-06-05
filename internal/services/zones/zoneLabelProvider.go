package zones

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
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
	configuration config.GeneratorConfig) models.NeutralZonePlans {
	var plans models.NeutralZonePlans
	maxNeutral := max(0, len(this.zoneLabels)-configuration.PlayerCount)
	castleZoneCastleCount := helpers.Clamp(configuration.ZoneConfiguration.NeutralZoneCastles, 1, 4)

	add := func(requested int, quality models.NeutralZoneQuality, castleCount int) {
		count := helpers.Clamp(requested, 0, 30)                // TODO: Clamp up to labelCount - playerCount
		for i := 0; i < count && len(plans) < maxNeutral; i++ { // TODO: Is plans length needed?
			plans.AddPlan(this.zoneLabels[configuration.PlayerCount+len(plans)], quality, castleCount)
		}
	}

	advanced := configuration.ZoneConfiguration.Advanced
	advancedTotal := advanced.NeutralLowNoCastleCount + advanced.NeutralLowCastleCount +
		advanced.NeutralMediumNoCastleCount + advanced.NeutralMediumCastleCount +
		advanced.NeutralHighNoCastleCount + advanced.NeutralHighCastleCount
	if advancedTotal > 0 {
		add(advanced.NeutralLowNoCastleCount, models.QualityLow, 0)
		add(advanced.NeutralLowCastleCount, models.QualityLow, castleZoneCastleCount)
		add(advanced.NeutralMediumNoCastleCount, models.QualityMedium, 0)
		add(advanced.NeutralMediumCastleCount, models.QualityMedium, castleZoneCastleCount)
		add(advanced.NeutralHighNoCastleCount, models.QualityHigh, 0)
		add(advanced.NeutralHighCastleCount, models.QualityHigh, castleZoneCastleCount)
	} else {
		castleCount := helpers.Clamp(configuration.ZoneConfiguration.NeutralZoneCastles, 0, 4)
		add(configuration.ZoneConfiguration.NeutralZoneCount, models.QualityMedium, castleCount)
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
	neutralZones models.NeutralZonePlans) string {
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
	neutralZones []models.NeutralZonePlan,
	isRing bool) []string {
	neutralLabels := linq.FromSlice(neutralZones).
		SelectString(func(x models.NeutralZonePlan) string { return x.Label }).
		ToSlice()

	if configuration.Topology == config.TopologyBalanced {
		separationCount := 0
		if configuration.MinNeutralZonesBetweenPlayers > 0 && configuration.CanHonorNeutralSeparation() {
			separationCount = configuration.MinNeutralZonesBetweenPlayers
		}
		if isRing {
			return this.CreateBalancedRingZoneLabels(playerLabels, neutralZones, separationCount)
		}
		return this.CreateBalancedChainZoneLabels(playerLabels, neutralZones, separationCount)
	}

	separationCount := configuration.MinNeutralZonesBetweenPlayers
	if separationCount <= 0 || configuration.RandomPortals || !configuration.CanHonorNeutralSeparation() {
		return append(append([]string{}, playerLabels...), neutralLabels...)
	}

	var orderedLabels []string
	queue := make([]string, len(neutralLabels))
	copy(queue, neutralLabels)
	qi := 0
	for i, pl := range playerLabels {
		orderedLabels = append(orderedLabels, pl)
		needsSep := isRing || i < len(playerLabels)-1
		if !needsSep {
			continue
		}
		for j := 0; j < separationCount && qi < len(queue); j++ {
			orderedLabels = append(orderedLabels, queue[qi])
			qi++
		}
	}
	for qi < len(queue) {
		orderedLabels = append(orderedLabels, queue[qi])
		qi++
	}
	if len(orderedLabels) == 0 {
		return append(append([]string{}, playerLabels...), neutralLabels...)
	}

	return orderedLabels
}

func (this *ZoneLabelProvider) CreateBalancedRingZoneLabels(
	playerLetters []string,
	neutralZones []models.NeutralZonePlan,
	separationCount int) []string {
	if len(playerLetters) == 0 {
		return this.CreateBalancedNeutralRingZoneLabels(neutralZones, 1)
	}

	if len(neutralZones) == 0 {
		return append([]string{}, playerLetters...)
	}

	caps := utils.GetEvenGapCapacities(len(playerLetters), len(neutralZones), separationCount)
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, caps, false)
	var ordered []string
	for i, pl := range playerLetters {
		ordered = append(ordered, pl)
		for _, nz := range utils.OrderNeutralsWithinGap(gaps[i]) {
			ordered = append(ordered, nz.Label)
		}
	}
	return ordered
}

func (this *ZoneLabelProvider) CreateBalancedChainZoneLabels(
	playerLabels []string,
	neutralZones []models.NeutralZonePlan,
	minSep int) []string {
	if len(playerLabels) == 0 {
		return linq.FromSlice(neutralZones).
			SelectString(func(x models.NeutralZonePlan) string { return x.Label }).
			ToSlice()
	}

	gapCount := len(playerLabels) + 1
	capacities := make([]int, gapCount)
	remaining := len(neutralZones)
	reqInterior := max(0, len(playerLabels)-1) * minSep
	if minSep > 0 && len(neutralZones) >= reqInterior {
		for i := 1; i < gapCount-1; i++ {
			capacities[i] = minSep
		}
		remaining -= reqInterior
	}
	// Distribute extra neutrals only into interior gaps so that the first
	// and last positions of the chain are always player zones. Degenerate cases (0 or 1
	// player) fall back to even distribution across every gap
	interiorGapCount := max(0, gapCount-2)
	if interiorGapCount > 0 {
		extras := utils.GetEvenGapCapacities(interiorGapCount, remaining, 0)
		for i := 1; i < gapCount-1; i++ {
			capacities[i] += extras[i-1]
		}
	} else {
		extras := utils.GetEvenGapCapacities(gapCount, remaining, 0)
		for i := range gapCount {
			capacities[i] += extras[i]
		}
	}
	gaps := utils.AssignNeutralZonesToGaps(neutralZones, capacities, true)
	orderedLabels := linq.FromSlice(utils.OrderEdgeGap(gaps[0], true)).
		SelectString(func(x models.NeutralZonePlan) string { return x.Label }).
		ToSlice()
	for i, pl := range playerLabels {
		orderedLabels = append(orderedLabels, pl)
		gap := gaps[i+1]
		trailing := i == len(playerLabels)-1
		var g []models.NeutralZonePlan
		if trailing {
			g = utils.OrderEdgeGap(gap, false)
		} else {
			g = utils.OrderNeutralsWithinGap(gap)
		}
		for _, nz := range g {
			orderedLabels = append(orderedLabels, nz.Label)
		}
	}
	if len(orderedLabels) == 0 {
		nl := make([]string, len(neutralZones))
		for i, nz := range neutralZones {
			nl[i] = nz.Label
		}
		return append(append([]string{}, playerLabels...), nl...)
	}

	return orderedLabels
}

func (this *ZoneLabelProvider) CreateBalancedNeutralRingZoneLabels(
	neutralZones []models.NeutralZonePlan,
	playerCount int) []string {
	if len(neutralZones) < 2 {
		labels := make([]string, len(neutralZones))
		for index, zonePlan := range neutralZones {
			labels[index] = zonePlan.Label
		}
		return labels
	}

	caps := utils.GetEvenGapCapacities(max(1, playerCount), len(neutralZones), 0)
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
	neutralZones []models.NeutralZonePlan) models.ZoneAdjacency {
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
	case config.TopologyDefault, config.TopologyBalanced:
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
