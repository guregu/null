package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"gopkg.in/nullbio/null.v6/convert"
)

// String is a nullable string. It supports SQL and JSON serialization.
type String struct {
	String string
	Valid  bool
}

// StringFrom creates a new String that will never be blank.
func StringFrom(s string) String {
	return NewString(s, true)
}

// StringFromPtr creates a new String that be null if s is nil.
func StringFromPtr(s *string) String {
	if s == nil {
		return NewString("", false)
	}
	return NewString(*s, true)
}

// NewString creates a new String
func NewString(s string, valid bool) String {
	return String{
		String: s,
		Valid:  valid,
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *String) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		s.String = ""
		s.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &s.String); err != nil {
		return err
	}

	s.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return NullBytes, nil
	}
	return json.Marshal(s.String)
}

// MarshalText implements encoding.TextMarshaler.
func (s String) MarshalText() ([]byte, error) {
	if !s.Valid {
		return []byte{}, nil
	}
	return []byte(s.String), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *String) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		s.Valid = false
		return nil
	}

	s.String = string(text)
	s.Valid = true
	return nil
}

// SetValid changes this String's value and also sets it to be non-null.
func (s *String) SetValid(v string) {
	s.String = v
	s.Valid = true
}

// Ptr returns a pointer to this String's value, or a nil pointer if this String is null.
func (s String) Ptr() *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// IsZero returns true for null strings, for potential future omitempty support.
func (s String) IsZero() bool {
	return !s.Valid
}

// Scan implements the Scanner interface.
func (s *String) Scan(value interface{}) error {
	if value == nil {
		s.String, s.Valid = "", false
		return nil
	}
	s.Valid = true
	return convert.ConvertAssign(&s.String, value)
}

// Value implements the driver Valuer interface.
func (s String) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.String, nil
}
