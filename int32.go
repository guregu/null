package null

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Int32 is an nullable int32.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Int32 struct {
	sql.NullInt32
}

// NewInt32 creates a new Int32
func NewInt32(i int32, valid bool) Int32 {
	return Int32{
		NullInt32: sql.NullInt32{
			Int32: i,
			Valid: valid,
		},
	}
}

// Int32From creates a new Int32 that will always be valid.
func Int32From(i int32) Int32 {
	return NewInt32(i, true)
}

// Int32FromPtr creates a new Int32 that be null if i is nil.
func Int32FromPtr(i *int32) Int32 {
	if i == nil {
		return NewInt32(0, false)
	}
	return NewInt32(*i, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Int32) ValueOrZero() int32 {
	if !i.Valid {
		return 0
	}
	return i.Int32
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number, string, and null input.
// 0 will not be considered a null Int32.
func (i *Int32) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		i.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &i.Int32); err != nil {
		var typeError *json.UnmarshalTypeError
		if errors.As(err, &typeError) {
			// special case: accept string input
			if typeError.Value != "string" {
				return fmt.Errorf("null: JSON input is invalid type (need int or string): %w", err)
			}
			var str string
			if err := json.Unmarshal(data, &str); err != nil {
				return fmt.Errorf("null: couldn't unmarshal number string: %w", err)
			}
			n, err := strconv.ParseInt(str, 10, 32)
			if err != nil {
				return fmt.Errorf("null: couldn't convert string to int: %w", err)
			}
			i.Int32 = int32(n)
			i.Valid = true
			return nil
		}
		return fmt.Errorf("null: couldn't unmarshal JSON: %w", err)
	}

	i.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int32 if the input is blank.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int32) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	i64, err := strconv.ParseInt(string(text), 10, 32)
	if err != nil {
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}
	i.Int32 = int32(i64)
	i.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Int32 is null.
func (i Int32) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(int64(i.Int32), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Int32 is null.
func (i Int32) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int32), 10)), nil
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

// IsZero returns true for invalid Int32s, for future omitempty support (Go 1.4?)
// A non-null Int32 with a 0 value will not be considered zero.
func (i Int32) IsZero() bool {
	return !i.Valid
}

// Equal returns true if both int32s have the same value or are both null.
func (i Int32) Equal(other Int32) bool {
	return i.Valid == other.Valid && (!i.Valid || i.Int32 == other.Int32)
}
