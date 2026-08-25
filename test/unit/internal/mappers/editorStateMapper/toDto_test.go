package editorStateMapper_test

import (
	"reflect"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenACompleteModelIsMapped_MapsEveryTransferFieldWithoutSchemaVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	model := test_helpers.NewAllFieldsEditorStateModel()
	mapper := mappers.NewEditorStateMapper()

	// Act
	dto := mapper.ToDto(model)

	// Assert
	assert.Equal(t, map[string]any{
		"TemplateName":  dto.TemplateName,
		"MapSize":       dto.MapSize,
		"Topology":      dto.Topology,
		"ManualZones":   dto.ManualZones,
		"SchemaVersion": reflect.ValueOf(dto).FieldByName("SchemaVersion").IsValid(),
	}, map[string]any{
		"TemplateName":  model.TemplateName,
		"MapSize":       model.MapSize,
		"Topology":      model.Topology,
		"ManualZones":   model.ManualZones,
		"SchemaVersion": false,
	})
}
