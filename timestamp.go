package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Timestamp is a nullable time.Time. It supports SQL and JSON serialization.
// It will marshal to null if null.
type Timestamp struct {
	Time  time.Time
	Valid bool
}

// Scan implements the Scanner interface.
func (t *Timestamp) Scan(value interface{}) error {
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
func (t Timestamp) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// NewTimestamp creates a new Timestamp.
func NewTimestamp(t time.Time, valid bool) Timestamp {
	return Timestamp{
		Time:  t,
		Valid: valid,
	}
}

// TimestampFrom creates a new Timestamp that will always be valid.
func TimestampFrom(t time.Time) Timestamp {
	return NewTimestamp(t, true)
}

// TimestampFromPtr creates a new Timestamp that will be null if t is nil.
func TimestampFromPtr(t *time.Time) Timestamp {
	if t == nil {
		return NewTimestamp(time.Time{}, false)
	}
	return NewTimestamp(*t, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (t Timestamp) ValueOrZero() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(t.Time.Unix(), 10)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string, object (e.g. pq.NullTime and friends)
// and null input.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch x := v.(type) {
	case float64:
		t.Time = time.Unix(int64(x), 0)
	case map[string]interface{}:
		ti, tiOK := x["Time"].(float64)
		valid, validOK := x["Valid"].(bool)
		if !tiOK || !validOK {
			return fmt.Errorf(`json: unmarshalling object into Go value of type null.Timestamp requires key "Time" to be of type float64 and key "Valid" to be of type bool; found %T and %T, respectively`, x["Time"], x["Valid"])
		}
		t.Time = time.Unix(int64(ti), 0)
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

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Timestamp is null.
func (t Timestamp) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(t.Time.Unix(), 10)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null int64 Unix timestamp to time.Time if the input is a blank or not an time.Time.
func (t *Timestamp) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	v, err := strconv.ParseInt(str, 0, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(v, 0)
	t.Valid = true
	return nil
}

// SetValid changes this Timestamp's value and sets it to be non-null.
func (t *Timestamp) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

// Ptr returns a pointer to this Timestamp's value, or a nil pointer if this Timestamp is null.
func (t Timestamp) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
