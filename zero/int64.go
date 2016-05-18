package zero

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// Int64 is a nullable int64.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
type Int64 struct {
	sql.NullInt64
}

// NewInt64 creates a new Int64
func NewInt64(i int64, valid bool) Int64 {
	return Int64{
		NullInt64: sql.NullInt64{
			Int64: i,
			Valid: valid,
		},
	}
}

// Int64From creates a new Int64 that will be null if zero.
func Int64From(i int64) Int64 {
	return NewInt64(i, i != 0)
}

// Int64FromPtr creates a new Int64 that be null if i is nil.
func Int64FromPtr(i *int64) Int64 {
	if i == nil {
		return NewInt64(0, false)
	}
	n := NewInt64(*i, true)
	return n
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will be considered a null Int64.
// It also supports unmarshalling a sql.NullInt64.
func (i *Int64) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int64, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Int64)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullInt64)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type zero.Int64", reflect.TypeOf(v).Name())
	}
	i.Valid = (err == nil) && (i.Int64 != 0)
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int64 if the input is a blank, zero, or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	i.Int64, err = strconv.ParseInt(string(text), 10, 64)
	i.Valid = (err == nil) && (i.Int64 != 0)
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode 0 if this Int64 is null.
func (i Int64) MarshalJSON() ([]byte, error) {
	n := i.Int64
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(n, 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Int64 is null.
func (i Int64) MarshalText() ([]byte, error) {
	n := i.Int64
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(n, 10)), nil
}

// SetValid changes this Int64's value and also sets it to be non-null.
func (i *Int64) SetValid(n int64) {
	i.Int64 = n
	i.Valid = true
}

// Ptr returns a pointer to this Int64's value, or a nil pointer if this Int64 is null.
func (i Int64) Ptr() *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

// IsZero returns true for null or zero Int64s, for future omitempty support (Go 1.4?)
func (i Int64) IsZero() bool {
	return !i.Valid || i.Int64 == 0
}
