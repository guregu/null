package internal

import "fmt"

func TypeName[T any]() string {
	return fmt.Sprintf("%T", *(new(T)))
}
