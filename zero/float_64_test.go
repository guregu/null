package zero

import (
	"encoding/json"
	"errors"
	"math"
	"testing"
)

var (
	float64JSON       = []byte(`1.2345`)
	float64StringJSON = []byte(`"1.2345"`)
	float64BlankJSON  = []byte(`""`)
	nullFloat64JSON   = []byte(`{"Float64":1.2345,"Valid":true}`)
)

func TestFloat64From(t *testing.T) {
	f := Float64From(1.2345)
	assertFloat64(t, f, "Float64From()")

	zero := Float64From(0)
	if zero.Valid {
		t.Error("Float64From(0)", "is valid, but should be invalid")
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
	assertFloat64(t, f, "float json")

	var sf Float64
	err = json.Unmarshal(float64StringJSON, &sf)
	maybePanic(err)
	assertFloat64(t, sf, "string float json")

	var nf Float64
	err = json.Unmarshal(nullFloat64JSON, &nf)
	if err == nil {
		panic("err should not be nil")
	}

	var blank Float64
	err = json.Unmarshal(float64BlankJSON, &blank)
	if err == nil {
		panic("err should not be nil")
	}

	var zero Float64
	err = json.Unmarshal(zeroJSON, &zero)
	maybePanic(err)
	assertNullFloat64(t, zero, "zero json")

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
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
	assertNullFloat64(t, invalid, "invalid json")
}

func TestTextUnmarshalFloat64(t *testing.T) {
	var f Float64
	err := f.UnmarshalText([]byte("1.2345"))
	maybePanic(err)
	assertFloat64(t, f, "UnmarshalText() float")

	var zero Float64
	err = zero.UnmarshalText([]byte("0"))
	maybePanic(err)
	assertNullFloat64(t, zero, "UnmarshalText() zero float")

	var blank Float64
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullFloat64(t, blank, "UnmarshalText() empty float")

	var null Float64
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullFloat64(t, null, `UnmarshalText() "null"`)

	var invalid Float64
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		panic("expected error")
	}
}

func TestMarshalFloat64(t *testing.T) {
	f := Float64From(1.2345)
	data, err := json.Marshal(f)
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty json marshal")

	// invalid values should be encoded as 0
	null := NewFloat64(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "0", "null json marshal")
}

func TestMarshalFloat64Text(t *testing.T) {
	f := Float64From(1.2345)
	data, err := f.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty text marshal")

	// invalid values should be encoded as zero
	null := NewFloat64(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "0", "null text marshal")
}

func TestFloat64Pointer(t *testing.T) {
	f := Float64From(1.2345)
	ptr := f.Ptr()
	if *ptr != 1.2345 {
		t.Errorf("bad %s Float64: %#v ≠ %v\n", "pointer", ptr, 1.2345)
	}

	null := NewFloat64(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s Float64: %#v ≠ %s\n", "nil pointer", ptr, "nil")
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
	if !zero.IsZero() {
		t.Errorf("IsZero() should be true")
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
	assertFloat64(t, f, "scanned float")

	var null Float64
	err = null.Scan(nil)
	maybePanic(err)
	assertNullFloat64(t, null, "scanned null")
}

func TestFloat64InfNaN(t *testing.T) {
	nan := NewFloat64(math.NaN(), true)
	_, err := nan.MarshalJSON()
	if err == nil {
		t.Error("expected error for NaN, got nil")
	}

	inf := NewFloat64(math.Inf(1), true)
	_, err = inf.MarshalJSON()
	if err == nil {
		t.Error("expected error for Inf, got nil")
	}
}

func TestFloat64ValueOrZero(t *testing.T) {
	valid := NewFloat64(1.2345, true)
	if valid.ValueOrZero() != 1.2345 {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := NewFloat64(1.2345, false)
	if invalid.ValueOrZero() != 0 {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestFloat64Equal(t *testing.T) {
	f1 := NewFloat64(10, false)
	f2 := NewFloat64(10, false)
	assertFloat64EqualIsTrue(t, f1, f2)

	f1 = NewFloat64(10, false)
	f2 = NewFloat64(20, false)
	assertFloat64EqualIsTrue(t, f1, f2)

	f1 = NewFloat64(10, true)
	f2 = NewFloat64(10, true)
	assertFloat64EqualIsTrue(t, f1, f2)

	f1 = NewFloat64(10, false)
	f2 = NewFloat64(0, true)
	assertFloat64EqualIsTrue(t, f1, f2)

	f1 = NewFloat64(10, true)
	f2 = NewFloat64(10, false)
	assertFloat64EqualIsFalse(t, f1, f2)

	f1 = NewFloat64(10, false)
	f2 = NewFloat64(10, true)
	assertFloat64EqualIsFalse(t, f1, f2)

	f1 = NewFloat64(10, true)
	f2 = NewFloat64(20, true)
	assertFloat64EqualIsFalse(t, f1, f2)
}

func assertFloat64(t *testing.T, f Float64, from string) {
	if f.Float64 != 1.2345 {
		t.Errorf("bad %s float: %f ≠ %f\n", from, f.Float64, 1.2345)
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

func assertFloat64EqualIsTrue(t *testing.T, a, b Float64) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf("Equal() of Float64{%v, Valid:%t} and Float64{%v, Valid:%t} should return true", a.Float64, a.Valid, b.Float64, b.Valid)
	}
}

func assertFloat64EqualIsFalse(t *testing.T, a, b Float64) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf("Equal() of Float64{%v, Valid:%t} and Float64{%v, Valid:%t} should return false", a.Float64, a.Valid, b.Float64, b.Valid)
	}
}
