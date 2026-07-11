package registry

type gatePlacements struct {
	Center string
}

// GetGatePlacementValues returns the gate placement values used for
//
//	variants.connections.gatePlacement
func GetGatePlacementValues() gatePlacements {
	return gatePlacements{
		Center: "Center",
	}
}
