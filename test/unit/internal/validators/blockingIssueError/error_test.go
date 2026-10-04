package blockingIssueError_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/validators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTopologyIsRetired_ErrorTextIsTheMessageVerbatim(t *testing.T) {
	t.Parallel()
	// Arrange
	rejection := rejectionFor(t, "HubAndSpoke")

	// Act
	actual := rejection.Error()

	// Assert
	assert.Equal(t,
		`topology "Hub" (saved as "HubAndSpoke") has been retired and is no longer supported; `+
			"re-create the template with a supported topology",
		actual)
}

// rejectionFor returns the error the validator rejects the given topology with.
func rejectionFor(t *testing.T, topology config.MapTopology) error {
	t.Helper()
	state := editor_state_model.NewDefaultEditorStateModel()
	state.Topology = topology
	issues := validators.NewEditorStateValidator().Validate(&state)
	require.Len(t, issues, 1)
	require.Error(t, issues[0].Rejection())
	return issues[0].Rejection()
}
