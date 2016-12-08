package zero

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"
)

var (
	int64JSON     = []byte(`12345`)
	nullInt64JSON = []byte(`{"Int64":12345,"Valid":true}`)
	zeroJSON      = []byte(`0`)
)

func TestInt64From(t *testing.T) {
	i := Int64From(12345)
	assertInt64(t, i, "Int64From())")

	zero := Int64From(0)
	if zero.Valid {
		t.Error("Int64From(0)", "is valid, but should be invalid")
	}
}

func TestInt64FromPtr(t *testing.T) {
	n := int64(12345)
	iptr := &n
	i := Int64FromPtr(iptr)
	assertInt64(t, i, "Int64FromPtr())")

	null := Int64FromPtr(nil)
	assertNullInt64(t, null, "Int64FromPtr(nil)")
}

func TestUnmarshalInt64(t *testing.T) {
	var i Int64
	err := json.Unmarshal(int64JSON, &i)
	maybePanic(err)
	assertInt64(t, i, "int json")

	var ni Int64
	err = json.Unmarshal(nullInt64JSON, &ni)
	maybePanic(err)
	assertInt64(t, ni, "sql.NullInt64 json")

	var zero Int64
	err = json.Unmarshal(zeroJSON, &zero)
	maybePanic(err)
	assertNullInt64(t, zero, "zero json")

	var null Int64
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullInt64(t, null, "null json")

	var badType Int64
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullInt64(t, badType, "wrong type json")

	var invalid Int64
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullInt64(t, invalid, "invalid json")
}

func TestUnmarshalNonIntegerNumber64(t *testing.T) {
	var i Int64
	err := json.Unmarshal(floatJSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int")
	}
}

func TestUnmarshalInt64Overflow(t *testing.T) {
	int64Overflow := uint64(math.MaxInt64)

	// Max int64 should decode successfully
	var i Int64
	err := json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	int64Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int64")
	}
}

func TestTextUnmarshalInt64(t *testing.T) {
	var i Int64
	err := i.UnmarshalText([]byte("12345"))
	maybePanic(err)
	assertInt64(t, i, "UnmarshalText() int")

	var zero Int64
	err = zero.UnmarshalText([]byte("0"))
	maybePanic(err)
	assertNullInt64(t, zero, "UnmarshalText() zero int")

	var blank Int64
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullInt64(t, blank, "UnmarshalText() empty int")

	var null Int64
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullInt64(t, null, `UnmarshalText() "null"`)
}

func TestMarshalInt64(t *testing.T) {
	i := Int64From(12345)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "12345", "non-empty json marshal")

	// invalid values should be encoded as 0
	null := NewInt64(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "0", "null json marshal")
}

func TestMarshalInt64Text(t *testing.T) {
	i := Int64From(12345)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "12345", "non-empty text marshal")

	// invalid values should be encoded as zero
	null := NewInt64(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "0", "null text marshal")
}

func TestInt64Pointer(t *testing.T) {
	i := Int64From(12345)
	ptr := i.Ptr()
	if *ptr != 12345 {
		t.Errorf("bad %s int: %#v ≠ %d\n", "pointer", ptr, 12345)
	}

	null := NewInt64(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestInt64IsZero(t *testing.T) {
	i := Int64From(12345)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewInt64(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewInt64(0, true)
	if !zero.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestInt64Scan(t *testing.T) {
	var i Int64
	err := i.Scan(12345)
	maybePanic(err)
	assertInt64(t, i, "scanned int")

	var null Int64
	err = null.Scan(nil)
	maybePanic(err)
	assertNullInt64(t, null, "scanned null")
}

func TestInt64SetValid(t *testing.T) {
	change := NewInt64(0, false)
	assertNullInt64(t, change, "SetValid()")
	change.SetValid(12345)
	assertInt64(t, change, "SetValid()")
}

func assertInt64(t *testing.T, i Int64, from string) {
	if i.Int64 != 12345 {
		t.Errorf("bad %s int: %d ≠ %d\n", from, i.Int64, 12345)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullInt64(t *testing.T, i Int64, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
