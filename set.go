package uuid

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Set of UUIDs
// Is a map[UUID]struct{} underneath.
// Implements the database/sql.Scanner and database/sql/driver.Valuer interfaces
// with the nil map value used as SQL NULL
type Set map[UUID]struct{}

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

func (set Set) Slice() Slice {
	s := make(Slice, len(set))
	i := 0
	for id := range set {
		s[i] = id
		i++
	}
	return s
}

func (set Set) SortedSlice() Slice {
	s := set.Slice()
	s.Sort()
	return s
}

func (set Set) AddSlice(s Slice) {
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

func (set Set) Contains(id UUID) bool {
	_, has := set[id]
	return has
}

func (set Set) Delete(id UUID) {
	delete(set, id)
}

func (set Set) DeleteAll() {
	for id := range set {
		delete(set, id)
	}
}

func (set Set) DeleteSlice(s Slice) {
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
		if !other.Contains(id) {
			diff.Add(id)
		}
	}
	for id := range other {
		if !set.Contains(id) {
			diff.Add(id)
		}
	}
	return diff
}

// Scan implements the database/sql.Scanner interface
// with the nil map value used as SQL NULL.
// Does *set = make(Set) if *set == nil
// so it can be used with an not initialized Set variable
func (set *Set) Scan(value interface{}) (err error) {
	switch x := value.(type) {
	case string:
		return set.scanBytes([]byte(x))

	case []byte:
		return set.scanBytes(x)

	case nil:
		*set = nil
		return nil
	}

	return errors.Errorf("Can't scan value '%#v' of type %T as uuid.Set", value, value)
}

func (set *Set) scanBytes(src []byte) (err error) {
	if src == nil {
		*set = nil
		return nil
	}

	if len(src) < 2 || src[0] != '{' || src[len(src)-1] != '}' {
		return errors.Errorf("Can't parse %#v as uuid.Set", string(src))
	}

	ids := make(Slice, 0, 16)

	elements := bytes.Split(src[1:len(src)-1], []byte{','})
	for _, elem := range elements {
		elem = bytes.Trim(elem, `'"`)
		id, err := FromString(string(elem))
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}

	if *set == nil {
		*set = make(Set)
	} else {
		set.DeleteAll()
	}
	set.AddSlice(ids)

	return nil
}

// Value implements the driver database/sql/driver.Valuer interface
// with the nil map value used as SQL NULL
func (set Set) Value() (driver.Value, error) {
	if set == nil {
		return nil, nil
	}

	var b strings.Builder
	b.WriteByte('{')
	for i, id := range set.SortedSlice() {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteByte('"')
		b.WriteString(id.String())
		b.WriteByte('"')
	}
	b.WriteByte('}')

	return b.String(), nil
}
