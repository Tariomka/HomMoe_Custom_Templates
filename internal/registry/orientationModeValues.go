package registry

type orientationModes struct {
	BoundingCircle        string
	MinimalBoundingSquare string
}

// GetOrientationModeValues returns the available orientation mode values used for
//
//	variants.orientation.mode
func GetOrientationModeValues() orientationModes {
	return orientationModes{
		BoundingCircle:        "BoundingCircle",
		MinimalBoundingSquare: "MinimalBoundingSquare",
	}
}
