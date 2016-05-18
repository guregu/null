package zero

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/pobri19/null-extended/convert"
)

// NullUint64 is a replica of sql.NullInt64 for uint64 types.
type NullUint64 struct {
	Uint64 uint64
	Valid  bool
}

// Uint64 is an nullable uint64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Uint64 struct {
	NullUint64
}

// NewUint64 creates a new Uint64
func NewUint64(i uint64, valid bool) Uint64 {
	return Uint64{
		NullUint64: NullUint64{
			Uint64: i,
			Valid:  valid,
		},
	}
}

// Uint64From creates a new Uint64 that will be null if zero.
func Uint64From(i uint64) Uint64 {
	return NewUint64(i, i != 0)
}

// Uint64FromPtr creates a new Uint64 that be null if i is nil.
func Uint64FromPtr(i *uint64) Uint64 {
	if i == nil {
		return NewUint64(0, false)
	}
	n := NewUint64(*i, true)
	return n
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will be considered a null Uint64.
// It also supports unmarshalling a sql.NullUint64.
func (i *Uint64) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to uint64, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Uint64)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullUint64)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type zero.Uint64", reflect.TypeOf(v).Name())
	}
	i.Valid = (err == nil) && (i.Uint64 != 0)
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint64 if the input is a blank, zero, or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Uint64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	i.Uint64, err = strconv.ParseUint(string(text), 10, 64)
	i.Valid = (err == nil) && (i.Uint64 != 0)
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode 0 if this Uint64 is null.
func (i Uint64) MarshalJSON() ([]byte, error) {
	n := i.Uint64
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatUint(n, 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Uint64 is null.
func (i Uint64) MarshalText() ([]byte, error) {
	n := i.Uint64
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatUint(n, 10)), nil
}

// SetValid changes this Uint64's value and also sets it to be non-null.
func (i *Uint64) SetValid(n uint64) {
	i.Uint64 = n
	i.Valid = true
}

// Ptr returns a pointer to this Uint64's value, or a nil pointer if this Uint64 is null.
func (i Uint64) Ptr() *uint64 {
	if !i.Valid {
		return nil
	}
	return &i.Uint64
}

// IsZero returns true for null or zero Uint64s, for future omitempty support (Go 1.4?)
func (i Uint64) IsZero() bool {
	return !i.Valid || i.Uint64 == 0
}

// Scan implements the Scanner interface.
func (n *NullUint64) Scan(value interface{}) error {
	if value == nil {
		n.Uint64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Uint64, value)
}

// Value implements the driver Valuer interface.
func (n NullUint64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Uint64, nil
}
