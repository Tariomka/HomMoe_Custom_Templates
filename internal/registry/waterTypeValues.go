package registry

type waterTypes struct {
	WaterDirt   string
	WaterSand   string
	WaterDeath  string
	WaterSnow   string
	WaterFallen string
	Lava        string
	WaterGrass  string
}

// GetWaterTypeValues returns the different water types used for
//
//	variants.border.waterType
func GetWaterTypeValues() waterTypes {
	return waterTypes{
		WaterDirt:   "water dirt",
		WaterSand:   "water sand",
		WaterDeath:  "water death",
		WaterSnow:   "water snow",
		WaterFallen: "water fallen",
		Lava:        "lava",
		WaterGrass:  "water grass",
	}
}
