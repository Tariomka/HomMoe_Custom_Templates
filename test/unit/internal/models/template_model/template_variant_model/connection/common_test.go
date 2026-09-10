package connection_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_common_model"
)

// newPortalPlacementRules builds the portal-endpoint placement rules that a
// non-Portal connection may legitimately carry.
func newPortalPlacementRules() []template_common_model.PlacementRule {
	return []template_common_model.PlacementRule{
		{Type: "MainObject", Weight: 1},
	}
}
