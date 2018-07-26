package null

import (
	"encoding/json"
	"testing"
	"time"
)

var (
	timestampString     = "1356124881"
	timestampJSON       = []byte(timestampString)
	timestampValue      = time.Unix(1356124881, 0)
	timestampObject     = []byte(`{"Time":1356124881,"Valid":true}`)
	timestampNullObject = []byte(`{"Time":0,"Valid":false}`)
)

func TestUnmarshalTimestampJSON(t *testing.T) {
	var ti Timestamp
	err := json.Unmarshal(timestampJSON, &ti)
	maybePanic(err)
	assertTimestamp(t, ti, "UnmarshalJSON() json")

	var null Timestamp
	err = json.Unmarshal(nullTimeJSON, &null)
	maybePanic(err)
	assertNullTimestamp(t, null, "null time json")

	var fromObject Timestamp
	err = json.Unmarshal(timestampObject, &fromObject)
	maybePanic(err)
	assertTimestamp(t, fromObject, "time from object json")

	var nullFromObj Timestamp
	err = json.Unmarshal(timestampNullObject, &nullFromObj)
	maybePanic(err)
	assertNullTimestamp(t, nullFromObj, "null from object json")

	var invalid Timestamp
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullTimestamp(t, invalid, "invalid from object json")

	var bad Timestamp
	err = json.Unmarshal(badObject, &bad)
	if err == nil {
		t.Errorf("expected error: bad object")
	}
	assertNullTimestamp(t, bad, "bad from object json")

	var wrongType Timestamp
	err = json.Unmarshal(timeJSON, &wrongType)
	if err == nil {
		t.Errorf("expected error: wrong type JSON")
	}
	assertNullTimestamp(t, wrongType, "wrong type object json")
}

func TestUnmarshalTimestampText(t *testing.T) {
	ti := TimestampFrom(timestampValue)
	txt, err := ti.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, timestampString, "marshal text")

	var unmarshal Timestamp
	err = unmarshal.UnmarshalText(txt)
	maybePanic(err)
	assertTimestamp(t, unmarshal, "unmarshal text")

	var null Timestamp
	err = null.UnmarshalText(nullJSON)
	maybePanic(err)
	assertNullTimestamp(t, null, "unmarshal null text")
	txt, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, string(nullJSON), "marshal null text")

	var invalid Timestamp
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		t.Error("expected error")
	}
	assertNullTimestamp(t, invalid, "bad string")
}

func TestMarshalTimestamp(t *testing.T) {
	ti := TimestampFrom(timestampValue)
	data, err := json.Marshal(ti)
	maybePanic(err)
	assertJSONEquals(t, data, string(timestampJSON), "non-empty json marshal")

	ti.Valid = false
	data, err = json.Marshal(ti)
	maybePanic(err)
	assertJSONEquals(t, data, string(nullJSON), "null json marshal")
}

func TestTimestampFrom(t *testing.T) {
	ti := TimestampFrom(timestampValue)
	assertTimestamp(t, ti, "TimeFrom() time.Time")
}

func TestTimestampFromPtr(t *testing.T) {
	ti := TimestampFromPtr(&timestampValue)
	assertTimestamp(t, ti, "TimeFromPtr() time")

	null := TimestampFromPtr(nil)
	assertNullTimestamp(t, null, "TimeFromPtr(nil)")
}

func TestTimestampSetValid(t *testing.T) {
	var ti time.Time
	change := NewTimestamp(ti, false)
	assertNullTimestamp(t, change, "SetValid()")
	change.SetValid(timestampValue)
	assertTimestamp(t, change, "SetValid()")
}

func TestTimestampPointer(t *testing.T) {
	ti := TimestampFrom(timestampValue)
	ptr := ti.Ptr()
	if *ptr != timestampValue {
		t.Errorf("bad %s time: %#v ≠ %v\n", "pointer", ptr, timestampValue)
	}

	var nt time.Time
	null := NewTimestamp(nt, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s time: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestTimestampScanValue(t *testing.T) {
	var ti Timestamp
	err := ti.Scan(timestampValue)
	maybePanic(err)
	assertTimestamp(t, ti, "scanned time")
	if v, err := ti.Value(); v != timestampValue || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var null Timestamp
	err = null.Scan(nil)
	maybePanic(err)
	assertNullTimestamp(t, null, "scanned null")
	if v, err := null.Value(); v != nil || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var wrong Timestamp
	err = wrong.Scan(int64(42))
	if err == nil {
		t.Error("expected error")
	}
	assertNullTimestamp(t, wrong, "scanned wrong")
}

func TestTimestampValueOrZero(t *testing.T) {
	valid := TimestampFrom(timestampValue)
	if valid.ValueOrZero() != valid.Time || valid.ValueOrZero().IsZero() {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := valid
	invalid.Valid = false
	if !invalid.ValueOrZero().IsZero() {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func assertTimestamp(t *testing.T, ti Timestamp, from string) {
	if ti.Time != timestampValue {
		t.Errorf("bad %v time: %v ≠ %v\n", from, ti.Time, timestampValue)
	}
	if !ti.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullTimestamp(t *testing.T, ti Timestamp, from string) {
	if ti.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
