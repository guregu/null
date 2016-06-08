package zero

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/nullbio/null.v4/convert"
)

type NullUint8 struct {
	Uint8 uint8
	Valid bool
}

// Uint8 is a nullable uint8.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
type Uint8 struct {
	NullUint8
}

// NewUint8 creates a new Uint8
func NewUint8(i uint8, valid bool) Uint8 {
	return Uint8{
		NullUint8: NullUint8{
			Uint8: i,
			Valid: valid,
		},
	}
}

// Uint8From creates a new Uint8 that will be null if zero.
func Uint8From(i uint8) Uint8 {
	return NewUint8(i, i != 0)
}

// Uint8FromPtr creates a new Uint8 that be null if i is nil.
func Uint8FromPtr(i *uint8) Uint8 {
	if i == nil {
		return NewUint8(0, false)
	}
	n := NewUint8(*i, true)
	return n
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will be considered a null Uint8.
// It also supports unmarshalling a sql.NullUint8.
func (i *Uint8) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to uint8, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Uint8)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullUint8)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type zero.Uint8", reflect.TypeOf(v).Name())
	}
	i.Valid = (err == nil) && (i.Uint8 != 0)
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint8 if the input is a blank, zero, or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Uint8) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseUint(string(text), 10, 8)
	i.Uint8 = uint8(res)
	i.Valid = (err == nil) && (i.Uint8 != 0)
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode 0 if this Uint8 is null.
func (i Uint8) MarshalJSON() ([]byte, error) {
	n := i.Uint8
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatUint(uint64(n), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Uint8 is null.
func (i Uint8) MarshalText() ([]byte, error) {
	n := i.Uint8
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatUint(uint64(n), 10)), nil
}

// SetValid changes this Uint8's value and also sets it to be non-null.
func (i *Uint8) SetValid(n uint8) {
	i.Uint8 = n
	i.Valid = true
}

// Ptr returns a pointer to this Uint8's value, or a nil pointer if this Uint8 is null.
func (i Uint8) Ptr() *uint8 {
	if !i.Valid {
		return nil
	}
	return &i.Uint8
}

// IsZero returns true for null or zero Uint8s, for future omitempty support (Go 1.4?)
func (i Uint8) IsZero() bool {
	return !i.Valid || i.Uint8 == 0
}

// Scan implements the Scanner interface.
func (n *NullUint8) Scan(value interface{}) error {
	if value == nil {
		n.Uint8, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Uint8, value)
}

// Value implements the driver Valuer interface.
func (n NullUint8) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return int64(n.Uint8), nil
}
