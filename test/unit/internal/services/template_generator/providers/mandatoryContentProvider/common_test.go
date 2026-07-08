package mandatoryContentProvider_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// groupContent returns the content items of the mandatory-content group with
// the given name, or nil when no such group exists.
func groupContent(groups []entities.MandatoryContent, name string) []entities.MandatoryContentItem {
	for _, group := range groups {
		if group.Name == name {
			return group.Content
		}
	}
	return nil
}

// itemSids returns the SIDs of the given content items in order.
func itemSids(items []entities.MandatoryContentItem) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, item.SID)
	}
	return out
}

// countGroupsNamed counts how many groups carry the given name.
func countGroupsNamed(groups []entities.MandatoryContent, name string) int {
	count := 0
	for _, group := range groups {
		if group.Name == name {
			count++
		}
	}
	return count
}

// groupNames returns the names of the given groups in order.
func groupNames(groups []entities.MandatoryContent) []string {
	var names []string
	for _, group := range groups {
		names = append(names, group.Name)
	}
	return names
}
