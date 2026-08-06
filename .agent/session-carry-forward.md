# Session Carry-Forward — Batch 10 (Duplication cleanup)

## 1. Session goal

Deliver **Batch 10** of the `todo/review-opus5-08-04.md` remediation: review
findings §3.1, §3.2, §3.3, §3.4 and §5.3 (duplicate code + positional
constructor calls).

## 2. Fixes applied

- **§3.1** — all fifteen hard-coded `WithGuardWeeklyIncrement(0.15)` call sites
  across twelve files now read
  `WithGuardWeeklyIncrement(common_connections.GetGuardWeeklyIncrements().Standard)`.
  The distinct `0.20` in [internal/services/zones/zoneFactory.go](internal/services/zones/zoneFactory.go)
  was deliberately left alone.
- **§3.3** — the two verbatim spell-label helpers (`bannedSpellLabel` in
  [app/gui/panels/bonusesPanel.go](app/gui/panels/bonusesPanel.go) and
  `spellNameAndSchool` in [app/gui/dialogs/bonusPickerDialog.go](app/gui/dialogs/bonusPickerDialog.go))
  were deleted in favour of one exported
  `constants.SpellNameAndSchool` in the new
  [app/gui/constants/spellLabel.go](app/gui/constants/spellLabel.go).
- **§3.4** — verified only. The owner had already extracted
  `newBaseButtonWidget` + `newButtonInset` in
  [app/gui/widgets/buttonWidget.go](app/gui/widgets/buttonWidget.go) in commit
  `0311318`. The headless GUI snapshot suite passes with zero diff and the two
  `dupl` lint findings are gone.
- **§3.2 + §5.3** — see below.

## 3. Features added / changed

- **`...CreationRequest` struct family** (§5.3). `ZoneFactory.CreateSpawnZone`
  and `TopologyBase.CreateNeutralZone` no longer take 7–9 positional arguments;
  they take a request struct. For consistency the three pre-existing creation
  structs were renamed into the same family (files renamed per AGENTS.md §4.1):
  `NeutralZoneCreation` → `NeutralZoneCreationRequest`, `HubZoneCreation` →
  `HubZoneCreationRequest`, `NeutralLikeZoneCreation` →
  `NeutralLikeZoneCreationRequest`.
- **`TopologyBase.CreateClusterZone`** (§3.2). One exported helper on
  [topologyBase.go](internal/services/template_generator/providers/topology/base/topologyBase.go)
  owns the spawn/neutral choice, the `Player%d` naming and the `FirstOrDefault`
  plan lookup. Per owner decision it was applied to **all ten** duplicated call
  sites, not only the four tournament cluster services.
- **Behaviour note for the owner (§3.3).** Both original spell-label copies
  ended with `if label == "" { label = spell.School }`. That branch is provably
  a no-op — `GetSpellSchoolDisplayName` returns one of five non-empty constants
  or the raw `schoolType`, so `label` is empty only when `spell.School` is
  empty, in which case the assignment changes nothing. It was **dropped**
  rather than carried into the new file as an uncoverable branch.

## 4. File modifications

**Created**

| File | Summary |
| --- | --- |
| [app/gui/constants/spellLabel.go](app/gui/constants/spellLabel.go) | `SpellNameAndSchool(sid) (name, school)` |
| [internal/models/spawnZoneCreationRequest.go](internal/models/spawnZoneCreationRequest.go) | Request struct for `ZoneFactory.CreateSpawnZone` |
| [internal/models/topologyNeutralZoneCreationRequest.go](internal/models/topologyNeutralZoneCreationRequest.go) | Request struct for `TopologyBase.CreateNeutralZone` |
| internal/models/neutralZoneCreationRequest.go | Renamed from `neutralZoneCreation.go` |
| internal/models/hubZoneCreationRequest.go | Renamed from `hubZoneCreation.go` |
| internal/models/neutralLikeZoneCreationRequest.go | Renamed from `neutralLikeZoneCreation.go` |

**Edited (production)**

- [internal/services/zones/zoneFactory.go](internal/services/zones/zoneFactory.go) — `CreateSpawnZone` now takes `models.SpawnZoneCreationRequest`; §3.1.
- [internal/services/zones/castleFactory.go](internal/services/zones/castleFactory.go), [internal/services/zones/zoneFactoryNeutralLike.go](internal/services/zones/zoneFactoryNeutralLike.go) — §3.1 / rename fallout.
- [.../topology/base/topologyBase.go](internal/services/template_generator/providers/topology/base/topologyBase.go) — new `CreateClusterZone`; both wrappers take request structs.
- [.../topology/base/topologyConnectionService.go](internal/services/template_generator/providers/topology/base/topologyConnectionService.go) — §3.1 (3 sites).
- `.../topology/{chain,geometricHub,hub,ring,web}Topology.go` and `.../topology/positionedTopologyBuilder.go` — converted to `CreateClusterZone`; unused `linq` imports dropped from chain/positioned/ring.
- `.../topology/tournament_variant/{balanced,chain,hub,ring}ClusterService.go` — converted to `CreateClusterZone`.
- [internal/services/connection_editor/zoneEditorService.go](internal/services/connection_editor/zoneEditorService.go) — rename fallout.
- [app/gui/panels/bonusesPanel.go](app/gui/panels/bonusesPanel.go), [app/gui/dialogs/bonusPickerDialog.go](app/gui/dialogs/bonusPickerDialog.go) — §3.3.

**Edited (docs)**

- [todo/review-opus5-08-04.md](todo/review-opus5-08-04.md) — §3.1, §3.2, §3.3, §3.4, §5.3 marked `✅ FIXED` in place; §12 item 10 marked done.

## 5. Tests added or updated

**Added**

- [test/unit/app/gui/constants/spellLabel/spellNameAndSchool_test.go](test/unit/app/gui/constants/spellLabel/spellNameAndSchool_test.go) — 4 tests (known SID name, known SID school, unknown SID sentence-cased name, unknown SID generic `"Spell"`).
- [test/unit/.../topologyBase/createClusterZone_test.go](test/unit/internal/services/template_generator/providers/topology/base/topologyBase/createClusterZone_test.go) — 4 tests (spawn player index, spawn name, neutral plan lookup by label, hold-city win condition).

**Updated**

- `.../topologyBase/common_test.go` — added `newSpawnRequest` / `newNeutralRequest` builders.
- `.../topologyBase/createSpawnZone_test.go` (7 calls), `.../topologyBase/createNeutralZone_test.go` (8 calls), `test/unit/internal/services/zones/zoneFactory/{createSpawnZone,createNeutralZone,createHubZone}_test.go` — migrated to the request structs.

**Last verification run — all green:**

| Check | Result |
| --- | --- |
| `go build ./...` | pass |
| `go vet -tags=integration_test ./...` | pass |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test -count=1 ./test/unit/...` | pass |
| `go test -tags=integration_test -count=1 ./test/integration/...` | pass |
| `go test -tags='integration_test,gui' -count=1 ./test/integration/gui/...` | pass (zero snapshot diff — this is the §3.4 verification) |
| Unit coverage `-coverpkg=./internal/...,./app/...` | **65.6%** (baseline 65.3% — improved) |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | **0 issues** (baseline 42; the 2 `dupl` findings are gone) |

## 6. Git status snapshot

Branch: **`AD/refactoring-07-21`**. Last commits: `0311318 "Task 10 Resolved"`
(owner), `b6a3cd1 "Batch 9 Done"`.

**Nothing is staged.** `git status --short`:

```
 M app/gui/dialogs/bonusPickerDialog.go
 M app/gui/panels/bonusesPanel.go
 D internal/models/hubZoneCreation.go
 D internal/models/neutralLikeZoneCreation.go
 D internal/models/neutralZoneCreation.go
 M internal/services/connection_editor/zoneEditorService.go
 M internal/services/template_generator/providers/topology/base/topologyBase.go
 M internal/services/template_generator/providers/topology/base/topologyConnectionService.go
 M internal/services/template_generator/providers/topology/chainTopology.go
 M internal/services/template_generator/providers/topology/geometricHubTopology.go
 M internal/services/template_generator/providers/topology/hubTopology.go
 M internal/services/template_generator/providers/topology/positionedTopologyBuilder.go
 M internal/services/template_generator/providers/topology/ringTopology.go
 M internal/services/template_generator/providers/topology/tournament_variant/balancedClusterService.go
 M internal/services/template_generator/providers/topology/tournament_variant/chainClusterService.go
 M internal/services/template_generator/providers/topology/tournament_variant/hubClusterService.go
 M internal/services/template_generator/providers/topology/tournament_variant/ringClusterService.go
 M internal/services/template_generator/providers/topology/webTopology.go
 M internal/services/zones/castleFactory.go
 M internal/services/zones/zoneFactory.go
 M internal/services/zones/zoneFactoryNeutralLike.go
 M test/unit/.../topologyBase/{common,createNeutralZone,createSpawnZone}_test.go
 M test/unit/internal/services/zones/zoneFactory/{createHubZone,createNeutralZone,createSpawnZone}_test.go
 M todo/review-opus5-08-04.md
?? app/gui/constants/spellLabel.go
?? internal/models/{hubZone,neutralLikeZone,neutralZone,spawnZone,topologyNeutralZone}CreationRequest.go
?? test/unit/app/gui/constants/spellLabel/
?? test/unit/.../topologyBase/createClusterZone_test.go
```

The three ` D` + `??` pairs under `internal/models/` are the file renames (done
with `git mv`, then `git restore --staged internal/models/` so nothing stayed
staged).

## 7. Rejections / things not done

- `WithGuardWeeklyIncrement(0.20)` in `zoneFactory.go` — out of §3.1's scope, a
  genuinely different value.
- `NewSegmentButtonWidget` was **not** folded into `newBaseButtonWidget`; it
  uses `layout.UniformInset(constants.DefaultPaddingSmall)`, so sharing the
  body would move pixels.
- `bannedItemLabel` in `bonusesPanel.go` was left untouched (not part of §3.3).
- The `if label == "" { label = spell.School }` no-op branch was dropped, not
  ported (see §3 above).

## 8. Open questions

- None blocking Batch 10. Still-open owner decisions for later batches:
  §1.1 (transactionality), §1.5 (ceilings), §1.8 (output-directory persistence
  shape), §2.2 (regeneration-policy refactor scope). Optional leftover: §5.4
  (`.gitignore` blanket-ignores `/*.txt`).
- FYI: the lint baseline dropped from 42 issues to **0** during this batch. The
  40 `gochecknoglobals` findings are no longer reported — likely a `.golangci`
  configuration change in the owner's `0311318`. Worth a glance.

## 9. Next recommended actions

1. Owner reviews and commits Batch 10.
2. **Batch 11 — Coverage PR.** §6.2 (`internal/handlers` mirrored tests,
   starting with `stateHandler` and `previewHandler`), §6.4 (the two 0%
   catalogues `bannableItems.go` / `valueOverrideSids.go`).
3. **Batch 12 — Product decisions.** Only ⚠ §1.8 remains.
4. **Batch 13 — Large refactors.** §2.1 (extract filesystem policy) → unblocks
   §2.5; then §2.2; §2.6 opportunistically. Needs a plan file under `plans/`
   per AGENTS.md §4.7.

## 10. Carry-forward prompt

> Read `AGENTS.md` first. Hard rules, one line each: never modify `data/`,
> `internal/entities/template/` or `internal/registry/`; keep everything
> cross-platform (Windows + Linux, `path/filepath`, PowerShell chains with `;`);
> every change ships with tests and must not drop coverage; durable
> multi-session work gets a plan file under `plans/`; **never stage and never
> commit** — the owner reviews and commits.
>
> We are remediating the 46-finding review in `todo/review-opus5-08-04.md`,
> which defines 13 PR-sized batches in §12. Findings are marked `✅ FIXED` /
> `❌ WILL NOT FIX` **in place** in the review document — do not create a
> separate plan file for this.
>
> Workflow for every batch, without exception: (1) ask the owner whether the
> batch should be done at all; (2) if declined, document in the review file why
> it should not be attempted in future; (3) ask all clarifying questions up
> front; (4) implement; (5) rewrite `.agent/session-carry-forward.md`; (6) stop
> and wait for owner review.
>
> Batches 1–10 are done. Batch 10 (duplication cleanup: §3.1, §3.2, §3.3, §3.4,
> §5.3) is complete and awaiting owner review — nothing is staged. Next up is
> Batch 11 (coverage: §6.2, §6.4). Full handoff detail is in
> `./.agent/session-carry-forward.md`.
