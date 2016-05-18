package zero

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/pobri19/null-extended/convert"
)

type NullInt struct {
	Int   int
	Valid bool
}

// Int is a nullable int.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
type Int struct {
	NullInt
}

// NewInt creates a new Int
func NewInt(i int, valid bool) Int {
	return Int{
		NullInt: NullInt{
			Int:   i,
			Valid: valid,
		},
	}
}

// IntFrom creates a new Int that will be null if zero.
func IntFrom(i int) Int {
	return NewInt(i, i != 0)
}

// IntFromPtr creates a new Int that be null if i is nil.
func IntFromPtr(i *int) Int {
	if i == nil {
		return NewInt(0, false)
	}
	n := NewInt(*i, true)
	return n
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will be considered a null Int.
// It also supports unmarshalling a sql.NullInt.
func (i *Int) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Int)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullInt)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type zero.Int", reflect.TypeOf(v).Name())
	}
	i.Valid = (err == nil) && (i.Int != 0)
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank, zero, or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseInt(string(text), 10, 0)
	i.Int = int(res)
	i.Valid = (err == nil) && (i.Int != 0)
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode 0 if this Int is null.
func (i Int) MarshalJSON() ([]byte, error) {
	n := i.Int
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(int64(n), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Int is null.
func (i Int) MarshalText() ([]byte, error) {
	n := i.Int
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(int64(n), 10)), nil
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

// IsZero returns true for null or zero Ints, for future omitempty support (Go 1.4?)
func (i Int) IsZero() bool {
	return !i.Valid || i.Int == 0
}

// Scan implements the Scanner interface.
func (n *NullInt) Scan(value interface{}) error {
	if value == nil {
		n.Int, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Int, value)
}

// Value implements the driver Valuer interface.
func (n NullInt) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int, nil
}
