package null

import (
	"encoding"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"testing"

	"github.com/guregu/null/v5/internal"
)

var (
	intJSON       = []byte(`123`)
	intStringJSON = []byte(`"123"`)
)

type nullint interface {
	Int | Int32 | Int16 | Byte
	IsZero() bool
	value() (int64, bool)
}

func TestIntFrom(t *testing.T) {
	testIntFrom(t, IntFrom)
	testIntFrom(t, Int32From)
	testIntFrom(t, Int16From)
	testIntFrom(t, ByteFrom)
}

func testIntFrom[N nullint, V internal.Integer](t *testing.T, from func(V) N) {
	t.Run(internal.TypeName[N](), func(t *testing.T) {
		i := from(123)
		assertInt(t, i, "from(123)")

		zero := from(0)
		_, valid := zero.value()
		if !valid {
			t.Error("from(0)", "is invalid, but should be valid")
		}
	})
}

func TestIntFromPtr(t *testing.T) {
	testIntFromPtr(t, IntFromPtr)
	testIntFromPtr(t, Int32FromPtr)
	testIntFromPtr(t, Int16FromPtr)
	testIntFromPtr(t, ByteFromPtr)
}

func testIntFromPtr[N nullint, V internal.Integer](t *testing.T, fromPtr func(*V) N) {
	t.Run(internal.TypeName[N](), func(t *testing.T) {
		n := V(123)
		iptr := &n
		i := fromPtr(iptr)
		assertInt(t, i, "fromPtr()")

		null := fromPtr(nil)
		assertNullInt(t, null, "fromPtr(nil)")
	})
}

func TestUnmarshalInt(t *testing.T) {
	testUnmarshalInt[Int](t)
	testUnmarshalInt[Int32](t)
	testUnmarshalInt[Int16](t)
	testUnmarshalInt[Byte](t)
}

func testUnmarshalInt[N nullint](t *testing.T) {
	t.Run(internal.TypeName[N](), func(t *testing.T) {
		var i N
		err := json.Unmarshal(intJSON, &i)
		maybePanic(err)
		assertInt(t, i, "int json")

		var si N
		err = json.Unmarshal(intStringJSON, &si)
		maybePanic(err)
		assertInt(t, si, "int string json")

		var bi N
		err = json.Unmarshal(floatBlankJSON, &bi)
		maybePanic(err)
		assertNullInt(t, bi, "blank json")

		var null N
		err = json.Unmarshal(nullJSON, &null)
		maybePanic(err)
		assertNullInt(t, null, "null json")

		var badType N
		err = json.Unmarshal(boolJSON, &badType)
		if err == nil {
			panic("err should not be nil")
		}
		assertNullInt(t, badType, "wrong type json")

		var invalid N
		err = json.Unmarshal(invalidJSON, &invalid)
		var syntaxError *json.SyntaxError
		if !errors.As(err, &syntaxError) {
			t.Errorf("expected wrapped json.SyntaxError, not %T", err)
		}
		assertNullInt(t, invalid, "invalid json")
	})
}

func TestUnmarshalNonIntegerNumber(t *testing.T) {
	var i Int
	err := json.Unmarshal(floatJSON, &i)
	if err == nil {
		panic("err should be present; non-internal.Integer number coerced to int")
	}
}

func TestUnmarshalIntOverflow(t *testing.T) {
	testUnmarshalIntOverflow[Int, int64](t, math.MaxInt64)
	testUnmarshalIntOverflow[Int32, int32](t, math.MaxInt32)
	testUnmarshalIntOverflow[Int16, int16](t, math.MaxInt16)
	testUnmarshalIntOverflow[Byte, byte](t, math.MaxUint8)
}

func testUnmarshalIntOverflow[N nullint, V internal.Integer](t *testing.T, max V) {
	t.Run(internal.TypeName[N](), func(t *testing.T) {
		overflow := uint64(max)

		// Max int64 should decode successfully
		var i N
		err := json.Unmarshal([]byte(strconv.FormatUint(overflow, 10)), &i)
		maybePanic(err)

		// Attempt to overflow
		overflow++
		err = json.Unmarshal([]byte(strconv.FormatUint(overflow, 10)), &i)
		if err == nil {
			t.Error("err should be present but isn't; decoded value overflows")
		}
	})
}

func TestTextUnmarshalInt(t *testing.T) {
	testTextUnmarshalInt(t, (*Int).UnmarshalText)
	testTextUnmarshalInt(t, (*Int32).UnmarshalText)
	testTextUnmarshalInt(t, (*Int16).UnmarshalText)
	testTextUnmarshalInt(t, (*Byte).UnmarshalText)
}

func testTextUnmarshalInt[N nullint](t *testing.T, unmarshal func(*N, []byte) error) {
	t.Run(internal.TypeName[N](), func(t *testing.T) {
		var i N
		err := unmarshal(&i, []byte("123"))
		maybePanic(err)
		assertInt(t, i, "unmarshal int")

		var blank N
		err = unmarshal(&blank, []byte(""))
		maybePanic(err)
		assertNullInt(t, blank, "unmarshal empty int")

		var null N
		err = unmarshal(&null, []byte("null"))
		maybePanic(err)
		assertNullInt(t, null, `unmarshal "null"`)

		var invalid N
		err = unmarshal(&invalid, []byte("hello world"))
		if err == nil {
			panic("expected error")
		}
	})
}

func TestMarshalInt(t *testing.T) {
	testMarshalInt(t, NewInt)
	testMarshalInt(t, NewInt32)
	testMarshalInt(t, NewInt16)
	testMarshalInt(t, NewByte)
}

func testMarshalInt[N interface{ ValueOrZero() V }, V internal.Integer](t *testing.T, newInt func(V, bool) N) {
	t.Run(internal.TypeName[N](), func(t *testing.T) {
		i := newInt(123, true)
		data, err := json.Marshal(i)
		maybePanic(err)
		assertJSONEquals(t, data, "123", "non-empty json marshal")

		// invalid values should be encoded as null
		null := newInt(0, false)
		data, err = json.Marshal(null)
		maybePanic(err)
		assertJSONEquals(t, data, "null", "null json marshal")
	})
}

func TestMarshalIntText(t *testing.T) {
	testMarshalIntText(t, NewInt)
	testMarshalIntText(t, NewInt32)
	testMarshalIntText(t, NewInt16)
	testMarshalIntText(t, NewByte)
}

func testMarshalIntText[N encoding.TextMarshaler, V internal.Integer](t *testing.T, newInt func(V, bool) N) {
	t.Run(internal.TypeName[N](), func(t *testing.T) {
		i := newInt(123, true)
		data, err := i.MarshalText()
		maybePanic(err)
		assertJSONEquals(t, data, "123", "non-empty text marshal")

		// invalid values should be encoded as null
		null := newInt(0, false)
		data, err = null.MarshalText()
		maybePanic(err)
		assertJSONEquals(t, data, "", "null text marshal")
	})
}

func TestIntPointer(t *testing.T) {
	testIntPointer(t, NewInt)
	testIntPointer(t, NewInt32)
	testIntPointer(t, NewInt16)
	testIntPointer(t, NewByte)
}

func testIntPointer[N interface{ Ptr() *V }, V internal.Integer](t *testing.T, newInt func(V, bool) N) {
	t.Run(internal.TypeName[N](), func(t *testing.T) {
		i := newInt(123, true)
		ptr := i.Ptr()
		if *ptr != 123 {
			t.Errorf("bad %s int: %#v ≠ %d\n", "pointer", ptr, 123)
		}

		null := newInt(0, false)
		ptr = null.Ptr()
		if ptr != nil {
			t.Errorf("bad %s int: %#v ≠ %s\n", "nil pointer", ptr, "nil")
		}
	})
}

func TestIntIsZero(t *testing.T) {
	testIntIsZero(t, NewInt)
	testIntIsZero(t, NewInt32)
	testIntIsZero(t, NewInt16)
	testIntIsZero(t, NewByte)
}

func testIntIsZero[N nullint, V internal.Integer](t *testing.T, newInt func(V, bool) N) {
	t.Run(internal.TypeName[N](), func(t *testing.T) {
		i := newInt(123, true)
		if i.IsZero() {
			t.Errorf("IsZero() should be false")
		}

		null := newInt(0, false)
		if !null.IsZero() {
			t.Errorf("IsZero() should be true")
		}

		zero := newInt(0, true)
		if zero.IsZero() {
			t.Errorf("IsZero() should be false")
		}
	})
}

func TestIntSetValid(t *testing.T) {
	testIntSetValid(t, NewInt, (*Int).SetValid)
	testIntSetValid(t, NewInt32, (*Int32).SetValid)
	testIntSetValid(t, NewInt16, (*Int16).SetValid)
	testIntSetValid(t, NewByte, (*Byte).SetValid)
}

func testIntSetValid[N nullint, V internal.Integer](t *testing.T, newInt func(V, bool) N, setValid func(*N, V)) {
	t.Run(internal.TypeName[N](), func(t *testing.T) {
		change := newInt(0, false)
		assertNullInt(t, change, "SetValid()")
		setValid(&change, 123)
		assertInt(t, change, "SetValid()")
	})
}

func TestIntScan(t *testing.T) {
	testIntScan(t, (*Int).Scan)
	testIntScan(t, (*Int32).Scan)
	testIntScan(t, (*Int16).Scan)
	testIntScan(t, (*Byte).Scan)
}

func testIntScan[N nullint](t *testing.T, scan func(*N, any) error) {
	t.Run(internal.TypeName[N](), func(t *testing.T) {
		var i N
		err := scan(&i, 123)
		maybePanic(err)
		assertInt(t, i, "scanned int")

		var null N
		err = scan(&null, nil)
		maybePanic(err)
		assertNullInt(t, null, "scanned null")
	})
}

func TestIntValueOrZero(t *testing.T) {
	testIntValueOrZero(t, NewInt)
	testIntValueOrZero(t, NewInt32)
	testIntValueOrZero(t, NewInt16)
	testIntValueOrZero(t, NewByte)
}

func testIntValueOrZero[N interface{ ValueOrZero() V }, V internal.Integer](t *testing.T, newInt func(V, bool) N) {
	t.Run(internal.TypeName[N](), func(t *testing.T) {
		valid := newInt(123, true)
		if valid.ValueOrZero() != 123 {
			t.Error("unexpected ValueOrZero", valid.ValueOrZero())
		}

		invalid := newInt(123, false)
		if invalid.ValueOrZero() != 0 {
			t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
		}
	})
}

func TestIntEqual(t *testing.T) {
	testIntEqual(t, NewInt)
	testIntEqual(t, NewInt32)
	testIntEqual(t, NewInt16)
	testIntEqual(t, NewByte)
}

func testIntEqual[N interface{ Equal(N) bool }, V internal.Integer](t *testing.T, newInt func(V, bool) N) {
	t.Run(internal.TypeName[N](), func(t *testing.T) {
		int1 := newInt(10, false)
		int2 := newInt(10, false)
		assertIntEqualIsTrue(t, int1, int2)

		int1 = newInt(10, false)
		int2 = newInt(20, false)
		assertIntEqualIsTrue(t, int1, int2)

		int1 = newInt(10, true)
		int2 = newInt(10, true)
		assertIntEqualIsTrue(t, int1, int2)

		int1 = newInt(10, true)
		int2 = newInt(10, false)
		assertIntEqualIsFalse(t, int1, int2)

		int1 = newInt(10, false)
		int2 = newInt(10, true)
		assertIntEqualIsFalse(t, int1, int2)

		int1 = newInt(10, true)
		int2 = newInt(20, true)
		assertIntEqualIsFalse(t, int1, int2)
	})
}

func assertInt(t *testing.T, i interface{ value() (int64, bool) }, from string) {
	t.Helper()
	n, valid := i.value()
	if n != 123 {
		t.Errorf("bad %s int: %d ≠ %d\n", from, n, 123)
	}
	if !valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullInt(t *testing.T, i interface{ value() (int64, bool) }, from string) {
	t.Helper()
	_, valid := i.value()
	if valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertIntEqualIsTrue[N interface{ Equal(N) bool }](t *testing.T, a, b N) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf("Equal() of %#v and %#v should return true", a, b)
	}
}

func assertIntEqualIsFalse[N interface{ Equal(N) bool }](t *testing.T, a, b N) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf("Equal() of %#v and %#v should return false", a, b)
	}
}
