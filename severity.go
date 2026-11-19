package goshdarnit

// Severity ranks how offensive a listed word is. A higher value is worse.
type Severity int

const (
	// SeverityMild covers words most audiences accept and ordinary words that
	// can also be insults: "damn", "crap", "idiot", "queen", clinical anatomy.
	SeverityMild Severity = iota + 1
	// SeverityModerate covers ordinary profanity and crude insults: "fuck",
	// "shit", "bitch", "dick".
	SeverityModerate
	// SeveritySevere covers slurs, explicit sexual terms and sexual insults.
	SeveritySevere
)

// String returns "mild", "moderate" or "severe".
func (s Severity) String() string {
	switch s {
	case SeverityMild:
		return "mild"
	case SeverityModerate:
		return "moderate"
	case SeveritySevere:
		return "severe"
	default:
		return "unknown"
	}
}

// wordSeverity maps every listed word to its severity. It is set at package
// load.
var wordSeverity map[string]Severity

// buildWordSeverity gives each listed word its group's severity. Words that
// collapse to the same pattern ("spaz", "spazz") match the same text, so they
// share the highest severity among them.
func buildWordSeverity() map[string]Severity {
	collapsedSeverity := make(map[string]Severity, len(profanityList))
	for _, group := range wordGroups {
		for _, word := range group.words {
			collapsed, _ := collapseRepeats(word)
			collapsedSeverity[collapsed] = max(collapsedSeverity[collapsed], group.severity)
		}
	}

	severities := make(map[string]Severity, len(profanityList))
	for _, word := range profanityList {
		collapsed, _ := collapseRepeats(word)
		severities[word] = collapsedSeverity[collapsed]
	}
	return severities
}

// SeverityOf returns the severity of a listed word, such as one returned by
// FindProfanity. It returns 0 for a word that is not listed.
func SeverityOf(word string) Severity {
	return wordSeverity[word]
}

// FindProfanityAtLeast is FindProfanity limited to words whose severity is
// at least minSeverity. Returns nil if none is found.
func FindProfanityAtLeast(text string, minSeverity Severity) []string {
	var result []string
	for _, word := range FindProfanity(text) {
		if wordSeverity[word] >= minSeverity {
			result = append(result, word)
		}
	}
	return result
}

// IsProfaneAtLeast returns true if the text contains a word whose severity
// is at least minSeverity.
func IsProfaneAtLeast(text string, minSeverity Severity) bool {
	return len(FindProfanityAtLeast(text, minSeverity)) > 0
}
