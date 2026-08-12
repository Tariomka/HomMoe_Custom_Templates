package common_connections

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

func GetGuardWeeklyIncrements() models.GuardWeeklyIncrement {
	return models.GuardWeeklyIncrement{
		Slow:     0.05,
		Normal:   0.10,
		Standard: 0.15,
		Fast:     0.20,
		VeryFast: 0.25,
	}
}

func GetGuardWeeklyIncrementList() []data.Tuple[string, float64] {
	increments := GetGuardWeeklyIncrements()
	return []data.Tuple[string, float64]{
		data.NewTuple("Slow (5%)", increments.Slow),
		data.NewTuple("Normal (10%)", increments.Normal),
		data.NewTuple("Standard (15%)", increments.Standard),
		data.NewTuple("Fast (20%)", increments.Fast),
		data.NewTuple("Very Fast (25%)", increments.VeryFast),
	}
}

func GetGuardWeeklyIncrementLabels() []string {
	incrementList := GetGuardWeeklyIncrementList()
	return linq.FromSlice(incrementList).
		SelectString(func(t data.Tuple[string, float64]) string { return t.Key }).
		ToSlice()
}

func GetGuardWeeklyIncrementValues() []float64 {
	increments := GetGuardWeeklyIncrements()
	return []float64{
		increments.Slow,
		increments.Normal,
		increments.Standard,
		increments.Fast,
		increments.VeryFast,
	}
}
