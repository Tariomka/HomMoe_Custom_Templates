package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTemplateNameIsEmpty_ReturnsNoTemplateNameError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	stateDto := editor_state_model.NewDefaultEditorStateModel()
	stateDto.TemplateName = ""

	// Act
	_, err := handler.GenerateTemplate(toDto(stateDto))

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoTemplateName)
}

func TestWhenStateIsDefault_ReturnsNoError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	stateDto := editor_state_model.NewDefaultEditorStateModel()

	// Act
	_, err := handler.GenerateTemplate(toDto(stateDto))

	// Assert
	assert.NoError(t, err)
}

func TestWhenStateIsDefault_ReturnsGeneratedTemplate(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	stateDto := editor_state_model.NewDefaultEditorStateModel()

	// Act
	loadDto, err := handler.GenerateTemplate(toDto(stateDto))

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, loadDto.Template)
}

func TestWhenStateCarriesCustomName_GeneratedTemplateUsesThatName(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	stateDto := editor_state_model.NewDefaultEditorStateModel()
	stateDto.TemplateName = gofakeit.ProductName()

	// Act
	loadDto, err := handler.GenerateTemplate(toDto(stateDto))

	// Assert
	require.NoError(t, err)
	assert.Equal(t, stateDto.TemplateName, loadDto.Template.Name)
}

func TestWhenStateIsDefault_GeneratedTemplateHasOneVariant(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	stateDto := editor_state_model.NewDefaultEditorStateModel()

	// Act
	loadDto, err := handler.GenerateTemplate(toDto(stateDto))

	// Assert
	require.NoError(t, err)
	assert.Len(t, loadDto.Template.Variants, 1)
}

func TestWhenPlayerCountIsAboveMaximum_ReturnsValidationWarning(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	stateDto := editor_state_model.NewDefaultEditorStateModel()
	stateDto.PlayerCount = 50

	// Act
	loadDto, err := handler.GenerateTemplate(toDto(stateDto))

	// Assert
	require.NoError(t, err)
	assert.Equal(t, []string{"playerCount 50 is outside [2, 8]"}, loadDto.Warnings)
}
