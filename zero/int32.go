package zero

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"

	"github.com/axiomzen/null/sql"

	"gopkg.in/mgo.v2/bson"
)

// Int32 is a nullable int32.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
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

// Int32From creates a new Int32 that will be null if zero.
func Int32From(i int32) Int32 {
	return NewInt32(i, i != 0)
}

// Int32FromPtr creates a new Int32 that be null if i is nil.
func Int32FromPtr(i *int32) Int32 {
	if i == nil {
		return NewInt32(0, false)
	}
	n := NewInt32(*i, true)
	return n
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will be considered a null Int32.
// It also supports unmarshalling a sql.NullInt32.
func (i *Int32) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float32, float64:
		// Unmarshal again, directly to int32, to avoid intermediate float32
		err = json.Unmarshal(data, &i.Int32)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullInt32)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type zero.Int32", reflect.TypeOf(v).Name())
	}
	i.Valid = (err == nil) && (i.Int32 != 0)
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int32 if the input is a blank, zero, or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int32) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	var v int64
	v, err = strconv.ParseInt(string(text), 10, 32)
	i.Int32 = int32(v)
	i.Valid = (err == nil) && (i.Int32 != 0)
	return err
}

// UnmarshalXML implments the xml.Unmarshaler interface
func (i *Int32) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var x *int32
	if err := d.DecodeElement(&x, &start); err != nil {
		return err
	}
	if x != nil {
		i.Int32 = *x
		i.Valid = i.Int32 != 0
	} else {
		i.Valid = false
	}
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode 0 if this Int32 is null.
func (i Int32) MarshalJSON() ([]byte, error) {
	n := i.Int32
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(int64(n), 10)), nil
}

// MarshalXML implements the xml.Marshaler interface
func (i Int32) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	n := i.Int32
	if !i.Valid {
		n = 0
	}
	return e.EncodeElement(n, start)
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Int32 is null.
func (i Int32) MarshalText() ([]byte, error) {
	n := i.Int32
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(int64(n), 10)), nil
}

// GetBSON implements bson.Getter.
func (i Int32) GetBSON() (interface{}, error) {
	n := i.Int32
	if !i.Valid {
		n = 0
	}
	// TODO: do we need a nil pointer to a string?
	return n, nil
}

// SetBSON implements bson.Setter.
func (i *Int32) SetBSON(raw bson.Raw) error {
	var ii int32
	err := raw.Unmarshal(&ii)

	if err == nil {
		*i = Int32{sql.NullInt32{Int32: ii, Valid: ii != 0}}
	} else {
		*i = Int32{sql.NullInt32{Valid: false}}
	}
	return nil
}

// GetValue implements the compare.Valuable interface
func (i Int32) GetValue() reflect.Value {
	if i.Valid {
		return reflect.ValueOf(i.Int32)
	}
	// or just nil?
	return reflect.ValueOf(0)
}

// LoremDecode implements lorem.Decoder
func (i *Int32) LoremDecode(tag, example string) error {
	i.SetValid(rand.Int31())
	return nil
}

// SetValid changes this Int32's value and also sets it to be non-null.
func (i *Int32) SetValid(n int32) {
	i.Int32 = n
	i.Valid = n != 0
}

// Ptr returns a pointer to this Int32's value, or a nil pointer if this Int32 is null.
func (i Int32) Ptr() *int32 {
	if !i.Valid {
		return nil
	}
	return &i.Int32
}

// IsZero returns true for null or zero Ints, for future omitempty support (Go 1.4?)
func (i Int32) IsZero() bool {
	return !i.Valid || i.Int32 == 0
}
