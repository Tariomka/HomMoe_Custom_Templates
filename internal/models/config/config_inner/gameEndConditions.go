package config_inner

type GameEndConditions struct {
	VictoryCondition string // "win_condition_1"..._6
	LostStartCity    bool
	LostStartCityDay int
	LostStartHero    bool
	CityHold         bool
	CityHoldDays     int
}
