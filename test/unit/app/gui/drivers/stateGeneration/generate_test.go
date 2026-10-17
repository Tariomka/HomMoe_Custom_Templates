package stateGeneration_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWhenGenerationSucceeds_LastTemplateIsStored(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	state := drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),

		false)

	// Act
	state.Generate()

	// Assert
	assert.Equal(t, &template, state.GetLastTemplate())
}

func TestWhenGenerationSucceeds_StatusReportsGeneratedTemplate(t *testing.T) {
	t.Parallel()
	// Arrange
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	state := drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),

		false)

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
	state := drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),

		false)

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
	state := drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),

		false)

	// Act
	state.Generate()

	// Assert
	assert.Nil(t, state.GetLastTemplate())
}

func TestWhenEditsArePendingConversion_GenerationConvertsThemOnce(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, _ := newPendingConversionState()
	convertCount := 0
	handlerMock.ConvertManualEditsFunc = func(
		request dtos.ManualEditConvertRequestDto) (dtos.ManualEditCaptureDto, error) {
		convertCount++
		return dtos.ManualEditCaptureDto{ManualEdits: request.ManualEdits}, nil
	}

	// Act
	state.Generate()

	// Assert
	assert.Equal(t, 1, convertCount)
}

func TestWhenConversionSucceeds_TheConvertedEditsAreStored(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, pendingEdits := newPendingConversionState()
	converted := editor_state_model.ManualEditSettings{ManualZones: pendingEdits.ManualZones[:1]}
	handlerMock.ConvertManualEditsFunc = func(dtos.ManualEditConvertRequestDto) (dtos.ManualEditCaptureDto, error) {
		return dtos.ManualEditCaptureDto{ManualEdits: converted}, nil
	}

	// Act
	state.Generate()

	// Assert
	stateData := state.GetStateData()
	assert.Equal(t, converted.Clone(), stateData.ManualEditSettings)
}

func TestWhenConversionSucceeds_TheGeneratorVersionIsTheCurrentOne(t *testing.T) {
	t.Parallel()
	// Arrange
	state, _, _ := newPendingConversionState()

	// Act
	state.Generate()

	// Assert
	stateData := state.GetStateData()
	assert.Equal(t, common_generator.GeneratorOutputVersion, stateData.GeneratorVersion)
}

func TestWhenConversionSucceeds_TheDocumentBecomesUnsaved(t *testing.T) {
	t.Parallel()
	// Arrange
	state, _, _ := newPendingConversionState()
	require.False(t, state.IsUnsaved())

	// Act
	state.Generate()

	// Assert
	assert.True(t, state.IsUnsaved())
}

func TestWhenConversionRunsWithoutALoadedFile_NoDialogIsOpened(t *testing.T) {
	t.Parallel()
	// Arrange
	state, _, _ := newPendingConversionState()

	// Act
	state.Generate()

	// Assert
	assert.False(t, state.GetDialogHost().IsOpen())
}

func TestWhenConversionFails_TheStateStillNeedsDerivation(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, _ := newPendingConversionState()
	expectConversionFailure(handlerMock)

	// Act
	state.Generate()

	// Assert
	stateData := state.GetStateData()
	assert.True(t, stateData.NeedsDerivation())
}

func TestWhenConversionFails_TheStatusReportsAnErrorMentioningTheConversion(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, _ := newPendingConversionState()
	expectConversionFailure(handlerMock)

	// Act
	state.Generate()

	// Assert
	message, _ := state.GetStatus()
	assert.Contains(t, message, "could not be converted")
}

func TestWhenConversionFails_TheStatusIsAnError(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, _ := newPendingConversionState()
	expectConversionFailure(handlerMock)

	// Act
	state.Generate()

	// Assert
	_, isError := state.GetStatus()
	assert.True(t, isError)
}

func TestWhenConversionFails_ExpansionIsSkipped(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, _ := newPendingConversionState()
	expectConversionFailure(handlerMock)
	expandCount := 0
	handlerMock.ExpandManualEditsFunc = func(request dtos.ManualEditExpandRequestDto) dtos.ManualEditExpansionDto {
		expandCount++
		return decodeStoredEdits(request)
	}

	// Act
	state.Generate()

	// Assert
	assert.Equal(t, 0, expandCount)
}

func TestWhenTheBaselineConversionFails_TheRetryUsesTheCurrentSettings(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock := newBaselineState()
	densityAtApply := state.GetStateData().ResourceDensityPercent
	state.UpdateState(func(settings *editor_state_model.EditorState) {
		settings.ResourceDensityPercent = densityAtApply + 50
	})
	var baselineDensities []int
	handlerMock.ConvertManualEditsFunc = func(
		request dtos.ManualEditConvertRequestDto) (dtos.ManualEditCaptureDto, error) {
		baselineDensities = append(baselineDensities, request.Baseline.ResourceDensityPercent)
		if request.Baseline.ResourceDensityPercent == densityAtApply {
			return dtos.ManualEditCaptureDto{}, gofakeit.ErrorValidation()
		}
		return dtos.ManualEditCaptureDto{ManualEdits: request.ManualEdits}, nil
	}

	// Act
	state.Generate()

	// Assert
	assert.Equal(t, []int{densityAtApply, densityAtApply + 50}, baselineDensities)
}

func TestWhenExpansionReportsWarnings_TheyAppearInTheStatus(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, _ := newPendingConversionState()
	state.UpdateState(func(settings *editor_state_model.EditorState) {
		settings.GeneratorVersion = common_generator.GeneratorOutputVersion
	})
	warning := gofakeit.LoremIpsumSentence(4)
	handlerMock.ExpandManualEditsFunc = func(request dtos.ManualEditExpandRequestDto) dtos.ManualEditExpansionDto {
		expansion := decodeStoredEdits(request)
		expansion.Warnings = []string{warning}
		return expansion
	}

	// Act
	state.Generate()

	// Assert
	message, _ := state.GetStatus()
	assert.Contains(t, message, warning)
}

func TestWhenExpansionChangesTheEdits_TheReturnedEditsAreStored(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock, pendingEdits := newPendingConversionState()
	state.UpdateState(func(settings *editor_state_model.EditorState) {
		settings.GeneratorVersion = common_generator.GeneratorOutputVersion
	})
	expanded := editor_state_model.ManualEditSettings{ManualZones: pendingEdits.ManualZones[:1]}
	handlerMock.ExpandManualEditsFunc = func(request dtos.ManualEditExpandRequestDto) dtos.ManualEditExpansionDto {
		expansion := decodeStoredEdits(request)
		expansion.ManualEdits = expanded
		expansion.Changed = true
		return expansion
	}

	// Act
	state.Generate()

	// Assert
	stateData := state.GetStateData()
	assert.Equal(t, expanded.Clone(), stateData.ManualEditSettings)
}

func TestWhenConversionFollowsABaseline_ExpansionReceivesTheBaselineToCurrentCastleChanges(t *testing.T) {
	t.Parallel()
	// Arrange
	state, handlerMock := newBaselineState()
	state.UpdateState(func(settings *editor_state_model.EditorState) {
		settings.PlayerZoneCastles++
	})
	var castleChanges editor_state_model.CastleSettingChanges
	handlerMock.ExpandManualEditsFunc = func(request dtos.ManualEditExpandRequestDto) dtos.ManualEditExpansionDto {
		castleChanges = request.CastleChanges
		return decodeStoredEdits(request)
	}

	// Act
	state.Generate()

	// Assert
	assert.True(t, castleChanges.PlayerCastles)
}

func TestWhenTheEditsWereSavedWithAnOlderGeneratorVersion_TheStatusNotesTheVersionChange(t *testing.T) {
	t.Parallel()
	// Arrange
	state, _, _ := newPendingConversionState()
	state.UpdateState(func(settings *editor_state_model.EditorState) { settings.GeneratorVersion = 99 })

	// Act
	state.Generate()

	// Assert
	message, _ := state.GetStatus()
	assert.Contains(t, message, "Saved with generator version 99")
}

func TestWhenTheEditsWereSavedWithAnOlderGeneratorVersion_TheVersionBecomesTheCurrentOne(t *testing.T) {
	t.Parallel()
	// Arrange
	state, _, _ := newPendingConversionState()
	state.UpdateState(func(settings *editor_state_model.EditorState) { settings.GeneratorVersion = 99 })

	// Act
	state.Generate()

	// Assert
	stateData := state.GetStateData()
	assert.Equal(t, common_generator.GeneratorOutputVersion, stateData.GeneratorVersion)
}

// newPendingConversionState returns a never-generated State holding an
// all-overridden manual snapshot of the default template (generator version
// 0), a handler that accepts every update, and that snapshot.
func newPendingConversionState() (
	*drivers.State, *test_helpers.TemplateHandlerMock, editor_state_model.ManualEditSettings) {
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplateModel()
	updatedTemplate := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	handlerMock.On("UpdateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, nil)
	state := drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),
		false)
	variant := template.Variants[0]
	pendingEdits := editor_state_model.ManualEditSettingsFromSnapshot(variant.Zones, variant.Connections)
	state.UpdateState(func(settings *editor_state_model.EditorState) {
		settings.ManualEditSettings = pendingEdits
		settings.GeneratorVersion = 0
	})
	return state, handlerMock, pendingEdits
}

// newBaselineState returns a generated State whose manual snapshot is pending
// conversion against the settings at Apply, because capturing failed then.
func newBaselineState() (*drivers.State, *test_helpers.TemplateHandlerMock) {
	handlerMock := &test_helpers.TemplateHandlerMock{}
	template := test_helpers.GetDefaultTemplateModel()
	updatedTemplate := test_helpers.GetDefaultTemplateModel()
	handlerMock.On("GenerateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &template}, nil)
	handlerMock.On("UpdateTemplate", mock.Anything).Return(dtos.TemplateLoadDto{Template: &updatedTemplate}, nil)
	state := drivers.NewUIState(
		handlerMock,
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),
		false)
	state.Generate()
	handlerMock.CaptureManualEditsFunc = func(dtos.ManualEditCaptureRequestDto) (dtos.ManualEditCaptureDto, error) {
		return dtos.ManualEditCaptureDto{}, gofakeit.ErrorValidation()
	}
	variant := template.Variants[0]
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{Zones: variant.Zones, Connections: variant.Connections})
	handlerMock.CaptureManualEditsFunc = nil
	return state, handlerMock
}

func expectConversionFailure(handlerMock *test_helpers.TemplateHandlerMock) {
	handlerMock.ConvertManualEditsFunc = func(dtos.ManualEditConvertRequestDto) (dtos.ManualEditCaptureDto, error) {
		return dtos.ManualEditCaptureDto{}, gofakeit.ErrorValidation()
	}
}

// decodeStoredEdits mirrors the mock's default expansion, for tests that wrap it.
func decodeStoredEdits(request dtos.ManualEditExpandRequestDto) dtos.ManualEditExpansionDto {
	zones, connections := request.State.ToSnapshot()
	return dtos.ManualEditExpansionDto{
		Zones:       zones,
		Connections: connections,
		ManualEdits: request.State.ManualEditSettings,
	}
}
