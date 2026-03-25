package app

import (
	"strings"
	"testing"

	"charm.land/bubbles/v2/viewport"
	"github.com/chmouel/lazyworktree/internal/app/services"
	"github.com/chmouel/lazyworktree/internal/config"
)

// TestBuildStatusContentParsesFiles tests that buildStatusContent parses git status correctly.
func TestBuildStatusContentParsesFiles(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 2
	m.state.ui.statusViewport = viewport.New(viewport.WithWidth(40), viewport.WithHeight(10))

	// Simulated git status --porcelain=v2 output
	statusRaw := `1 .M N... 100644 100644 100644 abc123 abc123 modified.go
1 M. N... 100644 100644 100644 def456 def456 staged.go
? untracked.txt
1 A. N... 100644 100644 100644 ghi789 ghi789 added.go
1 .D N... 100644 100644 100644 jkl012 jkl012 deleted.go`

	m.setStatusFiles(parseStatusFiles(statusRaw))
	m.rebuildStatusContentWithHighlight()

	if len(m.state.data.statusFiles) != 5 {
		t.Fatalf("expected 5 status files, got %d", len(m.state.data.statusFiles))
	}

	// Check first file (modified)
	if m.state.data.statusFiles[0].Filename != "modified.go" {
		t.Fatalf("expected filename 'modified.go', got %q", m.state.data.statusFiles[0].Filename)
	}
	if m.state.data.statusFiles[0].Status != ".M" {
		t.Fatalf("expected status '.M', got %q", m.state.data.statusFiles[0].Status)
	}
	if m.state.data.statusFiles[0].IsUntracked {
		t.Fatal("expected IsUntracked to be false for modified file")
	}

	// Check untracked file
	if m.state.data.statusFiles[2].Filename != "untracked.txt" {
		t.Fatalf("expected filename 'untracked.txt', got %q", m.state.data.statusFiles[2].Filename)
	}
	if !m.state.data.statusFiles[2].IsUntracked {
		t.Fatal("expected IsUntracked to be true for untracked file")
	}
}

// TestBuildStatusContentCleanTree tests that clean working tree is handled.

// TestBuildStatusContentCleanTree tests that clean working tree is handled.
func TestBuildStatusContentCleanTree(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 2
	m.state.ui.statusViewport = viewport.New(viewport.WithWidth(40), viewport.WithHeight(10))
	m.state.data.statusFiles = []StatusFile{{Filename: "old.go", Status: ".M"}}
	m.state.data.statusFileIndex = 5

	m.setStatusFiles(parseStatusFiles(""))
	m.rebuildStatusContentWithHighlight()

	if len(m.state.data.statusFiles) != 0 {
		t.Fatalf("expected 0 status files for clean tree, got %d", len(m.state.data.statusFiles))
	}
	if m.state.data.statusFileIndex != 0 {
		t.Fatalf("expected statusFileIndex reset to 0, got %d", m.state.data.statusFileIndex)
	}
	if !strings.Contains(m.statusContent, "Clean working tree") {
		t.Fatalf("expected 'Clean working tree' in result, got %q", m.statusContent)
	}
}

// TestRenderStatusFilesHighlighting tests that selected file is highlighted.
func TestRenderStatusFilesHighlighting(t *testing.T) {
	SetIconProvider(&NerdFontV3Provider{})
	cfg := config.DefaultConfig()
	cfg.WorktreeDir = t.TempDir()
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 2
	m.state.ui.statusViewport = viewport.New(viewport.WithWidth(40), viewport.WithHeight(10))
	m.setStatusFiles([]StatusFile{
		{Filename: "file1.go", Status: ".M", IsUntracked: false},
		{Filename: "file2.go", Status: ".M", IsUntracked: false},
	})
	m.state.services.statusTree.Index = 1

	result := m.renderStatusFiles()

	// The result should contain both filenames
	if !strings.Contains(result, "file1.go") {
		t.Fatalf("expected result to contain 'file1.go', got %q", result)
	}
	if !strings.Contains(result, "file2.go") {
		t.Fatalf("expected result to contain 'file2.go', got %q", result)
	}
	icon := deviconForName("file1.go", false)
	if icon == "" {
		t.Fatalf("expected devicon for file1.go, got empty string")
	}
	if !strings.Contains(result, icon) {
		t.Fatalf("expected result to contain devicon %q, got %q", icon, result)
	}

	// Result should have multiple lines
	lines := strings.Split(result, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestRenderStatusFilesIconsDisabled(t *testing.T) {
	SetIconProvider(&NerdFontV3Provider{})
	cfg := config.DefaultConfig()
	cfg.WorktreeDir = t.TempDir()
	cfg.IconSet = "text"
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 2
	m.state.ui.statusViewport = viewport.New(viewport.WithWidth(40), viewport.WithHeight(10))
	m.setStatusFiles([]StatusFile{
		{Filename: "file1.go", Status: ".M", IsUntracked: false},
	})
	m.state.services.statusTree.Index = 0

	result := m.renderStatusFiles()
	icon := deviconForName("file1.go", false)
	if icon != "" && strings.Contains(result, icon) {
		t.Fatalf("expected icons disabled, got %q in %q", icon, result)
	}
}

// TestStatusTreeIndexClamping tests that statusTreeIndex is clamped to valid range.

// TestStatusTreeIndexClamping tests that statusTreeIndex is clamped to valid range.
func TestStatusTreeIndexClamping(t *testing.T) {
	cfg := &config.AppConfig{
		WorktreeDir: t.TempDir(),
	}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 2
	m.state.ui.statusViewport = viewport.New(viewport.WithWidth(40), viewport.WithHeight(10))

	// Set index out of range before parsing
	m.state.services.statusTree.Index = 100

	statusRaw := `1 .M N... 100644 100644 100644 abc123 abc123 file1.go
1 .M N... 100644 100644 100644 abc123 abc123 file2.go`

	m.setStatusFiles(parseStatusFiles(statusRaw))
	m.rebuildStatusContentWithHighlight()

	// Index should be clamped to last valid index
	if m.state.services.statusTree.Index != 1 {
		t.Fatalf("expected statusTreeIndex clamped to 1, got %d", m.state.services.statusTree.Index)
	}

	// Test negative index
	m.state.services.statusTree.Index = -5
	m.setStatusFiles(parseStatusFiles(statusRaw))
	m.rebuildStatusContentWithHighlight()

	if m.state.services.statusTree.Index != 0 {
		t.Fatalf("expected statusTreeIndex clamped to 0, got %d", m.state.services.statusTree.Index)
	}
}

// TestMouseScrollNavigatesFiles tests that mouse scroll navigates tree items in pane 2.

// TestBuildStatusTreeEmpty tests building tree from empty file list.
func TestBuildStatusTreeEmpty(t *testing.T) {
	tree := services.BuildStatusTree([]StatusFile{})
	if tree == nil {
		t.Fatal("expected non-nil tree root")
	} else if tree.Path != "" {
		t.Errorf("expected empty root path, got %q", tree.Path)
	}
	if len(tree.Children) != 0 {
		t.Errorf("expected no children for empty input, got %d", len(tree.Children))
	}
}

// TestBuildStatusTreeFlatFiles tests tree with files at root level.

// TestBuildStatusTreeFlatFiles tests tree with files at root level.
func TestBuildStatusTreeFlatFiles(t *testing.T) {
	files := []StatusFile{
		{Filename: "README.md", Status: ".M"},
		{Filename: "main.go", Status: "M."},
	}
	tree := services.BuildStatusTree(files)

	if len(tree.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(tree.Children))
	}

	// Should be sorted alphabetically
	if tree.Children[0].Path != "README.md" {
		t.Errorf("expected first child README.md, got %q", tree.Children[0].Path)
	}
	if tree.Children[1].Path != "main.go" {
		t.Errorf("expected second child main.go, got %q", tree.Children[1].Path)
	}

	// Both should be files, not directories
	for _, child := range tree.Children {
		if child.IsDir() {
			t.Errorf("expected %q to be a file, not directory", child.Path)
		}
		if child.File == nil {
			t.Errorf("expected %q to have File pointer", child.Path)
		}
	}
}

// TestBuildStatusTreeNestedDirs tests tree with nested directory structure.

// TestBuildStatusTreeNestedDirs tests tree with nested directory structure.
func TestBuildStatusTreeNestedDirs(t *testing.T) {
	files := []StatusFile{
		{Filename: "internal/app/app.go", Status: ".M"},
		{Filename: "internal/app/handlers.go", Status: ".M"},
		{Filename: "internal/git/git.go", Status: "M."},
		{Filename: "README.md", Status: ".M"},
	}
	tree := services.BuildStatusTree(files)

	// Root should have 2 children: internal (dir) and README.md (file)
	// After compression, internal/app and internal/git are separate
	if len(tree.Children) != 2 {
		t.Fatalf("expected 2 root children, got %d", len(tree.Children))
	}

	// Directories should come before files
	if tree.Children[0].Path != "internal" && !strings.HasPrefix(tree.Children[0].Path, "internal") {
		t.Errorf("expected first child to be internal dir, got %q", tree.Children[0].Path)
	}
	if tree.Children[1].Path != "README.md" {
		t.Errorf("expected second child to be README.md, got %q", tree.Children[1].Path)
	}
}

// TestBuildStatusTreeDirsSortedBeforeFiles tests that directories appear before files.

// TestBuildStatusTreeDirsSortedBeforeFiles tests that directories appear before files.
func TestBuildStatusTreeDirsSortedBeforeFiles(t *testing.T) {
	files := []StatusFile{
		{Filename: "zebra.txt", Status: ".M"},
		{Filename: "aaa/file.go", Status: ".M"},
		{Filename: "alpha.txt", Status: ".M"},
	}
	tree := services.BuildStatusTree(files)

	if len(tree.Children) != 3 {
		t.Fatalf("expected 3 children, got %d", len(tree.Children))
	}

	// First should be the directory (aaa), then files alphabetically
	if !tree.Children[0].IsDir() {
		t.Error("expected first child to be a directory")
	}
	if tree.Children[0].Path != "aaa" {
		t.Errorf("expected first child aaa, got %q", tree.Children[0].Path)
	}
	if tree.Children[1].IsDir() {
		t.Error("expected second child to be a file")
	}
	if tree.Children[2].IsDir() {
		t.Error("expected third child to be a file")
	}
}

// TestCompressStatusTreeSingleChild tests compression of single-child directory chains.

// TestCompressStatusTreeSingleChild tests compression of single-child directory chains.
func TestCompressStatusTreeSingleChild(t *testing.T) {
	files := []StatusFile{
		{Filename: "a/b/c/file.go", Status: ".M"},
	}
	tree := services.BuildStatusTree(files)

	// After compression, a/b/c should be one node, not three nested nodes
	flat := services.FlattenStatusTree(tree, map[string]bool{}, 0)

	// Should have: a/b/c (dir) + file.go (file) = 2 nodes
	if len(flat) != 2 {
		t.Fatalf("expected 2 flattened nodes after compression, got %d", len(flat))
	}

	if flat[0].Path != "a/b/c" {
		t.Errorf("expected compressed path a/b/c, got %q", flat[0].Path)
	}
	if !flat[0].IsDir() {
		t.Error("expected first node to be a directory")
	}
	if flat[1].Path != "a/b/c/file.go" {
		t.Errorf("expected file path a/b/c/file.go, got %q", flat[1].Path)
	}
}

// TestFlattenStatusTreeCollapsed tests that collapsed directories hide children.

// TestFlattenStatusTreeCollapsed tests that collapsed directories hide children.
func TestFlattenStatusTreeCollapsed(t *testing.T) {
	files := []StatusFile{
		{Filename: "dir/file1.go", Status: ".M"},
		{Filename: "dir/file2.go", Status: ".M"},
		{Filename: "root.go", Status: ".M"},
	}
	tree := services.BuildStatusTree(files)

	// Without collapse: should see dir + 2 files + root.go = 4 nodes
	flatOpen := services.FlattenStatusTree(tree, map[string]bool{}, 0)
	if len(flatOpen) != 4 {
		t.Fatalf("expected 4 nodes when expanded, got %d", len(flatOpen))
	}

	// With dir collapsed: should see dir + root.go = 2 nodes
	collapsed := map[string]bool{"dir": true}
	flatClosed := services.FlattenStatusTree(tree, collapsed, 0)
	if len(flatClosed) != 2 {
		t.Fatalf("expected 2 nodes when collapsed, got %d", len(flatClosed))
	}

	if flatClosed[0].Path != "dir" {
		t.Errorf("expected first node to be dir, got %q", flatClosed[0].Path)
	}
	if flatClosed[1].Path != "root.go" {
		t.Errorf("expected second node to be root.go, got %q", flatClosed[1].Path)
	}
}

// TestStatusTreeNodeHelpers tests IsDir and Name helper methods.

// TestStatusTreeNodeHelpers tests IsDir and Name helper methods.
func TestStatusTreeNodeHelpers(t *testing.T) {
	fileNode := &StatusTreeNode{
		Path: "internal/app/app.go",
		File: &StatusFile{Filename: "internal/app/app.go", Status: ".M"},
	}
	dirNode := &StatusTreeNode{
		Path:     "internal/app",
		Children: []*StatusTreeNode{},
	}

	if fileNode.IsDir() {
		t.Error("file node should not be a directory")
	}
	if !dirNode.IsDir() {
		t.Error("dir node should be a directory")
	}

	if fileNode.Name() != "app.go" {
		t.Errorf("expected file name app.go, got %q", fileNode.Name())
	}
	if dirNode.Name() != "app" {
		t.Errorf("expected dir name app, got %q", dirNode.Name())
	}
}

// TestFlattenStatusTreeDepth tests that depth is correctly calculated.

// TestFlattenStatusTreeDepth tests that depth is correctly calculated.
func TestFlattenStatusTreeDepth(t *testing.T) {
	files := []StatusFile{
		{Filename: "dir/subdir/file.go", Status: ".M"},
		{Filename: "root.go", Status: ".M"},
	}
	tree := services.BuildStatusTree(files)
	flat := services.FlattenStatusTree(tree, map[string]bool{}, 0)

	// After compression: dir/subdir (depth 0), file.go (depth 1), root.go (depth 0)
	if len(flat) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(flat))
	}

	// Root level nodes should have depth 0
	if flat[0].Depth != 0 {
		t.Errorf("expected dir/subdir depth 0, got %d", flat[0].Depth)
	}
	// File inside dir should have depth 1
	if flat[1].Depth != 1 {
		t.Errorf("expected file.go depth 1, got %d", flat[1].Depth)
	}
	// Root file should have depth 0
	if flat[2].Depth != 0 {
		t.Errorf("expected root.go depth 0, got %d", flat[2].Depth)
	}
}

// TestDirectoryToggleUpdatesFlat tests that toggling directory collapse updates flattened list.

// TestDirectoryToggleUpdatesFlat tests that toggling directory collapse updates flattened list.
func TestDirectoryToggleUpdatesFlat(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	m.state.view.FocusedPane = 2
	m.state.ui.statusViewport = viewport.New(viewport.WithWidth(40), viewport.WithHeight(10))
	m.state.view.WindowWidth = 100
	m.state.view.WindowHeight = 30

	m.setStatusFiles([]StatusFile{
		{Filename: "dir/file1.go", Status: ".M"},
		{Filename: "dir/file2.go", Status: ".M"},
	})

	initialCount := len(m.state.services.statusTree.TreeFlat)
	if initialCount != 3 { // dir + 2 files
		t.Fatalf("expected 3 nodes initially, got %d", initialCount)
	}

	// Collapse the directory
	m.state.services.statusTree.CollapsedDirs["dir"] = true
	m.rebuildStatusTreeFlat()

	if len(m.state.services.statusTree.TreeFlat) != 1 { // just the dir
		t.Fatalf("expected 1 node after collapse, got %d", len(m.state.services.statusTree.TreeFlat))
	}

	// Expand again
	m.state.services.statusTree.CollapsedDirs["dir"] = false
	m.rebuildStatusTreeFlat()

	if len(m.state.services.statusTree.TreeFlat) != 3 {
		t.Fatalf("expected 3 nodes after expand, got %d", len(m.state.services.statusTree.TreeFlat))
	}
}
