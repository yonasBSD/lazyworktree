package config

import "strings"

// TmuxCommand represents a configured tmux session layout.
type TmuxCommand struct {
	SessionName string
	Attach      bool
	OnExists    string
	Windows     []TmuxWindow
}

// TmuxWindow represents a tmux window configuration.
type TmuxWindow struct {
	Name    string
	Command string
	Cwd     string
}

func parseTmuxCommand(data map[string]any) *TmuxCommand {
	cmd := &TmuxCommand{
		SessionName: getString(data, "session_name"),
		Attach:      coerceBool(data["attach"], true),
		OnExists:    strings.ToLower(getString(data, "on_exists")),
	}
	if cmd.OnExists == "" {
		cmd.OnExists = "switch"
	}

	if windows, ok := data["windows"].([]any); ok {
		for _, w := range windows {
			if wData, ok := w.(map[string]any); ok {
				cmd.Windows = append(cmd.Windows, TmuxWindow{
					Name:    getString(wData, "name"),
					Command: getString(wData, "command"),
					Cwd:     getString(wData, "cwd"),
				})
			}
		}
	}
	if len(cmd.Windows) == 0 {
		cmd.Windows = []TmuxWindow{
			{
				Name:    "shell",
				Command: "",
				Cwd:     "",
			},
		}
	}
	return cmd
}
