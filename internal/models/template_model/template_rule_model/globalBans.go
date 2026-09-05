package template_rule_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
)

type GlobalBans struct {
	template.GlobalBans
}

func (this GlobalBans) Clone() GlobalBans {
	clone := this
	clone.Items = slices.Clone(this.Items)
	clone.Magics = slices.Clone(this.Magics)
	clone.Heroes = slices.Clone(this.Heroes)
	return clone
}

func ToGlobalBansModel(entity template.GlobalBans) GlobalBans {
	return GlobalBans{GlobalBans: entity}
}

func ToGlobalBansEntity(model GlobalBans) template.GlobalBans {
	return model.GlobalBans
}
