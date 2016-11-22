package null

import (
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/axiomzen/null/format"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/pg.v4/types"
)

// Time is a nullable time.Time. It supports SQL and JSON serialization.
// It will marshal to null if null.
type Time struct {
	Time  time.Time
	Valid bool
}

///////////////////////////////
// PG interfaces
///////////////////////////////

// IsZero implements the orm.isZeroer interface
func (t Time) IsZero() bool {
	return !t.Valid
}

// AppendValue implements the types.ValueAppender interface
func (t Time) AppendValue(b []byte, quote int) ([]byte, error) {
	if !t.Valid {
		return types.AppendNull(b, quote), nil
	}
	// may write a zero time
	return types.AppendTime(b, t.Time, quote), nil
}

///////////////////////////////
// sql interfaces
///////////////////////////////

// Scan implements the sql.Scanner interface.
func (t *Time) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case time.Time:
		t.Time = x
	case nil:
		t.Valid = false
		return nil
	case []byte:
		b, _ := value.([]byte)
		t.Time, err = types.ParseTime(b)
	default:
		err = fmt.Errorf("null: cannot scan type %T into null.Time: %v", value, value)
	}
	t.Valid = err == nil
	return err
}

// Value implements the driver Valuer interface.
func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

///////////////////////////////
// methods
///////////////////////////////

// NewTime creates a new Time.
func NewTime(t time.Time, valid bool) Time {
	return Time{
		Time:  t,
		Valid: valid,
	}
}

// TimeFrom creates a new Time that will always be valid.
func TimeFrom(t time.Time) Time {
	return NewTime(t, true)
}

// TimeFromPtr creates a new Time that will be null if t is nil.
func TimeFromPtr(t *time.Time) Time {
	if t == nil {
		return NewTime(time.Time{}, false)
	}
	return NewTime(*t, true)
}

// SetValid changes this Time's value and sets it to be non-null.
func (t *Time) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

// Ptr returns a pointer to this Time's value, or a nil pointer if this Time is null.
func (t Time) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

///////////////////////////////
// json
///////////////////////////////

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		//return json.Marshal(nil)
		return []byte("null"), nil
	}
	//return t.Time.MarshalJSON()
	f := format.GetTimeFormat()
	b := make([]byte, 0, len(f)+2)
	b = append(b, '"')
	b = t.Time.AppendFormat(b, f)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string, object (e.g. pq.NullTime and friends)
// and null input.
func (t *Time) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		//err = t.Time.UnmarshalJSON(data)
		//RFC3339
		t.Time, err = time.Parse(`"`+format.GetTimeFormat()+`"`, string(data))
		t.Valid = err == nil
		return err
	case map[string]interface{}:
		ti, tiOK := x["Time"].(string)
		valid, validOK := x["Valid"].(bool)
		if !tiOK || !validOK {
			return fmt.Errorf(`json: unmarshalling object into Go value of type null.Time requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively`, x["Time"], x["Valid"])
		}

		err = t.UnmarshalText([]byte(ti))
		t.Valid = valid
		return err
		//err = t.Time.UnmarshalText(ti)
		//t.Valid = valid
		//return err
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Time", reflect.TypeOf(v).Name())
	}
	t.Valid = err == nil
	return err
}

///////////////////////////////
// xml
///////////////////////////////

// MarshalXML implements the xml.Marshaler interface
func (t Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if t.Valid {
		// to string?
		return e.EncodeElement(t.Time.Format(format.GetTimeFormat()), start)
	}
	return e.EncodeElement(nil, start)
}

// UnmarshalXML implments the xml.Unmarshaler interface
func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var s string
	err := d.DecodeElement(&s, &start)
	if err != nil {
		return err
	}
	t.Time, err = time.Parse(format.GetTimeFormat(), s)
	if err != nil {
		t.Valid = false
	} else {
		t.Valid = true
	}
	return nil
}

///////////////////////////////
// bson
///////////////////////////////

// GetBSON implements bson.Getter.
func (t Time) GetBSON() (interface{}, error) {
	if t.Valid {
		return t.Time, nil
	}
	var tt *time.Time
	return tt, nil
}

// SetBSON implements bson.Setter.
func (t *Time) SetBSON(raw bson.Raw) error {
	var ti time.Time
	err := raw.Unmarshal(&ti)

	if err == nil {
		*t = Time{Time: ti, Valid: true}
	} else {
		*t = Time{Valid: false}
	}
	return nil
}

///////////////////////////////
// text
///////////////////////////////

// MarshalText implements the encoding.TextMarshaler interface.
func (t Time) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	//return t.Time.MarshalText()
	f := format.GetTimeFormat()
	b := make([]byte, 0, len(f))
	return t.Time.AppendFormat(b, f), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *Time) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	// if err := t.Time.UnmarshalText(text); err != nil {
	// 	return err
	// }
	// Fractional seconds are handled implicitly by Parse.
	var err error
	t.Time, err = time.Parse(format.GetTimeFormat(), str)
	if err != nil {
		return err
	}

	t.Valid = true
	return nil
}

///////////////////////////////
// compare
///////////////////////////////

// GetValue implements the compare.Valuable interface
func (t Time) GetValue() reflect.Value {
	if t.Valid {
		return reflect.ValueOf(t.Time)
	}
	// or just nil?
	return reflect.ValueOf(nil)
}

///////////////////////////////
// lorem
///////////////////////////////

// LoremDecode implements lorem.Decoder
func (t *Time) LoremDecode(tag, example string) error {
	// fill ourselves with a random time
	// so now we can set valid right or not
	t.SetValid(time.Unix(rand.Int63(), 0))
	return nil
}
