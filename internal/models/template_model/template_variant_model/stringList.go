package template_variant_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
)

type StringList []string

func ToStringListModel(entity template.StringList) StringList {
	return StringList(entity)
}

func ToStringListEntity(model StringList) template.StringList {
	return template.StringList(model)
}
