# Batch I — `EditorStateDto` rework (backlog §2.1, folding in §1.5)

Split the 72-field god-DTO into behaviour-free entity groups under a model that
owns its behaviour, make `EditorStateDto` a thin versioned persistence shell,
switch the runtime type from the DTO to the model, and stop the panels
deep-cloning the whole state every frame.

Backlog items closed by this work: [§2.1](../backlog/backlog-opus5.md) and
[§1.5](../backlog/backlog-opus5.md).

## For Future Agents

> **⚠ Read the doctrine sections newest-first: [§0.7](#07-owner-rework-2026-08-25--the-shape-that-actually-shipped),
> then [§0.6](#06-owner-decisions-2026-08-24-round-3--superseded-in-part-by-07),
> then [§0.5](#05-owner-decisions-2026-08-22-round-2--these-settle-phases-812),
> then [§0.4](#04-layering-doctrine--entity--model--dto).** Each supersedes the
> ones after it. In one line each:
>
> - **§0.7 (2026-08-25):** the **Model owns the structure**; `EditorStateDto` is
>   just `struct { editor_state_model.EditorState }`. *"Redefinition is expected
>   in Models, but it should never happen in DTOs."*
> - **§0.6 (2026-08-24):** decisions 4, 5, 6 still stand (Phase 11 folded in;
>   regeneration split; `ValidateEditorState` stays on the Model). Decisions
>   1, 2, 3, 7 are **void**.
> - **§0.5 (2026-08-22):** `app/` **may** hold a Model; only the *crossing* into
>   `internal/` must be a DTO.
> - **§0.4:** the original Entity/Model/DTO table — read it last and only for
>   background.
>
> On 2026-08-22 the owner corrected the premise this plan was built on:
> `EditorStateDto` is **not** the persisted `.gen.json` shape — that is an
> **Entity's** job. §0 rows "DTO shape" and "Wire shape", the whole of Phase 5,
> and Phase 4's grep gate are **superseded**.
>
> **State as of 2026-08-31:** Phases 1–4, 6 and 8–10 are **complete**; 1–4 and
> 8–10 are committed through `2fc6b13`, **Phase 6 is uncommitted in the working
> tree** awaiting the owner's review. Phases 5 and 11 are **superseded — do not
> execute**. Remaining order: **12 → 7**.

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

Read [AGENTS.md](../../AGENTS.md) first. Hard rules that bite in this batch:
never touch `data/`, `internal/entities/template/` or `internal/registry/`;
never stage and never commit; delete with `Remove-Item`, never `git rm`;
never run a bulk in-place rewrite over the repository; chain PowerShell with `;`,
never `&&`; unit coverage must not drop below **72.5 %** (currently **73.9 %**).

> **Doc paths moved (2026-08-25).** Plans now live in `.agent/plans/` and the
> backlog/review docs in `.agent/backlog/`. That makes this file two levels deep,
> so every repo-relative link here is `../../`, and backlog links are `../backlog/`.

---

## 0. Design decisions (settled with the owner 2026-08-21 — do not relitigate)

> ⚠ **Historical.** Several rows below were overtaken by §0.5 and §0.7 — the
> "DTO shape", "Wire shape", "Runtime type" and "Content defaults" rows in
> particular. Read §0.7 for what actually shipped; this table is kept so the
> reasoning behind the earlier rounds is not lost.

| Decision | Choice |
| --- | --- |
| How groups attach | **Anonymous embedding** of the 9 entity groups into `EditorStateModel`. Field promotion means `state.MapSize` keeps compiling at ~1,000 call sites, and `encoding/json` flattens anonymous embedded structs for free. |
| json tags | Live on the **entity leaf fields**, same as `internal/entities/template/` already does for `.rmg.json`. |
| Entity behaviour | Entities are **strictly behaviour-free value groups**. All logic lives on the model. Per AGENTS.md §4.6 pure data structs need no unit tests. |
| Runtime type | The **model**, everywhere. `EditorStateDto` appears only at the load/save boundary (`internal/repositories`, `internal/services/file_service`). |
| DTO shape | `EditorStateDto { SchemaVersion int; EditorState editor_state_model.EditorStateModel }` — a **named** field, not embedded. |
| Wire shape | Stays **flat**. `MarshalJSON`/`UnmarshalJSON` on the DTO hoist the model's fields using an internal anonymous-embedding alias, so `schemaVersion` sits as a sibling of the existing 72 keys. |
| `schemaVersion` | Always written as `1`. On load, `0` is normalised to `1` through an explicit migration hook. |
| Entity package | `internal/entities/editor_state/` (package `editor_state`). The owner decided on 2026-08-11 and re-confirmed on 2026-08-21 to put these under `internal/entities/` rather than `internal/models/`. AGENTS.md §4.4 carries this carve-out explicitly since Phase 3 — do not "fix" it. |
| Model package | `internal/models/editor_state_model/` (package `editor_state_model`) — renamed from `editor_state` to avoid a package-name collision with the entities. |
| DTO package | The three `editorState*Dto.go` files move into the existing `internal/dtos/editor_state_dto/`, alongside its current internals. |
| Content defaults | `DefaultPlayerZoneContentRows` and its helpers move to a new `internal/common/common_zone_contents/`, with `Get`/`get` prefixes, mirroring [neutralZoneProfile.go](../../internal/common/common_zones/neutralZoneProfile.go). |
| §1.5 | **In scope**, as its own phase, using option 3 (per-panel view structs) in `app/gui/models/`. Acceptance: any measurable allocation improvement, recorded here. |
| Round-trip guard | A committed golden `.gen.json` with every field populated, under `test/test_helpers/testdata/`. |

### 0.1 Field → group map (72 fields, 9 groups)

Verify the count is exactly 72 against
[editorStateDto.go](../../internal/dtos/editorStateDto.go#L15-L94) before writing
any struct; a prior automated census reported 76 and was wrong.

| Group (file) | Fields |
| --- | --- |
| `templateIdentity.go` (2) | `TemplateName`, `GameMode` |
| `mapSettings.go` (2) | `MapSize`, `ExperimentalMapSizes` |
| `playerSettings.go` (4) | `PlayerCount`, `HeroCountMin`, `HeroCountMax`, `HeroCountIncrement` |
| `neutralZoneSettings.go` (11) | `NeutralZoneCount`, `SpawnAbandonedOutposts`, `AbandonedOutpostCount`, `NeutralLowestNoCastleCount`, `NeutralLowestCastleCount`, `NeutralLowNoCastleCount`, `NeutralLowCastleCount`, `NeutralMediumNoCastleCount`, `NeutralMediumCastleCount`, `NeutralHighNoCastleCount`, `NeutralHighCastleCount` |
| `castleSettings.go` (10) | `AdvancedMode`, `PlayerOwnedCastles`, `PlayerZoneCastles`, `NeutralZoneCastles`, `HubZoneCastles`, `NeutralLowestCastlesPerZone`, `NeutralLowCastlesPerZone`, `NeutralMediumCastlesPerZone`, `NeutralHighCastlesPerZone`, `MatchPlayerCastleFactions` |
| `generationSettings.go` (15) | `PlayerZoneSize`, `NeutralZoneSize`, `HubZoneSize`, `GuardRandomization`, `Topology`, `RandomPortals`, `MaxPortalConnections`, `SpawnRemoteFootholds`, `RemoteFootholdCount`, `GenerateRoads`, `NoDirectPlayerConn`, `ResourceDensityPercent`, `StructureDensityPercent`, `NeutralStackStrengthPercent`, `BorderGuardStrengthPercent` |
| `gameRuleSettings.go` (16) | `VictoryCondition`, `FactionLawXpPercent`, `AstrologyXpPercent`, `LostStartCity`, `LostStartCityDay`, `LostStartHero`, `CityHold`, `CityHoldDays`, `GladiatorArena`, `GladiatorArenaDaysDelayStart`, `GladiatorArenaCountDay`, `Tournament`, `TournamentFirstTournamentDay`, `TournamentInterval`, `TournamentPointsToWin`, `TournamentSaveArmy` |
| `contentSettings.go` (10) | `BannedItems`, `BannedMagics`, `ValueOverridesText`, `Bonuses`, `PlayerZoneContentRows`, `LowestNeutralContentRows`, `LowNeutralContentRows`, `MediumNeutralContentRows`, `HighNeutralContentRows`, `HubZoneContentRows` |
| `manualEditSettings.go` (2) | `ManualZones`, `ManualConnections` |

Owner rulings folded in: `AdvancedMode` sits with `castleSettings` because it
gates castle diffing; all three zone sizes sit in `generationSettings`; there is
**no** `hubZoneSettings` group — `HubZoneSize` → `generationSettings`,
`HubZoneCastles` → `castleSettings`.

### 0.2 Phase ordering rationale

The build stays green at every phase boundary. Phase 3 has the DTO embed the
model **anonymously and temporarily**, so field promotion survives and no
consumer's *field accesses* change; its promoted **methods** do change shape,
which is why Phase 3 also keeps DTO-signature shims (hazard 1). Only after
Phase 4 moves every consumer onto the model does Phase 5 give the DTO its named
field, at which point nothing reads DTO fields any more.

### 0.3 Known hazards

Items 1–6 were found by an independent review of this plan's first draft and
had each falsified something it originally claimed. Treat them as load-bearing.

1. **Promoted methods change signatures, so Phase 3 needs shims.** With the model
   embedded, `dto.Clone()` resolves to the model's method and returns an
   `EditorStateModel`, and `EqualsIgnoringManualEdits`/`LayoutDefiningOptionsChanged`/
   `DiffCastleSettings` start demanding `*EditorStateModel`. Today's callers assign
   the clone to a DTO and pass DTO pointers —
   [editorState.go](../../app/gui/models/editorState.go#L26) and
   [regenerationDecisionService.go](../../internal/services/editor/regenerationDecisionService.go#L29).
   Phase 3 must therefore keep **explicit DTO-signature shim methods** that shadow
   the promoted ones and delegate inward; Phase 4 deletes them.
2. **`omitempty` already conflates nil and empty — do not "restore" a
   distinction that was never persisted.** All six content-row slices and both
   manual-edit slices carry `omitempty`
   ([editorStateDto.go](../../internal/dtos/editorStateDto.go#L85-L94)), so nil and
   `[]` both serialise to an absent key, and an absent key then keeps whatever
   `NewDefaultEditorStateDto` seeded. `EqualsIgnoringManualEdits` distinguishes
   nil from empty **in memory only**. Batch I must preserve this behaviour
   exactly and characterise it in a test — it must not try to make the
   distinction survive the file.
3. **JSON key order will change, so byte-diffing is not a valid gate.**
   `encoding/json` emits fields in declaration/embedding order, and regrouping
   reorders them. The format stays flat with the same key set; only the textual
   order moves. Every round-trip assertion must compare **parsed objects**, never
   bytes. `.gen.json` is only ever read by `json.Unmarshal`, so the reordering is
   cosmetic — but the owner's existing files will be rewritten in the new order
   on their next save.
4. **Phase 1 is not a pure package-clause move.**
   [editorStateDto.go](../../internal/dtos/editorStateDto.go#L4) imports
   `editor_state_dto` and qualifies `editor_state_dto.ManualZoneSave` /
   `CastleSettingChanges`; once it lives *in* that package those qualifiers must
   be stripped. Four sibling DTOs in the parent package also embed or carry
   `EditorStateDto` and need the new import —
   [regenerationDecisionRequestDto.go](../../internal/dtos/regenerationDecisionRequestDto.go#L13),
   [templateUpdateDto.go](../../internal/dtos/templateUpdateDto.go#L13),
   [castleSettingsReapplyRequestDto.go](../../internal/dtos/castleSettingsReapplyRequestDto.go#L11)
   and `manualEditDecisionDto.go`.
5. **Carrier DTOs are part of the runtime swap.** Those same four, plus
   `EditorStateSaveDto` and `EditorStateValidationDto`, hold editor state inside
   another DTO. Phase 4 changes what they carry; it is not only an interface
   signature change. `stateHandler` forwards the very value that `file_service`
   persists ([stateHandler.go](../../internal/handlers/stateHandler.go#L27)), so the
   model→DTO wrap on save and DTO→model unwrap on load need a named owner.
6. **The two reflection guards need *different* treatment.**
   [clone_test.go](../../test/unit/internal/dtos/editorStateDto/clone_test.go#L120-L132)
   already walks recursively through `assertNoSharedStorage` and keeps working
   unchanged. [equalsIgnoringManualEdits_test.go](../../test/unit/internal/dtos/editorStateDto/equalsIgnoringManualEdits_test.go#L211-L230)
   enumerates **top-level fields only** and mutates them by index — after
   grouping it will fail loudly on a struct kind rather than fail silently. It
   must be rewritten to enumerate recursive leaf paths, ignore only the two
   manual-edit leaves, exclude `SchemaVersion`, and mutate the model rather than
   the DTO shell.
7. ~~**Named literal.**~~ **Falsified in Phase 3.** The plan expected
   `NewDefaultEditorStateDto`'s ~32 named fields to become a nested literal. Go
   1.27 accepts **promoted fields as composite-literal keys**, so the literal
   stays flat — and `modernize`'s `embedlit` rule rewrites the nested form back
   to the flat one, so flat is also what the linter demands. Eliding only the
   type (`Group: {…}`) is invalid Go.
8. **Load seeds defaults.** [editorStateRepository.go](../../internal/repositories/editorStateRepository.go#L26-L29)
   unmarshals *over* `NewDefaultEditorStateDto()`, which is how absent keys keep
   their defaults. `UnmarshalJSON` in Phase 5 must preserve that merge semantic —
   it must not zero the receiver first.
9. **`MarshalJSON` recursion.** The alias type inside `MarshalJSON` must be a
   locally declared type that embeds the model — not the DTO — or the marshaller
   calls itself forever.
10. **Import cycle — the shared types cannot stay in `editor_state_dto`.**
    Found while starting Phase 2; the plan's first draft missed it entirely.
    `ManualZoneSave`, `ManualConnectionSave` and `CastleSettingChanges` live in
    `internal/dtos/editor_state_dto/`, but Phase 2's `manualEditSettings.go`
    must hold the first two and Phase 3's model must return the third. Combined
    with the settled DTO shape (`editor_state_dto` → `editor_state_model`) and
    the model → entities edge, that closes a cycle:

    ```
    editor_state_dto ──embeds──▶ editor_state_model ──embeds──▶ entities/editor_state
           ▲                                                             │
           └──────────── ManualZoneSave / CastleSettingChanges ◀─────────┘
    ```

    Only the third edge is removable, so those three types **must** move out.
    **Owner ruling (2026-08-21):** split them by behaviour.
    - The two pure structs (fields + json tags only, no methods) go to
      `internal/entities/editor_state/`, keeping entities behaviour-free per §0.
    - `CastleSettingChanges` (it has `Any()`) and two wrapper models embedding
      the pure structs — `ManualZoneSaveModel`, `ManualConnectionSaveModel`,
      carrying `Clone()` — go to `internal/models/editor_state_model/`.
    - The four `To…`/`From…` converters keep returning the **entity** slice, so
      no assignment site needs an unwrap loop. The six private deep-clone
      helpers for `entities.Zone`/`entities.Connection` stay private in
      `editor_state_model`.

    This runs as **Phase 2 step 0**, before the group extraction, and partly
    undoes Phase 1's consolidation — that is expected, not a regression.
    `internal/entities` must remain the base layer: it imports nothing from
    `models`, `dtos` or `helpers`, and `entities/editor_state` must not either.

---

## 0.4 Layering doctrine — Entity / Model / DTO

**Owner ruling, 2026-08-22. This supersedes any contrary statement earlier in
this plan.** It is the standard enterprise three-layer split used in C# and Java
projects; see [Baeldung — Entities vs DTOs](https://www.baeldung.com/java-entity-vs-dto)
and [DevsDaily — Entity, Model, ViewModel and DTO in C#](https://devsdaily.com/understanding-entity-model-viewmodel-and-dto-in-csharp/).

| Layer | Type | Declared in | May be used by | Logic |
| --- | --- | --- | --- | --- |
| **Database** | Entity | `internal/entities/` | `internal/repositories/` and `internal/mappers/` only | none beyond (de)serialisation — i.e. json tags |
| **Service** | Model | `internal/models/` | `internal/services/`, `internal/validators/`, `internal/mappers/`, `internal/repositories/` | **all business logic lives here** |
| **Consumer** | DTO | `internal/dtos/` | `internal/handlers/`, and in `app/` **only at handler call sites** | none |

The rules, stated exactly:

- An **Entity** is literally the boundary to external storage — files on disk,
  rows in a database. It has no logic beyond serialising to and from its storage
  format (here, JSON).
- An Entity **may be embedded in, or held as a field of, a Model**. In that case
  it is the *Model* that is being used, not the Entity, and that is allowed. An
  Entity must never be passed around by itself above the repository.
- A **Model** carries the business logic, is used in services, and **never leaves
  the backend**. `app/` must not see a Model.
- A **DTO** is the contract between the frontend (`app/gui/`) and the backend
  (`internal/`). It has no innate logic and exists only to cross the API
  boundary — which, in this project, is `internal/handlers/`. A DTO in
  `internal/` may therefore be used **only** in `internal/handlers/`, and in
  `app/` only at the handler call site.
- **Dependency direction is DTO → Model → Entity and never the reverse.** An
  Entity that imports `internal/models` is a defect, not a shortcut.

### 0.4.1 Request/response flow

The consuming UI creates a **request DTO**; each layer converts it down, and the
answer travels back up. Conversion happens at exactly two seams — the handler
(DTO ⇄ Model) and the repository (Model ⇄ Entity).

```mermaid
sequenceDiagram
    participant GUI as app/gui
    participant H as internal/handlers
    participant S as internal/services
    participant R as internal/repositories
    participant Disk as .gen.json / .rmg.json

    GUI->>H: request DTO
    Note over H: maps DTO → Model
    H->>S: Model
    S->>R: Model
    Note over R: save: maps Model → Entity
    R->>Disk: Entity
    Disk-->>R: Entity
    Note over R: load: maps Entity → Model
    R-->>S: Model
    S-->>H: Model
    Note over H: maps Model → DTO
    H-->>GUI: response DTO
```

- **Save.** The GUI builds a request DTO → the handler maps it to a Model → the
  service works on the Model → the repository maps the Model to an Entity and
  writes it.
- **Load.** The repository reads an Entity and **returns a Model** to the service
  → the service returns the Model to the handler → the handler maps it to a
  response DTO → the GUI renders the DTO.

### 0.4.2 What this batch got wrong

Phases 1–4 were built on the premise that `EditorStateDto` is the persisted
`.gen.json` shape. **It is not, and was never meant to be** — that is an Entity's
job, and correcting it is the whole point of the rework. Four concrete defects
follow from the mistake:

1. **The persisted shape is a DTO.** `internal/repositories/editorStateRepository.go`
   is typed `IFileRepository[editor_state_dto.EditorStateDto]`, so the storage
   boundary speaks the consumer type. §0's "DTO shape" and "Wire shape" rows and
   the whole of **Phase 5** are built on this and are superseded.
2. **Phase 4 pushed `app/` onto the Model.** It is backwards: the GUI must hold
   the DTO. Phase 4's own grep gate ("`EditorStateDto` only in
   `editor_state_dto`, `repositories`, `file_service`") encodes the inverted
   rule and must be replaced by the §0.4 table.
3. **The dependency arrow is reversed in the entity layer.**
   [contentSettings.go](../../internal/entities/editor_state/contentSettings.go)
   imports `internal/models` and `internal/models/config`;
   [generationSettings.go](../../internal/entities/editor_state/generationSettings.go)
   imports `internal/models/config`. Entities depend on Models today.
4. **AGENTS.md §4.4 contradicts the doctrine.** It currently permits packages
   outside `internal/` to use `internal/models/` and `internal/entities/` for
   "data typing". Under §0.4 only `internal/dtos/` is permitted there.

**Scope note — the breach is repository-wide, not Batch I's alone.** 36 files
under `app/` import `internal/models` or `internal/entities` (Phase 4 added 8 of
them; the other 28 predate this batch). Fixing all of them is a separate batch.
Phases 8–12 below correct **editor state** and put a mechanical guard in place;
the rest is recorded as a backlog item in Phase 12.

### 0.4.3 Editor-state flows in the target state (after Phases 8–12)

These are the concrete shapes of §0.4.1 for editor state, as they must look once
the correction phases are done. Solid arrows are the request path, dotted arrows
the response path; the label on each arrow is **the type that crosses it**.

Invariants visible in all three — use them as the review checklist:

- There are exactly **two conversion seams**: `internal/handlers` (DTO ⇄ Model)
  and `internal/repositories` (Model ⇄ Entity). Nowhere else converts.
- `app/` never names a Model or an Entity. `internal/repositories` never names a
  DTO. `internal/services` and `internal/validators` never name a DTO.
- The Entity exists only inside the repository; it never crosses that subgraph's
  border.

#### Save

```mermaid
flowchart TD
    subgraph app["app/gui — consumer layer"]
        A1["panel.SaveToState()<br/>mutates the working EditorStateDto"]
        A2["drivers.State.handleSaveState(path)"]
        A3["EditorStateSaveDto{ State, OutputPath }"]
    end

    subgraph handlers["internal/handlers — API boundary"]
        H1["IStatePersistenceHandler.SaveState(EditorStateSaveDto)"]
        H2["mappers.EditorStateMapper.ToModel(dto)"]
    end

    subgraph services["internal/services — service layer"]
        S1["file_service.SaveSettings(path, *EditorState)"]
    end

    subgraph repos["internal/repositories — database layer"]
        R1["EditorStateRepository.Save(dir, name, EditorState)"]
        R2["toEntity(model)<br/>EditorStateEntity{ SchemaVersion: 1, 9 groups }"]
        R3["json.Marshal → atomicFileWriter"]
    end

    Disk[("&lt;TemplateName&gt;.gen.json")]

    A1 --> A2 --> A3
    A3 -->|"DTO"| H1 --> H2
    H2 -->|"Model"| S1
    S1 -->|"Model"| R1 --> R2
    R2 -->|"Entity"| R3 --> Disk
    Disk -.->|"savedPath string"| A2
```

#### Load

```mermaid
flowchart TD
    subgraph app["app/gui — consumer layer"]
        A1["drivers.State.handleLoadState(path)"]
        A2["models.EditorState.OverrideState(EditorStateDto)"]
        A3["panel.LoadFromState()"]
    end

    subgraph handlers["internal/handlers — API boundary"]
        H1["IStatePersistenceHandler.LoadState(path, fixIssues)"]
        H2["mappers.EditorStateMapper.ToDto(model)"]
        H3["envelope DTO<br/>{ State EditorStateDto, Warnings []string }"]
    end

    subgraph services["internal/services — service layer"]
        S1["file_service.LoadSettingsFile(path)"]
        V1["validators.Validate(*EditorState) → []ValidationIssue"]
        V2["issue.Fix(*EditorState) when fixIssues"]
    end

    subgraph repos["internal/repositories — database layer"]
        R1["EditorStateRepository.Load(path)"]
        R2["json.Unmarshal into a defaults-seeded EditorStateEntity<br/>merge into the receiver, never zero it"]
        R3["migrate SchemaVersion 0 → 1"]
        R4["toModel(entity)"]
    end

    Disk[("&lt;path&gt;.gen.json")]

    A1 --> H1 --> S1 --> R1 --> R2
    Disk --> R2
    R2 --> R3 --> R4
    R4 -.->|"Model"| S1
    S1 -.->|"Model"| H1
    H1 --> V1 --> V2
    V2 -.->|"Model"| H2 --> H3
    H3 -.->|"DTO"| A1 --> A2 --> A3
```

#### Update (a single edited setting)

Runs on the frame path, so it is the flow most sensitive to the Phase 11
decision — see the note under the diagram.

```mermaid
flowchart TD
    subgraph app["app/gui — consumer layer"]
        A1["widget changed → panel.SaveToState()"]
        A2["drivers.State.UpdateState(func(*EditorStateDto))"]
        A3["models.EditorState.UpdateCurrentState<br/>copy current DTO, apply the mutation"]
        A4["store returned DTO as current"]
        A5["WasStateChanged() → unsaved / confirmExit"]
    end

    subgraph handlers["internal/handlers — API boundary"]
        H1["IStateValidationHandler.ValidateEditorState(EditorStateDto, fixIssues)"]
        H2["mappers.EditorStateMapper.ToModel(dto)"]
        H3["mappers.EditorStateMapper.ToDto(model)"]
        H4["EditorStateValidationDto{ State, Warnings }"]
        H5["IStateComparisonHandler.EqualsIgnoringManualEdits(previous, current)"]
    end

    subgraph services["internal/services — service layer"]
        V1["EditorStateValidator.Validate(*EditorState)"]
        V2["issue.Fix(*EditorState) when fixIssues"]
        V3["normalizeInactiveNeutralCounts(*EditorState)"]
        M1["EditorState.EqualsIgnoringManualEdits(*EditorState)"]
    end

    A1 --> A2 --> A3
    A3 -->|"DTO"| H1 --> H2
    H2 -->|"Model"| V1 --> V2 --> V3
    V3 -.->|"Model"| H3 --> H4
    H4 -.->|"DTO"| A4 --> A5
    A5 -->|"DTO"| H5 --> M1
    M1 -.->|"bool"| A5
```

Notes on the update flow, all of which are **open until the owner answers Phase
11's decision**:

- `H5` / `M1` are drawn as a handler round trip because
  `EqualsIgnoringManualEdits` is Model logic and `app/` may not call it directly.
  The same applies to `LayoutDefiningOptionsChanged`, `DiffCastleSettings` and
  `HasManualEdits`. `IStateComparisonHandler` is an **invented name** — no such
  interface exists yet.
- The "copy current DTO" step in `A3` is drawn as staying inside `app/`. That is
  Phase 11 option (b): copying a logic-free struct is view-layer plumbing, not
  business logic. Under option (a) it becomes another handler call and the
  diagram gains a second round trip per edit.
- Every arrow crossing the `app` ⇄ `handlers` border on this path costs a
  DTO ⇄ Model conversion **per frame**, which is exactly what Phase 6 is trying
  to reduce. Phase 11 must be decided with Phase 6's benchmark in view, not
  separately.

---

## 0.5 Owner decisions, 2026-08-22 (round 2) — these settle Phases 8–12

Answers to the four blocking questions in the carry-forward, plus the follow-ups
they raised. **These override §0.4 where they conflict, and §0.4.3's diagrams are
partly invalidated by decision 1 — see the note under it.**

### 0.5.1 `app/` keeps the Model; the DTO is the *crossing*, not the storage

**AGENTS.md §4.4 is authoritative over plan §0.4 on this point.** The owner's
committed amendment (`a264d0b`) states it directly: *"if `app/` would be detached
from the project, it would have a separate model implementation, but here
`app/.../` can use the `internal/models/` instead of duplicating code"*, and its
sequence diagram has the GUI **storing a Model**, mapping *stored Model → DTO* on
the request and *DTO → Model and stores it* on the response.

So the rule is **not** "`app/` never sees a Model". It is:

> `app/` may **hold** a Model as its working state, but every **crossing** into
> `internal/` is a DTO. The handler boundary speaks DTO in both directions.

Consequences, all load-bearing:

- **Phase 4 is correct as committed.** `app/gui/models.EditorState` keeps
  `*editor_state_model.EditorState` for `current`/`previous`/`next`.
- **Phase 11 shrinks drastically.** It no longer retypes the panels or moves the
  five state methods behind handlers. It becomes: *the handler call sites pass
  DTOs*. Rewritten below.
- **Phase 12 drops the rule `internal/models/** must not be imported by app/**`.**
  The entity and DTO rules stand.
- **§0.4's line "A Model never leaves the backend. `app/` must not see one" is
  superseded.** Left in place for the record; read it together with this section.
- The 28-file "residue" the carry-forward flagged (`config`, `neutral_zone`,
  `preview`, `entities.Zone`) is **no longer a breach** for the model half. Only
  the `internal/entities` imports under `app/` remain questionable, and they are
  still out of scope for Batch I.

### 0.5.2 Phase 8 — the entity layer gets its own leaf types

The owner's ruling, verbatim in intent: move `MapTopology`, `BonusEntry`,
`ZoneContentRowSave` (renamed **`ZoneContentRow`**) and `ContentRuleRowSave`
(renamed **`ContentRuleRow`**) into `internal/entities/editor_state/`. **Do not
move attached methods** — behaviour that came off a type becomes a **function**.

- `BonusPresetType` is dragged along because it is `BonusEntry.PresetType`'s
  type. Its constants move with it.
- Behaviour lifted off those types becomes plain functions in
  `internal/helpers/`, following the existing
  [road_helpers/roadType.go](../../internal/helpers/road_helpers/roadType.go) /
  [zone_helpers/zoneNameType.go](../../internal/helpers/zone_helpers/zoneNameType.go)
  convention — package per domain, value as the first argument:

  | Was | Becomes |
  | --- | --- |
  | `BonusEntry.GetHash() string` | `config_helpers.GetHash(BonusEntry) string` |
  | `BonusPresetType.String() string` | `config_helpers.GetString(BonusPresetType) string` |
  | `BonusPresetType.IsResource() bool` | `config_helpers.IsResource(BonusPresetType) bool` |
  | `ZoneContentRowSave.Clone()` / `.Normalized()` | `editor_state_helpers.CloneZoneContentRow(...)` / `NormalizeZoneContentRow(...)` |
  | `ContentRuleRowSave.Clone()` | `editor_state_helpers.CloneContentRuleRow(...)` |

- **No wrapper structs.** Wrappers were considered and rejected: the entity
  groups hold `[]ZoneContentRow` and `[]BonusEntry` as *slices*, and Go slices do
  not interconvert, so a wrapper type would force wrap/unwrap at ~85 call sites
  (`ContentSettings.Bonuses`, `GeneratorConfig.Bonuses`, `bonusesPanel.bonuses`,
  the picker callback, `cloneContentRows`, and the six content-row fields).
  Functions over entity values keep one type end to end and put the logic in the
  same layer a wrapper would have.

> ⚠ **Reversed by §0.7.** The owner chose wrapper structs after all: every model
> group and leaf type now wraps its entity, `Clone()` / `Normalized()` are back
> as **methods** on the model rows (delegating to the `editor_state_helpers`
> functions below, which survive and still operate on entity values), and the
> two groups whose slices carry model elements — `ContentSettings` and
> `ManualEditSettings` — are re-declared with explicit converters. The wrap/unwrap
> cost the paragraph above predicted was paid, once, in Phase 10.
>
> `MapTopology` also ended up in **`internal/entities/topology/`**, not in
> `internal/entities/editor_state/` as planned here.

**Alias façades absorb the churn.** The old spellings keep working because their
declaring packages re-export the entity types:

- [internal/models/config/types.go](../../internal/models/config/types.go) —
  `MapTopology`, `BonusEntry`, `BonusPresetType` and all their constants
  re-point at `editor_state.*`. **182 files import `config`; none of them
  change.** *(Still true, except `MapTopology` now aliases
  `topology.MapTopology` and `BonusEntry`/`BonusPresetType` alias the
  `editor_state_model` wrappers rather than the entities.)*
- ~~`internal/models/` re-exports `ZoneContentRow` / `ContentRuleRow`.~~
  **Deleted by §0.7.6.** `internal/models/types.go` and
  `internal/models/editor_state_model/types.go` are gone; the 56 files that used
  those aliases now name `editor_state_model.ZoneContentRow` /
  `.ContentRuleRow` directly.

This is deliberate and symmetrical with the decision below: an alias names the
*same* type, so it is not a layering violation — `models.ZoneContentRow` **is**
`editor_state.ZoneContentRow`. It also keeps the entity package's name out of 77
service/app files, which is *more* §0.4-compliant than importing it there. The
façades are deletable later, one caller at a time, as their own batch.

> ⚠ **The row-type façade was deleted in Phase 10, not "later".** Only the
> `config` façade survives.

### 0.5.3 `internal/entities/types.go` → `internal/entities/template/types.go`

Requested by the owner, but the file is `package entities`, so moving it renames
the package for **426 importers** (282 `test/unit`, 85 `internal/services`, 15
`test/test_helpers`, 11 `app/gui`, 11 `internal/dtos`, …). A 426-file identifier
rename inside Batch I would collide with AGENTS.md §2.6 and bury the layering
work.

**Chosen: the zero-churn form.** `internal/entities/template/types.go` becomes
the canonical alias home, and `internal/entities/types.go` shrinks to re-alias
it. Aliases of aliases name the identical type, so nothing else recompiles and
both spellings stay valid. The root façade can then be retired incrementally in
a dedicated batch.

> **§2.1 note.** This adds a file to the read-only
> `internal/entities/template/`. The owner gave explicit approval for exactly
> this file. **No existing file in that directory is modified**, so the
> `.rmg.json` schema and game compatibility are untouched.

### 0.5.4 Phase 9 — Entity ⇄ Model mapping lives in `internal/mappers/`

`internal/mappers/editorStateEntityMapper.go`, per AGENTS.md §4.4 ("data mappers
and converters"). This needs a **carve-out to §0.4**, which says Entities are used
only in `internal/repositories/`: the permitted namers of an Entity are
`internal/repositories/`, `internal/models/`, `internal/entities/`,
`internal/helpers/*_helpers/` and now `internal/mappers/`. Phase 12's checker
must encode that list, not §0.4's shorter one.

### 0.5.5 Phase 6 / Phase 11 frame path — deferred, decide on numbers

Deliberately **not** answered up front. Phase 6's allocation benchmark runs
first and the measurement decides. Under decision 0.5.1 the pressure largely
evaporates anyway — `app/` keeps the Model, so no per-frame DTO ⇄ Model
conversion is introduced on the update path, and §0.4.3's update flowchart
overstates the cost.

---

## Phase 1: Baseline guard + DTO package move
Status: Complete   <!-- Not started | In progress | Complete -->

Freeze today's on-disk format as a committed artefact **before any code changes**,
then perform the pure package move on its own so the noisy import diff never
mixes with a behaviour change.

- [x] Confirm the field count is exactly 72 and record any discrepancy here.
      **72 confirmed, no discrepancy** — counted by hand on
      [editorStateDto.go](../../internal/dtos/editorStateDto.go) (24 + 36 + 4 + 6 + 2)
      and re-confirmed mechanically by the fixture's top-level key count assertion.
- [x] Record the **before** coverage number (AGENTS.md §2.3) so every later phase
      has a baseline to compare against. **Baseline = 73.6 %** on a clean tree.
      This supersedes the 72.8 % quoted in the carry-forward (the Go 1.27 batch
      raised it); the 72.5 % floor still governs.
- [x] Write a throwaway program (or a temporary test) that builds an
      `EditorStateDto` with **every** field set to a distinctive non-zero value,
      including all six content-row slices, bonuses, manual zones and manual
      connections, and marshal it with the **current** code.
      Builder kept permanently as
      [allFieldsEditorState.go](../../test/test_helpers/allFieldsEditorState.go)
      (`NewAllFieldsEditorStateDto`); the throwaway `cmd/tmpfixturegen` generator
      wrote the fixture through the **real** `EditorStateRepository.Save` and was
      then deleted.
- [x] Add the result as `test/test_helpers/testdata/editorState_v0_flat.gen.json`
      — a new **unstaged** file; do not stage it and do not commit it. Name it
      `_v0_` because it has no `schemaVersion`: it is the legacy-shape fixture
      that Phase 5's migration hook must keep loading.
- [x] Add a golden round-trip test that loads the fixture and asserts every one of
      the 72 fields survives with its distinctive value. Compare **parsed values,
      never bytes** (hazard 3). This test must stay green, unchanged in intent,
      through every later phase.
      Added as
      [editorStateWireFormat_integration_test.go](../../test/integration/editorStateWireFormat_integration_test.go)
      — **untagged** (production APIs only, no `*_testexports.go`, no GPU), so it
      runs in a plain `go test ./test/...`. Three gates: fixture unmarshals into a
      **zero** DTO and equals the builder; the fixture has exactly 72 top-level
      keys; marshalling the builder equals the fixture **as parsed maps**.
- [x] Move [editorStateDto.go](../../internal/dtos/editorStateDto.go),
      [editorStateSaveDto.go](../../internal/dtos/editorStateSaveDto.go) and
      [editorStateValidationDto.go](../../internal/dtos/editorStateValidationDto.go)
      into `internal/dtos/editor_state_dto/`.
- [x] Strip the now-invalid self-qualifiers inside the moved file
      (`editor_state_dto.ManualZoneSave` → `ManualZoneSave`, likewise
      `CastleSettingChanges`) and drop the self-import (hazard 4).
      Four sites, exactly as surveyed.
- [x] Add the new import to the **three** sibling DTOs left in `internal/dtos/`
      that carry editor state: `regenerationDecisionRequestDto.go`,
      `templateUpdateDto.go`, `castleSettingsReapplyRequestDto.go`.
      **Plan correction:** `manualEditDecisionDto.go` needs **no** change — it
      only references `*editor_state_dto.CastleSettingChanges` and already
      imports that package. §0.3 hazard 4 lists it in error. The compiler
      confirmed this: after the move it reported exactly these three files.
- [x] Update imports and qualifiers (`dtos.EditorStateDto` →
      `editor_state_dto.EditorStateDto`) across the 12 production packages and the
      test tree. Do this file by file — **no bulk in-place rewrite** (AGENTS.md §2.6).
      **104 files, 509 sites.** Owner approved the mechanism: substitution of the
      five fixed, unambiguous tokens over an **explicit enumerated file list**
      (§2.6's own "explicit list from `gofmt -l`" pattern), never a repo-wide sweep.
- [x] Move the mirrored unit-test folders to
      `test/unit/internal/dtos/editor_state_dto/…` per AGENTS.md §4.6.
      `test/unit/internal/dtos/editorStateDto/` →
      `test/unit/internal/dtos/editor_state_dto/editorStateDto/` (7 files).

### Verification Plan
- `go build ./...` clean; `go vet ./...` and
  `go vet -tags='integration_test,gui' ./...` clean.
- `go run ./cmd/testlayoutcheck .` prints `test-layout check passed`.
- `go test ./test/unit/... -count=1` and
  `go test -tags=integration_test ./test/integration/... -count=1` pass.
- The new golden test passes.
- Coverage report re-run; total ≥ the number recorded above and ≥ 72.5 %.
- Diff review: apart from the fixture and its test, every change must be an
  import, a package clause, a qualifier, or the self-qualifier strip — **no
  behaviour change**.

### Phase Summary
**Complete.** All verification gates green:

| Gate | Result |
| --- | --- |
| `go build ./...` | clean |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | clean |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test ./test/... -count=1` (untagged, incl. the new golden test) | exit 0 |
| `go test -tags=integration_test ./test/integration/... -count=1` | exit 0 |
| `go test ./test/unit/... -count=1` | exit 0, 164 packages |
| Coverage | **73.7 %** (baseline 73.6 %, floor 72.5 %) |
| `gofmt -l .` | empty |
| `golangci-lint-v2 run ./...` | 50 issues, **all pre-existing `funcorder`**; down from 71 |

What exists now:

- `internal/dtos/editor_state_dto/` owns `editorStateDto.go`, `editorStateSaveDto.go`
  and `editorStateValidationDto.go` alongside the three types that were already
  there. `internal/dtos/` keeps only the carrier DTOs, which now qualify.
- The wire format is frozen by
  [editorState_v0_flat.gen.json](../../test/test_helpers/testdata/editorState_v0_flat.gen.json),
  written through the **real** `EditorStateRepository.Save`, and guarded by three
  untagged tests in
  [editorStateWireFormat_integration_test.go](../../test/integration/editorStateWireFormat_integration_test.go).
  The builder [allFieldsEditorState.go](../../test/test_helpers/allFieldsEditorState.go)
  is permanent — Phase 5 regenerates the `_v1_` fixture from the same builder.

Decisions and gotchas for whoever picks this up:

- **The golden test is untagged on purpose.** It touches production APIs only, so
  it runs in a plain `go test ./test/...`. Do not add `integration_test` to it.
- **Zero-value unmarshal is deliberate.** The fixture is loaded into a *zero*
  `EditorStateDto`, not over `NewDefaultEditorStateDto()`. That is what proves the
  file alone carries all 72 values; loading over defaults would let a dropped
  field pass by inheriting its default. Do not "fix" this to match
  `editorStateRepository.Load`.
- **Every bool in the fixture is `true`** for the same reason — a `false` would
  pass vacuously against a zero struct.
- **Comparisons are parsed objects, never bytes** (hazard 3), so Phase 2–3
  regrouping may reorder keys freely without touching this test.
- The 72-key count is now asserted mechanically, so Phase 2's field→group split
  cannot silently lose or add a key.
- Three files — [state.go](../../app/gui/drivers/state.go),
  [editorState.go](../../app/gui/models/editorState.go),
  [layoutPanelZones.go](../../app/gui/panels/layoutPanelZones.go) — carry extra
  reflow in the diff. They were **already** non-`gofmt`-clean at HEAD (single-line
  function bodies), and the longer `editor_state_dto.` qualifier pushed those
  lines past the `golines` limit, so a reflow was unavoidable. It was applied by
  the project's own *"Go: Run Linter"* `--fix` task, not by an ad-hoc `gofmt -w`.
- `coverage.txt`, `coverage.html` and `lcov.info` are regenerated tracked
  artefacts, expected in the diff.
- Nothing is staged. Nothing is committed.

---

## Phase 2: Extract the entity groups
Status: Complete

- [x] **Step 0 — break the import cycle (hazard 10).** Move `ManualZoneSave` and
      `ManualConnectionSave` (pure structs) to `internal/entities/editor_state/`;
      move `CastleSettingChanges`, the two wrapper models and the four converters
      to `internal/models/editor_state_model/`. Update consumers and move the
      three mirrored unit-test folders. Must land before the groups are written.
- [x] Create `internal/entities/editor_state/` with the 9 files from §0.1, one
      struct per file, camelCase file names (AGENTS.md §4.1).
- [x] Move each field **with its existing json tag verbatim**. Do not rename, do
      not reorder within a group, do not change a type.
- [x] Create `internal/common/common_zone_contents/` and move
      `DefaultPlayerZoneContentRows` → `GetDefaultPlayerZoneContentRows`, with the
      four private helpers renamed to a `get` prefix.
- [x] Move its unit test to the mirrored
      `test/unit/internal/common/common_zone_contents/…` folder.
- [x] Leave `EditorStateDto` untouched this phase — the entities exist but nothing
      uses them yet, so the phase is purely additive.

### Verification Plan
- `go build ./...` clean, full unit + integration suites pass unchanged.
- `go run ./cmd/testlayoutcheck .` passes.
- Manual check: the 9 structs together declare exactly 72 tagged fields, and the
  tag set is byte-identical to the DTO's.
- Coverage ≥ 72.5 % and not below the Phase 1 baseline.
- **Owner gate:** the §0.1 field→group table is signed off before Phase 3 starts.

### Phase Summary

**Step 0 — the import cycle (hazard 10).** Writing the groups would have closed
`editor_state_dto → editor_state_model → entities/editor_state → editor_state_dto`.
The first edge is a settled §0 decision and the second is the point of the batch,
so only the third was removable. Owner ruling: behaviour-free structs
(`ManualZoneSave`, `ManualConnectionSave`) go to `internal/entities/editor_state/`;
`CastleSettingChanges`, the two wrapper models (`ManualZoneSaveModel`,
`ManualConnectionSaveModel`, anonymous embeds with `Clone()`) and the four
converters go to `internal/models/editor_state_model/`. Converters return the
**entity** slice; the six deep-clone helpers stay **private** in
`editor_state_model`. This partly undoes Phase 1's consolidation, and it makes the
layering rule load-bearing: **`internal/entities` must import nothing from
`models`/`dtos`/`helpers`.**

**Groups.** `internal/entities/editor_state/` now holds the 9 structs from §0.1
(`templateIdentity` 2, `mapSettings` 2, `playerSettings` 4, `neutralZoneSettings`
11, `castleSettings` 10, `generationSettings` 15, `gameRuleSettings` 16,
`contentSettings` 10, `manualEditSettings` 2) — **72 tagged fields**, every tag
verbatim.

**Tag verification.** The plan called for comparing the tag set against the DTO at
HEAD, but Phase 1 is now committed (`586fc18 "Batch I wip"`) and that path no
longer exists. Compared instead against the frozen wire-format fixture
`test/test_helpers/testdata/editorState_v0_flat.gen.json` — a stronger check:
entity tag count 72, fixture key count 72, sets equal, 0 duplicates.

**Content rows.** `GetDefaultPlayerZoneContentRows` (+ four `get`-prefixed private
helpers and `const guardedRuleName`) moved to
`internal/common/common_zone_contents/`; both callers (`editorStateDto.go`,
`app/gui/dialogs/zoneContentDialog.go`) updated; the unit test moved to the
mirrored folder. `editorStateDto.go` net `14 +/112 -`.

**Gates.** `go build` clean; `go vet -tags='integration_test,gui'` clean;
`gofmt -l .` empty; `testlayoutcheck` passed; unit, default (frozen golden test
included) and tagged-integration suites all exit 0 — re-run after the lint pass.
Coverage **73.6 %** (floor 72.5 %). Lint `--fix` took 85 issues → **21, all
pre-existing `funcorder`**; `git status --short` counts were identical before and
after, confirming no drive-by files. Overall diff: 32 files, +85/−1007.

**Gotchas worth remembering.** (a) `embedlit`/modernize collapses
`Model{Embedded: X{...}}` to `Model{{...}}`, which orphaned the
`entities/editor_state` import in both `clone_test.go` files. (b) PowerShell
`-replace '…' , '$1' + $type` parses `$1ManualZoneSave` as a *group name*, so the
replacement silently no-ops — use a separate blanket pass. (c) The plan's hazard 4
is also inaccurate: it lists `manualEditDecisionDto.go` as a Phase 1 change that
was never needed.

**Next:** Phase 3 is blocked on the owner gate — the §0.1 field→group table must be
signed off first.

---

## Phase 3: Introduce the model, DTO delegates
Status: Complete

- [x] Create `internal/models/editor_state_model/editorStateModel.go` —
      `EditorStateModel` anonymously embedding the 9 entity structs.
- [x] Move the behaviour onto it: `Clone`, `LayoutDefiningOptionsChanged`,
      `DiffCastleSettings`, `EqualsIgnoringManualEdits`, `HasManualEdits`, the
      private comparison helpers, and a `NewDefaultEditorStateModel` factory
      carrying the current defaults.
- [x] Reduce `EditorStateDto` to `struct { editor_state_model.EditorStateModel }`
      — **anonymous, temporarily** — so every `state.MapSize` selector and the flat
      json output are unchanged and the whole repository still compiles untouched.
- [x] Keep **DTO-signature shim methods** on the DTO that shadow the promoted
      ones (hazard 1): `Clone() EditorStateDto`,
      `EqualsIgnoringManualEdits(*EditorStateDto) bool`,
      `LayoutDefiningOptionsChanged(*EditorStateDto) bool`,
      `DiffCastleSettings(*EditorStateDto) CastleSettingChanges`. Each delegates
      to the embedded model. Without these, Phase 3 does **not** build.
- [x] `NewDefaultEditorStateDto` delegates to `NewDefaultEditorStateModel`.
- [x] Move the DTO behaviour tests to
      `test/unit/internal/models/editor_state_model/editorStateModel/<method>_test.go`,
      one file per public method (AGENTS.md §4.6).
- [x] Leave the recursive clone guard alone; rewrite the equality guard to walk
      recursive **leaf paths** instead of top-level fields (hazard 6), and add a
      case proving it still trips when a field is added to a nested group but
      omitted from the comparison.

### Verification Plan
- `go build ./...` clean; full unit + integration suites pass **with no test
  changes other than the moves and the equality guard rewrite**.
- Phase 1's golden test passes untouched — this is the proof the key set is
  intact. Expect the key **order** to change; that is allowed (hazard 3).
- Coverage ≥ 72.5 % and not below the Phase 1 baseline.

### Phase Summary

**Complete.** The owner signed off the §0.1 field→group table before this phase
started, as the Phase 2 gate required.

What exists now:

- [editorStateModel.go](../../internal/models/editor_state_model/editorState.go)
  holds `EditorStateModel` — the nine entity groups embedded anonymously — plus
  every method that used to live on the DTO (`NewDefaultEditorStateModel`,
  `Clone`, `LayoutDefiningOptionsChanged`, `DiffCastleSettings`,
  `EqualsIgnoringManualEdits`, `HasManualEdits`, the three scalar-comparison
  helpers and the content-row/pointer helpers). Bodies were moved **verbatim**;
  only the receiver and parameter types changed.
- [editorStateDto.go](../../internal/dtos/editor_state_dto/editorStateDto.go) is now
  41 lines: the embedded model, `NewDefaultEditorStateDto` delegating to
  `NewDefaultEditorStateModel`, and the four DTO-signature shims. `HasManualEdits`
  needs no shim — its signature does not mention the state type, so the promoted
  method is used as-is.

**Hazard 7 is obsolete on Go 1.27 — and the linter enforces the flat form.**
The plan assumed `NewDefaultEditorStateDto` would have to become a *nested*
literal. Go 1.27 accepts **promoted fields as keys in a composite literal**, so
`EditorStateModel{TemplateName: …, MapSize: …}` compiles unchanged, and
`modernize`'s `embedlit` rule actively rewrites the nested form back to the flat
one. Both `NewDefaultEditorStateModel` and
[allFieldsEditorState.go](../../test/test_helpers/allFieldsEditorState.go) therefore
kept their original flat literals. Note the inverse trap: `Group: {…}` (eliding
only the type) is **not** valid Go and does not compile — it is either the full
`Group: pkg.Group{…}` or the flat promoted keys.

**Test moves.** The five behaviour test files moved to
`test/unit/internal/models/editor_state_model/editorStateModel/`, with
`newDefaultEditorStateDto_test.go` renamed to `newDefaultEditorStateModel_test.go`.
The DTO folder was repopulated with **thin delegation tests** for the four shims
and the factory, so §4.6's one-file-per-public-method rule still holds while the
shims exist; Phase 4 deletes both the shims and those files.

**Equality drift guard (hazard 6), as rewritten.** It now enumerates recursive
**leaf paths** (`"MapSettings.MapSize"`), descending only through *anonymous*
fields, and mutates through `FieldByIndex`. `SchemaVersion` needs no exclusion:
the guard walks the model, and Phase 5 puts the version on the DTO shell. The
"proof it still trips" case is
`TestWhenTheLeafWalkIsBuilt_ReachesEveryFieldOfEveryGroup`, a single
`assert.Equal` on a group→field-count map — it fails both when the walk stops at
a group (nested fields unreached, so an uncompared one could hide) and when any
group gains or loses a field. It is the §0.1 table asserted mechanically.

**Unexpected side effect worth knowing.** Two `//nolint:govet` directives in
[getCurrentState_test.go](../../test/unit/app/gui/models/editorState/getCurrentState_test.go)
became *unused* and were removed: `govet`'s `unusedwrite` analyzer no longer
fires on `copyOfState.TemplateName = …` now that the field is promoted through
two levels of embedding.

**Gates.**

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | exit 0 |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test -count=1 ./test/unit/...` | exit 0 |
| `go test -count=1 ./test/...` (incl. the frozen golden test, untouched) | exit 0 |
| `go test -tags=integration_test -count=1 ./test/integration/...` | exit 0 |
| Coverage | **73.7 %** (Phase 1 baseline 73.6 %, floor 72.5 %) |
| `gofmt -l .` | empty |
| `golangci-lint-v2 run ./...` | **0 issues** |

The GPU suite was not run: no rendering code changed. Phase 4's plan requires it.
Nothing was staged and nothing was committed.

**Also amended:** AGENTS.md §4.4 now carries the editor-state carve-out
explicitly (open question 2 from the carry-forward), so the
`internal/entities/editor_state/` placement no longer contradicts the table.

---

## Phase 4: Swap the runtime type to the model
Status: Complete

Pure type-name substitution. Because fields stay promoted, **no field-access site
changes** — only signatures and variable types.

**Owner rulings (2026-08-22), settled before the phase started:**

- `EditorStateSaveDto` and `EditorStateValidationDto` **stay in
  `internal/dtos/editor_state_dto/`**. The package is not "the persistence
  shell" — it is the DTO package for `EditorState(Model)` operations (load,
  save, validate). Only their `State` field type changes. Do not propose moving
  them to `internal/dtos/` again.
- The model↔DTO conversion gets a **named API** on the DTO:
  `NewEditorStateDto(model)` and `(*EditorStateDto).Model()`. `file_service`
  calls those, so Phase 5's reshape to a named field rewrites those two bodies
  and nothing else. This is hazard 5's "named owner".

**Plan correction:** `ManualEditDecisionDto` needs **no** change. Phase 2 moved
`CastleSettingChanges` to `editor_state_model`, so it already carries no editor
state. The carrier list below is over-inclusive by one — the same class of error
as §0.3 hazard 4's `manualEditDecisionDto.go` entry.

- [x] `app/gui/models.EditorState`: `current`/`previous`/`next` become
      `*editor_state_model.EditorStateModel`; update `GetCurrentState`,
      `OverrideState`, `SetNextState`, `UpdateCurrentState`, `GetPreviousState`,
      `GetNextState`.
- [x] `internal/handlers/handler_interfaces`: `IStateValidationHandler` and the
      other DTO-carrying signatures take/return the model.
- [x] `internal/validators`: `ValidationIssue.Fix(*EditorStateModel)`; the 45
      validated fields and 46 issue categories are unchanged.
- [x] `internal/mappers/generatorConfigMapper.go`: `FromEditorState` takes the
      model. All 69 field reads stay as written.
- [x] `app/gui/drivers`, `app/gui/editor`, `app/gui/panels`,
      `internal/services/editor`: update the types they pass through.
- [x] `internal/repositories` and `internal/services/file_service` keep the
      **DTO** — they are the load/save boundary. Name the owner of the conversion
      explicitly: `file_service` wraps model→DTO on save and unwraps DTO→model on
      load, and nothing above it sees a DTO (hazard 5).
- [x] Convert the **carrier DTOs** that hold editor state inside another DTO to
      carry the model: `EditorStateSaveDto`, `EditorStateValidationDto`,
      `RegenerationDecisionRequestDto`, `TemplateUpdateDto`,
      `CastleSettingsReapplyRequestDto`, ~~`ManualEditDecisionDto`~~. This is the
      item most likely to be under-estimated — it is not just interface signatures.
- [x] Delete the Phase 3 shim methods once every caller is on the model.
- [x] Update every mock under [test/test_helpers](../../test/test_helpers) and the
      affected unit/integration tests.

### Verification Plan
- `go build ./...`, both `go vet` variants, `testlayoutcheck` all clean.
- Full unit + integration suites pass.
- `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1`
  passes **without `-update`** — the GPU suite proves the GUI still behaves.
- Phase 1's golden test still passes.
- Round-trip non-aliasing test: load a file, mutate the returned model, reload,
  assert the second read is unaffected.
- `grep` check: `EditorStateDto` appears only in `internal/dtos/editor_state_dto`,
  `internal/repositories` and `internal/services/file_service`.
- Coverage ≥ 72.5 % and not below the Phase 1 baseline.

### Phase Summary

**Complete.** 99 files, **+530/−628** — the phase is net-negative because the four
Phase 3 shims and their five delegation tests went away.

**The boundary now has a name.** `file_service` is the only thing that converts:
[fileService.go](../../internal/services/file_service/fileService.go) unwraps with
`dto.Model()` on load and wraps with `editor_state_dto.NewEditorStateDto(*model)`
on save. Both helpers live in
[editorStateDto.go](../../internal/dtos/editor_state_dto/editorStateDto.go) and are
the *only* two places that know how the DTO stores the model, so Phase 5's
reshape to a named field rewrites those two bodies and nothing else.
`IFileService` is model-facing on both methods; `repositories` still speaks
`IFileRepository[EditorStateDto]`.

**Grep gate passes.** `EditorStateDto` now appears in production code only in
`internal/dtos/editor_state_dto` (8), `internal/repositories` (6) and
`internal/services/file_service` (3). One stale comment in
[zoneContentRowSave.go](../../internal/models/zoneContentRowSave.go) was retargeted
at `EditorStateModel`; it was the only other hit.

**Two plan items were wrong and are corrected above.** `ManualEditDecisionDto`
needed no change — Phase 2 had already moved `CastleSettingChanges` to
`editor_state_model`, so it carries no editor state. And the owner ruled that
`EditorStateSaveDto` / `EditorStateValidationDto` **stay** in `editor_state_dto`:
they are the DTOs for `EditorState(Model)` *operations*, not for the persisted
object, so the package is not "the persistence shell" and must not be treated as
one.

**Mechanism.** The type swap was a substitution of two fixed, unambiguous tokens
(`editor_state_dto.EditorStateDto` → `editor_state_model.EditorStateModel`,
`…NewDefaultEditorStateDto` → `…NewDefaultEditorStateModel`) over an **explicit
enumerated file list** — the same owner-approved mechanism as Phase 1, never a
repo-wide sweep. Five paths were deliberately **excluded** because they must keep
the DTO: `test_helpers/allFieldsEditorState.go`,
`test/integration/editorStateWireFormat_integration_test.go`,
`test/unit/…/editorStateRepository/`, `test/unit/…/editorStateDto/` and
`test/unit/…/fileService/`. The last of those is genuinely **mixed** — the
repository mock is typed on the DTO while the service API is model-facing — and
was rewritten by hand.

**New tests.**
[editorStateRoundTrip_integration_test.go](../../test/integration/editorStateRoundTrip_integration_test.go)
is the non-aliasing gate the verification plan asked for: save through the real
`InitializeGuiHandler()` graph, then prove that mutating a returned model
(scalar *and* slice element) leaves a reload untouched, that mutating the source
model after a save does not reach the file, and that a full round-trip returns
the all-fields model unchanged. It is **untagged** — production APIs only, no GPU
— so it runs in a plain `go test ./test/...`, matching the Phase 1 golden test.
`NewEditorStateDto` and `Model` got the two unit files §4.6 requires.

**Naming cleanup, deliberately partial.** Parameters literally named `stateDto`
that now carry a *model* were renamed to `state` in production
(`guiHandler`, `stateHandler`, `templateHandler`, `templateHandlerInterface`) and
in `TemplateHandlerMock`; `stateFiles.handleLoadState`'s `dto` became `loaded`.
**~55 test-local `dto` / `stateDto` identifiers were left alone on purpose**: in
most of those closures the enclosing scope already binds `state` to the driver
`State`, so a blind rename would shadow it. Renaming them is cosmetic and is
recorded in [test_observations.md](../backlog/test_observations.md) rather than
smuggled into this diff.

**Gates.**

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | exit 0 |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test -count=1 ./test/unit/...` | exit 0 |
| `go test -count=1 ./test/...` (incl. the frozen golden test, untouched) | exit 0 |
| `go test -tags=integration_test -count=1 ./test/integration/...` | exit 0 |
| `go test -tags='integration_test,gui' -count=1 ./test/integration/gui/...` | exit 0, **no `-update`** |
| Grep gate | DTO confined to the three permitted packages |
| Coverage | **73.7 %** (Phase 1 baseline 73.6 %, floor 72.5 %) |
| `gofmt -l .` | empty |
| `golangci-lint-v2 run ./...` | **0 issues** |

`coverage.txt` / `coverage.html` / `lcov.info` regenerate identically and do not
dirty the tree. Nothing was staged and nothing was committed.

**For Phase 5.** The frozen golden test still passes **untouched**, which is the
proof the wire format did not move. Phase 5 only has to: reshape the DTO to
`{ SchemaVersion int; EditorState EditorStateModel }`, rewrite the bodies of
`NewEditorStateDto` / `Model` / `NewDefaultEditorStateDto`, add
`MarshalJSON`/`UnmarshalJSON`, and point the repository at the new field.
`test_helpers.NewAllFieldsEditorStateModel()` already exists for the `_v1_`
fixture.

---

## Phase 5: Versioned persistence shell
Status: **Superseded — do not execute.** See §0.4 and Phase 9.

> This phase puts `SchemaVersion`, `MarshalJSON` and `UnmarshalJSON` on the
> **DTO**, i.e. it makes the consumer contract the file format. Under the §0.4
> layering doctrine that is an **Entity's** job. Every item below is still
> wanted — versioning, the migration hook, the merge-into-receiver unmarshal, the
> `_v0_`/`_v1_` fixtures, the `omitempty` characterisation — but it belongs on
> `EditorStateEntity`. **Phase 9 carries all of it forward**; hazards 8 and 9
> apply there unchanged. Kept here only so the reasoning behind each item is not
> lost.

- [ ] Reshape the DTO to
      `struct { SchemaVersion int; EditorState editor_state_model.EditorStateModel }`
      — named field, no embedding.
- [ ] Add `MarshalJSON` using a **locally declared** alias type that anonymously
      embeds the model (never the DTO — hazard 9) so the output stays flat with
      `schemaVersion` as a sibling key. Always write `1`.
- [ ] Add `UnmarshalJSON` that **merges into the existing receiver** (preserving
      the seed-defaults-then-overlay semantic of
      [editorStateRepository.go](../../internal/repositories/editorStateRepository.go#L26-L29)),
      then calls an explicit migration hook normalising `SchemaVersion == 0` → `1`.
- [ ] Point `internal/repositories` and `internal/services/file_service` at the
      new field.
- [ ] Generate the current-writer golden as `editorState_v1_flat.gen.json` and
      keep **both** unstaged fixtures: `_v0_` proves legacy files still load,
      `_v1_` proves the current writer's output round-trips.
- [ ] Add tests: legacy fixture loads and lands at version 1; the two fixtures
      have the same key set apart from `schemaVersion` (**compare parsed objects,
      not bytes**); an absent content-row key still yields the seeded default.
- [ ] Add a characterisation test pinning today's `omitempty` behaviour — nil and
      empty both serialise to an absent key and both reload as the seeded default
      (hazard 2). This documents the limitation rather than pretending to fix it.

### Verification Plan
- All gates from Phase 4, plus: marshal the default state and compare the
  **parsed object** against the pre-batch output — the only permitted delta is
  the added `schemaVersion` key. Key order may differ.
- Load each real `.gen.json` in [output/](../../output/) and assert no error and a
  non-default template name. These are the owner's own legacy files and are the
  realest migration evidence available.
- Coverage ≥ 72.5 % and not below the Phase 1 baseline.

### Phase Summary
_(write when phase completes)_

---

## Phase 6: Stop the per-frame whole-state clone (backlog §1.5)
Status: **Complete** (2026-08-31).

> **Refreshed 2026-08-26.** The two numbers this phase used to quote are both
> stale: the batch-D figures below come from a different tree *and* different
> hardware, and Phase 10's `9,543 allocs/op` measurement describes the §0.6
> design that was subsequently deleted. **Measure the current tree first and
> treat that as the only baseline.** The last known-good commit to measure from
> is `2fc6b13`.

Historical, for context only — batch D on `BenchmarkEditorWindow_TabCycling`:
**3.01 M ns/op · 1,435 KB/op · 6,640 allocs/op** (pre-batch-D was 4,676 allocs/op).
Acceptance: any measurable improvement, recorded below.

- [x] Re-measure the baseline on this machine first — every earlier number in
      this plan is from a superseded tree.
- [x] ~~Add read-only per-panel view structs in `app/gui/models/`, each carrying
      only what its panel renders.~~ **Not built — the profile said no.** See the
      Phase Summary: no panel reads the whole state per frame except one click
      handler, which now reads it on the click instead.
- [x] Convert the five cloning `Layout` paths — of the five, only
      `handleZoneContentDialogClicks` in
      [layoutPanelZones.go](../../app/gui/panels/layoutPanelZones.go) was actually
      per-frame; the four `LoadFromState` sites are event-driven.
- [x] Keep the clone on any path that hands state **out** of the model — the
      §1.1 aliasing fix must not be undone.
- [x] Unit tests for each new view struct's projection — replaced by tests for
      `linq.SelectSlice`, the projection helper that carries the whole win.

**The measurement this phase exists to take.** `ValidateEditorState` was kept on
the Model precisely so this phase could price the per-frame path without a DTO
conversion in the way (§0.6 decision 6). The hot path is
`app/gui/models.EditorState.UpdateCurrentState`, which deep-`Clone()`s and then
calls the validator on every panel write.

### Verification Plan
- `go test -v -tags='integration_test,gui' -bench=BenchmarkEditorWindow_TabCycling -run=xxx ./test/performance/... -benchmem -benchtime=20x -timeout=120s`
  — record before/after ns/op, B/op, allocs/op in the Phase Summary. Run locally
  on the real GPU; never in CI.
- GPU integration suite passes without `-update`. If a golden genuinely must
  change, `-run`-scope it tightly enough not to match neighbouring tests.
- Full unit + integration suites pass; coverage ≥ 72.5 %.

### Phase Summary

**Result — `BenchmarkEditorWindow_TabCycling`, headless, `-benchtime=50x -count=6`,
baseline taken on the clean tree at `540f7eb` (= `2fc6b13` plus three doc-only
commits):**

| | allocs/op | B/op | ns/op (median of 6) |
| --- | --- | --- | --- |
| Before | 12,689–12,691 | 1,045.0 KB | 4.05 M |
| After | 4,772–4,774 | 720.0 KB | 3.57 M |
| Δ | **−62.4 %** | **−31.1 %** | **−11.8 %** |

**What the profile actually said.** The plan's bullet list was written against a
batch-D-era reading of the tree and pointed at the wrong code. An
`alloc_objects` profile (`-benchtime=200x -memprofile`) attributed **75.2 % of
every allocation in the benchmark** to `editor_state_model.EditorState.Clone`,
split across its callers: `UpdateCurrentState` 37.4 %, `stateHandler.
ValidateEditorState` 29.3 %, `GetCurrentState` 24.1 %, `GetPreviousState` 9.2 %.
Inside `Clone`, 97 % of that went into `CloneZoneContentRows`.

The cost was not the clone *count* but the clone *mechanism*. Every row-slice
clone ran `linq.FromSlice(x).Select(f).ToSlice()`, a lazy chain that allocates
its three closures and boxes `ToSlice`'s accumulator **before it looks at the
source**, then regrows the result with `append`. `EditorState.Clone` runs eight
such chains, and six of the eight are empty on a default state — including
`editor_state_helpers.CloneZoneContentRow`'s clone of the embedded entity's
`Rules`, which §0.7.4 deliberately keeps nil. Roughly a quarter of all
allocations in the whole benchmark were chains projecting nothing.

**The four changes.**

1. **`linq.SelectSlice`** ([slice.go](../../internal/helpers/linq/slice.go)) —
   the eager equivalent of `FromSlice(...).Select(...).ToSlice()`. Returns `nil`
   for an empty source without allocating and sizes the result exactly once.
   Semantics match the lazy chain exactly, including empty → `nil`.
2. **The five clone helpers on the frame path use it**: `CloneContentRuleRows` /
   `CloneZoneContentRows` in both
   [editor_state_helpers](../../internal/helpers/editor_state_helpers/) and
   [editor_state_model](../../internal/models/editor_state_model/), plus the
   `ManualZones` / `ManualConnections` chains in `EditorState.Clone`. The two
   `linq` chains there collapse to method expressions, so the per-element
   closure is gone as well.
3. **`handleZoneContentDialogClicks` stopped cloning the whole state per frame.**
   It cloned the state to read six row slices that only matter when one of six
   buttons was clicked. `openZoneContentDialog` now takes a *getter* instead of a
   row slice and reads the state itself, so the clone happens on the click. The
   `switch` and its short-circuit over `Clicked(gtx)` are untouched — a click
   left unread this frame must still be readable next frame.
4. **Nothing else.** The clones in `UpdateCurrentState`, `ValidateEditorState`,
   `GetCurrentState`, `GetPreviousState` and `AutoRegenerate` all hand state
   across a boundary, so they stay (§1.1). The `UpdateCurrentState` →
   `ValidateEditorState` pair is a genuine double deep-clone and is the largest
   remaining item, but collapsing it means letting `updateFunc` mutate through to
   the live `current`, which is exactly the aliasing the bullet above forbids.
   With each clone now ~9× cheaper the case for taking that risk is weak; it is
   left as a note, not a task.

**Why no per-panel view structs.** The bullet assumed the panels read whole state
every frame. They do not: `LoadFromState` runs from the panel constructors and
from `Window.load()` (toolbar new/open), never from `Layout`. The one genuine
per-frame reader was the click handler in change 3, and the fix there is to not
read at all rather than to read into a new type. Building six view structs would
have added types and mapping code for zero measured gain, which §3.1 forbids.

**Where the remaining allocations are.** Re-profiled after the change: the state
path is down from ~1.72 M sampled objects to ~0.20 M (−88 %), and what is left
is one `make` per content row per clone — the irreducible cost of a deep clone.
Everything above it in the profile is now Gio rendering (`GetPanelWidget`,
`newSectionHeaderWidget`, glyph shaping), which is a different backlog item.

**Gates** — all green: `go build ./...`, both `go vet` variants,
`gofmt -l ./app ./internal ./test ./cmd` empty, `go run ./cmd/testlayoutcheck .`
passed, `wire diff` exit 0, `go test ./test/unit/...`, `go test ./test/...`
(untagged), `go test -tags=integration_test ./test/integration/...`, and the GPU
suite **without `-update`** (35.0 s). Unit coverage **73.9 %**, unchanged, floor
72.5 %. `golangci-lint-v2` reports **1 issue** — the owner's `godox` TODO in
[stateHandler.go](../../internal/handlers/stateHandler.go), i.e. the current
baseline, no regression.

---

## Phase 7: Cleanup, docs and gates
Status: Not started — **runs last, after Phase 12** (see the correction phases).

- [ ] Update the architecture description in [README.md](../../README.md#L233),
      which still says the GUI collects input into `dtos.EditorStateDto`.
- [ ] Rewrite backlog [§2.1](../backlog/backlog-opus5.md) and
      [§1.5](../backlog/backlog-opus5.md) as self-contained ✅ FIXED entries; update
      the header counts and the §8 batch table row **I**.
- [ ] Update [test_observations.md](../backlog/test_observations.md) if any gap
      opened or closed.
- [ ] Fix the stale `todo/` links in **AGENTS.md §4.6** (two of them) — the docs
      moved to `.agent/backlog/` and those links now dangle. Same for the
      relative links inside `.agent/backlog/review-opus5-08-04.md`, which is now
      two levels deep and still uses single-`../` paths.
- [ ] Resolve the `godox` TODO in
      [stateHandler.go](../../internal/handlers/stateHandler.go) or accept it,
      so the lint baseline returns to **0 issues**.
- [ ] Delete this plan file once the backlog entries are self-contained
      (doc-lifecycle rule), using `Remove-Item`.
- [ ] Write `./.agent/session-carry-forward.md` per AGENTS.md §5.2.

> The former first bullet, *"remove any delegation shim left from Phases 3–5"*,
> is already done: the DTO carries no shims — it embeds the model and declares
> nothing else.

### Verification Plan
- Every gate green in one run: `go build ./...`, both `go vet` variants,
  `go run ./cmd/testlayoutcheck .`, unit, integration, GPU integration,
  `golangci-lint-v2 run ./... --issues-exit-code=0` at **0 issues**,
  `gofmt -l ./app ./internal ./test ./cmd` empty, coverage ≥ 72.5 %.
- `git status --short` reviewed and reported to the owner. **Nothing staged,
  nothing committed.**

### Phase Summary
_(write when phase completes)_

---

# Correction phases (added 2026-08-22)

Phases 8–12 exist because Phases 1–5 were designed around the premise that
`EditorStateDto` is the persisted `.gen.json` shape. **§0.4 settles that it is
not.** Read §0.4 and §0.4.2 before touching any of them.

**Revised phase order: 8 → 9 → 10 → 6 → 12 → 7.** Phase 5 is superseded by
Phase 9. Phase 6 (the per-frame clone) must run **after** Phase 11, because
Phase 11 changes what the panels hold and would otherwise invalidate Phase 6's
view structs. Phase 7 (cleanup, docs and gates) now runs **last, after Phase 12**,
not after Phase 6.

## Decision required before Phase 8 — what happens to the unstaged Phase 4 diff

**Resolved 2026-08-22: Option A.** The owner committed the Phase 4 diff as
`a264d0b "Batch I wip 4"` (111 files, +1641/−941), together with the plan edits
and the AGENTS.md §4.4 amendment. The working tree is clean and Phase 4 stands.
See [§0.5.1](#051-app-keeps-the-model-the-dto-is-the-crossing-not-the-storage) —
the amendment also makes Phase 4's *direction* correct, so Phase 11 no longer
reverses it.

The original options are kept below for the record.

| Option | Consequence |
| --- | --- |
| **A — keep it, pivot forward** | Phases 10–11 rewrite ~80 of the same files a second time. Keeps the Phase 3 shim removal and the round-trip tests without redoing them. |
| **B — revert the Phase 4 diff to `0ca3b6d`, start at Phase 8** | Avoids the double-churn in the GUI and ~70 test files. Costs a redo of the shim deletion (small) — **save [editorStateRoundTrip_integration_test.go](../../test/integration/editorStateRoundTrip_integration_test.go) and the two new DTO unit tests before reverting**, they are still wanted. |

**Recommendation: B.** At `0ca3b6d` the GUI and handlers already name
`EditorStateDto` at exactly the call sites §0.4 wants them to. Only the DTO's
*shape* is wrong there, and Phase 10 fixes the shape in one package. Phase 4
moved those call sites to the Model, so keeping it means moving them back — the
same ~80 files touched twice for no net gain. **Only the owner may revert**;
`git restore`/`git checkout` on a working tree the agent did not create is
destructive and out of bounds without an explicit instruction.

Phases 8–12 are written to be valid under either option.

---

## Phase 8: Restore the dependency direction in the entity layer
Status: Complete

Entities must import nothing from `internal/models`. Today two of the nine
groups do (§0.4.2 defect 3), so the arrow points the wrong way at the bottom of
the stack. Nothing above can be laid out correctly until this is fixed.

Design settled in [§0.5.2](#052-phase-8--the-entity-layer-gets-its-own-leaf-types)
and [§0.5.3](#053-internalentitiestypesgo--internalentitiestemplatetypesgo).
The offenders are exactly two files:
[contentSettings.go](../../internal/entities/editor_state/contentSettings.go)
(`models.ZoneContentRowSave`, `config.BonusEntry`) and
[generationSettings.go](../../internal/entities/editor_state/generationSettings.go)
(`config.MapTopology`).

**8a — template alias façade (§0.5.3).**

- [x] Add `internal/entities/template/types.go` (`package template`) holding the
      alias block currently in [internal/entities/types.go](../../internal/entities/types.go).
      Omit `RmgTemplate` — it is *declared* in that package already, so aliasing
      it there is a redeclaration.
- [x] Shrink `internal/entities/types.go` to re-alias `template.*`. Both
      spellings then name the identical type; no other file changes.

**8b — entity leaf types.** New files under `internal/entities/editor_state/`,
all behaviour-free (fields, json tags, constants — no methods):

- [x] `mapTopology.go` — `type MapTopology string` + the 11 `Topology*` constants.
- [x] `bonusPresetType.go` — `type BonusPresetType int` + the 11 `Bonus*` constants.
- [x] `bonusEntry.go` — `BonusEntry{PresetType, ReceiverFilter, Param, Param2}`.
- [x] `contentRuleRow.go` — `ContentRuleRow`, renamed from `ContentRuleRowSave`.
- [x] `zoneContentRow.go` — `ZoneContentRow`, renamed from `ZoneContentRowSave`;
      `Rules []ContentRuleRow`.
- [x] **json tags copied byte-for-byte.** A renamed *type* must not move a *key*;
      Phase 1's frozen golden proves it.

**8c — helper functions (§0.5.2 table).**

- [x] `internal/helpers/config_helpers/bonusEntry.go` — `GetHash`.
- [x] `internal/helpers/config_helpers/bonusPresetType.go` — `GetString`, `IsResource`.
- [x] `internal/helpers/editor_state_helpers/zoneContentRow.go` —
      `CloneZoneContentRow`, `CloneZoneContentRows`, `NormalizeZoneContentRow`.
- [x] `internal/helpers/editor_state_helpers/contentRuleRow.go` —
      `CloneContentRuleRow`, `CloneContentRuleRows`.

**8d — re-point the façades, delete the originals.**

- [x] [config/types.go](../../internal/models/config/types.go): `MapTopology`,
      `BonusEntry`, `BonusPresetType` and every `Bonus*` / `Topology*` constant
      re-point at `editor_state.*`.
- [x] `Remove-Item` `config_inner/mapTopology.go`, `bonusEntry.go`,
      `bonusPresetType.go`.
- [x] `internal/models/zoneContentRow.go` + `contentRuleRow.go` replace the two
      `*RowSave.go` files with alias declarations.
- [x] Semantic rename (`vscode_renameSymbol`, **not** a text sweep — §2.6)
      `ZoneContentRowSave` → `ZoneContentRow`, `ContentRuleRowSave` →
      `ContentRuleRow` across all 77 referencing files.

**8e — the two entity groups.**

- [x] `contentSettings.go` and `generationSettings.go` end with **no import
      block at all**.

**8f — method call sites** (the only behavioural edits):
`editorStateModel.go` `cloneContentRows`, `app/gui/utils/models.go`
`CloneRuleRows`, `services/bonuses/bonusEntryService.go` `GetHash`,
`bonusPickerDialog.go` and `bonusesPanel.go` `IsResource`, plus any
`.Normalized()` caller.

**8g — tests** relocate to mirror the new impl paths per §4.6: the
`config_inner/bonusEntry` and `config_inner/bonusPresetType` folders move under
`test/unit/internal/helpers/config_helpers/`, and the two `*RowSave` folders
under `test/unit/internal/helpers/editor_state_helpers/`. The json round-trip
tests move to the entity path — they characterise the persisted format and are
worth keeping despite §4.6's "pure data structs need no tests".

### Verification Plan
- No file under `internal/entities/` imports `internal/models`,
  `internal/dtos`, `internal/services`, `internal/handlers` or
  `internal/helpers`. Assert it mechanically, not by eye — this is Phase 12's
  rule, brought forward for this phase.
- Phase 1's frozen golden test passes **unchanged**: splitting a row type must
  not move a json key. Parsed objects, never bytes (hazard 3).
- `go build ./...`, both `go vet` variants, full unit + integration suites,
  coverage ≥ 72.5 %.

### Phase Summary

**Done 2026-08-23. All gates green.** The entity layer no longer depends on the
model layer: `go list` over all **9** packages under `internal/entities/`
reports **zero** imports of `internal/models`, `internal/dtos`,
`internal/services`, `internal/handlers` or `internal/helpers`. That is the
whole point of the phase and it is now mechanically true.

**Shape of the change.** Five behaviour-free leaf types were added to
`internal/entities/editor_state/` — `MapTopology`, `BonusPresetType`,
`BonusEntry`, `ContentRuleRow`, `ZoneContentRow` — with json tags copied
verbatim. The behaviour that used to hang off them became functions in two new
helper packages, per the `road_helpers`/`zone_helpers` convention:
`internal/helpers/config_helpers` (`GetHash`, `GetString`, `IsResource`) and
`internal/helpers/editor_state_helpers` (`CloneZoneContentRow(s)`,
`NormalizeZoneContentRow`, `CloneContentRuleRow(s)`).

**Churn was absorbed by alias façades, exactly as §0.5.2 planned.**
`internal/models/config/types.go` re-points `MapTopology` / `BonusEntry` /
`BonusPresetType` and all 22 constants at the entity types, so **none of the 182
files importing `config` changed**. `internal/models/{zoneContentRow,
contentRuleRow}.go` re-export the two rows, so the 77 files referencing them
only lost the `Save` suffix — applied by `vscode_renameSymbol` (290 edits, ~70
files), never by a text sweep.

**Only six production call sites changed behaviourally**, all method → function:
`editorStateModel.cloneContentRows`, `mandatoryContentItemMapper.FromRows`,
`bonusEntryService` (×2), `bonusPickerDialog`, `bonusesPanel`,
`zoneContentDialog`.

**Deleted:** `config_inner/{mapTopology,bonusEntry,bonusPresetType}.go` and
`internal/models/{zoneContentRowSave,contentRuleRowSave}.go`.

**Tests relocated** to mirror the new impl paths: the two `config_inner` folders
moved under `test/unit/internal/helpers/config_helpers/`, the two `*RowSave`
clone/normalize folders under
`test/unit/internal/helpers/editor_state_helpers/`, and the two json round-trip
files under `test/unit/internal/entities/editor_state/` — they characterise the
persisted format and are deliberately kept despite §4.6's "pure data structs
need no tests". Four emptied folders were removed.

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet -tags='integration_test,gui' ./...` | exit 0 |
| `gofmt -l .` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test ./test/unit/...` | pass |
| `go test ./test/...` (untagged) | pass |
| `go test -tags=integration_test ./test/integration/...` | pass |
| GPU suite, **no `-update`** | pass (24.0s) |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Unit coverage | **73.7 %** — unchanged, floor 72.5 % |

**8a note.** `internal/entities/template/types.go` was added as the canonical
alias home and `internal/entities/types.go` now re-aliases it; all 426 importers
were left untouched. The owner separately simplified
`internal/entities/template/rmgTemplate.go` to use the new unqualified names and
staged it. **The agent is not permitted to modify `rmgTemplate.go` or any
`template` subpackage** — that file is the owner's to change.

**Carried into Phase 12:** `.golangci.yml` has a `config-inner-private` depguard
rule but nothing yet guards the entity layer. Phase 12's checker should encode
the §0.5.4 namer list and can lean on `go list` for the import graph, as used to
verify this phase.

**Not done here, and deliberately:** `internal/models` still holds the two row
aliases and `internal/entities/types.go` still exists. Both are façades kept to
bound this diff; retiring them is its own batch.

---

## Phase 9: Make the persisted shape an Entity (carries Phase 5 forward)

Status: Complete

Every item of the superseded Phase 5 is still wanted; it moves onto the Entity.
Hazards 2, 3, 8 and 9 apply here unchanged.

- [x] Add `internal/entities/editor_state/editorStateEntity.go`:
      `EditorStateEntity { SchemaVersion int; <the 9 groups> }`. Behaviour-free
      apart from (de)serialisation, which §0.4 explicitly allows an Entity.
- [x] `MarshalJSON` via a **locally declared** alias type that embeds the groups
      — never the entity itself (hazard 9) — so the wire format stays flat with
      `schemaVersion` as a sibling key. Always write `1`.
- [x] `UnmarshalJSON` that **merges into the existing receiver**, preserving the
      seed-defaults-then-overlay semantic of
      [editorStateRepository.go](../../internal/repositories/editorStateRepository.go#L26-L29)
      (hazard 8), then an explicit migration hook normalising `0` → `1`.
- [x] Retype `EditorStateRepository`: it reads and writes the **Entity**, and
      maps Entity ⇄ Model at its own boundary — `Load` returns a Model, `Save`
      takes a Model (§0.4.1).
- [x] **Owner decision — resolved (§0.5.4):** the Entity ⇄ Model mapping lives in
      `internal/mappers/editorStateEntityMapper.go`, per AGENTS.md §4.4. §0.4's
      "Entities are used only in repositories" gains a carve-out: the permitted
      namers of an Entity are `internal/repositories/`, `internal/models/`,
      `internal/entities/`, `internal/helpers/*_helpers/` and
      `internal/mappers/`. Phase 12's checker encodes that list.
- [x] `internal/services/file_service` speaks **Models only** — delete the
      model↔DTO conversion Phase 4 put there.
- [x] Delete `NewEditorStateDto` and `(*EditorStateDto).Model()`. They exist only
      to make the DTO a persistence shell.
- [x] Regenerate the golden as `editorState_v1_flat.gen.json` and keep **both**
      unstaged fixtures: `_v0_` proves legacy files still load, `_v1_` proves the
      current writer round-trips.
- [x] Tests: the legacy fixture loads and lands at version 1; the two fixtures
      have the same key set apart from `schemaVersion` (**parsed objects, not
      bytes**); an absent content-row key still yields the seeded default; and the
      `omitempty` characterisation test pinning hazard 2.

### Verification Plan
- Load every real `.gen.json` in [output/](../../output/) and assert no error and a
  non-default template name — the owner's own legacy files are the realest
  migration evidence available.
- `EditorStateDto` no longer appears anywhere in `internal/repositories` or
  `internal/services`.
- Marshal the default state and compare the **parsed object** against the
  pre-batch output; the only permitted delta is the added `schemaVersion` key.
- Full unit + integration suites, coverage ≥ 72.5 %.

### Phase Summary

**Done 2026-08-23. All gates green.** The `.gen.json` file is now an Entity. The
storage boundary speaks `EditorStateEntity` to the disk and
`editor_state_model.EditorState` to everything above it, and nothing under
`internal/repositories` or `internal/services` names a DTO any more (verified by
a tree scan, not by eye).

**What exists now.**

- [editorStateEntity.go](../../internal/entities/editor_state/editorState.go)
  — the nine groups embedded anonymously plus `SchemaVersion`, with
  `MarshalJSON` / `UnmarshalJSON` and a private `migrateSchemaVersion` hook.
  `CurrentEditorStateSchemaVersion = 1`.
- [editorStateEntityMapper.go](../../internal/mappers/editorStateEntityMapper.go)
  + its interface — `ToEntity` / `ToModel`. Both sides carry the *same* nine
  group types, so the mapping is a regrouping, not a copy: nine field
  assignments each, and the slices stay shared with the argument. Registered in
  `InfrastructureSet`; `wire gen` re-run, `wire diff` clean.
- `NewEditorStateRepository(mapper)` returns
  `IFileRepository[editor_state_model.EditorState]`. `Load` decodes **over the
  mapped defaults**, `Save` maps then writes. `file_service` lost its
  `editor_state_dto` import entirely.
- `EditorStateDto` keeps only its embedded model and `NewDefaultEditorStateDto`.
  It is now referenced by **no production code at all** — that is expected and
  temporary; Phase 10 rebuilds it as the consumer contract.

**Hazard 9, resolved by type definition rather than by embedding.** The plan said
"a locally declared alias type that embeds the groups". A `type
editorStateFields EditorStateEntity` definition is the same thing without
duplicating the nine-field list: it has the entity's layout but **not** its
methods, so neither codec can re-enter itself. Embedding `EditorStateEntity`
instead would promote `MarshalJSON` and recurse forever.

**Hazard 8 held.** `UnmarshalJSON` seeds `editorStateFields` *from the receiver*
before decoding, so an absent key keeps whatever the caller put there. The
repository relies on that: it decodes over
`ToEntity(NewDefaultEditorStateModel())`, exactly as the old code decoded over
`NewDefaultEditorStateDto()`.

**Two lint rules changed the struct's shape — worth knowing before "fixing" it
back.**

- `embeddedstructfieldcheck` forbids a regular field before embedded ones, so
  `SchemaVersion` is declared **last** and `schemaVersion` is therefore the
  **last** key in the file, not the first. Purely cosmetic; every test compares
  parsed objects (hazard 3).
- `recvcheck` flags the value/pointer receiver mix that a JSON codec pair
  inherently has. It carries a `//nolint:recvcheck` with the reason: `MarshalJSON`
  **must** take a value so that both a state and a `*state` serialise through it
  — with a pointer receiver, `json.Marshal(entity)` would silently skip the
  version stamp.

**Fixtures.** `editorState_v1_flat.gen.json` was written through the **real**
`EditorStateRepository.Save` by a throwaway `cmd/tmpfixturegen`, which was then
deleted. Both fixtures are kept and both are untracked. `_v0_` (72 keys, no
`schemaVersion`) proves a legacy file still loads and lands at version 1; `_v1_`
(73 keys) proves the current writer round-trips.
`test_helpers.NewAllFieldsEditorStateDto` became
**`NewAllFieldsEditorStateEntity`**, built by running the all-fields model
through the real mapper — so a group the mapper forgets to carry fails the
fixture comparison instead of being mirrored by a hand-written expectation.

**Tests added.**

```text
test/unit/internal/entities/editor_state/editorStateEntity/{marshalJSON,unmarshalJSON}_test.go
test/unit/internal/mappers/editorStateEntityMapper/{newEditorStateEntityMapper,toEntity,toModel}_test.go
```

plus two repository cases covering the seeded defaults, and a rewritten
[editorStateWireFormat_integration_test.go](../../test/integration/editorStateWireFormat_integration_test.go)
(still untagged) with six cases across the two fixtures.
`test/unit/internal/dtos/editor_state_dto/editorStateDto/{model,newEditorStateDto}_test.go`
were deleted with their subjects.

**Hazard 2 characterised, not fixed.** `TestWhenContentRowsAreNil_TheKeyIsNotWritten`
and `TestWhenContentRowsAreEmpty_TheKeyIsNotWrittenEither` pin that `omitempty`
conflates the two on disk, and
`TestWhenStateFileCarriesEmptyContentRows_TheSeededDefaultRowsSurviveToo` pins
that either one reloads as the default. The nil-versus-empty distinction
`EqualsIgnoringManualEdits` draws is in-memory only and stays that way.

**Real-file migration evidence.** All three legacy files in [output/](../../output/)
— `Buggy preview`, `Custom Template`, `Custom Template2` — load without error
and keep their template names, player counts and 14 content rows. Checked with a
throwaway `cmd/tmplegacyload`, deleted afterwards.

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | exit 0 |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 |
| `go test ./test/unit/...` | exit 0 |
| `go test ./test/...` (untagged) | exit 0 |
| `go test -tags=integration_test ./test/integration/...` | exit 0 |
| GPU suite, **no `-update`** | exit 0 |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Unit coverage | **73.7 %** — unchanged, floor 72.5 % |

**Two pre-existing `gocritic unlambda` findings** in the Phase 8 helpers
(`CloneContentRuleRows`, `CloneZoneContentRows`) were cleared in passing so the
lint baseline stays at 0; both are one-line `Select(func(x) { return f(x) })` →
`Select(f)`.

Nothing was staged and nothing was committed.

---

## Phase 10: Rebuild `EditorStateDto` as a consumer contract
Status: Complete

**Rewritten 2026-08-24 (§0.6), then reworked by the owner 2026-08-25 (§0.7).**
After Phase 9 the DTO has no persistence job left. It becomes what §0.4 says it
is: the contract between `app/gui/` and `internal/handlers/`.

> ⚠ **Read §0.7 before §0.6.** The owner rebuilt the DTO layer after §0.6 was
> written; §0.7 supersedes §0.6 decisions **1, 2, 3 and 7** outright. Decisions
> **4, 5 and 6** survived and are still in force. §0.6 is kept for the record so
> the rejected shape is not proposed a third time.

### 0.7 Owner rework, 2026-08-25 — the shape that actually shipped

The owner's verdict on §0.6's design: *"you made the editor state model be a
wrapper around the entity and the dto a full structural rewrite, but it should
have been the other way around."*

**The rule, in the owner's words:** *"Redefinition is expected in Models, but it
should never happen in DTOs."*

1. **The Model is the structural layer; the DTO is thin.**
   `EditorStateDto` is now just

   ```go
   type EditorStateDto struct { editor_state_model.EditorState }
   ```

   The nine `*Dto` group structs, `CastleSettingChangesDto`, and
   `internal/mappers/editorStateMapper.go` + its interface are **deleted**.
   §0.6 decisions 1 and 2 are void. **A DTO embedding a Model is intended here**
   — do not "fix" it, and do not let Phase 12's checker flag it.
2. **Model groups embed their entity group**, one struct per file in
   `internal/models/editor_state_model/`:
   `type CastleSettings struct { editor_state.CastleSettings }`, and likewise
   for `BonusEntry`, `ContentRuleRow`, `ManualZoneSave`, `ManualConnectionSave`.
   Go 1.27 promoted-field literal keys keep
   `editor_state_model.BonusEntry{PresetType: x}` compiling, and `embedlit`
   enforces that flat form.
3. **Two groups cannot embed** — `ContentSettings` and `ManualEditSettings` —
   because their slices must carry *model* element types and Go slices do not
   interconvert. They are re-declared and get explicit converters.
4. **`ZoneContentRow` embeds the entity but shadows `Rules`** with
   `[]ContentRuleRow`. The embedded entity's `Rules` is deliberately left nil so
   there is one source of truth; `ToZoneContentRowModels` lifts the rules out and
   nils the embedded slice, `ToZoneContentRowEntities` folds them back in.
5. **Converter naming: name the layer explicitly** — `ToXModel(s)` /
   `ToXEntity(ies)`. This is not stylistic: `ToManualZoneSaves` /
   `FromManualZoneSaves` already exist and convert a **different axis** (live
   `[]entities.Zone` ⇄ save). The new `ToManualZoneSaveModels` /
   `ToManualZoneSaveEntities` convert the **storage** axis. Confusing the two
   silently produces the wrong shape.
6. **The alias façades for the row types are gone.** `internal/models/types.go`
   and `internal/models/editor_state_model/types.go` are deleted; 56 files now
   name `editor_state_model.ZoneContentRow` / `.ContentRuleRow` directly. §0.5.2's
   claim that `internal/models/` re-exports these two **no longer holds**.
7. **`internal/common/common_zone_contents` no longer exists.**
   `editor_state_model` → `common_zone_contents` → `editor_state_model` was a
   real import cycle once the defaults returned the model row. It was broken by
   moving `defaultPlayerZoneContentRows.go` (and its test folder) **into**
   `editor_state_model`; `GetDefaultPlayerZoneContentRows()` returns
   `[]editor_state_model.ZoneContentRow`. §0.5.2's "content defaults" row and the
   Phase 2 bullets that name that package are historical.
8. **No mapper DI.** §0.6 decision 7 is void: `NewUIState`, `NewWindow`,
   `NewStateHandler`, `NewTemplateHandler`, `NewZoneEditorHandler` and
   `NewRegenerationHandler` are back to their pre-Phase-10 argument lists.
   `app/` therefore imports **no** mapper — see the open question in Phase 12
   about the now-unused depguard relaxation.
9. **`MapTopology` lives in `internal/entities/topology/`**, not in
   `internal/entities/editor_state/` as §0.5.2 planned. `config.MapTopology`
   aliases it from there.

### 0.6 Owner decisions, 2026-08-24 (round 3) — **superseded in part by §0.7**

The owner's Phase 9 rework moved the ground under the original Phase 10 text
(the entity is now `editor_state.EditorState`, it has **no** codec methods, and
the model **embeds** it), so the shape was re-decided from scratch.

1. ~~**DTO shape — nine group DTOs, converted as structs.**~~ **Void (§0.7.1).**
   `EditorStateDto`
   anonymously embeds nine mirror structs (`TemplateIdentityDto` …
   `ManualEditSettingsDto`), so all 72 fields stay promoted and `dto.MapSize`
   reads the same as `state.MapSize`. **Go ignores struct tags when deciding
   type identity for a conversion**, so the mapper is nine conversions per
   direction — and a field added to an entity group but not mirrored in its DTO
   is a **compile error**, not a silent drop. This is the compile-time
   equivalent of the reflection drift guards, and it is why the 72-flat-field
   form was rejected: it is exactly the "can silently drop one" shape the owner
   rejected for the entity mapper in Phase 9.
   `SchemaVersion` is **not** on the DTO — it is persistence metadata and no
   consumer concern.
2. ~~**Leaf types are reused through the §0.5.2 alias façade**, never copied.~~
   **Void (§0.7.6).**
   `config.BonusEntry`, `config.MapTopology`, `models.ZoneContentRow` already
   *are* the entity types, so conversion identity holds and the ~85 + ~77 call
   sites §0.5.2 counted stay untouched. Two aliases were missing and are added:
   `models.ManualZoneSave` and `models.ManualConnectionSave`.
3. ~~**`app/` may import `internal/mappers`.**~~ **Void in practice (§0.7.8)** —
   the rule was relaxed and AGENTS.md §4.4 amended, but the mapper that needed
   it no longer exists, so nothing under `app/` uses the permission. Phase 12
   decides whether to keep or revert it.
4. **Phase 11 is folded in and superseded.** Changing a handler signature breaks
   `app/` in the same compile, so the crossing swap and the call-site rewrite
   cannot be separate phases. **(Still in force.)**
5. **The regeneration envelopes are split in two.** The service calls Model
   methods (`EqualsIgnoringManualEdits`, `LayoutDefiningOptionsChanged`,
   `DiffCastleSettings`) on the request's fields, so a behaviour-free DTO would
   break it. The DTO stays the *crossing* type; the handler maps it to new
   model-side types and the service takes those. This also closes a real §0.4
   breach — `internal/services/editor` consumed DTOs directly. **(Still in
   force: `internal/models/regeneration/` exists and the service no longer
   imports `internal/dtos`.)**
6. **Every crossing swaps except `ValidateEditorState`.** It sits on the
   per-frame path (`app/gui/models.EditorState.UpdateCurrentState` calls it on
   every panel write), so it stays Model-only and `EditorStateValidationDto`
   keeps carrying a Model. **Deliberate, documented exception** — do not "fix"
   it without a measurement. **(Still in force.)**
7. ~~**The new mapper is fully DI'd**~~ **Void (§0.7.8).** matching
   `IEditorStateEntityMapper`:
   struct + interface + wire set, threaded `program.go` →
   `editor.NewWindow(…)` → `drivers.NewUIState(…)`. The signature churn in the
   window and its test helpers is accepted.

### The eight crossings that swap

`LoadState`, `SaveState`, `GenerateTemplate`, `UpdateTemplate`,
`ReapplyCastleSettings`, `DecideRegeneration`, `DecideManualEditReapplication`,
`GetZoneEditorOptions`.

**All eight survived the §0.7 rework** — only the DTO's *shape* changed.
`ValidateEditorState` is the ninth and stays on the Model.

### Work (as landed — the §0.6 list it replaced is in the git history)

- [x] `EditorStateDto` embeds `editor_state_model.EditorState`; the nine `*Dto`
      group structs, `CastleSettingChangesDto` and the DTO⇄Model mapper are gone.
- [x] Model groups wrap their entity group, one struct per file in
      `internal/models/editor_state_model/`; `ContentSettings` and
      `ManualEditSettings` are re-declared with model element types and carry
      `ToXModel` / `ToXEntity` converters, as do `ZoneContentRow`,
      `ContentRuleRow`, `BonusEntry`, `ManualZoneSave`, `ManualConnectionSave`.
- [x] `internal/mappers/editorStateEntityMapper.go` carries all nine groups,
      including the two that need converters.
- [x] Regeneration types moved below the DTO line into
      **`internal/models/regeneration/`** (`DecisionRequest`, `Decision`,
      `ManualEditDecision`, `NextStateAction`); `dtos.NextStateAction` is an
      alias so `app/`'s switches still compile.
- [x] The eight crossings swapped in `handler_interfaces`, the handler
      implementations and the envelope DTOs.
- [x] `internal/models/types.go` and `internal/models/editor_state_model/types.go`
      deleted; 56 files repointed onto the real model types.
- [x] `internal/common/common_zone_contents` folded into `editor_state_model`,
      breaking the import cycle.
- [x] `internal/mappers` dropped from depguard `no-services-from-app` and
      AGENTS.md §4.4 amended — **now unused; Phase 12 decides whether to revert.**
- [x] Residual DTO-below-handlers audit recorded below — list only.

### Verification Plan
- `internal/services`, `internal/validators` and `internal/repositories` name no
  DTO type in any exported signature **except** the three picker/bonus/content
  services recorded in the audit.
- The frozen fixtures and the untagged
  [editorStateWireFormat_integration_test.go](../../test/integration/editorStateWireFormat_integration_test.go)
  pass **unchanged** — the DTO is off the persistence path entirely.
- Full unit + default + tagged integration suites; coverage ≥ 72.5 %; lint 0;
  `wire diff` clean.
- **GPU suite passes without `-update`** — no rendering change is intended.
- `BenchmarkEditorWindow_TabCycling` before/after from a detached worktree, to
  price the crossing rather than assume it.

### Phase Summary

**Landed 2026-08-24, then reworked by the owner 2026-08-25 (§0.7) and finished
2026-08-26. Committed as `2fc6b13 "Batch I wip 8"`; the working tree is clean.**
Eight of the nine editor-state crossings speak `EditorStateDto`;
`ValidateEditorState` stays on the Model by design.

**The shape that shipped.** `EditorStateDto` embeds
`editor_state_model.EditorState` — the Model owns the structure, the DTO is the
thin crossing wrapper. The nine `*Dto` group structs and the DTO⇄Model mapper
built on 2026-08-24 were **deleted** in the rework; do not rebuild them. On the
model side each group wraps its entity group, except `ContentSettings` and
`ManualEditSettings`, which must be re-declared because their slices carry model
element types and Go slices do not interconvert. `ZoneContentRow` shadows the
entity's `Rules` with `[]ContentRuleRow` and leaves the embedded slice nil.

**The regeneration §0.4 breach is closed.** `internal/services/editor` no longer
imports `internal/dtos` at all: the service takes `regeneration.DecisionRequest`
and the handler maps the crossing DTO onto it. The regeneration types could not
live in `internal/models` — `internal/models` ← `common_zone_contents` ←
`editor_state_model` was a real cycle at the time — so they sit in
**`internal/models/regeneration/`**. `dtos.NextStateAction` and its three
constants are **aliases** of that enum, so every `dtos.NextStateClear` in `app/`
kept compiling and `dtos.RegenerationDecisionDto(decision)` is a plain struct
conversion in the handler.

**Two nil bugs were found and fixed during the rework.**
`drivers.getPreviousStateDto` / `getNextStateDto` and
`regenerationHandler.DecideRegeneration` / `DecideManualEditReapplication`
dereferenced possibly-nil states. **nil is load-bearing** here: the decision
service reads a nil `Previous` as "first generation" and a nil `Next` as
"unarmed debounce". Wrapping with `new(...)` would have destroyed the signal
even without the panic, because it always yields a non-nil pointer. Both paths
now go through nil-preserving helpers. Any future refactor of these crossings
must preserve nil.

**The entity-mapper drift guard changed shape.** It used to marshal the model and
compare against the marshalled entity, which only worked while the model embedded
the entity verbatim; the two re-declared groups carry no json tags. It is now a
`ToModel(ToEntity(x))` **round trip** — tag-independent, and still fails if the
mapper drops a group.

**`ValidateEditorState` was deliberately left on the Model** (§0.6 decision 6),
so `EditorStateValidationDto` still carries one. It is the per-frame path.

**depguard changed, and the change is now unused.** `internal/mappers` was
removed from `no-services-from-app` and AGENTS.md §4.4 amended, but the mapper
that needed the permission was deleted in the rework — nothing under `app/`
imports `internal/mappers` today. `test/unit/architecture/dependency/dependency_test.go`
carries the matching allow-list entry. **Phase 12 must decide: keep the
permission, or revert both.**

**`app/` imports no editor-state entity package.** Eleven files under `app/` do
import the base `internal/entities` package (`entities.Zone`,
`entities.Connection` in the zone editor, preview and drivers) — that is the
residue Phase 12 records as a backlog entry, and it is out of scope for Batch I.

`test/unit/internal/dtos/editor_state_dto/` was deleted — a DTO that only embeds
the model has no behaviour to test, per §4.6. The mapper test folder built for
the deleted DTO⇄Model mapper was deleted with it.

### Audit — DTOs still used below the handler layer (list only, as instructed)

`internal/services/editor` is now **clean**, as are `internal/validators`,
`internal/repositories` and `internal/mappers`. Three services remain:

| Package | DTOs it names |
| --- | --- |
| `internal/services/bonuses` | `BonusCompositionRequestDto`, `BonusCompositionResultDto`, `ExistingBonusesDto` |
| `internal/services/pickers` | `PickerItemDto`, `PickerSpellDto`, `PickerEntryDto`, `PickerRowDto` |
| `internal/services/zone_content` | `ContentRuleCompositionRequestDto`, `ContentRuleCompositionResultDto` |

Each needs the same treatment the regeneration service just got: a model-side
request/result pair with the handler mapping onto it.

**Separately — DTOs that carry Models, all intended under §0.7.1:**
`EditorStateDto` (embeds it), `EditorStateValidationDto.State`,
`CastleSettingsReapplyRequestDto.Changes` and
`ManualEditDecisionDto.ReapplyWithCastleChanges` (both
`editor_state_model.CastleSettingChanges`). Phase 12's checker must **not** flag
these.

### Gates

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | exit 0 |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 |
| `go test ./test/unit/...` | exit 0 |
| `go test ./test/...` (untagged) | exit 0 |
| `go test -tags=integration_test ./test/integration/...` | exit 0 |
| GPU suite, **no `-update`** | exit 0 (40.3 s) |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Unit coverage | **73.7 %** — unchanged, floor 72.5 % |

**Re-run after the §0.7 rework (2026-08-26):** every gate above green again, and
coverage rose to **73.9 %**. Lint finished at **1 issue** — a `godox` TODO the
owner left in [stateHandler.go](../../internal/handlers/stateHandler.go), which is
theirs to resolve, not a regression. Treat **1 godox issue** as the current lint
baseline until that TODO is cleared.

The frozen fixtures and the untagged `editorStateWireFormat_integration_test.go`
passed **unchanged** through both rounds, as expected: the DTO is off the
persistence path entirely.

### Benchmark — the crossing is free

`BenchmarkEditorWindow_TabCycling`, `-benchtime=50x -count=6`, baseline taken
from a detached `git worktree` at `2106c70`:

| | allocs/op | B/op | ns/op (median) |
| --- | --- | --- | --- |
| Before | 9,519–9,521 | 907.5 KB | ~4.79 M |
| After | 9,543–9,545 | 927.2 KB | ~4.96 M |

**+24 allocs/op (+0.25 %)** and +2.1 % bytes; ns/op moves less than the spread
between repeats of either build, so it is noise. Excluding `ValidateEditorState`
kept the per-frame path untouched — the 24 allocations are the generate and
regenerate calls, which happen once per action, not once per frame.

> ⚠ **These numbers describe the §0.6 design, which was then deleted.** The
> §0.7 shape has no DTO⇄Model conversion at all — the DTO embeds the model — so
> the +24 allocations are very likely gone. **Phase 6 must re-measure from
> scratch** rather than treat this table as its baseline.

---

## Phase 11: DTOs at the handler crossings
Status: **Superseded — do not execute.** Folded into Phase 10 (§0.6 decision 4).

Changing a handler signature breaks `app/` in the same compile, so the crossing
swap and the `app/` call-site rewrite could not be separated. Its one genuinely
independent item — *"`app/` stops importing `internal/entities/editor_state`"* —
turned out to be **already true**: a tree scan on 2026-08-24 found no file under
`app/` importing that package, and a re-scan on 2026-08-26 confirmed it still
holds. The rest is Phase 10's work list.

> One correction the folded text must not carry forward: the note below claiming
> `internal/mappers/` is "already permitted" in `app/` was **wrong**. depguard
> denied it; Phase 10 relaxed the rule (§0.6 decision 3) — though the rework then
> removed the only consumer, so the permission now sits unused (§0.7.8).

The original body is kept below for the record only.

**Rewritten 2026-08-22 under [§0.5.1](#051-app-keeps-the-model-the-dto-is-the-crossing-not-the-storage).**
The original Phase 11 moved `app/` off the Model entirely; AGENTS.md §4.4 says
`app/` **may** hold a Model and only the *crossing* into `internal/` must be a
DTO. So this phase no longer reverses Phase 4 — it finishes it.

- [ ] `app/gui/models.EditorState` **keeps** `*editor_state_model.EditorState`
      for `current` / `previous` / `next`. No change.
- [ ] The panels' closures **keep** `UpdateState(func(*editor_state_model.EditorState))`.
      No change.
- [ ] Every **handler call** from `app/` passes and receives the DTO produced by
      Phase 10, with the GUI mapping stored Model → request DTO and response DTO
      → stored Model at that call site — exactly the two `Note over GUI` steps in
      AGENTS.md §4.4's diagram. The sites are the `IStatePersistenceHandler`,
      `IStateValidationHandler`, `ITemplateHandler`, `IZoneEditorHandler` and
      `IRegenerationHandler` calls in [app/gui/drivers/](../../app/gui/drivers/).
- [ ] The five state methods (`Clone`, `EqualsIgnoringManualEdits`,
      `LayoutDefiningOptionsChanged`, `DiffCastleSettings`, `HasManualEdits`)
      **stay callable from `app/`** — they are Model methods on a Model `app/`
      legitimately holds. The original bullet moving them behind handlers is
      withdrawn, and with it the per-frame conversion cost.
- [ ] `app/` stops importing `internal/entities/editor_state`. The `internal/models`
      imports stay.
- [ ] Re-point `app/gui/editor/window_testexports.go` and the mocks only where
      handler signatures actually changed.

> Because the GUI maps at the call site, the mapper from Phase 10 must be
> reachable from `app/` — `internal/mappers/` is already permitted there under
> AGENTS.md §4.4.

### Verification Plan
- No file under `app/` imports `internal/entities/editor_state`.
- Every editor-state handler signature in `handler_interfaces` names a DTO, and
  no `app/` → `internal/` call passes a Model across the boundary.
- **GPU suite passes without `-update`** — the rendering must be proven unchanged.
- Round-trip non-aliasing tests still pass; full unit + integration suites;
  coverage ≥ 72.5 %.

### Phase Summary
_(write when phase completes)_

---

## Phase 12: Enforce the layering mechanically
Status: Not started — **next phase to run** (Phase 6 is done).

§0.4 was violated silently for four phases because nothing checked it. A doctrine
without a gate is a suggestion.

> **Refreshed 2026-08-26 against the post-§0.7 tree.** The rule list below is
> current; the counts and exceptions were re-verified.

- [ ] Extend [cmd/testlayoutcheck](../../cmd/testlayoutcheck) (or add a sibling
      `cmd/layeringcheck`) with import rules that fail the build:
      - `internal/entities/**` must not import `internal/models`, `internal/dtos`,
        `internal/services`, `internal/handlers` or `internal/helpers`.
      - `internal/entities/**` must be imported only by `internal/repositories/**`,
        `internal/models/**`, `internal/entities/**`, `internal/mappers/**` and
        `internal/helpers/*_helpers/**` (§0.5.4).
      - `internal/dtos/**` must be imported only by `internal/handlers/**`,
        `internal/dtos/**` and `app/**`.
      - ~~`internal/models/**` must not be imported by `app/**`.~~ **Withdrawn**
        under §0.5.1 — AGENTS.md §4.4 permits it explicitly.
      - ⚠ **Do not add a "a DTO must not name a Model" rule.** Under §0.7.1 that
        is the intended shape: `EditorStateDto` *embeds*
        `editor_state_model.EditorState`, and `EditorStateValidationDto.State`,
        `CastleSettingsReapplyRequestDto.Changes` and
        `ManualEditDecisionDto.ReapplyWithCastleChanges` all deliberately carry
        Models.
- [ ] Seed the checker with an **explicit, shrinking allow-list** of the
      pre-existing violations so it can be turned on before they are all fixed.
      An empty allow-list is the end state, not the starting condition.
      **The allow-list content, verified 2026-08-26:**
      - `internal/services/bonuses`, `internal/services/pickers` and
        `internal/services/zone_content` import `internal/dtos` (the audit in
        Phase 10's summary lists the exact types). `internal/services/editor`,
        `internal/validators`, `internal/repositories` and `internal/mappers` are
        **clean** — do not allow-list them.
      - **11 files under `app/` import the base `internal/entities` package**
        (`entities.Zone` / `entities.Connection` in the zone editor, the drivers
        and `app/gui/models/editorState.go`). **No file under `app/` imports
        `internal/entities/editor_state`, `/topology` or `/template`.**
- [ ] **Decide the fate of the unused depguard permission.** Phase 10 removed
      `internal/mappers` from `no-services-from-app` (and added the matching entry
      to `test/unit/architecture/dependency/dependency_test.go`) so the GUI could
      reach the DTO⇄Model mapper. §0.7 deleted that mapper, so **nothing under
      `app/` imports `internal/mappers` today**. Either revert both edits, or keep
      the permission and say why in AGENTS.md §4.4. Do not leave it undecided.
- [ ] Add a VS Code task for it and wire it into the Phase 7 gate list.
- [ ] Amend **AGENTS.md §4.4**: fold the §0.4 table, the §0.5.1 crossing rule, the
      §0.5.4 Entity-namer list and the §0.7 Model-owns-the-structure rule into the
      instructions so the doctrine outlives this plan file. Do **not** forbid
      `app/` → `internal/models`.
- [ ] Record the residual breach as its own backlog entry: the 11 `app/` files
      that import `internal/entities` directly. Under §0.5.1 the `internal/models`
      imports are **no longer** a breach, so the residue is far smaller than the
      28 files an earlier carry-forward recorded. Out of scope for Batch I.

### Verification Plan
- `go run ./cmd/testlayoutcheck .` (or the new checker) exits 0 with the
  allow-list, and exits 1 when a deliberate test violation is introduced —
  prove the gate actually trips.
- Every Phase 7 gate green.

### Phase Summary
_(write when phase completes)_

---

## Final Recap
_(write when all phases complete)_

## Deployment Plan
_(write when all phases complete)_
