package constants

import "strings"

// IsHubLabel reports whether a zone label denotes a hub - either the single
// shared hub ("Hub") or a tournament per-player hub ("Hub-A").
func IsHubLabel(label string) bool {
	return label == HubZoneName || strings.HasPrefix(label, HubZonePrefix)
}
