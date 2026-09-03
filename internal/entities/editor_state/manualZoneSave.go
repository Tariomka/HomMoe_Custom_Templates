package editor_state

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// ManualZoneSave persists a zone edited in the manual zone editor. The zone's
// ManualPosition is captured separately because entities.Zone deliberately
// omits it from JSON (json:"-"), yet that normalized position is the essential
// piece of a hand-made layout and must survive a save/load round-trip.
type ManualZoneSave struct {
	Zone           entities.Zone `json:"zone"`
	ManualPosition *[2]float64   `json:"manualPosition,omitempty"`

	// Quality is the raw ordinal of the zone's recorded tier, absent when the
	// tier was never recorded. It is an int8 rather than a neutral_zone.Quality
	// because an entity may not import internal/models, and a POINTER because
	// the enum counts from iota - 1: with omitempty a plain field would drop
	// every Plastic zone (ordinal 0) back to "absent".
	Quality *int8 `json:"quality,omitempty"`
}
