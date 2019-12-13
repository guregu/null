package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
	"github.com/philpearl/plenc"
)

// Time is a nullable time.Time. It supports SQL and JSON serialization.
// It will marshal to null if null.
type Time struct {
	Time  time.Time
	Valid bool
}

// Scan implements the Scanner interface.
func (t *Time) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case time.Time:
		t.Time = x
	case nil:
		t.Valid = false
		return nil
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

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return nullLiteral, nil
	}
	return t.Time.MarshalJSON()
}

func (t Time) MarshalEasyJSON(w *jwriter.Writer) {
	if !t.Valid {
		w.RawString("null")
		return
	}
	w.Buffer.EnsureSpace(len(time.RFC3339Nano) + 2)
	w.Buffer.AppendByte('"')
	w.Buffer.Buf = t.Time.UTC().AppendFormat(w.Buffer.Buf, time.RFC3339Nano)
	w.Buffer.AppendByte('"')
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
		err = t.Time.UnmarshalJSON(data)
	case map[string]interface{}:
		ti, tiOK := x["Time"].(string)
		valid, validOK := x["Valid"].(bool)
		if !tiOK || !validOK {
			return fmt.Errorf(`json: unmarshalling object into Go value of type null.Time requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively`, x["Time"], x["Valid"])
		}
		err = t.Time.UnmarshalText([]byte(ti))
		t.Valid = valid
		return err
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Time", reflect.TypeOf(v).Name())
	}
	t.Valid = err == nil
	return err
}

// UnmarshalEasyJSON is an easy-JSON specific decoder, that should be more efficient than the standard one.
// We expect the value to be either `null` or `true`, but we also unmarshal if we receive
// `{"Valid":true,"Bool":false}`
func (i *Time) UnmarshalEasyJSON(w *jlexer.Lexer) {
	if w.IsNull() {
		w.Skip()
		i.Valid = false
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
			case "time", "Time":
				t, err := time.Parse(time.RFC3339, w.UnsafeString())
				if err != nil {
					w.AddError(err)
				}
				i.Time = t
			case "valid", "Valid":
				i.Valid = w.Bool()
			}
			w.WantComma()
		}
		return
	}
	t, err := time.Parse(time.RFC3339, w.UnsafeString())
	if err != nil {
		w.AddError(err)
	}
	i.Time = t
	i.Valid = (w.Error() == nil)
}

func (t Time) MarshalText() ([]byte, error) {
	if !t.Valid {
		return nullLiteral, nil
	}
	return t.Time.MarshalText()
}

func (t *Time) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	if err := t.Time.UnmarshalText(text); err != nil {
		return err
	}
	t.Valid = true
	return nil
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

// ΦλSize determines how many bytes are needed to encode this value
func (t Time) ΦλSize() (size int) {
	if !t.Valid {
		return 0
	}
	var pt plenc.Time
	pt.Set(t.Time)
	return pt.ΦλSize()
}

// ΦλAppend encodes example by appending to data. It returns the final slice
func (t Time) ΦλAppend(data []byte) []byte {
	if !t.Valid {
		return data
	}
	var pt plenc.Time
	pt.Set(t.Time)
	return pt.ΦλAppend(data)
}

// ΦλUnmarshal decodes a plenc encoded value
func (t *Time) ΦλUnmarshal(data []byte) (int, error) {
	// There's no tag within the encoding. If we're being asked to decode, then this value field must be present
	// within the encoded data,
	t.Valid = true
	var pt plenc.Time
	n, err := pt.ΦλUnmarshal(data)
	t.Time = pt.Standard()
	return n, err
}
