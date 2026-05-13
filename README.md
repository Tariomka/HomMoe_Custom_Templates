# Heroes of Might and Magic: Olden Era — Custom Templates

A desktop GUI for designing and generating `.rmg.json` random map templates for
**Heroes of Might and Magic: Olden Era**. Written in Go using the
[Gio](https://gioui.org) immediate-mode UI toolkit.

The app lets you configure every knob the game's RMG exposes (players, map
size, topology, zone counts, victory conditions, neutral zone quality
distribution, mandatory content, etc.), persist your work as a `.gen.json`
settings file, and emit a ready-to-drop-in `.rmg.json` template.

## Project Structure

```
.
├── main.go                         # GUI entry point (boots gioui app)
├── go.mod                          # Go module + gioui.org dependency
├── data/
│   ├── ExampleTemplates/           # 57 reference .rmg.json templates
│   └── GameData/GeneratorData/     # Game configuration & content pools
│       ├── generator_config.json
│       ├── generator_environment_assets.json
│       ├── generator_stats_config.json
│       ├── content_lists/
│       ├── content_pools/          # quality tiers t2..t5
│       ├── encounter_templates/
│       └── zone_layouts/
├── internal/
│   ├── constants/                  # Stable string IDs (content, includes, groups)
│   ├── gui/                        # Gio UI: window, tabs, widgets, theme, state
│   ├── helpers/                    # Small shared utilities (e.g. SID lookup)
│   ├── models/
│   │   ├── settingsFile.go         # .gen.json persistence model
│   │   ├── sidMapping.go
│   │   ├── types.go                # Re-exports of generator/template types
│   │   ├── zoneContentItem.go
│   │   ├── generator/              # Inputs to the generator (settings, topology, …)
│   │   └── template/               # Output schema for .rmg.json files
│   └── services/
│       ├── template_generator.go   # Core generation logic (Generate)
│       └── zone_content_manager.go # Mandatory content per zone tier
├── test/
│   ├── models/template/            # RmgTemplateModel round-trip tests
│   └── services/                   # Generator / topology / zone tests
└── bin/                            # Pre-built binaries (gitignored target)
```

## Features

- **Gio desktop GUI** (`gioui.org v0.9.0`) with four tabs:
  1. **Map Setup** — template name, game mode, players, map size, topology
  2. **Generation Options** — roads, portals, footholds, player isolation,
     advanced neutral-zone count breakdown by quality / castle presence
  3. **Game Rules** — victory condition, hero counts, faction-laws &
     astrology XP, gladiator arena, tournament rules, lost-city / city-hold
  4. **Zone Content (EXP)** — extra mandatory content seeded into player zones
- **Settings persistence**: load/save `.gen.json` files (`models.SettingsFile`)
- **Template generation**: emit `.rmg.json` files compatible with the in-game RMG
- **Five topologies**: Random (default), Ring, Hub-and-Spoke, Chain, Shared Web
- **Map sizes**: 64–240 tiles, plus an experimental range up to 512
- **Players**: 2–8
- **Quality-tiered neutral zones**: Low / Medium / High, with optional
  fine-grained per-tier counts split by "with castle" / "without castle"
- **Mandatory-content seeding** through `services.ZoneContentManager`

## Building & Running

Requires **Go 1.25.8** or later (see `go.mod`). On Windows, Gio additionally
needs a working CGO-free build — no extra toolchain is required.

```powershell
# Build the GUI
go build -o bin/template-gui.exe .

# Run directly without producing a binary
go run .

# Run tests
go test ./test/...
```

A pre-built `bin/template-gui.exe` is checked in for convenience.

## Workflow

1. Launch the GUI (`go run .` or `bin\template-gui.exe`).
2. Configure the template across the four tabs.
3. **Save** / **Save As** writes a `.gen.json` settings file (your inputs).
4. **Pick output folder**, then **Generate** writes
   `<TemplateName>.rmg.json` into that folder.
5. Drop the `.rmg.json` into the game's templates folder and pick it from
   the in-game RMG screen.

The toolbar also has shortcut buttons to open the bundled example templates
folder, and to jump out to the project's Discord / GitHub / patch-notes pages.

## Topologies

| Topology      | Constant                          | Shape                                                        |
|---------------|-----------------------------------|--------------------------------------------------------------|
| Random        | `generator.TopologyRandom`        | Default. Random placement / Delaunay-style connections.      |
| Ring          | `generator.TopologyDefault`       | Players in a circle, each connected to neighbours.           |
| Hub-and-Spoke | `generator.TopologyHubAndSpoke`   | All players connect through a central hub neutral zone.      |
| Chain         | `generator.TopologyChain`         | Linear arrangement of zones.                                 |
| Shared Web    | `generator.TopologySharedWeb`     | Players connected through shared neutral zones.              |

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

1. **GUI** (`internal/gui/`) — Gio widgets, four-tab layout, file dialogs,
   binds widget state into `models.SettingsFile`.
2. **Models** (`internal/models/`)
   - `SettingsFile` — `.gen.json` persistence schema (the editor's source of truth).
   - `models/generator/` — `GeneratorSettings`, `MapTopology`,
     `ZoneConfiguration`, `AdvancedSettings`, `HeroSettings`,
     `GameEndConditions`, `GladiatorArenaRules`, `TournamentRules`,
     `NeutralZoneQuality`.
   - `models/template/` — On-disk `.rmg.json` schema (`RmgTemplateModel`,
     `Variant`, `Zone`, `GameRules`, `ContentPool`, `MandatoryContent`, …).
   - `models/types.go` — Re-export aliases used by the GUI.
3. **Services** (`internal/services/`)
   - `Generate(*GeneratorSettings) (*RmgTemplateModel, error)` — the entry
     point that turns settings into a template.
   - `ZoneContentManager` — builds per-tier mandatory content for player
     zones and Low / Medium / High neutral zones, plus content count limits.
4. **Constants** (`internal/constants/`) — stable string IDs for content
   items, include lists and content groups.

### Generation Flow

```
GUI widgets
   │   (translated by gui.settingsFileToGenerator)
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
                RmgTemplateModel  ──►  <Name>.rmg.json
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
  round-trip tests under `test/models/template/`.

## Related

- The original C# WPF editor that this project was ported from. See
  `MIGRATION.md` and `CONVERSION_SUMMARY.md` for the migration history.
- `data/ExampleTemplates/` — 57 reference templates from the game itself,
  invaluable when reasoning about the on-disk schema.

## License

See the main project repository for license information.
