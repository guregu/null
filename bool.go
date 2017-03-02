package null

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"math/rand"
	"reflect"

	"gopkg.in/mgo.v2/bson"
)

// Bool is a nullable bool.
// It does not consider false values to be null.
// It will decode to null, not false, if null.
type Bool struct {
	sql.NullBool
}

// NewBool creates a new Bool
func NewBool(b bool, valid bool) Bool {
	return Bool{
		NullBool: sql.NullBool{
			Bool:  b,
			Valid: valid,
		},
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
// It supports number and null input.
// 0 will not be considered a null Bool.
// It also supports unmarshalling a sql.NullBool.
func (b *Bool) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case bool:
		b.Bool = x
	case map[string]interface{}:
		err = json.Unmarshal(data, &b.NullBool)
	case nil:
		b.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Bool", reflect.TypeOf(v).Name())
	}
	b.Valid = err == nil
	return err
}

// UnmarshalXML implments the xml.Unmarshaler interface
func (b *Bool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var x *bool
	if err := d.DecodeElement(&x, &start); err != nil {
		return err
	}
	if x != nil {
		b.Valid = true
		b.Bool = *x
	} else {
		b.Valid = false
	}
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Bool if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (b *Bool) UnmarshalText(text []byte) error {
	str := string(text)
	switch str {
	case "", "null":
		b.Valid = false
		return nil
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
// It will encode null if this Bool is null.
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Bool is null.
func (b Bool) MarshalText() ([]byte, error) {
	if !b.Valid {
		return []byte{}, nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// MarshalXML implements the xml.Marshaler interface
func (b Bool) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if b.Valid {
		return e.EncodeElement(b.Bool, start)
	}
	return e.EncodeElement(nil, start)
}

// GetBSON implements bson.Getter.
func (b Bool) GetBSON() (interface{}, error) {
	if b.Valid {
		return b.Bool, nil
	}
	// TODO: do we need a nil pointer to a string?
	return nil, nil
}

// SetBSON implements bson.Setter.
func (b *Bool) SetBSON(raw bson.Raw) error {
	var bb bool
	err := raw.Unmarshal(&bb)

	if err == nil {
		*b = Bool{sql.NullBool{Bool: bb, Valid: true}}
	} else {
		*b = Bool{sql.NullBool{Valid: false}}
	}
	return nil
}

// GetValue implements the compare.Valuable interface
func (b Bool) GetValue() reflect.Value {
	if b.Valid {
		return reflect.ValueOf(b.Bool)
	}
	// or just nil?
	return reflect.ValueOf(nil)
}

// LoremDecode implements lorem.Decoder
func (b *Bool) LoremDecode(tag, example string) error {
	b.SetValid(rand.Int()%2 == 0)
	return nil
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
// A non-null Bool with a 0 value will not be considered zero.
func (b Bool) IsZero() bool {
	return !b.Valid
}
