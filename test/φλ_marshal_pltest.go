package test

import (
	"github.com/philpearl/plenc"
)

// TODO: missing types
// slice of numeric ()
// slice of other
// pointers
// TODO: option whether top-level type is a pointer for marshaler

// ΦλSize works out how many bytes are needed to encode pltest
func (e *pltest) ΦλSize() (size int) {
	if e == nil {
		return 0
	}

	size += e.A.ΦλSizeFull(1)

	size += e.B.ΦλSizeFull(2)

	size += e.C.ΦλSizeFull(3)

	if s := e.D.ΦλSize(); s != 0 {
		size += plenc.SizeTag(plenc.WTLength, 4)
		size += plenc.SizeVarUint(uint64(s))
		size += s
	}

	if s := e.E.ΦλSize(); s != 0 {
		size += plenc.SizeTag(plenc.WTLength, 5)
		size += plenc.SizeVarUint(uint64(s))
		size += s
	}

	size += e.A1.ΦλSizeFull(6)

	size += e.B1.ΦλSizeFull(7)

	size += e.C1.ΦλSizeFull(8)

	if s := e.D1.ΦλSize(); s != 0 {
		size += plenc.SizeTag(plenc.WTLength, 9)
		size += plenc.SizeVarUint(uint64(s))
		size += s
	}

	if s := e.E1.ΦλSize(); s != 0 {
		size += plenc.SizeTag(plenc.WTLength, 10)
		size += plenc.SizeVarUint(uint64(s))
		size += s
	}

	size += e.A2.ΦλSizeFull(11)

	size += e.B2.ΦλSizeFull(12)

	size += e.C2.ΦλSizeFull(13)

	if s := e.D2.ΦλSize(); s != 0 {
		size += plenc.SizeTag(plenc.WTLength, 14)
		size += plenc.SizeVarUint(uint64(s))
		size += s
	}

	if s := e.E2.ΦλSize(); s != 0 {
		size += plenc.SizeTag(plenc.WTLength, 15)
		size += plenc.SizeVarUint(uint64(s))
		size += s
	}

	return size
}

// ΦλAppend encodes pltest by appending to data. It returns the final slice
func (e *pltest) ΦλAppend(data []byte) []byte {

	data = e.A.ΦλAppendFull(data, 1)

	data = e.B.ΦλAppendFull(data, 2)

	data = e.C.ΦλAppendFull(data, 3)

	if s := e.D.ΦλSize(); s != 0 {
		data = plenc.AppendTag(data, plenc.WTLength, 4)
		data = plenc.AppendVarUint(data, uint64(s))
		data = e.D.ΦλAppend(data)
	}

	if s := e.E.ΦλSize(); s != 0 {
		data = plenc.AppendTag(data, plenc.WTLength, 5)
		data = plenc.AppendVarUint(data, uint64(s))
		data = e.E.ΦλAppend(data)
	}

	data = e.A1.ΦλAppendFull(data, 6)

	data = e.B1.ΦλAppendFull(data, 7)

	data = e.C1.ΦλAppendFull(data, 8)

	if s := e.D1.ΦλSize(); s != 0 {
		data = plenc.AppendTag(data, plenc.WTLength, 9)
		data = plenc.AppendVarUint(data, uint64(s))
		data = e.D1.ΦλAppend(data)
	}

	if s := e.E1.ΦλSize(); s != 0 {
		data = plenc.AppendTag(data, plenc.WTLength, 10)
		data = plenc.AppendVarUint(data, uint64(s))
		data = e.E1.ΦλAppend(data)
	}

	data = e.A2.ΦλAppendFull(data, 11)

	data = e.B2.ΦλAppendFull(data, 12)

	data = e.C2.ΦλAppendFull(data, 13)

	if s := e.D2.ΦλSize(); s != 0 {
		data = plenc.AppendTag(data, plenc.WTLength, 14)
		data = plenc.AppendVarUint(data, uint64(s))
		data = e.D2.ΦλAppend(data)
	}

	if s := e.E2.ΦλSize(); s != 0 {
		data = plenc.AppendTag(data, plenc.WTLength, 15)
		data = plenc.AppendVarUint(data, uint64(s))
		data = e.E2.ΦλAppend(data)
	}

	return data
}
