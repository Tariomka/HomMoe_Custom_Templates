package content_rules

import (
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
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
	Distance DistanceVariation
}

// NewRuleDistanceToRoad creates a road-distance rule, defaulting to Medium when
// no distance is supplied.
func NewRuleDistanceToRoad(distance *DistanceVariation) *RuleDistanceToRoad {
	resolved := defaultDistanceVariation()
	if distance != nil {
		resolved = *distance
	}
	return &RuleDistanceToRoad{Distance: resolved}
}

func (this *RuleDistanceToRoad) Name() string        { return RuleDistanceToRoadName }
func (this *RuleDistanceToRoad) Description() string { return RuleDistanceToRoadDescription }
func (this *RuleDistanceToRoad) Marker() string      { return RuleDistanceToRoadMarker }

func (this *RuleDistanceToRoad) DisplayText() string {
	return fmt.Sprintf("%s: %s", this.Name(), this.Distance.Name)
}

func (this *RuleDistanceToRoad) Apply(item *entities.MandatoryContentItem) {
	item.Rules = append(item.Rules, placement_rule.NewPlacementRuleBuilder().
		BuildRoadRule(placement_rule.Distance{Min: this.Distance.Min, Max: this.Distance.Max}, 1))
}

func (this *RuleDistanceToRoad) SerializeToRowSave() models.ContentRuleRowSave {
	return models.ContentRuleRowSave{
		Name:         this.Name(),
		DistanceName: this.Distance.Name,
	}
}
