package placement_rule

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_distances"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

type PlacementRuleBuilder struct {
	item entities.PlacementRule
}

func NewPlacementRuleBuilder() *PlacementRuleBuilder {
	return &PlacementRuleBuilder{}
}

func (this *PlacementRuleBuilder) WithTypeRoad() *PlacementRuleBuilder {
	return this.withType(registry.GetRuleTypeValues().Road)
}
func (this *PlacementRuleBuilder) WithTypeCrossroads() *PlacementRuleBuilder {
	return this.withType(registry.GetRuleTypeValues().Crossroads)
}
func (this *PlacementRuleBuilder) WithTypeMainObject() *PlacementRuleBuilder {
	return this.withType(registry.GetRuleTypeValues().MainObject)
}

func (this *PlacementRuleBuilder) WithDistance(distance models.DistancePreset) *PlacementRuleBuilder {
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

func (this *PlacementRuleBuilder) BuildRoadRule(
	distance models.DistancePreset,
	weight int,
) entities.PlacementRule {
	return this.
		WithTypeRoad().
		WithDistance(distance).
		WithWeight(weight).
		Build()
}

func (this *PlacementRuleBuilder) BuildCrossroadsRule(
	distance models.DistancePreset,
	weight int,
) entities.PlacementRule {
	return this.
		WithTypeCrossroads().
		WithDistance(distance).
		WithWeight(weight).
		Build()
}

func (this *PlacementRuleBuilder) BuildCastleRule(
	distance models.DistancePreset,
	weight int,
) entities.PlacementRule {
	return this.
		WithTypeMainObject().
		WithArgs("0").
		WithDistance(distance).
		WithWeight(weight).
		Build()
}

func (this *PlacementRuleBuilder) BuildNearCastleRule(weight int) entities.PlacementRule {
	return this.
		WithTypeMainObject().
		WithArgs("0").
		WithDistance(getPortalPlacementNearDistance()).
		WithWeight(weight).
		Build()
}

func (this *PlacementRuleBuilder) BuildNearCrossroadsRule(weight int) entities.PlacementRule {
	return this.
		WithTypeCrossroads().
		WithDistance(getPortalPlacementNearDistance()).
		WithWeight(weight).
		Build()
}

func getPortalPlacementNearDistance() models.DistancePreset {
	distance, _ := common_distances.GetPortalPlacementDistancePreset("Near")
	return distance
}

func (this *PlacementRuleBuilder) withType(ruleType string) *PlacementRuleBuilder {
	this.item.Type = ruleType
	return this
}
