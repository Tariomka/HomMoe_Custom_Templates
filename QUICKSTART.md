# Quick Start — Olden Era Custom Templates GUI

This is a desktop GUI app, not a CLI. Launching it opens a window where
you configure a template across four tabs, preview the layout in a side
panel, then click **Generate** + **Save Template** to write a `.rmg.json`
file.

## 1. Run the App

Requires Go **1.25.8+**.

```powershell
# Run directly
go run .

# Or build a binary first
go build .
.\hommoe_custom_templates.exe
```

For iterative work the project ships an [air](https://github.com/air-verse/air)
config ([.air.toml](.air.toml)):

```powershell
air
```

## 2. The Window

The window has four regions:

- **Toolbar** (top): `New`, `Open…`, `Save`, `Save As…`, with the current
  settings-file path on the right (a trailing `*` marks unsaved edits).
- **Tabs** (centre-left): the four configuration tabs listed below.
- **Preview** (centre-right): live render of the most recently generated
  template, with a `Refresh` button. The preview PNG is written automatically
  when you **Save Template**.
- **Footer** (bottom): output folder picker, `Browse…`, `Reveal`,
  `Generate Template`, `Save Template`, plus a status line. The output folder
  is pre-filled from your Steam install when it can be located.

### Tabs

#### 1. General
- **Template** — template name (used as the output filename), player count
  (2–8), map size, and a toggle for non-official large sizes (>240).
- **Hero Restrictions** — game mode (`Classic` / `SingleHero`) and hero count
  min / max / increment (hidden for `SingleHero`).
- **Rules** — faction-laws XP % and astrology XP % (25–200), a victory
  condition (Standard / Lost Starting City / Hold City / Tournament) and the
  matching win/loss toggles: lose-on-lost-city (+ grace days),
  lose-on-lost-hero, hold-a-city (+ days to hold), gladiator arena (+ delay /
  count days) and tournament (first day, interval, points, save army).

#### 2. Layout
- **Topology** — pick one of ten layouts; a description line explains the
  selected one.
- **Manual zone editor** — opens the visual editor for the last generated
  template (generate first).
- **Connectivity** — roads, random portals (+ max connections), remote
  footholds, player isolation, matching player castle factions, and minimum
  neutrals between players.
- **Zone sizes** — player / neutral size multipliers, plus hub size and hub
  castles for the Hub topology.
- **Difficulty & Density** — resource / structure density, neutral stack
  strength, border-guard strength and guard randomization.
- **Advanced zone control** — when enabled, set total neutral zones, castles
  per zone, and per-tier (Low / Medium / High) counts split by with / without
  castle.

#### 3. Zone Content
Seed mandatory content per zone tier — **Player**, **Low**, **Medium**,
**High Neutral** and **Hub** — across mines, utility structures, treasures,
unit recruitment, resource banks and hero-improvement groups. Each row can
open a **Manage Rules** dialog for placement constraints (distance to
road/town, guarded, solo encounter, variant).

#### 4. Bonuses & Bans (experimental)
Add game-start bonuses, banned items, banned spells and guard-value overrides.
Entries are added through picker dialogs and removed per row. Effects only
apply at generation time.

## 3. Save / Load Settings

Your editor state is persisted as `.gen.json` files (the
`dtos.EditorStateDto` model, handled by `services.SaveSettingsFile` /
`services.LoadSettingsFile`).

- **Save / Save As…** — write the current widget state to disk.
- **Open…** — load a `.gen.json` file back into the widgets.
- **New** — reset to defaults.

The toolbar's right-hand label always shows the active `.gen.json` path
and an asterisk when there are unsaved changes.

## 4. Generate a Template

1. The output folder is pre-filled from your Steam install when it can be
   found; otherwise pick one in the footer (`Browse…`).
2. Click **Generate Template** — this builds the template in memory and
   refreshes the preview panel.
3. Click **Save Template** — writes `<TemplateName>.rmg.json` plus a preview
   `<TemplateName>.png` into the chosen folder.
4. **Reveal** opens the output folder in your file explorer.

Drop the resulting `.rmg.json` into the game's templates directory and
pick it from the in-game Random Map Generator screen.

## 5. Programmatic Use

If you want to call the generator from Go without the GUI:

```go
package main

import (
    "log"

    "github.com/Tariomka/hommoe_custom_templates/internal/models/config"
    "github.com/Tariomka/hommoe_custom_templates/internal/services"
    "github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
)

func main() {
    cfg := config.NewGeneratorConfig() // sensible defaults
    cfg.TemplateName = "Programmatic Map"
    cfg.PlayerCount = 4
    cfg.MapSize = 160
    cfg.Topology = config.TopologyDefault // Ring

    template := template_generator.NewTemplateGenerator(cfg).Generate()

    out, err := services.WriteTemplate(".", template)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("wrote", out)
}
```

`config.NewGeneratorConfig()` matches the GUI's defaults
(2 players, size 160, topology Circles, Classic mode, etc.).

## 6. Map Sizes

Map size is the raw tile count (`sizeX = sizeZ = MapSize`). Common
presets the slider snaps to:

| Tiles | Notes                                  |
|-------|----------------------------------------|
|  64   | Tiny                                   |
|  96   | Small                                  |
| 128   | Medium                                 |
| 160   | Large (default)                        |
| 192   | Extra Large                            |
| 240   | Maximum non-experimental               |
| 256+  | Experimental (enable the checkbox)     |
| 512   | Largest experimental size              |

## 7. Topologies at a Glance

- **Circles** (`TopologyCircles`) — default; concentric rings sorted by zone tier.
- **Random** — positions and connections are randomised.
- **Ring** (`TopologyDefault`) — players in a circle, neighbours connected.
- **Hub** (`TopologyHubAndSpoke`) — all players connect through a central hub zone.
- **Chain** — players strung in a line.
- **Shared Web** — players linked through shared neutral zones.
- **Square** (`TopologySquare`) — players line the edges of a square; neutrals on the edges and inside.
- **Geometric** (`TopologyGeometric`) — symmetric geometric shapes built around a central zone.
- **Cross** (`TopologyCross`) — zones and connections radiate from a centre into cross-shaped arms.
- **Fractal** (`TopologyFractal`) — every player anchors a self-similar fractal that branches inward; low neutral zones sit nearest the player and high zones weave together at the centre, so no two players border directly.

## 8. Troubleshooting

- **Window doesn't open** — Gio needs a working graphics stack. If
  running over Remote Desktop or in a headless container, run on a real
  desktop session instead.
- **`go build` complains about Go version** — install Go 1.25.8 or later
  (`go version`).
- **`Save Template` is disabled** — click **Generate Template** first;
  the button only enables once a template is in memory.
- **Nothing happens after Generate** — check the footer status line for
  errors and verify the template name is non-empty.
- **Template won't load in the game** — verify the JSON is valid
  (`Get-Content map.rmg.json | python -m json.tool`) and compare against
  files in [data/ExampleTemplates/](data/ExampleTemplates).

## 9. Reference

- [README.md](README.md) — architecture, project layout, build/test commands.
- [data/ExampleTemplates/](data/ExampleTemplates) — 57 reference templates
  shipped with the game.
