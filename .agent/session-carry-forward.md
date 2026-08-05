# Session Carry-Forward — Batch 6 (DI PR)

## 1. Session goal

Work through `todo/review-opus5-08-04.md` batch by batch, asking go/no-go and
clarifying questions before each batch, documenting rejections in place, and
stopping for owner review after every batch. This document covers **Batch 6 —
the DI PR (§2.3 + §2.4)**.

## 2. Fixes applied

- **§2.3** — [mandatoryContentItemMapper.go](../internal/mappers/mandatoryContentItemMapper.go)
  no longer builds its own `ContentRuleService`; it takes a
  `content_rules.IContentRuleService`.
- **§2.3 (extra, required)** — [generatorConfigMapper.go](../internal/mappers/generatorConfigMapper.go)
  built its own `MandatoryContentItemMapper`, which would have dead-ended the
  injection one level up. `NewConfigMapper` now takes the mapper.
- **§2.4** — [topologyBase.go](../internal/services/template_generator/providers/topology/base/topologyBase.go)
  no longer builds a `ZoneLabelProvider` or a connection service; both are
  injected.
- **§2.4 (extra, owner-approved)** — `NewTournamentTopologyService` built its own
  four `IClusterService` instances. Same anti-pattern, not named in the review;
  now injected.

**Net effect on the object graph** (verify in
[wire_gen.go](../internal/composition/wire_gen.go)): the application now
allocates **one** `zoneLabelProvider`, **one** `topologyConnectionService` and
**one** `contentRuleService`. Previously a full graph carried **16** label
providers, **16** connection services and **2** rule services.

## 3. Features added / changed

- **`IZoneLabelProvider`** — new interface in
  [zoneLabelProviderInterface.go](../internal/services/zones/zoneLabelProviderInterface.go)
  (8 methods), per AGENTS.md §4.2.1.
- **`IContentRuleService`** — new interface in
  [contentRuleServiceInterface.go](../internal/services/content_rules/contentRuleServiceInterface.go)
  (7 methods). Both are bound with `wire.Bind` in `providerSets.go`.
- **`TopologyConnectionService` exported.** It was unexported with unexported
  methods; wire generates into `internal/composition` and cannot name an
  unexported type. Its four externally consumed methods are now exported;
  `createBridgeConnection` and `buildZoneAdjacency` stay private, and
  `GetBorderGuardValue` moved above them to satisfy the `funcorder` linter.
- **Constructor signatures.** 17 topology constructors gained
  `zoneLabelProvider zone_services.IZoneLabelProvider` and
  `connectionService *base.TopologyConnectionService`.
  `NewTournamentTopologyService` gained four more cluster-service parameters.
- **`contentRuleHandler`** widened from `*ContentRuleService` to the interface.

## 4. File modifications

**New (5 + 1 folder):**

| File | Purpose |
| --- | --- |
| `internal/services/zones/zoneLabelProviderInterface.go` | `IZoneLabelProvider` |
| `internal/services/content_rules/contentRuleServiceInterface.go` | `IContentRuleService` |
| `test/test_helpers/configMapper.go` | builds the wired `GeneratorConfigMapper` |
| `test/test_helpers/contentRuleServiceMock.go` | reusable testify mock |
| `test/test_helpers/tournamentTopologyDependencies.go` | 8-value spread helper |
| `test/unit/.../base/topologyConnectionService/` | constructor test folder |

**Modified — production (26):** `providerSets.go`, `topologyServiceProvider.go`,
`wire_gen.go` (regenerated), `contentRuleHandler.go`, both mappers,
`templateGenerator.go`, `topologyBase.go`, `topologyConnectionService.go`, and
the 17 topology / cluster / positioned-builder constructors.

**Modified — tests (14):** the three `test_helpers` files, one integration test,
three `guiHandler` tests, four mapper tests, two tournament-topology tests and
`newTopologyBase_test.go`.

## 5. Tests added or updated

- `TestWhenRowsAreMapped_UsesTheInjectedContentRuleService` — proves the injected
  service is the one consulted (uses the new mock). This is the test the review
  asked for and it was **impossible before** the interface existed.
- `TestWhenBaseIsConstructed_RetainsTheInjectedZoneLabelProvider` — `assert.Same`
  on the injected instance.
- `TestWhenServiceIsCreated_ReturnsInstance` — constructor test for the newly
  exported `NewTopologyConnectionService`.
- ~30 existing call sites updated mechanically for the new signatures.

**Verification — all green:**

| Check | Result |
| --- | --- |
| `go build ./...` | pass |
| `go vet -tags=integration_test ./...` | pass |
| `go test -count=1 ./test/unit/...` | `unit=0` |
| `go test -tags=integration_test ./test/integration/...` | `integration=0` |
| `go test -tags='integration_test,gui' ./test/integration/gui/...` | `gui=0` |
| Unit coverage | **65.0%** — identical to the pre-batch baseline |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | **42 issues** (40 `gochecknoglobals`, 2 `dupl`) — unchanged baseline |

## 6. Git status snapshot

Branch `AD/refactoring-07-21`. **Nothing staged by the agent.** 42 modified
files, 5 new files and 1 new test folder — see §4. Batch 5's files no longer
appear in `git status`, confirming the owner committed them.

## 7. Rejections / things not done

- **Injecting `topologyConnectionService` without exporting it** — impossible.
  Wire generates into `internal/composition` and cannot name an unexported type.
  The owner chose to export rather than keep it internal.
- **Per-method test folders for the four newly exported connection-service
  methods** — not added. They keep their existing coverage through the
  `topologyBase` suites, which remain the intended entry point; adding parallel
  folders would duplicate those tests without adding coverage. Only the new
  constructor got a dedicated folder.
- **Concrete-type injection** — considered and rejected by the owner; it fixes
  the DI graph but does not deliver the review's stated stubbing benefit.

## 8. Open questions

None for Batch 6.

Still blocked from earlier batches: **§7.1** (are direct pushes to master
intentional?), **§9.1** (public-API documentation decision), **§2.7**
(finish or remove the gladiator-arena preview), **§1.8** (output-dir
persistence: `.gen.json` vs machine-local), **§2.2** (scope of extracting
regeneration policy from `app/gui/drivers/`).

## 9. Next recommended actions

Batch 7 — **Test-policy PR**: §6.3 (three `internal/services/content_rules`
tests import `app/gui/constants`) **then** §6.5 (depguard scope + CI tag check).
Order matters or CI goes red.

Then batches 8–13 per review §12: CI/security posture, docs, duplication
cleanup, coverage, product decisions, large refactors.

## 10. Carry-forward prompt

> Read `AGENTS.md` first. Hard rules, one line each: never modify `data/`,
> `internal/entities/template/` or `internal/registry/`; keep everything
> cross-platform (Windows + Linux, `path/filepath`, PowerShell chains with `;`);
> every change ships with tests and must not drop coverage; durable multi-session
> work gets a plan file under `plans/`; **never stage and never commit** — the
> owner reviews and commits.
>
> We are remediating the 46-finding review in `todo/review-opus5-08-04.md`,
> which defines 13 PR-sized batches in §12. Batches 1–6 are done: Security,
> Correctness, Durability, Input-validation, Performance and the DI PR. Findings
> are marked `✅ FIXED` / `❌ WILL NOT FIX` **in place** in the review document —
> do not create a separate plan file for this.
>
> Workflow for every batch, without exception: (1) ask the owner whether the
> batch should be done at all; (2) if declined, document in the review file why
> it should not be attempted in future; (3) ask all clarifying questions up
> front; (4) implement; (5) rewrite `.agent/session-carry-forward.md`; (6) stop
> and wait for owner review.
>
> Next up is Batch 7, the Test-policy PR: §6.3 before §6.5, or CI goes red.
>
> Useful gotchas: `wire gen` writes success to stderr, so PowerShell shows a
> `NativeCommandError` even when it worked. Never pipe `go test` through
> `Select-Object -First N` — it kills the upstream process and fakes an exit
> code 1; redirect to a temp file and use `Select-String`. Adding the first test
> that imports a previously untested package can *lower* total `-coverpkg`
> coverage by enlarging the denominator, and CI hard-fails on any decrease.
>
> See `.agent/session-carry-forward.md` for the full handoff.
