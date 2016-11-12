package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"gopkg.in/nullbio/null.v6/convert"
)

// Int is an nullable int.
type Int struct {
	Int   int
	Valid bool
}

// NewInt creates a new Int
func NewInt(i int, valid bool) Int {
	return Int{
		Int:   i,
		Valid: valid,
	}
}

// IntFrom creates a new Int that will always be valid.
func IntFrom(i int) Int {
	return NewInt(i, true)
}

// IntFromPtr creates a new Int that be null if i is nil.
func IntFromPtr(i *int) Int {
	if i == nil {
		return NewInt(0, false)
	}
	return NewInt(*i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *Int) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, NullBytes) {
		i.Valid = false
		i.Int = 0
		return nil
	}

	var x int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	i.Int = int(x)
	i.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *Int) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseInt(string(text), 10, 0)
	i.Valid = err == nil
	if i.Valid {
		i.Int = int(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
func (i Int) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return NullBytes, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
func (i Int) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int), 10)), nil
}

// SetValid changes this Int's value and also sets it to be non-null.
func (i *Int) SetValid(n int) {
	i.Int = n
	i.Valid = true
}

// Ptr returns a pointer to this Int's value, or a nil pointer if this Int is null.
func (i Int) Ptr() *int {
	if !i.Valid {
		return nil
	}
	return &i.Int
}

// IsZero returns true for invalid Ints, for future omitempty support (Go 1.4?)
func (i Int) IsZero() bool {
	return !i.Valid
}

// Scan implements the Scanner interface.
func (i *Int) Scan(value interface{}) error {
	if value == nil {
		i.Int, i.Valid = 0, false
		return nil
	}
	i.Valid = true
	return convert.ConvertAssign(&i.Int, value)
}

// Value implements the driver Valuer interface.
func (i Int) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return int64(i.Int), nil
}
