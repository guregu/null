package null

import (
	"encoding/json"
	"testing"
)

var (
	stringJSON      = []byte(`"test"`)
	blankStringJSON = []byte(`""`)
	nullStringJSON  = []byte(`{"String":"test","Valid":true}`)
	nullJSON        = []byte(`null`)
)

type stringInStruct struct {
	Test String `json:"test,omitempty"`
}

func TestStringFrom(t *testing.T) {
	str := StringFrom("test")
	assert(t, str, "StringFrom() string")

	null := StringFrom("")
	assertNull(t, null, "StringFrom() empty string")
}

func TestUnmarshalString(t *testing.T) {
	var str String
	err := json.Unmarshal(stringJSON, &str)
	maybePanic(err)
	assert(t, str, "string json")

	var ns String
	err = json.Unmarshal(nullStringJSON, &ns)
	maybePanic(err)
	assert(t, ns, "null string object json")

	var blank String
	err = json.Unmarshal(blankStringJSON, &blank)
	maybePanic(err)
	assertNull(t, blank, "blank string json")

	var null String
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNull(t, null, "null json")
}

func TestTextUnmarshalString(t *testing.T) {
	var str String
	err := str.UnmarshalText([]byte("test"))
	maybePanic(err)
	assert(t, str, "TextUnmarshal() string")

	var null String
	err = null.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNull(t, null, "TextUnmarshal() empty string")
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

func TestPointer(t *testing.T) {
	str := StringFrom("test")
	ptr := str.Pointer()
	if *ptr != "test" {
		t.Errorf("bad %s string: %#v ≠ %s\n", "pointer", ptr, "test")
	}

	null := StringFrom("")
	ptr = null.Pointer()
	if ptr != nil {
		t.Errorf("bad %s: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestIsZero(t *testing.T) {
	str := StringFrom("test")
	if str.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := StringFrom("")
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestScan(t *testing.T) {
	var str String
	err := str.Scan("test")
	maybePanic(err)
	assert(t, str, "scanned string")

	var null String
	err = null.Scan(nil)
	maybePanic(err)
	assertNull(t, null, "scanned null")
}

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func assert(t *testing.T, s String, from string) {
	if s.String != "test" {
		t.Errorf("bad %s string: %s ≠ %s\n", from, s.String, "test")
	}
	if !s.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNull(t *testing.T, s String, from string) {
	if s.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertJSONEquals(t *testing.T, data []byte, cmp string, from string) {
	if string(data) != cmp {
		t.Errorf("bad %s data: %s ≠ %s\n", from, data, cmp)
	}
}