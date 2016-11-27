package intSet

import (
	"bytes"
	"fmt"
)

// wordLength posesses the length of
// a single word in the set
const wordLength = 64

// IntSet is a container for representing
// a set of small unsigned integers
type IntSet struct {
	words []int64
}

// Has checks whether x is present in the set
func (s *IntSet) Has(x int) bool {
	word, bit := x/wordLength, x%wordLength

	return word < len(s.words) &&
		s.words[word]&(1<<uint(bit)) != 0
}

// Add adds a new word to the set
func (s *IntSet) Add(x int) {
	word, bit := x/wordLength, x%wordLength

	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}

	s.words[word] |= 1 << uint(bit)
}

// Union unites the both sets
func (s *IntSet) Union(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// String for user friendly output
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')

	for i, word := range s.words {
		if word == 0 {
			continue
		}

		for j := 0; j < wordLength; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() >= len("{") {
					buf.WriteByte(' ')
				}

				fmt.Fprintf(&buf, "%d", wordLength*i+j)
			}
		}
	}

	buf.WriteString(" }")
	return buf.String()
}

// Len returns the length of a set
func (s *IntSet) Len() (len uint) {
	for i, word := range s.words {
		if word == 0 {
			continue
		}

		for j := 0; j < wordLength; j++ {
			if s.words[i]&(1<<uint(j)) != 0 {
				len++
			}
		}
	}

	return
}

// Remove removes a word from the set
func (s *IntSet) Remove(x int) {
	if !s.Has(x) {
		return
	}

	for i, word := range s.words {
		if word == 0 {
			continue
		}

		for j := 0; j < wordLength; j++ {
			if wordLength*i+j == x {
				var mask int64 = 1 << uint(j)
				s.words[i] ^= mask
			}
		}
	}
}

// Clear deletes all the content of a set
func (s *IntSet) Clear() {
	s.words = s.words[:0]
}

// Copy returns a copy of a set
func (s *IntSet) Copy() *IntSet {
	t := IntSet{make([]int64, len(s.words), cap(s.words))}
	copy(t.words, s.words)

	return &t
}

// AddAll adds a sequence of numbers to the set
func (s *IntSet) AddAll(t ...int) {
	for _, x := range t {
		s.Add(x)
	}
}
