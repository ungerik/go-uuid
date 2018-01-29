package uuid

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
)

// NullUUID can be used with the standard sql package to represent a
// UUID value that can be NULL in the database
type NullUUID struct {
	UUID  UUID
	Valid bool
}

// NullUUIDFrom creates a new valid NullUUID
func NullUUIDFrom(u UUID) NullUUID {
	return NewNullUUID(u, true)
}

// NullUUIDFromString creates a new valid NullUUID
func NullUUIDFromString(s string) (n NullUUID, err error) {
	n.UUID, err = FromString(s)
	if err != nil {
		return NullUUID{}, err
	}
	n.Valid = true
	return n, nil
}

// NullUUIDFromBytes creates a new valid NullUUID
func NullUUIDFromBytes(s []byte) (n NullUUID, err error) {
	n.UUID, err = FromBytes(s)
	if err != nil {
		return NullUUID{}, err
	}
	n.Valid = true
	return n, nil
}

// NullUUIDFromPtr creates a new NullUUID that be null if ptr is nil.
func NullUUIDFromPtr(ptr *UUID) NullUUID {
	if ptr == nil {
		return NullUUID{}
	}
	return NewNullUUID(*ptr, true)
}

// NewNullUUID creates a new NullUUID
func NewNullUUID(u UUID, valid bool) NullUUID {
	return NullUUID{
		UUID:  u,
		Valid: valid,
	}
}

// SetValid changes this NullUUID's value and also sets it to be non-null.
func (u *NullUUID) SetValid(v UUID) {
	u.UUID = v
	u.Valid = true
}

// Ptr returns a pointer to this NullUUID's value, or a nil pointer if this NullUUID is null.
func (u NullUUID) Ptr() *UUID {
	if !u.Valid {
		return nil
	}
	return &u.UUID
}

// Value implements the driver.Valuer interface.
func (u NullUUID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	// Delegate to UUID Value function
	return u.UUID.Value()
}

// Scan implements the sql.Scanner interface.
func (u *NullUUID) Scan(src interface{}) error {
	if src == nil {
		u.UUID, u.Valid = Nil, false
		return nil
	}

	// Delegate to UUID Scan function
	u.Valid = true
	return u.UUID.Scan(src)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input does not produce a null UUID.
// It also supports unmarshalling a sql.NullString.
func (u *NullUUID) UnmarshalJSON(data []byte) (err error) {
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		u.UUID, err = FromString(x)
	case map[string]interface{}:
		var n sql.NullString
		err = json.Unmarshal(data, &n)
		if n.Valid {
			u.UUID, err = FromString(n.String)
		}
	case nil:
		u.UUID = Nil
		u.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type uuid.NullString", reflect.TypeOf(v).Name())
	}
	u.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if Valid == false.
func (u NullUUID) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(u.UUID.String())
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string when this String is null.
func (u NullUUID) MarshalText() ([]byte, error) {
	if !u.Valid {
		return []byte{}, nil
	}
	return []byte(u.UUID.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null String if the input is a blank string.
func (u *NullUUID) UnmarshalText(text []byte) (err error) {
	if len(text) == 0 {
		u.UUID = Nil
		u.Valid = false
		return nil
	}
	u.UUID, err = FromBytes(text)
	u.Valid = err == nil
	return err
}

// String returns the UUID as string if Valid == true,
// or "null" if Valid == false.
func (u NullUUID) String() string {
	if !u.Valid {
		return "null"
	}
	return u.UUID.String()
}
