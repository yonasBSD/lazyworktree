package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseContainerCommand(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    map[string]any
		wantNil  bool
		wantImg  string
		wantRT   string
		wantWD   string
		wantMnts int
		wantEnv  int
		wantArgs int
	}{
		{
			name:    "missing image returns nil",
			input:   map[string]any{"runtime": "docker"},
			wantNil: true,
		},
		{
			name:    "empty map returns nil",
			input:   map[string]any{},
			wantNil: true,
		},
		{
			name:    "minimal with image only",
			input:   map[string]any{"image": "golang:1.22"},
			wantImg: "golang:1.22",
		},
		{
			name: "full config",
			input: map[string]any{
				"image":       "node:20",
				"runtime":     "podman",
				"working_dir": "/src",
				"mounts": []any{
					map[string]any{"source": "/tmp", "target": "/cache", "read_only": true},
					map[string]any{"source": "/home", "target": "/home"},
				},
				"env":        map[string]any{"NODE_ENV": "test", "CI": "true"},
				"extra_args": []any{"--network=host", "--privileged"},
			},
			wantImg:  "node:20",
			wantRT:   "podman",
			wantWD:   "/src",
			wantMnts: 2,
			wantEnv:  2,
			wantArgs: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := parseContainerCommand(tt.input)
			if tt.wantNil {
				assert.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			assert.Equal(t, tt.wantImg, got.Image)
			assert.Equal(t, tt.wantRT, got.Runtime)
			assert.Equal(t, tt.wantWD, got.WorkingDir)
			assert.Len(t, got.Mounts, tt.wantMnts)
			assert.Len(t, got.Env, tt.wantEnv)
			assert.Len(t, got.ExtraArgs, tt.wantArgs)
		})
	}
}

func TestParseContainerCommandMountReadOnly(t *testing.T) {
	t.Parallel()
	input := map[string]any{
		"image": "alpine",
		"mounts": []any{
			map[string]any{"source": "/tmp", "target": "/cache", "read_only": true},
			map[string]any{"source": "/home", "target": "/home"},
		},
	}
	got := parseContainerCommand(input)
	require.NotNil(t, got)
	require.Len(t, got.Mounts, 2)
	assert.True(t, got.Mounts[0].ReadOnly)
	assert.False(t, got.Mounts[1].ReadOnly)
}

func TestParseContainerCommandMountOptions(t *testing.T) {
	t.Parallel()
	input := map[string]any{
		"image": "alpine",
		"mounts": []any{
			map[string]any{"source": "/tmp/data", "target": "/data", "options": "z"},
			map[string]any{"source": "/tmp/cache", "target": "/cache"},
		},
	}
	got := parseContainerCommand(input)
	require.NotNil(t, got)
	require.Len(t, got.Mounts, 2)
	assert.Equal(t, "z", got.Mounts[0].Options)
	assert.Empty(t, got.Mounts[1].Options)
}

func TestParseContainerCommandArgsAndInteractive(t *testing.T) {
	t.Parallel()
	input := map[string]any{
		"image":       "alpine",
		"args":        []any{"--flag1", "--flag2"},
		"interactive": true,
	}
	got := parseContainerCommand(input)
	require.NotNil(t, got)
	assert.Equal(t, []string{"--flag1", "--flag2"}, got.Args)
	assert.True(t, got.Interactive)
}

func TestParseContainerCommandInteractiveDefaults(t *testing.T) {
	t.Parallel()
	input := map[string]any{
		"image": "alpine",
	}
	got := parseContainerCommand(input)
	require.NotNil(t, got)
	assert.Nil(t, got.Args)
	assert.False(t, got.Interactive)
}

func TestParseCustomCommandsWithContainer(t *testing.T) {
	t.Parallel()
	input := map[string]any{
		"custom_commands": map[string]any{
			PaneUniversal: map[string]any{
				"C": map[string]any{
					"command":     "go test ./...",
					"description": "Run tests in container",
					"show_output": true,
					"container": map[string]any{
						"image": "golang:1.22",
					},
				},
			},
		},
	}
	cmds, _ := parseCustomCommands(input)
	require.Contains(t, cmds[PaneUniversal], "C")
	require.NotNil(t, cmds[PaneUniversal]["C"].Container)
	assert.Equal(t, "golang:1.22", cmds[PaneUniversal]["C"].Container.Image)
	assert.Equal(t, "go test ./...", cmds[PaneUniversal]["C"].Command)
}

func TestParseCustomCommandsContainerOnly(t *testing.T) {
	t.Parallel()
	input := map[string]any{
		"custom_commands": map[string]any{
			PaneUniversal: map[string]any{
				"D": map[string]any{
					"description": "Container shell",
					"container": map[string]any{
						"image": "alpine",
					},
				},
			},
		},
	}
	cmds, _ := parseCustomCommands(input)
	require.Contains(t, cmds[PaneUniversal], "D")
	require.NotNil(t, cmds[PaneUniversal]["D"].Container)
	assert.Equal(t, "alpine", cmds[PaneUniversal]["D"].Container.Image)
}
