package null

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"
)

var (
	uint32JSON = []byte(`4294967294`)
)

func TestUint32From(t *testing.T) {
	i := Uint32From(4294967294)
	assertUint32(t, i, "Uint32From()")

	zero := Uint32From(0)
	if !zero.Valid {
		t.Error("Uint32From(0)", "is invalid, but should be valid")
	}
}

func TestUint32FromPtr(t *testing.T) {
	n := uint32(4294967294)
	iptr := &n
	i := Uint32FromPtr(iptr)
	assertUint32(t, i, "Uint32FromPtr()")

	null := Uint32FromPtr(nil)
	assertNullUint32(t, null, "Uint32FromPtr(nil)")
}

func TestUnmarshalUint32(t *testing.T) {
	var i Uint32
	err := json.Unmarshal(uint32JSON, &i)
	maybePanic(err)
	assertUint32(t, i, "uint32 json")

	var null Uint32
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullUint32(t, null, "null json")

	var badType Uint32
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullUint32(t, badType, "wrong type json")

	var invalid Uint32
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullUint32(t, invalid, "invalid json")
}

func TestUnmarshalNonUintegerNumber32(t *testing.T) {
	var i Uint32
	err := json.Unmarshal(float64JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to uint32")
	}
}

func TestUnmarshalUint32Overflow(t *testing.T) {
	uint32Overflow := int64(math.MaxUint32)

	// Max uint32 should decode successfully
	var i Uint32
	err := json.Unmarshal([]byte(strconv.FormatUint(uint64(uint32Overflow), 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	uint32Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(uint64(uint32Overflow), 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows uint32")
	}
}

func TestTextUnmarshalUint32(t *testing.T) {
	var i Uint32
	err := i.UnmarshalText([]byte("4294967294"))
	maybePanic(err)
	assertUint32(t, i, "UnmarshalText() uint32")

	var blank Uint32
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullUint32(t, blank, "UnmarshalText() empty uint32")
}

func TestMarshalUint32(t *testing.T) {
	i := Uint32From(4294967294)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "4294967294", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewUint32(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalUint32Text(t *testing.T) {
	i := Uint32From(4294967294)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "4294967294", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewUint32(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestUint32Pointer(t *testing.T) {
	i := Uint32From(4294967294)
	ptr := i.Ptr()
	if *ptr != 4294967294 {
		t.Errorf("bad %s uint32: %#v ≠ %d\n", "pointer", ptr, 4294967294)
	}

	null := NewUint32(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s uint32: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestUint32IsZero(t *testing.T) {
	i := Uint32From(4294967294)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewUint32(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewUint32(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestUint32SetValid(t *testing.T) {
	change := NewUint32(0, false)
	assertNullUint32(t, change, "SetValid()")
	change.SetValid(4294967294)
	assertUint32(t, change, "SetValid()")
}

func TestUint32Scan(t *testing.T) {
	var i Uint32
	err := i.Scan(4294967294)
	maybePanic(err)
	assertUint32(t, i, "scanned uint32")

	var null Uint32
	err = null.Scan(nil)
	maybePanic(err)
	assertNullUint32(t, null, "scanned null")
}

func assertUint32(t *testing.T, i Uint32, from string) {
	if i.Uint32 != 4294967294 {
		t.Errorf("bad %s uint32: %d ≠ %d\n", from, i.Uint32, 4294967294)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUint32(t *testing.T, i Uint32, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
