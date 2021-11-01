package null

import (
	"encoding/json"
	"math"
	"reflect"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
)

var (
	intJSON       = []byte(`12345`)
	intJSONString = []byte(`"12345"`)
	nullIntJSON   = []byte(`  { "Int64":12345,"Valid":true}`)
)

func TestIntFrom(t *testing.T) {
	i := IntFrom(12345)
	assertInt(t, i, "IntFrom()")

	zero := IntFrom(0)
	if !zero.Valid {
		t.Error("IntFrom(0)", "is invalid, but should be valid")
	}
}

func TestIntFromPtr(t *testing.T) {
	n := int64(12345)
	iptr := &n
	i := IntFromPtr(iptr)
	assertInt(t, i, "IntFromPtr()")

	null := IntFromPtr(nil)
	assertNullInt(t, null, "IntFromPtr(nil)")
}

func TestIntUnmarshal(t *testing.T) {
	tests := []struct {
		in             []byte
		exp            Int
		expErrType     reflect.Type
		expErrTypeEasy reflect.Type
	}{
		{
			in:  intJSON,
			exp: IntFrom(12345),
		},
		{
			in:  intJSONString,
			exp: IntFrom(12345),
		},
		{
			in: []byte(` "12345"  	 `),
			exp: IntFrom(12345),
		},
		{
			in:  nullIntJSON,
			exp: IntFrom(12345),
		},
		{
			in: nullJSON,
		},
		{
			in:             boolJSON,
			expErrType:     reflect.TypeOf((*strconv.NumError)(nil)),
			expErrTypeEasy: reflect.TypeOf((*jlexer.LexerError)(nil)),
		},
		{
			in:             invalidJSON,
			expErrType:     reflect.TypeOf((*json.SyntaxError)(nil)),
			expErrTypeEasy: reflect.TypeOf((*jlexer.LexerError)(nil)),
		},
		{
			in:             []byte(`{"Int64":true,"Valid":true}`),
			expErrType:     reflect.TypeOf((*json.UnmarshalTypeError)(nil)),
			expErrTypeEasy: reflect.TypeOf((*jlexer.LexerError)(nil)),
		},
	}

	for _, test := range tests {
		t.Run(string(test.in), func(t *testing.T) {
			var i Int
			err := json.Unmarshal(test.in, &i)
			if err != nil {
				if test.expErrType == nil {
					t.Fatal(err)
				}
				if reflect.TypeOf(err) != test.expErrType {
					t.Fatalf("error %s(%T) is not of type %s", err, err, test.expErrType)
				}

			} else if test.expErrType != nil {
				t.Fatal("expected an error")
			}
			if diff := cmp.Diff(test.exp, i); diff != "" {
				t.Fatalf("result not as expected. %s", diff)
			}
		})

		t.Run(string(test.in)+"_easyjson", func(t *testing.T) {
			var i Int
			var err error
			allocs := testing.AllocsPerRun(10, func() {
				err = easyjson.Unmarshal(test.in, &i)
			})
			if err != nil {
				if test.expErrTypeEasy == nil {
					t.Fatal(err)
				}
				if reflect.TypeOf(err) != test.expErrTypeEasy {
					t.Fatalf("error %s(%T) is not of type %s", err, err, test.expErrTypeEasy)
				}

			} else if test.expErrTypeEasy != nil {
				t.Fatal("expected an error")
			}
			if test.expErrTypeEasy == nil && allocs > 0 {
				t.Fatalf("easyjson made %.0f allocations unmarshalling %T from: %s", allocs, i, test.in)
			}
			if diff := cmp.Diff(test.exp, i); diff != "" {
				t.Fatalf("result not as expected. %s", diff)
			}
		})
	}
}

func BenchmarkIntUnmarshal(b *testing.B) {
	tests := [][]byte{
		intJSON,
		intJSONString,
		[]byte(` "12345"  	 `),
		nullIntJSON,
		nullJSON,
	}

	for _, test := range tests {
		b.Run(string(test), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var ii Int
				if err := json.Unmarshal(test, &ii); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("easy "+string(test), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				w := &jlexer.Lexer{Data: test}
				var ii Int
				ii.UnmarshalEasyJSON(w)
				if err := w.Error(); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func TestUnmarshalNonIntegerNumber(t *testing.T) {
	var i Int
	err := json.Unmarshal(floatJSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int")
	}
}

func TestUnmarshalInt64Overflow(t *testing.T) {
	int64Overflow := uint64(math.MaxInt64)

	// Max int64 should decode successfully
	var i Int
	err := json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	int64Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int64")
	}
}

func TestTextUnmarshalInt(t *testing.T) {
	var i Int
	err := i.UnmarshalText([]byte("12345"))
	maybePanic(err)
	assertInt(t, i, "UnmarshalText() int")

	var blank Int
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullInt(t, blank, "UnmarshalText() empty int")

	var null Int
	err = null.UnmarshalText(nullLiteral)
	maybePanic(err)
	assertNullInt(t, null, `UnmarshalText() "null"`)

	var invalid Int
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		panic("expected error")
	}
}

func TestMarshalInt(t *testing.T) {
	i := IntFrom(12345)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "12345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewInt(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalIntText(t *testing.T) {
	i := IntFrom(12345)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "12345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewInt(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestIntPointer(t *testing.T) {
	i := IntFrom(12345)
	ptr := i.Ptr()
	if *ptr != 12345 {
		t.Errorf("bad %s int: %#v ≠ %d\n", "pointer", ptr, 12345)
	}

	null := NewInt(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestIntIsZero(t *testing.T) {
	i := IntFrom(12345)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewInt(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewInt(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestIntSetValid(t *testing.T) {
	change := NewInt(0, false)
	assertNullInt(t, change, "SetValid()")
	change.SetValid(12345)
	assertInt(t, change, "SetValid()")
}

func TestIntScan(t *testing.T) {
	var i Int
	err := i.Scan(12345)
	maybePanic(err)
	assertInt(t, i, "scanned int")

	var null Int
	err = null.Scan(nil)
	maybePanic(err)
	assertNullInt(t, null, "scanned null")
}

func TestIntValueOrZero(t *testing.T) {
	valid := NewInt(12345, true)
	if valid.ValueOrZero() != 12345 {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := NewInt(12345, false)
	if invalid.ValueOrZero() != 0 {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestIntEqual(t *testing.T) {
	int1 := NewInt(10, false)
	int2 := NewInt(10, false)
	assertIntEqualIsTrue(t, int1, int2)

	int1 = NewInt(10, false)
	int2 = NewInt(20, false)
	assertIntEqualIsTrue(t, int1, int2)

	int1 = NewInt(10, true)
	int2 = NewInt(10, true)
	assertIntEqualIsTrue(t, int1, int2)

	int1 = NewInt(10, true)
	int2 = NewInt(10, false)
	assertIntEqualIsFalse(t, int1, int2)

	int1 = NewInt(10, false)
	int2 = NewInt(10, true)
	assertIntEqualIsFalse(t, int1, int2)

	int1 = NewInt(10, true)
	int2 = NewInt(20, true)
	assertIntEqualIsFalse(t, int1, int2)
}

func assertInt(t *testing.T, i Int, from string) {
	if i.Int64 != 12345 {
		t.Errorf("bad %s int: %d ≠ %d\n", from, i.Int64, 12345)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullInt(t *testing.T, i Int, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertIntEqualIsTrue(t *testing.T, a, b Int) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf("Equal() of Int{%v, Valid:%t} and Int{%v, Valid:%t} should return true", a.Int64, a.Valid, b.Int64, b.Valid)
	}
}

func assertIntEqualIsFalse(t *testing.T, a, b Int) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf("Equal() of Int{%v, Valid:%t} and Int{%v, Valid:%t} should return false", a.Int64, a.Valid, b.Int64, b.Valid)
	}
}
