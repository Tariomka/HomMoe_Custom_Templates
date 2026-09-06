package template_variant_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
)

type StringList []string

func ToStringListModel(entity template_entity.StringList) StringList {
	return StringList(entity)
}

func ToStringListEntity(model StringList) template_entity.StringList {
	return template_entity.StringList(model)
}
