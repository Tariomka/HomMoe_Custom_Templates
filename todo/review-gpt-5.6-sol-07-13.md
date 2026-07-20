# Project Review — 2026-07-13 (GPT-5.6 Sol)

Senior/principal-level review of the full repository, verified at commit
`3ac96111109ccf28c90b2386cbea5d4f840efb1c` on 2026-07-13.

Toolchain: `go1.26.5 windows/amd64`; root module Go 1.26.5; tools module Go
1.26.3; Gio v0.10.0; golangci-lint v2.12.2 (built with Go 1.26.3). Configured
lint baseline: **84 `gochecknoglobals` issues** (the review's original "50"
was golangci's default `max-issues-per-linter: 50` display cap, not the real
count — see §8). Unit coverage baseline: **62.2%**.

> **Post-review verification (2026-07-14, in-repo audit at the same commit):**
> every claim below was re-checked against the working tree. Items marked
> **❌ INVALIDATED** are factually wrong, were explicitly skipped/deferred by
> the owner before this review, or duplicate existing tracking files — do not
> action them. Invalidated: §1.3 (owner-deferred), §6.2/§6.3 (already tracked
> in test_observations.md), §7.2 (owner declined the pin), §8 (wrong data,
> corrected in place). Verified accurate against code: §1.1, §1.2, §3.1, §4.1,
> §6.4, §6.5 (with location caveat), §7.1, §7.3, §9.1, §9.2.

This review supersedes the deleted historical `todo/review-fable-07-12.md`
(final pre-deletion version recovered from Git), and consolidates the current
[todo/backlog.md](backlog.md), [todo/test_observations.md](test_observations.md),
repository memories, and owner decisions. It does not supersede those two live
tracking files; §0 dispositions them explicitly.

Severity legend: 🔴 High (bug / correctness / user-visible) · 🟠 Medium
(architecture, performance, CI gap) · 🟡 Low (readability, hygiene) ·
⚪ Informational / no action.

---

## 0. Disposition of prior reviews, observations, backlog, and memories

### 0.1 Fixed since the historical review ✅

Every numbered item in the historical review is covered below; combined rows
mean the historical item itself pointed to the same implementation.

| Prior item | Current disposition and evidence |
| --- | --- |
| §1.1 + §6.3 lost road mutations and regressions | **Fixed.** Current topology code mutates indexed zones and the mirrored `topologyBase` suites cover bridge/fallback roads; build, stress tests, and configured lint pass. |
| §1.2 exit path / confirmation reset | **Fixed.** [stateFiles.go](../app/gui/drivers/stateFiles.go#L53-L64) calls the injected close callback, while save and mutation paths reset confirmation; `stateExit_integration_test.go` covers the workflow. |
| §1.3 game mode not restored | **Fixed.** [generalPanel.go](../app/gui/panels/generalPanel.go#L115-L118) selects the saved mode instead of forcing index zero. |
| §1.4 duplicate bridge-name loop | **Fixed.** The topology rewrite guarantees progress and the relevant stress suites pass 20 consecutive runs. |
| §1.5 unknown victory coercion | **Fixed.** [victoryConditions.go](../app/gui/constants/victoryConditions.go#L45-L56) returns `(Victory, bool)` and [generalPanel.go](../app/gui/panels/generalPanel.go#L126-L134) reports the fallback. |
| §1.6 `UpdateTemplate` aliases `Variants` | **Fixed.** [guiHandler.go](../internal/handlers/guiHandler.go#L61-L64) clones the slice header before replacing variant data. |
| §1.7 Steam-path / glob error shape | **Fixed.** [io.go](../internal/helpers/io.go#L142-L157) wraps glob errors and returns `ErrTemplatesDirNotFound`; platform-specific registry lookup is present. |
| §1.8 file-explorer append/complexity | **Fixed.** [fileExplorerDialog.go](../app/gui/dialogs/fileExplorerDialog.go#L541-L560) uses `slices.Concat`; confirmation is split into focused methods. |
| §1.9 G115 alpha conversion | **Fixed/cleared.** Fresh configured lint has no `gosec` finding. |
| §2.1 validators package | **Fixed.** [editorStateValidator.go](../internal/validators/editorStateValidator.go#L28-L38) validates loaded state and carries executable fixes; handler load applies them. |
| §2.2 + §6.1 `drivers.State` seam/god file/tests | **Fixed.** State is split by concern, depends on `interfaces.ITemplateHandler`, and has mirrored driver unit tests plus integration coverage. |
| §2.3 zone-editor god file | **Fixed.** Methods are split across dialog/canvas/snap/connection/zone files. |
| §2.4 layout/preview oversized functions | **Fixed.** Layout panel is split by concern and preview canvas helpers are extracted. |
| §2.5 + §3.1 topology duplication | **Fixed.** Shared fixed-layout construction is centralized; fresh `dupl` output is empty. |
| §2.6 mapper provider allocation | **Fixed.** `GeneratorConfigMapper` owns and reuses its content provider. |
| §3.2 row factories | **Fixed.** Reusable slider and checkbox row widgets are used throughout panels. |
| §3.3 geometry-center parameters | **Fixed.** Call sites use the shared `layoutCenter` constant; fresh `unparam` output is empty. |
| §4.1 per-frame reflection | **Fixed.** [editorStateDto.go](../internal/dtos/editorStateDto.go#L176-L186) uses a hand-written comparison, with parity/tripwire tests. |
| §4.2 static tab allocations | **Fixed.** [window.go](../app/gui/editor/window.go#L69-L78) caches static `FlexChild` wrappers. |
| §4.3 zone-editor geometry every frame | **Fixed.** Geometry is gated by dirty state/canvas size and mutation paths mark it dirty. |
| §4.4 preview pixel fill | **Fixed.** Connector brushes use bulk drawing; fresh performance and unit suites pass. |
| §5.1 oversized functions | **Fixed.** Fresh configured lint has no `funlen` or `gocognit` output. |
| §5.2 magic serialized strings | **Fixed.** Fresh configured lint has no `goconst` output. |
| §5.5 TODO purge | **Fixed as scoped.** Fresh configured lint has no `godox`; permanent future work moved to backlog/owner records. |
| §5.6 Steam-path comment | **No longer actionable.** The retained inline note is accurate and has no runtime effect. |
| §6.2 helper IO testability | **Fixed.** Dedicated mirrored tests cover VDF traversal, path construction, and glob outcomes; current file coverage is 16.7% because platform/host branches remain. |
| §6.4 coverage posture | **Accepted and re-baselined.** CI still enforces a 60% floor; the current measured 62.2% file-level posture is detailed in §6.1. |
| §7.1 vulnerability scan | **Fixed.** [pr-validation.yml](../.github/workflows/pr-validation.yml#L55-L75) runs `govulncheck` on PRs. |
| §7.2 race detector | **Fixed.** PR validation has a dedicated Linux `go test -race ./test/unit/...` job. |
| §7.3 full build / gated compile | **Fixed.** CI runs `go build ./...` and `go vet -tags=integration_test ./...`. |
| §7.4 dual-OS validation | **Fixed.** A Windows build/unit job is present. |
| §7.5 release hardening except action pin | **Mostly fixed.** `-trimpath`, version injection, tag checkout, checksums, concurrency, and artifact validation are present; the remaining mutable release action is carried to §7.2. |
| §7.6 tracked output artifacts | **Fixed.** No ignored file is tracked; no binary/generated output artifact is tracked. |
| §7.7 Dependabot | **Fixed.** [.github/dependabot.yml](../.github/dependabot.yml#L1-L20) updates root Go and Actions dependencies weekly. |

### 0.2 Invalidated / accepted / owner-controlled ✖

| Prior or memory item | Disposition |
| --- | --- |
| Historical §3.4 duplicated registry lookup globals | **Owner's responsibility.** Do not convert or edit automatically; current lint accounting is documented in §8. |
| Historical §5.3 exported functions returning private types | **Owner's responsibility.** The owner retained this API style and already renamed `selectedIDs`; no new finding. |
| Historical §5.4 default `slog` logger | **Accepted configuration.** The owner disabled the relevant no-global behavior; fresh configured lint is clean for `sloglint`. |
| Backlog: preview sub-pixel `Vec2` | **Deliberate future work.** No current correctness failure; retain in [backlog.md](backlog.md#L5-L9). |
| Backlog: dead Chain/Ring `createTopologyAdjacency` branches | **Explicit owner retention.** Removal was tried and rolled back; do not re-report or remove. |
| Backlog: replace schema `[2]float64` with `Vec2` | **Protected-owner decision.** It touches `internal/entities/template`; retain in backlog only. |
| Test observations: Gio widget/layout files | **Accepted integration territory.** Pure formatters and business services are unit-tested; UI rendering is not a unit-test requirement under AGENTS.md. |
| Test observations: State host/dialog-only branches | **Accepted where genuinely callback/Gio-bound.** Current seams and integration exports cover file load/save and manual regeneration without widening production API. |
| Memory claims that GUI has four tabs / Zone Content tab / footer panel | **Stale memory, not source truth.** Current [window.go](../app/gui/editor/window.go#L32-L38) has three tabs and zone content is opened from Layout & Zones dialogs. |
| Memory claims 86%–92% coverage, Go 1.26.3 root, Gio v0.9.0, old service files | **Stale.** Current measured coverage is 62.2%; root is Go 1.26.5/Gio v0.10.0; old service-root implementations are deleted. |
| Resource-density `/200` candidate | **Validated compatibility behavior, not a finding.** [generationTuning.go](../internal/models/generationTuning.go#L31-L36) matches the parallel C# reference implementation and current tests intentionally pin it. |
| `go mod tidy -diff` failure on this Windows checkout | **Environment/EOL artifact, not dependency drift.** Root shows 58 removed/58 added lines and tools 985/985; both removed/added multisets are identical. CI checks a fresh Linux checkout. |

### 0.3 Carried forward or newly re-verified ❗

| Source item | New section |
| --- | --- |
| Historical §2.7 internal tests import GUI constants | §3.1 |
| Historical §7.5 third-party release action pin | §7.2 |
| Test observation: file explorer has no interaction scenario | §6.2 |
| Test observation: manual editor has no interaction scenario | §6.3 |
| Memory: manual editor reset knowingly incomplete | §1.3 |
| Current persistence durability gap | §1.1 |
| Current integration-tag inconsistency | §1.2 |
| Current live-preview hot path | §4.1 |
| Current catalog coverage gap | §6.4 |
| Current direct-push CI coverage | §7.1 |
| Current docs/toolchain/UI drift | §9.1–§9.2 |

---

## 1. Bugs & correctness

### 1.1 🔴 Persistent user files are replaced non-atomically

**Evidence.** [fileService.go](../internal/services/file_service/fileService.go#L41-L49)
and [fileService.go](../internal/services/file_service/fileService.go#L64-L75) write JSON
directly to the final path:

```go
return os.WriteFile(filepath, data, fileReadWritePermission)
...
if err = os.WriteFile(out, data, fileReadWritePermission); err != nil {
    return "", err
}
```

PNG export likewise opens/truncates the final file before encoding at
[fileService.go](../internal/services/file_service/fileService.go#L91-L100).

**Why it is wrong.** `os.WriteFile` and `os.Create` truncate an existing file
before all bytes are durable. A crash, power loss, disk-full condition, or PNG
encode error can destroy the last valid `.gen.json`, `.rmg.json`, or preview and
leave a partial file. Settings are the user's authored work; template save is a
multi-file operation whose JSON can succeed while PNG fails.

**Fix.** Add one private same-directory atomic-write helper: `os.CreateTemp`,
write/encode, `Chmod(0o644)`, `Sync`, `Close`, then replace the destination.
Remove the temp file on every failure. Handle Windows replacement explicitly:
`os.Rename` cannot replace every existing destination on Windows, so use a
cross-platform replacement helper with rollback semantics rather than deleting
the only valid destination first. Consider saving template JSON and PNG to temps
before committing either, or document/report the partial-success contract.

**Tests.** Extend the existing mirrored files:
`test/unit/internal/services/file_service/fileService/saveSettings_test.go`,
`saveTemplate_test.go`, and `savePreviewImage_test.go`. Add tests proving an
existing destination remains byte-identical when an injected write/encode/commit
step fails and that no temp file remains. If fault injection would exist only
for tests, instead extract a small file-operations interface used by production
and test implementations; do not expose test-only APIs.

### 1.2 🟠 One integration test bypasses the repository's test gate

**Evidence.** [rmgTemplateModel_test.go](../test/integration/rmgTemplateModel_test.go#L1-L13)
starts directly with `package integration_test`; all four sibling integration
files and both performance files start with `//go:build integration_test`.
`go list` confirms the default package contains one external test, while the
tagged package contains five. AGENTS.md explicitly says every file in these
directories carries the tag.

**Why it is wrong.** `go test ./test/...` unexpectedly reads and decodes all 57
bundled templates. The default suite no longer excludes the gated directories,
so local/CI timing and the documented split are false. A future test-only
accessor used in this file would also fail in the default build.

**Fix.** Add `//go:build integration_test` plus the required blank line before
the package declaration. No production change is required.

**Tests.** Add a repository-policy test or CI script that enumerates
`test/integration/*_test.go` and `test/performance/*_test.go` and fails unless
each first non-blank line is the build constraint. Exact suggested location:
`test/unit/repository/testLayout/integrationBuildTags_test.go` (repository-level
policy exception to the production mirror), or implement the check directly in
CI if AGENTS.md should remain code-test-only.

### 1.3 🟡 ❌ INVALIDATED — “Reset to generated” resets only to the dialog's opening snapshot

> **❌ INVALIDATED (2026-07-14): explicitly skipped previously.** The mismatch
> is real but was knowingly deferred by the owner: the in-code comment reads
> "wont add todo so the llm does not trigger" and repo memory records it as a
> "known issue, deliberately NOT a TODO" (accepted during review item §4.3,
> commit 1645ee7). This is owner-initiated territory — agents must not action
> or re-report it unless the owner raises it. Original text kept for reference.

**Evidence.** The UI admits the mismatch in
[zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go#L245-L255):

```go
// Reset only resets current edits, not all manual edits, need to fix eventually...
NewButtonWidget(theme, "Reset to generated", ...)
```

The constructor snapshots the zones it receives as `originalZones`; Layout &
Zones passes `lastTemplate.Variants[0]`, which may already be the reapplied
manual snapshot. [resetToOriginal](../app/gui/dialogs/zoneEditorDialog.go#L413-L433)
therefore restores only the dialog-opening values.

**Why it is wrong.** After applying manual edits, closing, and reopening the
editor, the button says it will return to generated topology but cannot do so.
It returns to the already-edited snapshot, misleading the user and preventing
an in-dialog discard of all persisted manual changes.

**Fix.** Decide the intended operation and make label/data agree:

- If it only discards changes made since opening, rename it to “Reset current
  session” and replace the evasive comment with accurate documentation.
- If “Reset to generated” is intended, add a State-level operation that clears
  `ManualZones`/`ManualConnections`, regenerates from current settings, and
  supplies that fresh variant to the dialog. This changes authoritative manual
  state, so require an explicit confirmation before destructive reset.

**Tests.** Add a headless interaction scenario in
`test/integration/manualZoneReset_integration_test.go`: generate → apply manual
edit → reopen → reset → assert either opening-snapshot semantics (renamed UI) or
fresh generator zones and cleared manual DTO fields. **Owner decision required:**
choose the button's contract before implementation.

---

## 2. Architecture

No new production layering violation, dead package, or god-object regression was
verified. `go list` confirms current production packages compile; depguard has
no production issue. Persistence durability (§1.1) is the only missing service
abstraction with correctness impact.

Verified non-issues:

- `GUIHandler.UpdateTemplate` clones `Variants` before replacement.
- `drivers.State` uses an interface seam and concern-based method files.
- Legacy root `internal/services/previewLayout.go`, `previewRenderer.go`,
  `templateWriter.go`, and `settingsFileLoader.go` are absent.
- Protected schema/data/registry directories were not modified or treated as
  refactoring targets.

---

## 3. Duplicate code and dependency direction

### 3.1 🟡 Three internal-service tests import the GUI layer

**Evidence.** [getVariantsForContent_test.go](../test/unit/internal/services/content_rules/variantMappingManager/getVariantsForContent_test.go#L1-L9),
[getVariantForContentById_test.go](../test/unit/internal/services/content_rules/variantMappingManager/getVariantForContentById_test.go#L1-L9),
and [createRuleFromSavedRule_test.go](../test/unit/internal/services/content_rules/contentRuleManager/createRuleFromSavedRule_test.go#L1-L9)
import `app/gui/constants` to obtain content IDs. The production depguard rule
explicitly excludes tests in [.golangci.yml](../.golangci.yml#L135-L142).

**Why it is wrong.** Tests for an internal service are now coupled to a UI
catalog and can keep the wrong dependency direction alive even if production
code remains clean. Moving or changing GUI presentation data can break business
logic tests unrelated to the UI.

**Fix.** For these behavioral tests, use the actual wire SID literals in local
`const` declarations or construct `models.SidMapping` inputs directly. Avoid
editing protected `internal/registry`; do not add production exports solely for
tests. Add a second depguard scope for `test/unit/internal/**` denying `app/**`,
or remove the blanket `!$test` exception from the existing rule and add narrow
exceptions only where architectural integration is intentional.

**Tests.** No new behavior test is needed. Run the three existing mirrored test
folders and lint after the import replacement.

No current `dupl` finding exists; the former topology and widget-row duplicates
are resolved.

---

## 4. Performance

### 4.1 🟠 Live preview recomputes topology layout on every frame

**Evidence.** [previewPanel.go](../app/gui/panels/previewPanel.go#L148-L170) calls:

```go
previewLayout := this.layoutService.BuildPreviewLayout(
    template, this.state.GetStateData().Topology, float64(canvasSize.X))
```

inside the returned per-frame Gio widget. For Random/scatter layouts,
[layoutScatter.go](../internal/services/preview_service/layoutScatter.go#L122-L134)
can execute 500 O(zones² + edges×zones) relaxation iterations.

**Why it is wrong.** Idle repaints, pointer movement, status changes, and dialog
frames repeat deterministic layout work even though template, topology, and
canvas side did not change. This is reasoned (not profiler-measured), but the
hot path and upper bound are concrete; it can consume CPU and reduce UI
responsiveness on dense random layouts.

**Fix.** Cache `preview.Layout` in `PreviewPanel` using a key of template
identity/revision, topology, and canvas side. Invalidate when a new generated or
manually updated template is installed, topology changes, or the square side
changes. Prefer a monotonically increasing template revision exposed by State
over deep hashing the full template every frame. Keep the stateless
`PreviewLayoutService` API unchanged unless the same cache is needed elsewhere.

**Tests.** Add a call-count test around an injected narrow layout-builder
interface in `test/unit/app/gui/panels/previewPanel/previewLayoutCache_test.go`
only if the cache key can be exercised without a Gio shaper; otherwise add a
headless benchmark/assertion to
`test/performance/window_preview_idle_test.go` that renders repeated unchanged
frames and proves only one build, then changes topology/side/template and proves
one rebuild per key change. Record before/after benchmark timing and allocations.

---

## 5. Readability & maintainability

The configured lint run found no `funlen`, `gocognit`, `goconst`, `godox`,
`dupl`, or formatter issue (re-verified 2026-07-14 with issue caps disabled).
The deliberately hidden “wont add todo” comment is an explicit owner decision,
not a tracking lapse — see §1.3 (invalidated).

The remaining 50 linter findings are handled exhaustively in §8.

---

## 6. Testing

### 6.1 ⚪ Current measured coverage posture

The deduplicated unit profile contains 178 production files. **48 files are
below 80%; 35 are at 0%.** The 62.2% total remains above CI's 60% floor. Most
zeros are Gio rendering files intentionally assigned to integration coverage,
so a blanket percentage chase is not recommended.

Current files below 80%:

| Coverage | Files |
| --- | --- |
| 0% — Gio components/dialogs/editor/widgets/themes (accepted integration territory) | `app/gui/components/dropdownSelector.go`, `segmentButtonGroup.go`; dialogs `bonusPickerDialog.go`, `pickerDialog.go`, `ruleDialog.go`, `zoneContent.go`, `zoneContentDialog.go`, all five `zoneEditor*` files; `app/gui/drivers/tab.go`; `app/gui/themes/theme.go`; `app/gui/utils/draw.go`; 15 files under `app/gui/widgets/`. |
| 0% — actionable pure catalogs | `app/gui/constants/bannableItems.go`, `app/gui/constants/valueOverrideSids.go` (see §6.4). |
| 0% — platform/protected accepted | `internal/helpers/steamPath_windows.go`; protected `internal/registry/factionTypeValues.go`, `orientationModeValues.go`. |
| 2.9%–39.3% | `app/gui/drivers/dialogHost.go` 2.9%; protected `internal/registry/mapObjectArtifactValues.go` 5.9%; `internal/helpers/io.go` 16.7%; `app/gui/dialogs/fileExplorerDialog.go` 20.4%; `app/gui/drivers/stateFiles.go` 30.4%; protected `internal/registry/spellSchoolTypeValues.go` 33.3%; `app/gui/constants/spells.go` 39.3%. |
| 46.4%–77.3% | protected `internal/registry/contentIncludeListValues.go` 46.4%; `app/gui/drivers/state.go` 65.7%; protected `internal/registry/guardedContentPoolValues.go` 66.7%, `unguardedContentPoolValues.go` 66.7%; `app/gui/drivers/stateManualEdits.go` 68.4%; protected schema `internal/entities/template/template_rule/winConditions.go` 77.3%. |

### 6.2 🟠 ❌ INVALIDATED — File explorer's critical workflows have no interaction test

> **❌ INVALIDATED (2026-07-14): already explicitly tracked as deferred future
> work.** [test_observations.md](test_observations.md) — the intentional-gap
> registry mandated by AGENTS.md §4.6 — already records this exact gap
> ("synthetic-click coverage via the test/performance AppRunner pattern is
> possible future work"). The review's own §0.2 accepts that registry as
> authoritative, then contradicts itself here. Nothing new; the item stays
> where it is already tracked and is picked up only when the owner schedules it.

**Evidence.** [test_observations.md](test_observations.md#L20-L25) explicitly says
open/save/overwrite flows have no integration scenario. The current integration
suite renders the window and programmatically loads files, but never drives
`handleConfirm`, `confirmOverwrite`, folder navigation, filtering, or callbacks.
The file contains 289 statements and measures 20.4% coverage.

**Why it is wrong.** This dialog is the only user path for opening settings,
Save As, choosing output directories, and browsing output. Click ordering is
Gio-sensitive; a dead confirm button, wrong overwrite prompt, or navigation
regression can pass all current tests.

**Fix.** Extend the existing headless `input.Router` driver with stable semantic
lookup or dialog-level synthetic clicks. Cover Open selection, Save automatic
`.gen.json` suffix, existing-file overwrite cancel/confirm, folder selection,
hidden toggle, and unreadable-directory error retention. Do not expose private
methods as production API.

**Tests.** Add `test/integration/fileExplorer_integration_test.go` with the
`integration_test` tag. Split scenarios into named tests with one logical
assertion each per AGENTS.md.

### 6.3 🟠 ❌ INVALIDATED — Manual zone editor has no end-to-end pointer/property scenario

> **❌ INVALIDATED (2026-07-14): already explicitly tracked as deferred future
> work.** [test_observations.md](test_observations.md) records this exact gap
> for the five zone-editor dialog files ("future work — no integration scenario
> exists yet"), per AGENTS.md §4.6. Additionally, the proposed reset-contract
> scenario depends on §1.3, which is owner-deferred. Nothing new to action.

**Evidence.** [test_observations.md](test_observations.md#L27-L38) records that
canvas drag/connect, add/delete, snapping, and property panels have no
integration scenario. Business services are well unit-tested, but the five Gio
dialog files are 0% in the unit profile. Existing integration tests invoke
State-level manual snapshots rather than user pointer/click flows.

**Why it is wrong.** The editor's behavior depends on previous-frame geometry,
input routing order, sticky modes, dropdown updates, and dirty-cache
invalidation. Those are exactly the failure classes pure connection-editor
unit tests cannot catch. The known reset mismatch (§1.3) also survived because
no UI scenario exercises reopening/resetting.

**Fix.** Use the existing performance `AppRunner` input techniques in an
integration helper: open the dialog from Layout & Zones, select/drag a zone,
add and delete a connection, toggle snap, edit one property, Apply, and verify
`CurrentState().ManualZones/ManualConnections` plus rendered stability. Keep
coordinates calibrated from semantics/known canvas bounds rather than fixed
screen pixels where possible.

**Tests.** Add `test/integration/manualZoneEditor_integration_test.go`; add the
separate reset-contract scenario from §1.3 after the owner decision.

### 6.4 🟡 Two pure UI catalogs are entirely untested

**Evidence.** `bannableItems.go` (110 statements) and `valueOverrideSids.go`
(40 statements) both measure 0%. They are pure catalog builders/filter/sort
logic, unlike Gio-bound files. [bannableItems.go](../app/gui/constants/bannableItems.go#L24-L59)
concatenates more than 30 groups; [valueOverrideSids.go](../app/gui/constants/valueOverrideSids.go#L9-L27)
builds, excludes, and sorts selectable SIDs.

**Why it is wrong.** An omitted group, duplicate SID, accidental empty ID,
incorrect category, or exclusion bug silently removes or mislabels picker
choices. Protected registry constants need no tests, but the UI composition
logic is not protected data and has behavior.

**Fix/tests.** Add:

- `test/unit/app/gui/constants/bannableItems/getBannableItems_test.go` covering
  non-empty IDs/names/categories, SID uniqueness, representative entries from
  each group, and stable sorted/filter lookup behavior exposed by public APIs.
- `test/unit/app/gui/constants/valueOverrideSids/getValueOverrideSidsWithExclusions_test.go`
  covering sort order, exclusion, no duplicates, no empty SIDs, and caller
  input not mutated.

Do not snapshot the entire protected registry; assert invariants and a small
representative set.

### 6.5 🟡 Repository test-layout enforcement is absent

**Evidence.** [rmgTemplateModel_test.go](../test/integration/rmgTemplateModel_test.go#L1-L13)
omits the required tag, while the three tests linked in §3.1 import the GUI
layer; [AGENTS.md](../AGENTS.md#L297-L327) requires the opposite. Both policy
violations compile and passed CI.

**Why it is wrong.** Written conventions alone did not prevent either
regression, and standard compilation/lint excludes or accepts both cases.

**Fix/tests.** Add the repository-policy test or CI script proposed in §1.2
and extend it to reject `app/**` imports beneath `test/unit/internal/**`. This
single executable guard covers both policy classes without changing runtime
code.

> **Verification caveat (2026-07-14):** the finding is valid, but the suggested
> `test/unit/repository/testLayout/` location violates the AGENTS.md §4.6
> mirror convention (unit folders mirror production files — the review itself
> admits it is an "exception"). If adopted, implement as a CI script or a
> depguard scope instead of a unit-test folder.

---

## 7. CI/CD, security, and dependencies

### 7.1 🟠 Direct pushes to `master` skip tidy, lint, vulnerability, race, and coverage gates

**Evidence.** [pr-validation.yml](../.github/workflows/pr-validation.yml#L3-L8)
runs on both pushes and PRs, but `check-go-mod`, lint, vulnerability, race, and
coverage jobs are guarded with `github.event_name == 'pull_request'`; only
build, Windows checks, unit tests, and integration/performance run on direct
pushes.

**Why it is wrong.** If branch protection is bypassed or an administrative
push lands directly, master can receive untidy modules, lint violations,
reachable vulnerabilities, races, or coverage regressions. The workflow name
“PR Tests” obscures that it also handles pushes.

**Fix.** Prefer branch protection requiring PRs and all required checks. For
defense in depth, remove PR-only conditions from read-only gates (tidy, lint,
vulnerability, race, coverage floor); only the PR-commenting coverage report
needs PR context. Split coverage enforcement from comment publishing so push
runs can enforce the numeric floor without `pull-requests: write`. Rename the
workflow to validation.

**Tests.** Validate both events with `actionlint` and inspect a test branch push.
No Go test is required. **Owner/repository-settings decision required:** confirm
whether direct pushes are intentionally allowed.

### 7.2 🟠 ❌ INVALIDATED — Release publication still trusts a mutable third-party action tag

> **❌ INVALIDATED (2026-07-14): explicitly declined previously.** Commit
> `838306a` (historical review §7.5 hardening) added the SHA-pinned line
> **commented out** and deliberately kept `@v2` live — an explicit owner
> choice, already acknowledged by this review's own §0.1 ("carried"). This is
> owner-controlled supply-chain policy; agents must not flip the pin. Revisit
> only at the owner's initiative.

**Evidence.** [release.yml](../.github/workflows/release.yml#L87-L93) contains a
commented SHA but executes:

```yaml
uses: softprops/action-gh-release@v2
```

**Why it is wrong.** A mutable tag can move or be compromised after review; the
release job has `contents: write` and publishes executable binaries. Official
Actions also use major tags, but this third-party publishing action is the
highest-risk remaining pin from historical §7.5.

**Fix.** Verify the current upstream commit for the desired v2 release and pin
that full SHA, retaining a version comment. Configure Dependabot to keep SHA
pins updated (current Actions ecosystem update already exists). Optionally pin
all Actions for a uniform supply-chain policy.

**Verification.** Run `actionlint`; trigger a dry-run or draft release on a test
tag and verify checksums/artifacts.

### 7.3 🟡 Workflow permissions are not least-privilege by default

**Evidence.** [pr-validation.yml](../.github/workflows/pr-validation.yml#L1-L14)
has no top-level `permissions`; only the coverage job narrows its token. Most
jobs only need source read access.

**Why it is wrong.** Repository/org defaults determine each job's token scope.
Defaults can change, and third-party actions should not receive permissions they
do not need.

**Fix.** Add top-level `permissions: contents: read`; keep the existing narrow
coverage override (`actions: read`, `pull-requests: write`) and release job's
`contents: write`. Revalidate reusable composite-action behavior.

### 7.4 ⚪ Dependency and filesystem security posture

- Local `govulncheck` was unavailable and the user declined ephemeral install;
  [CI runs the official scanner](../.github/workflows/pr-validation.yml#L55-L75).
  No local zero-vulnerability claim is made.
- Root and tools module content is tidy; Windows `tidy -diff` output is
  line-ending-only (§0.2). CI currently checks only the root module. The tools
  module is intentionally excluded from Dependabot by prior owner decision.
- Filename sanitization removes Windows-invalid separators, and the in-app save
  dialog strips directory components before joining the current directory.
- No command injection surface was found in template generation; release shell
  interpolation uses repository-controlled tag/ref contexts, but tags should
  remain subject to repository release controls.

---

## 8. Linter disposition — ❌ ORIGINAL DATA INVALIDATED; corrected below

> **❌ INVALIDATED (2026-07-14): the review's count and table were wrong.**
> "50 issues" is golangci-lint's default `max-issues-per-linter: 50` display
> cap, not the real total. The 50 reported findings are a **nondeterministic
> sample** — consecutive runs on the identical tree return different file
> groupings — so the review's "exhaustive" table (and its jab that "the
> previous memory grouping was stale") is unreliable: it missed 33 findings
> (e.g. all of `app/gui/constants`, `connection_editor`, `bonusPickerDialog`)
> and miscounted `variant_content`. An uncapped run
> (`--max-issues-per-linter=0 --max-same-issues=0`) reports **84
> `gochecknoglobals` findings**. Corrected grouping:

| Group | Count | Notes |
| --- | ---: | --- |
| `app/gui/constants` (contentIds 1, gameModes 2, legend 1, mapSizes 3, roadDistances 1, spells 2) | 10 | UI display-name catalogs. |
| `app/gui/dialogs/bonusPickerDialog.go` | 4 | Dialog-level catalogs. |
| `app/gui/utils/lookupSid.go` | 1 | `allSidMappings` immutable lookup. |
| `internal/common` (mapSizes 3, topologies 2) | 5 | Exported + private catalogs. |
| `internal/models/config/generatorConfig.go` | 2 | |
| `internal/registry/unguardedContentPoolValues.go` | 6 | **Protected (§2.1) — agents must not edit.** |
| `internal/services/asset_provider/assetProvider.go` | 2 | |
| `internal/services/builders/placement_rule` (distance 5, placementRuleBuilder 1) | 6 | Exported distance presets. |
| `internal/services/builders/variant_content` | 10 | Immutable registry aliases. |
| `internal/services/connection_editor` (connectionEditor 5, zoneEditor 2) | 7 | Guard presets / labels. |
| `internal/services/content_rules` (distancePresets 6, variantMappingManager 7) | 13 | Rule vocabulary. |
| `internal/services/preview_service/previewGeneratorService.go` | 1 | |
| `internal/services/template_generator/providers` (common.go 16, topology/base/topologyBase.go 1) | 17 | Historical §3.4 territory. |
| **Total** | **84** | No runtime bug is implied. All of it remains **owner's responsibility** per historical §3.4 / §0.2 — agents skip. |

CI intentionally disables `gochecknoglobals` (plus `godox`, `dupl`, `unparam`)
as warning-class linters. The uncapped 2026-07-14 run confirms no other linter,
formatter, complexity, security, or duplicate finding — that part of the
original claim stands.

---

## 9. Documentation and developer experience

### 9.1 🟠 QUICKSTART's programmatic example cannot compile

**Evidence.** [QUICKSTART.md](../QUICKSTART.md#L112-L138) imports
`internal/services` and calls `services.WriteTemplate`, but that package/file no
longer exists. Current persistence API is `file_service.NewFileService().SaveTemplate`.

**Why it is wrong.** A user copying the advertised non-GUI example gets a
compile error. Because the module uses `internal`, the example is only usable
inside this module (or its parent tree), which is also not explained.

**Fix.** Update imports and save call to the current `file_service` API, or add
a supported public package if external programmatic use is truly a product
surface. Add a compile-checked example under
`test/unit/docs/programmaticUse/programmaticUse_test.go` or a small `examples/`
package built by CI so docs cannot drift again. **Owner decision required:**
choose internal-only example versus supported external API.

### 9.2 🟡 README, QUICKSTART, AGENTS, and memory describe a deleted UI/tree

**Evidence.** Current [window.go](../app/gui/editor/window.go#L32-L38) has three
tabs: General, Layout & Zones, Bonuses & Bans. Zone content is opened via per-zone
dialogs. Output/status/generate/save controls live in PreviewPanel. In contrast:

- [README.md](../README.md#L35-L72) lists deleted `internal/constants`, root
  preview/template/settings service files, a Zone Content panel, and footer.
- [README.md](../README.md#L79-L100) claims Gio v0.9.0 and four tabs; current
  [go.mod](../go.mod#L3-L6) is Go 1.26.5 / Gio v0.10.0.
- [README.md](../README.md#L220-L244) gives nonexistent `test/services`,
  `test/models`, and `test/models/template` paths.
- [QUICKSTART.md](../QUICKSTART.md#L3-L40) claims four tabs, Refresh, and a footer;
  [QUICKSTART.md](../QUICKSTART.md#L8-L12) requires Go 1.25.8.
- [QUICKSTART.md](../QUICKSTART.md#L99-L105) says Reveal opens the OS explorer,
  but current State opens the in-app Browse dialog.
- [AGENTS.md](../AGENTS.md#L7-L15) still says Go 1.26.3 / Gio v0.9.0.

**Why it is wrong.** Onboarding commands fail, contributors search for deleted
files, and users cannot reconcile the guide with the UI. AGENTS.md is
operational policy, so its stale snapshot can misdirect future agents.

**Fix.** In one docs-only change, generate the tree from current directories;
update three-tab dialog-based workflow, PreviewPanel controls, in-app Reveal,
Go/Gio versions, current service packages, and mirrored `test/unit` commands.
Update repository memory after docs land. Consider a small docs check for `go`
version and paths/commands.

**Tests.** Run every documented Go command in a clean checkout. Compile the
programmatic example per §9.1. No production change is needed.

---

## 10. Verified non-issues and rejected candidates

- Build, vet, default tests, gated integration/performance tests, and both
  suspicious stress suites passed; no flaky failure reproduced in 20 runs.
- `UpdateTemplate` checks variant length and clones the variant slice.
- Save/load callback ordering preserves loaded panel state; integration tests
  cover the former clobber regression.
- Unsaved state is set on actual updates and exit confirmation resets.
- Resource-density `/200` exactly matches the C# reference; do not “fix” to
  `/100` without a product decision.
- `linq.ToMap` nil-map and zone-editor range-copy concerns from old memory are
  stale/non-issues.
- No tracked ignored files, executables, generated output, or profiler artifacts
  were found.
- Root Go version matches validation CI (1.26.5); tools module is separately
  pinned to 1.26.3 and builds through root tooling workflows as intended.
- Protected game data/schema/registry were treated as authoritative, not
  refactoring targets.
- Deliberate `integration_test` architecture itself is valid; only the one
  missing file tag (§1.2) is wrong.

---

## 11. Suggested execution order

1. **Durability/correctness PR:** §1.1 atomic persistence, with failure-path
   tests. Decide whether template+PNG commit is transactional.
2. **Small test-policy PR:** §1.2 build tag + §3.1 test dependency direction +
   §6.5 executable policy checks.
3. ~~**Manual-editor product decision:** choose §1.3 semantics, then implement it
   with §6.3 interaction coverage.~~ ❌ Invalidated — §1.3 is owner-deferred;
   §6.3 already tracked in test_observations.md.
4. ~~**UI integration coverage:** §6.2 file explorer, then remaining §6.3 pointer
   workflows.~~ ❌ Invalidated — both already tracked as deliberate future work
   in test_observations.md; owner schedules them.
5. **Measured performance PR:** benchmark §4.1, add cache/revision key, prove
   invalidation and before/after timing.
6. **CI/security PR:** §7.3 permissions only. ~~§7.2 SHA pin~~ ❌ invalidated —
   owner explicitly declined the pin (commit `838306a`). Handle §7.1 only after
   confirming branch-protection/direct-push policy.
7. **Docs-only PR:** §9.1–§9.2 plus compile-checked example. Update repository
   memory afterward.
8. **Optional quality PR:** §6.4 pure catalogs. Leave §8 global-catalog cleanup
   to the owner as already decided.

---

## 12. Measured baselines

| Check | Result |
| --- | --- |
| Reviewed revision | `3ac96111109ccf28c90b2386cbea5d4f840efb1c` (`AD/refactoring-07-13`) |
| `go build ./...` | Pass |
| `go vet -tags=integration_test ./...` | Pass |
| `go test ./test/... -count=1` | Pass; unexpectedly includes one integration test (§1.2) |
| Tagged integration + performance suites | Pass |
| Integration stress `-count=20` | Pass |
| Generator/driver suspicious suites `-count=20` | Pass |
| Unit coverage | **62.2%**; 178 files, 48 below 80%, 35 at 0% |
| Configured lint | ~~**50**~~ **84**, all `gochecknoglobals` (review's 50 was the default `max-issues-per-linter` display cap; see §8) |
| Vulnerability scan | CI-covered; not run locally because executable absent and ephemeral install was declined |
| Root/tools tidy | Semantically tidy; local Windows `-diff` is line-ending-only (58/58 and 985/985 identical lines) |
| Tracked ignored files | 0 |
| Working tree before review document | Clean |

Finding count: ~~**1 🔴, 7 🟠, 6 🟡, 2 ⚪**~~ — after the 2026-07-14
verification pass, the actionable set is **1 🔴 (§1.1), 4 🟠 (§1.2, §4.1,
§7.1, §9.1), 5 🟡 (§3.1, §6.4, §6.5, §7.3, §9.2), 2 ⚪**. Invalidated:
§1.3, §6.2, §6.3, §7.2, plus §8's original data (corrected in place).
