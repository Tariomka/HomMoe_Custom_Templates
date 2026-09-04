package connection_editor

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// IZoneEditorService is the contract of the zone-level behaviour of the manual
// zone editor.
type IZoneEditorService interface {
	// EnsureConnectionNames assigns a unique name to every nameless connection,
	// in place.
	EnsureConnectionNames(connections []template_model.Connection)

	// RebuildZoneConnectionRoads recomputes each zone's connection and castle
	// roads to match the current connection list and main objects.
	RebuildZoneConnectionRoads(zones []template_model.Zone, connections []template_model.Connection)

	// RebuildCastleRoads regenerates only the zone's castle<->castle roads,
	// preserving every other road.
	RebuildCastleRoads(zone *template_model.Zone)

	// NextFreeZoneLabel returns the first generator label not used by any zone,
	// or "" when the pool is exhausted.
	NextFreeZoneLabel(zones []template_model.Zone) string

	// NewDefaultNeutralZone builds a manually-added neutral zone with the same
	// builder the generator uses, recording the requested quality on it.
	NewDefaultNeutralZone(
		label string,
		quality neutral_zone.Quality,
		castleCount int,
		generateRoads bool,
		tuning models.GenerationTuning) template_model.Zone

	// CountZoneCastles returns the number of City main objects in the zone.
	CountZoneCastles(zone template_model.Zone) int

	// ApplyNeutralZoneQuality re-applies the quality profile, records it on the
	// zone and rebuilds the zone's castles for the requested count.
	ApplyNeutralZoneQuality(
		zone *template_model.Zone,
		quality neutral_zone.Quality,
		castleCount int,
		tuning models.GenerationTuning)

	// CanDeleteZone reports whether the zone may be removed in the editor.
	CanDeleteZone(zoneName string, playerZoneNames map[string]bool) bool

	// RemoveZone returns the zone and connection lists without the named zone
	// and without any connection referencing it.
	RemoveZone(
		zones []template_model.Zone,
		connections []template_model.Connection,
		zoneName string) ([]template_model.Zone, []template_model.Connection)

	// FindOpenPosition returns a normalized position that maximizes the
	// distance to the occupied positions.
	FindOpenPosition(occupied [][2]float64) [2]float64
}
