package null

import (
	"database/sql/driver"

	"gopkg.in/nullbio/null.v4/convert"
)

// NullJSON is a nullable byte slice.
type NullJSON struct {
	JSON  []byte
	Valid bool
}

// JSON is a nullable []byte.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
type JSON struct {
	NullJSON
}

// NewJSON creates a new JSON
func NewJSON(b []byte, valid bool) JSON {
	return JSON{
		NullJSON: NullJSON{
			JSON:  b,
			Valid: valid,
		},
	}
}

// JSONFrom creates a new JSON that will be null if len zero.
func JSONFrom(b []byte) JSON {
	return NewJSON(b, len(b) != 0)
}

// JSONFromPtr creates a new JSON that be null if len zero.
func JSONFromPtr(b *[]byte) JSON {
	if b == nil || len(*b) == 0 {
		return NewJSON(nil, false)
	}
	n := NewJSON(*b, true)
	return n
}

// UnmarshalJSON implements json.Unmarshaler.
// JSON UnmarshalJSON is different in that it only
// unmarshals sql.NullJSON defined as JSON objects,
// It supports all JSON types.
// It also supports unmarshalling a sql.NullJSON.
func (j *JSON) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		j.JSON = nil
		j.Valid = false
	} else {
		j.JSON = append(j.JSON[0:0], data...)
		j.Valid = true
	}

	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null JSON if the input is blank.
// It will return an error if the input is not an integer, blank, or "null".
func (j *JSON) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		j.Valid = false
	} else {
		j.JSON = append(j.JSON[0:0], text...)
		j.Valid = true
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if the JSON is invalid.
func (j JSON) MarshalJSON() ([]byte, error) {
	if !j.Valid {
		return []byte("null"), nil
	}
	return j.JSON, nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode nil if the JSON is invalid.
func (j JSON) MarshalText() ([]byte, error) {
	if !j.Valid {
		return nil, nil
	}
	return j.JSON, nil
}

// SetValid changes this JSON's value and also sets it to be non-null.
func (j *JSON) SetValid(n []byte) {
	j.JSON = n
	j.Valid = true
}

// Ptr returns a pointer to this JSON's value, or a nil pointer if this JSON is null.
func (j JSON) Ptr() *[]byte {
	if !j.Valid {
		return nil
	}
	return &j.JSON
}

// IsZero returns true for null or zero JSON's, for future omitempty support (Go 1.4?)
func (j JSON) IsZero() bool {
	return !j.Valid
}

// Scan implements the Scanner interface.
func (n *NullJSON) Scan(value interface{}) error {
	if value == nil {
		n.JSON, n.Valid = []byte{}, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.JSON, value)
}

// Value implements the driver Valuer interface.
func (n NullJSON) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.JSON, nil
}
