package template

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/stretchr/testify/assert"
)

// TestRmgTemplate_RoundTripAllExamples decodes every bundled example template,
// re-encodes it, and decodes again to verify the model captures every field.
func TestRmgTemplate_RoundTripAllExamples(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", "..", "..", "data", "ExampleTemplates"))
	assert.NoError(t, err, "resolve example dir")

	entries, err := os.ReadDir(root)
	assert.NoError(t, err, "read example dir")

	count := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".rmg.json") {
			continue
		}
		name := e.Name()
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(root, name)
			raw, err := os.ReadFile(path)
			assert.NoError(t, err, "read file")

			var tpl template.RmgTemplate
			dec := json.NewDecoder(strings.NewReader(string(raw)))
			dec.DisallowUnknownFields()
			err = dec.Decode(&tpl)
			assert.NoError(t, err, func() string {
				offset := dec.InputOffset()
				line := 1 + strings.Count(string(raw[:offset]), "\n")
				return fmt.Sprintf("decode %s near line %d (offset %d): %v", name, line, offset, err)
			}())

			// Re-encode and decode again to confirm the model is self-consistent.
			out, err := json.Marshal(&tpl)
			assert.NoError(t, err, "re-encode")
			var tpl2 template.RmgTemplate
			err = json.Unmarshal(out, &tpl2)
			assert.NoError(t, err, "re-decode")
		})
		count++
	}

	assert.NotEqual(t, count, 0, "no .rmg.json files found")
}
