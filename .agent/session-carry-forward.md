# Session Carry-Forward — Batch 5 (Performance)

## 1. Session goal

Work through the senior code review at
[todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md) in PR-sized batches,
pausing after each batch for the owner to review. This document covers
**Batch 5 — the performance PR (§4.1)**. Batches 1–4 are already committed.

## 2. Fixes applied

- **§4.1 🟠 — the live preview rebuilt the full topology layout on every frame.**
  It now rebuilds only when its inputs change. At the shipped default topology
  (Random) with 8 players and 16 neutral zones the panel was burning ~2.5 ms and
  ~400 allocations per frame, on idle repaints and even while a modal covered
  the preview entirely.

## 3. Features added / changed

- New `models.PreviewLayoutCache` in
  [app/gui/models/previewLayoutCache.go](../app/gui/models/previewLayoutCache.go).
  `Get(templateRevision, topology, canvasSide, build)` returns the memoized
  layout and calls `build` only when the key differs. A failed build is **not**
  cached, so the next frame retries.
- The key is an unexported struct in its own file
  ([previewLayoutCacheKey.go](../app/gui/models/previewLayoutCacheKey.go)), so
  callers pass the three values rather than constructing a data type.
- `drivers.State` gained `templateRevision uint64` and `GetTemplateRevision()`.
  **All three** writers of `lastTemplate` now go through one private
  `setLastTemplate`, chosen by the owner over bumping the counter inline, so it
  cannot drift from the template on screen.
- `PreviewPanel` owns a cache instance and wraps its `BuildPreviewLayout` call
  in the `build` closure. No constructor signature changed, so wire is untouched.

### Behaviour that is deliberately unchanged

`BuildPreviewLayout` is deterministic (no RNG anywhere in `preview_service`), so
caching cannot alter what is drawn. The GPU snapshot suite confirms this.

## 4. File modifications

**Created**

- [app/gui/models/previewLayoutCache.go](../app/gui/models/previewLayoutCache.go)
- [app/gui/models/previewLayoutCacheKey.go](../app/gui/models/previewLayoutCacheKey.go)
- [test/unit/app/gui/models/previewLayoutCache/get_test.go](../test/unit/app/gui/models/previewLayoutCache/get_test.go)
- [test/unit/app/gui/models/previewLayoutCache/newPreviewLayoutCache_test.go](../test/unit/app/gui/models/previewLayoutCache/newPreviewLayoutCache_test.go)
- [test/unit/app/gui/drivers/state/getTemplateRevision_test.go](../test/unit/app/gui/drivers/state/getTemplateRevision_test.go)
- [test/performance/preview_layout_test.go](../test/performance/preview_layout_test.go)

**Edited**

- [app/gui/drivers/state.go](../app/gui/drivers/state.go) — revision field,
  `GetTemplateRevision`, `setLastTemplate`.
- [app/gui/drivers/stateGeneration.go](../app/gui/drivers/stateGeneration.go) —
  `applyGeneratedTemplate` and `clearGeneratedState` route through the setter.
- [app/gui/drivers/stateManualEdits.go](../app/gui/drivers/stateManualEdits.go) —
  `handleUpdateTemplate` routes through the setter.
- [app/gui/panels/previewPanel.go](../app/gui/panels/previewPanel.go) — owns and
  consults the cache.
- [todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md) — §4.1 marked
  `✅ FIXED` with the three deviations recorded; §12 item 5 ticked.

**Deleted** — none.

## 5. Tests added or updated

- `previewLayoutCache/get_test.go` — one build across N identical keys; one
  rebuild per key change (revision, topology, canvas side, three separate
  tests); the cached layout is returned on a hit; the rebuilt layout is returned
  on a miss; a failing build surfaces its error and is not cached.
- `previewLayoutCache/newPreviewLayoutCache_test.go` — constructor returns a
  usable cache.
- `state/getTemplateRevision_test.go` — revision starts at zero and advances on
  generation, regeneration, manual edits and reset.
- `test/performance/preview_layout_test.go` — the benchmark from the review,
  reinstated permanently and **untagged** (see §7).

**Verification results**

| Check | Result |
| --- | --- |
| `go build ./...` | pass |
| `go vet -tags=integration_test ./...` | pass |
| `go test -count=1 ./test/unit/...` | exit 0 |
| `go test -tags=integration_test ./test/integration/...` | exit 0 |
| `go test -tags='integration_test,gui' ./test/integration/gui/...` | exit 0 (snapshots unchanged) |
| Unit coverage | **65.0%** — unchanged, no drop |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | **42 issues — the pre-existing baseline** (2 `dupl`, 40 `gochecknoglobals`) |

Benchmark baseline recorded at `-benchtime=50x`, canvas side 600, Intel Core
Ultra 7 165H — Random/8p/16n at **2 553 776 ns/op, 41 264 B/op, 397 allocs/op**.
The benchmark measures the *uncached service* on purpose, so the numbers stay
comparable as the cache evolves.

## 6. Git status snapshot

Branch `AD/refactoring-07-21`. Nothing staged by the agent (AGENTS.md §2.5).

```
 D .agent/session-carry-forward.md   (rewritten below)
 M app/gui/drivers/state.go
 M app/gui/drivers/stateGeneration.go
 M app/gui/drivers/stateManualEdits.go
 M app/gui/panels/previewPanel.go
 M todo/review-opus5-08-04.md
?? app/gui/models/previewLayoutCache.go
?? app/gui/models/previewLayoutCacheKey.go
?? test/performance/preview_layout_test.go
?? test/unit/app/gui/drivers/state/getTemplateRevision_test.go
?? test/unit/app/gui/models/previewLayoutCache/
```

Batches 1–4 are already committed by the owner.

## 7. Rejections / things declined

- **Template pointer identity** as the cache key was declined in favour of the
  explicit revision counter.
- **Bumping the counter inline at each of the three assignment sites** was
  declined in favour of a single private setter.
- **`internal/services/preview_service/`** was the owner's first choice for the
  cache but is *architecturally forbidden*:
  [dependency_test.go](../test/unit/architecture/dependency/dependency_test.go)
  fails any import from `app/**` into `internal/services`. Do not re-attempt it
  without also changing that allow-list, which exists so the GUI talks only to
  handlers.
- **`app/gui/panels/`** (the review's own fallback, and the owner's second
  choice) was abandoned for a non-obvious reason worth remembering: that package
  had never been part of a unit-test binary, so adding a single test for it
  pulled **284 untested GUI functions** into the `-coverpkg` denominator and
  dropped the reported total from 65.0% to **60.1%**. CI has a hard
  *"Fail if coverage drops"* gate plus a 60.0% floor
  ([pr-validation.yml](../.github/workflows/pr-validation.yml)), so that
  placement would have failed the build. `app/gui/models/` is already covered by
  the `drivers` tests, so the total is unaffected.
  **General lesson: adding the first test that imports a previously untested
  package can lower total coverage even though nothing regressed.**
- **Caching the zone-editor canvas** was requested but turned out to be
  unnecessary:
  [zoneEditorCanvas.go](../app/gui/dialogs/zoneEditorCanvas.go) already guards
  the identical call behind `geometryDirty`/`geometrySide`, raised by every
  mutator. That guard is finer-grained than a layout key — it also catches zone
  drags, which the key cannot express — so it was left alone.
- **Tagging the reinstated benchmark `integration_test`**, as §4.1 instructed,
  was declined and the review corrected. AGENTS.md §4.6.1 restricts that tag to
  files consuming `*_testexports.go`; the benchmark uses production APIs and
  needs no GPU.
- **The zone-size "1.5 vs 2.0" discrepancy** reported during Batch 3 is *not* a
  bug and must not be "fixed": `MultiplierFormatter` has base 0.5 and scale 1.5
  over a slider value in `[0, 1]`, so both the label and persistence yield
  `[0.5, 2.0]`.

## 8. Open questions

None for Batch 5. Blocking questions still owed for later batches:

- **§7.1** (Batch 8): are direct pushes to `master` intentional, or should the
  branch be protected?
- **§9.1** (Batch 9): is any of `internal/` intended as a public API? Decides
  how much package documentation is written.
- **§2.7** (Batch 12): finish or remove the gladiator-arena preview (6 dead
  PNGs, 2 dead enum values). **If removed, drop its two validator entries
  `gladiatorArenaDaysDelayStart` and `gladiatorArenaCountDay` added in Batch 4.**
- **§1.8** (Batch 12): persist the output directory (a) in `.gen.json` or
  (b) machine-local. The review recommends (b).
- **§2.2** (Batch 13): the owner must confirm the scope of extracting
  regeneration policy out of `app/gui/drivers/`.

## 9. Next recommended actions

1. Owner reviews Batch 5 and commits.
2. **Batch 6 — DI PR.** §2.3 (`NewMandatoryContentItemMapper` builds its own
   `ContentRuleService`) + §2.4 (`NewTopologyBase` builds its own collaborators).
   Parameterise both constructors, register providers, regenerate wire.
3. **Batch 7 — Test-policy PR.** §6.3 *then* §6.5. Order matters or CI goes red.
4. **Batch 8 — CI/security posture.** §7.2, §7.3, §7.4, §8.3. ⚠ §7.1 blocked.
5. **Batch 9 — Docs.** §9.1–§9.6. ⚠ §9.1 blocked; §9.5 must follow §2.7. Also
   refresh `/memories/repo/` — the "four tabs" note is stale, there are three
   (General, Layout & Zones, Bonuses & Bans) — and AGENTS.md §1 wrongly says
   "single module" (§9.6).
6. **Batch 10 — Duplication cleanup.** §3.1, §3.3, §3.4 (verify via the GUI
   snapshot suite), then §3.2 + §5.3 together.
7. **Batch 11 — Coverage.** §6.2 (`internal/handlers` has no unit tests — start
   with `stateHandler`, `previewHandler`), §6.4.
8. **Batch 12 — Product decisions.** §2.7 and §1.8, both blocked on the owner.
9. **Batch 13 — Large refactors, plan first per AGENTS.md §4.7.** §2.1 → unblocks
   §2.5; then §2.2; §2.6 opportunistically.

Blockers recorded in review §12: *§6.5 after §6.3 · §2.5 after §2.1 · §5.1 folds
into §1.1 or §2.1 · §9.5 after §2.7 · §3.2 with §5.3.*

## 10. Carry-forward prompt

> Read `AGENTS.md` first and follow it strictly. In particular:
> **§2.1** never modify `data/`, `internal/entities/template/` or
> `internal/registry/` — they are authoritative game data.
> **§2.2** everything must build and run on Windows *and* Linux; use
> `path/filepath`, and chain PowerShell commands with `;`, never `&&`.
> **§2.3** every non-trivial change ships with tests, and total unit coverage
> must not drop.
> **§2.4** durable multi-session work gets a plan file under `plans/`.
> **§2.5** never stage and never commit — the owner reviews and commits himself.
>
> We are remediating the 46-finding review in `todo/review-opus5-08-04.md` in
> the PR-sized batches listed in its §12. Findings are marked `✅ FIXED` or
> `❌ WILL NOT FIX` **in place** inside that document — do not create a separate
> plan file for it.
>
> **Workflow for every batch, without exception:** ask the owner whether the
> batch should be done at all; if he declines, document *why it should not be
> attempted in the future* in the review file; ask every clarifying question
> **before** implementing; implement; rewrite
> `.agent/session-carry-forward.md`; then **stop and wait for his review**.
>
> Batches 1–5 are done (security, correctness, durability, input validation,
> performance). Batch 6 is the DI PR: §2.3 and §2.4, then regenerate wire.
>
> Read `.agent/session-carry-forward.md` for the full handoff, including the
> open questions that block Batches 8, 9, 12 and 13, and the verification
> command set.
