package zero

import (
	"database/sql"
	"strconv"

	"github.com/guregu/null/v5/internal"
)

// Int32 is a nullable int32.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
type Int32 struct {
	sql.NullInt32
}

// NewInt32 creates a new Int32
func NewInt32(i int32, valid bool) Int32 {
	return Int32{
		NullInt32: sql.NullInt32{
			Int32: i,
			Valid: valid,
		},
	}
}

// Int32From creates a new Int32 that will be null if zero.
func Int32From(i int32) Int32 {
	return NewInt32(i, i != 0)
}

// Int32FromPtr creates a new Int32 that be null if i is nil.
func Int32FromPtr(i *int32) Int32 {
	if i == nil {
		return NewInt32(0, false)
	}
	n := NewInt32(*i, true)
	return n
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Int32) ValueOrZero() int32 {
	if !i.Valid {
		return 0
	}
	return i.Int32
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will be considered a null Int32.
func (i *Int32) UnmarshalJSON(data []byte) error {
	err := internal.UnmarshalIntJSON(data, &i.Int32, &i.Valid, 32, strconv.ParseInt)
	if err != nil {
		return err
	}
	i.Valid = i.Int32 != 0
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int32 if the input is a blank, or zero.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int32) UnmarshalText(text []byte) error {
	err := internal.UnmarshalIntText(text, &i.Int32, &i.Valid, 32, strconv.ParseInt)
	if err != nil {
		return err
	}
	i.Valid = i.Int32 != 0
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode 0 if this Int32 is null.
func (i Int32) MarshalJSON() ([]byte, error) {
	n := i.Int32
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(int64(n), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Int32 is null.
func (i Int32) MarshalText() ([]byte, error) {
	n := i.Int32
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(int64(n), 10)), nil
}

// SetValid changes this Int32's value and also sets it to be non-null.
func (i *Int32) SetValid(n int32) {
	i.Int32 = n
	i.Valid = true
}

// Ptr returns a pointer to this Int32's value, or a nil pointer if this Int32 is null.
func (i Int32) Ptr() *int32 {
	if !i.Valid {
		return nil
	}
	return &i.Int32
}

// IsZero returns true for null or zero Int32s, for future omitempty support (Go 1.4?)
func (i Int32) IsZero() bool {
	return !i.Valid || i.Int32 == 0
}

// Equal returns true if both ints have the same value or are both either null or zero.
func (i Int32) Equal(other Int32) bool {
	return i.ValueOrZero() == other.ValueOrZero()
}

func (i Int32) value() (int64, bool) {
	return int64(i.Int32), i.Valid
}
