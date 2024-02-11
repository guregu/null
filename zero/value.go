//go:build go1.22

package zero

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Value[T comparable] struct {
	sql.Null[T]
}

// NewValue creates a new Value.
func NewValue[T comparable](t T, valid bool) Value[T] {
	return Value[T]{
		Null: sql.Null[T]{
			V:     t,
			Valid: valid,
		},
	}
}

// ValueFrom creates a new Value that will always be valid.
func ValueFrom[T comparable](t T) Value[T] {
	var zero T
	return NewValue(t, t != zero)
}

// ValueFromPtr creates a new Value that will be null if t is nil.
func ValueFromPtr[T comparable](t *T) Value[T] {
	var zero T
	if t == nil {
		return NewValue(zero, false)
	}
	return NewValue(*t, *t != zero)
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
// It will encode null if this value is null or zero.
func (t Value[T]) MarshalJSON() ([]byte, error) {
	var zero T
	if !t.Valid || t.V == zero {
		return []byte("null"), nil
	}
	return json.Marshal(t.V)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input.
func (t *Value[T]) UnmarshalJSON(data []byte) error {
	var zero T
	if len(data) > 0 && data[0] == 'n' {
		t.Valid = false
		t.V = zero
		return nil
	}

	if err := json.Unmarshal(data, &t.V); err != nil {
		return fmt.Errorf("null: couldn't unmarshal JSON: %w", err)
	}

	t.Valid = t.V != zero
	return nil
}

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

// IsZero returns true for invalid or zero Values, hopefully for future omitempty support.
func (t Value[T]) IsZero() bool {
	var zero T
	return !t.Valid || t.V == zero
}

// Equal returns true if both Value objects encode the same value or are both null.
func (t Value[T]) Equal(other Value[T]) bool {
	return t.ValueOrZero() == other.ValueOrZero()
}
