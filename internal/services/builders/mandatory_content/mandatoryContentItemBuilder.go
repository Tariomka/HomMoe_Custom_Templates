package mandatory_content

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

type MandatoryContentItemBuilder struct {
	item entities.MandatoryContentItem
}

func NewContentItemBuilder(sid string) *MandatoryContentItemBuilder {
	return &MandatoryContentItemBuilder{item: entities.MandatoryContentItem{SID: sid}}
}

func (this *MandatoryContentItemBuilder) WithName(name string) *MandatoryContentItemBuilder {
	this.item.Name = name
	return this
}
func (this *MandatoryContentItemBuilder) WithGuarded() *MandatoryContentItemBuilder {
	this.item.IsGuarded = true
	return this
}
func (this *MandatoryContentItemBuilder) WithMine() *MandatoryContentItemBuilder {
	this.item.IsMine = true
	return this
}
func (this *MandatoryContentItemBuilder) WithSoloEncounter() *MandatoryContentItemBuilder {
	this.item.SoloEncounter = true
	return this
}
func (this *MandatoryContentItemBuilder) WithRules(rules ...entities.PlacementRule) *MandatoryContentItemBuilder {
	this.item.Rules = append(this.item.Rules, rules...)
	return this
}
func (this *MandatoryContentItemBuilder) WithRulesCallback(
	callback func() []entities.PlacementRule) *MandatoryContentItemBuilder {
	return this.WithRules(callback()...)
}
func (this *MandatoryContentItemBuilder) Build() entities.MandatoryContentItem { return this.item }
