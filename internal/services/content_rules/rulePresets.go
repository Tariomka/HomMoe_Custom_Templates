package content_rules

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/placement_rule"
)

// Commonly used placement rulesets, ported from the C# RulePresets helper.
// These reuse the existing placement_rule builder so the generated
// PlacementRule output is byte-for-byte identical to the rest of the generator.

// RoadDistance builds a "Road" placement rule for the given distance band.
func RoadDistance(distance DistanceVariation, weight int) template.PlacementRule {
	return placement_rule.NewPlacementRuleBuilder().
		BuildRoadRule(placement_rule.Distance{Min: distance.Min, Max: distance.Max}, weight)
}

// TownDistance builds a "MainObject" (nearest town) placement rule.
func TownDistance(distance DistanceVariation, weight int) template.PlacementRule {
	return placement_rule.NewPlacementRuleBuilder().
		WithTypeMainObject().
		WithArgs("0").
		WithDistance(placement_rule.Distance{Min: distance.Min, Max: distance.Max}).
		WithWeight(weight).
		Build()
}

// CrossroadsDistance builds a "Crossroads" placement rule.
func CrossroadsDistance(distance DistanceVariation, weight int) template.PlacementRule {
	return placement_rule.NewPlacementRuleBuilder().
		BuildCrossroadsRule(placement_rule.Distance{Min: distance.Min, Max: distance.Max}, weight)
}
