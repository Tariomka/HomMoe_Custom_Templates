package topology

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

type topologyBase struct {
}

func (this *topologyBase) someMethod() {
	// Method implementation
}

func (this *topologyBase) getSpawnZone(
	letter, player string,
	ringConns []string, castleCount int, matchFactions bool, zoneSize float64,
	spawnFootholds,
	generateRoads bool, tuning models.GenerationTuning) template.Zone {
	mainObjects := []template.MainObject{
		{
			Type: "Spawn", Spawn: player, RemoveGuardIfHasOwner: true,
			GuardChance:              1,
			GuardValue:               tuning.ScaleByNeutralGuardStrength(5000),
			GuardWeeklyIncrement:     0.10,
			BuildingsConstructionSid: "default_buildings_construction",
			Placement:                "Uniform", PlacementArgs: []string{"true", "0.7", "0"},
		},
	}
	for i := 1; i < castleCount; i++ {
		// TODO: add player own castles
		mo := template.MainObject{
			Type: "City", GuardChance: 1, GuardValue: tuning.ScaleByNeutralGuardStrength(2500),
			GuardWeeklyIncrement:     0.10,
			BuildingsConstructionSid: "poor_buildings_construction",
			Placement:                "Uniform", PlacementArgs: []string{"false", "-0.8", "3"},
		}
		if matchFactions {
			mo.Faction = &template.TypedRef{Type: "Match", Args: []string{"0"}}
		} else {
			mo.Faction = &template.TypedRef{Type: "Random", Args: []string{}}
		}
		mainObjects = append(mainObjects, mo)
	}
	cp0 := 0
	biome := template.TypedRef{Type: "MatchMainObject", Args: []string{"0"}}
	var roads []template.Road
	if castleCount > 0 {
		roads = buildOuterZoneRoads(ringConns, castleCount, spawnFootholds, generateRoads)
	} else {
		roads = buildConnectorZoneRoads(ringConns, generateRoads)
	}
	return template.Zone{
		Name: "Spawn-" + letter, Size: normalizeZoneSize(zoneSize), Layout: spawnLayoutName,
		GuardCutoffValue: 2000, GuardRandomization: tuning.GuardRandomization,
		GuardMultiplier: scaleGuardMultiplier(1.0, tuning), GuardWeeklyIncrement: 0.20,
		GuardReactionDistribution: []int{60, 20, 10, 10, 2, 0}, DiplomacyModifier: -0.5,
		GuardedContentPool: cp(t2Guarded), UnguardedContentPool: cp(t2Unguarded), ResourcesContentPool: cp(generalResourcesPoor),
		MandatoryContent:             template.StringList{"mandatory_content_side_" + letter},
		ContentCountLimits:           buildSideContentLimits(),
		GuardedContentValue:          scaleStructureValue(200000*tuning.ContentScale, tuning),
		GuardedContentValuePerArea:   scaleStructureValue(2000*math.Sqrt(tuning.ContentScale), tuning),
		UnguardedContentValue:        scaleStructureValue(50000*tuning.ContentScale, tuning),
		UnguardedContentValuePerArea: scaleStructureValue(400*math.Sqrt(tuning.ContentScale), tuning),
		ResourcesValue:               scaleResourceValue(80000*tuning.ContentScale, tuning),
		ResourcesValuePerArea:        scaleResourceValue(600*math.Sqrt(tuning.ContentScale), tuning),
		MainObjects:                  mainObjects, ZoneBiome: biome, ContentBiome: biome, MetaObjectsBiome: biome,
		CrossroadsPosition: &cp0, Roads: roads,
	}
}

func scaleValue(value, multiplier float64) int {
	return max(0, int(value*multiplier))
}
