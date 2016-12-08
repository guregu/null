package null

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"
)

var (
	int32JSON     = []byte(`12345`)
	nullInt32JSON = []byte(`{"Int32":12345,"Valid":true}`)
)

func TestInt32From(t *testing.T) {
	i := Int32From(12345)
	assertInt32(t, i, "Int32From()")

	zero := Int32From(0)
	if !zero.Valid {
		t.Error("Int32From(0)", "is invalid, but should be valid")
	}
}

func TestInt32FromPtr(t *testing.T) {
	n := int32(12345)
	iptr := &n
	i := Int32FromPtr(iptr)
	assertInt32(t, i, "Int32FromPtr()")

	null := Int32FromPtr(nil)
	assertNullInt32(t, null, "Int32FromPtr(nil)")
}

func TestUnmarshalInt32(t *testing.T) {
	var i Int32
	err := json.Unmarshal(int32JSON, &i)
	maybePanic(err)
	assertInt32(t, i, "int json")

	var ni Int32
	err = json.Unmarshal(nullInt32JSON, &ni)
	maybePanic(err)
	assertInt32(t, ni, "sq.NullInt32 json")

	var null Int32
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullInt32(t, null, "null json")

	var badType Int32
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullInt32(t, badType, "wrong type json")

	var invalid Int32
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullInt32(t, invalid, "invalid json")
}

func TestUnmarshalNonIntegerNumberInt32(t *testing.T) {
	var i Int32
	err := json.Unmarshal(floatJSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int")
	}
}

func TestUnmarshalInt32Overflow(t *testing.T) {
	int32Overflow := uint32(math.MaxInt32)

	// Max int32 should decode successfully
	var i Int32
	err := json.Unmarshal([]byte(strconv.FormatUint(uint64(int32Overflow), 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	int32Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(uint64(int32Overflow), 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int32")
	}
}

func TestTextUnmarshalInt32(t *testing.T) {
	var i Int32
	err := i.UnmarshalText([]byte("12345"))
	maybePanic(err)
	assertInt32(t, i, "UnmarshalText() int")

	var blank Int32
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullInt32(t, blank, "UnmarshalText() empty int")

	var null Int32
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullInt32(t, null, `UnmarshalText() "null"`)
}

func TestMarshalInt32(t *testing.T) {
	i := Int32From(12345)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "12345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewInt32(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalInt32Text(t *testing.T) {
	i := Int32From(12345)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "12345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewInt32(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestInt32Pointer(t *testing.T) {
	i := Int32From(12345)
	ptr := i.Ptr()
	if *ptr != 12345 {
		t.Errorf("bad %s int: %#v ≠ %d\n", "pointer", ptr, 12345)
	}

	null := NewInt32(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestInt32IsZero(t *testing.T) {
	i := Int32From(12345)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewInt32(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewInt32(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestInt32SetValid(t *testing.T) {
	change := NewInt32(0, false)
	assertNullInt32(t, change, "SetValid()")
	change.SetValid(12345)
	assertInt32(t, change, "SetValid()")
}

func TestInt32Scan(t *testing.T) {
	var i Int32
	err := i.Scan(12345)
	maybePanic(err)
	assertInt32(t, i, "scanned int")

	var null Int32
	err = null.Scan(nil)
	maybePanic(err)
	assertNullInt32(t, null, "scanned null")
}

func assertInt32(t *testing.T, i Int32, from string) {
	if i.Int32 != 12345 {
		t.Errorf("bad %s int: %d ≠ %d\n", from, i.Int32, 12345)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullInt32(t *testing.T, i Int32, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
