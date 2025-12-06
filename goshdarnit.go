package goshdarnit

// CensorMode controls how profanity is censored.
type CensorMode int

const (
	// CensorAll replaces all characters with asterisks.
	CensorAll CensorMode = iota
	// CensorKeepFirst keeps the first character visible.
	CensorKeepFirst
	// CensorKeepFirstLast keeps the first and last characters visible.
	CensorKeepFirstLast
)

// profanityMatcher is the global Aho-Corasick automaton, initialized once at package load.
var profanityMatcher *ahoCorasick

var collapsedToOriginal map[string]string

func init() {
	collapsedPatterns := make([]string, 0, len(profanityList))
	collapsedToOriginal = make(map[string]string, len(profanityList))
	seen := make(map[string]struct{})

	for _, pattern := range profanityList {
		collapsed, _ := collapseRepeats(pattern)
		// Avoid duplicates (e.g., "ass" and "as" might both become "as")
		if _, exists := seen[collapsed]; !exists {
			seen[collapsed] = struct{}{}
			collapsedPatterns = append(collapsedPatterns, collapsed)
			collapsedToOriginal[collapsed] = pattern
		}
	}

	profanityMatcher = newAhoCorasick(collapsedPatterns)
}

// IsProfane returns true if the text contains any profanity.
// The function handles various evasion techniques and uses word boundary
// detection to avoid false positives.
func IsProfane(text string) bool {
	if len(text) == 0 {
		return false
	}

	matches := findMatches(text, profanityMatcher)
	return len(matches) > 0
}

// Censor replaces profanity in the text with asterisks.
// The mode parameter controls which characters to reveal:
//   - CensorAll: replaces all characters with asterisks
//   - CensorKeepFirst: keeps the first character visible (e.g., "f***")
//   - CensorKeepFirstLast: keeps first and last characters visible (e.g., "f**k")
//
// The returned string preserves the original length.
func Censor(text string, mode CensorMode) string {
	if len(text) == 0 {
		return text
	}

	matches := findMatches(text, profanityMatcher)
	if len(matches) == 0 {
		return text
	}

	// Sort and merge overlapping matches
	matches = mergeOverlapping(matches)

	// Build the censored string
	result := make([]byte, 0, len(text))
	lastEnd := 0

	for _, m := range matches {
		// Add text before this match
		if m.origStart > lastEnd {
			result = append(result, text[lastEnd:m.origStart]...)
		}

		// Censor this match
		segment := text[m.origStart:m.origEnd]
		mask := buildAsteriskMask(segment, int(mode))
		result = append(result, mask...)

		lastEnd = m.origEnd
	}

	// Add any remaining text
	if lastEnd < len(text) {
		result = append(result, text[lastEnd:]...)
	}

	return string(result)
}

// CensorWithDefault is a convenience function that censors with CensorAll mode.
func CensorWithDefault(text string) string {
	return Censor(text, CensorAll)
}

func mergeOverlapping(matches []matchInfo) []matchInfo {
	if len(matches) <= 1 {
		return matches
	}

	// Sort by start position (simple insertion sort for typically small slices)
	for i := 1; i < len(matches); i++ {
		j := i
		for j > 0 && matches[j].origStart < matches[j-1].origStart {
			matches[j], matches[j-1] = matches[j-1], matches[j]
			j--
		}
	}

	// Merge overlapping
	result := make([]matchInfo, 0, len(matches))
	current := matches[0]

	for i := 1; i < len(matches); i++ {
		if matches[i].origStart <= current.origEnd {
			// Overlapping, extend current
			if matches[i].origEnd > current.origEnd {
				current.origEnd = matches[i].origEnd
			}
		} else {
			// Not overlapping, save current and start new
			result = append(result, current)
			current = matches[i]
		}
	}
	result = append(result, current)

	return result
}

// FindProfanity returns a slice of profane words discovered in the txet.
// Returns nil if no profanity is found.
func FindProfanity(text string) []string {
	if len(text) == 0 {
		return nil
	}

	matches := findMatches(text, profanityMatcher)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]struct{})
	var result []string

	for _, m := range matches {
		// Map collapsed pattern back to original
		original := m.pattern
		if orig, exists := collapsedToOriginal[m.pattern]; exists {
			original = orig
		}
		if _, exists := seen[original]; !exists {
			seen[original] = struct{}{}
			result = append(result, original)
		}
	}

	return result
}
