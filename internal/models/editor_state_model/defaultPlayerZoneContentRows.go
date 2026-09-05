package editor_state_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

const guardedRuleName = "Guarded"

// GetDefaultPlayerZoneContentRows returns the historical default
// mandatory-content rows seeded into every player zone: the six basic mines
// plus an alchemy lab, a couple of guarded treasures, random hires and
// resource banks.
func GetDefaultPlayerZoneContentRows() []ZoneContentRow {
	rows := getPlayerZoneMineRows()
	rows = append(rows, getPlayerZoneTreasureRows()...)
	rows = append(rows, getPlayerZoneHireAndBankRows()...)
	return rows
}

// getPlayerZoneMineRows returns the six basic mines plus an alchemy lab:
// the wood, ore and gold mines near the town, the rest next to a road.
func getPlayerZoneMineRows() []ZoneContentRow {
	interactable := registry.GetMapObjectAllInteractableValues()
	nearTown := ContentRuleRow{Name: "Distance to town", DistanceName: "Near"}
	nextToRoad := ContentRuleRow{Name: "Distance to road", DistanceName: "Next To"}
	return []ZoneContentRow{
		getGuardedMineRow(interactable.WoodMine, nearTown),
		getGuardedMineRow(interactable.OreMine, nearTown),
		getGuardedMineRow(interactable.GoldMine, nearTown),
		getGuardedMineRow(interactable.CrystalMine, nextToRoad),
		getGuardedMineRow(interactable.MercuryMine, nextToRoad),
		getGuardedMineRow(interactable.GemstoneMine, nextToRoad),
		getGuardedMineRow(interactable.AlchemyLab, nextToRoad),
	}
}

// getGuardedMineRow builds a single guarded mine row with the given
// distance rule.
func getGuardedMineRow(sid string, distanceRule ContentRuleRow) ZoneContentRow {
	return ZoneContentRow{
		Sid:    sid,
		Count:  1,
		IsMine: true,
		Rules: []ContentRuleRow{
			{Name: guardedRuleName, IsGuarded: new(true)},
			distanceRule,
		},
	}
}

// getPlayerZoneTreasureRows returns the guarded treasures: a pandora box
// and an epic random item.
func getPlayerZoneTreasureRows() []ZoneContentRow {
	resources := registry.GetMapObjectResourceValues()
	randomItems := registry.GetMapObjectRandomItemValues()

	return []ZoneContentRow{
		{
			Sid:   resources.PandoraBox,
			Count: 1,
			Rules: []ContentRuleRow{{Name: guardedRuleName, IsGuarded: new(true)}},
		},
		{
			Sid:   randomItems.RandomItemEpic,
			Count: 1,
			Rules: []ContentRuleRow{{Name: guardedRuleName, IsGuarded: new(true)}},
		},
	}
}

// getPlayerZoneHireAndBankRows returns the guarded random-hire and
// resource-bank group rows.
func getPlayerZoneHireAndBankRows() []ZoneContentRow {
	randomHires := registry.GetMandatoryContentRandomHiresBuildingValues()
	basicRandomHires := registry.GetMandatoryContentBasicRandomHiresBuildingValues()
	basicResourceBanks := registry.GetMandatoryContentBasicResourceBanksBuildingValues()

	trueVal := true // To not allocate multiple times
	return []ZoneContentRow{
		{
			Sid:     randomHires.RandomHiresLowTier,
			Count:   2,
			IsGroup: true,
			Rules:   []ContentRuleRow{{Name: guardedRuleName, IsGuarded: &trueVal}},
		},
		{
			Sid:     randomHires.RandomHiresHighTier,
			Count:   1,
			IsGroup: true,
			Rules:   []ContentRuleRow{{Name: guardedRuleName, IsGuarded: &trueVal}},
		},
		{
			Sid:     basicRandomHires.BasicRandomHires,
			Count:   1,
			IsGroup: true,
			Rules:   []ContentRuleRow{{Name: guardedRuleName, IsGuarded: &trueVal}},
		},
		{
			Sid:     basicResourceBanks.BasicResourceBanksTier1,
			Count:   2,
			IsGroup: true,
			Rules:   []ContentRuleRow{{Name: guardedRuleName, IsGuarded: &trueVal}},
		},
		{
			Sid:     basicResourceBanks.BasicResourceBanksTier2,
			Count:   1,
			IsGroup: true,
			Rules:   []ContentRuleRow{{Name: guardedRuleName, IsGuarded: &trueVal}},
		},
	}
}
