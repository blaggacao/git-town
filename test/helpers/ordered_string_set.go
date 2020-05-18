package helpers

import "strings"

// OrderedStringSet is a Set for strings
// that provides the accumulated strings in the order they were received.
// The zero value is a fully functional empty set.
type OrderedStringSet struct {
	elements []string
}

// NewOrderedStringSet provides instances of OrderedStringSet
// populated with the given elements.
func NewOrderedStringSet(elements ...string) OrderedStringSet {
	return OrderedStringSet{elements: elements}
}

// Add provides a new Set with the given element added.
// If the element already exists, it re-uses the existing element and does not append a new one.
func (set OrderedStringSet) Add(elements ...string) OrderedStringSet {
	result := set
	for e := range elements {
		if !set.Contains(elements[e]) {
			result = OrderedStringSet{elements: append(result.elements, elements[e])}
		}
	}
	return result
}

// Contains indicates whether this Set contains the given string.
func (set OrderedStringSet) Contains(text string) bool {
	for _, element := range set.elements {
		if text == element {
			return true
		}
	}
	return false
}

// Slice provides the elements of this set in the order they were received.
func (set OrderedStringSet) Slice() []string {
	return set.elements
}

func (set OrderedStringSet) String() string {
	return strings.Join(set.elements, ", ")
}
