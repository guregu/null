package zero

import (
	"encoding/json"
	"testing"
)

var (
	float32JSON     = []byte(`1.2345`)
	nullFloat32JSON = []byte(`{"Float32":1.2345,"Valid":true}`)
)

func TestFloat32From(t *testing.T) {
	f := Float32From(1.2345)
	assertFloat32(t, f, "Float32From()")

	zero := Float32From(0)
	if zero.Valid {
		t.Error("Float32From(0)", "is valid, but should be invalid")
	}
}

func TestFloat32FromPtr(t *testing.T) {
	n := float32(1.2345)
	iptr := &n
	f := Float32FromPtr(iptr)
	assertFloat32(t, f, "Float32FromPtr()")

	null := Float32FromPtr(nil)
	assertNullFloat32(t, null, "Float32FromPtr(nil)")
}

func TestUnmarshalFloat32(t *testing.T) {
	var f Float32
	err := json.Unmarshal(float32JSON, &f)
	maybePanic(err)
	assertFloat32(t, f, "float json")

	var nf Float32
	err = json.Unmarshal(nullFloat32JSON, &nf)
	maybePanic(err)
	assertFloat32(t, nf, "sql.NullFloat32 json")

	var zero Float32
	err = json.Unmarshal(zeroJSON, &zero)
	maybePanic(err)
	assertNullFloat32(t, zero, "zero json")

	var null Float32
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullFloat32(t, null, "null json")

	var badType Float32
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullFloat32(t, badType, "wrong type json")

	var invalid Float32
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullFloat32(t, invalid, "invalid json")
}

func TestTextUnmarshalFloat32(t *testing.T) {
	var f Float32
	err := f.UnmarshalText([]byte("1.2345"))
	maybePanic(err)
	assertFloat32(t, f, "UnmarshalText() float")

	var zero Float32
	err = zero.UnmarshalText([]byte("0"))
	maybePanic(err)
	assertNullFloat32(t, zero, "UnmarshalText() zero float")

	var blank Float32
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullFloat32(t, blank, "UnmarshalText() empty float")

	var null Float32
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullFloat32(t, null, `UnmarshalText() "null"`)
}

func TestMarshalFloat32(t *testing.T) {
	f := Float32From(1.2345)
	data, err := json.Marshal(f)
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty json marshal")

	// invalid values should be encoded as 0
	null := NewFloat32(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "0", "null json marshal")
}

func TestMarshalFloat32Text(t *testing.T) {
	f := Float32From(1.2345)
	data, err := f.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty text marshal")

	// invalid values should be encoded as zero
	null := NewFloat32(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "0", "null text marshal")
}

func TestFloat32Pointer(t *testing.T) {
	f := Float32From(1.2345)
	ptr := f.Ptr()
	if *ptr != 1.2345 {
		t.Errorf("bad %s Float32: %#v ≠ %v\n", "pointer", ptr, 1.2345)
	}

	null := NewFloat32(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s Float32: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestFloat32IsZero(t *testing.T) {
	f := Float32From(1.2345)
	if f.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewFloat32(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewFloat32(0, true)
	if !zero.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestFloat32SetValid(t *testing.T) {
	change := NewFloat32(0, false)
	assertNullFloat32(t, change, "SetValid()")
	change.SetValid(1.2345)
	assertFloat32(t, change, "SetValid()")
}

func TestFloat32Scan(t *testing.T) {
	var f Float32
	err := f.Scan(1.2345)
	maybePanic(err)
	assertFloat32(t, f, "scanned float")

	var null Float32
	err = null.Scan(nil)
	maybePanic(err)
	assertNullFloat32(t, null, "scanned null")
}

func assertFloat32(t *testing.T, f Float32, from string) {
	if f.Float32 != 1.2345 {
		t.Errorf("bad %s float: %f ≠ %f\n", from, f.Float32, 1.2345)
	}
	if !f.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullFloat32(t *testing.T, f Float32, from string) {
	if f.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
