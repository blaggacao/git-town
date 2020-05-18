package helpers_test

import (
	"testing"

	"github.com/git-town/git-town/test/helpers"
	"github.com/stretchr/testify/assert"
)

func TestOrderedStringSet(t *testing.T) {
	set := helpers.NewOrderedStringSet("one")
	set = set.Add("two")
	set = set.Add("two")
	set = set.Add("two", "three")
	assert.Equal(t, []string{"one", "two", "three"}, set.Slice())
	assert.Equal(t, "one, two, three", set.String())
}
