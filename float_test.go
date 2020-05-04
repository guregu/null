package null

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
)

var (
	floatJSON           = []byte(`1.2345`)
	floatJSONString     = []byte(`"1.2345"`)
	nullFloatJSON       = []byte(`{"Float64":1.2345,"Valid":true}`)
	nullFloatJSONString = []byte(`{"Float64":"1.2345","Valid":true}`)
)

func TestFloatFrom(t *testing.T) {
	f := FloatFrom(1.2345)
	assertFloat(t, f, "FloatFrom()")

	zero := FloatFrom(0)
	if !zero.Valid {
		t.Error("FloatFrom(0)", "is invalid, but should be valid")
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
	tests := []struct {
		in             []byte
		exp            Float
		expErrType     reflect.Type
		expErrTypeEasy reflect.Type
	}{
		{
			in:  floatJSON,
			exp: FloatFrom(1.2345),
		},
		{
			in:  floatJSONString,
			exp: FloatFrom(1.2345),
		},
		{
			in: []byte(` "1.2345"  	 `),
			exp: FloatFrom(1.2345),
		},
		{
			in:  nullFloatJSON,
			exp: FloatFrom(1.2345),
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
	}

	for _, test := range tests {
		t.Run(string(test.in), func(t *testing.T) {
			var f Float
			err := json.Unmarshal(test.in, &f)
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
			if diff := cmp.Diff(test.exp, f); diff != "" {
				t.Fatalf("result not as expected. %s", diff)
			}
		})

		t.Run(string(test.in)+"_easyjson", func(t *testing.T) {
			var f Float
			err := easyjson.Unmarshal(test.in, &f)
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
			if diff := cmp.Diff(test.exp, f); diff != "" {
				t.Fatalf("result not as expected. %s", diff)
			}
		})
	}
}

func BenchmarkFloatUnmarshal(b *testing.B) {
	tests := [][]byte{
		floatJSON,
		floatJSONString,
		[]byte(` "1.2345"  	 `),
		nullFloatJSON,
		nullJSON,
	}

	for _, test := range tests {
		b.Run(string(test), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var ii Float
				if err := json.Unmarshal(test, &ii); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("easy "+string(test), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				w := &jlexer.Lexer{Data: test}
				var ii Float
				ii.UnmarshalEasyJSON(w)
				if err := w.Error(); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func TestTextUnmarshalFloat(t *testing.T) {
	var f Float
	err := f.UnmarshalText([]byte("1.2345"))
	maybePanic(err)
	assertFloat(t, f, "UnmarshalText() float")

	var blank Float
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullFloat(t, blank, "UnmarshalText() empty float")

	var null Float
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullFloat(t, null, `UnmarshalText() "null"`)
}

func TestMarshalFloat(t *testing.T) {
	f := FloatFrom(1.2345)
	data, err := json.Marshal(f)
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewFloat(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalFloatText(t *testing.T) {
	f := FloatFrom(1.2345)
	data, err := f.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewFloat(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestFloatPointer(t *testing.T) {
	f := FloatFrom(1.2345)
	ptr := f.Ptr()
	if *ptr != 1.2345 {
		t.Errorf("bad %s float: %#v ≠ %v\n", "pointer", ptr, 1.2345)
	}

	null := NewFloat(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s float: %#v ≠ %s\n", "nil pointer", ptr, "nil")
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
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
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
