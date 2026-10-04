package common_topologies

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

// retiredTopologyLabels maps the saved IDs of retired topologies to their former labels.
// It exists only so loading such a state can name what was retired.
var retiredTopologyLabels = map[config.MapTopology]string{ //nolint:gochecknoglobals // Immutable lookup.
	"Default":     "Ring",
	"HubAndSpoke": "Hub",
	"Chain":       "Chain",
	"SharedWeb":   "Shared Web",
}

// GetRetiredTopologyLabel returns the former label of a retired topology ID.
func GetRetiredTopologyLabel(topology config.MapTopology) (string, bool) {
	label, ok := retiredTopologyLabels[topology]
	return label, ok
}
