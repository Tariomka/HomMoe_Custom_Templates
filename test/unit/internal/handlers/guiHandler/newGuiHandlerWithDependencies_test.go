package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenDependencyIsMissing_ReturnsError(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		remove        func(*handlers.GUIHandlerDependencies)
		expectedError string
	}{
		{
			name: "WhenTemplateWorkflowIsMissing_ReturnsError",
			remove: func(dependencies *handlers.GUIHandlerDependencies) {
				dependencies.TemplateWorkflow = nil
			},
			expectedError: "template workflow handler is required",
		},
		{
			name: "WhenStatePersistenceIsMissing_ReturnsError",
			remove: func(dependencies *handlers.GUIHandlerDependencies) {
				dependencies.StatePersistence = nil
			},
			expectedError: "state persistence handler is required",
		},
		{
			name: "WhenTemplatePersistenceIsMissing_ReturnsError",
			remove: func(dependencies *handlers.GUIHandlerDependencies) {
				dependencies.TemplatePersistence = nil
			},
			expectedError: "template persistence handler is required",
		},
		{
			name: "WhenPreviewIsMissing_ReturnsError",
			remove: func(dependencies *handlers.GUIHandlerDependencies) {
				dependencies.Preview = nil
			},
			expectedError: "preview handler is required",
		},
		{
			name: "WhenContentRuleIsMissing_ReturnsError",
			remove: func(dependencies *handlers.GUIHandlerDependencies) {
				dependencies.ContentRule = nil
			},
			expectedError: "content rule handler is required",
		},
		{
			name: "WhenZoneEditorIsMissing_ReturnsError",
			remove: func(dependencies *handlers.GUIHandlerDependencies) {
				dependencies.ZoneEditor = nil
			},
			expectedError: "zone editor handler is required",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			stub := &handlerDependenciesStub{}
			dependencies := stub.dependencies()
			testCase.remove(&dependencies)

			// Act
			_, err := handlers.NewGuiHandlerWithDependencies(dependencies)

			// Assert
			assert.EqualError(t, err, testCase.expectedError)
		})
	}
}

func TestWhenDependenciesAreProvided_ReturnsHandler(t *testing.T) {
	t.Parallel()
	// Arrange
	stub := &handlerDependenciesStub{}

	// Act
	handler, err := handlers.NewGuiHandlerWithDependencies(stub.dependencies())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, handler)
}

func TestWhenGenerateTemplateCalled_DelegatesToTemplateWorkflow(t *testing.T) {
	t.Parallel()
	// Arrange
	stub := &handlerDependenciesStub{}
	handler, err := handlers.NewGuiHandlerWithDependencies(stub.dependencies())
	require.NoError(t, err)

	// Act
	_, _ = handler.GenerateTemplate(dtos.EditorStateDto{})

	// Assert
	assert.True(t, stub.templateWorkflowCalled)
}

func TestWhenSaveStateCalled_DelegatesToStatePersistence(t *testing.T) {
	t.Parallel()
	// Arrange
	stub := &handlerDependenciesStub{}
	handler, err := handlers.NewGuiHandlerWithDependencies(stub.dependencies())
	require.NoError(t, err)

	// Act
	_, _ = handler.SaveState(dtos.EditorStateSaveDto{})

	// Assert
	assert.True(t, stub.statePersistenceCalled)
}

func TestWhenSaveTemplateCalled_DelegatesToTemplatePersistence(t *testing.T) {
	t.Parallel()
	// Arrange
	stub := &handlerDependenciesStub{}
	handler, err := handlers.NewGuiHandlerWithDependencies(stub.dependencies())
	require.NoError(t, err)

	// Act
	_, _ = handler.SaveTemplate(dtos.TemplateSaveDto{})

	// Assert
	assert.True(t, stub.templatePersistenceCalled)
}

func TestWhenBuildPreviewCalled_DelegatesToPreview(t *testing.T) {
	t.Parallel()
	// Arrange
	stub := &handlerDependenciesStub{}
	handler, err := handlers.NewGuiHandlerWithDependencies(stub.dependencies())
	require.NoError(t, err)

	// Act
	_, _ = handler.BuildPreviewLayout(dtos.PreviewLayoutRequestDto{})

	// Assert
	assert.True(t, stub.previewCalled)
}

func TestWhenContentOptionsRequested_DelegatesToContentRule(t *testing.T) {
	t.Parallel()
	// Arrange
	stub := &handlerDependenciesStub{}
	handler, err := handlers.NewGuiHandlerWithDependencies(stub.dependencies())
	require.NoError(t, err)

	// Act
	handler.GetContentRuleEditorOptions(models.SidMapping{})

	// Assert
	assert.True(t, stub.contentRuleCalled)
}

func TestWhenCastleCountRequested_DelegatesToZoneEditor(t *testing.T) {
	t.Parallel()
	// Arrange
	stub := &handlerDependenciesStub{}
	handler, err := handlers.NewGuiHandlerWithDependencies(stub.dependencies())
	require.NoError(t, err)

	// Act
	handler.CountZoneCastles(entities.Zone{})

	// Assert
	assert.True(t, stub.zoneEditorCalled)
}
