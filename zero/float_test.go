package zero

import (
	"encoding/json"
	"errors"
	"math"
	"testing"
)

var (
	floatJSON       = []byte(`1.2345`)
	floatStringJSON = []byte(`"1.2345"`)
	floatBlankJSON  = []byte(`""`)
	nullFloatJSON   = []byte(`{"Float64":1.2345,"Valid":true}`)
)

func TestFloatFrom(t *testing.T) {
	f := FloatFrom(1.2345)
	assertFloat(t, f, "FloatFrom()")

	zero := FloatFrom(0)
	if zero.Valid {
		t.Error("FloatFrom(0)", "is valid, but should be invalid")
	}
}

func TestFloatFromPtr(t *testing.T) {
	n := float64(1.2345)
	iptr := &n
	f := FloatFromPtr(iptr)
	assertFloat(t, f, "FloatFromPtr()")

	null := FloatFromPtr(nil)
	assertNullFloat(t, null, "FloatFromPtr(nil)")
}

func TestUnmarshalFloat(t *testing.T) {
	var f Float
	err := json.Unmarshal(floatJSON, &f)
	maybePanic(err)
	assertFloat(t, f, "float json")

	var sf Float
	err = json.Unmarshal(floatStringJSON, &sf)
	maybePanic(err)
	assertFloat(t, sf, "string float json")

	var nf Float
	err = json.Unmarshal(nullFloatJSON, &nf)
	if err == nil {
		panic("err should not be nil")
	}

	var blank Float
	err = json.Unmarshal(floatBlankJSON, &blank)
	if err == nil {
		panic("err should not be nil")
	}

	var zero Float
	err = json.Unmarshal(zeroJSON, &zero)
	maybePanic(err)
	assertNullFloat(t, zero, "zero json")

	var null Float
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullFloat(t, null, "null json")

	var badType Float
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullFloat(t, badType, "wrong type json")

	var invalid Float
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
	assertNullFloat(t, invalid, "invalid json")
}

func TestTextUnmarshalFloat(t *testing.T) {
	var f Float
	err := f.UnmarshalText([]byte("1.2345"))
	maybePanic(err)
	assertFloat(t, f, "UnmarshalText() float")

	var zero Float
	err = zero.UnmarshalText([]byte("0"))
	maybePanic(err)
	assertNullFloat(t, zero, "UnmarshalText() zero float")

	var blank Float
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullFloat(t, blank, "UnmarshalText() empty float")

	var null Float
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullFloat(t, null, `UnmarshalText() "null"`)

	var invalid Float
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		panic("expected error")
	}
}

func TestMarshalFloat(t *testing.T) {
	f := FloatFrom(1.2345)
	data, err := json.Marshal(f)
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty json marshal")

	// invalid values should be encoded as 0
	null := NewFloat(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "0", "null json marshal")
}

func TestMarshalFloatText(t *testing.T) {
	f := FloatFrom(1.2345)
	data, err := f.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty text marshal")

	// invalid values should be encoded as zero
	null := NewFloat(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "0", "null text marshal")
}

func TestFloatPointer(t *testing.T) {
	f := FloatFrom(1.2345)
	ptr := f.Ptr()
	if *ptr != 1.2345 {
		t.Errorf("bad %s Float: %#v ≠ %v\n", "pointer", ptr, 1.2345)
	}

	null := NewFloat(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s Float: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestFloatIsZero(t *testing.T) {
	f := FloatFrom(1.2345)
	if f.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewFloat(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewFloat(0, true)
	if !zero.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestFloatSetValid(t *testing.T) {
	change := NewFloat(0, false)
	assertNullFloat(t, change, "SetValid()")
	change.SetValid(1.2345)
	assertFloat(t, change, "SetValid()")
}

func TestFloatScan(t *testing.T) {
	var f Float
	err := f.Scan(1.2345)
	maybePanic(err)
	assertFloat(t, f, "scanned float")

	var null Float
	err = null.Scan(nil)
	maybePanic(err)
	assertNullFloat(t, null, "scanned null")
}

func TestFloatInfNaN(t *testing.T) {
	nan := NewFloat(math.NaN(), true)
	_, err := nan.MarshalJSON()
	if err == nil {
		t.Error("expected error for NaN, got nil")
	}

	inf := NewFloat(math.Inf(1), true)
	_, err = inf.MarshalJSON()
	if err == nil {
		t.Error("expected error for Inf, got nil")
	}
}

func TestFloatValueOrZero(t *testing.T) {
	valid := NewFloat(1.2345, true)
	if valid.ValueOrZero() != 1.2345 {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := NewFloat(1.2345, false)
	if invalid.ValueOrZero() != 0 {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestFloatEqual(t *testing.T) {
	f1 := NewFloat(10, false)
	f2 := NewFloat(10, false)
	assertFloatEqualIsTrue(t, f1, f2)

	f1 = NewFloat(10, false)
	f2 = NewFloat(20, false)
	assertFloatEqualIsTrue(t, f1, f2)

	f1 = NewFloat(10, true)
	f2 = NewFloat(10, true)
	assertFloatEqualIsTrue(t, f1, f2)

	f1 = NewFloat(10, false)
	f2 = NewFloat(0, true)
	assertFloatEqualIsTrue(t, f1, f2)

	f1 = NewFloat(10, true)
	f2 = NewFloat(10, false)
	assertFloatEqualIsFalse(t, f1, f2)

	f1 = NewFloat(10, false)
	f2 = NewFloat(10, true)
	assertFloatEqualIsFalse(t, f1, f2)

	f1 = NewFloat(10, true)
	f2 = NewFloat(20, true)
	assertFloatEqualIsFalse(t, f1, f2)
}

func assertFloat(t *testing.T, f Float, from string) {
	if f.Float64 != 1.2345 {
		t.Errorf("bad %s float: %f ≠ %f\n", from, f.Float64, 1.2345)
	}
	if !f.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullFloat(t *testing.T, f Float, from string) {
	if f.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertFloatEqualIsTrue(t *testing.T, a, b Float) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf("Equal() of Float{%v, Valid:%t} and Float{%v, Valid:%t} should return true", a.Float64, a.Valid, b.Float64, b.Valid)
	}
}

func assertFloatEqualIsFalse(t *testing.T, a, b Float) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf("Equal() of Float{%v, Valid:%t} and Float{%v, Valid:%t} should return false", a.Float64, a.Valid, b.Float64, b.Valid)
	}
}
