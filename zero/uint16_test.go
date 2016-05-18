package zero

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"
)

var (
	uint16JSON     = []byte(`65534`)
	nullUint16JSON = []byte(`{"Uint16":65534,"Valid":true}`)
	zeroU16JSON    = []byte(`0`)
)

func TestUint16From(t *testing.T) {
	i := Uint16From(65534)
	assertUint16(t, i, "Uint16From()")

	zero := Uint16From(0)
	if zero.Valid {
		t.Error("Uint16From(0)", "is valid, but should be invalid")
	}
}

func TestUint16FromPtr(t *testing.T) {
	n := uint16(65534)
	iptr := &n
	i := Uint16FromPtr(iptr)
	assertUint16(t, i, "Uint16FromPtr()")

	null := Uint16FromPtr(nil)
	assertNullUint16(t, null, "Uint16FromPtr(nil)")
}

func TestUnmarshalUint16(t *testing.T) {
	var i Uint16
	err := json.Unmarshal(uint16JSON, &i)
	maybePanic(err)
	assertUint16(t, i, "int json")

	var ni Uint16
	err = json.Unmarshal(nullUint16JSON, &ni)
	maybePanic(err)
	assertUint16(t, ni, "sql.NullUint16 json")

	var zero Uint16
	err = json.Unmarshal(zeroU16JSON, &zero)
	maybePanic(err)
	assertNullUint16(t, zero, "zero json")

	var null Uint16
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullUint16(t, null, "null json")

	var badType Uint16
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullUint16(t, badType, "wrong type json")

	var invalid Uint16
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullUint16(t, invalid, "invalid json")
}

func TestUnmarshalNonIntegerNumberU16(t *testing.T) {
	var i Uint16
	err := json.Unmarshal(floatJSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int")
	}
}

func TestUnmarshalUint16Overflow(t *testing.T) {
	int64Overflow := uint64(math.MaxUint16)

	// Max int64 should decode successfully
	var i Uint16
	err := json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	int64Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int64")
	}
}

func TestTextUnmarshalUint16(t *testing.T) {
	var i Uint16
	err := i.UnmarshalText([]byte("65534"))
	maybePanic(err)
	assertUint16(t, i, "UnmarshalText() int")

	var zero Uint16
	err = zero.UnmarshalText([]byte("0"))
	maybePanic(err)
	assertNullUint16(t, zero, "UnmarshalText() zero int")

	var blank Uint16
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullUint16(t, blank, "UnmarshalText() empty int")

	var null Uint16
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullUint16(t, null, `UnmarshalText() "null"`)
}

func TestMarshalUint16(t *testing.T) {
	i := Uint16From(65534)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "65534", "non-empty json marshal")

	// invalid values should be encoded as 0
	null := NewUint16(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "0", "null json marshal")
}

func TestMarshalUint16Text(t *testing.T) {
	i := Uint16From(65534)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "65534", "non-empty text marshal")

	// invalid values should be encoded as zero
	null := NewUint16(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "0", "null text marshal")
}

func TestUint16Pointer(t *testing.T) {
	i := Uint16From(65534)
	ptr := i.Ptr()
	if *ptr != 65534 {
		t.Errorf("bad %s int: %#v ≠ %d\n", "pointer", ptr, 65534)
	}

	null := NewUint16(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestUint16IsZero(t *testing.T) {
	i := Uint16From(65534)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewUint16(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewUint16(0, true)
	if !zero.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestUint16Scan(t *testing.T) {
	var i Uint16
	err := i.Scan(65534)
	maybePanic(err)
	assertUint16(t, i, "scanned int")

	var null Uint16
	err = null.Scan(nil)
	maybePanic(err)
	assertNullUint16(t, null, "scanned null")
}

func TestUint16SetValid(t *testing.T) {
	change := NewUint16(0, false)
	assertNullUint16(t, change, "SetValid()")
	change.SetValid(65534)
	assertUint16(t, change, "SetValid()")
}

func assertUint16(t *testing.T, i Uint16, from string) {
	if i.Uint16 != 65534 {
		t.Errorf("bad %s int: %d ≠ %d\n", from, i.Uint16, 65534)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUint16(t *testing.T, i Uint16, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
