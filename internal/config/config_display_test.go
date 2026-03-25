package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyntaxThemeForUITheme(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		inputTheme string
		want       string
	}{
		{name: "default dracula", inputTheme: "dracula", want: "Dracula"},
		{name: "dracula-light", inputTheme: "dracula-light", want: "\"Monokai Extended Light\""},
		{name: "kanagawa", inputTheme: "kanagawa", want: "Dracula"},
		{name: "narna", inputTheme: "narna", want: "\"OneHalfDark\""},
		{name: "clean-light", inputTheme: "clean-light", want: "GitHub"},
		{name: "solarized-dark", inputTheme: "solarized-dark", want: "\"Solarized (dark)\""},
		{name: "solarized-light", inputTheme: "solarized-light", want: "\"Solarized (light)\""},
		{name: "gruvbox-dark", inputTheme: "gruvbox-dark", want: "\"Gruvbox Dark\""},
		{name: "gruvbox-light", inputTheme: "gruvbox-light", want: "\"Gruvbox Light\""},
		{name: "nord", inputTheme: "nord", want: "\"Nord\""},
		{name: "monokai", inputTheme: "monokai", want: "\"Monokai Extended\""},
		{name: "catppuccin-mocha", inputTheme: "catppuccin-mocha", want: "\"Catppuccin Mocha\""},
		{name: "catppuccin-latte", inputTheme: "catppuccin-latte", want: "\"Catppuccin Latte\""},
		{name: "rose-pine-dawn", inputTheme: "rose-pine-dawn", want: "GitHub"},
		{name: "one-light", inputTheme: "one-light", want: "\"OneHalfLight\""},
		{name: "everforest-light", inputTheme: "everforest-light", want: "\"Gruvbox Light\""},
		{name: "unknown falls back to dracula", inputTheme: "unknown", want: "Dracula"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, SyntaxThemeForUITheme(tt.inputTheme))
		})
	}
}

func TestNormalizeThemeName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "lowercase", input: "dracula", want: "dracula"},
		{name: "kanagawa", input: "kanagawa", want: "kanagawa"},
		{name: "uppercase", input: "NARNA", want: "narna"},
		{name: "trimmed whitespace", input: "  clean-light  ", want: "clean-light"},
		{name: "unknown theme", input: "invalid", want: ""},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NormalizeThemeName(tt.input))
		})
	}
}

func TestParseLayoutSizes(t *testing.T) {
	t.Parallel()
	data := map[string]any{
		"layout_sizes": map[string]any{
			"worktrees":      45,
			"info":           30,
			"git_status":     40,
			"commit":         30,
			"notes":          25,
			"agent_sessions": 20,
		},
	}
	cfg, err := parseConfig(data)
	require.NoError(t, err)
	require.NotNil(t, cfg.LayoutSizes)
	assert.Equal(t, 45, cfg.LayoutSizes.Worktrees)
	assert.Equal(t, 30, cfg.LayoutSizes.Info)
	assert.Equal(t, 40, cfg.LayoutSizes.GitStatus)
	assert.Equal(t, 30, cfg.LayoutSizes.Commit)
	assert.Equal(t, 25, cfg.LayoutSizes.Notes)
	assert.Equal(t, 20, cfg.LayoutSizes.AgentSessions)
}

func TestParseLayoutSizesPartial(t *testing.T) {
	t.Parallel()
	data := map[string]any{
		"layout_sizes": map[string]any{
			"worktrees": 60,
			"info":      40,
		},
	}
	cfg, err := parseConfig(data)
	require.NoError(t, err)
	require.NotNil(t, cfg.LayoutSizes)
	assert.Equal(t, 60, cfg.LayoutSizes.Worktrees)
	assert.Equal(t, 40, cfg.LayoutSizes.Info)
	assert.Equal(t, 0, cfg.LayoutSizes.GitStatus)
	assert.Equal(t, 0, cfg.LayoutSizes.Commit)
	assert.Equal(t, 0, cfg.LayoutSizes.Notes)
	assert.Equal(t, 0, cfg.LayoutSizes.AgentSessions)
}

func TestParseLayoutSizesValidation(t *testing.T) {
	t.Parallel()
	data := map[string]any{
		"layout_sizes": map[string]any{
			"worktrees": 150,
			"info":      -5,
			"commit":    0,
			"notes":     "50",
		},
	}
	cfg, err := parseConfig(data)
	require.NoError(t, err)
	require.NotNil(t, cfg.LayoutSizes)
	assert.Equal(t, 100, cfg.LayoutSizes.Worktrees, "values > 100 should be clamped to 100")
	assert.Equal(t, 0, cfg.LayoutSizes.Info, "negative values should stay at 0 (unset)")
	assert.Equal(t, 0, cfg.LayoutSizes.Commit, "zero should stay at 0 (unset)")
	assert.Equal(t, 50, cfg.LayoutSizes.Notes, "string coercion should work")
}

func TestParseLayoutSizesNil(t *testing.T) {
	t.Parallel()
	data := map[string]any{}
	cfg, err := parseConfig(data)
	require.NoError(t, err)
	assert.Nil(t, cfg.LayoutSizes, "nil when layout_sizes not present")
}

func TestParseLayoutSizesOverride(t *testing.T) {
	t.Parallel()
	cfg := DefaultConfig()
	cfg.LayoutSizes = &LayoutSizes{
		Worktrees: 55,
		Info:      30,
		GitStatus: 40,
		Commit:    30,
	}
	err := cfg.ApplyCLIOverrides([]string{"lw.layout_sizes.worktrees=70", "lw.layout_sizes.notes=20"})
	require.NoError(t, err)
	assert.Equal(t, 70, cfg.LayoutSizes.Worktrees)
	assert.Equal(t, 30, cfg.LayoutSizes.Info, "unchanged fields preserved")
	assert.Equal(t, 20, cfg.LayoutSizes.Notes)
}
