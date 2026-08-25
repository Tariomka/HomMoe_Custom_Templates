package templateHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenCastleSettingsAreReapplied_ReturnsTheRequestedZones(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	zones := []entities.Zone{{Name: gofakeit.Word()}}
	fixture.mapper.On("FromEditorState", mock.Anything).Return(namedConfiguration())
	fixture.manualReapply.On("ApplyCastleSettingChanges", mock.Anything, mock.Anything, mock.Anything).Return()

	// Act
	reappliedZones := fixture.handler.ReapplyCastleSettings(dtos.CastleSettingsReapplyRequestDto{Zones: zones})

	// Assert
	assert.Equal(t, zones, reappliedZones)
}

func TestWhenCastleSettingsAreReapplied_PassesTheMappedConfigurationToTheReapplyService(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	zones := []entities.Zone{{Name: gofakeit.Word()}}
	changes := editor_state_model.CastleSettingChanges{PlayerCastles: true}
	configuration := namedConfiguration()
	fixture.mapper.On("FromEditorState", state).Return(configuration)
	fixture.manualReapply.On("ApplyCastleSettingChanges", zones, changes, configuration).Return()

	// Act
	_ = fixture.handler.ReapplyCastleSettings(dtos.CastleSettingsReapplyRequestDto{
		Zones:       zones,
		Changes:     fixture.editorStateMapper.ToCastleSettingChangesDto(changes),
		EditorState: fixture.editorStateMapper.ToDto(state),
	})

	// Assert
	fixture.manualReapply.AssertCalled(t, "ApplyCastleSettingChanges", zones, changes, configuration)
}
