package interfaces

type IBackend interface {
	ITemplateWorkflowHandler
	IStatePersistenceHandler
	IStateValidationHandler
	IPreviewHandler
	IContentRuleHandler
	IZoneEditorHandler
}
