package config_inner

import (
	"crypto/sha256"
	"encoding/json"
)

// BonusEntry is the editor-side view-model for a single configurable game-start bonus.
type BonusEntry struct {
	PresetType     BonusPresetType `json:"presetType"`
	ReceiverFilter string          `json:"receiverFilter"`   // "start_hero" or "all_heroes".
	Param          string          `json:"param,omitempty"`  // Spell sid / item sid / numeric value depending on type.
	Param2         string          `json:"param2,omitempty"` // For Spell: "1" = free, "0" = normal. Unused for other types.
}

func (this BonusEntry) GetHash() string {
	data, err := json.Marshal(this)
	if err != nil {
		return this.PresetType.String() + "|" + this.ReceiverFilter + "|" + this.Param + "|" + this.Param2
	}

	hash := sha256.Sum256(data)
	return string(hash[:])
}
