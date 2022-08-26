package null

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgtype"
)

// Interval is an extension of pgtype.Interval that adds JSON serialization support.
// It will marshal to null if the status is "Null" or "Undefined".
type Interval struct {
	pgtype.Interval
}

// IntervalFrom creates a new Interval from the given pgtype.Interval.
func IntervalFrom(i pgtype.Interval) Interval {
	return Interval{i}
}

// IntervalFromPtr creates a new Interval that is null if i is nil.
func IntervalFromPtr(i *pgtype.Interval) Interval {
	if i == nil {
		return Interval{Interval: pgtype.Interval{Status: pgtype.Null}}
	}
	return IntervalFrom(*i)
}

// ValueOrZero returns the inner value if present, otherwise a zero pgtype.Interval.
func (i Interval) ValueOrZero() pgtype.Interval {
	if i.Status != pgtype.Present {
		return pgtype.Interval{}
	}
	return i.Interval
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *Interval) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		i.Status = pgtype.Null
		return nil
	}

	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return json.Unmarshal(data, &i.Interval)
	}

	if err := i.DecodeText(nil, data[1:len(data)-1]); err != nil {
		return fmt.Errorf("null: couldn't unmarshal JSON: %w", err)
	}

	i.Status = pgtype.Present
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Interval has the status "Null" or "Undefined".
func (i Interval) MarshalJSON() ([]byte, error) {
	switch i.Status {
	case pgtype.Null, pgtype.Undefined:
		return []byte("null"), nil
	case pgtype.Present:
		r, err := i.EncodeText(nil, []byte{})
		if err != nil {
			return nil, err
		}
		return json.Marshal(string(r))
	}
	return nil, fmt.Errorf("Unsupported Interval status %d", i.Status)
}

// Ptr returns a pointer to this Interval's value, or a nil pointer if this Interval doesn't have the "Present" status.
func (i Interval) Ptr() *pgtype.Interval {
	if i.Status != pgtype.Present {
		return nil
	}
	return &i.Interval
}

// IsZero returns true for Interval with "Null" status, for potential future omitempty support.
func (i Interval) IsZero() bool {
	return i.Status == pgtype.Undefined
}

// Equal returns true if both Intervals have the same value or are both null.
func (i Interval) Equal(other Interval) bool {
	return i.Status == other.Status && i.Months == other.Months && i.Days == other.Days && i.Microseconds == other.Microseconds
}
