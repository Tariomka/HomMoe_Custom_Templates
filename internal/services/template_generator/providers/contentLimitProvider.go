package providers

import (
	"fmt"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

type ContentLimitProvider struct{}

func NewContentLimitProvider() *ContentLimitProvider {
	return &ContentLimitProvider{}
}

func (this *ContentLimitProvider) CreateContentCountLimits(settings config.GeneratorConfig) []template.ContentCountLimit {
	sidLimits := this.createDefaultContentLimits()

	// Lift limits when any mandatory-content list (player or neutral or
	// hub) requests more of a given SID than the default cap.
	sidCounts := map[string]int{}
	tally := func(items []template.MandatoryContentItem) {
		for _, item := range items {
			if item.SID != "" {
				sidCounts[strings.ToLower(item.SID)]++
			}
		}
	}
	tally(settings.PlayerZoneMandatoryContent)
	tally(settings.LowNeutralMandatoryContent)
	tally(settings.MediumNeutralMandatoryContent)
	tally(settings.HighNeutralMandatoryContent)
	tally(settings.HubZoneMandatoryContent)
	for i := range sidLimits {
		if count, ok := sidCounts[strings.ToLower(sidLimits[i].SID)]; ok {
			if count > sidLimits[i].MaxCount {
				sidLimits[i].MaxCount = count
			}
		}
	}

	var limits []template.ContentCountLimit
	limits = append(limits, template.ContentCountLimit{Name: "content_limits_side", Limits: sidLimits})
	limits = append(limits, template.ContentCountLimit{Name: "content_limits_side_0_0", Limits: sidLimits})
	for a := 1; a <= 5; a++ {
		for b := a + 1; b <= 6; b++ {
			limits = append(limits, template.ContentCountLimit{
				Name:   fmt.Sprintf("content_limits_side_%d_%d", a, b),
				Limits: sidLimits,
			})
		}
	}
	return limits
}

func (this *ContentLimitProvider) createDefaultContentLimits() []template.ContentLimit {
	return []template.ContentLimit{
		{SID: "black_tower", MaxCount: 0},
		{SID: constants.ContentIds.Fountain.Sid, MaxCount: 2},
		{SID: constants.ContentIds.Fountain2.Sid, MaxCount: 2},
		{SID: constants.ContentIds.ManaWell.Sid, MaxCount: 2},
		{SID: constants.ContentIds.BeerFountain.Sid, MaxCount: 2},
		{SID: constants.ContentIds.Market.Sid, MaxCount: 1},
		{SID: constants.ContentIds.Forge.Sid, MaxCount: 2},
		{SID: constants.ContentIds.Stables.Sid, MaxCount: 1},
		{SID: constants.ContentIds.Watchtower.Sid, MaxCount: 2},
		{SID: constants.ContentIds.WindRose.Sid, MaxCount: 1},
		{SID: constants.ContentIds.QuixsPath.Sid, MaxCount: 2},
		{SID: constants.ContentIds.CrystalTrail.Sid, MaxCount: 3},
		{SID: constants.ContentIds.MysteriousStone.Sid, MaxCount: 2},
		{SID: constants.ContentIds.University.Sid, MaxCount: 2},
		{SID: constants.ContentIds.WiseOwl.Sid, MaxCount: 4},
		{SID: constants.ContentIds.CelestialSphere.Sid, MaxCount: 2},
		{SID: constants.ContentIds.PileOfBooks.Sid, MaxCount: 2},
		{SID: constants.ContentIds.InsarasEye.Sid, MaxCount: 2},
		{SID: constants.ContentIds.TearOfTruth.Sid, MaxCount: 3},
		{SID: constants.ContentIds.TreeOfAbundance.Sid, MaxCount: 2},
		{SID: constants.ContentIds.HuntsmansCamp.Sid, MaxCount: 2},
		{SID: constants.ContentIds.ShadyDen.Sid, MaxCount: 2},
		{SID: constants.ContentIds.RandomHire1.Sid, MaxCount: 6},
		{SID: constants.ContentIds.RandomHire2.Sid, MaxCount: 6},
		{SID: constants.ContentIds.RandomHire3.Sid, MaxCount: 6},
		{SID: constants.ContentIds.RandomHire4.Sid, MaxCount: 6},
		{SID: constants.ContentIds.RandomHire5.Sid, MaxCount: 6},
		{SID: constants.ContentIds.RandomHire6.Sid, MaxCount: 6},
		{SID: constants.ContentIds.RandomHire7.Sid, MaxCount: 6},
		{SID: constants.ContentIds.Arena.Sid, MaxCount: 2},
		{SID: constants.ContentIds.SacrificialShrine.Sid, MaxCount: 2},
		{SID: constants.ContentIds.Chimerologist.Sid, MaxCount: 2},
		{SID: constants.ContentIds.Circus.Sid, MaxCount: 2},
		{SID: constants.ContentIds.InfernalCirque.Sid, MaxCount: 2},
		{SID: constants.ContentIds.FlatteringMirror.Sid, MaxCount: 2},
		{SID: constants.ContentIds.FickleShrine.Sid, MaxCount: 1},
		{SID: constants.ContentIds.PointOfBalance.Sid, MaxCount: 3},
		{SID: constants.ContentIds.PandoraBox.Sid, MaxCount: 4},
		{SID: constants.ContentIds.RitualPyre.Sid, MaxCount: 3},
		{SID: constants.ContentIds.BorealCall.Sid, MaxCount: 3},
		{SID: constants.ContentIds.JoustingRange.Sid, MaxCount: 1},
		{SID: constants.ContentIds.UnforgottenGrave.Sid, MaxCount: 1},
		{SID: constants.ContentIds.PetrifiedMemorial.Sid, MaxCount: 1},
		{SID: constants.ContentIds.TheGorge.Sid, MaxCount: 1},
	}
}
