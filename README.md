## null [![GoDoc](https://godoc.org/github.com/guregu/null?status.svg)](https://godoc.org/github.com/guregu/null) [![Coverage](http://gocover.io/_badge/github.com/guregu/null)](http://gocover.io/github.com/guregu/null)
null is a library with opinions on how to deal with nullable SQL and JSON values

There are two packages, `null`, and `nuller`. 

Types in `null` are treated like zero values in Go: blank string input will produce a null `null.String`, and null Strings will JSON encode to `""`. If you need zero and null treated the same, use these.

Types in `nuller` will only be considered null on null input, and will JSON encode to `null`. If you need zero and null be considered separate values, use these.

All types implement `sql.Scanner` and `driver.Valuer`, so you can use this library in place of `sql.NullXXX`. All types also implement: `encoding.TextMarshaler`, `encoding.TextUnmarshaler`, `json.Marshaler`, and `json.Unmarshaler`. 

#### null.String
A nullable string.

Will marshal to a blank string if null. Blank string input produces a null String. In other words, null values and empty values are considered equivalent. Can unmarshal from `sql.NullString` JSON input. 

#### null.Int
A nullable int64.

Will marshal to 0 if null. Blank string or 0 input produces a null Int. In other words, null values and empty values are considered equivalent. Can unmarshal from `sql.NullInt64` JSON input. 

#### nuller.String
An even nuller nullable string. 

Unlike `null.String`, `nuller.String` will marshal to null if null. Zero (blank) input will not produce a null String. Can unmarshal from `sql.NullString` JSON input. 

#### nuller.Int
An even nuller nullable int64. 

Unlike `null.Int`, `nuller.Int` will marshal to null if null. Zero input will not produce a null Int. Can unmarshal from `sql.NullInt64` JSON input. 

### Bugs
`json`'s `",omitempty"` struct tag does not work correctly right now. It will never omit a null or empty String. This should be [fixed in Go 1.4](https://code.google.com/p/go/issues/detail?id=4357).


### License
BSD