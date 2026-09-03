package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/mock"
)

// ZoneEditorServiceMock is a testify mock of
// connection_editor.IZoneEditorService, used to unit-test collaborators
// without the real zone/road rebuilding.
type ZoneEditorServiceMock struct {
	mock.Mock
}

func (this *ZoneEditorServiceMock) EnsureConnectionNames(connections []entities.Connection) {
	this.Called(connections)
}

func (this *ZoneEditorServiceMock) RebuildZoneConnectionRoads(
	zones []template_model.Zone,
	connections []entities.Connection) {
	this.Called(zones, connections)
}

func (this *ZoneEditorServiceMock) RebuildCastleRoads(zone *template_model.Zone) {
	this.Called(zone)
}

func (this *ZoneEditorServiceMock) NextFreeZoneLabel(zones []template_model.Zone) string {
	arguments := this.Called(zones)
	return arguments.String(0)
}

func (this *ZoneEditorServiceMock) NewDefaultNeutralZone(
	label string,
	quality neutral_zone.Quality,
	castleCount int,
	generateRoads bool,
	tuning models.GenerationTuning) template_model.Zone {
	arguments := this.Called(label, quality, castleCount, generateRoads, tuning)
	zone, _ := arguments.Get(0).(template_model.Zone)
	return zone
}

func (this *ZoneEditorServiceMock) CountZoneCastles(zone template_model.Zone) int {
	arguments := this.Called(zone)
	return arguments.Int(0)
}

func (this *ZoneEditorServiceMock) ApplyNeutralZoneQuality(
	zone *template_model.Zone,
	quality neutral_zone.Quality,
	castleCount int,
	tuning models.GenerationTuning) {
	this.Called(zone, quality, castleCount, tuning)
}

func (this *ZoneEditorServiceMock) CanDeleteZone(zoneName string, playerZoneNames map[string]bool) bool {
	arguments := this.Called(zoneName, playerZoneNames)
	return arguments.Bool(0)
}

func (this *ZoneEditorServiceMock) RemoveZone(
	zones []template_model.Zone,
	connections []entities.Connection,
	zoneName string) ([]template_model.Zone, []entities.Connection) {
	arguments := this.Called(zones, connections, zoneName)
	remainingZones, _ := arguments.Get(0).([]template_model.Zone)
	remainingConnections, _ := arguments.Get(1).([]entities.Connection)
	return remainingZones, remainingConnections
}

func (this *ZoneEditorServiceMock) FindOpenPosition(occupied [][2]float64) [2]float64 {
	arguments := this.Called(occupied)
	position, _ := arguments.Get(0).([2]float64)
	return position
}
