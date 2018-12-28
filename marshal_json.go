package null

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
)

type IsZeroer interface {
	IsZero() bool
}

func MarshalJSON(s interface{}) ([]byte, error) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return nil, nil
	}

	buf := bytes.NewBufferString("{")

	// written fields count
	j := 0
	for i := 0; i < v.NumField(); i++ {
		valueField := v.Field(i)
		structField := v.Type().Field(i)

		fieldName := structField.Name
		allowSkip := false
		skipped := false
		if v, ok := structField.Tag.Lookup("json"); ok {
			parts := strings.Split(v, ",")
			fieldName = parts[0]
			if fieldName == "-" {
				skipped = true
			} else {
				for _, p := range parts {
					if p == "omitempty" {
						allowSkip = true
						continue
					}
				}
			}
		}

		if allowSkip {
			if v, ok := valueField.Interface().(IsZeroer); ok {
				if v.IsZero() {
					skipped = true
				}
			} else {
				// if value is zero-value
				if valueField.Interface() == reflect.Zero(valueField.Type()).Interface() {
					skipped = true
				}
			}
		}

		if !skipped {
			if i > 0 && j > 0 {
				buf.WriteString(",")
			}

			b, err := json.Marshal(valueField.Interface())
			if err != nil {
				return nil, err
			}
			buf.WriteString(`"` + fieldName + `":` + string(b))
			j++
		}
	}

	if _, err := buf.WriteString("}"); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
