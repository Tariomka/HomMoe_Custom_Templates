package template_content_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_common_model"
)

type MandatoryContentItem struct {
	SID       string
	Name      string
	IsMine    bool
	IsGuarded bool
	Rules     []template_common_model.PlacementRule

	Variant *int
	Owner   string
	Road    *bool

	GuardValue int

	IncludeLists []string
	Content      []WeightedContent

	DesignatedEncounter *bool
	SoloEncounter       bool
}

func (this MandatoryContentItem) Clone() MandatoryContentItem {
	clone := this
	clone.Rules = helpers.MapSlice(this.Rules, template_common_model.PlacementRule.Clone)
	clone.Variant = helpers.ClonePointer(this.Variant)
	clone.IncludeLists = slices.Clone(this.IncludeLists)
	clone.Content = slices.Clone(this.Content)
	clone.DesignatedEncounter = helpers.ClonePointer(this.DesignatedEncounter)
	clone.Road = helpers.ClonePointer(this.Road)
	return clone
}
