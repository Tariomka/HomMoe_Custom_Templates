package template_rule_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

type BonusList []Bonus

func (this BonusList) Clone() BonusList {
	return helpers.MapSlice(this, Bonus.Clone)
}

func ToBonusListModel(entity template.BonusList) BonusList {
	return helpers.MapSlice(entity, ToBonusModel)
}

func ToBonusListEntity(model BonusList) template.BonusList {
	return helpers.MapSlice(model, ToBonusEntity)
}
