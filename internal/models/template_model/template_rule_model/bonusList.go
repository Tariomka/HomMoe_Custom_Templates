package template_rule_model

import "github.com/Tariomka/hommoe_custom_templates/internal/helpers"

type BonusList []Bonus

func (this BonusList) Clone() BonusList {
	return helpers.MapSlice(this, Bonus.Clone)
}
