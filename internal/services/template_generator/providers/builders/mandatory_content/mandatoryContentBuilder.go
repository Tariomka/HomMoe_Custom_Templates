package mandatory_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/placement_rule"
)

type MandatoryContentBuilder struct {
	item template.MandatoryContentItem
}

func NewContentBuilder(sid string) *MandatoryContentBuilder {
	return &MandatoryContentBuilder{
		item: template.MandatoryContentItem{SID: sid},
	}
}

func (this *MandatoryContentBuilder) WithName(name string) *MandatoryContentBuilder {
	this.item.Name = name
	return this
}
func (this *MandatoryContentBuilder) WithGuarded() *MandatoryContentBuilder {
	this.item.IsGuarded = true
	return this
}
func (this *MandatoryContentBuilder) WithMine() *MandatoryContentBuilder {
	this.item.IsMine = true
	return this
}
func (this *MandatoryContentBuilder) WithSoloEncounter() *MandatoryContentBuilder {
	this.item.SoloEncounter = true
	return this
}
func (this *MandatoryContentBuilder) WithRules(rules ...template.PlacementRule) *MandatoryContentBuilder {
	this.item.Rules = append(this.item.Rules, rules...)
	return this
}
func (this *MandatoryContentBuilder) WithRulesCallback(
	callback func() []template.PlacementRule) *MandatoryContentBuilder {
	return this.WithRules(callback()...)
}
func (this *MandatoryContentBuilder) Build() template.MandatoryContentItem { return this.item }

func (this *MandatoryContentBuilder) WithRoadDistance(distance placement_rule.Distance) *MandatoryContentBuilder { // TODO: probably not needed
	return this.WithRules(placement_rule.NewPlacementRuleBuilder().BuildRoadRule(distance, 1))
}
