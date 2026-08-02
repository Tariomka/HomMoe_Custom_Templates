package handlers

type GUIHandlerDependencies struct {
	TemplateWorkflow    ITemplateWorkflow
	StatePersistence    IStatePersistence
	TemplatePersistence ITemplatePersistence
	Preview             IPreview
	ContentRule         IContentRule
	ZoneEditor          IZoneEditor
}
