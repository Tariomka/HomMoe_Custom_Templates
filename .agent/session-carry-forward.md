# Session Carry-Forward — Geometric Hub fix round (PLANNING session, 2026-07-15)

## 1. Session goal

Resume after the Geometric Hub / zone-tiers implementation session, collect the
user's review findings, and produce a detailed execution plan for the fix round.
**Plan-only session — NO implementation was done** (user's explicit instruction:
"write the plan in detail but do not start implementing").

## 2. Fixes applied

None (planning session). No production or test code was touched.

## 3. Features added / changed

- Created [plans/geometric-hub-fixes.md](../plans/geometric-hub-fixes.md) — the
  5-phase execution plan for the fix round, containing every user-confirmed
  decision, the interior-polygon connection spec (k=1..k≥5 table), the
  regular-hexagon geometry spec, the reference-image table, the code map, and
  per-phase verification commands.
- Updated repo memory (`/memories/repo/conventions.md`, GeometricHub section)
  with a "FIX ROUND PLANNED" entry summarizing the confirmed decisions.

## 4. File modifications (this session)

- `plans/geometric-hub-fixes.md` — CREATED (the plan; untracked `??` in git).
- `/memories/repo/conventions.md` — memory update (not a repo file).
- Nothing else. All other dirty files in `git status` are the PREVIOUS
  session's uncommitted Geometric Hub + tiers implementation.

## 5. Tests added or updated

None this session. Last known suite state (from the previous session, verified
on this workspace): `go test ./test/... -count=1` fails ONLY in
`test/unit/app/gui/utils/buttonPositionLogger` (5 tests, pre-existing on clean
HEAD, caused by the owner's intended extra `"====== New Frame ======"` log
record — fixing the expectations is Phase 1 of the new plan). Coverage 63.5%,
lint at the 84-gochecknoglobals baseline.

## 6. Git status snapshot

Branch: `AD/refactoring-07-13`. ~57 staged/modified files + 3 new PNG assets +
2 new topology source files + 3 new topology test files + `tools/platinumgen/`
— ALL from the previous session's (still uncommitted) Geometric Hub + tiers
feature. New untracked: `plans/geometric-hub-fixes.md`. `todo/promt.md` has
unstaged owner edits — do not touch. The next session inherits this dirty tree;
the fix round builds on top of it (the feature is NOT yet committed).

## 7. Rejections / things the user declined

- My proposed "chord model" for k≥3 interiors (middle interiors touching no
  stable) — REPLACED by the user's regular-k-gon spec (see plan).
- Implementing anything this session — deferred to the next session.
- Corner-split slot class — removed from the design entirely (was part of the
  original confirmed spec, now obsolete).

## 8. Open questions

None blocking. All scope questions were resolved via four interview rounds;
the plan's "User-confirmed decisions (do not re-ask)" section is exhaustive.
Implementation-time freedoms explicitly delegated to the implementer: exact
geometry constants (tuned visually against the reference PNGs) and the k-gon
circumradius.

## 9. Next recommended actions

1. Execute [plans/geometric-hub-fixes.md](../plans/geometric-hub-fixes.md)
   Phase 1 (buttonPositionLogger test expectations — quick, makes the whole
   suite green).
2. Phase 2 (platinum arena sprite — gold swords composite).
3. Phase 3 (growth ladder + interior polygon connection graph) — the big one;
   re-spec the topology tests, validate `-count=20`.
4. Phase 4 (regular-hexagon geometry + per-player-count scale) — visual
   verification against `output/implementation/*.png`.
5. Phase 5 (final verification, repo-memory update, user GUI smoke).
6. Still open from the PREVIOUS plan (optional): independent model review of
   the topology/profile changes (plans/geometric-hub-topology.md Phase 6).

## 10. Carry-forward prompt (paste into the next session)

> Read `AGENTS.md` first and follow it. Hard rules in one sentence each: never
> modify `data/`, `internal/entities/template/` or `internal/registry/`
> (read-only game data); keep everything cross-platform (Windows + Linux,
> `path/filepath`, PowerShell `;` chaining); every code change ships with unit
> tests under `test/unit/` mirroring the impl path and must not drop coverage.
>
> **Task**: EXECUTE the plan in `plans/geometric-hub-fixes.md` — read it in
> full before doing anything; it is the source of truth and is self-contained
> (user-confirmed decisions marked "do not re-ask", interior-polygon
> connection spec table, regular-hexagon geometry spec, reference-image table,
> code map, per-phase verification commands). Work phase by phase, marking
> checkboxes, setting phase Status, writing each Phase Summary, and running
> each phase's Verification Plan before moving on.
>
> **The 5 phases in one line each**: (1) fix the 5 buttonPositionLogger test
> expectations (the extra "====== New Frame ======" record is intended; do
> NOT change the logger); (2) regenerate `neutral_highest_arena.png` with GOLD
> swords by compositing the `gladiator_arena.png` master @0.95 onto the
> platinum bubble in `tools/platinumgen`; (3) remove corner splits from the
> Geometric Hub growth ladder and rebuild interior connections onto the
> regular-k-gon spec (only x1/x2 hub-portal; k=3's x3 connects BOTH stables —
> 4-link exception); (4) rework `geometricHubLayout.go` positions onto
> regular-hexagon math (corner=s, stable=√3·s, player=2·s; ±30°/±60° ideal
> angles, sector-fraction fallback for P≥5) with players CLOSER to the hub for
> 2–4P, tuned visually against `output/implementation/*.png` (the 180°
> rotation vs `One for All.png` is deliberate — keep it); (5) full
> verification: build, unit + gated integration/performance suites, coverage
> no-drop vs 63.5%, lint at the 84-gochecknoglobals baseline, repo-memory
> update, then hand the GUI smoke test to the user.
>
> **Context you inherit**: branch `AD/refactoring-07-13` has the ENTIRE
> previous Geometric Hub + Lowest/Highest tiers feature still uncommitted
> (~60 files; see `plans/geometric-hub-topology.md` for what it is) — the fix
> round builds on top of it. The only failing tests repo-wide are the 5
> buttonPositionLogger ones that Phase 1 fixes. Repo memory
> (`/memories/repo/conventions.md`) has a "GeometricHub FIX ROUND PLANNED"
> entry mirroring the key decisions. Full handoff details in
> `./.agent/session-carry-forward.md`.
>
> Start with Phase 1 and proceed autonomously; only stop for the Phase 5 user
> smoke test or if a verification gate fails in a way the plan doesn't cover.
