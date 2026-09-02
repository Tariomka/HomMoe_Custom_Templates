package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/mock"
)

// ZoneTierServiceMock is a testify mock of zone_interfaces.IZoneTierService,
// used to unit-test collaborators without the real pool-based detection.
type ZoneTierServiceMock struct {
	mock.Mock
}

func (this *ZoneTierServiceMock) GetQuality(zone entities.Zone) neutral_zone.Quality {
	arguments := this.Called(zone)
	quality, _ := arguments.Get(0).(neutral_zone.Quality)
	return quality
}

func (this *ZoneTierServiceMock) ResolveQuality(zone template_model.Zone) neutral_zone.Quality {
	arguments := this.Called(zone)
	quality, _ := arguments.Get(0).(neutral_zone.Quality)
	return quality
}

func (this *ZoneTierServiceMock) GetGuardQuality(
	zoneName string,
	zones []entities.Zone,
	playerNames []string) neutral_zone.Quality {
	arguments := this.Called(zoneName, zones, playerNames)
	quality, _ := arguments.Get(0).(neutral_zone.Quality)
	return quality
}

func (this *ZoneTierServiceMock) GetConnectionGuardQuality(
	zoneA, zoneB string,
	zones []entities.Zone,
	playerNames []string) neutral_zone.Quality {
	arguments := this.Called(zoneA, zoneB, zones, playerNames)
	quality, _ := arguments.Get(0).(neutral_zone.Quality)
	return quality
}
