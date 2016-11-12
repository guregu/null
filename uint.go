package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"gopkg.in/nullbio/null.v6/convert"
)

// Uint is an nullable uint.
type Uint struct {
	Uint  uint
	Valid bool
}

// NewUint creates a new Uint
func NewUint(i uint, valid bool) Uint {
	return Uint{
		Uint:  i,
		Valid: valid,
	}
}

// UintFrom creates a new Uint that will always be valid.
func UintFrom(i uint) Uint {
	return NewUint(i, true)
}

// UintFromPtr creates a new Uint that be null if i is nil.
func UintFromPtr(i *uint) Uint {
	if i == nil {
		return NewUint(0, false)
	}
	return NewUint(*i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *Uint) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		u.Valid = false
		u.Uint = 0
		return nil
	}

	var x uint64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	u.Uint = uint(x)
	u.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *Uint) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseUint(string(text), 10, 0)
	u.Valid = err == nil
	if u.Valid {
		u.Uint = uint(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (u Uint) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return NullBytes, nil
	}
	return []byte(strconv.FormatUint(uint64(u.Uint), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
func (u Uint) MarshalText() ([]byte, error) {
	if !u.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(u.Uint), 10)), nil
}

// SetValid changes this Uint's value and also sets it to be non-null.
func (u *Uint) SetValid(n uint) {
	u.Uint = n
	u.Valid = true
}

// Ptr returns a pointer to this Uint's value, or a nil pointer if this Uint is null.
func (u Uint) Ptr() *uint {
	if !u.Valid {
		return nil
	}
	return &u.Uint
}

// IsZero returns true for invalid Uints, for future omitempty support (Go 1.4?)
func (u Uint) IsZero() bool {
	return !u.Valid
}

// Scan implements the Scanner interface.
func (u *Uint) Scan(value interface{}) error {
	if value == nil {
		u.Uint, u.Valid = 0, false
		return nil
	}
	u.Valid = true
	return convert.ConvertAssign(&u.Uint, value)
}

// Value implements the driver Valuer interface.
func (u Uint) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return int64(u.Uint), nil
}
