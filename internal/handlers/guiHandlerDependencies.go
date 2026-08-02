package handlers

type GUIHandlerDependencies struct {
	TemplateWorkflow    ITemplateWorkflowOperations
	StatePersistence    IStatePersistenceOperations
	TemplatePersistence ITemplatePersistenceOperations
	Preview             IPreviewOperations
	ContentRule         IContentRuleOperations
	ZoneEditor          IZoneEditorOperations
}
