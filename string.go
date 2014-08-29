// Package null provides an opinionated yet reasonable way of handling null values.
package null

import (
	"database/sql"
	"encoding/json"
)

// String is a nullable string.
type String struct {
	sql.NullString
}

// StringFrom creates a new String that will be null if s is blank.
func StringFrom(s string) String {
	return NewString(s, s != "")
}

// NewString creates a new String
func NewString(s string, valid bool) String {
	return String{
		NullString: sql.NullString{
			String: s,
			Valid:  s != "",
		},
	}
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input produces a null String.
// It also supports unmarshalling a sql.NullString.
func (s *String) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	json.Unmarshal(data, &v)
	switch v.(type) {
	case string:
		err = json.Unmarshal(data, &s.String)
	case map[string]interface{}:
		err = json.Unmarshal(data, &s.NullString)
	case nil:
		s.Valid = false
		return nil
	}
	s.Valid = (err == nil) && (s.String != "")
	return err
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string when this String is null.
func (s String) MarshalText() ([]byte, error) {
	if !s.Valid {
		return []byte{}, nil
	}
	return []byte(s.String), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null String if the input is a blank string.
func (s *String) UnmarshalText(text []byte) error {
	s.String = string(text)
	s.Valid = s.String != ""
	return nil
}

// Pointer returns a pointer to this String's value, or a nil pointer if this String is null.
func (s String) Pointer() *string {
	if s.String == "" {
		return nil
	}
	return &s.String
}

// IsZero returns true for null strings, for future omitempty support. (Go 1.4?)
func (s String) IsZero() bool {
	return !s.Valid || s.String == ""
}
