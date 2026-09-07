package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type TemplateMapper struct{}

func NewTemplateMapper() ITemplateMapper {
	return &TemplateMapper{}
}

func (this *TemplateMapper) ToModel(entity template_entity.RmgTemplate) template_model.Template {
	return template_model.Template{
		Name:                entity.Name,
		GameMode:            entity.GameMode,
		Description:         entity.Description,
		DisplayWinCondition: entity.DisplayWinCondition,
		MapSize:             entity.SizeX,
		ValueOverrides:      this.GetValueOverrideModelList(entity.ValueOverrides),
		Orientation: helpers.MapPointer(
			entity.Orientation,
			func(entity template_entity.Orientation) template_model.Orientation {
				return template_model.Orientation{Orientation: entity}
			}),
		Border:    helpers.MapPointer(entity.Border, this.ToBorderModel),
		GameRules: this.ToGameRulesModel(entity.GameRules),
		GlobalBans: helpers.MapPointer(
			entity.GlobalBans,
			func(entity template_entity.GlobalBans) template_model.GlobalBans {
				return template_model.GlobalBans{GlobalBans: entity}
			}),
		Variants:           this.GetVariantModelList(entity.Variants),
		ZoneLayouts:        this.GetZoneLayoutDefModelList(entity.ZoneLayouts),
		MandatoryContent:   this.GetMandatoryContentModelList(entity.MandatoryContent),
		ContentCountLimits: this.GetContentCountLimitModelList(entity.ContentCountLimits),
		ContentPools:       this.GetContentPoolModelList(entity.ContentPools),
		ContentLists:       this.GetContentListModelList(entity.ContentLists),
	}
}

func (this *TemplateMapper) ToEntity(model template_model.Template) template_entity.RmgTemplate {
	return template_entity.RmgTemplate{
		Name:                model.Name,
		GameMode:            model.GameMode,
		Description:         model.Description,
		DisplayWinCondition: model.DisplayWinCondition,
		SizeX:               model.MapSize,
		SizeZ:               model.MapSize,
		ValueOverrides:      this.GetValueOverrideEntityList(model.ValueOverrides),
		Orientation: helpers.MapPointer(
			model.Orientation,
			func(model template_model.Orientation) template_entity.Orientation { return model.Orientation }),
		Border:    helpers.MapPointer(model.Border, this.ToBorderEntity),
		GameRules: this.ToGameRulesEntity(model.GameRules),
		GlobalBans: helpers.MapPointer(
			model.GlobalBans,
			func(model template_model.GlobalBans) template_entity.GlobalBans { return model.GlobalBans }),
		Variants:           this.GetVariantEntityList(model.Variants),
		ZoneLayouts:        this.GetZoneLayoutDefEntityList(model.ZoneLayouts),
		MandatoryContent:   this.GetMandatoryContentEntityList(model.MandatoryContent),
		ContentCountLimits: this.GetContentCountLimitEntityList(model.ContentCountLimits),
		ContentPools:       this.GetContentPoolEntityList(model.ContentPools),
		ContentLists:       this.GetContentListEntityList(model.ContentLists),
	}
}

func (this *TemplateMapper) ToMandatoryContentItemModel(
	entity template_entity.MandatoryContentItem) template_model.MandatoryContentItem {
	return template_model.MandatoryContentItem{
		SID:                 entity.SID,
		Name:                entity.Name,
		IsMine:              entity.IsMine,
		IsGuarded:           entity.IsGuarded,
		Rules:               this.GetPlacementRuleModelList(entity.Rules),
		Variant:             entity.Variant,
		Owner:               entity.Owner,
		GuardValue:          entity.GuardValue,
		IncludeLists:        entity.IncludeLists,
		Content:             this.GetWeightedContentModelList(entity.Content),
		DesignatedEncounter: entity.DesignatedEncounter,
		SoloEncounter:       entity.SoloEncounter,
		Road:                entity.Road,
	}
}

func (this *TemplateMapper) ToMandatoryContentItemEntity(
	model template_model.MandatoryContentItem) template_entity.MandatoryContentItem {
	return template_entity.MandatoryContentItem{
		SID:                 model.SID,
		Name:                model.Name,
		IsMine:              model.IsMine,
		IsGuarded:           model.IsGuarded,
		Rules:               this.GetPlacementRuleEntityList(model.Rules),
		Variant:             model.Variant,
		Owner:               model.Owner,
		GuardValue:          model.GuardValue,
		IncludeLists:        model.IncludeLists,
		Content:             this.GetWeightedContentEntityList(model.Content),
		DesignatedEncounter: model.DesignatedEncounter,
		SoloEncounter:       model.SoloEncounter,
		Road:                model.Road,
	}
}

func (this *TemplateMapper) ToContentLimitModel(entity template_entity.ContentLimit) template_model.ContentLimit {
	return template_model.ContentLimit{
		SID:          entity.SID,
		IncludeLists: entity.IncludeLists,
		Content:      this.GetWeightedContentModelList(entity.Content),
		Variant:      entity.Variant,
		MaxCount:     entity.MaxCount,
	}
}

func (this *TemplateMapper) ToContentLimitEntity(model template_model.ContentLimit) template_entity.ContentLimit {
	return template_entity.ContentLimit{
		SID:          model.SID,
		IncludeLists: model.IncludeLists,
		Content:      this.GetWeightedContentEntityList(model.Content),
		Variant:      model.Variant,
		MaxCount:     model.MaxCount,
	}
}

func (this *TemplateMapper) ToGameRulesModel(entity template_entity.GameRules) template_model.GameRules {
	return template_model.GameRules{
		HeroCountMin:       entity.HeroCountMin,
		HeroCountMax:       entity.HeroCountMax,
		HeroCountIncrement: entity.HeroCountIncrement,
		HeroHireBan:        entity.HeroHireBan,
		EncounterHoles:     entity.EncounterHoles,
		TournamentRules:    entity.TournamentRules,
		Bonuses: helpers.MapSlice(entity.Bonuses, func(entity template_entity.Bonus) template_model.Bonus {
			return template_model.Bonus{Bonus: entity}
		}),
		WinConditions:                        template_model.WinConditions{WinConditions: entity.WinConditions},
		GladiatorArena:                       entity.GladiatorArena,
		GladiatorArenaRegistrationStartWork:  entity.GladiatorArenaRegistrationStartWork,
		GladiatorArenaRegistrationStartFight: entity.GladiatorArenaRegistrationStartFight,
		GladiatorArenaDaysDelayStart:         entity.GladiatorArenaDaysDelayStart,
		GladiatorArenaCountDay:               entity.GladiatorArenaCountDay,
		ChampionSelectRule:                   entity.ChampionSelectRule,
		GlobalBans: helpers.MapPointer(entity.GlobalBans,
			func(entity template_entity.GlobalBans) template_model.GlobalBans {
				return template_model.GlobalBans{GlobalBans: entity}
			}),
		FactionLawsExpModifier: entity.FactionLawsExpModifier,
		AstrologyExpModifier:   entity.AstrologyExpModifier,
	}
}

func (this *TemplateMapper) ToGameRulesEntity(model template_model.GameRules) template_entity.GameRules {
	return template_entity.GameRules{
		HeroCountMin:       model.HeroCountMin,
		HeroCountMax:       model.HeroCountMax,
		HeroCountIncrement: model.HeroCountIncrement,
		HeroHireBan:        model.HeroHireBan,
		EncounterHoles:     model.EncounterHoles,
		TournamentRules:    model.TournamentRules,
		Bonuses: helpers.MapSlice(
			model.Bonuses,
			func(model template_model.Bonus) template_entity.Bonus { return model.Bonus }),
		WinConditions:                        model.WinConditions.WinConditions,
		GladiatorArena:                       model.GladiatorArena,
		GladiatorArenaRegistrationStartWork:  model.GladiatorArenaRegistrationStartWork,
		GladiatorArenaRegistrationStartFight: model.GladiatorArenaRegistrationStartFight,
		GladiatorArenaDaysDelayStart:         model.GladiatorArenaDaysDelayStart,
		GladiatorArenaCountDay:               model.GladiatorArenaCountDay,
		ChampionSelectRule:                   model.ChampionSelectRule,
		GlobalBans: helpers.MapPointer(
			model.GlobalBans,
			func(model template_model.GlobalBans) template_entity.GlobalBans { return model.GlobalBans }),
		FactionLawsExpModifier: model.FactionLawsExpModifier,
		AstrologyExpModifier:   model.AstrologyExpModifier,
	}
}

func (this *TemplateMapper) ToBorderModel(entity template_entity.Border) template_model.Border {
	return template_model.Border{
		CornerRadius:   entity.CornerRadius,
		ObstaclesWidth: entity.ObstaclesWidth,
		ObstaclesNoise: this.GetNoiseModelList(entity.ObstaclesNoise),
		WaterWidth:     entity.WaterWidth,
		WaterNoise:     this.GetNoiseModelList(entity.WaterNoise),
		WaterType:      entity.WaterType,
	}
}

func (this *TemplateMapper) ToBorderEntity(model template_model.Border) template_entity.Border {
	return template_entity.Border{
		CornerRadius:   model.CornerRadius,
		ObstaclesWidth: model.ObstaclesWidth,
		ObstaclesNoise: this.GetNoiseEntityList(model.ObstaclesNoise),
		WaterWidth:     model.WaterWidth,
		WaterNoise:     this.GetNoiseEntityList(model.WaterNoise),
		WaterType:      model.WaterType,
	}
}

func (this *TemplateMapper) ToVariantModel(entity template_entity.Variant) template_model.Variant {
	return template_model.Variant{
		Orientation: template_model.Orientation{Orientation: entity.Orientation},
		Border:      this.ToBorderModel(entity.Border),
		Zones:       this.GetZoneModelList(entity.Zones),
		Connections: this.GetConnectionModelList(entity.Connections),
	}
}

func (this *TemplateMapper) ToVariantEntity(model template_model.Variant) template_entity.Variant {
	return template_entity.Variant{
		Orientation: model.Orientation.Orientation,
		Border:      this.ToBorderEntity(model.Border),
		Zones:       this.GetZoneEntityList(model.Zones),
		Connections: this.GetConnectionEntityList(model.Connections),
	}
}

func (this *TemplateMapper) ToConnectionModel(entity template_entity.Connection) template_model.Connection {
	return template_model.Connection{
		Name:                     entity.Name,
		From:                     entity.From,
		To:                       entity.To,
		ConnectionType:           entity.ConnectionType,
		SimTurnSquad:             entity.SimTurnSquad,
		Road:                     entity.Road,
		GuardZone:                entity.GuardZone,
		GuardEscape:              entity.GuardEscape,
		GuardValue:               entity.GuardValue,
		GuardRandomization:       entity.GuardRandomization,
		GuardWeeklyIncrement:     entity.GuardWeeklyIncrement,
		GatePlacement:            entity.GatePlacement,
		Length:                   entity.Length,
		GuardMatchGroup:          entity.GuardMatchGroup,
		PortalPlacementRulesFrom: this.GetPlacementRuleModelList(entity.PortalPlacementRulesFrom),
		PortalPlacementRulesTo:   this.GetPlacementRuleModelList(entity.PortalPlacementRulesTo),
	}
}

func (this *TemplateMapper) ToConnectionEntity(model template_model.Connection) template_entity.Connection {
	return template_entity.Connection{
		Name:                     model.Name,
		From:                     model.From,
		To:                       model.To,
		ConnectionType:           model.ConnectionType,
		SimTurnSquad:             model.SimTurnSquad,
		Road:                     model.Road,
		GuardZone:                model.GuardZone,
		GuardEscape:              model.GuardEscape,
		GuardValue:               model.GuardValue,
		GuardRandomization:       model.GuardRandomization,
		GuardWeeklyIncrement:     model.GuardWeeklyIncrement,
		GatePlacement:            model.GatePlacement,
		Length:                   model.Length,
		GuardMatchGroup:          model.GuardMatchGroup,
		PortalPlacementRulesFrom: this.GetPlacementRuleEntityList(model.PortalPlacementRulesFrom),
		PortalPlacementRulesTo:   this.GetPlacementRuleEntityList(model.PortalPlacementRulesTo),
	}
}

func (this *TemplateMapper) ToZoneModel(entity template_entity.Zone) template_model.Zone {
	return template_model.Zone{
		Name:                      entity.Name,
		Size:                      entity.Size,
		Layout:                    entity.Layout,
		GuardCutoffValue:          entity.GuardCutoffValue,
		GuardRandomization:        entity.GuardRandomization,
		GuardMultiplier:           entity.GuardMultiplier,
		GuardWeeklyIncrement:      entity.GuardWeeklyIncrement,
		GuardReactionDistribution: entity.GuardReactionDistribution,
		DiplomacyModifier:         entity.DiplomacyModifier,
		EncounterHolesSettings: helpers.MapPointer(
			entity.EncounterHolesSettings,
			func(entity template_entity.EncounterHolesSettings) template_model.EncounterHolesSettings {
				return template_model.EncounterHolesSettings{EncounterHolesSettings: entity}
			}),
		RandomHireEnableWeeklyUnitIncrement: entity.RandomHireEnableWeeklyUnitIncrement,
		RandomHireInitialUnitIncrement:      entity.RandomHireInitialUnitIncrement,
		GuardedContentPool:                  entity.GuardedContentPool,
		UnguardedContentPool:                entity.UnguardedContentPool,
		ResourcesContentPool:                entity.ResourcesContentPool,
		MandatoryContent:                    template_model.StringList(entity.MandatoryContent),
		ContentCountLimits:                  template_model.StringList(entity.ContentCountLimits),
		GuardedContentValue:                 entity.GuardedContentValue,
		GuardedContentValuePerArea:          entity.GuardedContentValuePerArea,
		UnguardedContentValue:               entity.UnguardedContentValue,
		UnguardedContentValuePerArea:        entity.UnguardedContentValuePerArea,
		ResourcesValue:                      entity.ResourcesValue,
		ResourcesValuePerArea:               entity.ResourcesValuePerArea,
		MainObjects:                         this.GetMainObjectModelList(entity.MainObjects),
		ZoneBiome:                           template_model.TypedRef{TypedRef: entity.ZoneBiome},
		ContentBiome:                        template_model.TypedRef{TypedRef: entity.ContentBiome},
		MetaObjectsBiome:                    template_model.TypedRef{TypedRef: entity.MetaObjectsBiome},
		CrossroadsPosition:                  entity.CrossroadsPosition,
		Roads:                               this.GetRoadModelList(entity.Roads),
	}
}

func (this *TemplateMapper) ToZoneEntity(model template_model.Zone) template_entity.Zone {
	return template_entity.Zone{
		Name:                      model.Name,
		Size:                      model.Size,
		Layout:                    model.Layout,
		GuardCutoffValue:          model.GuardCutoffValue,
		GuardRandomization:        model.GuardRandomization,
		GuardMultiplier:           model.GuardMultiplier,
		GuardWeeklyIncrement:      model.GuardWeeklyIncrement,
		GuardReactionDistribution: model.GuardReactionDistribution,
		DiplomacyModifier:         model.DiplomacyModifier,
		EncounterHolesSettings: helpers.MapPointer(
			model.EncounterHolesSettings,
			func(model template_model.EncounterHolesSettings) template_entity.EncounterHolesSettings {
				return model.EncounterHolesSettings
			}),
		RandomHireEnableWeeklyUnitIncrement: model.RandomHireEnableWeeklyUnitIncrement,
		RandomHireInitialUnitIncrement:      model.RandomHireInitialUnitIncrement,
		GuardedContentPool:                  model.GuardedContentPool,
		UnguardedContentPool:                model.UnguardedContentPool,
		ResourcesContentPool:                model.ResourcesContentPool,
		MandatoryContent:                    template_entity.StringList(model.MandatoryContent),
		ContentCountLimits:                  template_entity.StringList(model.ContentCountLimits),
		GuardedContentValue:                 model.GuardedContentValue,
		GuardedContentValuePerArea:          model.GuardedContentValuePerArea,
		UnguardedContentValue:               model.UnguardedContentValue,
		UnguardedContentValuePerArea:        model.UnguardedContentValuePerArea,
		ResourcesValue:                      model.ResourcesValue,
		ResourcesValuePerArea:               model.ResourcesValuePerArea,
		MainObjects:                         this.GetMainObjectEntityList(model.MainObjects),
		ZoneBiome:                           model.ZoneBiome.TypedRef,
		ContentBiome:                        model.ContentBiome.TypedRef,
		MetaObjectsBiome:                    model.MetaObjectsBiome.TypedRef,
		CrossroadsPosition:                  model.CrossroadsPosition,
		Roads:                               this.GetRoadEntityList(model.Roads),
	}
}

func (this *TemplateMapper) ToRoadModel(entity template_entity.Road) template_model.Road {
	return template_model.Road{
		Type:                 entity.Type,
		From:                 template_model.TypedRef{TypedRef: entity.From},
		To:                   template_model.TypedRef{TypedRef: entity.To},
		Road:                 entity.Road,
		SimTurnSquad:         entity.SimTurnSquad,
		GuardValue:           entity.GuardValue,
		GuardWeeklyIncrement: entity.GuardWeeklyIncrement,
	}
}

func (this *TemplateMapper) ToRoadEntity(model template_model.Road) template_entity.Road {
	return template_entity.Road{
		Type:                 model.Type,
		From:                 model.From.TypedRef,
		To:                   model.To.TypedRef,
		Road:                 model.Road,
		SimTurnSquad:         model.SimTurnSquad,
		GuardValue:           model.GuardValue,
		GuardWeeklyIncrement: model.GuardWeeklyIncrement,
	}
}

func (this *TemplateMapper) ToMainObjectModel(entity template_entity.MainObject) template_model.MainObject {
	return template_model.MainObject{
		Type:                     entity.Type,
		Spawn:                    entity.Spawn,
		Owner:                    entity.Owner,
		RemoveGuardIfHasOwner:    entity.RemoveGuardIfHasOwner,
		GuardChance:              entity.GuardChance,
		GuardValue:               entity.GuardValue,
		GuardRandomization:       entity.GuardRandomization,
		GuardWeeklyIncrement:     entity.GuardWeeklyIncrement,
		BuildingsConstructionSid: entity.BuildingsConstructionSid,
		Faction: helpers.MapPointer(
			entity.Faction,
			func(entity template_entity.TypedRef) template_model.TypedRef {
				return template_model.TypedRef{TypedRef: entity}
			}),
		Factions:                  entity.Factions,
		Placement:                 entity.Placement,
		PlacementArgs:             entity.PlacementArgs,
		HoldCityWinCon:            entity.HoldCityWinCon,
		IsKeyObject:               entity.IsKeyObject,
		EnableWeeklyUnitIncrement: entity.EnableWeeklyUnitIncrement,
		InitialUnitIncrement:      entity.InitialUnitIncrement,
	}
}

func (this *TemplateMapper) ToMainObjectEntity(model template_model.MainObject) template_entity.MainObject {
	return template_entity.MainObject{
		Type:                     model.Type,
		Spawn:                    model.Spawn,
		Owner:                    model.Owner,
		RemoveGuardIfHasOwner:    model.RemoveGuardIfHasOwner,
		GuardChance:              model.GuardChance,
		GuardValue:               model.GuardValue,
		GuardRandomization:       model.GuardRandomization,
		GuardWeeklyIncrement:     model.GuardWeeklyIncrement,
		BuildingsConstructionSid: model.BuildingsConstructionSid,
		Faction: helpers.MapPointer(
			model.Faction,
			func(model template_model.TypedRef) template_entity.TypedRef { return model.TypedRef }),
		Factions:                  model.Factions,
		Placement:                 model.Placement,
		PlacementArgs:             model.PlacementArgs,
		HoldCityWinCon:            model.HoldCityWinCon,
		IsKeyObject:               model.IsKeyObject,
		EnableWeeklyUnitIncrement: model.EnableWeeklyUnitIncrement,
		InitialUnitIncrement:      model.InitialUnitIncrement,
	}
}

func (this *TemplateMapper) ToZoneLayoutDefModel(entity template_entity.ZoneLayoutDef) template_model.ZoneLayoutDef {
	return template_model.ZoneLayoutDef{
		Name:                  entity.Name,
		ObstaclesFill:         entity.ObstaclesFill,
		ObstaclesFillVoid:     entity.ObstaclesFillVoid,
		LakesFill:             entity.LakesFill,
		MinLakeArea:           entity.MinLakeArea,
		ElevationClusterScale: entity.ElevationClusterScale,
		ElevationModes:        this.GetElevationModeModelList(entity.ElevationModes),
		RoadClusterArea:       entity.RoadClusterArea,
		GuardedEncounterResourceFractions: template_model.GuardedEncounterResourceFractions{
			GuardedEncounterResourceFractions: entity.GuardedEncounterResourceFractions,
		},
		AmbientPickupDistribution: template_model.AmbientPickupDistribution{
			AmbientPickupDistribution: entity.AmbientPickupDistribution,
		},
	}
}

func (this *TemplateMapper) ToZoneLayoutDefEntity(model template_model.ZoneLayoutDef) template_entity.ZoneLayoutDef {
	return template_entity.ZoneLayoutDef{
		Name:                              model.Name,
		ObstaclesFill:                     model.ObstaclesFill,
		ObstaclesFillVoid:                 model.ObstaclesFillVoid,
		LakesFill:                         model.LakesFill,
		MinLakeArea:                       model.MinLakeArea,
		ElevationClusterScale:             model.ElevationClusterScale,
		ElevationModes:                    this.GetElevationModeEntityList(model.ElevationModes),
		RoadClusterArea:                   model.RoadClusterArea,
		GuardedEncounterResourceFractions: model.GuardedEncounterResourceFractions.GuardedEncounterResourceFractions,
		AmbientPickupDistribution:         model.AmbientPickupDistribution.AmbientPickupDistribution,
	}
}

func (this *TemplateMapper) GetElevationModeModelList(
	entities []template_entity.ElevationMode) []template_model.ElevationMode {
	return helpers.MapSlice(entities, func(entity template_entity.ElevationMode) template_model.ElevationMode {
		return template_model.ElevationMode{ElevationMode: entity}
	})
}

func (this *TemplateMapper) GetElevationModeEntityList(
	models []template_model.ElevationMode) []template_entity.ElevationMode {
	return helpers.MapSlice(models, func(model template_model.ElevationMode) template_entity.ElevationMode {
		return model.ElevationMode
	})
}

func (this *TemplateMapper) GetZoneLayoutDefModelList(
	entities []template_entity.ZoneLayoutDef) []template_model.ZoneLayoutDef {
	return helpers.MapSlice(entities, this.ToZoneLayoutDefModel)
}

func (this *TemplateMapper) GetZoneLayoutDefEntityList(
	models []template_model.ZoneLayoutDef) []template_entity.ZoneLayoutDef {
	return helpers.MapSlice(models, this.ToZoneLayoutDefEntity)
}

func (this *TemplateMapper) GetMainObjectModelList(entities []template_entity.MainObject) []template_model.MainObject {
	return helpers.MapSlice(entities, this.ToMainObjectModel)
}

func (this *TemplateMapper) GetMainObjectEntityList(models []template_model.MainObject) []template_entity.MainObject {
	return helpers.MapSlice(models, this.ToMainObjectEntity)
}
func (this *TemplateMapper) GetRoadModelList(entities []template_entity.Road) []template_model.Road {
	return helpers.MapSlice(entities, this.ToRoadModel)
}

func (this *TemplateMapper) GetRoadEntityList(models []template_model.Road) []template_entity.Road {
	return helpers.MapSlice(models, this.ToRoadEntity)
}

func (this *TemplateMapper) GetZoneModelList(entities []template_entity.Zone) []template_model.Zone {
	return helpers.MapSlice(entities, this.ToZoneModel)
}

func (this *TemplateMapper) GetZoneEntityList(models []template_model.Zone) []template_entity.Zone {
	return helpers.MapSlice(models, this.ToZoneEntity)
}

func (this *TemplateMapper) GetPlacementRuleModelList(
	entities []template_entity.PlacementRule) []template_model.PlacementRule {
	return helpers.MapSlice(entities, func(entity template_entity.PlacementRule) template_model.PlacementRule {
		return template_model.PlacementRule{PlacementRule: entity}
	})
}

func (this *TemplateMapper) GetPlacementRuleEntityList(
	models []template_model.PlacementRule) []template_entity.PlacementRule {
	return helpers.MapSlice(models, func(model template_model.PlacementRule) template_entity.PlacementRule {
		return model.PlacementRule
	})
}

func (this *TemplateMapper) GetConnectionModelList(entities []template_entity.Connection) []template_model.Connection {
	return helpers.MapSlice(entities, this.ToConnectionModel)
}

func (this *TemplateMapper) GetConnectionEntityList(models []template_model.Connection) []template_entity.Connection {
	return helpers.MapSlice(models, this.ToConnectionEntity)
}

func (this *TemplateMapper) GetNoiseModelList(entities []template_entity.Noise) []template_model.Noise {
	return helpers.MapSlice(entities, func(entity template_entity.Noise) template_model.Noise {
		return template_model.Noise{Noise: entity}
	})
}

func (this *TemplateMapper) GetNoiseEntityList(models []template_model.Noise) []template_entity.Noise {
	return helpers.MapSlice(models, func(model template_model.Noise) template_entity.Noise { return model.Noise })
}

func (this *TemplateMapper) GetVariantModelList(entities []template_entity.Variant) []template_model.Variant {
	return helpers.MapSlice(entities, this.ToVariantModel)
}

func (this *TemplateMapper) GetVariantEntityList(models []template_model.Variant) []template_entity.Variant {
	return helpers.MapSlice(models, this.ToVariantEntity)
}

func (this *TemplateMapper) GetContentCountLimitModelList(
	entities []template_entity.ContentCountLimit) []template_model.ContentCountLimit {
	return helpers.MapSlice(entities, func(entity template_entity.ContentCountLimit) template_model.ContentCountLimit {
		return template_model.ContentCountLimit{
			Name:   entity.Name,
			Limits: this.GetContentLimitModelList(entity.Limits),
		}
	})
}

func (this *TemplateMapper) GetContentCountLimitEntityList(
	models []template_model.ContentCountLimit) []template_entity.ContentCountLimit {
	return helpers.MapSlice(models, func(model template_model.ContentCountLimit) template_entity.ContentCountLimit {
		return template_entity.ContentCountLimit{
			Name:   model.Name,
			Limits: this.GetContentLimitEntityList(model.Limits),
		}
	})
}

func (this *TemplateMapper) GetContentLimitModelList(
	entities []template_entity.ContentLimit) []template_model.ContentLimit {
	return helpers.MapSlice(entities, this.ToContentLimitModel)
}

func (this *TemplateMapper) GetContentLimitEntityList(
	models []template_model.ContentLimit) []template_entity.ContentLimit {
	return helpers.MapSlice(models, this.ToContentLimitEntity)
}

func (this *TemplateMapper) GetWeightedContentModelList(
	entities []template_entity.WeightedContent) []template_model.WeightedContent {
	return helpers.MapSlice(entities, func(entity template_entity.WeightedContent) template_model.WeightedContent {
		return template_model.WeightedContent{WeightedContent: entity}
	})
}

func (this *TemplateMapper) GetWeightedContentEntityList(
	models []template_model.WeightedContent) []template_entity.WeightedContent {
	return helpers.MapSlice(models, func(model template_model.WeightedContent) template_entity.WeightedContent {
		return model.WeightedContent
	})
}

func (this *TemplateMapper) GetMandatoryContentItemModelList(
	entities []template_entity.MandatoryContentItem) []template_model.MandatoryContentItem {
	return helpers.MapSlice(entities, this.ToMandatoryContentItemModel)
}

func (this *TemplateMapper) GetMandatoryContentItemEntityList(
	models []template_model.MandatoryContentItem) []template_entity.MandatoryContentItem {
	return helpers.MapSlice(models, this.ToMandatoryContentItemEntity)
}

func (this *TemplateMapper) GetMandatoryContentModelList(
	entities []template_entity.MandatoryContent) []template_model.MandatoryContent {
	return helpers.MapSlice(entities, func(entity template_entity.MandatoryContent) template_model.MandatoryContent {
		return template_model.MandatoryContent{
			Name:    entity.Name,
			Content: this.GetMandatoryContentItemModelList(entity.Content),
		}
	})
}

func (this *TemplateMapper) GetMandatoryContentEntityList(
	models []template_model.MandatoryContent) []template_entity.MandatoryContent {
	return helpers.MapSlice(models, func(model template_model.MandatoryContent) template_entity.MandatoryContent {
		return template_entity.MandatoryContent{
			Name:    model.Name,
			Content: this.GetMandatoryContentItemEntityList(model.Content),
		}
	})
}

func (this *TemplateMapper) GetValueOverrideModelList(
	entities []template_entity.ValueOverride) []template_model.ValueOverride {
	return helpers.MapSlice(entities, func(entity template_entity.ValueOverride) template_model.ValueOverride {
		return template_model.ValueOverride{ValueOverride: entity}
	})
}

func (this *TemplateMapper) GetValueOverrideEntityList(
	models []template_model.ValueOverride) []template_entity.ValueOverride {
	return helpers.MapSlice(models, func(model template_model.ValueOverride) template_entity.ValueOverride {
		return model.ValueOverride
	})
}

func (this *TemplateMapper) GetContentListModelList(
	entities []template_entity.ContentList) []template_model.ContentList {
	return helpers.MapSlice(entities, func(entity template_entity.ContentList) template_model.ContentList {
		return template_model.ContentList(entity)
	})
}

func (this *TemplateMapper) GetContentListEntityList(
	models []template_model.ContentList) []template_entity.ContentList {
	return helpers.MapSlice(models, func(model template_model.ContentList) template_entity.ContentList {
		return template_entity.ContentList(model)
	})
}

func (this *TemplateMapper) GetContentPoolModelList(
	entities []template_entity.ContentPool) []template_model.ContentPool {
	return helpers.MapSlice(entities, func(entity template_entity.ContentPool) template_model.ContentPool {
		return template_model.ContentPool(entity)
	})
}

func (this *TemplateMapper) GetContentPoolEntityList(
	models []template_model.ContentPool) []template_entity.ContentPool {
	return helpers.MapSlice(models, func(model template_model.ContentPool) template_entity.ContentPool {
		return template_entity.ContentPool(model)
	})
}
