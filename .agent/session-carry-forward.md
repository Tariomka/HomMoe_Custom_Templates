# Session carry-forward — 2026-09-03

## 1. Session goal

Execute **phase 5**, the last phase of
[.agent/plans/batch-j-zone-tier-source-of-truth.md](plans/batch-j-zone-tier-source-of-truth.md):
shrink `entityNamerAllowList`, prove the shrink is real, then write the batch's
durable records (backlog §2.2/§2.6/§8, the plan's Final Recap and Deployment
Plan).

**Batch J is now complete.** All five phases are done and the plan carries its
Final Recap and Deployment Plan.

## 2. Fixes applied

None. Phase 5 is a measurement and a records pass — no production code changed.

## 3. Features added / changed

None. The one code change is a **tightened architecture constraint**:
[test/unit/architecture/dependency/layering_test.go](../test/unit/architecture/dependency/layering_test.go)
no longer allows `app/gui/editor` or `internal/services/preview_service` to name
an `internal/entities` type. Both packages had already stopped doing so —
`preview_service` because phase 4 moved it onto `template_model` end to end — so
this makes an existing fact enforceable rather than changing behaviour.

## 4. File modifications

| File | Change |
| --- | --- |
| [test/unit/architecture/dependency/layering_test.go](../test/unit/architecture/dependency/layering_test.go) | Removed two `entityNamerAllowList` entries; added a comment naming the two **permanent** entries (`file_service`, `template_generator`) so a future agent does not try to clean them. |
| [.agent/backlog/backlog-opus5.md](backlog/backlog-opus5.md) | §2.2 → ✅ DONE with the full record; §2.6 retitled and recounted (113 files / 23 packages → **84 / 21**), steps 3–4 rewritten; §8 row **J** filled in; coverage refreshed **73.8 % → 74.3 %** in all three places (header baseline, §8 coverage note, §9 gate table); done-count 14 → **15**, open 7 → **6**. |
| [.agent/plans/batch-j-zone-tier-source-of-truth.md](plans/batch-j-zone-tier-source-of-truth.md) | Phase 5 → Complete with checklist, Phase Summary and verification table; **Final Recap** and **Deployment Plan** written; the header's stale 73.9 % coverage figure corrected. |
| `.agent/session-carry-forward.md` | This file (rewritten). |

Two files were **temporarily mutated and reverted** for the proof —
`internal/services/preview_service/previewLayoutService.go` and
`app/gui/editor/window.go`. `git status` confirms neither is modified.

## 5. Tests added or updated

No new tests. The existing gate
`TestWhenEntityConsumersAreScanned_OnlyPermittedPackagesNameAnEntity` now covers
two more packages.

Last full run, all green:

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | clean |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 |
| `go test ./test/unit/... -count=1` (with coverage) | exit 0 — **74.3 %** |
| `go test ./test/... -count=1` | exit 0 |
| `go test -tags=integration_test ./test/integration/... -count=1` | exit 0 |
| `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | exit 0, **no `-update`**, no golden modified |
| `golangci-lint-v2 run ./...` | **0 issues** |

## 6. Git status snapshot

Branch **`AD/fixing_some_stuff_08-12`**, head `8aa5f26 Batch J wip 5`. Phases
1–4 are committed — the carry-forward this session started from was stale in
saying phase 4 was uncommitted.

```
 M .agent/backlog/backlog-opus5.md
 M .agent/plans/batch-j-zone-tier-source-of-truth.md
 M test/unit/architecture/dependency/layering_test.go
```

**Nothing staged.** Zero `*.golden`. The author stages and commits (AGENTS §2.5).

## 7. Rejections / things not done

- **`internal/dtos` was not removed from the allow-list.** All five files name
  `entities.Connection`; phase 3 left connections on the entity on purpose,
  because a connection carries no tier. Removing it would have meant moving
  connections for the sake of a shorter list.
- **`app/gui/dialogs`, `app/gui/drivers`, `app/gui/models` were not removed** —
  same reason: their zones are `template_model.Zone` now, their connections and
  templates are not.
- **`internal/services/file_service` and the `template_generator` tree were not
  touched, permanently.** `file_service` writes `.rmg.json`; the generator
  assembles the entity that
  `TestWhenDefaultConfiguration_ReturnsGoldenTemplate` compares. Both are now
  recorded as exempt-by-decision in backlog §2.6 step 4 **and** in the
  allow-list comment.
- **No sweep, no bulk rewrite, no staging.** Only removals from the allow-list,
  as the plan requires.

## 8. Open questions

- **Backlog §2.6 step 2 is a decision, not work**: does the `.rmg.json`
  vocabulary (`entities.Connection`, `entities.RmgTemplate`) earn a documented
  carve-out the way `internal/helpers/data` has one, or is naming it below the
  repositories a genuine breach? Everything left on the entity allow-list except
  the two permanent seams hinges on that answer.
- **Two benchmark baselines disagree in the record** — phase 2 measured
  `BenchmarkEditorWindow_TabCycling` at ~5,699 allocs/op, backlog §1.4 records
  6,640 after the clone batch. Neither is a batch J regression; the numbers were
  taken on different trees and nobody has reconciled them.

## 9. Next recommended actions

1. **Owner reviews and commits batch J.** Follow the plan's Deployment Plan — in
   particular step 6, the in-app smoke test: generate, re-tier a neutral zone in
   the editor, Apply, save the `.gen.json`, reload, confirm the Quality dropdown
   still shows the chosen tier. That is the one path the batch exists to make
   reliable.
2. Delete the transient docs once it lands (owner's doc-lifecycle rule): the
   plan file and this carry-forward. Backlog §2.2 is self-contained and is the
   surviving record.
3. Take **backlog §2.6 step 2** as its own scoped decision before any further
   allow-list work.
4. Optional: reconcile the two `TabCycling` benchmark baselines.

## 10. Carry-forward prompt

> Read `AGENTS.md` first. The hard rules, one line each: never modify `data/`,
> `internal/registry/` or anything under `internal/entities/template/` —
> `internal/entities/editor_state/` is *not* protected; everything must build and
> run on Windows and Linux (`path/filepath`; chain PowerShell with `;`, never
> `&&`); every change ships with tests and unit coverage must not drop below
> 72.5 % (currently **74.3 %**), lint baseline **0 issues**; **never stage and
> never commit** — `Move-Item` not `git mv`, `Remove-Item` not `git rm`; never
> change where `.rmg.json` is written and never persist the output directory;
> never run a bulk in-place rewrite and **never round-trip a `.go` file through
> `Get-Content`/`Set-Content`** — use `gofmt -r` on an explicit file list and
> verify insertions == deletions per file.
>
> **Batch J (zone tier single source of truth, backlog §2.2) is COMPLETE** — all
> five phases, records written, every gate green, no golden moved. Its plan at
> `.agent/plans/batch-j-zone-tier-source-of-truth.md` carries the Final Recap and
> Deployment Plan; `.agent/backlog/backlog-opus5.md` §2.2 is the durable,
> self-contained record and survives the plan's deletion. Only phase 5 is
> uncommitted (three modified files, nothing staged); phases 1–4 are committed
> through `8aa5f26` on `AD/fixing_some_stuff_08-12`.
>
> Standing traps this codebase punishes: **nil is load-bearing** three times over
> — nil `Previous` = first generation, nil `Next` = unarmed debounce, nil
> `Zone.Quality` = "infer it"; the persisted tier is `*int8` because `omitempty`
> on a plain `int8` would silently drop every Plastic zone (ordinal 0), guarded
> by a mutation-verified test in
> `test/integration/manualZoneTierPersistence_integration_test.go`; the two
> frozen fixtures under `test/test_helpers/testdata/` and the untagged
> `editorStateWireFormat_integration_test.go` must keep passing unchanged and
> compare **parsed objects, never bytes**; `cmd/testlayoutcheck` matches
> test-only export names tree-wide, so grep any new accessor name first; a file
> gets `//go:build integration_test` **only** if it calls a `*_testexports.go`
> accessor; `helpers.MapSlice` / `helpers.MapPointer` preserve nil-vs-empty where
> `linq.SelectSlice` does not; and `golangci-lint --fix` wraps a long signature
> as `param,\n) Ret {` where the house style is `param) Ret {` — restyle by hand
> after a `--fix`.
>
> Next up is **backlog §2.6 step 2**, which is a *decision* and not a sweep:
> whether `internal/dtos` and `internal/handlers` naming `entities.Connection` /
> `entities.RmgTemplate` is a breach at all, or whether the `.rmg.json`
> vocabulary deserves a documented carve-out. Note that two entries on
> `entityNamerAllowList` are **permanent by decision** —
> `internal/services/file_service` (it writes `.rmg.json`) and
> `internal/services/template_generator` plus its topology tree (it assembles the
> entity the golden-template test compares). Do not try to clean either.
>
> Full handoff in `./.agent/session-carry-forward.md`.
