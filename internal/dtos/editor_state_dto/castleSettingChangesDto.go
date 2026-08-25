package editor_state_dto

// CastleSettingChangesDto reports which castle options moved between two editor
// states. It mirrors the model type field for field so the mapper can convert
// the two structs directly.
type CastleSettingChangesDto struct {
	PlayerCastles bool
	NeutralSimple bool
	NeutralLowest bool
	NeutralLow    bool
	NeutralMedium bool
	NeutralHigh   bool
	Hub           bool
}
