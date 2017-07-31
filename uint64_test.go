package null

import (
	"encoding/json"
	"testing"
)

var (
	uint64JSON = []byte(`18446744073709551614`)
)

func TestUint64From(t *testing.T) {
	i := Uint64From(18446744073709551614)
	assertUint64(t, i, "Uint64From()")

	zero := Uint64From(0)
	if !zero.Valid {
		t.Error("Uint64From(0)", "is invalid, but should be valid")
	}
}

func TestUint64FromPtr(t *testing.T) {
	n := uint64(18446744073709551614)
	iptr := &n
	i := Uint64FromPtr(iptr)
	assertUint64(t, i, "Uint64FromPtr()")

	null := Uint64FromPtr(nil)
	assertNullUint64(t, null, "Uint64FromPtr(nil)")
}

func TestUnmarshalUint64(t *testing.T) {
	var i Uint64
	err := json.Unmarshal(uint64JSON, &i)
	maybePanic(err)
	assertUint64(t, i, "uint64 json")

	var null Uint64
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullUint64(t, null, "null json")

	var badType Uint64
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullUint64(t, badType, "wrong type json")

	var invalid Uint64
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullUint64(t, invalid, "invalid json")
}

func TestUnmarshalNonUintegerNumber64(t *testing.T) {
	var i Uint64
	err := json.Unmarshal(float64JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to uint64")
	}
}

func TestTextUnmarshalUint64(t *testing.T) {
	var i Uint64
	err := i.UnmarshalText([]byte("18446744073709551614"))
	maybePanic(err)
	assertUint64(t, i, "UnmarshalText() uint64")

	var blank Uint64
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullUint64(t, blank, "UnmarshalText() empty uint64")
}

func TestMarshalUint64(t *testing.T) {
	i := Uint64From(18446744073709551614)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "18446744073709551614", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewUint64(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalUint64Text(t *testing.T) {
	i := Uint64From(18446744073709551614)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "18446744073709551614", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewUint64(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestUint64Pointer(t *testing.T) {
	i := Uint64From(18446744073709551614)
	ptr := i.Ptr()
	if *ptr != 18446744073709551614 {
		t.Errorf("bad %s uint64: %#v ≠ %d\n", "pointer", ptr, uint64(18446744073709551614))
	}

	null := NewUint64(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s uint64: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestUint64IsZero(t *testing.T) {
	i := Uint64From(18446744073709551614)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewUint64(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewUint64(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestUint64SetValid(t *testing.T) {
	change := NewUint64(0, false)
	assertNullUint64(t, change, "SetValid()")
	change.SetValid(18446744073709551614)
	assertUint64(t, change, "SetValid()")
}

func TestUint64Scan(t *testing.T) {
	var i Uint64
	err := i.Scan(uint64(18446744073709551614))
	maybePanic(err)
	assertUint64(t, i, "scanned uint64")

	var null Uint64
	err = null.Scan(nil)
	maybePanic(err)
	assertNullUint64(t, null, "scanned null")
}

func assertUint64(t *testing.T, i Uint64, from string) {
	if i.Uint64 != 18446744073709551614 {
		t.Errorf("bad %s uint64: %d ≠ %d\n", from, i.Uint64, uint64(18446744073709551614))
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUint64(t *testing.T, i Uint64, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
