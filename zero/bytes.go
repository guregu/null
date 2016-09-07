package zero

import (
	"database/sql/driver"

	"gopkg.in/nullbio/null.v5/convert"
)

// NullBytes is a nullable byte slice.
type NullBytes struct {
	Bytes []byte
	Valid bool
}

// Bytes is a nullable []byte.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
type Bytes struct {
	NullBytes
}

// NewBytes creates a new Bytes
func NewBytes(b []byte, valid bool) Bytes {
	return Bytes{
		NullBytes: NullBytes{
			Bytes: b,
			Valid: valid,
		},
	}
}

// BytesFrom creates a new Bytes that will be null if len zero.
func BytesFrom(b []byte) Bytes {
	return NewBytes(b, len(b) != 0)
}

// BytesFromPtr creates a new Bytes that be null if len zero.
func BytesFromPtr(b *[]byte) Bytes {
	if b == nil || len(*b) == 0 {
		return NewBytes(nil, false)
	}
	n := NewBytes(*b, true)
	return n
}

// UnmarshalJSON implements json.Unmarshaler.
// If data is len 0 or nil, it will unmarshal to JSON null.
// If not, it will copy your data slice into Bytes.
func (b *Bytes) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		b.Bytes = []byte("null")
		b.Valid = false
	} else {
		b.Bytes = append(b.Bytes[0:0], data...)
		b.Valid = true
	}

	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to nil if the text is nil or len 0.
func (b *Bytes) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		b.Bytes = nil
		b.Valid = false
	} else {
		b.Bytes = append(b.Bytes[0:0], text...)
		b.Valid = true
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if the Bytes is nil.
func (b Bytes) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	return b.Bytes, nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode nil if the Bytes is invalid.
func (b Bytes) MarshalText() ([]byte, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Bytes, nil
}

// SetValid changes this Bytes's value and also sets it to be non-null.
func (b *Bytes) SetValid(n []byte) {
	b.Bytes = n
	b.Valid = true
}

// Ptr returns a pointer to this Bytes's value, or a nil pointer if this Bytes is null.
func (b Bytes) Ptr() *[]byte {
	if !b.Valid {
		return nil
	}
	return &b.Bytes
}

// IsZero returns true for null or zero Bytes's, for future omitempty support (Go 1.4?)
func (b Bytes) IsZero() bool {
	return !b.Valid || b.Bytes == nil || len(b.Bytes) == 0
}

// Scan implements the Scanner interface.
func (n *NullBytes) Scan(value interface{}) error {
	if value == nil {
		n.Bytes, n.Valid = []byte{}, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Bytes, value)
}

// Value implements the driver Valuer interface.
func (n NullBytes) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bytes, nil
}
