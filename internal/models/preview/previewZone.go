package preview

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
)

// Zone is one zone laid out on the preview canvas.
type Zone struct {
	Name    string
	Label   string
	Center  image.Point
	Type    ZoneType
	Quality neutralZone.Quality
	Owner   int
	Castles int
}

func (this Zone) HasCastles() bool {
	return this.Castles > 0
}

type ZoneType uint8

const (
	ZoneTypeUnknown ZoneType = iota
	ZoneTypePlayer
	ZoneTypeNeutralZone
	ZoneTypeHub
)
