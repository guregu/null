## null [![GoDoc](https://godoc.org/github.com/guregu/null?status.svg)](https://godoc.org/github.com/guregu/null) [![Coverage](http://gocover.io/_badge/github.com/guregu/null)](http://gocover.io/github.com/guregu/null)
`import "gopkg.in/guregu/null.v2"`

null is a library with reasonable options for dealing with nullable SQL and JSON values

There are two packages: `null` and its subpackage `zero`. 

Types in `null` will only be considered null on null input, and will JSON encode to `null`. If you need zero and null be considered separate values, use these.

Types in `zero` are treated like zero values in Go: blank string input will produce a null `zero.String`, and null Strings will JSON encode to `""`. If you need zero and null treated the same, use these.

All types implement `sql.Scanner` and `driver.Valuer`, so you can use this library in place of `sql.NullXXX`. All types also implement: `encoding.TextMarshaler`, `encoding.TextUnmarshaler`, `json.Marshaler`, and `json.Unmarshaler`. 

#### zero.String
A nullable string.

Will marshal to a blank string if null. Blank string input produces a null String. In other words, null values and empty values are considered equivalent. Can unmarshal from `sql.NullString` JSON input. 

#### zero.Int
A nullable int64.

Will marshal to 0 if null. Blank string or 0 input produces a null Int. In other words, null values and empty values are considered equivalent. Can unmarshal from `sql.NullInt64` JSON input. 

#### zero.Float
A nullable float64.

Will marshal to 0 if null. Blank string or 0 input produces a null Float. In other words, null values and empty values are considered equivalent. Can unmarshal from `sql.NullFloat64` JSON input. 

#### zero.Bool
A nullable bool.

Will marshal to false if null. Blank string or false input produces a null Float. In other words, null values and empty values are considered equivalent. Can unmarshal from `sql.NullBool` JSON input. 

#### null.String
An even nuller nullable string. 

Unlike `zero.String`, `null.String` will marshal to null if null. Zero (blank) input will not produce a null String. Can unmarshal from `sql.NullString` JSON input. 

#### null.Int
An even nuller nullable int64. 

Unlike `zero.Int`, `null.Int` will marshal to null if null. Zero input will not produce a null Int. Can unmarshal from `sql.NullInt64` JSON input. 

#### null.Float
An even nuller nullable float64. 

Unlike `zero.Float`, `null.Float` will marshal to null if null. Zero input will not produce a null Float. Can unmarshal from `sql.NullFloat64` JSON input. 

#### null.Bool
An even nuller nullable float64. 

Unlike `zero.Bool`, `null.Bool` will marshal to null if null. False input will not produce a null Bool. Can unmarshal from `sql.NullBool` JSON input. 

#### null.Time
An even nuller nullable time.Time. 

Unlike `null.Time` will marshal to null if null. Zero input will not produce a null time.Time. Can unmarshal from `pq.NullTime` JSON input. 

### Bugs
`json`'s `",omitempty"` struct tag does not work correctly right now. It will never omit a null or empty String. This should be [fixed eventually](https://github.com/golang/go/issues/4357).

### License
BSD