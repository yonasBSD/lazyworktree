package app

import (
	"fmt"
	"sort"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	appscreen "github.com/chmouel/lazyworktree/internal/app/screen"
)

type worktreeTagStat struct {
	Tag   string
	Count int
}

func buildWorktreeFilterQueryFromTag(tag string) string {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return ""
	}
	return "tag:" + tag
}

func currentExactTagFilter(query string) string {
	parsed := parseWorktreeFilterQuery(query)
	if len(parsed.textTerms) > 0 || len(parsed.tagTerms) != 1 {
		return ""
	}
	return parsed.tagTerms[0]
}

func (m *Model) worktreeTagStats() []worktreeTagStat {
	counts := make(map[string]int)
	for _, wt := range m.state.data.worktrees {
		note, ok := m.getWorktreeNote(wt.Path)
		if !ok || len(note.Tags) == 0 {
			continue
		}
		for _, tag := range note.Tags {
			counts[tag]++
		}
	}

	stats := make([]worktreeTagStat, 0, len(counts))
	for tag, count := range counts {
		stats = append(stats, worktreeTagStat{Tag: tag, Count: count})
	}

	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Count != stats[j].Count {
			return stats[i].Count > stats[j].Count
		}
		return stats[i].Tag < stats[j].Tag
	})
	return stats
}

func (m *Model) showBrowseWorktreeTags() tea.Cmd {
	stats := m.worktreeTagStats()
	items := make([]appscreen.SelectionItem, len(stats))
	for i, stat := range stats {
		desc := fmt.Sprintf("%d worktrees", stat.Count)
		if stat.Count == 1 {
			desc = "1 worktree"
		}
		items[i] = appscreen.SelectionItem{
			ID:          stat.Tag,
			Label:       stat.Tag,
			Description: desc,
		}
	}

	scr := appscreen.NewListSelectionScreen(
		items,
		"Browse worktree tags",
		"Filter tags...",
		"No worktree tags found.",
		m.state.view.WindowWidth,
		m.state.view.WindowHeight,
		currentExactTagFilter(m.state.services.filter.FilterQuery),
		m.theme,
	)
	scr.FooterHint = "Enter applies exact tag filter"
	scr.OnSelect = func(item appscreen.SelectionItem) tea.Cmd {
		query := buildWorktreeFilterQueryFromTag(item.ID)
		m.state.view.ShowingFilter = true
		m.state.view.ShowingSearch = false
		m.setFilterTarget(filterTargetWorktrees)
		m.setFilterQuery(filterTargetWorktrees, query)
		m.state.ui.filterInput.SetValue(query)
		m.state.ui.filterInput.CursorEnd()
		m.state.ui.filterInput.Focus()
		m.updateTable()
		return textinput.Blink
	}
	scr.OnCancel = func() tea.Cmd {
		return nil
	}

	m.state.ui.screenManager.Push(scr)
	return textinput.Blink
}
