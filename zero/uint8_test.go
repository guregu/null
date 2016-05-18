package zero

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"
)

var (
	uint8JSON     = []byte(`254`)
	nullUint8JSON = []byte(`{"Uint8":254,"Valid":true}`)
	zeroU8JSON    = []byte(`0`)
)

func TestUint8From(t *testing.T) {
	i := Uint8From(254)
	assertUint8(t, i, "Uint8From()")

	zero := Uint8From(0)
	if zero.Valid {
		t.Error("Uint8From(0)", "is valid, but should be invalid")
	}
}

func TestUint8FromPtr(t *testing.T) {
	n := uint8(254)
	iptr := &n
	i := Uint8FromPtr(iptr)
	assertUint8(t, i, "Uint8FromPtr()")

	null := Uint8FromPtr(nil)
	assertNullUint8(t, null, "Uint8FromPtr(nil)")
}

func TestUnmarshalUint8(t *testing.T) {
	var i Uint8
	err := json.Unmarshal(uint8JSON, &i)
	maybePanic(err)
	assertUint8(t, i, "int json")

	var ni Uint8
	err = json.Unmarshal(nullUint8JSON, &ni)
	maybePanic(err)
	assertUint8(t, ni, "sql.NullUint8 json")

	var zero Uint8
	err = json.Unmarshal(zeroU8JSON, &zero)
	maybePanic(err)
	assertNullUint8(t, zero, "zero json")

	var null Uint8
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullUint8(t, null, "null json")

	var badType Uint8
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullUint8(t, badType, "wrong type json")

	var invalid Uint8
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullUint8(t, invalid, "invalid json")
}

func TestUnmarshalNonIntegerNumberU8(t *testing.T) {
	var i Uint8
	err := json.Unmarshal(floatJSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int")
	}
}

func TestUnmarshalUint8Overflow(t *testing.T) {
	int64Overflow := uint64(math.MaxUint8)

	// Max int64 should decode successfully
	var i Uint8
	err := json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	int64Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int64")
	}
}

func TestTextUnmarshalUint8(t *testing.T) {
	var i Uint8
	err := i.UnmarshalText([]byte("254"))
	maybePanic(err)
	assertUint8(t, i, "UnmarshalText() int")

	var zero Uint8
	err = zero.UnmarshalText([]byte("0"))
	maybePanic(err)
	assertNullUint8(t, zero, "UnmarshalText() zero int")

	var blank Uint8
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullUint8(t, blank, "UnmarshalText() empty int")

	var null Uint8
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullUint8(t, null, `UnmarshalText() "null"`)
}

func TestMarshalUint8(t *testing.T) {
	i := Uint8From(254)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "254", "non-empty json marshal")

	// invalid values should be encoded as 0
	null := NewUint8(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "0", "null json marshal")
}

func TestMarshalUint8Text(t *testing.T) {
	i := Uint8From(254)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "254", "non-empty text marshal")

	// invalid values should be encoded as zero
	null := NewUint8(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "0", "null text marshal")
}

func TestUint8Pointer(t *testing.T) {
	i := Uint8From(254)
	ptr := i.Ptr()
	if *ptr != 254 {
		t.Errorf("bad %s int: %#v ≠ %d\n", "pointer", ptr, 254)
	}

	null := NewUint8(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestUint8IsZero(t *testing.T) {
	i := Uint8From(254)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewUint8(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewUint8(0, true)
	if !zero.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestUint8Scan(t *testing.T) {
	var i Uint8
	err := i.Scan(254)
	maybePanic(err)
	assertUint8(t, i, "scanned int")

	var null Uint8
	err = null.Scan(nil)
	maybePanic(err)
	assertNullUint8(t, null, "scanned null")
}

func TestUint8SetValid(t *testing.T) {
	change := NewUint8(0, false)
	assertNullUint8(t, change, "SetValid()")
	change.SetValid(254)
	assertUint8(t, change, "SetValid()")
}

func assertUint8(t *testing.T, i Uint8, from string) {
	if i.Uint8 != 254 {
		t.Errorf("bad %s int: %d ≠ %d\n", from, i.Uint8, 254)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUint8(t *testing.T, i Uint8, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
