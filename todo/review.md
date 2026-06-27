# Senior Engineering Code Review — HomMoe Custom Templates

**Scope:** Whole repository (~27,100 LOC, 213 source + 20 test Go files).
**Stack:** Go 1.25.8 single module, Gio v0.10.0 immediate-mode GUI.
**Date:** 2026-06-28
**Reviewer perspective:** senior software engineer — correctness, architecture, performance, testability, CI/CD.

> Note on read-only areas: `data/`, `internal/entities/template/`, and `internal/registry/` are authoritative game-data/schema per [.github/AGENTS.md](.github/AGENTS.md). Nothing in this review proposes editing their *contents* — only how other packages depend on them.

---

## 0. Executive Summary

The project is in good shape for a hobby/desktop tool: the generation pipeline is layered (GUI → handlers → services → providers), the topology system is cleanly extensible, and there is a solid block of generation tests. The main weaknesses are:

| Theme | Severity | Headline |
|-------|----------|----------|
| Bugs | **High** | A nil-map panic in `linq.ToMap`, an unchecked type assertion that can panic on a corrupt Steam VDF, and hard-coded `C:` Windows paths. |
| Architecture | **High** | `internal/services` imports `app/gui` (dependency inversion / layering violation, already flagged with a `TODO`). |
| Tests | **Medium** | Handlers, mappers, file I/O, helpers, and the GUI state machine are essentially untested (~0–5%). A stale `oldTests/` tree shadows newer tests. |
| CI/CD | **Medium** | No static analysis, no `go vet`, no lint, no race detector, no coverage, no vuln scan; release lacks checksums/provenance. |
| Duplication | **Medium** | ~50–60 near-identical labeled-slider/checkbox rows; a 1,364-line golden fixture; duplicated package-global vars. |
| Perf | **Low/Medium** | Per-frame allocations and `reflect.DeepEqual` in the Gio Layout hot path. |

A prioritized action plan is in [§8](#8-prioritized-action-plan).

---

## 1. Bugs (with fixes)

### 1.1 — `linq.ToMap()` writes to a nil map → guaranteed panic  **[High]**
[internal/helpers/linq/map.go](internal/helpers/linq/map.go)

```go
func (this QueryMap[TKey, TValue]) ToMap() map[TKey]TValue {
	var result map[TKey]TValue          // nil map
	this.Iterate(func(key TKey, value TValue) bool {
		result[key] = value             // panic: assignment to entry in nil map
		return true
	})
	return result
}
```

Writing to a `nil` map panics at runtime. It is currently *dead code* (no caller — verified by searching for `.ToMap()`), which is the only reason it has not blown up. This is a latent landmine: the first person to use this helper gets a crash.

**Fix:**
```go
func (this QueryMap[TKey, TValue]) ToMap() map[TKey]TValue {
	result := make(map[TKey]TValue)
	this.Iterate(func(key TKey, value TValue) bool {
		result[key] = value
		return true
	})
	return result
}
```
Also note `Query.ToSlice()` returns `nil` for an empty query (acceptable) — keep that consistent and add a unit test that exercises both empty and non-empty cases.

### 1.2 — Unchecked type assertion can panic on a malformed Steam VDF  **[High]**
[internal/helpers/io.go](internal/helpers/io.go#L93)

```go
func getBasePath(vdfContent map[string]any) string {
	for _, data := range vdfContent["libraryfolders"].(map[string]any) { // panics if key missing/!map
```

If `libraryfolders.vdf` is corrupt, partially written, or from a future Steam format where the top key differs/changes type, `vdfContent["libraryfolders"]` is `nil` (or another type) and the type assertion panics — taking the whole app down at startup (it runs from `NewUIState`).

**Fix:** use the comma-ok form and fail gracefully:
```go
libs, ok := vdfContent["libraryfolders"].(map[string]any)
if !ok {
	return ""
}
for _, data := range libs {
	...
}
```

### 1.3 — Hard-coded `C:` drive and `USERNAME` break cross-platform/edge installs  **[High]**
[internal/helpers/io.go](internal/helpers/io.go#L11-L24)

```go
windowsSteamPath = filepath.Join("C:", "Program Files (x86)", "Steam")
windowsUserPath  = filepath.Join("C:", "Users", os.Getenv("USERNAME"))
```

Problems:
- Windows installed on `D:`/`E:`, or Steam installed to a non-default path → never found.
- Profiles redirected (domain machines, OneDrive Known Folder Move) or where `%USERNAME%` ≠ profile folder name → wrong path.
- These are **package-level `var`s evaluated at init**, so they cannot be overridden in a test and are evaluated even on Linux.

The AGENTS guide explicitly requires `os.UserConfigDir`/`os.UserHomeDir` ([.github/AGENTS.md](.github/AGENTS.md) §2.2), which this violates.

**Fix:**
- Replace `windowsUserPath` with `os.UserHomeDir()` (returns the real profile dir; respects `%USERPROFILE%`).
- For the Steam path, read the registry on Windows (`HKCU\Software\Valve\Steam\SteamPath`, fallback `HKLM\...\WOW6432Node\Valve\Steam`) before falling back to the hard-coded default. At minimum, probe both `Program Files (x86)` and the install drive from `%ProgramFiles(x86)%`.
- Move these from package-level `var` into a function so they are evaluated lazily and are testable.

### 1.4 — `resolveGlob` / `os.Stat` error handling only covers “not exist”  **[Medium]**
[internal/helpers/io.go](internal/helpers/io.go#L48-L60), [internal/helpers/io.go](internal/helpers/io.go#L122)

- `FindOldenEraTemplatesDir`: `if _, err := os.Stat(directory); os.IsNotExist(err) { return "", err }` — a *permission* error (non-`IsNotExist`) falls through and the directory is returned as if valid.
- `resolveGlob` returns `("", nil)` when there are zero matches, and `("", err)` at the bottom where `err` may be `nil`. The `("", nil)` shape is then interpreted by the caller ([app/gui/drivers/state.go](app/gui/drivers/state.go#L88-L100)) as “not found, non-error” — acceptable, but the mixed semantics are fragile.

**Fix:** distinguish “not found” from real errors explicitly (e.g. return a sentinel `ErrTemplatesDirNotFound`), and wrap underlying errors with `%w` so callers/tests can assert with `errors.Is`.

### 1.5 — `NewUIState` can silently end up with an empty output path  **[Low]**
[app/gui/drivers/state.go](app/gui/drivers/state.go#L88-L101)

If `FindOldenEraTemplatesDir` returns `""` **and** `os.Getwd()` also errors, `templateDir` stays `""` and is set as the output path with no clear feedback. Set a definite fallback (e.g. `"."` or `os.TempDir()`) and surface a status message.

### 1.6 — `Exit()` calls `os.Exit(0)` directly and bypasses the unsaved-changes guard  **[Low]**
[app/gui/drivers/state.go](app/gui/drivers/state.go#L211-L219)

The unsaved-changes check is commented out and `os.Exit(0)` skips any deferred cleanup and makes the path untestable. Either restore the guard or route exit through the Gio window close so the app can flush state. Track `unsaved` properly (it is currently set in few places).

### 1.7 — Unknown victory condition silently coerced to the first entry  **[Low]**
[app/gui/constants/victoryConditions.go](app/gui/constants/victoryConditions.go#L50)

```go
return VictoryConditions[0] // TODO: probably should return empty Victory... suck it
```

Masks bad/loaded data by returning a valid-looking default. Prefer returning `(VictoryCondition, bool)` or a zero value the caller can detect, and log/raise a status when an unknown condition is encountered.

### 1.8 — `GUIHandler.UpdateTemplate` mutates the caller’s template in place  **[Low]**
[internal/handlers/guiHandler.go](internal/handlers/guiHandler.go#L44-L70)

It assigns into `templateDto.Template.Variants[0]` and rebuilds `MandatoryContent` on the passed pointer, so the input and output alias. The adjacent `// TODO: might not be needed` on `RebuildZoneConnectionRoads` signals uncertainty about the road-rebuild. Document the in-place contract (or operate on a copy) and resolve the TODO with a test that pins the road-rebuild behavior.

---

## 2. Performance / Optimization

These are mostly Gio **Layout hot-path** issues — `Layout` runs every frame, so allocations and recomputation there add up.

### 2.1 — `reflect.DeepEqual` on the whole state every frame  **[Medium]**
[internal/dtos/editorStateDto.go](internal/dtos/editorStateDto.go#L108-L115), called from [app/gui/drivers/state.go](app/gui/drivers/state.go#L268-L285)

`EqualsIgnoringManualEdits` copies the entire `EditorStateDto` (incl. slice headers) and runs `reflect.DeepEqual` on every frame via `AutoRegenerate`. Correct, but reflection allocates and walks every content-row slice each frame.

**Options:** generate/maintain a hand-written field comparison (the struct is large but flat), or compute a cheap hash/fingerprint of the change-relevant fields and compare that. At minimum, short-circuit on the scalar fields (which already exist as `LayoutDefiningOptionsChanged`/`zoneCountOptionsChanged`) before falling back to a deep compare of the content-row slices only.

### 2.2 — Per-frame slice allocations in Layout  **[Low/Medium]**
- [app/gui/editor/window.go](app/gui/editor/window.go#L70-L78) — `getTabsWidget` allocates `make([]layout.FlexChild, 0)` and rebuilds the tab children **every frame**, even though tabs are fixed. Build once in `NewWindow` (or cache and only rebuild when selection changes).
- [app/gui/components/dropdownSelector.go](app/gui/components/dropdownSelector.go) — `children = append(...)` inside `Layout` reallocates for large dropdowns each frame.
- [app/gui/dialogs/ruleDialog.go](app/gui/dialogs/ruleDialog.go) — `buildContentWidgets()` rebuilds the whole widget tree (and `rows` slice) each frame.
- [app/gui/dialogs/zoneEditorDialog.go](app/gui/dialogs/zoneEditorDialog.go) — `recomputeGeometry()` rebuilds edge geometry (`groups`/`order`/`edges`) every frame; it only needs to recompute when zones/connections change. Add a dirty flag.

### 2.3 — Repeated config mapping per frame  **[Low]**
[app/gui/drivers/state.go](app/gui/drivers/state.go#L118-L122) — `GetGeneratorConfig()` calls `mapper.FromEditorState`, which constructs a fresh `providers.NewMandatoryContentProvider()` and re-parses bonuses/banned lists ([internal/mappers/generatorConfigMapper.go](internal/mappers/generatorConfigMapper.go#L22-L48)). When called on hot paths this re-allocates each time. Hold the provider on the mapper, and only remap when the state actually changed.

### 2.4 — Micro
- [internal/helpers/math.go](internal/helpers/math.go#L37-L46) `BoolToInt` can be the idiomatic one-liner; the “compiler optimized” comment is misleading. Minor.

---

## 3. Architecture / Clean Architecture / Separation of Concerns

### 3.1 — Layering violation: `internal/services` imports `app/gui`  **[High]**
[internal/services/content_rules/variantMappingManager.go](internal/services/content_rules/variantMappingManager.go#L6)

```go
import (
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants" // TODO: This should not exist
)
```

The dependency arrow points *upward* (domain/service layer → UI layer). This:
- inverts the intended layering (GUI → handlers → services),
- makes `internal/services/content_rules` impossible to test without dragging in Gio/UI constants,
- risks an import cycle the moment `app/gui/constants` needs anything from services.

`ContentIds` is **domain data** (content SIDs), not UI display data. It belongs under `internal/` (e.g. `internal/constants` or `internal/registry`), with `app/gui/constants` re-exporting/decorating it for display. The repo memory already documents that `app/gui/constants` mixes “display-name catalogs” with SID references — the SID half should live below the UI line.

**Fix:** move `ContentIds` (and any other pure-domain catalog imported by `internal/`) down to `internal/constants` or `internal/registry`; have the GUI depend on that. Add an import-boundary lint (see [§7.4](#74-enforce-architecture-with-a-linter)) to prevent regressions.

### 3.2 — God objects / files doing too much  **[Medium]**
| File | LOC | Concern |
|------|-----|---------|
| [app/gui/dialogs/zoneEditorDialog.go](app/gui/dialogs/zoneEditorDialog.go) | 1,045 | canvas rendering + pointer handling + property edit/writeback + parsing/validation. Split into `ZoneEditorDialog` (UI), `ZoneEditorModel` (state/logic), `ZoneEditorRenderer` (canvas). The `sync*/writeback*` functions are pure logic and belong in a service. |
| [internal/services/previewLayout.go](internal/services/previewLayout.go) | 1,024 | many layout strategies (manual/rings/fixed/scatter/ring/hub) in one file. Split per strategy behind a small interface. |
| [app/gui/panels/layoutPanel.go](app/gui/panels/layoutPanel.go) | 451 | ~50 widget fields + UI build + state sync. Extract a `LayoutPanelState` and a widget-factory. |
| [app/gui/dialogs/pickerDialog.go](app/gui/dialogs/pickerDialog.go) | 414 | generic picker + spell-specific + item-specific behavior. Strategy object per picker type. |
| [test/test_helpers/defaultTemplate.go](test/test_helpers/defaultTemplate.go) | 1,364 | see [§5.2](#52-the-1364-line-golden-fixture). |

### 3.3 — Non-idiomatic `this` receivers project-wide  **[Low, but pervasive]**
Every method uses `func (this *T)`. Go style (and `revive`/`golint`) flags `this`/`self` as generic receiver names. This is a deliberate house style, so treat it as optional — but if you ever want to adopt `golangci-lint` with `revive`, either disable that rule explicitly or do a one-time rename. Document the decision in [.github/AGENTS.md](.github/AGENTS.md).

### 3.4 — Duplicated package-global registry lookups  **[Low]**
`winConditions = registry.GetWinningConditionValues()` is declared as a package global in both [internal/dtos/editorStateDto.go](internal/dtos/editorStateDto.go#L11-L14) and [internal/mappers/generatorConfigMapper.go](internal/mappers/generatorConfigMapper.go#L10-L12). Fine functionally, but it’s a smell that the same derived constant is recomputed in multiple layers. Consider a single shared accessor.

### 3.5 — Doc drift between AGENTS.md and code  **[Low]**
[.github/AGENTS.md](.github/AGENTS.md#L8) says Gio `v0.9.0` and references `app/gui/program.go` as the entry, while [go.mod](go.mod) pins `gioui.org v0.10.0` and `main.go` calls `application.StartApplication()`. Keep the operating guide in sync.

---

## 4. Duplicate Code / Refactoring Opportunities

### 4.1 — Labeled slider/checkbox row construction (~50–60 occurrences)  **[Medium]**
[app/gui/panels/layoutPanel.go](app/gui/panels/layoutPanel.go#L229-L343), [app/gui/panels/generalPanel.go](app/gui/panels/generalPanel.go#L191-L312), and others repeat:

```go
widgets.NewLabeledRowWidget(theme, "Label", width,
    widgets.NewLabeledSliderWidget(theme, &slider, utils.RoundedRangeString(...)))
```

Extract `NewLabeledSliderRow(theme, label, width, &slider, min, max)` and `NewLabeledCheckboxRow(...)` helpers. This removes dozens of lines and makes the panels read like a form schema.

### 4.2 — `oldTests/` shadows newer tests  **[Medium]**
[test/services/oldTests/](test/services/oldTests/) (`services_test.go` 1,074 lines, plus `previewLayout/`, `previewRenderer/`, `templateWriter/`, `zoneContentManager/`, `settingsFileLoader/`) use raw `t.Errorf` and overlap with the newer testify-based suites under `test/services/template_generator/…`. Audit for any *unique* coverage, port what’s valuable, and delete the rest. Maintaining two parallel suites is a liability.

### 4.3 — Geometry-driven topologies already deduped (good)  ✅
`circlesTopology`, `squareTopology`, `geometricTopology`, `crossTopology` reuse `RandomTopologyService`’s connection/zone creation via shared `geometryHelpers.go`. This is the right pattern — keep it; just make sure new topologies follow it rather than copy-pasting.

---

## 5. Test Coverage Gaps

Rough coverage today: **~30–40% of non-GUI logic**, **~5% of GUI**. The generation/topology core is well covered; the seams around it are not.

### 5.1 — Untested critical paths  **[High]**
- [internal/handlers/guiHandler.go](internal/handlers/guiHandler.go) — all 6 entry points (`GenerateTemplate`/`UpdateTemplate`/`SaveTemplate`/`LoadState`/`SaveState`) including their error branches (`ErrNoTemplateName`, `ErrProvidedTemplateInvalid`, `ErrZonesMissing`, `ErrNothingToSave`, `ErrNoOutputPath`). **0 tests.**
- [internal/mappers/generatorConfigMapper.go](internal/mappers/generatorConfigMapper.go) — 40+ field DTO→config transform, bonus/banned-list parsing, `CityHold` win-condition coupling. **0 tests.**
- [internal/helpers/io.go](internal/helpers/io.go) — Steam VDF parsing, glob resolution, cross-platform path logic. **0 tests** (and currently panics on bad input — see §1.2). Refactor for injectable FS/env, then table-test malformed/missing/permission cases.
- [internal/helpers/](internal/helpers/) — `string.go` (`SanitizeFilename`), `slice.go` (`GetUniqueElements`), `math.go` (`Clamp`/`Scale`/`RoundWithPrecision`), `linq/`. Pure functions, trivial to table-test; `linq.ToMap` would have caught §1.1.
- [internal/services/templateWriter.go](internal/services/templateWriter.go) and [internal/services/previewRenderer.go](internal/services/previewRenderer.go) — file/PNG I/O (use `t.TempDir()`).
- [internal/dtos/editorStateDto.go](internal/dtos/editorStateDto.go) — JSON **round-trip** tests (marshal→unmarshal→equal) to catch tag typos and default drift; plus `EqualsIgnoringManualEdits`, `LayoutDefiningOptionsChanged`, `zoneCountOptionsChanged`.
- [app/gui/drivers/state.go](app/gui/drivers/state.go) — the debounce/auto-regen state machine (`AutoRegenerate`, `performAutoRegen`, manual-edit persist/restore). Logic is decoupled enough to unit-test with a fake clock without rendering.

### 5.2 — The 1,364-line golden fixture  **[Medium]**
[test/test_helpers/defaultTemplate.go](test/test_helpers/defaultTemplate.go) hand-builds an entire `RmgTemplate`. Any schema/field change breaks compilation and produces unreadable diffs; comments like “Most likely this is random” suggest it is partly stale. Prefer:
- a `testdata/*.golden.json` file compared with a normalized marshal (regenerated behind a `-update` flag), **or**
- a builder + targeted assertions on the fields that matter, instead of one monolithic snapshot.

### 5.3 — Assertion depth / table-driven style  **[Low]**
Many generator tests assert a single field (e.g. only `Name`). Convert the topology/win-condition matrix to table-driven tests and assert structural invariants (zone counts, connectivity, no player↔player edges where forbidden, roads present when enabled). The connection-editor suite (`T015`/`T019`… IDs) is a good model — extend that rigor.

### 5.4 — Cross-platform test concerns  **[Medium]**
`io.go` is OS-sensitive but untested on either OS; CI runs Linux-only. Inject the filesystem/registry/env so the Windows path logic is exercised on Linux CI, and add a Windows runner to the matrix (see [§7.1](#71-add-a-lint--static-analysis-job)).

---

## 6. CI/CD Pipeline

Current state: [.github/workflows/pr-validation.yml](.github/workflows/pr-validation.yml) does `go build .` + `go test -v ./test/...`; [.github/workflows/release.yml](.github/workflows/release.yml) cross-builds Win/Linux and publishes a release. `go vet ./...` already passes locally; no other static analysis is wired up. No `staticcheck`/`golangci-lint` installed.

### 6.1 — PR validation: missing quality gates  **[Medium]**
Add to the PR workflow:
1. **`go vet ./...`** — zero-cost, already green; lock it in.
2. **`golangci-lint`** (see [§7](#7-static-analysis--linting-golang)) — `govet`, `staticcheck`, `errcheck`, `ineffassign`, `unused`, `gosimple`, `misspell`, `gocritic`, `revive`.
3. **Race detector + coverage:** `go test -race -covermode=atomic -coverprofile=coverage.out ./...` and upload/gate coverage (Codecov or a threshold check). The GUI uses Gio with frame scheduling — `-race` is valuable.
4. **`govulncheck ./...`** — known-vuln scan of deps + stdlib.
5. **Build the *whole* module** (`go build ./...`), not just `.`, so all packages compile in CI.
6. **gofmt/goimports check** — *blocked today* because the repo is CRLF and `gofmt -l .` flags all 233 files. First add a [.gitattributes](.gitattributes) (`*.go text eol=lf` or `text=auto`) and normalize, *then* enable a formatting gate. Until then, a format check is impractical (documented in repo memory).
7. **Test the full module** (`./...`) rather than only `./test/...`, so example/integration packages outside `test/` are covered.

Example lint+vet job:
```yaml
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.25.8', cache: true }
      - run: go vet ./...
      - uses: golangci/golangci-lint-action@v6
        with: { version: latest }
      - run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...
```

### 6.2 — OS matrix  **[Medium]**
Tests run Linux-only, yet the most OS-sensitive code (`io.go`) targets Windows. Add `windows-latest` (and ideally `macos-latest`) to the test matrix. Gio system deps only need installing on Linux (guard with `if: runner.os == 'Linux'`).

### 6.3 — Release hardening  **[Low/Medium]**
[.github/workflows/release.yml](.github/workflows/release.yml):
- Add `-trimpath` and inject version via `-ldflags "-X main.version=..."` for reproducible, identifiable binaries.
- Publish **SHA-256 checksums** (and optionally an SBOM via `anchore/sbom-action`, and SLSA provenance) alongside binaries.
- Add a `concurrency:` group to cancel superseded tag builds.
- Pin third-party actions to commit SHAs (supply-chain hardening) — currently pinned to major tags only.
- The Linux build uses `CGO_ENABLED=1` (needed for Gio) but the runner installs X11/Wayland dev libs only in `pr-validation`; confirm the release job has them too if CGO needs them at link time, otherwise the Linux artifact may fail to link.

### 6.4 — Dependency hygiene  **[Low]**
Add **Dependabot** ([.github/dependabot.yml](.github/dependabot.yml)) for `gomod` and `github-actions`. Add a `go mod verify` / `go mod tidy -diff` (Go ≥1.25) check to fail PRs that leave `go.mod`/`go.sum` dirty.

### 6.5 — Toolchain pin  **[Low]**
[go.mod](go.mod) declares `go 1.25.8`; the local toolchain is `go1.26.3`. Consider adding a `toolchain go1.25.8` directive (or bump deliberately) so contributors and CI build with the same compiler.

---

## 7. Static Analysis & Linting (Go)

Go has a rich, mostly-free static-analysis ecosystem. Recommended adoption order:

### 7.1 — Add a lint + static-analysis job
- **`go vet`** — built in, already green.
- **`staticcheck`** (Honnef) — would have flagged the nil-map write in §1.1 (`SA5000`/`SA4000` family) and dead code.
- **`golangci-lint`** as the umbrella runner. Suggested [.golangci.yml](.golangci.yml):
```yaml
run:
  timeout: 5m
linters:
  enable:
    - govet
    - staticcheck
    - errcheck        # flags the swallowed os.Mkdir / Stat errors in §1.4 / dialogs
    - ineffassign
    - unused          # flags dead code like linq.ToMap if it stays unused
    - gosimple
    - misspell
    - gocritic
    - revive
    - bodyclose
    - nilerr
issues:
  exclude-rules:
    - linters: [revive]
      text: "receiver name should be"   # house style uses `this`; opt out explicitly
```

### 7.2 — `govulncheck`
Run `govulncheck ./...` in CI to catch vulnerable module/stdlib versions (low noise — only reports vulns on reachable code paths).

### 7.3 — `errcheck`
Several errors are dropped (`os.Mkdir` in the file explorer, `os.Stat` non-`IsNotExist` cases in `io.go`). `errcheck` will enumerate them so they can be handled or explicitly ignored with `_ =`.

### 7.4 — Enforce architecture with a linter
Add an import-boundary check so `internal/**` can never import `app/gui/**` again (prevents §3.1 from regressing). Options: `depguard` (built into golangci-lint) or `go-arch-lint`. Example `depguard` rule:
```yaml
linters-settings:
  depguard:
    rules:
      no-ui-from-internal:
        files: ["**/internal/**"]
        deny:
          - pkg: "github.com/Tariomka/hommoe_custom_templates/app/gui"
            desc: "internal/ must not depend on the UI layer"
```

### 7.5 — `gofmt`/`goimports`
Blocked by CRLF (see §6.1). Land [.gitattributes](.gitattributes) + a one-time normalization first.

---

## 8. Prioritized Action Plan

**P0 — correctness (do first):**
1. Fix `linq.ToMap` nil-map panic (§1.1) + add a unit test.
2. Guard the VDF type assertion (§1.2) + test malformed/missing input.
3. Replace hard-coded `C:`/`USERNAME` paths with `os.UserHomeDir`/registry lookup (§1.3).

**P1 — architecture & safety net:**
4. Move `ContentIds` (and any domain catalog) out of `app/gui/constants`; remove the `internal → app/gui` import; add a `depguard` boundary rule (§3.1, §7.4).
5. Add `golangci-lint` + `go vet` + `govulncheck` + `-race`/coverage to PR CI (§6.1, §7).
6. Add tests for handlers, mapper, and `io.go` (with injectable FS/env) (§5.1).

**P2 — maintainability:**
7. Extract labeled-row widget helpers (§4.1); cache per-frame allocations (§2.2).
8. Delete/port `oldTests/` (§4.2); replace the 1,364-line fixture with a `-update` golden file (§5.2).
9. Split `zoneEditorDialog.go` and `previewLayout.go` (§3.2).
10. Add `.gitattributes` + normalize line endings, then enable a gofmt gate (§6.1).

**P3 — polish:**
11. Release hardening: checksums, `-trimpath`, version ldflags, Dependabot, pinned action SHAs (§6.3–6.4).
12. Resolve the scattered `TODO`s (victory-condition fallback §1.7, road-rebuild §1.8, `MinNeutralZonesBetweenPlayers` dead slider, etc.).
13. Decide on `this`-receiver policy and document it (§3.3).

---

## Appendix A — Notable TODOs found in source
- [internal/handlers/guiHandler.go](internal/handlers/guiHandler.go#L52) — `// TODO: might not be needed` (road rebuild).
- [internal/services/content_rules/variantMappingManager.go](internal/services/content_rules/variantMappingManager.go#L6) — `// TODO: This should not exist` (UI import).
- [app/gui/constants/victoryConditions.go](app/gui/constants/victoryConditions.go#L50) — fallback to first condition.
- [app/gui/panels/layoutPanel.go](app/gui/panels/layoutPanel.go#L164), [app/gui/panels/generalPanel.go](app/gui/panels/generalPanel.go#L147) — `// TODO: check .Update(gtx)`.
- [app/gui/panels/layoutPanel.go](app/gui/panels/layoutPanel.go#L259) — commented-out `MinNeutralZonesBetweenPlayers` slider (dead).
- [internal/services/zones/zoneLabelProvider.go](internal/services/zones/zoneLabelProvider.go#L40-L41) — clamp/length TODOs.
- [app/gui/drivers/state.go](app/gui/drivers/state.go#L250) — `// TODO: add validator for state updates`.

## Appendix B — What’s already good
- Clear layered pipeline (GUI → handlers → services → providers) with DTO seams.
- Topology system is genuinely extensible and the geometry helpers avoid copy-paste.
- Cross-platform file handling uses `path/filepath` consistently (outside the hard-coded roots in §1.3).
- `go vet ./...` is clean.
- Connection-editor and mandatory-content suites are well-structured and pin specific bug fixes.
- Sensible error sentinels in [internal/common/editorErrors.go](internal/common/editorErrors.go).
