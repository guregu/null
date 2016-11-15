package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gopkg.in/nullbio/null.v6/convert"
)

// Bool is a nullable bool.
type Bool struct {
	Bool  bool
	Valid bool
}

// NewBool creates a new Bool
func NewBool(b bool, valid bool) Bool {
	return Bool{
		Bool:  b,
		Valid: valid,
	}
}

// BoolFrom creates a new Bool that will always be valid.
func BoolFrom(b bool) Bool {
	return NewBool(b, true)
}

// BoolFromPtr creates a new Bool that will be null if f is nil.
func BoolFromPtr(b *bool) Bool {
	if b == nil {
		return NewBool(false, false)
	}
	return NewBool(*b, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Bool) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		b.Bool = false
		b.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &b.Bool); err != nil {
		return err
	}

	b.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Bool) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		b.Valid = false
		return nil
	}

	str := string(text)
	switch str {
	case "true":
		b.Bool = true
	case "false":
		b.Bool = false
	default:
		b.Valid = false
		return errors.New("invalid input:" + str)
	}
	b.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return NullBytes, nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// MarshalText implements encoding.TextMarshaler.
func (b Bool) MarshalText() ([]byte, error) {
	if !b.Valid {
		return []byte{}, nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// SetValid changes this Bool's value and also sets it to be non-null.
func (b *Bool) SetValid(v bool) {
	b.Bool = v
	b.Valid = true
}

// Ptr returns a pointer to this Bool's value, or a nil pointer if this Bool is null.
func (b Bool) Ptr() *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

// IsZero returns true for invalid Bools, for future omitempty support (Go 1.4?)
func (b Bool) IsZero() bool {
	return !b.Valid
}

// Scan implements the Scanner interface.
func (b *Bool) Scan(value interface{}) error {
	if value == nil {
		b.Bool, b.Valid = false, false
		return nil
	}
	b.Valid = true
	return convert.ConvertAssign(&b.Bool, value)
}

// Value implements the driver Valuer interface.
func (b Bool) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Bool, nil
}
