package null

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

// Time is a nullable time.Time. It supports SQL and JSON serialization.
// It will marshal to null if null.
type Time struct {
	sql.NullTime
}

// Value implements the driver Valuer interface.
func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// NewTime creates a new Time.
func NewTime(t time.Time, valid bool) Time {
	return Time{
		NullTime: sql.NullTime{
			Time:  t,
			Valid: valid,
		},
	}
}

// T creates a new Time that will always be valid.
func T(t time.Time) Time {
	return TimeFrom(t)
}

// TimeFrom creates a new Time that will always be valid.
func TimeFrom(t time.Time) Time {
	return NewTime(t, true)
}

// TimeFromPtr creates a new Time that will be null if t is nil.
func TimeFromPtr(t *time.Time) Time {
	if t == nil {
		return NewTime(time.Time{}, false)
	}
	return NewTime(*t, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (t Time) ValueOrZero() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return nullLiteral, nil
	}
	return t.Time.MarshalJSON()
}

func (t Time) MarshalEasyJSON(w *jwriter.Writer) {
	if !t.Valid {
		w.RawString("null")
		return
	}
	w.Buffer.EnsureSpace(len(time.RFC3339Nano) + 2)
	w.Buffer.AppendByte('"')
	w.Buffer.Buf = t.Time.UTC().AppendFormat(w.Buffer.Buf, time.RFC3339Nano)
	w.Buffer.AppendByte('"')
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input.
func (t *Time) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullLiteral) {
		t.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &t.Time); err != nil {
		return fmt.Errorf("null: couldn't unmarshal JSON: %w", err)
	}

	t.Valid = true
	return nil
}

// UnmarshalEasyJSON is an easy-JSON specific decoder, that should be more efficient than the standard one.
func (i *Time) UnmarshalEasyJSON(w *jlexer.Lexer) {
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
			case "time", "Time":
				t, err := time.Parse(time.RFC3339, w.UnsafeString())
				if err != nil {
					w.AddError(err)
				}
				i.Time = t
			case "valid", "Valid":
				i.Valid = w.Bool()
			}
			w.WantComma()
		}
		return
	}
	t, err := time.Parse(time.RFC3339, w.UnsafeString())
	if err != nil {
		w.AddError(err)
	}
	i.Time = t
	i.Valid = (w.Error() == nil)
}

// MarshalText implements encoding.TextMarshaler.
// It returns an empty string if invalid, otherwise time.Time's MarshalText.
func (t Time) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte{}, nil
	}
	return t.Time.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It has backwards compatibility with v3 in that the string "null" is considered equivalent to an empty string
// and unmarshaling will succeed. This may be removed in a future version.
func (t *Time) UnmarshalText(text []byte) error {
	str := string(text)
	// allowing "null" is for backwards compatibility with v3
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	if err := t.Time.UnmarshalText(text); err != nil {
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}
	t.Valid = true
	return nil
}

// SetValid changes this Time's value and sets it to be non-null.
func (t *Time) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

// Ptr returns a pointer to this Time's value, or a nil pointer if this Time is null.
func (t Time) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func (t Time) IsDefined() bool {
	return !t.IsZero()
}

// IsZero returns true for invalid Times, hopefully for future omitempty support.
// A non-null Time with a zero value will not be considered zero.
func (t Time) IsZero() bool {
	return !t.Valid
}

// Equal returns true if both Time objects encode the same time or are both null.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 CEST and 4:00 UTC are Equal.
func (t Time) Equal(other Time) bool {
	return t.Valid == other.Valid && (!t.Valid || t.Time.Equal(other.Time))
}

// ExactEqual returns true if both Time objects are equal or both null.
// ExactEqual returns false for times that are in different locations or
// have a different monotonic clock reading.
func (t Time) ExactEqual(other Time) bool {
	return t.Valid == other.Valid && (!t.Valid || t.Time == other.Time)
}
