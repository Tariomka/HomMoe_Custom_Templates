package constants

import "github.com/Tariomka/hommoe_custom_templates/internal/models"

// ContentItemGroup categorizes models.SidMappings by zone-content kind.
// Mirrors Services/ContentManagement/ContentItemGroup.cs.
//
//nolint:gochecknoglobals // semantic registry
var ContentItemGroup = struct {
	Mines         []models.SidMapping
	Treasures     []models.SidMapping
	HireBuildings []models.SidMapping
	ResourceBanks []models.SidMapping
}{
	Mines: []models.SidMapping{
		ContentIds.MineWood,
		ContentIds.MineOre,
		ContentIds.MineGold,
		ContentIds.MineMercury,
		ContentIds.MineCrystals,
		ContentIds.MineGemstones,
		ContentIds.AlchemyLab,
	},
	Treasures: []models.SidMapping{
		ContentIds.MythicScrollBox,
		ContentIds.PandoraBox,
		ContentIds.RandomItemCommon,
		ContentIds.RandomItemEpic,
		ContentIds.RandomItemLegendary,
	},
	HireBuildings: []models.SidMapping{
		ContentIds.RandomHire1,
		ContentIds.RandomHire2,
		ContentIds.RandomHire3,
		ContentIds.RandomHire4,
		ContentIds.RandomHire5,
		ContentIds.RandomHire6,
		ContentIds.RandomHire7,
		IncludeListIds.RandomHiresLowTier,
		IncludeListIds.RandomHiresHighTier,
		IncludeListIds.RandomHiresAllTier,
	},
	ResourceBanks: []models.SidMapping{
		IncludeListIds.ResourceBanksTier1,
		IncludeListIds.ResourceBanksTier2,
	},
}
