package stateGeneration_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenGenerationSucceeds_LastTemplateIsStored(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplate()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	state := drivers.NewUIStateWithHandler(handlerMock)

	// Act
	state.Generate()

	// Assert
	assert.Equal(t, &template, state.GetLastTemplate())
}

func TestWhenGenerationSucceeds_StatusReportsGeneratedTemplate(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplate()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	state := drivers.NewUIStateWithHandler(handlerMock)

	// Act
	state.Generate()

	// Assert
	message, _ := state.GetStatus()
	assert.Contains(t, message, "generated with latest changes")
}

func TestWhenGenerationFails_ErrorStatusIsSet(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}
	handlerMock.On("GenerateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{}, gofakeit.ErrorValidation())
	state := drivers.NewUIStateWithHandler(handlerMock)

	// Act
	state.Generate()

	// Assert
	_, isError := state.GetStatus()
	assert.True(t, isError)
}

func TestWhenGenerationFails_NoTemplateIsStored(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}
	handlerMock.On("GenerateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{}, gofakeit.ErrorValidation())
	state := drivers.NewUIStateWithHandler(handlerMock)

	// Act
	state.Generate()

	// Assert
	assert.Nil(t, state.GetLastTemplate())
}
