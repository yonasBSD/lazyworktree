package screen

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/chmouel/lazyworktree/internal/models"
	"github.com/chmouel/lazyworktree/internal/theme"
)

func TestIssueSelectionScreenFilterToggle(t *testing.T) {
	issues := []*models.IssueInfo{
		{Number: 1, Title: "First"},
		{Number: 2, Title: "Second"},
	}

	scr := NewIssueSelectionScreen(issues, 80, 30, theme.Dracula(), true)
	if scr.FilterActive {
		t.Fatal("expected filter to be inactive by default")
	}

	next, _ := scr.Update(tea.KeyPressMsg{Code: 'f', Text: string('f')})
	nextScr, ok := next.(*IssueSelectionScreen)
	if !ok || nextScr == nil {
		t.Fatal("expected Update to return issue selection screen after f")
	}
	scr = nextScr
	if !scr.FilterActive {
		t.Fatal("expected filter to be active after f")
	}

	next, _ = scr.Update(tea.KeyPressMsg{Code: '2', Text: string('2')})
	nextScr, ok = next.(*IssueSelectionScreen)
	if !ok || nextScr == nil {
		t.Fatal("expected Update to return issue selection screen after typing")
	}
	scr = nextScr
	if len(scr.Filtered) != 1 || scr.Filtered[0].Number != 2 {
		t.Fatalf("expected filtered results to include only #2, got %v", scr.Filtered)
	}

	next, _ = scr.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	nextScr, ok = next.(*IssueSelectionScreen)
	if !ok || nextScr == nil {
		t.Fatal("expected Update to return issue selection screen after Esc")
	}
	scr = nextScr
	if scr.FilterActive {
		t.Fatal("expected filter to be inactive after Esc")
	}
	if len(scr.Filtered) != 1 || scr.Filtered[0].Number != 2 {
		t.Fatalf("expected filter to remain applied after Esc, got %v", scr.Filtered)
	}
}

func TestIssueSelectionScreenRanksNumberAndTitleMatches(t *testing.T) {
	issues := []*models.IssueInfo{
		{Number: 45, Title: "Open browser page"},
		{Number: 451, Title: "Browse worktree files"},
	}

	scr := NewIssueSelectionScreen(issues, 80, 30, theme.Dracula(), true)

	scr.FilterInput.SetValue("45")
	scr.applyFilter()
	if len(scr.Filtered) != 2 {
		t.Fatalf("expected two issue matches, got %d", len(scr.Filtered))
	}
	if scr.Filtered[0].Number != 45 {
		t.Fatalf("expected exact number match first, got #%d", scr.Filtered[0].Number)
	}

	scr.FilterInput.SetValue("browse")
	scr.applyFilter()
	if scr.Filtered[0].Number != 451 {
		t.Fatalf("expected stronger title match first, got #%d", scr.Filtered[0].Number)
	}
	if scr.Cursor != 0 {
		t.Fatalf("expected cursor to reset to first ranked issue, got %d", scr.Cursor)
	}
}
