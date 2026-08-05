# Project Review — 2026-08-04 (Claude Opus 5)

Senior/principal-level review of the whole repository, verified at commit
`687f47d6cff07dd2f42239c796dd8dad5385931a` on branch `AD/refactoring-07-21`,
working tree clean.

**Toolchain measured for this review:**

| Item | Value |
| --- | --- |
| `go version` | `go1.26.5 windows/amd64` |
| Root module | `go 1.26.5` |
| `tools/` module | `go 1.26.3` (separate module) |
| Gio | `gioui.org v0.10.0` |
| `golangci-lint-v2 version` | `v2.12.2`, built with `go1.26.3` |
| Lint issues (uncapped, configured set) | **42** — 40 `gochecknoglobals`, 2 `dupl` |
| Unit coverage | **64.7%** (`go tool cover -func`); 213 files, 49 below 80%, 36 at 0% |
| `govulncheck` (symbol scan) | **1 reachable vulnerability**, 2 at module level |

**Supersession.** This document **supersedes
[todo/review-gpt-5.6-sol-07-13.md](review-gpt-5.6-sol-07-13.md)** in full
(which itself superseded the deleted `todo/review-fable-07-12.md`). Once this
review is accepted, `review-gpt-5.6-sol-07-13.md` can be deleted; every one of
its open items is dispositioned in §0 and, where still valid, restated here
with fresh evidence.

This document does **not** supersede the two live tracking files
[todo/backlog.md](backlog.md) and [todo/test_observations.md](test_observations.md).
They remain authoritative for owner future-work and intentional test gaps; §0.2
dispositions their contents so nothing is silently re-reported.

**Severity legend:**
🔴 **High** — bug / correctness / user-visible data loss ·
🟠 **Medium** — architecture, performance, CI/security gap ·
🟡 **Low** — readability, hygiene, docs ·
⚪ **Informational** — no action required, recorded for completeness.

**Finding count: 3 🔴 · 17 🟠 · 22 🟡 · 4 ⚪ (46 total).**

---

## 0. Disposition of prior reviews, backlog, observations, and memories

### 0.1 Fixed / no longer reproducible ✅

The prior review's own §0.1 table (35 rows carried from the deleted
`review-fable-07-12.md`) was re-spot-checked. All rows I could re-verify remain
fixed; the ones with fresh evidence gathered this session:

| Prior item | Re-verification at `687f47d` |
| --- | --- |
| Historical §1.5 unknown victory coercion | **Still fixed.** `constants.GetVictoryCondition` returns `(Victory, bool)`; a `.gen.json` with an unknown value is caught by [editorStateValidator.go](../internal/validators/editorStateValidator.go#L115-L133) and reset to `Standard` with a warning. |
| Historical §1.6 `UpdateTemplate` aliases `Variants` | **Still fixed.** |
| Historical §4.1 per-frame reflection | **Still fixed.** [editorStateDto.go](../internal/dtos/editorStateDto.go#L186-L201) uses the hand-rolled `EqualsIgnoringManualEdits`. |
| Historical §4.2 static tab allocations | **Still fixed.** [window.go](../app/gui/editor/window.go#L74-L86) caches `tabChildren`. |
| Historical §7.7 Dependabot | **Still fixed.** [.github/dependabot.yml](../.github/dependabot.yml) covers root `gomod` + `github-actions`, weekly, grouped. |
| Historical §7.6 tracked output artifacts | **Still fixed.** `git ls-files --ignored --exclude-standard -c` → 0; no tracked `*.exe/dll/so/prof/out/test/html/info/txt`. |
| Prior §0.2 "`go mod tidy -diff` fails on Windows (EOL artifact)" | **No longer reproducible.** `go mod tidy -diff` now exits **0** for both the root module and `tools/`. Item closed. |
| Prior §1.5 / historical zone-editor quality-index panic | **Fixed.** [zoneEditorZoneProps.go](../app/gui/dialogs/zoneEditorZoneProps.go#L60-L66) now uses `SelectByName(quality.GetName())` instead of indexing `QualityLabels`. |
| Prior §8 lint baseline 84 `gochecknoglobals` | **Improved.** Uncapped run now reports **40**; `internal/registry`, `content_rules`, `connection_editor`, `placement_rule`, `common` and `generatorConfig` globals are no longer reported. |
| Prior §6.1 coverage 62.2% | **Improved to 64.7%.** |

### 0.2 Invalidated / accepted / owner-controlled ✖ — do not re-report

| Item | Disposition |
| --- | --- |
| Prior **§1.3** zone-editor "Reset to generated" semantics | **Owner-deferred.** The in-code comment at [zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go#L213) still reads `// Reset only resets current edits, not all manual edits, need to fix eventually. wont add todo so the llm does not trigger`. Confirmed present and untouched. Agents must not action it. It is also the **only** TODO/FIXME/HACK-class comment in the non-test Go tree. |
| Prior **§6.2** file-explorer interaction test | **Already tracked** in [test_observations.md](test_observations.md). Not re-reported. |
| Prior **§6.3** manual zone-editor interaction test | **Already tracked** in [test_observations.md](test_observations.md). Not re-reported. |
| Prior **§7.2** release action SHA pin | **Owner declined** (commit `838306a` left the pin commented out deliberately). Note only: the action has since moved to [release.yml](../.github/workflows/release.yml#L94) `softprops/action-gh-release@v3`; still a mutable major tag, still owner territory. |
| Prior **§8** `gochecknoglobals` catalogue cleanup | **Owner's responsibility** (historical §3.4). CI disables the linter deliberately. §10 gives the full current inventory for the record only. |
| Prior §0.2 historical §5.3 exported funcs returning private types | **Owner-retained API style.** |
| Prior §0.2 historical §5.4 default `slog` logger global | **Accepted configuration** (`sloglint` clean). See §1.9 for a *different*, behavioural aspect of the same code that is a new finding. |
| Prior §0.2 resource-density `/200` | **Validated compatibility behaviour**, matches the C# reference. Not a finding. |
| Prior §7.4 "tools module intentionally excluded from Dependabot" | **Owner decision.** §7.4 below therefore reports only the *CI build/test/lint* gap for `tools/`, not the Dependabot gap. |
| Backlog: preview sub-pixel `Vec2` | Deliberate future work — [backlog.md](backlog.md). |
| Backlog: remove `[2]float64` from template entities | Protected directory (§2.1 AGENTS.md) — owner-only. |
| Backlog: `createTopologyAdjacency` dead Chain/Ring branches | **Explicit owner retention** — removal was implemented and then rolled back. Do not re-report or remove. |
| Backlog: zone tier property on entities; consolidate road distances; "app should only use entities/models/handlers/commons"; common generation values; rework `EditorStateDto`; rename `template` → `template_entity` | Owner future-work. §2.2 below overlaps the "app should only use…" item; it is raised as a *current* layering finding with evidence, not as a backlog restatement. |
| test_observations: Gio widgets/dialogs/panels at 0% | **Accepted integration territory** under AGENTS.md §4.6. §6.4 raises only the two *non-Gio* catalogs. |
| test_observations: `drivers.State` dialog-bound branches, `topologyConnectionService` private policy, unreachable defensive branches (`connectInteriorStables` len==0, `providePreviewGenerator` err!=nil) | Accepted, already registered. |
| Repo memory: "GUI has four tabs / Zone Content tab / footer panel" | **Stale memory.** [window.go](../app/gui/editor/window.go#L36-L40) has three tabs. Docs still repeat it — see §9.2. |
| Repo memory: "coverage 86–92%, Go 1.26.3 root, Gio v0.9.0" | **Stale.** Current: 64.7%, Go 1.26.5, Gio v0.10.0. |
| Repo memory: "arena assets embedded but not decoded/drawn" | **Confirmed still true this session** — promoted to a real finding with evidence at §2.7 (it is a shipped-binary cost plus a documented-but-absent feature, not just a note). |

### 0.3 Carried forward — re-verified against source this session ❗

| Prior item | Status at `687f47d` | New section |
| --- | --- | --- |
| §1.1 non-atomic persistent writes | Still open, verbatim | §1.1 |
| §1.2 `rmgTemplateModel_test.go` missing build tag | Rejected — see §6.1 | §6.1 |
| §3.1 internal-service tests import `app/gui/constants` | Still open; **file paths changed** (`variantMappingCatalog`, `contentRuleService`) | §6.3 |
| §4.1 live preview rebuilds layout every frame | Still open; now **profiler-measured** and routed through `previewHandler` | §4.1 |
| §6.4 `bannableItems.go` / `valueOverrideSids.go` at 0% | Still 0% | §6.4 |
| §6.5 no executable test-layout enforcement | Still absent | §6.5 |
| §7.1 direct pushes to master skip gates | Still open | §7.1 |
| §7.3 no top-level workflow `permissions:` | Still open | §7.2 |
| §9.1 QUICKSTART programmatic example cannot compile | Still open, now with **three** distinct compile errors | §9.1 |
| §9.2 README/QUICKSTART/AGENTS describe a deleted UI | Still open; AGENTS.md drift changed shape (versions now correct, module count and task labels wrong) | §9.2, §9.6 |

---

## 1. Bugs & correctness

### 1.1 ✅ FIXED 🔴 Persistent user files are replaced non-atomically

**Fixed 2026-08-05.** The owner replaced this item's proposed design (a private
`AtomicFileWriter` inside `internal/services/file_service/`) with a repository
layer, reusing the pre-existing but unimplemented
[fileRepositoryInterface.go](../internal/repositories/fileRepositoryInterface.go):

- [atomicFileWriter.go](../internal/repositories/atomicFileWriter.go) — private
  to `internal/repositories`. `MkdirAll` → write `{dir}/TEMP-{name}{ext}` →
  `Sync` → `Close` (close error wins when the encode succeeded) → bounded
  `os.Rename` retry (5 attempts, 20 ms apart) onto `{dir}/{name}{ext}`. Every
  failure path removes the temporary file and leaves the destination untouched.
- [editorStateRepository.go](../internal/repositories/editorStateRepository.go),
  [templateRepository.go](../internal/repositories/templateRepository.go),
  [previewRepository.go](../internal/repositories/previewRepository.go) own the
  extension (`.gen.json` / `.rmg.json` / `.png`) and the encoding.
- [fileService.go](../internal/services/file_service/fileService.go) is now a
  pure controller: it decides the directory and the requested name and hands
  them to a repository, which sanitizes the name (falling back to
  `Generated_Template`) inside the writer. It gained
  `SaveTemplateWithPreview` — a single call that skips the preview when the
  image is `nil` — and `SaveTemplate` / `SavePreviewImage` were removed.

Point 2's Windows caveat is honoured by the retry, except that a final failure
removes the temporary file rather than leaving it: the destination still holds
the previous valid contents, so the half-written copy is the one to lose.

Point 4 was decided by the owner in favour of **keeping the documented
partial-success contract**: a preview failure still returns the template path so
[stateGeneration.go](../app/gui/drivers/stateGeneration.go#L79-L85) can report it.

Two deliberate behaviour changes follow from routing everything through the
repositories, both approved by the owner:

- Editor-state saves now create missing directories instead of failing.
- The Save As dialog's file name is discarded; the state is written as
  `{TemplateName}.gen.json`. `SaveState` therefore returns the path actually
  written and `drivers.State` records **that**, not the requested path.

**Tests.** The fault injection needed no test-only seam: `math.NaN()` in a float
field makes `json.MarshalIndent` fail, and a zero-sized `image.RGBA` makes
`png.Encode` fail. Each repository has a mirror folder under
`test/unit/internal/repositories/` covering "destination byte-identical after a
failed save", "no `TEMP-*` residue", the extension, directory creation and the
round trip. The `FileService` folder was rewritten against mocked
`IFileRepository[T]` collaborators and now asserts only routing decisions.
`atomicFileWriter` is private and is covered indirectly — recorded in
[test_observations.md](test_observations.md).

**Original finding follows.**

**Evidence.** [fileService.go](../internal/services/file_service/fileService.go#L45-L53):

```go
// SaveSettings writes settings object to the given filepath as JSON.
func (this *FileService) SaveSettings(filePath string, editorState *dtos.EditorStateDto) error {
	data, err := json.MarshalIndent(editorState, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, fileReadWritePermission)
}
```

and [fileService.go](../internal/services/file_service/fileService.go#L67-L75):

```go
	out := filepath.Join(directory, safeName+templateExtension)
	data, err := json.MarshalIndent(template, "", "  ")
	...
	if err = os.WriteFile(out, data, fileReadWritePermission); err != nil {
		return "", err
	}
```

**Why it is wrong.** `os.WriteFile` opens with `O_TRUNC`: the previous contents
are destroyed *before* the new bytes are durable. A crash, power loss, or
disk-full between truncate and completion leaves a zero-length or half-written
`.gen.json` — the user's authored editor state — with no recoverable copy. The
`.rmg.json` output has the same exposure.

**Fix.**

1. Add a private helper in `internal/services/file_service/` (own file per
   AGENTS.md §4.1), e.g. `atomicFileWriter.go` with an `AtomicFileWriter`
   struct: `os.CreateTemp(filepath.Dir(dest), ".tmp-*")` → write → `Chmod(0o644)`
   → `Sync()` → `Close()` → replace destination → `os.Remove(temp)` on every
   error path.
2. **Windows caveat:** `os.Rename` over an existing destination works on modern
   Windows Go runtimes, but fails if the destination is open by another process
   (a common case: the game or an editor holding the `.rmg.json`). Wrap the
   rename with a bounded retry and, on final failure, leave the temp file in
   place and return an error naming it — never delete the only valid copy.
3. Route `SaveSettings`, `SaveTemplate` and `SavePreviewImage` (§1.6) through it.
4. **If investigation shows** the two-file template+PNG write must be
   transactional (JSON succeeds, PNG fails), write *both* temps first and commit
   both, or explicitly document the partial-success contract that
   [stateGeneration.go](../app/gui/drivers/stateGeneration.go#L79-L85) already
   surfaces to the user.

**Tests to add** (AGENTS.md §4.6 mirror layout):

- `test/unit/internal/services/file_service/atomicFileWriter/write_test.go` —
  destination byte-identical when the write step fails; no `.tmp-*` residue;
  destination replaced on success; permissions `0o644`.
- Extend `test/unit/internal/services/file_service/fileService/saveSettings_test.go`
  and `saveTemplate_test.go` with "pre-existing destination survives a failed
  save" cases.
- Fault injection must come from a narrow interface consumed by production code
  (e.g. `IFileWriter`), **not** from a test-only export.

**Owner-decision flag:** ⚠ choose whether template JSON + preview PNG commit
transactionally or keep the documented partial-success contract.

---

### 1.2 ✅ FIXED 🔴 `SaveAs` records the new path even when the save failed

**Fixed 2026-08-05.** `handleSaveState` now returns `bool` (mirroring
`handleLoadState`) and the `SaveAs` callback only assigns `this.currentPath`
when that write succeeded. The two statement-position callers (`Save` and the
`integration_test`-gated `SaveStateToFile` export) needed no change.

The regression test could **not** be written as a unit test: the decision lives
in the closure handed to `dialogs.NewSaveFileDialog`, which stores it in the
unexported `onSave` field, and `onSave` normally fires only from
`confirmSelection` / `confirmOverwrite` — both of which require a
`layout.Context`. Per this item's own fallback clause the scenario went to
[stateSaveAs_integration_test.go](../test/integration/stateSaveAs_integration_test.go)
(`TestWhenSaveAsFails_CurrentPathIsNotRecorded` /
`TestWhenSaveAsSucceeds_CurrentPathIsRecorded`), backed by two new
`*_testexports.go` accessors — `DialogHost.GetTopDialog` and
`FileExplorerDialog.ConfirmSave` — which is exactly what AGENTS.md §4.6.1
means by the `integration_test` tag. The unit-level gap is recorded in
[test_observations.md](test_observations.md).

**Evidence.** [stateFiles.go](../app/gui/drivers/stateFiles.go#L39-L49):

```go
func (this *State) SaveAs(templateName string) {
	dir, err := os.Getwd() // Editor state by default is saved in the same directory as the executable
	if err != nil {
		dir = this.suggestDirectory()
	}
	defaultName := helpers.SanitizeFilename(strings.TrimSpace(templateName)) + configFileExtension
	this.dialogs.Open(dialogs.NewSaveFileDialog(dir, defaultName, func(path string) {
		this.handleSaveState(path)
		this.currentPath = path
	}))
}
```

`handleSaveState` returns nothing and swallows the outcome
([stateFiles.go](../app/gui/drivers/stateFiles.go#L85-L97)):

```go
func (this *State) handleSaveState(path string) {
	if _, err := this.handler.SaveState(dtos.EditorStateSaveDto{...}); err != nil {
		this.SetStatus(fmt.Sprintf("Save failed: %v.", err), true)
		return
	}

	this.unsaved = false
	...
}
```

**Why it is wrong.** On a failed *Save As* (read-only directory, invalid path,
disk full), the status line says "Save failed" but `currentPath` is now set to a
path that holds no file. Subsequent `Save`
([stateFiles.go](../app/gui/drivers/stateFiles.go#L30-L37)) takes the
`currentPath != ""` branch and silently retargets the same broken path instead
of re-prompting. Worse, the toolbar
([toolbar.go](../app/gui/editor/toolbar.go#L78-L85)) now displays that path as
the active file *without* the unsaved `*` marker only after a later successful
save — meaning the user believes their work is filed under a location that was
never written. Contrast with `Load`, where the equivalent bug was already fixed:
[stateFiles.go](../app/gui/drivers/stateFiles.go#L23-L27) only runs `onLoaded`
when `handleLoadState` returns `true`.

**Fix.** Make `handleSaveState` return `bool` (mirroring `handleLoadState`) and
gate the assignment:

```go
	this.dialogs.Open(dialogs.NewSaveFileDialog(dir, defaultName, func(path string) {
		if this.handleSaveState(path) {
			this.currentPath = path
		}
	}))
```

Update the sole other caller at [stateFiles.go](../app/gui/drivers/stateFiles.go#L36)
to ignore the result.

**Tests to add.**
`test/unit/app/gui/drivers/stateFiles/saveAs_test.go` (new file in the existing
mirrored folder): with a stub `IGuiHandler` whose `SaveState` returns an error,
assert `GetCurrentPath()` is unchanged and `IsUnsaved()` stays `true`; with a
succeeding stub, assert the path is recorded. The dialog callback is reachable
because `DialogHost` exposes the active dialog — **if investigation shows** the
save-file dialog cannot be driven headlessly from a unit test, add the scenario
to `test/integration/stateSaveAs_integration_test.go` with the `integration_test`
tag instead, and record the unit-level gap in
[test_observations.md](test_observations.md).

---

### 1.3 ✅ FIXED 🟠 `WasLayoutChanged` dereferences a possibly-nil `previous` state

**Fixed 2026-08-05.** `WasLayoutChanged` now guards with
`this.HasPreviousState() && ...`. The owner also asked for the follow-on
cleanup, so `ShouldReapplyManualEdits` dropped its now-redundant
`!this.HasPreviousState() ||` short-circuit and reads as `HasManualEdits() &&
!WasLayoutChanged()`. Regression test:
`TestWhenNoPreviousStateExists_ReportsLayoutNotChanged` in
[wasLayoutChanged_test.go](../test/unit/app/gui/models/editorState/wasLayoutChanged_test.go).

**Evidence.** [editorState.go](../app/gui/models/editorState.go#L92-L94):

```go
func (this *EditorState) WasLayoutChanged() bool {
	return this.previous.LayoutDefiningOptionsChanged(this.current)
}
```

Every sibling guards ([editorState.go](../app/gui/models/editorState.go#L84-L98)):

```go
func (this *EditorState) WasStateChanged() bool {
	return this.HasPreviousState() && !this.previous.EqualsIgnoringManualEdits(this.current)
}
...
func (this *EditorState) WasLayoutUnchanged() bool {
	return this.HasPreviousState() && !this.previous.LayoutDefiningOptionsChanged(this.current)
}
```

`previous` is explicitly `nil` after `OverrideState`
([editorState.go](../app/gui/models/editorState.go#L25-L29)) — i.e. after every
`New` and every settings-file load — and after `ResetPreviousState()`.

**Why it is wrong.** `LayoutDefiningOptionsChanged` is a pointer method that
immediately reads `this.PlayerCount`
([editorStateDto.go](../internal/dtos/editorStateDto.go#L147-L156)), so a nil
receiver panics and takes the whole GUI down. The exported method is currently
safe only by accident, and only along one path:
`ResetNextStateIfLayoutChanged` ([editorState.go](../app/gui/models/editorState.go#L53-L60))
does **not** guard, and is saved solely by `AutoRegenerate` checking
`HasPreviousState()` several lines earlier in a *different file*
([stateGeneration.go](../app/gui/drivers/stateGeneration.go#L36-L47)). The one
caller that does guard, `ShouldReapplyManualEdits`
([editorState.go](../app/gui/models/editorState.go#L139-L144)), relies on `||`
short-circuit — a refactor that reorders that expression is an instant crash.

**Fix.** Add the guard to the method itself so the invariant lives with the data:

```go
func (this *EditorState) WasLayoutChanged() bool {
	return this.HasPreviousState() && this.previous.LayoutDefiningOptionsChanged(this.current)
}
```

Then simplify `ShouldReapplyManualEdits` to `return !this.WasLayoutChanged()`
after the `HasManualEdits` check — **verify** this preserves the documented
"right after a load there is no previous state; the stored edits are then
trusted" semantics (it does: no previous ⇒ `WasLayoutChanged()` false ⇒ reapply).

**Test that would have caught this.** Add
`test/unit/app/gui/models/editorState/wasLayoutChanged_test.go` →
`TestWhenNoPreviousStateExists_ReportsLayoutNotChanged`: construct via
`newEditorState()` **without** calling `SnapshotCurrentState`, call
`WasLayoutChanged()`, assert `false`. Both existing tests in that file snapshot
first, which is exactly why the nil path was never exercised.

---

### 1.4 🟠 Editor-state copies are shallow, so snapshots alias live slices

**Evidence.** `EditorStateDto` holds nine slice fields
([editorStateDto.go](../internal/dtos/editorStateDto.go#L81-L98)): `Bonuses`,
six `[]models.ZoneContentRowSave`, `ManualZones`, `ManualConnections`. Three
places copy the struct by value and treat the result as an independent snapshot:

[editorState.go](../app/gui/models/editorState.go#L31-L33):

```go
func (this *EditorState) GetCurrentState() dtos.EditorStateDto {
	return *this.current
}
```

[editorState.go](../app/gui/models/editorState.go#L41-L45):

```go
func (this *EditorState) SnapshotCurrentState() {
	previousState := *this.current
	this.previous = &previousState
	// this.next = nil
}
```

[stateHandler.go](../internal/handlers/stateHandler.go#L59-L74) —
`ValidateEditorState(stateDto dtos.EditorStateDto, ...)` takes the DTO **by
value** and returns `dtos.EditorStateValidationDto{State: stateDto, ...}`.

**Why it is wrong.** A struct copy duplicates slice *headers*, not backing
arrays. `this.previous` therefore shares element storage with `this.current`.
Change detection compares them element-wise
([editorStateDto.go](../internal/dtos/editorStateDto.go#L186-L201) →
`contentRowSlicesEqual` / `slices.Equal`), so **any in-place element write makes
the change invisible to `WasStateChanged()`** — the editor would not mark the
file unsaved and `AutoRegenerate` would not regenerate. The same aliasing leaks
the live editor state out of `GetCurrentState()` to every panel and to
`previewPanel`.

This is currently **latent, not live**: every writer replaces the whole slice
rather than an element —
[layoutPanelZones.go](../app/gui/panels/layoutPanelZones.go#L133-L148) assigns
`s.PlayerZoneContentRows = rows` etc.,
[bonusesPanel.go](../app/gui/panels/bonusesPanel.go#L248) reslices with a
capacity-limited three-index slice, and
[editorState.go](../app/gui/models/editorState.go#L116-L119) reassigns
`ManualZones`/`ManualConnections`. One in-place edit anywhere reintroduces the
bug with no compiler or lint signal.

**Fix.** Give `EditorStateDto` an explicit deep copy and use it at the three
sites above:

```go
// Clone returns a copy that shares no backing array with the receiver.
func (this *EditorStateDto) Clone() EditorStateDto { ... slices.Clone each slice field ... }
```

Place it beside the existing methods in
[editorStateDto.go](../internal/dtos/editorStateDto.go). Note the *elements*
(`config.BonusEntry`, `models.ZoneContentRowSave`,
`editor_state_dto.ManualZoneSave`) must themselves be checked for nested slices
— **if investigation shows** they contain slices (e.g. rule rows inside a zone
content row), clone recursively and add the tripwire test below.

**Tests to add.**

- `test/unit/internal/dtos/editorStateDto/clone_test.go` — mutate every slice
  field of the clone in place, assert the original is unchanged (one assertion
  per field, per AGENTS.md "one unit per test").
- `test/unit/app/gui/models/editorState/snapshotCurrentState_test.go` →
  `TestWhenSnapshotTakenAndContentRowMutatedInPlace_ReportsStateChanged` — this
  is the regression test that locks the behaviour in.

---

### 1.5 ✅ FIXED 🟠 Loaded `.gen.json` integer counts have no upper bound

**Fixed 2026-08-05.** `nonNegativeIntFields()` / `validateNonNegativeFields()`
are gone; all twenty count fields moved into `rangedIntFields()` in
[editorStateValidator.go](../internal/validators/editorStateValidator.go), which
is now a single declarative table of every bounded integer. Ceilings are named
constants taken from the matching editor slider:

| Field(s) | Range |
| --- | --- |
| `neutralZoneCount` | 0..16 |
| `abandonedOutpostCount`, `playerOwnedCastles`, `playerCastles`, `neutralCastles`, `hubCastles`, `remoteFootholdCount` | 0..4 |
| `maxPortalConns` | 0..32 |
| the eight neutral tier counts | 0..8 |
| the four castles-per-zone | 0..4 |
| `playerZoneSize`, `neutralZoneSize`, `hubZoneSize` | 0.5..2.0 |
| `guardRandomization` | 0.0..0.5 |
| `lostStartCityDay`, `cityHoldDays`, `gladiatorArenaCountDay` | 1..30 |
| `gladiatorArenaDaysDelayStart` | 1..60 |
| `tournamentFirstTournamentDay`, `tournamentInterval` | 3..30 |
| `tournamentPointsToWin` | 1..10 |

**Owner decisions recorded:**

- Ceilings accepted as proposed (they mirror the GUI slider maxima).
- The twenty fields that previously lived in `nonNegativeIntFields()` **keep a
  floor of 0** even where their slider starts at 1, so states that were valid
  before this change keep loading without new warnings. The newly-validated
  rules-tab fields have no such history and take their slider minimum.
- Zone sizes are `0.5..2.0`: `MultiplierFormatter` has base 0.5 and scale 1.5
  over a slider value in `[0, 1]`, so label and persistence agree. The `1.5`
  discrepancy reported during implementation was incorrect — do not "fix" it.
- The `gladiatorArena*` fields are validated even though §2.7 may later remove
  the feature; if the feature goes, its two entries go with it.

Floats are handled by a `rangedFloatFields()` sibling backed by new
[floatField.go](../internal/validators/floatField.go) and
[rangedFloatField.go](../internal/validators/rangedFloatField.go), reusing
`helpers.Clamp`. `validateTopology` mirrors `validateGameMode` and falls back to
`config.TopologyRandom`.

The negative-count message changed from `"%s %d is negative"` to
`"%s %d is outside [0, N]"`; the guiHandler test asserting the old wording was
updated.

**Original finding follows.**

#### 1.5 original text

**Evidence.** [editorStateValidator.go](../internal/validators/editorStateValidator.go#L62-L77)
only clamps at zero:

```go
	for _, field := range this.nonNegativeIntFields() {
		currentValue := *field.value(state)
		if currentValue >= 0 {
			continue
		}
```

and the twenty fields in
[editorStateValidator.go](../internal/validators/editorStateValidator.go#L172-L211)
— `neutralZoneCount`, `abandonedOutpostCount`, all eight neutral tier counts,
all four castles-per-zone, `hubCastles`, `remoteFootholdCount`,
`maxPortalConns` — carry no maximum. `playerCount` and the percentage fields
*are* range-checked
([editorStateValidator.go](../internal/validators/editorStateValidator.go#L134-L170)),
which shows the mechanism exists and simply was not applied here.

Additionally **not validated at all**: `PlayerZoneSize`, `NeutralZoneSize`,
`HubZoneSize`, `GuardRandomization` (all `float64`), `LostStartCityDay`,
`CityHoldDays`, `GladiatorArena*`, `Tournament*`, and `Topology`.

**Why it is wrong.** `.gen.json` is a plain-text file users hand-edit and share.
`{"neutralZoneCount": 100000000}` passes validation untouched and reaches the
generator, which allocates one `entities.Zone` per zone — an unbounded
allocation on the UI thread with no cancellation. A negative `playerZoneSize`
produces geometrically invalid templates the game must reject. This is the
classic "trust the file on disk" input-validation gap; the validator exists
precisely to be the boundary.

**Fix.** Convert the twenty entries from `nonNegativeIntFields()` to
`rangedIntFields()` with defensible ceilings derived from the UI sliders (the
GUI already constrains them — read the `Max` of the corresponding slider in
[layoutPanelZones.go](../app/gui/panels/layoutPanelZones.go) and
[generalPanel.go](../app/gui/panels/generalPanel.go) and reuse those numbers as
named constants). Add a `rangedFloatFields()` sibling for the four float fields.
Add a `validateTopology` mirroring `validateGameMode`.
Every clamp already produces a user-visible warning through the existing
`ValidationIssue.Message` pipeline, so no new UI work is needed.

**Tests to add.** In the existing folder
`test/unit/internal/validators/editorStateValidator/validate_test.go`: one test
per newly-bounded field asserting the issue message and the clamped value
(table-driven with named `t.Run` subtests is acceptable here per AGENTS.md
§4.6). Use `gofakeit` to fuzz values above the ceiling.

**Owner-decision flag:** ⚠ the exact ceilings are a product decision — propose
them, get confirmation, then implement.

---

### 1.6 ✅ FIXED 🟠 PNG export truncates the destination before encoding and discards the close error

**Fixed 2026-08-05** as part of §1.1. `SavePreviewImage` no longer exists; the
preview is written by
[previewRepository.go](../internal/repositories/previewRepository.go) through the
shared `atomicFileWriter`, so `png.Encode` writes into `TEMP-{name}.png` and the
destination is only replaced once encoding succeeded. The writer uses exactly the
named-return `defer` this item prescribes, so a `Close` error wins whenever the
encode itself succeeded.

The "close failure is propagated" test was **not** written: forcing `Close` to
fail requires either a test-only seam in production code (forbidden by
AGENTS.md §4.6) or filling a real filesystem. The truncation half is covered by
`TestWhenEncodingFailsOverAnExistingPreview_LeavesTheDestinationUntouched` in
`test/unit/internal/repositories/previewRepository/save_test.go`, and the gap is
recorded in [test_observations.md](test_observations.md).

**Original finding follows.**

**Evidence.** [fileService.go](../internal/services/file_service/fileService.go#L95-L107):

```go
	out := filepath.Join(directory, safeName+pngExtension)
	file, err := os.Create(out)
	if err != nil {
		return "", err
	}

	defer file.Close()
	if err = png.Encode(file, previewImage); err != nil {
		return "", err
	}

	return out, nil
```

**Why it is wrong.** Two distinct defects:

1. `os.Create` truncates the previous preview *before* `png.Encode` runs. An
   encode failure (or disk-full mid-encode) leaves a corrupt PNG where a valid
   one used to be. This is the §1.1 problem in its most reachable form, since
   `SavePreviewImage` is called on every *Save Template*.
2. `defer file.Close()` discards the close error. For a buffered write path a
   failed `Close` is where an out-of-space condition surfaces, so
   `SavePreviewImage` can return `(path, nil)` for a truncated file. `errcheck`
   does not flag it because the configured exclusions permit deferred `Close`.

**Fix.** Route through the §1.1 atomic writer and replace the bare defer with:

```go
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
```

using named return values — or drop the defer entirely and close explicitly on
both paths, which the atomic-writer helper makes natural.

**Tests to add.** `test/unit/internal/services/file_service/fileService/savePreviewImage_test.go`
(extend): a pre-existing PNG is byte-identical after an injected encode failure;
a close failure is propagated as an error.

---

### 1.7 ✅ FIXED 🟡 Malformed guard-value override lines are dropped without telling the user

**Fixed 2026-08-05.** `CreateValueOverrides` now returns
`([]entities.ValueOverride, []string)`; the line parsing moved into a
`parseValueOverride` helper that returns either the override or the warning
explaining the rejection. Blank lines stay silent, and line numbers count the
raw source lines so they match what the user sees in the editor:

- `line 4: 'foo:30000' is not sid=value` — no `=`, or `=` in first position
- `line 4: 'foo=abc' has a non-numeric value`

The "empty sid" case the original finding lists is **unreachable**: the line is
trimmed before parsing, so an empty prefix implies `=` at index 0, which the
first branch already rejects. That branch was written, proven dead by a failing
test, and removed rather than left as untestable code.

**Plumbing.** `TemplateGenerator.Generate()` now returns
`(*entities.RmgTemplate, []string)`. `templateHandler` concatenates the warnings
onto `validation.Warnings`, so they reach `dto.Warnings` and the existing
status-bar counter. This was first done as a separate `GenerateWithWarnings()`
method to avoid touching ~70 test call sites, but the owner rejected carrying two
entry points and had the signature changed outright.

**Original finding follows.**

#### 1.7 original text

**Evidence.** [gameRulesProvider.go](../internal/services/template_generator/providers/gameRulesProvider.go#L38-L57):

```go
func (this *GameRulesProvider) CreateValueOverrides(configuration config.GeneratorConfig) []entities.ValueOverride {
	var overrides []entities.ValueOverride
	for line := range strings.SplitSeq(configuration.ValueOverridesText, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		equals := strings.Index(line, "=")
		if equals <= 0 {
			continue
		}
		sid := strings.TrimSpace(line[:equals])
		if sid == "" {
			continue
		}
		guardValue, err := strconv.Atoi(strings.TrimSpace(line[equals+1:]))
		if err != nil {
			continue
		}
```

**Why it is wrong.** Four silent `continue`s. A user who types
`dragon_utopia = 30 000` or `dragon_utopia:30000` gets a template that generates
successfully, reports no problem, and simply ignores the override. The generator
already has a warnings channel — `GenerateTemplate` returns `dto.Warnings` and
[stateGeneration.go](../app/gui/drivers/stateGeneration.go#L113-L116) surfaces
the count — so the reporting mechanism exists and is unused here. The SID itself
is *deliberately* not validated (free-text by design), but a syntactically
unparseable line is a user error, not a design choice.

**Fix.** Change the signature to
`CreateValueOverrides(configuration config.GeneratorConfig) ([]entities.ValueOverride, []string)`
returning one warning per rejected line (`"line 4: 'foo:30000' is not
sid=value"`), and thread it into the existing warnings slice built by the
template generator. Keep the blank-line skip silent.

**Tests to add.** `test/unit/internal/services/template_generator/providers/gameRulesProvider/createValueOverrides_test.go`
(extend): one test per rejection class asserting the warning text, plus one
asserting blank lines produce no warning.

---

### 1.8 🟡 The chosen output directory is never saved or restored

**Evidence.** The output directory lives only as a Gio widget on the driver
([state.go](../app/gui/drivers/state.go#L34)):

```go
	outputPath   widget.Editor
```

It is seeded once from Steam detection at
[state.go](../app/gui/drivers/state.go#L62-L74), written by the picker at
[stateFiles.go](../app/gui/drivers/stateFiles.go#L75), and read at
[stateGeneration.go](../app/gui/drivers/stateGeneration.go#L71). It appears
**nowhere** in `EditorStateDto`
([editorStateDto.go](../internal/dtos/editorStateDto.go#L15-L98)) — I grepped
`outputPath` across `app/gui/**` and the only five hits are the five lines
above.

**Why it is wrong.** A user who deliberately points the output at, say, a
modding workspace rather than the detected Steam directory loses that choice on
every restart, silently reverting to the auto-detected path. It also fails the
"UI state survives a save/load round-trip" property: `Save` → quit → `Load`
restores everything *except* where the template will actually be written.

**Fix.** Two viable shapes — pick one with the owner:

- **(a)** Add `OutputDirectory string \`json:"outputDirectory,omitempty"\`` to
  `EditorStateDto`, load it in the preview panel's `LoadFromState` and write it
  in `SaveToState`. This makes the directory travel with the shared `.gen.json`,
  which may be undesirable across machines.
- **(b)** Persist it as machine-local preference next to the settings file
  (`os.UserConfigDir()` + `filepath.Join`, cross-platform per AGENTS.md §2.2),
  outside the `.gen.json`.

**(b) is recommended** — a path from another user's machine is worse than no
path. **If investigation shows** the owner wants portability, choose (a) and add
a fallback when the stored directory no longer exists.

**Tests to add.** For (b): `test/unit/internal/services/file_service/fileService/`
tests for the new preference read/write; plus a round-trip test in
`test/integration/editorState_integration_test.go`.

**Owner-decision flag:** ⚠ (a) vs (b) is a product decision.

---

### 1.9 🟡 A fatal window error is logged to a discard handler, then the process exits silently

**Evidence.** [program.go](../app/gui/program.go#L25-L47):

```go
func eventLoop(version string) {
	window := getAndConfigureWindow(version)
	windowLayout := editor.NewWindow(composition.InitializeGuiHandler())
	...
		case app.DestroyEvent:
			if event.Err != nil {
				slog.Error("Window destroyed with error", slog.String("error", event.Err.Error()))
				os.Exit(1)
			}
			os.Exit(0)
```

but `getAndConfigureWindow` — called on the line *above* — installs a discard
logger at [program.go](../app/gui/program.go#L56):

```go
	slog.SetDefault(slog.New(slog.DiscardHandler))
```

Logging is only re-enabled by the opt-in `-with-logging` flag
([program.go](../app/gui/program.go#L68-L71)).

**Why it is wrong.** When Gio fails to create or maintain the window (no GPU,
missing X/Wayland libraries on Linux, driver fault), the user sees the app
vanish with exit code 1 and **zero diagnostic output**, because the only error
message goes to `DiscardHandler`. This is precisely the failure a first-time
Linux user will hit (§9.7), and it is unreportable.

**Fix.** Write the fatal path to stderr unconditionally, independent of the
configured `slog` default:

```go
		case app.DestroyEvent:
			if event.Err != nil {
				fmt.Fprintln(os.Stderr, "Window destroyed with error:", event.Err)
				slog.Error("Window destroyed with error", slog.String("error", event.Err.Error()))
				os.Exit(1)
			}
```

Note `depguard` denies `log` in non-main files
([.golangci.yml](../.golangci.yml)) but not `fmt`, so this stays lint-clean.
Secondary: `os.Exit` inside the loop skips deferred cleanup — currently harmless
(nothing is deferred), but returning an error to `main.go` would be cleaner if
the owner wants a testable bootstrap.

**Tests.** No unit test is practical for the Gio event loop; record the gap in
[test_observations.md](test_observations.md) alongside the existing
`app/gui/program.go` entry.

---

## 2. Architecture & layering

### 2.1 🟠 The file-explorer dialog implements filesystem policy inside the GUI layer

**Evidence.** [fileExplorerDialog.go](../app/gui/dialogs/fileExplorerDialog.go)
makes **sixteen** direct `os`/`filepath` calls, including directory creation:

| Line | Call |
| --- | --- |
| [493](../app/gui/dialogs/fileExplorerDialog.go#L493) | `os.Stat(path)` (overwrite detection) |
| [524](../app/gui/dialogs/fileExplorerDialog.go#L524) | `os.ReadDir(dir)` |
| [603](../app/gui/dialogs/fileExplorerDialog.go#L603) | `filepath.Dir(this.currentDir)` |
| [645](../app/gui/dialogs/fileExplorerDialog.go#L645), [653](../app/gui/dialogs/fileExplorerDialog.go#L653) | `filepath.Base` / `filepath.Join` (name sanitisation) |
| [671-672](../app/gui/dialogs/fileExplorerDialog.go#L671-L672) | `filepath.Join` + `os.Mkdir(target, 0o750)` |
| [706](../app/gui/dialogs/fileExplorerDialog.go#L706), [720-729](../app/gui/dialogs/fileExplorerDialog.go#L720-L729), [738-744](../app/gui/dialogs/fileExplorerDialog.go#L738-L744) | root/home/cwd discovery |

AGENTS.md §4.5 is explicit: code under `app/gui/` "must contain only rendering
logic"; the depguard rule `no-services-from-app` enforces the *import* direction
but cannot see raw `os` usage.

**Why it is wrong.** Directory enumeration, hidden-file filtering, overwrite
detection, name sanitisation and directory creation are persistence policy, not
drawing. Consequences today: the file is 20.4% covered and 653 LOC because none
of that logic is reachable without a Gio context (§6 / test_observations.md);
`os.Mkdir(target, 0o750)` here disagrees with `folderPermission = 0o755` in
[fileService.go](../internal/services/file_service/fileService.go#L14-L18)
(§5.1); and a future TUI/web front-end (`app/tui/`, `app/web/` both exist as
placeholders) would have to reimplement all of it.

**Fix.** Extract a `internal/services/file_system/` service — e.g.
`directoryBrowserService.go` exposing `ListEntries(dir string, showHidden bool)
([]models.DirectoryEntry, error)`, `CreateDirectory(parent, name string) (string,
error)`, `ResolveStartDirectory(preferred string) string`,
`WouldOverwrite(path string) bool`, `SanitizeChildName(name string) (string,
bool)` — reached from the dialog through a new `handler_interfaces.IFileSystemHandler`
(the `app → internal/handlers/handler_interfaces` direction is already the
established, depguard-approved seam: nine such imports exist today and all are
clean). The dialog keeps only widget composition and event handling.

Do this **after** §1.1, so the new service and the atomic writer share one
permission constant.

**Tests to add.** A full mirrored folder
`test/unit/internal/services/file_system/directoryBrowserService/` with one file
per public method — `listEntries_test.go`, `createDirectory_test.go`,
`resolveStartDirectory_test.go`, `wouldOverwrite_test.go`,
`sanitizeChildName_test.go` — using `t.TempDir()` so they stay cross-platform.
This converts the largest untested file in the repo into ordinary unit-testable
logic and directly retires part of the deferred item in
[test_observations.md](test_observations.md).

---

### 2.2 🟠 Generation and manual-edit policy lives in the GUI driver

**Evidence.** [stateGeneration.go](../app/gui/drivers/stateGeneration.go#L30-L65)
implements the regeneration state machine (immediate vs. debounced vs. waiting),
and [stateGeneration.go](../app/gui/drivers/stateGeneration.go#L88-L120) decides
manual-edit reapplication ordering:

```go
	// The reapply decision and the castle-option diff both compare against the
	// state of the LAST generation, so they must be taken before
	// applyGeneratedTemplate snapshots the current state.
	reapplyManual := this.innerState.ShouldReapplyManualEdits()
	castleChanges := this.innerState.CastleSettingsChangedSinceGeneration()
	this.applyGeneratedTemplate(dto.Template)
	if reapplyManual && this.hasTemplateVariants() {
		this.reapplyManualEdits(castleChanges)
	} else if !reapplyManual {
		this.innerState.ClearManualEdits()
	}
```

with the supporting predicates in
[editorState.go](../app/gui/models/editorState.go#L53-L60) and
[editorState.go](../app/gui/models/editorState.go#L136-L150), and the castle
propagation in
[stateManualEdits.go](../app/gui/drivers/stateManualEdits.go).

**Why it is wrong.** This is business logic by AGENTS.md §4.5's own test — it
decides *what* to do, not *how to draw*. The comment quoted above documents a
subtle ordering invariant that is enforced nowhere but by comment. It is
untestable without a driver harness (`stateManualEdits.go` sits at 70%,
`stateFiles.go` at 30.4%, `state.go` at 79.4%), and it is the single largest
reason the backlog item "*anything inside app should only use entities, models,
handlers and commons, not services*" cannot be completed.

**Fix.** Move the decision logic behind the existing `IGuiHandler` seam: add a
`RegenerationDecisionService` under `internal/services/editor/` that takes
`(previous, current *dtos.EditorStateDto, now time.Time, pendingSince time.Time)`
and returns a small `dtos.RegenerationDecisionDto{Action, RedrawAt,
ReapplyManual, CastleChanges}`. The driver then becomes a thin dispatcher.
`EditorState` keeps only the three DTO pointers plus accessors.

This is a **large, multi-session refactor** — write a plan per AGENTS.md §2.4 /
§4.7 in `plans/extract-regeneration-policy.md` before starting, and do it in at
least two phases (predicates first, then the state machine).

**Tests to add.** `test/unit/internal/services/editor/regenerationDecisionService/`
with one file per public method; the debounce timing becomes deterministic
because `now` is a parameter. Existing driver tests shrink to dispatch checks.

**Owner-decision flag:** ⚠ large refactor, overlaps a backlog item — confirm
scope and sequencing before starting.

---

### 2.3 🟠 `NewMandatoryContentItemMapper` constructs its own collaborator, bypassing wire

**Evidence.** [mandatoryContentItemMapper.go](../internal/mappers/mandatoryContentItemMapper.go#L10-L16):

```go
type MandatoryContentItemMapper struct {
	contentRuleService *content_rules.ContentRuleService
}

func NewMandatoryContentItemMapper() *MandatoryContentItemMapper {
	return &MandatoryContentItemMapper{contentRuleService: content_rules.NewContentRuleService()}
}
```

**Why it is wrong.** Every other `internal/` constructor takes its dependencies
as parameters so `wire` can build the graph
([wire_gen.go](../internal/composition/wire_gen.go), providers declared in
[providerSets.go](../internal/composition/providerSets.go)). This one hard-wires
a concrete `*ContentRuleService`, so: the DI graph no longer describes the real
object graph; a second `ContentRuleService` instance exists alongside the wired
one; and the mapper cannot be unit-tested against a stub rule service.

**Fix.**

```go
func NewMandatoryContentItemMapper(contentRuleService *content_rules.ContentRuleService) *MandatoryContentItemMapper {
	return &MandatoryContentItemMapper{contentRuleService: contentRuleService}
}
```

Add the provider to the appropriate set in
[providerSets.go](../internal/composition/providerSets.go) and regenerate with
the *"Go: Generate wire injectors"* task (`wire gen ./internal/composition/...`,
AGENTS.md §4.6.2 — never pass `-tags=wireinject` to build/test). Fix the call
sites the compiler surfaces.

**Tests to add.** `test/unit/internal/mappers/mandatoryContentItemMapper/` gains
a constructor test asserting the injected instance is retained; existing mapper
tests can then inject a stub instead of exercising the real rule catalogue.

---

### 2.4 🟠 `NewTopologyBase` constructs its own collaborators

**Evidence.** [topologyBase.go](../internal/services/template_generator/providers/topology/base/topologyBase.go#L14-L29)
creates `zones.NewZoneLabelProvider()` and `newTopologyConnectionService(...)`
inside the constructor rather than receiving them.

**Why it is wrong.** Same class as §2.3: every topology service embedding
`TopologyBase` gets its own label provider and connection service, the wire
graph is incomplete, and the twelve topology services cannot be tested against a
deterministic label provider. `zoneLabelProvider.go` is 228 LOC of real logic.

**Fix.** Same shape as §2.3 — parameterise the constructor, register providers,
regenerate wire. Do §2.3 and §2.4 in one PR since both touch
[providerSets.go](../internal/composition/providerSets.go).

**Tests to add.** Extend
`test/unit/internal/services/template_generator/providers/topology/base/topologyBase/`
with a constructor test; the existing topology suites can then stub labelling.

---

### 2.5 🟡 `fileExplorerDialog.go` is a god object

**Evidence.** 653 LOC, 29 struct fields, 35 functions/methods — the largest
non-protected file in the repository. It is simultaneously an open dialog, a
save dialog, a folder picker, a read-only browser, a directory creator, and a
filesystem adapter.

**Fix.** Blocked on §2.1: once filesystem policy moves out, split the remainder
by mode into sibling files per AGENTS.md §4.1 —
`fileExplorerDialogEntries.go` (list rendering), `fileExplorerDialogToolbar.go`,
`fileExplorerDialogConfirm.go`. Do **not** split before §2.1, or the same
policy is merely scattered across more files.

**Tests.** Covered by §2.1's new service tests plus the already-tracked
integration scenario in [test_observations.md](test_observations.md).

---

### 2.6 🟡 `zoneEditorDialog` spreads ~58 fields across five state structs

**Evidence.** [zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go#L31-L61)
embeds five state structs (canvas, connection properties, zone properties, snap,
content) totalling roughly 58 fields and 42 methods across
`zoneEditorDialog.go` (479 LOC), `zoneEditorCanvas.go` (454),
`zoneEditorConnectionProps.go`, `zoneEditorSnap.go`, `zoneEditorZoneProps.go`.
All five sit at 0% unit coverage.

**Why it is worth recording.** The five-file split already satisfies AGENTS.md
§4.1, so this is **not** a style finding; the issue is that the *state* is
undifferentiated — geometry cache, selection, drag state, dirty flags, and
widget handles all live in one reachable blob, which is why the reset-semantics
bug (§0.2, owner-deferred) could exist undetected.

**Fix.** Low priority. When the owner next touches this dialog, extract the
pure-geometry and selection state into `internal/services/connection_editor/`
(where `zoneEditorService.go` already lives at 256 LOC) so the dialog holds only
widget handles. Not worth a standalone PR.

---

### 2.7 🟡 The gladiator-arena preview feature is half-landed: six dead embedded assets and two dead enum values

**Evidence.** Six arena sprites are compiled into every binary via
[assetProvider.go](../internal/services/asset_provider/assetProvider.go#L25-L26):

```go
	//go:embed assets/*.png
	assetFileSystem embed.FS
```

`internal/services/asset_provider/assets/` contains `gladiator_arena.png`,
`neutral_none_arena.png`, `neutral_low_arena.png`, `neutral_medium_arena.png`,
`neutral_high_arena.png`, `neutral_highest_arena.png`. None is referenced:
[assetProvider.go](../internal/services/asset_provider/assetProvider.go#L32-L37)
lists exactly ten neutral asset names and no arena variant, and the provider
exposes only `DrawBackground`, `DrawPlayerZone`, `DrawNeutralZone`.

Correspondingly, [previewConnection.go](../internal/models/preview/previewConnection.go#L21-L26)
declares:

```go
const (
	ConnectionTypeDirect ConnectionType = iota
	ConnectionTypePortal
	ConnectionTypeGladiatorArena
	ConnectionTypeProximity
)
```

but [previewLayoutService.go](../internal/services/preview_service/previewLayoutService.go#L194-L203)
only ever assigns `ConnectionTypeDirect` or `ConnectionTypePortal`. A
repository-wide grep finds `preview.ConnectionTypeGladiatorArena` at its
declaration and **nowhere else**.

Meanwhile [docs/gladiator-arena-marker.md](../docs/gladiator-arena-marker.md)
documents crossed-swords markers for exactly this case (and points at a
package that does not exist — §9.5).

**Why it is wrong.** Users generating a Gladiator-Arena template see the arena
connection rendered identically to a plain direct connection, contradicting the
in-repo documentation. The six unused PNGs are dead payload in every shipped
binary, and the two unused enum values are dead code that `exhaustive` cannot
catch because the producing code uses an `if`, not a `switch`.

**Fix.** Either finish or remove — this is an owner call:

- **Finish:** add the six names to the neutral/arena asset tables, add a
  `DrawArenaZone`, and set `ConnectionTypeGladiatorArena` /
  `ConnectionTypeProximity` in
  [previewLayoutService.go](../internal/services/preview_service/previewLayoutService.go#L194-L203)
  from `registry.GetConnectionTypeValues().GladiatorArena` / `.Proximity`
  (`registry` is read-only — read the constants, do not edit them). Then render
  the marker in `previewGeneratorService.go`.
- **Remove:** delete the six PNGs and the two enum constants.

**Tests to add (finish path).**
`test/unit/internal/services/preview_service/previewLayoutService/buildPreviewLayout_test.go`
gains `TestWhenConnectionTypeIsGladiatorArena_MarksPreviewConnectionAsArena`,
mirroring the existing Portal test at line 259; plus an asset-provider test
asserting every embedded PNG is decoded (which is the invariant that would have
caught this).

**Owner-decision flag:** ⚠ finish vs. remove is a product decision.

---

## 3. Duplicate code

### 3.1 🟠 The standard guard-weekly-increment is hard-coded at fifteen production call sites

**Evidence.** The value is already modelled —
[guardWeeklyIncrement.go](../internal/common/common_connections/guardWeeklyIncrement.go#L13):

```go
		Standard: 0.15,
```

yet fifteen production sites write the literal instead:

| File | Lines |
| --- | --- |
| [topologyConnectionService.go](../internal/services/template_generator/providers/topology/base/topologyConnectionService.go#L64) | 64, 116, 193 |
| [chainTopology.go](../internal/services/template_generator/providers/topology/chainTopology.go#L134) | 134 |
| [geometricHubTopology.go](../internal/services/template_generator/providers/topology/geometricHubTopology.go#L138) | 138, 160 |
| [hubTopology.go](../internal/services/template_generator/providers/topology/hubTopology.go#L128) | 128, 138 |
| [positionedTopologyBuilder.go](../internal/services/template_generator/providers/topology/positionedTopologyBuilder.go#L152) | 152 |
| [ringTopology.go](../internal/services/template_generator/providers/topology/ringTopology.go#L131) | 131 |
| [webTopology.go](../internal/services/template_generator/providers/topology/webTopology.go#L175) | 175, 196 |
| [balancedClusterService.go](../internal/services/template_generator/providers/topology/tournament_variant/balancedClusterService.go#L332) | 332 |
| [chainClusterService.go](../internal/services/template_generator/providers/topology/tournament_variant/chainClusterService.go#L98) | 98 |
| [hubClusterService.go](../internal/services/template_generator/providers/topology/tournament_variant/hubClusterService.go#L105) | 105 |
| [ringClusterService.go](../internal/services/template_generator/providers/topology/tournament_variant/ringClusterService.go#L144) | 144 |
| [castleFactory.go](../internal/services/zones/castleFactory.go#L48) | 48 |

all as `WithGuardWeeklyIncrement(0.15).`

**Why it is wrong.** Changing the standard increment — a plausible balance tweak
— requires finding fifteen unmarked float literals across twelve files, with no
compiler help and no single place that says "this is the standard". `goconst`
does not fire because it only tracks strings.

**Fix.** Replace every site with
`WithGuardWeeklyIncrement(common_connections.GetGuardWeeklyIncrements().Standard)`.
Purely mechanical; verify with a repository-wide grep for `0.15` afterwards
(only the definition and test fixtures should remain).

**Tests.** No new behaviour; the existing suites already pin `0.15` in the
expected values (e.g.
`test/unit/internal/services/template_generator/providers/topology/base/topologyBase/createMissingConnections_test.go`)
and must continue to pass unchanged — that is the verification.

---

### 3.2 🟡 The four tournament cluster services duplicate zone construction verbatim

**Evidence.** [ringClusterService.go](../internal/services/template_generator/providers/topology/tournament_variant/ringClusterService.go#L97-L108):

```go
		if index == 0 {
			zones = append(zones, this.CreateSpawnZone(
				label, fmt.Sprintf("Player%d", playerIndex+1), myConns,
				configuration.ZoneConfiguration.PlayerZoneCastles, configuration.MatchPlayerCastleFactions,
				configuration.ZoneConfiguration.PlayerZoneSize, tuning.RemoteFootholdCount,
				configuration.GenerateRoads, tuning))
		} else {
			zones = append(zones, this.CreateNeutralZone(
				linq.FromSlice(allNeutralZonePlans).
					FirstOrDefault(func(x neutral_zone.Plan) bool { return x.Label == label }),
				myConns, configuration.ZoneConfiguration.NeutralZoneSize,
				tuning.RemoteFootholdCount, configuration.GenerateRoads, tuning, false))
		}
```

and [hubClusterService.go](../internal/services/template_generator/providers/topology/tournament_variant/hubClusterService.go#L72-L86)
— byte-identical apart from the connection-name variable. The same block recurs
at [chainClusterService.go](../internal/services/template_generator/providers/topology/tournament_variant/chainClusterService.go#L67)
and [balancedClusterService.go](../internal/services/template_generator/providers/topology/tournament_variant/balancedClusterService.go#L279).

**Why it is wrong.** Nine and eight positional arguments respectively, repeated
four times. A new parameter on `CreateSpawnZone` means four identical edits, and
an argument transposed in one copy is silently wrong (Go cannot catch it — the
neighbouring `float64` and `int` parameters differ, but `PlayerZoneSize` and a
future size parameter would not).

**Fix.** Add one shared method on the common cluster base (the type that already
provides `CreateSpawnZone`/`CreateNeutralZone`):

```go
func (this *ClusterBase) createClusterZone(
	configuration config.GeneratorConfig, label string, connectionNames []string,
	playerIndex int, isSpawn bool, tuning models.GenerationTuning,
	allNeutralZonePlans neutral_zone.Plans) entities.Zone
```

and replace the four blocks with a call. Locate the base type first — it is
whatever `this.CreateSpawnZone` resolves to from all four services.

**Tests.** The four topology suites under
`test/unit/internal/services/template_generator/providers/topology/tournament_variant/`
must pass unchanged; add
`test/unit/.../clusterBase/createClusterZone_test.go` covering the spawn and
neutral branches.

---

### 3.3 🟡 The spell-label helper exists twice, verbatim, under two names

**Evidence.** [bonusesPanel.go](../app/gui/panels/bonusesPanel.go#L386-L397):

```go
// bannedSpellLabel returns the display name and school label for a banned spell.
func bannedSpellLabel(sid string) (name, school string) {
	if spell, ok := constants.FindSpell(sid); ok {
		label := constants.SpellSchoolDisplayNames[spell.School]
		if label == "" {
			label = spell.School
		}
		return spell.Name, label
	}

	return constants.SidToDisplayName(sid), "Spell"
}
```

[bonusPickerDialog.go](../app/gui/dialogs/bonusPickerDialog.go#L448-L459):

```go
// spellNameAndSchool resolves a spell SID to its display name and school
// label, with a sentence-case fallback for unknown SIDs.
func spellNameAndSchool(sid string) (name, school string) {
	if spell, ok := constants.FindSpell(sid); ok {
		label := constants.SpellSchoolDisplayNames[spell.School]
		if label == "" {
			label = spell.School
		}
		return spell.Name, label
	}
	return constants.SidToDisplayName(sid), "Spell"
}
```

`dupl` does not flag it because the body is below the token threshold.

**Fix.** Move one copy to `app/gui/constants/` beside `FindSpell` — e.g. a new
`spellLabel.go` exporting `SpellNameAndSchool(sid string) (name, school string)`
— and delete both private copies. `app/gui/constants` is *not* a protected
directory (only `internal/registry` is), so this is safe.

**Tests to add.** `test/unit/app/gui/constants/spellLabel/spellNameAndSchool_test.go`
— known SID returns name + display school; known SID with an unmapped school
falls back to the raw school; unknown SID returns sentence-cased name and
`"Spell"`. This also lifts `app/gui/constants/spells.go` off 39.3% coverage.

---

### 3.4 🟡 `buttonWidget.go` duplicates its render body (the repository's only `dupl` findings)

**Evidence.** Both current `dupl` issues:

```
app/gui/widgets/buttonWidget.go:33: 33-47 lines are duplicate of `app/gui/widgets/buttonWidget.go:67-81`
app/gui/widgets/buttonWidget.go:67: 67-81 lines are duplicate of `app/gui/widgets/buttonWidget.go:33-47`
```

[buttonWidget.go](../app/gui/widgets/buttonWidget.go#L33-L47) (inside
`NewButtonWidget`) and
[buttonWidget.go](../app/gui/widgets/buttonWidget.go#L67-L81) (inside
`NewToggleButtonWidget`) share an identical `material.Clickable(...)` body; the
functions differ only in how `textColor` / `backgroundColor` / `borderColor` are
derived. `NewSegmentButtonWidget` at
[buttonWidget.go](../app/gui/widgets/buttonWidget.go#L86-L90) begins the same
pattern a third time.

**Fix.** Extract a private
`renderButton(theme *material.Theme, label string, button *widget.Clickable, colors buttonColors) layout.Widget`
in the same file (it is the file's own primary struct-free helper, so no new
file is required per AGENTS.md §4.1) and have all three constructors compute
their `buttonColors` then delegate.

**Tests.** `app/gui/widgets/buttonWidget.go` is 0% and registered in
[test_observations.md](test_observations.md) as Gio-bound; do not add unit
tests. Verification is the existing GUI snapshot suite —
run *"Go: Run UI Integration tests headlessly"* and confirm no snapshot diff,
which proves the extraction is pixel-neutral.

---

## 4. Performance

### 4.1 🟠 The live preview rebuilds the full topology layout on every frame — measured 2.1 ms at default settings

**Evidence.** [previewPanel.go](../app/gui/panels/previewPanel.go#L149-L172) —
the call is inside the returned per-frame widget closure, with no cache and no
revision key:

```go
func (this *PreviewPanel) getPreviewCanvasWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		...
		response, err := this.previewHandler.BuildPreviewLayout(dtos.PreviewLayoutRequestDto{
			Template:   template,
			Topology:   this.state.GetStateData().Topology,
			CanvasSide: float64(canvasSize.X),
		})
```

For scatter topologies that reaches
[layoutScatter.go](../internal/services/preview_service/layoutScatter.go#L128-L137):

```go
func relaxPasses(px, py []float64, adj [][]int, zoneRadius float64) {
	...
	for range 500 {
		pushed := pushApartPass(px, py, minDist)
		nudged := nudgeOffEdgesPass(px, py, adj, edgeClear)
```

with `pushApartPass` being O(zones²) per pass.

**Measured this session** (scratch benchmark, since deleted; Intel Core Ultra 7
165H, `-benchtime=200x`, canvas side 600):

| Case | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| Random, 2 players, 0 neutrals (**app default**) | 1 284 | 832 | 14 |
| Random, 4 players, 8 neutrals | 561 107 | 13 456 | 203 |
| **Random, 8 players, 16 neutrals** | **2 090 102** | 40 378 | 391 |
| Ring, 8 players, 16 neutrals | 15 303 | 17 256 | 284 |
| Circles, 8 players, 16 neutrals | 36 654 | 37 000 | 338 |

**Why it is wrong.** `TopologyRandom` is the shipped default
([editorStateDto.go](../internal/dtos/editorStateDto.go#L118)). At 8 players /
16 neutrals the preview burns **2.1 ms and 391 allocations every single frame**
— about 13% of a 60 fps frame budget — recomputing a deterministic result from
unchanged inputs. This runs on idle repaints, pointer motion, status-line
changes and while a modal dialog covers the preview entirely. Doubling the zone
count quadruples the cost (0.56 ms → 2.09 ms), so larger templates degrade
super-linearly.

**Fix.** Cache `preview.Layout` in `PreviewPanel` keyed on
`(templateRevision, topology, canvasSide)`:

1. Add a monotonically-increasing revision counter to `drivers.State`, bumped in
   `applyGeneratedTemplate` ([stateGeneration.go](../app/gui/drivers/stateGeneration.go#L124-L128)),
   `clearGeneratedState`, and wherever manual edits install a new variant.
   Expose `GetTemplateRevision() uint64`. Prefer this over hashing the template.
2. In `PreviewPanel`, store `cachedKey` + `cachedLayout`; rebuild only on key
   change.
3. Leave `PreviewLayoutService` stateless and unchanged.

**Tests to add.**

- Reinstate the benchmark permanently as
  `test/performance/preview_layout_test.go` (with `//go:build integration_test`)
  so the before/after numbers are reproducible and future regressions visible.
- Cache-correctness: `test/unit/app/gui/panels/previewPanel/previewLayoutCache_test.go`
  against a stub `handler_interfaces.IPreviewHandler` counting calls — assert one
  build across N identical keys, and exactly one rebuild per key change (three
  separate tests: revision change, topology change, canvas-side change).
  **If investigation shows** the panel cannot be constructed without a Gio
  shaper, extract the key/cache into a small `previewLayoutCache.go` struct in
  `app/gui/panels/` and unit-test that struct directly.

---

## 5. Readability & maintainability

### 5.1 ✅ FIXED 🟡 Directory-creation permissions disagree between the two layers

**Fixed 2026-08-05** as part of §1.1, using this item's fallback rather than
waiting for §2.1. The two permission values moved into a new
[internal/constants/filePermissions.go](../internal/constants/filePermissions.go)
(`FolderPermission = 0o755`, `FileReadWritePermission = 0o644`) — the package
AGENTS.md §4.4 already prescribes for constants but which did not exist yet.
`internal/repositories` and
[fileExplorerDialog.go](../app/gui/dialogs/fileExplorerDialog.go) both consume it,
so the in-app explorer now creates `0o755` directories like every other path.
`internal/constants` had to be added to the allow-list in
`test/unit/architecture/dependency/dependency_test.go` so `app/*` may import it.

**Original finding follows.**

**Evidence.** [fileService.go](../internal/services/file_service/fileService.go#L14-L18):

```go
const (
	folderPermission        = 0o755
	fileReadWritePermission = 0o644
```

vs. [fileExplorerDialog.go](../app/gui/dialogs/fileExplorerDialog.go#L672):

```go
	if err := os.Mkdir(target, 0o750); err != nil {
```

**Why it is wrong.** A directory the user creates through the in-app explorer is
`0o750`; the same directory created implicitly by `SaveTemplate`'s `os.MkdirAll`
is `0o755`. On Linux that is a visible behavioural difference (group/other read
access). Neither value is wrong; having both, undocumented, is.

**Fix.** Fold into §2.1: the extracted filesystem service owns one
`folderPermission` constant used by both paths. If §2.1 is deferred, at minimum
export the constant from `internal/services/file_service` and use it in the
dialog.

---

### 5.2 ✅ FIXED 🟡 Commented-out code left in the state machine

**Fixed 2026-08-05.** The owner chose to **activate** the line rather than
delete it, so `SnapshotCurrentState` now really does clear `next`, with a
one-line comment stating why (a fresh snapshot supersedes any debounced state
still waiting to be applied). Behaviourally this is a no-op today — both
production call sites are in
[stateGeneration.go](../app/gui/drivers/stateGeneration.go), and on every path
that reaches them `next` is either already `nil` (`AutoRegenerate` clears it
first) or would be cleared by the next frame's
`ResetNextStateIfStateWasNotChanged` (manual `Generate()` button) — but the
invariant is now explicit instead of accidental. The full unit suite passes
unchanged.


**Evidence.** [editorState.go](../app/gui/models/editorState.go#L41-L45):

```go
func (this *EditorState) SnapshotCurrentState() {
	previousState := *this.current
	this.previous = &previousState
	// this.next = nil
}
```

**Why it is wrong.** In a hand-rolled three-pointer state machine whose ordering
invariants are already subtle (§1.3, §1.4, §2.2), a commented-out mutation of
one of the three pointers is genuinely ambiguous: a reader cannot tell whether
clearing `next` here is a pending fix or a rejected one. `godox` cannot see it
because it carries no TODO marker.

**Fix.** Delete the line and, if the decision matters, replace it with a
sentence explaining why `next` deliberately survives a snapshot. Verify against
[stateGeneration.go](../app/gui/drivers/stateGeneration.go#L30-L65), which is the
only consumer of the debounce pointer.

---

### 5.3 🟡 Nine-argument positional constructor calls in the cluster/topology services

**Evidence.** `CreateSpawnZone` is called with nine positional arguments at four
sites (see §3.2), `CreateNeutralZone` with eight; both interleave `string`,
`[]string`, `int`, `bool`, `float64`, `int`, `bool` and a struct.

**Why it is wrong.** Two adjacent `bool` parameters (`MatchPlayerCastleFactions`,
`GenerateRoads`) and two adjacent numeric parameters can be transposed without a
compile error. `funlen`/`gocognit` do not catch parameter-count problems and
`golangci-lint`'s argument-count linters are not enabled.

**Fix.** Introduce a request struct (`models.SpawnZoneRequest` /
`models.NeutralZoneRequest`) in `internal/models/`. Best executed together with
§3.2, since the shared helper already has to name every argument once.

**Tests.** Existing topology suites act as the regression net; no new tests
required beyond §3.2's.

---

### 5.4 ⚪ `.gitignore` blanket-ignores every top-level `.txt`

**Evidence.** [.gitignore](../.gitignore) contains, alongside the explicit
`coverage.txt`:

```
/*.txt
```

**Why it is recorded.** Any future top-level `.txt` (a `NOTICE.txt`, a third-party
licence bundle, a `requirements.txt` for a helper script) would be silently
untracked with no error. `coverage.txt` is already listed by name one line
earlier, so the wildcard adds no needed coverage.

**Fix (optional).** Delete the `/*.txt` line. Verify with
`git status --short` that nothing new appears (it will not — the tree is clean
and no other top-level `.txt` exists).

---

## 6. Testing

### 6.1 ❌ WILL NOT FIX 🟠 One integration test still bypasses the `integration_test` gate

**Rejected 2026-08-05 — do not re-attempt.** This finding rests on a version of
AGENTS.md §4.6.1 that no longer exists. The rule was rewritten (2026-08-05,
after this review was authored) and now reads: *"A file gets
`//go:build integration_test` **if and only if** it (or another file it shares a
package with) references an accessor declared in a `*_testexports.go`
implementation file. That is the whole rule."* — followed by *"The tag is NOT a
label for 'this is an integration/performance test.'"* and an explicit
instruction not to blanket-apply it to everything under
[test/integration/](../test/integration/).

[rmgTemplateModel_test.go](../test/integration/rmgTemplateModel_test.go)
references **no** test-only export (verified by grepping every accessor declared
in `window_testexports.go` / `state_testexports.go`); it only decodes
`.rmg.json` files from `data/` through production APIs. Adding the tag would
therefore violate the current AGENTS.md, and would also hide the test from the
default `go test ./test/...` run for no benefit. The file is correctly untagged.

**Residual, unaddressed.** The secondary observation stands: this file makes the
default suite decode all bundled example templates, so the "fast" run is slower
than it looks. That is a runtime concern to solve with `testing.Short()` or a
sampled corpus if it ever becomes painful — not with a build tag.


**Evidence.** First non-blank line of each file in the two gated directories:

```
test/integration/editorState_integration_test.go          //go:build integration_test
test/integration/main_integration_test.go                 //go:build integration_test
test/integration/manualCastleReapply_integration_test.go  //go:build integration_test
test/integration/rmgTemplateModel_test.go                 package integration_test   <-- VIOLATION
test/integration/stateExit_integration_test.go            //go:build integration_test
test/integration/window_render_integration_test.go        //go:build integration_test
test/integration/gui/*.go                                 //go:build integration_test && gui   (4 files)
test/performance/*.go                                     //go:build integration_test          (2 files)
```

`go list` confirms `test/integration` reports `untagged_xtest=1, ignored=5`.

**Why it is wrong.** AGENTS.md §4.6.1 states every file in these two directories
carries the tag, and that a plain `go test ./...` must skip them entirely.
Because of this one file, the default suite loads and decodes all 57 bundled
`.rmg.json` templates — so local and CI timings for the "fast" suite are wrong,
and a future `integration_test` accessor added to this file would break
`go build`/`go vet` in the untagged configuration.

**Fix.** Add the two lines before `package integration_test`:

```go
//go:build integration_test

package integration_test
```

No production change. Verify with
`go list -f '{{.Dir}} {{len .XTestGoFiles}}' ./test/integration/...` before and
after.

---

### 6.2 🟠 The entire `internal/handlers` package has no unit tests

**Evidence.** `internal/handlers/` contains `zoneEditorHandler.go` (144 LOC),
`contentRuleHandler.go` (124), `templateHandler.go` (123),
`stateHandler.go` (89), `previewHandler.go` (31), plus `guiHandler.go`. Under
`test/unit/internal/handlers/` only a `guiHandler/` folder exists — the other
five files have no mirrored folder, contrary to AGENTS.md §4.6 ("Code that is
exercised indirectly by other tests still requires its **own** test folder").

**Why it is wrong.** This is the entire boundary between the GUI and the
services — every load, save, generate, preview and zone-edit request crosses it.
It contains real logic, not pass-throughs: `stateHandler.LoadState` trims and
rejects empty paths, `ValidateEditorState` applies fixes conditionally, and
`normalizeInactiveNeutralCounts` ([stateHandler.go](../internal/handlers/stateHandler.go#L76-L89))
**silently zeroes eight user-supplied fields** depending on `AdvancedMode`
without emitting a warning — behaviour that no test asserts and no user is told
about. `previewHandler.BuildPreviewLayout` synthesises a template from loose
zones/connections ([previewHandler.go](../internal/handlers/previewHandler.go#L20-L31))
with an untested nil/empty matrix.

**Fix.** Add mirrored folders with one file per public method:

```
test/unit/internal/handlers/stateHandler/loadState_test.go
test/unit/internal/handlers/stateHandler/saveState_test.go
test/unit/internal/handlers/stateHandler/validateEditorState_test.go
test/unit/internal/handlers/previewHandler/buildPreviewLayout_test.go
test/unit/internal/handlers/templateHandler/...
test/unit/internal/handlers/contentRuleHandler/...
test/unit/internal/handlers/zoneEditorHandler/...
```

Package name `<fileName>_test`; `t.Parallel()` everywhere; `testify` +
`gofakeit` only. Start with `stateHandler` (highest risk, feeds §1.5) and
`previewHandler` (smallest, feeds §4.1's stub).

Cover explicitly: empty/whitespace path → `ErrNoOutputPath`; nil state →
`ErrNothingToSave`; `fixIssues=false` returns warnings but an unmodified state;
`AdvancedMode` true/false each zeroing the correct field set;
`BuildPreviewLayout` with nil template + nil zones + nil connections.

---

### 6.3 🟡 Three internal-service unit tests import the GUI layer

**Evidence.** All three import `app/gui/constants` at line 6:

- [createRuleFromSavedRule_test.go](../test/unit/internal/services/content_rules/contentRuleService/createRuleFromSavedRule_test.go#L1-L10)
- [getVariantForContentById_test.go](../test/unit/internal/services/content_rules/variantMappingCatalog/getVariantForContentById_test.go#L1-L10)
- [getVariantsForContent_test.go](../test/unit/internal/services/content_rules/variantMappingCatalog/getVariantsForContent_test.go#L1-L10)

(The prior review reported this against `variantMappingManager/` and
`contentRuleManager/`; those folders were renamed, the violation moved with
them.) The production depguard rule `no-ui-from-internal` exempts tests via
`!$test`, so nothing catches it.

**Why it is wrong.** Business-logic tests are coupled to a UI catalogue: a
presentational change in `app/gui/constants` can break `content_rules` tests
that have nothing to do with the UI, and it keeps the wrong dependency direction
alive as a legitimate-looking pattern.

**Fix.** Replace the imports with local `const` SID literals, or construct
`models.SidMapping` values directly in each test's Arrange section. Do not add
production exports for tests, and do not touch the protected
`internal/registry`. Then close the hole per §6.5.

---

### 6.4 🟡 Two pure, non-Gio catalogues remain at 0% coverage

**Evidence.** From the deduplicated profile:
`app/gui/constants/bannableItems.go` — **110 statements, 0.0%**;
`app/gui/constants/valueOverrideSids.go` — **40 statements, 0.0%**.
(For context `app/gui/constants/spells.go` is 39.3%.) Unlike the widgets and
dialogs registered in [test_observations.md](test_observations.md), these are
plain slice construction, filtering and sorting with no `layout.Context`.

**Fix / tests.**

- `test/unit/app/gui/constants/bannableItems/getBannableItems_test.go` — no empty
  ID/name/category; SIDs unique; one representative entry from each of the
  ~30 groups present; `FindBannableItem` returns `(item, true)` for a known SID
  and `(_, false)` otherwise.
- `test/unit/app/gui/constants/valueOverrideSids/getValueOverrideSidsWithExclusions_test.go`
  — sorted output; excluded SIDs absent; no duplicates; no empty SIDs; the
  caller's input slice is not mutated.

Assert invariants and a small representative set — do not snapshot the whole
protected registry.

---

### 6.5 🟡 No executable enforcement of the repository's test-layout rules

**Evidence.** §6.1 and §6.3 are both *written* rules (AGENTS.md §4.6, §4.6.1)
that were violated, compiled cleanly, and passed CI. `.golangci.yml`'s depguard
rules all carry `!$test`, so no linter inspects test files' imports.

**Fix.** Add a depguard scope rather than a unit-test folder — the prior
review's suggested `test/unit/repository/testLayout/` location itself violates
the mirror convention (its own verification pass flagged this). In
[.golangci.yml](../.golangci.yml) add:

```yaml
      test-unit-internal-no-gui:
        files:
          - 'test/unit/internal/**'
        deny:
          - pkg: github.com/Tariomka/hommoe_custom_templates/app
            desc: internal-layer tests must not depend on the GUI layer
```

For the build-tag rule, add a small CI step in
[pr-validation.yml](../.github/workflows/pr-validation.yml) that greps the first
non-blank line of every `test/integration/**/*_test.go` and
`test/performance/*_test.go` and fails if it is not the constraint. Keep it in
CI, not in Go, so it does not need a mirror folder.

**Verification.** After adding the depguard scope, `golangci-lint-v2 run ./...`
must report exactly the three §6.3 files until they are fixed, then zero.

---

### 6.6 ⚪ The integration suite depends on a hand-maintained golden template

**Evidence.** [defaultTemplate.go](../test/test_helpers/defaultTemplate.go#L12-L31)
loads `test/test_helpers/defaultTemplate.json` and applies imperative
post-processing; nine test files consume it.

**Why it is recorded, not flagged.** A committed golden is a legitimate choice
and the file is small enough to review in a diff. The risk is only that a
generator change requires a hand-edit that is easy to get subtly wrong. If the
owner ever finds themselves editing it by hand more than occasionally, add a
`-update` regeneration flag mirroring the existing GUI snapshot task
(*"Go: Update UI Integration tests snapshots"*). No action now.

---

## 7. CI/CD

### 7.1 🟠 Direct pushes to `master` skip tidy, lint, vulnerability, race and coverage gates

**Evidence.** [pr-validation.yml](../.github/workflows/pr-validation.yml#L4-L7)
triggers on both `push: master` and `pull_request: master`, but six of ten jobs
are PR-gated:

| Job | Line | Gate |
| --- | --- | --- |
| `check-go-mod` | [21](../.github/workflows/pr-validation.yml#L21) | PR only |
| `run-gci-lint` | [43](../.github/workflows/pr-validation.yml#L43) | PR only |
| `run-vulnerability-scan` | [66](../.github/workflows/pr-validation.yml#L66) | PR only |
| `run-race-tests` | [154](../.github/workflows/pr-validation.yml#L154) | PR only |
| `code_coverage` (enforces the 60.0% floor) | [172](../.github/workflows/pr-validation.yml#L172) | PR only |
| `run-gui-integration-tests` | [224](../.github/workflows/pr-validation.yml#L224) | PR only |

Only `check-build`, `check-windows`, `run-unit-tests` and
`run-integration-tests` run on a direct push.

**Why it is wrong.** §8.1 below is the concrete proof: a reachable dependency
vulnerability is present on the current branch, and the only job that would
catch it does not run on push. The same applies to lint regressions, data races
and coverage drops.

**Fix.** Prefer branch protection requiring a PR plus all checks. For defence in
depth, drop the `if:` from the four read-only gates (`check-go-mod`,
`run-gci-lint`, `run-vulnerability-scan`, `run-race-tests`) — none needs PR
context. Split `code_coverage` into a floor-enforcement step (runs always, no
special token) and a comment-publishing step (PR only, keeps
`pull-requests: write`). Rename the workflow from "PR Tests", which
misdescribes it.

**Verification.** `actionlint`, then push a scratch branch and confirm the job
matrix. **Owner-decision flag:** ⚠ confirm whether direct pushes to master are
intentionally permitted before changing the trigger shape.

---

### 7.2 🟠 Workflows have no top-level `permissions:`

**Evidence.** [pr-validation.yml](../.github/workflows/pr-validation.yml#L1-L16)
declares `name`, `on` and `env` but no `permissions`. The only narrowing is
inside `code_coverage` at
[pr-validation.yml](../.github/workflows/pr-validation.yml#L175). In
[release.yml](../.github/workflows/release.yml#L77) the publish job narrows to
`contents: write`, but the build job does not.

**Why it is wrong.** Every unnarrowed job inherits the repository/organisation
default token scope, which may be read-write. Nine jobs that only need to read
source (and one that runs a third-party linter action and a third-party coverage
action) receive whatever the org default grants, and org defaults change outside
this repo's control.

**Fix.** Add to both workflow files, immediately after `on:`:

```yaml
permissions:
  contents: read
```

Keep the existing per-job overrides (`actions: read` + `pull-requests: write`
for coverage; `contents: write` for release publishing). Re-run the composite
`setup-steps` action afterwards to confirm `go mod download` still succeeds.

---

### 7.3 🟡 `actions/setup-go` version drift between the workflows and the composite action

**Evidence.** [pr-validation.yml](../.github/workflows/pr-validation.yml#L30)
uses `actions/setup-go@v7`, while the composite action every other job calls —
[setup-steps/action.yml](../.github/workflows/setup-steps/action.yml#L14) — uses
`actions/setup-go@v6`.

**Why it is wrong.** Nine of ten jobs run on `@v6` and exactly one on `@v7`, so
toolchain-setup and caching behaviour differ between `check-build` and every
other job. Dependabot's `github-actions` ecosystem groups only *minor and patch*
updates ([dependabot.yml](../.github/dependabot.yml)), so the major bump was
applied by hand in one place and missed in the other — and will stay split.

**Fix.** Bump the composite to `actions/setup-go@v7`, or pin both to `@v6`.
Verify the cache still restores by checking the "Setup Go" step output on the
next run.

---

### 7.4 🟡 The `tools/` module is never built, tested, linted or tidy-checked in CI

**Evidence.** Every CI Go command targets the root module —
`go build ./...`, `go vet -tags=integration_test ./...`,
`go test ./test/...`, and `check-go-mod` — and none runs inside `tools/`.
[tools/go.mod](../tools/go.mod) is a separate module:

```
module github.com/Tariomka/hommoe_custom_templates/tools

go 1.26.3
```

pinning `wire`, `golangci-lint` and `gcov2lcov` through `tool` directives, and
declaring a **different Go version** from the root module's `1.26.5` and from
CI's `GO_VERSION: 1.26.5`.

**Why it is wrong.** The module that pins the repository's own linter and code
generator is unverified. If `tools/go.sum` drifts, the failure surfaces only
when a developer runs `wire gen` locally. (Dependabot's exclusion of `tools/` is
a recorded owner decision — §0.2 — so this finding is about CI verification, not
dependency updates.)

**Fix.** Add one step to the existing `check-go-mod` job:

```yaml
      - name: Check tools module is tidy
        working-directory: tools
        run: go mod tidy -diff
```

Both modules currently pass this locally (verified: exit 0 for each). Separately,
decide whether `tools/go.mod` should track `1.26.5`; a lower `go` directive is
valid but the divergence is undocumented.

---

### 7.5 ⚪ CI lints with four linters disabled

**Evidence.** [pr-validation.yml](../.github/workflows/pr-validation.yml#L61):

```yaml
          args: --disable=godox,dupl,unparam,gochecknoglobals # these are the same linters marked as warning level in .golangci.yml
```

**Why it is recorded, not flagged.** This is deliberate and matches the
`severity: warning` entries in [.golangci.yml](../.golangci.yml). The
consequence worth knowing: all **42** locally-reported issues (§10) are in the
disabled set, so CI's lint job reports zero and the two `dupl` findings in §3.4
never gate a merge. No change recommended; recorded so the §10 numbers are not
mistaken for CI failures.

---

## 8. Security & dependencies

### 8.1 ✅ FIXED 🔴 A reachable vulnerability is present in `golang.org/x/text`

**Fixed 2026-08-05.** `golang.org/x/text` bumped to `v0.39.0` (and
`golang.org/x/sys` to `v0.46.0` as a transitive consequence of `go mod tidy`).
`govulncheck ./...` and `govulncheck -scan module` both report
*No vulnerabilities found*; `go build ./...` and `go test ./test/unit/...` pass;
`go mod tidy -diff` exits 0.

**Evidence.** `govulncheck` symbol scan, run this session against the working
tree:

```
Vulnerability #1: GO-2026-5970
    Infinite loop on invalid input in golang.org/x/text
  Module: golang.org/x/text
    Found in: golang.org/x/text@v0.38.0
    Fixed in: golang.org/x/text@v0.39.0
    Example traces found:
      #1: app/gui/program.go:21:10: gui.StartApplication calls app.Main, which eventually calls norm.Form.Bytes
      #2: app/gui/program.go:21:10: gui.StartApplication calls app.Main, which eventually calls norm.Form.IsNormalString
      #3: app/gui/program.go:21:10: gui.StartApplication calls app.Main, which eventually calls norm.Form.QuickSpan
      #4: app/gui/program.go:21:10: gui.StartApplication calls app.Main, which eventually calls norm.Form.String

Your code is affected by 1 vulnerability from 1 module.
```

The dependency is indirect —
[go.mod](../go.mod#L26): `golang.org/x/text v0.38.0 // indirect` — pulled in
through Gio's text handling, and the trace starts at the application's own entry
point.

**Why it is wrong.** Text normalisation is reachable from the running GUI (Gio
normalises text input and font data), so malformed input can hang the process in
an infinite loop — a denial of service on the user's own machine. It is
**reachable**, not merely present. The gate that would have caught it does not
run on this branch's pushes (§7.1).

**Fix.**

```powershell
go get golang.org/x/text@v0.39.0
go mod tidy
go build ./...
go test ./test/unit/... -count=1
```

Then re-run `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` and confirm
the symbol scan is clean.

**Verification / prevention.** Combine with §7.1 so `run-vulnerability-scan`
executes on pushes too. No Go test applies.

---

### 8.2 ✅ FIXED 🟠 A second known vulnerability is present but not currently called

**Fixed 2026-08-05.** `golang.org/x/net` bumped to `v0.56.0` in the same change
as §8.1.

**Evidence.** `govulncheck -scan module`:

```
Vulnerability #2: GO-2026-5942
    Parsing an invalid SVCB or HTTPS RR can panic in golang.org/x/net/dns/dnsmessage
  Module: golang.org/x/net
    Found in: golang.org/x/net@v0.55.0
    Fixed in: golang.org/x/net@v0.56.0
```

[go.mod](../go.mod#L25): `golang.org/x/net v0.55.0 // indirect`.

**Why it matters.** The symbol scan shows no current call path, so this is not
an active exploit surface. It should still be bumped in the same PR as §8.1 —
leaving a known-vulnerable module in `go.mod` means the next `govulncheck` run
after any new import can flip to "affected" with no warning.

**Fix.** `go get golang.org/x/net@v0.56.0; go mod tidy`, in the same PR as §8.1.

---

### 8.3 🟡 The vulnerability gate only runs on pull requests

**Evidence.** `run-vulnerability-scan` is gated at
[pr-validation.yml](../.github/workflows/pr-validation.yml#L66) with
`if: ${{ github.event_name == 'pull_request' }}`, and there is no scheduled run.

**Why it is wrong.** A vulnerability disclosed *after* a PR merges is never
detected, because nothing re-scans `master` on a schedule. GO-2026-5970 (§8.1)
is exactly this case — the code did not change, the advisory did.

**Fix.** Two independent changes:
1. Remove the PR gate (see §7.1) so pushes are scanned.
2. Add `schedule: - cron: '0 6 * * 1'` to the workflow triggers (or a small
   dedicated `security-scan.yml`) so `master` is re-scanned weekly regardless of
   activity. Pair with `permissions: contents: read` from §7.2.

---

## 9. Documentation & developer experience

### 9.1 🟠 The QUICKSTART programmatic example cannot compile — three separate errors

**Evidence.** [QUICKSTART.md](../QUICKSTART.md#L111-L133):

```go
    cfg.Topology = config.TopologyDefault // Ring

    template := template_generator.NewTemplateGenerator(cfg).Generate()

    out, err := services.WriteTemplate(".", template)
```

1. `config.TopologyDefault` does not exist. The Ring topology constant is
   `TopologyRing` with the *value* `"Default"`
   ([mapTopology.go](../internal/models/config/config_inner/mapTopology.go#L6-L7)).
2. `template_generator.NewTemplateGenerator` does not take one argument — it
   takes eight injected providers (see the wire providers in
   [providerSets.go](../internal/composition/providerSets.go)); the tests use
   the `test_helpers.NewTemplateGenerator(configuration)` wrapper.
3. `services.WriteTemplate` does not exist. The current API is
   `(*file_service.FileService).SaveTemplate(directory string, template *entities.RmgTemplate) (string, error)`
   ([fileService.go](../internal/services/file_service/fileService.go#L58)).

**Why it is wrong.** The section is titled "Programmatic Use" and is the only
non-GUI entry point documented. A user copying it gets three compile errors, and
nothing explains that the whole API sits behind `internal/` and is therefore
unusable from outside the module anyway.

**Fix.** Replace the snippet with a working one built on the composition root —
`composition.InitializeGuiHandler()` returns an `IGuiHandler` exposing
`GenerateTemplate` and `SaveTemplate`, which is the real supported path — and
add a sentence stating the packages are `internal/` and the example must live
inside this module.

**Prevention.** Make the example compile-checked: add
`test/unit/docs/programmaticUse/programmaticUse_test.go` containing the exact
code from the doc inside a test function. **If investigation shows** the mirror
convention forbids a `docs/` folder under `test/unit/`, put it in a top-level
`examples/programmatic/main.go` package instead and let `go build ./...` (already
in CI) catch drift.

**Owner-decision flag:** ⚠ whether external programmatic use is a supported
product surface (which would justify a public package) or the example should
simply be marked internal-only.

---

### 9.2 🟠 README and QUICKSTART describe a UI that does not exist

**Evidence.** Verified against source this session:

| Doc claim | Reality |
| --- | --- |
| [QUICKSTART.md](../QUICKSTART.md#L3-L5) "across four tabs"; [QUICKSTART.md](../QUICKSTART.md#L31) "the four configuration tabs"; §"3. Zone Content" tab | [window.go](../app/gui/editor/window.go#L36-L40) builds **three**: `General`, `Layout & Zones`, `Bonuses & Bans`. There is no `zoneContentPanel.go`; zone content is a dialog opened from Layout & Zones. |
| [README.md](../README.md#L80-L97) "four configuration tabs" incl. `Zone Content` | Same. |
| [QUICKSTART.md](../QUICKSTART.md#L30) toolbar `New, Open…, Save, Save As…` | [toolbar.go](../app/gui/editor/toolbar.go#L64-L74) renders `New`, **`Load`**, `Save`, `Save As`, **`Exit`**. The doc has the wrong label for one button and omits another. |
| [QUICKSTART.md](../QUICKSTART.md#L33-L35) preview panel has a `Refresh` button | [previewPanel.go](../app/gui/panels/previewPanel.go#L55-L75) has `Browse`, `Reveal`, `Generate`, `Save Template`. No `Refresh` exists anywhere. |
| [QUICKSTART.md](../QUICKSTART.md#L36-L39) "**Footer** (bottom): output folder picker…" | There is no footer region. Those controls are inside `PreviewPanel`. |
| [QUICKSTART.md](../QUICKSTART.md#L36) / [QUICKSTART.md](../QUICKSTART.md#L104) "`Generate Template`" button; "**Reveal** opens the output folder in your file explorer" | The button is labelled `Generate`; `Reveal` opens the **in-app** browser — [stateFiles.go](../app/gui/drivers/stateFiles.go#L81-L83) calls `dialogs.NewBrowseDialog`. |
| [QUICKSTART.md](../QUICKSTART.md#L79-L81) "`services.SaveSettingsFile` / `services.LoadSettingsFile`" | The methods are `(*file_service.FileService).SaveSettings` / `.LoadSettingsFile` ([fileService.go](../internal/services/file_service/fileService.go#L31-L53)). |
| [README.md](../README.md#L46) source tree | Lists `General / Layout / Zone Content / Bonuses & Bans / preview / footer`. Real panels: `generalPanel.go`, `layoutPanel.go`, `bonusesPanel.go`, `previewPanel.go`. |

**Why it is wrong.** A new user cannot map the guide onto the window; a
contributor searches for files that do not exist. This has now survived two
consecutive reviews.

**Fix.** One docs-only PR. Generate the tree section from the real directory
listing; rewrite the tab and control sections from
[window.go](../app/gui/editor/window.go),
[toolbar.go](../app/gui/editor/toolbar.go) and
[previewPanel.go](../app/gui/panels/previewPanel.go); correct the persistence
API names. Update repository memory afterwards (`/memories/repo/`) so the stale
"four tabs" note (§0.2) is corrected at the source.

**Tests.** Run every documented command in a clean checkout; compile the §9.1
example.

---

### 9.3 🟡 QUICKSTART states the wrong minimum Go version

**Evidence.** [QUICKSTART.md](../QUICKSTART.md#L10): "Requires Go **1.25.8+**."
[go.mod](../go.mod#L3) declares `go 1.26.5`, and CI pins `GO_VERSION: 1.26.5`
([pr-validation.yml](../.github/workflows/pr-validation.yml#L12)).

**Why it is wrong.** A user on Go 1.25.8 following the instructions gets a
toolchain error at `go run .`, not a clear "upgrade Go" message.

**Fix.** Change to `1.26.5+`. Add it to the docs checklist in the §9.2 PR.

---

### 9.4 🟡 QUICKSTART undercounts the topologies

**Evidence.** [QUICKSTART.md](../QUICKSTART.md#L56): "pick one of ten layouts".
[mapTopology.go](../internal/models/config/config_inner/mapTopology.go#L6-L18)
declares **eleven**: Ring, HubAndSpoke, Chain, SharedWeb, Random, Circles,
Square, Geometric, Cross, Fractal, GeometricHub. (README's own topology table is
correct with eleven — the two docs disagree with each other.)

**Fix.** Change "ten" to "eleven", or better, drop the count and point at the
README table so only one place needs maintenance.

---

### 9.5 🟡 `docs/gladiator-arena-marker.md` points at a package that does not exist

**Evidence.** [gladiator-arena-marker.md](../docs/gladiator-arena-marker.md#L111-L131)
locates the arena sprites under `internal/services/previewassets/`.
A file search for `internal/services/**/previewassets/**` returns **no files**.
The assets are at `internal/services/asset_provider/assets/`, embedded at
[assetProvider.go](../internal/services/asset_provider/assetProvider.go#L25-L26).

**Fix.** Correct the path. Do this together with §2.7, since the same document
describes the rendering behaviour that is not implemented — whichever way §2.7 is
resolved (finish or remove), this document must be updated to match.

---

### 9.6 🟡 AGENTS.md claims a single module and names VS Code tasks that do not exist

**Evidence.**

- [AGENTS.md](../AGENTS.md#L10): "**Language / Toolchain:** Go 1.26.5, single
  module `github.com/Tariomka/hommoe_custom_templates`." There are **two**
  modules: the root and [tools/go.mod](../tools/go.mod) (`go 1.26.3`). AGENTS.md
  §4.6.2 itself later refers to "the `tool` directive of `tools/go.mod`",
  contradicting §1.
- [AGENTS.md](../AGENTS.md) §4.6.1 tells agents to use the VS Code tasks
  *"go: test (default, no integration_test)"* and *"go: test
  integration+performance (integration_test)"*. Neither label exists in
  `.vscode/tasks.json`. The actual labels are `Go: Run Unit tests`,
  `Go: Run Integration tests`, `Go: Run UI Integration tests headlessly`,
  `Go: Run UI Integration tests in headed mode`,
  `Go: Update UI Integration tests snapshots`, `Go: Run Performance tests`,
  `Go: Run Performance tests with profiling`,
  `Go: Generate code coverage report`, `Go: Run Linter`,
  `Go: Get Linter Results`, `Go: Generate wire injectors`.

**Why it is wrong.** AGENTS.md is operational policy for every future agent
session. An agent told to run a nonexistent task either fails or improvises, and
"single module" causes agents to skip `tools/` entirely (§7.4).

**Fix.** Correct both statements in place. §1 becomes "two modules: root
(Go 1.26.5) and `tools/` (Go 1.26.3, tool dependencies only)". Replace the task
names in §4.6.1 and §7 with the real labels. Note the §7 quick-reference row for
the coverage report already uses the correct label, so only §4.6.1 is wrong
there.

---

### 9.7 ⚪ Linux build prerequisites are undocumented

**Evidence.** [setup-steps/action.yml](../.github/workflows/setup-steps/action.yml)
installs sixteen system packages before any Linux build (`libgles2-mesa-dev`,
`libegl1-mesa-dev`, `libffi-dev`, `libxkbcommon-dev`, `libxkbcommon-x11-dev`,
`libvulkan-dev`, `libwayland-dev`, `libx11-dev`, `libx11-xcb-dev`, `libxcb1-dev`,
`libxcursor-dev`, `libxfixes-dev`, `libxrandr-dev`, `libxinerama-dev`,
`libxi-dev`, `xorg-dev`). Neither [README.md](../README.md) nor
[QUICKSTART.md](../QUICKSTART.md) mentions any of them.

**Why it is recorded.** AGENTS.md §2.2 mandates Linux support, and the release
workflow ships a Linux binary with `CGO_ENABLED=1`. A Linux contributor running
`go run .` hits a cgo link error with no guidance — and, because of §1.9, may
get no message at all.

**Fix (optional, fold into §9.2).** Add a "Building on Linux" subsection to
QUICKSTART with the `apt-get install` line copied from the composite action.

---

## 10. Linter disposition — all 42 current issues

Command: `golangci-lint-v2 run ./... --issues-exit-code=0 --max-issues-per-linter=0 --max-same-issues=0`.
Total **42**. No `funlen`, `gocognit`, `cyclop`, `maintidx`, `goconst`, `godox`,
`unparam`, `gosec`, `errcheck`, `depguard`, `exhaustive` or formatter finding
exists. CI disables all four linters below (§7.5), so CI's lint job is green.

| Linter | Count | Location | Disposition |
| --- | ---: | --- | --- |
| `dupl` | 2 | [buttonWidget.go](../app/gui/widgets/buttonWidget.go#L33) 33-47 ↔ 67-81 | **Actionable — §3.4.** The repository's only duplication findings. |
| `gochecknoglobals` | 6 | [providerSets.go](../internal/composition/providerSets.go#L26) lines 26, 35, 48, 58, 65, 74 — `ZoneSet`, `GenerationSet`, `EditorSet`, `InfrastructureSet`, `HandlerSet`, `GuiHandlerSet` | **Unavoidable.** `wire.NewSet` values must be package-level. Candidate for a `nolint` with reason, or a path exclusion for `internal/composition/`. |
| `gochecknoglobals` | 16 | [common.go](../internal/services/template_generator/providers/common.go#L6) lines 6–21 — `buildingObjects`, `championSelectRules`, `gameModes`, `heroBuffBuildings`, `magicBuildings`, `nonContentObjects`, `randomUnitBanks`, `resourceObjects`, `ruleTypes`, `t1GuardedResourceBanks`, `t1StatsAndSkillsBuildings`, `t2StatsAndSkillsBuildings`, `unitBanks`, `visionBuildings`, `winConditionValues`, `zoneLayouts` | **Owner's responsibility** (historical §3.4). Immutable registry aliases. |
| `gochecknoglobals` | 6 | `app/gui/constants` — [contentIds.go](../app/gui/constants/contentIds.go#L132) `ContentIDs`; [gameModes.go](../app/gui/constants/gameModes.go#L5) `GameModeValues`, `GameModes`; [legend.go](../app/gui/constants/legend.go#L15) `LegendRows`; [spells.go](../app/gui/constants/spells.go#L22) `SpellSchoolDisplayNames`, `KnownSpells` | **Owner's responsibility.** UI display catalogues. |
| `gochecknoglobals` | 4 | [bonusPickerDialog.go](../app/gui/dialogs/bonusPickerDialog.go#L24) lines 24, 32, 46, 48 — `receiversFilters`, `bonusTypeOptions`, `bonusReceiverOptions`, `bonusResourceDefaults` | **Owner's responsibility.** Dialog catalogues. |
| `gochecknoglobals` | 4 | `internal/services/builders/variant_content` — [mainObjectBuilder.go](../internal/services/builders/variant_content/mainObjectBuilder.go#L9) `castleQualities`; [typedRefBuilder.go](../internal/services/builders/variant_content/typedRefBuilder.go#L9) `roadConnTypes`, `biomeTypes`; [zoneBuilder.go](../internal/services/builders/variant_content/zoneBuilder.go#L8) `layoutValues` | **Owner's responsibility.** Registry aliases. |
| `gochecknoglobals` | 2 | [assetProvider.go](../internal/services/asset_provider/assetProvider.go#L30) `loadAssetProvider`, `neutralAssetNames` | **Justified.** `loadAssetProvider` is a `sync.OnceValues` memo whose package-level lifetime is the point (see the in-code comment). Leave. |
| `gochecknoglobals` | 1 | [lookupSid.go](../app/gui/utils/lookupSid.go#L30) `allSidMappings` | **Owner's responsibility.** |
| `gochecknoglobals` | 1 | [previewGeneratorService.go](../internal/services/preview_service/previewGeneratorService.go#L27) `connectorLineColor` | **Owner's responsibility.** Immutable colour constant. |
| **Total** | **42** | | **2 actionable (§3.4); 40 owner-controlled or justified.** |

---

## 11. Verified non-issues

Checked this session and found correct — recorded so they are not re-examined:

**Layering and dependencies**

- No production file under `internal/**` imports `app/**`.
- No production file under `app/**` imports `internal/services/**` or
  `internal/handlers` (the concrete package). All nine `app → internal/handlers`
  imports target `handler_interfaces` and are depguard-approved:
  [ruleDialog.go](../app/gui/dialogs/ruleDialog.go#L16),
  [zoneContent.go](../app/gui/dialogs/zoneContent.go#L19),
  [zoneContentDialog.go](../app/gui/dialogs/zoneContentDialog.go#L13),
  [zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go#L21),
  [state.go](../app/gui/drivers/state.go#L17),
  [window.go](../app/gui/editor/window.go#L17),
  [editorState.go](../app/gui/models/editorState.go#L7),
  [layoutPanel.go](../app/gui/panels/layoutPanel.go#L13),
  [previewPanel.go](../app/gui/panels/previewPanel.go#L22).
- `internal/services/preview_service` uses the standard `image`/`color` packages,
  not Gio types — the service layer is genuinely UI-free.
- Protected directories (`data/`, `internal/entities/template/`,
  `internal/registry/`) were read only; nothing in this review proposes editing
  them.

**Correctness**

- Unknown topology values from a `.gen.json` are safe:
  [topologyServiceLookup.go](../internal/services/template_generator/providers/topologyServiceLookup.go#L56-L62)
  falls back to Ring, and
  [topologies.go](../internal/common/common_topologies/topologies.go#L139-L160)
  falls back to `descriptorValues.Default` / `topologies[0]` for both
  type-lookup and index-lookup. No panic path.
- `Load` correctly runs `onLoaded` only on success
  ([stateFiles.go](../app/gui/drivers/stateFiles.go#L23-L27)) — the historical
  panel-clobber regression is still fixed.
- `MapSize` from a `.gen.json` is snapped to the nearest valid size with a
  warning ([editorStateValidator.go](../internal/validators/editorStateValidator.go#L90-L102)).
- `GameMode` and `VictoryCondition` are validated against the registry with
  documented fallbacks.
- `zoneEditorZoneProps.go` no longer indexes `QualityLabels` by index (previously
  a panic path); it uses `SelectByName`.
- `layoutPanelZones.go` guards `lastTemplate`/`Variants[0]` before indexing;
  `hasTemplateVariants()` guards the same in `stateGeneration.go`.
- `AssetProvider` initialisation uses `sync.OnceValues`, so concurrent
  `NewAssetProvider` callers are race-free.
- `relaxPasses` is bounded at 500 iterations — no unbounded loop.
- The bridge-naming loop in `topologyConnectionService.go` always makes progress.
- `mandatoryContentProvider.go` deep-clones `Rules` before mutation.

**Tooling and hygiene**

- `go build ./...` — pass.
- `go vet -tags=integration_test ./...` — pass.
- `go test ./test/... -count=1` — pass (but see §6.1 for what it wrongly includes).
- `go test -tags=integration_test ./test/integration/... ./test/performance/...` — pass.
- `go test -tags='integration_test,gui' ./test/integration/gui/...` — pass, 0.78 s.
- `go mod tidy -diff` — **exit 0 for both the root and `tools/` modules**. The
  prior review's Windows/EOL discrepancy no longer reproduces.
- `git ls-files --ignored --exclude-standard -c` → 0 tracked ignored files.
- No tracked binaries, profiles, coverage artefacts or generated output.
- [.gitattributes](../.gitattributes) forces `*.go text eol=lf`, so the LF
  convention is enforced at checkout — no formatter CRLF findings.
- `wire_gen.go` is committed with the inverse `!wireinject` constraint, exactly
  as AGENTS.md §4.6.2 requires; `wire.go` correctly shows as excluded.
- The `integration_test` tag scoping is sound — only the one untagged file
  (§6.1) violates it.
- The single `// Reset only resets current edits…` comment is the **only**
  TODO/FIXME/HACK-class marker in the non-test Go tree.

---

## 12. Suggested execution order

Bugs first, PR-sized batches, blockers noted. Items keep these numbers
permanently — mark them `✅ FIXED` in place as they land.

1. **Security PR.** ✅ FIXED (2026-08-05)
   §8.1 + §8.2 — bump `golang.org/x/text` to v0.39.0 and `golang.org/x/net` to
   v0.56.0, `go mod tidy`, re-run `govulncheck`. Verify with `go build ./...`
   and the unit suite.

2. **Correctness PR.** ✅ FIXED (2026-08-05)
   §1.2 (`SaveAs` path), §1.3 (`WasLayoutChanged` nil guard + the
   `ShouldReapplyManualEdits` simplification), §5.2 (the commented-out line was
   activated rather than deleted). §6.1 (build tag) ❌ **rejected** — it
   contradicts the current AGENTS.md §4.6.1; see that item for the full
   rationale. §1.2's regression test landed in the integration suite, not
   `saveAs_test.go`, because the callback is not reachable headlessly.

3. **Durability PR.** ✅ FIXED (2026-08-05)
   §1.1 + §1.6 + §5.1. The owner replaced this item's `atomicFileWriter.go`
   inside `file_service` with a repository layer: `internal/repositories`
   gained a private `atomicFileWriter` (TEMP file → `Sync` → `Close` → retried
   `os.Rename`) plus concrete editor-state / template / preview repositories,
   and `FileService` became a pure controller with a single
   `SaveTemplateWithPreview` entry point. Transactionality decision: **keep the
   documented partial-success contract**. §5.1's constants went to a new
   `internal/constants` package instead of waiting for §2.1. Two approved
   behaviour changes: saves now create missing directories, and editor state is
   written as `{TemplateName}.gen.json` regardless of the name typed in the
   Save As dialog.

4. ✅ **Input-validation PR.** §1.5 (bounded counts) + §1.7 (override warnings).
   ⚠ Owner decision on the numeric ceilings first. Depends on nothing.

5. **Performance PR.** §4.1 — reinstate the benchmark, add the revision-keyed
   cache, record before/after. Independent of everything above.

6. **DI PR.** §2.3 + §2.4 — parameterise both constructors, register providers,
   `wire gen ./internal/composition/...`. Single PR because both touch
   `providerSets.go`.

7. **Test-policy PR.** §6.3 (drop GUI imports from the three tests) then §6.5
   (depguard scope + CI tag check). §6.5 must land *after* §6.3 or CI goes red.

8. **CI/security-posture PR.** §7.2 (top-level `permissions`), §7.3 (setup-go
   drift), §7.4 (tools tidy check), §8.3 (scheduled scan).
   ⚠ §7.1 only after the owner confirms the direct-push policy.

9. **Docs PR.** §9.1–§9.6 (+ optional §9.7) in one pass, then update repository
   memory. §9.5 must agree with whatever §2.7 decides.
   ⚠ §9.1 depends on the owner's public-API decision.

10. **Duplication cleanup PR.** §3.1 (mechanical, 15 sites), §3.3 (spell helper
    + its new tests), §3.4 (button widget — verify via GUI snapshots), then
    §3.2 + §5.3 together.

11. **Coverage PR.** §6.2 (`internal/handlers` mirrored tests — start with
    `stateHandler` and `previewHandler`), §6.4 (the two catalogues).

12. **Product decisions, then implementation.**
    ⚠ §2.7 (finish or remove the arena preview), ⚠ §1.8 (output-directory
    persistence shape).

13. **Large refactors — plan first per AGENTS.md §4.7.**
    §2.1 (extract filesystem policy) → unblocks §2.5. Then §2.2 (extract
    regeneration policy), which is multi-session and overlaps a backlog item.
    §2.6 opportunistically, whenever the zone editor is next touched.

**Blockers summary:** §6.5 after §6.3 · §2.5 after §2.1 · §5.1 folds into §1.1 or
§2.1 · §9.5 after §2.7 · §3.2 with §5.3.
**Owner decisions required before implementation:** §1.1 (transactionality),
§1.5 (ceilings), §1.8 (persistence shape), §2.2 (refactor scope), §2.7
(finish/remove), §7.1 (push policy), §9.1 (public API).

---

## 13. Measured baselines

| Check | Result |
| --- | --- |
| Reviewed revision | `687f47d6cff07dd2f42239c796dd8dad5385931a` (branch `AD/refactoring-07-21`) |
| Working tree before this document | Clean (`git status --short` empty) |
| `go version` | `go1.26.5 windows/amd64` |
| Root module / `tools` module | `go 1.26.5` / `go 1.26.3` |
| CI `GO_VERSION` | `1.26.5` (matches root) |
| `golangci-lint-v2 version` | `v2.12.2`, built with `go1.26.3` |
| `go build ./...` | **Pass** |
| `go vet -tags=integration_test ./...` | **Pass** |
| `go test ./test/... -count=1` | **Pass** — wrongly includes one integration file (§6.1) |
| `go test -tags=integration_test ./test/integration/... ./test/performance/...` | **Pass**; performance reports "no tests to run" (benchmarks only) |
| `go test -tags='integration_test,gui' ./test/integration/gui/...` | **Pass**, 0.782 s |
| Unit coverage (total) | **64.7%** (CI floor 60.0%) |
| Coverage detail (deduplicated profile) | 213 files; **49 below 80%**, **36 at 0%** |
| Lint, uncapped, configured set | **42** — 40 `gochecknoglobals`, 2 `dupl` (was 84 at the prior review) |
| Lint as CI runs it (4 linters disabled) | **0** |
| `govulncheck` symbol scan | **1 affected**: GO-2026-5970 (`x/text` v0.38.0, reachable from `program.go:21`) |
| `govulncheck -scan module` | **2**: GO-2026-5970 (`x/text`), GO-2026-5942 (`x/net` v0.55.0, not called) |
| `go mod tidy -diff` (root / tools) | **Exit 0 / Exit 0** |
| `git ls-files --ignored --exclude-standard -c` | **0** |
| Tracked binaries / artefacts | **0** |
| Production Go files | 339 non-test files, 24 872 LOC |
| Preview layout, Random 8p/16n, canvas 600 | **2.09 ms**, 40 378 B, 391 allocs — **per frame** (§4.1) |
| Preview layout, Ring 8p/16n | 15.3 µs, 17 256 B, 284 allocs |
| Preview layout, Random 2p (app default) | 1.28 µs, 832 B, 14 allocs |

Scratch benchmark used for the §4.1 measurements was created at
`test/performance/zzscratch_preview_layout_test.go` and **deleted**; the tree is
clean apart from this document. §4.1's fix plan proposes reinstating it
permanently as `test/performance/preview_layout_test.go`.
