package registry

type ruleTypes struct {
	Crossroads string
	MainObject string
	Road       string
	Sid        string
}

var ruleTypeValues = ruleTypes{
	Crossroads: "Crossroads",
	MainObject: "MainObject",
	Road:       "Road",
	Sid:        "Sid",
}

// GetRuleTypeValues returns the rule types used for
//
//	mandatoryContent.content.rules.type
func GetRuleTypeValues() ruleTypes {
	return ruleTypeValues
}
