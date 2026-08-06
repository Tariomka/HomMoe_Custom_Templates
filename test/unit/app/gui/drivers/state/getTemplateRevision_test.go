package state_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newRevisionState returns a State whose handler always generates the default
// template, plus its mock so manual-edit expectations can be added.
func newRevisionState() (*drivers.State, *test_helpers.TemplateHandlerMock) {
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplate()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	return drivers.NewUIState(handlerMock, test_helpers.NewFileSystemHandler(), false), handlerMock
}

func TestWhenNothingWasGenerated_TemplateRevisionIsZero(t *testing.T) {
	t.Parallel()
	// Arrange
	state, _ := newRevisionState()

	// Act
	actual := state.GetTemplateRevision()

	// Assert
	assert.Zero(t, actual)
}

func TestWhenTemplateIsGenerated_TemplateRevisionAdvances(t *testing.T) {
	t.Parallel()
	// Arrange
	state, _ := newRevisionState()
	before := state.GetTemplateRevision()

	// Act
	state.Generate()

	// Assert
	assert.Greater(t, state.GetTemplateRevision(), before)
}

func TestWhenTemplateIsRegenerated_TemplateRevisionAdvancesAgain(t *testing.T) {
	t.Parallel()
	// Arrange
	state, _ := newRevisionState()
	state.Generate()
	before := state.GetTemplateRevision()

	// Act
	state.Generate()

	// Assert
	assert.Greater(t, state.GetTemplateRevision(), before)
}

func TestWhenManualEditsAreApplied_TemplateRevisionAdvances(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock := newRevisionState()
	state.Generate()
	template := state.GetLastTemplate()
	updatedTemplate := test_helpers.GetDefaultTemplate()
	updatedTemplate.Name = gofakeit.ProductName()
	handlerMock.On("UpdateTemplate", mock.Anything).
		Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, nil)
	before := state.GetTemplateRevision()

	// Act
	state.ApplyEditedZones(template.Variants[0].Zones, template.Variants[0].Connections)

	// Assert
	assert.Greater(t, state.GetTemplateRevision(), before)
}

func TestWhenStateIsReset_TemplateRevisionAdvances(t *testing.T) {
	t.Parallel()
	// Arrange
	state, _ := newRevisionState()
	state.Generate()
	before := state.GetTemplateRevision()

	// Act
	state.Reset()

	// Assert
	assert.Greater(t, state.GetTemplateRevision(), before)
}
