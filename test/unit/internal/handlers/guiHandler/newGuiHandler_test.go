package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenHandlerIsConstructed_ReturnsNonNilHandler(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	handler := newProductionGuiHandler()

	// Assert
	assert.NotNil(t, handler)
}

func TestWhenCollaboratorsAreProvided_ReturnsHandler(t *testing.T) {
	t.Parallel()
	// Arrange
	stub := &handlerDependenciesStub{}

	// Act
	handler := stub.newHandler()

	// Assert
	assert.NotNil(t, handler)
}

func TestWhenGenerateTemplateCalled_DelegatesToTemplateHandler(t *testing.T) {
	t.Parallel()
	// Arrange
	stub := &handlerDependenciesStub{}
	handler := stub.newHandler()
	require.NotNil(t, handler)

	// Act
	_, _ = handler.GenerateTemplate(editor_state_dto.EditorStateDto{})

	// Assert
	assert.True(t, stub.templateWorkflowCalled)
}

func TestWhenSaveTemplateCalled_DelegatesToTemplateHandler(t *testing.T) {
	t.Parallel()
	// Arrange
	stub := &handlerDependenciesStub{}
	handler := stub.newHandler()
	require.NotNil(t, handler)

	// Act
	_, _ = handler.SaveTemplate(dtos.TemplateSaveDto{})

	// Assert
	assert.True(t, stub.templatePersistenceCalled)
}

func TestWhenSaveStateCalled_DelegatesToStateHandler(t *testing.T) {
	t.Parallel()
	// Arrange
	stub := &handlerDependenciesStub{}
	handler := stub.newHandler()
	require.NotNil(t, handler)

	// Act
	_, _ = handler.SaveState(editor_state_dto.EditorStateSaveDto{})

	// Assert
	assert.True(t, stub.statePersistenceCalled)
}

func TestWhenBuildPreviewLayoutCalled_DelegatesToPreviewHandler(t *testing.T) {
	t.Parallel()
	// Arrange
	stub := &handlerDependenciesStub{}
	handler := stub.newHandler()
	require.NotNil(t, handler)

	// Act
	_, _ = handler.BuildPreviewLayout(dtos.PreviewLayoutRequestDto{})

	// Assert
	assert.True(t, stub.previewCalled)
}

func TestWhenContentRuleOptionsRequested_DelegatesToContentRuleHandler(t *testing.T) {
	t.Parallel()
	// Arrange
	stub := &handlerDependenciesStub{}
	handler := stub.newHandler()
	require.NotNil(t, handler)

	// Act
	handler.GetContentRuleEditorOptions(models.SidMapping{})

	// Assert
	assert.True(t, stub.contentRuleCalled)
}

func TestWhenCastleCountRequested_DelegatesToZoneEditorHandler(t *testing.T) {
	t.Parallel()
	// Arrange
	stub := &handlerDependenciesStub{}
	handler := stub.newHandler()
	require.NotNil(t, handler)

	// Act
	handler.CountZoneCastles(entities.Zone{})

	// Assert
	assert.True(t, stub.zoneEditorCalled)
}

func TestWhenSpellCountLabelRequested_DelegatesToBonusHandler(t *testing.T) {
	t.Parallel()
	// Arrange
	stub := &handlerDependenciesStub{}
	handler := stub.newHandler()
	require.NotNil(t, handler)

	// Act
	handler.GetSpellCountLabel(0)

	// Assert
	assert.True(t, stub.bonusCalled)
}
