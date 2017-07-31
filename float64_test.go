package null

import (
	"encoding/json"
	"testing"
)

var (
	float64JSON = []byte(`1.2345`)
)

func TestFloat64From(t *testing.T) {
	f := Float64From(1.2345)
	assertFloat64(t, f, "Float64From()")

	zero := Float64From(0)
	if !zero.Valid {
		t.Error("Float64From(0)", "is invalid, but should be valid")
	}
}

func TestFloat64FromPtr(t *testing.T) {
	n := float64(1.2345)
	iptr := &n
	f := Float64FromPtr(iptr)
	assertFloat64(t, f, "Float64FromPtr()")

	null := Float64FromPtr(nil)
	assertNullFloat64(t, null, "Float64FromPtr(nil)")
}

func TestUnmarshalFloat64(t *testing.T) {
	var f Float64
	err := json.Unmarshal(float64JSON, &f)
	maybePanic(err)
	assertFloat64(t, f, "float64 json")

	var null Float64
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullFloat64(t, null, "null json")

	var badType Float64
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullFloat64(t, badType, "wrong type json")

	var invalid Float64
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
}

func TestTextUnmarshalFloat64(t *testing.T) {
	var f Float64
	err := f.UnmarshalText([]byte("1.2345"))
	maybePanic(err)
	assertFloat64(t, f, "UnmarshalText() float64")

	var blank Float64
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullFloat64(t, blank, "UnmarshalText() empty float64")
}

func TestMarshalFloat64(t *testing.T) {
	f := Float64From(1.2345)
	data, err := json.Marshal(f)
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewFloat64(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalFloat64Text(t *testing.T) {
	f := Float64From(1.2345)
	data, err := f.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewFloat64(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestFloat64Pointer(t *testing.T) {
	f := Float64From(1.2345)
	ptr := f.Ptr()
	if *ptr != 1.2345 {
		t.Errorf("bad %s float64: %#v ≠ %v\n", "pointer", ptr, 1.2345)
	}

	null := NewFloat64(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s float64: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestFloat64IsZero(t *testing.T) {
	f := Float64From(1.2345)
	if f.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewFloat64(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewFloat64(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestFloat64SetValid(t *testing.T) {
	change := NewFloat64(0, false)
	assertNullFloat64(t, change, "SetValid()")
	change.SetValid(1.2345)
	assertFloat64(t, change, "SetValid()")
}

func TestFloat64Scan(t *testing.T) {
	var f Float64
	err := f.Scan(1.2345)
	maybePanic(err)
	assertFloat64(t, f, "scanned float64")

	var null Float64
	err = null.Scan(nil)
	maybePanic(err)
	assertNullFloat64(t, null, "scanned null")
}

func assertFloat64(t *testing.T, f Float64, from string) {
	if f.Float64 != 1.2345 {
		t.Errorf("bad %s float64: %f ≠ %f\n", from, f.Float64, 1.2345)
	}
	if !f.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullFloat64(t *testing.T, f Float64, from string) {
	if f.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
