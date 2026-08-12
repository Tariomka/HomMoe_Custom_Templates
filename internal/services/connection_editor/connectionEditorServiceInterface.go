package connection_editor

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// IConnectionEditorService is the contract of the connection-level behaviour of
// the Zone Connection Editor.
type IConnectionEditorService interface {
	// NewDefaultConnection builds a user-added direct connection between two
	// zones with guard settings derived from the pair's connection quality.
	NewDefaultConnection(
		from string,
		to string,
		zones []entities.Zone,
		playerZoneNames map[string]bool) entities.Connection

	// FindIsolatedZones returns the names of zones that no connection touches.
	FindIsolatedZones(zones []entities.Zone, connections []entities.Connection) []string

	// ComputeHasErrors reports whether any connection references a zone that
	// does not exist.
	ComputeHasErrors(zones []entities.Zone, connections []entities.Connection) bool

	// HasDuplicateName reports whether another connection already uses the
	// current connection's name.
	HasDuplicateName(connections []entities.Connection, current *entities.Connection) bool
}
