package registry

type waterTypes struct {
	WaterGrass string
}

var waterTypeValues = waterTypes{
	WaterGrass: "water grass",
}

// GetWaterTypeValues returns the different water types used for
//
//	variants.border.waterType
func GetWaterTypeValues() waterTypes {
	return waterTypeValues
}
