# Batch P — Connection follows Zone (backlog §2.6 steps 2 + 3)

Close backlog §2.6 step 2 with the ruling **"`entities.Connection` in the DTO,
handler and GUI layers is a breach, not `.rmg.json` vocabulary"**, then act on it:
migrate `entities.Connection` → `template_model.Connection` across
`internal/dtos`, `internal/handlers`, `internal/services/connection_editor` and
`app/gui/**`, and move the runtime-only `IsUserAdded` flag off the frozen schema
entity onto the model. Removes **7 packages** from `entityNamerAllowList`,
closing steps 2 and 3 together.

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is
done, set its status to `Complete` and write its **Phase Summary** (what was
done, key decisions, anything needed to continue with zero context); run the
phase's **Verification Plan** and record the result before moving on. When all
phases are done, fill in **Final Recap** and **Deployment Plan**.

**Read `AGENTS.md` first.** Never modify `data/`, `internal/registry/` or
`internal/entities/template/` **without explicit owner approval** — phase 4 of
this plan needs exactly one such approval and must not proceed without it.
Never stage, never commit. Chain PowerShell with `;`. Never round-trip a `.go`
file through `Get-Content`/`Set-Content`; use `gofmt -r` on an explicit file
list and verify insertions == deletions per file.

### The ruling, and why (this is the durable answer to step 2)

The backlog framed step 2 as "`entities.Connection` / `entities.RmgTemplate` in
`internal/dtos` and `internal/handlers`". That framing was stale. Measured
2026-09-03:

- `entities.RmgTemplate` is **not named anywhere** in `internal/dtos`,
  `internal/handlers` or `app/gui` — batch J already removed it.
- The entire residual breach in those layers is **one type**,
  `entities.Connection`, in **8 files**.

Four findings decided it against a carve-out:

1. **The model twin already exists.** Batch J built
   [connection.go](../../internal/models/template_model/template_variant_model/connection.go)
   — field-identical, with `ToConnectionModels` / `ToConnectionEntities` already
   re-exported from `template_model/converters.go`. Nothing to design; the twin
   is sitting unused.
2. **The seams are already half-migrated.**
   [connectionEditorServiceInterface.go](../../internal/services/connection_editor/connectionEditorServiceInterface.go)
   reads `FindIsolatedZones(zones []template_model.Zone, connections []entities.Connection)`,
   and [zoneEditorGeometryService.go](../../internal/services/connection_editor/zoneEditorGeometryService.go)
   already calls `template_model.ToConnectionModels` **mid-service** to feed the
   preview. A service converting halfway through its own call is the smell.
3. **`IsUserAdded` is editor state inside the frozen schema entity.**
   [connection.go](../../internal/entities/template/template_variant/connection.go)
   carries it as `json:"-"`, "runtime-only". That is precisely the pollution the
   model layer exists to absorb — and `entities/editor_state.ManualConnectionSave`
   already carries a **sidecar** `IsUserAdded` field *because* the entity cannot
   serialize it. The workaround proves the field is misplaced.
4. **The base package is already slated for deletion.**
   [types.go](../../internal/entities/types.go): *"This needs to be removed and
   `internal/entities/template` should be used instead."* A permanent carve-out
   would bless a façade that is on its way out.

Owner ruled **A — Connection follows Zone**, with `IsUserAdded` moved to the
model only. Record this in backlog §2.6 in phase 5; it must survive this plan's
deletion.

### The precedent to copy

Everything here mirrors what batch J did for `Zone`. When in doubt, open the
zone-side twin of the file you are editing and follow it exactly:

| Connection work | Zone precedent to copy |
| --- | --- |
| `editor_state_model.To/FromManualConnectionSaves` | `To/FromManualZoneSaves` in [manualZoneSave.go](../../internal/models/editor_state_model/manualZoneSave.go) |
| Re-attaching applied connections in `UpdateTemplate` | `updated.Variants[0].Zones = zones` in [templateHandler.go](../../internal/handlers/templateHandler.go) |
| DTO field type | `Zones []template_model.Zone`, already in all five DTOs |

### Non-negotiable invariants

- **Wire format must not change.** `template.Connection.IsUserAdded` is
  `json:"-"`, so removing it emits and consumes byte-identical JSON. The two
  frozen fixtures under `test/test_helpers/testdata/` carry `isUserAdded` at the
  `ManualConnectionSave` level (line 224 in both), **not** inside `connection` —
  verify, do not assume. Neither fixture may be edited.
- **Conversion happens at exactly two seams** (AGENTS §4.4.1 rule 4):
  `internal/handlers` (DTO ⇄ Model) and `internal/repositories` /
  `editor_state_model` (Model ⇄ Entity). `NewDefaultConnection` converting the
  builder's output is a third, tolerated only because the builder is shared with
  the generator — flag it in the phase summary, do not spread the pattern.
- **`app/` must never name an entity.** It may hold a Model (AGENTS §4.4.1
  rule 2). That is the whole point of the batch.
- No golden may move. Run the GPU suite **without** `-update`.
- Coverage floor **72.5 %** (baseline **74.3 %**); lint baseline **0 issues**.

---

## Phase 1: Baseline and the `IsUserAdded` regression guard
Status: Complete

The one real behaviour risk in this batch is that `UpdateTemplate` round-trips
the template through the entity, so once `IsUserAdded` leaves the entity the
flag silently dies on every Apply — the exact failure mode batch J fixed for the
zone tier. Pin it before touching anything.

- [x] Record the baseline: `go build ./...`, `go vet ./...`, `gofmt -l`,
      `go run ./cmd/testlayoutcheck .`, unit tests with coverage, integration,
      GPU integration, `golangci-lint-v2 run ./...`. Write the numbers here.
- [x] Confirm the frozen fixtures place `isUserAdded` on `ManualConnectionSave`
      and never inside `connection` (`test/test_helpers/testdata/editorState_v0_flat.gen.json`
      and `editorState_v1_flat.gen.json`, line 224 in both).
- [x] Confirm `test_helpers.allFieldsManualConnections()` does **not** set
      `IsUserAdded` inside the `entities.Connection` literal — if it does, the
      phase 4 field removal will not compile.
- [x] Add a unit test under
      `test/unit/internal/handlers/templateHandler/` asserting a
      connection with `IsUserAdded: true` **still carries the flag** on the
      returned `TemplateLoadDto.Template`. It must pass on today's tree.
- [x] Mutation-verify that guard: temporarily blank `IsUserAdded` in
      `ToConnectionModel`, watch the new test fail, revert.

### Verification Plan
- `go test ./test/unit/internal/handlers/... -count=1` → exit 0, new test present and passing.
- Mutation step fails the new test and nothing else surprising; `git status --short` clean of the mutation afterwards.

### Phase Summary

**Baseline recorded (branch `AD/fixing_some_stuff_08-12`, batch J committed):**

| Gate | Result |
| --- | --- |
| `go build ./...` / `go vet ./...` | exit 0 |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test ./test/unit/... -count=1` + coverage | exit 0 — **74.3 %** |
| `go test -tags=integration_test ./test/integration/...` | exit 0 |
| `go test -tags='integration_test,gui' ./test/integration/gui/...` | exit 0, `goldens=0` |
| `golangci-lint-v2 run ./...` | **0 issues** |

**Fixture assumptions verified, not assumed.** Both frozen fixtures put
`"isUserAdded": true` on the `ManualConnectionSave` object; the nested
`"connection"` object has no such key. `allFieldsManualConnections()` sets the
flag only on the sidecar, never inside the `entities.Connection` literal — so
the phase 4 field removal will compile and the wire format cannot move.

**Guard added:** `TestWhenAnAppliedConnectionIsUserAdded_KeepsTheFlagThroughTheEntityRoundTrip`
in [updateTemplate_test.go](../../test/unit/internal/handlers/templateHandler/updateTemplate_test.go).
⚠ **Plan correction:** the plan sketched a nested
`templateHandler/updateTemplate/` folder. That is wrong — AGENTS §4.6 puts the
folder at the *implementation file* name and the file at the *function* name, so
it belongs in the existing `templateHandler/updateTemplate_test.go`.

**Mutation-verified:** blanking `IsUserAdded` in `ToConnectionModel` fails the
new test with `Should be true` at line 91 and nothing else. Reverted;
`git diff --numstat` afterwards shows only `updateTemplate_test.go` (19/0), so
`connection.go` is untouched.

**Note on the working tree:** `.agent/plans/batch-j-…md` (deleted) and
`.agent/plans/batch-k-…md` (added) are **staged by the owner**. Per AGENTS §2.5
they were left exactly as found — nothing staged or unstaged by the agent.

## Phase 2: The type swap
Status: Complete

One atomic change — the compiler will not tolerate a partial swap, since
handlers call the services and `app/` calls the handlers. Use `gofmt -r
'entities.Connection -> template_model.Connection'` on an **explicit file list**,
then fix imports by hand and verify insertions == deletions per file.

Production files (21):

- [ ] `internal/dtos` (5): `templateUpdateDto.go`, `zoneEditorGeometryRequestDto.go`,
      `zoneEditorMutationDto.go`, `zoneEditorRemoveRequestDto.go`, `zoneEditorZonesDto.go`.
      All five already import `template_model`; the `entities` import drops out.
- [ ] `internal/handlers` (3): `guiHandler.go`, `zoneEditorHandler.go`,
      `handler_interfaces/zoneEditorHandlerInterface.go`. Pure passthrough — no
      conversion needed once the services move.
- [ ] `internal/services/connection_editor` (6): `connectionEditorService.go` +
      `Interface`, `zoneEditorService.go` + `Interface`, `zoneEditorGeometryService.go` +
      `Interface`. Delete the now-redundant `template_model.ToConnectionModels`
      call at `zoneEditorGeometryService.go:69` — the input is already a model.
- [ ] `NewDefaultConnection` returns `template_model.Connection`: keep the
      `variant_content` builder for the field values, then
      `template_model.ToConnectionModel(...)`. One-line comment saying why the
      conversion is here.
- [ ] `app/gui` (6): `dialogs/zoneEditorCanvas.go`,
      `dialogs/zoneEditorConnectionPropertiesState.go`, `dialogs/zoneEditorDialog.go`,
      `dialogs/zoneEditorInteractionState.go`, `drivers/stateManualEdits.go`,
      `models/editorState.go`.
- [ ] `internal/models/editor_state_model/manualConnectionSave.go`:
      `ToManualConnectionSaves([]template_model.Connection)` /
      `FromManualConnectionSaves() []template_model.Connection`, converting with
      `template_model.ToConnectionEntity` / `ToConnectionModel`. Mirror
      `manualZoneSave.go` exactly. `cloneConnection` moves onto the model type.
- [ ] `internal/handlers/templateHandler.go` `UpdateTemplate`:
      `newTemplate.Variants[0].Connections = template_model.ToConnectionEntities(connections)`,
      **and add the positional re-attach** `updated.Variants[0].Connections = connections`
      next to the existing zone re-attach, with a comment. This is what makes
      phase 4 safe by construction.

Test-side (mocks first, they gate compilation):

- [ ] `test/test_helpers`: `connectionEditorServiceMock.go`, `zoneEditorServiceMock.go`,
      `zoneEditorGeometryServiceMock.go`, `templateHandlerMock.go`.
      `allFieldsEditorState.go` keeps `entities.Connection` — it builds the
      **entity** save and must stay entity-side.
- [ ] Unit tests under `test/unit/internal/services/connection_editor/**`,
      `test/unit/internal/handlers/**`, `test/unit/app/gui/**`.
- [ ] Integration: `editorState_integration_test.go`,
      `manualCastleReapply_integration_test.go`, `gui/zoneEditorDialog_integration_test.go`.

### Verification Plan
- `go build ./...` and `go vet -tags='integration_test,gui' ./...` → clean.
- `gofmt -l ./app ./internal ./test ./cmd` → empty.
- `go test ./test/unit/... -count=1`, `go test -tags=integration_test ./test/integration/... -count=1`, `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` → all exit 0.
- `git status --short -- '*.golden'` → **zero lines**. GPU suite run without `-update`.
- `git diff --numstat` per file: insertions == deletions on every pure-swap file.

### Phase Summary

**Done. 62 files, +260/−236, all gates green.**

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | exit 0 both |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | passed |
| `go test ./test/unit/... -count=1` | exit 0 — coverage **74.3 %**, unchanged |
| `go test -tags=integration_test ./test/integration/...` | exit 0 |
| `go test -tags='integration_test,gui' ./test/integration/gui/...` | exit 0, **no `-update`** |
| `git status --short -- '*.golden'` | `goldens=0` |
| `git status --short -- 'test/test_helpers/testdata/*'` | `fixtures=0` |
| `golangci-lint-v2 run ./...` | **0 issues** |

**Technique.** Two `gofmt -r` rules on explicit file lists —
`entities.Connection -> template_model.Connection` and, on test files only,
`template_model.ToConnectionEntities(a) -> a`. Insertions == deletions on every
pure-swap file (the one 22/3 file carries the 19-line phase-1 test). Imports
were then fixed by hand with the compiler as the oracle; **no file was
round-tripped through `Get-Content`/`Set-Content`**.

⚠ **The second rule was applied to test files only.** `templateHandler.go`
genuinely needs `ToConnectionEntities`, so sweeping it repo-wide would have
silently broken the entity assignment.

**Result — the packages are entity-free, verified by grep:** `internal/dtos`,
`internal/handlers` and `internal/services/connection_editor` have **zero**
`internal/entities` imports in production, and `app/gui/**` names **zero**
entity types *including its tests*. (One hit remains in
`test/unit/internal/handlers/templateHandler/updateTemplate_test.go` for
`entities.MandatoryContent`; `productionRoots` is `app`/`internal`/`cmd`, so
`test/` is outside the gate entirely.)

**Substantive edits, beyond the mechanical swap:**

1. **`template_model/converters.go` gained `ToConnectionEntity` / `ToConnectionModel`** —
   only the plural forms were re-exported. Needed by the two real seams.
2. **`editor_state_model.ManualConnectionSave`** now takes and returns
   `[]template_model.Connection`, converting with the singular converters. The
   save's `Connection` field **stays an entity** — that is the Model⇄Entity seam
   doing its job. Mirrors `manualZoneSave.go` line for line.
3. **`templateHandler.UpdateTemplate`** converts on the way in
   (`ToConnectionEntities`) and now **re-attaches the applied connections
   positionally** on the way out, beside the existing zone re-attach. This is
   what makes phase 4 safe by construction, and the phase-1 guard covers it.
4. **`NewDefaultConnection`** returns a model, converting the shared builder's
   entity output. ⚠ This is a **third** conversion seam, outside handlers and
   repositories. It is tolerated because `variant_content.NewConnectionBuilder`
   is shared with the generator; it is commented in place. **Do not spread the
   pattern** — phase 4 removes half its reason to exist.

**Four conversions turned out to be pure ceremony and were deleted**, which is
the batch paying for itself: `zoneEditorGeometryService.go:69`,
`stateManualEdits.PreviewBaseZones`, `layoutPanelZones.go` and
`manualZoneTierPersistence_integration_test.go:102` all called
`ToConnection*` to cross a boundary that no longer exists.

**Three files were NOT in the planned list and were found by the compiler:**
`app/gui/panels/layoutPanelZones.go`,
`test/unit/app/gui/drivers/state/getTemplateRevision_test.go` and
`test/unit/app/gui/drivers/stateManualEdits/previewBaseZones_test.go` — none
*name* an entity, so no grep for `entities.Connection` would have found them;
they called `ToConnectionEntities`. **Lesson: when migrating a type, grep for
its converters as well as its name.**

**Two tests were genuinely mixed and were hand-edited, not swept** —
`toManualConnectionSaves_test.go` and `fromManualConnectionSaves_test.go`. The
save literal keeps `entities.Connection`; only the live connection becomes a
model. A blind sweep would have "fixed" the entity side and destroyed what the
test is for.

**One lint fix:** `golines` flagged the widened
`GUIHandler.CreateZoneEditorConnection` signature. Restyled by hand to
`param) Ret {`, the house style — not the `param,\n) Ret {` a `--fix` produces.

## Phase 3: Shrink the allow-list
Status: Complete

- [x] Remove from `entityNamerAllowList` in
      [layering_test.go](../../test/unit/architecture/dependency/layering_test.go):
      `app/gui/dialogs`, `app/gui/drivers`, `app/gui/models`, `internal/dtos`,
      `internal/handlers`, `internal/handlers/handler_interfaces`,
      `internal/services/connection_editor`.
- [x] Prove each removal is real, the way batch J did: re-introduce an
      `internal/entities` reference in the package, watch
      `TestWhenEntityConsumersAreScanned_OnlyPermittedPackagesNameAnEntity` name
      that exact file, revert. Confirm clean after each.
- [x] Update the allow-list comment: **exactly one** permanent entry,
      `file_service`. **Only ever remove entries.**

### Verification Plan
- `go test ./test/unit/architecture/... -count=1` → exit 0 with the seven entries gone.
- Each mutation reports the reintroduced file and only that file; clean afterwards.

### Phase Summary

**`entityNamerAllowList` went from 21 entries to 14.** Removed:
`app/gui/dialogs`, `app/gui/drivers`, `app/gui/models`, `internal/dtos`,
`internal/handlers`, `internal/handlers/handler_interfaces`,
`internal/services/connection_editor`. `go test ./test/unit/architecture/...`
is green with them gone.

**Mutation proof.** The rule keys on the **import path**, not on naming a type
(`isEntityImport` is the guard, `findUnlistedNamers` filters by file path), so a
blank `_ "…/internal/entities"` import is a sufficient and minimal mutation. All
seven were introduced **simultaneously**, one file per package — which proves
each entry independently *and* proves no entry was masking another. The gate
failed with exactly seven entries, one per package:

```
app/gui/dialogs/zoneEditorInteractionState.go
app/gui/drivers/stateManualEdits.go
app/gui/models/editorState.go
internal/dtos/zoneEditorMutationDto.go
internal/handlers/handler_interfaces/zoneEditorHandlerInterface.go
internal/handlers/zoneEditorHandler.go
internal/services/connection_editor/connectionEditorService.go
```

Reverted; a repo-wide grep for `_ "…/internal/entities"` returns nothing.
⚠ `git status` **cannot** confirm this revert — all seven files are already
modified by phase 2. Grep for the mutation itself, not the file's dirty flag.

**Comment updated** — and ‼ **corrected after owner review.** The first version
of this comment said *two* entries were permanent (`file_service` **and**
`template_generator` + topology). **That is wrong and directly contradicted a
settled decision**: `.agent/memories/settled-decisions.md` §"Entity confinement"
and `.agent/memories/architecture.md` both record that **exactly one** entry is
permanent, `file_service`, and that the generator's entity usage is **debt**.
The claim came from an earlier carry-forward that had already been corrected on
2026-09-03; it was re-propagated without checking the memory. The comment now
says one, and names the generator explicitly as debt so the error cannot repeat.
**Read `settled-decisions.md` before writing anything about what is permanent.**

## Phase 4: `IsUserAdded` moves to the model — ⚠ NEEDS OWNER APPROVAL
Status: Complete — **approved by the owner 2026-09-04**

- [x] Propose and get approval for removing
      `IsUserAdded bool \`json:"-"\`` (and its comment) from
      `internal/entities/template/template_variant/connection.go`.
- [x] Drop `IsUserAdded` from `ToConnectionModel` / `ToConnectionEntity`.
- [x] Remove `ConnectionBuilder.WithIsUserAdded()` and delete
      `withIsUserAdded_test.go` (`Remove-Item`). `NewDefaultConnection` sets the
      flag on the model after conversion instead.
- [x] Refresh the two stale comments referencing the entity field.
- [x] Confirm the phase 1 guard still passes.

### Verification Plan
- Full gate set: build, vet (both tag sets), `gofmt -l`, `testlayoutcheck`, unit + coverage, integration, GPU integration, `golangci-lint-v2 run ./...` at **0 issues**.
- Fixtures unmodified.
- Coverage ≥ 72.5 %; record the actual figure.

### Phase Summary

**Approved and applied.** The `.rmg.json` schema mirror is now a pure mirror:
the only `IsUserAdded` left under `internal/entities/` is the editor-state
sidecar `editor_state.ManualConnectionSave.IsUserAdded`, which is where editor
state belongs. `WithIsUserAdded` no longer exists anywhere in the repo.

| Gate | Result |
| --- | --- |
| `go build ./...` / `go vet` both tag sets | exit 0 |
| `gofmt -l` | empty |
| `go run ./cmd/testlayoutcheck .` | passed |
| `wire diff ./internal/composition/...` | exit 0 |
| `go test ./test/unit/... -count=1` + coverage | exit 0 — **74.3 %**, unchanged from baseline |
| integration / GPU integration | exit 0, no `-update` |
| `git status --short -- '*.golden'` / `testdata/*` | `goldens=0`, `fixtures=0` |
| `golangci-lint-v2 run ./...` | **0 issues** |

**The wire format did not move, and that was verified rather than assumed.**
The field was `json:"-"`, so removal changes no byte. `fixtures=0` after the
untagged wire-format test and the whole integration suite confirms it.

**⚠ The important proof: the phase-1 guard is now load-bearing, not vacuous.**
Before this phase it passed either way, because the entity carried the flag
through the round trip. With the field gone, deleting
`updated.Variants[0].Connections = connections` from `UpdateTemplate` makes
`TestWhenAnAppliedConnectionIsUserAdded_KeepsTheFlagThroughTheEntityRoundTrip`
**fail** — verified by mutation and reverted. That single line in
[templateHandler.go](../../internal/handlers/templateHandler.go) is the only
thing keeping a hand-drawn connection marked as user-added across an Apply.
Do not "simplify" it away.

**One test was deleted, not migrated:**
`TestWhenSaveFlagDiffersFromEmbeddedConnectionFlag_SaveFlagWins` asserted that
the save's flag beats the one embedded in the entity connection. The entity no
longer has a flag to disagree with, so the scenario is unconstructible. The
surviving `TestWhenSavesCarryUserAddedFlags_RestoresEachFlagOntoConnection`
still fails if `FromManualConnectionSaves` stops copying the flag, so nothing
went uncovered.

**`NewDefaultConnection`'s third conversion seam is now half-justified.** It
still converts the shared builder's entity output, but it must also set
`IsUserAdded` on the model afterwards, because the builder can no longer express
it. If the generator is ever moved onto models (backlog §2.6 step 4), this
whole function collapses to a model literal and the seam disappears.

## Phase 5: Durable records
Status: Complete

- [x] `.agent/backlog/backlog-opus5.md` §2.6: steps 2 and 3 **✅ DONE**, ruling
      and its four findings written self-contained, breach recounted
      (84/21 → **64/14**), heading retitled, header and §8 updated.
- [x] `.agent/memories/template-model.md`: "Remaining entity usage" refreshed;
      new "Connection (batch P)" section with the three traps.
- [x] `.agent/memories/settled-decisions.md`: "Base `internal/entities` gets NO
      vocabulary carve-out (ruled 2026-09-04)".
- [x] **Final Recap** and **Deployment Plan** written; carry-forward rewritten.

### Verification Plan
- Every §2.6 claim is checkable without this plan file open.
- Final full gate run recorded in the Final Recap.

### Phase Summary

⚠ **Batch letter corrected mid-phase: this is batch P, not K.** §8 of the
backlog already used **⚠ K** for the owner-gated group (§2.2 Branch A, §2.4,
§2.5, §6.1). Letters through O were taken, so P was the next free one; the plan
file was renamed with `Move-Item` and every reference updated. **Check §8 for a
free letter before naming a batch.**

The recount was done by script rather than by trusting the old figure: walk
`app`/`internal`/`cmd`, exclude `test`, keep files importing
`internal/entities`, then subtract the permitted namers
(`repositories|models|entities|mappers`, `helpers/*_helpers`). Result **64 files
in 14 packages**, which matches `entityNamerAllowList` entry-for-entry — a
useful cross-check that the list has no dead entries.

## Final Recap

**Backlog §2.6 steps 2 and 3 are closed.** The ruling: base `internal/entities`
gets **no carve-out** — naming a `.rmg.json` schema type below the repositories
is a genuine breach. Recorded in backlog §2.6 step 2 (self-contained, survives
this file) and in `settled-decisions.md` so it is not re-litigated.

The ruling was then acted on in the same batch:

| Measure | Before | After |
| --- | --- | --- |
| Entity breach | 84 files / 21 packages | **64 / 14** |
| `entityNamerAllowList` | 21 entries | **14** |
| `internal/dtos` naming entities | 5 files | **0** |
| `internal/handlers` (+interfaces) | 3 files | **0** |
| `internal/services/connection_editor` | 6 files | **0** |
| `app/gui/**` (incl. tests) | 6 files | **0** |
| Coverage | 74.3 % | **74.3 %** |
| Lint | 0 issues | **0 issues** |

**Five phases.** (1) Baseline plus a mutation-verified guard on the one real
risk. (2) The type swap — 62 files, `gofmt -r` on explicit lists, insertions ==
deletions per file, compiler as the oracle for imports. (3) Seven allow-list
entries removed, all seven mutation-proved simultaneously. (4) Owner-approved
removal of `IsUserAdded` from the protected schema entity. (5) Records.

**What the batch actually bought, beyond a shorter list.** Four
`ToConnection*` calls were deleted outright because they crossed a boundary that
no longer exists, and the `.rmg.json` mirror stopped carrying an editor-only
flag it could not even serialize. The layering rule now polices the DTO, handler
and GUI layers for real rather than by exception.

**Three things a future agent should not have to rediscover:**

1. **`updated.Variants[0].Connections = connections` in `UpdateTemplate` is
   load-bearing.** With `IsUserAdded` off the entity, that line is the only
   thing carrying the flag across an Apply. Mutation-verified: remove it and
   the phase-1 guard fails. Its zone twin beside it does the same for the tier.
2. **When migrating a type, grep its converters as well as its name.** Three
   files called `ToConnectionEntities` without ever naming `entities.Connection`
   — no name-based grep would have found them, only the compiler did.
3. **`git status` cannot verify a mutation revert** when the file is already
   dirty from the same batch. Grep for the mutation itself.

**Final gate run — all green, nothing staged by the agent:**

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `-tags='integration_test,gui'` | exit 0 both |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 |
| `go test ./test/unit/... -count=1` + coverage | exit 0 — **74.3 %** (floor 72.5 %) |
| `go test -tags=integration_test ./test/integration/...` | exit 0 |
| `go test -tags='integration_test,gui' ./test/integration/gui/...` | exit 0, **no `-update`** |
| `git status --short -- '*.golden'` | `goldens=0` |
| `git status --short -- 'test/test_helpers/testdata/*'` | `fixtures=0` |
| `golangci-lint-v2 run ./...` | **0 issues** |

## Deployment Plan

1. **Review the one protected-directory edit first.**
   `internal/entities/template/template_variant/connection.go` loses only
   `IsUserAdded bool \`json:"-"\`` and its comment. Confirm nothing else in
   `internal/entities/template/` changed:
   `git diff --stat -- internal/entities/template/`.
2. **Confirm the wire format is untouched** — `git status --short --
   'test/test_helpers/testdata/*'` must be empty, and
   `git status --short -- '*.golden'` must be empty. Neither may be regenerated.
3. **Run the full gate set** (the table above). The GPU suite must be run
   **without** `-update`.
4. **Read the two load-bearing lines** in
   [templateHandler.go](../../internal/handlers/templateHandler.go): the zone
   and connection re-attaches after `ToModel`. They look redundant and are not.
5. **In-app smoke test — this is the path the batch exists to protect:**
   generate a template; open the zone editor; **draw a new connection by hand**;
   Apply; confirm the connection still renders with its user-added styling
   (`zoneEditorCanvas` and the connection properties panel both branch on
   `IsUserAdded`). Then save the `.gen.json`, reload it, reopen the editor and
   confirm the connection is *still* marked user-added. Also re-tier a neutral
   zone in the same session to confirm batch J's guarantee still holds.
6. **Load a legacy `.gen.json`** written before this batch. It must open with no
   error and no connection wrongly marked user-added — `isUserAdded` was always
   an `omitempty` sidecar, so absent means false.
7. Commit. Then delete the transient docs per the owner's doc-lifecycle rule:
   this plan and `.agent/session-carry-forward.md`. **Backlog §2.6 is the
   surviving record** and is written to stand alone.
