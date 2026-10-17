package v2EditorStateRepository_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state/editor_state_v2"
	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
	"github.com/stretchr/testify/require"
)

func newRepository() repositories.IFileRepository[editor_state_v2.EditorState] {
	return repositories.NewV2EditorStateRepository()
}

func writeStateFile(t *testing.T, body string) string {
	t.Helper()

	statePath := filepath.Join(t.TempDir(), "state.gen.json")
	require.NoError(t, os.WriteFile(statePath, []byte(body), 0o644))

	return statePath
}
