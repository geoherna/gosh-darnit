
<p align="center">
<img width="500" height="508" alt="Emoji-angry-default" src="https://github.com/user-attachments/assets/3af0e6fe-9fdb-4e8f-9d41-5cf17881e43b" />
</p>

# gosh-darnit

[![CI](https://github.com/geoherna/gosh-darnit/actions/workflows/ci.yml/badge.svg)](https://github.com/geoherna/gosh-darnit/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/geoherna/gosh-darnit)](https://goreportcard.com/report/github.com/geoherna/gosh-darnit)
[![codecov](https://codecov.io/gh/geoherna/gosh-darnit/branch/main/graph/badge.svg)](https://codecov.io/gh/geoherna/gosh-darnit)
[![Go Reference](https://pkg.go.dev/badge/github.com/geoherna/gosh-darnit.svg)](https://pkg.go.dev/github.com/geoherna/gosh-darnit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A fast, efficient Go library for profanity detection and censorship.

## Features

- **Fast**: Uses the Aho-Corasick algorithm to match all patterns in a single pass
- **Smart word boundaries**: Prevents false positives like "bass", "analyst", "assist", "Scunthorpe"
- **Evasion resistant**: Handles common obfuscation techniques:
  - Leetspeak: `@ss`, `sh1t`, `fvck`, `a$$`
  - Unicode homoglyphs: Cyrillic, Greek, fullwidth characters
  - Zero-width characters: U+200B, U+200C, U+200D, U+FEFF
  - Repeated characters: `fuuuuck`, `shiiiit`
  - NFKC Unicode normalization
- **Flexible censoring**: Multiple modes for replacing profanity
- **Zero external dependencies**: Only uses Go standard library + `golang.org/x/text`

## Installation

```bash
go get github.com/geoherna/gosh-darnit
```

## Usage

### Basic Detection

```go
package main

import (
	"fmt"
	"github.com/geoherna/gosh-darnit"
)

func main() {
    // Check if text contains profanity
    if goshdarnit.IsProfane("What the fuck?") {
        fmt.Println("Profanity detected!")
    }

    // Find which words matched
    words := goshdarnit.FindProfanity("This is some shit")
    fmt.Println("Found:", words) // ["shit"]
}
```

### Censoring

```go
package main

import (
	"fmt"
	"github.com/geoherna/gosh-darnit"
)

func main() {
    text := "What the fuck is this shit?"

    // Replace all characters with asterisks
    fmt.Println(goshdarnit.Censor(text, goshdarnit.CensorAll))
    // Output: "What the **** is this ****?"

    // Keep first character visible
    fmt.Println(goshdarnit.Censor(text, goshdarnit.CensorKeepFirst))
    // Output: "What the f*** is this s***?"

    // Keep first and last characters visible
    fmt.Println(goshdarnit.Censor(text, goshdarnit.CensorKeepFirstLast))
    // Output: "What the f**k is this s**t?"
}
```

### Evasion Detection

The library automatically handles common evasion techniques:

```go
// Leetspeak
goshdarnit.IsProfane("@ss")      // true (@ -> a)
goshdarnit.IsProfane("sh1t")     // true (1 -> i)
goshdarnit.IsProfane("fvck")     // true (v -> u)
goshdarnit.IsProfane("a$$")      // true ($ -> s)

// Repeated characters
goshdarnit.IsProfane("fuuuuck")  // true
goshdarnit.IsProfane("shiiiit")  // true

// Unicode homoglyphs (Cyrillic 'а' looks like Latin 'a')
goshdarnit.IsProfane("аss")      // true
```

### False Positive Prevention

Word boundary detection prevents common false positives:

```go
goshdarnit.IsProfane("The bass is great")     // false
goshdarnit.IsProfane("She's an analyst")      // false
goshdarnit.IsProfane("I need to assist you")  // false
goshdarnit.IsProfane("Scunthorpe is a town")  // false
goshdarnit.IsProfane("The shitake mushrooms") // false
goshdarnit.IsProfane("Assess the situation")  // false
goshdarnit.IsProfane("Classic movie")         // false
```

## API Reference

### Functions

| Function | Description |
|----------|-------------|
| `IsProfane(text string) bool` | Returns true if text contains profanity |
| `ContainsProfanity(text string) bool` | Alias for `IsProfane` |
| `Censor(text string, mode CensorMode) string` | Replaces profanity with asterisks |
| `CensorWithDefault(text string) string` | Censors with `CensorAll` mode |
| `FindProfanity(text string) []string` | Returns list of matched profane words |

### Censor Modes

| Mode | Example | Description |
|------|---------|-------------|
| `CensorAll` | `****` | Replace all characters |
| `CensorKeepFirst` | `f***` | Keep first character visible |
| `CensorKeepFirstLast` | `f**k` | Keep first and last characters visible |

## Performance

Benchmarks on Apple M4 Max:

| Benchmark | Time | Allocations |
|-----------|------|-------------|
| CleanShort | ~766ns | 8 allocs |
| ProfaneShort | ~839ns | 9 allocs |
| Leetspeak | ~847ns | 9 allocs |
| RepeatedChars | ~1.0µs | 11 allocs |
| MixedText | ~2.5µs | 14 allocs |

Run benchmarks yourself:

```bash
go test -bench=. -benchmem
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details.

