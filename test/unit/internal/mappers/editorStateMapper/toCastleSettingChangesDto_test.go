package editorStateMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenCastleSettingChangesAreMappedToDto_CopiesEveryField(t *testing.T) {
	t.Parallel()
	// Arrange
	changes := editor_state_model.CastleSettingChanges{
		PlayerCastles: true,
		NeutralSimple: true,
		NeutralLowest: true,
		NeutralLow:    true,
		NeutralMedium: true,
		NeutralHigh:   true,
		Hub:           true,
	}

	// Act
	actual := mappers.NewEditorStateMapper().ToCastleSettingChangesDto(changes)

	// Assert
	assert.Equal(t, editor_state_dto.CastleSettingChangesDto(changes), actual)
}
