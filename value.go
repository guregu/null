//go:build go1.22

package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// Value represents a value that may be null.
type Value[T any] struct {
	sql.Null[T]
}

// NewValue creates a new Value.
func NewValue[T any](t T, valid bool) Value[T] {
	return Value[T]{
		Null: sql.Null[T]{
			V:     t,
			Valid: valid,
		},
	}
}

// ValueFrom creates a new Value that will always be valid.
func ValueFrom[T any](t T) Value[T] {
	return NewValue(t, true)
}

// ValueFromPtr creates a new Value that will be null if t is nil.
func ValueFromPtr[T any](t *T) Value[T] {
	if t == nil {
		var zero T
		return NewValue(zero, false)
	}
	return NewValue(*t, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (t Value[T]) ValueOrZero() T {
	if !t.Valid {
		var zero T
		return zero
	}
	return t.V
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this value is null.
func (t Value[T]) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.V)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input.
func (t *Value[T]) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == 'n' {
		t.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &t.V); err != nil {
		return fmt.Errorf("null: couldn't unmarshal JSON: %w", err)
	}

	t.Valid = true
	return nil
}

/*
// MarshalText implements encoding.TextMarshaler.
// It returns an empty string if invalid, otherwise T's MarshalText.
func (t Value[T]) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte{}, nil
	}
	if tm, ok := any(t.V).(encoding.TextMarshaler); ok {
		return tm.MarshalText()
	}

	rv := reflect.ValueOf(t.V)
	if !rv.IsValid() {
		return []byte{}, nil
	}

try:
	switch rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() {
			return []byte{}, nil
		}
		rv = rv.Elem()
		goto try
	case reflect.String:
		return []byte(rv.String()), nil
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return []byte(strconv.FormatInt(rv.Int(), 10)), nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return []byte(strconv.FormatUint(rv.Uint(), 10)), nil
	case reflect.Float32, reflect.Float64:
		return []byte(strconv.FormatFloat(rv.Float(), 'f', -1, rv.Type().Bits())), nil

		// 	case reflect.Slice:
		// 		if rv.IsNil() {
		// 			return []byte{}, nil
		// 		}
		// 		if rv.Type().Elem().Kind() == reflect.Uint8 {
		// 			return rv.Bytes(), nil
		// 		}
		//
	}

	return t.Value.MarshalText()
}
*/

// SetValid changes this Value's value and sets it to be non-null.
func (t *Value[T]) SetValid(v T) {
	t.V = v
	t.Valid = true
}

// Ptr returns a pointer to this Value's value, or a nil pointer if this Value is null.
func (t Value[T]) Ptr() *T {
	if !t.Valid {
		return nil
	}
	return &t.V
}

// IsZero returns true for invalid Values, hopefully for future omitempty support.
// A non-null Value with a zero value will not be considered zero.
func (t Value[T]) IsZero() bool {
	return !t.Valid
}

/*
// Equal returns true if both Value objects encode the same value or are both null.
func (t Value[T]) Equal(other Value[T]) bool {
	return t.Valid == other.Valid && (t.V == other.V)
}
*/
