package null

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"
)

var (
	int8JSON = []byte(`126`)
)

func TestInt8From(t *testing.T) {
	i := Int8From(126)
	assertInt8(t, i, "Int8From()")

	zero := Int8From(0)
	if !zero.Valid {
		t.Error("Int8From(0)", "is invalid, but should be valid")
	}
}

func TestInt8FromPtr(t *testing.T) {
	n := int8(126)
	iptr := &n
	i := Int8FromPtr(iptr)
	assertInt8(t, i, "Int8FromPtr()")

	null := Int8FromPtr(nil)
	assertNullInt8(t, null, "Int8FromPtr(nil)")
}

func TestUnmarshalInt8(t *testing.T) {
	var i Int8
	err := json.Unmarshal(int8JSON, &i)
	maybePanic(err)
	assertInt8(t, i, "int8 json")

	var null Int8
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullInt8(t, null, "null json")

	var badType Int8
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullInt8(t, badType, "wrong type json")

	var invalid Int8
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullInt8(t, invalid, "invalid json")
}

func TestUnmarshalNonIntegerNumber8(t *testing.T) {
	var i Int8
	err := json.Unmarshal(float64JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int8")
	}
}

func TestUnmarshalInt8Overflow(t *testing.T) {
	int8Overflow := uint8(math.MaxInt8)

	// Max int8 should decode successfully
	var i Int8
	err := json.Unmarshal([]byte(strconv.FormatUint(uint64(int8Overflow), 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	int8Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(uint64(int8Overflow), 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int8")
	}
}

func TestTextUnmarshalInt8(t *testing.T) {
	var i Int8
	err := i.UnmarshalText([]byte("126"))
	maybePanic(err)
	assertInt8(t, i, "UnmarshalText() int8")

	var blank Int8
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullInt8(t, blank, "UnmarshalText() empty int8")
}

func TestMarshalInt8(t *testing.T) {
	i := Int8From(126)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "126", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewInt8(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalInt8Text(t *testing.T) {
	i := Int8From(126)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "126", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewInt8(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestInt8Pointer(t *testing.T) {
	i := Int8From(126)
	ptr := i.Ptr()
	if *ptr != 126 {
		t.Errorf("bad %s int8: %#v ≠ %d\n", "pointer", ptr, 126)
	}

	null := NewInt8(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int8: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestInt8IsZero(t *testing.T) {
	i := Int8From(126)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewInt8(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewInt8(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestInt8SetValid(t *testing.T) {
	change := NewInt8(0, false)
	assertNullInt8(t, change, "SetValid()")
	change.SetValid(126)
	assertInt8(t, change, "SetValid()")
}

func TestInt8Scan(t *testing.T) {
	var i Int8
	err := i.Scan(126)
	maybePanic(err)
	assertInt8(t, i, "scanned int8")

	var null Int8
	err = null.Scan(nil)
	maybePanic(err)
	assertNullInt8(t, null, "scanned null")
}

func assertInt8(t *testing.T, i Int8, from string) {
	if i.Int8 != 126 {
		t.Errorf("bad %s int8: %d ≠ %d\n", from, i.Int8, 126)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullInt8(t *testing.T, i Int8, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
