package common_connections

import "github.com/Tariomka/hommoe_custom_templates/internal/registry"

func GetConnectionTypes() []string {
	connectionTypes := registry.GetConnectionTypeValues()
	return []string{connectionTypes.Direct, connectionTypes.Portal}
}
