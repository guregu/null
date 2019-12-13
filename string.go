// Package null contains SQL types that consider zero input and null input as separate values,
// with convenient support for JSON and text marshaling.
// Types in this package will always encode to their null value if null.
// Use the zero subpackage if you want zero values and null to be treated the same.
package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
	"github.com/philpearl/plenc"
)

// String is a nullable string. It supports SQL and JSON serialization.
// It will marshal to null if null. Blank string input will be considered null.
type String struct {
	sql.NullString
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
		NullString: sql.NullString{
			String: s,
			Valid:  valid,
		},
	}
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input produces a null String.
// It also supports unmarshalling a sql.NullString.
func (s *String) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		s.String = x
	case map[string]interface{}:
		err = json.Unmarshal(data, &s.NullString)
	case nil:
		s.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.String", reflect.TypeOf(v).Name())
	}
	s.Valid = (err == nil) && (s.String != "")
	return err
}

// UnmarshalEasyJSON is an easy-JSON specific decoder, that should be more efficient than the standard one.
// We expect the value to be either `null` or `"a string"`, but we also unmarshal if we receive
// `{"Valid":true,"String":"a string"}`
func (s *String) UnmarshalEasyJSON(w *jlexer.Lexer) {
	if w.IsNull() {
		w.Skip()
		s.Valid = false
		return
	}
	if w.IsDelim('{') {
		w.Skip()
		for !w.IsDelim('}') {
			key := w.UnsafeString()
			w.WantColon()
			if w.IsNull() {
				w.Skip()
				w.WantComma()
				continue
			}
			switch key {
			case "string", "String":
				s.String = w.String()
			case "valid", "Valid":
				s.Valid = w.Bool()
			}
			w.WantComma()
		}
		return
	}
	s.String = w.String()
	s.Valid = (w.Error() == nil)
}

func (s String) MarshalEasyJSON(w *jwriter.Writer) {
	if !s.Valid {
		w.RawString("null")
		return
	}

	w.String(s.String)
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this String is null.
func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return nullLiteral, nil
	}
	return json.Marshal(s.String)
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null String if the input is a blank string.
func (s *String) UnmarshalText(text []byte) error {
	s.String = string(text)
	s.Valid = s.String != ""
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

// IsZero returns true for null or empty strings, for future omitempty support. (Go 1.4?)
// Will return false s if blank but non-null.
func (s String) IsZero() bool {
	return !s.Valid
}

// ΦλSize determines how many bytes are needed to encode this value
func (s String) ΦλSize() (size int) {
	if !s.Valid {
		return 0
	}
	// We're going to cheat and assume this will always only include a single value. So we won't do the tag
	// thing
	return plenc.SizeString(s.String)
}

// ΦλAppend encodes example by appending to data. It returns the final slice
func (s String) ΦλAppend(data []byte) []byte {
	if !s.Valid {
		return data
	}
	return plenc.AppendString(data, s.String)
}

// ΦλUnmarshal decodes a plenc encoded value
func (s *String) ΦλUnmarshal(data []byte) (int, error) {
	// There's no tag within the encoding. If we're being asked to decode, then this value field must be present
	// within the encoded data,
	s.Valid = true
	var n int
	s.String, n = plenc.ReadString(data)
	return n, nil
}
