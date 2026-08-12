package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/stretchr/testify/mock"
)

// ConnectionEditorServiceMock is a testify mock of
// connection_editor.IConnectionEditorService, used to unit-test collaborators
// without the real connection graph analysis.
type ConnectionEditorServiceMock struct {
	mock.Mock
}

func (this *ConnectionEditorServiceMock) NewDefaultConnection(
	from string,
	to string,
	zones []entities.Zone,
	playerZoneNames map[string]bool) entities.Connection {
	arguments := this.Called(from, to, zones, playerZoneNames)
	connection, _ := arguments.Get(0).(entities.Connection)
	return connection
}

func (this *ConnectionEditorServiceMock) FindIsolatedZones(
	zones []entities.Zone,
	connections []entities.Connection) []string {
	arguments := this.Called(zones, connections)
	names, _ := arguments.Get(0).([]string)
	return names
}

func (this *ConnectionEditorServiceMock) ComputeHasErrors(
	zones []entities.Zone,
	connections []entities.Connection) bool {
	arguments := this.Called(zones, connections)
	return arguments.Bool(0)
}

func (this *ConnectionEditorServiceMock) HasDuplicateName(
	connections []entities.Connection,
	current *entities.Connection) bool {
	arguments := this.Called(connections, current)
	return arguments.Bool(0)
}
