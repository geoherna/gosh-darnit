package goshdarnit

// ahoCorasick implements the Aho-Corasick string matching algorithm.
// It builds a finite state machine from a set of patterns and can match
// all patterns simultaneously in a single pass through the text.
type ahoCorasick struct {
	root     *acNode
	patterns []string
}

// acNode represents a node in the Aho-Corasick automaton.
type acNode struct {
	children   map[rune]*acNode
	fail       *acNode
	output     []int // indices into patterns slice
	depth      int
	isTerminal bool
}

// acMatch represents a match found by the Aho-Corasick automaton.
type acMatch struct {
	patternIndex int // index into patterns slice
	start        int // start position in text (byte offset)
	end          int // end position in text (byte offset, exclusive)
}

// newAhoCorasick creates a new Aho-Corasick automaton from the given patterns.
func newAhoCorasick(patterns []string) *ahoCorasick {
	ac := &ahoCorasick{
		root:     newACNode(0),
		patterns: patterns,
	}
	ac.buildTrie()
	ac.buildFailLinks()
	return ac
}

// newACNode creates a new node at the given depth.
func newACNode(depth int) *acNode {
	return &acNode{
		children: make(map[rune]*acNode),
		depth:    depth,
	}
}

// buildTrie builds the trie structure from the patterns.
func (ac *ahoCorasick) buildTrie() {
	for i, pattern := range ac.patterns {
		node := ac.root
		for _, r := range pattern {
			if _, exists := node.children[r]; !exists {
				node.children[r] = newACNode(node.depth + 1)
			}
			node = node.children[r]
		}
		node.isTerminal = true
		node.output = append(node.output, i)
	}
}

// buildFailLinks builds the failure links using BFS.
func (ac *ahoCorasick) buildFailLinks() {
	queue := make([]*acNode, 0)

	// Initialize fail links for depth-1 nodes to root
	for _, child := range ac.root.children {
		child.fail = ac.root
		queue = append(queue, child)
	}

	// BFS to build fail links for remaining nodes
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for r, child := range current.children {
			queue = append(queue, child)

			// Find the fail link for this child
			failNode := current.fail
			for failNode != nil {
				if next, exists := failNode.children[r]; exists {
					child.fail = next
					break
				}
				failNode = failNode.fail
			}
			if child.fail == nil {
				child.fail = ac.root
			}

			// Merge output from fail link (for overlapping patterns)
			if child.fail.output != nil {
				child.output = append(child.output, child.fail.output...)
			}
		}
	}
}

// Search finds all matches of the patterns in the text.
// The callback receives the pattern index, start byte position, and end byte position.
// The callback should return true to continue searching, flase to stop.
func (ac *ahoCorasick) Search(text string, callback func(match acMatch) bool) {
	node := ac.root
	bytePos := 0

	for i, r := range text {
		// Follow fail links until we find a match or reach root
		for node != ac.root {
			if _, exists := node.children[r]; exists {
				break
			}
			node = node.fail
		}

		// Try to transition on the current character
		if next, exists := node.children[r]; exists {
			node = next
		} else {
			node = ac.root
		}

		// Report any matches at this position
		if len(node.output) > 0 {
			runeLen := len(string(r))
			endPos := i + runeLen
			for _, patternIdx := range node.output {
				patternLen := len(ac.patterns[patternIdx])
				startPos := endPos - patternLen
				match := acMatch{
					patternIndex: patternIdx,
					start:        startPos,
					end:          endPos,
				}
				if !callback(match) {
					return
				}
			}
		}

		bytePos = i + len(string(r))
	}
	_ = bytePos // suppress unused warning
}

// SearchAll returns all matches found in the text.
func (ac *ahoCorasick) SearchAll(text string) []acMatch {
	var matches []acMatch
	ac.Search(text, func(match acMatch) bool {
		matches = append(matches, match)
		return true
	})
	return matches
}

// HasMatch returns true if any pattern matches in the text.
func (ac *ahoCorasick) HasMatch(text string) bool {
	found := false
	ac.Search(text, func(match acMatch) bool {
		found = true
		return false // stop after first match
	})
	return found
}
