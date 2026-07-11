package registry

type placements struct {
	Center     string
	Connection string
	NearZone   string
	Uniform    string
}

// GetPlacementValues returns the placement values used for
//
//	variants.zones.mainObjects.placement
func GetPlacementValues() placements {
	return placements{
		Center:     "Center",
		Connection: "Connection",
		NearZone:   "NearZone",
		Uniform:    "Uniform",
	}
}
