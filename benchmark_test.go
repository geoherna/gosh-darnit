package goshdarnit

import (
	"strings"
	"testing"
)

// Sample texts for benchmarking
var (
	cleanShort    = "Hello, how are you today?"
	cleanLong     = strings.Repeat("This is a perfectly clean sentence without any bad words. ", 100)
	profaneShort  = "What the fuck is going on?"
	profaneLong   = strings.Repeat("This is some text with shit and fuck scattered throughout. ", 100)
	leetspeak     = "wh4t th3 fvck 1s g01ng 0n?"
	repeatedChars = "What the fuuuuuuuuuck is happening?"
	mixedText     = "Hello world! Some shit happened, but it's all good now. The analyst was helpful."
)

// BenchmarkIsProfane benchmarks the IsProfane function
func BenchmarkIsProfane(b *testing.B) {
	benchmarks := []struct {
		name string
		text string
	}{
		{"CleanShort", cleanShort},
		{"CleanLong", cleanLong},
		{"ProfaneShort", profaneShort},
		{"ProfaneLong", profaneLong},
		{"Leetspeak", leetspeak},
		{"RepeatedChars", repeatedChars},
		{"MixedText", mixedText},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = IsProfane(bm.text)
			}
		})
	}
}

// BenchmarkCensor benchmarks the Censor function
func BenchmarkCensor(b *testing.B) {
	benchmarks := []struct {
		name string
		text string
	}{
		{"CleanShort", cleanShort},
		{"CleanLong", cleanLong},
		{"ProfaneShort", profaneShort},
		{"ProfaneLong", profaneLong},
		{"Leetspeak", leetspeak},
		{"RepeatedChars", repeatedChars},
		{"MixedText", mixedText},
	}

	modes := []struct {
		name string
		mode CensorMode
	}{
		{"All", CensorAll},
		{"KeepFirst", CensorKeepFirst},
		{"KeepFirstLast", CensorKeepFirstLast},
	}

	for _, bm := range benchmarks {
		for _, m := range modes {
			b.Run(bm.name+"/"+m.name, func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_ = Censor(bm.text, m.mode)
				}
			})
		}
	}
}

// BenchmarkNormalization benchmarks the text normalization
func BenchmarkNormalization(b *testing.B) {
	benchmarks := []struct {
		name string
		text string
	}{
		{"ASCII", "Hello, this is plain ASCII text!"},
		{"Unicode", "Héllo, thïs hás únïcödé chàracters!"},
		{"Leetspeak", "H3ll0, th1s h4s l33tsp34k!"},
		{"Cyrillic", "Неllo, thіs hаs Сyrilliс homоglyphs!"},
		{"ZeroWidth", "He\u200Bll\u200Co, th\u200Dis has ze\uFEFFro-width!"},
		{"Mixed", "H3ll\u200Bо, th1s hаs 3v3ryth1ng!"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = normalizeText(bm.text)
			}
		})
	}
}

// BenchmarkCollapseRepeats benchmarks the repeat collapsing
func BenchmarkCollapseRepeats(b *testing.B) {
	benchmarks := []struct {
		name string
		text string
	}{
		{"NoRepeats", "abcdefghijklmnop"},
		{"SomeRepeats", "heeeellooo wooooorld"},
		{"ManyRepeats", "aaaaaabbbbbbccccccddddddeeeeee"},
		{"Long", strings.Repeat("heeeello ", 100)},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, info := collapseRepeats(bm.text)
				_ = info
			}
		})
	}
}

// BenchmarkAhoCorasick benchmarks the raw Aho-Corasick search
func BenchmarkAhoCorasick(b *testing.B) {
	// Create a small automaton for controlled testing
	patterns := []string{"foo", "bar", "baz", "foobar", "foobaz"}
	ac := newAhoCorasick(patterns)

	benchmarks := []struct {
		name string
		text string
	}{
		{"NoMatch", "the quick brown fox jumps over the lazy dog"},
		{"SingleMatch", "the foo is here"},
		{"MultiMatch", "foo and bar and baz are here"},
		{"OverlapMatch", "foobar and foobaz"},
		{"Long", strings.Repeat("some text with foo and bar ", 100)},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = ac.SearchAll(bm.text)
			}
		})
	}
}

// BenchmarkAhoCorasickHasMatch benchmarks early exit on first match
func BenchmarkAhoCorasickHasMatch(b *testing.B) {
	benchmarks := []struct {
		name string
		text string
	}{
		{"CleanShort", cleanShort},
		{"CleanLong", cleanLong},
		{"ProfaneShort", profaneShort},
		{"ProfaneLong", profaneLong},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			// Normalize first for fair comparison
			nt := normalizeText(bm.text)
			collapsed, _ := collapseRepeats(nt.normalized)
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = profanityMatcher.HasMatch(collapsed)
			}
		})
	}
}

// BenchmarkBuildAutomaton benchmarks building the Aho-Corasick automaton
func BenchmarkBuildAutomaton(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = newAhoCorasick(profanityList)
	}
}

