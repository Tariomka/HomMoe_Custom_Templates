# Geometric Hub Fixes: regular hexagons, interior polygons, platinum arena, logger tests

Post-review fix round for the Geometric Hub topology delivered by
[plans/geometric-hub-topology.md](geometric-hub-topology.md) (read its
"Phase Summary" sections for implementation facts — slot model, plan
assignment, file map). Scope confirmed with the user on 2026-07-15;
all decisions below are user-confirmed — do NOT re-ask.

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

Read `AGENTS.md` first. Hard rules: never touch `data/`, `internal/entities/template/`,
`internal/registry/` (read-only game data); stay cross-platform (Windows+Linux,
`path/filepath` only); every change ships with unit tests under `test/unit/`
mirroring the impl path (§4.6) and must not drop coverage (§2.3); receivers are
named `this`; one struct per file, camelCase file names. Repo gotchas: .go files
are LF; `tmp/` is deleted by the user's `air` watcher (scratch tools go under
`cmd/`, delete after); sprites use `image.NRGBA` (premultiplied `color.RGBA`
corrupts edges); topology tests must be validated with `-count=20`; `funcorder`
lint places unexported methods below exported ones.

### User-confirmed decisions (do not re-ask)

1. **buttonPositionLogger**: the extra `"====== New Frame ======"` record emitted
   once per `LogButtonPositions` call is INTENDED behavior (owner added it).
   Update the 5 failing test expectations; do NOT change the logger.
2. **Platinum arena sprite**: swords must stay GOLD. Build
   `neutral_highest_arena.png` by compositing the gold `gladiator_arena.png`
   master onto the recolored platinum bubble (same 0.95-scale technique the
   other `neutral_*_arena.png` variants used), instead of recoloring
   `neutral_high_arena.png` wholesale.
3. **Corner splits removed ENTIRELY**. New growth ladder: hub → stables
   (2/player) → merged corners (1/gap) → interiors (unbounded, round-robin per
   hexagon). Delete split-corner code, constants and tests. Degradation ladder
   below full stables (shared stables / hub-only star) is unchanged.
4. **Interior polygon scheme** (replaces the old rule-8 chaining entirely — no
   "one unused stable + all earlier interiors" logic remains). See spec below.
5. **Regular-hexagon geometry**: derive positions from true regular-hexagon
   ratios (corner=s, stable=√3·s, player=2·s from hub), exact for 3P, best
   approximation otherwise; implementer tunes constants VISUALLY against the
   reference PNGs. Player zones must move CLOSER to the hub for 2–4 players;
   the current distance is good for ≥5 players.
6. **Orientation**: the current implementation is rotated 180° vs the
   `One for All.png` inspiration — this is DELIBERATE and must be kept. When
   comparing 1:1 against that inspiration image, flip it 180°.

### Reference images (in `output/implementation/`, with matching .rmg.json)

| File | What it shows |
|------|---------------|
| `One for All.png` | Inspiration (its center "player 1" plays the hub role; 180° rotated vs our layout). Perfect regular hexagons + shared corner zones (empty circles). |
| `test layout One for All.png` | CURRENT 3P output — hexagon angles irregular (finding 1), players too far from hub (finding 2). |
| `3Hexes-2HexCentralZones.png` | CURRENT 3P + 15 neutrals — WRONG: corner splits fired instead of a 2nd interior. |
| `3Hexes-2HexCentralZones-RoughExpectation.png` | EXPECTED 3P + 15 neutrals: 3 hexagons with merged corners + 2 interiors per hexagon (k=2 chain). |
| `3Hexes-2HexCentralZones-Plus1Central.png` | CURRENT (broken) 16-neutral output: crossing edges, interiors with 4 links, unequal spacing. |
| `4 Player Hub.png`, `5 Player Hub.png`, `6 Player Hub.png`, `Gladiator Hub.png` | User's hand-made shape references for higher player counts. |

### Interior polygon spec (k = interiors in one hexagon)

Placement: the k interiors form a regular k-gon centered at the hexagon center
(midpoint of hub→player segment, i.e. radius s along the player axis for a
regular hexagon with hub vertex at map center). Orientation: the x1—x2 edge
faces the HUB — vertex angles relative to the hub direction (as seen from the
polygon center) are `±π/k + m·2π/k`. Spacing must come out visually equal
(user finding: current 2 interiors are not equidistant from each other and
from the stables they connect to). Circumradius: implementer-tuned so
sL—x1 ≈ x1—x2 ≈ x2—sR for k=2 and the k-gon reads as a clean shape inside
the hexagon.

Connections (sL/sR = the hexagon's two stable zones; ring edges are Direct,
hub links are Portal):

| k | Edges | Hub portals |
|---|-------|-------------|
| 1 | sL—x1, x1—sR (splits hexagon into 3 rhombuses; unchanged rule 7) | x1 |
| 2 | sL—x1, x1—x2, x2—sR; x1,x2 evenly spaced on the sL→sR chord, nearer hub | x1, x2 |
| 3 | Triangle ring x1—x2, x2—x3, x3—x1; sL—x1, sL—x3, x2—sR, x3—sR. x3 sits toward the player and has 4 connections (BOTH stables) — the one allowed exception; x3 does NOT touch the hub | x1, x2 only |
| 4 | Square ring x1—x2, x2—x4, x4—x3, x3—x1 (NO diagonals: never x1—x4 nor x2—x3); sL—x1, sL—x3, x2—sR, x4—sR. x3,x4 form the player-side chain sL—x3—x4—sR | x1, x2 only |
| ≥5 | Regular k-gon ring edges only (no diagonals). sL connects x1 + its nearest non-hub-side vertex; sR connects x2 + its nearest non-hub-side vertex. When one vertex is nearest for both stables (odd small k, e.g. k=3), it connects to both | x1, x2 only |

x1 = hub-side-left vertex, x2 = hub-side-right. Interiors still take the
highest-quality plans (pop order unchanged); the label→slot order within
`hexInteriors[i]` must be made deterministic and documented (which slot index
is x1, x2, x3…) so tier placement is reproducible.

### Regular-hexagon geometry spec

Hexagon i (player i at angle `θ_i = -π/2 + i·2π/P`): vertices Hub(r=0),
cL(r=s, θ−60°), sL(r=√3·s, θ−30°), P(r=2s, θ), sR(√3·s, θ+30°), cR(s, θ+60°),
where s = hexagon side. Shared merged corner between hexagons i,i+1 sits at the
gap mid-angle `θ_i + sector/2` at radius s (exact for P=3 where mid = θ+60°;
an approximation otherwise). Angular offsets: use the absolute ±30°/±60° ideal
when it fits inside the sector; for P≥5 the sector (≤72°) is too narrow —
fall back to sector-fraction offsets (current behavior, which the user said
looks good for >4P). Suggested form: `offset = min(idealAbsolute,
fraction·sector)` — implementer's choice, tuned visually.

Scale: player radius = 2s. For 2–4 players REDUCE the current
`geoHubPlayerRadius` (0.46) so players sit closer to the hub (compare
`test layout One for All.png` vs the flipped `One for All.png`); for ≥5
players keep ≈0.46. A per-player-count scale function (or two constants with
a P threshold) is acceptable. Derived radii follow the ratios: stable =
player·(√3/2) ≈ 0.866·player; corner = player·0.5; hexagon center (interior
polygon center) = s along θ = player·0.5 as well (same radius as the corners,
but on the player axis instead of the gap angles). Verify all against the
images, not just the math: the rendered shape must read as regular hexagons
for 3P.

Current constants being replaced (top of
[geometricHubLayout.go](../internal/services/template_generator/providers/topology/geometricHubLayout.go)):
`geoHubPlayerRadius=0.46, geoHubStableRadius=0.41, geoHubStableAngleFraction=0.35,
geoHubMergedCornerRadius=0.22, geoHubSplitCornerRadius=0.27,
geoHubSplitAngleFraction=0.15, geoHubInteriorRadius=0.24,
geoHubExtraInteriorRadius=0.28, geoHubExtraAngleFraction=0.16` — the three
split/extra-fan constants disappear with their features.

### Code map (everything lives in these files)

- [internal/services/template_generator/providers/topology/geometricHubLayout.go](../internal/services/template_generator/providers/topology/geometricHubLayout.go)
  — ALL changes concentrate here: `distributeGeometricHubSlots` (drop the
  cornerSplits round), `assignPlans` (corner rounds 2→1; interior slot→vertex
  order), `computePositions` (regular-hexagon math, polygon vertices, per-P
  scale), `buildEdges`/`buildInteriorEdges` (polygon ring + stable links +
  x1/x2-only portals), `hexagonStables` (still needed), `extraInteriorAngleOffset`/
  `extraInteriorRadius` (replaced by polygon vertex math).
- [geometricHubTopology.go](../internal/services/template_generator/providers/topology/geometricHubTopology.go)
  — consumes `directEdges`/`hubPortalLabels`/`positions`; NO changes expected
  (verify only).
- Tests: [test/unit/internal/services/template_generator/providers/topology/geometricHubTopology/](../test/unit/internal/services/template_generator/providers/topology/geometricHubTopology/)
  (16 tests; corner-split tests get deleted/re-specced, polygon tests added).
- [tools/platinumgen/main.go](../tools/platinumgen/main.go) — arena composite
  change (Phase 2). NOTE: `tools/assetgen` no longer exists; platinumgen must
  implement the 0.95-scale composite itself.
- [test/unit/app/gui/utils/buttonPositionLogger/logButtonPositions_test.go](../test/unit/app/gui/utils/buttonPositionLogger/logButtonPositions_test.go)
  — Phase 1. The logger ([app/gui/utils/buttonPositionLogger.go](../app/gui/utils/buttonPositionLogger.go))
  emits `this.logger.Debug("====== New Frame ======")` at the top of
  `LogButtonPositions` (line ~32) — stays as-is.

---

## Phase 1: buttonPositionLogger test expectations
Status: Complete

The 5 failures (`TestWhenOpsContainMultipleButtons_LogsEachButton`,
`TestWhenButtonHasNoLabel_SkipsButton`, `TestWhenOpsContainOffsetButton_LogsAbsoluteCenter`,
`TestWhenOpsContainLabeledButton_LogsLabel`, `TestWhenOpsContainNoButtons_LogsNothing`)
all break the same way: one extra leading record with message
`"====== New Frame ======"` per `LogButtonPositions` call.

- [x] In [logButtonPositions_test.go](../test/unit/app/gui/utils/buttonPositionLogger/logButtonPositions_test.go),
  add a small helper that filters the captured `recordingHandler.records` down
  to button records (message == "Button position"); switch the 5 tests'
  assertions onto the filtered slice (counts/emptiness stay as originally
  intended).
- [x] Add one NEW test asserting the intended frame-marker behavior: a single
  `"====== New Frame ======"` DEBUG record is emitted per `LogButtonPositions`
  call even when ops contain no buttons (name per §4.6, e.g.
  `TestWhenLoggingFrame_EmitsSingleNewFrameRecord`).
- [x] Do NOT touch [buttonPositionLogger.go](../app/gui/utils/buttonPositionLogger.go).

### Verification Plan
- `go test -count=1 ./test/unit/app/gui/utils/buttonPositionLogger/...` — all green. ✅ RESULT: ok, 1.988s.
- `go test ./test/... -count=1` — no other suite regressed; expected result:
  ZERO failing tests repo-wide (this was the only failing package). ✅ RESULT: zero FAIL lines, whole suite ok.

### Phase Summary
Added `buttonRecords(handler)` filter helper (message == "Button position") and
switched the 5 broken tests onto it; original count/emptiness semantics kept.
Added `TestWhenLoggingFrame_EmitsSingleNewFrameRecord` asserting exactly one
`"====== New Frame ======"` DEBUG record per call with empty ops. Logger
untouched. Full default suite now green repo-wide.

## Phase 2: Platinum arena sprite — gold swords
Status: Complete

- [x] Rework [tools/platinumgen/main.go](../tools/platinumgen/main.go):
  `neutral_highest.png` and `neutral_highest_castle.png` keep the current
  recolor path; `neutral_highest_arena.png` becomes recolor(`neutral_high.png`)
  + composite of `gladiator_arena.png` (96×96, solid/fully-opaque interior,
  gold swords) scaled to 0.95 and centered at (48,48), alpha-over. Use
  bilinear scaling and `image.NRGBA` end-to-end (NOT `color.RGBA` —
  premultiplication corrupts gold edges to green; see repo memory).
- [x] Regenerate the three `neutral_highest*.png` under
  [internal/services/asset_provider/assets/](../internal/services/asset_provider/assets/)
  (`go run ./platinumgen` from `tools/`, or from repo root — path detection
  handles both).
- [x] Visually verify (view the PNG): swords GOLD (matching
  `neutral_high_arena.png` / `gladiator_arena.png` sword color), bubble icy
  platinum, sword tips poking ~4px past the outer ring like the other arena
  variants; `neutral_highest.png`/`_castle.png` byte-identical or visually
  unchanged vs current.
- [x] Tests: arena assets are embedded but NOT loaded/drawn by the asset
  provider (consistent with the other `*_arena` assets — out of scope, same
  as the previous session). platinumgen is a `tools/` module outside
  `-coverpkg` — no unit tests required (registry §4.6 scope rules do not
  cover tools/), but note the tool change in the phase summary.

### Verification Plan
- `go build ./...` and `go test ./test/... -count=1` still green (assets are
  go:embed'd — a corrupt PNG would fail the existing decode tests for
  neutral_highest/neutral_highest_castle). ✅ RESULT: build clean, full suite green.
- Visual check of `neutral_highest_arena.png` (vision or the asciiview
  technique from repo memory). ✅ RESULT: gold crossed swords (same hue as
  `neutral_high_arena.png`) on icy platinum bubble, tips extending past the ring.

### Phase Summary
Split `main()` in [tools/platinumgen/main.go](../tools/platinumgen/main.go):
`neutral_highest.png`/`_castle.png` keep the recolor path;
`neutral_highest_arena.png` is now `composeArenaFile` = recolorToPlatinum(`neutral_high.png`)
+ `compositeScaledOver` of `gladiator_arena.png` at `arenaGlyphScale=0.95`,
centered at (47.5,47.5) pixel-center. New helpers: `decodePNG`, `toNRGBA`,
`sampleBilinear` (alpha-weighted bilinear, transparent outside bounds),
`pixelAt`, `blendOver` (source-over in non-premultiplied NRGBA — no
premultiplication anywhere, per repo gotcha). All three assets regenerated
(`go run ./platinumgen` from `tools/`; the tool is a separate module, so it
must run from `tools/` or repo root with its own module context — repo root
`go run ./tools/platinumgen` does NOT work). Visual + build + full-suite
verification passed. No unit tests: tools/ module is outside `-coverpkg`.

## Phase 3: Growth ladder + interior polygon connection graph
Status: Complete

Graph/topology changes only (positions follow in Phase 4). All in
[geometricHubLayout.go](../internal/services/template_generator/providers/topology/geometricHubLayout.go).

- [x] `distributeGeometricHubSlots`: remove the `cornerSplits` round; budget
  flows stables(cap 2/gap) → merged corners(cap 1/gap) → interiors
  (round-robin per hexagon, uncapped). `gapCorners[j]` now holds 0 or 1 label.
- [x] `assignPlans`: corner loop becomes a single round; interior pop order
  stays highest-first but define + document the slot→vertex mapping (which
  placement index is x1/x2/x3… — recommended: assignment order == x1..xk so
  the highest plans land on the hub-facing pair, matching "interiors nearest
  hub get the best plans" tier intent; confirm behavior in tests).
- [x] `buildInteriorEdges`: replace entirely with the polygon spec table
  (k=1 unchanged; k=2 chain; k≥3 ring edges + stable links + the k=3 x3
  double-stable exception; hub portals for x1/x2 only — k=1 portals x1).
- [x] Remove now-dead helpers/constants (`geoHubSplit*`, `geoHubExtra*`,
  `extraInteriorAngleOffset`, `extraInteriorRadius` — Phase 4 replaces the
  position side; coordinate the two phases if implemented together).
  (`geoHubSplit*` + split-corner position branch removed in Phase 3;
  `geoHubExtra*`/`extraInterior*` position helpers deliberately left for
  Phase 4, which replaces the whole position side.)
- [x] `buildEdges` gap-chain logic: unchanged except corners are single-label;
  the "corners present → corners portal to hub, else stables" rule stays; the
  degradation ladder (shared stable / star) stays byte-for-byte.
- [x] Re-spec unit tests in
  [geometricHubTopology/](../test/unit/internal/services/template_generator/providers/topology/geometricHubTopology/):
  DELETE corner-split tests; UPDATE the 15-neutral/3P expectation (now 2
  interiors per hexagon, merged corners kept); ADD: k=2 chain edges + both
  portals; k=3 triangle ring + x3 4-connection exception + x3 has NO hub
  portal; k=4 square ring without diagonals + player-side chain; k=5 pentagon
  ring + nearest-vertex stable links; interiors-take-highest-plans slot
  mapping; no zone exceeds its allowed connection count; dangling-edge check
  still passes. Keep the star/shared-stable degradation tests untouched
  (behavior unchanged).

### Verification Plan
- `go build ./...`; `go test -count=20 ./test/unit/internal/services/template_generator/providers/topology/...`. ✅ RESULT: build clean, all 20 topology packages ok at -count=20.
- `go test ./test/... -count=1` green; coverage no-drop
  (`go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...`
  then `go tool cover '-func=coverage.txt'`; baseline 63.5%). ✅ RESULT: suite green, total 63.5% (no drop).
- Graph spot-check vs `3Hexes-2HexCentralZones-RoughExpectation.rmg.json`
  (3P, 15 neutrals): same adjacency structure (zone names may differ). ✅ RESULT:
  reference decoded to 3 gap chains player—stable—corner—stable—player + per-hexagon
  k=2 chain sL—x1—x2—sR + hub portals = 3 corners + 6 interiors — identical to the
  implementation's structure (covered by the new 15-neutral/3P test).

### Phase Summary
`distributeGeometricHubSlots` ladder is now stables(2/gap) → merged corners
(1/gap) → interiors (uncapped round-robin); `cornerSplits` round and the
`geoHubSplit*` constants + split-corner position branch are gone. `assignPlans`
corners take a single round; interiors are dealt highest-first round-robin and
`hexInteriors[i][m]` is DOCUMENTED as polygon vertex x(m+1) (x1 hub-left,
x2 hub-right, x3/x4 next left/right pair, odd k puts the last vertex on the
player axis). `buildInteriorEdges` rewritten: hub portals only for
`interiors[:min(k,2)]`; ring edges via new `interiorAngularOrder` helper
(odd vertices descend the left side reversed + even ascend the right → cyclic
neighbors; k=1 self-edge and k=2 duplicate wrap are absorbed by addDirect's
dedup/self guards); stable links via new `connectInteriorStables` (sL—x1,
sR—x2, and for k≥3 sL—x3, sR—x(min(4,k)) — which yields the k=3 both-stables
x3 exception naturally). `buildEdges`/degradation ladder untouched.
Tests: deleted the 2 corner-split tests; replaced the rule-8 chaining test with
10 new tests covering k=2/k=3/k=4/k=5 rings, portals, the ≤4-connection
invariant, highest-plans slot mapping, and the 3P/15-neutral hub-portal set.
`geoHubExtra*`/`extraInterior*` remain ONLY as interim interior-position
providers until Phase 4 replaces `computePositions`.

## Phase 4: Regular-hexagon geometry + per-player-count scale
Status: Complete

- [x] `computePositions`: implement the regular-hexagon spec (radii ratios
  2s/√3s/s from hub; ±30°/±60° ideal angles with sector-fraction fallback for
  P≥5; merged corner at gap mid-angle radius s).
- [x] Interior polygon vertex positions: regular k-gon centered at radius s
  along the player axis, x1—x2 edge facing hub, circumradius tuned so k=2
  spacing sL—x1—x2—sR is even and k≥3 shapes read cleanly (no overlap with
  stables/corners at typical preview scale).
- [x] Per-player-count player radius: 2–4P noticeably closer to hub (tune
  against flipped `One for All.png`), ≥5P ≈ current 0.46. Keep constants in
  the top-of-file const block, named `geoHub*`.
- [x] Unit tests: 3P positions form regular hexagons (equal side lengths
  within epsilon — hub→corner ≈ corner→stable ≈ stable→player); player radius
  smaller for 3P than for 6P; k-gon vertices equidistant from polygon center;
  k=2 even chain spacing; positions stay inside [0,1] for extreme configs
  (8P, many interiors).
- [x] Visual sanity: scratch `cmd/geohubpreview` tool (same technique as the
  previous session — delete after) generating PNGs for at least: 3P+9M,
  3P+12M, 3P+15M (vs RoughExpectation), 3P+16M, 2P+6M, 4P+12M (vs
  `4 Player Hub.png`), 6P (vs `6 Player Hub.png`). Compare with the reference
  images (remember the 180° flip for `One for All.png` only).

### Verification Plan
- Same commands as Phase 3 (build, -count=20 topology, full suite, coverage
  no-drop) + `golangci-lint-v2 run ./... --issues-exit-code=0` at the
  pre-existing baseline (84 gochecknoglobals, uncapped). ✅ RESULT: build clean,
  -count=20 green, full suite green, coverage 63.5% (no drop), lint uncapped =
  exactly 84 gochecknoglobals (nothing new).
- Generated preview PNGs visually match the reference expectations; record
  which configs were compared in the phase summary. ✅ RESULT: see summary.

### Phase Summary
`computePositions` rebuilt on regular-hexagon math: `playerRadius =
geoHubPlayerRadiusFor(P)` (0.38 for 2–4P via `geoHubPlayerRadiusClose` /
`geoHubClosePlayerMaximum=4`, 0.46 for ≥5P), `hexagonSide s = playerRadius/2`,
stables at `√3/2·playerRadius` with angular offset `min(π/6,
0.35·sector)` (`geoHubStableIdealOffset` + existing fraction — exact ±30° for
P≤4, sector-fraction fallback for P≥5), merged corners at gap mid-angle radius
s, single-stable gaps at mid-angle (degradation unchanged). New
`computeInteriorPositions`: k-gon centered at `circlePoint(θ, s)`; k=1 sits
exactly on the center (rule 7); k≥2 vertices at hub-relative angles
±(2⌈j/2⌉−1)·π/k (odd j on the sL side — matches `connectInteriorStables`),
circumradius `geoHubInteriorPolygonFactor=0.357·s` (exact even-k=2-chain
solution ≈0.3568). Removed `geoHubStableRadius/MergedCornerRadius/
InteriorRadius/ExtraInteriorRadius/ExtraAngleFraction` consts +
`extraInteriorAngleOffset`/`extraInteriorRadius` helpers. Added 5 position
tests (regular-hexagon side spread, 3P<6P player radius, k=3 vertex
equidistance, k=2 even chain, 8P+64-neutral [0,1] bounds) + `positionOf`/
`distanceBetween`/`spreadOf` helpers in common_test.go. Visual check via
scratch `cmd/geohubpreview` (DELETED after): 2P+6M, 3P+9M (regular hexagons ==
flipped One for All), 3P+12M (rule-7 rhombuses), 3P+15M (== RoughExpectation
structure), 3P+16M (clean k=3 triangle, no crossing edges — the old broken
case), 4P+12M, 6P+18M (petal shapes match the hand-made references).

## Phase 5: Final verification & handoff
Status: Complete (pending user smoke test)

- [x] Full suite: `go build ./...`; `go test ./test/... -count=1`;
  `go test -tags=integration_test ./test/integration/... ./test/performance/... -count=1`.
- [x] Coverage: no drop vs 63.5% baseline; new/changed functions covered per
  branch.
- [x] Lint: `golangci-lint-v2 run ./... --issues-exit-code=0` — nothing new
  beyond the 84-gochecknoglobals baseline.
- [x] Expectation: after Phase 1 the repo has ZERO failing unit tests — a
  clean `go test ./test/... -count=1` run.
- [x] Update repo memory ([/memories/repo/conventions.md] GeometricHub + zone
  tier entries: splits removed, polygon interior rule, regular-hexagon
  geometry, platinum arena composite, buttonPositionLogger tests fixed).
- [ ] USER: GUI smoke test (`go run .`) — walk the zone-count ladder for
  2/3/4/6 players, confirm hexagon shapes, interior polygons, closer 2–4P
  players, gold swords on the platinum arena sprite.

### Verification Plan
- All commands above green; user visually accepts the shapes. ✅ RESULT (agent
  side): build clean; unit suite zero failures; integration suite ok (2.6s);
  performance package compiles ok ("no tests to run" — benchmarks only);
  coverage 63.5% == baseline; uncapped lint = exactly 84 gochecknoglobals.
  Per-function coverage of geometricHubLayout.go: everything 100% except
  `fillRoundRobin` 90.9% (pre-existing safety break) and
  `connectInteriorStables` 90.9% (unreachable defensive stables==0 guard,
  registered in todo/test_observations.md). ⏳ USER smoke test outstanding.

### Phase Summary
All automated gates green. Repo memory conventions.md updated (GeometricHub
entry rewritten to the implemented state, fix-round entry switched from
PLANNED to IMPLEMENTED with the vertex-mapping/geometry/tool facts, stale
"pre-existing buttonPositionLogger failures" note replaced with "suite is
clean now"). Unreachable defensive branch documented in
todo/test_observations.md. The only open item is the USER GUI smoke test.

## Final Recap
All five phases delivered on top of the uncommitted Geometric Hub + tiers
feature (branch `AD/refactoring-07-13`):
1. buttonPositionLogger tests fixed (filter helper + new frame-marker test;
   logger untouched) — repo suite now fully green.
2. `neutral_highest_arena.png` regenerated with GOLD swords: platinumgen now
   composites the `gladiator_arena.png` master @0.95 onto the recolored
   platinum bubble (new bilinear NRGBA source-over compositor in the tool).
3. Growth ladder simplified (stables → merged corners → uncapped interiors),
   corner splits deleted, interior connections rebuilt on the regular-k-gon
   spec (ring edges only, x1/x2 hub portals, k=3 x3 both-stables exception),
   vertex order deterministic and documented (`hexInteriors[i][m]` = x(m+1)).
4. Positions rebuilt on regular-hexagon math (player 2s / stable √3s /
   corner s; ±30° ideal with sector-fraction fallback; k-gon interiors at
   center s with circumradius 0.357s; players closer to hub for 2–4P) —
   visually verified against all reference images via a scratch preview tool
   (deleted after use).
5. Full verification: build/unit/integration/performance green, coverage
   63.5% (no drop), lint at the 84-gochecknoglobals baseline, repo memory
   updated. Files touched: geometricHubLayout.go, tools/platinumgen/main.go,
   3 regenerated PNGs, logger + geo hub test files (incl. common_test.go),
   todo/test_observations.md, this plan.

## Deployment Plan
1. USER: GUI smoke test (`go run .`) — walk the zone-count ladder for 2/3/4/6
   players on the Geometric Hub topology; confirm hexagon shapes, interior
   polygons, closer 2–4P players, gold swords on the platinum arena sprite.
2. Review the diff (`git status` / `git diff`) — note the branch also carries
   the entire previous session's uncommitted feature (~60 files).
3. Commit — suggested message: "Fix Geometric Hub geometry (regular hexagons,
   interior polygons), gold arena swords on platinum sprite, logger test
   expectations" (or fold into the feature commit, owner's choice).
4. Push; CI validates on Linux (build + tests; coverage gate ≥60% and
   no-decrease both satisfied at 63.5%).
5. No data migration; read-only dirs untouched; `.gen.json`/`.rmg.json`
   schemas unchanged.
