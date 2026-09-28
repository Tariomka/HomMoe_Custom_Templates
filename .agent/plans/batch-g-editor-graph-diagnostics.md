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
Status: Not started

- [ ] Record the environment: `go version`, `go env GOOS GOARCH`, `$env:GOFLAGS` (must be
      empty), branch and HEAD.
- [ ] Capture the baseline on the unchanged code:
  - `go build ./...`;
  - `go test -count=1 -p=2 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...`,
    then `go tool cover '-func=coverage.txt'` total (74.6% on record);
  - `go test ./test/... -count=1`;
  - `go test -tags='integration_test,gui' ./test/integration/... -count=1`, recording
    root and GUI package times;
  - `go run ./cmd/testlayoutcheck .`, `gofmt -l` over the Go tree (never `data/`), and
    `golangci-lint-v2 run ./... --issues-exit-code=1` (run `golangci-lint-v2 cache clean`
    first if stale `nolintlint` hits appear).
- [ ] Add `test/performance/zoneEditorGraph_test.go` (untagged, `package performance_test`):
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
- [ ] Verify the file needs no build tag: it references production APIs only
      (`dialogs.NewZoneEditorDialog`, `Body`, the composition root and services).
- [ ] Run each benchmark `-count=6` on the unchanged code and record the median ns/op,
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
_(write when phase completes)_

## Phase 2: Graph summary cache
Status: Not started

- [ ] New `app/gui/dialogs/zoneEditorGraphState.go` with one struct,
      `zoneEditorGraphState`:
      `graph dtos.ZoneEditorGraphDto` and `graphDirty bool`. Embed it in `ZoneEditorDialog`
      beside the other state structs.
- [ ] A private `markGraphDirty()` on the dialog. Call it in each of the seven structural
      mutators listed in Context, **in addition to** their existing `geometryDirty = true`,
      and never as a replacement. The constructor path runs through `setEditingSet`, so the
      first frame is dirty.
- [ ] `layoutStatus`: when `graphDirty`, rebuild with
      `DescribeZoneEditorGraph(this.zones, derefConnections(this.working))`, store it and
      clear the flag. Otherwise use the cached value. Replace `len(connections)` with
      `len(this.working)`, so the idle path allocates no connection slice.
- [ ] Do **not** touch `geometryDirty`, `ensureGeometry` or hit testing.
- [ ] Keep the connection slice passed to the handler a fresh clone (`derefConnections`),
      never a live pointer view.

### Verification Plan
- `go build ./...` and `go test ./test/unit/... -count=1` pass.
- The GUI suite passes unchanged: `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1`
  (no golden churn is expected; the status text is identical frame for frame once the
  next frame is laid out).
- A code search shows that every structural mutator calls `markGraphDirty` and that no
  code reads `geometryDirty` for the summary.

### Phase Summary
_(write when phase completes)_

## Phase 3: Immediate redraw for late status changes
Status: Not started

- [ ] New `app/gui/dialogs/zoneEditorStatusKey.go` with the unexported comparable value
      type `zoneEditorStatusKey` `{hint string; addMode, addZoneMode bool; zoneCount,
      connectionCount int; graphDirty bool}` (its own file per AGENTS §4.1). Add a
      `shownStatus zoneEditorStatusKey` field to `zoneEditorGraphState`.
- [ ] A private `statusKey()` method on the dialog builds the key from the **full** dialog
      state, whichever branch the status line draws.
- [ ] `layoutStatus` stores `shownStatus = this.statusKey()` unconditionally, after any
      rebuild and before branching. This guarantees convergence: an unchanged state always
      matches, even in the error branch or while a hint hides the counts.
- [ ] At the end of `Body`, on the normal return only (the Apply and Cancel returns close
      the dialog), compare `this.statusKey()` with `shownStatus`. If they differ, call
      `gtx.Execute(op.InvalidateCmd{})`. This captures canvas presses and releases (add or
      delete connection, add zone, add-mode exits, hints), zone and connection deletes from
      the canvas, and side-panel quality writeback. No allocation: string comparison of
      `hint` plus scalar fields.
- [ ] Toolbar-button mutations run before the toolbar layout, so they do not trigger a
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
_(write when phase completes)_

## Phase 4: Tests
Status: Not started

Unit tests (untagged, mirrored layout, triple-A, `t.Parallel()`, gofakeit where it fits):

- [ ] `test/unit/internal/services/connection_editor/connectionEditorService/findIsolatedZones_test.go`, add:
  - no zones, returns `nil` (`assert.Nil`, locking the contract before Phase 5);
  - no connections, returns every zone in input order;
  - a zone referenced only as `To`, is not isolated;
  - a connection to an unknown zone, does not un-isolate a known zone;
  - duplicate zone names that are both unreferenced, both reported (the current behaviour,
    locked before any rewrite);
  - a self-loop connection, the zone is not isolated.
- [ ] `test/unit/internal/handlers/zoneEditorHandler/describeZoneEditorGraph_test.go`, add:
  - a `nil` isolated result gives count 0.
  (Argument forwarding is already locked by the exact-argument `On(...)` expectations.)

GUI integration tests (tagged `integration_test && gui`, because they use test exports
and the window):

- [ ] Test helper `test/test_helpers/integration_common/graphDescriptionCounter.go`: one
      struct, `GraphDescriptionCounter`, which embeds `handler_interfaces.IGuiHandler`,
      overrides `DescribeZoneEditorGraph` to count calls and delegate, and exposes
      `Calls() int` and `Last() dtos.ZoneEditorGraphDto` (the last result). Both are
      guarded by a mutex, because in headed mode the render goroutine writes them while
      the test reads them. Tag it `//go:build integration_test` for consistency with every
      other `integration_common` file (testlayoutcheck only inspects `_test.go` files).
- [ ] `AppRunner`: add `NewAppRunnerWithGuiHandler(tb, gui handler_interfaces.IGuiHandler)`.
      It and `NewAppRunnerWithFileSystem` share one private constructor taking both
      handlers, rather than duplicating the body.
- [ ] New `test/integration/gui/zoneEditorGraphCache_integration_test.go`, driven through a
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
- [ ] Invalidate tests, only if the Phase 3 spike succeeded. They run **dialog-direct**
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
- [ ] No new goldens are expected. If a plain run produces `.failure` files, inspect each
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
_(write when phase completes)_

## Phase 5: Optional O(zones + connections) isolation rewrite (G3)
Status: Not started

- [ ] Implement the candidate in `FindIsolatedZones`: build a
      `map[string]struct{}` of every `From`/`To` sized `2*len(connections)`, then collect
      unreferenced zones in input order. Keep the `nil` result for no isolated zones.
- [ ] Run `BenchmarkConnectionEditorService_FindIsolatedZones` with `-count=6` against the
      Phase 1 medians.
- [ ] Apply the G3 rule: lower ns/op at 24 and 40 zones, and within ±5% at 4 and 12 zones.
  - If it passes, keep it, and the Phase 4 unit tests must still pass unchanged.
  - If it fails, revert the function to its Phase 1 text exactly and record the numbers.

### Verification Plan
- The unit tests for `FindIsolatedZones` and `DescribeZoneEditorGraph` pass either way.
- The decision and the medians are recorded in the Phase Summary.

### Phase Summary
_(write when phase completes)_

## Phase 6: Final verification and review
Status: Not started

- [ ] Rerun the full Phase 1 baseline set and all three benchmarks (`-count=6`), and
      record before/after medians per size.
- [ ] `gofmt -l` clean, lint 0 issues, and `testlayoutcheck` passes.
- [ ] Independent implementation review by Claude Opus 5.5, with its findings resolved or
      explicitly declined with a reason.
- [ ] Mark §3.1 FIXED in place in the review, with a Progress paragraph, and update the §9
      `G` row and the progress line (17 fixed, 15 remaining). Do not renumber, and do not
      edit §8.
- [ ] Update the handoff for the next session.

### Verification Plan
- Every Phase 1 command passes, and coverage is ≥ 74.6% (or any drop is explained and
  owner-approved).
- The acceptance criterion G2 is demonstrated by the Phase 4 call-count tests.

### Phase Summary
_(write when phase completes)_

## Final Recap
_(write when all phases complete: summary of the entire piece of work)_

## Deployment Plan
_(write when all phases complete: step-by-step deployment instructions)_
