package null

import (
	"encoding/json"
	//"fmt"
	"testing"
	"time"

	"github.com/axiomzen/null/format"
)

var (
	timeString    = "2012-12-21T21:21:21Z"
	timeJSON      = []byte(`"` + timeString + `"`)
	nullTimeJSON  = []byte(`null`)
	timeValue, _  = time.Parse(time.RFC3339, timeString)
	timeObject    = []byte(`{"Time":"2012-12-21T21:21:21Z","Valid":true}`)
	nullObject    = []byte(`{"Time":"0001-01-01T00:00:00Z","Valid":false}`)
	timeObjectXML = []byte(`<Time>2012-12-21T21:21:21Z</Time>`)
	badObject     = []byte(`{"hello": "world"}`)

	//	RFC3339     = "2006-01-02T15:04:05Z07:00"
	//	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	customTimeFormat = "2006-01-02T15:04:05Z0700"
	//customTimeFormat = "2006-01-02T15:04:05-07:00"
	//customTimeString = time.Now().Format(customTimeFormat)
	//customTimeString = "2012-12-21T21:21:21-02:00"
	customTimeString   = "2012-12-21T21:21:21-0200"
	customTimeJSON     = []byte(`"` + customTimeString + `"`)
	customTimeValue, _ = time.Parse(customTimeFormat, customTimeString)
	customTimeObject   = []byte(`{"Time":"` + customTimeString + `","Valid":true}`)
	customNullObject   = []byte(`{"Time":"0001-01-01T00:00:00+0000","Valid":false}`)
	//customNullObject   = []byte(`{"Time":"0001-01-01T00:00:00+00:00","Valid":false}`)
)

func testUnmarshalTimeJSON(t *testing.T, f string, jsonStr []byte, value time.Time, to []byte) {
	format.SetTimeFormat(f)

	var ti Time
	err := json.Unmarshal(jsonStr, &ti)
	maybePanic(err)
	assertTime(t, ti, "UnmarshalJSON() json", value)

	var null Time
	err = json.Unmarshal(nullTimeJSON, &null)
	maybePanic(err)
	assertNullTime(t, null, "null time json")

	var fromObject Time
	err = json.Unmarshal(to, &fromObject)
	maybePanic(err)
	assertTime(t, fromObject, "time from object json", value)

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
}

func TestUnmarshalTimeJSON(t *testing.T) {
	testUnmarshalTimeJSON(t, time.RFC3339Nano, timeJSON, timeValue, timeObject)
}

func TestUnmarshalCustomTimeJSON(t *testing.T) {
	testUnmarshalTimeJSON(t, customTimeFormat, customTimeJSON, customTimeValue, customTimeObject)
}

//t *testing.T, f string, jsonStr []byte, value time.Time, to []byte
func testUnmarshalTimeText(t *testing.T, f string, timeStr string, value time.Time) {
	format.SetTimeFormat(f)

	ti := TimeFrom(value)
	txt, err := ti.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, txt, timeStr, "marshal text")

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
	assertJSONEquals(t, txt, string(nullJSON), "marshal null text")

	var invalid Time
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		t.Error("expected error")
	}
	assertNullTime(t, invalid, "bad string")
}

func TestUnmarshalTimeText(t *testing.T) {
	testUnmarshalTimeText(t, time.RFC3339Nano, timeString, timeValue)
}

func TestUnmarshalCustomTimeText(t *testing.T) {
	testUnmarshalTimeText(t, customTimeFormat, customTimeString, customTimeValue)
}

func testMarshalTime(t *testing.T, f string, value time.Time, tJSON []byte) {
	format.SetTimeFormat(f)

	ti := TimeFrom(value)
	data, err := json.Marshal(ti)
	maybePanic(err)
	assertJSONEquals(t, data, string(tJSON), "non-empty json marshal")

	ti.Valid = false
	data, err = json.Marshal(ti)
	maybePanic(err)
	assertJSONEquals(t, data, string(nullJSON), "null json marshal")
}

func TestMarshalTime(t *testing.T) {
	testMarshalTime(t, time.RFC3339Nano, timeValue, timeJSON)
}

func TestMarshalCustomTime(t *testing.T) {
	testMarshalTime(t, customTimeFormat, customTimeValue, customTimeJSON)
}

func testTimeFrom(t *testing.T, f string, value time.Time) {
	format.SetTimeFormat(f)

	ti := TimeFrom(value)
	assertTime(t, ti, "TimeFrom() time.Time", value)
}

func TestTimeFrom(t *testing.T) {
	testTimeFrom(t, time.RFC3339Nano, timeValue)
}

func TestCustomTimeFrom(t *testing.T) {
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

func TestCustomTimeFromPtr(t *testing.T) {
	testTimeFromPtr(t, customTimeFormat, customTimeValue)
}

func testTimeSetValid(t *testing.T, f string, value time.Time) {
	format.SetTimeFormat(f)

	var ti time.Time
	change := NewTime(ti, false)
	assertNullTime(t, change, "SetValid()")
	change.SetValid(value)
	assertTime(t, change, "SetValid()", value)
}

func TestTimeSetValid(t *testing.T) {
	testTimeSetValid(t, time.RFC3339Nano, timeValue)
}

func TestCustomTimeSetValid(t *testing.T) {
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
	null := NewTime(nt, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s time: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestTimePointer(t *testing.T) {
	testTimePointer(t, time.RFC3339Nano, timeValue)
}

func TestCustomTimePointer(t *testing.T) {
	testTimePointer(t, customTimeFormat, customTimeValue)
}

func testTimeScanValue(t *testing.T, f string, value time.Time) {
	format.SetTimeFormat(f)

	var ti Time
	err := ti.Scan(value)
	maybePanic(err)
	assertTime(t, ti, "scanned time", value)
	if v, err := ti.Value(); v != value || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var null Time
	err = null.Scan(nil)
	maybePanic(err)
	assertNullTime(t, null, "scanned null")
	if v, err := null.Value(); v != nil || err != nil {
		t.Error("bad value or err:", v, err)
	}

	var wrong Time
	err = wrong.Scan(int64(42))
	if err == nil {
		t.Error("expected error")
	}
	assertNullTime(t, wrong, "scanned wrong")
}

func TestTimeScanValue(t *testing.T) {
	testTimeScanValue(t, time.RFC3339Nano, timeValue)
}

func TestCustomTimeScanValue(t *testing.T) {
	testTimeScanValue(t, customTimeFormat, customTimeValue)
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
