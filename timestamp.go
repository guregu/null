package null

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	pq "github.com/lib/pq"
)

// Timestamp is helper an even nuller nullable Time.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Timestamp struct {
	pq.NullTime
}

// NewTime creates a new Time
func NewTimestamp(t time.Time, valid bool) Timestamp {
	return Timestamp{
		NullTime: pq.NullTime{
			Time:  t,
			Valid: valid,
		},
	}
}

// TimestampFrom creates a new Timestamp that will always be valid.
func TimestampFrom(t time.Time) Timestamp {
	return NewTimestamp(t, true)
}

// TimestampFromPtr creates a new Time that be null if i is nil.
func TimestampFromPtr(t *time.Time) Timestamp {
	if t == nil {
		return NewTimestamp(time.Time{}, false)
	}
	return NewTimestamp(*t, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports null.Timestamp JSON or nil values
// It also supports unmarshalling a pq.NullTime
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	var i struct {
		Time int64
	}
	var j struct {
		Time  int64
		Valid bool
	}
	var s string
	json.Unmarshal(data, &v)
	switch v.(type) {
	case float64:
		err = json.Unmarshal(data, &i.Time)
		if err == nil {
			t.Time = time.Unix(i.Time, 0)
		}
	case map[string]interface{}:
		err = json.Unmarshal(data, &j)
		if (err == nil) && j.Valid {
			t.Time = time.Unix(j.Time, 0)
		}
	case string:
		err = json.Unmarshal(data, &s)
		t.Time, _ = time.Parse(time.RFC3339, s)
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Timestamp", reflect.TypeOf(v).Name())
	}
	t.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null int64 Unix timestamp to time.Time if the input is a blank or not an time.Time.
func (t *Timestamp) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	var (
		err error
		i   int64
	)
	i, err = strconv.ParseInt(str, 0, 64)
	if err == nil {
		t.Time = time.Unix(i, 0)
	}
	t.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(t.Time.Unix(), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Timestamp is null.
func (t Timestamp) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(t.Time.Unix(), 10)), nil
}

// SetValid changes this Timestamp's value and also sets it to be non-null.
func (t *Timestamp) SetValid(n time.Time) {
	t.Time = n
	t.Valid = true
}

// Ptr returns a pointer to this Timestamp's value, or a nil pointer if this Timestamp is null.
func (t Timestamp) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// IsZero returns true for invalid Timestamps, for future omitempty support (Go 1.4?)
func (t Timestamp) IsZero() bool {
	return !t.Valid
}
