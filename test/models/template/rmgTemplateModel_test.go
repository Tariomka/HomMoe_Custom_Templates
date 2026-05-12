package template

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

// TestRmgTemplate_RoundTripAllExamples decodes every bundled example template,
// re-encodes it, and decodes again to verify the model captures every field.
func TestRmgTemplate_RoundTripAllExamples(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", "..", "..", "data", "ExampleTemplates"))
	if err != nil {
		t.Fatalf("resolve example dir: %v", err)
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("read example dir %s: %v", root, err)
	}

	count := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".rmg.json") {
			continue
		}
		name := e.Name()
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(root, name)
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read: %v", err)
			}

			var tpl template.RmgTemplateModel
			dec := json.NewDecoder(strings.NewReader(string(raw)))
			dec.DisallowUnknownFields()
			if err := dec.Decode(&tpl); err != nil {
				offset := dec.InputOffset()
				line := 1 + strings.Count(string(raw[:offset]), "\n")
				t.Fatalf("decode %s near line %d (offset %d): %v", name, line, offset, err)
			}

			// Re-encode and decode again to confirm the model is self-consistent.
			out, err := json.Marshal(&tpl)
			if err != nil {
				t.Fatalf("re-encode: %v", err)
			}
			var tpl2 template.RmgTemplateModel
			if err := json.Unmarshal(out, &tpl2); err != nil {
				t.Fatalf("re-decode: %v", err)
			}
		})
		count++
	}

	if count == 0 {
		t.Fatalf("no .rmg.json files found under %s", root)
	}
}
