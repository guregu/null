package zero

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/nullbio/null.v4/convert"
)

type NullUint struct {
	Uint  uint
	Valid bool
}

// Uint is a nullable uint.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
type Uint struct {
	NullUint
}

// NewUint creates a new Uint
func NewUint(i uint, valid bool) Uint {
	return Uint{
		NullUint: NullUint{
			Uint:  i,
			Valid: valid,
		},
	}
}

// UintFrom creates a new Uint that will be null if zero.
func UintFrom(i uint) Uint {
	return NewUint(i, i != 0)
}

// UintFromPtr creates a new Uint that be null if i is nil.
func UintFromPtr(i *uint) Uint {
	if i == nil {
		return NewUint(0, false)
	}
	n := NewUint(*i, true)
	return n
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will be considered a null Uint.
// It also supports unmarshalling a sql.NullUint.
func (i *Uint) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to uint, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Uint)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullUint)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type zero.Uint", reflect.TypeOf(v).Name())
	}
	i.Valid = (err == nil) && (i.Uint != 0)
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint if the input is a blank, zero, or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Uint) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseUint(string(text), 10, 0)
	i.Uint = uint(res)
	i.Valid = (err == nil) && (i.Uint != 0)
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode 0 if this Uint is null.
func (i Uint) MarshalJSON() ([]byte, error) {
	n := i.Uint
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatUint(uint64(n), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Uint is null.
func (i Uint) MarshalText() ([]byte, error) {
	n := i.Uint
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatUint(uint64(n), 10)), nil
}

// SetValid changes this Uint's value and also sets it to be non-null.
func (i *Uint) SetValid(n uint) {
	i.Uint = n
	i.Valid = true
}

// Ptr returns a pointer to this Uint's value, or a nil pointer if this Uint is null.
func (i Uint) Ptr() *uint {
	if !i.Valid {
		return nil
	}
	return &i.Uint
}

// IsZero returns true for null or zero Uints, for future omitempty support (Go 1.4?)
func (i Uint) IsZero() bool {
	return !i.Valid || i.Uint == 0
}

// Scan implements the Scanner interface.
func (n *NullUint) Scan(value interface{}) error {
	if value == nil {
		n.Uint, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Uint, value)
}

// Value implements the driver Valuer interface.
func (n NullUint) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return int64(n.Uint), nil
}
