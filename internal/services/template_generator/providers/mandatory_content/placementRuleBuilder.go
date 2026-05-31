package mandatory_content

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
func (this *PlacementRuleBuilder) WithArg(argument any) *PlacementRuleBuilder {
	this.item.Args = append(this.item.Args, argument)
	return this
}
func (this *PlacementRuleBuilder) WithArgs(arguments []any) *PlacementRuleBuilder {
	this.item.Args = append(this.item.Args, arguments...)
	return this
}
func (this *PlacementRuleBuilder) Build() template.PlacementRule { return this.item }
