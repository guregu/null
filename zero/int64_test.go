package zero

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"testing"
)

var (
	intJSON       = []byte(`12345`)
	intStringJSON = []byte(`"12345"`)
	nullIntJSON   = []byte(`{"Int64":12345,"Valid":true}`)
	zeroJSON      = []byte(`0`)
)

func TestIntFrom(t *testing.T) {
	i := Int64From(12345)
	assertInt(t, i, "Int64From()")

	zero := Int64From(0)
	if zero.Valid {
		t.Error("Int64From(0)", "is valid, but should be invalid")
	}
}

func TestIntFromPtr(t *testing.T) {
	n := int64(12345)
	iptr := &n
	i := Int64FromPtr(iptr)
	assertInt(t, i, "Int64FromPtr()")

	null := Int64FromPtr(nil)
	assertNullInt(t, null, "Int64FromPtr(nil)")
}

func TestUnmarshalInt(t *testing.T) {
	var i Int64
	err := json.Unmarshal(intJSON, &i)
	maybePanic(err)
	assertInt(t, i, "int json")

	var si Int64
	err = json.Unmarshal(intStringJSON, &si)
	maybePanic(err)
	assertInt(t, si, "int string json")

	var ni Int64
	err = json.Unmarshal(nullIntJSON, &ni)
	if err == nil {
		panic("expected error")
	}

	var bi Int64
	err = json.Unmarshal(float64BlankJSON, &bi)
	if err == nil {
		panic("expected error")
	}

	var zero Int64
	err = json.Unmarshal(zeroJSON, &zero)
	maybePanic(err)
	assertNullInt(t, zero, "zero json")

	var null Int64
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullInt(t, null, "null json")

	var badType Int64
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullInt(t, badType, "wrong type json")

	var invalid Int64
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
	assertNullInt(t, invalid, "invalid json")
}

func TestUnmarshalNonIntegerNumber(t *testing.T) {
	var i Int64
	err := json.Unmarshal(float64JSON, &i)
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

func TestTextUnmarshalInt(t *testing.T) {
	var i Int64
	err := i.UnmarshalText([]byte("12345"))
	maybePanic(err)
	assertInt(t, i, "UnmarshalText() int")

	var zero Int64
	err = zero.UnmarshalText([]byte("0"))
	maybePanic(err)
	assertNullInt(t, zero, "UnmarshalText() zero int")

	var blank Int64
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullInt(t, blank, "UnmarshalText() empty int")

	var null Int64
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullInt(t, null, `UnmarshalText() "null"`)

	var invalid Int64
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		panic("expected error")
	}
}

func TestMarshalInt(t *testing.T) {
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

func TestMarshalIntText(t *testing.T) {
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

func TestIntPointer(t *testing.T) {
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

func TestIntIsZero(t *testing.T) {
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

func TestIntScan(t *testing.T) {
	var i Int64
	err := i.Scan(12345)
	maybePanic(err)
	assertInt(t, i, "scanned int")

	var null Int64
	err = null.Scan(nil)
	maybePanic(err)
	assertNullInt(t, null, "scanned null")
}

func TestIntSetValid(t *testing.T) {
	change := NewInt64(0, false)
	assertNullInt(t, change, "SetValid()")
	change.SetValid(12345)
	assertInt(t, change, "SetValid()")
}

func TestIntValueOrZero(t *testing.T) {
	valid := NewInt64(12345, true)
	if valid.ValueOrZero() != 12345 {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := NewInt64(12345, false)
	if invalid.ValueOrZero() != 0 {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestIntEqual(t *testing.T) {
	int1 := NewInt64(10, false)
	int2 := NewInt64(10, false)
	assertIntEqualIsTrue(t, int1, int2)

	int1 = NewInt64(10, false)
	int2 = NewInt64(20, false)
	assertIntEqualIsTrue(t, int1, int2)

	int1 = NewInt64(10, true)
	int2 = NewInt64(10, true)
	assertIntEqualIsTrue(t, int1, int2)

	int1 = NewInt64(0, true)
	int2 = NewInt64(10, false)
	assertIntEqualIsTrue(t, int1, int2)

	int1 = NewInt64(10, true)
	int2 = NewInt64(10, false)
	assertIntEqualIsFalse(t, int1, int2)

	int1 = NewInt64(10, false)
	int2 = NewInt64(10, true)
	assertIntEqualIsFalse(t, int1, int2)

	int1 = NewInt64(10, true)
	int2 = NewInt64(20, true)
	assertIntEqualIsFalse(t, int1, int2)
}

func assertInt(t *testing.T, i Int64, from string) {
	if i.Int64 != 12345 {
		t.Errorf("bad %s int: %d ≠ %d\n", from, i.Int64, 12345)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullInt(t *testing.T, i Int64, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertIntEqualIsTrue(t *testing.T, a, b Int64) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf("Equal() of Int64{%v, Valid:%t} and Int64{%v, Valid:%t} should return true", a.Int64, a.Valid, b.Int64, b.Valid)
	}
}

func assertIntEqualIsFalse(t *testing.T, a, b Int64) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf("Equal() of Int64{%v, Valid:%t} and Int64{%v, Valid:%t} should return false", a.Int64, a.Valid, b.Int64, b.Valid)
	}
}
