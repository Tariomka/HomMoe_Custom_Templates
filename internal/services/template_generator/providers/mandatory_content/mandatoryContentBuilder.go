package mandatory_content

import "github.com/Tariomka/hommoe_custom_templates/internal/models/template"

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
func (this *MandatoryContentBuilder) WithRoadDistance(distance Distance) *MandatoryContentBuilder {
	return this.WithRule(ruleRoadDistance(distance, 1))
}
func (this *MandatoryContentBuilder) WithRule(rule template.PlacementRule) *MandatoryContentBuilder {
	this.item.Rules = append(this.item.Rules, rule)
	return this
}
func (this *MandatoryContentBuilder) WithRules(rules []template.PlacementRule) *MandatoryContentBuilder {
	this.item.Rules = append(this.item.Rules, rules...)
	return this
}
func (this *MandatoryContentBuilder) Build() template.MandatoryContentItem { return this.item }
