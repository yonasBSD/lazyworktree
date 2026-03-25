package app

import (
	"path/filepath"
	"strings"
	"testing"

	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/chmouel/lazyworktree/internal/app/state"
	"github.com/chmouel/lazyworktree/internal/config"
	"github.com/chmouel/lazyworktree/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestHandlePageDownUpOnStatusPane(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 2
	m.state.ui.statusViewport = viewport.New(viewport.WithWidth(10), viewport.WithHeight(2))
	m.state.ui.statusViewport.SetContent(strings.Repeat("line\n", 10))

	start := m.state.ui.statusViewport.YOffset()
	_, _ = m.handlePageDown(tea.KeyPressMsg{Code: tea.KeyPgDown})
	if m.state.ui.statusViewport.YOffset() <= start {
		t.Fatalf("expected YOffset to increase, got %d", m.state.ui.statusViewport.YOffset())
	}

	m.state.ui.statusViewport.SetYOffset(2)
	_, _ = m.handlePageUp(tea.KeyPressMsg{Code: tea.KeyPgUp})
	if m.state.ui.statusViewport.YOffset() >= 2 {
		t.Fatalf("expected YOffset to decrease, got %d", m.state.ui.statusViewport.YOffset())
	}
}

func TestHandlePageDownUpOnInfoPane(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 1
	m.state.ui.infoViewport = viewport.New(viewport.WithWidth(40), viewport.WithHeight(2))
	m.state.ui.infoViewport.SetContent(strings.Repeat("line\n", 20))

	start := m.state.ui.infoViewport.YOffset()
	_, _ = m.handlePageDown(tea.KeyPressMsg{Code: tea.KeyPgDown})
	assert.Greater(t, m.state.ui.infoViewport.YOffset(), start, "PageDown should increase YOffset")

	m.state.ui.infoViewport.SetYOffset(5)
	_, _ = m.handlePageUp(tea.KeyPressMsg{Code: tea.KeyPgUp})
	assert.Less(t, m.state.ui.infoViewport.YOffset(), 5, "PageUp should decrease YOffset")
}

func TestInfoViewportScrollDownUp(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 1
	m.state.ui.infoViewport = viewport.New(viewport.WithWidth(40), viewport.WithHeight(2))
	m.state.ui.infoViewport.SetContent(strings.Repeat("line\n", 20))

	// j/k should scroll when no CI checks
	m.state.data.filteredWts = []*models.WorktreeInfo{{Path: "/tmp/wt", Branch: "main"}}
	m.state.data.selectedIndex = 0

	start := m.state.ui.infoViewport.YOffset()
	_, _ = m.handleNavigationDown(tea.KeyPressMsg{Code: tea.KeyDown})
	assert.Greater(t, m.state.ui.infoViewport.YOffset(), start, "j should scroll down in info pane without CI checks")

	m.state.ui.infoViewport.SetYOffset(3)
	_, _ = m.handleNavigationUp(tea.KeyPressMsg{Code: tea.KeyUp})
	assert.Less(t, m.state.ui.infoViewport.YOffset(), 3, "k should scroll up in info pane without CI checks")
}

func TestInfoViewportGotoBottom(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 1
	m.state.ui.infoViewport = viewport.New(viewport.WithWidth(40), viewport.WithHeight(2))
	m.state.ui.infoViewport.SetContent(strings.Repeat("line\n", 20))

	_, _ = m.handleBuiltInKey(tea.KeyPressMsg{Code: -1, Text: "G"})
	assert.Positive(t, m.state.ui.infoViewport.YOffset(), "G should scroll to bottom in info pane")
}

// TestStatusFileNavigation tests j/k navigation through status tree items in pane 2.
func TestStatusFileNavigation(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 2
	m.state.ui.statusViewport = viewport.New(viewport.WithWidth(40), viewport.WithHeight(10))

	// Set up status files using setStatusFiles to build tree
	m.setStatusFiles([]StatusFile{
		{Filename: "file1.go", Status: ".M", IsUntracked: false},
		{Filename: "file2.go", Status: "M.", IsUntracked: false},
		{Filename: "file3.go", Status: " ?", IsUntracked: true},
	})
	m.state.services.statusTree.Index = 0

	// Test navigation down with j
	_, _ = m.handleNavigationDown(tea.KeyPressMsg{Code: 'j', Text: string('j')})
	if m.state.services.statusTree.Index != 1 {
		t.Fatalf("expected statusTreeIndex 1 after j, got %d", m.state.services.statusTree.Index)
	}

	// Test navigation down again
	_, _ = m.handleNavigationDown(tea.KeyPressMsg{Code: 'j', Text: string('j')})
	if m.state.services.statusTree.Index != 2 {
		t.Fatalf("expected statusTreeIndex 2 after second j, got %d", m.state.services.statusTree.Index)
	}

	// Test boundary - should not go past last item
	_, _ = m.handleNavigationDown(tea.KeyPressMsg{Code: 'j', Text: string('j')})
	if m.state.services.statusTree.Index != 2 {
		t.Fatalf("expected statusTreeIndex to stay at 2, got %d", m.state.services.statusTree.Index)
	}

	// Test navigation up with k
	_, _ = m.handleNavigationUp(tea.KeyPressMsg{Code: 'k', Text: string('k')})
	if m.state.services.statusTree.Index != 1 {
		t.Fatalf("expected statusTreeIndex 1 after k, got %d", m.state.services.statusTree.Index)
	}

	// Navigate to first item
	_, _ = m.handleNavigationUp(tea.KeyPressMsg{Code: 'k', Text: string('k')})
	if m.state.services.statusTree.Index != 0 {
		t.Fatalf("expected statusTreeIndex 0 after second k, got %d", m.state.services.statusTree.Index)
	}

	// Test boundary - should not go below 0
	_, _ = m.handleNavigationUp(tea.KeyPressMsg{Code: 'k', Text: string('k')})
	if m.state.services.statusTree.Index != 0 {
		t.Fatalf("expected statusTreeIndex to stay at 0, got %d", m.state.services.statusTree.Index)
	}
}

func TestLogPaneCtrlJMovesNextCommit(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 3
	m.state.ui.logTable.Focus()
	m.state.data.filteredWts = []*models.WorktreeInfo{
		{Path: t.TempDir(), Branch: testFeat},
	}
	m.state.data.selectedIndex = 0
	m.state.data.logEntries = []commitLogEntry{
		{sha: "abc123", authorInitials: "ab", message: "first"},
		{sha: "def456", authorInitials: "de", message: "second"},
	}
	m.state.ui.logTable.SetRows([]table.Row{
		{"abc123", "ab", "first"},
		{"def456", "de", "second"},
	})
	m.state.ui.logTable.SetCursor(0)

	updated, cmd := m.handleBuiltInKey(tea.KeyPressMsg{Code: 'j', Mod: tea.ModCtrl})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.ui.logTable.Cursor() != 1 {
		t.Fatalf("expected log cursor at 1, got %d", m.state.ui.logTable.Cursor())
	}
	if cmd == nil {
		t.Fatal("expected command to be returned")
	}
	// The command now returns a commitFilesLoadedMsg instead of calling execProcess
	// since openCommitView shows the files screen first
	msg := cmd()
	if _, ok := msg.(commitFilesLoadedMsg); !ok {
		t.Fatalf("expected commitFilesLoadedMsg, got %T", msg)
	}
}

// TestStatusFileNavigationEmptyList tests navigation with no status files.
func TestStatusFileNavigationEmptyList(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 2
	m.state.ui.statusViewport = viewport.New(viewport.WithWidth(40), viewport.WithHeight(10))
	m.setStatusFiles(nil)
	m.state.services.statusTree.Index = 0

	// Should not panic with empty list
	_, _ = m.handleNavigationDown(tea.KeyPressMsg{Code: 'j', Text: string('j')})
	if m.state.services.statusTree.Index != 0 {
		t.Fatalf("expected statusTreeIndex to stay at 0, got %d", m.state.services.statusTree.Index)
	}

	_, _ = m.handleNavigationUp(tea.KeyPressMsg{Code: 'k', Text: string('k')})
	if m.state.services.statusTree.Index != 0 {
		t.Fatalf("expected statusTreeIndex to stay at 0, got %d", m.state.services.statusTree.Index)
	}
}

// TestStatusFileEnterShowsDiff tests that Enter on pane 2 triggers showFileDiff.

func TestZoomPaneToggle(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 0

	if m.state.view.ZoomedPane != -1 {
		t.Fatalf("expected zoomedPane to start at -1, got %d", m.state.view.ZoomedPane)
	}

	// Press = to zoom pane 0
	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: '=', Text: string('=')})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.view.ZoomedPane != 0 {
		t.Fatalf("expected zoomedPane to be 0 after zoom, got %d", m.state.view.ZoomedPane)
	}

	// Press = again to unzoom
	updated, _ = m.handleBuiltInKey(tea.KeyPressMsg{Code: '=', Text: string('=')})
	updatedModel, ok = updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.view.ZoomedPane != -1 {
		t.Fatalf("expected zoomedPane to be -1 after unzoom, got %d", m.state.view.ZoomedPane)
	}
}

func TestZoomPaneExitsOnPaneKeys(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 0
	m.state.view.ZoomedPane = 0

	// Press 2 to switch to pane 2 and exit zoom
	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: '2', Text: string('2')})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.view.ZoomedPane != -1 {
		t.Fatalf("expected zoomedPane to be -1 after pressing 2, got %d", m.state.view.ZoomedPane)
	}
	if m.state.view.FocusedPane != 1 {
		t.Fatalf("expected focusedPane to be 1, got %d", m.state.view.FocusedPane)
	}
}

func TestZoomPaneExitsOnTabKey(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 0
	m.state.view.ZoomedPane = 0

	// Press tab to cycle panes and exit zoom
	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: tea.KeyTab})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.view.ZoomedPane != -1 {
		t.Fatalf("expected zoomedPane to be -1 after pressing tab, got %d", m.state.view.ZoomedPane)
	}
	if m.state.view.FocusedPane != 1 {
		t.Fatalf("expected focusedPane to be 1 after tab, got %d", m.state.view.FocusedPane)
	}
}

func TestTabChangesDefaultLayoutGeometryWithFocus(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.WindowWidth = 120
	m.state.view.WindowHeight = 40
	m.state.view.FocusedPane = 0
	m.state.data.statusFilesAll = []StatusFile{{Filename: "file.go", Status: ".M"}}

	before := m.computeLayout()

	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: tea.KeyTab})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	after := m.computeLayout()
	if before.leftWidth == after.leftWidth && before.rightWidth == after.rightWidth &&
		before.rightTopHeight == after.rightTopHeight && before.rightMiddleHeight == after.rightMiddleHeight &&
		before.rightBottomHeight == after.rightBottomHeight {
		t.Fatalf("expected pane geometry to change after tab when auto resize is enabled: before=%+v after=%+v", before, after)
	}
}

func TestTabLeavingGitStatusRefreshesCachedStatusContent(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 2
	m.state.ui.statusViewport = viewport.New(viewport.WithWidth(40), viewport.WithHeight(10))
	m.setStatusFiles([]StatusFile{
		{Filename: "file1.go", Status: ".M"},
		{Filename: "file2.go", Status: ".M"},
	})
	m.state.services.statusTree.Index = 1
	m.rebuildStatusContentWithHighlight()

	focusedContent := m.statusContent

	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: tea.KeyTab})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.view.FocusedPane != 3 {
		t.Fatalf("expected focusedPane to be 3 after leaving git status with tab, got %d", m.state.view.FocusedPane)
	}
	if m.statusContent == focusedContent {
		t.Fatal("expected cached status content to refresh after leaving git status pane")
	}
	if m.statusContent != m.renderStatusFiles() {
		t.Fatal("expected cached status content to match current render after pane switch")
	}
}

func TestZoomPaneExitsOnBracketKey(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 1
	m.state.view.ZoomedPane = 1

	// Press [ to cycle back and exit zoom
	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: '[', Text: string('[')})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.view.ZoomedPane != -1 {
		t.Fatalf("expected zoomedPane to be -1 after pressing [, got %d", m.state.view.ZoomedPane)
	}
	if m.state.view.FocusedPane != 0 {
		t.Fatalf("expected focusedPane to be 0 after [, got %d", m.state.view.FocusedPane)
	}
}

func TestPaneKey1ToggleZoom(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 0
	m.state.view.ZoomedPane = -1

	// Press 1 while on pane 0, not zoomed - should zoom
	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: '1', Text: string('1')})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.view.FocusedPane != 0 {
		t.Fatalf("expected focusedPane to remain 0, got %d", m.state.view.FocusedPane)
	}
	if m.state.view.ZoomedPane != 0 {
		t.Fatalf("expected zoomedPane to be 0 after toggle, got %d", m.state.view.ZoomedPane)
	}

	// Press 1 again while on pane 0, already zoomed - should unzoom
	updated, _ = m.handleBuiltInKey(tea.KeyPressMsg{Code: '1', Text: string('1')})
	updatedModel, ok = updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.view.FocusedPane != 0 {
		t.Fatalf("expected focusedPane to remain 0, got %d", m.state.view.FocusedPane)
	}
	if m.state.view.ZoomedPane != -1 {
		t.Fatalf("expected zoomedPane to be -1 after unzoom, got %d", m.state.view.ZoomedPane)
	}
}

func TestPaneKey2ToggleZoom(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 1
	m.state.view.ZoomedPane = -1

	// Press 2 while on pane 1, not zoomed - should zoom
	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: '2', Text: string('2')})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.view.FocusedPane != 1 {
		t.Fatalf("expected focusedPane to remain 1, got %d", m.state.view.FocusedPane)
	}
	if m.state.view.ZoomedPane != 1 {
		t.Fatalf("expected zoomedPane to be 1 after toggle, got %d", m.state.view.ZoomedPane)
	}

	// Press 2 again while on pane 1, already zoomed - should unzoom
	updated, _ = m.handleBuiltInKey(tea.KeyPressMsg{Code: '2', Text: string('2')})
	updatedModel, ok = updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.view.FocusedPane != 1 {
		t.Fatalf("expected focusedPane to remain 1, got %d", m.state.view.FocusedPane)
	}
	if m.state.view.ZoomedPane != -1 {
		t.Fatalf("expected zoomedPane to be -1 after unzoom, got %d", m.state.view.ZoomedPane)
	}
}

func TestPaneKey3ToggleZoom(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 2
	m.state.view.ZoomedPane = -1
	// Add status files so git status pane is available
	m.state.data.statusFilesAll = []StatusFile{{Filename: "file.go", Status: ".M"}}

	// Press 3 while on pane 2, not zoomed - should zoom
	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: '3', Text: string('3')})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.view.FocusedPane != 2 {
		t.Fatalf("expected focusedPane to remain 2, got %d", m.state.view.FocusedPane)
	}
	if m.state.view.ZoomedPane != 2 {
		t.Fatalf("expected zoomedPane to be 2 after toggle, got %d", m.state.view.ZoomedPane)
	}

	// Press 3 again while on pane 2, already zoomed - should unzoom
	updated, _ = m.handleBuiltInKey(tea.KeyPressMsg{Code: '3', Text: string('3')})
	updatedModel, ok = updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.view.FocusedPane != 2 {
		t.Fatalf("expected focusedPane to remain 2, got %d", m.state.view.FocusedPane)
	}
	if m.state.view.ZoomedPane != -1 {
		t.Fatalf("expected zoomedPane to be -1 after unzoom, got %d", m.state.view.ZoomedPane)
	}
}

func TestPaneKeyCrossPaneSwitching(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 0
	m.state.view.ZoomedPane = 0
	// Add status files so git status pane is available
	m.state.data.statusFilesAll = []StatusFile{{Filename: "file.go", Status: ".M"}}

	// Press 2 while on pane 0 (zoomed) - should switch to pane 1 and exit zoom
	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: '2', Text: string('2')})
	updatedModel, ok := updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.view.FocusedPane != 1 {
		t.Fatalf("expected focusedPane to be 1, got %d", m.state.view.FocusedPane)
	}
	if m.state.view.ZoomedPane != -1 {
		t.Fatalf("expected zoomedPane to be -1 after switching panes, got %d", m.state.view.ZoomedPane)
	}

	// Now press 3 while on pane 1 (not zoomed) - should switch to pane 2 and remain unzoomed
	updated, _ = m.handleBuiltInKey(tea.KeyPressMsg{Code: '3', Text: string('3')})
	updatedModel, ok = updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.view.FocusedPane != 2 {
		t.Fatalf("expected focusedPane to be 2, got %d", m.state.view.FocusedPane)
	}
	if m.state.view.ZoomedPane != -1 {
		t.Fatalf("expected zoomedPane to remain -1, got %d", m.state.view.ZoomedPane)
	}

	// Press 1 while on pane 2 (not zoomed) - should switch to pane 0
	updated, _ = m.handleBuiltInKey(tea.KeyPressMsg{Code: '1', Text: string('1')})
	updatedModel, ok = updated.(*Model)
	if !ok {
		t.Fatalf("expected updated model, got %T", updated)
	}
	m = updatedModel

	if m.state.view.FocusedPane != 0 {
		t.Fatalf("expected focusedPane to be 0, got %d", m.state.view.FocusedPane)
	}
	if m.state.view.ZoomedPane != -1 {
		t.Fatalf("expected zoomedPane to remain -1, got %d", m.state.view.ZoomedPane)
	}
}

// Message handler tests

func TestHandleGotoTop(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")

	t.Run("goto top on worktree pane", func(t *testing.T) {
		m.state.view.FocusedPane = 0
		m.state.data.worktrees = []*models.WorktreeInfo{
			{Path: filepath.Join(cfg.WorktreeDir, "wt1"), Branch: "branch1"},
			{Path: filepath.Join(cfg.WorktreeDir, "wt2"), Branch: "branch2"},
			{Path: filepath.Join(cfg.WorktreeDir, "wt3"), Branch: "branch3"},
		}
		m.state.ui.worktreeTable.SetWidth(100)
		m.updateTable()
		m.updateTableColumns(m.state.ui.worktreeTable.Width())
		if len(m.state.ui.worktreeTable.Rows()) == 0 {
			t.Fatal("table has no rows after updateTable")
		}
		m.state.ui.worktreeTable.SetCursor(2)
		_, cmd := m.handleGotoTop()
		if m.state.ui.worktreeTable.Cursor() != 0 {
			t.Errorf("expected cursor at top (0), got %d", m.state.ui.worktreeTable.Cursor())
		}
		if cmd == nil {
			t.Error("expected command to be returned")
		}
	})

	t.Run("goto top on status pane", func(t *testing.T) {
		m.state.view.FocusedPane = 2
		m.state.services.statusTree.TreeFlat = []*StatusTreeNode{
			{Path: "file1.txt", File: &StatusFile{Filename: "file1.txt"}},
			{Path: "dir1", Children: []*StatusTreeNode{}},
			{Path: "file2.txt", File: &StatusFile{Filename: "file2.txt"}},
		}
		m.state.services.statusTree.Index = 2
		_, cmd := m.handleGotoTop()
		if m.state.services.statusTree.Index != 0 {
			t.Errorf("expected statusTreeIndex to be 0, got %d", m.state.services.statusTree.Index)
		}
		if cmd != nil {
			t.Error("expected nil command for status pane")
		}
	})

	t.Run("goto top on log pane", func(t *testing.T) {
		m.state.view.FocusedPane = 3
		m.state.data.logEntries = []commitLogEntry{
			{sha: "abc123", message: "commit 1"},
			{sha: "def456", message: "commit 2"},
			{sha: "ghi789", message: "commit 3"},
		}
		m.setLogEntries(m.state.data.logEntries, false)
		m.updateLogColumns(m.state.ui.logTable.Width())
		m.state.ui.logTable.SetCursor(2)
		_, cmd := m.handleGotoTop()
		if m.state.ui.logTable.Cursor() != 0 {
			t.Errorf("expected cursor at top, got %d", m.state.ui.logTable.Cursor())
		}
		// Command may or may not be nil - just verify the cursor moved
		_ = cmd
	})
}

func TestHandleGotoBottom(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")

	t.Run("goto bottom on worktree pane", func(t *testing.T) {
		m.state.view.FocusedPane = 0
		m.state.ui.worktreeTable.SetCursor(0)
		m.state.data.filteredWts = []*models.WorktreeInfo{
			{Path: "wt1", Branch: "branch1"},
			{Path: "wt2", Branch: "branch2"},
			{Path: "wt3", Branch: "branch3"},
		}
		m.updateTable()
		_, cmd := m.handleGotoBottom()
		expectedBottom := len(m.state.data.filteredWts) - 1
		if m.state.ui.worktreeTable.Cursor() != expectedBottom {
			t.Errorf("expected cursor at bottom (%d), got %d", expectedBottom, m.state.ui.worktreeTable.Cursor())
		}
		if cmd == nil {
			t.Error("expected command to be returned")
		}
	})

	t.Run("goto bottom on status pane", func(t *testing.T) {
		m.state.view.FocusedPane = 2
		m.state.services.statusTree.TreeFlat = []*StatusTreeNode{
			{Path: "file1.txt", File: &StatusFile{Filename: "file1.txt"}},
			{Path: "dir1", Children: []*StatusTreeNode{}},
			{Path: "file2.txt", File: &StatusFile{Filename: "file2.txt"}},
		}
		m.state.services.statusTree.Index = 0
		_, cmd := m.handleGotoBottom()
		expectedBottom := len(m.state.services.statusTree.TreeFlat) - 1
		if m.state.services.statusTree.Index != expectedBottom {
			t.Errorf("expected statusTreeIndex to be %d, got %d", expectedBottom, m.state.services.statusTree.Index)
		}
		// Command may or may not be nil - just verify the index changed
		_ = cmd
	})

	t.Run("goto bottom on empty status pane", func(t *testing.T) {
		m.state.view.FocusedPane = 2
		m.state.services.statusTree.TreeFlat = []*StatusTreeNode{}
		m.state.services.statusTree.Index = 0
		_, cmd := m.handleGotoBottom()
		if m.state.services.statusTree.Index != 0 {
			t.Errorf("expected statusTreeIndex to remain 0, got %d", m.state.services.statusTree.Index)
		}
		_ = cmd
	})

	t.Run("goto bottom on log pane", func(t *testing.T) {
		m.state.view.FocusedPane = 3
		m.state.ui.logTable.SetCursor(0)
		m.state.data.logEntries = []commitLogEntry{
			{sha: "abc123", message: "commit 1"},
			{sha: "def456", message: "commit 2"},
			{sha: "ghi789", message: "commit 3"},
		}
		m.setLogEntries(m.state.data.logEntries, false)
		m.updateLogColumns(m.state.ui.logTable.Width())
		_, cmd := m.handleGotoBottom()
		expectedBottom := len(m.state.data.logEntries) - 1
		if m.state.ui.logTable.Cursor() != expectedBottom {
			t.Errorf("expected cursor at bottom (%d), got %d", expectedBottom, m.state.ui.logTable.Cursor())
		}
		// Command may or may not be nil - just verify the cursor moved
		_ = cmd
	})
}

func TestHandleNextFolder(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 2

	t.Run("empty status tree", func(t *testing.T) {
		m.state.services.statusTree.TreeFlat = []*StatusTreeNode{}
		m.state.services.statusTree.Index = 0
		_, cmd := m.handleNextFolder()
		if cmd != nil {
			t.Error("expected nil command for empty tree")
		}
	})

	t.Run("find next folder", func(t *testing.T) {
		m.state.services.statusTree.TreeFlat = []*StatusTreeNode{
			{Path: "file1.txt", File: &StatusFile{Filename: "file1.txt"}},
			{Path: "dir1", Children: []*StatusTreeNode{}},
			{Path: "file2.txt", File: &StatusFile{Filename: "file2.txt"}},
			{Path: "dir2", Children: []*StatusTreeNode{}},
		}
		m.state.services.statusTree.Index = 0
		_, cmd := m.handleNextFolder()
		if m.state.services.statusTree.Index != 1 {
			t.Errorf("expected statusTreeIndex to be 1 (dir1), got %d", m.state.services.statusTree.Index)
		}
		if cmd != nil {
			t.Error("expected nil command")
		}
	})

	t.Run("no next folder found", func(t *testing.T) {
		m.state.services.statusTree.TreeFlat = []*StatusTreeNode{
			{Path: "file1.txt", File: &StatusFile{Filename: "file1.txt"}},
			{Path: "dir1", Children: []*StatusTreeNode{}},
		}
		m.state.services.statusTree.Index = 1 // Already at last folder
		_, cmd := m.handleNextFolder()
		if m.state.services.statusTree.Index != 1 {
			t.Errorf("expected statusTreeIndex to remain 1, got %d", m.state.services.statusTree.Index)
		}
		if cmd != nil {
			t.Error("expected nil command")
		}
	})
}

func TestHandlePrevFolder(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 2

	t.Run("empty status tree", func(t *testing.T) {
		m.state.services.statusTree.TreeFlat = []*StatusTreeNode{}
		m.state.services.statusTree.Index = 0
		_, cmd := m.handlePrevFolder()
		if cmd != nil {
			t.Error("expected nil command for empty tree")
		}
	})

	t.Run("find previous folder", func(t *testing.T) {
		m.state.services.statusTree.TreeFlat = []*StatusTreeNode{
			{Path: "dir1", Children: []*StatusTreeNode{}},
			{Path: "file1.txt", File: &StatusFile{Filename: "file1.txt"}},
			{Path: "dir2", Children: []*StatusTreeNode{}},
			{Path: "file2.txt", File: &StatusFile{Filename: "file2.txt"}},
		}
		m.state.services.statusTree.Index = 3
		_, cmd := m.handlePrevFolder()
		if m.state.services.statusTree.Index != 2 {
			t.Errorf("expected statusTreeIndex to be 2 (dir2), got %d", m.state.services.statusTree.Index)
		}
		if cmd != nil {
			t.Error("expected nil command")
		}
	})

	t.Run("no previous folder found", func(t *testing.T) {
		m.state.services.statusTree.TreeFlat = []*StatusTreeNode{
			{Path: "dir1", Children: []*StatusTreeNode{}},
			{Path: "file1.txt", File: &StatusFile{Filename: "file1.txt"}},
		}
		m.state.services.statusTree.Index = 0 // Already at first folder
		_, cmd := m.handlePrevFolder()
		if m.state.services.statusTree.Index != 0 {
			t.Errorf("expected statusTreeIndex to remain 0, got %d", m.state.services.statusTree.Index)
		}
		if cmd != nil {
			t.Error("expected nil command")
		}
	})
}

func TestInfoPaneJKScrollsViewport(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 1
	m.state.data.selectedIndex = 0
	m.state.data.filteredWts = []*models.WorktreeInfo{
		{Path: testWorktreePath, Branch: "feat"},
	}

	// Add CI checks to cache — j/k should still scroll, not navigate CI checks
	checks := []*models.CICheck{
		{Name: "build", Conclusion: "success", Link: "https://github.com/owner/repo/actions/runs/123"},
		{Name: "test", Conclusion: "failure", Link: "https://github.com/owner/repo/actions/runs/456"},
	}
	m.cache.ciCache.Set("feat", checks)
	m.ciCheckIndex = -1

	// Press j — should scroll viewport, not change ciCheckIndex
	_, _ = m.handleNavigationDown(tea.KeyPressMsg{Code: 'j', Text: string('j')})
	if m.ciCheckIndex != -1 {
		t.Fatalf("expected ciCheckIndex to remain -1 after j, got %d", m.ciCheckIndex)
	}

	// Press k — should scroll viewport, not change ciCheckIndex
	_, _ = m.handleNavigationUp(tea.KeyPressMsg{Code: 'k', Text: string('k')})
	if m.ciCheckIndex != -1 {
		t.Fatalf("expected ciCheckIndex to remain -1 after k, got %d", m.ciCheckIndex)
	}
}

func TestNextPaneWithoutNotes(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.data.filteredWts = []*models.WorktreeInfo{{Path: "/tmp/wt", Branch: "main"}}
	m.state.data.selectedIndex = 0
	// Add status files so git status pane is included
	m.state.data.statusFilesAll = []StatusFile{{Filename: "file.go", Status: ".M"}}

	// Without notes, cycle is 0 -> 1 -> 2 -> 3 -> 0
	assert.Equal(t, 1, m.nextPane(0, 1))
	assert.Equal(t, 2, m.nextPane(1, 1))
	assert.Equal(t, 3, m.nextPane(2, 1))
	assert.Equal(t, 0, m.nextPane(3, 1))

	// Reverse
	assert.Equal(t, 3, m.nextPane(0, -1))
	assert.Equal(t, 0, m.nextPane(1, -1))
}

func TestNextPaneWithNotes(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	wt := &models.WorktreeInfo{Path: "/tmp/wt-notes", Branch: "feat"}
	m.state.data.filteredWts = []*models.WorktreeInfo{wt}
	m.state.data.selectedIndex = 0
	m.worktreeNotes[worktreeNoteKey(wt.Path)] = models.WorktreeNote{Note: "hello"}
	// Add status files so git status pane is included
	m.state.data.statusFilesAll = []StatusFile{{Filename: "file.go", Status: ".M"}}

	// With notes, cycle is 0 -> 1 -> 2 -> 3 -> 4 -> 0
	assert.Equal(t, 1, m.nextPane(0, 1))
	assert.Equal(t, 2, m.nextPane(1, 1))
	assert.Equal(t, 3, m.nextPane(2, 1))
	assert.Equal(t, 4, m.nextPane(3, 1))
	assert.Equal(t, 0, m.nextPane(4, 1))

	// Reverse
	assert.Equal(t, 4, m.nextPane(0, -1))
	assert.Equal(t, 3, m.nextPane(4, -1))
}

func TestPaneKey5FocusNotesPane(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	wt := &models.WorktreeInfo{Path: "/tmp/wt-k5", Branch: "feat"}
	m.state.data.filteredWts = []*models.WorktreeInfo{wt}
	m.state.data.selectedIndex = 0
	m.worktreeNotes[worktreeNoteKey(wt.Path)] = models.WorktreeNote{Note: "note text"}

	m.state.view.FocusedPane = 0
	m.state.view.ZoomedPane = -1

	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: '5', Text: "5"})
	m = updated.(*Model)

	assert.Equal(t, 4, m.state.view.FocusedPane)
	assert.Equal(t, -1, m.state.view.ZoomedPane)

	// Press 5 again to zoom
	updated, _ = m.handleBuiltInKey(tea.KeyPressMsg{Code: '5', Text: "5"})
	m = updated.(*Model)

	assert.Equal(t, 4, m.state.view.FocusedPane)
	assert.Equal(t, 4, m.state.view.ZoomedPane)
}

func TestPaneKey5NoopWithoutNote(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	wt := &models.WorktreeInfo{Path: "/tmp/wt-nk5", Branch: "feat"}
	m.state.data.filteredWts = []*models.WorktreeInfo{wt}
	m.state.data.selectedIndex = 0

	m.state.view.FocusedPane = 0
	m.state.view.ZoomedPane = -1

	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: '5', Text: "5"})
	m = updated.(*Model)

	// Should remain on pane 0 since no note exists
	assert.Equal(t, 0, m.state.view.FocusedPane)
}

func TestNotesPaneScrollDownUp(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 4
	m.state.ui.notesViewport.SetWidth(40)
	m.state.ui.notesViewport.SetHeight(2)
	m.state.ui.notesViewport.SetContent(strings.Repeat("line\n", 20))

	start := m.state.ui.notesViewport.YOffset()
	_, _ = m.handleNavigationDown(tea.KeyPressMsg{Code: tea.KeyDown})
	assert.Greater(t, m.state.ui.notesViewport.YOffset(), start, "j should scroll notes down")

	m.state.ui.notesViewport.SetYOffset(5)
	_, _ = m.handleNavigationUp(tea.KeyPressMsg{Code: tea.KeyUp})
	assert.Less(t, m.state.ui.notesViewport.YOffset(), 5, "k should scroll notes up")
}

func TestNotesPanePageDownUp(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 4
	m.state.ui.notesViewport.SetWidth(40)
	m.state.ui.notesViewport.SetHeight(2)
	m.state.ui.notesViewport.SetContent(strings.Repeat("line\n", 20))

	start := m.state.ui.notesViewport.YOffset()
	_, _ = m.handlePageDown(tea.KeyPressMsg{Code: tea.KeyPgDown})
	assert.Greater(t, m.state.ui.notesViewport.YOffset(), start, "PageDown should scroll notes")

	m.state.ui.notesViewport.SetYOffset(5)
	_, _ = m.handlePageUp(tea.KeyPressMsg{Code: tea.KeyPgUp})
	assert.Less(t, m.state.ui.notesViewport.YOffset(), 5, "PageUp should scroll notes up")
}

func TestNotesPaneGotoTopBottom(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 4
	m.state.ui.notesViewport.SetWidth(40)
	m.state.ui.notesViewport.SetHeight(2)
	m.state.ui.notesViewport.SetContent(strings.Repeat("line\n", 20))

	m.handleGotoBottom()
	assert.Positive(t, m.state.ui.notesViewport.YOffset(), "G should go to bottom")

	m.handleGotoTop()
	assert.Equal(t, 0, m.state.ui.notesViewport.YOffset(), "gg should go to top")
}

func TestTabCycleIncludesNotesPaneWhenNoteExists(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	wt := &models.WorktreeInfo{Path: "/tmp/wt-tab", Branch: "feat"}
	m.state.data.filteredWts = []*models.WorktreeInfo{wt}
	m.state.data.selectedIndex = 0
	m.worktreeNotes[worktreeNoteKey(wt.Path)] = models.WorktreeNote{Note: "a note"}
	// Add status files so git status pane is included
	m.state.data.statusFilesAll = []StatusFile{{Filename: "file.go", Status: ".M"}}
	m.state.view.FocusedPane = 0
	m.state.view.ZoomedPane = -1

	// Tab from pane 0 should go to pane 1 first; notes are last in the cycle.
	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(*Model)
	assert.Equal(t, 1, m.state.view.FocusedPane)

	// Continue through info -> git status -> commit -> notes.
	updated, _ = m.handleBuiltInKey(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(*Model)
	assert.Equal(t, 2, m.state.view.FocusedPane)
	updated, _ = m.handleBuiltInKey(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(*Model)
	assert.Equal(t, 3, m.state.view.FocusedPane)
	updated, _ = m.handleBuiltInKey(tea.KeyPressMsg{Code: tea.KeyTab})
	m = updated.(*Model)
	assert.Equal(t, 4, m.state.view.FocusedPane)
}

func TestNextPaneWithNotesAndAgentSessions(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	wt := &models.WorktreeInfo{Path: "/tmp/wt-agent-notes", Branch: "feat"}
	m.state.data.filteredWts = []*models.WorktreeInfo{wt}
	m.state.data.selectedIndex = 0
	m.worktreeNotes[worktreeNoteKey(wt.Path)] = models.WorktreeNote{Note: "hello"}
	m.state.data.statusFilesAll = []StatusFile{{Filename: "file.go", Status: ".M"}}
	seedClaudeAgentSession(t, m, wt.Path, true)

	assert.Equal(t, 1, m.nextPane(0, 1))
	assert.Equal(t, 2, m.nextPane(1, 1))
	assert.Equal(t, 3, m.nextPane(2, 1))
	assert.Equal(t, 4, m.nextPane(3, 1))
	assert.Equal(t, 5, m.nextPane(4, 1))
	assert.Equal(t, 0, m.nextPane(5, 1))
}

func TestPaneKey6FocusAgentSessionsPane(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	wt := &models.WorktreeInfo{Path: "/tmp/wt-agent-pane", Branch: "feat"}
	m.state.data.filteredWts = []*models.WorktreeInfo{wt}
	m.state.data.selectedIndex = 0
	seedClaudeAgentSession(t, m, wt.Path, true)

	m.state.view.FocusedPane = paneWorktrees
	m.state.view.ZoomedPane = -1

	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: '6', Text: "6"})
	m = updated.(*Model)
	assert.Equal(t, paneAgentSessions, m.state.view.FocusedPane)
	assert.Equal(t, -1, m.state.view.ZoomedPane)

	updated, _ = m.handleBuiltInKey(tea.KeyPressMsg{Code: '6', Text: "6"})
	m = updated.(*Model)
	assert.Equal(t, paneAgentSessions, m.state.view.FocusedPane)
	assert.Equal(t, paneAgentSessions, m.state.view.ZoomedPane)
}

func TestPaneKey6RevealsOfflineAgentSessions(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	wt := &models.WorktreeInfo{Path: "/tmp/wt-agent-offline", Branch: "feat"}
	m.state.data.filteredWts = []*models.WorktreeInfo{wt}
	m.state.data.selectedIndex = 0
	seedClaudeAgentSession(t, m, wt.Path, false)

	assert.False(t, m.hasAgentSessionsForSelectedWorktree())
	assert.True(t, m.hasAnyAgentSessionsForSelectedWorktree())

	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: '6', Text: "6"})
	m = updated.(*Model)
	assert.True(t, m.state.view.ShowAllAgentSessions)
	assert.Equal(t, paneAgentSessions, m.state.view.FocusedPane)
	assert.Len(t, m.state.data.agentSessions, 1)
	assert.False(t, m.state.data.agentSessions[0].IsOpen)
}

func TestAgentSessionsPaneToggleShowAll(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	wt := &models.WorktreeInfo{Path: "/tmp/wt-agent-toggle", Branch: "feat"}
	m.state.data.filteredWts = []*models.WorktreeInfo{wt}
	m.state.data.selectedIndex = 0
	seedClaudeAgentSession(t, m, wt.Path, false)
	m.state.view.ShowAllAgentSessions = true
	m.refreshSelectedWorktreeAgentSessionsPane()
	m.state.view.FocusedPane = paneAgentSessions

	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: 'A', Text: "A"})
	m = updated.(*Model)
	assert.False(t, m.state.view.ShowAllAgentSessions)
	assert.Equal(t, paneWorktrees, m.state.view.FocusedPane)
	assert.Empty(t, m.state.data.agentSessions)
}

func TestAgentSessionsPaneNavigationKeepsSelectionVisible(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	wt := &models.WorktreeInfo{Path: "/tmp/wt-agent-scroll", Branch: "feat"}
	m.state.data.filteredWts = []*models.WorktreeInfo{wt}
	m.state.data.selectedIndex = 0
	m.state.view.FocusedPane = paneAgentSessions
	m.state.ui.agentSessionsViewport.SetWidth(80)
	m.state.ui.agentSessionsViewport.SetHeight(2)
	seedClaudeAgentSessions(t, m, wt.Path, []bool{true, true, true, true})

	assert.Equal(t, 0, m.state.ui.agentSessionsViewport.YOffset())

	_, _ = m.handleNavigationDown(tea.KeyPressMsg{Code: tea.KeyDown})
	assert.Equal(t, 1, m.state.data.agentSessionIndex)
	assert.Equal(t, 1, m.state.ui.agentSessionsViewport.YOffset())

	_, _ = m.handleNavigationDown(tea.KeyPressMsg{Code: tea.KeyDown})
	assert.Equal(t, 2, m.state.data.agentSessionIndex)
	assert.Equal(t, 3, m.state.ui.agentSessionsViewport.YOffset())

	_, _ = m.handleNavigationUp(tea.KeyPressMsg{Code: tea.KeyUp})
	assert.Equal(t, 1, m.state.data.agentSessionIndex)
	assert.Equal(t, 2, m.state.ui.agentSessionsViewport.YOffset())
}

func TestAgentSessionsPanePageAndGotoAdjustViewport(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	wt := &models.WorktreeInfo{Path: "/tmp/wt-agent-page", Branch: "feat"}
	m.state.data.filteredWts = []*models.WorktreeInfo{wt}
	m.state.data.selectedIndex = 0
	m.state.view.FocusedPane = paneAgentSessions
	m.state.ui.agentSessionsViewport.SetWidth(80)
	m.state.ui.agentSessionsViewport.SetHeight(3)
	seedClaudeAgentSessions(t, m, wt.Path, []bool{true, true, true, true, true})

	_, _ = m.handlePageDown(tea.KeyPressMsg{Code: tea.KeyPgDown})
	assert.Equal(t, 3, m.state.data.agentSessionIndex)
	assert.Equal(t, 4, m.state.ui.agentSessionsViewport.YOffset())

	_, _ = m.handlePageUp(tea.KeyPressMsg{Code: tea.KeyPgUp})
	assert.Equal(t, 0, m.state.data.agentSessionIndex)
	assert.Equal(t, 0, m.state.ui.agentSessionsViewport.YOffset())

	_, _ = m.handleGotoBottom()
	assert.Equal(t, 4, m.state.data.agentSessionIndex)
	assert.Equal(t, 6, m.state.ui.agentSessionsViewport.YOffset())

	_, _ = m.handleGotoTop()
	assert.Equal(t, 0, m.state.data.agentSessionIndex)
	assert.Equal(t, 0, m.state.ui.agentSessionsViewport.YOffset())
}

func TestAgentSessionsPaneRefreshClampsViewportAfterShrink(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	wt := &models.WorktreeInfo{Path: "/tmp/wt-agent-shrink", Branch: "feat"}
	m.state.data.filteredWts = []*models.WorktreeInfo{wt}
	m.state.data.selectedIndex = 0
	m.state.view.FocusedPane = paneAgentSessions
	m.state.ui.agentSessionsViewport.SetWidth(80)
	m.state.ui.agentSessionsViewport.SetHeight(2)
	seedClaudeAgentSessions(t, m, wt.Path, []bool{true, true, true, true})

	_, _ = m.handleGotoBottom()
	assert.Equal(t, 5, m.state.ui.agentSessionsViewport.YOffset())

	seedClaudeAgentSessions(t, m, wt.Path, []bool{true})
	assert.Equal(t, 0, m.state.data.agentSessionIndex)
	assert.Equal(t, 0, m.state.ui.agentSessionsViewport.YOffset())
}

func TestHasNoteForSelectedWorktree(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")

	wt := &models.WorktreeInfo{Path: "/tmp/wt-has", Branch: "feat"}
	m.state.data.filteredWts = []*models.WorktreeInfo{wt}
	m.state.data.selectedIndex = 0

	assert.False(t, m.hasNoteForSelectedWorktree())

	m.worktreeNotes[worktreeNoteKey(wt.Path)] = models.WorktreeNote{Note: "hello"}
	assert.True(t, m.hasNoteForSelectedWorktree())

	// Empty note should not count
	m.worktreeNotes[worktreeNoteKey(wt.Path)] = models.WorktreeNote{Note: "  "}
	assert.False(t, m.hasNoteForSelectedWorktree())
}

func TestNextPaneWithoutGitStatus(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.data.filteredWts = []*models.WorktreeInfo{{Path: "/tmp/wt", Branch: "main"}}
	m.state.data.selectedIndex = 0
	// No status files => clean working tree

	// Without git status, cycle is 0 -> 1 -> 3 -> 0 (pane 2 skipped)
	assert.Equal(t, 1, m.nextPane(0, 1))
	assert.Equal(t, 3, m.nextPane(1, 1))
	assert.Equal(t, 0, m.nextPane(3, 1))

	// Reverse
	assert.Equal(t, 3, m.nextPane(0, -1))
	assert.Equal(t, 1, m.nextPane(3, -1))
	assert.Equal(t, 0, m.nextPane(1, -1))
}

func TestNextPaneWithGitStatus(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.data.filteredWts = []*models.WorktreeInfo{{Path: "/tmp/wt", Branch: "main"}}
	m.state.data.selectedIndex = 0
	m.state.data.statusFilesAll = []StatusFile{{Filename: "file.go", Status: ".M"}}

	// With git status, cycle is 0 -> 1 -> 2 -> 3 -> 0
	assert.Equal(t, 1, m.nextPane(0, 1))
	assert.Equal(t, 2, m.nextPane(1, 1))
	assert.Equal(t, 3, m.nextPane(2, 1))
	assert.Equal(t, 0, m.nextPane(3, 1))

	// Reverse
	assert.Equal(t, 3, m.nextPane(0, -1))
	assert.Equal(t, 0, m.nextPane(1, -1))
}

func TestPaneKey3IgnoredWhenClean(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.data.filteredWts = []*models.WorktreeInfo{{Path: "/tmp/wt", Branch: "main"}}
	m.state.data.selectedIndex = 0

	m.state.view.FocusedPane = 0
	m.state.view.ZoomedPane = -1

	// Press 3 with clean status — should do nothing
	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: '3', Text: "3"})
	um := updated.(*Model)

	assert.Equal(t, 0, um.state.view.FocusedPane)
	assert.Equal(t, -1, um.state.view.ZoomedPane)
}

func TestPaneKey3WorksWhenDirty(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.data.filteredWts = []*models.WorktreeInfo{{Path: "/tmp/wt", Branch: "main"}}
	m.state.data.selectedIndex = 0
	m.state.data.statusFilesAll = []StatusFile{{Filename: "file.go", Status: ".M"}}

	m.state.view.FocusedPane = 0
	m.state.view.ZoomedPane = -1

	// Press 3 with dirty status — should focus pane 2
	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: '3', Text: "3"})
	um := updated.(*Model)

	assert.Equal(t, 2, um.state.view.FocusedPane)
}

func TestFocusResetWhenStatusBecomesClean(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.data.filteredWts = []*models.WorktreeInfo{{Path: "/tmp/wt", Branch: "main"}}
	m.state.data.selectedIndex = 0

	// Start with dirty files and focus on git status pane
	m.setStatusFiles([]StatusFile{{Filename: "file.go", Status: ".M"}})
	m.state.view.FocusedPane = 2
	m.state.view.ZoomedPane = 2

	// Status becomes clean
	m.setStatusFiles(nil)

	// Focus should reset to pane 0 and zoom should clear
	assert.Equal(t, 0, m.state.view.FocusedPane)
	assert.Equal(t, -1, m.state.view.ZoomedPane)
}

func TestHasGitStatus(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")

	assert.False(t, m.hasGitStatus())

	m.state.data.statusFilesAll = []StatusFile{{Filename: "file.go", Status: ".M"}}
	assert.True(t, m.hasGitStatus())

	m.state.data.statusFilesAll = nil
	assert.False(t, m.hasGitStatus())
}

func TestHKeyDecrementsResizeOffset(t *testing.T) {
	t.Parallel()
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.view.WindowWidth = 120
	m.state.view.WindowHeight = 40
	m.state.data.filteredWts = []*models.WorktreeInfo{{Path: "/tmp/wt", Branch: "main"}}
	m.state.data.selectedIndex = 0

	assert.Equal(t, 0, m.state.view.ResizeOffset)

	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: 'h', Text: "h"})
	um := updated.(*Model)
	assert.Equal(t, -resizeStep, um.state.view.ResizeOffset)

	// Press again
	updated, _ = um.handleBuiltInKey(tea.KeyPressMsg{Code: 'h', Text: "h"})
	um = updated.(*Model)
	assert.Equal(t, -resizeStep*2, um.state.view.ResizeOffset)
}

func TestLKeyIncrementsResizeOffset(t *testing.T) {
	t.Parallel()
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.view.WindowWidth = 120
	m.state.view.WindowHeight = 40
	m.state.data.filteredWts = []*models.WorktreeInfo{{Path: "/tmp/wt", Branch: "main"}}
	m.state.data.selectedIndex = 0

	assert.Equal(t, 0, m.state.view.ResizeOffset)

	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: 'l', Text: "l"})
	um := updated.(*Model)
	assert.Equal(t, resizeStep, um.state.view.ResizeOffset)

	// Press again
	updated, _ = um.handleBuiltInKey(tea.KeyPressMsg{Code: 'l', Text: "l"})
	um = updated.(*Model)
	assert.Equal(t, resizeStep*2, um.state.view.ResizeOffset)
}

func TestResizeOffsetClampedAtBounds(t *testing.T) {
	t.Parallel()
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.view.WindowWidth = 120
	m.state.view.WindowHeight = 40
	m.state.data.filteredWts = []*models.WorktreeInfo{{Path: "/tmp/wt", Branch: "main"}}
	m.state.data.selectedIndex = 0

	t.Run("clamp at negative bound", func(t *testing.T) {
		m.state.view.ResizeOffset = -78
		updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: 'h', Text: "h"})
		um := updated.(*Model)
		assert.Equal(t, -80, um.state.view.ResizeOffset)
	})

	t.Run("clamp at positive bound", func(t *testing.T) {
		m.state.view.ResizeOffset = 78
		updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: 'l', Text: "l"})
		um := updated.(*Model)
		assert.Equal(t, 80, um.state.view.ResizeOffset)
	})
}

func TestLayoutToggleResetsResizeOffset(t *testing.T) {
	t.Parallel()
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.view.WindowWidth = 120
	m.state.view.WindowHeight = 40
	m.state.view.ResizeOffset = 20

	updated, _ := m.handleBuiltInKey(tea.KeyPressMsg{Code: 'L', Text: "L"})
	um := updated.(*Model)
	assert.Equal(t, 0, um.state.view.ResizeOffset)
	assert.Equal(t, state.LayoutTop, um.state.view.Layout)
}

func TestZoomPane2ResetsWhenClean(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.data.filteredWts = []*models.WorktreeInfo{{Path: "/tmp/wt", Branch: "main"}}
	m.state.data.selectedIndex = 0
	m.state.view.WindowWidth = 120
	m.state.view.WindowHeight = 40

	// Zoom on pane 2 with no status files
	m.state.view.ZoomedPane = 2
	m.state.view.FocusedPane = 2

	layout := m.computeLayout()
	m.renderBody(layout)

	// Zoom should have been reset
	assert.Equal(t, -1, m.state.view.ZoomedPane)
}
