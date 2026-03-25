package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTmuxCommandDefaults(t *testing.T) {
	t.Parallel()

	got := parseTmuxCommand(map[string]any{
		"session_name": "${REPO_NAME}_wt_$WORKTREE_NAME",
	})

	require.NotNil(t, got)
	assert.Equal(t, "${REPO_NAME}_wt_$WORKTREE_NAME", got.SessionName)
	assert.True(t, got.Attach)
	assert.Equal(t, "switch", got.OnExists)
	require.Len(t, got.Windows, 1)
	assert.Equal(t, "shell", got.Windows[0].Name)
}

func TestParseTmuxCommandWindows(t *testing.T) {
	t.Parallel()

	got := parseTmuxCommand(map[string]any{
		"session_name": "${REPO_NAME}_wt_$WORKTREE_NAME",
		"attach":       false,
		"on_exists":    "kill",
		"windows": []any{
			map[string]any{
				"name":    "shell",
				"command": "zsh",
				"cwd":     "$WORKTREE_PATH",
			},
			map[string]any{
				"name":    "git",
				"command": "lazygit",
			},
		},
	})

	require.NotNil(t, got)
	assert.False(t, got.Attach)
	assert.Equal(t, "kill", got.OnExists)
	require.Len(t, got.Windows, 2)
	assert.Equal(t, TmuxWindow{Name: "shell", Command: "zsh", Cwd: "$WORKTREE_PATH"}, got.Windows[0])
	assert.Equal(t, TmuxWindow{Name: "git", Command: "lazygit", Cwd: ""}, got.Windows[1])
}
