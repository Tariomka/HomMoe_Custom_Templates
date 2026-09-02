package template_rule_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
)

type Bonus struct{ template.Bonus }

func ToBonusModel(entity template.Bonus) Bonus {
	return Bonus{Bonus: entity}
}

func ToBonusEntity(model Bonus) template.Bonus {
	return model.Bonus
}
