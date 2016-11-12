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

// Uint32 is an nullable uint32.
type Uint32 struct {
	Uint32 uint32
	Valid  bool
}

// NewUint32 creates a new Uint32
func NewUint32(i uint32, valid bool) Uint32 {
	return Uint32{
		Uint32: i,
		Valid:  valid,
	}
}

// Uint32From creates a new Uint32 that will always be valid.
func Uint32From(i uint32) Uint32 {
	return NewUint32(i, true)
}

// Uint32FromPtr creates a new Uint32 that be null if i is nil.
func Uint32FromPtr(i *uint32) Uint32 {
	if i == nil {
		return NewUint32(0, false)
	}
	return NewUint32(*i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *Uint32) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		u.Valid = false
		u.Uint32 = 0
		return nil
	}

	var x uint64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x > math.MaxUint32 {
		return fmt.Errorf("json: %d overflows max uint32 value", x)
	}

	u.Uint32 = uint32(x)
	u.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *Uint32) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseUint(string(text), 10, 32)
	u.Valid = err == nil
	if u.Valid {
		u.Uint32 = uint32(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (u Uint32) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return NullBytes, nil
	}
	return []byte(strconv.FormatUint(uint64(u.Uint32), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
func (u Uint32) MarshalText() ([]byte, error) {
	if !u.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(u.Uint32), 10)), nil
}

// SetValid changes this Uint32's value and also sets it to be non-null.
func (u *Uint32) SetValid(n uint32) {
	u.Uint32 = n
	u.Valid = true
}

// Ptr returns a pointer to this Uint32's value, or a nil pointer if this Uint32 is null.
func (u Uint32) Ptr() *uint32 {
	if !u.Valid {
		return nil
	}
	return &u.Uint32
}

// IsZero returns true for invalid Uint32's, for future omitempty support (Go 1.4?)
func (u Uint32) IsZero() bool {
	return !u.Valid
}

// Scan implements the Scanner interface.
func (u *Uint32) Scan(value interface{}) error {
	if value == nil {
		u.Uint32, u.Valid = 0, false
		return nil
	}
	u.Valid = true
	return convert.ConvertAssign(&u.Uint32, value)
}

// Value implements the driver Valuer interface.
func (u Uint32) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return uint64(u.Uint32), nil
}
