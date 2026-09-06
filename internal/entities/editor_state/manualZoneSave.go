package editor_state

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
)

// ManualZoneSave persists a zone edited in the manual zone editor. The three
// positions travel beside the zone because entities.Zone has nowhere to put
// them: they describe where the editor drew the zone, not anything the game
// reads out of a .rmg.json.
//
// Every position is a POINTER because nil is load-bearing - it means the stamp
// was never taken, which is not the same as a stamp at the origin.
type ManualZoneSave struct {
	Zone entities.Zone `json:"zone"`

	// GeneratorPosition and GeneratorRing are where the generator placed the
	// zone; the preview layout falls back to a computed layout without them.
	GeneratorPosition *data.Vec2[float64] `json:"generatorPosition,omitempty"`
	GeneratorRing     *int                `json:"generatorRing,omitempty"`

	// ManualPosition is where the user dragged the zone, and is the essential
	// piece of a hand-made layout.
	ManualPosition *data.Vec2[float64] `json:"manualPosition,omitempty"`

	// Quality is the raw ordinal of the zone's recorded tier, absent when the
	// tier was never recorded. It is an int8 rather than a neutral_zone.Quality
	// because an entity may not import internal/models, and a POINTER because
	// the enum counts from iota - 1: with omitempty a plain field would drop
	// every Plastic zone (ordinal 0) back to "absent".
	Quality *int8 `json:"quality,omitempty"`
}
