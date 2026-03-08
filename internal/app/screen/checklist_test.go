package screen

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/chmouel/lazyworktree/internal/theme"
)

func TestChecklistScreenFilterToggle(t *testing.T) {
	items := []ChecklistItem{
		{ID: "one", Label: "One"},
		{ID: "two", Label: "Two"},
	}

	scr := NewChecklistScreen(items, "Test", "Filter...", "No items", 80, 30, theme.Dracula())
	if scr.FilterActive {
		t.Fatal("expected filter to be inactive by default")
	}

	next, _ := scr.Update(tea.KeyPressMsg{Code: 'f', Text: string('f')})
	nextScr, ok := next.(*ChecklistScreen)
	if !ok || nextScr == nil {
		t.Fatal("expected Update to return checklist screen after f")
	}
	scr = nextScr
	if !scr.FilterActive {
		t.Fatal("expected filter to be active after f")
	}

	next, _ = scr.Update(tea.KeyPressMsg{Code: 't', Text: string('t')})
	nextScr, ok = next.(*ChecklistScreen)
	if !ok || nextScr == nil {
		t.Fatal("expected Update to return checklist screen after typing")
	}
	scr = nextScr
	if len(scr.Filtered) != 1 || scr.Filtered[0].ID != "two" {
		t.Fatalf("expected filtered results to include only 'two', got %v", scr.Filtered)
	}

	next, _ = scr.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	nextScr, ok = next.(*ChecklistScreen)
	if !ok || nextScr == nil {
		t.Fatal("expected Update to return checklist screen after Esc")
	}
	scr = nextScr
	if scr.FilterActive {
		t.Fatal("expected filter to be inactive after Esc")
	}
	if len(scr.Filtered) != 1 || scr.Filtered[0].ID != "two" {
		t.Fatalf("expected filter to remain applied after Esc, got %v", scr.Filtered)
	}
}

func TestChecklistScreenRanksBetterMatchesFirst(t *testing.T) {
	items := []ChecklistItem{
		{ID: "browser", Label: "Open PR", Description: "Open PR in browser"},
		{ID: "browse", Label: "Browse files", Description: "Inspect files"},
	}

	scr := NewChecklistScreen(items, "Test", "Filter...", "No items", 80, 30, theme.Dracula())
	scr.FilterInput.SetValue("browse")
	scr.applyFilter()

	if len(scr.Filtered) != 2 {
		t.Fatalf("expected both items to match, got %d", len(scr.Filtered))
	}
	if scr.Filtered[0].ID != "browse" {
		t.Fatalf("expected label match first, got %q", scr.Filtered[0].ID)
	}
	if scr.Cursor != 0 {
		t.Fatalf("expected cursor to reset to first ranked item, got %d", scr.Cursor)
	}
}
