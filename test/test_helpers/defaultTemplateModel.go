package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// GetDefaultTemplateModel is the golden template lifted to the service layer,
// for the tests that work with what the generator and drivers.State actually
// hold. It goes through the real mapper, so a converter that drops a field
// fails those tests too.
func GetDefaultTemplateModel() template_model.Template {
	return mappers.NewTemplateMapper().ToModel(GetDefaultTemplate())
}
