package content_rules

import (
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// Rule metadata constants for the guarded rule.
const (
	RuleGuardedName        = "Guarded"
	RuleGuardedDescription = "Forces the content item to be guarded or unguarded, regardless of the default behavior."
	RuleGuardedMarker      = "G"
)

// RuleGuarded forces the content item to be guarded or unguarded.
type RuleGuarded struct {
	IsGuarded bool
}

// NewRuleGuarded creates a guarded rule with the supplied state.
func NewRuleGuarded(isGuarded bool) *RuleGuarded {
	return &RuleGuarded{IsGuarded: isGuarded}
}

func (this *RuleGuarded) Name() string { return RuleGuardedName }

func (this *RuleGuarded) Description() string { return RuleGuardedDescription }

// Marker shows "G" when guarded and "!G" when explicitly unguarded.
func (this *RuleGuarded) Marker() string {
	if this.IsGuarded {
		return RuleGuardedMarker
	}

	return "!" + RuleGuardedMarker
}

func (this *RuleGuarded) DisplayText() string {
	return fmt.Sprintf("%s: %t", this.Name(), this.IsGuarded)
}

func (this *RuleGuarded) Apply(item *entities.MandatoryContentItem) { item.IsGuarded = this.IsGuarded }

func (this *RuleGuarded) SerializeToRowSave() editor_state_model.ContentRuleRow {
	value := this.IsGuarded
	return editor_state_model.ContentRuleRow{
		Name:      this.Name(),
		IsGuarded: &value,
	}
}
