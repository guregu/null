package zero

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

// Bool is a nullable bool. False input is considered null.
// JSON marshals to false if null.
// Considered null to SQL unmarshaled from a false value.
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

// BoolFrom creates a new Bool that will be null if false.
func BoolFrom(b bool) Bool {
	return NewBool(b, b)
}

// BoolFromPtr creates a new Bool that be null if b is nil.
func BoolFromPtr(b *bool) Bool {
	if b == nil {
		return NewBool(false, false)
	}
	return NewBool(*b, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// "false" will be considered a null Bool.
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
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type zero.Bool", reflect.TypeOf(v).Name())
	}
	b.Valid = (err == nil) && b.Bool
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Bool if the input is a false or not a bool.
// It will return an error if the input is not a float, blank, or "null".
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
	b.Valid = b.Bool
	return nil
}

// UnmarshalXML implments the xml.Unmarshaler interface
func (b *Bool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var x *bool
	if err := d.DecodeElement(&x, &start); err != nil {
		return err
	}
	if x != nil {
		b.Bool = *x
		b.Valid = b.Bool
	} else {
		b.Valid = false
	}
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode false if this Bool is null.
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid || !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a false if this Bool is null.
func (b Bool) MarshalText() ([]byte, error) {
	if !b.Valid || !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// MarshalXML implements the xml.Marshaler interface
func (b Bool) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	n := b.Bool
	if !b.Valid {
		n = false
	}
	return e.EncodeElement(n, start)
}

// GetBSON implements bson.Getter.
func (b Bool) GetBSON() (interface{}, error) {
	if b.Valid {
		return b.Bool, nil
	}
	return false, nil
}

// SetBSON implements bson.Setter.
func (b *Bool) SetBSON(raw bson.Raw) error {
	var bb bool
	err := raw.Unmarshal(&bb)

	if err == nil {
		*b = Bool{sql.NullBool{Bool: bb, Valid: bb}}
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
	return reflect.ValueOf(false)
}

// LoremDecode implements lorem.Decoder
func (b *Bool) LoremDecode(tag, example string) error {
	b.SetValid(rand.Int()%2 == 0)
	return nil
}

// SetValid changes this Bool's value and also sets it to be non-null.
func (b *Bool) SetValid(v bool) {
	b.Bool = v
	b.Valid = v
}

// Ptr returns a poBooler to this Bool's value, or a nil poBooler if this Bool is null.
func (b Bool) Ptr() *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

// IsZero returns true for null or zero Bools, for future omitempty support (Go 1.4?)
func (b Bool) IsZero() bool {
	return !b.Valid || !b.Bool
}
