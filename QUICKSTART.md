# Quick Start — Olden Era Custom Templates GUI

This is a desktop GUI app, not a CLI. Launching it opens a window where you
configure a template across three tabs, preview the layout in a side panel,
then click **Generate** + **Save Template** to write a `.rmg.json` file.

## 1. Run the App

Requires Go **1.27.0+**.

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

### Building on Linux

Gio links against the system graphics and input stack, so a Linux build needs
development headers first. On Debian/Ubuntu — this is exactly what CI installs,
see [.github/workflows/setup-steps/action.yml](.github/workflows/setup-steps/action.yml):

```bash
sudo apt-get update
sudo apt-get install -y \
  libgles2-mesa-dev libegl1-mesa-dev libffi-dev \
  libxkbcommon-dev libxkbcommon-x11-dev libvulkan-dev libwayland-dev \
  libx11-dev libx11-xcb-dev libxcb1-dev libxcursor-dev libxfixes-dev \
  libxrandr-dev libxinerama-dev libxi-dev xorg-dev
```

Windows needs no extra packages.

## 2. The Window

The window has three regions:

- **Toolbar** (top): `New`, `Load`, `Save`, `Save To`, `Exit`, with the current
  settings-file path on the right (a trailing `*` marks unsaved edits).
- **Tabs** (left): the three configuration tabs listed below.
- **Preview** (right): live render of the most recently generated template,
  the output-directory picker (`Browse`, `Reveal`), a status line, and the
  `Generate` / `Save Template` buttons. The output folder is pre-filled from
  your Steam install when it can be located, and the preview PNG is written
  automatically when you **Save Template**.

There is no separate footer — everything below the preview image belongs to
the preview panel.

### Tabs

#### 1. General
- **Template** — template name (used as the output filename), player count
  (2–8), map size, and a toggle for non-official large sizes (>240).
- **Hero Restrictions** — game mode (`Classic` / `SingleHero`) and hero count
  min / max / increment (hidden for `SingleHero`).
- **Rules / Conditions** — faction-laws XP % and astrology XP % (25–200), a
  victory condition (Standard / Lost Starting City / Guardian Arena / Hold
  City / Tournament) and the matching win/loss toggles: lose-on-lost-city
  (+ grace days), lose-on-lost-hero, hold-a-city (+ days to hold), gladiator
  arena (+ delay / count days) and tournament (first day, interval, points,
  save army).

#### 2. Layout & Zones
- **Topology** — pick one of the layouts listed in the
  [README topology table](README.md#topologies); a description line explains
  the selected one.
- **Connectivity** — roads, random portals (+ max connections), remote
  footholds, abandoned outposts, player isolation, matching player castle
  factions, and minimum neutrals between players.
- **Zone sizes** — player / neutral size multipliers, plus hub size for the
  hub topologies.
- **Difficulty & Density** — resource / structure density, neutral stack
  strength, border-guard strength and guard randomization.
- **Manual zone editing** — opens the visual editor for the last generated
  template (generate first).
- **Zones** — player castles per zone and either a single **Total neutral
  zones** slider or, with **Advanced zone control** enabled, per-tier
  (Lowest / Low / Medium / High) counts split by with / without castle, plus
  a **Hub** sub-section for hub topologies.
- **Zone content** — the `Edit zone content...` buttons open a per-tier dialog
  (Player, Lowest / Low / Medium / High Neutral, Hub) where you seed mandatory
  content across mines, utility structures, treasures, unit recruitment,
  resource banks and hero-improvement groups. Each row can open a
  **Manage Rules** dialog for placement constraints (distance to road/town,
  guarded, solo encounter, variant).

#### 3. Bonuses & Bans (experimental)
Add game-start bonuses, banned items, banned spells and guard-value overrides.
Entries are added through picker dialogs and removed per row. Effects only
apply at generation time.

## 3. Save / Load Settings

Your editor state is persisted as `.gen.json` files (the `dtos.EditorStateDto`
model, handled by `file_service.FileService.SaveSettings` /
`file_service.FileService.LoadSettingsFile`).

- **Save** — write the current widget state back to the active `.gen.json`,
  or ask for a folder if there is not one yet.
- **Save To** — pick the *folder* to write to. The filename is not yours to
  choose: it is always derived from the template name, and the dialog shows
  the name it will use. An unnamed template cannot be saved — name it on the
  **General** tab first.
- **Load** — load a `.gen.json` file back into the widgets.
- **New** — reset to defaults.

The toolbar's right-hand label always shows the active `.gen.json` path
and an asterisk when there are unsaved changes.

## 4. Generate a Template

1. The output folder is pre-filled from your Steam install when it can be
   found; otherwise pick one in the preview panel (`Browse`).
2. Click **Generate** — this builds the template in memory and refreshes the
   preview panel.
3. Click **Save Template** — writes `<TemplateName>.rmg.json` plus a preview
   `<TemplateName>.png` into the chosen folder.
4. **Reveal** opens the output folder in the app's own browse dialog.

Drop the resulting `.rmg.json` into the game's templates directory and
pick it from the in-game Random Map Generator screen.

## 5. Building Another Front-End

`internal/` is not a public API — Go's `internal` rule means no other module
can import it. It exists so that *this* repository can grow front-ends beyond
the Gio GUI: [app/tui/](app/tui) and [app/web/](app/web) are placeholders for
exactly that.

Every front-end talks to the same seam. The composition root builds the whole
object graph and hands back a single interface:

```go
package main

import (
    "log"

    "github.com/Tariomka/hommoe_custom_templates/internal/composition"
    "github.com/Tariomka/hommoe_custom_templates/internal/dtos"
)

func main() {
    backend := composition.InitializeGuiHandler()

    state := dtos.NewDefaultEditorStateDto()
    state.TemplateName = "Programmatic Map"
    state.PlayerCount = 4
    state.MapSize = 160

    loaded, err := backend.GenerateTemplate(state)
    if err != nil {
        log.Fatal(err)
    }

    savedPath, err := backend.SaveTemplate(dtos.TemplateSaveDto{
        Template:   loaded.Template,
        Topology:   state.Topology,
        OutputPath: ".",
    })
    if err != nil {
        log.Fatal(err)
    }
    log.Println("wrote", savedPath)
}
```

`handler_interfaces.IGuiHandler` is the whole contract a front-end needs:

| Interface             | What it covers                                    |
|-----------------------|---------------------------------------------------|
| `ITemplateHandler`    | generate, update, re-apply castle settings, save  |
| `IStateHandler`       | validate, load and save `.gen.json` editor state  |
| `IPreviewHandler`     | preview layout and PNG rendering                  |
| `IContentRuleHandler` | per-row zone-content placement rules and catalogs |
| `IZoneEditorHandler`  | manual zone and connection editing                |

Rules for a new front-end: it may only render and collect input, it exchanges
`internal/dtos` types with the handlers, and it never constructs services
itself — add providers to [internal/composition](internal/composition) and
regenerate the injector instead. `dtos.NewDefaultEditorStateDto()` matches the
GUI's defaults (2 players, size 160, topology Random, Classic mode, etc.).

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

## 7. Topologies

See the [topology table in the README](README.md#topologies) for the full list
and the shape each one produces.

## 8. Troubleshooting

- **Window doesn't open** — Gio needs a working graphics stack. If
  running over Remote Desktop or in a headless container, run on a real
  desktop session instead. On Linux, make sure the development packages listed
  under [Building on Linux](#building-on-linux) are installed.
- **`go build` complains about Go version** — install Go 1.27.0 or later
  (`go version`).
- **`Save Template` is disabled** — click **Generate** first; the button only
  enables once a template is in memory.
- **Nothing happens after Generate** — check the preview panel's status line
  for errors and verify the template name is non-empty.
- **Template won't load in the game** — verify the JSON is valid
  (`Get-Content map.rmg.json | python -m json.tool`) and compare against
  files in [data/ExampleTemplates/](data/ExampleTemplates).

## 9. Reference

- [README.md](README.md) — architecture, project layout, build/test commands.
- [data/ExampleTemplates/](data/ExampleTemplates) — 57 reference templates
  shipped with the game.
