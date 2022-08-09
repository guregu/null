package null

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"testing"
)

var (
	int16JSON       = []byte(`12345`)
	int16StringJSON = []byte(`"12345"`)
	nullInt16JSON   = []byte(`{"Int16":12345,"Valid":true}`)
)

func TestInt16From(t *testing.T) {
	i := Int16From(12345)
	assertInt16(t, i, "Int16From()")

	zero := Int16From(0)
	if !zero.Valid {
		t.Error("Int16From(0)", "is invalid, but should be valid")
	}
}

func TestInt16FromPtr(t *testing.T) {
	n := int16(12345)
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

	var si Int16
	err = json.Unmarshal(int16StringJSON, &si)
	maybePanic(err)
	assertInt16(t, si, "int string json")

	var ni Int16
	err = json.Unmarshal(nullInt16JSON, &ni)
	if err == nil {
		panic("err should not be nill")
	}

	var bi Int16
	err = json.Unmarshal(floatBlankJSON, &bi)
	if err == nil {
		panic("err should not be nill")
	}

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
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
	assertNullInt16(t, invalid, "invalid json")
}

func TestUnmarshalNonInt16egerNumber(t *testing.T) {
	var i Int16
	err := json.Unmarshal(floatJSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int")
	}
}

func TestUnmarshalInt16Overflow(t *testing.T) {
	int16Overflow := uint64(math.MaxInt16)

	// Max int16 should decode successfully
	var i Int16
	err := json.Unmarshal([]byte(strconv.FormatUint(int16Overflow, 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	int16Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(int16Overflow, 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int16")
	}
}

func TestTextUnmarshalInt16(t *testing.T) {
	var i Int16
	err := i.UnmarshalText([]byte("12345"))
	maybePanic(err)
	assertInt16(t, i, "UnmarshalText() int")

	var blank Int16
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullInt16(t, blank, "UnmarshalText() empty int")

	var null Int16
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullInt16(t, null, `UnmarshalText() "null"`)

	var invalid Int16
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		panic("expected error")
	}
}

func TestMarshalInt16(t *testing.T) {
	i := Int16From(12345)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "12345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewInt16(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalInt16Text(t *testing.T) {
	i := Int16From(12345)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "12345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewInt16(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestInt16Pointer(t *testing.T) {
	i := Int16From(12345)
	ptr := i.Ptr()
	if *ptr != 12345 {
		t.Errorf("bad %s int: %#v ≠ %d\n", "pointer", ptr, 12345)
	}

	null := NewInt16(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestInt16IsZero(t *testing.T) {
	i := Int16From(12345)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewInt16(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewInt16(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestInt16SetValid(t *testing.T) {
	change := NewInt16(0, false)
	assertNullInt16(t, change, "SetValid()")
	change.SetValid(12345)
	assertInt16(t, change, "SetValid()")
}

func TestInt16Scan(t *testing.T) {
	var i Int16
	err := i.Scan(12345)
	maybePanic(err)
	assertInt16(t, i, "scanned int")

	var null Int16
	err = null.Scan(nil)
	maybePanic(err)
	assertNullInt16(t, null, "scanned null")
}

func TestInt16ValueOrZero(t *testing.T) {
	valid := NewInt16(12345, true)
	if valid.ValueOrZero() != 12345 {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := NewInt16(12345, false)
	if invalid.ValueOrZero() != 0 {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestInt16Equal(t *testing.T) {
	int1 := NewInt16(10, false)
	int2 := NewInt16(10, false)
	assertInt16EqualIsTrue(t, int1, int2)

	int1 = NewInt16(10, false)
	int2 = NewInt16(20, false)
	assertInt16EqualIsTrue(t, int1, int2)

	int1 = NewInt16(10, true)
	int2 = NewInt16(10, true)
	assertInt16EqualIsTrue(t, int1, int2)

	int1 = NewInt16(10, true)
	int2 = NewInt16(10, false)
	assertInt16EqualIsFalse(t, int1, int2)

	int1 = NewInt16(10, false)
	int2 = NewInt16(10, true)
	assertInt16EqualIsFalse(t, int1, int2)

	int1 = NewInt16(10, true)
	int2 = NewInt16(20, true)
	assertInt16EqualIsFalse(t, int1, int2)
}

func assertInt16(t *testing.T, i Int16, from string) {
	if i.Int16 != 12345 {
		t.Errorf("bad %s int: %d ≠ %d\n", from, i.Int16, 12345)
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

func assertInt16EqualIsTrue(t *testing.T, a, b Int16) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf("Equal() of Int16{%v, Valid:%t} and Int16{%v, Valid:%t} should return true", a.Int16, a.Valid, b.Int16, b.Valid)
	}
}

func assertInt16EqualIsFalse(t *testing.T, a, b Int16) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf("Equal() of Int16{%v, Valid:%t} and Int16{%v, Valid:%t} should return false", a.Int16, a.Valid, b.Int16, b.Valid)
	}
}
