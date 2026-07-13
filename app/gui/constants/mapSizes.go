package constants

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/common"
)

// MapSize and the map-size lists moved to internal/constants so that
// non-GUI code (e.g. internal/validators) can use them; these forwarders
// keep the GUI-facing API unchanged.
type MapSize = common.MapSize

var (
	BaseMapSizes     = common.BaseMapSizes
	ExpandedMapSizes = common.ExpandedMapSizes
	AllMapSizes      = common.AllMapSizes
)

func GetMapSize(size int) MapSize {
	return common.GetMapSize(size)
}

func GetMapSizes(withExperimental bool) []MapSize {
	return common.GetMapSizes(withExperimental)
}
