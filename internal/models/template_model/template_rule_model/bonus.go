package template_rule_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
)

type Bonus struct{ template_entity.Bonus }

func (this Bonus) Clone() Bonus {
	clone := this
	clone.Parameters = slices.Clone(this.Parameters)
	return clone
}
