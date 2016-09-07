package null

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

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

// Unmarshal will unmarshal your JSON stored in
// your JSON object and store the result in the
// value pointed to by dest.
func (j JSON) Unmarshal(dest interface{}) error {
	if dest == nil {
		return errors.New("destination is nil, not a valid pointer to an object")
	}

	// Call our implementation of
	// JSON MarshalJSON through json.Marshal
	// to get the value of the JSON object
	res, err := json.Marshal(j)
	if err != nil {
		return err
	}

	return json.Unmarshal(res, dest)
}

// UnmarshalJSON implements json.Unmarshaler.
// If data is len 0 or nil, it will unmarshal to JSON null.
// If not, it will copy your data slice into JSON.
func (j *JSON) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		j.JSON = []byte("null")
	} else {
		j.JSON = append(j.JSON[0:0], data...)
	}

	j.Valid = true

	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to nil if the text is nil or len 0.
func (j *JSON) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		j.JSON = nil
		j.Valid = false
	} else {
		j.JSON = append(j.JSON[0:0], text...)
		j.Valid = true
	}

	return nil
}

// Marshal will marshal the passed in object,
// and store it in the JSON member on the JSON object.
func (j *JSON) Marshal(obj interface{}) error {
	res, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	// Call our implementation of
	// JSON UnmarshalJSON through json.Unmarshal
	// to set the result to the JSON object
	return json.Unmarshal(res, j)
}

// MarshalJSON implements json.Marshaler.
// It will encode null if the JSON is nil.
func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j.JSON) == 0 || j.JSON == nil {
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
