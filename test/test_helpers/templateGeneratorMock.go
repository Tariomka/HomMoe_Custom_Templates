package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/mock"
)

// TemplateGeneratorMock is a testify mock of
// template_generator.ITemplateGenerator, used to unit-test collaborators
// without running a real generation.
type TemplateGeneratorMock struct {
	mock.Mock
}

func (this *TemplateGeneratorMock) SetConfiguration(configuration *config.GeneratorConfig) {
	this.Called(configuration)
}

func (this *TemplateGeneratorMock) Generate() (*template_model.Template, []string) {
	arguments := this.Called()
	generated, _ := arguments.Get(0).(*template_model.Template)
	warnings, _ := arguments.Get(1).([]string)
	return generated, warnings
}
