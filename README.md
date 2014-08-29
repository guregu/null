## null [![GoDoc](https://godoc.org/github.com/guregu/null?status.svg)](https://godoc.org/github.com/guregu/null) [![Coverage](http://gocover.io/_badge/github.com/guregu/null)](http://gocover.io/github.com/guregu/null)
null is a library with opinions on how to deal with nullable SQL and JSON values

### String
A nullable string. Implements `sql.Scanner`, `encoding.Marshaler` and `encoding.TextUnmarshaler`, providing support for JSON and XML. 

Will marshal to a blank string if null. Blank string input produces a null String. In other words, null values and empty values are considered equivalent.

`UnmarshalJSON` supports `sql.NullString` input. 

### Bugs
`json`'s `",omitempty"` struct tag does not work correctly right now. It will never omit a null or empty String. This should be [fixed in Go 1.4](https://code.google.com/p/go/issues/detail?id=4357).


### License
BSD