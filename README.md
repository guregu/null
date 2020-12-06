## null
`import "gopkg.in/guregu/null.v4"`

null is a library with reasonable options for dealing with nullable SQL and JSON values

There are two packages: `null` and its subpackage `zero`. 

Types in `null` will only be considered null on null input, and will JSON encode to `null`. If you need zero and null be considered separate values, use these.

Types in `zero` are treated like zero values in Go: blank string input will produce a null `zero.String`, and null Strings will JSON encode to `""`. Zero values of these types will be considered null to SQL. If you need zero and null treated the same, use these.

All types implement `sql.Scanner` and `driver.Valuer`, so you can use this library in place of `sql.NullXXX`.
All types also implement: `encoding.TextMarshaler`, `encoding.TextUnmarshaler`, `json.Marshaler`, and `json.Unmarshaler`. A null object's `MarshalText` will return a blank string.

### null package

`import "github.com/yeetpay/null"`

#### null.String
Nullable string.

Marshals to JSON null if SQL source data is null. Zero (blank) input will not produce a null String.

#### null.Int64
Nullable int64. 

Marshals to JSON null if SQL source data is null. Zero input will not produce a null Int63.

#### null.Int32
Nullable int32. 

Marshals to JSON null if SQL source data is null. Zero input will not produce a null Int32.


#### null.Float64
Nullable float64. 

Marshals to JSON null if SQL source data is null. Zero input will not produce a null Float.

#### null.Bool
Nullable bool. 

Marshals to JSON null if SQL source data is null. False input will not produce a null Bool.

#### null.Time

Marshals to JSON null if SQL source data is null. Zero input will not produce a null Time.

### zero package

`import "github.com/yeetpay/null/zero"`

#### zero.String
Nullable string.

Will marshal to a blank string if null. Blank string input produces a null String. Null values and zero values are considered equivalent.

#### zero.Int64
Nullable int64.

Will marshal to 0 if null. 0 produces a null Int64. Null values and zero values are considered equivalent. 

#### zero.Int32
Nullable int32.

Will marshal to 0 if null. 0 produces a null Int32. Null values and zero values are considered equivalent. 


#### zero.Float
Nullable float64.

Will marshal to 0.0 if null. 0.0 produces a null Float. Null values and zero values are considered equivalent. 

#### zero.Bool
Nullable bool.

Will marshal to false if null. `false` produces a null Float. Null values and zero values are considered equivalent.

#### zero.Time

Will marshal to the zero time if null. Uses `time.Time`'s marshaler.

### License
BSD
