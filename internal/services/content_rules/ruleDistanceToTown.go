package content_rules

import (
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
)

// Rule metadata constants for the distance-to-town rule.
const (
	RuleDistanceToTownName        = "Distance to town"
	RuleDistanceToTownDescription = "Distance to the nearest town from the content item."
	RuleDistanceToTownMarker      = "T"
)

// RuleDistanceToTown constrains how far the content spawns from the nearest town.
type RuleDistanceToTown struct {
	Distance models.DistancePreset
}

// NewRuleDistanceToTown creates a town-distance rule, defaulting to Medium when
// no distance is supplied.
func NewRuleDistanceToTown(distance *models.DistancePreset) *RuleDistanceToTown {
	resolved := defaultDistancePreset()
	if distance != nil {
		resolved = *distance
	}
	return &RuleDistanceToTown{Distance: resolved}
}

func (this *RuleDistanceToTown) Name() string { return RuleDistanceToTownName }

func (this *RuleDistanceToTown) Description() string { return RuleDistanceToTownDescription }

func (this *RuleDistanceToTown) Marker() string { return RuleDistanceToTownMarker }

func (this *RuleDistanceToTown) DisplayText() string {
	return fmt.Sprintf("%s: %s", this.Name(), this.Distance.Name)
}

func (this *RuleDistanceToTown) Apply(item *template_model.MandatoryContentItem) {
	item.Rules = append(item.Rules, placement_rule.NewPlacementRuleBuilder().BuildCastleRule(this.Distance, 1))
}

func (this *RuleDistanceToTown) SerializeToRowSave() editor_state_model.ContentRuleRow {
	return editor_state_model.ContentRuleRow{
		Name:         this.Name(),
		DistanceName: this.Distance.Name,
	}
}
