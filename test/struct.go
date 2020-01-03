package test

import "github.com/unravelin/null"

//go:generate easyjson -pkg -no_std_marshalers

//easyjson:json
type pltest struct {
	A  null.Bool   `plenc:"1"`
	B  null.Float  `plenc:"2"`
	C  null.Int    `plenc:"3"`
	D  null.String `plenc:"4"`
	E  null.Time   `plenc:"5"`
	A1 null.Bool   `plenc:"6"`
	B1 null.Float  `plenc:"7"`
	C1 null.Int    `plenc:"8"`
	D1 null.String `plenc:"9"`
	E1 null.Time   `plenc:"10"`
	A2 null.Bool   `plenc:"11"`
	B2 null.Float  `plenc:"12"`
	C2 null.Int    `plenc:"13"`
	D2 null.String `plenc:"14"`
	E2 null.Time   `plenc:"15"`
}
