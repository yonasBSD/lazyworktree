package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseKeybindings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    map[string]any
		wantLen  int
		wantKeys map[string]map[string]string
	}{
		{
			name: "valid pane-scoped input",
			input: map[string]any{
				"keybindings": map[string]any{
					"universal": map[string]any{
						"G": "git-lazygit",
						"F": "git-fetch",
					},
					"worktrees": map[string]any{
						"x": "worktree-delete",
					},
				},
			},
			wantLen: 2,
			wantKeys: map[string]map[string]string{
				"universal": {"G": "git-lazygit", "F": "git-fetch"},
				"worktrees": {"x": "worktree-delete"},
			},
		},
		{
			name: "whitespace trimming",
			input: map[string]any{
				"keybindings": map[string]any{
					"universal": map[string]any{
						"  G  ": "  git-lazygit  ",
					},
				},
			},
			wantKeys: map[string]map[string]string{
				"universal": {"G": "git-lazygit"},
			},
		},
		{
			name: "flat map returns empty",
			input: map[string]any{
				"keybindings": map[string]any{
					"G": "git-lazygit",
				},
			},
			wantLen: 0,
		},
		{
			name:    "empty input",
			input:   map[string]any{},
			wantLen: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := parseKeybindings(tc.input)
			if tc.wantKeys != nil {
				for pane, bindings := range tc.wantKeys {
					for key, actionID := range bindings {
						assert.Equal(t, actionID, result[pane][key], "pane=%s key=%s", pane, key)
					}
				}
			}
			if tc.wantLen > 0 {
				assert.Len(t, result, tc.wantLen)
			} else if tc.wantKeys == nil {
				assert.Empty(t, result)
			}
		})
	}
}

func TestKeybindingsConfigAllForPane(t *testing.T) {
	t.Parallel()

	kb := KeybindingsConfig{
		PaneUniversal: {"G": "git-lazygit", "F": "git-fetch"},
		PaneWorktrees: {"G": "worktree-delete", "x": "worktree-prune"},
	}

	t.Run("universal only", func(t *testing.T) {
		t.Parallel()
		result := kb.AllForPane(PaneUniversal)
		assert.Equal(t, "git-lazygit", result["G"])
		assert.Equal(t, "git-fetch", result["F"])
		assert.Empty(t, result["x"])
	})

	t.Run("pane overrides universal", func(t *testing.T) {
		t.Parallel()
		result := kb.AllForPane(PaneWorktrees)
		assert.Equal(t, "worktree-delete", result["G"], "pane-specific should override universal")
		assert.Equal(t, "git-fetch", result["F"], "universal should be inherited")
		assert.Equal(t, "worktree-prune", result["x"], "pane-specific key should be present")
	})

	t.Run("pane with no specific bindings", func(t *testing.T) {
		t.Parallel()
		result := kb.AllForPane(PaneStatus)
		assert.Equal(t, "git-lazygit", result["G"])
		assert.Equal(t, "git-fetch", result["F"])
	})
}
