# Batch I: docs (review §7.1, §7.2)

Fix the review's §7.1 (stale README/QUICKSTART descriptions) and §7.2 (stale test
observations), every other inaccuracy the 2026-10-04 audit found in those files, two stale Go
comments, and the AGENTS.md interface example. Documentation and comments only: no
behaviour change.

## For Future Agents
As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

**Starting state** (2026-10-04): branch `AD/pbi_resolution`, HEAD `879e41a` ("docs"), with
Batch H committed as `a140e0a`. The only working-tree change is the review's Batch H close-out
edit in `.agent/backlog/review-gpt-6-astra-09-07.md`.

**Hard rules in force:** never stage, unstage, commit, push, stash, switch branches or touch
worktrees; never touch `data/`, `internal/entities/template_entity/`, `internal/registry/`,
generated Wire or the output directory; no bulk rewrites; no repo-committed check scripts
(throwaway checks only, as in Batch H).

### Owner decisions (2026-10-04)

| ID | Decision |
| --- | --- |
| I1 | §7.1 fixes every audited inaccuracy in README/QUICKSTART. It is not a rewrite: correct text stays verbatim. |
| I2 | README names Gio without a version and links [go.mod](../../go.mod). |
| I3 | Generation Flow: `drivers.State` → `handler_interfaces.IGuiHandler` → `GUIHandler` → `templateHandler.GenerateTemplate` → `TemplateGenerator.Generate` → `*template_model.Template` (+ warnings) → preview, and `templateHandler.SaveTemplate` → `FileService.SaveTemplateWithPreview` → `TemplateMapper.ToEntity` (the Model→Entity seam) → `TemplateRepository.Save` → `.rmg.json`. |
| I4 | The fixed QUICKSTART §5 snippet is compile-checked once from a throwaway `tmp/quickstartcheck/main.go` (gitignored), then deleted. It is never run. |
| I5 | test_observations: delete entries whose limitation no longer exists, correct stale paths and symbols, and label historical figures as historical. |
| I6 | The integration_common entry describes the layout (tagged helpers vs the untagged `snapshot/` subpackage) instead of listing every file. |
| I7 | Fix the two stale Go comments (`templateHandlerMock.go`, `logButtonPositions_test.go`). |
| I8 | Update the AGENTS.md §4.2.1 example from the non-existent `IBackend` to the real `IGuiHandler`. |
| I9 | Verification is docs checks only: throwaway link/anchor check, the snippet compile check, `go vet` + `gofmt -l` on the two comment-touched packages. No coverage, unit or testlayoutcheck run. |
| I10 | Reviews by GPT-6.1 Sol (copy is written by the main agent, per the AGENTS.md §3.4 taste rule). |

**Plan review** (GPT-6.1 Sol): APPROVE WITH CHANGES; all 7 findings applied:
- an owner-commit gate before marking items fixed;
- the `-run` check done by source inspection instead of a test run;
- the full detection mechanism described per platform;
- the QUICKSTART persistence layering corrected;
- the button semantics helper separated from the constructor limitation;
- the `reapplyManualEdits` location and rationale corrected;
- the fixture-tested io.go Linux path acknowledged.

**Owner plan approval:** 2026-10-04, "Approved, please proceed".

**Out of scope:** topology retirement wording (Batch K; "Eleven topologies" is accurate today),
the QUICKSTART map-size table (labelled "common presets", accurate), restructuring or
restyling any document, owner_findings.md, and the review's §8.

## Phase 1: README (§7.1)
Status: Complete

- [x] Features: `gioui.org v0.10.0` → Gio named without a version, linked to [go.mod](../../go.mod) (I2).
- [x] Features, Steam auto-detection: replace the bullet's mechanism and fallback with what the
      startup path (`FindGameTemplateDirectory` → `FindOldenEraTemplatesDir(false)`) does:
      - on Windows it globs the user-profile
        `AppData/LocalLow/Unfrozen/HeroesOldenEra/users/*/my_map_templates`, with no Steam parsing;
      - on Linux/Steam Deck it reads Steam's `libraryfolders.vdf` and uses the Proton prefix's
        `my_map_templates`;
      - the install `map_templates` folder is not used at startup;
      - on failure the output folder stays empty with an error status, and export is refused
        until the user picks the game's templates folder (session-only picker).
- [x] Workflow steps 1 and 6 (and QUICKSTART §4's closing line, Phase 2): the template is
      written straight into the detected templates folder, so "drop it into the game's
      templates folder" becomes "pick it in-game". Step 1 drops "via Steam".
- [x] Building & Running: `HOT_RELOAD=1` → `.air.toml` passes `-minimized` (and `-with-logging`).
- [x] Project Structure tree:
  - `entities/ … (template/)` → `template_entity/`;
  - `Thin GuiHandler facade` → `GUIHandler`;
  - add the missing `services/` packages (`bonuses/`, `editor/`, `file_system/`, `zone_content/`);
  - add the missing `models/` packages (`content_rule_model/`, `neutral_zone/`, `preview/`,
    `regeneration/`, `template_model/`);
  - add `data/GameData/DB/` and `data/Images/`.
  Verify each name with `list_dir` before writing.
- [x] Game Modes: drop "generator currently always emits `Classic`" and "(reserved)". State
      that the selected mode is emitted, and that `SingleHero` applies its hero rules.
      Re-check `gameRulesProvider.go` before writing.
- [x] Generation Flow diagram per I3. Replace the `app/gui/interfaces.IBackend` line with
      `drivers.State` calling `handler_interfaces.IGuiHandler`.
- [x] Testing: `-run TestWhenStateIsSaved` on `file_service` → a path where the name exists
      (`./test/unit/internal/repositories/...`), or a real `file_service` test name.
- [x] License: "MIT, see `[LICENSE](LICENSE)`" (the README's own link).

### Verification Plan
- Throwaway PowerShell link check: every relative link/image in README resolves (`Test-Path`
  after URL-decoding, anchor stripped). In-file anchors match a heading slug.
- `Select-String README.md -Pattern 'IBackend|HOT_RELOAD|v0\.10\.0|working directory|always emits|See the main project'` → no hits.
- Source inspection (no test run, I9): the documented `-run` pattern matches at least one
  `func Test…` declaration under the documented package path.

### Phase Summary
Done 2026-10-04. README now:
- names Gio by link and points to go.mod for the version;
- describes per-platform template-folder detection and the refused export on failure
  (the bullet is renamed "Game templates folder auto-detection");
- describes the air flags;
- updates workflow steps 1 and 6;
- states that the selected game mode is emitted, with `SingleHero`'s hero-hire ban and
  starting-hero loss rule (`gameRulesProvider.go` L28/L108);
- redraws the Generation Flow per I3, ending in `TemplateMapper.ToEntity` (the seam, inside
  `FileService.SaveTemplateWithPreview`) and `TemplateRepository.Save`;
- points the `-run` example at `./test/unit/internal/repositories/...`;
- states MIT with a link to LICENSE.

The tree adds `data/Images/`, `GameData/DB/`, 4 services and 5 model packages (names
listed from disk) and fixes `template_entity/` and `GUIHandler`.

Verification:
- A throwaway link checker (`%TEMP%\batchi-linkcheck.ps1`, outside the repo; deleted at
  close-out) reports README's 14 relative links all OK. Its negative control caught a missing
  path and a missing anchor.
- The stale-term grep has no hits.
- `TestWhenStateIsSaved_*` exists twice in `repositories/editorStateRepository/save_test.go`.

## Phase 2: QUICKSTART (§7.1)
Status: Complete

- [x] §3: the persistence sentence distinguishes three things: the editor state is the
      `editor_state_model.EditorState` model, which crosses the handler boundary wrapped in
      `editor_state_dto.EditorStateDto`; `file_service.FileService.SaveSettings` /
      `LoadSettingsFile` map it to and from the `.gen.json` entity
      (`internal/entities/editor_state`).
- [x] §4 step 1 and closing line: the folder is detected (not "from your Steam install" on
      Windows), and the template is written into it, so "drop it into the game's templates
      directory" becomes "pick it in-game".
- [x] §5 snippet:
  - `state := editor_state_dto.EditorStateDto{EditorState: editor_state_model.NewDefaultEditorStateModel()}`;
  - output path from `composition.InitializeFileSystemHandler().FindGameTemplateDirectory()`
    (exit on error) instead of `"."`, matching AGENTS.md §2.7;
  - imports updated.
- [x] §5 table: the six embedded interfaces of `IGuiHandler` (`ITemplateHandler`,
      `IStateHandler`, `IPreviewHandler`, `IZoneContentHandler` (embeds `IContentRuleHandler`),
      `IZoneEditorHandler`, `IBonusHandler`). Note the two standalone seams with their own
      injectors: `IFileSystemHandler` and `IRegenerationHandler`. The "whole contract" sentence
      is adjusted to match.
- [x] §5 prose: `dtos.NewDefaultEditorStateDto()` → `editor_state_model.NewDefaultEditorStateModel()`.
      "exchanges `internal/dtos` types" stays (DTOs cross the seam; the default state is a
      model wrapped in the DTO).

### Verification Plan
- Copy the snippet verbatim into `tmp/quickstartcheck/main.go`; `go vet ./tmp/quickstartcheck`
  and `go build -o <temp>.exe ./tmp/quickstartcheck` both succeed; never run the binary;
  delete the folder and the binary.
- Link/anchor check as in Phase 1, for QUICKSTART.
- `Select-String QUICKSTART.md -Pattern 'NewDefaultEditorStateDto|dtos\.EditorStateDto|OutputPath: "\."'` → no hits.

### Phase Summary
Done 2026-10-04. In QUICKSTART:
- §3 separates the model, the DTO wrapper and the `.gen.json` entity; checked against
  `FileService.SaveSettings` / `LoadSettingsFile`, which map through `editorStateMapper`.
- §4 describes the detected folder, the session-only picker on failure, and "pick it in-game".
- The §5 snippet wraps `NewDefaultEditorStateModel()` in `editor_state_dto.EditorStateDto`
  and writes to `composition.InitializeFileSystemHandler().FindGameTemplateDirectory()`, with
  a one-line comment citing the game's folder rule.
- The §5 table lists `IGuiHandler`'s six embeds and names the two standalone injectors
  (`InitializeFileSystemHandler`, `InitializeRegenerationHandler`, verified in
  `internal/composition/wire.go`). "A single interface" became "interfaces".

Verification:
- The snippet was extracted verbatim from the markdown into `tmp/quickstartcheck/main.go`.
  `go vet` and `go build` (to a temp exe) both exit 0. It was never run, and the folder and
  exe were deleted.
- All 11 relative links are OK, and the stale-term grep (plus "from your Steam install" and
  "Drop the resulting") has no hits.

## Phase 3: test_observations.md (§7.2)
Status: Complete

- [x] buttonWidget entry: separate the helper from the constructor limitation. The button
      constructors still need `layout.Context` + a text shaper. The semantics helper is now the
      exported `utils.AddButtonSemantics` (app/gui/utils/buttonPositionLogger.go), which needs
      only ops, a label and dimensions; keep the accurate note that the logger tests replay
      equivalent ops.
- [x] zoneEditorDialog entry: `groupConnectionsByPair` now lives in `internal/models/preview`
      (connectionCurveLayout.go, the shared curve builder). The ">=92.9%" figure is labelled
      as measured 2026-08-08.
- [x] drivers.State intro: `NewUIState(handler, false)` → the 4-argument
      `NewUIState(handler, fileSystem, regeneration, findTemplateDir)` with the handler mocks.
- [x] Delete the `templateDir == ""` fallback bullet (detection is now injected via
      `IFileSystemHandler.FindGameTemplateDirectory` and unit-tested). Keep `GetOutputPathWidget`
      and repair the orphaned "(returns a Gio widget)" sentence.
- [x] stateFiles bullet: `suggestDirectory` → `getWorkingDirectory`. Check against the existing
      `getWorkingDirectory` bullet so the two stay consistent.
- [x] `reapplyManualEdits` bullet: it lives in app/gui/drivers/stateManualEdits.go and delegates
      through mockable handlers. Inspect `test/unit/app/gui/drivers/` for existing coverage of
      the castle branch. If covered, delete the bullet (I5). Otherwise replace the obsolete
      "entangled with the real mapper" rationale with the true remaining reason. No tests are
      added in this batch.
- [x] io.go entry: the Linux path is fixture-tested (test/unit/internal/helpers/io/ supplies a
      temp HOME, `libraryfolders.vdf` and the Proton directory, including success). Narrow the
      limitation to what is still host-dependent: the Windows registry and user-profile branches,
      plus anything else the inspection confirms. Mention `io_other.go`.
- [x] integration_common entry (I6): `integration_test`-tagged helpers need `editor.Window` and
      a GPU context; the untagged `snapshot/` subpackage (comparer, difference, masker, store)
      has unit tests under `test/unit/test/test_helpers/integration_common/snapshot/`. Drop
      `tabCalibration.go` and the flat file names.
- [x] Allocation entry: scope "never run in CI" to the benchmark, which needs both tags while the
      CI performance job uses only `integration_test`. State that GUI integration tests do run on
      PRs (Mesa/Xvfb job) and that the figures are historical (2026-08-31). Keep the deliberate
      absence of an automated allocation threshold.
- [x] Disambiguate "Batch I Phase 4/6" (2026-08) as the earlier backlog's Batch I, not this one.
- [x] Label the remaining dated coverage figures (75.1→74.9, 6911/9246→9261) as historical
      where not already dated.

### Verification Plan
- `Select-String .agent/backlog/test_observations.md -Pattern 'suggestDirectory|tabCalibration|snapshotComparer|NewUIState\(handler, false\)|private\s+`addButtonSemantics`'` → no hits.
- Every repo path named in the file exists (throwaway check of backtick/plain paths that look
  like `app/…`, `internal/…`, `test/…`).

### Phase Summary
Done 2026-10-04. Every listed item was applied. Decisions taken from inspection:
- **stateFiles bullet:** `suggestDirectory` was simply dropped, not renamed, because
  `getWorkingDirectory` is called from the public `Load`/`SaveTo`. Its one untested branch is
  already documented by the separate `currentPath != ""` bullet.
- **`reapplyManualEdits`:** no unit test covers the castle branch (no `ReapplyCastleSettings`,
  `ReapplyWithCastleChanges` or `DecideManualEditReapplication` in `test/unit/app/gui/drivers/`).
  It is nonetheless reachable through the public `Generate` with the mocked
  `IRegenerationHandler` and `IGuiHandler`, so the bullet now says it is not a limitation and
  the missing unit test is recorded, not assigned. No test was added (I9).
- **io.go:** Windows user-profile glob tests run only on Windows, and VDF/Proton tests only
  elsewhere; each skips on the other platform. What stays host-dependent is the registry
  lookup and the Windows Steam-path fallbacks. `io_other.go` was confirmed as the
  non-Windows no-op.
- **CI claim:** checked against `pr-validation.yml`. `run-gui-integration-tests` runs on
  `pull_request` under `xvfb-run` with Mesa and `-tags=integration_test,gui`, while the
  performance step uses `-tags integration_test` only.

Verification: the stale-term grep (plus "entangled with the real") has no hits. A path
existence check over 39 path-like tokens found 37 real paths, all existing; the other 2 are
symbol references in untouched text (`app/gui/drivers.State`,
`internal/validators/editorStateValidator`).

## Phase 4: Go comments and AGENTS.md
Status: Complete

- [x] [templateHandlerMock.go](../../test/test_helpers/templateHandlerMock.go#L16):
      `interfaces.IBackend` → `handler_interfaces.IGuiHandler`.
- [x] [logButtonPositions_test.go](../../test/unit/app/gui/utils/buttonPositionLogger/logButtonPositions_test.go#L227):
      `widgets.addButtonSemantics` → `utils.AddButtonSemantics`.
- [x] AGENTS.md §4.2.1: the bullet's examples `IDialog`, `IPanel`, `IGuiHandler`; the code
      block shows the real `IGuiHandler` with its six embeds.

### Verification Plan
- `gofmt -l` on both Go files → empty; `go vet ./test/test_helpers/ ./test/unit/app/gui/utils/buttonPositionLogger/` → clean.
- `Select-String -Path AGENTS.md,test/test_helpers/templateHandlerMock.go -Pattern 'IBackend'` → no hits.
- AGENTS.md block matches [guiHandlerInterface.go](../../internal/handlers/handler_interfaces/guiHandlerInterface.go) exactly.

### Phase Summary
Done 2026-10-04. Two one-line comment edits and the AGENTS.md §4.2.1 example now name the
real `IGuiHandler` with its six embeds, in the order of `guiHandlerInterface.go`.
`TemplateHandlerMock` satisfies `IGuiHandler`: it is passed as `NewUIState`'s first argument.

Verification:
- `gofmt -l` is empty and `go vet` on both packages exits 0.
- No `IBackend` remains in AGENTS.md, README, QUICKSTART or the mock.
- The AGENTS.md link check reports only the two pre-existing literal `[path](path)`
  placeholders in §6, which are formatting examples, not links.

## Phase 5: Review and close-out
Status: Complete

- [x] Independent implementation review (GPT-6.1 Sol): factual accuracy of every changed
      sentence against source; apply or record findings.
- [x] Update `.agent/session-carry-forward.md`.
- [x] **Gate: owner commit.** Only after the owner confirms the commit (the agent never
      commits): mark §7.1 and §7.2 `✅ FIXED` in place in the review, complete §9 row I, and
      refresh the progress line (23 fixed / 9 remaining: 0 High, 6 Medium, 3 Low). Never
      renumber; §8 untouched.

### Verification Plan
- Review verdict recorded; `git status --short` lists only the intended files.

### Phase Summary
Partial, 2026-10-04. The implementation review (GPT-6.1 Sol) returned APPROVE WITH CHANGES
with 7 low findings, all applied:
1. README shows the PNG on the `SaveTemplate` branch (`templateHandler` calls
   `CreatePreviewImage`), and QUICKSTART's `IPreviewHandler` row says "preview panel layout".
2. `builders/` is described as "template model builders" (they return `template_model` types).
3. QUICKSTART's preview-panel paragraph no longer says "from your Steam install".
4. The test-observations `handleSaveState` bullet explains that `Save` calls it directly once
   `currentPath` is set by a dialog callback.
5. The `reapplyManualEdits` bullet requires at least one castle change flag.
6. The README composition root names all three injectors and drops the false
   "constructed exactly once" claim.
7. This plan's Phase 3 status and a duplicate summary placeholder were fixed.

After the fixes, the README/QUICKSTART links are re-checked OK and the greps are clean. The
handoff link check found the deleted Batch H plan link; it is now plain text. The temp link
checker is deleted. `git status --short` lists exactly the 8 intended files.

**Gate passed 2026-10-04:** the owner reviewed and committed everything as `dd2b8bf`
("Batch I", pushed). Before committing, the owner removed `TemplateHandlerMock`'s doc comment
entirely and shortened three comments in `logButtonPositions_test.go`, including the one this
batch edited. The review now marks §7.1 and §7.2 `✅ FIXED`, row I is complete, and the
progress line reads 23 fixed / 9 remaining (0 High, 6 Medium, 3 Low; recounted from the
headings).

## Final Recap
Batch I brought the user-facing docs and the test-observation registry in line with the
code, with no behaviour change:
- **README and QUICKSTART:** every audited inaccuracy is fixed. The most consequential were
  a false working-directory fallback that contradicted the output-folder hard rule, a
  non-compiling QUICKSTART example that wrote templates to `"."`, a non-existent `IBackend`
  and a pipeline diagram ending in an entity.
- **test_observations.md:** reconciled with the current tests and CI.
- **AGENTS.md:** the interface example names the real `IGuiHandler`.

Verification was docs-only by owner decision (links, a one-off snippet compile, vet). The
coverage baseline is unchanged at 74.5% (6915 / 9263).

## Deployment Plan
Nothing to deploy: documentation and comments only, committed and pushed by the owner as
`dd2b8bf`. No migration, Wire regeneration, dependency or output-path change.
