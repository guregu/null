package null

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"strconv"
	"unsafe"

	"github.com/philpearl/plenc"
)

// Int is an nullable int64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Int struct {
	sql.NullInt64
}

// NewInt creates a new Int
func NewInt(i int64, valid bool) Int {
	return Int{
		NullInt64: sql.NullInt64{
			Int64: i,
			Valid: valid,
		},
	}
}

// IntFrom creates a new Int that will always be valid.
func IntFrom(i int64) Int {
	return NewInt(i, true)
}

// IntFromPtr creates a new Int that be null if i is nil.
func IntFromPtr(i *int64) Int {
	if i == nil {
		return NewInt(0, false)
	}
	return NewInt(*i, true)
}

type basicInt Int

type stringInt struct {
	Int64 string
	Valid bool
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Int.
// It also supports unmarshalling a sql.NullInt64.
func (i *Int) UnmarshalJSON(data []byte) error {
	var err error
	// Golden path is being passed a integer or null
	if bytes.Compare(data, []byte("null")) == 0 {
		i.Valid = false
		return nil
	}
	// BQ sends numbers as strings. We can strip quotes on simple strings
	if data[0] == '"' {
		data = bytes.Trim(data, `"`)
	}
	if data[0] == '{' {
		// We've been sent a structure. This is not our main-line as we encode
		// to a simple int
		var ii basicInt
		err = json.Unmarshal(data, &ii)
		if err != nil {
			// Try a string version
			var si stringInt
			err = json.Unmarshal(data, &si)
			if err == nil {
				i.Valid = si.Valid
				if si.Valid {
					i.Int64, err = strconv.ParseInt(si.Int64, 10, 64)
					i.Valid = (err == nil)
				}
			}
		} else {
			*i = Int(ii)
		}
	} else {
		i.Int64, err = strconv.ParseInt(*(*string)(unsafe.Pointer(&data)), 10, 64)
		i.Valid = (err == nil)
	}

	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	i.Int64, err = strconv.ParseInt(string(text), 10, 64)
	i.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Int is null.
func (i Int) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return nullLiteral, nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Int is null.
func (i Int) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

// SetValid changes this Int's value and also sets it to be non-null.
func (i *Int) SetValid(n int64) {
	i.Int64 = n
	i.Valid = true
}

// Ptr returns a pointer to this Int's value, or a nil pointer if this Int is null.
func (i Int) Ptr() *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

// IsZero returns true for invalid Ints, for future omitempty support (Go 1.4?)
// A non-null Int with a 0 value will not be considered zero.
func (i Int) IsZero() bool {
	return !i.Valid
}

// ΦλSizeFull determines how many bytes are needed to encode this value
func (i Int) ΦλSizeFull(index int) (size int) {
	if !i.Valid {
		return 0
	}
	// We're going to cheat and assume this will always only include a single value. So we won't do the tag
	// thing
	return plenc.SizeTag(plenc.WTVarInt, index) + plenc.SizeVarInt(i.Int64)
}

// ΦλAppendFull encodes example by appending to data. It returns the final slice
func (i Int) ΦλAppendFull(data []byte, index int) []byte {
	if !i.Valid {
		return data
	}
	data = plenc.AppendTag(data, plenc.WTVarInt, index)
	return plenc.AppendVarInt(data, i.Int64)
}

// ΦλUnmarshal decodes a plenc encoded value
func (i *Int) ΦλUnmarshal(data []byte) (int, error) {
	// There's no tag within the encoding. If we're being asked to decode, then this value field must be present
	// within the encoded data,
	i.Valid = true
	var n int
	i.Int64, n = plenc.ReadVarInt(data)
	return n, nil
}
