package null

import (
	"encoding/json"
	"testing"

	"github.com/jackc/pgtype"
)

var (
	intervalJSONString = []byte(`"1 year 2 mons 1 day 03:02:01.123456"`)
	intervalJSON       = []byte(`{"Months": 14, "Days": 1, "Microseconds": 10921123456, "Status": 2}`)
	intervalNullJSON   = []byte("null")
)

func TestUnmarshalIntervalFromString(t *testing.T) {
	var i Interval
	err := json.Unmarshal(intervalJSONString, &i)
	maybePanic(err)
	if i.Status != pgtype.Present {
		t.Errorf("bad %s interval status: %d ≠ %d\n", intervalJSONString, i.Status, pgtype.Present)
	}
	if i.Months != 14 {
		t.Errorf("bad %s interval months: %d ≠ %d\n", intervalJSONString, i.Months, 14)
	}
	if i.Days != 1 {
		t.Errorf("bad %s interval days: %d ≠ %d\n", intervalJSONString, i.Days, 1)
	}
	if i.Microseconds != 10921123456 {
		t.Errorf("bad %s interval microseconds: %d ≠ %d\n", intervalJSONString, i.Microseconds, 10921123456)
	}

	var i2 Interval
	err = json.Unmarshal(intervalNullJSON, &i2)
	maybePanic(err)
	if i2.Status != pgtype.Null {
		t.Errorf("bad %s interval status: %d ≠ %d\n", intervalNullJSON, i2.Status, pgtype.Null)
	}
}

func TestUnmarshalIntervalFromObject(t *testing.T) {
	var i Interval
	err := json.Unmarshal(intervalJSON, &i)
	maybePanic(err)
	if i.Status != pgtype.Present {
		t.Errorf("bad %s interval status: %d ≠ %d\n", intervalJSONString, i.Status, pgtype.Present)
	}
	if i.Months != 14 {
		t.Errorf("bad %s interval months: %d ≠ %d\n", intervalJSONString, i.Months, 14)
	}
	if i.Days != 1 {
		t.Errorf("bad %s interval days: %d ≠ %d\n", intervalJSONString, i.Days, 1)
	}
	if i.Microseconds != 10921123456 {
		t.Errorf("bad %s interval microseconds: %d ≠ %d\n", intervalJSONString, i.Microseconds, 10921123456)
	}
}

func TestMarshalInterval(t *testing.T) {
	i := Interval{
		Interval: pgtype.Interval{
			Microseconds: 10921123456,
			Days:         1,
			Months:       14,
			Status:       pgtype.Present,
		},
	}

	expected := `"14 mon 1 day 03:02:01.123456"`
	res, err := json.Marshal(i)
	maybePanic(err)
	if string(res) != expected {
		t.Errorf("bad %s interval: %s ≠ %s\n", intervalJSONString, string(res), expected)
	}

	i = Interval{}
	expected = "null"
	res, err = json.Marshal(i)
	maybePanic(err)
	if string(res) != expected {
		t.Errorf("bad %s interval: %s ≠ %s\n", intervalJSONString, string(res), expected)
	}

	i = Interval{Interval: pgtype.Interval{Status: pgtype.Undefined}}
	expected = "null"
	res, err = json.Marshal(i)
	maybePanic(err)
	if string(res) != expected {
		t.Errorf("bad %s interval: %s ≠ %s\n", intervalJSONString, string(res), expected)
	}
}
