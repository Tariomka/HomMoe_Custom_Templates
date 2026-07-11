package integration_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/stretchr/testify/require"
)

// TestRmgTemplate_RoundTripAllExamples decodes every bundled example template,
// re-encodes it, and decodes again to verify the model captures every field.
func TestRmgTemplate_RoundTripAllExamples(t *testing.T) {
	t.Parallel()
	root, err := filepath.Abs(filepath.Join("..", "..", "data", "ExampleTemplates"))
	require.NoError(t, err, "resolve example dir")

	entries, err := os.ReadDir(root)
	require.NoError(t, err, "read example dir: "+root)

	count := 0
	for _, entity := range entries {
		if entity.IsDir() || !strings.HasSuffix(entity.Name(), ".rmg.json") {
			continue
		}
		name := entity.Name()
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(root, name)
			raw, err := os.ReadFile(path)
			require.NoError(t, err, "read file: "+path)

			var tpl entities.RmgTemplate
			dec := json.NewDecoder(strings.NewReader(string(raw)))
			dec.DisallowUnknownFields()
			require.NoError(t, dec.Decode(&tpl), func() string {
				offset := dec.InputOffset()
				line := 1 + strings.Count(string(raw[:offset]), "\n")
				return fmt.Sprintf("decode %s near line %d (offset %d): %v", name, line, offset, err)
			}())

			// Re-encode and decode again to confirm the model is self-consistent.
			out, err := json.Marshal(&tpl)
			require.NoError(t, err, "re-encode: "+path)
			var tpl2 entities.RmgTemplate
			err = json.Unmarshal(out, &tpl2)
			require.NoError(t, err, "re-decode: "+path)
		})
		count++
	}

	require.NotEqual(t, 0, count, "no .rmg.json files found")
}
