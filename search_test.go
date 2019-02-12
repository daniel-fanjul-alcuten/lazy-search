package lsearch

import (
	"fmt"
	"sort"
	"testing"
)

func TestSearch(t *testing.T) {
	x := Index{}
	x.Index(1, []rune("foo"))
	x.Index(2, []rune("bar"))
	x.Index(3, []rune("baz"))
	type N struct {
		p []rune
		s string
	}
	for _, n := range []N{
		{[]rune(""), "[1 2 3]"},
		{[]rune("a"), "[2 3]"},
		{[]rune("b"), "[2 3]"},
		{[]rune("ba"), "[2 3]"},
		{[]rune("bab"), "[]"},
		{[]rune("bar"), "[2]"},
		{[]rune("baz"), "[3]"},
		{[]rune("br"), "[2]"},
		{[]rune("bz"), "[3]"},
		{[]rune("f"), "[1]"},
		{[]rune("fa"), "[]"},
		{[]rune("faa"), "[]"},
		{[]rune("fo"), "[1]"},
		{[]rune("foo"), "[1]"},
		{[]rune("fr"), "[]"},
		{[]rune("fz"), "[]"},
		{[]rune("o"), "[1]"},
		{[]rune("r"), "[2]"},
		{[]rune("z"), "[3]"},
		{[]rune{}, "[1 2 3]"},
		{nil, "[1 2 3]"},
	} {
		docs := x.Search(n.p).Docs(make([]int, 0, 3))
		sort.Ints(docs)
		if s := fmt.Sprint(docs); s != n.s {
			t.Error(n.p, n.s, s)
		}
	}
}
