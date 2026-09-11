# Carry-forward: Batch C closed, Batch E discovery is next

Date: 2026-09-11.

**Supersession note.** This revision replaces every earlier handoff, including the one
that described Phase 5 as open. **Batch C is fully closed.** Two things in §8 below are
now obsolete: its phase-status wording ("Phase 3 is complete", "Phase 4 remains
unstarted") and its pending engine-validation paragraph. The owner ran the generated
templates in the game on 2026-09-11 and reported: "Everything in-game looks good, the
road values are correct in engine." That closed the last open question of the batch.
Section 8 is nevertheless kept **verbatim** by owner request, contradictions included —
do not edit it. Its **retained later-scope decisions remain binding**; only the two
items named above are superseded. Batches A/D, B and C are all closed: no revisit, no
re-verification, no lookup of retired documents.

## 1. Session goal

Close out Batch C Phase 5 after the owner's in-game confirmation, then prepare the
handoff and a start prompt for the next batch. Batch-wide verification and the
independent whole-batch review were already green; the owner's engine report was the
only remaining gate and it came back positive.

**No production code, tests, snapshots or generated files changed.** This session was
documentation and handoff only.

## 2. Fixes applied

All four implementation phases are owner-committed on `AD/road_and_graph_invariants`:

- Phase 1 — `07ca8b6`: explicit-Portal and road classification on the connection model;
  endpoint-based, road-gated public topology repair that keeps graph repairs, guards,
  the impossible-isolation fallback and collision-safe edge names.
- Phase 2 — `abf0d36`: `RoadPolicyService` reconciliation for authoritative non-Portal
  road flags, zone-scoped mandatory-content validation, eligible approach restoration
  and invalid generated-road removal.
- Phase 3 — `a574a9e`: pending editor create/type changes obey the road checkbox before
  Apply through the handler/service layer; shared road-state colors, Preview legend and
  selected-edge width; road display state kept separate from effective preview type.
- Phase 4 — `0a55306`: PNG roadless strokes at 50% once per edge using one reusable
  local `image.Alpha` mask, with the opaque raster path pinned to pre-change
  fingerprints.

Phase 5 applied no fixes. The independent review raised no blocker, and the owner's
engine result required no corrective change either.

## 3. Features added / changed

Settled behavior contract, verified by tests, by independent review and now by the
owner's in-game check — **preserve all of it**:

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

Documentation only. Nothing in the source tree was created, edited or deleted:

- The completed Batch C plan: closed out with its final phase status and verification
  record. It is deliberately not linked or named here — that completed batch record is
  retained for owner retirement, and the next batch must not be pointed at it.
- [review backlog](backlog/review-gpt-6-astra-09-07.md): the surviving review, and the
  document the next batch works from.
- [handoff](session-carry-forward.md): this document; §8 retained verbatim.
- Ignored [settled decisions](memories/settled-decisions.md): final visual/policy
  decisions, now including the owner's positive engine result.

No Go files, tests, snapshots, goldens, assets, dependencies, generated Wire or
protected trees changed.

## 5. Tests added or updated

No tests were added or modified. The Batch C verification matrix stands as recorded on
Windows, 2026-09-11; nothing was re-run in this session because nothing changed:

- Fresh full unit coverage with `-count=1`: **PASS, 75.1%**, matching the Phase 4
  baseline. Every function changed by the batch is at **100%**.
- `go build ./...`: **PASS**.
- `go test ./test/...`: **PASS**.
- `go test -tags=integration_test ./test/integration/...`: **PASS**.
- `go test -tags='integration_test,gui' ./test/integration/...`: **PASS** with actual
  execution — GUI **29.469 s**, integration **2.899 s**.
- `go run ./cmd/testlayoutcheck .`: **PASS**. Report-only lint: **0 issues** plus the
  three existing unused-exclusion warnings. `gofmt -l` clean on every changed Go file.
- Independent Claude Opus 5 whole-batch review: **APPROVED, no blockers**. That review
  is complete — do not schedule a rerun without an actual code change to justify it.

The engine evidence, stated precisely: acceptance is an **owner report** from running
the game. The assistant ran nothing in-game, built nothing on Linux or a native Steam
Deck target, and this is **not** evidence that in-game bonuses/bans or hero-hire
behavior is fixed. Those remain open in the review.

## 6. Git status snapshot

Branch `AD/road_and_graph_invariants`, HEAD `0a55306`. Phases 1–4 are owner-committed at
`07ca8b6`, `abf0d36`, `a574a9e` and `0a55306`. The Phase 5 documentation is **staged by
the owner**; this session's further documentation edits sit on top of it, unstaged.

The agent performed no Git mutation of any kind: no staging, unstaging, commit, push,
stash, branch or worktree operation, and the owner's index was left exactly as found.
No tooling was installed.

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

The next batch is **Batch E**, taken from §9 of the surviving review,
[review-gpt-6-astra-09-07.md](backlog/review-gpt-6-astra-09-07.md): findings **§1.5**,
**§1.11** and **§1.12** — effective tournament/arena manual invalidation, incident guard
preset recalculation after quality changes, and enforcing two players for tournament
mode. There is **no Batch E plan yet and no implementation approval**; do not invent
either.

1. Read [AGENTS.md](../AGENTS.md), this handoff, and the surviving review at
  [backlog/review-gpt-6-astra-09-07.md](backlog/review-gpt-6-astra-09-07.md).
2. Treat the next session as **read-only discovery**. Inspect how the arena/manual-edit
  path, the guard preset recalculation path and the player-count path behave today.
  Change nothing and add nothing; run only read-only inspection and, if useful, the
  existing verification commands.
3. Resolve product questions with the owner: §1.5 manual-edit invalidation versus arena recomputation and effective-mode aliases; §1.11 custom guard values, preset identity and stronger opposite endpoints; §1.12 restoring the previous non-tournament player count and correction/validation of loaded 3–8 player tournament states.
4. Only after the owner confirms scope, write a durable plan under `.agent/plans/` and
  obtain explicit owner approval before writing any code.
5. Out of scope for Batch E, do not pull it in: topology retirement/redesign (Batch K,
  review §2.3), which is coordination-implications-only; and any opportunistic GUI/PNG
  geometry consolidation, DTO, schema or package work.
6. Preserve §8 verbatim. Its Phase 3/Phase 4 status wording and its engine-validation
  paragraph are obsolete; its later-scope decisions are still binding.

**Deployment plan.** Nothing is deployed and none is claimed. The Windows source build
and the full test matrix pass; Linux and native Steam Deck builds were not run. Batch C
needs no schema migration, no Wire regeneration, no dependency installation and no
output-path change. The owner confirmed the in-game result and performs all staging,
commits, tagging and release.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first, then [this handoff](session-carry-forward.md)
> and the surviving review at
> [backlog/review-gpt-6-astra-09-07.md](backlog/review-gpt-6-astra-09-07.md).
>
> **Batch C is fully closed** as of 2026-09-11: Phases 1–4 are owner-committed on
> `AD/road_and_graph_invariants` (`07ca8b6`, `abf0d36`, `a574a9e`, `0a55306`, HEAD
> `0a55306`), and the owner confirmed in-game that the road values are correct in the
> engine. Batches A/D and B are closed too. Do not revisit, re-verify, reimplement or
> look up any retired batch document; the completed batch record is retained for owner
> retirement and is not your working document.
>
> **Next up is Batch E**, review §9: findings §1.5, §1.11 and §1.12 — effective
> tournament/arena manual invalidation, incident guard preset recalculation after
> quality changes, and enforcing two players for tournament mode. **There is no Batch E
> plan and no implementation approval.** Your session is **read-only discovery plus
> product questions**, nothing more: inspect current behavior, then settle with the
> owner whether an effective-mode change clears manual edits or recomputes arena
> placement (and how aliases and effective-mode transitions behave), what preset
> identity means once guard values are customized and how a stronger opposite endpoint
> affects recalculation, and whether turning tournament mode off restores
> the remembered previous non-tournament player count and how a loaded 3–8 player
> template is corrected or rejected. Confirm scope, then write a durable plan and get
> explicit owner approval before writing any code.
>
> Out of scope: topology retirement/redesign (Batch K, review §2.3), coordination
> implications only; and any opportunistic GUI/PNG geometry consolidation, DTO, schema
> or package work.
>
> **Preserve the settled behavior:** explicit Portal flags including `false` and `nil`
> plus valid Portal approaches; internal roads independent of the between-zone checkbox;
> roads and content preserved on a `nil` state; sources cloned before mutation; one
> reusable local mask applying 50% once per edge over an unchanged opaque raster path;
> the legend on Preview only; and a machine-detected game output directory that is never
> persisted.
>
> **Hard rules:** never touch protected `data/`, the template schema or the registry;
> keep everything cross-platform; never change or persist the output directory; test
> nontrivial changes and check coverage; never stage, unstage, commit, push, stash or
> switch branches, and leave the owner's staged changes alone; never bulk-rewrite or
> hand-edit generated Wire; never enable global `integration_test`, `gui` or
> `wireinject` tags, and never add fake unit seams; keep plans durable and resumable.
> Section 8 of the handoff is verbatim historical text — its phase-status and
> engine-validation wording is superseded, its later-scope decisions are not.
