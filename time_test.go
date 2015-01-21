package null

import (
	"encoding/json"
	"testing"
	"time"
)

var (
	timeTest, _  = time.Parse(time.RFC3339, "2012-12-21T21:21:21Z")
	nullTimeJSON = []byte(`{"Time":"2012-12-21T21:21:21Z","Valid":true}`)
	nilJSON      = []byte(`null`)
)

func TestTimeFrom(t *testing.T) {
	i := TimeFrom(timeTest)
	assertTime(t, i, "TimeFrom()")
}

func TestTimeFromPtr(t *testing.T) {
	iptr := &timeTest
	i := TimeFromPtr(iptr)
	assertTime(t, i, "TimeFromPtr()")

	null := TimeFromPtr(nil)
	assertNullTime(t, null, "TimeFromPtr(nil)")
}

func TestUnmarshalTime(t *testing.T) {
	var ni Time
	err := json.Unmarshal(nullTimeJSON, &ni)
	maybePanic(err)
	assertTime(t, ni, "NullTime json")

	var nl Time
	err = json.Unmarshal(nilJSON, &nl)
	maybePanic(err)
	assertNullTime(t, nl, "null json")

	var badType Time
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullTime(t, badType, "wrong type json")

}

func TestTextUnmarshalTime(t *testing.T) {
	var i Time
	err := i.UnmarshalText([]byte("2012-12-21T21:21:21Z"))
	maybePanic(err)
	assertTime(t, i, "UnmarshalText() time")

	var blank Time
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullTime(t, blank, "UnmarshalText() empty time")

	var null Time
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullTime(t, null, `UnmarshalText() "null"`)

}

func TestMarshalTime(t *testing.T) {
	i := TimeFrom(timeTest)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "2012-12-21T21:21:21Z", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewTime(time.Time{}, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null json marshal")
}

func TestMarshalTimeText(t *testing.T) {
	i := TimeFrom(timeTest)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "2012-12-21T21:21:21Z", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewTime(time.Time{}, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestTimePointer(t *testing.T) {
	i := TimeFrom(timeTest)
	ptr := i.Ptr()
	if *ptr != timeTest {
		t.Errorf("bad %s time.Time: %#v ≠ %s\n", "pointer", ptr, timeTest)
	}

	null := NewTime(time.Time{}, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s time.Time: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestTimeIsZero(t *testing.T) {
	i := TimeFrom(timeTest)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewTime(time.Time{}, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestTimeSetValid(t *testing.T) {
	change := NewTime(time.Time{}, false)
	assertNullTime(t, change, "SetValid()")
	change.SetValid(timeTest)
	assertTime(t, change, "SetValid()")
}

func TestTimeScan(t *testing.T) {
	var i Time
	err := i.Scan(timeTest)
	maybePanic(err)
	assertTime(t, i, "scanned time.Time")

	var null Time
	err = null.Scan(nil)
	maybePanic(err)
	assertNullTime(t, null, "scanned null")
}

func assertTime(t *testing.T, i Time, from string) {
	if i.Time != timeTest {
		t.Errorf("bad %s time.Time: %d ≠ %d\n", from, i.Time, timeTest)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullTime(t *testing.T, i Time, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
