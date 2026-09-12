# Carry-forward: Batch E Phase 1 complete, Phase 2 next

Date: 2026-09-12.

**Current authority:** the owner-approved [Batch E plan](plans/batch-e-effective-modes-and-guards.md).
The owner committed it in `b9dd616` on `AD/modes_and_guard_propagation`, then asked to
finish Phase 1 and update this handoff. Phase 1 is now complete; Phases 2-4 have not
started. Product decisions and implementation-plan approval are settled. Do not repeat
discovery or approval questions unless a concrete scope change is needed.

**Supersession note.** Batches A/D, B and C remain fully closed: no revisit, separate
re-verification or lookup of retired documents. Batch C's engine gate closed on the
owner's 2026-09-11 report: "Everything in-game looks good, the road values are correct
in engine." Section 8 below is kept **verbatim**, contradictions included. Its Batch C
phase-status and pending-engine wording is superseded; its final paragraph's pending
Batch E questions (§1.5/§1.11/§1.12) are now also settled by this plan. Its other retained
later-scope decisions remain binding. Do not edit §8.

## 1. Session goal

Discover Batch E, settle product questions, write/review the durable plan, then complete
Phase 1 after the owner committed it. Latest instruction: "the plan is commited,
proceed with finishing Phase 1 and update carry forward".

**No production code, tests, snapshots or generated Wire changed.** Phase 1 closeout
ran fresh unit coverage and regenerated ignored reports, then updated the plan and
this handoff. Phase 2 implementation was deliberately not started.

## 2. Fixes applied

No application fixes yet. Review findings §1.5, §1.11 and §1.12 remain unimplemented.

- [Batch E plan](plans/batch-e-effective-modes-and-guards.md): approved product contract,
  implementation seams, tests, independent review and owner authorization recorded;
  Phase 1 completed with fresh per-target coverage.
- [This handoff](session-carry-forward.md): next-session entry point moved to Phase 2;
  historical §8 retained unchanged.

## 3. Features added / changed

### Approved Batch E behavior, not yet implemented

- Either effective tournament or effective arena boolean changing clears the entire
  manual snapshot and regenerates immediately. Warn on actual discard, mark dirty and
  re-arm Exit; no confirmation dialog. Alias-only changes preserve edits. Effective
  tournament is checkbox OR Tournament victory; effective arena is checkbox OR
  FinalBattle (Guardian Arena in UI). Preserve existing selector reset behavior.
- Clear during state update before generation can fail or consume the transition.
  Validator count fixes do not silently own snapshot clearing; a domain operation
  supplies the driver with an actual discard/correction outcome. Keep warnings visible
  through automatic regeneration. No arena reconciliation or retrospective repair.
- Actual neutral-quality changes recalculate every incident edge type. Match its old
  numeric value against the old stronger-endpoint table in order; carry that named
  preset to the new stronger-endpoint table. Default is a named preset too. GuardZone
  does not choose the table. Same-quality/castle-only edits do not recalculate.
- Unmatched numbers stay unchanged and display Custom; exact typed matches count as
  presets, even after reload. Keep existing Plastic/Bronze tables, including the
  generated Plastic 10,000 value remaining Custom. No hidden preset metadata. Preserve
  unrelated edges, other fields, source ownership and pending Apply/Cancel isolation.
- Effective tournament forces/locks two players. Turning it off leaves two, with no
  previous-count memory. Invalid tournament loads (any count other than two) correct
  to two, clear manual snapshots, warn, become unsaved and re-arm Exit. No file rewrite
  until Save. Valid two-player loads retain edits; unrelated load policy is unchanged.

### Settled behavior to preserve

- Explicit Portal flags are preserved exactly, including `false` and `nil`. With
  settings present, road policy overwrites `Road` on every non-explicit-Portal
  connection, including custom, imported, Default, empty, arena and proximity edges.
- Internal castle/object/foothold roads are independent of the between-zone road
  checkbox. Roads-off never disables them and never removes valid Portal approaches;
  roads-on restores eligible target sets rather than connector-record equality.
- A `nil` `EditorState` preserves existing roads and content: policy runs only when the
  state is present.
- Final cleanup is zone-scoped and treats nil/empty mandatory content as authoritative,
  removing only confirmed invalid `MainObject`, incident `Connection` or named
  `MandatoryContent` references.
- The arena marker is never an anchor; spawn anchoring and rebasing of shifted imported
  non-arena anchors stay intact. Sources are cloned before mutation — no shared backing
  arrays are written through.
- Display classification: explicit `false` is roadless, explicit `true` is roaded, and
  `nil` is roaded only for explicit Portals.
- The editor legend stays removed by owner decision. Preview keeps its four-entry key
  (`Road`, `No road`, `Portal`, `Portal without road`).
- PNG roadless strokes apply 50% once per edge through one reusable local mask; the
  opaque raster path, geometry, dashes and clipping are unchanged.
- The output directory stays machine-detected, with an explicit session-only picker
  escape hatch. It is never persisted, and no fallback authorizes an unrelated
  directory.

## 4. File modifications

After the owner's plan commit, this closeout edited only:

- [Batch E plan](plans/batch-e-effective-modes-and-guards.md): Phase 1 complete, branch/
  commit evidence, measured total and target-function baseline, Phase 2 next.
- [Handoff](session-carry-forward.md): this resumable record, with §8 preserved verbatim.
- Ignored [coverage.txt](../coverage.txt), [coverage.html](../coverage.html) and
  [lcov.info](../lcov.info): regenerated by the existing coverage task.

No Go files, tests, snapshots, goldens, assets, dependencies, generated Wire, protected
trees or surviving review were edited. No additional memory/summary files were created.

## 5. Tests added or updated

No tests added or modified. Fresh Batch E pre-change evidence, **Windows/amd64,
Go 1.27.0, 2026-09-12, source at `b9dd616`**, with empty `GOFLAGS`:

- Existing coverage task: `go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...`
  **PASS, 186 unit packages**, followed by HTML/LCOV generation.
- `go tool cover '-func=coverage.txt'`: **PASS, 75.1% total statements**. This is the
  measured Batch E baseline, not a copied historical total. Plan Phase 1 records the
  target-function table and profile SHA-256.
- Target domain/model/validator/handler functions: 100.0% statement coverage. Existing
  driver `handleGenerateTemplate`: 81.0%; load callback and targeted dialog functions:
  0.0% in unit coverage, requiring planned integration tests. General panel has no
  entries in this profile, not a measured 0% or 100%. These numbers do not prove the
  missing Batch E contracts or authorize fake unit seams.
- The full `go test ./test/...` suite was **not rerun** in this closeout. Its last
  recorded result remains PASS on 2026-09-11 from the prior handoff, not current-session
  evidence. Build, tagged integration/GUI, lint, layout checker, Wire generation and
  Linux execution likewise were not run in this Phase 1 closeout.
- Independent Claude Opus 5 **plan** review: revised plan APPROVED on 2026-09-12.
  No Batch E implementation exists to review yet.

Before implementation, check that source/toolchain/coverage inputs still match this
baseline. Refresh only if inputs changed or evidence is unavailable. No separate
closed-batch verification is requested. No assistant in-game or native Linux/Steam
Deck execution is claimed; the historical owner report covers road values only.

## 6. Git status snapshot

Branch **`AD/modes_and_guard_propagation`**, HEAD **`b9dd616` (`Init`)**. Read-only plan
history confirmed that same commit contains the approved plan, with no plan diff at
the start of Phase 1 closeout. The owner created/committed the branch; the agent did
not switch branches or commit anything.

Initial worktree: clean. Initial staged paths: none. Final status has exactly two
unstaged modified files: the plan and handoff listed in §4; coverage reports are ignored.
No source changes are inherited. The complete `git ls-files --stage` snapshot was
compared before/after, alongside the original UTF-8 §8 text.

No staging, unstaging, commit, push, stash, branch or worktree mutation was performed.
No tooling was installed. Final checks PASS: entire index unchanged; §8 exact text
and line endings unchanged; only these two tracked documents modified; local links
resolve; `git diff --check` clean; both documents have no editor diagnostics.

## 7. Rejections / things the user declined

- Do not restore the editor legend or its deleted test. Do not reopen settled Batch C
  approval questions, reimplement a committed phase, or re-verify a closed batch.
- Topology retirement/redesign (Batch K, review §2.3) is explicitly **out of scope** for
  the next batch. Note coordination implications only; do not start it.
- No opportunistic work: no GUI/PNG geometry consolidation, no DTO removal, no schema
  changes, no package renames, no allocation tuning, no output-path changes.
- Do not claim engine defaults or in-game outcomes the assistant did not observe. The
  positive result is the owner's report and covers road values only.
- The earlier §8 comparison mismatch was a PowerShell 5.1 text-decoding artifact, not a
  content change. Use explicit UTF-8 for Git stdout
  (`ProcessStartInfo.StandardOutputEncoding`) and for disk reads whenever verifying §8;
  no §8 content was ever altered.
- Optional observation, **not** assigned work: a service comment still names the legacy
  rebuild entry point rather than the current reconciliation finalizer. Cosmetic, out of
  scope, not a task.
- Owner chose full snapshot clearing, not arena recomputation; no remembered player
  count; no persistent guard preset identity. Do not reopen those alternatives.
- First plan review rejected unapproved direct GeneratorConfig rejection. The approved
  plan enforces two players at editor-state validation/application generation only;
  direct generator/provider compatibility behavior and topology algorithms stay intact.
- First plan review also required visible dirty/warning outcomes separate from count
  validation, pre-generation invalidation, safe selected-pointer replacement and ordered
  Default matching. Revised plan was approved; do not undo those safeguards.
- Discovery initially used unsupported lowercase model names and two nonexistent guessed
  paths; corrected display names/searches succeeded. These failures changed no files.

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

**No unresolved Batch E product questions or Phase 1 blockers.** Written implementation
approval is recorded; the owner committed the plan. Findings §1.5/§1.11/§1.12 remain
open until implementation, verification and owner commit. This phase did not fix them.

1. Read [AGENTS.md](../AGENTS.md), this handoff, the approved
  [Batch E plan](plans/batch-e-effective-modes-and-guards.md) and relevant findings in
  the [surviving review](backlog/review-gpt-6-astra-09-07.md).
2. Inspect current worktree/index read-only, respecting any owner changes made since
  this snapshot. Confirm source/toolchain still matches the recorded 75.1% baseline.
  Do not rerun discovery/plan review or refresh unchanged evidence without cause.
3. Begin **Phase 2: effective modes and two-player state lifecycle**. Mark it In progress
  before editing. Read targets and callers; implement/tests together. Invalidation
  occurs before generation; count correction alone must not silently discard edits.
  Load reads once through the existing non-fixing path, followed by fixing validation
  and a model-based correction outcome; warnings must survive the next frame.
4. Complete/verify Phase 2 and record its summary before starting Phase 3 guard
  propagation/Custom UI. Preserve pointer ownership and Apply/Cancel boundaries.
  Phase 4 requires combined verification, coverage comparison and independent review.
5. Do not pull in Batch K, GUI/PNG geometry, DTO cleanup, schema or package work. Only
  narrow functional quality-request/response plumbing is approved. Regenerate Wire
  when constructor/provider changes require it, never by hand.
6. Preserve §8 verbatim; apply the supersession note above. Carry the plan's remaining
  work forward rather than marking the batch complete after one phase.

**Deployment plan.** Nothing deployed. Phase 1 is a pre-change baseline and documentation
closeout. Future implementation follows the approved plan's verification/deployment
sections; the owner alone stages, commits and releases. No schema migration, dependency
installation or output-directory change is planned. Native Linux/Steam Deck execution
remains unmeasured in this session.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first, then [this handoff](session-carry-forward.md),
> the approved [Batch E plan](plans/batch-e-effective-modes-and-guards.md) and relevant
> findings §1.5/§1.11/§1.12 in the
> [surviving review](backlog/review-gpt-6-astra-09-07.md).
>
> **Resume Batch E at Phase 2.** Scope and written plan are owner-approved; the plan
> was committed at `b9dd616` on `AD/modes_and_guard_propagation`. Phase 1 is complete:
> initially clean worktree/index, fresh 186-package unit PASS, Windows Go 1.27.0,
> total statement coverage 75.1%, with target baselines in the plan. No Batch E code
> has been implemented. Check current source/toolchain and inherited changes read-only
> before using that baseline; refresh only if inputs changed or evidence is unavailable.
> Do not repeat settled questions or the approved plan review. Mark Phase 2 In progress,
> implement its tests and behavior, verify it, and record its summary before moving on.
>
> **Settled scope:** effective tournament/arena transitions clear the full manual
> snapshot before generation, warn and dirty/re-arm Exit; alias-only changes preserve
> edits. No arena reconciliation. Quality changes carry exact-matched old-table guard
> presets to the new stronger-endpoint table on all incident edge types; unmatched
> numbers stay Custom. Keep Plastic/Bronze tables unchanged. Tournament forces/locks
> two players and off stays two; invalid tournament loads correct to two, discard
> snapshots, warn and become unsaved without rewriting the file. Valid two-player
> loads retain edits. Preserve warnings through regeneration and GUI pointer ownership.
>
> Batches A/D, B and C are closed, with Batch C's road engine acceptance reported by
> the owner. Do not retrieve retired documents or reopen/re-verify those batches.
> Batch K topology retirement, direct GeneratorConfig rejection, geometry consolidation,
> DTO cleanup, schema and package work are out of scope. Normal regression suites for
> actual Batch E changes still apply.
>
> **Hard rules:** never modify protected data, template schema or registry. Keep
> Windows/Linux compatibility. Never change or persist the machine-detected output
> directory. Test nontrivial logic and check before/after coverage. Never stage, unstage,
> commit, push, stash, switch branches or manipulate worktrees; preserve owner changes.
> Never bulk-rewrite or hand-edit generated Wire. Never enable global integration_test,
> gui or wireinject tags, or introduce fake unit seams. Keep the plan durable/resumable.
>
> Preserve explicit Portal Road false/nil and valid approaches, independent internal
> roads, nil-state content/road preservation, source cloning, one reusable local
> half-opacity edge mask over unchanged opaque rasterization, and the Preview-only legend.
> Handoff §8 stays verbatim; its old phase/engine and pending Batch E wording is
> superseded, while unrelated later-scope decisions remain binding. This handoff and
> the approved plan contain the full continuation context.
