package goshdarnit

import (
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

// zeroWidthChars contains Unicode zero-width and invisible characters to strip
var zeroWidthChars = map[rune]struct{}{
	'\u200B': {}, // Zero-width space
	'\u200C': {}, // Zero-width non-joiner
	'\u200D': {}, // Zero-width joiner
	'\uFEFF': {}, // Byte order mark / zero-width no-break space
	'\u00AD': {}, // Soft hyphen
	'\u200E': {}, // Left-to-right mark
	'\u200F': {}, // Right-to-left mark
	'\u2060': {}, // Word joiner
	'\u2061': {}, // Function application
	'\u2062': {}, // Invisible times
	'\u2063': {}, // Invisible separator
	'\u2064': {}, // Invisible plus
	'\u180E': {}, // Mongolian vowel separator
	'\u034F': {}, // Combining grapheme joiner
}

// leetspeakMap maps leetspeak characters to their ASCII equivalents
var leetspeakMap = map[rune]rune{
	'@': 'a',
	'4': 'a',
	'8': 'b',
	'3': 'e',
	'!': 'i',
	'1': 'i',
	'|': 'i',
	'0': 'o',
	'$': 's',
	'5': 's',
	'7': 't',
	'+': 't',
	'v': 'u', // common substitution for u
	'V': 'u',
}

// homoglyphMap maps common Unicode homoglyphs to ASCII equivalents
// This includes Cyrillic and other look-alike characters
var homoglyphMap = map[rune]rune{
	// Cyrillic homoglyphs
	'а': 'a', // Cyrillic а
	'А': 'a', // Cyrillic А
	'е': 'e', // Cyrillic е
	'Е': 'e', // Cyrillic Е
	'ё': 'e', // Cyrillic ё
	'Ё': 'e', // Cyrillic Ё
	'о': 'o', // Cyrillic о
	'О': 'o', // Cyrillic О
	'р': 'p', // Cyrillic р
	'Р': 'p', // Cyrillic Р
	'с': 'c', // Cyrillic с
	'С': 'c', // Cyrillic С
	'у': 'y', // Cyrillic у
	'У': 'y', // Cyrillic У
	'х': 'x', // Cyrillic х
	'Х': 'x', // Cyrillic Х
	'і': 'i', // Cyrillic і (Ukrainian)
	'І': 'i', // Cyrillic І (Ukrainian)
	'ї': 'i', // Cyrillic ї (Ukrainian)
	'Ї': 'i', // Cyrillic Ї (Ukrainian)
	'ј': 'j', // Cyrillic ј (Serbian)
	'Ј': 'j', // Cyrillic Ј (Serbian)
	'ӏ': 'l', // Cyrillic ӏ
	'Ӏ': 'l', // Cyrillic Ӏ
	'К': 'k', // Cyrillic К
	'к': 'k', // Cyrillic к
	'М': 'm', // Cyrillic М
	'м': 'm', // Cyrillic м
	'Н': 'h', // Cyrillic Н (looks like H)
	'Т': 't', // Cyrillic Т
	'т': 't', // Cyrillic т
	'В': 'b', // Cyrillic В (looks like B)
	'в': 'b', // Cyrillic в
	'Ѕ': 's', // Cyrillic Ѕ
	'ѕ': 's', // Cyrillic ѕ
	'Ԁ': 'd', // Cyrillic Ԁ
	'ԁ': 'd', // Cyrillic ԁ

	// Greek homoglyphs
	'Α': 'a', // Greek Alpha
	'α': 'a',
	'Β': 'b', // Greek Beta
	'β': 'b',
	'Ε': 'e', // Greek Epsilon
	'ε': 'e',
	'Η': 'h', // Greek Eta
	'η': 'n',
	'Ι': 'i', // Greek Iota
	'ι': 'i',
	'Κ': 'k', // Greek Kappa
	'κ': 'k',
	'Μ': 'm', // Greek Mu
	'μ': 'u',
	'Ν': 'n', // Greek Nu
	'ν': 'v',
	'Ο': 'o', // Greek Omicron
	'ο': 'o',
	'Ρ': 'p', // Greek Rho
	'ρ': 'p',
	'Τ': 't', // Greek Tau
	'τ': 't',
	'Υ': 'y', // Greek Upsilon
	'υ': 'u',
	'Χ': 'x', // Greek Chi
	'χ': 'x',

	// Fullwidth characters
	'Ａ': 'a', 'ａ': 'a',
	'Ｂ': 'b', 'ｂ': 'b',
	'Ｃ': 'c', 'ｃ': 'c',
	'Ｄ': 'd', 'ｄ': 'd',
	'Ｅ': 'e', 'ｅ': 'e',
	'Ｆ': 'f', 'ｆ': 'f',
	'Ｇ': 'g', 'ｇ': 'g',
	'Ｈ': 'h', 'ｈ': 'h',
	'Ｉ': 'i', 'ｉ': 'i',
	'Ｊ': 'j', 'ｊ': 'j',
	'Ｋ': 'k', 'ｋ': 'k',
	'Ｌ': 'l', 'ｌ': 'l',
	'Ｍ': 'm', 'ｍ': 'm',
	'Ｎ': 'n', 'ｎ': 'n',
	'Ｏ': 'o', 'ｏ': 'o',
	'Ｐ': 'p', 'ｐ': 'p',
	'Ｑ': 'q', 'ｑ': 'q',
	'Ｒ': 'r', 'ｒ': 'r',
	'Ｓ': 's', 'ｓ': 's',
	'Ｔ': 't', 'ｔ': 't',
	'Ｕ': 'u', 'ｕ': 'u',
	'Ｖ': 'v', 'ｖ': 'v',
	'Ｗ': 'w', 'ｗ': 'w',
	'Ｘ': 'x', 'ｘ': 'x',
	'Ｙ': 'y', 'ｙ': 'y',
	'Ｚ': 'z', 'ｚ': 'z',
	'０': 'o', '１': 'i', '２': 'z', '３': 'e',
	'４': 'a', '５': 's', '６': 'b', '７': 't',
	'８': 'b', '９': 'g',

	// Other common homoglyphs
	'ℓ': 'l', // Script small l
	'ℒ': 'l', // Script capital L
	'ℐ': 'i', // Script capital I
	'ℑ': 'i', // Black-letter capital I
	'ℊ': 'g', // Script small g
	'ℋ': 'h', // Script capital H
	'ℌ': 'h', // Black-letter capital H
	'ℍ': 'h', // Double-struck capital H
	'ℕ': 'n', // Double-struck capital N
	'ℙ': 'p', // Double-struck capital P
	'ℚ': 'q', // Double-struck capital Q
	'ℛ': 'r', // Script capital R
	'ℜ': 'r', // Black-letter capital R
	'ℝ': 'r', // Double-struck capital R
	'ℤ': 'z', // Double-struck capital Z
	'ℨ': 'z', // Black-letter capital Z
	'ℬ': 'b', // Script capital B
	'ℭ': 'c', // Black-letter capital C
	'ℰ': 'e', // Script capital E
	'ℱ': 'f', // Script capital F
	'ℳ': 'm', // Script capital M

	// Enclosed/circled letters - map to regular letters
	'Ⓐ': 'a', 'ⓐ': 'a', 'Ⓑ': 'b', 'ⓑ': 'b',
	'Ⓒ': 'c', 'ⓒ': 'c', 'Ⓓ': 'd', 'ⓓ': 'd',
	'Ⓔ': 'e', 'ⓔ': 'e', 'Ⓕ': 'f', 'ⓕ': 'f',
	'Ⓖ': 'g', 'ⓖ': 'g', 'Ⓗ': 'h', 'ⓗ': 'h',
	'Ⓘ': 'i', 'ⓘ': 'i', 'Ⓙ': 'j', 'ⓙ': 'j',
	'Ⓚ': 'k', 'ⓚ': 'k', 'Ⓛ': 'l', 'ⓛ': 'l',
	'Ⓜ': 'm', 'ⓜ': 'm', 'Ⓝ': 'n', 'ⓝ': 'n',
	'Ⓞ': 'o', 'ⓞ': 'o', 'Ⓟ': 'p', 'ⓟ': 'p',
	'Ⓠ': 'q', 'ⓠ': 'q', 'Ⓡ': 'r', 'ⓡ': 'r',
	'Ⓢ': 's', 'ⓢ': 's', 'Ⓣ': 't', 'ⓣ': 't',
	'Ⓤ': 'u', 'ⓤ': 'u', 'Ⓥ': 'v', 'ⓥ': 'v',
	'Ⓦ': 'w', 'ⓦ': 'w', 'Ⓧ': 'x', 'ⓧ': 'x',
	'Ⓨ': 'y', 'ⓨ': 'y', 'Ⓩ': 'z', 'ⓩ': 'z',
}

// normalizedText holds the result of text normalization, tracking the mapping
// between original positions and normalized positions.
type normalizedText struct {
	original   string
	normalized string
	// posMap maps normalized byte positions to original byte positions
	posMap []int
}

// normalizeText applies all normalization steps to the input text and returns
// a normalizedText struct that tracks position mappings.
func normalizeText(text string) *normalizedText {
	// Step 1: NFKC normalization
	nfkc := norm.NFKC.String(text)

	// Pre-allocate with capacity estimation
	result := make([]byte, 0, len(nfkc))
	posMap := make([]int, 0, len(nfkc))

	// Track byte position in original NFKC-normalized string
	origPos := 0

	for _, r := range nfkc {
		runeLen := utf8.RuneLen(r)
		startPos := origPos

		// Step 2: Skip zero-width characters
		if _, isZeroWidth := zeroWidthChars[r]; isZeroWidth {
			origPos += runeLen
			continue
		}

		// Step 3: Convert to lowercase
		r = unicode.ToLower(r)

		// Step 4: Apply homoglyph mapping
		if mapped, exists := homoglyphMap[r]; exists {
			r = mapped
		}

		// Step 5: Apply leetspeak mapping
		if mapped, exists := leetspeakMap[r]; exists {
			r = mapped
		}

		// Encode the rune
		var buf [4]byte
		n := utf8.EncodeRune(buf[:], r)
		for i := 0; i < n; i++ {
			result = append(result, buf[i])
			posMap = append(posMap, startPos)
		}

		origPos += runeLen
	}

	return &normalizedText{
		original:   text,
		normalized: string(result),
		posMap:     posMap,
	}
}

// collapsedPosInfo holds position mapping information for collapsed text
type collapsedPosInfo struct {
	startPos []int // maps collapsed byte index to start position in original
	endPos   []int // maps collapsed byte index to end position in original (after all repeats)
}

// collapseRepeats collapses repeated characters in a string for matching.
// For example, "fuuuuck" becomes "fuck".
// Returns the collapsed string and position mapping info.
func collapseRepeats(text string) (string, *collapsedPosInfo) {
	if len(text) == 0 {
		return "", &collapsedPosInfo{}
	}

	result := make([]byte, 0, len(text))
	startPos := make([]int, 0, len(text))
	endPos := make([]int, 0, len(text))

	var lastRune rune
	lastRuneSet := false
	bytePos := 0
	lastCharEndPos := 0

	for _, r := range text {
		runeLen := utf8.RuneLen(r)
		nextPos := bytePos + runeLen

		// Skip if this is a repeat of the last character
		if lastRuneSet && r == lastRune {
			// Update the end position of the last character to include this repeat
			lastCharEndPos = nextPos
			bytePos = nextPos
			continue
		}

		// Update end positions for previous character if needed
		if lastRuneSet && len(endPos) > 0 {
			// Fill in end positions for all bytes of the last rune
			lastRuneLen := utf8.RuneLen(lastRune)
			for i := len(endPos) - lastRuneLen; i < len(endPos); i++ {
				endPos[i] = lastCharEndPos
			}
		}

		// Add this character
		var buf [4]byte
		n := utf8.EncodeRune(buf[:], r)
		for i := 0; i < n; i++ {
			result = append(result, buf[i])
			startPos = append(startPos, bytePos)
			endPos = append(endPos, nextPos) // will be updated if there are repeats
		}

		lastRune = r
		lastRuneSet = true
		lastCharEndPos = nextPos
		bytePos = nextPos
	}

	// Update end positions for the final character
	if lastRuneSet && len(endPos) > 0 {
		lastRuneLen := utf8.RuneLen(lastRune)
		for i := len(endPos) - lastRuneLen; i < len(endPos); i++ {
			endPos[i] = lastCharEndPos
		}
	}

	return string(result), &collapsedPosInfo{startPos: startPos, endPos: endPos}
}

// isWordBoundary returns true if the position is at a word boundary.
// A word boundary is at the start/end of the string, or adjacent to a non-word character.
func isWordBoundary(text string, pos int) bool {
	if pos < 0 || pos > len(text) {
		return true
	}
	if pos == 0 || pos == len(text) {
		return true
	}

	// Check the character at the boundary
	r, _ := utf8.DecodeRuneInString(text[pos:])
	return !isWordChar(r)
}

// isWordBoundaryBefore checks if there's a word boundary before the given position.
func isWordBoundaryBefore(text string, pos int) bool {
	if pos <= 0 {
		return true
	}

	// Find the rune before this position
	r, _ := utf8.DecodeLastRuneInString(text[:pos])
	return !isWordChar(r)
}

// isWordBoundaryAfter checks if there's a word boundary after the given position.
func isWordBoundaryAfter(text string, pos int) bool {
	if pos >= len(text) {
		return true
	}

	r, _ := utf8.DecodeRuneInString(text[pos:])
	return !isWordChar(r)
}

// isWordChar returns true if the rune is a word character (letter or digit).
func isWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// buildAsteriskMask creates an asterisk mask of the given length.
// mode controls which characters to reveal:
//
//	0 = all asterisks
//	1 = reveal first character
//	2 = reveal first and last characters
func buildAsteriskMask(originalSegment string, mode int) string {
	runes := []rune(originalSegment)
	length := len(runes)

	if length == 0 {
		return ""
	}

	result := make([]rune, length)
	for i := range result {
		result[i] = '*'
	}

	switch mode {
	case 1:
		if length >= 1 {
			result[0] = runes[0]
		}
	case 2:
		if length >= 1 {
			result[0] = runes[0]
		}
		if length >= 2 {
			result[length-1] = runes[length-1]
		}
	}

	return string(result)
}

// matchInfo holds information about a profanity match in text.
type matchInfo struct {
	// Start position in original text
	origStart int
	// End position in original text (exclusive)
	origEnd int
	// The pattern that matched
	pattern string
}

// findMatches finds all profanity matches in the given text using the global automaton.
// It handles normalization, repeat collapsing, and word boundary checking.
func findMatches(text string, ac *ahoCorasick) []matchInfo {
	if len(text) == 0 {
		return nil
	}

	// First normalize the text
	nt := normalizeText(text)

	// Then collapse repeats on the normalized text
	collapsed, collapsedPos := collapseRepeats(nt.normalized)

	// Search in the collapsed text
	matches := ac.SearchAll(collapsed)
	if len(matches) == 0 {
		return nil
	}

	var results []matchInfo

	for _, m := range matches {
		// Map from collapsed position back to normalized position
		normalizedStart := 0
		normalizedEnd := len(nt.normalized)

		if m.start < len(collapsedPos.startPos) {
			normalizedStart = collapsedPos.startPos[m.start]
		}
		// For end position, use the endPos which accounts for repeated characters
		if m.end > 0 && m.end <= len(collapsedPos.endPos) {
			normalizedEnd = collapsedPos.endPos[m.end-1]
		} else if m.end >= len(collapsedPos.endPos) {
			normalizedEnd = len(nt.normalized)
		}

		// Map back to original text positions
		// The posMap maps normalized positions to NFKC positions
		origStart := 0
		origEnd := len(text)

		if normalizedStart < len(nt.posMap) {
			origStart = nt.posMap[normalizedStart]
		}

		// Find the original end position
		if normalizedEnd > 0 && normalizedEnd <= len(nt.posMap) {
			origEnd = nt.posMap[normalizedEnd-1]
			// Add the length of the last rune in original
			if origEnd < len(nt.original) {
				_, runeLen := utf8.DecodeRuneInString(nt.original[origEnd:])
				origEnd += runeLen
			}
		} else if normalizedEnd >= len(nt.posMap) {
			// Match extends to end
			origEnd = len(nt.original)
		}

		// For NFKC normalized text, map back to actual original if needed
		if nt.original != text {
			// NFKC changed the text, use approximation
			if origStart > len(text) {
				origStart = len(text)
			}
			if origEnd > len(text) {
				origEnd = len(text)
			}
		}

		// Check word boundaries on the ORIGINAL text, not normalized
		// This is crucial because leetspeak mapping changes punctuation to letters
		if !isWordBoundaryBefore(text, origStart) ||
			!isWordBoundaryAfter(text, origEnd) {
			continue
		}

		// Get the collapsed pattern that matched
		collapsedPattern := ac.patterns[m.patternIndex]

		// Validate that the match is legitimate (not a false positive due to collapsing)
		// Get the normalized segment that matched
		normalizedSegment := nt.normalized[normalizedStart:normalizedEnd]

		// Get the original (uncollapsed) pattern for validation
		originalPattern := collapsedPattern
		if orig, exists := collapsedToOriginal[collapsedPattern]; exists {
			originalPattern = orig
		}

		if !isValidMatch(normalizedSegment, originalPattern) {
			continue
		}

		results = append(results, matchInfo{
			origStart: origStart,
			origEnd:   origEnd,
			pattern:   collapsedPattern,
		})
	}

	return results
}

// isValidMatch checks if a normalized segment is a valid match for a pattern.
// This prevents false positives where legitimate words collapse to match profanity patterns.
// A match is valid if the segment is the pattern with optional repeated characters added.
func isValidMatch(normalizedSegment, pattern string) bool {
	// Quick check: segment must be at least as long as pattern
	if len(normalizedSegment) < len(pattern) {
		return false
	}

	// If they're equal, it's definitely a match
	if normalizedSegment == pattern {
		return true
	}

	// Check if segment is pattern with some characters repeated
	// Walk through pattern and segment together, allowing repeats
	segmentRunes := []rune(normalizedSegment)
	patternRunes := []rune(pattern)

	si := 0 // segment index
	pi := 0 // pattern index

	for pi < len(patternRunes) && si < len(segmentRunes) {
		if segmentRunes[si] == patternRunes[pi] {
			// Characters match, advance pattern
			lastMatch := segmentRunes[si]
			si++
			pi++

			// Skip any repeated characters in segment
			for si < len(segmentRunes) && segmentRunes[si] == lastMatch {
				si++
			}
		} else {
			// Characters don't match - not a valid extension
			return false
		}
	}

	// Valid if we consumed the entire pattern
	// Any remaining characters in segment after pattern is consumed means
	// the segment has extra characters that aren't just repeats
	return pi == len(patternRunes) && si == len(segmentRunes)
}


