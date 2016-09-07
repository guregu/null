package zero

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/nullbio/null.v5/convert"
)

type NullFloat32 struct {
	Float32 float32
	Valid   bool
}

// Float32 is a nullable float32. Zero input will be considered null.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
type Float32 struct {
	NullFloat32
}

// NewFloat32 creates a new Float32
func NewFloat32(f float32, valid bool) Float32 {
	return Float32{
		NullFloat32: NullFloat32{
			Float32: f,
			Valid:   valid,
		},
	}
}

// Float32From creates a new Float32 that will be null if zero.
func Float32From(f float32) Float32 {
	return NewFloat32(f, f != 0)
}

// Float32FromPtr creates a new Float32 that be null if f is nil.
func Float32FromPtr(f *float32) Float32 {
	if f == nil {
		return NewFloat32(0, false)
	}
	return NewFloat32(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will be considered a null Float32.
// It also supports unmarshalling a sql.NullFloat32.
func (f *Float32) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		f.Float32 = float32(x)
	case map[string]interface{}:
		err = json.Unmarshal(data, &f.NullFloat32)
	case nil:
		f.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type zero.Float32", reflect.TypeOf(v).Name())
	}
	f.Valid = (err == nil) && (f.Float32 != 0)
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Float32 if the input is a blank, zero, or not a float.
// It will return an error if the input is not a float, blank, or "null".
func (f *Float32) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		f.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseFloat(string(text), 32)
	f.Float32 = float32(res)
	f.Valid = (err == nil) && (f.Float32 != 0)
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Float32 is null.
func (f Float32) MarshalJSON() ([]byte, error) {
	n := f.Float32
	if !f.Valid {
		n = 0
	}
	return []byte(strconv.FormatFloat(float64(n), 'f', -1, 32)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Float32 is null.
func (f Float32) MarshalText() ([]byte, error) {
	n := f.Float32
	if !f.Valid {
		n = 0
	}
	return []byte(strconv.FormatFloat(float64(n), 'f', -1, 32)), nil
}

// SetValid changes this Float32's value and also sets it to be non-null.
func (f *Float32) SetValid(v float32) {
	f.Float32 = v
	f.Valid = true
}

// Ptr returns a poFloater to this Float32's value, or a nil poFloater if this Float32 is null.
func (f Float32) Ptr() *float32 {
	if !f.Valid {
		return nil
	}
	return &f.Float32
}

// IsZero returns true for null or zero Float32's, for future omitempty support (Go 1.4?)
func (f Float32) IsZero() bool {
	return !f.Valid || f.Float32 == 0
}

// Scan implements the Scanner interface.
func (n *NullFloat32) Scan(value interface{}) error {
	if value == nil {
		n.Float32, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Float32, value)
}

// Value implements the driver Valuer interface.
func (n NullFloat32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return float64(n.Float32), nil
}
