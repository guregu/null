package null

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"
	"time"
)

var (
	testTimestamp = time.Unix(1419196236, 0)
	//nullJSON     = []byte(`null`)
	timestampJSON     = []byte(`1419196236`)
	nullTimestampJSON = []byte(`{"Time":1419196236,"Valid":true}`)
)

func testTimestampstampFrom(t *testing.T) {
	i := TimestampFrom(testTimestamp)
	assertTimestamp(t, i, "TimestampFrom()")

	zero := TimestampFrom(time.Time{})
	if !zero.Valid {
		t.Error("TimestampFrom(0)", "is invalid, but should be valid")
	}
}

func testTimestampstampFromPtr(t *testing.T) {
	n := testTimestamp
	iptr := &n
	i := TimestampFromPtr(iptr)
	assertTimestamp(t, i, "TimestampFromPtr()")

	null := TimestampFromPtr(nil)
	assertNullTimestamp(t, null, "TimestampFromPtr(nil)")
}

func TestUnmarshalTimestamp(t *testing.T) {
	var i Timestamp
	err := json.Unmarshal(timestampJSON, &i)
	maybePanic(err)
	assertTimestamp(t, i, "int json")

	var ni Timestamp
	err = json.Unmarshal(nullTimestampJSON, &ni)
	maybePanic(err)
	assertTimestamp(t, ni, "pq.NullTime json")

	var null Timestamp
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullTimestamp(t, null, "null json")

	var badType Timestamp
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullTimestamp(t, badType, "wrong type json")
}

func TestUnmarshalNonTimestampNumber(t *testing.T) {
	var i Timestamp
	err := json.Unmarshal(floatJSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to time.Time")
	}
}

func TestUnmarshalTimestampOverflow(t *testing.T) {
	int64Overflow := uint64(math.MaxInt64)

	// Max int64 should decode successfully
	var i Timestamp
	err := json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	int64Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int64")
	}
}

func TestTextUnmarshalTimestamp(t *testing.T) {
	var i Timestamp
	err := i.UnmarshalText(timestampJSON)
	maybePanic(err)
	assertTimestamp(t, i, "UnmarshalText() time.Time")

	var blank Timestamp
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullTimestamp(t, blank, "UnmarshalText() empty time.Time")

	var null Timestamp
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullTimestamp(t, null, `UnmarshalText() "null"`)
}

func TestMarshalTimestamp(t *testing.T) {
	i := TimestampFrom(testTimestamp)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "1419196236", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewTimestamp(time.Time{}, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalTimestampText(t *testing.T) {
	i := TimestampFrom(testTimestamp)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "1419196236", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewTimestamp(time.Time{}, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func testTimestampstampPointer(t *testing.T) {
	i := TimestampFrom(testTimestamp)
	ptr := i.Ptr()
	if *ptr != testTimestamp {
		t.Errorf("bad %s time.Time: %#v ≠ %s\n", "pointer", ptr, 12345)
	}

	null := NewTimestamp(time.Time{}, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s time.Time: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func testTimestampstampIsZero(t *testing.T) {
	i := TimestampFrom(testTimestamp)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewTimestamp(time.Time{}, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewTimestamp(time.Time{}, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func testTimestampstampSetValid(t *testing.T) {
	change := NewTimestamp(time.Time{}, false)
	assertNullTimestamp(t, change, "SetValid()")
	change.SetValid(testTimestamp)
	assertTimestamp(t, change, "SetValid()")
}

func testTimestampstampScan(t *testing.T) {
	var i Timestamp
	err := i.Scan(testTimestamp)
	maybePanic(err)
	assertTimestamp(t, i, "scanned time.Time")

	var null Timestamp
	err = null.Scan(nil)
	maybePanic(err)
	assertNullTimestamp(t, null, "scanned null")
}

func assertTimestamp(t *testing.T, i Timestamp, from string) {
	if i.Time != testTimestamp {
		t.Errorf("bad %s time.Time: %d ≠ %d\n", from, i.Time, testTimestamp)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullTimestamp(t *testing.T, i Timestamp, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
