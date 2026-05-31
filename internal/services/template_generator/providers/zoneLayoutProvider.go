package providers

import "github.com/Tariomka/hommoe_custom_templates/internal/models/template"

type ZoneLayoutProvider struct{}

func NewZoneLayoutProvider() *ZoneLayoutProvider {
	return &ZoneLayoutProvider{}
}

func (this *ZoneLayoutProvider) CreateZoneLayouts() []template.ZoneLayoutDef {
	return []template.ZoneLayoutDef{
		this.createZoneLayout("zone_layout_spawns", 0.24, 0.48, 0.30, 16, 0.16, 160, -0.30, 0.4, []int{20, 2, 1}),
		this.createZoneLayout("zone_layout_sides", 0.36, 0.50, 0.25, 16, 0.128, 128, -0.30, 0.3, []int{20, 2, 1}),
		this.createZoneLayout("zone_layout_treasure_zone", 0.50, 0.50, 0.45, 12, 0.12, 96, -0.30, 0.3, []int{12, 3, 1}),
		this.createZoneLayout("zone_layout_center", 0.56, 0.60, 0.30, 10, 0.128, 96, -0.25, 0.3, []int{12, 4, 1}),
	}
}

func (this *ZoneLayoutProvider) createZoneLayout(
	zoneName string,
	obsFill, obsFillVoid, lakesFill float64,
	minLake int,
	elevScale float64,
	roadCluster int,
	roadAttraction, ambientNoise float64,
	groupWeights []int) template.ZoneLayoutDef {
	return template.ZoneLayoutDef{
		Name:                  zoneName,
		ObstaclesFill:         obsFill,
		ObstaclesFillVoid:     obsFillVoid,
		LakesFill:             lakesFill,
		MinLakeArea:           minLake,
		ElevationClusterScale: elevScale,
		ElevationModes: []template.ElevationMode{
			{Weight: 2, MinElevatedFraction: 0.2, MaxElevatedFraction: 0.4},
			{Weight: 1, MinElevatedFraction: 0.6, MaxElevatedFraction: 0.8},
		},
		RoadClusterArea: roadCluster,
		GuardedEncounterResourceFractions: template.GuardedEncounterResourceFractions{
			CountBounds: []int{},
			Fractions:   []float64{0.66},
		},
		AmbientPickupDistribution: template.AmbientPickupDistribution{
			Repulsion: 1.0, Noise: ambientNoise, RoadAttraction: roadAttraction,
			ObstacleAttraction: 0, GroupSizeWeights: groupWeights,
		},
	}
}
