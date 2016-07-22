package zero

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strconv"
)

// Int is a nullable int64.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
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

// IntFrom creates a new Int that will be null if zero.
func IntFrom(i int64) Int {
	return NewInt(i, i != 0)
}

// IntFromPtr creates a new Int that be null if i is nil.
func IntFromPtr(i *int64) Int {
	if i == nil {
		return NewInt(0, false)
	}
	n := NewInt(*i, true)
	return n
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will be considered a null Int.
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
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type zero.Int", reflect.TypeOf(v).Name())
	}
	i.Valid = (err == nil) && (i.Int64 != 0)
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank, zero, or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	i.Int64, err = strconv.ParseInt(string(text), 10, 64)
	i.Valid = (err == nil) && (i.Int64 != 0)
	return err
}

// UnmarshalXML implments the xml.Unmarshaler interface
func (i *Int) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var x *int64
	if err := d.DecodeElement(&x, &start); err != nil {
		return err
	}
	if x != nil {
		i.Int64 = *x
		i.Valid = i.Int64 != 0
	} else {
		i.Valid = false
	}
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode 0 if this Int is null.
func (i Int) MarshalJSON() ([]byte, error) {
	n := i.Int64
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(n, 10)), nil
}

// MarshalXML implements the xml.Marshaler interface
func (i Int) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	n := i.Int64
	if !i.Valid {
		n = 0
	}
	return e.EncodeElement(n, start)
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Int is null.
func (i Int) MarshalText() ([]byte, error) {
	n := i.Int64
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(n, 10)), nil
}

// GetBSON implements bson.Getter.
func (i Int) GetBSON() (interface{}, error) {
	n := i.Int64
	if !i.Valid {
		n = 0
	}
	// TODO: do we need a nil pointer to a string?
	return n, nil
}

// SetBSON implements bson.Setter.
func (i *Int) SetBSON(raw bson.Raw) error {
	var ii int64
	err := raw.Unmarshal(&ii)

	if err == nil {
		*i = Int{sql.NullInt64{Int64: ii, Valid: ii != 0}}
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
	return reflect.ValueOf(0)
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

// IsZero returns true for null or zero Ints, for future omitempty support (Go 1.4?)
func (i Int) IsZero() bool {
	return !i.Valid || i.Int64 == 0
}
