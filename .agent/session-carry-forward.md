# Carry-forward: Batch H implemented, awaiting owner review and commit

Date: 2026-10-04. **Batch H (review §6.2, §6.3, §6.5: CI/tooling hardening) is
implemented, verified on Windows, and approved by independent review. It is uncommitted.**
The owner approved the plan (committed by the owner as `25b4c98` "H Plan") and asked for every
phase to run without stopping, with their review at the end.

**§6.2, §6.3 and §6.5 are marked FIXED ahead of the commit** in the
[surviving review](backlog/review-gpt-6-astra-09-07.md), at the owner's request. Their
Progress paragraphs say "uncommitted, awaiting owner commit". Row H of §9 is complete, and
the progress line reads **21 fixed, 11 remaining** (0 High, 6 Medium, 5 Low). If the owner
rejects any part, revert the matching review entries.

The plan, [batch-h-ci-tooling-hardening.md](plans/batch-h-ci-tooling-hardening.md), has all
four phases Complete, plus its Final Recap and Deployment Plan. It records the owner decisions
H1–H12. §3 below restates them, along with J1–J9 and G1–G9a.

**Section 8 is preserved verbatim**, contradictions and all. Its Batch C phase/engine
wording, its "pending decisions" wording for §1.5, §1.11, §1.12 and §2.1, its §2.2
bullet, and its tooling pending-decision wording (§6.2/§6.3/§6.5) are **superseded**. Everything
else §8 retains as later scope remains **binding**. Do not edit §8.

Batches A/D, B, C, E, F, G and J are closed and merged to `master` (G and J as PR #49,
`ec42217`). Do not revisit or re-verify closed batches. Batch H is on `AD/pbi_resolution`
as uncommitted working-tree changes.

## 1. Session goal

Plan, implement and verify Batch H from review §9: §6.2 (one linter version), §6.3 (LF
policy for module/checksum files) and §6.5 (release tag interpolated into shell). The owner
picked H. Previous session's goal, for context: Batch J (§2.2), closed.

### Previous session goal (Batch J, closed)

Plan, implement, verify and close Batch J from review §9: §2.2, where
`internal/services/zone_content` consumed and built DTOs under an accepted allow-list
exception the owner had reopened. The owner picked J over H in this session. This closing
turn marked the finding fixed and handed this file forward.

## 2. Fixes applied

### Batch H, uncommitted, marked FIXED pending the owner's commit

- **§6.3:** [.gitattributes](../.gitattributes) adds `go.mod text eol=lf` and
  `go.sum text eol=lf`. Only the four module/checksum working copies were re-materialized
  (the index was already LF). Windows `go mod tidy -diff` now exits 0 in both modules.
- **§6.2:** [tools/go.mod](../tools/go.mod) pins `golangci-lint/v2 v2.13.1`. MVS/tidy changed
  only indirect requirements: 47 bumps, 2 added, 2 removed. `wire`, `gcov2lcov`, `go 1.27.0`
  and the `tool` block are unchanged. The [PR lint job](../.github/workflows/pr-validation.yml)
  reads the version with a fail-fast `go list -m` step, so `tools/go.mod` is the single
  source.
- **§6.5:** [release.yml](../.github/workflows/release.yml) has a new `validate` job. It reads
  the tag only through `env`, enforces
  `^v[0-9]+\.[0-9]+(\.[0-9]+)?(-[0-9A-Za-z.]+)?$`, and outputs `tag` and `prerelease`.
  - The build injects `main.version` from `env.VERSION`.
  - The release takes `tag_name`, `name` and `prerelease` from the validated outputs.
  - The checkout ref and concurrency group are unchanged, and there is no existence check.

### Batch J, committed in `3c0ad87` and **marked FIXED**. Closed; do not reopen.

- **§2.2, the whole zone-content exception removed.**
  - New [content_rule_model](../internal/models/content_rule_model/) package: `ContentRuleKey`
    (5 consts) and `ContentRuleEditorKind` (3 consts) moved from `internal/dtos`, plus
    `ContentRuleOption{Key, Name}`, flat `ContentRuleComposition` and
    `ContentRuleDescription`.
  - [zoneContentEditorService.go](../internal/services/zone_content/zoneContentEditorService.go):
    `ComposeContentRule(ContentRuleComposition) (ContentRuleRow, bool)`;
    `GetDefaultContentRules([]ContentRuleOption)`; markers and display name take
    `[]ContentRuleDescription`. No `dtos` import; `validRule` deleted.
  - [zoneContentHandler.go](../internal/handlers/zoneContentHandler.go) does every DTO ⇄
    Model conversion and builds `ContentRuleCompositionResultDto{Rule, Valid}`.
  - `internal/services/zone_content` removed from `dtoNamerAllowList` in
    [layering_test.go](../test/unit/architecture/dependency/layering_test.go). Bonuses stays.

### Batch G, committed in `35e0fab` and **marked FIXED** (prior session). Closed.

All are committed in `35e0fab` and **marked FIXED**. They are closed; do not reopen them.

- **§3.1, the graph summary cache.**
  [zoneEditorGraphState.go](../app/gui/dialogs/zoneEditorGraphState.go) caches the
  `ZoneEditorGraphDto` behind its own `graphDirty` flag, which is never shared with
  `geometryDirty`.
  - All seven structural mutators call `markGraphDirty`: `setEditingSet`,
    `addConnection`, `deleteConnection`, `undoSessionEdits`, `addZoneAt`, `deleteZone`
    and `applyQualityMutation`.
  - `layoutStatus` in [zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go)
    rebuilds only when the flag is set, still through a fresh `derefConnections` copy.
  - Idle frames make **0** handler calls, and each structural edit makes exactly **1**.
- **The one-frame-late status line** (a pre-existing issue the owner chose to fix). The
  status line is drawn before canvas input and side-panel writeback. The key in
  [zoneEditorStatusKey.go](../app/gui/dialogs/zoneEditorStatusKey.go) captures the hint,
  both add modes, the zone and connection counts, and the dirty flag.
  `requestLateStatusRedraw` compares it at the end of `Body` and issues
  `op.InvalidateCmd` when it changed. The key is stored unconditionally after any rebuild,
  so redraws always settle, even in the error branch.

## 3. Features added / changed

### Delivered in Batch H, uncommitted

| ID | Decision |
| --- | --- |
| H1 | Single linter version v2.13.1. |
| H2 | `tools/go.mod` is the single source; CI reads it. |
| H3 | `go get` + `tidy` in `tools/` only; accept the indirect bumps MVS forces. |
| H4 | `eol=lf` for `go.mod`/`go.sum`; the agent refreshes only those four working copies. |
| H5 | Tag format `v` + 2–3 numeric parts + optional `-[0-9A-Za-z.]+` prerelease. |
| H6 | Bash validation step in a `validate` job; tag only via `env`. |
| H7 | A `-` suffix publishes as a GitHub prerelease. |
| H8 | Checkout ref, concurrency group unchanged; no tag-existence check. |
| H9 | Local tag cases are a throwaway one-off run (Git for Windows bash); nothing committed for it; no WSL. |
| H10 | No build/unit/coverage/testlayoutcheck baseline or rerun: no Go code changes (owner override). |
| H11 | Windows-only verification; Linux via the owner's PR CI. |
| H12 | Reviews by GPT-6.1 Sol. |

### Delivered in Batch J, committed and settled

The owner decisions, restated here in case the plan is deleted:

| ID | Decision |
| --- | --- |
| J1 | Remove the **whole** zone-content exception (all four DTO-bearing service methods); keep the bonuses exception. |
| J2 | The composition unit is `ContentRuleRow`, not `ZoneContentRow`; the service returns `(ContentRuleRow, bool)`. |
| J3 | `ContentRuleKey` and `ContentRuleEditorKind` live in the model package, not `internal/common`, because models may import only entities, helpers and registry. No `dtos` alias remains. |
| J4 | New `internal/models/content_rule_model/`, one type per file. |
| J5 | `ContentRuleComposition` is flat (`Key, Name, DistanceNames, DistanceIndex, IsGuarded, IsSoloEncounter, VariantIDs, VariantIndex`); the request DTO embeds it. |
| J6 | `ContentRuleDescription{Key, DisplayText, Marker, VariantLabel, Valid, SavedRule}`; the description DTO embeds it. |
| J7 | `ContentRuleOption{Key, Name}`; the option DTO embeds it and keeps `Description, Marker, EditorKind, EditorLabel`. The "default only when Guarded is offered" check stays. |
| J8 | `ContentRuleCompositionResultDto{Rule, Valid}` stays; handler interfaces, `ContentRuleEditorOptionsDto` and `ContentRuleVariantOptionDto` are unchanged. |
| J9 | The catalogue and describe logic stay in `contentRuleHandler`; business logic in a handler is recorded, not assigned. |

- **Go 1.27 promoted-field literals.** The plan assumed nested literals
  (`ContentRuleOption: content_rule_model.ContentRuleOption{...}`); Go 1.27 accepts promoted
  fields directly and lint's `modernize/embedlit` rejects the nested form. Literals are flat;
  only literals wrapping an existing model value keep the embedded key.
- **No behaviour change.** No GUI golden moved; Wire, handler interfaces and constructors are
  unchanged.

### Delivered in Batch G, committed and settled

The owner decisions, restated here because the plan has been deleted:

| ID | Decision |
| --- | --- |
| G1 | Measure first, then fix. |
| G2 | Acceptance is deterministic: idle frames make 0 `DescribeZoneEditorGraph` calls, and each structural edit makes exactly 1. Benchmarks are recorded with no percentage gate and no global allocation threshold. |
| G3 | The O(z+c) isolation rewrite is adopted only if it is faster at 24 and 40 zones and within ±5% at 4 and 12. |
| G4 | The cache lives in the dialog with its own dirty flag. The handler stays stateless. |
| G5 | Handler calls are counted by a test-helper wrapper embedding `IGuiHandler`, with no new test exports. |
| G6 | There is an untagged, GPU-free idle-frame benchmark. |
| G7/G8 | Late status changes (hint, both add modes, counts, diagnostics) trigger an immediate `op.InvalidateCmd`. The layout order is unchanged. |
| G9/G9a | Benchmarks run at 4/12/24/40 zones on **deterministic** fixtures, because generation uses the unseeded global `rand`. |

- **The O(zones + connections) `FindIsolatedZones` rewrite was measured and rejected**
  under G3. It was faster only at 24 zones (−19%) and 40 zones (−33%). It was slower at
  4 zones (+292%) and 12 zones (+67%), and added 3 allocations per call. The nested loop
  is unchanged. Do not re-propose it without new measurements.
- **New untagged benchmarks** in
  [zoneEditorGraph_test.go](../test/performance/zoneEditorGraph_test.go) run on fixed
  4/12/24/40-zone graphs (G9a, deterministic because generation is unseeded): the
  handler, the service, and a GPU-free dialog idle frame.
- **Idle-frame result:** B/op fell from 4457–18478 to 3049–3073, now flat across graph
  sizes, and allocations fell by 1–4 per frame. ns/op at 40 zones fell from 215030 to
  190195. The other ns/op changes were within run-to-run noise: unchanged code varied up
  to 9%.

### Delivered in Batch F, committed and settled

The shared curve builder (`preview.ConnectionCurveLayout.Build`, 21px spacing, obstacle
bending everywhere), effective-portal classification, the type dropdown and live-list hit
tests are committed, merged, and must not be reopened. The Batch F plan file is no longer
in the tree; the merged code is the reference.

### Settled behavior that must stay preserved

- The road tri-state is unchanged: explicit `false` is roadless, explicit `true` is roaded,
  and `nil` is roaded only for explicit Portals (`IsExplicitPortal`). Effective
  classification never feeds `HasRoad` or road policy.
- Explicit Portal flags are preserved exactly, `false` and `nil` included. With settings
  present, road policy overwrites `Road` on every non-explicit-Portal connection,
  including custom, imported, Default, empty, arena and proximity edges.
- Internal castle/object/foothold roads are independent of the between-zone road checkbox.
  Roads-off never disables them and never removes valid Portal approaches; roads-on
  restores eligible target sets rather than connector-record equality.
- A `nil` `EditorState` preserves existing roads and content: policy runs only when the
  state is present.
- Final cleanup is zone-scoped and treats nil/empty mandatory content as authoritative,
  removing only confirmed invalid `MainObject`, incident `Connection` or named
  `MandatoryContent` references.
- The arena marker is never an anchor; spawn anchoring and rebasing of shifted imported
  non-arena anchors stay intact. Sources are cloned before mutation - no shared backing
  arrays are written through.
- The editor legend stays removed by owner decision. Preview keeps its four-entry key
  (`Road`, `No road`, `Portal`, `Portal without road`).
- PNG roadless strokes apply 50% once per edge through one reusable local mask; the opaque
  raster path, dashes and clipping are unchanged.
- The output directory stays machine-detected, with an explicit session-only picker escape
  hatch. It is never persisted, and no fallback authorizes an unrelated directory.
- All Batch E behavior (effective-mode invalidation, the two-player tournament lock, guard
  propagation and Custom) is committed and settled. Do not reopen it.
- All Batch F behavior above is committed and settled. The shared curve builder stays the
  single source of curves for the editor, the Preview and the PNG; do not reintroduce
  per-renderer curve math.
- All Batch G behavior is committed and settled. The graph cache keeps its own dirty flag:
  never fold it into `geometryDirty`, which hit tests clear mid-input. Any new edit that
  replaces the zone or connection lists must call `markGraphDirty`. The status key must
  stay built from the full state, and never from only the fields the branch drew,
  otherwise the redraw never settles.
- All Batch J behavior is committed and settled. `internal/services/zone_content` must not
  name a DTO again, and the DTO allow-list only ever shrinks. New content-rule shapes go in
  `content_rule_model`; conversion stays in the handlers.

## 4. File modifications

**Batch H, uncommitted** (plan committed by the owner as `25b4c98`):
- [.gitattributes](../.gitattributes): `go.mod`/`go.sum` `text eol=lf` and a reworded comment.
- [tools/go.mod](../tools/go.mod), [tools/go.sum](../tools/go.sum): golangci-lint v2.13.1
  plus its forced indirect changes.
- [pr-validation.yml](../.github/workflows/pr-validation.yml): the `lint-version` step and
  `version: ${{ steps.lint-version.outputs.version }}`.
- [release.yml](../.github/workflows/release.yml): the `validate` job, `needs`, `env.VERSION`,
  and validated release inputs plus `prerelease`.
- [Batch H plan](plans/batch-h-ci-tooling-hardening.md): all phases, recap and deployment.
- [Surviving review](backlog/review-gpt-6-astra-09-07.md): §6.2, §6.3, §6.5 marked FIXED
  (pending commit), row H, progress line. §8 untouched, nothing renumbered.
- This handoff.

`go.mod` and `go.sum` at the root were re-materialized to LF in the working tree, with no
content or Git change. **Not touched:** any Go code, tests, `data/`, the schema, the
registry, Wire, the output path, `.golangci.yml`, README/QUICKSTART/AGENTS.md, and `wire`
and `gcov2lcov` versions.

**Previous closing turn (Batch J, merged):** documentation only, in three files:
- [Surviving review](backlog/review-gpt-6-astra-09-07.md): §2.2 marked `✅ FIXED` in place
  with a Progress paragraph (J9 recorded, not assigned), the §2 architecture-inventory
  bullet updated, the §9 row J completed, and the progress line refreshed to 18 fixed /
  14 remaining. Nothing was renumbered, and the review's §8 hash is unchanged.
- The Batch J plan (since deleted by the owner): Phase 4, Final Recap and
  Deployment Plan marked complete/CLOSED.
- This handoff.

**Batch J inventory, committed in `3c0ad87`** (plan in `6649def`):

**Production (new):** [content_rule_model/](../internal/models/content_rule_model/):
`contentRuleKey.go`, `contentRuleEditorKind.go`, `contentRuleOption.go`,
`contentRuleComposition.go`, `contentRuleDescription.go`.

**Production (deleted):** `internal/dtos/contentRuleKey.go`,
`internal/dtos/contentRuleEditorKind.go`.

**Production (edited):**
- `internal/dtos/`: `contentRuleOptionDto.go`, `contentRuleCompositionRequestDto.go`,
  `contentRuleDescriptionDto.go` embed the models.
- [zoneContentEditorService.go](../internal/services/zone_content/zoneContentEditorService.go)
  and its interface: model signatures.
- [zoneContentHandler.go](../internal/handlers/zoneContentHandler.go): conversions.
- [contentRuleHandler.go](../internal/handlers/contentRuleHandler.go) and
  [ruleDialog.go](../app/gui/dialogs/ruleDialog.go): model constants; the dialog fills
  `Key`/`Name` instead of `Option`.

**Test helpers:** `test/test_helpers/zoneContentEditorServiceMock.go` (new signatures) and
`templateHandlerMock.go` (`ComposeContentRule` wraps the comma-ok pair).

**Tests:** listed in §5, plus the allow-list edit in `layering_test.go`. No goldens changed.

**Not touched in Batch J:** `data/`, `internal/entities/template_entity/`,
`internal/registry/`, the output path, generated Wire, handler interfaces, `guiHandler.go`,
the bonuses service, dependencies, `.agent/backlog/test_observations.md` and
`.agent/backlog/owner_findings.md` (its O15 entry stays, as fixed entries have in earlier
batches; the review's §0 is the disposition of record).

**Batch G inventory, committed in `35e0fab`:**

**Production (new):**
- [zoneEditorGraphState.go](../app/gui/dialogs/zoneEditorGraphState.go): the cached graph,
  `graphDirty`, `shownStatus` and `markGraphDirty`.
- [zoneEditorStatusKey.go](../app/gui/dialogs/zoneEditorStatusKey.go): the comparable
  status key.

**Production (edited):**
- [zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go): embeds the state, adds
  the cached `layoutStatus`, `statusKey`, `requestLateStatusRedraw`, and 6 mutator marks.
- [zoneEditorZoneProps.go](../app/gui/dialogs/zoneEditorZoneProps.go):
  `applyQualityMutation` marks the graph (gofmt realigned one trailing comment).

**Test helpers:**
- [appRunner.go](../test/test_helpers/integration_common/appRunner.go): a private
  `newAppRunner` shared by `NewAppRunnerWithFileSystem` and the new
  `NewAppRunnerWithGuiHandler`, plus `ImmediateRedrawRequested`.
- [graphDescriptionCounter.go](../test/test_helpers/integration_common/graphDescriptionCounter.go)
  (new, `integration_test`).

**Tests:** listed in §5. No goldens changed.

**Agent docs:**
- the Batch G plan (first version `a587617`, completed in `35e0fab`; since deleted by the
  owner).
- [test_observations.md](backlog/test_observations.md): records the GUI-only statements of
  the graph cache.
- `.agent/memories/gui-and-tests.md` (gitignored): redraw assertions, handler wrappers,
  and dialog-direct limits.
- This handoff.

**Not touched in Batch G:** `data/`, `internal/entities/template_entity/`,
`internal/registry/`, the output path, generated Wire, `connectionEditorService.go`
(reverted byte-identical), dependencies and `.agent/backlog/owner_findings.md`.

## 5. Tests added or updated

### Batch H (uncommitted)

No Go tests were added or run: no Go code changed, and the owner waived the build, unit,
coverage and testlayoutcheck runs (H10). The coverage baseline stays **74.5% (6915 / 9263)**,
from Batch J. Verification of record (Windows/amd64, Go 1.27.0):
- `git ls-files --eol` on the four module/checksum files: `i/lf w/lf attr/text eol=lf`.
- `go mod tidy -diff`: root and `tools/` both exit 0, before and after the bump.
- `go list -m` in `tools/` returns `v2.13.1`. A temporary build's `go version -m` matches the
  installed binary's module sum `h1:RuM4Ocl…`.
- Report-only lint: **0 issues** (only the three pre-existing unused-exclusion warnings).
- The lint-version step was run once in Git bash. Success prints `version=v2.13.1`; an
  unknown module and an empty value both exit 1 under `bash --noprofile --norc -eo pipefail`.
- The release `validate` step body was extracted from the YAML and run against 25 crafted
  cases (4 stable, 3 prerelease, 18 rejects including quotes, `;`, `$()`, a backtick,
  newlines, a tab and whitespace) and all 26 existing tags. All pass. Temp files deleted.
- actionlint v1.7.12: both workflows clean. A negative control with a misspelled output
  property is caught. shellcheck is not installed.
- Reviews (GPT-6.1 Sol): plan APPROVE WITH CHANGES (4 findings, all applied);
  implementation **APPROVE**, no findings.
- Not run: Linux tidy, the PR lint job and `release.yml` at runtime (PR CI and the next real
  release do that).

### Batch J, previous closing turn

**That closing turn added no tests and reran nothing** beyond confirming that the committed
`3c0ad87` builds and passes `go test ./test/unit/...`. The Batch J verification below is the
evidence of record.

### Batch J (`3c0ad87`)

**Unit tests** (mirrored layout, triple-A, `t.Parallel()`, one assertion each):
- `zoneContentEditorService/composeContentRule_test.go`, rewritten on the model: unknown key
  rejected and empty row; distance `-1` / past-the-end rejected for both distance kinds;
  out-of-range indices return an empty row; exact rows for road, town, variant, and guarded
  and solo with `true` and `false` (stored as a non-nil `false`, not nil); all five valid
  kinds accepted.
- The other three service test files migrated to model inputs, same scenarios.
- `zoneContentHandler/`: the embedded composition reaches the service; accept →
  `{rule, true}`; reject → `{rule, false}` passed through; only Key/Name reach
  `GetDefaultContentRules`; descriptions reach the service unwrapped and in order.
- `contentRuleHandler/describeContentRule_test.go`: a table over all five rule names for
  `contentRuleKeyFromName` (was 71.4%, now 100%).
- `contentRuleHandler` and `guiHandler` tests: constant renames only.

**Architecture gate:** a temporary blank `dtos` import in the service made
`TestWhenDtoConsumersAreScanned_OnlyTheApiBoundaryAndAppNameADto` fail on that file; the
file hash matched after removal and the test passed again.

**Baseline** (HEAD `33bff1f`): build PASS; coverage **74.5%** (6909 / 9261); `./test/...`
PASS; tagged run PASS (root 3.429s, GUI 26.150s); testlayoutcheck PASS; gofmt clean; lint 0.

**Final** (Windows/amd64, Go 1.27.0, empty `GOFLAGS`): build and `go vet` under no tag,
`integration_test` and `integration_test,gui` PASS; coverage **74.5%** (6915 / 9263: +2
handler statements, +6 covered; every function in the three touched files at 100%);
`./test/...` PASS; tagged run PASS (root 4.103s, GUI 33.104s), no `.failure` files;
testlayoutcheck PASS; gofmt clean; lint **0 issues** after `--fix` (34 findings, all in this
batch's files: gofmt/gci on the new files, golines, 18 `modernize/embedlit`).

**Reviews** (Claude Opus 5.5):
- Plan: APPROVE WITH CHANGES; all 10 findings applied.
- Implementation: APPROVE; 3 nits, all applied.

### Batch G (`35e0fab`, prior session)

**Unit tests** (mirrored layout, triple-A, `t.Parallel()`):
- `connectionEditorService/findIsolatedZones_test.go`: 6 new tests. They lock the `nil`
  result, input order, `To`-only references, unknown endpoints, duplicate names and
  self-loops.
- `zoneEditorHandler/describeZoneEditorGraph_test.go`: 1 new test for a `nil` isolated
  result.

**GUI integration** (`integration_test && gui`):
[zoneEditorGraphCache_integration_test.go](../test/integration/gui/zoneEditorGraphCache_integration_test.go),
with 20 tests:
- Call counts:
  - opening the editor makes 1 call, and idle frames make 0;
  - 8 structural edits make 1 each;
  - 4 non-structural edits make 0 each.
- 2 isolation correctness tests.
- 3 window-harness redraw tests: the positive case, an empty-canvas negative control, and
  convergence.
- 1 dialog-direct error-state convergence test.

Each group was shown to fail under a temporary production mutation, reverted afterwards.

**Baseline** (unchanged code, HEAD `a587617`):
- build: PASS;
- coverage: **74.6%**;
- `go test ./test/...`: PASS;
- tagged run: PASS, root 6.453s, GUI 47.423s (machine load);
- `testlayoutcheck`: PASS; `gofmt`: clean;
- lint: 0 issues.

**Final** (Windows/amd64, Go 1.27.0, empty `GOFLAGS`):
- build: PASS;
- `-p=2` coverage: PASS, **74.5%**. Covered statements are unchanged at 6911, and the
  denominator grew 9246 → 9261 with 15 GUI-only dialog statements. **Accepted by the
  owner with the close-out. 74.5% is the new baseline for the next batch.**
- `go test ./test/...`: PASS;
- tagged run: PASS, root 3.322s, GUI 32.380s;
- the new GUI tests at `-count=5`: PASS;
- `testlayoutcheck`: PASS; `gofmt -l`: clean;
- lint: **0 issues**, after fixing 4 new ones: `funcorder` ×2, `gochecknoglobals` and
  `funlen`.

**Reviews** (Claude Opus 5.5):
- Plan: APPROVE WITH CHANGES. All 12 findings were applied, and G9a was owner-approved.
- Implementation: APPROVE WITH CHANGES. The 1 minor finding and 3 of the 4 nits were
  applied. Nit 5 was declined with a reason, and one out-of-scope observation was
  recorded (see §7).

**Not measured:** native Linux, Steam Deck, the race detector (cgo is off locally), and
in-game behaviour.

## 6. Git status snapshot

The branch is **`AD/pbi_resolution`**, HEAD **`25b4c98` ("H Plan")**, on top of `9d302a5`
("init", AGENTS.md) and `master` `ec42217` (Batches G and J, PR #49). It has no upstream.

`git status --short` shows only Batch H's unstaged edits:

```text
 M .agent/backlog/review-gpt-6-astra-09-07.md
 M .agent/plans/batch-h-ci-tooling-hardening.md
 M .agent/session-carry-forward.md
 M .gitattributes
 M .github/workflows/pr-validation.yml
 M .github/workflows/release.yml
 M tools/go.mod
 M tools/go.sum
```

`.agent/memories/` and `tmp/` are gitignored.

**The assistant performed no staging, unstaging, commit, push, stash, branch switch or
worktree change.** The only Git command that wrote files was the H4-approved
`git checkout -- <path>` restore of the four unmodified module/checksum files.

## 7. Rejections / things the user declined

- **Batch H.**
  - Declined alternatives:
    - downgrading CI to v2.12.2, or moving to the latest linter;
    - hard-coding the linter version in both places;
    - strict `vX.Y.Z` (it rejects `v0.1`) and a loose charset;
    - a Go helper for tag validation;
    - checking out `refs/tags/<tag>` explicitly;
    - an explicit tag-existence check;
    - grouping concurrency by tag;
    - any WSL use.
  - The owner waived the build, unit, coverage and testlayoutcheck baseline and reruns,
    because no Go code changed.
  - The owner does not want local test scripts for the workflow in the repository; the
    validation lives only in the pipeline.
  - Recorded, not assigned: other static `${{ matrix.* }}`/`${{ env.* }}` interpolations in
    `run:` blocks are not user-controlled and were left alone.
- **Batch J.**
  - Rejected designs: keeping `ContentRuleKey` in `internal/common` (breaks models'
    import rule), a `dtos` type alias, nesting the option inside `ContentRuleComposition`,
    `*ContentRuleRow` or a result-struct model instead of comma-ok, deleting the result
    DTO, narrowing the descriptions to plain strings, and dropping the Guarded check in
    `GetDefaultContentRules`.
  - Moving the `contentRuleHandler` catalogue/describe logic into `content_rules` was
    declined for this batch (J9): recorded, not assigned.
  - The plan's nested-literal instruction was superseded by Go 1.27 promoted-field literals
    (see §3). An interim `funlen` helper extraction in `contentRuleHandler.go` was reverted.
- **Batch G.**
  - The O(z+c) isolation rewrite failed G3 and was reverted, with the numbers recorded.
  - Dialog-direct invalidate tests were planned but proved infeasible: the canvas offset
    depends on the toolbar height. The window harness was used instead, with a settle loop
    and a negative control.
  - Implementation review nit 5 (the open-editor test's Act only reads the counter) was
    declined, because opening the editor is the action.
  - Recorded, not assigned: the toolbar's "Delete selected" enabled state is still one
    frame late after a canvas selection, as at HEAD, because `hasSelection` is not in the
    status key. It is out of §3.1 scope.

- **Batch F, settled by the owner.**
  - Editor-only obstacle bending, and dropping obstacle bending entirely, were both
    rejected. Bending now applies everywhere.
  - The editor's 18px spacing was rejected in favour of 21px.
  - Always bending a single curve to the positive side was rejected in favour of the shorter
    detour, with positive only on a tie.
  - A display-only type change was rejected: leaving Portal clears the rules.
  - A case-sensitive preview shape check was rejected.
  - Keeping the current editor portal width was rejected: portals are drawn thinner.
  - Listing every registry type, or keeping two items with a silent Direct fallback, were
    both rejected in favour of Direct, Portal, plus the unlisted stored type.
- **Declined review items.**
  - Plan review (GPT-6 Sol): the "hidden pair members shift slots" finding was declined.
    Every member of a pair shares both endpoint names, so a pair is either fully visible or
    fully skipped.
  - Implementation review nit 3 (curves built twice per editor rebuild) was noted without
    change.
- **Do not reuse a blanket `-update` snapshot run.** It rewrote about 280 unrelated goldens,
  which were restored. Accept only the `.failure` files a plain run produces.
- **Still in force from earlier batches:**
  - Do not restore the editor legend or its deleted test.
  - Do not reopen Batch C or Batch E decisions: arena recomputation, remembered player
    counts, persisted guard preset identity, the Plastic/Bronze tables, `GeneratorConfig`
    player rejection.
  - Keep the Batch E safeguards.
  - Topology retirement (Batch K) stays out of scope.
  - No opportunistic DTO, schema, package, allocation or output-path work.
  - Never claim unobserved engine outcomes.
- **Recorded, not assigned:**
  - `contentRuleHandler` holds the content-rule option catalogue and describe logic
    (business logic in a handler; J9).
  - A service comment still names the legacy rebuild entry point.
  - A zero-edge tournament count remark.
  - The defensive selected-pointer rebinding branch.
  - Double curve construction per editor rebuild.

## 8. Open questions and confirmed later scope

Batch C policy/scope approval remains settled. **Phase 3 is complete with no blockers.**
Owner Phase 3 authorization, 2026-09-10: "Changes reviewed, you can proceed".
Phase 4 remains unstarted and belongs to a new session. Do not repeat settled questions.

Non-blocking review observations, outside closeout scope:

- `ConnectionLegendRow` remains an export used only by `LegendRows` after editor legend
  removal. No behavior issue; optional simplification is not required for Phase 4.
- Pre-existing editor dropdown normalization: its non-`WasUpdated` writeback can turn
  types not offered by the Direct/Portal dropdown (for example Proximity) into Direct
  on selection. This predates Phase 3; road policy still stamps correctly. Investigate
  separately if requested, do not fix opportunistically in PNG work. See
  [writebackProps](../app/gui/dialogs/zoneEditorConnectionProps.go#L105-L119).

The owner will validate true/false/nil engine behavior after changes. Existing engine
evidence is inconclusive: schema is optional bool; shipped examples use both values
for direct and portal connections; examples do not prove defaults.

Retained owner decisions for later work:

- §2.3/O08: remove Ring, Hub, Chain, Shared Web and corresponding tournament builders; reject retired saved IDs without modifying the current document; surviving tournament fallback cases use balanced generation. Preserve Geometric Hub and shared hub-zone concepts.
- §2.4/O11: investigate live state and persistence first; compare exact-base reconstruction against regenerated untouched zones plus deltas. Cover identity, randomness, upgrades, defaults, derived fields, migration and measured size/performance. Owner approves feasibility before representation changes.
- §2.5/O12: structured entries following BonusEntry; bans carry SID, overrides SID/GuardValue, preserve semantics and Variant=-1. Legacy/UI text parsing stays at boundaries.
- §2.6/O13: section current groups except TemplateIdentity, MapSettings and SchemaOptions remain flat. Coordinate migration with §2.5 and approved §2.4 outcome; retain supported legacy loading. Section keys need planning.
- §2.7/O14: all six zone-content accessors on drivers.State, replacing panel callbacks with zone/tier selection; retain validation, dirty tracking and snapshot isolation.
- §2.8/O16: rename services/zones to zone_services and nested interfaces to zone_service_interfaces, updating imports/tests/links/generated wiring without behavior change.
- §2.9/O18: General/Layout/Bonuses named subpackages with private cohesive sections and public aliases/constructor compatibility; Preview unchanged.
- §2.10/O19: public GeneratorConfig lookup preserving Highest→Hub, unknown→nil and read-only borrowed-row semantics.
- §4.1/O20: permitted production Vec2 audit, equivalent methods, justified tested additions, numerical and hot-path behavior preserved.
- §2.2: zone-content DTO removal reopened; bonuses exception remains accepted. Replacement row/result API and full-exception versus single-seam scope still need approval.
- O09 presets remain deferred and uninvestigated.

Other pending decisions: arena/manual invalidation and effective-mode aliases (§1.5), custom guards/preset identity (§1.11), previous non-tournament player count (§1.12), GUI/PNG curve agreement (§2.1), tooling version/EOL/release tags (§6.2/§6.3/§6.5). Review §10 retains in-game bonuses/bans, hero-hire and historical preview validation. Previous tidy differences were Windows checksum EOL-only, not dependency drift.

## 9. Next recommended actions

**Batch H is implemented and awaits the owner's review and commit.** If the owner asks
for changes, apply them, re-run only the affected Batch H checks, and keep the review's
Progress paragraphs truthful. Once committed, replace "uncommitted, awaiting owner commit"
with the commit hash in §6.2, §6.3, §6.5 and row H. After merge, the PR's "PR Tests" and
"Tools Module" runs are the Linux evidence, and the next legitimate release exercises
`release.yml`.

After that, the next unit of work is the **owner's choice** among the remaining batches of
the surviving review's §9 table (11 items left):
- **I: docs,** §7.1 and §7.2. This is next in table order.
- **K–P:** topology retirement, the compact-state investigation, the persistence format,
  panel/state organization, naming and lookup, and the Vec2 audit. The binding scope for
  each is in §8's retained decisions.

Routing for the next session, in order:
1. Read [AGENTS.md](../AGENTS.md), this handoff, then the chosen items and their §9 row.
2. Inspect the current code yourself, because review line numbers predate Batches F–J and H.
3. Ask the owner every open decision in the items, summarize the scope, and get approval.
4. Write a new durable plan under `.agent/plans/`, get an independent review (GPT-6.1 Sol)
   and explicit plan approval, and capture a fresh baseline (coverage 74.5%, 6915 / 9263)
   before the first Go edit.

**Deployment.** Nothing is deployed and nothing is authorized to be. The owner alone
stages, commits, merges and releases. Batch H's deployment steps are in its plan. No
schema migration, Wire regeneration or output-directory change is required or pending.
Native Linux and Steam Deck execution remain unmeasured.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first, then [this handoff](session-carry-forward.md). It
> is self-contained.
>
> **Batch H is IMPLEMENTED, uncommitted, awaiting the owner's review.** It is on
> `AD/pbi_resolution` (HEAD `25b4c98`, the owner's "H Plan" commit; `master` `ec42217`
> already contains Batches G and J). It pins golangci-lint v2.13.1 in `tools/go.mod` as
> CI's single source, adds LF attributes for `go.mod`/`go.sum`, and validates the release
> tag through `env` in a new `validate` job. Review §6.2, §6.3 and §6.5 are marked
> **FIXED** ahead of the commit at the owner's request (21 fixed, 11 remaining). The plan
> [batch-h-ci-tooling-hardening.md](plans/batch-h-ci-tooling-hardening.md) records H1–H12,
> the verification and the deployment steps. §3 of this handoff restates H1–H12, J1–J9 and
> G1–G9a. Do not reimplement or re-review Batch H unless the owner asks for changes, and
> never revisit Batches A–G or J.
>
> **State:** the working tree holds Batch H's eight unstaged edits (see §6). The assistant
> performed no Git mutation beyond the approved restore of four unmodified module files;
> preserve the owner's state exactly.
>
> **Next:** apply owner feedback on Batch H if any. Then the owner picks the next batch
> from review §9 (I, docs, is next in order). Follow the gates in §9: read, inspect, ask,
> summarize, plan, independent review, plan approval, fresh baseline.
>
> **Verification on record:** Batch H is Windows-only. Both tidy dry-runs pass, the lint
> pin is confirmed, lint shows 0 issues, every tag case passes and actionlint is clean.
> Batch J's code baseline (Windows/amd64, Go 1.27.0, empty `GOFLAGS`):
> - build and tagged `go vet`: PASS;
> - coverage: PASS, 74.5% (6915 / 9263);
> - `go test ./test/...`: PASS;
> - tagged `integration_test,gui` run: PASS, root 4.103s, GUI 33.104s;
> - `testlayoutcheck`: PASS; `gofmt`: clean; lint: 0 issues.
>
> **Native Linux, Steam Deck and the race detector are UNAVAILABLE.** No in-game result is
> claimed.
>
> **Go 1.27 literals:** promoted fields of embedded structs are set directly in composite
> literals; `modernize/embedlit` rejects the nested form.
>
> **GUI snapshots:** never keep a blanket `-update`, because it rewrites about 280 goldens.
> Accept only the `.failure` files a plain run produces.
>
> **Hard rules:**
> - Never modify protected data, the template schema or the registry.
> - Keep Windows/Linux compatibility.
> - Never change or persist the machine-detected output directory.
> - Test nontrivial logic, and check coverage before and after.
> - Never stage, unstage, commit, push, stash, switch branches or manipulate worktrees;
>   preserve owner changes.
> - Never bulk-rewrite or hand-edit generated Wire.
> - Never enable global `integration_test`, `gui` or `wireinject` tags, and never
>   introduce fake unit seams.
> - Keep plans durable and resumable.
>
> **Out of scope** unless the owner selects that batch: Batch K topology retirement,
> direct `GeneratorConfig` rejection, DTO cleanup, schema and package work, allocation
> tuning, and every settled alternative in §7. §8 stays verbatim. Its Batch C wording,
> its pending-decision wording for §1.5, §1.11, §1.12, §2.1 and §6.2/§6.3/§6.5, and its
> §2.2 bullet are superseded; its retained later-scope decisions remain binding.
>
> **Preserve:**
> - explicit Portal Road `false`/`nil` and valid approaches;
> - the tri-state road display (`nil` is roaded only for explicit Portals, and effective
>   classification never feeds roads);
> - independent internal roads, nil-state content/road preservation, and source cloning;
> - the single reusable half-opacity PNG edge mask, and the Preview-only legend;
> - the shared curve builder, as the single source of curves;
> - hit tests resolving against the live connection list;
> - the graph cache's own dirty flag (`markGraphDirty` on every edit that replaces the
>   zone or connection lists), and a status key built from the full state;
> - `internal/services/zone_content` naming no DTO, with conversion in the handlers;
> - `tools/go.mod` as the single linter-version source, LF module/checksum files, and the
>   release tag reaching shell only through `env` after validation.
>
> This handoff contains the full continuation context.
