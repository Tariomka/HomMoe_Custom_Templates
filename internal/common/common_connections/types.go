package common_connections

import "github.com/Tariomka/hommoe_custom_templates/internal/registry"

// GetConnectionTypes lists the connection types a user may add in the
// editor. Proximity is generator-only and intentionally excluded.
func GetConnectionTypes() []string {
	connectionTypes := registry.GetConnectionTypeValues()
	return []string{
		connectionTypes.Direct,
		connectionTypes.Portal,
	}
}
