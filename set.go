package uuid

import (
	"bytes"
	"fmt"
	"sort"
)

func Less(a, b UUID) bool {
	return bytes.Compare(a[:], b[:]) < 0
	// for i := len(a) - 1; i >= 0; i-- {
	// 	if a[i] < b[i] {
	// 		return true
	// 	}
	// 	if a[i] > b[i] {
	// 		return false
	// 	}
	// }
	// return false
}

type Set map[UUID]struct{}

func SetFromSlice(s []UUID) Set {
	set := make(Set)
	set.AddSlice(s)
	return set
}

func (set Set) String() string {
	return fmt.Sprintf("set%v", set.SortedSlice())
}

// GetOne returns a random UUID from the set or Nil if the set is empty.
// Most useful to get the only UUID in a set of size one.
func (set Set) GetOne() UUID {
	for id := range set {
		return id
	}
	return Nil
}

func (set Set) Slice() []UUID {
	s := make([]UUID, len(set))
	i := 0
	for id := range set {
		s[i] = id
		i++
	}
	return s
}

func (set Set) SortedSlice() []UUID {
	s := set.Slice()
	sort.Slice(s, func(i, j int) bool { return Less(s[i], s[j]) })
	return s
}

func (set Set) SortedStringSlice() []string {
	s := make([]string, len(set))
	for i, id := range set.SortedSlice() {
		s[i] = id.String()
	}
	return s
}

func (set Set) AddSlice(s []UUID) {
	for _, id := range s {
		set[id] = struct{}{}
	}
}

func (set Set) AddSet(other Set) {
	for id := range other {
		set[id] = struct{}{}
	}
}

func (set Set) Add(id UUID) {
	set[id] = struct{}{}
}

func (set Set) Has(id UUID) bool {
	_, has := set[id]
	return has
}

func (set Set) Delete(id UUID) {
	delete(set, id)
}

func (set Set) DeleteSlice(s []UUID) {
	for _, id := range s {
		delete(set, id)
	}
}

func (set Set) DeleteSet(other Set) {
	for id := range other {
		delete(set, id)
	}
}

func (set Set) Clone() Set {
	clone := make(Set)
	clone.AddSet(set)
	return clone
}

func (set Set) Diff(other Set) Set {
	diff := make(Set)
	for id := range set {
		if !other.Has(id) {
			diff.Add(id)
		}
	}
	for id := range other {
		if !set.Has(id) {
			diff.Add(id)
		}
	}
	return diff
}
