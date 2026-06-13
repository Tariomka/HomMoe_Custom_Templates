package previewRenderer_test

import (
	"image"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

func simpleTemplate(name string) *entities.RmgTemplate {
	return &entities.RmgTemplate{
		Name: name,
		Variants: []entities.Variant{{
			Zones: []entities.Zone{
				{Name: "Spawn-A", MainObjects: []entities.MainObject{{Type: "Spawn", Spawn: "Player1"}}},
				{Name: "Spawn-B", MainObjects: []entities.MainObject{{Type: "Spawn", Spawn: "Player2"}}},
				{Name: "Neutral-C"},
			},
			Connections: []entities.Connection{
				{From: "Spawn-A", To: "Neutral-C", ConnectionType: "Direct"},
				{From: "Neutral-C", To: "Spawn-B", ConnectionType: "Direct"},
			},
		}},
	}
}

// ── RenderPreviewImage ───────────────────────────────────────────────

func TestRenderPreviewImage_ReturnsImageOfRequestedSize(t *testing.T) {
	img := services.RenderPreviewImage(simpleTemplate("T"), config.TopologyDefault, 400)
	if img == nil {
		t.Fatal("nil image")
	}
	if img.Bounds() != image.Rect(0, 0, 400, 400) {
		t.Errorf("bounds = %v", img.Bounds())
	}
}

func TestRenderPreviewImage_EmptyTemplateReturnsBackgroundOnly(t *testing.T) {
	img := services.RenderPreviewImage(&entities.RmgTemplate{}, config.TopologyDefault, 100)
	if img == nil {
		t.Fatal("nil image")
	}
	if img.Bounds().Dx() != 100 {
		t.Errorf("size = %d", img.Bounds().Dx())
	}
}

func TestRenderPreviewImage_PortalConnectionRenders(t *testing.T) {
	tmpl := simpleTemplate("T")
	tmpl.Variants[0].Connections = append(tmpl.Variants[0].Connections,
		entities.Connection{From: "Spawn-A", To: "Spawn-B", ConnectionType: "Portal"},
	)
	img := services.RenderPreviewImage(tmpl, config.TopologyDefault, 300)
	if img == nil {
		t.Fatal("nil image")
	}
}

func TestRenderPreviewImage_CastleBadgeDrawn(t *testing.T) {
	tmpl := simpleTemplate("T")
	// Add an extra city on player A so it has multiple castles and a badge.
	tmpl.Variants[0].Zones[0].MainObjects = append(tmpl.Variants[0].Zones[0].MainObjects,
		entities.MainObject{Type: "City"})
	img := services.RenderPreviewImage(tmpl, config.TopologyDefault, 300)
	if img == nil {
		t.Fatal("nil image")
	}
}

func TestRenderPreviewImage_HubZoneRendered(t *testing.T) {
	tmpl := &entities.RmgTemplate{
		Variants: []entities.Variant{{
			Zones: []entities.Zone{
				{Name: "Hub"},
				{Name: "Spawn-A", MainObjects: []entities.MainObject{{Type: "Spawn", Spawn: "Player1"}}},
			},
			Connections: []entities.Connection{{From: "Hub", To: "Spawn-A", ConnectionType: "Direct"}},
		}},
	}
	img := services.RenderPreviewImage(tmpl, config.TopologyDefault, 300)
	if img == nil {
		t.Fatal("nil image")
	}
}

// ── WritePreviewPNG ──────────────────────────────────────────────────

func TestWritePreviewPNG_WritesFileWithSanitisedName(t *testing.T) {
	dir := t.TempDir()
	path, err := services.WritePreviewPNG(dir, simpleTemplate("My/Preview"), config.TopologyDefault, 200)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, "My_Preview.png")
	if path != want {
		t.Errorf("path = %q, want %q", path, want)
	}
	if info, err := os.Stat(path); err != nil || info.Size() == 0 {
		t.Errorf("file missing/empty: %v %v", info, err)
	}
}

func TestWritePreviewPNG_EmptyNameFallback(t *testing.T) {
	dir := t.TempDir()
	path, err := services.WritePreviewPNG(dir, simpleTemplate(""), config.TopologyDefault, 100)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(path, "Generated_Template.png") {
		t.Errorf("path = %q", path)
	}
}

func TestWritePreviewPNG_CreatesMissingDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "preview")
	if _, err := services.WritePreviewPNG(dir, simpleTemplate("T"), config.TopologyDefault, 100); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("dir not created: %v", err)
	}
}

func TestWritePreviewPNG_MkdirError(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := services.WritePreviewPNG(filepath.Join(blocker, "x"), simpleTemplate("T"), config.TopologyDefault, 100); err == nil {
		t.Error("expected mkdir error")
	}
}
