package app

import (
	"strings"
	"testing"
)

func TestCuratedIconsIncludeTaskCategories(t *testing.T) {
	t.Parallel()

	expectedLabels := []string{
		"Todo / Inbox",
		"Feature / Enhancement",
		"Bug / Defect",
		"Fix / Patch",
		"Refactor / Cleanup",
		"Chore / Maintenance",
		"Docs / Writing",
		"Tests / Verification",
		"Review / Inspect",
		"Blocked / Attention",
		"Urgent / Priority",
		"Release / Launch",
	}

	for _, expected := range expectedLabels {
		if !curatedIconLabelExists(expected) {
			t.Fatalf("expected curated icon label containing %q", expected)
		}
	}
}

func TestCuratedIconsKeepDefaultFolderResetOption(t *testing.T) {
	t.Parallel()

	if len(curatedIcons) == 0 {
		t.Fatal("expected curated icon list to be non-empty")
	}

	first := curatedIcons[0]
	if first.ID != "" {
		t.Fatalf("expected first icon id to reset to default folder, got %q", first.ID)
	}
	if !strings.Contains(first.Label, "Default Folder") {
		t.Fatalf("expected first icon label to mention default folder, got %q", first.Label)
	}
}

func curatedIconLabelExists(expected string) bool {
	for _, item := range curatedIcons {
		if strings.Contains(item.Label, expected) {
			return true
		}
	}
	return false
}
