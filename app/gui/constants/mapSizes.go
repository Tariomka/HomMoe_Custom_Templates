package constants

import (
	internal_constants "github.com/Tariomka/hommoe_custom_templates/internal/constants"
)

// MapSize and the map-size lists moved to internal/constants so that
// non-GUI code (e.g. internal/validators) can use them; these forwarders
// keep the GUI-facing API unchanged.
type MapSize = internal_constants.MapSize

var (
	BaseMapSizes     = internal_constants.BaseMapSizes
	ExpandedMapSizes = internal_constants.ExpandedMapSizes
	AllMapSizes      = internal_constants.AllMapSizes
)

func GetMapSize(size int) MapSize {
	return internal_constants.GetMapSize(size)
}

func GetMapSizes(withExperimental bool) []MapSize {
	return internal_constants.GetMapSizes(withExperimental)
}
