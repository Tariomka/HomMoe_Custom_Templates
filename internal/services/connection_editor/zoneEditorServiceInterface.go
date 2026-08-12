package connection_editor

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

// IZoneEditorService is the contract of the zone-level behaviour of the manual
// zone editor.
type IZoneEditorService interface {
	// EnsureConnectionNames assigns a unique name to every nameless connection,
	// in place.
	EnsureConnectionNames(connections []entities.Connection)

	// RebuildZoneConnectionRoads recomputes each zone's connection and castle
	// roads to match the current connection list and main objects.
	RebuildZoneConnectionRoads(zones []entities.Zone, connections []entities.Connection)

	// RebuildCastleRoads regenerates only the zone's castle<->castle roads,
	// preserving every other road.
	RebuildCastleRoads(zone *entities.Zone)

	// NextFreeZoneLabel returns the first generator label not used by any zone,
	// or "" when the pool is exhausted.
	NextFreeZoneLabel(zones []entities.Zone) string

	// NewDefaultNeutralZone builds a manually-added neutral zone with the same
	// builder the generator uses.
	NewDefaultNeutralZone(
		label string,
		quality neutral_zone.Quality,
		castleCount int,
		generateRoads bool,
		tuning models.GenerationTuning) entities.Zone

	// CountZoneCastles returns the number of City main objects in the zone.
	CountZoneCastles(zone entities.Zone) int

	// ApplyNeutralZoneQuality re-applies the quality profile and rebuilds the
	// zone's castles for the requested count.
	ApplyNeutralZoneQuality(
		zone *entities.Zone,
		quality neutral_zone.Quality,
		castleCount int,
		tuning models.GenerationTuning)

	// CanDeleteZone reports whether the zone may be removed in the editor.
	CanDeleteZone(zoneName string, playerZoneNames map[string]bool) bool

	// RemoveZone returns the zone and connection lists without the named zone
	// and without any connection referencing it.
	RemoveZone(
		zones []entities.Zone,
		connections []entities.Connection,
		zoneName string) ([]entities.Zone, []entities.Connection)

	// FindOpenPosition returns a normalized position that maximizes the
	// distance to the occupied positions.
	FindOpenPosition(occupied [][2]float64) [2]float64
}
