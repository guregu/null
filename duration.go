package null

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	errBadDurationSQLType = errors.New("incompatible SQL type for Scan() into null.Duration")
)

// Duration is a nullable time.Duration. It supports SQL and JSON serialization.
// It will marshal to null if null.
type Duration struct {
	Duration time.Duration
	Valid    bool
}

// Scan implements the Scanner interface.
// Supports the ´interval` type for PostgreSQL
func (d *Duration) Scan(value interface{}) error {
	var source string
	switch value.(type) {
	case string:
		source = value.(string)
	case []byte:
		source = string(value.([]byte))
	default:
		return errBadDurationSQLType
	}

	split := strings.Split(source, ":")
	if len(split) != 3 {
		return errBadDurationSQLType
	}

	hours, err := strconv.Atoi(split[0])
	if err != nil {
		return errBadDurationSQLType
	}

	minutes, err := strconv.Atoi(split[1])
	if err != nil {
		return errBadDurationSQLType
	}

	seconds, err := strconv.Atoi(split[2])
	if err != nil {
		return errBadDurationSQLType
	}

	duration := (time.Hour * time.Duration(hours)) +
		(time.Minute * time.Duration(minutes)) +
		(time.Second * time.Duration(seconds))

	*d = Duration{
		Duration: duration,
		Valid:    true,
	}

	return nil
}

// Value implements the driver Valuer interface.
func (d Duration) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Duration, nil
}

// NewDuration creates a new Duration.
func NewDuration(d time.Duration, valid bool) Duration {
	return Duration{
		Duration: d,
		Valid:    valid,
	}
}

// DurationFrom creates a new Duration that will always be valid.
func DurationFrom(d time.Duration) Duration {
	return NewDuration(d, true)
}

// DurationFromPtr creates a new Duration that will be null if d is nil.
func DurationFromPtr(d *time.Duration) Duration {
	if d == nil {
		return NewDuration(time.Duration(0), false)
	}
	return NewDuration(*d, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (d Duration) ValueOrZero() time.Duration {
	if !d.Valid {
		return time.Duration(0)
	}

	return d.Duration
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this duration is null
// or use the string representation if it's valid.
func (d Duration) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(d.Duration.String())
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input.
func (d *Duration) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		duration, err := time.ParseDuration(x)
		if err != nil {
			return err
		}

		d.Duration = duration

	case nil:
		d.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Duration", reflect.TypeOf(v).Name())
	}

	d.Valid = (err == nil)
	return err
}

// SetValid changes this Time's value and sets it to be non-null.
func (d *Duration) SetValid(v time.Duration) {
	d.Duration = v
	d.Valid = true
}

// Ptr returns a pointer to this Duration's value, or a nil pointer if this Duration is null.
func (d *Duration) Ptr() *time.Duration {
	if !d.Valid {
		return nil
	}

	return &d.Duration
}
