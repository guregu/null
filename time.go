package null

import (
	"encoding/json"
	pq "github.com/lib/pq"
	"reflect"
	"time"
	"fmt"
)

// Time is an even nuller nullable Time.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Time struct {
	pq.NullTime
}

// NewTime creates a new Time
func NewTime(t time.Time, valid bool) Time {
	return Time{
		NullTime: pq.NullTime{
			Time:  t,
			Valid: valid,
		},
	}
}

// TimeFrom creates a new Time that will always be valid.
func TimeFrom(t time.Time) Time {
	return NewTime(t, true)
}

// TimeFromPtr creates a new Time that be null if i is nil.
func TimeFromPtr(t *time.Time) Time {
	if t == nil {
		return NewTime(time.Time{}, false)
	}
	return NewTime(*t, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports null.Time JSON or nil values
// It also supports unmarshalling a pq.NullTime
func (t *Time) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	json.Unmarshal(data, &v)
	switch v.(type) {
	case map[string]interface{}:
		err = json.Unmarshal(data, &t.NullTime)
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Time", reflect.TypeOf(v).Name())
	}
	t.Valid = err == nil
	return err

}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank or not an time.Time RFC3339.
// It will return an error if the input is not an time.Time RFC3339, blank, or "null".
func (t *Time) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	var err error
	t.Time, err = time.Parse(time.RFC3339, str)
	t.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode time.Time RFC3339 or null if this Time is null
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return []byte(t.Time.Format(time.RFC3339)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Time is null.
func (t Time) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte{}, nil
	}
	return []byte(t.Time.Format(time.RFC3339)), nil
}

// SetValid changes this Time's value and also sets it to be non-null.
func (t *Time) SetValid(n time.Time) {
	t.Time = n
	t.Valid = true
}

// Ptr returns a pointer to this Time's value, or a nil pointer if this Time is null.
func (t Time) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// IsZero returns true for invalid Times, for future omitempty support (Go 1.4?)
func (t Time) IsZero() bool {
	return !t.Valid
}
