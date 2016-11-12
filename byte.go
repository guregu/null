package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Byte is an nullable int.
type Byte struct {
	Byte  byte
	Valid bool
}

// NewByte creates a new Byte
func NewByte(b byte, valid bool) Byte {
	return Byte{
		Byte:  b,
		Valid: valid,
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

// UnmarshalJSON implements json.Unmarshaler.
func (b *Byte) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(data, NullBytes) {
		b.Valid = false
		b.Byte = 0
		return nil
	}

	var x string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if len(x) > 1 {
		return errors.New("json: cannot convert to byte, text len is greater than one")
	}

	b.Byte = x[0]
	b.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Byte) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		b.Valid = false
		return nil
	}

	if len(text) > 1 {
		return errors.New("text: cannot convert to byte, text len is greater than one")
	}

	b.Valid = true
	b.Byte = text[0]
	return nil
}

// MarshalJSON implements json.Marshaler.
func (b Byte) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return NullBytes, nil
	}
	return []byte{'"', b.Byte, '"'}, nil
}

// MarshalText implements encoding.TextMarshaler.
func (b Byte) MarshalText() ([]byte, error) {
	if !b.Valid {
		return []byte{}, nil
	}
	return []byte{b.Byte}, nil
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
func (b Byte) IsZero() bool {
	return !b.Valid
}

// Scan implements the Scanner interface.
func (b *Byte) Scan(value interface{}) error {
	if value == nil {
		b.Byte, b.Valid = 0, false
		return nil
	}

	val := value.(string)
	if len(val) == 0 {
		b.Valid = false
		b.Byte = 0
		return nil
	}

	b.Valid = true
	b.Byte = byte(val[0])
	return nil
}

// Value implements the driver Valuer interface.
func (b Byte) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return []byte{b.Byte}, nil
}
