package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/mock"
)

// MandatoryContentProviderMock is a testify mock of
// provider_interfaces.IMandatoryContentProvider, used to unit-test
// collaborators without the real content catalogue.
type MandatoryContentProviderMock struct {
	mock.Mock
}

func (this *MandatoryContentProviderMock) CreateContents(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans) []entities.MandatoryContent {
	arguments := this.Called(configuration, playerLabels, neutralZones)
	contents, _ := arguments.Get(0).([]entities.MandatoryContent)
	return contents
}

func (this *MandatoryContentProviderMock) CreateContentsForZones(
	configuration config.GeneratorConfig,
	zones []template_model.Zone) []entities.MandatoryContent {
	arguments := this.Called(configuration, zones)
	contents, _ := arguments.Get(0).([]entities.MandatoryContent)
	return contents
}
