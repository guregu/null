package test

import "github.com/unravelin/null"

//go:generate plencgen -pkg . -type pltest
//go:generate easyjson -pkg -no_std_marshalers

//easyjson:json
type pltest struct {
	A  null.Bool
	B  null.Float
	C  null.Int
	D  null.String
	E  null.Time
	A1 null.Bool
	B1 null.Float
	C1 null.Int
	D1 null.String
	E1 null.Time
	A2 null.Bool
	B2 null.Float
	C2 null.Int
	D2 null.String
	E2 null.Time
}
