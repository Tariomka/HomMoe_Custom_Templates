package zoneContentManager_test

import (
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

func cfg() *config.GeneratorConfig { return config.NewGeneratorConfig() }

// ── RowsToMandatoryContent ───────────────────────────────────────────

func TestRowsToMandatoryContent_EmptyInputReturnsNil(t *testing.T) {
	if out := services.RowsToMandatoryContent(nil); out != nil {
		t.Errorf("got %v, want nil", out)
	}
}

func TestRowsToMandatoryContent_BlankSidIsSkipped(t *testing.T) {
	rows := []models.ZoneContentRowSave{{Sid: "", Count: 3}}
	if out := services.RowsToMandatoryContent(rows); len(out) != 0 {
		t.Errorf("blank sid should be skipped, got %d items", len(out))
	}
}

func TestRowsToMandatoryContent_ExpandsByCount(t *testing.T) {
	rows := []models.ZoneContentRowSave{{Sid: "a", Count: 3}}
	out := services.RowsToMandatoryContent(rows)
	if len(out) != 3 {
		t.Errorf("got %d items, want 3", len(out))
	}
}

func TestRowsToMandatoryContent_ZeroCountNormalisesToOne(t *testing.T) {
	rows := []models.ZoneContentRowSave{{Sid: "a", Count: 0}}
	out := services.RowsToMandatoryContent(rows)
	if len(out) != 1 {
		t.Errorf("got %d items, want 1 (Normalised should floor at 1)", len(out))
	}
}

func TestRowsToMandatoryContent_SidPopulatesSID(t *testing.T) {
	rows := []models.ZoneContentRowSave{{Sid: "mine_gold", Count: 1, IsMine: true, IsGuarded: true}}
	out := services.RowsToMandatoryContent(rows)
	if out[0].SID != "mine_gold" || !out[0].IsMine || !out[0].IsGuarded {
		t.Errorf("bad item: %+v", out[0])
	}
}

func TestRowsToMandatoryContent_GroupRoutesToIncludeLists(t *testing.T) {
	rows := []models.ZoneContentRowSave{{Sid: "group_x", Count: 1, IsGroup: true}}
	out := services.RowsToMandatoryContent(rows)
	if out[0].SID != "" {
		t.Errorf("group row should leave SID empty, got %q", out[0].SID)
	}
	if len(out[0].IncludeLists) != 1 || out[0].IncludeLists[0] != "group_x" {
		t.Errorf("expected IncludeLists=[group_x], got %v", out[0].IncludeLists)
	}
}

func TestRowsToMandatoryContent_NearCastleAddsRule(t *testing.T) {
	rows := []models.ZoneContentRowSave{{Sid: "x", Count: 1, NearCastle: true}}
	out := services.RowsToMandatoryContent(rows)
	if len(out[0].Rules) != 1 || out[0].Rules[0].Type != "MainObject" {
		t.Errorf("expected one MainObject rule, got %v", out[0].Rules)
	}
}

func TestRowsToMandatoryContent_RoadDistanceNextTo(t *testing.T) {
	out := services.RowsToMandatoryContent([]models.ZoneContentRowSave{{Sid: "x", Count: 1, RoadDistance: "Next To"}})
	if len(out[0].Rules) != 1 || out[0].Rules[0].Type != "Road" {
		t.Errorf("expected Road rule, got %v", out[0].Rules)
	}
}

func TestRowsToMandatoryContent_RoadDistanceNear(t *testing.T) {
	out := services.RowsToMandatoryContent([]models.ZoneContentRowSave{{Sid: "x", Count: 1, RoadDistance: "Near"}})
	if out[0].Rules[0].Type != "Road" {
		t.Errorf("expected Road rule")
	}
}

func TestRowsToMandatoryContent_RoadDistanceMedium(t *testing.T) {
	out := services.RowsToMandatoryContent([]models.ZoneContentRowSave{{Sid: "x", Count: 1, RoadDistance: "Medium"}})
	if out[0].Rules[0].Type != "Road" {
		t.Errorf("expected Road rule")
	}
}

func TestRowsToMandatoryContent_RoadDistanceFar(t *testing.T) {
	out := services.RowsToMandatoryContent([]models.ZoneContentRowSave{{Sid: "x", Count: 1, RoadDistance: "Far"}})
	if out[0].Rules[0].Type != "Road" {
		t.Errorf("expected Road rule")
	}
}

func TestRowsToMandatoryContent_RoadDistanceVeryFar(t *testing.T) {
	out := services.RowsToMandatoryContent([]models.ZoneContentRowSave{{Sid: "x", Count: 1, RoadDistance: "Very Far"}})
	if out[0].Rules[0].Type != "Road" {
		t.Errorf("expected Road rule")
	}
}

func TestRowsToMandatoryContent_RoadDistanceAnyAddsNoRule(t *testing.T) {
	out := services.RowsToMandatoryContent([]models.ZoneContentRowSave{{Sid: "x", Count: 1, RoadDistance: "Any"}})
	if len(out[0].Rules) != 0 {
		t.Errorf("Any distance should add no rule, got %v", out[0].Rules)
	}
}

func TestRowsToMandatoryContent_RoadDistanceUnknownLabelAddsNoRule(t *testing.T) {
	out := services.RowsToMandatoryContent([]models.ZoneContentRowSave{{Sid: "x", Count: 1, RoadDistance: "Whatever"}})
	if len(out[0].Rules) != 0 {
		t.Errorf("unknown distance should add no rule, got %v", out[0].Rules)
	}
}

func TestRowsToMandatoryContent_NearCastleAndRoadDistanceTogether(t *testing.T) {
	out := services.RowsToMandatoryContent([]models.ZoneContentRowSave{{Sid: "x", Count: 1, NearCastle: true, RoadDistance: "Near"}})
	if len(out[0].Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(out[0].Rules))
	}
}

// ── StripNearCastleRules ─────────────────────────────────────────────

func TestStripNearCastleRules_RemovesMainObjectIndexZero(t *testing.T) {
	items := []template.MandatoryContentItem{{
		Rules: []template.PlacementRule{
			{Type: "MainObject", Args: []any{"0"}},
			{Type: "Road", Args: []any{}},
		},
	}}
	stripped := services.StripNearCastleRules(items)
	if len(stripped[0].Rules) != 1 || stripped[0].Rules[0].Type != "Road" {
		t.Errorf("expected only Road rule to remain, got %v", stripped[0].Rules)
	}
}

func TestStripNearCastleRules_KeepsMainObjectIndexOne(t *testing.T) {
	items := []template.MandatoryContentItem{{
		Rules: []template.PlacementRule{
			{Type: "MainObject", Args: []any{"1"}},
		},
	}}
	stripped := services.StripNearCastleRules(items)
	if len(stripped[0].Rules) != 1 {
		t.Errorf("MainObject index 1 should be kept, got %v", stripped[0].Rules)
	}
}

func TestStripNearCastleRules_KeepsMainObjectWithNonStringArg(t *testing.T) {
	items := []template.MandatoryContentItem{{
		Rules: []template.PlacementRule{
			{Type: "MainObject", Args: []any{0}}, // int, not string
		},
	}}
	stripped := services.StripNearCastleRules(items)
	if len(stripped[0].Rules) != 1 {
		t.Errorf("non-string arg should be kept, got %v", stripped[0].Rules)
	}
}

func TestStripNearCastleRules_KeepsMainObjectWithNoArgs(t *testing.T) {
	items := []template.MandatoryContentItem{{
		Rules: []template.PlacementRule{
			{Type: "MainObject", Args: nil},
		},
	}}
	stripped := services.StripNearCastleRules(items)
	if len(stripped[0].Rules) != 1 {
		t.Errorf("MainObject with no args should be kept, got %v", stripped[0].Rules)
	}
}

func TestStripNearCastleRules_EmptyRulesShortCircuit(t *testing.T) {
	items := []template.MandatoryContentItem{{Rules: nil}}
	stripped := services.StripNearCastleRules(items)
	if stripped[0].Rules != nil {
		t.Errorf("expected nil rules preserved, got %v", stripped[0].Rules)
	}
}

// ── BuildPlayerZoneMandatoryContent ──────────────────────────────────

func TestBuildPlayerZoneMandatoryContent_FootholdAddedWhenEnabled(t *testing.T) {
	s := cfg()
	s.SpawnRemoteFootholds = true
	out := services.BuildPlayerZoneMandatoryContent(s)
	if len(out) == 0 || out[0].Name != "name_remote_foothold_1" {
		t.Errorf("expected foothold prepended, got %+v", out)
	}
}

func TestBuildPlayerZoneMandatoryContent_NoFootholdWhenDisabled(t *testing.T) {
	s := cfg()
	s.SpawnRemoteFootholds = false
	out := services.BuildPlayerZoneMandatoryContent(s)
	for _, c := range out {
		if c.Name == "name_remote_foothold_1" {
			t.Error("foothold present despite disabled")
		}
	}
}

func TestBuildPlayerZoneMandatoryContent_AppendsUserRows(t *testing.T) {
	s := cfg()
	s.SpawnRemoteFootholds = false
	s.PlayerZoneMandatoryContent = []template.MandatoryContentItem{{SID: "u1"}, {SID: "u2"}}
	out := services.BuildPlayerZoneMandatoryContent(s)
	if len(out) != 2 {
		t.Errorf("got %d, want 2", len(out))
	}
}

// ── BuildLowNeutralMandatoryContent ──────────────────────────────────

func TestBuildLowNeutralMandatoryContent_CastleCountPreservesNearCastleRule(t *testing.T) {
	s := cfg()
	s.SpawnRemoteFootholds = false
	s.LowNeutralMandatoryContent = []template.MandatoryContentItem{{
		SID:   "x",
		Rules: []template.PlacementRule{{Type: "MainObject", Args: []any{"0"}}},
	}}
	out := services.BuildLowNeutralMandatoryContent(s, 1)
	if len(out[0].Rules) != 1 {
		t.Errorf("with castle, MainObject[0] rule should be kept")
	}
}

func TestBuildLowNeutralMandatoryContent_NoCastleStripsNearCastleRule(t *testing.T) {
	s := cfg()
	s.SpawnRemoteFootholds = false
	s.LowNeutralMandatoryContent = []template.MandatoryContentItem{{
		SID:   "x",
		Rules: []template.PlacementRule{{Type: "MainObject", Args: []any{"0"}}},
	}}
	out := services.BuildLowNeutralMandatoryContent(s, 0)
	if len(out[0].Rules) != 0 {
		t.Errorf("without castle, MainObject[0] rule should be stripped")
	}
}

func TestBuildLowNeutralMandatoryContent_FootholdAdded(t *testing.T) {
	s := cfg()
	s.SpawnRemoteFootholds = true
	out := services.BuildLowNeutralMandatoryContent(s, 1)
	if len(out) == 0 || out[0].Name != "name_remote_foothold_1" {
		t.Errorf("foothold not prepended: %+v", out)
	}
}

// ── BuildMediumNeutralMandatoryContent ───────────────────────────────

func TestBuildMediumNeutralMandatoryContent_StripsWhenNoCastle(t *testing.T) {
	s := cfg()
	s.SpawnRemoteFootholds = false
	s.MediumNeutralMandatoryContent = []template.MandatoryContentItem{{
		Rules: []template.PlacementRule{{Type: "MainObject", Args: []any{"0"}}},
	}}
	out := services.BuildMediumNeutralMandatoryContent(s, 0)
	if len(out[0].Rules) != 0 {
		t.Errorf("expected rule stripped")
	}
}

func TestBuildMediumNeutralMandatoryContent_KeepsWhenCastle(t *testing.T) {
	s := cfg()
	s.SpawnRemoteFootholds = false
	s.MediumNeutralMandatoryContent = []template.MandatoryContentItem{{
		SID:   "x",
		Rules: []template.PlacementRule{{Type: "MainObject", Args: []any{"0"}}},
	}}
	out := services.BuildMediumNeutralMandatoryContent(s, 2)
	if len(out[0].Rules) != 1 {
		t.Errorf("expected rule kept")
	}
}

// ── BuildHighNeutralMandatoryContent ─────────────────────────────────

func TestBuildHighNeutralMandatoryContent_StripsWhenNoCastle(t *testing.T) {
	s := cfg()
	s.SpawnRemoteFootholds = false
	s.HighNeutralMandatoryContent = []template.MandatoryContentItem{{
		Rules: []template.PlacementRule{{Type: "MainObject", Args: []any{"0"}}},
	}}
	out := services.BuildHighNeutralMandatoryContent(s, 0)
	if len(out[0].Rules) != 0 {
		t.Errorf("expected rule stripped")
	}
}

func TestBuildHighNeutralMandatoryContent_KeepsWhenCastle(t *testing.T) {
	s := cfg()
	s.SpawnRemoteFootholds = false
	s.HighNeutralMandatoryContent = []template.MandatoryContentItem{{
		SID:   "x",
		Rules: []template.PlacementRule{{Type: "MainObject", Args: []any{"0"}}},
	}}
	out := services.BuildHighNeutralMandatoryContent(s, 1)
	if len(out[0].Rules) != 1 {
		t.Errorf("expected rule kept")
	}
}

// ── BuildHubZoneMandatoryContent ─────────────────────────────────────

func TestBuildHubZoneMandatoryContent_StripsWhenNoCastle(t *testing.T) {
	s := cfg()
	s.SpawnRemoteFootholds = false
	s.HubZoneMandatoryContent = []template.MandatoryContentItem{{
		Rules: []template.PlacementRule{{Type: "MainObject", Args: []any{"0"}}},
	}}
	out := services.BuildHubZoneMandatoryContent(s, 0)
	if len(out[0].Rules) != 0 {
		t.Errorf("expected rule stripped")
	}
}

func TestBuildHubZoneMandatoryContent_KeepsWhenCastle(t *testing.T) {
	s := cfg()
	s.SpawnRemoteFootholds = false
	s.HubZoneMandatoryContent = []template.MandatoryContentItem{{
		SID:   "x",
		Rules: []template.PlacementRule{{Type: "MainObject", Args: []any{"0"}}},
	}}
	out := services.BuildHubZoneMandatoryContent(s, 1)
	if len(out[0].Rules) != 1 {
		t.Errorf("expected rule kept")
	}
}

// ── BuildAllContentCountLimits ───────────────────────────────────────

func TestBuildAllContentCountLimits_DefaultGroupCount(t *testing.T) {
	limits := services.BuildAllContentCountLimits(cfg())
	if len(limits) != 17 {
		t.Errorf("expected 17 limit groups, got %d", len(limits))
	}
}

func TestBuildAllContentCountLimits_UserCountLiftsCap(t *testing.T) {
	s := cfg()
	for range 6 {
		s.PlayerZoneMandatoryContent = append(s.PlayerZoneMandatoryContent,
			template.MandatoryContentItem{SID: constants.ContentIds.PandoraBox.Sid})
	}
	limits := services.BuildAllContentCountLimits(s)
	for _, g := range limits {
		for _, l := range g.Limits {
			if strings.EqualFold(l.SID, constants.ContentIds.PandoraBox.Sid) && l.MaxCount < 6 {
				t.Errorf("pandora cap not lifted: %d", l.MaxCount)
				return
			}
		}
	}
}

func TestBuildAllContentCountLimits_UserCountBelowCapKeepsDefault(t *testing.T) {
	s := cfg()
	// One pandora — default cap is 4, should stay at 4.
	s.PlayerZoneMandatoryContent = []template.MandatoryContentItem{{SID: constants.ContentIds.PandoraBox.Sid}}
	limits := services.BuildAllContentCountLimits(s)
	for _, g := range limits {
		if g.Name != "content_limits_side" {
			continue
		}
		for _, l := range g.Limits {
			if strings.EqualFold(l.SID, constants.ContentIds.PandoraBox.Sid) {
				if l.MaxCount != 4 {
					t.Errorf("pandora cap = %d, want 4", l.MaxCount)
				}
				return
			}
		}
	}
}

func TestBuildAllContentCountLimits_BaseGroupsPresent(t *testing.T) {
	limits := services.BuildAllContentCountLimits(cfg())
	names := map[string]bool{}
	for _, g := range limits {
		names[g.Name] = true
	}
	for _, want := range []string{"content_limits_side", "content_limits_side_0_0", "content_limits_side_1_2", "content_limits_side_5_6"} {
		if !names[want] {
			t.Errorf("missing group %q", want)
		}
	}
}
