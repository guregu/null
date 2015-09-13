package zero

import (
	"encoding/json"
	"testing"
	"time"
)

var (
	timeString    = "2012-12-21T21:21:21Z"
	timeJSON      = []byte(`"` + timeString + `"`)
	zeroTimeJSON  = []byte(`"0001-01-01T00:00:00Z"`)
	blankTimeJSON = []byte(`null`)
	timeValue, _  = time.Parse(time.RFC3339, timeString)
	timeMap       = []byte(`{"Time":"2012-12-21T21:21:21Z","Valid":true}`)
	invalidMap    = []byte(`{"Time":"0001-01-01T00:00:00Z","Valid":false}`)
)

func TestUnmarshalTimeString(t *testing.T) {
	var ti Time
	err := json.Unmarshal(timeMap, &ti)
	maybePanic(err)
	assertTime(t, ti, "UnmarshalJSON() json")

	var blank Time
	err = json.Unmarshal(blankTimeJSON, &blank)
	maybePanic(err)
	assertNullTime(t, blank, "blank time json")

	var zero Time
	err = json.Unmarshal(zeroTimeJSON, &zero)
	maybePanic(err)
	assertNullTime(t, zero, "zero time json")

	var fromMap Time
	err = json.Unmarshal(timeMap, &fromMap)
	maybePanic(err)
	assertTime(t, fromMap, "map time json")

	var invalid Time
	err = json.Unmarshal(invalidMap, &invalid)
	maybePanic(err)
	assertNullTime(t, invalid, "map invalid time json")
}

func TestMarshalTime(t *testing.T) {
	ti := TimeFrom(timeValue)
	data, err := json.Marshal(ti)
	maybePanic(err)
	assertJSONEquals(t, data, string(timeJSON), "non-empty json marshal")

	null := TimeFromPtr(nil)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, string(zeroTimeJSON), "empty json marshal")
}

func TestTimeFrom(t *testing.T) {
	ti := TimeFrom(timeValue)
	assertTime(t, ti, "TimeFrom() time.Time")

	var nt time.Time
	null := TimeFrom(nt)
	assertNullTime(t, null, "TimeFrom() empty time.Time")
}

func TestTimeFromPtr(t *testing.T) {
	ti := TimeFromPtr(&timeValue)
	assertTime(t, ti, "TimeFromPtr() time")

	null := TimeFromPtr(nil)
	assertNullTime(t, null, "TimeFromPtr(nil)")
}

func TestTimeSetValid(t *testing.T) {
	var ti time.Time
	change := TimeFrom(ti)
	assertNullTime(t, change, "SetValid()")
	change.SetValid(timeValue)
	assertTime(t, change, "SetValid()")
}

func TestTimePointer(t *testing.T) {
	ti := TimeFrom(timeValue)
	ptr := ti.Ptr()
	if *ptr != timeValue {
		t.Errorf("bad %s time: %#v ≠ %v\n", "pointer", ptr, timeValue)
	}

	var nt time.Time
	null := TimeFrom(nt)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s time: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestTimeScan(t *testing.T) {
	var ti Time
	err := ti.Scan(timeValue)
	maybePanic(err)
	assertTime(t, ti, "scanned time")

	var null Time
	err = null.Scan(nil)
	maybePanic(err)
	assertNullTime(t, null, "scanned null")
}

func TestTimeValue(t *testing.T) {
	ti := TimeFrom(timeValue)
	v, err := ti.Value()
	maybePanic(err)
	if ti.Time != timeValue {
		t.Errorf("bad time.Time value: %v ≠ %v", ti.Time, timeValue)
	}

	var nt time.Time
	zero := TimeFrom(nt)
	v, err = zero.Value()
	maybePanic(err)
	if v != nil {
		t.Errorf("bad %s time.Time value: %v ≠ %v", "zero", v, nil)
	}
}

func assertTime(t *testing.T, ti Time, from string) {
	if ti.Time != timeValue {
		t.Errorf("bad %v time: %v ≠ %v\n", from, ti.Time, timeValue)
	}
	if !ti.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullTime(t *testing.T, ti Time, from string) {
	if ti.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
