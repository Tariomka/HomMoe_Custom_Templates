package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenNeedsDerivationIsAsked_AnswersFromEditsAndGeneratorVersion(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name             string
		hasEdits         bool
		generatorVersion int
		expected         bool
	}{
		{"EditsWithoutGeneratorVersion_NeedDerivation", true, 0, true},
		{"EditsWithGeneratorVersion_DoNotNeedDerivation", true, 1, false},
		{"NoEditsWithoutGeneratorVersion_DoNotNeedDerivation", false, 0, false},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			state := editor_state_model.NewDefaultEditorStateModel()
			state.GeneratorVersion = testCase.generatorVersion
			if testCase.hasEdits {
				state.ManualZones = editor_state_model.ManualEditSettingsFromSnapshot(
					[]template_model.Zone{{Name: "Neutral-A"}}, nil).ManualZones
			}

			// Act
			needsDerivation := state.NeedsDerivation()

			// Assert
			assert.Equal(t, testCase.expected, needsDerivation)
		})
	}
}
