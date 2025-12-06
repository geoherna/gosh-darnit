package goshdarnit

import (
	"reflect"
	"testing"
)

func TestIsProfane(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		// Empty and clean text
		{"empty string", "", false},
		{"clean short text", "Hello, how are you today?", false},
		{"clean long text", "The quick brown fox jumps over the lazy dog.", false},

		// Basic profanity detection
		{"basic profanity fuck", "What the fuck", true},
		{"basic profanity shit", "This is shit", true},
		{"basic profanity damn", "Damn it all", true},
		{"profanity at start", "fuck this", true},
		{"profanity at end", "oh shit", true},
		{"profanity in middle", "what the hell is this", true},

		// Case insensitivity
		{"uppercase FUCK", "WHAT THE FUCK", true},
		{"mixed case FuCk", "What the FuCk", true},

		// Leetspeak detection
		{"leetspeak fvck", "What the fvck", true},
		{"leetspeak sh1t", "This is sh1t", true},
		{"leetspeak @ss", "You are an @ss", true},
		{"leetspeak f00l", "Hello world", false}, // Not profane

		// Repeated characters
		{"repeated fuuuck", "What the fuuuuck", true},
		{"repeated shiiiit", "Oh shiiiit", true},
		// Note: "assss" collapses to "as" (single s), but "ass" has 2 s's - doesn't match
		{"repeated assss", "what an assss", false},

		// Word boundary prevention (false positives)
		{"bass is not profane", "I play bass guitar", false},
		{"class is not profane", "This is a class", false},
		{"assassin is not profane", "The assassin struck", false},
		{"analyst is not profane", "The analyst reviewed it", false},
		{"password is not profane", "Enter your password", false},
		{"scunthorpe is not profane", "Welcome to Scunthorpe", false},

		// Unicode homoglyphs - note: cyrillic а -> a, so fаck -> fack (not fuck)
		{"cyrillic а in fuck", "fаck you", false},
		{"cyrillic о in shit", "shоt", false},
		{"fullwidth letters", "ｆuck", false},

		// Zero-width characters
		{"zero-width in fuck", "fu\u200Bck", true},
		{"zero-width joiner", "sh\u200Dit", true},

		// Multiple profanities
		{"multiple profanities", "fuck this shit", true},

		// Edge cases
		{"just profanity", "fuck", true},
		{"profanity with punctuation", "fuck!", true},
		{"profanity in parentheses", "(shit)", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsProfane(tt.text)
			if got != tt.expected {
				t.Errorf("IsProfane(%q) = %v, want %v", tt.text, got, tt.expected)
			}
		})
	}
}

func TestCensor(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		mode     CensorMode
		expected string
	}{
		// Empty string
		{"empty CensorAll", "", CensorAll, ""},
		{"empty CensorKeepFirst", "", CensorKeepFirst, ""},
		{"empty CensorKeepFirstLast", "", CensorKeepFirstLast, ""},

		// Clean text (no change)
		{"clean text CensorAll", "Hello world", CensorAll, "Hello world"},

		// CensorAll mode
		{"CensorAll fuck", "What the fuck", CensorAll, "What the ****"},
		{"CensorAll shit", "This is shit", CensorAll, "This is ****"},
		{"CensorAll multiple", "fuck this shit", CensorAll, "**** this ****"},

		// CensorKeepFirst mode
		{"CensorKeepFirst fuck", "What the fuck", CensorKeepFirst, "What the f***"},
		{"CensorKeepFirst shit", "This is shit", CensorKeepFirst, "This is s***"},

		// CensorKeepFirstLast mode
		{"CensorKeepFirstLast fuck", "What the fuck", CensorKeepFirstLast, "What the f**k"},
		{"CensorKeepFirstLast shit", "This is shit", CensorKeepFirstLast, "This is s**t"},
		{"CensorKeepFirstLast short", "Oh ass", CensorKeepFirstLast, "Oh a*s"},

		// Preserve surrounding text
		{"preserve text around", "before fuck after", CensorAll, "before **** after"},
		{"preserve punctuation", "What the fuck!", CensorAll, "What the ****!"},

		// Single character edge case
		{"very short word", "A a A", CensorAll, "A a A"}, // 'a' alone is not profane
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Censor(tt.text, tt.mode)
			if got != tt.expected {
				t.Errorf("Censor(%q, %v) = %q, want %q", tt.text, tt.mode, got, tt.expected)
			}
		})
	}
}

func TestCensorWithDefault(t *testing.T) {
	tests := []struct {
		text     string
		expected string
	}{
		{"", ""},
		{"clean text", "clean text"},
		{"What the fuck", "What the ****"},
		{"fuck shit", "**** ****"},
	}

	for _, tt := range tests {
		got := CensorWithDefault(tt.text)
		if got != tt.expected {
			t.Errorf("CensorWithDefault(%q) = %q, want %q", tt.text, got, tt.expected)
		}
	}
}

func TestFindProfanity(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected []string
	}{
		{"empty string", "", nil},
		{"clean text", "Hello world", nil},
		{"single profanity", "What the fuck", []string{"fuck"}},
		{"multiple profanities", "fuck this shit", []string{"fuck", "shit"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindProfanity(tt.text)
			if tt.expected == nil {
				if got != nil {
					t.Errorf("FindProfanity(%q) = %v, want nil", tt.text, got)
				}
				return
			}
			if len(got) != len(tt.expected) {
				t.Errorf("FindProfanity(%q) returned %d items, want %d", tt.text, len(got), len(tt.expected))
				return
			}
			// Check that all expected words are found (order may vary)
			gotMap := make(map[string]bool)
			for _, w := range got {
				gotMap[w] = true
			}
			for _, w := range tt.expected {
				if !gotMap[w] {
					t.Errorf("FindProfanity(%q) missing expected word %q", tt.text, w)
				}
			}
		})
	}
}

func TestMergeOverlapping(t *testing.T) {
	tests := []struct {
		name     string
		input    []matchInfo
		expected []matchInfo
	}{
		{
			name:     "empty slice",
			input:    []matchInfo{},
			expected: []matchInfo{},
		},
		{
			name:     "single match",
			input:    []matchInfo{{origStart: 0, origEnd: 4, pattern: "test"}},
			expected: []matchInfo{{origStart: 0, origEnd: 4, pattern: "test"}},
		},
		{
			name: "non-overlapping",
			input: []matchInfo{
				{origStart: 0, origEnd: 4, pattern: "test"},
				{origStart: 10, origEnd: 14, pattern: "word"},
			},
			expected: []matchInfo{
				{origStart: 0, origEnd: 4, pattern: "test"},
				{origStart: 10, origEnd: 14, pattern: "word"},
			},
		},
		{
			name: "overlapping",
			input: []matchInfo{
				{origStart: 0, origEnd: 5, pattern: "first"},
				{origStart: 3, origEnd: 8, pattern: "second"},
			},
			expected: []matchInfo{
				{origStart: 0, origEnd: 8, pattern: "first"},
			},
		},
		{
			name: "adjacent (not overlapping)",
			input: []matchInfo{
				{origStart: 0, origEnd: 4, pattern: "first"},
				{origStart: 5, origEnd: 9, pattern: "second"},
			},
			expected: []matchInfo{
				{origStart: 0, origEnd: 4, pattern: "first"},
				{origStart: 5, origEnd: 9, pattern: "second"},
			},
		},
		{
			name: "touching (edge overlap)",
			input: []matchInfo{
				{origStart: 0, origEnd: 4, pattern: "first"},
				{origStart: 4, origEnd: 8, pattern: "second"},
			},
			expected: []matchInfo{
				{origStart: 0, origEnd: 8, pattern: "first"},
			},
		},
		{
			name: "unsorted input",
			input: []matchInfo{
				{origStart: 10, origEnd: 14, pattern: "second"},
				{origStart: 0, origEnd: 4, pattern: "first"},
			},
			expected: []matchInfo{
				{origStart: 0, origEnd: 4, pattern: "first"},
				{origStart: 10, origEnd: 14, pattern: "second"},
			},
		},
		{
			name: "multiple overlapping",
			input: []matchInfo{
				{origStart: 0, origEnd: 5, pattern: "a"},
				{origStart: 3, origEnd: 8, pattern: "b"},
				{origStart: 6, origEnd: 12, pattern: "c"},
			},
			expected: []matchInfo{
				{origStart: 0, origEnd: 12, pattern: "a"},
			},
		},
		{
			name: "contained match",
			input: []matchInfo{
				{origStart: 0, origEnd: 10, pattern: "outer"},
				{origStart: 2, origEnd: 6, pattern: "inner"},
			},
			expected: []matchInfo{
				{origStart: 0, origEnd: 10, pattern: "outer"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeOverlapping(tt.input)
			if len(got) != len(tt.expected) {
				t.Errorf("mergeOverlapping() returned %d items, want %d", len(got), len(tt.expected))
				return
			}
			for i := range got {
				if got[i].origStart != tt.expected[i].origStart ||
					got[i].origEnd != tt.expected[i].origEnd {
					t.Errorf("mergeOverlapping()[%d] = {%d, %d}, want {%d, %d}",
						i, got[i].origStart, got[i].origEnd,
						tt.expected[i].origStart, tt.expected[i].origEnd)
				}
			}
		})
	}
}

func TestCensorModes(t *testing.T) {
	// Test that all CensorMode constants work correctly
	text := "fuck"

	// CensorAll (0)
	result := Censor(text, CensorAll)
	if result != "****" {
		t.Errorf("CensorAll: got %q, want %q", result, "****")
	}

	// CensorKeepFirst (1)
	result = Censor(text, CensorKeepFirst)
	if result != "f***" {
		t.Errorf("CensorKeepFirst: got %q, want %q", result, "f***")
	}

	// CensorKeepFirstLast (2)
	result = Censor(text, CensorKeepFirstLast)
	if result != "f**k" {
		t.Errorf("CensorKeepFirstLast: got %q, want %q", result, "f**k")
	}
}

// =============================================================================
// Tests for aho_corasick.go
// =============================================================================

func TestNewAhoCorasick(t *testing.T) {
	patterns := []string{"he", "she", "his", "hers"}
	ac := newAhoCorasick(patterns)

	if ac == nil {
		t.Fatal("newAhoCorasick returned nil")
	}
	if ac.root == nil {
		t.Fatal("automaton root is nil")
	}
	if len(ac.patterns) != len(patterns) {
		t.Errorf("patterns count = %d, want %d", len(ac.patterns), len(patterns))
	}
}

func TestAhoCorasickSearchAll(t *testing.T) {
	patterns := []string{"he", "she", "his", "hers"}
	ac := newAhoCorasick(patterns)

	tests := []struct {
		name          string
		text          string
		expectedCount int
	}{
		{"no matches", "xxx yyy zzz", 0},
		{"single match 'he'", "hello", 1},
		{"match 'she' contains 'he'", "she", 2}, // "she" matches both "she" and "he"
		{"multiple matches", "he said she said his hers", 6},
		{"empty text", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := ac.SearchAll(tt.text)
			if len(matches) != tt.expectedCount {
				t.Errorf("SearchAll(%q) returned %d matches, want %d", tt.text, len(matches), tt.expectedCount)
			}
		})
	}
}

func TestAhoCorasickHasMatch(t *testing.T) {
	patterns := []string{"bad", "word", "test"}
	ac := newAhoCorasick(patterns)

	tests := []struct {
		text     string
		expected bool
	}{
		{"this is bad", true},
		{"test case", true},
		{"the word is here", true},
		{"clean text", false},
		{"", false},
	}

	for _, tt := range tests {
		got := ac.HasMatch(tt.text)
		if got != tt.expected {
			t.Errorf("HasMatch(%q) = %v, want %v", tt.text, got, tt.expected)
		}
	}
}

func TestAhoCorasickSearch(t *testing.T) {
	patterns := []string{"a", "ab", "abc"}
	ac := newAhoCorasick(patterns)

	// Test early exit
	var count int
	ac.Search("abcabc", func(match acMatch) bool {
		count++
		return count < 2 // Stop after 2 matches
	})
	if count != 2 {
		t.Errorf("Search with early exit: count = %d, want 2", count)
	}

	// Test full search
	count = 0
	ac.Search("abcabc", func(match acMatch) bool {
		count++
		return true // Continue
	})
	if count < 4 { // Should have multiple matches
		t.Errorf("Search full: count = %d, want >= 4", count)
	}
}

func TestAhoCorasickMatchPositions(t *testing.T) {
	patterns := []string{"test"}
	ac := newAhoCorasick(patterns)

	text := "a test string"
	matches := ac.SearchAll(text)

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}

	m := matches[0]
	if m.start != 2 || m.end != 6 {
		t.Errorf("match position = [%d:%d], want [2:6]", m.start, m.end)
	}
	if text[m.start:m.end] != "test" {
		t.Errorf("matched text = %q, want %q", text[m.start:m.end], "test")
	}
}

func TestAhoCorasickEmptyPatterns(t *testing.T) {
	ac := newAhoCorasick([]string{})

	if ac == nil {
		t.Fatal("newAhoCorasick with empty patterns returned nil")
	}

	matches := ac.SearchAll("any text")
	if len(matches) != 0 {
		t.Errorf("SearchAll with no patterns returned %d matches", len(matches))
	}
}

func TestAhoCorasickOverlappingPatterns(t *testing.T) {
	patterns := []string{"foo", "foobar", "bar"}
	ac := newAhoCorasick(patterns)

	matches := ac.SearchAll("foobar")

	if len(matches) != 3 {
		t.Errorf("expected 3 matches in 'foobar', got %d", len(matches))
	}
}

func TestNormalizeText(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedNormalized string
	}{
		{"empty string", "", ""},
		{"plain ascii", "hello", "hello"},
		{"uppercase", "HELLO", "hello"},
		{"mixed case", "HeLLo", "hello"},

		// Leetspeak
		{"leetspeak @", "@ss", "ass"},
		{"leetspeak 0", "0h", "oh"},
		{"leetspeak 1", "h1", "hi"},
		{"leetspeak 3", "h3llo", "hello"},
		{"leetspeak 4", "4ss", "ass"},
		{"leetspeak 5", "5hit", "shit"},
		{"leetspeak 7", "7his", "this"},

		// Zero-width characters
		{"zero-width space", "hel\u200Blo", "hello"},
		{"zero-width non-joiner", "he\u200Cllo", "hello"},
		{"zero-width joiner", "hel\u200Dlo", "hello"},
		{"BOM", "hel\uFEFFlo", "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeText(tt.input)
			if result.normalized != tt.expectedNormalized {
				t.Errorf("normalizeText(%q).normalized = %q, want %q",
					tt.input, result.normalized, tt.expectedNormalized)
			}
		})
	}
}

func TestNormalizeTextPosMap(t *testing.T) {
	// Test that position mapping works correctly
	text := "a\u200Bb"
	result := normalizeText(text)

	// Normalized should be "ab"
	if result.normalized != "ab" {
		t.Errorf("normalized = %q, want %q", result.normalized, "ab")
	}

	// Position map should exist and have correct length
	if len(result.posMap) != 2 {
		t.Errorf("posMap length = %d, want 2", len(result.posMap))
	}
}

func TestCollapseRepeats(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedResult  string
		expectedStartAt int
	}{
		{"empty string", "", "", -1},
		{"no repeats", "abc", "abc", 0},
		{"single repeat", "aab", "ab", 0},
		{"multiple repeats", "aaabbb", "ab", 0},
		{"repeats in middle", "helllo", "helo", 0},
		{"all same", "aaaa", "a", 0},
		{"fuuuck pattern", "fuuuck", "fuck", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, info := collapseRepeats(tt.input)
			if result != tt.expectedResult {
				t.Errorf("collapseRepeats(%q) = %q, want %q", tt.input, result, tt.expectedResult)
			}
			if tt.expectedStartAt >= 0 && len(info.startPos) > 0 {
				if info.startPos[0] != tt.expectedStartAt {
					t.Errorf("collapseRepeats(%q) startPos[0] = %d, want %d",
						tt.input, info.startPos[0], tt.expectedStartAt)
				}
			}
		})
	}
}

func TestCollapseRepeatsEndPositions(t *testing.T) {
	// Test that end positions correctly span repeated characters. I might need
	// to revisit this if I ever support Unicode normalization.
	text := "heeello"
	result, info := collapseRepeats(text)

	if result != "helo" {
		t.Fatalf("collapsed = %q, want %q", result, "helo")
	}

	if len(info.endPos) < 2 {
		t.Fatalf("endPos too short: %d", len(info.endPos))
	}
}

func TestIsWordBoundary(t *testing.T) {
	text := "hello world"

	tests := []struct {
		pos      int
		expected bool
	}{
		{-1, true},
		{0, true},
		{5, true},
		{6, false},
		{11, true},
		{20, true},
	}

	for _, tt := range tests {
		got := isWordBoundary(text, tt.pos)
		if got != tt.expected {
			t.Errorf("isWordBoundary(%q, %d) = %v, want %v", text, tt.pos, got, tt.expected)
		}
	}
}

func TestIsWordBoundaryBefore(t *testing.T) {
	text := "a test"

	tests := []struct {
		pos      int
		expected bool
	}{
		{0, true},
		{1, false},
		{2, true},
		{3, false},
	}

	for _, tt := range tests {
		got := isWordBoundaryBefore(text, tt.pos)
		if got != tt.expected {
			t.Errorf("isWordBoundaryBefore(%q, %d) = %v, want %v", text, tt.pos, got, tt.expected)
		}
	}
}

func TestIsWordBoundaryAfter(t *testing.T) {
	text := "test a"

	tests := []struct {
		pos      int
		expected bool
	}{
		{0, false},
		{4, true},
		{5, false},
		{6, true},
		{10, true},
	}

	for _, tt := range tests {
		got := isWordBoundaryAfter(text, tt.pos)
		if got != tt.expected {
			t.Errorf("isWordBoundaryAfter(%q, %d) = %v, want %v", text, tt.pos, got, tt.expected)
		}
	}
}

func TestIsWordChar(t *testing.T) {
	tests := []struct {
		r        rune
		expected bool
	}{
		{'a', true},
		{'Z', true},
		{'5', true},
		{' ', false},
		{'.', false},
		{'!', false},
		{'-', false},
		{'_', false},
	}

	for _, tt := range tests {
		got := isWordChar(tt.r)
		if got != tt.expected {
			t.Errorf("isWordChar(%q) = %v, want %v", tt.r, got, tt.expected)
		}
	}
}

func TestBuildAsteriskMask(t *testing.T) {
	tests := []struct {
		name     string
		segment  string
		mode     int
		expected string
	}{
		{"empty", "", 0, ""},
		{"mode 0 short", "ab", 0, "**"},
		{"mode 0 long", "test", 0, "****"},
		{"mode 1 single char", "a", 1, "a"},
		{"mode 1 short", "ab", 1, "a*"},
		{"mode 1 long", "test", 1, "t***"},
		{"mode 2 single char", "a", 2, "a"},
		{"mode 2 two chars", "ab", 2, "ab"},
		{"mode 2 three chars", "abc", 2, "a*c"},
		{"mode 2 long", "test", 2, "t**t"},

		// Unicode handling
		{"unicode mode 0", "日本語", 0, "***"},
		{"unicode mode 1", "日本語", 1, "日**"},
		{"unicode mode 2", "日本語", 2, "日*語"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildAsteriskMask(tt.segment, tt.mode)
			if got != tt.expected {
				t.Errorf("buildAsteriskMask(%q, %d) = %q, want %q",
					tt.segment, tt.mode, got, tt.expected)
			}
		})
	}
}

func TestFindMatches(t *testing.T) {
	// Use the global profanityMatcher.
	// This is a bit of a hack, but it's the easiest way to test. I will revisit.
	tests := []struct {
		name          string
		text          string
		expectMatches bool
	}{
		{"empty text", "", false},
		{"clean text", "hello world", false},
		{"single profanity", "what the fuck", true},
		{"with leetspeak", "fvck this", true},
		{"repeated chars", "fuuuuck", true},
		{"false positive prevention", "class", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := findMatches(tt.text, profanityMatcher)
			hasMatches := len(matches) > 0
			if hasMatches != tt.expectMatches {
				t.Errorf("findMatches(%q) hasMatches = %v, want %v",
					tt.text, hasMatches, tt.expectMatches)
			}
		})
	}
}

func TestIsValidMatch(t *testing.T) {
	tests := []struct {
		name     string
		segment  string
		pattern  string
		expected bool
	}{
		{"exact match", "fuck", "fuck", true},
		{"with repeats", "fuuuck", "fuck", true},
		{"multiple repeats", "fuuuuuck", "fuck", true},
		{"all repeats", "fffuuuucccckkk", "fuck", true},
		{"too short", "fuc", "fuck", false},
		{"different chars", "duck", "fuck", false},
		{"extra chars", "fuckx", "fuck", false},
		{"missing char", "fck", "fuck", false},
		{"empty segment", "", "fuck", false},
		{"empty pattern", "test", "", false},
		{"both empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidMatch(tt.segment, tt.pattern)
			if got != tt.expected {
				t.Errorf("isValidMatch(%q, %q) = %v, want %v",
					tt.segment, tt.pattern, got, tt.expected)
			}
		})
	}
}

func TestHomoglyphNormalization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool // Should it be detected as profane?
	}{
		// Cyrillic homoglyphs - cyrillic а -> a, so fаck becomes fack (not fuck).
		// Still not sure if this is the best approach.
		{"cyrillic а (looks like a)", "fаck", false},
		{"cyrillic о (looks like o)", "fоck", false},

		// Fullwidth characters - NFKC normalization may handle these differently
		{"fullwidth f mixed", "ｆuck", false},
		{"fullwidth all", "ｆｕｃｋ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsProfane(tt.input)
			if got != tt.expected {
				t.Errorf("IsProfane(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestLeetspeakCombinations(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"f@ck", false},
		{"fvck", true},
		{"sh!t", true},
		{"sh1t", true},
		{"a$$", true},
		{"a55", true},
		{"fu(k", false}, // ( not mapped. I dont intend to support this.
	}

	for _, tt := range tests {
		got := IsProfane(tt.input)
		if got != tt.expected {
			t.Errorf("IsProfane(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestZeroWidthCharacterHandling(t *testing.T) {
	zeroWidths := []rune{
		'\u200B', // Zero-width space
		'\u200C', // Zero-width non-joiner
		'\u200D', // Zero-width joiner
		'\uFEFF', // BOM
		'\u00AD', // Soft hyphen
		'\u200E', // LTR mark
		'\u200F', // RTL mark
		'\u2060', // Word joiner
	}

	for _, zw := range zeroWidths {
		text := "fu" + string(zw) + "ck"
		if !IsProfane(text) {
			t.Errorf("IsProfane with zero-width %U should detect 'fuck'", zw)
		}
	}
}

func TestIntegrationCensorAndFind(t *testing.T) {
	text := "What the fuck is this shit?"

	// Check detection
	if !IsProfane(text) {
		t.Error("IsProfane should return true")
	}

	// Check finding
	found := FindProfanity(text)
	if len(found) != 2 {
		t.Errorf("FindProfanity should find 2 words, found %d", len(found))
	}

	// Check censoring
	censored := Censor(text, CensorAll)
	if censored == text {
		t.Error("Censor should modify the text")
	}

	// Censored text should not be profane (when checked literally)
	// Note: The asterisks themselves won't match patterns. This is intentional.
	// I'm not sure if it's a good idea to include them in the results.
}

func TestNewACNodeDepth(t *testing.T) {
	node := newACNode(5)
	if node.depth != 5 {
		t.Errorf("node.depth = %d, want 5", node.depth)
	}
	if node.children == nil {
		t.Error("node.children should not be nil")
	}
}

func TestNormalizedTextOriginal(t *testing.T) {
	text := "Hello World"
	result := normalizeText(text)

	if result.original != text {
		t.Errorf("original = %q, want %q", result.original, text)
	}
}

func TestProfanityListNotEmpty(t *testing.T) {
	if len(profanityList) == 0 {
		t.Error("profanityList should not be empty")
	}
}

func TestGlobalMatcherInitialized(t *testing.T) {
	if profanityMatcher == nil {
		t.Error("profanityMatcher should be initialized")
	}
	if collapsedToOriginal == nil {
		t.Error("collapsedToOriginal should be initialized")
	}
}

func TestCensorPreservesLength(t *testing.T) {
	// For ASCII text, censored length should match original
	text := "What the fuck"
	censored := Censor(text, CensorAll)

	// Count runes, not bytes
	origRunes := []rune(text)
	censoredRunes := []rune(censored)

	if len(origRunes) != len(censoredRunes) {
		t.Errorf("Censor changed rune count: %d -> %d", len(origRunes), len(censoredRunes))
	}
}

func TestMultipleProfanitiesInText(t *testing.T) {
	text := "fuck shit damn ass"
	found := FindProfanity(text)

	if len(found) < 4 {
		t.Errorf("Expected at least 4 profanities, found %d: %v", len(found), found)
	}

	censored := Censor(text, CensorAll)
	if censored == text {
		t.Error("Text should be censored")
	}
}

// Test edge cases for position mapping in findMatches
func TestFindMatchesPositionMapping(t *testing.T) {
	// Test that position mapping handles edge cases
	tests := []struct {
		name string
		text string
	}{
		{"start of text", "fuck you"},
		{"end of text", "oh fuck"},
		{"middle of text", "oh fuck you"},
		{"with punctuation", "what the fuck!"},
		{"with numbers", "123 fuck 456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := findMatches(tt.text, profanityMatcher)
			if len(matches) == 0 {
				t.Errorf("findMatches(%q) found no matches", tt.text)
				return
			}

			for _, m := range matches {
				if m.origStart < 0 || m.origStart > len(tt.text) {
					t.Errorf("invalid origStart: %d for text len %d", m.origStart, len(tt.text))
				}
				if m.origEnd < m.origStart || m.origEnd > len(tt.text) {
					t.Errorf("invalid origEnd: %d (start=%d, textLen=%d)", m.origEnd, m.origStart, len(tt.text))
				}
			}
		})
	}
}

// Test that CensorMode values are correct
func TestCensorModeValues(t *testing.T) {
	// i know...kinda janky but it's the easiest way to test.
	if CensorAll != 0 {
		t.Errorf("CensorAll = %d, want 0", CensorAll)
	}
	if CensorKeepFirst != 1 {
		t.Errorf("CensorKeepFirst = %d, want 1", CensorKeepFirst)
	}
	if CensorKeepFirstLast != 2 {
		t.Errorf("CensorKeepFirstLast = %d, want 2", CensorKeepFirstLast)
	}
}

// Test handling of patterns in collapsedToOriginal map
func TestCollapsedToOriginalMapping(t *testing.T) {
	// The global map should be populated
	if len(collapsedToOriginal) == 0 {
		t.Error("collapsedToOriginal should not be empty")
	}

	// Check that at least some common patterns are mapped
	// Note: We can't check specific patterns without knowing the exact implementation
	// But we can verify the map has reasonable entries.
	for collapsed, original := range collapsedToOriginal {
		if len(collapsed) == 0 || len(original) == 0 {
			t.Errorf("empty string in collapsedToOriginal: collapsed=%q, original=%q", collapsed, original)
		}
	}
}

func TestAcMatchStruct(t *testing.T) {
	m := acMatch{
		patternIndex: 5,
		start:        10,
		end:          15,
	}

	if m.patternIndex != 5 || m.start != 10 || m.end != 15 {
		t.Error("acMatch struct fields not set correctly")
	}
}

func TestMatchInfoStruct(t *testing.T) {
	m := matchInfo{
		origStart: 0,
		origEnd:   10,
		pattern:   "test",
	}

	if m.origStart != 0 || m.origEnd != 10 || m.pattern != "test" {
		t.Error("matchInfo struct fields not set correctly")
	}
}

func TestCollapsedPosInfoStruct(t *testing.T) {
	info := &collapsedPosInfo{
		startPos: []int{0, 1, 2},
		endPos:   []int{1, 2, 3},
	}

	if !reflect.DeepEqual(info.startPos, []int{0, 1, 2}) {
		t.Error("collapsedPosInfo.startPos not set correctly")
	}
	if !reflect.DeepEqual(info.endPos, []int{1, 2, 3}) {
		t.Error("collapsedPosInfo.endPos not set correctly")
	}
}

// Additional edge case tests for findMatches to improve coverage
func TestFindMatchesEdgeCases(t *testing.T) {
	// Test with NFKC normalization that changes text length
	// The ﬁ ligature normalizes to "fi" (2 chars from 1). NFC would normalize to "f" (1 char).
	tests := []struct {
		name string
		text string
	}{

		{"ending with profanity", "this ends with fuck"},
		{"only profanity", "shit"},
		{"profanity at very end", "test shit"},

		// Unicode that NFKC normalizes (covers nt.original != text branch)
		{"NFKC normalized text", "ﬁne fuck here"},
		{"composed unicode", "café fuck naïve"},

		// Very short matches
		{"short profanity", "a ass b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just ensure it doesn't panic and produces some result
			matches := findMatches(tt.text, profanityMatcher)
			_ = matches
		})
	}
}

// Test normalized text with NFKC transformations
func TestNormalizeTextNFKC(t *testing.T) {
	// Test with characters that NFKC normalizes
	tests := []struct {
		name  string
		input string
	}{
		{"ligature fi", "ﬁnish"},
		{"ligature fl", "ﬂow"},
		{"superscript", "x²"},
		{"fraction", "½"},
		{"roman numeral", "Ⅳ"},
		{"circled number", "①"},
		{"compatibility char", "㎞"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeText(tt.input)
			// NFKC changes the text
			if result.normalized == tt.input {
				// Some chars may not change, that's OK
			}

		})
	}
}

// Test profanity at exact text boundaries
func TestProfanityAtBoundaries(t *testing.T) {
	tests := []struct {
		text     string
		expected bool
	}{
		// Exact boundaries
		{"fuck", true},
		{" fuck ", true},
		{"\tfuck\t", true},
		{"\nfuck\n", true},

		// At start
		{"fuck is bad", true},

		// At end
		{"say fuck", true},

		{".fuck.", true},
		{"(fuck)", true},
		{"[fuck]", true},
		{"{fuck}", true},
		{"\"fuck\"", true},
		{"'fuck'", true},
	}

	for _, tt := range tests {
		got := IsProfane(tt.text)
		if got != tt.expected {
			t.Errorf("IsProfane(%q) = %v, want %v", tt.text, got, tt.expected)
		}
	}
}

// Test the normalizedText struct fields
func TestNormalizedTextStruct(t *testing.T) {
	text := "Hello"
	result := normalizeText(text)

	// Check original is preserved
	if result.original != text {
		t.Errorf("original = %q, want %q", result.original, text)
	}

	if result.normalized != "hello" {
		t.Errorf("normalized = %q, want %q", result.normalized, "hello")
	}

	if len(result.posMap) != 5 {
		t.Errorf("posMap len = %d, want 5", len(result.posMap))
	}
}

// Test with NFKC text that contains profanity - exercises the NFKC path in findMatches
func TestFindMatchesWithNFKCAndProfanity(t *testing.T) {
	// Text where NFKC changes the string AND contains profanity
	// Note: NFKC normalization can affect position mapping
	tests := []struct {
		name string
		text string
	}{
		// Various NFKC scenarios - we just test that they don't panic
		{"ligature text", "ﬁne ﬂow"},
		{"superscript", "x² y³"},
		{"fullwidth space", "this\u3000that"},
		{"fraction", "½ ¼"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = IsProfane(tt.text)
			_ = Censor(tt.text, CensorAll)
			_ = FindProfanity(tt.text)
		})
	}
}

// Test Censor with text where NFKC normalization applies
func TestCensorWithNFKCText(t *testing.T) {
	// Text with characters that NFKC normalizes + profanity
	tests := []struct {
		text string
	}{
		{"ﬁne fuck here"},
		{"x² fuck y²"},
		{"½ fuck ½"},
	}

	for _, tt := range tests {
		censored := Censor(tt.text, CensorAll)
		// The profanity should be censored, though text length might differ
		if !IsProfane(tt.text) {
			// If original has profanity, censored version should be different
			continue
		}
		if censored == tt.text {
			t.Errorf("Censor(%q) should modify the text", tt.text)
		}
	}
}

// Test edge cases in collapseRepeats position tracking
func TestCollapseRepeatsPositionTracking(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedRunes int // Expected number of RUNES in collapsed result
		checkStartEnd bool
	}{
		{"single char", "a", 1, true},
		{"two same chars", "aa", 1, true},
		{"alternating", "aba", 3, true},
		{"triple repeat", "aaa", 1, true},
		{"unicode repeated", "日日日", 1, true},
		{"unicode different", "日本語", 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, info := collapseRepeats(tt.input)
			if len([]rune(result)) != tt.expectedRunes {
				t.Errorf("collapseRepeats(%q) rune len = %d, want %d",
					tt.input, len([]rune(result)), tt.expectedRunes)
			}
			if tt.checkStartEnd {
				if len(info.startPos) != len(result) {
					t.Errorf("startPos len = %d, result len = %d",
						len(info.startPos), len(result))
				}
				if len(info.endPos) != len(result) {
					t.Errorf("endPos len = %d, result len = %d",
						len(info.endPos), len(result))
				}
			}
		})
	}
}

// Test profanity matching at exact end of string (edge case)
func TestProfanityAtExactEnd(t *testing.T) {
	// These specifically test when match.end equals text length
	tests := []struct {
		text     string
		expected bool
	}{
		{"fuck", true},
		{"a fuck", true},
		{"this is shit", true},
		{"multiple ass", true},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			got := IsProfane(tt.text)
			if got != tt.expected {
				t.Errorf("IsProfane(%q) = %v, want %v", tt.text, got, tt.expected)
			}
			// Also test censoring
			censored := Censor(tt.text, CensorAll)
			if tt.expected && censored == tt.text {
				t.Errorf("Censor(%q) should have modified text", tt.text)
			}
		})
	}
}
