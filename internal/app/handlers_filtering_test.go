package app

import (
	"path/filepath"
	"testing"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/chmouel/lazyworktree/internal/config"
	"github.com/chmouel/lazyworktree/internal/models"
)

func TestFilterEnterClosesWithoutSelecting(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
		SortMode:    "path",
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 0

	m.state.data.worktrees = []*models.WorktreeInfo{
		{Path: filepath.Join(cfg.WorktreeDir, "b-worktree"), Branch: testFeat},
		{Path: filepath.Join(cfg.WorktreeDir, "a-worktree"), Branch: testFeat},
	}
	m.state.services.filter.FilterQuery = testFeat
	m.state.ui.filterInput.SetValue(testFeat)
	m.updateTable()
	m.state.view.ShowingFilter = true
	m.state.ui.filterInput.Focus()
	m.state.ui.worktreeTable.SetCursor(1)
	m.state.data.selectedIndex = 1

	updated, cmd := m.handleKeyMsg(tea.KeyPressMsg{Code: tea.KeyEnter})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if cmd != nil {
		t.Fatal("expected no command to be returned")
	}
	if m.state.view.ShowingFilter {
		t.Fatal("expected filter to be closed")
	}
	if m.selectedPath != "" {
		t.Fatalf("expected selected path to remain empty, got %q", m.selectedPath)
	}
}

func TestFilterAltNPMovesSelectionAndFills(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
		SortMode:    "path",
	}
	m := NewModel(cfg, "")

	wt1Path := filepath.Join(cfg.WorktreeDir, testWt1)
	wt2Path := filepath.Join(cfg.WorktreeDir, testWt2)
	m.state.data.worktrees = []*models.WorktreeInfo{
		{Path: wt1Path, Branch: "feat-one"},
		{Path: wt2Path, Branch: "feat-two"},
	}
	m.state.services.filter.FilterQuery = testFeat
	m.state.ui.filterInput.SetValue(testFeat)
	m.updateTable()
	m.state.view.ShowingFilter = true
	m.state.ui.filterInput.Focus()
	m.state.ui.worktreeTable.SetCursor(0)
	m.state.data.selectedIndex = 0

	updated, _ := m.handleKeyMsg(tea.KeyPressMsg{Code: 'n', Mod: tea.ModAlt})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.ui.filterInput.Value() != testWt2 || m.state.services.filter.FilterQuery != testWt2 {
		t.Fatalf("expected filter query to match selected worktree, got %q", m.state.services.filter.FilterQuery)
	}
	if len(m.state.data.filteredWts) != 1 || m.state.data.filteredWts[0].Path != wt2Path {
		t.Fatalf("expected filtered worktree %q, got %v", wt2Path, m.state.data.filteredWts)
	}

	updated, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: 'p', Mod: tea.ModAlt})
	updatedModel, ok = updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.ui.filterInput.Value() != testWt1 || m.state.services.filter.FilterQuery != testWt1 {
		t.Fatalf("expected filter query to match selected worktree, got %q", m.state.services.filter.FilterQuery)
	}
	if len(m.state.data.filteredWts) != 1 || m.state.data.filteredWts[0].Path != wt1Path {
		t.Fatalf("expected filtered worktree %q, got %v", wt1Path, m.state.data.filteredWts)
	}
}

func TestFilterArrowKeysNavigateWithoutFilling(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
		SortMode:    "path",
	}
	m := NewModel(cfg, "")

	wt1Path := filepath.Join(cfg.WorktreeDir, testWt1)
	wt2Path := filepath.Join(cfg.WorktreeDir, testWt2)
	m.state.data.worktrees = []*models.WorktreeInfo{
		{Path: wt1Path, Branch: "feat-one"},
		{Path: wt2Path, Branch: "feat-two"},
	}
	m.state.services.filter.FilterQuery = testFeat
	m.state.ui.filterInput.SetValue(testFeat)
	m.updateTable()
	m.state.view.ShowingFilter = true
	m.state.ui.filterInput.Focus()
	m.state.ui.worktreeTable.SetCursor(0)
	m.state.data.selectedIndex = 0

	updated, _ := m.handleKeyMsg(tea.KeyPressMsg{Code: tea.KeyDown})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.ui.filterInput.Value() != testFeat || m.state.services.filter.FilterQuery != testFeat {
		t.Fatalf("expected filter query unchanged, got %q", m.state.services.filter.FilterQuery)
	}

	updated, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: tea.KeyUp})
	updatedModel, ok = updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.ui.filterInput.Value() != testFeat || m.state.services.filter.FilterQuery != testFeat {
		t.Fatalf("expected filter query unchanged, got %q", m.state.services.filter.FilterQuery)
	}
}

func TestFilterEmptyEnterSelectsCurrent(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
		SortMode:    "path",
	}
	m := NewModel(cfg, "")

	wt1Path := filepath.Join(cfg.WorktreeDir, testWt1)
	wt2Path := filepath.Join(cfg.WorktreeDir, testWt2)
	m.state.data.worktrees = []*models.WorktreeInfo{
		{Path: wt1Path, Branch: "feat-one"},
		{Path: wt2Path, Branch: "feat-two"},
	}
	m.state.services.filter.FilterQuery = ""
	m.state.ui.filterInput.SetValue("")
	m.updateTable()
	m.state.view.ShowingFilter = true
	m.state.ui.filterInput.Focus()
	m.state.ui.worktreeTable.SetCursor(1)
	m.state.data.selectedIndex = 1

	updated, _ := m.handleKeyMsg(tea.KeyPressMsg{Code: tea.KeyEnter})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.view.ShowingFilter {
		t.Fatal("expected filter to be closed")
	}
	if m.state.data.selectedIndex != 1 {
		t.Fatalf("expected selectedIndex to remain 1, got %d", m.state.data.selectedIndex)
	}
}

func TestFilterCtrlCExitsFilter(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
		SortMode:    "path",
	}
	m := NewModel(cfg, "")

	wt1Path := filepath.Join(cfg.WorktreeDir, testWt1)
	m.state.data.worktrees = []*models.WorktreeInfo{
		{Path: wt1Path, Branch: "feat-one"},
	}
	m.state.services.filter.FilterQuery = "something"
	m.state.ui.filterInput.SetValue("something")
	m.updateTable()
	m.state.view.ShowingFilter = true
	m.state.ui.filterInput.Focus()

	updated, _ := m.handleKeyMsg(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.view.ShowingFilter {
		t.Fatal("expected filter to be closed after Ctrl+C")
	}
	if m.state.ui.filterInput.Focused() {
		t.Fatal("expected filter input to be blurred")
	}
}

func TestFilterStatusNarrowsList(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 2
	m.state.ui.statusViewport = viewport.New(viewport.WithWidth(40), viewport.WithHeight(10))
	m.setStatusFiles([]StatusFile{
		{Filename: "app.go", Status: ".M"},
		{Filename: "README.md", Status: ".M"},
	})

	updated, _ := m.handleKeyMsg(tea.KeyPressMsg{Code: 'f', Text: string('f')})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	_, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: 'r', Text: string('r')})
	_, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: 'e', Text: string('e')})
	_, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: 'a', Text: string('a')})
	_, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: 'd', Text: string('d')})

	if len(m.state.data.statusFiles) != 1 {
		t.Fatalf("expected 1 filtered status file, got %d", len(m.state.data.statusFiles))
	}
	if m.state.data.statusFiles[0].Filename != testReadme {
		t.Fatalf("expected %s, got %q", testReadme, m.state.data.statusFiles[0].Filename)
	}
}

func TestFilterEnterClosesWithoutSelectingItem(t *testing.T) {
	// Set default provider for testing
	SetIconProvider(&NerdFontV3Provider{})
	cfg := &config.AppConfig{
		WorktreeDir:      t.TempDir(),
		SortMode:         "path",
		SearchAutoSelect: false,
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 0

	wt1Path := filepath.Join(cfg.WorktreeDir, "srv-api")
	wt2Path := filepath.Join(cfg.WorktreeDir, "srv-auth")
	wt3Path := filepath.Join(cfg.WorktreeDir, "srv-worker")
	m.state.data.worktrees = []*models.WorktreeInfo{
		{Path: wt1Path, Branch: "feature/srv-api"},
		{Path: wt2Path, Branch: "feature/srv-auth"},
		{Path: wt3Path, Branch: "feature/srv-worker"},
	}

	// Apply filter for "srv"
	m.state.services.filter.FilterQuery = "srv"
	m.state.ui.filterInput.SetValue("srv")
	m.updateTable()
	m.state.view.ShowingFilter = true
	m.state.ui.filterInput.Focus()

	// Navigate to the second item (srv-auth)
	m.state.ui.worktreeTable.SetCursor(1)
	m.state.data.selectedIndex = 1

	// Press Enter - should exit filter without selecting
	updated, cmd := m.handleKeyMsg(tea.KeyPressMsg{Code: tea.KeyEnter})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if cmd != nil {
		t.Fatal("expected no command to be returned")
	}
	if m.state.view.ShowingFilter {
		t.Fatal("expected filter to be closed")
	}
	if m.selectedPath != "" {
		t.Fatalf("expected selected path to remain empty, got %q", m.selectedPath)
	}
}

func TestFilterNavigationThroughMultipleFilteredItems(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
		SortMode:    "path",
	}
	m := NewModel(cfg, "")

	// Create 5 worktrees, 3 of which match "srv" filter
	wt1Path := filepath.Join(cfg.WorktreeDir, "main")
	wt2Path := filepath.Join(cfg.WorktreeDir, "srv-api")
	wt3Path := filepath.Join(cfg.WorktreeDir, "frontend")
	wt4Path := filepath.Join(cfg.WorktreeDir, "srv-auth")
	wt5Path := filepath.Join(cfg.WorktreeDir, "srv-worker")

	m.state.data.worktrees = []*models.WorktreeInfo{
		{Path: wt1Path, Branch: "main", IsMain: true},
		{Path: wt2Path, Branch: "feature/srv-api"},
		{Path: wt3Path, Branch: "feature/frontend"},
		{Path: wt4Path, Branch: "feature/srv-auth"},
		{Path: wt5Path, Branch: "feature/srv-worker"},
	}

	// Apply filter for "srv"
	m.state.services.filter.FilterQuery = "srv"
	m.state.ui.filterInput.SetValue("srv")
	m.updateTable()
	m.state.view.ShowingFilter = true
	m.state.ui.filterInput.Focus()
	m.state.ui.worktreeTable.SetCursor(0)
	m.state.data.selectedIndex = 0

	// Verify we have exactly 3 filtered items
	if len(m.state.data.filteredWts) != 3 {
		t.Fatalf("expected 3 filtered items, got %d", len(m.state.data.filteredWts))
	}

	// Navigate down through all filtered items
	for i := 0; i < 2; i++ {
		updated, _ := m.handleKeyMsg(tea.KeyPressMsg{Code: tea.KeyDown})
		updatedModel, ok := updated.(*Model)
		if !ok {
			t.Fatalf("expected updated model, got %T", updated)
		}
		m = updatedModel
	}

	// Should be at the last filtered item (index 2)
	cursor := m.state.ui.worktreeTable.Cursor()
	if cursor != 2 {
		t.Fatalf("expected cursor at index 2, got %d", cursor)
	}

	// Try to navigate down again - should stay at last item
	updated, _ := m.handleKeyMsg(tea.KeyPressMsg{Code: tea.KeyDown})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	cursor = m.state.ui.worktreeTable.Cursor()
	if cursor != 2 {
		t.Fatalf("expected cursor to stay at index 2, got %d", cursor)
	}

	// Navigate back up
	for i := 0; i < 2; i++ {
		updated, _ := m.handleKeyMsg(tea.KeyPressMsg{Code: tea.KeyUp})
		updatedModel, ok := updated.(*Model)
		if !ok {
			t.Fatalf("expected updated model, got %T", updated)
		}
		m = updatedModel
	}

	// Should be at the first filtered item (index 0)
	cursor = m.state.ui.worktreeTable.Cursor()
	if cursor != 0 {
		t.Fatalf("expected cursor at index 0, got %d", cursor)
	}

	// Try to navigate up again - should stay at first item
	updated, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: tea.KeyUp})
	updatedModel, ok = updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	cursor = m.state.ui.worktreeTable.Cursor()
	if cursor != 0 {
		t.Fatalf("expected cursor to stay at index 0, got %d", cursor)
	}
}

// TestStatusFileNavigation tests j/k navigation through status tree items in pane 2.

func TestFilterLogNarrowsList(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 3
	m.setLogEntries([]commitLogEntry{
		{sha: "abc123", authorInitials: "ab", message: "Fix bug in parser"},
		{sha: "def456", authorInitials: "de", message: "Add new feature"},
	}, false)

	updated, _ := m.handleKeyMsg(tea.KeyPressMsg{Code: 'f', Text: string('f')})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	_, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: 'f', Text: string('f')})
	_, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: 'i', Text: string('i')})
	_, _ = m.handleKeyMsg(tea.KeyPressMsg{Code: 'x', Text: string('x')})

	if len(m.state.data.logEntries) != 1 {
		t.Fatalf("expected 1 filtered commit, got %d", len(m.state.data.logEntries))
	}
	if m.state.data.logEntries[0].sha != "abc123" {
		t.Fatalf("expected commit abc123, got %q", m.state.data.logEntries[0].sha)
	}
}

// TestStatusFileNavigationEmptyList tests navigation with no status files.

func TestEscClearsWorktreeFilter(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 0
	m.state.services.filter.FilterQuery = testFilterQuery
	m.state.data.worktrees = []*models.WorktreeInfo{
		{Path: filepath.Join(cfg.WorktreeDir, "test-wt"), Branch: testFeat},
	}
	m.updateTable()

	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: tea.KeyEscape})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}

	if updatedModel.state.services.filter.FilterQuery != "" {
		t.Fatalf("expected filter to be cleared, got %q", updatedModel.state.services.filter.FilterQuery)
	}
}

func TestEscClearsStatusFilter(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 2
	m.state.services.filter.StatusFilterQuery = testFilterQuery

	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: tea.KeyEscape})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}

	if updatedModel.state.services.filter.StatusFilterQuery != "" {
		t.Fatalf("expected status filter to be cleared, got %q", updatedModel.state.services.filter.StatusFilterQuery)
	}
}

func TestEscClearsLogFilter(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 3
	m.state.services.filter.LogFilterQuery = testFilterQuery

	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: tea.KeyEscape})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}

	if updatedModel.state.services.filter.LogFilterQuery != "" {
		t.Fatalf("expected log filter to be cleared, got %q", updatedModel.state.services.filter.LogFilterQuery)
	}
}

func TestEscDoesNothingWhenNoFilter(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 0
	m.state.services.filter.FilterQuery = ""

	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: tea.KeyEscape})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}

	if updatedModel.state.services.filter.FilterQuery != "" {
		t.Fatalf("expected filter to remain empty, got %q", updatedModel.state.services.filter.FilterQuery)
	}
}

func TestHasActiveFilterForPane(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")

	// No filters active
	if m.hasActiveFilterForPane(0) {
		t.Fatal("expected no active filter for pane 0")
	}
	if m.hasActiveFilterForPane(1) {
		t.Fatal("expected no active filter for pane 1")
	}
	if m.hasActiveFilterForPane(2) {
		t.Fatal("expected no active filter for pane 2")
	}
	if m.hasActiveFilterForPane(3) {
		t.Fatal("expected no active filter for pane 3")
	}

	// Set worktree filter
	m.state.services.filter.FilterQuery = testFilterQuery
	if !m.hasActiveFilterForPane(0) {
		t.Fatal("expected active filter for pane 0")
	}

	// Pane 1 is info-only, no filter
	if m.hasActiveFilterForPane(1) {
		t.Fatal("expected no active filter for pane 1 (info-only)")
	}

	// Set status filter (pane 2)
	m.state.services.filter.StatusFilterQuery = testFilterQuery
	if !m.hasActiveFilterForPane(2) {
		t.Fatal("expected active filter for pane 2")
	}

	// Set log filter (pane 3)
	m.state.services.filter.LogFilterQuery = testFilterQuery
	if !m.hasActiveFilterForPane(3) {
		t.Fatal("expected active filter for pane 3")
	}

	// Whitespace-only should not count as active
	m.state.services.filter.FilterQuery = "   "
	if m.hasActiveFilterForPane(0) {
		t.Fatal("expected whitespace-only filter to not be active")
	}
}

func TestRestoreFocusAfterFilter(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")

	t.Run("restore focus to worktrees", func(t *testing.T) {
		m.state.view.FilterTarget = filterTargetWorktrees
		m.restoreFocusAfterFilter()
		if !m.state.ui.worktreeTable.Focused() {
			t.Error("expected worktreeTable to be focused")
		}
	})

	t.Run("restore focus to log", func(t *testing.T) {
		m.state.view.FilterTarget = filterTargetLog
		m.restoreFocusAfterFilter()
		if !m.state.ui.logTable.Focused() {
			t.Error("expected logTable to be focused")
		}
	})
}
