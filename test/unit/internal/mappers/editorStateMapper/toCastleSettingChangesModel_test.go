package editorStateMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenCastleSettingChangesDtoIsMapped_CopiesEveryField(t *testing.T) {
	t.Parallel()
	// Arrange
	dto := editor_state_dto.CastleSettingChangesDto{
		PlayerCastles: true,
		NeutralSimple: true,
		NeutralLowest: true,
		NeutralLow:    true,
		NeutralMedium: true,
		NeutralHigh:   true,
		Hub:           true,
	}

	// Act
	actual := mappers.NewEditorStateMapper().ToCastleSettingChangesModel(dto)

	// Assert
	assert.Equal(t, editor_state_model.CastleSettingChanges(dto), actual)
}
