package placement_rule

import "github.com/Tariomka/hommoe_custom_templates/internal/models/template"

type PlacementRuleBuilder struct {
	item template.PlacementRule
}

func NewPlacementRuleBuilder() *PlacementRuleBuilder {
	return &PlacementRuleBuilder{}
}

func (this *PlacementRuleBuilder) WithRoadType() *PlacementRuleBuilder {
	this.item.Type = "Road"
	return this
}
func (this *PlacementRuleBuilder) WithCrossroadsType() *PlacementRuleBuilder {
	this.item.Type = "Crossroads"
	return this
}
func (this *PlacementRuleBuilder) WithMainObjectType() *PlacementRuleBuilder {
	this.item.Type = "MainObject"
	return this
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
func (this *PlacementRuleBuilder) Build() template.PlacementRule { return this.item }

func (this *PlacementRuleBuilder) BuildRoadRule(distance Distance, weight int) template.PlacementRule {
	return this.
		WithRoadType().
		WithDistance(distance).
		WithWeight(weight).
		Build()
}

func (this *PlacementRuleBuilder) BuildCrossroadsRule(distance Distance, weight int) template.PlacementRule {
	return this.
		WithCrossroadsType().
		WithDistance(distance).
		WithWeight(weight).
		Build()
}

func (this *PlacementRuleBuilder) BuildNearCastleRule(weight int) template.PlacementRule {
	return this.
		WithMainObjectType().
		WithArgs("0").
		WithDistance(Distance{Min: 0.1, Max: 0.3}).
		WithWeight(weight).
		Build()
}
