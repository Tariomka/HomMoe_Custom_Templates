# Personal findings and manual backlog

---

Creating a Geometric Hub topology template, saving the editor and loading that state
produces a completely incorrect preview in the preview panel and in the preview png.
A pre-made tested editor state and its output before reloading can be found in
[Topology when loaded](../../output/research/Topology%20when%20loaded).

---

App portal connections are not being rendered in preview png - such an example can be found in
[Colosseum v3.png](../../output/research/Topology%20when%20loaded/Colosseum%20v3.png).
I haven't checked if this is applicable for all topologies or just Geometric Hub topology
(it was only tested on Geometric Hub).

---

Last time I tried, Bonuses and Bans did not apply in the generated template in the game.
More testing is required.

---

Some of the custom connections are not being applied even though it is rendered in the manual
editor dialog and the preview panel - this was observed on the preview png.
Easiest way to reproduce this is to add an additional connection between 2 zones that already have
a connection - the arched outer connection will not appear in the png file.

---

Need to create some sort of OCR tool for agents and me to use when the need arises to debug
integration test goldens vs failure images.

---

Need to block the ability to change player count when Tournament mode is selected. It should always
use 2 players only. Most of the current tournament topology variants will need to be removed, new
tournament specific topologies will need to be added and the tournament topology provider will
need to be reworked.

---

Need to check if hero hire ban is enabled when Single Hero, FinalBattle/Gladiator Arena or
Lost Start Hero options are activated. If it is not, need to fix it - this is a bug because in-game
this is the option that blocks players from recruiting additional heros.

---

Need to remove old obsolete topologies (Ring, Hub, Chain and Shared Web) these topologies do not add
any value and are not balanced.

---

Need to add Zone Content presets to be able to be applied, preferably on top of existing contents,
in the zone content editor dialogs for neutral zones. Maybe the same or a separate set of presets
should be added for player zones as well.

---

In Manual Zone Editor dialog, when a neutral zone quality is changed, it should recalculate and
reapply connection guard strength value between the edited zone and all connected zones. The
reapplication should be done on both value and preset - the value should be the same as the guard
preset, and the guard preset should be of the same tier as it was before the change, for example:
if changing from Silver to Gold, the preset should stay the same (Default (20000) -> Default (25000),
Medium (24000) -> Medium (48000), etc.) and the value should be updated appropriately (20000 -> 25000,
24000 -> 48000, etc.)

---

Need to optimize the editor state object and the values saved in the .gen.json - for example,
currently all of the connections and zones are saved if a single manual change occurs, even though
only the manual edits should be tracked (if a single zone is removed/changed/added, that one change
should be tracked, all other zones and connection should be applied from the topology).
Also need to check what values are being stored that can be extrapolated from other properties -
for example guardReactionDistribution, guardContentPool, unguardedContentPool, contentCountLimits
and other values never change within a zone and zone type can easily be extrapolated, default
content rows are always the same and only modifications could be tracked, etc.

---

BannedItems, BannedMagics and ValueOverrides should be lists, not strings in the .gen.json, just
like Bonuses.

---

Eventually need to make the editor state in the .gen.json not a flat option mashup, but sectioned,
possibly using the current separated structure of the editor state and the sectioned structure in
the json.

---
