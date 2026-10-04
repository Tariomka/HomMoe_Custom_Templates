package validationIssue_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenIssueIsBlocking_RejectionIsUnsupportedTopology(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.Topology = "SharedWeb"
	issue := singleIssue(t, &state)

	// Act
	actual := issue.Rejection()

	// Assert
	assert.ErrorIs(t, actual, common_errors.ErrUnsupportedTopology)
}

func TestWhenIssueIsBlocking_RejectionTextIsTheIssueMessage(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.Topology = "SharedWeb"
	issue := singleIssue(t, &state)

	// Act
	actual := issue.Rejection()

	// Assert
	assert.EqualError(t, actual, issue.Message)
}

func TestWhenIssueIsNotBlocking_HasNoRejection(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_model.NewDefaultEditorStateModel()
	state.Topology = ""
	issue := singleIssue(t, &state)

	// Act
	actual := issue.Rejection()

	// Assert
	assert.NoError(t, actual)
}
