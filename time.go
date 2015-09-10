package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
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
		var ti time.Time
		return NewTime(ti, false)
	}
	return NewTime(*t, true)
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports time, string and null input.
func (t *Time) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	json.Unmarshal(data, &v)
	switch x := v.(type) {
	case time.Time:
		t.Time = x
	case string:
		t.Time, err = time.Parse(time.RFC3339, x)
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Time", reflect.TypeOf(v).Name())
	}
	t.Valid = err == nil
	return err
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
