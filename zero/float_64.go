package zero

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

// Float64 is a nullable float64. Zero input will be considered null.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
type Float64 struct {
	sql.NullFloat64
}

// NewFloat64 creates a new Float64
func NewFloat64(f float64, valid bool) Float64 {
	return Float64{
		NullFloat64: sql.NullFloat64{
			Float64: f,
			Valid:   valid,
		},
	}
}

// Float64From creates a new Float64 that will be null if zero.
func Float64From(f float64) Float64 {
	return NewFloat64(f, f != 0)
}

// Float64FromPtr creates a new Float64 that be null if f is nil.
func Float64FromPtr(f *float64) Float64 {
	if f == nil {
		return NewFloat64(0, false)
	}
	return NewFloat64(*f, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (f Float64) ValueOrZero() float64 {
	if !f.Valid {
		return 0
	}
	return f.Float64
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will be considered a null Float64.
func (f *Float64) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.Float64); err != nil {
		var typeError *json.UnmarshalTypeError
		if errors.As(err, &typeError) {
			// special case: accept string input
			if typeError.Value != "string" {
				return fmt.Errorf("zero: JSON input is invalid type (need float or string): %w", err)
			}
			var str string
			if err := json.Unmarshal(data, &str); err != nil {
				return fmt.Errorf("zero: couldn't unmarshal number string: %w", err)
			}
			n, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return fmt.Errorf("zero: couldn't convert string to float: %w", err)
			}
			f.Float64 = n
			f.Valid = n != 0
			return nil
		}
		return fmt.Errorf("zero: couldn't unmarshal JSON: %w", err)
	}

	f.Valid = f.Float64 != 0
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Float64 if the input is blank or zero.
// It will return an error if the input is not a float, blank, or "null".
func (f *Float64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		f.Valid = false
		return nil
	}
	var err error
	f.Float64, err = strconv.ParseFloat(string(text), 64)
	if err != nil {
		return fmt.Errorf("zero: couldn't unmarshal text: %w", err)
	}
	f.Valid = f.Float64 != 0
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Float64 is null.
func (f Float64) MarshalJSON() ([]byte, error) {
	n := f.Float64
	if !f.Valid {
		n = 0
	}
	if math.IsInf(f.Float64, 0) || math.IsNaN(f.Float64) {
		return nil, &json.UnsupportedValueError{
			Value: reflect.ValueOf(f.Float64),
			Str:   strconv.FormatFloat(f.Float64, 'g', -1, 64),
		}
	}
	return []byte(strconv.FormatFloat(n, 'f', -1, 64)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Float64 is null.
func (f Float64) MarshalText() ([]byte, error) {
	n := f.Float64
	if !f.Valid {
		n = 0
	}
	return []byte(strconv.FormatFloat(n, 'f', -1, 64)), nil
}

// SetValid changes this Float64's value and also sets it to be non-null.
func (f *Float64) SetValid(v float64) {
	f.Float64 = v
	f.Valid = true
}

// Ptr returns a poFloater to this Float64's value, or a nil poFloater if this Float64 is null.
func (f Float64) Ptr() *float64 {
	if !f.Valid {
		return nil
	}
	return &f.Float64
}

// IsZero returns true for null or zero Floats, for future omitempty support (Go 1.4?)
func (f Float64) IsZero() bool {
	return !f.Valid || f.Float64 == 0
}

// Equal returns true if both floats have the same value or are both either null or zero.
// Warning: calculations using floating point numbers can result in different ways
// the numbers are stored in memory. Therefore, this function is not suitable to
// compare the result of a calculation. Use this method only to check if the value
// has changed in comparison to some previous value.
func (f Float64) Equal(other Float64) bool {
	return f.ValueOrZero() == other.ValueOrZero()
}
