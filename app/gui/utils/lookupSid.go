package utils

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// GetSidMappingByName returns the models.SidMapping whose Name matches name (case-insensitive).
func GetSidMappingByName(name string) (models.SidMapping, bool) {
	for _, mapping := range allSidMappings {
		if strings.EqualFold(mapping.Name, name) {
			return mapping, true
		}
	}
	return models.SidMapping{}, false
}

// GetSidMappingBySid returns the models.SidMapping whose Sid matches sid.
func GetSidMappingBySid(sid string) (models.SidMapping, bool) {
	for _, mapping := range allSidMappings {
		if mapping.Sid == sid {
			return mapping, true
		}
	}
	return models.SidMapping{}, false
}

var allSidMappings = []models.SidMapping{
	constants.ContentIds.AlchemyLab,
	constants.ContentIds.Arena,
	constants.ContentIds.BeerFountain,
	constants.ContentIds.BorealCall,
	constants.ContentIds.CelestialSphere,
	constants.ContentIds.Chimerologist,
	constants.ContentIds.Circus,
	constants.ContentIds.CollegeOfWonder,
	constants.ContentIds.CrystalTrail,
	constants.ContentIds.DragonUtopia,
	constants.ContentIds.EternalDragon,
	constants.ContentIds.FickleShrine,
	constants.ContentIds.FlatteringMirror,
	constants.ContentIds.Forge,
	constants.ContentIds.Fort,
	constants.ContentIds.Fountain,
	constants.ContentIds.Fountain2,
	constants.ContentIds.HuntsmansCamp,
	constants.ContentIds.InfernalCirque,
	constants.ContentIds.InsarasEye,
	constants.ContentIds.JoustingRange,
	constants.ContentIds.ManaWell,
	constants.ContentIds.Market,
	constants.ContentIds.MineCrystals,
	constants.ContentIds.MineGemstones,
	constants.ContentIds.MineGold,
	constants.ContentIds.MineMercury,
	constants.ContentIds.MineOre,
	constants.ContentIds.MineWood,
	constants.ContentIds.Mirage,
	constants.ContentIds.MontyHall,
	constants.ContentIds.MysteriousStone,
	constants.ContentIds.MysticalTower,
	constants.ContentIds.MythicScrollBox,
	constants.ContentIds.OrbObservatory,
	constants.ContentIds.PandoraBox,
	constants.ContentIds.PetrifiedMemorial,
	constants.ContentIds.PileOfBooks,
	constants.ContentIds.PointOfBalance,
	constants.ContentIds.Prison,
	constants.ContentIds.QuixsPath,
	constants.ContentIds.RandomHire1,
	constants.ContentIds.RandomHire2,
	constants.ContentIds.RandomHire3,
	constants.ContentIds.RandomHire4,
	constants.ContentIds.RandomHire5,
	constants.ContentIds.RandomHire6,
	constants.ContentIds.RandomHire7,
	constants.ContentIds.RandomItemCommon,
	constants.ContentIds.RandomItemEpic,
	constants.ContentIds.RandomItemLegendary,
	constants.ContentIds.RandomItemRare,
	constants.ContentIds.RemoteFoothold,
	constants.ContentIds.ResearchLaboratory,
	constants.ContentIds.RitualPyre,
	constants.ContentIds.SacrificialShrine,
	constants.ContentIds.ShadyDen,
	constants.ContentIds.Stables,
	constants.ContentIds.Tavern,
	constants.ContentIds.TearOfTruth,
	constants.ContentIds.TheGorge,
	constants.ContentIds.TownGate,
	constants.ContentIds.TreeOfAbundance,
	constants.ContentIds.TroglodyteThrone,
	constants.ContentIds.UnforgottenGrave,
	constants.ContentIds.University,
	constants.ContentIds.UnstableRuins,
	constants.ContentIds.Watchtower,
	constants.ContentIds.WindRose,
	constants.ContentIds.WiseOwl,

	constants.IncludeListIds.RandomHiresLowTier,
	constants.IncludeListIds.RandomHiresHighTier,
	constants.IncludeListIds.RandomHiresAllTier,
	constants.IncludeListIds.ResourceBanksTier1,
	constants.IncludeListIds.ResourceBanksTier2,
}
