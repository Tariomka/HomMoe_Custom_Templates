package registry

type guardedContentPools struct {
	Base     string
	Item     string
	Pandora  string
	Hire     string
	UnitBank string
	ResBank  string
	Stat     string
	Magic    string
}

var guardedT0PoolValues = guardedContentPools{
	Base:     "classic_template_pool_random_t0_base",
	Item:     "classic_template_pool_random_t0_item",
	Pandora:  "classic_template_pool_random_t0_pandora",
	Hire:     "classic_template_pool_random_t0_hire",
	UnitBank: "classic_template_pool_random_t0_unit_bank",
	ResBank:  "classic_template_pool_random_t0_res_bank",
	Stat:     "classic_template_pool_random_t0_stat",
	Magic:    "classic_template_pool_random_t0_magic",
}
var guardedT1PoolValues = guardedContentPools{
	Base:     "classic_template_pool_random_t1_base",
	Item:     "classic_template_pool_random_t1_item",
	Pandora:  "classic_template_pool_random_t1_pandora",
	Hire:     "classic_template_pool_random_t1_hire",
	UnitBank: "classic_template_pool_random_t1_unit_bank",
	ResBank:  "classic_template_pool_random_t1_res_bank",
	Stat:     "classic_template_pool_random_t1_stat",
	Magic:    "classic_template_pool_random_t1_magic",
}
var guardedT2PoolValues = guardedContentPools{
	Base:     "classic_template_pool_random_t2_base",
	Item:     "classic_template_pool_random_t2_item",
	Pandora:  "classic_template_pool_random_t2_pandora",
	Hire:     "classic_template_pool_random_t2_hire",
	UnitBank: "classic_template_pool_random_t2_unit_bank",
	ResBank:  "classic_template_pool_random_t2_res_bank",
	Stat:     "classic_template_pool_random_t2_stat",
	Magic:    "classic_template_pool_random_t2_magic",
}
var guardedT3PoolValues = guardedContentPools{
	Base:     "classic_template_pool_random_t3_base",
	Item:     "classic_template_pool_random_t3_item",
	Pandora:  "classic_template_pool_random_t3_pandora",
	Hire:     "classic_template_pool_random_t3_hire",
	UnitBank: "classic_template_pool_random_t3_unit_bank",
	ResBank:  "classic_template_pool_random_t3_res_bank",
	Stat:     "classic_template_pool_random_t3_stat",
	Magic:    "classic_template_pool_random_t3_magic",
}
var guardedT4PoolValues = guardedContentPools{
	Base:     "classic_template_pool_random_t4_base",
	Item:     "classic_template_pool_random_t4_item",
	Pandora:  "classic_template_pool_random_t4_pandora",
	Hire:     "classic_template_pool_random_t4_hire",
	UnitBank: "classic_template_pool_random_t4_unit_bank",
	ResBank:  "classic_template_pool_random_t4_res_bank",
	Stat:     "classic_template_pool_random_t4_stat",
	Magic:    "classic_template_pool_random_t4_magic",
}
var guardedT5PoolValues = guardedContentPools{
	Base:     "classic_template_pool_random_t5_base",
	Item:     "classic_template_pool_random_t5_item",
	Pandora:  "classic_template_pool_random_t5_pandora",
	Hire:     "classic_template_pool_random_t5_hire",
	UnitBank: "classic_template_pool_random_t5_unit_bank",
	ResBank:  "classic_template_pool_random_t5_res_bank",
	Stat:     "classic_template_pool_random_t5_stat",
	Magic:    "classic_template_pool_random_t5_magic",
}

// GetGuardedContentPoolT0List returns a list of all guarded tier 0 content pool IDs (excluding the base pool) used for
//
//	variants.zones.guardedContentPool
func GetGuardedContentPoolT0List() []string {
	return []string{
		guardedT0PoolValues.Item,
		guardedT0PoolValues.Pandora,
		guardedT0PoolValues.Hire,
		guardedT0PoolValues.UnitBank,
		guardedT0PoolValues.ResBank,
		guardedT0PoolValues.Stat,
		guardedT0PoolValues.Magic,
	}
}

// GetGuardedContentPoolT1List returns a list of all guarded tier 1 content pool IDs (excluding the base pool) used for
//
//	variants.zones.guardedContentPool
func GetGuardedContentPoolT1List() []string {
	return []string{
		guardedT1PoolValues.Item,
		guardedT1PoolValues.Pandora,
		guardedT1PoolValues.Hire,
		guardedT1PoolValues.UnitBank,
		guardedT1PoolValues.ResBank,
		guardedT1PoolValues.Stat,
		guardedT1PoolValues.Magic,
	}
}

// GetGuardedContentPoolT2List returns a list of all guarded tier 2 content pool IDs (excluding the base pool) used for
//
//	variants.zones.guardedContentPool
func GetGuardedContentPoolT2List() []string {
	return []string{
		guardedT2PoolValues.Item,
		guardedT2PoolValues.Pandora,
		guardedT2PoolValues.Hire,
		guardedT2PoolValues.UnitBank,
		guardedT2PoolValues.ResBank,
		guardedT2PoolValues.Stat,
		guardedT2PoolValues.Magic,
	}
}

// GetGuardedContentPoolT3List returns a list of all guarded tier 3 content pool IDs (excluding the base pool) used for
//
//	variants.zones.guardedContentPool
func GetGuardedContentPoolT3List() []string {
	return []string{
		guardedT3PoolValues.Item,
		guardedT3PoolValues.Pandora,
		guardedT3PoolValues.Hire,
		guardedT3PoolValues.UnitBank,
		guardedT3PoolValues.ResBank,
		guardedT3PoolValues.Stat,
		guardedT3PoolValues.Magic,
	}
}

// GetGuardedContentPoolT4List returns a list of all guarded tier 4 content pool IDs (excluding the base pool) used for
//
//	variants.zones.guardedContentPool
func GetGuardedContentPoolT4List() []string {
	return []string{
		guardedT4PoolValues.Item,
		guardedT4PoolValues.Pandora,
		guardedT4PoolValues.Hire,
		guardedT4PoolValues.UnitBank,
		guardedT4PoolValues.ResBank,
		guardedT4PoolValues.Stat,
		guardedT4PoolValues.Magic,
	}
}

// GetGuardedContentPoolT5List returns a list of all guarded tier 5 content pool IDs (excluding the base pool) used for
//
//	variants.zones.guardedContentPool
func GetGuardedContentPoolT5List() []string {
	return []string{
		guardedT5PoolValues.Item,
		guardedT5PoolValues.Pandora,
		guardedT5PoolValues.Hire,
		guardedT5PoolValues.UnitBank,
		guardedT5PoolValues.ResBank,
		guardedT5PoolValues.Stat,
		guardedT5PoolValues.Magic,
	}
}
