package zero

import (
	"bytes"
	"encoding/json"
	"testing"
)

var (
	jsonJSON = []byte(`"hello"`)
)

func TestJSONFrom(t *testing.T) {
	i := JSONFrom([]byte(`"hello"`))
	assertJSON(t, i, "JSONFrom()")

	zero := JSONFrom(nil)
	if zero.Valid {
		t.Error("JSONFrom(nil)", "is valid, but should be invalid")
	}

	zero = JSONFrom([]byte{})
	if zero.Valid {
		t.Error("JSONFrom([]byte{})", "is valid, but should be invalid")
	}
}

func TestJSONFromPtr(t *testing.T) {
	n := []byte(`"hello"`)
	iptr := &n
	i := JSONFromPtr(iptr)
	assertJSON(t, i, "JSONFromPtr()")

	null := JSONFromPtr(nil)
	assertNullJSON(t, null, "JSONFromPtr(nil)")
}

func TestUnmarshalJSON(t *testing.T) {
	var i JSON
	err := json.Unmarshal(bytesJSON, &i)
	maybePanic(err)
	assertJSON(t, i, "[]byte json")

	var ni JSON
	err = ni.UnmarshalJSON([]byte{})
	if ni.Valid == true {
		t.Errorf("expected Valid to be false, got true")
	}
	if !bytes.Equal(ni.JSON, []byte("null")) {
		t.Errorf("Expected JSON to be nil slice, but was not: %#v %#v", ni.JSON, []byte(`null`))
	}

	var null JSON
	err = null.UnmarshalJSON(nil)
	if null.Valid == true {
		t.Errorf("expected Valid to be false, got true")
	}
	if !bytes.Equal(null.JSON, []byte(`null`)) {
		t.Errorf("Expected JSON to be []byte nil, but was not: %#v %#v", null.JSON, []byte(`null`))
	}
}

func TestTextUnmarshalJSON(t *testing.T) {
	var i JSON
	err := i.UnmarshalText([]byte(`"hello"`))
	maybePanic(err)
	assertJSON(t, i, "UnmarshalText() []byte")

	var blank JSON
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullJSON(t, blank, "UnmarshalText() empty []byte")
}

func TestMarshalJSON(t *testing.T) {
	i := JSONFrom([]byte(`"hello"`))
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, `"hello"`, "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewJSON(nil, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalJSONText(t *testing.T) {
	i := JSONFrom([]byte(`"hello"`))
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, `"hello"`, "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewJSON(nil, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestJSONPointer(t *testing.T) {
	i := JSONFrom([]byte(`"hello"`))
	ptr := i.Ptr()
	if !bytes.Equal(*ptr, []byte(`"hello"`)) {
		t.Errorf("bad %s []byte: %#v ≠ %s\n", "pointer", ptr, `"hello"`)
	}

	null := NewJSON(nil, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s []byte: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestJSONIsZero(t *testing.T) {
	i := JSONFrom([]byte(`"hello"`))
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewJSON(nil, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewJSON(nil, true)
	if !zero.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	nz := NewJSON([]byte("thing"), true)
	if nz.IsZero() {
		t.Error("IsZero() should be false")
	}
}

func TestJSONSetValid(t *testing.T) {
	change := NewJSON(nil, false)
	assertNullJSON(t, change, "SetValid()")
	change.SetValid([]byte(`"hello"`))
	assertJSON(t, change, "SetValid()")
}

func TestJSONScan(t *testing.T) {
	var i JSON
	err := i.Scan(`"hello"`)
	maybePanic(err)
	assertJSON(t, i, "scanned []byte")

	var null JSON
	err = null.Scan(nil)
	maybePanic(err)
	assertNullJSON(t, null, "scanned null")
}

func assertJSON(t *testing.T, i JSON, from string) {
	if !bytes.Equal(i.JSON, []byte(`"hello"`)) {
		t.Errorf("bad %s []byte: %#v ≠ %#v\n", from, string(i.JSON), string([]byte(`"hello"`)))
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullJSON(t *testing.T, i JSON, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
