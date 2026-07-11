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
	constants.ContentIDs.AlchemyLab,
	constants.ContentIDs.Arena,
	constants.ContentIDs.BeerFountain,
	constants.ContentIDs.BorealCall,
	constants.ContentIDs.CelestialSphere,
	constants.ContentIDs.Chimerologist,
	constants.ContentIDs.Circus,
	constants.ContentIDs.CollegeOfWonder,
	constants.ContentIDs.CrystalTrail,
	constants.ContentIDs.DragonUtopia,
	constants.ContentIDs.EternalDragon,
	constants.ContentIDs.FickleShrine,
	constants.ContentIDs.FlatteringMirror,
	constants.ContentIDs.Forge,
	constants.ContentIDs.Fort,
	constants.ContentIDs.Fountain,
	constants.ContentIDs.Fountain2,
	constants.ContentIDs.HuntsmansCamp,
	constants.ContentIDs.InfernalCirque,
	constants.ContentIDs.InsarasEye,
	constants.ContentIDs.JoustingRange,
	constants.ContentIDs.ManaWell,
	constants.ContentIDs.Market,
	constants.ContentIDs.MineCrystals,
	constants.ContentIDs.MineGemstones,
	constants.ContentIDs.MineGold,
	constants.ContentIDs.MineMercury,
	constants.ContentIDs.MineOre,
	constants.ContentIDs.MineWood,
	constants.ContentIDs.Mirage,
	constants.ContentIDs.MontyHall,
	constants.ContentIDs.MysteriousStone,
	constants.ContentIDs.MysticalTower,
	constants.ContentIDs.MythicScrollBox,
	constants.ContentIDs.OrbObservatory,
	constants.ContentIDs.PandoraBox,
	constants.ContentIDs.PetrifiedMemorial,
	constants.ContentIDs.PileOfBooks,
	constants.ContentIDs.PointOfBalance,
	constants.ContentIDs.Prison,
	constants.ContentIDs.QuixsPath,
	constants.ContentIDs.RandomHire1,
	constants.ContentIDs.RandomHire2,
	constants.ContentIDs.RandomHire3,
	constants.ContentIDs.RandomHire4,
	constants.ContentIDs.RandomHire5,
	constants.ContentIDs.RandomHire6,
	constants.ContentIDs.RandomHire7,
	constants.ContentIDs.RandomItemCommon,
	constants.ContentIDs.RandomItemEpic,
	constants.ContentIDs.RandomItemLegendary,
	constants.ContentIDs.RandomItemRare,
	constants.ContentIDs.RemoteFoothold,
	constants.ContentIDs.ResearchLaboratory,
	constants.ContentIDs.RitualPyre,
	constants.ContentIDs.SacrificialShrine,
	constants.ContentIDs.ShadyDen,
	constants.ContentIDs.Stables,
	constants.ContentIDs.Tavern,
	constants.ContentIDs.TearOfTruth,
	constants.ContentIDs.TheGorge,
	constants.ContentIDs.TownGate,
	constants.ContentIDs.TreeOfAbundance,
	constants.ContentIDs.TroglodyteThrone,
	constants.ContentIDs.UnforgottenGrave,
	constants.ContentIDs.University,
	constants.ContentIDs.UnstableRuins,
	constants.ContentIDs.Watchtower,
	constants.ContentIDs.WindRose,
	constants.ContentIDs.WiseOwl,

	constants.IncludeListIDs.RandomHiresLowTier,
	constants.IncludeListIDs.RandomHiresHighTier,
	constants.IncludeListIDs.RandomHiresAllTier,
	constants.IncludeListIDs.ResourceBanksTier1,
	constants.IncludeListIDs.ResourceBanksTier2,
}
