package handlers

type GUIHandlerDependencies struct {
	TemplateWorkflow    TemplateWorkflowOperations
	StatePersistence    StatePersistenceOperations
	TemplatePersistence TemplatePersistenceOperations
	Preview             PreviewOperations
	ContentRule         ContentRuleOperations
	ZoneEditor          ZoneEditorOperations
}
