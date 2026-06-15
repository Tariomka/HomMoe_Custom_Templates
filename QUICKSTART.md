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
go build -o bin/template-gui.exe .
.\bin\template-gui.exe
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
  template, with `Refresh` and `Save PNG` buttons.
- **Footer** (bottom): output folder picker, `Browse…`, `Reveal`,
  `Generate Template`, `Save Template`, plus a status line.

### Tabs

#### 1. Map Setup
- **Template name** — used as the output filename (`<name>.rmg.json`).
- **Game mode** — `Classic` or `SingleHero` (only `Classic` is currently emitted).
- **Players** — 2..8.
- **Map size** — slider in tile units (64..240). Tick
  *"Allow experimental large map sizes (>240)"* to extend to 256..512.
- **Topology** — `Random` (default), `Ring`, `Hub`, `Chain`, `Shared Web`.

#### 2. Generation Options
- **Connectivity** — toggle roads, random portals (with max-portal slider),
  remote footholds, experimental balanced placement, player isolation,
  faction-matched player castles, minimum neutrals between players.
- **Advanced neutral zones** — when enabled, set per-quality counts split
  between *with castle* and *without castle* (Low / Medium / High).
  When disabled, just set the total `Neutral zone count`.
- **Zone sizes** — player / neutral / hub multipliers, hub castle count,
  guard randomization.
- **Content density** — resource and structure density percents
  (independent or via the unified `Content density` knob).

#### 3. Game Rules
- **Victory condition** — Standard / Lost Starting City / Hold City / Tournament.
  When Hold City is selected the generator also picks a neutral zone to
  host the city-hold target.
- **Heroes** — min, max, increment.
- **Faction laws XP %** and **Astrology XP %** — 25..200 (100 = baseline).
- **Lost-city / lost-hero / city-hold** day toggles.
- **Gladiator arena** — start delay and counter day.
- **Tournament** — first day, interval, points-to-win, save-army.

#### 4. Zone Content
Add extra mandatory content items to seed into player zones in addition
to the defaults from `services.ZoneContentManager`.

## 3. Save / Load Settings

Settings are persisted as `.gen.json` files (the
`models.SettingsFile` struct, handled by `services.settingsFileLoader`).

- **Save / Save As…** — write the current widget state to disk.
- **Open…** — load a `.gen.json` file back into the widgets.
- **New** — reset to defaults.

The toolbar's right-hand label always shows the active `.gen.json` path
and an asterisk when there are unsaved changes.

## 4. Generate a Template

1. Pick an output folder in the footer (`Browse…` button next to the path).
2. Click **Generate Template** — this builds the template in memory and
   refreshes the preview panel.
3. Click **Save Template** — writes `<TemplateName>.rmg.json` into the
   chosen folder.
4. **Reveal** opens the output folder in Explorer.

Drop the resulting `.rmg.json` into the game's templates directory and
pick it from the in-game Random Map Generator screen.

The preview panel's **Save PNG** button writes a snapshot of the current
preview canvas next to the template.

## 5. Programmatic Use

If you want to call the generator from Go without the GUI:

```go
package main

import (
    "encoding/json"
    "os"

    "github.com/Tariomka/hommoe_custom_templates/internal/models/generator"
    "github.com/Tariomka/hommoe_custom_templates/internal/services"
)

func main() {
    s := generator.NewGeneratorSettings() // sensible defaults
    s.TemplateName = "Programmatic Map"
    s.PlayerCount = 4
    s.MapSize = 160
    s.Topology = generator.TopologyDefault // Ring

    tmpl, err := services.Generate(s)
    if err != nil {
        panic(err)
    }

    data, _ := json.MarshalIndent(tmpl, "", "  ")
    _ = os.WriteFile("Programmatic Map.rmg.json", data, 0o644)
}
```

`generator.NewGeneratorSettings()` matches the GUI's defaults
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
