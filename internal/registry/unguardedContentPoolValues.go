package registry

type unguardedContentPools struct {
	Base     string
	Item     string
	Pandora  string
	Hire     string
	UnitBank string
	ResBank  string
	Stat     string
	Magic    string
}

var unguardedT0PoolValues = unguardedContentPools{
	Base:     "classic_template_pool_random_unguarded_t0_base",
	Item:     "classic_template_pool_random_unguarded_t0_item",
	Pandora:  "classic_template_pool_random_unguarded_t0_pandora",
	Hire:     "classic_template_pool_random_unguarded_t0_hire",
	UnitBank: "classic_template_pool_random_unguarded_t0_unit_bank",
	ResBank:  "classic_template_pool_random_unguarded_t0_res_bank",
	Stat:     "classic_template_pool_random_unguarded_t0_stat",
	Magic:    "classic_template_pool_random_unguarded_t0_magic",
}
var unguardedT1PoolValues = unguardedContentPools{
	Base:     "classic_template_pool_random_unguarded_t1_base",
	Item:     "classic_template_pool_random_unguarded_t1_item",
	Pandora:  "classic_template_pool_random_unguarded_t1_pandora",
	Hire:     "classic_template_pool_random_unguarded_t1_hire",
	UnitBank: "classic_template_pool_random_unguarded_t1_unit_bank",
	ResBank:  "classic_template_pool_random_unguarded_t1_res_bank",
	Stat:     "classic_template_pool_random_unguarded_t1_stat",
	Magic:    "classic_template_pool_random_unguarded_t1_magic",
}
var unguardedT2PoolValues = unguardedContentPools{
	Base:     "classic_template_pool_random_unguarded_t2_base",
	Item:     "classic_template_pool_random_unguarded_t2_item",
	Pandora:  "classic_template_pool_random_unguarded_t2_pandora",
	Hire:     "classic_template_pool_random_unguarded_t2_hire",
	UnitBank: "classic_template_pool_random_unguarded_t2_unit_bank",
	ResBank:  "classic_template_pool_random_unguarded_t2_res_bank",
	Stat:     "classic_template_pool_random_unguarded_t2_stat",
	Magic:    "classic_template_pool_random_unguarded_t2_magic",
}
var unguardedT3PoolValues = unguardedContentPools{
	Base:     "classic_template_pool_random_unguarded_t3_base",
	Item:     "classic_template_pool_random_unguarded_t3_item",
	Pandora:  "classic_template_pool_random_unguarded_t3_pandora",
	Hire:     "classic_template_pool_random_unguarded_t3_hire",
	UnitBank: "classic_template_pool_random_unguarded_t3_unit_bank",
	ResBank:  "classic_template_pool_random_unguarded_t3_res_bank",
	Stat:     "classic_template_pool_random_unguarded_t3_stat",
	Magic:    "classic_template_pool_random_unguarded_t3_magic",
}
var unguardedT4PoolValues = unguardedContentPools{
	Base:     "classic_template_pool_random_unguarded_t4_base",
	Item:     "classic_template_pool_random_unguarded_t4_item",
	Pandora:  "classic_template_pool_random_unguarded_t4_pandora",
	Hire:     "classic_template_pool_random_unguarded_t4_hire",
	UnitBank: "classic_template_pool_random_unguarded_t4_unit_bank",
	ResBank:  "classic_template_pool_random_unguarded_t4_res_bank",
	Stat:     "classic_template_pool_random_unguarded_t4_stat",
	Magic:    "classic_template_pool_random_unguarded_t4_magic",
}
var unguardedT5PoolValues = unguardedContentPools{
	Base:     "classic_template_pool_random_unguarded_t5_base",
	Item:     "classic_template_pool_random_unguarded_t5_item",
	Pandora:  "classic_template_pool_random_unguarded_t5_pandora",
	Hire:     "classic_template_pool_random_unguarded_t5_hire",
	UnitBank: "classic_template_pool_random_unguarded_t5_unit_bank",
	ResBank:  "classic_template_pool_random_unguarded_t5_res_bank",
	Stat:     "classic_template_pool_random_unguarded_t5_stat",
	Magic:    "classic_template_pool_random_unguarded_t5_magic",
}

// GetUnguardedContentPoolT0List returns a list of all unguarded tier 0 content pool IDs (excluding the base pool) used for
//
//	variants.zones.unguardedContentPool
func GetUnguardedContentPoolT0List() []string {
	return []string{
		unguardedT0PoolValues.Item,
		unguardedT0PoolValues.Pandora,
		unguardedT0PoolValues.Hire,
		unguardedT0PoolValues.UnitBank,
		unguardedT0PoolValues.ResBank,
		unguardedT0PoolValues.Stat,
		unguardedT0PoolValues.Magic,
	}
}

// GetUnguardedContentPoolT1List returns a list of all unguarded tier 1 content pool IDs (excluding the base pool) used for
//
//	variants.zones.unguardedContentPool
func GetUnguardedContentPoolT1List() []string {
	return []string{
		unguardedT1PoolValues.Item,
		unguardedT1PoolValues.Pandora,
		unguardedT1PoolValues.Hire,
		unguardedT1PoolValues.UnitBank,
		unguardedT1PoolValues.ResBank,
		unguardedT1PoolValues.Stat,
		unguardedT1PoolValues.Magic,
	}
}

// GetUnguardedContentPoolT2List returns a list of all unguarded tier 2 content pool IDs (excluding the base pool) used for
//
//	variants.zones.unguardedContentPool
func GetUnguardedContentPoolT2List() []string {
	return []string{
		unguardedT2PoolValues.Item,
		unguardedT2PoolValues.Pandora,
		unguardedT2PoolValues.Hire,
		unguardedT2PoolValues.UnitBank,
		unguardedT2PoolValues.ResBank,
		unguardedT2PoolValues.Stat,
		unguardedT2PoolValues.Magic,
	}
}

// GetUnguardedContentPoolT3List returns a list of all unguarded tier 3 content pool IDs (excluding the base pool) used for
//
//	variants.zones.unguardedContentPool
func GetUnguardedContentPoolT3List() []string {
	return []string{
		unguardedT3PoolValues.Item,
		unguardedT3PoolValues.Pandora,
		unguardedT3PoolValues.Hire,
		unguardedT3PoolValues.UnitBank,
		unguardedT3PoolValues.ResBank,
		unguardedT3PoolValues.Stat,
		unguardedT3PoolValues.Magic,
	}
}

// GetUnguardedContentPoolT4List returns a list of all unguarded tier 4 content pool IDs (excluding the base pool) used for
//
//	variants.zones.unguardedContentPool
func GetUnguardedContentPoolT4List() []string {
	return []string{
		unguardedT4PoolValues.Item,
		unguardedT4PoolValues.Pandora,
		unguardedT4PoolValues.Hire,
		unguardedT4PoolValues.UnitBank,
		unguardedT4PoolValues.ResBank,
		unguardedT4PoolValues.Stat,
		unguardedT4PoolValues.Magic,
	}
}

// GetUnguardedContentPoolT5List returns a list of all unguarded tier 5 content pool IDs (excluding the base pool) used for
//
//	variants.zones.unguardedContentPool
func GetUnguardedContentPoolT5List() []string {
	return []string{
		unguardedT5PoolValues.Item,
		unguardedT5PoolValues.Pandora,
		unguardedT5PoolValues.Hire,
		unguardedT5PoolValues.UnitBank,
		unguardedT5PoolValues.ResBank,
		unguardedT5PoolValues.Stat,
		unguardedT5PoolValues.Magic,
	}
}
