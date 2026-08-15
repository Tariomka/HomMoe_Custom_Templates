# Session carry-forward — 2026-08-15

## 1. Session goal

Implement **batch L** of `todo/backlog-opus5.md` — §5.4 items (a)–(c) (handler
file hygiene, a single home for the coordinate literals, and a narrowed snapshot
mask) plus §5.5 steps 1–2 (diagnose why GUI snapshots differ between local and
CI, and replace the tolerance that was hiding it) — before batch F.

## 2. Fixes applied

- **A real production font bug, found while diagnosing §5.5.**
  [theme.go](../app/gui/themes/theme.go) built its shaper without
  `text.NoSystemFonts()`, so any glyph missing from the bundled Go collection was
  resolved through whatever face the operating system happened to offer. `◆`
  (U+25C6), used in **every** section header, is one such glyph. On CI the
  substitute face is 1 px shorter and 4 px narrower, and because the shortfall
  accumulates down the panel, rows near the bottom were displaced by 3 px.
  Proof, from a per-band best-offset search on snapshot 1: the checkbox rows at
  y 430–580 scored **1.2957 %** at zero offset and **0.0029 % at dy = −3**.
  This also **disproves** the backlog's standing hypothesis that CI was capturing
  a half-rendered frame — no fix belongs in `captureScreenshot`.
- **Seven glyph substitutions** to characters the bundled fonts actually contain:
  `◆`→`♦` in [sectionWidget.go](../app/gui/widgets/sectionWidget.go#L35);
  `✕`→`×` in [bonusPickerDialog.go](../app/gui/dialogs/bonusPickerDialog.go#L355),
  [bonusesPanel.go](../app/gui/panels/bonusesPanel.go#L208) and
  [dialogHost.go](../app/gui/drivers/dialogHost.go#L144);
  `↺`→`←` in [zoneContentDialog.go](../app/gui/dialogs/zoneContentDialog.go#L94);
  `⚠`→`‼` in [zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go#L242)
  and twice in [stateManualEdits.go](../app/gui/drivers/stateManualEdits.go#L82).
- **The snapshot gate could not see a layout break.** The single mean-distance
  comparison at a 2 % threshold scored the three broken CI frames at 0.66 %,
  1.22 % and 0.66 % — comfortably passing while 3.4 % of the window was
  misplaced. A mean over 1.44 million pixels cannot distinguish "everything is
  slightly dimmer" from "a third of the frame moved".

## 3. Features added / changed

### Determinism

- `text.NoSystemFonts()` in `NewTheme`, so the UI renders identically on every
  machine rather than borrowing OS faces with different metrics.
- **Regression guard**: `TestWhenAppSourcesContainNonAsciiRunes_AllAreCoveredByTheBundledFonts`
  parses every file under `app/` with `go/parser`, walks STRING and CHAR
  `BasicLit` nodes, and fails on any rune `gofont.Collection()` cannot render.
  Verified by mutation — restoring `◆` fails with the exact file, line and column.

### Snapshot measurement (§5.5 step 2)

- [comparer.go](../test/test_helpers/integration_common/snapshot/comparer.go)
  now returns a
  [Difference](../test/test_helpers/integration_common/snapshot/difference.go)
  carrying **both** `MeanDistance` and `ChangedPixelFraction`, judged against
  three named constants — `DefaultMeanThreshold = 0.0025`,
  `DefaultPixelTolerance = 64`, `DefaultChangedPixelThreshold = 0.0005` —
  instead of one unexplained `0.02`. A pixel counts as changed when the largest
  of its three channel deltas exceeds the tolerance.
- The broken CI frames score **1.38 %, 2.83 % and 1.38 %** on the new fraction
  gate against its 0.05 % limit — i.e. it trips by 27–57×, where the mean gate
  passed.
- `Describe` names both measurements and both limits in the failure message, so
  a future CI artifact can be read without re-deriving anything.

### Harness (§5.4 a–c)

- `runnerHandler.go` → [baseHandler.go](../test/test_helpers/integration_common/baseHandler.go),
  and the 46-line first-person design essay replaced by a four-line pointer to
  backlog §5.4 d–f.
- All literals moved to
  [handlerCoordinates.go](../test/test_helpers/integration_common/handlerCoordinates.go),
  with a doc block recording **how each was measured** (tabs via
  `utils.ButtonPositionLogger`; masks via a two-unmasked-run diff).
- The one full-height 470 px mask replaced by three measured rectangles.
  **Masked area 423 000 → 208 000 px**, putting the canvas border, the legend,
  the template name, Browse/Reveal and the **Generate** and **Save Template**
  buttons back under test.

### Adjacent cleanups (owner-requested extras)

- `constants.GetHubSpokeConnectionNameFor(label)` — the three
  `HubZonePrefix + label` **connection** names moved out of
  [hubTopology.go](../internal/services/template_generator/providers/topology/hubTopology.go).
  It deliberately emits the same string as `GetHubZoneNameFor`; only the intent
  differs, and changing the string would change generated templates.
- `zone_helpers.IsClusterHubZoneName` (prefix match, per-cluster hubs) and
  `IsSharedHubZoneName` (exact `"Hub"`), used by
  [layoutRingHub.go](../internal/services/preview_service/layoutRingHub.go).
  The pre-existing `IsZoneNameHub` conflates the two and would have changed the
  preview layout.

## 4. File modifications

All of the below landed in commit **`00aaf9c` "Batch L tryout"** (51 files,
+1429/−162), staged and committed **by the owner**, and pushed.

**Production**

| File | Change |
| --- | --- |
| [theme.go](../app/gui/themes/theme.go) | `text.NoSystemFonts()` + doc comment explaining why |
| [sectionWidget.go](../app/gui/widgets/sectionWidget.go) | `◆` → `♦` |
| [bonusPickerDialog.go](../app/gui/dialogs/bonusPickerDialog.go), [bonusesPanel.go](../app/gui/panels/bonusesPanel.go), [dialogHost.go](../app/gui/drivers/dialogHost.go) | `✕` → `×` |
| [zoneContentDialog.go](../app/gui/dialogs/zoneContentDialog.go) | `↺` → `←` |
| [zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go) | `⚠` → `‼`, doc comment updated |
| [stateManualEdits.go](../app/gui/drivers/stateManualEdits.go) | two `⚠ Error: %v` → `‼ Error: %v` |
| [connectionNames.go](../internal/common/constants/connectionNames.go) | + `GetHubSpokeConnectionNameFor` |
| [hubTopology.go](../internal/services/template_generator/providers/topology/hubTopology.go) | three sites use the new builder |
| [zoneNameType.go](../internal/helpers/zone_helpers/zoneNameType.go) | + `IsClusterHubZoneName`, `IsSharedHubZoneName` |
| [layoutRingHub.go](../internal/services/preview_service/layoutRingHub.go) | uses the pair; dropped `strings`/`constants` imports |
| [fileExplorerDialogToolbar.go](../app/gui/dialogs/fileExplorerDialogToolbar.go) | `"Up"` → `"← Back"` — **the owner's own edit**, confirmed "Mine - leave it" |

**Test harness**

- **New** `baseHandler.go`, `handlerCoordinates.go`,
  `snapshot/difference.go`.
- **Deleted** `runnerHandler.go` (via `Remove-Item`, never `git rm`).
- **Rewritten** `snapshot/comparer.go`; **modified** `appRunnerSnapshots.go`.
- **Regenerated** all four goldens under
  `test/test_helpers/integration_common/snapshot/__snapshots__/window_snapshot_integration_test/`,
  **locally on a real GPU** — never in CI.

**Plan** — [plans/gui-test-harness-groundwork.md](../plans/gui-test-harness-groundwork.md),
six phases with summaries, Final Recap and Deployment Plan. Its Phase 4 update
(recording the green CI run) is the **only uncommitted change** in the tree.

## 5. Tests added or updated

- **New** `test/unit/app/gui/themes/theme/newTheme_test.go` — the AST glyph guard.
- **New** `test/unit/internal/common/constants/connectionNames/` — 18 files, one
  per builder, closing an AGENTS §4.6 gap (17 builders had no attributable test
  folder).
- **New** `isClusterHubZoneName_test.go`, `isSharedHubZoneName_test.go`.
- **New** `snapshot/comparer/describe_test.go`, `snapshot/difference/string_test.go`.
- **Updated** `snapshot/comparer/{compare,matches,newComparer}_test.go` for the
  two-gate design.

**Last full gate — all green:**

| check | result |
| --- | --- |
| `go build ./...` | clean |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | exit 0 |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go test ./test/unit/... -count=1` | pass |
| `go test -tags=integration_test ./test/integration/... -count=1` | pass |
| `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | pass |
| `golangci-lint-v2 run ./...` | **0 issues** |
| coverage | **72.9 %** (floor 72.5 %, unchanged) |
| **owner's CI pipeline** | **green — no `gui-snapshot-failures` artifact** |

## 6. Git status snapshot

Branch **`AD/fixing_some_stuff_08-12`**, in sync with
`origin/AD/fixing_some_stuff_08-12`.

```
 M plans/gui-test-harness-groundwork.md
```

That is the entire inheritance: the Phase 4 completion note written after CI went
green. Batch L itself is committed (`00aaf9c`). The owner separately committed
`ada8e39 "package updates"` — Go module bumps and workflow `setup-steps` changes,
**not part of this work**.

Recent history: `ada8e39` package updates → `00aaf9c` Batch L tryout →
`f908a72` Reset CI test changes → `892c161` Batch D done → … → `9d1b83b`
(`origin/master`).

## 7. Rejections / things the user declined

- **Running CI.** Standing instruction: *"Never run anything in CI yourself. Ask
  me to do it instead."* The owner supplied the failure artifacts and ran the
  pipeline.
- **CI as the reference.** *"never make CI the reference (do not generate goldens
  in CI)"* — CI is a software renderer; goldens come from a real GPU only.
- **§5.4 (a) half-done on purpose.** The backlog also asks that `NewHandler` stop
  returning an unexported `*baseHandler` from an exported function. Owner:
  *"Leave the return type; just rename the file."* Deferred until (d) introduces
  a handler contract to return.
- **The half-rendered-frame theory** in backlog §5.5 was rejected by evidence,
  not by preference — see §2 above.
- **Widening the mean gate** is explicitly off the table as a way to make CI
  green; only the fraction floor may be pinned, and only to a measured value.

## 8. Open questions

- **`todo/backlog-opus5.md` has not been updated.** §5.4 still describes
  `runnerHandler.go`, the 46-line comment and the 470 px mask; §5.5 still reads
  as open and still carries the disproved half-rendered-frame hypothesis. Items
  (a)–(c) and §5.5 steps 1–2 are done. Someone should decide whether the backlog
  is amended or left as a historical record.
- **The `air` terminal exited 1** and `.vscode/settings.json` showed as modified
  during the session. Neither belongs to this work; unexamined.
- **§5.4 (e) needs a decision when it lands**: *Allow non-official larger map
  sizes* expands the Map Size dropdown inline, shifting everything below it. The
  measured constants in `handlerCoordinates.go` assume the collapsed layout.

## 9. Next recommended actions

1. **Batch L is closed.** Commit the one outstanding plan edit, or leave it —
   nothing depends on it.
2. **Backlog §5.4 items (d)–(f)** — the next batch, in this order:
   (d) per-tab and per-dialog handlers (tab handlers embed `baseHandler`; dialog
   handlers must **not**, because a dialog disables the background);
   (e) handler state tracking for layout-shifting toggles;
   (f) `Scroll(point, delta)` on `AppRunner`, injecting `pointer.Scroll` through
   the same router as `ClickAt` — a hard blocker for §5.2.
   Respect the backlog's own scope guard: grow these one method at a time as
   §5.1–§5.3 need them. An unused handler method is untested test code.
3. **Then batch F**, per the backlog's §8 ordering.
4. Any new snapshot work: regenerate goldens **locally only**, and read the new
   two-part failure message before touching a threshold.

## 10. Carry-forward prompt

> **Read `AGENTS.md` first — it governs everything below.**
>
> Hard rules, one line each: never modify `data/`, `internal/entities/template/`
> or `internal/registry/` without my explicit approval; keep every change
> cross-platform (Windows + Linux, `path/filepath`, PowerShell chained with `;`
> never `&&`); every change ships with tests and must not drop unit coverage
> below 72.5 % (currently 72.9 %); durable multi-session work gets a plan file
> under `plans/`; **never stage and never commit** — I review and commit, and you
> delete files with `Remove-Item`, never `git rm`; never change where `.rmg.json`
> is written and never persist the output directory; never run a bulk in-place
> rewrite over the repository. Additionally: **never run CI yourself** — ask me —
> and **never generate snapshot goldens in CI**, which is a software renderer and
> must not become the reference.
>
> Where things stand: **batch L is complete and CI is green.** It fixed a real
> production font bug (`themes.NewTheme` lacked `text.NoSystemFonts()`, so `◆`
> fell back to an OS face and shifted the layout), replaced the single 2 % mean
> snapshot gate with a two-gate mean + changed-pixel-fraction comparer, renamed
> `runnerHandler.go` to `baseHandler.go` with the coordinates extracted to
> `handlerCoordinates.go`, and narrowed the snapshot mask from 423 k to 208 k px.
> Full detail — including the per-band measurements that proved the diagnosis —
> is in `plans/gui-test-harness-groundwork.md`.
>
> Next up is backlog `todo/backlog-opus5.md` §5.4 items **(d)–(f)**, then batch
> **F**. Note that §5.4 and §5.5 in the backlog have **not** been updated to
> reflect batch L, so parts of them are stale.
>
> Before starting, prompt me to confirm the item(s) and surface every open
> question first. The full handoff is in
> `./.agent/session-carry-forward.md` — read it before asking.
