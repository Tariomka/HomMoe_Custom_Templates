# Heroes of Might and Magic: Olden Era - Custom Templates

A desktop GUI for designing and generating `.rmg.json` random map templates for
**Heroes of Might and Magic: Olden Era**. Written in Go using the
[Gio](https://gioui.org) immediate-mode UI toolkit.

DISCLAIMER: This is a semi-AI generated rewrite of the original [Olden Era - Template Generator](https://github.com/KhanDevelopsGames/Olden-Era---Template-Generator) for personal use, developing on Windows and using on Linux (SteamDeck). All credit and inspiration goes to the people behind that project. Don't use this project, instead go to the original.

The app lets you configure every knob the game's RMG exposes (players, map
size, topology, zone counts, victory conditions, neutral zone quality
distribution, mandatory content, etc.), preview the resulting layout in a
side panel, persist your work as a `.gen.json` settings file, and emit a
ready-to-drop-in `.rmg.json` template.

## Project Structure

```
.
├── main.go                              # Process entry; calls internal.StartApplication
├── go.mod                               # Go module + gioui.org dependency
├── data/
│   ├── ExampleTemplates/                # 57 reference .rmg.json templates
│   └── GameData/GeneratorData/          # Game configuration & content pools
│       ├── generator_config.json
│       ├── generator_environment_assets.json
│       ├── generator_stats_config.json
│       ├── content_lists/
│       ├── content_pools/               # quality tiers t2..t5
│       ├── encounter_templates/
│       └── zone_layouts/
├── internal/
│   ├── program.go                       # Gio app bootstrap + event loop
│   ├── constants/                       # Stable string IDs (content, includes, groups, legend, UI)
│   ├── gui/
│   │   ├── window.go                    # Window layout: toolbar, tabs, preview, footer
│   │   ├── components/
│   │   │   ├── basicSetupPanel.go       # Tab 1: Map Setup
│   │   │   ├── generationPanel.go       # Tab 2: Generation Options
│   │   │   ├── rulesPanel.go            # Tab 3: Game Rules
│   │   │   ├── zoneContentPanel.go      # Tab 4: Zone Content
│   │   │   ├── previewPanel.go          # Live map preview sidebar
│   │   │   ├── footerPanel.go           # Output dir, Generate, Save Template
│   │   │   ├── toolbar.go               # New / Open / Save / Save As
│   │   │   ├── tab.go, panel.go, state.go
│   │   │   ├── content/                 # Dropdowns, segment buttons, zone-content rows
│   │   │   ├── themes/                  # Color palette + Gio material theme
│   │   │   └── widgets/                 # Reusable buttons, sliders, sections, etc.
│   │   └── utils/                       # Drawing, file IO, math helpers
│   ├── helpers/                         # Small shared utilities (e.g. SID lookup)
│   ├── models/
│   │   ├── settingsFile.go              # .gen.json persistence model
│   │   ├── sidMapping.go
│   │   ├── types.go                     # Re-exports of generator/template types
│   │   ├── zoneContentItem.go
│   │   ├── generator/                   # Inputs to the generator
│   │   └── template/                    # Output schema for .rmg.json files
│   └── services/
│       ├── templateGenerator.go         # Generate(*GeneratorSettings) -> *RmgTemplate
│       ├── templateWriter.go            # Marshal + write <Name>.rmg.json
│       ├── settingsFileLoader.go        # Load/save .gen.json settings files
│       ├── zoneContentManager.go        # Mandatory content per zone tier
│       ├── previewLayout.go             # Computes preview zone geometry
│       └── previewRenderer.go           # PNG export of the preview canvas
└── test/
    ├── models/template/                 # RmgTemplate round-trip tests
    └── services/                        # Generator / topology / zone tests
```

## Features

- **Gio desktop GUI** (`gioui.org v0.9.0`) with four configuration tabs and a
  live preview sidebar:
  1. **Map Setup** — template name, game mode, players, map size, topology
  2. **Generation Options** — roads, portals, footholds, player isolation,
     advanced neutral-zone count breakdown by quality / castle presence
  3. **Game Rules** — victory condition, hero counts, faction-laws &
     astrology XP, gladiator arena, tournament rules, lost-city / city-hold
  4. **Zone Content** — extra mandatory content seeded into player zones
- **Live preview panel** — renders the generated topology (zones,
  connections, hubs, portals) and can export a PNG snapshot.
- **Settings persistence**: load/save `.gen.json` files (`models.SettingsFile`).
- **Template generation**: emit `.rmg.json` files compatible with the in-game RMG.
- **Five topologies**: Random (default), Ring, Hub-and-Spoke, Chain, Shared Web.
- **Map sizes**: 64–240 tiles, plus an experimental range up to 512.
- **Players**: 2–8.
- **Quality-tiered neutral zones**: Low / Medium / High, with optional
  fine-grained per-tier counts split by "with castle" / "without castle".
- **Mandatory-content seeding** through `services.ZoneContentManager`.

## Building & Running

Requires **Go 1.25.8** or later (see [go.mod](go.mod)).

```powershell
# Run directly
go run .

# Build a binary
go build -o bin/template-gui.exe .
.\bin\template-gui.exe

# Run tests
go test ./test/...
```

Hot reload via [air](https://github.com/air-verse/air) is configured in
[.air.toml](.air.toml); set `HOT_RELOAD=1` to start the window minimized.

## Workflow

1. Launch the GUI (`go run .`).
2. Configure the template across the four tabs. The preview panel updates
   each time you click **Generate Template**.
3. **Save** / **Save As…** writes a `.gen.json` settings file (your inputs).
4. Pick an output folder in the footer, click **Generate Template** to
   compute the template (and refresh the preview), then click
   **Save Template** to write `<TemplateName>.rmg.json` into that folder.
5. Drop the `.rmg.json` into the game's templates folder and pick it from
   the in-game RMG screen.

## Topologies

| Topology      | Constant                          | Shape                                                        |
|---------------|-----------------------------------|--------------------------------------------------------------|
| Circles       | `generator.TopologyCircles`       | Default. Concentric rings sorted by zone tier.               |
| Random        | `generator.TopologyRandom`        | Random placement / Delaunay-style connections.               |
| Ring          | `generator.TopologyDefault`       | Players in a circle, each connected to neighbors.            |
| Hub-and-Spoke | `generator.TopologyHubAndSpoke`   | All players connect through a central hub neutral zone.      |
| Chain         | `generator.TopologyChain`         | Linear arrangement of zones.                                 |
| Shared Web    | `generator.TopologySharedWeb`     | Players connected through shared neutral zones.              |
| Square        | `generator.TopologySquare`        | Players line the edges of a square; neutrals on edges and inside. |
| Geometric     | `generator.TopologyGeometric`     | Symmetric geometric shapes built around a central zone.      |
| Cross         | `generator.TopologyCross`         | Zones and connections radiate from a centre into cross arms. |

## Game Modes & Victory Conditions

Game modes (UI exposes both, generator currently always emits `Classic`):

- `Classic`
- `SingleHero` (reserved)

Victory condition IDs (`SettingsFile.VictoryCondition`):

| ID                 | UI label             |
|--------------------|----------------------|
| `win_condition_1`  | Standard             |
| `win_condition_3`  | Lost Starting City   |
| `win_condition_5`  | Hold City            |
| `win_condition_6`  | Tournament           |

Independent toggles also exist for `lostStartCity`, `lostStartHero`,
`cityHold`, `gladiatorArena`, and `tournament`.

## Architecture

### Layers

1. **GUI** (`internal/gui/`) — Gio widgets, four-tab layout + preview
   sidebar + footer, file dialogs, binds widget state into
   `models.SettingsFile`.
2. **Models** (`internal/models/`)
   - `SettingsFile` — `.gen.json` persistence schema (the editor's source of truth).
   - `models/generator/` — `GeneratorSettings`, `MapTopology`,
     `ZoneConfiguration`, `AdvancedSettings`, `HeroSettings`,
     `GameEndConditions`, `GladiatorArenaRules`, `TournamentRules`,
     `NeutralZoneQuality`.
   - `models/template/` — On-disk `.rmg.json` schema (`RmgTemplate`,
     `Variant`, `Zone`, `GameRules`, `ContentPool`, `MandatoryContent`, …).
   - `models/types.go` — Re-export aliases used by the GUI.
3. **Services** (`internal/services/`)
   - `Generate(*GeneratorSettings) (*RmgTemplate, error)` — turns settings
     into a template.
   - `templateWriter.go` — serializes a template to disk.
   - `settingsFileLoader.go` — load/save `.gen.json` files.
   - `ZoneContentManager` — builds per-tier mandatory content for player
     zones and Low / Medium / High neutral zones, plus content count limits.
   - `BuildPreviewLayout` / preview renderer — compute zone geometry for
     the preview panel and export a PNG snapshot.
4. **Constants** (`internal/constants/`) — stable string IDs for content
   items, include lists, content groups, the preview legend and UI strings.

### Generation Flow

```
GUI widgets
   │   (translated by services.SettingsToGenerator)
   ▼
SettingsFile  ──►  GeneratorSettings
                        │
                        ▼
                services.Generate
                ├── buildNeutralZonePlan
                ├── buildTopologyAdjacency
                ├── buildGameRules
                ├── buildVariant / buildZones
                └── ZoneContentManager (mandatory content)
                        │
                        ▼
                   RmgTemplate
                   ├──► services.BuildPreviewLayout  (preview panel)
                   └──► services template writer     ──►  <Name>.rmg.json
```

## Testing

```powershell
# Full suite
go test ./test/...

# Just the generator / services tests
go test ./test/services/...

# A single test
go test ./test/services/ -run TestGenerate_DefaultSettings_Succeeds
```

## Notes

- Zones are labelled `A`–`AF` (up to 32 zones).
- Default guard randomization is `0.05`.
- Default connections-per-zone is `1`.
- Player zones are emitted first (`A` onwards), neutral zones follow.
- The generator and the example templates are kept in sync by the
  round-trip tests under [test/models/template](test/models/template).

## Related

- [data/ExampleTemplates/](data/ExampleTemplates) — 57 reference templates
  from the game itself, invaluable when reasoning about the on-disk schema.
- [QUICKSTART.md](QUICKSTART.md) — short guide for first-time users.

## License

See the main project repository for license information.
