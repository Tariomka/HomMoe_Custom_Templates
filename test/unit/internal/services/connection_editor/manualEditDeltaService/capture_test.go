package manualEditDeltaService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenNothingWasEdited_StoresNoZoneOverrides(t *testing.T) {
	t.Parallel()
	// Arrange
	request := captureRequest(newConfiguration(), unchanged)

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Empty(t, overriddenZoneNames(settings))
}

func TestWhenNothingWasEdited_StoresNoConnectionOverrides(t *testing.T) {
	t.Parallel()
	// Arrange
	request := captureRequest(newConfiguration(), unchanged)

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Empty(t, overriddenConnectionNames(settings))
}

func TestWhenAZoneWasMoved_StoresItsManualPosition(t *testing.T) {
	t.Parallel()
	// Arrange
	request := captureRequest(newConfiguration(),
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			zones[1].ManualPosition = new(data.NewVec2(0.9, 0.1))
			return zones, connections
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Equal(t, new(data.NewVec2(0.9, 0.1)), zoneEdit(settings, goldZoneName).ManualPosition)
}

func TestWhenAZoneWasMoved_StoresNoOverrideForIt(t *testing.T) {
	t.Parallel()
	// Arrange
	request := captureRequest(newConfiguration(),
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			zones[1].ManualPosition = new(data.NewVec2(0.9, 0.1))
			return zones, connections
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Empty(t, overriddenZoneNames(settings))
}

func TestWhenAFieldWasEdited_StoresItAsAnOverride(t *testing.T) {
	t.Parallel()
	// Arrange
	request := captureRequest(newConfiguration(),
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			zones[1].Size = 2.5
			return zones, connections
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Equal(t, new(2.5), zoneEdit(settings, goldZoneName).FieldOverrides.Size)
}

func TestWhenAFieldWasCleared_ListsItAsUnset(t *testing.T) {
	t.Parallel()
	// Arrange
	request := captureRequest(newConfiguration(),
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			zones[1].GuardedContentPool = nil
			return zones, connections
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Equal(t, []string{"guardedContentPool"}, zoneEdit(settings, goldZoneName).Unset)
}

func TestWhenAFieldWasEmptied_StoresTheEmptyValueAsAnOverride(t *testing.T) {
	t.Parallel()
	// Arrange
	request := captureRequest(newConfiguration(),
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			zones[1].GuardedContentPool = []string{}
			return zones, connections
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Equal(t, []string{}, zoneEdit(settings, goldZoneName).FieldOverrides.GuardedContentPool)
}

func TestWhenANeutralZoneWasReTiered_RecordsItsCastlesAsReprofiled(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	request := captureRequest(configuration,
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			test_helpers.NewZoneEditorService().ApplyNeutralZoneQuality(
				&zones[2], neutral_zone.QualityLow, 2, newTuning(configuration, len(zones)))
			return zones, connections
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Equal(t, new(2), zoneEdit(settings, silverZoneName).Castles)
}

func TestWhenANeutralZoneWasReTiered_DerivesItFromTheFactory(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	request := captureRequest(configuration,
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			zones[2] = factoryZone("H", neutral_zone.QualityLow, 2, configuration, len(zones))
			zones[2].Name = silverZoneName
			zones[2].GeneratorPosition = new(data.NewVec2(0.6, 0.5))
			return zones, connections
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Empty(t, overriddenZoneNames(settings))
}

func TestWhenOnlyTheCastleCountChanged_RecordsTheZoneAsReprofiled(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	request := captureRequest(configuration,
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			test_helpers.NewZoneEditorService().ApplyNeutralZoneQuality(
				&zones[1], neutral_zone.QualityHigh, 3, newTuning(configuration, len(zones)))
			return zones, connections
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Equal(t, new(3), zoneEdit(settings, goldZoneName).Castles)
}

func TestWhenAPlayerZoneChangedCastles_DoesNotRecordItAsReprofiled(t *testing.T) {
	t.Parallel()
	// Arrange
	request := captureRequest(newConfiguration(),
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			zones[0].MainObjects = append(zones[0].MainObjects, template_model.MainObject{Type: "City"})
			return zones, connections
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Nil(t, zoneEdit(settings, spawnZoneName).Castles)
}

func TestWhenAZoneWasAddedInTheEditor_RecordsItAsAdded(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	request := captureRequest(configuration,
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			added := factoryZone(addedZoneLabel, neutral_zone.QualityMedium, 1, configuration, len(zones)+1)
			return append(zones, added), connections
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.True(t, zoneEdit(settings, addedZoneName).IsAdded())
}

func TestWhenAnAddedZoneMatchesTheFactory_StoresNoOverrideForIt(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	request := captureRequest(configuration,
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			added := factoryZone(addedZoneLabel, neutral_zone.QualityMedium, 1, configuration, len(zones)+1)
			return append(zones, added), connections
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Empty(t, overriddenZoneNames(settings))
}

func TestWhenADeletedZoneLabelWasReAdded_RecordsTheZoneAsAdded(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	request := captureRequest(configuration,
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			zones[1] = factoryZone("G", neutral_zone.QualityMedium, 1, configuration, len(zones))
			return zones, connections
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.True(t, zoneEdit(settings, goldZoneName).IsAdded())
}

func TestWhenALegacySnapshotHasNoGeneratorPositions_KeepsAGeneratedNameAsGenerated(t *testing.T) {
	t.Parallel()
	// Arrange
	request := captureRequest(newConfiguration(),
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			for index := range zones {
				zones[index].GeneratorPosition = nil
			}
			return zones, connections
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.False(t, zoneEdit(settings, goldZoneName).IsAdded())
}

func TestWhenALegacySnapshotHasAnUnknownName_RecordsItAsAdded(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	request := captureRequest(configuration,
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			for index := range zones {
				zones[index].GeneratorPosition = nil
			}
			added := factoryZone(addedZoneLabel, neutral_zone.QualityMedium, 1, configuration, len(zones)+1)
			return append(zones, added), connections
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.True(t, zoneEdit(settings, addedZoneName).IsAdded())
}

func TestWhenAZoneWasDeleted_RecordsOnlyTheRemainingZonesInOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	request := captureRequest(newConfiguration(),
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			return []template_model.Zone{zones[2], zones[0]}, connections[:1]
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Equal(t, []string{silverZoneName, spawnZoneName},
		[]string{settings.ManualZones[0].Name, settings.ManualZones[1].Name})
}

func TestWhenTheLayoutIsNonDeterministic_StoresUnchangedRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	request := captureRequest(newConfiguration(), unchanged)
	request.NonDeterministic = true

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Equal(t, template_model.ToZoneEntity(request.Zones[1]).Roads,
		zoneEdit(settings, goldZoneName).FieldOverrides.Roads)
}

func TestWhenTheLayoutIsNonDeterministicAndRoadsAreNil_ListsRoadsAsUnset(t *testing.T) {
	t.Parallel()
	// Arrange
	request := captureRequest(newConfiguration(),
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			zones[0].Roads = nil
			return zones, connections
		})
	request.FreshZones[0].Roads = nil
	request.NonDeterministic = true

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Equal(t, []string{"roads"}, zoneEdit(settings, spawnZoneName).Unset)
}

func TestWhenAConnectionGuardWasEdited_StoresItAsAnOverride(t *testing.T) {
	t.Parallel()
	// Arrange
	request := captureRequest(newConfiguration(),
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			connections[0].GuardValue = 80000
			return zones, connections
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Equal(t, new(80000), settings.ManualConnections[0].FieldOverrides.GuardValue)
}

func TestWhenAUserAddedConnectionMatchesTheFactory_StoresNoOverride(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	request := captureRequest(configuration,
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			return zones, append(connections, userAddedConnection(silverZoneName, zones, configuration))
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Empty(t, overriddenConnectionNames(settings))
}

func TestWhenAUserAddedConnectionCopiesAGeneratedOne_DerivesItFromTheFactory(t *testing.T) {
	t.Parallel()
	// Arrange
	request := captureRequest(newConfiguration(),
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			connections[0].IsUserAdded = true
			return zones, connections
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Equal(t, []string{spawnToGoldName}, overriddenConnectionNames(settings))
}

func TestWhenAConnectionEndpointsNoLongerMatch_DerivesItFromTheFactory(t *testing.T) {
	t.Parallel()
	// Arrange
	request := captureRequest(newConfiguration(),
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			connections[1].From, connections[1].To = connections[1].To, connections[1].From
			return zones, connections
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Equal(t, []string{goldToSilverName}, overriddenConnectionNames(settings))
}

func TestWhenAUserAddedConnectionChangedType_DerivesTheTypeFromTheFactory(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := newConfiguration()
	request := captureRequest(configuration,
		func(zones []template_model.Zone, connections []template_model.Connection) (
			[]template_model.Zone, []template_model.Connection) {
			added := userAddedConnection(silverZoneName, zones, configuration)
			test_helpers.NewZoneEditorService().ChangeConnectionType(&added, "Portal", configuration.GenerateRoads)
			return zones, append(connections, added)
		})

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Empty(t, overriddenConnectionNames(settings))
}

func TestWhenTheLayoutIsNonDeterministic_DerivesGeneratedConnectionsFromTheFactory(t *testing.T) {
	t.Parallel()
	// Arrange
	request := captureRequest(newConfiguration(), unchanged)
	request.NonDeterministic = true

	// Act
	settings := newService().Capture(request)

	// Assert
	assert.Len(t, overriddenConnectionNames(settings), 2)
}

func TestWhenConnectionsAreCaptured_KeepsTheirStructure(t *testing.T) {
	t.Parallel()
	// Arrange
	request := captureRequest(newConfiguration(), unchanged)

	// Act
	settings := newService().Capture(request)

	// Assert
	edited := request.Connections[1]
	stored := settings.ManualConnections[1]
	assert.Equal(t,
		[]string{edited.Name, edited.From, edited.To, edited.ConnectionType},
		[]string{stored.Name, stored.From, stored.To, stored.ConnectionType})
}
