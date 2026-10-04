package validationIssue_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/validators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenIssueRejectsTheState_IsBlocking(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.Topology = "NotARealTopology"
	issue := singleIssue(t, &state)

	// Act
	actual := issue.IsBlocking()

	// Assert
	assert.True(t, actual)
}

func TestWhenIssueHasACorrection_IsNotBlocking(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.Topology = ""
	issue := singleIssue(t, &state)

	// Act
	actual := issue.IsBlocking()

	// Assert
	assert.False(t, actual)
}

// singleIssue validates a state that is expected to carry exactly one issue.
func singleIssue(t *testing.T, state *editor_state_model.EditorState) validators.ValidationIssue {
	t.Helper()
	issues := validators.NewEditorStateValidator().Validate(state)
	require.Len(t, issues, 1)
	return issues[0]
}
