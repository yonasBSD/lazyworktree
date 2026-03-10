package screen

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/chmouel/lazyworktree/internal/models"
	"github.com/chmouel/lazyworktree/internal/theme"
)

func TestNewPRSelectionScreen(t *testing.T) {
	prs := []*models.PRInfo{
		{Number: 1, Title: "First PR", Author: "user1"},
		{Number: 2, Title: "Second PR", Author: "user2"},
	}
	scr := NewPRSelectionScreen(prs, 80, 30, theme.Dracula(), true)

	if scr.Type() != TypePRSelect {
		t.Errorf("expected Type to be TypePRSelect, got %v", scr.Type())
	}

	if len(scr.FilteredPRs()) != 2 {
		t.Errorf("expected 2 filtered PRs, got %d", len(scr.FilteredPRs()))
	}

	if scr.Cursor != 0 {
		t.Errorf("expected cursor to start at 0, got %d", scr.Cursor)
	}
}

func TestPRSelectionScreenNavigation(t *testing.T) {
	prs := []*models.PRInfo{
		{Number: 1, Title: "First"},
		{Number: 2, Title: "Second"},
		{Number: 3, Title: "Third"},
	}
	scr := NewPRSelectionScreen(prs, 80, 30, theme.Dracula(), true)

	scr.Cursor = 1
	if scr.Cursor != 1 {
		t.Errorf("expected cursor to be 1, got %d", scr.Cursor)
	}

	pr, ok := scr.SelectedPR()
	if !ok || pr.Number != 2 {
		t.Error("expected to select second PR")
	}
}

func TestPRSelectionScreenFiltering(t *testing.T) {
	prs := []*models.PRInfo{
		{Number: 123, Title: "Add feature X"},
		{Number: 456, Title: "Fix bug Y"},
		{Number: 789, Title: "Update feature Z"},
	}
	scr := NewPRSelectionScreen(prs, 80, 30, theme.Dracula(), true)

	// Filter by title
	scr.FilterInput.SetValue("feature")
	scr.applyFilter()

	filtered := scr.FilteredPRs()
	if len(filtered) != 2 {
		t.Errorf("expected 2 filtered PRs matching 'feature', got %d", len(filtered))
	}

	// Filter by number
	scr.FilterInput.SetValue("456")
	scr.applyFilter()

	filtered = scr.FilteredPRs()
	if len(filtered) != 1 {
		t.Errorf("expected 1 filtered PR matching '456', got %d", len(filtered))
	}

	if filtered[0].Number != 456 {
		t.Errorf("expected filtered PR to have number 456, got %d", filtered[0].Number)
	}

	// Clear filter
	scr.FilterInput.SetValue("")
	scr.applyFilter()

	filtered = scr.FilteredPRs()
	if len(filtered) != 3 {
		t.Errorf("expected all 3 PRs after clearing filter, got %d", len(filtered))
	}
}

func TestPRSelectionScreenRanksNumberAndTitleMatches(t *testing.T) {
	prs := []*models.PRInfo{
		{Number: 45, Title: "Open browser page"},
		{Number: 451, Title: "Browse worktree files"},
	}
	scr := NewPRSelectionScreen(prs, 80, 30, theme.Dracula(), true)

	scr.FilterInput.SetValue("45")
	scr.applyFilter()
	filtered := scr.FilteredPRs()
	if len(filtered) != 2 {
		t.Fatalf("expected two PR matches, got %d", len(filtered))
	}
	if filtered[0].Number != 45 {
		t.Fatalf("expected exact number match first, got #%d", filtered[0].Number)
	}

	scr.FilterInput.SetValue("browse")
	scr.applyFilter()
	filtered = scr.FilteredPRs()
	if filtered[0].Number != 451 {
		t.Fatalf("expected stronger title match first, got #%d", filtered[0].Number)
	}
	if scr.Cursor != 0 {
		t.Fatalf("expected cursor to reset to first ranked PR, got %d", scr.Cursor)
	}
}

func TestPRSelectionScreenPrefersPRNumberOverTitleConflict(t *testing.T) {
	prs := []*models.PRInfo{
		{Number: 45, Title: "General cleanup"},
		{Number: 99, Title: "45 cleanup"},
	}
	scr := NewPRSelectionScreen(prs, 80, 30, theme.Dracula(), true)

	scr.FilterInput.SetValue("45")
	scr.applyFilter()

	filtered := scr.FilteredPRs()
	if len(filtered) != 2 {
		t.Fatalf("expected two PR matches, got %d", len(filtered))
	}
	if filtered[0].Number != 45 {
		t.Fatalf("expected PR number match to rank before title match, got #%d", filtered[0].Number)
	}
}

func TestPRSelectionScreenFilterToggle(t *testing.T) {
	prs := []*models.PRInfo{
		{Number: 1, Title: "First"},
		{Number: 2, Title: "Second"},
	}
	scr := NewPRSelectionScreen(prs, 80, 30, theme.Dracula(), true)
	if scr.FilterActive {
		t.Fatal("expected filter to be inactive by default")
	}

	next, _ := scr.Update(tea.KeyPressMsg{Code: 'f', Text: string('f')})
	nextScr, ok := next.(*PRSelectionScreen)
	if !ok || nextScr == nil {
		t.Fatal("expected Update to return PR selection screen after f")
	}
	scr = nextScr
	if !scr.FilterActive {
		t.Fatal("expected filter to be active after f")
	}

	next, _ = scr.Update(tea.KeyPressMsg{Code: '2', Text: string('2')})
	nextScr, ok = next.(*PRSelectionScreen)
	if !ok || nextScr == nil {
		t.Fatal("expected Update to return PR selection screen after typing")
	}
	scr = nextScr
	filtered := scr.FilteredPRs()
	if len(filtered) != 1 || filtered[0].Number != 2 {
		t.Fatalf("expected filtered results to include only #2, got %v", filtered)
	}

	next, _ = scr.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	nextScr, ok = next.(*PRSelectionScreen)
	if !ok || nextScr == nil {
		t.Fatal("expected Update to return PR selection screen after Esc")
	}
	scr = nextScr
	if scr.FilterActive {
		t.Fatal("expected filter to be inactive after Esc")
	}
	filtered = scr.FilteredPRs()
	if len(filtered) != 1 || filtered[0].Number != 2 {
		t.Fatalf("expected filter to remain applied after Esc, got %v", filtered)
	}
}

func TestPRSelectionScreenSelection(t *testing.T) {
	prs := []*models.PRInfo{
		{Number: 1, Title: "First", Branch: "branch1"},
		{Number: 2, Title: "Second", Branch: "branch2"},
	}
	scr := NewPRSelectionScreen(prs, 80, 30, theme.Dracula(), true)

	pr, ok := scr.SelectedPR()
	if !ok {
		t.Fatal("expected SelectedPR to return true")
	}
	if pr.Number != 1 {
		t.Errorf("expected selected PR to have number 1, got %d", pr.Number)
	}

	scr.Cursor = 1
	pr, ok = scr.SelectedPR()
	if !ok {
		t.Fatal("expected SelectedPR to return true")
	}
	if pr.Number != 2 {
		t.Errorf("expected selected PR to have number 2, got %d", pr.Number)
	}

	scr.Cursor = 99
	_, ok = scr.SelectedPR()
	if ok {
		t.Error("expected SelectedPR to return false for out of bounds cursor")
	}
}

func TestPRSelectionScreenCallbacks(t *testing.T) {
	prs := []*models.PRInfo{
		{Number: 1, Title: "First", Branch: "branch1"},
	}
	scr := NewPRSelectionScreen(prs, 80, 30, theme.Dracula(), true)

	selectCalled := false
	var selectedPR *models.PRInfo
	scr.OnSelectPR = func(pr *models.PRInfo) tea.Cmd {
		selectCalled = true
		selectedPR = pr
		return nil
	}

	result, _ := scr.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if result != nil {
		t.Error("expected screen to close (return nil) on Enter")
	}
	if !selectCalled {
		t.Error("expected OnSelectPR callback to be called")
	}
	if selectedPR == nil || selectedPR.Number != 1 {
		t.Error("expected selectedPR to be the first PR")
	}

	// Test OnCancel callback
	scr = NewPRSelectionScreen(prs, 80, 30, theme.Dracula(), true)
	cancelCalled := false
	scr.OnCancel = func() tea.Cmd {
		cancelCalled = true
		return nil
	}

	result, _ = scr.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if result != nil {
		t.Error("expected screen to close (return nil) on Esc")
	}
	if !cancelCalled {
		t.Error("expected OnCancel callback to be called")
	}
}

func TestPRSelectionScreenView(t *testing.T) {
	prs := []*models.PRInfo{
		{Number: 1, Title: "Test PR", Author: "testuser", CIStatus: "success"},
	}
	scr := NewPRSelectionScreen(prs, 80, 30, theme.Dracula(), true)

	view := scr.View()
	if view == "" {
		t.Error("expected View to return non-empty string")
	}

	if !strings.Contains(view, "Test PR") {
		t.Error("expected view to contain PR title")
	}
}

func TestPRSelectionScreenCIIconsUseProvider(t *testing.T) {
	previousProvider := currentIconProvider
	SetIconProvider(&testIconProvider{ciIcon: "CI!"})
	defer SetIconProvider(previousProvider)

	prs := []*models.PRInfo{
		{Number: 1, Title: "Test PR", Author: "testuser", CIStatus: "success"},
	}
	scr := NewPRSelectionScreen(prs, 80, 30, theme.Dracula(), true)

	view := scr.View()
	if !strings.Contains(view, "CI!") {
		t.Fatalf("expected view to include CI icon from provider, got %q", view)
	}
}

func TestPRSelectionScreenCIStatusColoring(t *testing.T) {
	prs := []*models.PRInfo{
		{Number: 1, Title: "Success PR", CIStatus: "success"},
		{Number: 2, Title: "Failure PR", CIStatus: "failure"},
		{Number: 3, Title: "Pending PR", CIStatus: "pending"},
		{Number: 4, Title: "Draft PR", IsDraft: true},
	}
	scr := NewPRSelectionScreen(prs, 100, 30, theme.Dracula(), true)

	view := scr.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
}

type testIconProvider struct {
	ciIcon string
}

func (p *testIconProvider) GetPRIcon() string {
	return "PR"
}

func (p *testIconProvider) GetIssueIcon() string {
	return "ISS"
}

func (p *testIconProvider) GetCIIcon(conclusion string) string {
	return p.ciIcon
}

func (p *testIconProvider) GetUIIcon(icon UIIcon) string {
	return ""
}

func TestPRSelectionScreenEmptyList(t *testing.T) {
	scr := NewPRSelectionScreen([]*models.PRInfo{}, 80, 30, theme.Dracula(), true)

	view := scr.View()
	if !strings.Contains(view, "No open PRs") {
		t.Error("expected view to show 'No open PRs' message")
	}

	_, ok := scr.SelectedPR()
	if ok {
		t.Error("expected SelectedPR to return false for empty list")
	}
}

func TestPRSelectionScreenNoMatchingFilter(t *testing.T) {
	prs := []*models.PRInfo{
		{Number: 1, Title: "First"},
	}
	scr := NewPRSelectionScreen(prs, 80, 30, theme.Dracula(), true)

	scr.FilterInput.SetValue("nonexistent")
	scr.applyFilter()

	view := scr.View()
	if !strings.Contains(view, "No PRs match your filter") {
		t.Error("expected view to show 'No PRs match' message")
	}
}

func TestPRSelectionScreenAttachedBranches(t *testing.T) {
	prs := []*models.PRInfo{
		{Number: 1, Title: "Available PR", Branch: "available-branch"},
		{Number: 2, Title: "Attached PR", Branch: "attached-branch"},
	}
	scr := NewPRSelectionScreen(prs, 100, 30, theme.Dracula(), true)

	scr.AttachedBranches = map[string]string{
		"attached-branch": "my-worktree",
	}

	wtName, attached := scr.isAttached(prs[0])
	if attached {
		t.Error("expected first PR to not be attached")
	}
	if wtName != "" {
		t.Errorf("expected empty worktree name for non-attached PR, got %q", wtName)
	}

	wtName, attached = scr.isAttached(prs[1])
	if !attached {
		t.Error("expected second PR to be attached")
	}
	if wtName != "my-worktree" {
		t.Errorf("expected worktree name 'my-worktree', got %q", wtName)
	}

	view := scr.View()
	if !strings.Contains(view, "(in: my-worktree)") {
		t.Error("expected view to show worktree info for attached PR")
	}
}

func TestPRSelectionScreenAttachedPRSelectable(t *testing.T) {
	prs := []*models.PRInfo{
		{Number: 1, Title: "Attached PR", Branch: "attached-branch"},
	}
	scr := NewPRSelectionScreen(prs, 100, 30, theme.Dracula(), true)

	scr.AttachedBranches = map[string]string{
		"attached-branch": "my-worktree",
	}

	selectCalled := false
	scr.OnSelectPR = func(pr *models.PRInfo) tea.Cmd {
		selectCalled = true
		return nil
	}

	result, _ := scr.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if result != nil {
		t.Error("expected screen to close when selecting attached PR")
	}
	if !selectCalled {
		t.Error("expected OnSelectPR callback to be called for attached PR")
	}
}

func TestPRSelectionScreenStatusMessageClearedOnNavigation(t *testing.T) {
	prs := []*models.PRInfo{
		{Number: 1, Title: "First PR", Branch: "branch1"},
		{Number: 2, Title: "Second PR", Branch: "branch2"},
	}
	scr := NewPRSelectionScreen(prs, 100, 30, theme.Dracula(), true)

	scr.StatusMessage = "Some error message"

	scr.Update(tea.KeyPressMsg{Code: tea.KeyDown})

	if scr.StatusMessage != "" {
		t.Errorf("expected status message to be cleared on navigation, got %q", scr.StatusMessage)
	}

	scr.StatusMessage = "Another message"

	scr.Update(tea.KeyPressMsg{Code: tea.KeyUp})

	if scr.StatusMessage != "" {
		t.Errorf("expected status message to be cleared on navigation, got %q", scr.StatusMessage)
	}
}
