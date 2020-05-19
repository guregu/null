package null

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"strconv"
	"unsafe"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

// Float is a nullable float64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Float struct {
	sql.NullFloat64
}

// NewFloat creates a new Float
func NewFloat(f float64, valid bool) Float {
	return Float{
		NullFloat64: sql.NullFloat64{
			Float64: f,
			Valid:   valid,
		},
	}
}

// FloatFrom creates a new Float that will always be valid.
func FloatFrom(f float64) Float {
	return NewFloat(f, true)
}

// FloatFromPtr creates a new Float that be null if f is nil.
func FloatFromPtr(f *float64) Float {
	if f == nil {
		return NewFloat(0, false)
	}
	return NewFloat(*f, true)
}

type basicFloat Float

type stringFloat struct {
	Float64 string
	Valid   bool
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Float.
// It also supports unmarshalling a sql.NullFloat64.
func (i *Float) UnmarshalJSON(data []byte) error {
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
		// to a simple float
		var ii basicFloat
		err = json.Unmarshal(data, &ii)
		if err != nil {
			// Try a string version
			var si stringFloat
			err = json.Unmarshal(data, &si)
			if err == nil {
				i.Valid = si.Valid
				if si.Valid {
					i.Float64, err = strconv.ParseFloat(si.Float64, 64)
					i.Valid = (err == nil)
				}
			}
		} else {
			*i = Float(ii)
		}
	} else {
		i.Float64, err = strconv.ParseFloat(*(*string)(unsafe.Pointer(&data)), 64)
		i.Valid = (err == nil)
	}

	return err
}

// UnmarshalEasyJSON is an easy-JSON specific decoder, that should be more efficient than the standard one.
func (i *Float) UnmarshalEasyJSON(w *jlexer.Lexer) {
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
			case "float64", "Float64":
				i.Float64 = w.Float64()
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
	f, err := strconv.ParseFloat(*(*string)(unsafe.Pointer(&data)), 64)
	if err != nil {
		w.AddError(&jlexer.LexerError{
			Reason: err.Error(),
			Data:   string(data),
		})
	}
	i.Float64 = f
	i.Valid = (w.Error() == nil)
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Float if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (f *Float) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		f.Valid = false
		return nil
	}
	var err error
	f.Float64, err = strconv.ParseFloat(string(text), 64)
	f.Valid = err == nil
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Float is null.
func (f Float) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return nullLiteral, nil
	}
	return []byte(strconv.FormatFloat(f.Float64, 'f', -1, 64)), nil
}

func (i Float) MarshalEasyJSON(w *jwriter.Writer) {
	if !i.Valid {
		w.RawString("null")
		return
	}
	w.Float64(i.Float64)
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Float is null.
func (f Float) MarshalText() ([]byte, error) {
	if !f.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatFloat(f.Float64, 'f', -1, 64)), nil
}

// SetValid changes this Float's value and also sets it to be non-null.
func (f *Float) SetValid(n float64) {
	f.Float64 = n
	f.Valid = true
}

// Ptr returns a pointer to this Float's value, or a nil pointer if this Float is null.
func (f Float) Ptr() *float64 {
	if !f.Valid {
		return nil
	}
	return &f.Float64
}

// IsZero returns true for invalid Floats, for future omitempty support (Go 1.4?)
// A non-null Float with a 0 value will not be considered zero.
func (f Float) IsZero() bool {
	return !f.Valid
}

func (f Float) IsDefined() bool {
	return f.Valid
}
