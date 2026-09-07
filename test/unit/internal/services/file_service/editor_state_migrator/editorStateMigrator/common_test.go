package editorStateMigrator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service/editor_state_migrator"
	"github.com/stretchr/testify/require"
)

func newMigrator() editor_state_migrator.IEditorStateMigrator {
	return editor_state_migrator.NewEditorStateMigrator(
		repositories.NewEditorStateRepository(),
		repositories.NewLegacyEditorStateRepository())
}

func writeStateFile(t *testing.T, body string) string {
	t.Helper()

	statePath := filepath.Join(t.TempDir(), "state.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte(body), 0o644))

	return statePath
}
