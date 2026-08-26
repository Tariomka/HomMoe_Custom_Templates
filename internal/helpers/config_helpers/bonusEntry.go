package config_helpers

import (
	"crypto/sha256"
	"encoding/json/v2"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

func GetHash(entry editor_state_model.BonusEntry) string {
	data, err := json.Marshal(entry)
	if err != nil {
		return GetString(entry.PresetType) + "|" + entry.ReceiverFilter + "|" + entry.Param + "|" + entry.Param2
	}

	hash := sha256.Sum256(data)
	return string(hash[:])
}
