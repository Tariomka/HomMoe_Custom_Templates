package services_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

func TestZoneContentManagerBuildsPlayerContent(t *testing.T) {
	manager := &services.ZoneContentManager{}
	content := manager.BuildPlayerZoneMandatoryContent(1.0)

	if len(content) == 0 {
		t.Errorf("Expected player zone content, got none")
	}
}

func TestContentItemBuilderCreatesItem(t *testing.T) {
	item := services.NewContentItem("test_sid").
		WithCount(5).
		Guarded().
		Build()

	if item.SID != "test_sid" {
		t.Errorf("Expected SID test_sid, got %s", item.SID)
	}

	if item.Count != 5 {
		t.Errorf("Expected count 5, got %d", item.Count)
	}

	if !item.Guarded {
		t.Errorf("Expected guarded=true, got false")
	}
}
