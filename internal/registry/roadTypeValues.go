package registry

type roadTypes struct {
	Dirt  string
	Stone string
}

// GetRoadTypeValues returns the road material types used for
//
//	variants.zones.roads.type
func GetRoadTypeValues() roadTypes {
	return roadTypes{
		Dirt:  "Dirt",
		Stone: "Stone",
	}
}
