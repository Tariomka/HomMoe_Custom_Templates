package content_rules

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// IContentRule is the polymorphic interface implemented by every content rule.
// Rules describe a configurable constraint in the UI, know how to serialize themselves,
// and know how to apply their effect to a final content item.
type IContentRule interface {
	// Name uniquely identifies the rule and matches the persisted ContentRuleRowSave.Name field.
	Name() string
	// Description is the long-form explanation shown in the UI.
	Description() string
	// Marker is the short badge (usually a single letter) shown on zone
	// content rows.
	Marker() string
	// DisplayText is the user-facing single-line representation of the rule.
	DisplayText() string
	// SerializeToRowSave projects the rule back to its persisted form.
	SerializeToRowSave() editor_state_model.ContentRuleRow
	// Apply mutates the final content item according to the rule.
	Apply(item *entities.MandatoryContentItem)
}
