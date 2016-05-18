package zero

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/pobri19/null-extended/convert"
)

type NullInt8 struct {
	Int8  int8
	Valid bool
}

// Int8 is a nullable int8.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
type Int8 struct {
	NullInt8
}

// NewInt8 creates a new Int8
func NewInt8(i int8, valid bool) Int8 {
	return Int8{
		NullInt8: NullInt8{
			Int8:  i,
			Valid: valid,
		},
	}
}

// Int8From creates a new Int8 that will be null if zero.
func Int8From(i int8) Int8 {
	return NewInt8(i, i != 0)
}

// Int8FromPtr creates a new Int8 that be null if i is nil.
func Int8FromPtr(i *int8) Int8 {
	if i == nil {
		return NewInt8(0, false)
	}
	n := NewInt8(*i, true)
	return n
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will be considered a null Int8.
// It also supports unmarshalling a sql.NullInt8.
func (i *Int8) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int8, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Int8)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullInt8)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type zero.Int8", reflect.TypeOf(v).Name())
	}
	i.Valid = (err == nil) && (i.Int8 != 0)
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int8 if the input is a blank, zero, or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int8) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseInt(string(text), 10, 8)
	i.Int8 = int8(res)
	i.Valid = (err == nil) && (i.Int8 != 0)
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode 0 if this Int8 is null.
func (i Int8) MarshalJSON() ([]byte, error) {
	n := i.Int8
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(int64(n), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Int8 is null.
func (i Int8) MarshalText() ([]byte, error) {
	n := i.Int8
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(int64(n), 10)), nil
}

// SetValid changes this Int8's value and also sets it to be non-null.
func (i *Int8) SetValid(n int8) {
	i.Int8 = n
	i.Valid = true
}

// Ptr returns a pointer to this Int8's value, or a nil pointer if this Int8 is null.
func (i Int8) Ptr() *int8 {
	if !i.Valid {
		return nil
	}
	return &i.Int8
}

// IsZero returns true for null or zero Int8s, for future omitempty support (Go 1.4?)
func (i Int8) IsZero() bool {
	return !i.Valid || i.Int8 == 0
}

// Scan implements the Scanner interface.
func (n *NullInt8) Scan(value interface{}) error {
	if value == nil {
		n.Int8, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Int8, value)
}

// Value implements the driver Valuer interface.
func (n NullInt8) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int8, nil
}
