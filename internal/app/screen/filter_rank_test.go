package screen

import "testing"

func TestRankedFieldMatchScorePrefersEarlierFields(t *testing.T) {
	scoreLabel, ok := rankedFieldMatchScore("browse", "Browse files", "Open in browser")
	if !ok {
		t.Fatal("expected label match to score")
	}

	scoreDesc, ok := rankedFieldMatchScore("browse", "Open item", "Open in browser")
	if !ok {
		t.Fatal("expected description match to score")
	}

	if scoreLabel >= scoreDesc {
		t.Fatalf("expected label match score %d to beat description score %d", scoreLabel, scoreDesc)
	}
}

func TestRankedFieldMatchScorePrefersWordPrefixOverSubstring(t *testing.T) {
	scorePrefix, ok := rankedFieldMatchScore("browse", "Browse files")
	if !ok {
		t.Fatal("expected prefix match to score")
	}

	scoreSubstring, ok := rankedFieldMatchScore("browse", "Open browser page")
	if !ok {
		t.Fatal("expected substring match to score")
	}

	if scorePrefix >= scoreSubstring {
		t.Fatalf("expected prefix score %d to beat substring score %d", scorePrefix, scoreSubstring)
	}
}
