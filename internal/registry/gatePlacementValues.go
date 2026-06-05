package registry

type gatePlacements struct {
	Center string
}

var gatePlacementValues = gatePlacements{
	Center: "Center",
}

// GetGatePlacementValues returns the gate placement values used for
//
//	variants.connections.gatePlacement
func GetGatePlacementValues() gatePlacements {
	return gatePlacementValues
}
