package zero

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/nullbio/null.v4/convert"
)

type NullInt32 struct {
	Int32 int32
	Valid bool
}

// Int32 is a nullable int32.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
type Int32 struct {
	NullInt32
}

// NewInt32 creates a new Int32
func NewInt32(i int32, valid bool) Int32 {
	return Int32{
		NullInt32: NullInt32{
			Int32: i,
			Valid: valid,
		},
	}
}

// Int32From creates a new Int32 that will be null if zero.
func Int32From(i int32) Int32 {
	return NewInt32(i, i != 0)
}

// Int32FromPtr creates a new Int32 that be null if i is nil.
func Int32FromPtr(i *int32) Int32 {
	if i == nil {
		return NewInt32(0, false)
	}
	n := NewInt32(*i, true)
	return n
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will be considered a null Int32.
// It also supports unmarshalling a sql.NullInt32.
func (i *Int32) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int32, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Int32)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullInt32)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type zero.Int32", reflect.TypeOf(v).Name())
	}
	i.Valid = (err == nil) && (i.Int32 != 0)
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int32 if the input is a blank, zero, or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int32) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseInt(string(text), 10, 32)
	i.Int32 = int32(res)
	i.Valid = (err == nil) && (i.Int32 != 0)
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode 0 if this Int32 is null.
func (i Int32) MarshalJSON() ([]byte, error) {
	n := i.Int32
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(int64(n), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Int32 is null.
func (i Int32) MarshalText() ([]byte, error) {
	n := i.Int32
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(int64(n), 10)), nil
}

// SetValid changes this Int32's value and also sets it to be non-null.
func (i *Int32) SetValid(n int32) {
	i.Int32 = n
	i.Valid = true
}

// Ptr returns a pointer to this Int32's value, or a nil pointer if this Int32 is null.
func (i Int32) Ptr() *int32 {
	if !i.Valid {
		return nil
	}
	return &i.Int32
}

// IsZero returns true for null or zero Int32s, for future omitempty support (Go 1.4?)
func (i Int32) IsZero() bool {
	return !i.Valid || i.Int32 == 0
}

// Scan implements the Scanner interface.
func (n *NullInt32) Scan(value interface{}) error {
	if value == nil {
		n.Int32, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Int32, value)
}

// Value implements the driver Valuer interface.
func (n NullInt32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return int64(n.Int32), nil
}
