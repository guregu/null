package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"gopkg.in/nullbio/null.v6/convert"
)

// Uint16 is an nullable uint16.
type Uint16 struct {
	Uint16 uint16
	Valid  bool
}

// NewUint16 creates a new Uint16
func NewUint16(i uint16, valid bool) Uint16 {
	return Uint16{
		Uint16: i,
		Valid:  valid,
	}
}

// Uint16From creates a new Uint16 that will always be valid.
func Uint16From(i uint16) Uint16 {
	return NewUint16(i, true)
}

// Uint16FromPtr creates a new Uint16 that be null if i is nil.
func Uint16FromPtr(i *uint16) Uint16 {
	if i == nil {
		return NewUint16(0, false)
	}
	return NewUint16(*i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *Uint16) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		u.Valid = false
		u.Uint16 = 0
		return nil
	}

	var x uint64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x > math.MaxUint16 {
		return fmt.Errorf("json: %d overflows max uint8 value", x)
	}

	u.Uint16 = uint16(x)
	u.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *Uint16) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseUint(string(text), 10, 16)
	u.Valid = err == nil
	if u.Valid {
		u.Uint16 = uint16(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (u Uint16) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return NullBytes, nil
	}
	return []byte(strconv.FormatUint(uint64(u.Uint16), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
func (u Uint16) MarshalText() ([]byte, error) {
	if !u.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(u.Uint16), 10)), nil
}

// SetValid changes this Uint16's value and also sets it to be non-null.
func (u *Uint16) SetValid(n uint16) {
	u.Uint16 = n
	u.Valid = true
}

// Ptr returns a pointer to this Uint16's value, or a nil pointer if this Uint16 is null.
func (u Uint16) Ptr() *uint16 {
	if !u.Valid {
		return nil
	}
	return &u.Uint16
}

// IsZero returns true for invalid Uint16's, for future omitempty support (Go 1.4?)
func (u Uint16) IsZero() bool {
	return !u.Valid
}

// Scan implements the Scanner interface.
func (u *Uint16) Scan(value interface{}) error {
	if value == nil {
		u.Uint16, u.Valid = 0, false
		return nil
	}
	u.Valid = true
	return convert.ConvertAssign(&u.Uint16, value)
}

// Value implements the driver Valuer interface.
func (u Uint16) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return int64(u.Uint16), nil
}
