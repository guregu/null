package zero

import (
	"database/sql"
	"strconv"

	"github.com/guregu/null/v5/internal"
)

// Byte is a nullable byte.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
type Byte struct {
	sql.NullByte
}

// NewByte creates a new Byte
func NewByte(i byte, valid bool) Byte {
	return Byte{
		NullByte: sql.NullByte{
			Byte:  i,
			Valid: valid,
		},
	}
}

// ByteFrom creates a new Byte that will be null if zero.
func ByteFrom(i byte) Byte {
	return NewByte(i, i != 0)
}

// ByteFromPtr creates a new Byte that be null if i is nil.
func ByteFromPtr(i *byte) Byte {
	if i == nil {
		return NewByte(0, false)
	}
	n := NewByte(*i, true)
	return n
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (b Byte) ValueOrZero() byte {
	if !b.Valid {
		return 0
	}
	return b.Byte
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will be considered a null Byte.
func (b *Byte) UnmarshalJSON(data []byte) error {
	err := internal.UnmarshalIntJSON(data, &b.Byte, &b.Valid, 8, strconv.ParseUint)
	if err != nil {
		return err
	}
	b.Valid = b.Byte != 0
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Byte if the input is a blank, or zero.
// It will return an error if the input is not an integer, blank, or "null".
func (b *Byte) UnmarshalText(text []byte) error {
	err := internal.UnmarshalIntText(text, &b.Byte, &b.Valid, 8, strconv.ParseUint)
	if err != nil {
		return err
	}
	b.Valid = b.Byte != 0
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode 0 if this Byte is null.
func (b Byte) MarshalJSON() ([]byte, error) {
	n := b.Byte
	if !b.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(int64(n), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Byte is null.
func (b Byte) MarshalText() ([]byte, error) {
	n := b.Byte
	if !b.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(int64(n), 10)), nil
}

// SetValid changes this Byte's value and also sets it to be non-null.
func (b *Byte) SetValid(n byte) {
	b.Byte = n
	b.Valid = true
}

// Ptr returns a pointer to this Byte's value, or a nil pointer if this Byte is null.
func (b Byte) Ptr() *byte {
	if !b.Valid {
		return nil
	}
	return &b.Byte
}

// IsZero returns true for null or zero Bytes, for future omitempty support (Go 1.4?)
func (b Byte) IsZero() bool {
	return !b.Valid || b.Byte == 0
}

// Equal returns true if both ints have the same value or are both either null or zero.
func (b Byte) Equal(other Byte) bool {
	return b.ValueOrZero() == other.ValueOrZero()
}

func (b Byte) value() (int64, bool) {
	return int64(b.Byte), b.Valid
}
