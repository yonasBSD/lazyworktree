package app

import (
	tea "charm.land/bubbletea/v2"

	appscreen "github.com/chmouel/lazyworktree/internal/app/screen"
)

var curatedIcons = []appscreen.SelectionItem{
	{ID: "", Label: " Default Folder", Description: "nerd-tree"},

	// Version Control / Git
	{ID: "󰲋", Label: "󰲋 Git Branch", Description: "nerd-tree"},
	{ID: "󰊢", Label: "󰊢 Git", Description: "nerd-tree"},
	{ID: "", Label: " Repo", Description: "nerd-tree"},
	{ID: "", Label: " Directory", Description: "nerd-tree"},
	{ID: "󰉖", Label: "󰉖 Folder Open", Description: "nerd-tree"},
	{ID: "", Label: " Git Branch Outline", Description: "nerd-tree"},
	{ID: "", Label: " Git Commit", Description: "nerd-tree"},
	{ID: "", Label: " Git Merge", Description: "nerd-tree"},
	{ID: "", Label: " Github", Description: "nerd-tree"},
	{ID: "", Label: " Gitlab", Description: "nerd-tree"},

	// Development / Languages
	{ID: "󰅪", Label: "󰅪 Code", Description: "nerd-tree"},
	{ID: "󰌠", Label: "󰌠 Python", Description: "nerd-tree"},
	{ID: "", Label: " Go", Description: "nerd-tree"},
	{ID: "", Label: " JavaScript", Description: "nerd-tree"},
	{ID: "", Label: " TypeScript", Description: "nerd-tree"},
	{ID: "", Label: " Rust", Description: "nerd-tree"},
	{ID: "", Label: " Java", Description: "nerd-tree"},
	{ID: "", Label: " C++", Description: "nerd-tree"},
	{ID: "", Label: " PHP", Description: "nerd-tree"},
	{ID: "", Label: " React", Description: "nerd-tree"},
	{ID: "󰡄", Label: "󰡄 Vue", Description: "nerd-tree"},
	{ID: "󰎙", Label: "󰎙 Node.js", Description: "nerd-tree"},

	// Infrastructure / DevOps
	{ID: "", Label: " Docker", Description: "nerd-tree"},
	{ID: "󱃾", Label: "󱃾 Kubernetes", Description: "nerd-tree"},
	{ID: "", Label: " Amazon Web Services", Description: "nerd-tree"},
	{ID: "󱇶", Label: "󱇶 Google Cloud", Description: "nerd-tree"},
	{ID: "󰊫", Label: "󰊫 Azure", Description: "nerd-tree"},
	{ID: "", Label: " Linux", Description: "nerd-tree"},
	{ID: "", Label: " Apple", Description: "nerd-tree"},
	{ID: "", Label: " Windows", Description: "nerd-tree"},
	{ID: "󰒋", Label: "󰒋 Server", Description: "nerd-tree"},
	{ID: "", Label: " Database", Description: "nerd-tree"},

	// State / Status
	{ID: "", Label: " Todo", Description: "nerd-tree"},
	{ID: "󰄵", Label: "󰄵 Done", Description: "nerd-tree"},
	{ID: "󰅖", Label: "󰅖 Cancelled / Closed", Description: "nerd-tree"},
	{ID: "󰏤", Label: "󰏤 Paused", Description: "nerd-tree"},
	{ID: "󰑐", Label: "󰑐 Working / In Progress", Description: "nerd-tree"},
	{ID: "󰲡", Label: "󰲡 Not Working / Broken", Description: "nerd-tree"},
	{ID: "󰥔", Label: "󰥔 Someday / Later", Description: "nerd-tree"},
	{ID: "󰒲", Label: "󰒲 Sleeping / Waiting", Description: "nerd-tree"},
	{ID: "󰀪", Label: "󰀪 Warning / Blocked", Description: "nerd-tree"},
	{ID: "󰗖", Label: "󰗖 Success", Description: "nerd-tree"},

	// Concepts / Actions
	{ID: "", Label: " Bug Outline", Description: "nerd-tree"},
	{ID: "󰜎", Label: "󰜎 Feature", Description: "nerd-tree"},
	{ID: "󰈙", Label: "󰈙 File / Document", Description: "nerd-tree"},
	{ID: "󰱦", Label: "󰱦 Tool", Description: "nerd-tree"},
	{ID: "󰏗", Label: "󰏗 Package Outline", Description: "nerd-tree"},
	{ID: "󰒓", Label: "󰒓 Settings / Config", Description: "nerd-tree"},
	{ID: "󰖟", Label: "󰖟 Globe / Web", Description: "nerd-tree"},
	{ID: "󰅧", Label: "󰅧 Cloud Outline", Description: "nerd-tree"},
	{ID: "󰠮", Label: "󰠮 Book / Manual", Description: "nerd-tree"},
	{ID: "󰙏", Label: "󰙏 Clock / Performance", Description: "nerd-tree"},
	{ID: "󰄲", Label: "󰄲 Checkbox / Task", Description: "nerd-tree"},
	{ID: "󰕥", Label: "󰕥 Shield / Security", Description: "nerd-tree"},
	{ID: "󰙎", Label: "󰙎 Link / API", Description: "nerd-tree"},
	{ID: "󰗡", Label: "󰗡 Bot / AI", Description: "nerd-tree"},

	// Emojis (State / Status)
	{ID: "✅", Label: "✅ Done", Description: "emoji"},
	{ID: "❌", Label: "❌ Cancelled / Failed", Description: "emoji"},
	{ID: "⏸️", Label: "⏸️ Paused", Description: "emoji"},
	{ID: "⏳", Label: "⏳ Working / Waiting", Description: "emoji"},
	{ID: "🛑", Label: "🛑 Stopped / Blocked", Description: "emoji"},
	{ID: "⚠️", Label: "⚠️ Warning", Description: "emoji"},
	{ID: "🎉", Label: "🎉 Success / Celebration", Description: "emoji"},
	{ID: "💡", Label: "💡 Idea / Todo", Description: "emoji"},

	// Emojis (General)
	{ID: "🚀", Label: "🚀 Rocket", Description: "emoji"},
	{ID: "💻", Label: "💻 Laptop", Description: "emoji"},
	{ID: "🔥", Label: "🔥 Fire", Description: "emoji"},
	{ID: "🐛", Label: "🐛 Bug", Description: "emoji"},
	{ID: "🌟", Label: "🌟 Star", Description: "emoji"},
	{ID: "⚡", Label: "⚡ Zap", Description: "emoji"},
	{ID: "📦", Label: "📦 Package", Description: "emoji"},
	{ID: "🛠", Label: "🛠 Tools", Description: "emoji"},
	{ID: "🚧", Label: "🚧 Construction", Description: "emoji"},
	{ID: "🎨", Label: "🎨 Palette", Description: "emoji"},
	{ID: "✨", Label: "✨ Sparkles", Description: "emoji"},
	{ID: "📚", Label: "📚 Documentation", Description: "emoji"},
	{ID: "🌐", Label: "🌐 Web", Description: "emoji"},
	{ID: "📱", Label: "📱 Mobile", Description: "emoji"},
}

// showSetWorktreeIcon shows a picker to select a custom icon for the current worktree.
func (m *Model) showSetWorktreeIcon() tea.Cmd {
	wt := m.selectedWorktree()
	if wt == nil {
		return nil
	}

	initialID := ""
	if note, ok := m.getWorktreeNote(wt.Path); ok && note.Icon != "" {
		initialID = note.Icon
	}

	scr := appscreen.NewListSelectionScreen(
		curatedIcons,
		"Set worktree icon",
		"Filter icons...",
		"No matching icons.",
		m.state.view.WindowWidth,
		m.state.view.WindowHeight,
		initialID,
		m.theme,
	)

	scr.OnSelect = func(item appscreen.SelectionItem) tea.Cmd {
		icon := item.ID
		// Special case: reset to default
		if icon == "" {
			icon = ""
		}
		m.setWorktreeIcon(wt.Path, icon)
		return nil // UI update happens via setWorktreeIcon which calls refreshSelectedWorktreeNotesPane
	}
	scr.OnCancel = func() tea.Cmd {
		return nil
	}

	m.state.ui.screenManager.Push(scr)
	return nil
}
