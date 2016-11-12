package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"gopkg.in/nullbio/null.v6/convert"
)

// NullBytes is a global byte slice of JSON null
var NullBytes = []byte("null")

// Bytes is a nullable []byte.
type Bytes struct {
	Bytes []byte
	Valid bool
}

// NewBytes creates a new Bytes
func NewBytes(b []byte, valid bool) Bytes {
	return Bytes{
		Bytes: b,
		Valid: valid,
	}
}

// BytesFrom creates a new Bytes that will be invalid if nil.
func BytesFrom(b []byte) Bytes {
	return NewBytes(b, b != nil)
}

// BytesFromPtr creates a new Bytes that will be invalid if nil.
func BytesFromPtr(b *[]byte) Bytes {
	if b == nil {
		return NewBytes(nil, false)
	}
	n := NewBytes(*b, true)
	return n
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Bytes) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		b.Valid = false
		b.Bytes = nil
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	b.Bytes = []byte(s)
	b.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
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
func (b Bytes) MarshalJSON() ([]byte, error) {
	if len(b.Bytes) == 0 || b.Bytes == nil {
		return NullBytes, nil
	}
	return b.Bytes, nil
}

// MarshalText implements encoding.TextMarshaler.
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
	return !b.Valid
}

// Scan implements the Scanner interface.
func (b *Bytes) Scan(value interface{}) error {
	if value == nil {
		b.Bytes, b.Valid = []byte{}, false
		return nil
	}
	b.Valid = true
	return convert.ConvertAssign(&b.Bytes, value)
}

// Value implements the driver Valuer interface.
func (b Bytes) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Bytes, nil
}
