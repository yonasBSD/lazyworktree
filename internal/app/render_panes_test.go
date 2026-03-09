package app

import (
	"strings"
	"testing"

	"github.com/chmouel/lazyworktree/internal/config"
	"github.com/chmouel/lazyworktree/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAggregateCIConclusion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		checks   []*models.CICheck
		expected string
	}{
		{
			name:     "all success",
			checks:   []*models.CICheck{{Conclusion: "success"}, {Conclusion: "success"}},
			expected: "success",
		},
		{
			name:     "failure takes priority",
			checks:   []*models.CICheck{{Conclusion: "success"}, {Conclusion: "failure"}, {Conclusion: "pending"}},
			expected: "failure",
		},
		{
			name:     "pending over success",
			checks:   []*models.CICheck{{Conclusion: "success"}, {Conclusion: "pending"}},
			expected: "pending",
		},
		{
			name:     "empty conclusion treated as pending",
			checks:   []*models.CICheck{{Conclusion: "success"}, {Conclusion: ""}},
			expected: "pending",
		},
		{
			name:     "all skipped",
			checks:   []*models.CICheck{{Conclusion: "skipped"}, {Conclusion: "cancelled"}},
			expected: "skipped",
		},
		{
			name:     "single failure",
			checks:   []*models.CICheck{{Conclusion: "failure"}},
			expected: "failure",
		},
		{
			name:     "single success",
			checks:   []*models.CICheck{{Conclusion: "success"}},
			expected: "success",
		},
		{
			name:     "skipped and success",
			checks:   []*models.CICheck{{Conclusion: "skipped"}, {Conclusion: "success"}},
			expected: "success",
		},
		{
			name:     "cancelled and pending",
			checks:   []*models.CICheck{{Conclusion: "cancelled"}, {Conclusion: ""}},
			expected: "pending",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := aggregateCIConclusion(tt.checks)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRenderNotesBoxWrapsStyledContentWithoutLeakingANSITails(t *testing.T) {
	cfg := &config.AppConfig{WorktreeDir: t.TempDir()}
	m := NewModel(cfg, "")
	wt := &models.WorktreeInfo{
		Path:   "/tmp/wt-notes-wrap",
		Branch: "feature/notes-wrap",
	}

	m.state.data.worktrees = []*models.WorktreeInfo{wt}
	m.state.data.filteredWts = []*models.WorktreeInfo{wt}
	m.state.data.selectedIndex = 0
	m.worktreeNotes[worktreeNoteKey(wt.Path)] = models.WorktreeNote{
		Note: "- Centralized GitHub App JWT generation by moving logic from app package to main provider.\n- Updated all call sites to utilize the new centralized provider location.\n- Linked Jira: [SRVKP-10952].",
	}
	m.notesContent = m.buildNotesContent(wt)

	rendered := m.renderNotesBox(36, 7)
	plain := stripTerminalSequences(rendered)

	if strings.Contains(plain, "38;2;") {
		t.Fatalf("expected wrapped notes to hide raw ANSI fragments, got %q", plain)
	}
	if !strings.Contains(plain, "Linked Jira: [SRVKP-10952].") {
		t.Fatalf("expected wrapped note content to stay visible, got %q", plain)
	}
}
