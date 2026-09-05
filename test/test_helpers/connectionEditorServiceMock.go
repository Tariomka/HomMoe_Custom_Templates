package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
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
	zones []template_model.Zone,
	playerZoneNames map[string]bool) template_model.Connection {
	arguments := this.Called(from, to, zones, playerZoneNames)
	connection, _ := arguments.Get(0).(template_model.Connection)
	return connection
}

func (this *ConnectionEditorServiceMock) FindIsolatedZones(
	zones []template_model.Zone,
	connections []template_model.Connection) []string {
	arguments := this.Called(zones, connections)
	names, _ := arguments.Get(0).([]string)
	return names
}

func (this *ConnectionEditorServiceMock) ComputeHasErrors(
	zones []template_model.Zone,
	connections []template_model.Connection) bool {
	arguments := this.Called(zones, connections)
	return arguments.Bool(0)
}

func (this *ConnectionEditorServiceMock) HasDuplicateName(
	connections []template_model.Connection,
	current *template_model.Connection) bool {
	arguments := this.Called(connections, current)
	return arguments.Bool(0)
}
