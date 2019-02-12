// Package lazy-search performs a type of approximate string matching.
//
// The pattern specifies the runes that must appear in the document and the
// order in which they must appear. They must not necessarily appear
// consecutively.
//
// Documents and patterns are []rune.
package lsearch

// Index stores the indexed structures but not the documents.
type Index struct {
	// d maps each rune to each document id to the indexes of these document that
	// contains the rune.
	d map[rune]map[int][]int
	// r maps each document id to the runes contained in each document.
	r map[int]map[rune]struct{}
}

// Unindex unindexes the document with the given id.
func (x *Index) Unindex(id int) {
	for r := range x.r[id] {
		delete(x.d[r], id)
		if len(x.d[r]) == 0 {
			delete(x.d, r)
		}
	}
	delete(x.r, id)
}

// Index assigns the given id to the given document, unindexes the old document
// if existed, and indexes the new one.
func (x *Index) Index(id int, doc []rune) {
	x.Unindex(id)
	if x.d == nil {
		x.d = make(map[rune]map[int][]int)
	}
	if x.r == nil {
		x.r = make(map[int]map[rune]struct{})
	}
	if x.r[id] == nil {
		x.r[id] = make(map[rune]struct{})
	}
	for i, r := range doc {
		if x.d[r] == nil {
			x.d[r] = make(map[int][]int)
		}
		x.d[r][id] = append(x.d[r][id], i)
		x.r[id][r] = struct{}{}
	}
}

// Search is a search of a pattern in an Index.
type Search struct {
	x *Index
	r map[int][]int
}

// Index returns the *Index it belongs to.
func (s Search) Index() *Index {
	return s.x
}

// Search returns a Search of the given pattern. The instance should not
// be used after the index is modified.
func (x *Index) Search(p []rune) Search {
	return Search{x, nil}.Search(p)
}

// Search returns a new Search of the pattern composed by concatenating
// p to the pattern of this Search. The instance should not be used after
// the index is modified.
func (s Search) Search(p []rune) Search {
	for _, r := range p {
		if s.r == nil {
			s.r = s.x.d[r]
			continue
		}
		sr := s.r
		s.r = make(map[int][]int)
		for id, js := range s.x.d[r] {
			is, ok := sr[id]
			if !ok {
				continue
			}
			for _, j := range js {
				for _, i := range is {
					if j > i {
						if s.r[id] == nil {
							s.r[id] = make([]int, 0, len(js))
						}
						s.r[id] = append(s.r[id], j)
						break
					}
				}
			}
		}
	}
	return s
}

// Docs appends the ids of the documents that match the current Search to ids
// without exceding its capacity and returns the resulting slice.
func (s Search) Docs(ids []int) []int {
	if s.r == nil {
		for id := range s.x.r {
			if len(ids) >= cap(ids) {
				break
			}
			ids = append(ids, id)
		}
		return ids
	}
	for id := range s.r {
		if len(ids) >= cap(ids) {
			break
		}
		ids = append(ids, id)
	}
	return ids
}
