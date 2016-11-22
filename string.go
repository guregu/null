// Package null contains SQL types that consider zero input and null input as separate values,
// with convenient support for JSON and text marshaling.
// Types in this package will always encode to their null value if null.
// Use the zero subpackage if you want zero values and null to be treated the same.
package null

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"

	"gopkg.in/mgo.v2/bson"
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
// It supports string and null input. Blank string input does not produce a null String.
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
	s.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this String is null.
func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
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

// MarshalXML implements the xml.Marshaler interface
func (s String) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if s.Valid {
		return e.EncodeElement(s.String, start)
	}
	return e.EncodeElement(nil, start)
}

// UnmarshalXML implments the xml.Unmarshaler interface
func (s *String) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var x *string
	if err := d.DecodeElement(&x, &start); err != nil {
		return err
	}
	if x != nil {
		s.Valid = true
		s.String = *x
	} else {
		s.Valid = false
	}
	return nil
}

// GetBSON implements bson.Getter.
func (s String) GetBSON() (interface{}, error) {
	if s.Valid {
		return s.String, nil
	}
	// TODO: do we need a nil pointer to a string?
	return nil, nil
}

// SetBSON implements bson.Setter.
func (s *String) SetBSON(raw bson.Raw) error {
	var str string
	err := raw.Unmarshal(&str)

	if err == nil {
		*s = String{sql.NullString{String: str, Valid: true}}
	} else {
		*s = String{sql.NullString{Valid: false}}
	}
	return nil
}

// GetValue implements the compare.Valuable interface
func (s String) GetValue() reflect.Value {
	if s.Valid {
		return reflect.ValueOf(s.String)
	}
	// or just nil?
	return reflect.ValueOf(nil)
}

// LoremDecode implements lorem.Decoder
func (s *String) LoremDecode(tag, example string) error {
	s.SetValid(example)
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
