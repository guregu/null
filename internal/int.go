package internal

import (
	"encoding/json"
	"fmt"
)

type Integer interface {
	int64 | int32 | int16 | byte
}

func UnmarshalIntJSON[T Integer, U int64 | uint64](data []byte, value *T, valid *bool, bits int, parse func(string, int, int) (U, error)) error {
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
		n, err := parse(str, 10, bits)
		if err != nil {
			return fmt.Errorf("null: couldn't convert string to int: %w", err)
		}
		*value = T(n)
		*valid = true
		return nil

	default:
		err := json.Unmarshal(data, value)
		*valid = err == nil
		return err
	}
}

func UnmarshalIntText[T Integer, U int64 | uint64](text []byte, value *T, valid *bool, bits int, parse func(string, int, int) (U, error)) error {
	str := string(text)
	if str == "" || str == "null" {
		*value = 0
		*valid = false
		return nil
	}
	n, err := parse(str, 10, bits)
	*value = T(n)
	if err != nil {
		*valid = false
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}
	*valid = true
	return nil
}
