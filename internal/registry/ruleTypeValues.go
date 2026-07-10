package registry

type ruleTypes struct {
	Crossroads string
	MainObject string
	Road       string
	Sid        string
}

// GetRuleTypeValues returns the rule types used for
//
//	mandatoryContent.content.rules.type
func GetRuleTypeValues() ruleTypes {
	return ruleTypes{
		Crossroads: "Crossroads",
		MainObject: "MainObject",
		Road:       "Road",
		Sid:        "Sid",
	}
}
