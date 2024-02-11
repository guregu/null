package zero

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/guregu/null/v5/internal"
)

func TestValue(t *testing.T) {
	testValue[string](t, "hello")
	testValue[uint32](t, 1337)
	testValue[uint64](t, 42)

	type myint int
	testValue[myint](t, 2)
}

func testValue[T comparable](t *testing.T, good T) {
	t.Run(internal.TypeName[Value[T]](), func(t *testing.T) {
		var zero T

		// valid Value[T]
		testValueValid[T](t, good)

		// invalid Value[T]
		testValueNull[T](t, zero, good)
	})
}

func testValueValid[T comparable](t *testing.T, value T) {
	valid := NewValue(value, true)
	if valid.IsZero() {
		t.Errorf("%#v IsZero() should be false", valid)
	}
	validVF := ValueFrom(value)
	if !reflect.DeepEqual(valid, validVF) {
		t.Errorf("%#v != %#v", valid, validVF)
	}
	validVFP := ValueFromPtr(&value)
	if !reflect.DeepEqual(valid, validVFP) {
		t.Errorf("%#v != %#v", valid, validVFP)
	}

	validp := valid.Ptr()
	if validp == nil {
		t.Errorf("%#v Ptr() shouldn't be nil", valid)
	}

	validVOZ := valid.ValueOrZero()
	if !reflect.DeepEqual(validVOZ, value) {
		t.Error("ValueOrZero() want:", value, "got:", validVOZ)
	}

	t.Run("MarshalJSON", func(t *testing.T) {
		wantJSON, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		got, err := json.Marshal(valid)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(wantJSON, got) {
			t.Error("unexpected json. want:", string(wantJSON), "got:", string(got))
		}

		t.Run("UnmarshalJSON", func(t *testing.T) {
			var want T
			if err := json.Unmarshal(wantJSON, &want); err != nil {
				t.Fatal(err)
			}
			var got Value[T]
			if err := json.Unmarshal(wantJSON, &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(want, got.V) {
				t.Error("bad unmarshal. want:", want, "got:", got)
			}
			if got.IsZero() {
				t.Errorf("%#v IsZero() should be false", got)
			}
		})
	})

	t.Run(fmt.Sprintf("Scan(%v)", value), func(t *testing.T) {
		var want sql.Null[T]
		if err := want.Scan(value); err != nil {
			t.Fatal(err)
		}
		var got Value[T]
		if err := got.Scan(value); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(want, got.Null) {
			t.Error("bad scan. want:", want, "got:", got)
		}
	})
}

func testValueNull[T comparable](t *testing.T, value T, good T) {
	var zero T
	var nilv *T

	null := NewValue(zero, false)
	if !null.IsZero() {
		t.Errorf("%v IsZero() should be true", null)
	}
	nullVFP := ValueFromPtr(nilv)
	if !reflect.DeepEqual(null, nullVFP) {
		t.Errorf("%#v != %#v", null, nullVFP)
	}
	if !null.Equal(nullVFP) {
		t.Errorf("!%#v.Equal(%#v)", null, nullVFP)
	}

	nullVFPZ := ValueFromPtr(new(T))
	if !reflect.DeepEqual(null, nullVFPZ) {
		t.Errorf("%#v != %#v", null, nullVFPZ)
	}
	if !null.Equal(nullVFPZ) {
		t.Errorf("!%#v.Equal(%#v)", null, nullVFPZ)
	}

	nullp := null.Ptr()
	if nullp != nil {
		t.Errorf("%#v Ptr() should be nil", null)
	}

	nullVOZ := null.ValueOrZero()
	if !reflect.DeepEqual(nullVOZ, zero) {
		t.Error("ValueOrZero() want:", zero, "got:", nullVOZ)
	}

	t.Run("MarshalJSON", func(t *testing.T) {
		wantJSON, err := json.Marshal(nilv)
		if err != nil {
			t.Fatal(err)
		}
		got, err := json.Marshal(null)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(wantJSON, got) {
			t.Error("unexpected json. want:", string(wantJSON), "got:", string(got))
		}

		t.Run("UnmarshalJSON", func(t *testing.T) {
			var want T
			if err := json.Unmarshal(wantJSON, &want); err != nil {
				t.Fatal(err)
			}
			var got Value[T]
			if err := json.Unmarshal(wantJSON, &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(want, got.V) {
				t.Error("bad unmarshal. want:", want, "got:", got)
			}
			if !got.IsZero() {
				t.Errorf("%#v IsZero() should be true", got)
			}
		})
	})

	t.Run("Scan(nil)", func(t *testing.T) {
		var want sql.Null[T]
		if err := want.Scan(nil); err != nil {
			t.Fatal(err)
		}
		var got Value[T]
		if err := got.Scan(nil); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(want, got.Null) {
			t.Error("bad scan. want:", want, "got:", got)
		}
	})

	t.Run(fmt.Sprintf("SetValid(%v)", zero), func(t *testing.T) {
		valid2 := null
		valid2.SetValid(zero)
		if !valid2.IsZero() {
			t.Errorf("%#v IsZero() should be true", valid2)
		}

		valid3 := null
		valid3.SetValid(good)
		if valid3.IsZero() {
			t.Errorf("%#v IsZero() should be false", valid2)
		}
	})
}
