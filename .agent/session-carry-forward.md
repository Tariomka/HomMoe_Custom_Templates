# Session carry-forward — GUI handler framework (backlog batch M)

> The previous session's handoff (batch L) is archived beside this file as
> [session-carry-forward-batch-L.md](session-carry-forward-batch-L.md).

## 1. Session goal

Build backlog [§5.4](../todo/backlog-opus5.md) items **(d)–(f)** — the typed
per-tab / per-dialog GUI handler tree, the `AppRunner` scroll seam, and
layout-shift tracking — standalone and ahead of batch F.

## 2. Fixes applied

- **Clicks did not reach the editor state.** A click is processed on the frame it
  is queued against, but a panel only writes widget values back in `SaveToState`
  on the *following* layout. Added `BaseHandler.commit()` in
  [baseHandler.go](../test/test_helpers/integration_common/baseHandler.go) and
  called it at the end of every state-mutating handler method, always **before**
  `VerifySnapshot`.
- **Phase 2's goldens were invalidated** by that fix (it inserted a frame into
  `ToggleAdvancedZoneControl`); regenerated locally on a real GPU.
- **Dialog handlers dropped shift state.** They originally held a bare
  `*AppRunner`, so `Close()` had to fabricate a fresh `BaseHandler`. They now hold
  `base *BaseHandler` and return it.
- **A test expectation was wrong, not the code.** `SelectMapSize(12)` yields
  `272`, not `288` — the 11 official sizes occupy rows 0–10.

## 3. Features added / changed

- **`*BaseHandler` exported** with a fluent fan-out into three tab handlers
  (embedding `*BaseHandler`, so shift state is shared) and two dialog handlers
  (**not** embedding — a dialog's scrim absorbs background clicks, so promoting
  the tab clicks would be a lie).
- **Toolbar methods** `ClickNew` / `ClickLoad` / `ClickSaveAs`. **`Exit` is never
  reachable** (`State.Exit()` calls `os.Exit(0)`); `Save` is excluded because it
  writes into the real detected output directory.
- **`AppRunner.Scroll(point, delta)`** injecting a `pointer.Scroll` event through
  the router, plus `BaseHandler.ScrollPanel(delta)`.
- **Layout-shift tracking.** `isExperimentalMapSizes` is stored on `BaseHandler`
  because it decides whether the map-size dropdown has 11 rows or 28.
  `isSingleHero` needs **no** bookkeeping — it shifts right-column widgets only
  and no handler coordinate sits below them. This *corrects the backlog's
  framing*, which called it the larger shifter.
- **Real scroll position** read from `widget.List` through three new
  `*_testexports.go` accessors, never accumulated from injected deltas
  (`layout.List` clamps at both ends).
- **`ToggleAdvancedZoneControl()`** — one method beyond the agreed scope, added
  because no panel otherwise has a scroll range worth testing (~18 px; ~386 px
  with it). **Flagged for review.**
- **All coordinates measured**, not derived, and recorded with their measurement
  notes in
  [handlerCoordinates.go](../test/test_helpers/integration_common/handlerCoordinates.go).

## 4. File modifications

**Created**

| File | Summary |
| --- | --- |
| [plans/gui-handler-framework.md](../plans/gui-handler-framework.md) | The as-built design; 5 phases, all Complete, with Final Recap and Deployment Plan. |
| [generalTabHandler.go](../test/test_helpers/integration_common/generalTabHandler.go) | `ToggleExperimentalMapSizes`, `OpenMapSizeSelector`, `SelectMapSize`, `SelectGameMode`. |
| [layoutAndZonesTabHandler.go](../test/test_helpers/integration_common/layoutAndZonesTabHandler.go) | `OpenZoneEditor`, `ToggleAdvancedZoneControl`. |
| [bonusesAndBansTabHandler.go](../test/test_helpers/integration_common/bonusesAndBansTabHandler.go) | Type only — the tab has no interactions in scope. |
| [fileExplorerHandler.go](../test/test_helpers/integration_common/fileExplorerHandler.go) | Reachability-only: `IsOpen` / `Close`. |
| [zoneEditorHandler.go](../test/test_helpers/integration_common/zoneEditorHandler.go) | Reachability-only: `IsOpen` / `Close`. |
| [layoutPanel_testexports.go](../app/gui/panels/layoutPanel_testexports.go) | `ScrollPosition()`. |
| [generalPanel_testexports.go](../app/gui/panels/generalPanel_testexports.go) | `ScrollPosition()`. |
| [tab_testexports.go](../app/gui/drivers/tab_testexports.go) | `GetPanel()`. |
| [handlerDialogReachability_integration_test.go](../test/integration/handlerDialogReachability_integration_test.go) | 6 tests. |
| [handlerScroll_integration_test.go](../test/integration/handlerScroll_integration_test.go) | 2 tests. |
| [handlerGeneralTab_integration_test.go](../test/integration/handlerGeneralTab_integration_test.go) | 4 tests. |
| 6 `.golden` files | `ScrollMatchesGoldens_1..3`, `MapSizeShiftMatchesGoldens_1..3`. |

**Modified**

| File | Summary |
| --- | --- |
| [baseHandler.go](../test/test_helpers/integration_common/baseHandler.go) | Exported; tab fan-out, toolbar, `ScrollPanel`, `commit()`, shift flag. |
| [appRunner.go](../test/test_helpers/integration_common/appRunner.go) | Added `Scroll` and `SelectedPanelScrollPosition`; removed the temporary probe accessor. |
| [handlerCoordinates.go](../test/test_helpers/integration_common/handlerCoordinates.go) | ~20 measured constants plus the measurement procedure. |
| [window_testexports.go](../app/gui/editor/window_testexports.go) | `SelectedPanelScrollPosition` + a local `scrollablePanel` interface. |
| [window_snapshot_integration_test.go](../test/integration/gui/window_snapshot_integration_test.go) | Two new snapshot tests. |
| [backlog-opus5.md](../todo/backlog-opus5.md) | §5.4 and §5.5 marked done and annotated; §8 table updated for L and M. |

**Deleted:** two temporary probe/diagnostic test files (with `Remove-Item`).

## 5. Tests added or updated

12 integration tests + 2 snapshot tests (6 goldens). Last full run:

| Suite | Result |
| --- | --- |
| `go test ./test/unit/... -count=1` | pass |
| `go test -tags=integration_test ./test/integration/... -count=1` | pass |
| `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | pass |

All backlog §9 gates pass: build, both vets, `testlayoutcheck`, `gofmt -l` empty,
`golangci-lint-v2` **0 issues**, `wire diff` no diff, coverage **72.9 %**
(floor 72.5 %, unchanged as predicted — nothing in this batch is inside
`-coverpkg`).

## 6. Git status snapshot

Branch **`AD/fixing_some_stuff_08-12`**. Nothing was staged or committed by the
agent.

```
 M app/gui/editor/window_testexports.go
AM plans/gui-handler-framework.md
 M test/integration/gui/window_snapshot_integration_test.go
MM test/test_helpers/integration_common/appRunner.go
 M test/test_helpers/integration_common/baseHandler.go
M  test/test_helpers/integration_common/handlerCoordinates.go
 M todo/backlog-opus5.md
?? app/gui/drivers/tab_testexports.go
?? app/gui/panels/generalPanel_testexports.go
?? app/gui/panels/layoutPanel_testexports.go
?? test/integration/handlerDialogReachability_integration_test.go
?? test/integration/handlerGeneralTab_integration_test.go
?? test/integration/handlerScroll_integration_test.go
?? test/test_helpers/integration_common/bonusesAndBansTabHandler.go
?? test/test_helpers/integration_common/fileExplorerHandler.go
?? test/test_helpers/integration_common/generalTabHandler.go
?? test/test_helpers/integration_common/layoutAndZonesTabHandler.go
?? test/test_helpers/integration_common/zoneEditorHandler.go
?? .../__snapshots__/window_snapshot_integration_test/TestWindowSnapshots_MapSizeShiftMatchesGoldens_{1,2,3}.golden
?? .../__snapshots__/window_snapshot_integration_test/TestWindowSnapshots_ScrollMatchesGoldens_{1,2,3}.golden
```

**Inherited, not created by this session:** `plans/gui-handler-framework.md`,
`handlerCoordinates.go` and part of `appRunner.go` are already **staged** (`A`/`M`
in the index column). The agent did not stage them and, per AGENTS.md §2.5, did
not unstage them either.

## 7. Rejections / things the user declined

- **Generating goldens in CI** — rejected outright. CI is a software renderer
  (llvmpipe under `xvfb-run`) and must never become the reference. Backlog §5.5
  step 2 is struck through accordingly.
- **Running CI at all** — the agent never does; ask the owner.
- **An `IBaseHandler` interface** — the owner chose to export the struct instead.
- Two recommendations were overridden in favour of a larger scope: build (d)–(f)
  standalone rather than growing them from F, and cover all three tabs plus both
  dialogs. The owner also chose "Load + Save As + New" over the recommended
  "Load + Save As only".

## 8. Open questions

1. **The staged files listed above** — was that intentional? Nothing further will
   be staged or unstaged without an instruction.
2. **`ToggleAdvancedZoneControl` is one method beyond the agreed handler set.**
   Keep it (it is the only way to get a meaningful scroll range) or drop it and
   accept an 18 px scroll test?
3. **Empty leftover directory**
   `test/test_helpers/integration_common/snapshot/__snapshots__/runnerHandler/`
   from batch L's rename. Untracked and harmless; delete it or leave it?

## 9. Next recommended actions

1. Review this batch and commit it (the agent will not).
2. Eyeball the six new goldens — `MapSizeShiftMatchesGoldens_3` is the one that
   proves the layout shift (28 dropdown rows rendered inline). They are PNGs
   despite the `.golden` extension.
3. Answer the three open questions above.
4. Run CI.
5. Start batch **E** (§4.1, "Save As" → "Save To"), then batch **F** (§5.3,
   file-explorer pointer and hidden-file tests) — both written against the
   handler API this batch settled.

## 10. Carry-forward prompt

> **Read `AGENTS.md` first — it governs everything.** Its hard rules, one line
> each: never modify `data/`, `internal/entities/template/` or
> `internal/registry/` without my explicit approval; keep every change
> cross-platform (Windows + Linux, `path/filepath`, PowerShell chained with `;`
> never `&&`); every change ships with tests and must not drop unit coverage
> below **72.5 %** (currently 72.9 %); durable multi-session work gets a plan file
> under `plans/`; **never stage and never commit** — I review and commit, and you
> delete files with `Remove-Item`, never `git rm`; never change where `.rmg.json`
> is written and never persist the output directory; never run a bulk in-place
> rewrite over the repository. Additionally: **never run CI yourself** — ask me —
> and **never generate snapshot goldens in CI**, which is a software renderer and
> must not become the reference.
>
> Backlog batch **M** (`todo/backlog-opus5.md` §5.4 (d)–(g)) is **complete** and
> awaiting my review on branch `AD/fixing_some_stuff_08-12`; the as-built design,
> including everything you must not rediscover, is in
> `plans/gui-handler-framework.md`. Batch **L** (§5.4 (a)–(c), §5.5) is also done
> — see `plans/gui-test-harness-groundwork.md`. Do not redo either.
>
> The full handoff, including three unanswered questions and the working-tree
> state, is in `./.agent/session-carry-forward.md` — read it before you touch
> anything.
>
> Next up is batch **E** (§4.1, "Save As" → "Save To"), then batch **F** (§5.3).
> Before starting, prompt me to confirm the item and surface every open question
> first.
