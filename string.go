package null

import (
	"database/sql"
	"encoding/json"
)

type String struct {
	sql.NullString
}

func StringFrom(s string) String {
	return String{
		NullString: sql.NullString{
			String: s,
			Valid:  s != "",
		},
	}
}

func (s *String) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	json.Unmarshal(data, &v)
	switch v.(type) {
	case string:
		err = json.Unmarshal(data, &s.String)
	case map[string]interface{}:
		err = json.Unmarshal(data, &s.NullString)
	case nil:
		s.Valid = false
		return nil
	}
	s.Valid = err == nil
	return err
}

func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String)
}

func (s String) Pointer() *string {
	if s.String == "" {
		return nil
	}
	return &s.String
}

// IsZero returns true for invalid strings, for future omitempty support (Go 1.4?)
func (s String) IsZero() bool {
	return !s.Valid
}
