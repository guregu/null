module github.com/unravelin/null/testdata/bench

go 1.17

replace github.com/unravelin/null => ../../

require (
	github.com/google/go-cmp v0.5.6
	github.com/google/gofuzz v1.2.0
	github.com/json-iterator/go v1.1.12
	github.com/mailru/easyjson v0.7.7
	github.com/philpearl/plenc v0.0.1
	github.com/unravelin/null v2.1.4+incompatible
)

require (
	github.com/josharian/intern v1.0.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
)
