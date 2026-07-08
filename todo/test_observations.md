# Unit Test Observations — code unreachable or untestable by unit tests

This file tracks implementation code that cannot be (fully) unit-tested
through public entry points, per the rule: never add helpers/seams to
implementation code just to make it testable. Each entry needs manual
investigation by the maintainer.

Format: `path` — reason — suggested action.

## Untestable / unreachable code

(populated during the unit test refactoring, phases 1-8)

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

(populated in phase 8)

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

(populated during the unit test refactoring)

- `internal/models/config/config_inner/bonusPresetType.go` `parseBonusPresetType` — private function with
  ZERO callers anywhere in the repo (the old `ParseBonusesJSON`/`SerializeBonuses` string round-trip it served
  was removed when `EditorStateDto.BonusesJSON` became `[]config_inner.BonusEntry` serialized via std json).
  Unreachable from any public entry point, stays 0%. Suggested action: delete the function (and the stale
  "see ParseBonusesJSON" comment in `internal/dtos/editorStateDto.go`).
