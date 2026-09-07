package preview

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

// Zone is one zone laid out on the preview canvas.
type Zone struct {
	Name    string
	Label   string
	Center  data.Vec2[float64]
	Type    ZoneType
	Quality neutral_zone.Quality
	Owner   int
	Castles int
	Arena   bool
}

func (this Zone) HasCastles() bool {
	return this.Castles > 0
}

func (this Zone) HasArena() bool {
	return this.Arena
}

type ZoneType uint8

const (
	ZoneTypeUnknown ZoneType = iota
	ZoneTypePlayer
	ZoneTypeNeutral
	ZoneTypeHub
)
