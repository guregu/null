## null-extended [![GoDoc](https://godoc.org/github.com/nullbio/null?status.svg)](https://godoc.org/github.com/nullbio/null) [![Coverage](http://gocover.io/_badge/github.com/nullbio/null)](http://gocover.io/github.com/nullbio/null)

Install:

`go get -u "gopkg.in/nullbio/null.v6"`

---

`import "gopkg.in/nullbio/null.v6"`

null-extended is a library with reasonable options for dealing with nullable SQL and JSON values

Types in `null` will only be considered null on null input, and will JSON encode to `null`.

All types implement `sql.Scanner` and `driver.Valuer`, so you can use this library in place of `sql.NullXXX`. All types also implement: `encoding.TextMarshaler`, `encoding.TextUnmarshaler`, `json.Marshaler`, `json.Unmarshaler` and `sql.Scanner`.

### null package

`import "gopkg.in/nullbio/null.v6"`

#### null.JSON
Nullable []byte.

Will marshal to JSON null if Invalid. []byte{} input will not produce an Invalid JSON, but []byte(nil) will. This should be used for storing raw JSON in the database.

Also has `null.JSON.Marshal` and `null.JSON.Unmarshal` helpers to marshal and unmarshal foreign objects.

#### null.Bytes
Nullable []byte.

Will marshal to JSON null if Invalid. []]byte{} input will not produce an Invalid Bytes, but []byte(nil) will. This should be used for storing binary data (bytea in PSQL for example) in the database.

#### null.String
Nullable string.

Marshals to JSON null if SQL source data is null. Zero (blank) input will not produce a null String. Can unmarshal from `sql.NullString` JSON input or string input.

#### null.Bool
Nullable bool.

Marshals to JSON null if SQL source data is null. False input will not produce a null Bool. Can unmarshal from `sql.NullBool` JSON input.

#### null.Time

Marshals to JSON null if SQL source data is null. Uses `time.Time`'s marshaler. Can unmarshal from `pq.NullTime` and similar JSON input.

#### null.Float32
Nullable float32.

Marshals to JSON null if SQL source data is null. Zero input will not produce a null Float32. Can unmarshal from `null.NullFloat32` JSON input.

#### null.Float64
Nullable float64.

Marshals to JSON null if SQL source data is null. Zero input will not produce a null Float64. Can unmarshal from `sql.NullFloat64` JSON input.

#### null.Int
Nullable int.

Marshals to JSON null if SQL source data is null. Zero input will not produce a null Int. Can unmarshal from `null.NullInt` JSON input.

#### null.Int8
Nullable int8.

Marshals to JSON null if SQL source data is null. Zero input will not produce a null Int8. Can unmarshal from `null.NullInt8` JSON input.

#### null.Int16
Nullable int16.

Marshals to JSON null if SQL source data is null. Zero input will not produce a null Int16. Can unmarshal from `null.NullInt16` JSON input.

#### null.Int32
Nullable int32.

Marshals to JSON null if SQL source data is null. Zero input will not produce a null Int32. Can unmarshal from `null.NullInt32` JSON input.

#### null.Int64
Nullable int64.

Marshals to JSON null if SQL source data is null. Zero input will not produce a null Int64. Can unmarshal from `sql.NullInt64` JSON input.

#### null.Uint
Nullable uint.

Marshals to JSON null if SQL source data is null. Zero input will not produce a null Uint. Can unmarshal from `null.NullUint` JSON input.

#### null.Uint8
Nullable uint8.

Marshals to JSON null if SQL source data is null. Zero input will not produce a null Uint8. Can unmarshal from `null.NullUint8` JSON input.

#### null.Uint16
Nullable uint16.

Marshals to JSON null if SQL source data is null. Zero input will not produce a null Uint16. Can unmarshal from `null.NullUint16` JSON input.

#### null.Uint32
Nullable int32.

Marshals to JSON null if SQL source data is null. Zero input will not produce a null Uint32. Can unmarshal from `null.NullUint32` JSON input.

#### null.Int64
Nullable uint64.

Marshals to JSON null if SQL source data is null. Zero input will not produce a null Uint64. Can unmarshal from `null.NullUint64` JSON input.

### Bugs
`json`'s `",omitempty"` struct tag does not work correctly right now. It will never omit a null or empty String. This might be [fixed eventually](https://github.com/golang/go/issues/4357).

### License
BSD
