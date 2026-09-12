# Carry-forward: Batch E complete, awaiting the owner's commit

Date: 2026-09-12. Authority: the owner-approved
[Batch E plan](plans/batch-e-effective-modes-and-guards.md), which carries the full
product contract, implementation map and verification ledger. This handoff is the current
status record; it replaces every earlier status statement this file used to carry.

**Section 8 is preserved verbatim**, contradictions and all. Its Batch C phase/engine
wording and its closing "pending Batch E questions (§1.5/§1.11/§1.12)" wording are
**superseded** - those questions were settled by the approved plan and are now
implemented. Everything else §8 retains as later scope remains **binding**. Do not edit §8.

Batches A/D, B and C are closed. Batch C's engine gate closed on the owner's 2026-09-11
report: "Everything in-game looks good, the road values are correct in engine." Do not
revisit, re-verify or retrieve retired documents from those batches.

## 1. Session goal

Finish Batch E **Phase 4** and update this handoff.

Phase 4 is complete, so **the whole batch is complete**. The combined cross-phase flow is
now covered by six real-input GUI tests, every Windows gate passes, native Linux is
recorded unavailable, and the owner accepted the implementation on 2026-09-12:
*"Changes reviewed, all seems fine, you can finish up with Phase 4 and update Session
Carry Forward"*. The assistant made no Git change; the batch closes on the owner's commit.

## 2. Fixes applied

All three Batch E findings now have implemented, tested and independently reviewed
behavior. **None is marked fixed in the surviving review** - the review's protocol
requires the owner's commit of this final verification and documentation work first.

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

Exactly **three** files differ from owner commit `1b4658c`, and all three are unstaged.
Nothing was staged or committed by the assistant.

- [test/integration/gui/zoneEditorProperties_integration_test.go](../test/integration/gui/zoneEditorProperties_integration_test.go),
  **+171 lines**: six combined cross-phase tests, described in §5. The only code change.
- [Batch E plan](plans/batch-e-effective-modes-and-guards.md): authority block updated to
  all four phases complete, Phase 4 status/checkboxes/summary written, the batch-final
  ledger column filled with measured results, Final Recap and Deployment Plan completed.
  Phase 1/2/3 evidence is preserved as historical.
- [Handoff](session-carry-forward.md): this record, with §8 preserved verbatim.

Ignored [coverage.txt](../coverage.txt), [coverage.html](../coverage.html) and
[lcov.info](../lcov.info) were regenerated by the coverage run. They are local artifacts,
not tracked files.

**Nothing else was touched:** no production code, no protected data, template schema or
registry, no output path, no opaque raster path, no generated Wire, no topology code, no
dependencies, snapshots or goldens, and no `.agent/backlog/` review or owner-findings
document. No new memory or summary files were created.

The production files delivered by earlier phases are already committed at `1b4658c` and
are inventoried in the plan's implementation map and phase summaries. That inventory is
not repeated here.

## 5. Tests added or updated

**This turn:** six combined real-input tests in
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
21 real-input GUI cases and two persistence regressions. Those are catalogued in the plan
and are not repeated here. No fake unit seams, no new `*_testexports.go`, no global tags.

**Measured this turn, Windows/amd64, Go 1.27.0, empty `GOFLAGS`:**

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
  output and exit code 1 this turn. Nothing was installed; no native Linux, Steam Deck or
  engine result is claimed.
- **No coverage profile fingerprint is claimed.** The reports were regenerated but their
  SHA-256 was not measured this turn. The Phase 3 value
  `23B78F2CF15B73902398D8A21A755FD684D3668486E56D8CFA770BE2BE7F9F31` is historical.
- Independent **Claude Opus 5 review of the six new tests: APPROVED**. The plan, Phase 2
  and Phase 3 implementation reviews were approved earlier.

## 6. Git status snapshot

Branch **`AD/modes_and_guard_propagation`**, in sync with
`origin/AD/modes_and_guard_propagation`. HEAD is the owner's commit **`1b4658c`
("Phase 3 and partial 4")**. The turn started clean.

`git status --short` currently reports a single entry:

```text
 M test/integration/gui/zoneEditorProperties_integration_test.go
```

Plus the two untracked-to-review documents in §4, which live under `.agent/` and are
modified in place. Nothing is staged; the staged diff is empty. **The assistant performed
no Git mutation of any kind** - no staging, unstaging, commit, push, stash, branch switch
or worktree change - and every owner change is preserved exactly as found.
`.agent/backlog/owner_findings.md` and the surviving review were not read, edited or
re-indexed.

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
- Topology retirement/redesign (Batch K, review §2.3) is explicitly out of scope. No
  opportunistic work either: no GUI/PNG geometry consolidation, no DTO removal, no schema
  changes, no package renames, no allocation tuning, no output-path changes.
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

**Batch E is complete, with no open product questions and no blockers.** What remains is
the owner's own commit protocol.

1. **Owner:** review the three unstaged changes listed in §4, then stage and commit them.
2. **Only after that commit exists**, mark findings §1.5, §1.11 and §1.12 in the
   [surviving review](backlog/review-gpt-6-astra-09-07.md), following that document's own
   marking protocol. Do not edit the review before the commit, and never read, edit or
   re-index `.agent/backlog/owner_findings.md`.
3. Anything past that is a new batch. Take the next item from §8's retained later-scope
   list and plan it separately; do not fold it into Batch E.
4. If Batch E ever needs re-verification, rerun the batch-final command list in the plan's
   verification ledger. Native Linux stays unavailable until a machine with `go` and
   `pkg-config` exists; record it as unavailable, never as passing.

**Deployment.** Nothing is deployed and nothing is authorized to be. The owner alone
stages, commits and releases. No schema migration, dependency installation, Wire
regeneration or output-directory change is required. Native Linux and Steam Deck
execution remains unmeasured.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first, then [this handoff](session-carry-forward.md) and
> the approved [Batch E plan](plans/batch-e-effective-modes-and-guards.md).
>
> **Batch E is COMPLETE.** All four phases are done, the owner accepted the implementation
> on 2026-09-12, and findings §1.5, §1.11 and §1.12 all have implemented, tested and
> independently reviewed behavior. Do not reimplement, re-verify or re-review any phase,
> and do not repeat settled product questions.
>
> **State:** branch `AD/modes_and_guard_propagation`, in sync with origin, HEAD at the
> owner's `1b4658c` ("Phase 3 and partial 4"). Exactly three files are unstaged: the GUI
> test `test/integration/gui/zoneEditorProperties_integration_test.go` (+171 lines, six
> combined cross-phase tests), plus the plan and this handoff. Nothing is staged and the
> assistant performed no Git mutation.
>
> **Next action:** the owner commits those three files. Only then mark §1.5, §1.11 and
> §1.12 in the [surviving review](backlog/review-gpt-6-astra-09-07.md) per its protocol.
> Never touch `.agent/backlog/owner_findings.md`.
>
> **Verification on record, Windows/amd64, Go 1.27.0, empty `GOFLAGS`:** `go build ./...`;
> `go test -p=2 -count=1` unit run with `-coverpkg=./internal/...,./app/...`;
> `go tool cover -func` at **74.9%**, an owner-accepted exception rather than a restored
> 75.1% baseline; `go test ./test/...`; `go test -tags='integration_test,gui'
> ./test/integration/...` with root 3.297s and GUI 28.558s; `go run ./cmd/testlayoutcheck .`;
> clean `gofmt -l` on the changed file; and zero-issue
> `golangci-lint-v2 run ./... --issues-exit-code=1`. Generated Wire is already current
> because no constructor changed - do not rerun the generator needlessly. **Native
> Linux/Steam Deck is UNAVAILABLE:** the WSL probe found no `go` and no `pkg-config`,
> nothing was installed, and no native or in-game result may be claimed. No coverage
> profile fingerprint was measured this turn.
>
> **Hard rules:** never modify protected data, the template schema or the registry. Keep
> Windows/Linux compatibility. Never change or persist the machine-detected output
> directory. Test nontrivial logic and check before/after coverage. Never stage, unstage,
> commit, push, stash, switch branches or manipulate worktrees; preserve owner changes.
> Never bulk-rewrite or hand-edit generated Wire. Never enable global `integration_test`,
> `gui` or `wireinject` tags, and never introduce fake unit seams. Keep plans durable and
> resumable.
>
> **Out of scope:** Batch K topology retirement, direct `GeneratorConfig` rejection,
> geometry consolidation, DTO cleanup, schema and package work, and every settled
> alternative in §7. §8 stays verbatim: its Batch C phase/engine wording and its closing
> "pending Batch E questions" wording are superseded, while its retained later-scope
> decisions remain binding and are the source for the next batch.
>
> Preserve explicit Portal Road false/nil and valid approaches, independent internal
> roads, nil-state content/road preservation, source cloning, one reusable local
> half-opacity edge mask over unchanged opaque rasterization, and the Preview-only legend.
> This handoff and the approved plan contain the full continuation context.
