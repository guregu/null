package zero

import (
	"encoding/json"
	"errors"
	"testing"
)

var (
	stringJSON      = []byte(`"test"`)
	blankStringJSON = []byte(`""`)
	nullStringJSON  = []byte(`{"String":"test","Valid":true}`)

	nullJSON    = []byte(`null`)
	invalidJSON = []byte(`:)`)
)

type stringInStruct struct {
	Test String `json:"test,omitempty"`
}

func TestStringFrom(t *testing.T) {
	str := StringFrom("test")
	assertStr(t, str, "StringFrom() string")

	null := StringFrom("")
	assertNullStr(t, null, "StringFrom() empty string")
}

func TestUnmarshalString(t *testing.T) {
	var str String
	err := json.Unmarshal(stringJSON, &str)
	maybePanic(err)
	assertStr(t, str, "string json")

	var ns String
	err = json.Unmarshal(nullStringJSON, &ns)
	if err == nil {
		panic("expected error")
	}

	var blank String
	err = json.Unmarshal(blankStringJSON, &blank)
	maybePanic(err)
	assertNullStr(t, blank, "blank string json")

	var null String
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullStr(t, null, "null json")

	var badType String
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullStr(t, badType, "wrong type json")

	var invalid String
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
	assertNullStr(t, invalid, "invalid json")
}

func TestTextUnmarshalString(t *testing.T) {
	var str String
	err := str.UnmarshalText([]byte("test"))
	maybePanic(err)
	assertStr(t, str, "UnmarshalText() string")

	var null String
	err = null.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullStr(t, null, "UnmarshalText() empty string")
}

func TestMarshalString(t *testing.T) {
	str := StringFrom("test")
	data, err := json.Marshal(str)
	maybePanic(err)
	assertJSONEquals(t, data, `"test"`, "non-empty json marshal")

	// invalid values should be encoded as an empty string
	null := StringFrom("")
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, `""`, "empty json marshal")
}

// Tests omitempty... broken until Go 1.4
// func TestMarshalStringInStruct(t *testing.T) {
// 	obj := stringInStruct{Test: StringFrom("")}
// 	data, err := json.Marshal(obj)
// 	maybePanic(err)
// 	assertJSONEquals(t, data, `{}`, "null string in struct")
// }

func TestStringPointer(t *testing.T) {
	str := StringFrom("test")
	ptr := str.Ptr()
	if *ptr != "test" {
		t.Errorf("bad %s string: %#v ≠ %s\n", "pointer", ptr, "test")
	}

	null := StringFrom("")
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s string: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestStringFromPointer(t *testing.T) {
	test := "test"
	testptr := &test
	str := StringFromPtr(testptr)
	assertStr(t, str, "StringFromPtr()")

	testptr = nil
	null := StringFromPtr(testptr)
	assertNullStr(t, null, "StringFromPtr()")

	ptr := null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s string: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestStringIsZero(t *testing.T) {
	str := StringFrom("test")
	if str.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := StringFrom("")
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	empty := NewString("", true)
	if !empty.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestStringScan(t *testing.T) {
	var str String
	err := str.Scan("test")
	maybePanic(err)
	assertStr(t, str, "scanned string")

	var null String
	err = null.Scan(nil)
	maybePanic(err)
	assertNullStr(t, null, "scanned null")
}

func TestStringSetValid(t *testing.T) {
	change := NewString("", false)
	assertNullStr(t, change, "SetValid()")
	change.SetValid("test")
	assertStr(t, change, "SetValid()")
}

func TestStringValueOrZero(t *testing.T) {
	valid := NewString("test", true)
	if valid.ValueOrZero() != "test" {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := NewString("test", false)
	if invalid.ValueOrZero() != "" {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestStringEqual(t *testing.T) {
	str1 := NewString("foo", false)
	str2 := NewString("foo", false)
	assertStringEqualIsTrue(t, str1, str2)

	str1 = NewString("foo", false)
	str2 = NewString("bar", false)
	assertStringEqualIsTrue(t, str1, str2)

	str1 = NewString("foo", true)
	str2 = NewString("foo", true)
	assertStringEqualIsTrue(t, str1, str2)

	str1 = NewString("", true)
	str2 = NewString("foo", false)
	assertStringEqualIsTrue(t, str1, str2)

	str1 = NewString("foo", true)
	str2 = NewString("foo", false)
	assertStringEqualIsFalse(t, str1, str2)

	str1 = NewString("foo", false)
	str2 = NewString("foo", true)
	assertStringEqualIsFalse(t, str1, str2)

	str1 = NewString("foo", true)
	str2 = NewString("bar", true)
	assertStringEqualIsFalse(t, str1, str2)
}

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func assertStr(t *testing.T, s String, from string) {
	if s.String != "test" {
		t.Errorf("bad %s string: %s ≠ %s\n", from, s.String, "test")
	}
	if !s.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullStr(t *testing.T, s String, from string) {
	if s.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertJSONEquals(t *testing.T, data []byte, cmp string, from string) {
	if string(data) != cmp {
		t.Errorf("bad %s data: %s ≠ %s\n", from, data, cmp)
	}
}

func assertStringEqualIsTrue(t *testing.T, a, b String) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf("Equal() of String{\"%v\", Valid:%t} and String{\"%v\", Valid:%t} should return true", a.String, a.Valid, b.String, b.Valid)
	}
}

func assertStringEqualIsFalse(t *testing.T, a, b String) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf("Equal() of String{\"%v\", Valid:%t} and String{\"%v\", Valid:%t} should return false", a.String, a.Valid, b.String, b.Valid)
	}
}
