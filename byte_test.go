package null

import (
	"encoding/json"
	"testing"
)

var (
	byteJSON = []byte(`"b"`)
)

func TestByteFrom(t *testing.T) {
	i := ByteFrom('b')
	assertByte(t, i, "ByteFrom()")

	zero := ByteFrom(0)
	if !zero.Valid {
		t.Error("ByteFrom(0)", "is invalid, but should be valid")
	}
}

func TestByteFromPtr(t *testing.T) {
	n := byte('b')
	iptr := &n
	i := ByteFromPtr(iptr)
	assertByte(t, i, "ByteFromPtr()")

	null := ByteFromPtr(nil)
	assertNullByte(t, null, "ByteFromPtr(nil)")
}

func TestUnmarshalByte(t *testing.T) {
	var null Byte
	err := json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullByte(t, null, "null json")

	var badType Byte
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullByte(t, badType, "wrong type json")

	var invalid Byte
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullByte(t, invalid, "invalid json")
}

func TestUnmarshalNonByteegerNumber(t *testing.T) {
	var i Byte
	err := json.Unmarshal(float64JSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int")
	}
}

func TestTextUnmarshalByte(t *testing.T) {
	var i Byte
	err := i.UnmarshalText([]byte("b"))
	maybePanic(err)
	assertByte(t, i, "UnmarshalText() int")

	var blank Byte
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullByte(t, blank, "UnmarshalText() empty int")
}

func TestMarshalByte(t *testing.T) {
	i := ByteFrom('b')
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, `"b"`, "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewByte(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalByteText(t *testing.T) {
	i := ByteFrom('b')
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "b", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewByte(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestBytePointer(t *testing.T) {
	i := ByteFrom('b')
	ptr := i.Ptr()
	if *ptr != 'b' {
		t.Errorf("bad %s int: %#v ≠ %d\n", "pointer", ptr, 'b')
	}

	null := NewByte(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestByteIsZero(t *testing.T) {
	i := ByteFrom('b')
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewByte(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewByte(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestByteSetValid(t *testing.T) {
	change := NewByte(0, false)
	assertNullByte(t, change, "SetValid()")
	change.SetValid('b')
	assertByte(t, change, "SetValid()")
}

func TestByteScan(t *testing.T) {
	var i Byte
	err := i.Scan("b")
	maybePanic(err)
	assertByte(t, i, "scanned int")

	var null Byte
	err = null.Scan(nil)
	maybePanic(err)
	assertNullByte(t, null, "scanned null")
}

func assertByte(t *testing.T, i Byte, from string) {
	if i.Byte != 'b' {
		t.Errorf("bad %s int: %d ≠ %d\n", from, i.Byte, 'b')
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullByte(t *testing.T, i Byte, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
