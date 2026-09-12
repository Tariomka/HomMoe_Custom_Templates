# Carry-forward: Batch E closed, Batch F not started

Date: 2026-09-12. Batch E is **closed**. The owner committed the work as **`2062932`
("Phase 4")**, findings §1.5, §1.11 and §1.12 are now marked **FIXED** in the surviving
[review](backlog/review-gpt-6-astra-09-07.md), and its §9 execution-order table records
batch **E as complete**. Nothing in Batch E is re-reviewed, re-verified or reimplemented.

The Batch E plan was **retired by the owner** and its deletion is staged. Do not restore
it, do not read it, do not re-index it or link to it. This handoff is deliberately
self-contained: every fact a fresh session needs is written out here, and no statement
below depends on that retired document.

**Section 8 is preserved verbatim**, contradictions and all. Its Batch C phase/engine
wording and its closing "pending Batch E questions (§1.5/§1.11/§1.12)" wording are
**superseded** - those questions were settled, implemented and are now marked fixed.
Everything else §8 retains as later scope remains **binding**, and it governs scope
decisions past Batch F. Do not edit §8.

Batches A/D, B, C and E are closed. Batch C's engine gate closed on the owner's 2026-09-11
report: "Everything in-game looks good, the road values are correct in engine." Do not
revisit, re-verify or retrieve retired documents from those batches.

## 1. Session goal

Documentation only: mark the Batch E findings fixed and hand this file forward to Batch F.

The owner's authorization was *"Mark this batch items fixed and update carry forward for
the next batch"*, following the earlier acceptance *"Changes reviewed, all seems fine, you
can finish up with Phase 4 and update Session Carry Forward"* and the commit that closed
the batch. No code, test, data or Git change belongs to this turn, and nothing executable
changed, so nothing was rerun.

## 2. Fixes applied

All three Batch E findings have implemented, tested and independently reviewed behavior,
and all three are now **marked FIXED in the surviving review** - the owner's commit
`2062932` satisfied the review's marking protocol. They are closed; do not reopen them.

- **§1.5, effective-mode invalidation.** Effective tournament and effective arena are
  resolved booleans on the domain state, and layout comparison judges those booleans
  rather than the checkbox/selector fields behind them. A transition in either clears the
  entire manual snapshot before anything - a generation attempt included - can observe the
  new state, warns on an actual discard, marks the document unsaved and re-arms Exit.
- **§1.12, tournament player count and loading.** The validator pins tournament states to
  two players in the correct fix order, and the slider itself is locked and moved back
  rather than only its outgoing value, so idle frames stop resubmitting a stale count.
  Loading a tournament file at any other count corrects it, drops the incompatible layout,
  warns, becomes unsaved and re-arms Exit without rewriting the file. The notice survives
  the automatic regeneration that follows it, and a failed generation keeps it for retry.
- **§1.11, guard propagation and explicit Custom.** An actual neutral-quality change
  remaps exact preset guard values on every incident edge, matching the old number in
  order against the old stronger-endpoint table and carrying that named tier - Default
  included - into the new one. Unmatched numbers stay put and display an explicit,
  nonnumeric Custom entry derived from the number, never persisted.

## 3. Features added / changed

### Delivered across Batch E, all owner-approved

- Either effective boolean changing clears the entire manual snapshot and regenerates
  immediately. Warn on an actual discard, mark dirty, re-arm Exit; no confirmation dialog.
  Alias-only representation changes preserve edits. Effective tournament is checkbox OR
  Tournament victory; effective arena is checkbox OR FinalBattle (Guardian Arena in UI).
  Existing selector reset behavior is preserved, and no arena reconciliation pass exists.
- Clearing happens during state update, before generation can fail or consume the
  transition. The validator's count fix does not own snapshot clearing; a domain operation
  returns a discard/correction outcome the driver acts on and shows, and that notice stays
  visible through the automatic regeneration rather than being replaced by a generic one.
- Effective tournament forces and locks two players; turning it off leaves two, with no
  previous-count memory. Invalid tournament loads correct to two, clear manual snapshots,
  warn, become unsaved and re-arm Exit, with no file rewrite until Save. Valid two-player
  loads retain edits and unrelated load policy is unchanged.
- Quality propagation matches the old number in order against the old stronger-endpoint
  table and carries the named tier - Default included - into the new one. GuardZone
  controls placement, not which table wins. Same-quality and castle-only edits recalculate
  nothing. Unmatched numbers stay put and read Custom; exact typed matches are presets,
  even after reload. Plastic/Bronze tables are unchanged, generated Plastic 10,000 staying
  Custom included. Unrelated edges, other fields, source ownership, connection order and
  Apply/Cancel isolation are preserved, with no hidden preset metadata.

### Settled behavior that must stay preserved

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
- Display classification: explicit `false` is roadless, explicit `true` is roaded, `nil`
  is roaded only for explicit Portals.
- The editor legend stays removed by owner decision. Preview keeps its four-entry key
  (`Road`, `No road`, `Portal`, `Portal without road`).
- PNG roadless strokes apply 50% once per edge through one reusable local mask; the opaque
  raster path, geometry, dashes and clipping are unchanged.
- The output directory stays machine-detected, with an explicit session-only picker escape
  hatch. It is never persisted, and no fallback authorizes an unrelated directory.

## 4. File modifications

**This turn is documentation only - exactly two files:**

- [Surviving review](backlog/review-gpt-6-astra-09-07.md): findings §1.5, §1.11 and §1.12
  marked `FIXED` in place following that document's own protocol, and the §9
  execution-order row for batch `E` set to complete. Numbering was not altered.
- [Handoff](session-carry-forward.md): this record, rewritten to the post-commit state,
  with §8 preserved byte-for-byte.

**Owner-owned working-tree state, left exactly as found:** the Batch E plan under
`.agent/plans/` is **staged for deletion** by the owner, who retired it. That deletion was
neither made nor reversed here, the retired file was not read, and nothing links to it.

**Nothing else was touched:** no production code, no tests, no protected data, template
schema or registry, no output path, no opaque raster path, no generated Wire, no topology
code, no dependencies, snapshots or goldens, and never `.agent/backlog/owner_findings.md`.
No new memory, plan or summary files were created, and no state-mutating Git command was
run.

All Batch E production and test code is committed at `2062932` and its predecessors. That
inventory is not repeated here and is not needed again.

## 5. Tests added or updated

**This turn added no tests and reran nothing** - it changed documentation only. The
verification recorded below belongs to the code committed at `2062932` and stands as the
current evidence of record.

Batch E's final test contribution was six combined real-input tests in
[zoneEditorProperties_integration_test.go](../test/integration/gui/zoneEditorProperties_integration_test.go).
Each applies a remapped Silver Medium to Gold Medium incident edge, then switches into
Guardian Arena or Tournament through the General victory selector, across idle frames and
a dialog reopen. They assert positive outcomes, not merely absences: the arena main object
exists on the regenerated Hub graph, the old Neutral-C zone and its remapped edge are
gone, the discard notice survives the generation status that follows, and the tournament
transition replaces both the Hub layout and the manual zone.

Scope was deliberately not duplicated. The existing load-correction, frame, save/reload
and Apply/Cancel persistence suites already cover their own halves and still pass; the
dirty-flag side is unit-covered, so these GUI tests check the visible warning instead.

Earlier phases contributed the unit folders for the effective predicates, the transition
and its outcome, the load override and the new service method, the tagged
`effectiveModes` and `generalPanelTournament` integration files, 46 focused service cases,
21 real-input GUI cases and two persistence regressions. No fake unit seams, no new
`*_testexports.go`, no global tags.

**Recorded verification, Windows/amd64, Go 1.27.0, empty `GOFLAGS`. Historical - measured
against the committed Batch E code, deliberately not rerun this turn:**

- `go build ./...` — PASS.
- `go test -p=2 -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...`
  — PASS. `-p=2` is retained because unbounded Windows linkers were interrupted
  (`0xc000013a`) in an earlier run; that incomplete profile was discarded and never
  reused. This is the baseline over the reviewed production code, unchanged since.
- `go tool cover '-func=coverage.txt'` — **74.9%** total statements against the recorded
  75.1% baseline. This is the owner's accepted exception, not a restored baseline. All
  three new service functions and both quality handler/facade methods are 100.0%.
- `go test ./test/...` — PASS.
- `go test -tags='integration_test,gui' ./test/integration/...` — PASS; root integration
  3.297s, GUI integration 28.558s. The single invocation covers both suites.
- `go run ./cmd/testlayoutcheck .` — PASS. `gofmt -l` over the changed file — clean.
- `golangci-lint-v2 run ./... --issues-exit-code=1` — PASS, zero issues.
- Wire: no constructor or provider changed, so the committed generated output is already
  current and was deliberately not regenerated.
- **Native Linux: UNAVAILABLE.** The probe
  `wsl.exe -d Ubuntu -- sh -lc 'command -v go; command -v pkg-config'` returned empty
  output and exit code 1 when it was run. Nothing was installed; no native Linux, Steam
  Deck or engine result is claimed, then or now.
- **No coverage profile fingerprint is claimed.** The reports were regenerated at
  verification time but their SHA-256 was never measured. The Phase 3 value
  `23B78F2CF15B73902398D8A21A755FD684D3668486E56D8CFA770BE2BE7F9F31` is historical.
- Independent **Claude Opus 5 review of the six new tests: APPROVED**. The plan, Phase 2
  and Phase 3 implementation reviews were approved earlier.

## 6. Git status snapshot

Branch **`AD/modes_and_guard_propagation`**, **one commit ahead of**
`origin/AD/modes_and_guard_propagation`. HEAD is the owner's commit **`2062932`
("Phase 4")**, which supersedes `1b4658c` ("Phase 3 and partial 4") and carries all Batch
E code.

`git status --short` currently reports two entries:

```text
 M .agent/backlog/review-gpt-6-astra-09-07.md
D  .agent/plans/batch-e-effective-modes-and-guards.md
```

The modified review is this turn's documentation edit. The staged deletion is the
**owner's** - they retired the Batch E plan. Preserve it exactly: do not restore the file,
do not unstage the deletion, do not read the retired document.

**The assistant performed no Git mutation of any kind** - no staging, unstaging, commit,
push, stash, branch switch or worktree change - and every owner change is preserved as
found. `.agent/backlog/owner_findings.md` was not read, edited or re-indexed.

## 7. Rejections / things the user declined

- Do not restore the editor legend or its deleted test. Do not reopen settled Batch C
  approval questions, reimplement a completed phase, or re-verify a closed batch.
- Forbidden alternatives, all settled by the owner: arena recomputation or reconciliation
  instead of full snapshot clearing; a remembered previous player count, session or
  persistent; persistent guard preset identity or any new persisted metadata; fixing the
  Plastic/Bronze table discrepancy; and rejecting player counts inside the direct
  `GeneratorConfig` generator, which the first plan review threw out as unapproved. Two
  players are enforced at editor-state validation and application generation only,
  leaving direct generator/provider compatibility and topology algorithms intact.
- Safeguards the revised plan added; do not undo them: visible dirty/warning outcomes
  separate from count validation, invalidation before a generation can consume the
  transition, index-safe selected-pointer replacement, ordered Default matching, and
  `noteStateTransition` returning before merging an empty outcome so a no-op
  `UpdateState` cannot wipe a generation error or a just-saved message on the next idle
  frame. The load→failure→idle→save→idle→retry sequence is pinned by tests.
- Topology retirement/redesign (Batch K, review §2.3) remains explicitly out of scope. No
  opportunistic work either: no DTO removal, no schema changes, no package renames, no
  allocation tuning, no output-path changes. The GUI/PNG geometry-consolidation ban
  recorded here was a **Batch E exclusion**, not a permanent prohibition: review §2.1 is an
  *optional* Batch F item that may be worked only if the owner explicitly approves it into
  scope. It never becomes automatic.
- Do not claim engine defaults or in-game outcomes nobody observed. The only positive
  engine result is the owner's report, and it covers road values only.
- When verifying §8, use explicit UTF-8 for Git stdout
  (`ProcessStartInfo.StandardOutputEncoding`) and for disk reads. An earlier mismatch was
  a PowerShell 5.1 decoding artifact; no §8 content was ever altered. Note also that Git
  stores this file LF-normalized while the worktree is CRLF, so compare disk bytes
  against disk bytes, or LF-normalized text against `git show`.
- Recorded, not assigned: a service comment still names the legacy rebuild entry point
  rather than the current reconciliation finalizer, and a separately observed zero-edge
  tournament count remark is unrelated to Batch E. Neither is a task and neither reopens
  anything.
- The non-nil selected-pointer rebinding branch is defensive and unreachable from current
  real zone-selection input; it is source-reviewed, never claimed as covered. The preset
  tables contain no duplicate values, so first-match precedence is explicit in code but
  not experimentally distinguishable with current data.
- Tooling friction worth knowing, not work to redo: delegated exploration once returned
  silently and its results were re-verified by hand; discovery once used unsupported
  lowercase model names and two guessed paths that do not exist. No files were harmed and
  nothing needs replaying.

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

**Batch E is closed, with no open Batch E questions and no blockers.** The next unit of
work is **Batch F, editor geometry**, taken from the surviving review's §9
execution-order table — *not* from §8, whose retained decisions govern later scope.
Nothing has been planned or implemented for Batch F, and no code may be written before the
gates below are cleared.

Batch F scope, from the review's `F` row:

- **§1.13, deterministic obstacle ties.** The obstacle-bulge selection in the editor
  geometry service takes a winner only on strictly greater magnitude, so two
  equal-magnitude obstacles on opposite sides resolve by map iteration order and identical
  input yields two different curves. Determinism comes first in the batch.
- **§1.14, stale indices from batched pointer events.** The zone-editor canvas drains
  every queued pointer event before rebuilding geometry, while edge lookup indexes the
  live connection slice by a cached index that deletion compacts. Geometry-to-object
  identity must survive input batching. Real input routing only, never unit-only exports.
- **§1.15, effective portal classification mismatch.** The editor recognizes only a
  literal `Portal` connection type, while the preview additionally treats
  `PortalPlacementRules` From/To connections as portals. Classification belongs in one
  model-level place, fed through the existing geometry/handler seam.
- **§2.1, shared curve geometry — OPTIONAL.** Editor and preview build parallel-edge
  curves twice, with different spacing and obstacle deflection in only one of them. This
  is **not automatically in scope**: it enters Batch F only on the owner's explicit
  approval, and the intentional visual differences must be decided before anything moves.

Routing for the next session, in order:

1. Read [AGENTS.md](../AGENTS.md), this handoff, then the surviving review's §1.13, §1.14,
   §1.15, optional §2.1, and its §9 `F` row.
2. **Inspect the controlling code and tests locally before trusting any of it.** The
   review's evidence is historical; line numbers and behavior must be re-confirmed against
   current source. The files are
   [zoneEditorGeometryService.go](../internal/services/connection_editor/zoneEditorGeometryService.go),
   [zoneEditorCanvas.go](../app/gui/dialogs/zoneEditorCanvas.go),
   [zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go),
   [previewLayoutService.go](../internal/services/preview_service/previewLayoutService.go),
   [buildGeometry_test.go](../test/unit/internal/services/connection_editor/zoneEditorGeometryService/buildGeometry_test.go),
   [buildPreviewLayout_test.go](../test/unit/internal/services/preview_service/previewLayoutService/buildPreviewLayout_test.go)
   and [zoneEditorPointer_integration_test.go](../test/integration/gui/zoneEditorPointer_integration_test.go).
3. **Ask the owner** the genuinely undecided product questions: the preferred tie-break
   direction for §1.13; for §1.15, whether changing a connection type should clear the old
   placement rules or merely display the effective type; whether exact editor/PNG curve
   agreement is a requirement; and whether optional §2.1 is in scope at all. Do not guess
   any of these.
4. Summarize the resulting scope back to the owner and obtain approval.
5. Only then write a **new durable Batch F plan** under `.agent/plans/`, have it
   independently reviewed by Claude Opus 5, obtain explicit plan approval, and capture a
   fresh build/unit/coverage/lint baseline before the first edit.

Constraints Batch F inherits and must not quietly break: the settled road tri-state
display contract stays as is — explicit `false` roadless, explicit `true` roaded, `nil`
roaded only for explicit Portals — and that `nil` rule does **not** widen automatically
just because classification becomes "effective". Intentional editor/preview geometry
differences stay until consolidation is separately approved. Coverage must hold at or
above the review's comparable Windows floor with zero lint issues.

**Deployment.** Nothing is deployed and nothing is authorized to be. The owner alone
stages, commits and releases. No schema migration, dependency installation, Wire
regeneration or output-directory change is required or pending. Native Linux and Steam
Deck execution remains unmeasured.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first, then [this handoff](session-carry-forward.md). It
> is self-contained; there is no plan to open.
>
> **Batch E is CLOSED.** All four phases are done, the owner accepted the implementation
> and committed it, and findings §1.5, §1.11 and §1.12 are marked **FIXED** in the
> [surviving review](backlog/review-gpt-6-astra-09-07.md), whose §9 table shows batch `E`
> complete. Do not reimplement, re-verify or re-review any of it, and do not reopen settled
> product questions. The Batch E plan was retired by the owner and its deletion is staged:
> never restore, read or relink it.
>
> **State:** branch `AD/modes_and_guard_propagation`, **one commit ahead of origin**, HEAD
> at the owner's `2062932` ("Phase 4"). The working tree holds the owner's staged deletion
> of the retired plan plus the review edit from the closing turn. The assistant performed
> no Git mutation; preserve owner state exactly.
>
> **Next work is Batch F, editor geometry**, per review §9's `F` row: §1.13 deterministic
> obstacle tie-breaking, §1.14 stale edge indices from batched pointer events, §1.15
> editor-versus-preview portal classification, and **optionally** §2.1 shared curve
> geometry — optional means it enters scope only with the owner's explicit approval, never
> automatically. **No Batch F plan or implementation exists yet, and none may be written
> before the gates below.**
>
> **Do this in order:** read review §1.13, §1.14, §1.15, optional §2.1 and §9; then
> **inspect the current source and tests yourself**, because the review's evidence is
> historical and its line numbers may have moved —
> `internal/services/connection_editor/zoneEditorGeometryService.go`,
> `app/gui/dialogs/zoneEditorCanvas.go`, `app/gui/dialogs/zoneEditorDialog.go`,
> `internal/services/preview_service/previewLayoutService.go`, and the tests
> `test/unit/internal/services/connection_editor/zoneEditorGeometryService/buildGeometry_test.go`,
> `test/unit/internal/services/preview_service/previewLayoutService/buildPreviewLayout_test.go`,
> `test/integration/gui/zoneEditorPointer_integration_test.go`. Then **ask the owner**: the
> preferred tie direction; whether a type change clears old placement rules or only
> displays the effective type; whether exact GUI/PNG curve agreement is required; and
> whether §2.1 is in scope. Summarize the scope, get approval, write a **new durable Batch
> F plan**, get an independent Opus 5 review and explicit plan approval, and capture a fresh
> baseline before the first edit.
>
> **Verification on record, Windows/amd64, Go 1.27.0, empty `GOFLAGS`** (historical, for
> the committed Batch E code — not rerun): `go build ./...`; `go test -p=2 -count=1` unit
> run with `-coverpkg=./internal/...,./app/...`; `go tool cover -func` at **74.9%**, an
> owner-accepted GUI-only decrease against the 75.1% baseline rather than a restored one;
> `go test ./test/...`; `go test -tags='integration_test,gui' ./test/integration/...` with
> root 3.297s and GUI 28.558s; `go run ./cmd/testlayoutcheck .`; clean `gofmt -l` on changed
> files; and zero-issue `golangci-lint-v2 run ./... --issues-exit-code=1`. Generated Wire is
> already current because no constructor changed — do not rerun the generator needlessly.
> **Native Linux/Steam Deck is UNAVAILABLE:** the WSL probe found no `go` and no
> `pkg-config`, nothing was installed, and no native or in-game result may be claimed. No
> coverage profile fingerprint was ever measured.
>
> **Hard rules:** never modify protected data, the template schema or the registry. Keep
> Windows/Linux compatibility. Never change or persist the machine-detected output
> directory. Test nontrivial logic and check before/after coverage. Never stage, unstage,
> commit, push, stash, switch branches or manipulate worktrees; preserve owner changes.
> Never bulk-rewrite or hand-edit generated Wire. Never enable global `integration_test`,
> `gui` or `wireinject` tags, and never introduce fake unit seams. Keep plans durable and
> resumable.
>
> **Out of scope:** Batch K topology retirement, direct `GeneratorConfig` rejection, DTO
> cleanup, schema and package work, allocation tuning, and every settled alternative in §7.
> §8 stays verbatim: its Batch C phase/engine wording and its closing "pending Batch E
> questions" wording are superseded, while its retained later-scope decisions remain binding
> and govern batches after F.
>
> Preserve explicit Portal Road false/nil and valid approaches, the settled tri-state road
> display classification, independent internal roads, nil-state content/road preservation,
> source cloning, one reusable local half-opacity edge mask over unchanged opaque
> rasterization, and the Preview-only legend. This handoff contains the full continuation
> context.
