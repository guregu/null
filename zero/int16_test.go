package zero

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"
)

var (
	int16JSON     = []byte(`32766`)
	nullInt16JSON = []byte(`{"Int16":32766,"Valid":true}`)
	zero16JSON    = []byte(`0`)
)

func TestInt16From(t *testing.T) {
	i := Int16From(32766)
	assertInt16(t, i, "Int16From()")

	zero := Int16From(0)
	if zero.Valid {
		t.Error("Int16From(0)", "is valid, but should be invalid")
	}
}

func TestInt16FromPtr(t *testing.T) {
	n := int16(32766)
	iptr := &n
	i := Int16FromPtr(iptr)
	assertInt16(t, i, "Int16FromPtr()")

	null := Int16FromPtr(nil)
	assertNullInt16(t, null, "Int16FromPtr(nil)")
}

func TestUnmarshalInt16(t *testing.T) {
	var i Int16
	err := json.Unmarshal(int16JSON, &i)
	maybePanic(err)
	assertInt16(t, i, "int json")

	var ni Int16
	err = json.Unmarshal(nullInt16JSON, &ni)
	maybePanic(err)
	assertInt16(t, ni, "sql.NullInt16 json")

	var zero Int16
	err = json.Unmarshal(zero16JSON, &zero)
	maybePanic(err)
	assertNullInt16(t, zero, "zero json")

	var null Int16
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullInt16(t, null, "null json")

	var badType Int16
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullInt16(t, badType, "wrong type json")

	var invalid Int16
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullInt16(t, invalid, "invalid json")
}

func TestUnmarshalNonIntegerNumber16(t *testing.T) {
	var i Int16
	err := json.Unmarshal(floatJSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int")
	}
}

func TestUnmarshalInt16Overflow(t *testing.T) {
	int64Overflow := uint64(math.MaxInt16)

	// Max int64 should decode successfully
	var i Int16
	err := json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	int64Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int64")
	}
}

func TestTextUnmarshalInt16(t *testing.T) {
	var i Int16
	err := i.UnmarshalText([]byte("32766"))
	maybePanic(err)
	assertInt16(t, i, "UnmarshalText() int")

	var zero Int16
	err = zero.UnmarshalText([]byte("0"))
	maybePanic(err)
	assertNullInt16(t, zero, "UnmarshalText() zero int")

	var blank Int16
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullInt16(t, blank, "UnmarshalText() empty int")

	var null Int16
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullInt16(t, null, `UnmarshalText() "null"`)
}

func TestMarshalInt16(t *testing.T) {
	i := Int16From(32766)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "32766", "non-empty json marshal")

	// invalid values should be encoded as 0
	null := NewInt16(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "0", "null json marshal")
}

func TestMarshalInt16Text(t *testing.T) {
	i := Int16From(32766)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "32766", "non-empty text marshal")

	// invalid values should be encoded as zero
	null := NewInt16(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "0", "null text marshal")
}

func TestInt16Pointer(t *testing.T) {
	i := Int16From(32766)
	ptr := i.Ptr()
	if *ptr != 32766 {
		t.Errorf("bad %s int: %#v ≠ %d\n", "pointer", ptr, 32766)
	}

	null := NewInt16(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestInt16IsZero(t *testing.T) {
	i := Int16From(32766)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewInt16(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewInt16(0, true)
	if !zero.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestInt16Scan(t *testing.T) {
	var i Int16
	err := i.Scan(32766)
	maybePanic(err)
	assertInt16(t, i, "scanned int")

	var null Int16
	err = null.Scan(nil)
	maybePanic(err)
	assertNullInt16(t, null, "scanned null")
}

func TestInt16SetValid(t *testing.T) {
	change := NewInt16(0, false)
	assertNullInt16(t, change, "SetValid()")
	change.SetValid(32766)
	assertInt16(t, change, "SetValid()")
}

func assertInt16(t *testing.T, i Int16, from string) {
	if i.Int16 != 32766 {
		t.Errorf("bad %s int: %d ≠ %d\n", from, i.Int16, 32766)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullInt16(t *testing.T, i Int16, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
