package config

import "strings"

// Pane name constants used for keybinding and command dispatch.
const (
	PaneWorktrees     = "worktrees"
	PaneInfo          = "info"
	PaneStatus        = "status"
	PaneLog           = "log"
	PaneNotes         = "notes"
	PaneAgentSessions = "agent_sessions"
	PaneUniversal     = "universal"
)

// KeybindingsConfig maps pane names to key→actionID maps. "universal" applies to all panes.
type KeybindingsConfig map[string]map[string]string

// AllForPane returns a merged key→actionID map: universal bindings first,
// then pane-specific bindings override. Mirrors CustomCommandsConfig.AllForPane.
func (k KeybindingsConfig) AllForPane(paneName string) map[string]string {
	result := make(map[string]string)
	for key, id := range k[PaneUniversal] {
		result[key] = id
	}
	if paneName != PaneUniversal {
		for key, id := range k[paneName] {
			result[key] = id
		}
	}
	return result
}

// Lookup returns the action ID for the given pane and key,
// checking pane-specific first, then universal.
func (k KeybindingsConfig) Lookup(paneName, key string) (string, bool) {
	if m, ok := k[paneName]; ok {
		if id, ok := m[key]; ok {
			return id, true
		}
	}
	if m, ok := k[PaneUniversal]; ok {
		if id, ok := m[key]; ok {
			return id, true
		}
	}
	return "", false
}

func parseKeybindings(data map[string]any) KeybindingsConfig {
	raw, ok := data["keybindings"].(map[string]any)
	if !ok {
		return make(KeybindingsConfig)
	}
	result := make(KeybindingsConfig)
	for pane, val := range raw {
		paneMap, ok := val.(map[string]any)
		if !ok {
			continue
		}
		bindings := make(map[string]string)
		for key, actionVal := range paneMap {
			if actionID, ok := actionVal.(string); ok {
				actionID = strings.TrimSpace(actionID)
				if actionID != "" {
					bindings[strings.TrimSpace(key)] = actionID
				}
			}
		}
		if len(bindings) > 0 {
			result[pane] = bindings
		}
	}
	return result
}
