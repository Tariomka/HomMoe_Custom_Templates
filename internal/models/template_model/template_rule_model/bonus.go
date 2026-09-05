package template_rule_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
)

type Bonus struct{ template.Bonus }

func (this Bonus) Clone() Bonus {
	clone := this
	clone.Parameters = slices.Clone(this.Parameters)
	return clone
}

func ToBonusModel(entity template.Bonus) Bonus {
	return Bonus{Bonus: entity}
}

func ToBonusEntity(model Bonus) template.Bonus {
	return model.Bonus
}
