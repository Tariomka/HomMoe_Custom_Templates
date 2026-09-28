# Batch G: Zone-editor graph diagnostics (review §3.1)

Measure the per-frame cost of the zone editor's status line, then cache its graph
summary behind a dedicated dirty flag. Also make every late status update trigger an
immediate redraw, and optionally rewrite isolation detection to O(zones + connections),
but only if the benchmark justifies it. Scope is review §3.1 of
[review-gpt-6-astra-09-07.md](../backlog/review-gpt-6-astra-09-07.md) and nothing else.

## For Future Agents
As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

Read [AGENTS.md](../../AGENTS.md) and [the handoff](../session-carry-forward.md) first.
The owner stages and commits; never run any Git mutation. Never keep a blanket GUI
`-update`: accept only the `.failure` files a plain run produces.

## Context (verified 2026-09-28 on `AD/performance`, HEAD `f78caeb`)

Batch F is merged to `master` as `f78caeb` (#48). The owner created `AD/performance` from
it, and the tree was clean when this plan was written.

- [layoutStatus](../../app/gui/dialogs/zoneEditorDialog.go) runs on every toolbar layout.
  It calls `derefConnections(this.working)`, which allocates a slice, and then
  `zoneHandler.DescribeZoneEditorGraph`. That calls `ComputeHasErrors`, which allocates a
  map, and `FindIsolatedZones`, which is a nested zones × connections loop with no map
  ([connectionEditorService.go](../../internal/services/connection_editor/connectionEditorService.go)).
- `Body` order is: click polling, then the toolbar (which lays out the status), then the
  canvas (`handlePointer`, then the side update, then `ensureGeometry`), then the side
  panel (layout, then `writebackProps` / `writebackZoneProps`). Any change made by canvas
  input or side-panel writeback therefore reaches the status line only on the **next**
  frame, and nothing requests that frame.
- Diagnostics depend only on the zone **name** set and on connection **`From`/`To`**.
  The mutators that can change those (the structural mutators) are:
  1. `setEditingSet` (the constructor and `revertToBase`);
  2. `addConnection`;
  3. `deleteConnection`;
  4. `undoSessionEdits`;
  5. `addZoneAt`;
  6. `deleteZone`;
  7. `applyQualityMutation`, which replaces both lists.

  `moveDraggedZone`, `ensureManualPositions` and `writebackProps`/`writebackZoneProps`
  outside the quality branch never change names or endpoints. `ChangeZoneEditorConnectionType`
  clones and keeps `From`/`To`.
- `geometryDirty` is cleared by `ensureGeometry`, including from hit tests in the middle
  of input handling. **The summary cache must never read or share that flag.**
- The label pool (`constants.GetZoneLabels`) holds 32 labels, so the realistic upper bound
  is about 40 zones (32 neutrals plus up to 8 spawns).
- Gio v0.10.2: `gtx.Execute(op.InvalidateCmd{})` sets the router's wakeup, and
  `input.Router.WakeupTime()` returns `(time.Time{}, true)` until it is read. Pending
  events also force `true`, so tests need a negative control (see Phase 3).
- The GUI dialog tests construct `dialogs.NewZoneEditorDialog` with
  `composition.InitializeGuiHandler()`. The `AppRunner` builds its window from the same
  root, and `NewAppRunnerWithFileSystem` is the precedent for injecting one handler.

## Owner decisions (2026-09-28)

| ID | Decision |
| --- | --- |
| G1 | Measure first, then implement the cache (the threshold G2 is met by construction). |
| G2 | The acceptance criterion is deterministic: idle frames make **0** `DescribeZoneEditorGraph` calls, and each structural mutation makes exactly **1** call on the following frame. Benchmarks are recorded before and after with **no percentage gate** and no global allocation threshold. |
| G3 | The O(zones + connections) `FindIsolatedZones` rewrite is adopted **only** if its ns/op is lower at 24 and 40 zones, and no worse (within **±5%**) at 4 and 12 zones. Otherwise the numbers are recorded and the nested loop stays. |
| G4 | The cache lives in the dialog, in its own state struct with its own dirty flag. The handler stays stateless. |
| G5 | Handler calls are counted in the GUI suite through a wrapper in the test helpers that embeds the real `IGuiHandler`. No new test exports. |
| G6 | Add an untagged, GPU-free idle-frame benchmark of `ZoneEditorDialog.Body` beside the handler benchmark. |
| G7 | Fix the late status: when anything the status line shows changed after the toolbar was laid out, `Body` requests an immediate redraw with `op.InvalidateCmd`. Layout order and Batch F's input-before-side-update ordering stay unchanged. |
| G8 | Late-status scope is **everything** the line shows: hint, add-connection mode, add-zone mode, zone count, connection count and the diagnostics. |
| G9 | Benchmark sizes are 4 / 12 / 24 / 40 zones: three generated templates (2p+2n, 4p+8n, 8p+16n) plus one fixed 40-zone fixture. |
| G9a | **Owner-approved amendment to G9 (2026-09-28):** all four fixtures are built deterministically at the same zone counts (spawns and neutrals as in 2p+2n, 4p+8n, 8p+16n and 8p+32n), with fixed connection lists. Reason (plan review finding 3): generation uses the unseeded global `rand`, so separate before and after runs would benchmark different connection sets and orders, and the nested loop's early `break` makes its cost order-dependent. That would invalidate the ±5% G3 check. |

The binding constraints from review §3.1: do not reintroduce the rejected global Gio
allocation threshold, and do not introduce clone-free live state pointers.

## Out of scope

- Everything else in review §2, §4, §6 and §7; Batch K topology retirement; DTO, schema
  and package work.
- Allocation tuning beyond this path, including `ComputeHasErrors`'s map (it now runs only
  on dirty frames).
- Reordering `Body` layout, moving canvas input before the toolbar, and all settled
  Batch C, E and F behaviour (the shared curve builder, hit tests against the live list,
  tri-state roads, the Preview-only legend, the PNG mask).
- Protected paths (`data/`, `internal/entities/template_entity/`, `internal/registry/`),
  the output directory, and generated Wire (no constructor changes are planned).

## Plan review record

- 2026-09-28, Claude Opus 5.5: **APPROVE WITH CHANGES**, 3 major, 7 minor and 2 nits.
  It confirmed that the seven-mutator list is complete, that `gtx.Execute(op.InvalidateCmd{})`
  wakes both the real window and the aux router, and that `app/gui/dialogs` pulls in no
  `gioui.org/app` or `gioui.org/gpu`, so the idle-frame benchmark is GPU-free. Every
  finding was applied:
  1. A redraw loop was possible: the key is now built from the full state and stored
     unconditionally, and a convergence test was added.
  2. `AppRunner` wakeups are confounded by other sources: the invalidate tests are now
     dialog-direct with a negative control, and `WakeupRequested` was dropped.
  3. Nondeterministic fixtures would undermine G3: fixtures are now deterministic, which
     needs owner approval as G9a. The benchtime is time-based, and every benchmark
     reports allocations.
  4. Revert-to-base is reachable through the window.
  5. The benchmark must not touch the `RecomputeGeometry` test export.
  6. `zoneEditorStatusKey` gets its own file.
  7. The counter is mutex-guarded and tagged.
  8. The call-count delta procedure is now precise.
  9. The tests were split, one redundant test dropped, and tests added for the open and
     refused-delete cases.
  10. The Phase 4 coverage wording now matches Phase 6.
  11. Spawns are named with `GetPlayerZoneNameFor`.
  12. The runner constructors share one private constructor.
- 2026-09-28, owner: plan and G9a approved ("You can capture the baseline and proceed. Go
  through all of the Phases and I will review everything at the end").

## Phase 1: Baseline and benchmarks (no production edits)
Status: Complete

- [x] Record the environment: `go version`, `go env GOOS GOARCH`, `$env:GOFLAGS` (must be
      empty), branch and HEAD.
- [x] Capture the baseline on the unchanged code:
  - `go build ./...`;
  - `go test -count=1 -p=2 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...`,
    then `go tool cover '-func=coverage.txt'` total (74.6% on record);
  - `go test ./test/... -count=1`;
  - `go test -tags='integration_test,gui' ./test/integration/... -count=1`, recording
    root and GUI package times;
  - `go run ./cmd/testlayoutcheck .`, `gofmt -l` over the Go tree (never `data/`), and
    `golangci-lint-v2 run ./... --issues-exit-code=1` (run `golangci-lint-v2 cache clean`
    first if stale `nolintlint` hits appear).
- [x] Add `test/performance/zoneEditorGraph_test.go` (untagged, `package performance_test`):
  - A shared deterministic fixture builder (G9a) returning `(zones, connections)` for the
    four sizes. Spawns are named with `constants.GetPlayerZoneNameFor` and neutrals with
    `constants.GetNeutralZoneNameFor`, both over `constants.GetZoneLabels`, so the
    dialog's `IsZoneNamePlayer` classification holds. Connections are a fixed ring plus
    chords, no randomness, with every zone referenced. Each `b.Run` `require`s the exact
    zone and connection counts before the timer loop.
  - `BenchmarkZoneEditorHandler_DescribeZoneEditorGraph`: one `b.Run` per size, through
    `composition.InitializeGuiHandler()`, with `b.ReportAllocs()` and `b.Loop()`, and a
    `require` sanity check on the result after the loop.
  - `BenchmarkConnectionEditorService_FindIsolatedZones`: one `b.Run` per size, through
    `connection_editor.NewConnectionEditorService`, with `b.ReportAllocs()`. This is the
    G3 comparison benchmark.
  - `BenchmarkZoneEditorDialog_IdleFrame`: one `b.Run` per size, with `b.ReportAllocs()`.
    It builds the dialog with `GetZoneEditorOptions` plus `NewZoneEditorDialog` **only**,
    and uses no identifier declared in any `*_testexports.go` file (in particular not
    `RecomputeGeometry`). It uses a local layout context identical to `newDialogContext`
    (`op.Ops`, a persistent `input.Router`, exact 1000×720 constraints, PxPerDp 1). Two
    warm-up `Body` frames build the geometry; then, in `b.Loop()`, it resets the ops and
    calls `Body` and `router.Frame`. No window, no GPU, no rasterizing, and nothing
    selected, so no editor carets are laid out.
- [x] Verify the file needs no build tag: it references production APIs only
      (`dialogs.NewZoneEditorDialog`, `Body`, the composition root and services).
- [x] Run each benchmark `-count=6` on the unchanged code and record the median ns/op,
      B/op and allocs/op per size in the Phase Summary. Use time-based benchtime, because
      the service call takes tens of nanoseconds:
      `go test -bench='ZoneEditor|FindIsolatedZones' -run=xxx ./test/performance/... -benchtime=1s -benchmem -count=6 -timeout=900s`.

### Verification Plan
- The baseline commands above all pass, and coverage is recorded.
- The new benchmarks compile untagged and run: `go vet ./test/performance/...` is clean,
  and the benchmark command prints results for all four sizes of all three benchmarks.
- `go test ./test/... -count=1` still passes (benchmarks do not run, and the file adds no
  tests).
- `go run ./cmd/testlayoutcheck .` passes.

### Phase Summary
**Complete, 2026-09-28.** Environment: Go 1.27.0 windows/amd64, empty `GOFLAGS`, branch
`AD/performance`, HEAD `a587617` (the owner's plan commit on top of `f78caeb`). CPU: Intel
Core Ultra 7 165H.

Baseline on the unchanged code:
- `go build ./...`: PASS.
- `-p=2` unit coverage run: PASS, **74.6%**.
- `go test ./test/... -count=1`: PASS.
- Tagged `integration_test,gui` run: PASS, root 6.453s, GUI 47.423s. This is slower than
  Batch F's 3.4s/27.0s, and the tests are unchanged, so it is likely machine load.
- `testlayoutcheck`: PASS. `gofmt -l`: clean. Lint: 0 issues.

Added [zoneEditorGraph_test.go](../../test/performance/zoneEditorGraph_test.go), untagged.
The fixtures are deterministic (G9a). Spawns come first, then neutrals, pinned on a circle,
with a ring plus half-chords, so every zone is referenced. The sizes are 4/6, 12/18, 24/36
and 40/60 (zones/connections). `go vet` is clean, `testlayoutcheck` passes, and the untagged
run reports no tests to run.

Baseline medians of 6 runs at `-benchtime=1s` (the raw output is in the gitignored
`tmp/batch-g-bench-before.txt`):

| Benchmark | 4 | 12 | 24 | 40 |
| --- | ---: | ---: | ---: | ---: |
| DescribeZoneEditorGraph ns/op | 209.65 | 1044.5 | 3180.5 | 7074 |
| DescribeZoneEditorGraph B/op, allocs | 0, 0 | 456, 3 | 936, 3 | 1832, 3 |
| FindIsolatedZones ns/op | 73.09 | 502.35 | 2016 | 4664.5 |
| FindIsolatedZones B/op, allocs | 0, 0 | 0, 0 | 0, 0 | 0, 0 |
| Dialog idle frame ns/op | 70224 | 104624.5 | 151050.5 | 215030 |
| Dialog idle frame B/op, allocs | 4457, 58 | 7611, 63 | 12204, 67 | 18478, 67 |

At 4 zones `ComputeHasErrors`'s map stays on the stack, so there are no allocations; from
12 zones up it escapes.

## Phase 2: Graph summary cache
Status: Complete

- [x] New `app/gui/dialogs/zoneEditorGraphState.go` with one struct,
      `zoneEditorGraphState`:
      `graph dtos.ZoneEditorGraphDto` and `graphDirty bool`. Embed it in `ZoneEditorDialog`
      beside the other state structs.
- [x] A private `markGraphDirty()` on the dialog. Call it in each of the seven structural
      mutators listed in Context, **in addition to** their existing `geometryDirty = true`,
      and never as a replacement. The constructor path runs through `setEditingSet`, so the
      first frame is dirty.
- [x] `layoutStatus`: when `graphDirty`, rebuild with
      `DescribeZoneEditorGraph(this.zones, derefConnections(this.working))`, store it and
      clear the flag. Otherwise use the cached value. Replace `len(connections)` with
      `len(this.working)`, so the idle path allocates no connection slice.
- [x] Do **not** touch `geometryDirty`, `ensureGeometry` or hit testing.
- [x] Keep the connection slice passed to the handler a fresh clone (`derefConnections`),
      never a live pointer view.

### Verification Plan
- `go build ./...` and `go test ./test/unit/... -count=1` pass.
- The GUI suite passes unchanged: `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1`
  (no golden churn is expected; the status text is identical frame for frame once the
  next frame is laid out).
- A code search shows that every structural mutator calls `markGraphDirty` and that no
  code reads `geometryDirty` for the summary.

### Phase Summary
**Complete, 2026-09-28.** `zoneEditorGraphState` (the cached `ZoneEditorGraphDto`,
`graphDirty`, and `shownStatus` from Phase 3) is embedded in `ZoneEditorDialog`.
`markGraphDirty` is promoted from the state struct and called in all seven structural
mutators, beside their unchanged `geometryDirty = true`: `setEditingSet`,
`addConnection`, `deleteConnection`, `undoSessionEdits`, `addZoneAt`, `deleteZone` and
`applyQualityMutation`. `layoutStatus` rebuilds only when dirty, still through
`derefConnections` (a fresh value copy), and counts connections with `len(this.working)`.
Build and unit suite PASS. The existing GUI suite passed unchanged (53.9s) with no
`.failure` files.

## Phase 3: Immediate redraw for late status changes
Status: Complete

- [x] New `app/gui/dialogs/zoneEditorStatusKey.go` with the unexported comparable value
      type `zoneEditorStatusKey` `{hint string; addMode, addZoneMode bool; zoneCount,
      connectionCount int; graphDirty bool}` (its own file per AGENTS §4.1). Add a
      `shownStatus zoneEditorStatusKey` field to `zoneEditorGraphState`.
- [x] A private `statusKey()` method on the dialog builds the key from the **full** dialog
      state, whichever branch the status line draws.
- [x] `layoutStatus` stores `shownStatus = this.statusKey()` unconditionally, after any
      rebuild and before branching. This guarantees convergence: an unchanged state always
      matches, even in the error branch or while a hint hides the counts.
- [x] At the end of `Body`, on the normal return only (the Apply and Cancel returns close
      the dialog), compare `this.statusKey()` with `shownStatus`. If they differ, call
      `gtx.Execute(op.InvalidateCmd{})`. This captures canvas presses and releases (add or
      delete connection, add zone, add-mode exits, hints), zone and connection deletes from
      the canvas, and side-panel quality writeback. No allocation: string comparison of
      `hint` plus scalar fields.
- [x] Toolbar-button mutations run before the toolbar layout, so they do not trigger a
      redraw. That is intended.

### Verification Plan
- `go build ./...` and the unit suite pass.
- Observability spike first, in the **same dialog-direct harness** the Phase 4 invalidate
  tests use (record the outcome in the Phase Summary). With one persistent router and no
  events, drain `WakeupTime()`, lay out a frame, and confirm the next `WakeupTime()` is
  `false`. If Gio or a widget forces a wakeup on idle frames, the negative control is
  impossible. In that case, do **not** add a production seam; record the gap in
  [test_observations.md](../backlog/test_observations.md) and keep only the
  status-correctness and call-count tests.
- The GUI tests from Phase 4 for the invalidate pass.

### Phase Summary
**Complete, 2026-09-28.** The key lives in its own file. `statusKey()` builds it from
the full state, `layoutStatus` stores it unconditionally after any rebuild, and
`requestLateStatusRedraw(gtx)` (extracted from `Body` to satisfy `funlen`) issues
`op.InvalidateCmd` when the end-of-`Body` key differs.

**Deviation: the invalidate tests run through the window harness, not dialog-direct.**
A dialog-direct test cannot place canvas input: the canvas box's offset depends on the
toolbar height, and only the window harness has measured canvas constants. The spike
therefore ran in the window harness instead:
- `AppRunner.ImmediateRedrawRequested()` counts only zero-time wakeups, so caret blinks
  are excluded.
- `settleRedraws` lays out frames (bounded at 3s) until the button-ink animations from
  opening the editor stop.
- The empty-canvas right-click control then reports **no** redraw, so the signal is
  isolated.

The error-state convergence test stays dialog-direct and idle, because
`ApplyEditedZones` rejects connections to missing zones, so the window cannot reach
that state.

Temporary mutation checks, all reverted:
- Disabling the redraw failed both redraw tests.
- Recording a drawn-fields-only key failed the error-state test.

## Phase 4: Tests
Status: Complete

Unit tests (untagged, mirrored layout, triple-A, `t.Parallel()`, gofakeit where it fits):

- [x] `test/unit/internal/services/connection_editor/connectionEditorService/findIsolatedZones_test.go`, add:
  - no zones, returns `nil` (`assert.Nil`, locking the contract before Phase 5);
  - no connections, returns every zone in input order;
  - a zone referenced only as `To`, is not isolated;
  - a connection to an unknown zone, does not un-isolate a known zone;
  - duplicate zone names that are both unreferenced, both reported (the current behaviour,
    locked before any rewrite);
  - a self-loop connection, the zone is not isolated.
- [x] `test/unit/internal/handlers/zoneEditorHandler/describeZoneEditorGraph_test.go`, add:
  - a `nil` isolated result gives count 0.
  (Argument forwarding is already locked by the exact-argument `On(...)` expectations.)

GUI integration tests (tagged `integration_test && gui`, because they use test exports
and the window):

- [x] Test helper `test/test_helpers/integration_common/graphDescriptionCounter.go`: one
      struct, `GraphDescriptionCounter`, which embeds `handler_interfaces.IGuiHandler`,
      overrides `DescribeZoneEditorGraph` to count calls and delegate, and exposes
      `Calls() int` and `Last() dtos.ZoneEditorGraphDto` (the last result). Both are
      guarded by a mutex, because in headed mode the render goroutine writes them while
      the test reads them. Tag it `//go:build integration_test` for consistency with every
      other `integration_common` file (testlayoutcheck only inspects `_test.go` files).
- [x] `AppRunner`: add `NewAppRunnerWithGuiHandler(tb, gui handler_interfaces.IGuiHandler)`.
      It and `NewAppRunnerWithFileSystem` share one private constructor taking both
      handlers, rather than duplicating the body.
- [x] New `test/integration/gui/zoneEditorGraphCache_integration_test.go`, driven through a
      window built with the counter. Every helper already ends in `NextFrame`, so each
      count test works like this: open the editor, lay out one more `NextFrame`, take the
      baseline count, run the helper, lay out one more `NextFrame`, and assert the delta.
      This also proves nothing rebuilds on every frame.
  - Opening the editor makes exactly 1 call (the constructor path).
  - Idle: several `NextFrame`s with no input, delta 0.
  - Each structural mutator gives delta exactly 1, one test each:
    - add a connection by dragging (`DragFromZoneTo` in add mode);
    - delete a connection with a right click (`RightClickEdge`);
    - delete a connection through Delete selected;
    - add a zone (add-zone mode plus `ClickCanvasAt`);
    - delete a **neutral** zone;
    - Undo;
    - revert to base (`ClickRevertToBase`, which the window wires to
      `state.PreviewBaseZones`);
    - change a neutral zone's quality to a **different** quality (`SelectZoneQuality`;
      `WasUpdated` only fires on a change).
  - Delta 0, one test each: dragging a zone, typing a guard value, changing the
    connection type, and a refused spawn-zone delete.
  - Correctness, two tests: deleting a zone's only connection makes `Last()` report one
    more isolated zone; reconnecting it makes `Last()` report 0 isolated zones.
- [x] Invalidate tests, only if the Phase 3 spike succeeded. **Implemented through the
      window harness instead (see the Phase 3 deviation).** They run **dialog-direct**
      in the same test file, with one persistent router and ops (not through `AppRunner`,
      whose window, dialog host, button ink and pending events all set wakeups). Toolbar
      buttons use the programmatic `Click()` exports, which create no ink; canvas input
      is queued at `CanvasOrigin()` plus the square-local position. Drain `WakeupTime()`
      before acting. One test each:
    - right-clicking an edge (a late structural change) leaves a wakeup requested;
    - right-clicking empty canvas (the same gesture shape, no status change) leaves none;
    - the frame after the right-click on an edge leaves none (convergence);
    - convergence in the error state: with a connection to a missing zone and a hint set,
      the frame after a late mutation leaves none.
- [x] No new goldens are expected. If a plain run produces `.failure` files, inspect each
      one before accepting it; never use a blanket `-update`.

### Verification Plan
- `go test ./test/unit/... -count=1` passes.
- `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` passes, and
  the new tests pass with `-count=5`.
- `go run ./cmd/testlayoutcheck .` passes.
- Coverage run (same command as the baseline): the total is **≥ 74.6%**, or any drop is
  explained and owner-approved. The new dialog statements sit inside
  `-coverpkg=./app/...` but only the GUI suite reaches them, so a small dip is possible.
  Every touched non-GUI function is at 100%. Record the dialog files as GUI-only in
  test_observations.md if they are not already.

### Phase Summary
**Complete, 2026-09-28.**

Unit tests:
- 6 new `FindIsolatedZones` tests, which lock the `nil` contract, input order,
  `To`-only references, unknown endpoints, duplicate names and self-loops.
- 1 new handler test for a `nil` isolated result.

GUI integration: 20 new tests in
[zoneEditorGraphCache_integration_test.go](../../test/integration/gui/zoneEditorGraphCache_integration_test.go):
- Call counts:
  - opening the editor makes 1 call; idle frames make 0;
  - 8 structural edits make 1 each (draw, right-click delete, delete selected, place
    zone, delete Hub, Undo, revert to base, quality);
  - 4 non-structural edits make 0 each (drag zone, guard value, connection type,
    refused spawn delete).
- 2 isolation correctness tests.
- 3 window-harness redraw tests (positive, negative control, convergence).
- 1 dialog-direct error-state convergence test.

Helpers:
- `GraphDescriptionCounter` (`integration_test`, mutex-guarded).
- `NewAppRunnerWithGuiHandler` and `NewAppRunnerWithFileSystem` now share `newAppRunner`.
- `AppRunner.ImmediateRedrawRequested`.

Mutation check: making `layoutStatus` rebuild on every frame failed all 14 call-count
tests. Dropping `markGraphDirty` from `applyQualityMutation` failed the quality test.
All mutations were reverted.

Results:
- New tests at `-count=5`: PASS.
- Full tagged suite: PASS.
- `testlayoutcheck`: PASS.
- No `.failure` files.
- Coverage 74.6% -> **74.5%**: covered statements are unchanged at 6911, while the
  denominator grew 9246 -> 9261 with 15 GUI-only dialog statements. Every touched non-GUI
  function is at 100%. This is recorded in
  [test_observations.md](../backlog/test_observations.md) and needs owner acceptance.

## Phase 5: Optional O(zones + connections) isolation rewrite (G3)
Status: Complete (rewrite rejected)

- [x] Implement the candidate in `FindIsolatedZones`: build a
      `map[string]struct{}` of every `From`/`To` sized `2*len(connections)`, then collect
      unreferenced zones in input order. Keep the `nil` result for no isolated zones.
- [x] Run `BenchmarkConnectionEditorService_FindIsolatedZones` with `-count=6` against the
      Phase 1 medians.
- [x] Apply the G3 rule: lower ns/op at 24 and 40 zones, and within ±5% at 4 and 12 zones.
  - If it passes, keep it, and the Phase 4 unit tests must still pass unchanged.
  - If it fails, revert the function to its Phase 1 text exactly and record the numbers.

### Verification Plan
- The unit tests for `FindIsolatedZones` and `DescribeZoneEditorGraph` pass either way.
- The decision and the medians are recorded in the Phase Summary.

### Phase Summary
**Complete, 2026-09-28: rewrite REJECTED under G3.** Medians of 6 runs (the raw output
is in the gitignored `tmp/batch-g-bench-rewrite.txt`):

| Zones | Nested loop ns/op | Map rewrite ns/op | Change | Rewrite B/op, allocs |
| ---: | ---: | ---: | ---: | --- |
| 4 | 73.09 | 286.3 | +292% | 456, 3 |
| 12 | 502.35 | 838.45 | +67% | 1832, 3 |
| 24 | 2016 | 1628 | −19% | 3496, 3 |
| 40 | 4664.5 | 3138 | −33% | 6568, 3 |

The rewrite is faster only at 24 and 40 zones, and it adds 3 allocations per call, so the
±5% rule at 4 and 12 zones fails. The function was reverted; `git diff` on
`connectionEditorService.go` is empty. The unit tests passed against both versions.

## Phase 6: Final verification and review
Status: Complete (awaiting owner review)

- [x] Rerun the full Phase 1 baseline set and all three benchmarks (`-count=6`), and
      record before/after medians per size.
- [x] `gofmt -l` clean, lint 0 issues, and `testlayoutcheck` passes.
- [x] Independent implementation review by Claude Opus 5.5, with its findings resolved or
      explicitly declined with a reason.
- [ ] Mark §3.1 FIXED in place in the review, with a Progress paragraph, and update the §9
      `G` row and the progress line (17 fixed, 15 remaining). Do not renumber, and do not
      edit §8. **Deferred:** the review's fix-session protocol marks an item FIXED only
      after the owner commits, and the owner chose to review everything at the end.
- [x] Update the handoff for the next session.

### Verification Plan
- Every Phase 1 command passes, and coverage is ≥ 74.6% (or any drop is explained and
  owner-approved).
- The acceptance criterion G2 is demonstrated by the Phase 4 call-count tests.

### Phase Summary
**Complete, 2026-09-28, awaiting owner review.** Final verification (Windows/amd64, Go
1.27.0, empty `GOFLAGS`):
- build: PASS;
- `-p=2` unit coverage: PASS, **74.5%** (see Phase 4);
- `go test ./test/...`: PASS;
- tagged `integration_test,gui`: PASS (root 3.322s, GUI 32.380s on the final code);
- `testlayoutcheck`: PASS; `gofmt -l`: clean;
- lint: **0 issues**, after fixing 4 new findings: `funcorder` ×2 and
  `gochecknoglobals` in the benchmark file, and `funlen` on `Body`.

Before/after medians of 6 runs (the raw output is in the gitignored
`tmp/batch-g-bench-after.txt`; `DescribeZoneEditorGraph` and `FindIsolatedZones` are
unchanged code, so their deltas are run-to-run noise):

| Dialog idle frame | 4 | 12 | 24 | 40 |
| --- | ---: | ---: | ---: | ---: |
| ns/op before → after | 70224 → 75499 | 104624.5 → 100634 | 151050.5 → 140862.5 | 215030 → 190195 |
| B/op before → after | 4457 → 3049 | 7611 → 3057 | 12204 → 3073 | 18478 → 3073 |
| allocs before → after | 58 → 57 | 63 → 59 | 67 → 63 | 67 → 63 |

Idle-frame bytes per frame are now flat in graph size. The unchanged service varied up to
9% between runs, so the ns/op movement at 4 zones is noise.

The independent implementation review (Claude Opus 5.5) returned **APPROVE WITH
CHANGES**: 1 minor and 4 nits.
- Applied: the error-state test now `require`s that `HasErrors` is set, the graph-state
  doc was corrected, and two helper docs were tightened.
- Declined, nit 5: the open-editor test's Act only reads the counter. Opening the editor
  is the action, and the test reads clearly.
- Noted, not fixed (out of §3.1 scope): the toolbar's "Delete selected" enabled state is
  still one frame late after a canvas selection, as at HEAD, because `hasSelection` is not
  in the status key.

The reviewer ran the build, vet, lint, the layout check, and the new GUI tests ×3.

## Final Recap
Batch G (review §3.1) is implemented and verified, and awaits owner review and commit.
- The zone editor's status line now describes the graph only after a structural edit,
  instead of on every frame, through a dialog-owned cache with its own dirty flag.
  Nothing ties it to `geometryDirty`.
- Any status change made after the toolbar was drawn now requests the next frame
  immediately, fixing the pre-existing one-frame-late status line.
- The O(z+c) isolation rewrite was measured and rejected under G3.
- Untagged benchmarks now cover the handler, the service and a GPU-free idle frame.
- 20 GUI tests and 7 unit tests were added, each shown to fail under a corresponding
  temporary mutation.

## Deployment Plan
Nothing is deployed. The owner:
1. Reviews the working-tree changes on `AD/performance` (listed in the handoff).
2. Decides whether to accept coverage 74.5% (covered statements unchanged, 15 GUI-only
   statements added).
3. Stages and commits.
4. After the commit, a follow-up turn marks §3.1 `✅ FIXED` in the review (Progress
   paragraph, §9 `G` row, progress line 17 fixed / 15 remaining) and closes this plan.

No schema migration, Wire regeneration, dependency change or output-path change is
involved.
