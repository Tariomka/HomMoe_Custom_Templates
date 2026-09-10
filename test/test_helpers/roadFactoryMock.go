package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/mock"
)

// RoadFactoryMock is a testify mock of zone_interfaces.IRoadFactory, used to
// hand a collaborator road candidates the real factory would never build.
type RoadFactoryMock struct {
	mock.Mock
}

func (this *RoadFactoryMock) CreateConnectorZoneRoads(
	connectionNames []string,
	generateRoads bool) []template_model.Road {
	arguments := this.Called(connectionNames, generateRoads)
	roads, _ := arguments.Get(0).([]template_model.Road)
	return roads
}

func (this *RoadFactoryMock) CreateOuterZoneRoads(
	connectionNames []string,
	mainObjectCount int,
	footholdCount int,
	generateRoads bool) []template_model.Road {
	arguments := this.Called(connectionNames, mainObjectCount, footholdCount, generateRoads)
	roads, _ := arguments.Get(0).([]template_model.Road)
	return roads
}
