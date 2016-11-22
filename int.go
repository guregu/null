package null

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"

	"gopkg.in/mgo.v2/bson"
)

// Int is an nullable int64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Int struct {
	sql.NullInt64
}

// NewInt creates a new Int
func NewInt(i int64, valid bool) Int {
	return Int{
		NullInt64: sql.NullInt64{
			Int64: i,
			Valid: valid,
		},
	}
}

// IntFrom creates a new Int that will always be valid.
func IntFrom(i int64) Int {
	return NewInt(i, true)
}

// IntFromPtr creates a new Int that be null if i is nil.
func IntFromPtr(i *int64) Int {
	if i == nil {
		return NewInt(0, false)
	}
	return NewInt(*i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Int.
// It also supports unmarshalling a sql.NullInt64.
func (i *Int) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int64, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Int64)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullInt64)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

// UnmarshalXML implments the xml.Unmarshaler interface
func (i *Int) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var x *int64
	if err := d.DecodeElement(&x, &start); err != nil {
		return err
	}
	if x != nil {
		i.Valid = true
		i.Int64 = *x
	} else {
		i.Valid = false
	}
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	i.Int64, err = strconv.ParseInt(string(text), 10, 64)
	i.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Int is null.
func (i Int) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

// MarshalXML implements the xml.Marshaler interface
func (i Int) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if i.Valid {
		return e.EncodeElement(i.Int64, start)
	}
	return e.EncodeElement(nil, start)
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Int is null.
func (i Int) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

// GetBSON implements bson.Getter.
func (i Int) GetBSON() (interface{}, error) {
	if i.Valid {
		return i.Int64, nil
	}
	// TODO: do we need a nil pointer to a string?
	return nil, nil
}

// SetBSON implements bson.Setter.
func (i *Int) SetBSON(raw bson.Raw) error {
	var ii int64
	err := raw.Unmarshal(&ii)

	if err == nil {
		*i = Int{sql.NullInt64{Int64: ii, Valid: true}}
	} else {
		*i = Int{sql.NullInt64{Valid: false}}
	}
	return nil
}

// GetValue implements the compare.Valuable interface
func (i Int) GetValue() reflect.Value {
	if i.Valid {
		return reflect.ValueOf(i.Int64)
	}
	// or just nil?
	return reflect.ValueOf(nil)
}

// LoremDecode implements lorem.Decoder
func (i *Int) LoremDecode(tag, example string) error {
	i.SetValid(rand.Int63())
	return nil
}

// SetValid changes this Int's value and also sets it to be non-null.
func (i *Int) SetValid(n int64) {
	i.Int64 = n
	i.Valid = true
}

// Ptr returns a pointer to this Int's value, or a nil pointer if this Int is null.
func (i Int) Ptr() *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

// IsZero returns true for invalid Ints, for future omitempty support (Go 1.4?)
// A non-null Int with a 0 value will not be considered zero.
func (i Int) IsZero() bool {
	return !i.Valid
}
