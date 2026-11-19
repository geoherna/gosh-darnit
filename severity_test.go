package goshdarnit

import (
	"slices"
	"testing"
)

func TestSeverityOf(t *testing.T) {
	tests := []struct {
		word string
		want Severity
	}{
		{"damn", SeverityMild},
		{"crap", SeverityMild},
		{"idiot", SeverityMild},
		{"queen", SeverityMild},
		{"fuck", SeverityModerate},
		{"shit", SeverityModerate},
		{"pornography", SeverityModerate},
		{"faggot", SeveritySevere},
		{"coon", SeveritySevere},
		{"retard", SeveritySevere},
		{"whore", SeveritySevere},
		{"porn", SeveritySevere},
		{"bob", 0},
		{"", 0},
	}

	for _, tt := range tests {
		if got := SeverityOf(tt.word); got != tt.want {
			t.Errorf("SeverityOf(%q) = %v, want %v", tt.word, got, tt.want)
		}
	}
}

func TestEveryListedWordHasOneSeverity(t *testing.T) {
	seen := make(map[string]Severity)
	for _, group := range wordGroups {
		if group.severity < SeverityMild || group.severity > SeveritySevere {
			t.Errorf("group starting %q has severity %d", group.words[0], group.severity)
		}
		for _, word := range group.words {
			if previous, listed := seen[word]; listed {
				t.Errorf("%q is listed twice (%v and %v)", word, previous, group.severity)
			}
			seen[word] = group.severity
		}
	}

	for _, word := range profanityList {
		if SeverityOf(word) == 0 {
			t.Errorf("SeverityOf(%q) = 0, want a severity", word)
		}
	}
}

func TestWordsThatCollapseTogetherShareSeverity(t *testing.T) {
	byCollapsed := make(map[string]Severity)
	for _, word := range profanityList {
		collapsed, _ := collapseRepeats(word)
		severity := SeverityOf(word)
		if previous, ok := byCollapsed[collapsed]; ok && previous != severity {
			t.Errorf("%q collapses to %q with severity %v, but another word there has %v", word, collapsed, severity, previous)
		}
		byCollapsed[collapsed] = severity
	}

	if SeverityOf("spaz") != SeverityOf("spazz") {
		t.Errorf("spaz and spazz collapse to one pattern but have different severities")
	}
}

func TestFindProfanityAtLeast(t *testing.T) {
	text := "damn this shit, you retard"
	tests := []struct {
		name        string
		text        string
		minSeverity Severity
		want        []string
	}{
		{"mild finds every tier", text, SeverityMild, []string{"damn", "shit", "retard"}},
		{"moderate drops mild", text, SeverityModerate, []string{"shit", "retard"}},
		{"severe keeps only slurs", text, SeveritySevere, []string{"retard"}},
		{"no severe word", "fix this fucking chart", SeveritySevere, nil},
		{"clean text", "Hello world", SeverityMild, nil},
		{"empty text", "", SeverityMild, nil},
		{"leetspeak slur", "n1gg3r", SeveritySevere, []string{"nigger"}},
		{"repeated letters", "fuuuuck", SeverityModerate, []string{"fuck"}},
		{"collapsed pair", "spazzzz", SeveritySevere, []string{"spaz"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindProfanityAtLeast(tt.text, tt.minSeverity)
			if tt.want == nil {
				if got != nil {
					t.Errorf("FindProfanityAtLeast(%q, %v) = %v, want nil", tt.text, tt.minSeverity, got)
				}
				return
			}
			slices.Sort(got)
			want := slices.Sorted(slices.Values(tt.want))
			if !slices.Equal(got, want) {
				t.Errorf("FindProfanityAtLeast(%q, %v) = %v, want %v", tt.text, tt.minSeverity, got, want)
			}
		})
	}
}

func TestIsProfaneAtLeast(t *testing.T) {
	tests := []struct {
		text        string
		minSeverity Severity
		want        bool
	}{
		{"what the hell", SeverityMild, true},
		{"what the hell", SeverityModerate, false},
		{"this is bullshit", SeverityModerate, true},
		{"this is bullshit", SeveritySevere, false},
		{"you faggot", SeveritySevere, true},
		{"Hello world", SeverityMild, false},
	}

	for _, tt := range tests {
		if got := IsProfaneAtLeast(tt.text, tt.minSeverity); got != tt.want {
			t.Errorf("IsProfaneAtLeast(%q, %v) = %v, want %v", tt.text, tt.minSeverity, got, tt.want)
		}
	}
}

func TestSeverityString(t *testing.T) {
	tests := map[Severity]string{
		SeverityMild:     "mild",
		SeverityModerate: "moderate",
		SeveritySevere:   "severe",
		0:                "unknown",
	}
	for severity, want := range tests {
		if got := severity.String(); got != want {
			t.Errorf("Severity(%d).String() = %q, want %q", int(severity), got, want)
		}
	}
}
