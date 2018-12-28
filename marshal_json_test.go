package null

import (
	"encoding/json"
	"testing"
)

type testStruct struct {
	Id01 int32  `json:"id01"`
	Id02 int32  `json:"id02"`
	Id03 Int    `json:"id03"`
	Id04 Int    `json:"id04"`
	Id05 Int    `json:"id05"`
	Id06 Int    `json:"id06"`
	Id07 Int    `json:"-"`
	Id08 Int    `json:"id08,omitempty"`
	Id09 Int    `json:"id09,omitempty"`
	Id10 string `json:"id10,omitempty"`
	Id11 string `json:"-"`
}

func (ts testStruct) MarshalJSON() ([]byte, error) {
	return MarshalJSON(ts)
}

func TestMarshalJSON1(t *testing.T) {
	var ts testStruct
	ts.Id01 = 0
	ts.Id02 = 1
	ts.Id03 = NewInt(0, false)
	ts.Id04 = NewInt(0, true)
	ts.Id05 = NewInt(123, false)
	ts.Id06 = NewInt(123, true)
	ts.Id07 = NewInt(123, true)
	ts.Id09 = NewInt(123, true)
	ts.Id11 = "test"

	data, err := json.Marshal(ts)
	maybePanic(err)
	assertJSONEquals(t, data, `{"id01":0,"id02":1,"id03":null,"id04":0,"id05":null,"id06":123,"id09":123}`, "struct json marshal")
}

type testString string

func (ts testString) MarshalJSON() ([]byte, error) {
	return MarshalJSON(ts)
}

func TestMarshalJSON2(t *testing.T) {
	var ts testString
	data, err := json.Marshal(ts)
	assertErr(t, err, "TestMarshalJSON2 with string")
	assertJSONEquals(t, data, "", "string json marshal")
}

func assertErr(t *testing.T, err error, from string) {
	if err == nil {
		t.Error(from, "doesn't have an error, but should")
	}
}
