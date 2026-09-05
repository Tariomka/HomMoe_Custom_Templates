package content_rules

import (
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
)

// Rule metadata constants for the distance-to-road rule.
const (
	RuleDistanceToRoadName        = "Distance to road"
	RuleDistanceToRoadDescription = "Distance to the nearest road from the content item."
	RuleDistanceToRoadMarker      = "R"
)

// RuleDistanceToRoad constrains how far the content spawns from the nearest road.
type RuleDistanceToRoad struct {
	Distance models.DistancePreset
}

// NewRuleDistanceToRoad creates a road-distance rule, defaulting to Medium when
// no distance is supplied.
func NewRuleDistanceToRoad(distance *models.DistancePreset) *RuleDistanceToRoad {
	resolved := defaultDistancePreset()
	if distance != nil {
		resolved = *distance
	}
	return &RuleDistanceToRoad{Distance: resolved}
}

func (this *RuleDistanceToRoad) Name() string { return RuleDistanceToRoadName }

func (this *RuleDistanceToRoad) Description() string { return RuleDistanceToRoadDescription }

func (this *RuleDistanceToRoad) Marker() string { return RuleDistanceToRoadMarker }

func (this *RuleDistanceToRoad) DisplayText() string {
	return fmt.Sprintf("%s: %s", this.Name(), this.Distance.Name)
}

func (this *RuleDistanceToRoad) Apply(item *template_model.MandatoryContentItem) {
	item.Rules = append(item.Rules, placement_rule.NewPlacementRuleBuilder().BuildRoadRule(this.Distance, 1))
}

func (this *RuleDistanceToRoad) SerializeToRowSave() editor_state_model.ContentRuleRow {
	return editor_state_model.ContentRuleRow{
		Name:         this.Name(),
		DistanceName: this.Distance.Name,
	}
}
