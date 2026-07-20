package preview

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

// Zone is one zone laid out on the preview canvas.
type Zone struct {
	Name    string
	Label   string
	Center  image.Point
	Type    ZoneType
	Quality neutral_zone.Quality
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
	ZoneTypeNeutral
	ZoneTypeHub
)
