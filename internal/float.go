package internal

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func UnmarshalFloatJSON(data []byte, value *float64, valid *bool) error {
	if len(data) == 0 {
		return fmt.Errorf("UnmarshalJSON: no data")
	}

	switch data[0] {
	case 'n':
		*value = 0
		*valid = false
		return nil

	case '"':
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return fmt.Errorf("null: couldn't unmarshal number string: %w", err)
		}
		n, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return fmt.Errorf("null: couldn't convert string to int: %w", err)
		}
		*value = n
		*valid = true
		return nil

	default:
		err := json.Unmarshal(data, value)
		*valid = err == nil
		return err
	}
}
