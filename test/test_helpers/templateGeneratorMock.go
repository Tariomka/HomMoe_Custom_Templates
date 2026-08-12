package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
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

func (this *TemplateGeneratorMock) Generate() (*entities.RmgTemplate, []string) {
	arguments := this.Called()
	template, _ := arguments.Get(0).(*entities.RmgTemplate)
	warnings, _ := arguments.Get(1).([]string)
	return template, warnings
}
