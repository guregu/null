package null

import (
	"database/sql"
	"strconv"

	"github.com/guregu/null/v5/internal"
)

// Byte is an nullable byte.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Byte struct {
	sql.NullByte
}

// NewByte creates a new Byte.
func NewByte(b byte, valid bool) Byte {
	return Byte{
		NullByte: sql.NullByte{
			Byte:  b,
			Valid: valid,
		},
	}
}

// ByteFrom creates a new Byte that will always be valid.
func ByteFrom(b byte) Byte {
	return NewByte(b, true)
}

// ByteFromPtr creates a new Byte that be null if i is nil.
func ByteFromPtr(b *byte) Byte {
	if b == nil {
		return NewByte(0, false)
	}
	return NewByte(*b, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (b Byte) ValueOrZero() byte {
	if !b.Valid {
		return 0
	}
	return b.Byte
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number, string, and null input.
// 0 will not be considered a null Byte.
func (b *Byte) UnmarshalJSON(data []byte) error {
	return internal.UnmarshalIntJSON(data, &b.Byte, &b.Valid, 8, strconv.ParseUint)
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Byte if the input is blank.
// It will return an error if the input is not an integer, blank, or "null".
func (b *Byte) UnmarshalText(text []byte) error {
	return internal.UnmarshalIntText(text, &b.Byte, &b.Valid, 8, strconv.ParseUint)
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Byte is null.
func (b Byte) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(int64(b.Byte), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Byte is null.
func (b Byte) MarshalText() ([]byte, error) {
	if !b.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(b.Byte), 10)), nil
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

// IsZero returns true for invalid Bytes, for future omitempty support (Go 1.4?)
// A non-null Byte with a 0 value will not be considered zero.
func (b Byte) IsZero() bool {
	return !b.Valid
}

// Equal returns true if both ints have the same value or are both null.
func (b Byte) Equal(other Byte) bool {
	return b.Valid == other.Valid && (!b.Valid || b.Byte == other.Byte)
}

func (b Byte) value() (int64, bool) {
	return int64(b.Byte), b.Valid
}
