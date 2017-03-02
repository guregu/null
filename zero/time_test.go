package zero

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/axiomzen/null/format"
)

var (
	timeString   = "2012-12-21T21:21:21Z"
	timeJSON     = []byte(`"` + timeString + `"`)
	zeroTimeStr  = "0001-01-01T00:00:00Z"
	zeroTimeJSON = []byte(`"` + zeroTimeStr + `"`)
	timeValue, _ = time.Parse(time.RFC3339, timeString)
	timeObject   = []byte(`{"Time":"` + timeString + `","Valid":true}`)
	nullObject   = []byte(`{"Time":"` + zeroTimeStr + `","Valid":false}`)

	blankTimeJSON = []byte(`null`)
	badObject     = []byte(`{"hello": "world"}`)

	customTimeFormat   = "2006-01-02T15:04:05Z0700"
	customTimeString   = "2012-12-21T21:21:21-0200"
	customTimeJSON     = []byte(`"` + customTimeString + `"`)
	customZeroTimeStr  = "0001-01-01T00:00:00Z"
	customZeroTimeJSON = []byte(`"` + customZeroTimeStr + `"`)
	customTimeValue, _ = time.Parse(customTimeFormat, customTimeString)
	customTimeObject   = []byte(`{"Time":"` + customTimeString + `","Valid":true}`)
	customNullObject   = []byte(`{"Time":"` + customZeroTimeStr + `","Valid":false}`)
)

func testUnmarshalTimeJSON(t *testing.T, f string, to []byte, value time.Time, zeroT []byte) {
	format.SetTimeFormat(f)
	var ti Time
	err := json.Unmarshal(to, &ti)
	maybePanic(err)
	assertTime(t, ti, "UnmarshalJSON() json", value)

	var blank Time
	err = json.Unmarshal(blankTimeJSON, &blank)
	maybePanic(err)
	assertNullTime(t, blank, "blank time json")

	var zero Time
	err = json.Unmarshal(zeroT, &zero)
	maybePanic(err)
	assertNullTime(t, zero, "zero time json")

	var fromObject Time
	err = json.Unmarshal(to, &fromObject)
	maybePanic(err)
	assertTime(t, fromObject, "map time json", value)

	var null Time
	err = json.Unmarshal(nullObject, &null)
	maybePanic(err)
	assertNullTime(t, null, "map null time json")

	var nullFromObj Time
	err = json.Unmarshal(nullObject, &nullFromObj)
	maybePanic(err)
	assertNullTime(t, nullFromObj, "null from object json")

	var invalid Time
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullTime(t, invalid, "invalid from object json")

	var bad Time
	err = json.Unmarshal(badObject, &bad)
	if err == nil {
		t.Errorf("expected error: bad object")
	}
	assertNullTime(t, bad, "bad from object json")

	var wrongType Time
	err = json.Unmarshal(int64JSON, &wrongType)
	if err == nil {
		t.Errorf("expected error: wrong type JSON")
	}
	assertNullTime(t, wrongType, "wrong type object json")

	var wrongString Time
	err = json.Unmarshal(stringJSON, &wrongString)
	if err == nil {
		t.Errorf("expected error: wrong string JSON")
	}
	assertNullTime(t, wrongString, "wrong string object json")
}

func TestUnmarshalTimeJSON(t *testing.T) {
	testUnmarshalTimeJSON(t, time.RFC3339Nano, timeObject, timeValue, zeroTimeJSON)
}

func CustomTestUnmarshalTimeJSON(t *testing.T) {
	testUnmarshalTimeJSON(t, customTimeFormat, customTimeObject, customTimeValue, customZeroTimeJSON)
}

func testMarshalTime(t *testing.T, f string, value time.Time, tJSON []byte, zTJSON []byte) {
	format.SetTimeFormat(f)

	ti := TimeFrom(value)
	data, err := json.Marshal(ti)
	maybePanic(err)
	assertJSONEquals(t, data, string(tJSON), "non-empty json marshal")

	null := TimeFromPtr(nil)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, string(zTJSON), "empty json marshal")
}

func TestMarshalTime(t *testing.T) {
	testMarshalTime(t, time.RFC3339Nano, timeValue, timeJSON, zeroTimeJSON)
}

func CustomTestMarshalTime(t *testing.T) {
	testMarshalTime(t, customTimeFormat, customTimeValue, customTimeJSON, customZeroTimeJSON)
}

func testUnmarshalTimeText(t *testing.T, f string, value time.Time, tS string, zeroTStr string) {
	format.SetTimeFormat(f)
	ti := TimeFrom(value)
	txt, err := ti.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, tS, "marshal text")

	var unmarshal Time
	err = unmarshal.UnmarshalText(txt)
	maybePanic(err)
	assertTime(t, unmarshal, "unmarshal text", value)

	var null Time
	err = null.UnmarshalText(nullJSON)
	maybePanic(err)
	assertNullTime(t, null, "unmarshal null text")
	txt, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, zeroTStr, "marshal null text")

	var invalid Time
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		t.Error("expected error")
	}
	assertNullTime(t, invalid, "bad string")
}

func TestUnmarshalTimeText(t *testing.T) {
	testUnmarshalTimeText(t, time.RFC3339Nano, timeValue, timeString, zeroTimeStr)
}

func CustomTestUnmarshalTimeText(t *testing.T) {
	testUnmarshalTimeText(t, customTimeFormat, customTimeValue, customTimeString, customZeroTimeStr)
}

func testTimeFrom(t *testing.T, f string, value time.Time) {
	format.SetTimeFormat(f)
	ti := TimeFrom(value)
	assertTime(t, ti, "TimeFrom() time.Time", value)

	var nt time.Time
	null := TimeFrom(nt)
	assertNullTime(t, null, "TimeFrom() empty time.Time")
}

func TestTimeFrom(t *testing.T) {
	testTimeFrom(t, time.RFC3339Nano, timeValue)
}

func CustomTestTimeFrom(t *testing.T) {
	testTimeFrom(t, customTimeFormat, customTimeValue)
}

func testTimeFromPtr(t *testing.T, f string, value time.Time) {
	format.SetTimeFormat(f)
	ti := TimeFromPtr(&value)
	assertTime(t, ti, "TimeFromPtr() time", value)

	null := TimeFromPtr(nil)
	assertNullTime(t, null, "TimeFromPtr(nil)")
}

func TestTimeFromPtr(t *testing.T) {
	testTimeFromPtr(t, time.RFC3339Nano, timeValue)
}

func CustomTestTimeFromPtr(t *testing.T) {
	testTimeFromPtr(t, customTimeFormat, customTimeValue)
}

func testTimeSetValid(t *testing.T, f string, value time.Time) {
	format.SetTimeFormat(f)
	var ti time.Time
	change := TimeFrom(ti)
	assertNullTime(t, change, "SetValid()")
	change.SetValid(value)
	assertTime(t, change, "SetValid()", value)
}

func TestTimeSetValid(t *testing.T) {
	testTimeSetValid(t, time.RFC3339Nano, timeValue)
}

func CustomTestTimeSetValid(t *testing.T) {
	testTimeSetValid(t, customTimeFormat, customTimeValue)
}

func testTimePointer(t *testing.T, f string, value time.Time) {
	format.SetTimeFormat(f)
	ti := TimeFrom(value)
	ptr := ti.Ptr()
	if *ptr != value {
		t.Errorf("bad %s time: %#v ≠ %v\n", "pointer", ptr, value)
	}

	var nt time.Time
	null := TimeFrom(nt)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s time: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestTimePointer(t *testing.T) {
	testTimePointer(t, time.RFC3339Nano, timeValue)
}

func CustomTestTimePointer(t *testing.T) {
	testTimePointer(t, customTimeFormat, customTimeValue)
}

func testTimeScan(t *testing.T, f string, value time.Time) {
	format.SetTimeFormat(f)
	var ti Time
	err := ti.Scan(value)
	maybePanic(err)
	assertTime(t, ti, "scanned time", value)

	var null Time
	err = null.Scan(nil)
	maybePanic(err)
	assertNullTime(t, null, "scanned null")

	var wrong Time
	err = wrong.Scan(int64(42))
	if err == nil {
		t.Error("expected error")
	}
	assertNullTime(t, wrong, "scanned wrong")
}

func TestTimeScan(t *testing.T) {
	testTimeScan(t, time.RFC3339Nano, timeValue)
}

func CustomTestTimeScan(t *testing.T) {
	testTimeScan(t, customTimeFormat, customTimeValue)
}

func testTimeValue(t *testing.T, f string, value time.Time) {
	format.SetTimeFormat(f)
	ti := TimeFrom(value)
	v, err := ti.Value()
	maybePanic(err)
	if ti.Time != value {
		t.Errorf("bad time.Time value: %v ≠ %v", ti.Time, value)
	}

	var nt time.Time
	zero := TimeFrom(nt)
	v, err = zero.Value()
	maybePanic(err)
	if v != nil {
		t.Errorf("bad %s time.Time value: %v ≠ %v", "zero", v, nil)
	}
}

func TestTimeValue(t *testing.T) {
	testTimeValue(t, time.RFC3339Nano, timeValue)
}

func CustomTestTimeValue(t *testing.T) {
	testTimeValue(t, customTimeFormat, customTimeValue)
}

func assertTime(t *testing.T, ti Time, from string, val time.Time) {
	if ti.IsZero() != val.IsZero() {
		t.Errorf("%v IsZero() != %v\n", ti, val)
	}
	if !ti.Time.Equal(val) {
		t.Errorf("bad %v time: %v ≠ %v\n", from, ti.Time, val)
	}
	if !ti.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullTime(t *testing.T, ti Time, from string) {
	if ti.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
	if !ti.IsZero() {
		t.Error(from, "is not zero, but should be zero")
	}
}
