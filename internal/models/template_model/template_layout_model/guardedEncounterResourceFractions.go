package template_layout_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
)

type GuardedEncounterResourceFractions struct {
	template_entity.GuardedEncounterResourceFractions
}

func (this GuardedEncounterResourceFractions) Clone() GuardedEncounterResourceFractions {
	clone := this
	clone.CountBounds = slices.Clone(this.CountBounds)
	clone.Fractions = slices.Clone(this.Fractions)
	return clone
}
