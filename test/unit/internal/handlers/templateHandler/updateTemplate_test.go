package templateHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWhenUpdatedTemplateIsMissing_ReturnsProvidedTemplateInvalidError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()

	// Act
	_, err := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{})

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrProvidedTemplateInvalid)
}

func TestWhenUpdatedTemplateHasNoVariants_ReturnsProvidedTemplateInvalidError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()

	// Act
	_, err := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{Template: &template_model.Template{}})

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrProvidedTemplateInvalid)
}

func TestWhenTemplateIsUpdated_ReplacesTheFirstVariantsZones(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	zones := []template_model.Zone{{Name: gofakeit.Word()}}
	arrangeUpdateCollaborators(fixture, false)

	// Act
	loadDto, _ := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template: singleVariantTemplate(),
		Zones:    zones,
	})

	// Assert
	assert.Equal(t, zones, loadDto.Template.Variants[0].Zones)
}

func TestWhenTemplateIsUpdated_ReplacesTheFirstVariantsConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	connections := []template_model.Connection{{Name: gofakeit.Word()}}
	arrangeUpdateCollaborators(fixture, false)

	// Act
	loadDto, _ := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    singleVariantTemplate(),
		Connections: connections,
	})

	// Assert
	assert.Equal(t, connections, loadDto.Template.Variants[0].Connections)
}

// The user-added flag is editor-only state with no .rmg.json counterpart; an
// Apply must carry it on the model as-is.
func TestWhenAnAppliedConnectionIsUserAdded_KeepsTheFlagAcrossTheApply(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	connections := []template_model.Connection{{Name: gofakeit.Word(), IsUserAdded: true}}
	arrangeUpdateCollaborators(fixture, false)

	// Act
	loadDto, _ := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    singleVariantTemplate(),
		Connections: connections,
	})

	// Assert
	assert.True(t, loadDto.Template.Variants[0].Connections[0].IsUserAdded)
}

func TestWhenTemplateIsUpdated_LeavesTheSourceTemplateUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	source := singleVariantTemplate()
	originalZones := source.Variants[0].Zones
	arrangeUpdateCollaborators(fixture, false)

	// Act
	_, _ = fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template: source,
		Zones:    []template_model.Zone{{Name: gofakeit.Word()}},
	})

	// Assert
	assert.Equal(t, originalZones, source.Variants[0].Zones)
}

func TestWhenTemplateIsUpdated_NamesTheNamelessConnectionsDirectly(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	zones := []template_model.Zone{{Name: gofakeit.Word()}}
	connections := []template_model.Connection{{Name: gofakeit.Word()}}
	arrangeUpdateCollaborators(fixture, false)

	// Act
	_, _ = fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    singleVariantTemplate(),
		Zones:       zones,
		Connections: connections,
	})

	// Assert
	fixture.zoneEditor.AssertCalled(t, "EnsureConnectionNames", connections)
}

// Naming is the only zone-editor work left on this path: the settings-free
// road rebuild would reinterpret the road policy a second time.
func TestWhenTemplateIsUpdated_DoesNotRunTheSettingsFreeRoadRebuild(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	arrangeUpdateCollaborators(fixture, false)

	// Act
	_, _ = fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    singleVariantTemplate(),
		Zones:       []template_model.Zone{{Name: gofakeit.Word()}},
		Connections: []template_model.Connection{{Name: gofakeit.Word()}},
	})

	// Assert
	fixture.zoneEditor.AssertNotCalled(t, "RebuildZoneConnectionRoads", mock.Anything, mock.Anything)
}

func TestWhenNoEditorStateIsSupplied_NamesTheNamelessConnectionsAnyway(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	connections := []template_model.Connection{{Name: gofakeit.Word()}}
	arrangeUpdateCollaborators(fixture, false)

	// Act
	_, _ = fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    singleVariantTemplate(),
		Connections: connections,
	})

	// Assert
	fixture.zoneEditor.AssertCalled(t, "EnsureConnectionNames", connections)
}

func TestWhenNoEditorStateIsSupplied_KeepsTheExistingMandatoryContent(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	arrangeUpdateCollaborators(fixture, false)

	// Act
	_, _ = fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{Template: singleVariantTemplate()})

	// Assert
	fixture.contentProvider.AssertNotCalled(t, "CreateContentsForZones", mock.Anything, mock.Anything)
}

// There are no generation settings to interpret without the editor state, so
// the road policy - which would otherwise stamp flags and drop road targets it
// cannot resolve - must not run at all.
func TestWhenNoEditorStateIsSupplied_SkipsTheRoadPolicy(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	arrangeUpdateCollaborators(fixture, false)

	// Act
	_, _ = fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    singleVariantTemplate(),
		Zones:       []template_model.Zone{{Name: gofakeit.Word()}},
		Connections: []template_model.Connection{{Name: gofakeit.Word()}},
	})

	// Assert
	fixture.roadPolicy.AssertNotCalled(t, "Reconcile", mock.Anything)
}

func TestWhenNoEditorStateIsSupplied_KeepsTheSuppliedZoneRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	zones := []template_model.Zone{{Name: gofakeit.Word(), Roads: []template_model.Road{customContentRoad()}}}
	expected := append([]template_model.Road(nil), zones[0].Roads...)
	arrangeUpdateCollaborators(fixture, false)

	// Act
	loadDto, _ := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template: singleVariantTemplate(),
		Zones:    zones,
	})

	// Assert
	assert.Equal(t, expected, loadDto.Template.Variants[0].Zones[0].Roads)
}

func TestWhenNoEditorStateIsSupplied_KeepsTheSuppliedConnectionRoadFlags(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	roadless := false
	connections := []template_model.Connection{
		{Name: gofakeit.Word(), ConnectionType: "Direct", Road: &roadless},
		{Name: gofakeit.Word(), ConnectionType: "Direct"},
	}
	arrangeUpdateCollaborators(fixture, false)

	// Act
	loadDto, _ := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    singleVariantTemplate(),
		Connections: connections,
	})

	// Assert
	assert.Equal(t,
		[]*bool{&roadless, nil},
		roadFlagsOf(loadDto.Template.Variants[0].Connections))
}

func TestWhenNoEditorStateIsSupplied_KeepsTheTemplatesOwnMandatoryContent(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	source := singleVariantTemplate()
	source.MandatoryContent = []template_model.MandatoryContent{{Name: gofakeit.Word()}}
	arrangeUpdateCollaborators(fixture, false)

	// Act
	loadDto, _ := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{Template: source})

	// Assert
	assert.Equal(t, source.MandatoryContent, loadDto.Template.MandatoryContent)
}

func TestWhenEditorStateIsSupplied_RebuildsTheMandatoryContentFromTheFinalZones(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	zones := []template_model.Zone{{Name: gofakeit.Word()}}
	expected := []template_model.MandatoryContent{{Name: gofakeit.Word()}}
	configuration := namedConfiguration()
	arrangeUpdateCollaborators(fixture, false)
	fixture.mapper.On("FromEditorState", state).Return(configuration)
	fixture.contentProvider.On("CreateContentsForZones", *configuration, zones).Return(expected)

	// Act
	loadDto, _ := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    singleVariantTemplate(),
		Zones:       zones,
		EditorState: &editor_state_dto.EditorStateDto{EditorState: state},
	})

	// Assert
	assert.Equal(t, expected, loadDto.Template.MandatoryContent)
}

// The policy resolves content roads through the groups the zones end up with,
// so it has to see the content this very Apply just rebuilt.
func TestWhenEditorStateIsSupplied_HandsTheFinalContentToTheRoadPolicy(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	zones := []template_model.Zone{{Name: gofakeit.Word()}}
	rebuilt := []template_model.MandatoryContent{{Name: gofakeit.Word()}}
	configuration := namedConfiguration()
	arrangeUpdateCollaborators(fixture, false)
	fixture.mapper.On("FromEditorState", state).Return(configuration)
	fixture.contentProvider.On("CreateContentsForZones", *configuration, zones).Return(rebuilt)

	// Act
	_, _ = fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    singleVariantTemplate(),
		Zones:       zones,
		EditorState: &editor_state_dto.EditorStateDto{EditorState: state},
	})

	// Assert
	assert.Equal(t, rebuilt, capturedReconciliation(t, fixture).MandatoryContent)
}

func TestWhenEditorStateIsSupplied_HandsTheEffectiveSettingsToTheRoadPolicy(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	configuration := namedConfiguration()
	configuration.GenerateRoads = false
	configuration.SpawnRemoteFootholds = true
	configuration.RemoteFootholdCount = gofakeit.Number(1, 4)
	arrangeUpdateCollaborators(fixture, false)
	fixture.mapper.On("FromEditorState", state).Return(configuration)
	fixture.contentProvider.On("CreateContentsForZones", mock.Anything, mock.Anything).Return(nil)

	// Act
	_, _ = fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    singleVariantTemplate(),
		EditorState: &editor_state_dto.EditorStateDto{EditorState: state},
	})

	// Assert
	request := capturedReconciliation(t, fixture)
	assert.Equal(t,
		[]any{false, true, configuration.RemoteFootholdCount},
		[]any{request.GenerateRoads, request.SpawnRemoteFootholds, request.RemoteFootholdCount})
}

func TestWhenEditorStateIsSupplied_HandsTheFinalGraphToTheRoadPolicy(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	zones := []template_model.Zone{{Name: gofakeit.Word()}}
	connections := []template_model.Connection{{Name: gofakeit.Word()}}
	arrangeUpdateCollaborators(fixture, false)
	fixture.mapper.On("FromEditorState", state).Return(namedConfiguration())
	fixture.contentProvider.On("CreateContentsForZones", mock.Anything, mock.Anything).Return(nil)

	// Act
	_, _ = fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    singleVariantTemplate(),
		Zones:       zones,
		Connections: connections,
		EditorState: &editor_state_dto.EditorStateDto{EditorState: state},
	})

	// Assert
	request := capturedReconciliation(t, fixture)
	assert.Equal(t,
		[]any{zones, connections},
		[]any{request.Zones, request.Connections})
}

// Content first, policy second: the other order would reconcile against the
// content of the previous Apply.
func TestWhenEditorStateIsSupplied_RebuildsTheContentBeforeReconciling(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	var order []string
	fixture.zoneEditor.On("EnsureConnectionNames", mock.Anything).Return()
	fixture.connectionEditor.On("ComputeHasErrors", mock.Anything, mock.Anything).Return(false)
	fixture.mapper.On("FromEditorState", state).Return(namedConfiguration())
	fixture.contentProvider.On("CreateContentsForZones", mock.Anything, mock.Anything).
		Run(func(mock.Arguments) { order = append(order, "content") }).Return(nil)
	fixture.roadPolicy.On("Reconcile", mock.Anything).
		Run(func(mock.Arguments) { order = append(order, "policy") }).Return()

	// Act
	_, _ = fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    singleVariantTemplate(),
		EditorState: &editor_state_dto.EditorStateDto{EditorState: state},
	})

	// Assert
	assert.Equal(t, []string{"content", "policy"}, order)
}

func TestWhenTheUpdatedGraphHasErrors_ReturnsZonesMissingError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	arrangeUpdateCollaborators(fixture, true)

	// Act
	_, err := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{Template: singleVariantTemplate()})

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrZonesMissing)
}

func TestWhenTheUpdatedGraphIsSound_ReturnsNoError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	arrangeUpdateCollaborators(fixture, false)

	// Act
	_, err := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{Template: singleVariantTemplate()})

	// Assert
	assert.NoError(t, err)
}

// arrangeUpdateCollaborators stubs the collaborators every UpdateTemplate call
// reaches, with the requested graph-error verdict.
func arrangeUpdateCollaborators(fixture *templateHandlerFixture, hasErrors bool) {
	fixture.zoneEditor.On("EnsureConnectionNames", mock.Anything).Return()
	fixture.roadPolicy.On("Reconcile", mock.Anything).Return()
	fixture.connectionEditor.On("ComputeHasErrors", mock.Anything, mock.Anything).Return(hasErrors)
}

// capturedReconciliation returns the single request the handler handed the
// road policy.
func capturedReconciliation(
	t *testing.T,
	fixture *templateHandlerFixture) models.RoadReconciliationRequest {
	t.Helper()
	for _, call := range fixture.roadPolicy.Calls {
		if call.Method != "Reconcile" {
			continue
		}
		request, ok := call.Arguments.Get(0).(models.RoadReconciliationRequest)
		require.True(t, ok, "Reconcile was called with an unexpected argument type")
		return request
	}
	t.Fatal("the road policy was never asked to reconcile")
	return models.RoadReconciliationRequest{}
}

// customContentRoad is a hand-authored road whose target the stateless path
// has no content to validate against.
func customContentRoad() template_model.Road {
	return template_model.Road{
		From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
		To:   template_model.TypedRef{Type: "MandatoryContent", Args: []string{gofakeit.Word()}},
	}
}

func roadFlagsOf(connections []template_model.Connection) []*bool {
	flags := make([]*bool, 0, len(connections))
	for _, connection := range connections {
		flags = append(flags, connection.Road)
	}
	return flags
}

func singleVariantTemplate() *template_model.Template {
	return &template_model.Template{
		Name:     gofakeit.Word(),
		Variants: []template_model.Variant{{Zones: []template_model.Zone{{Name: gofakeit.Word()}}}},
	}
}
