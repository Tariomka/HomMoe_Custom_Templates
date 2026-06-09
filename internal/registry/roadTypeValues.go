package registry

type roadTypes struct {
	Dirt  string
	Stone string
}

var roadTypeValues = roadTypes{
	Dirt:  "Dirt",
	Stone: "Stone",
}

// GetRoadTypeValues returns the road material types used for
//
//	variants.zones.roads.type
func GetRoadTypeValues() roadTypes {
	return roadTypeValues
}
