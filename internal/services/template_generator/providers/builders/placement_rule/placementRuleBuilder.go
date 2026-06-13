package placement_rule

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

var ruleTypes = registry.GetRuleTypeValues()

type PlacementRuleBuilder struct {
	item entities.PlacementRule
}

func NewPlacementRuleBuilder() *PlacementRuleBuilder {
	return &PlacementRuleBuilder{}
}

func (this *PlacementRuleBuilder) WithTypeRoad() *PlacementRuleBuilder {
	return this.withType(ruleTypes.Road)
}
func (this *PlacementRuleBuilder) WithTypeCrossroads() *PlacementRuleBuilder {
	return this.withType(ruleTypes.Crossroads)
}
func (this *PlacementRuleBuilder) WithTypeMainObject() *PlacementRuleBuilder {
	return this.withType(ruleTypes.MainObject)
}
func (this *PlacementRuleBuilder) WithDistance(distance Distance) *PlacementRuleBuilder {
	this.item.TargetMin = distance.Min
	this.item.TargetMax = distance.Max
	return this
}
func (this *PlacementRuleBuilder) WithWeight(weight int) *PlacementRuleBuilder {
	this.item.Weight = weight
	return this
}
func (this *PlacementRuleBuilder) WithArgs(arguments ...any) *PlacementRuleBuilder {
	this.item.Args = append(this.item.Args, arguments...)
	return this
}
func (this *PlacementRuleBuilder) Build() entities.PlacementRule { return this.item }

func (this *PlacementRuleBuilder) BuildRoadRule(distance Distance, weight int) entities.PlacementRule {
	return this.
		WithTypeRoad().
		WithDistance(distance).
		WithWeight(weight).
		Build()
}

func (this *PlacementRuleBuilder) BuildCrossroadsRule(distance Distance, weight int) entities.PlacementRule {
	return this.
		WithTypeCrossroads().
		WithDistance(distance).
		WithWeight(weight).
		Build()
}

func (this *PlacementRuleBuilder) BuildNearCastleRule(weight int) entities.PlacementRule {
	return this.
		WithTypeMainObject().
		WithArgs("0").
		WithDistance(Distance{Min: 0.1, Max: 0.3}).
		WithWeight(weight).
		Build()
}

func (this *PlacementRuleBuilder) withType(ruleType string) *PlacementRuleBuilder {
	this.item.Type = ruleType
	return this
}
