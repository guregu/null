package null

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"unsafe"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
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

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Int) ValueOrZero() int64 {
	if !i.Valid {
		return 0
	}
	return i.Int64
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number, string, and null input.
// 0 will not be considered a null Int.
func (i *Int) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullLiteral) {
		i.Valid = false
		return nil
	}

	if data[0] == '{' {
		// Try the struct form of Int.
		type basicInt Int
		var ii basicInt
		if json.Unmarshal(data, &ii) == nil {
			*i = Int(ii)
			return nil
		}

		// Try a string version
		var si struct {
			Int64 string
			Valid bool
		}
		if err := json.Unmarshal(data, &si); err != nil {
			return err
		}
		i.Valid = si.Valid
		if !si.Valid {
			return nil
		}
		var err error
		i.Int64, err = strconv.ParseInt(si.Int64, 10, 64)
		i.Valid = (err == nil)
		return err
	}

	if data[0] == '"' {
		data = bytes.Trim(data, `"`)
	}
	var err error
	i.Int64, err = strconv.ParseInt(*(*string)(unsafe.Pointer(&data)), 10, 64)
	i.Valid = (err == nil)
	return err
}

// UnmarshalEasyJSON is an easy-JSON specific decoder, that should be more efficient than the standard one.
func (i *Int) UnmarshalEasyJSON(w *jlexer.Lexer) {
	if w.IsNull() {
		w.Skip()
		i.Valid = false
		return
	}
	if w.IsDelim('{') {
		w.Skip()
		for !w.IsDelim('}') {
			key := w.UnsafeString()
			w.WantColon()
			if w.IsNull() {
				w.Skip()
				w.WantComma()
				continue
			}
			switch key {
			case "int64", "Int64":
				v, err := w.JsonNumber().Int64()
				if err != nil {
					w.AddError(err)
					i.Valid = false
					return
				}
				i.Int64 = v
			case "valid", "Valid":
				i.Valid = w.Bool()
			}
			w.WantComma()
		}
		return
	}
	data := w.Raw()
	if data[0] == '"' {
		data = data[1 : len(data)-1]
	}
	ii, err := strconv.ParseInt(*(*string)(unsafe.Pointer(&data)), 10, 64)
	if err != nil {
		w.AddError(&jlexer.LexerError{
			Reason: err.Error(),
			Data:   string(data),
		})
	}
	i.Int64 = ii
	i.Valid = (err == nil)
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is blank.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	i.Int64, err = strconv.ParseInt(string(text), 10, 64)
	if err != nil {
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}
	i.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Int is null.
func (i Int) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return nullLiteral, nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

func (i Int) MarshalEasyJSON(w *jwriter.Writer) {
	if !i.Valid {
		w.RawString("null")
		return
	}
	w.Int64(i.Int64)
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

func (i Int) IsDefined() bool {
	return !i.IsZero()
}

// Equal returns true if both ints have the same value or are both null.
func (i Int) Equal(other Int) bool {
	return i.Valid == other.Valid && (!i.Valid || i.Int64 == other.Int64)
}
