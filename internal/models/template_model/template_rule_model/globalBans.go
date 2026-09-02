package template_rule_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
)

type GlobalBans struct {
	template.GlobalBans
}

func ToGlobalBansModel(entity template.GlobalBans) GlobalBans {
	return GlobalBans{GlobalBans: entity}
}

func ToGlobalBansEntity(model GlobalBans) template.GlobalBans {
	return model.GlobalBans
}
