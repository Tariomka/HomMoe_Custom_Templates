# Session carry-forward — 2026-08-12

## 1. Session goal

Work `todo/backlog-opus5.md` in the execution order of its §8: batch **A**
(§1.4), batch **B** (§1.2 → §1.3) and batch **C** (§3.1, §3.2, §3.4), one
PR-sized batch at a time, each reviewed and committed by the owner before the
next started.

## 2. Fixes applied

- **§1.4 — a fatal window error was written to a discard handler.**
  [app/gui/program.go](../app/gui/program.go) installs
  `slog.New(slog.DiscardHandler)` before the event loop unless `-with-logging`
  is passed, so the `app.DestroyEvent` failure path exited 1 with no output.
  The branch now writes to `os.Stderr` before the (still discarded) `slog.Error`.
  The optional "return the error to `main.go`" refactor was **not** done.
- **§1.2 — hub-touching connections were guarded as player borders.** Every hub
  topology substituted a *player* label as the guard anchor because the hub is
  deliberately absent from `neutral_zone.Plans`, so hub edges got the
  player-border value (30 000) instead of the top tier (35 000).
- **§1.3 — random portals ignored endpoint tiers**, using a flat
  `ScaleByBorderGuardStrength(25000)` regardless of what they connected.

## 3. Features added / changed

### Batch A (commit `ee812af`)

- stderr on the fatal window path; `app/gui/program.go` recorded as an
  intentional unit-test gap in
  [test_observations.md](../todo/test_observations.md).

### Batch B (commit `472aaa1`) — behaviour change, owner-approved 2026-08-11

- `GetBorderGuardValue` in
  [topologyConnectionService.go](../internal/services/template_generator/providers/topology/base/topologyConnectionService.go)
  collapsed from three branches to
  `max(labelQuality(a), labelQuality(b)).GetGuardValue()`, scaled once. The
  collapse is exact because `neutral_zone.QualityUnknown` is `-1` and therefore
  loses the `max` against every real tier, while still yielding the
  player-border value when both endpoints are players.
- A private `labelQuality` ranks a hub label as `QualityHighest`. The three hub
  callers ([hubTopology.go](../internal/services/template_generator/providers/topology/hubTopology.go),
  [geometricHubTopology.go](../internal/services/template_generator/providers/topology/geometricHubTopology.go),
  [hubClusterService.go](../internal/services/template_generator/providers/topology/tournament_variant/hubClusterService.go))
  now pass the hub instead of a player anchor → **every hub edge is 35 000**.
  The hub was *not* added to `neutral_zone.Plans`.
- `CreateRandomPortalConnections` gained a trailing
  `neutralZones neutral_zone.Plans` parameter (interface, `TopologyBase`
  pass-through, seven call sites) and now calls `GetBorderGuardValue` →
  **portals are guarded by the higher endpoint tier**, 10 000–35 000 instead of
  a flat 25 000. That deleted the literal, which closed **§3.3** as well.
- The golden generator tests needed **no** expectation churn — they do not pin
  hub or portal guard values. The backlog's warning to expect churn was wrong.

### Batch C (uncommitted at time of writing) — no behaviour change

- **Naming convention set by the owner:** a name **builder** is
  `Get<X>For(label)`. `Get` separates it from the constant of the same stem;
  `For` says it derives a new name rather than returning an existing one.
- [internal/common/constants/contentNames.go](../internal/common/constants/contentNames.go)
  — `HubContentName`, `NeutralContentPrefix`, `SideContentPrefix`,
  `GetNeutralContentNameFor`, `GetSideContentNameFor`. No `"mandatory_content_*"`
  literal is left in production code; test-side literals were left as literals
  on purpose so they still catch a value change.
- [internal/common/constants/zoneNames.go](../internal/common/constants/zoneNames.go)
  — `GetHubZoneNameFor`, `GetPlayerZoneNameFor`, `GetNeutralZoneNameFor`. The
  prefix constants stay exported. Every production `XZonePrefix + label` now
  goes through a builder.
- **`hubLabel.go` / `constants.IsHubLabel` deleted** (a batch-B artefact). The
  hub check uses the pre-existing
  [zone_helpers.IsZoneNameHub](../internal/helpers/zone_helpers/zoneNameType.go),
  and [zoneClassifier.go](../internal/services/zones/zoneClassifier.go) uses
  `IsZoneNamePlayer` instead of a raw `strings.HasPrefix`.
- [internal/common/common_distances/footholdDistancePresets.go](../internal/common/common_distances/footholdDistancePresets.go)
  — `GetFootholdDistancePresets()` returns a named struct (`Crossroads`,
  `NearMainCastle`, `NearSecondCastle`). Bounds unchanged; deliberately **not**
  added to the user-facing `GetContentDistancePresets` catalogue. The new `Name`
  fields cannot reach the output because `WithDistance` copies only `Min`/`Max`.
- **Owner extension, added during review:**
  [internal/common/constants/connectionNames.go](../internal/common/constants/connectionNames.go)
  — 19 *unexported* connection prefixes plus `Get*ConnectionNameFor` builders
  (pseudo, bridge, fallback, neutral ring, the five tournament families, web,
  geometric hub, portal, portal-hub, manual, chain, random, ring), wired through
  the chain/ring/web/geometric/positioned/tournament topologies. This is the
  same idea as the zone/content names applied to connection names, and it is why
  `chainTopology.go`, `ringTopology.go`, `positionedTopologyBuilder.go` and
  `positionedTopologyZoneDecorator.go` appear in the diff.

## 4. File modifications

### Committed (batches A and B)

| File | Change |
| --- | --- |
| [app/gui/program.go](../app/gui/program.go) | stderr on the fatal `DestroyEvent`; `fmt` import. |
| [todo/test_observations.md](../todo/test_observations.md) | `program.go` added to the Gio-UI gap list. |
| `.../topology/base/topologyConnectionService.go` | `GetBorderGuardValue` collapsed to `max`; new `labelQuality`; portal guard from the tier. |
| `.../topology/base/topologyConnectionServiceInterface.go`, `.../base/topologyBase.go` | `CreateRandomPortalConnections` signature + pass-through. |
| `.../topology/{chain,ring,hub,geometricHub,web,tournament}Topology.go`, `.../positionedTopologyBuilder.go` | portal call sites pass `neutralZones`. |
| `.../topology/{hub,geometricHub}Topology.go`, `.../tournament_variant/hubClusterService.go` | hub guard anchor deleted; hub label passed instead. |
| 5 unit-test files under `test/unit/.../topology/` | hub-tier and portal-tier guard tests. |

### Uncommitted (batch C — staged by the owner, not by the agent)

| File | Change |
| --- | --- |
| `internal/common/constants/contentNames.go` | **new** — mandatory-content name constants + builders. |
| `internal/common/constants/connectionNames.go` | **new** (owner) — connection-name prefixes + builders. |
| `internal/common/common_distances/footholdDistancePresets.go` | **new** — foothold placement bounds. |
| `internal/common/constants/zoneNames.go` | zone-name builders added. |
| `internal/common/constants/hubLabel.go` | **deleted** — superseded by `zone_helpers.IsZoneNameHub`. |
| `mandatoryContentProvider.go` | 6 name literals + 3 distance literals routed through the new helpers. |
| `topologyBase.go`, `topologyConnectionService.go`, `hubTopology.go`, `geometricHubTopology.go`, `webTopology.go`, `chainTopology.go`, `ringTopology.go`, `positionedTopologyBuilder.go`, `positionedTopologyZoneDecorator.go` | zone-/connection-name builders. |
| `tournament_variant/{balanced,chain,hub,ring}ClusterService.go` | zone-/connection-name builders. |
| `zones/zoneFactory.go`, `zones/zoneLabelProvider.go`, `zones/zoneClassifier.go`, `connection_editor/zoneEditorService.go` | builders; `zoneClassifier` uses `IsZoneNamePlayer` and dropped its `constants` import. |
| `test/unit/internal/common/constants/{zoneNames,contentNames}/` | **new** — 5 test files. |
| `test/unit/internal/common/constants/hubLabel/` | **deleted**. |
| `test/unit/internal/common/common_distances/footholdDistancePresets/` | **new** — 1 test file. |
| [todo/backlog-opus5.md](../todo/backlog-opus5.md) | §1.2/§1.3/§1.4/§3.1/§3.2/§3.3/§3.4 marked ✅ DONE with resolution notes; §0.1, §7, §8 and the coverage baselines updated. |

## 5. Tests added or updated

- **Batch B (6 files):** hub and tournament-hub cases in
  `getBorderGuardValue_test.go`; the portal suite threaded the new parameter and
  gained two guard-value tests; one variant-level test each in
  `hubTopology`, `geometricHubTopology` and `hubClusterService` pinning 35 000 on
  every hub-touching connection with `BorderGuardStrengthMultiplier = 1`.
- **Batch C (6 files, 1 deleted):** `zoneNames/` (3) and `contentNames/` (2)
  gofakeit-fuzzed builder tests that assert the literal prefix, so a constant
  edit still fails; `footholdDistancePresets/` pins the numeric bounds. The
  `hubLabel/` folder was removed with its implementation.

**Last full gate run (after batch C, including the owner's `connectionNames.go`):**

| Gate | Result |
| --- | --- |
| `go build ./...` | clean |
| `go vet -tags='integration_test,gui' ./...` | clean |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go test ./test/unit/... -count=1` | pass |
| `go test -tags=integration_test ./test/integration/... -count=1` | pass |
| `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | pass (GPU) |
| `golangci-lint-v2 run ./...` | **0 issues** |
| `wire diff ./internal/composition/...` | exit 0, no drift |
| Coverage | **72.6 %** (was 72.5 % through batch B; floor 72.5 %) |

## 6. Git status snapshot

Branch `AD/fixing_some_stuff_08-12`, in sync with
`origin/AD/fixing_some_stuff_08-12` at `472aaa1`.

```
472aaa1  Batch B done
ee812af  Batch A done
6eecc36  carry forward
```

The batch-C changes (24 files: 3 added, 2 deleted, 19 modified) are **staged by
the owner**, awaiting commit. `todo/backlog-opus5.md` and this file are
unstaged. The agent staged nothing (AGENTS.md §2.5); one accidental
`git rm --cached` during the `hubLabel.go` deletion was immediately reverted
with `git reset -q HEAD -- <path>` — **delete files with `Remove-Item`, never
`git rm`**.

## 7. Rejections / things the user declined

- Nothing was rejected outright. Two agent proposals were *redirected* by the
  owner during the batch-C scoping round:
  - the plan to move `IsHubLabel` into `zoneNames.go` → the owner had it
    **deleted** instead, in favour of the existing `zone_helpers.IsZoneName*`
    predicates;
  - `NeutralMandatoryContentName(label)` (the backlog's suggested name) →
    renamed to `GetNeutralContentNameFor(label)`; "Mandatory" adds nothing and
    the `Get…For` shape is now the house style for name builders.
- Standing decisions later sessions must not re-litigate: the output directory
  is **never** persisted (AGENTS.md §2.7); `todo/backlog.md` was deleted on
  purpose (§0.2 is its record); §2.2 Branch A, §2.4, §2.5 and §6.1 are
  owner-gated.

## 8. Open questions

- **[layoutRingHub.go](../internal/services/preview_service/layoutRingHub.go#L32-L43)
  was left with raw prefix checks.** It needs `HasPrefix(HubZonePrefix)` for
  per-player hubs *and* exact `== HubZoneName` for the shared hub;
  `zone_helpers.IsZoneNameHub` conflates the two and would change the preview
  layout. Converting it needs a precise pair of helpers — owner call.
- **Three `constants.HubZonePrefix + label` remain in
  [hubTopology.go](../internal/services/template_generator/providers/topology/hubTopology.go#L73).**
  They build *connection* names that happen to share the zone prefix. If the
  owner prefers, they could move to a `Get…ConnectionNameFor` builder in
  `connectionNames.go` for consistency.
- `internal/common/constants/connectionNames.go` has **no dedicated unit-test
  folder**; its builders are covered only indirectly through the topology
  tests. AGENTS.md §4.6 wants a mirrored folder per implementation file.
- §5.5 step 1 still needs the `gui-snapshot-failures` CI artifact; §5.4 (b)
  still wants one sentence of confirmation before batch L (detail in the
  2026-08-11 handoff, recoverable from git history).

## 9. Next recommended actions

1. Commit batch C (already staged).
2. Decide the three open naming/test questions in §8.
3. **Batch D — §1.1**, editor-state copies are shallow: `EditorStateDto` needs a
   `Clone`. Deep copy + regression tests, benchmark before/after. It is also the
   dependency for batch I (§2.1).
4. **Batch L** (§5.4 a–c, §5.5) before batch F, so the file-explorer tests are
   written against a settled GUI harness and a comparison that can actually fail.
5. Then E (§4.1 Save To) → F (§5.3) → G (§2.3, owner review of regenerated
   snapshots) → H (§5.1, §5.2) → I (§2.1, needs a `plans/` file) → J (§2.2 B).

## 10. Carry-forward prompt

> Read `AGENTS.md` first — it governs everything below.
>
> Hard rules, one line each: never modify `data/`, `internal/entities/template/`
> or `internal/registry/` without explicit owner approval; keep every change
> cross-platform (Windows + Linux, `path/filepath`, PowerShell chained with `;`
> never `&&`); every change ships with tests and must not drop unit coverage
> below 72.5 % (currently 72.6 %); durable multi-session work gets a plan file
> under `plans/`; never stage and never commit — the owner reviews and commits
> (delete files with `Remove-Item`, never `git rm`); never change where
> `.rmg.json` is written and never persist the output directory; never run a
> bulk in-place rewrite over the repository.
>
> Where work left off: `todo/backlog-opus5.md` is the authoritative backlog.
> Batches **A** (§1.4), **B** (§1.2 → §1.3, which also closed §3.3) and **C**
> (§3.1, §3.2, §3.4) are **done and marked ✅** — 7 of 19 items. A and B are
> committed (`ee812af`, `472aaa1`); C is staged and awaiting the owner's commit.
> Batch B changed generated guard values on purpose (hub edges 35 000, portals
> by endpoint tier) with owner approval. Batch C introduced the house naming
> convention for name builders: **`Get<X>For(label)`**, in
> `internal/common/constants/{zoneNames,contentNames,connectionNames}.go` and
> `common_distances/footholdDistancePresets.go`.
>
> Next up is **batch D (§1.1)** — a deep `Clone` for `EditorStateDto` — but read
> §8 of the backlog for the full order, and note that batch **L** should run
> before batch **F**.
>
> Before starting any batch, prompt the owner to confirm the item(s) and surface
> every open question first; they expect a scoping round before implementation.
>
> See `./.agent/session-carry-forward.md` for the full handoff, including the
> open questions and the gate results.
