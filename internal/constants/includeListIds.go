package constants

import "github.com/Tariomka/hommoe_custom_templates/internal/models"

// IncludeListIds enumerates the named include-list SIDs.
//
//nolint:gochecknoglobals // semantic registry
var IncludeListIds = struct {
	RandomHiresLowTier  models.SidMapping
	RandomHiresHighTier models.SidMapping
	RandomHiresAllTier  models.SidMapping
	ResourceBanksTier1  models.SidMapping
	ResourceBanksTier2  models.SidMapping
}{
	RandomHiresLowTier:  models.SidMapping{Sid: "content_list_building_random_hires_low_tier", Name: "Random Hires Low Tier"},
	RandomHiresHighTier: models.SidMapping{Sid: "content_list_building_random_hires_high_tier", Name: "Random Hires High Tier"},
	RandomHiresAllTier:  models.SidMapping{Sid: "basic_content_list_building_random_hires", Name: "Random Hires All Tier"},
	ResourceBanksTier1:  models.SidMapping{Sid: "basic_content_list_building_guarded_resource_banks_tier_1", Name: "Resource Banks T1"},
	ResourceBanksTier2:  models.SidMapping{Sid: "basic_content_list_building_guarded_resource_banks_tier_2", Name: "Resource Banks T2"},
}
