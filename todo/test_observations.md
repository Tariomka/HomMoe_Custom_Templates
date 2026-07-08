# Unit Test Observations — code unreachable or untestable by unit tests

This file tracks implementation code that cannot be (fully) unit-tested
through public entry points, per the rule: never add helpers/seams to
implementation code just to make it testable. Each entry needs manual
investigation by the maintainer.

Format: `path` — reason — suggested action.

## Untestable / unreachable code

- `internal/helpers/io.go` — entire file has NO unit tests. `FindOldenEraTemplatesDir` depends on
  host state that cannot be redirected: `USERNAME`/`HOME` are read at PACKAGE INIT into private vars
  (so `t.Setenv` cannot influence them), the Windows path is the hard-coded
  `C:\Program Files (x86)\Steam\steamapps\libraryfolders.vdf`, and every branch (VDF missing/parsed,
  app 3105440 present/absent, install dir exists) depends on the machine's Steam installation. Private
  helpers (`getVDFContent`, `getVDFFilePath`, `getBasePath`, `resolveGlob`) are only reachable through the
  same entry point. Suggested action: inject the env/paths via parameters or a config struct (also fixes the
  known `filepath.Join("C:", ...)` drive-relative bug), then add tests.
- `internal/services/template_generator/providers/topology/base/topologyBase.go`
  `buildNonAdjacentDerangement` — the deterministic-shift fallback requires 100 consecutive failed random
  attempts; practically unreachable in tests. Suggested action: none (safety net).
- `internal/services/template_generator/providers/topology/base/topologyBase.go` `CreateMissingConnections`
  — the duplicate-bridge-name `continue` guard is unreachable: hitting it would require a pre-existing
  non-Direct bridge with the same generated name, which would loop forever anyway. Suggested action: review
  the loop's exit conditions.
- `internal/services/content_rules/ruleVariant.go` `DisplayText` — the "Unforeseen Error" branch is
  unreachable via `NewRuleVariant` (constructor validates the id); covered only by constructing the exported
  struct directly. Suggested action: none.
- `internal/handlers/guiHandler.go` — `SaveState` marshal-error branch (DTO is always marshalable),
  `SaveTemplate` `previewGenerator == nil` skip (depends on embedded-asset init, not on inputs) and
  `GenerateTemplate`'s `ErrGeneratedTemplateInvalid` are unreachable via inputs. Suggested action: none
  (defensive guards).
- `internal/services/template_generator/providers/gameRulesProvider.go` — `entities.GameRules.TournamentRules`
  (bool) and `WinConditions.GladiatorArenaRegistrationStartWork` are never SET by the generator; they exist
  only for tolerant parsing of hand-authored templates, so generator-side tests cannot exercise them.
  Suggested action: none.
- `internal/entities/template/template_variant/connection.go` — `Connection.IsUserAdded` carries `json:"-"`
  (must never serialize). This contract has no public function to test against (pure struct tag); it is
  exercised implicitly by templateWriter/integration tests. The legacy test `TestIsUserAdded_IsNotSerialized`
  was dropped in the migration. Suggested action: keep the contract in mind when editing the schema.
- `internal/services/template_generator/providers/common.go` — no public API (package-private registry alias
  vars only); no test folder possible, covered transitively by every provider test.
- `internal/services/preview_service/assetFitter.go` — all identifiers private (`newAssetFitter`); covered
  indirectly through `PreviewGeneratorService.CreatePreviewImage` tests. Suggested action: none.
- `internal/services/connection_editor/zoneEditor.go` `FindOpenPosition` — on an empty board returns
  (0.9, 0.9): float accumulation makes the LAST tied corner win the strict `>` comparison. Not a bug, but the
  tie-breaking is floating-point-sensitive; tests pin the current behavior.
- `internal/services/connection_editor/connectionEditor.go` `WeeklyIncrementValues` — public VAR (not a
  function); pinned by a dedicated test file per the per-file coverage rule.

- `internal/entities/template/template_rule/winConditions.go` `MergeWinConditionsIfDoesNotExist` — the four
  error-return branches (marshal of destination/source, unmarshal into maps, marshal of merged map) are
  unreachable: `WinConditions` contains only JSON-safe field types, so `json.Marshal`/`json.Unmarshal` on it
  can never fail. Function capped at ~77% coverage. Suggested action: none (defensive guards).
- `internal/models/config/config_inner/bonusEntry.go` `GetHash` — the fallback branch when `json.Marshal`
  fails is unreachable: `BonusEntry` holds only string/int fields, so marshalling never errors. Capped at 80%.
  Suggested action: none (defensive guard).
- `internal/services/previewLayout.go` — several branches are unreachable via `BuildPreviewLayout` (the only
  public entry that drives the private layout code):
  - the `!ok { continue }` guard in the zone loop (line ~72): every dispatch path positions every zone
    (layoutMultiHub sweeps stragglers to the canvas centre), so a variant zone is never absent from Positions;
  - the `n == 0` guards in `layoutFixedPositions`, `layoutBalancedRings`, `layoutScatter` and
    `layoutRingOrHub`: BuildPreviewLayout early-returns when the variant has no zones and every internal
    caller passes a non-empty slice — consequently `scaledInt` (only called from those guards) stays 0%;
  - the `cnt == 0 { continue }` in the balanced-rings placement loop: every ring index maps to a present
    GeneratorRing tier which by construction holds at least one zone;
  - `relaxPasses` Pass B `elen2 < 1e-3 { continue }`: Pass A in the same iteration pushes any coincident
    connected pair apart to minDist (≥30 px) before Pass B runs;
  - `layoutMultiHub` `numHubs == 1` spoke-base branch: the function is only invoked with ≥2 hub zones;
  - `positionCentroid`/`minMax` empty-input guards: all callers pass ≥1 element.
  Suggested action: none (defensive guards).
- `internal/services/previewRenderer.go` — unreachable branches:
  - `WritePreviewPNG` `png.Encode` error return: encoding a valid RGBA into a freshly created file cannot
    fail without I/O fault injection;
  - `RenderPreviewImage` `NewAssetProvider` error return: preview assets are embedded PNGs that always decode;
  - the `allowed < 1 { continue }` guard in the border-fit loop: side is forced to 700 so allowed ≥ ~280;
  - `drawDashedQuadratic` `dashOn <= 0` / `dashOff < 0` guards: the sole caller passes fixed dashOn=9,
    dashOff=13 (scaled).
  Suggested action: none (defensive guards / fault-injection-only paths).
- `internal/services/template_generator/providers/topology/geometryHelpers.go` — every function/type in the
  file is private (`circlePoint`, `squarePerimeterPoint`, `nearestIndexInRange`, `pairBuilder`), so no
  dedicated test folder is possible; the helpers are covered indirectly through the topology service tests
  (circles/square/geometric/cross/fractal). Suggested action: none, unless the helpers are ever exported.
- `internal/services/template_generator/providers/topology/base/topologyBase.go` `CreateMissingConnections` /
  `CreateMissingPlayerConnections` — both append the fallback/bridge roads to `linq.FromSlice(zones).First(...)`
  results, which are VALUE COPIES of the zone structs; the road mutations are silently lost and never reach the
  caller's zones. Unit tests can therefore only assert the returned connections, not the road side-effect.
  Suggested action: investigate — either mutate `zones[i]` via index or document that roads are intentionally
  not added.

## Gio-UI-heavy files (covered by integration suite, not unit-testable)

These files require a Gio `layout.Context`/window/event loop to exercise and are validated by the gated
integration + performance suites (`go test -tags=integration_test ./test/integration/... ./test/performance/...`),
not by unit tests:

- `app/gui/program.go` — window/event loop entry point.
- `app/gui/widgets/*.go` (14 files) — `New*Widget` closures over `layout.Context`/`material.Theme`.
- `app/gui/dialogs/*.go` (9 files) — modal dialogs (bonus picker, file explorer, picker, rule, zone content,
  zone editor) built on Gio widgets/clip/paint.
- `app/gui/panels/*.go` (4 files) — generalPanel, layoutPanel, zoneContentPanel/bonusesPanel, previewPanel.
- `app/gui/components/*.go` (2 files) — dropdownSelector, segmentButtonGroup.
- `app/gui/editor/window.go`, `app/gui/editor/toolbar.go` — editor window composition/frame loop.
- `app/gui/drivers/tab.go`, `app/gui/drivers/dialogHost.go` — Gio widget wrappers.
- `app/gui/utils/draw.go` — clip/paint drawing helpers.
- `app/gui/themes/*`, `app/gui/constants/*`, `app/gui/interfaces/*` — pure data/color/interface catalogs
  (no logic to test).

- `app/gui/drivers/state.go` — evaluated for unit tests (Phase 7) and SKIPPED: the only constructor
  `NewUIState()` is machine-dependent — it probes the Steam library via `helpers.FindOldenEraTemplatesDir`
  (host filesystem + `USERNAME`/`HOME` read at package init) and falls back to `os.Getwd()`, so the initial
  status message and output path differ per host. A zero-value `drivers.State{}` bypasses invariants
  (nil `innerState`/`handler` panics in `AutoRegenerate`/`UpdateState`/`SetStatus`-adjacent flows), and adding
  a test-only constructor seam is forbidden. `AutoRegenerate` debounce and `GetStatus`/`SetStatus` are already
  exercised end-to-end by the gated integration suite
  (`test/integration/manualCastleReapply_integration_test.go` drives `AutoRegenerate(now)` /
  `AutoRegenerate(now+1s)`). Suggested action: none.


## Dead code found while testing

- `internal/services/settingsFileLoader.go` (`LoadSettingsFile`/`SaveSettingsFile`),
  `internal/services/templateWriter.go` (`WriteTemplate`) and
  `internal/services/previewRenderer.go` (`RenderPreviewImage`/`WritePreviewPNG`) — ZERO production callers.
  They are legacy duplicates superseded by `internal/services/file_service.FileService` and
  `internal/services/preview_service.PreviewGeneratorService` (used by `internal/handlers/guiHandler.go`).
  They still have unit tests (per-file coverage rule) but should be DELETED together with their test folders
  once confirmed obsolete.
- `internal/services/previewLayout.go` vs `internal/services/preview_service/previewLayoutService.go` —
  the preview-layout logic exists TWICE (near-identical `BuildPreviewLayout`, `ExtractZoneLetter`,
  `ClassifyZoneTier` + private layout dispatch). The old `services` copy is still LIVE
  (`app/gui/panels/previewPanel.go`, `app/gui/dialogs/zoneEditorDialog.go` call `services.BuildPreviewLayout`),
  while `preview_service` is used via `guiHandler`. Consolidate to one implementation, then drop the
  duplicate and its duplicated test folders (both are currently unit-tested separately).
- `internal/helpers/linq/map.go` `QueryMap.ToMap` — zero production callers; tested anyway per the
  full-coverage rule. Suggested action: delete or start using.

- `internal/models/config/config_inner/bonusPresetType.go` `parseBonusPresetType` — private function with
  ZERO callers anywhere in the repo (the old `ParseBonusesJSON`/`SerializeBonuses` string round-trip it served
  was removed when `EditorStateDto.BonusesJSON` became `[]config_inner.BonusEntry` serialized via std json).
  Unreachable from any public entry point, stays 0%. Suggested action: delete the function (and the stale
  "see ParseBonusesJSON" comment in `internal/dtos/editorStateDto.go`).
