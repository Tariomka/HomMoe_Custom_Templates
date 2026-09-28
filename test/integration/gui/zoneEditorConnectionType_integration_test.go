//go:build integration_test && gui

package gui_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
)

// The type dropdown lists Direct and Portal, but a loaded file can carry any
// type and portal placement rules on any connection. These tests seed such a
// connection through the state driver - the only way to reach one, since the
// editor's own controls cannot produce it - then drive the dropdown with real
// clicks and read back what Apply committed.

const (
	proximityConnectionTypeLabel = "Proximity"
	noneConnectionTypeLabel      = "(none)"
)

// connectionTypeRecord is the part of a connection the type dropdown decides.
type connectionTypeRecord struct {
	ConnectionType string
	HasPortalRules bool
}

func typeRecordOf(connection template_entity.Connection) connectionTypeRecord {
	return connectionTypeRecord{
		ConnectionType: connection.ConnectionType,
		HasPortalRules: len(connection.PortalPlacementRulesFrom) > 0 || len(connection.PortalPlacementRulesTo) > 0,
	}
}

// openZoneEditorWithEditedHubPortal commits the Geometric Hub layout as a manual
// edit with its Hub-to-Spawn-A portal rewritten by edit, opens the zone editor
// and selects that connection. The generated portal carries placement rules on
// both ends, so an edit that only changes the type leaves a rule-bearing
// connection behind. The between-zone road policy then stamps the road setting
// on every non-portal connection.
func openZoneEditorWithEditedHubPortal(
	t *testing.T,
	generateRoads bool,
	edit func(connection *template_model.Connection)) (
	*integration_common.AppRunner, *integration_common.ZoneEditorHandler) {
	t.Helper()
	runner := integration_common.NewAppRunner(t)
	tab := integration_common.NewHandler(runner).
		WithFixtureDirectory().
		ClickLayoutAndZonesTab().
		SelectTopology(geometricHubLayout)
	if !generateRoads {
		tab = tab.ToggleGenerateRoads()
	}
	runner.ApplyManualConnectionEdit(func(connections []template_model.Connection) {
		for index := range connections {
			if connections[index].Name == hubToSpawnAName {
				edit(&connections[index])
				return
			}
		}
		t.Fatalf("the generated layout has no connection called %q", hubToSpawnAName)
	})
	zoneEditor := tab.OpenZoneEditor()
	zoneEditor.ClickConnection(hubToSpawnAName)

	return runner, zoneEditor
}

func withType(connectionType string) func(connection *template_model.Connection) {
	return func(connection *template_model.Connection) { connection.ConnectionType = connectionType }
}

func withTypeAndNoRules(connectionType string) func(connection *template_model.Connection) {
	return func(connection *template_model.Connection) {
		connection.ConnectionType = connectionType
		connection.PortalPlacementRulesFrom = nil
		connection.PortalPlacementRulesTo = nil
	}
}

// Selecting a connection used to write the dropdown's fallback, Direct, over any
// type the dropdown does not list, on every frame the panel was drawn.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionWithAnUnlistedTypeIsSelected_ItsTypeSurvivesApply(t *testing.T) {
	for caseName, connectionType := range map[string]string{
		"WhenTypeIsProximity_KeepsProximity":           proximityConnectionTypeLabel,
		"WhenTypeIsGladiatorArena_KeepsGladiatorArena": "GladiatorArena",
		"WhenTypeIsEmpty_KeepsItEmpty":                 "",
		"WhenTypeIsALowerCasePortal_KeepsTheLowerCase": "portal",
	} {
		t.Run(caseName, func(t *testing.T) {
			// Arrange
			runner, zoneEditor := openZoneEditorWithEditedHubPortal(t, true, withTypeAndNoRules(connectionType))
			idleFrames(runner)

			// Act
			zoneEditor.ClickApply()

			// Assert
			assert.Equal(t, connectionType,
				manualConnectionSave(t, runner, hubToSpawnAName).Connection.ConnectionType)
		})
	}
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionHasAnUnlistedType_TheDropdownAlsoOffersIt(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditorWithEditedHubPortal(t, true, withTypeAndNoRules(proximityConnectionTypeLabel))

	// Act
	labels := zoneEditor.Dialog().ConnectionTypeLabels()

	// Assert
	assert.Equal(t,
		[]string{directConnectionTypeLabel, portalConnectionTypeLabel, proximityConnectionTypeLabel},
		labels)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionHasALowerCasedDirectType_TheDropdownOffersItAsItsOwnOption(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditorWithEditedHubPortal(t, true, withTypeAndNoRules("direct"))

	// Act
	labels := zoneEditor.Dialog().ConnectionTypeLabels()

	// Assert
	assert.Equal(t, []string{directConnectionTypeLabel, portalConnectionTypeLabel, "direct"}, labels)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAConnectionHasNoType_TheDropdownShowsNone(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditorWithEditedHubPortal(t, true, withTypeAndNoRules(""))

	// Act
	label := zoneEditor.Dialog().SelectedConnectionTypeLabel()

	// Assert
	assert.Equal(t, noneConnectionTypeLabel, label)
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenADirectConnectionCarriesPortalRules_TheDropdownShowsPortal(t *testing.T) {
	// Arrange
	_, zoneEditor := openZoneEditorWithEditedHubPortal(t, true, withType(directConnectionTypeLabel))

	// Act
	label := zoneEditor.Dialog().SelectedConnectionTypeLabel()

	// Assert
	assert.Equal(t, portalConnectionTypeLabel, label)
}

// Showing the effective type is display only: the stored type and the rules are
// untouched until the user actually picks a different type.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenADirectConnectionWithPortalRulesIsOnlySelected_ItIsAppliedUnchanged(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditorWithEditedHubPortal(t, true, withType(directConnectionTypeLabel))
	idleFrames(runner)

	// Act
	zoneEditor.ClickApply()

	// Assert
	assert.Equal(t,
		connectionTypeRecord{ConnectionType: directConnectionTypeLabel, HasPortalRules: true},
		typeRecordOf(manualConnectionSave(t, runner, hubToSpawnAName).Connection))
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenDirectIsPickedForAConnectionShownAsPortal_ItsPortalRulesAreCleared(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditorWithEditedHubPortal(t, true, withType(directConnectionTypeLabel))

	// Act
	zoneEditor.SelectConnectionType(directConnectionTypeLabel).ClickApply()

	// Assert
	assert.Equal(t,
		connectionTypeRecord{ConnectionType: directConnectionTypeLabel, HasPortalRules: false},
		typeRecordOf(manualConnectionSave(t, runner, hubToSpawnAName).Connection))
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenPortalIsPickedAgainAfterDirect_TheClearedRulesStayCleared(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditorWithEditedHubPortal(t, true, withType(directConnectionTypeLabel))
	zoneEditor.SelectConnectionType(directConnectionTypeLabel)

	// Act
	zoneEditor.SelectConnectionType(portalConnectionTypeLabel).ClickApply()

	// Assert
	assert.Equal(t,
		connectionTypeRecord{ConnectionType: portalConnectionTypeLabel, HasPortalRules: false},
		typeRecordOf(manualConnectionSave(t, runner, hubToSpawnAName).Connection))
}

// A Proximity connection carrying portal rules reads Portal, and picking its own
// listed type is a real change away from Portal, so the rules go.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheStoredTypeIsPickedForAConnectionShownAsPortal_ItsPortalRulesAreCleared(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditorWithEditedHubPortal(t, true, withType(proximityConnectionTypeLabel))

	// Act
	zoneEditor.SelectConnectionType(proximityConnectionTypeLabel).ClickApply()

	// Assert
	assert.Equal(t,
		connectionTypeRecord{ConnectionType: proximityConnectionTypeLabel, HasPortalRules: false},
		typeRecordOf(manualConnectionSave(t, runner, hubToSpawnAName).Connection))
}

//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenAPortalWithPlacementRulesIsSelected_ItKeepsItsRulesOnApply(t *testing.T) {
	// Arrange
	runner, zoneEditor := openZoneEditorWithEditedHubPortal(t, true, withType(portalConnectionTypeLabel))
	idleFrames(runner)

	// Act
	zoneEditor.ClickApply()

	// Assert
	assert.Equal(t,
		connectionTypeRecord{ConnectionType: portalConnectionTypeLabel, HasPortalRules: true},
		typeRecordOf(manualConnectionSave(t, runner, hubToSpawnAName).Connection))
}
