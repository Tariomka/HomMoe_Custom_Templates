package topology

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/placement_rule"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

var resourceContentPool = registry.GetResourcesContentPoolValues()

type topologyBase struct {
	zoneLabelProvider *zones.ZoneLabelProvider
}

func newTopologyBase() topologyBase {
	return topologyBase{
		zoneLabelProvider: zones.NewZoneLabelProvider(),
	}
}

func (this *topologyBase) CreateVariant(
	playerLabels []string,
	firstLabel string,
	zoneCount int,
	zones []template.Zone,
	connections []template.Connection) template.Variant {
	orientationBuilder := variant_content.NewOrientationBuilder().
		WithBaseAngleMin(45).
		WithBaseAngleMax(45).
		WithRandomAngleAmplitude(360)

	if zoneCount > 0 {
		orientationBuilder.WithRandomAngleStep(360 / zoneCount)
	}

	if slices.Contains(playerLabels, firstLabel) {
		orientationBuilder.WithZeroAngleZone("Spawn-" + firstLabel)
	} else {
		orientationBuilder.WithZeroAngleZone("Neutral-" + firstLabel)
	}

	return variant_content.NewVariantBuilder().
		WithOrientation(orientationBuilder.Build()).
		WithBorder(variant_content.NewBorderBuilder().
			WithCornerRadius(0).
			WithObstaclesWidth(3).
			WithObstaclesNoise(1, 12).
			WithWaterWidth(0).
			WithWaterNoise(1, 12).
			WithWaterTypeWaterGrass().
			Build()).
		WithZones(zones...).
		WithConnections(connections...).
		Build()
}

func (this *topologyBase) CreateSpawnZone(
	label, playerName string,
	connectionNames []string,
	castleCount int,
	matchFactions bool,
	zoneSize float64,
	spawnFootholds,
	generateRoads bool,
	tuning models.GenerationTuning) template.Zone {
	mainObjects := []template.MainObject{
		this.createPlayerSpawnCastle(playerName, tuning.ScaleByNeutralGuardStrength(5000)),
	}
	// TODO: add player owned castles
	mainObjects = append(mainObjects,
		this.createPlayerUnclaimedCastles(matchFactions, tuning.ScaleByNeutralGuardStrength(2500), castleCount-1)...)

	return variant_content.NewZoneBuilder().
		WithName("Spawn-" + label).
		WithSize(utils.NormalizeZoneSize(zoneSize)).
		WithLayoutSpawns().
		WithGuardCutoffValue(2000).
		WithGuardRandomization(tuning.GuardRandomization).
		WithGuardMultiplier(tuning.ScaleByNeutralGuardStrengthPrecise(1.0)).
		WithGuardWeeklyIncrement(0.20).
		WithGuardReactionDistribution([]int{60, 20, 10, 10, 2, 0}).
		WithDiplomacyModifier(-0.5).
		WithGuardedContentPool(registry.GetGuardedContentPoolT2List()).
		WithUnguardedContentPool(registry.GetUnguardedContentPoolT2List()).
		WithResourcesContentPool([]string{resourceContentPool.StartZonePoor}).
		WithMandatoryContent("mandatory_content_side_" + label).
		WithContentCountLimits(buildSideContentLimits()).
		WithGuardedContentValue(tuning.ScaleByStructureDensity(200000 * tuning.ContentScale)).
		WithGuardedContentValuePerArea(tuning.ScaleByStructureDensity(2000 * math.Sqrt(tuning.ContentScale))).
		WithUnguardedContentValue(tuning.ScaleByStructureDensity(50000 * tuning.ContentScale)).
		WithUnguardedContentValuePerArea(tuning.ScaleByStructureDensity(400 * math.Sqrt(tuning.ContentScale))).
		WithResourcesValue(tuning.ScaleByResourceDensity(80000 * tuning.ContentScale)).
		WithResourcesValuePerArea(tuning.ScaleByResourceDensity(600 * math.Sqrt(tuning.ContentScale))).
		WithMainObjects(mainObjects).
		WithBiomeMatchMainObject("0").
		WithCrossroadsPosition(0).
		WithRoads(this.createOuterZoneRoads(connectionNames, castleCount, spawnFootholds, generateRoads)).
		Build()
}

func (this *topologyBase) CreateNeutralZone(
	plan models.NeutralZonePlan,
	connectionNames []string,
	zoneSize float64,
	spawnFootholds, generateRoads bool,
	tuning models.GenerationTuning,
	isHoldCity bool) template.Zone {
	if isHoldCity && plan.CastleCount < 1 {
		plan.CastleCount = 1
	}
	profile := models.NewNeutralZoneProfile(plan.Quality)

	zoneBuilder := variant_content.NewZoneBuilder().
		WithName("Neutral-" + plan.Label).
		WithSize(utils.NormalizeZoneSize(zoneSize)).
		WithLayout(profile.Layout).
		WithGuardCutoffValue(2000).
		WithGuardRandomization(tuning.GuardRandomization).
		WithGuardMultiplier(tuning.ScaleByNeutralGuardStrengthPrecise(profile.GuardMultiplier)).
		WithGuardWeeklyIncrement(0.20).
		WithGuardReactionDistribution(profile.GuardReactionDistribution).
		WithDiplomacyModifier(-0.5).
		WithGuardedContentPool(profile.GuardedContentPool).
		WithUnguardedContentPool(profile.UnguardedContentPool).
		WithResourcesContentPool(profile.ResourcesContentPool).
		WithMandatoryContent("mandatory_content_neutral_" + plan.Label).
		WithContentCountLimits(buildSideContentLimits()).
		WithGuardedContentValue(tuning.ScaleByStructureDensity(float64(profile.GuardedContentValue) * tuning.ContentScale)).
		WithGuardedContentValuePerArea(tuning.ScaleByStructureDensity(float64(profile.GuardedContentValuePerArea) * math.Sqrt(tuning.ContentScale))).
		WithUnguardedContentValue(tuning.ScaleByStructureDensity(float64(profile.UnguardedContentValue) * tuning.ContentScale)).
		WithUnguardedContentValuePerArea(tuning.ScaleByStructureDensity(float64(profile.UnguardedContentValuePerArea) * math.Sqrt(tuning.ContentScale))).
		WithResourcesValue(tuning.ScaleByResourceDensity(float64(profile.ResourcesValue) * tuning.ContentScale)).
		WithResourcesValuePerArea(tuning.ScaleByResourceDensity(float64(profile.ResourcesValuePerArea) * math.Sqrt(tuning.ContentScale))).
		WithMainObjects(this.createNeutralZoneCastles(profile, tuning, plan.CastleCount, isHoldCity)).
		WithCrossroadsPosition(0).
		WithRoads(this.createOuterZoneRoads(connectionNames, plan.CastleCount, spawnFootholds, generateRoads))

	if plan.CastleCount > 0 {
		zoneBuilder = zoneBuilder.WithBiomeMatchMainObject("0")
	} else {
		zoneBuilder = zoneBuilder.WithBiomeMatchZone()
	}

	return zoneBuilder.Build()
}

func (this *topologyBase) CreateHubZone(
	connectionNames []string,
	tuning models.GenerationTuning,
	isHoldCity bool,
	size float64,
	castleCount int,
	generateRoads bool) template.Zone {
	if isHoldCity && castleCount < 1 {
		castleCount = 1
	}

	zoneBuilder := variant_content.NewZoneBuilder().
		WithName("Hub").
		WithSize(utils.NormalizeZoneSize(size)).
		WithLayoutCenter().
		WithGuardCutoffValue(2000).
		WithGuardRandomization(0.05).
		WithGuardMultiplier(tuning.ScaleByNeutralGuardStrengthPrecise(1.5)).
		WithGuardWeeklyIncrement(0.20).
		WithGuardReactionDistribution([]int{0, 10, 10, 20, 10, 0}).
		WithDiplomacyModifier(-0.5).
		WithGuardedContentPool(registry.GetGuardedContentPoolT3List()).
		WithUnguardedContentPool(registry.GetUnguardedContentPoolT3List()).
		WithResourcesContentPool([]string{resourceContentPool.StartZoneMedium}).
		WithContentCountLimits(buildSideContentLimits()).
		WithGuardedContentValue(tuning.ScaleByStructureDensity(300000 * tuning.ContentScale)).
		WithGuardedContentValuePerArea(tuning.ScaleByStructureDensity(2400 * math.Sqrt(tuning.ContentScale))).
		WithUnguardedContentValue(tuning.ScaleByStructureDensity(50000 * tuning.ContentScale)).
		WithUnguardedContentValuePerArea(tuning.ScaleByStructureDensity(600 * math.Sqrt(tuning.ContentScale))).
		WithResourcesValue(tuning.ScaleByResourceDensity(80000 * tuning.ContentScale)).
		WithResourcesValuePerArea(tuning.ScaleByResourceDensity(600 * math.Sqrt(tuning.ContentScale))).
		WithMainObjects(this.createHubZoneCastles(tuning, castleCount, isHoldCity)).
		WithCrossroadsPosition(0).
		WithRoads(this.createOuterZoneRoads(connectionNames, castleCount, false, generateRoads))

	if castleCount > 0 {
		zoneBuilder = zoneBuilder.WithBiomeMatchMainObject("0")
	} else {
		zoneBuilder = zoneBuilder.WithBiomeMatchZone()
	}

	return zoneBuilder.Build()
}

func (this *topologyBase) CreateRandomPortalConnections(
	playerLabels, orderedLabels []string,
	tuning models.GenerationTuning,
	maxCount int) []template.Connection {
	count := len(orderedLabels)
	if count < 2 {
		return nil
	}
	dest := buildNonAdjacentDerangement(count)
	indices := make([]int, count)
	for i := range indices {
		indices[i] = i
	}
	rand.Shuffle(len(indices), func(i, j int) { indices[i], indices[j] = indices[j], indices[i] })

	limit := min(count, maxCount)
	trueVal := true
	rule := placement_rule.NewPlacementRuleBuilder().BuildCrossroadsRule(placement_rule.DistanceNear, 2)
	var conns []template.Connection
	for i := range limit {
		idx := indices[i]
		from := orderedLabels[idx]
		to := orderedLabels[dest[idx]]
		fromZone := createZoneName(from, playerLabels)
		toZone := createZoneName(to, playerLabels)
		conns = append(conns, template.Connection{
			Name: fmt.Sprintf("Portal-%s-%s", from, to), From: fromZone, To: toZone,
			ConnectionType:           "Portal",
			PortalPlacementRulesFrom: []template.PlacementRule{rule},
			PortalPlacementRulesTo:   []template.PlacementRule{rule},
			Road:                     &trueVal, GuardValue: tuning.ScaleByBorderGuardStrength(25000), GuardWeeklyIncrement: 0.15,
		})
	}
	return conns
}

func (this *topologyBase) GetBorderGuardValue(
	labelA, labelB string,
	playerLabels []string,
	neutralZones models.NeutralZonePlans,
	tuning models.GenerationTuning) int {
	aIsPlayer := slices.Contains(playerLabels, labelA)
	bIsPlayer := slices.Contains(playerLabels, labelB)
	if aIsPlayer && bIsPlayer {
		return tuning.ScaleByBorderGuardStrength(30_000)
	}

	if !aIsPlayer && !bIsPlayer {
		qa := neutralZones.GetQuality(labelA)
		qb := neutralZones.GetQuality(labelB)
		higher := qa
		if int(qb) > int(qa) {
			higher = qb
		}
		return tuning.ScaleByBorderGuardStrength(higher.GetGuardValue())
	}

	neutralLabel := labelB
	if !aIsPlayer {
		neutralLabel = labelA
	}
	return tuning.ScaleByBorderGuardStrength(neutralZones.GetQuality(neutralLabel).GetGuardValue())
}

func (this *topologyBase) EnsurePlayerZonesConnected(
	playerLabels []string,
	zones []template.Zone,
	connections *[]template.Connection,
	tuning models.GenerationTuning) {
	if len(playerLabels) < 2 {
		return
	}
	connNames := map[string]bool{}
	for _, c := range *connections {
		if c.Name != "" {
			connNames[c.Name] = true
		}
	}
	for _, letter := range playerLabels {
		zn := "Spawn-" + letter
		z, ok := linq.FromSlice(zones).First(func(z template.Zone) bool { return z.Name == zn })
		if !ok {
			continue
		}
		hasConn := false
		for _, r := range z.Roads {
			if r.To.Type == "Connection" && len(r.To.Args) > 0 && connNames[r.To.Args[0]] {
				hasConn = true
				break
			}
		}
		if hasConn {
			continue
		}
		var partner string
		for _, pl := range playerLabels {
			if pl != letter {
				partner = pl
				break
			}
		}
		if partner == "" {
			continue
		}
		a, b := letter, partner
		if a > b {
			a, b = b, a
		}
		fn := "Fallback-" + a + "-" + b
		if connNames[fn] {
			continue
		}
		*connections = append(*connections, template.Connection{
			Name: fn, From: "Spawn-" + letter, To: "Spawn-" + partner,
			ConnectionType: "Direct", GuardZone: "Spawn-" + letter, SimTurnSquad: true,
			GuardValue: this.GetBorderGuardValue(letter, partner, playerLabels, nil, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: "fallback_guard_" + fn,
		})
		connNames[fn] = true
		for _, pl := range []string{letter, partner} {
			if pz, ok := linq.FromSlice(zones).First(func(z template.Zone) bool { return z.Name == "Spawn-"+pl }); ok {
				pz.Roads = append(pz.Roads, variant_content.NewRoadBuilder().
					WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).
					WithTo(variant_content.NewRefBuilder().BuildConnectionType(fn)).
					Build())
			}
		}
	}
}

func (this *topologyBase) EnsureFullConnectivity(
	playerLabels, allLabels []string,
	positions models.Positions,
	zones []template.Zone,
	connections []template.Connection,
	tuning models.GenerationTuning,
	neutralZones models.NeutralZonePlans) []template.Connection {
	if len(allLabels) <= 1 {
		return connections
	}

	adjacency := models.NewZoneIndexAdjacency(len(allLabels))
	// TODO: move out to a separate function
	zoneNameToIdx := map[string]int{}
	for i, l := range allLabels {
		zoneNameToIdx[createZoneName(l, playerLabels)] = i
	}
	for _, connection := range connections {
		if connection.ConnectionType != "Direct" && connection.ConnectionType != "Portal" {
			continue
		}
		indexA, okA := zoneNameToIdx[connection.From]
		indexB, okB := zoneNameToIdx[connection.To]
		if !okA || !okB {
			continue
		}
		adjacency.Link(indexA, indexB)
	}

	connNameSet := map[string]bool{}
	for _, connection := range connections {
		if connection.Name != "" {
			connNameSet[connection.Name] = true
		}
	}

	for {
		components := adjacency.FindIndexes(len(allLabels))
		bestIndexA, bestIndexB, ok := positions.GetShortestDistanceIndex(components)
		if !ok {
			break
		}

		labelA, labelB := allLabels[bestIndexA], allLabels[bestIndexB]
		if labelA > labelB {
			labelA, labelB = labelB, labelA
		}
		bridgeName := fmt.Sprintf("Bridge-%s-%s", labelA, labelB)
		if connNameSet[bridgeName] {
			continue
		}

		zoneFrom := createZoneName(allLabels[bestIndexA], playerLabels)
		zoneTo := createZoneName(allLabels[bestIndexB], playerLabels)
		connections = append(connections, template.Connection{
			Name: bridgeName, From: zoneFrom, To: zoneTo,
			ConnectionType: "Direct", GuardZone: zoneFrom, SimTurnSquad: true,
			GuardValue: this.GetBorderGuardValue(labelA, labelB, playerLabels, neutralZones, tuning), GuardWeeklyIncrement: 0.15,
			GuardMatchGroup: fmt.Sprintf("bridge_guard_%s-%s", labelA, labelB),
		})
		connNameSet[bridgeName] = true
		for _, zoneName := range []string{zoneFrom, zoneTo} {
			if zone, ok := linq.FromSlice(zones).First(func(x template.Zone) bool { return x.Name == zoneName }); ok {
				roadBuilder := variant_content.NewRoadBuilder().
					WithTo(variant_content.NewRefBuilder().BuildConnectionType(bridgeName))
				if len(zone.MainObjects) > 0 {
					zone.Roads = append(zone.Roads,
						roadBuilder.WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).Build())
				} else if len(zone.Roads) > 0 {
					existingConn := ""
					for _, r := range zone.Roads {
						if r.From.Type == "Connection" && len(r.From.Args) > 0 {
							existingConn = r.From.Args[0]
							break
						}
						if r.To.Type == "Connection" && len(r.To.Args) > 0 {
							existingConn = r.To.Args[0]
							break
						}
					}
					if existingConn != "" {
						zone.Roads = append(zone.Roads,
							roadBuilder.WithFrom(variant_content.NewRefBuilder().BuildConnectionType(existingConn)).Build())
					} else {
						zone.Roads = append(zone.Roads,
							roadBuilder.WithFrom(variant_content.NewRefBuilder().BuildConnectionType(bridgeName)).Build())
					}
				} else {
					zone.Roads = append(zone.Roads,
						roadBuilder.WithFrom(variant_content.NewRefBuilder().BuildConnectionType(bridgeName)).Build())
				}
			}
		}

		adjacency.Link(bestIndexA, bestIndexB)
	}

	return connections
}

func (this *topologyBase) createPlayerSpawnCastle(playerName string, guardValue int) template.MainObject {
	return variant_content.NewObjectBuilder().
		WithTypeSpawn().
		WithSpawn(playerName).
		WithNoGuardWhenOwned().
		WithGuardChance(1).
		WithGuardValue(guardValue).
		WithGuardWeeklyIncrement(0.10).
		WithCastleQualityDefault().
		WithPlacementUniform().
		WithPlacementArgs("true", "0.7", "0").
		Build()
}

// func (this *topologyBase) createPlayerOwnedCastles(playerIndex int, guardValue int, castleCount int) []template.MainObject {  }

func (this *topologyBase) createPlayerUnclaimedCastles(
	matchPlayerFaction bool,
	guardValue, castleCount int) []template.MainObject {
	var castles []template.MainObject
	for range castleCount {
		objectBuilder := variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardChance(1).
			WithGuardValue(guardValue).
			WithGuardWeeklyIncrement(0.10).
			WithCastleQualityPoor().
			WithPlacementUniform().
			WithPlacementArgs("false", "-0.8", "3")
		if matchPlayerFaction {
			objectBuilder = objectBuilder.WithFaction("Match", "0")
		} else {
			objectBuilder = objectBuilder.WithFaction("Random") // TODO: is this valid?
		}
		castles = append(castles, objectBuilder.Build())
	}
	return castles
}

func (this *topologyBase) createNeutralZoneCastles(
	profile models.NeutralZoneProfile,
	tuning models.GenerationTuning,
	castleCount int,
	isHoldCityZone bool) []template.MainObject {
	var castles []template.MainObject

	if castleCount > 0 {
		objectBuilder := variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardChance(1).
			WithGuardWeeklyIncrement(0.10).
			WithFaction("FromList")
		if isHoldCityZone {
			objectBuilder = objectBuilder.
				WithGuardValue(tuning.ScaleByBorderGuardStrength(max(profile.PrimaryCityGuardValue, 20_000))).
				WithCastleQualityUltraRich().
				WithPlacementCenter().
				WithHoldCityWinCon()
		} else {
			objectBuilder = objectBuilder.
				WithGuardValue(tuning.ScaleByBorderGuardStrength(profile.PrimaryCityGuardValue)).
				WithCastleQuality(profile.PrimaryBuildingsCSid).
				WithPlacementUniform().
				WithPlacementArgs("true", "0.8", "2")
		}
		castles = append(castles, objectBuilder.Build())
	}

	for range castleCount - 1 {
		castles = append(castles,
			variant_content.NewObjectBuilder().
				WithTypeCity().
				WithGuardChance(1).
				WithGuardValue(tuning.ScaleByBorderGuardStrength(profile.ExtraCityGuardValue)).
				WithGuardWeeklyIncrement(0.10).
				WithCastleQuality(profile.ExtraBuildingsCSid).
				WithFaction("FromList").
				WithPlacementUniform().
				WithPlacementArgs("false", "-0.8", "3").
				Build())
	}

	return castles
}

func (this *topologyBase) createHubZoneCastles(
	tuning models.GenerationTuning,
	castleCount int,
	isHoldCityZone bool) []template.MainObject {
	var castles []template.MainObject
	newCastleBuilder := func() *variant_content.MainObjectBuilder {
		return variant_content.NewObjectBuilder().
			WithTypeCity().
			WithGuardWeeklyIncrement(0.10).
			WithFaction("FromList")
	}
	buildHoldCityCastle := func(builder *variant_content.MainObjectBuilder) template.MainObject {
		return builder.
			WithGuardChance(1).
			WithGuardValue(tuning.ScaleByNeutralGuardStrength(25_000)).
			WithCastleQualityUltraRich().
			WithPlacementCenter().
			WithHoldCityWinCon().
			Build()
	}
	buildCastle := func(builder *variant_content.MainObjectBuilder) template.MainObject {
		return builder.
			WithGuardChance(0.5).
			WithGuardValue(tuning.ScaleByNeutralGuardStrength(16_000)).
			WithCastleQualityRich().
			WithPlacementUniform().
			WithPlacementArgs("true", "0.8", "2").
			Build()
	}

	if castleCount > 0 && isHoldCityZone {
		castles = append(castles, buildHoldCityCastle(newCastleBuilder()))
	} else if castleCount > 0 {
		castles = append(castles, buildCastle(newCastleBuilder()))
	}

	for range castleCount - 1 {
		castles = append(castles, buildCastle(newCastleBuilder()))
	}

	return castles
}

func (this *topologyBase) createOuterZoneRoads(
	connectionNames []string,
	castleCount int,
	includeFoothold, generateRoads bool) []template.Road {
	if !generateRoads {
		return nil
	}

	if castleCount == 0 {
		return this.createConnectorZoneRoads(connectionNames, generateRoads)
	}

	var roads []template.Road
	for i := 1; i < castleCount; i++ {
		roads = append(roads,
			variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).
				WithTo(variant_content.NewRefBuilder().BuildMainObjectType(fmt.Sprintf("%d", i))).
				Build())
	}
	if includeFoothold {
		roads = append(roads,
			variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).
				WithTo(variant_content.NewRefBuilder().BuildMandatoryContentType("name_remote_foothold_1")).
				Build())
	}
	for _, name := range connectionNames {
		roads = append(roads,
			variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).
				WithTo(variant_content.NewRefBuilder().BuildConnectionType(name)).
				Build())
	}
	return roads
}

func (this *topologyBase) createConnectorZoneRoads(connectionNames []string, generateRoads bool) []template.Road {
	if !generateRoads {
		return nil
	}

	distinctNames := helpers.GetUniqueElements(connectionNames)
	if len(distinctNames) == 0 {
		return nil
	}

	if len(distinctNames) == 1 {
		return []template.Road{
			variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildConnectionType(distinctNames[0])).
				WithTo(variant_content.NewRefBuilder().BuildConnectionType(distinctNames[0])).
				Build()}
	}
	var roads []template.Road
	for _, name := range distinctNames[1:] {
		roads = append(roads,
			variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildConnectionType(distinctNames[0])).
				WithTo(variant_content.NewRefBuilder().BuildConnectionType(name)).
				Build())
	}
	return roads
}

func buildSideContentLimits() template.StringList {
	var limits []string
	for a := 1; a <= 5; a++ {
		for b := a + 1; b <= 6; b++ {
			limits = append(limits, fmt.Sprintf("content_limits_side_%d_%d", a, b))
		}
	}
	return limits
}

func createZoneName(label string, playerLabels []string) string {
	if slices.Contains(playerLabels, label) {
		return "Spawn-" + label
	}
	return "Neutral-" + label
}

func buildNonAdjacentDerangement(count int) []int {
	dest := make([]int, count)
	for range 100 {
		candidates := make([]int, count)
		for i := range candidates {
			candidates[i] = i
		}
		rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
		valid := true
		used := make([]bool, count)
		for i := range count {
			found := -1
			for j := range candidates {
				if used[candidates[j]] {
					continue
				}
				candidate := candidates[j]
				if candidate != i && candidate != (i+1)%count && candidate != (i-1+count)%count {
					found = j
					break
				}
			}
			if found < 0 {
				for j := range candidates {
					if !used[candidates[j]] && candidates[j] != i {
						found = j
						break
					}
				}
			}
			if found < 0 {
				valid = false
				break
			}
			dest[i] = candidates[found]
			used[candidates[found]] = true
		}
		if valid {
			return dest
		}
	}

	shift := max(1, count/2)
	for i := range count {
		dest[i] = (i + shift) % count
	}
	return dest
}
