package test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	jsoniter "github.com/json-iterator/go"
	"github.com/mailru/easyjson"
	"github.com/unravelin/null"
)

var fuzzFuncs = []interface{}{
	func(a *null.Bool, c fuzz.Continue) {
		a.Valid = c.RandBool()
		if a.Valid {
			a.Bool = c.RandBool()
		}
	},
	func(a *null.Float, c fuzz.Continue) {
		a.Valid = c.RandBool()
		if a.Valid {
			c.Fuzz(&a.Float64)
		}
	},
	func(a *null.Int, c fuzz.Continue) {
		a.Valid = c.RandBool()
		if a.Valid {
			c.Fuzz(&a.Int64)
		}
	},
	func(a *null.String, c fuzz.Continue) {
		a.Valid = c.RandBool()
		if a.Valid {
			c.Fuzz(&a.String)
		}
	},
	func(a *null.Time, c fuzz.Continue) {
		a.Valid = c.RandBool()
		if a.Valid {
			c.Fuzz(&a.Time)
		}
	},
}

func TestPlenc(t *testing.T) {
	f := fuzz.New().Funcs(fuzzFuncs...)
	for i := 0; i < 100; i++ {
		var in, out pltest
		f.Fuzz(&in)

		s := in.ΦλSize()
		b := make([]byte, 0, s)
		b = in.ΦλAppend(b)

		n, err := out.ΦλUnmarshal(b)
		if err != nil {
			t.Fatal(err)
		}
		if n != s || len(b) != s {
			t.Errorf("unexpected lengths %d %d %d", n, s, len(b))
		}

		if diff := cmp.Diff(in, out); diff != "" {
			t.Fatalf("values differ. %s", diff)
		}
	}
}

func TestEasyjson(t *testing.T) {
	f := fuzz.New().Funcs(fuzzFuncs...)
	for i := 0; i < 100; i++ {
		var in, out pltest
		f.Fuzz(&in)
		data, err := easyjson.Marshal(&in)
		if err != nil {
			t.Fatal(err)
		}

		if err := easyjson.Unmarshal(data, &out); err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(in, out); diff != "" {
			t.Fatalf("values differ. %s", diff)
		}
	}
}

func BenchmarkSerialisation(b *testing.B) {
	f := fuzz.New().Funcs(fuzzFuncs...)

	var in pltest
	f.Fuzz(&in)

	b.Run("plenc", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				s := in.ΦλSize()
				bu := make([]byte, 0, s)
				bu = in.ΦλAppend(bu)

				var out pltest
				out.ΦλUnmarshal(bu)
			}
		})
	})

	b.Run("easyjson", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				data, _ := easyjson.Marshal(&in)
				var out pltest
				easyjson.Unmarshal(data, &out)
			}
		})
	})

	b.Run("json", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				data, _ := json.Marshal(&in)
				var out pltest
				json.Unmarshal(data, &out)
			}
		})
	})

	b.Run("jsoniter", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				data, _ := jsoniter.Marshal(&in)
				var out pltest
				jsoniter.Unmarshal(data, &out)
			}
		})
	})

}
