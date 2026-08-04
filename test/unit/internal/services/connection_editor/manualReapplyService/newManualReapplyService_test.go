package manualReapplyService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenDependenciesAreProvided_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	service := newManualReapplyService()

	// Assert
	assert.NotNil(t, service)
}

func TestWhenDependenciesAreNil_ReturnsUsableService(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewManualReapplyService(nil, nil, nil)
	configuration := config.NewGeneratorConfig()
	configuration.ZoneConfiguration.Advanced.NeutralMediumCastlesPerZone = 2
	zones := []entities.Zone{makeNeutralZone("G", neutral_zone.QualityMedium, 1)}

	// Act
	service.ApplyCastleSettingChanges(
		zones,
		editor_state_dto.CastleSettingChanges{NeutralMedium: true},
		configuration,
	)

	// Assert
	assert.Equal(t, 2, connection_editor.NewDefaultZoneEditorService().CountZoneCastles(zones[0]))
}
