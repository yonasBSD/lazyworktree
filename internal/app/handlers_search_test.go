package app

import (
	"path/filepath"
	"testing"

	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/chmouel/lazyworktree/internal/config"
	"github.com/chmouel/lazyworktree/internal/models"
)

func TestSearchWorktreeSelectsMatch(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
		SortMode:    "path",
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 0

	wt1Path := filepath.Join(cfg.WorktreeDir, "alpha")
	wt2Path := filepath.Join(cfg.WorktreeDir, "beta")
	m.state.data.worktrees = []*models.WorktreeInfo{
		{Path: wt1Path, Branch: "feat-one"},
		{Path: wt2Path, Branch: "feat-two"},
	}
	m.updateTable()
	m.state.ui.worktreeTable.SetCursor(0)
	m.state.data.selectedIndex = 0

	updated, _ := m.handleKeyMsg(tea.KeyPressMsg{Code: '/', Text: string('/')})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	_, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: 'b', Text: string('b')})

	if m.state.ui.worktreeTable.Cursor() != 1 {
		t.Fatalf("expected cursor to move to match, got %d", m.state.ui.worktreeTable.Cursor())
	}
}

func TestSearchLogSelectsNextMatch(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 3
	m.state.data.logEntries = []commitLogEntry{
		{sha: "abc123", authorInitials: "ab", message: "Fix bug in parser"},
		{sha: "def456", authorInitials: "de", message: "Add new feature"},
		{sha: "ghi789", authorInitials: "gh", message: "Fix tests"},
	}
	m.state.ui.logTable.SetRows([]table.Row{
		{"abc123", "ab", formatCommitMessage("Fix bug in parser")},
		{"def456", "de", formatCommitMessage("Add new feature")},
		{"ghi789", "gh", formatCommitMessage("Fix tests")},
	})

	updated, _ := m.handleKeyMsg(tea.KeyPressMsg{Code: '/', Text: string('/')})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	_, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: 'f', Text: string('f')})
	_, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: 'i', Text: string('i')})
	_, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: 'x', Text: string('x')})

	if m.state.ui.logTable.Cursor() != 0 {
		t.Fatalf("expected first match at cursor 0, got %d", m.state.ui.logTable.Cursor())
	}

	// Confirm search with Enter, then use n to advance to next match
	updated, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: tea.KeyEnter})
	updatedModel, ok = updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	_, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: 'n', Text: string('n')})
	if m.state.ui.logTable.Cursor() != 2 {
		t.Fatalf("expected next match at cursor 2, got %d", m.state.ui.logTable.Cursor())
	}
}

func TestSearchStatusSelectsMatch(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 2
	m.state.ui.statusViewport = viewport.New(viewport.WithWidth(40), viewport.WithHeight(10))
	// Note: tree sorts alphabetically, so README.md (R) comes before app.go (a)
	m.setStatusFiles([]StatusFile{
		{Filename: "app.go", Status: ".M"},
		{Filename: "README.md", Status: ".M"},
	})
	m.rebuildStatusContentWithHighlight()

	updated, _ := m.handleKeyMsg(tea.KeyPressMsg{Code: '/', Text: string('/')})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	// Search for "app" to find app.go which is at index 1 after sorting
	_, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: 'a', Text: string('a')})
	_, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: 'p', Text: string('p')})
	_, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: 'p', Text: string('p')})

	if m.state.services.statusTree.Index != 1 {
		t.Fatalf("expected statusTreeIndex 1, got %d", m.state.services.statusTree.Index)
	}
}

// TestRenderStatusFilesHighlighting tests that selected file is highlighted.

func TestClearSearchQuery(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.view.SearchTarget = searchTargetWorktrees
	m.setSearchQuery(searchTargetWorktrees, "test query")
	m.state.ui.filterInput.SetValue("test")

	m.clearSearchQuery()

	if m.state.services.filter.WorktreeSearchQuery != "" {
		t.Errorf("expected worktreeSearchQuery to be empty, got %q", m.state.services.filter.WorktreeSearchQuery)
	}
	if m.state.ui.filterInput.Value() != "" {
		t.Errorf("expected filterInput to be empty, got %q", m.state.ui.filterInput.Value())
	}
}

func TestRestoreFocusAfterSearch(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")

	t.Run("restore focus to worktrees", func(t *testing.T) {
		m.state.view.SearchTarget = searchTargetWorktrees
		m.restoreFocusAfterSearch()
		if !m.state.ui.worktreeTable.Focused() {
			t.Error("expected worktreeTable to be focused")
		}
	})

	t.Run("restore focus to log", func(t *testing.T) {
		m.state.view.SearchTarget = searchTargetLog
		m.restoreFocusAfterSearch()
		if !m.state.ui.logTable.Focused() {
			t.Error("expected logTable to be focused")
		}
	})
}

func TestCICheckNavigationDown(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 1
	m.state.data.selectedIndex = 0
	m.state.data.filteredWts = []*models.WorktreeInfo{
		{Path: testWorktreePath, Branch: "feat"},
	}

	// Add CI checks to cache
	checks := []*models.CICheck{
		{Name: "build", Conclusion: "success", Link: "https://github.com/owner/repo/actions/runs/123"},
		{Name: "test", Conclusion: "failure", Link: "https://github.com/owner/repo/actions/runs/456"},
		{Name: "lint", Conclusion: "success", Link: "https://github.com/owner/repo/actions/runs/789"},
	}
	m.cache.ciCache.Set("feat", checks)

	// Start navigating CI checks
	m.ciCheckIndex = 0

	// Navigate down with n
	_, _ = m.navigateCICheckDown()
	if m.ciCheckIndex != 1 {
		t.Fatalf("expected ciCheckIndex 1 after n, got %d", m.ciCheckIndex)
	}

	// Navigate down again
	_, _ = m.navigateCICheckDown()
	if m.ciCheckIndex != 2 {
		t.Fatalf("expected ciCheckIndex 2 after second n, got %d", m.ciCheckIndex)
	}

	// At last CI check, should stay at last check (no wrapping to file tree since it's a separate pane)
	_, _ = m.navigateCICheckDown()
	if m.ciCheckIndex != 2 {
		t.Fatalf("expected ciCheckIndex to stay at 2, got %d", m.ciCheckIndex)
	}
}

func TestCICheckNavigationUp(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 1
	m.state.data.selectedIndex = 0
	m.state.data.filteredWts = []*models.WorktreeInfo{
		{Path: testWorktreePath, Branch: "feat"},
	}

	// Add CI checks to cache
	checks := []*models.CICheck{
		{Name: "build", Conclusion: "success", Link: "https://github.com/owner/repo/actions/runs/123"},
		{Name: "test", Conclusion: "failure", Link: "https://github.com/owner/repo/actions/runs/456"},
		{Name: "lint", Conclusion: "success", Link: "https://github.com/owner/repo/actions/runs/789"},
	}
	m.cache.ciCache.Set("feat", checks)

	// Set up file tree
	m.setStatusFiles([]StatusFile{
		{Filename: "file1.go", Status: ".M", IsUntracked: false},
	})
	m.state.services.statusTree.Index = 0

	// Start at first CI check
	m.ciCheckIndex = 1

	// Navigate up with p
	_, _ = m.navigateCICheckUp()
	if m.ciCheckIndex != 0 {
		t.Fatalf("expected ciCheckIndex 0 after p, got %d", m.ciCheckIndex)
	}

	// At first CI check, should stay at 0
	_, _ = m.navigateCICheckUp()
	if m.ciCheckIndex != 0 {
		t.Fatalf("expected ciCheckIndex to stay at 0, got %d", m.ciCheckIndex)
	}

	// When ciCheckIndex is -1 (no CI check selected), navigating up should jump to last CI check
	m.ciCheckIndex = -1
	_, _ = m.navigateCICheckUp()
	if m.ciCheckIndex != 2 {
		t.Fatalf("expected ciCheckIndex 2 after jumping to last check, got %d", m.ciCheckIndex)
	}
}

func TestCICheckSelectionResetOnWorktreeChange(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.data.selectedIndex = 0
	m.state.data.filteredWts = []*models.WorktreeInfo{
		{Path: testWorktreePath, Branch: "feat"},
		{Path: "/other/path", Branch: "other"},
	}

	// Add CI checks to cache
	checks := []*models.CICheck{
		{Name: "build", Conclusion: "success", Link: "https://github.com/owner/repo/actions/runs/123"},
	}
	m.cache.ciCache.Set("feat", checks)
	m.ciCheckIndex = 0

	// Change worktree selection
	m.state.data.selectedIndex = 1
	m.updateDetailsView()

	// CI check selection should be reset
	if m.ciCheckIndex != -1 {
		t.Fatalf("expected ciCheckIndex -1 after worktree change, got %d", m.ciCheckIndex)
	}
}

func TestCICheckSelectionResetOnPaneSwitch(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 1
	m.state.data.selectedIndex = 0
	m.state.data.filteredWts = []*models.WorktreeInfo{
		{Path: testWorktreePath, Branch: "feat"},
	}

	// Add CI checks to cache
	checks := []*models.CICheck{
		{Name: "build", Conclusion: "success", Link: "https://github.com/owner/repo/actions/runs/123"},
	}
	m.cache.ciCache.Set("feat", checks)
	m.ciCheckIndex = 0

	// Switch to pane 0
	_, _ = m.handleBuiltInKey(tea.KeyPressMsg{Code: '1', Text: string('1')})

	// CI check selection should be reset
	if m.ciCheckIndex != -1 {
		t.Fatalf("expected ciCheckIndex -1 after pane switch, got %d", m.ciCheckIndex)
	}
}

func TestCICheckSelectionResetWhenUnavailable(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 1
	m.state.data.selectedIndex = 0
	m.state.data.filteredWts = []*models.WorktreeInfo{
		{Path: testWorktreePath, Branch: "feat"},
	}

	// Add CI checks to cache
	checks := []*models.CICheck{
		{Name: "build", Conclusion: "success", Link: "https://github.com/owner/repo/actions/runs/123"},
	}
	m.cache.ciCache.Set("feat", checks)
	m.ciCheckIndex = 5 // Out of bounds

	// Navigate with n - should reset invalid index to -1 then advance to first check (0)
	_, _ = m.navigateCICheckDown()
	if m.ciCheckIndex != 0 {
		t.Fatalf("expected ciCheckIndex 0 after invalid index reset, got %d", m.ciCheckIndex)
	}

	// Clear CI cache
	m.cache.ciCache.Clear()
	m.ciCheckIndex = 0

	// Navigate with n - should be no-op when no CI checks available
	_, _ = m.navigateCICheckDown()
	if m.ciCheckIndex != 0 {
		t.Fatalf("expected ciCheckIndex 0 when no CI checks (no-op), got %d", m.ciCheckIndex)
	}
}

func TestCICheckNavigationWithNoChecks(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 1
	m.state.ui.statusViewport = viewport.New(viewport.WithWidth(40), viewport.WithHeight(10))
	m.state.data.selectedIndex = 0
	m.state.data.filteredWts = []*models.WorktreeInfo{
		{Path: testWorktreePath, Branch: "feat"},
	}

	// No CI checks in cache - navigateCICheckDown should be a no-op
	m.ciCheckIndex = -1

	// Navigation with n should do nothing when no CI checks
	_, _ = m.navigateCICheckDown()
	if m.ciCheckIndex != -1 {
		t.Fatalf("expected ciCheckIndex -1 with no CI checks, got %d", m.ciCheckIndex)
	}
}

func TestInfoPaneNPNavigatesCIChecks(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 1
	m.state.data.selectedIndex = 0
	m.state.data.filteredWts = []*models.WorktreeInfo{
		{Path: testWorktreePath, Branch: "feat"},
	}

	checks := []*models.CICheck{
		{Name: "build", Conclusion: "success", Link: "https://github.com/owner/repo/actions/runs/123"},
		{Name: "test", Conclusion: "failure", Link: "https://github.com/owner/repo/actions/runs/456"},
		{Name: "lint", Conclusion: "success", Link: "https://github.com/owner/repo/actions/runs/789"},
	}
	m.cache.ciCache.Set("feat", checks)
	m.ciCheckIndex = -1

	// Press n — should navigate to first CI check
	_, _ = m.handleBuiltInKey(tea.KeyPressMsg{Code: 'n', Text: string('n')})
	if m.ciCheckIndex != 0 {
		t.Fatalf("expected ciCheckIndex 0 after n, got %d", m.ciCheckIndex)
	}

	// Press n again — should navigate to second CI check
	_, _ = m.handleBuiltInKey(tea.KeyPressMsg{Code: 'n', Text: string('n')})
	if m.ciCheckIndex != 1 {
		t.Fatalf("expected ciCheckIndex 1 after second n, got %d", m.ciCheckIndex)
	}

	// Press p — should navigate back to first CI check
	_, _ = m.handleBuiltInKey(tea.KeyPressMsg{Code: 'p', Text: string('p')})
	if m.ciCheckIndex != 0 {
		t.Fatalf("expected ciCheckIndex 0 after p, got %d", m.ciCheckIndex)
	}
}

func TestNKeySearchAdvanceOutsidePane1(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 0
	m.state.data.selectedIndex = 0
	m.state.data.filteredWts = []*models.WorktreeInfo{
		{Path: testWorktreePath, Branch: "feat"},
	}

	// Add CI checks to cache
	checks := []*models.CICheck{
		{Name: "build", Conclusion: "success", Link: "https://github.com/owner/repo/actions/runs/123"},
	}
	m.cache.ciCache.Set("feat", checks)
	m.ciCheckIndex = -1

	// Press n in pane 0 — should not navigate CI checks
	_, _ = m.handleBuiltInKey(tea.KeyPressMsg{Code: 'n', Text: string('n')})
	if m.ciCheckIndex != -1 {
		t.Fatalf("expected ciCheckIndex to remain -1 when not in pane 1, got %d", m.ciCheckIndex)
	}
}

// TestPRDataLoadedSyncAfterWorktreeReload verifies prDataLoaded is set when PR state is restored.
