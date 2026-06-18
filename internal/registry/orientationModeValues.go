package registry

type orientationModes struct {
	BoundingCircle        string
	MinimalBoundingSquare string
}

var orientationModeValues = orientationModes{
	BoundingCircle:        "BoundingCircle",
	MinimalBoundingSquare: "MinimalBoundingSquare",
}

// GetOrientationModeValues returns the available orientation mode values used for
//
//	variants.orientation.mode
func GetOrientationModeValues() orientationModes {
	return orientationModeValues
}
