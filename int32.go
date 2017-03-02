package null

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

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Int32.
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
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int32", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

// UnmarshalXML implments the xml.Unmarshaler interface
func (i *Int32) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var x *int32
	if err := d.DecodeElement(&x, &start); err != nil {
		return err
	}
	if x != nil {
		i.Valid = true
		i.Int32 = *x
	} else {
		i.Valid = false
	}
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int32 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int32) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	var res int64
	res, err = strconv.ParseInt(string(text), 10, 32)
	i.Int32 = int32(res)
	i.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Int32 is null.
func (i Int32) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(int64(i.Int32), 10)), nil
}

// MarshalXML implements the xml.Marshaler interface
func (i Int32) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if i.Valid {
		return e.EncodeElement(i.Int32, start)
	}
	return e.EncodeElement(nil, start)
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Int32 is null.
func (i Int32) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int32), 10)), nil
}

// GetBSON implements bson.Getter.
func (i Int32) GetBSON() (interface{}, error) {
	if i.Valid {
		return i.Int32, nil
	}
	// TODO: do we need a nil pointer to a string?
	return nil, nil
}

// SetBSON implements bson.Setter.
func (i *Int32) SetBSON(raw bson.Raw) error {
	var ii int32
	err := raw.Unmarshal(&ii)

	if err == nil {
		*i = Int32{sql.NullInt32{Int32: ii, Valid: true}}
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
	return reflect.ValueOf(nil)
}

// LoremDecode implements lorem.Decoder
func (i *Int32) LoremDecode(tag, example string) error {
	i.SetValid(rand.Int31())
	return nil
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

// IsZero returns true for invalid Ints, for future omitempty support (Go 1.4?)
// A non-null Int32 with a 0 value will not be considered zero.
func (i Int32) IsZero() bool {
	return !i.Valid
}
