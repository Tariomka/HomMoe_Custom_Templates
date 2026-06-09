package registry

type placements struct {
	Center     string
	Connection string
	NearZone   string
	Uniform    string
}

var placementValues = placements{
	Center:     "Center",
	Connection: "Connection",
	NearZone:   "NearZone",
	Uniform:    "Uniform",
}

// GetPlacementValues returns the placement values used for
//
//	variants.zones.mainObjects.placement
func GetPlacementValues() placements {
	return placementValues
}
