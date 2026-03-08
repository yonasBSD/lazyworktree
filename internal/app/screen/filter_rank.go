package screen

import (
	"sort"
	"strings"
)

const (
	noRankedFilterScore = int(^uint(0) >> 1)
	// fieldPenalty is applied per field index so earlier fields rank higher.
	fieldPenalty = 1000
)

type rankedSelection[T any] struct {
	item  T
	score int
	order int
}

func rankedFieldMatchScore(query string, fields ...string) (int, bool) {
	if query == "" {
		return 0, true
	}

	best := noRankedFilterScore
	for i, field := range fields {
		score, ok := rankedTextMatchScore(query, strings.ToLower(field))
		if !ok {
			continue
		}
		score += i * fieldPenalty
		if score < best {
			best = score
		}
	}

	if best == noRankedFilterScore {
		return 0, false
	}
	return best, true
}

// filterAndRank scores each item against query, sorts by score, and returns
// the filtered slice. The query must already be lowercased and trimmed.
func filterAndRank[T any](items []T, query string, fields func(T) []string) []T {
	scored := make([]rankedSelection[T], 0, len(items))
	for i := range items {
		score, ok := rankedFieldMatchScore(query, fields(items[i])...)
		if !ok {
			continue
		}
		scored = append(scored, rankedSelection[T]{
			item:  items[i],
			score: score,
			order: i,
		})
	}
	sortRankedSelections(scored)

	result := make([]T, len(scored))
	for i := range scored {
		result[i] = scored[i].item
	}
	return result
}

func rankedTextMatchScore(query, text string) (int, bool) {
	if query == "" || text == "" {
		return 0, false
	}

	if start := findExactWordMatch(text, query); start >= 0 {
		return start, true
	}
	if start := findWordPrefixMatch(text, query); start >= 0 {
		return 100 + start, true
	}
	if start := strings.Index(text, query); start >= 0 {
		return 200 + start, true
	}

	return 0, false
}

func sortRankedSelections[T any](items []rankedSelection[T]) {
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].score != items[j].score {
			return items[i].score < items[j].score
		}
		return items[i].order < items[j].order
	})
}
