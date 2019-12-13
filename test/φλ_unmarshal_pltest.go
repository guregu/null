package test

import (
	"fmt"
	"time"

	"github.com/philpearl/plenc"
)

var _ time.Time

func (e *pltest) ΦλUnmarshal(data []byte) (int, error) {

	var offset int
	for offset < len(data) {
		wt, index, n := plenc.ReadTag(data[offset:])
		if n == 0 {
			break
		}
		offset += n
		switch index {

		case 1:

			n, err := e.A.ΦλUnmarshal(data[offset:])
			if err != nil {
				return 0, fmt.Errorf("failed to unmarshal field %d A (Bool). %w", index, err)
			}

			offset += n

		case 2:

			n, err := e.B.ΦλUnmarshal(data[offset:])
			if err != nil {
				return 0, fmt.Errorf("failed to unmarshal field %d B (Float). %w", index, err)
			}

			offset += n

		case 3:

			n, err := e.C.ΦλUnmarshal(data[offset:])
			if err != nil {
				return 0, fmt.Errorf("failed to unmarshal field %d C (Int). %w", index, err)
			}

			offset += n

		case 4:

			s, n := plenc.ReadVarUint(data[offset:])
			offset += n
			n, err := e.D.ΦλUnmarshal(data[offset : offset+int(s)])
			if err != nil {
				return 0, fmt.Errorf("failed to unmarshal field %d D (String). %w", index, err)
			}

			offset += n

		case 5:

			s, n := plenc.ReadVarUint(data[offset:])
			offset += n
			n, err := e.E.ΦλUnmarshal(data[offset : offset+int(s)])
			if err != nil {
				return 0, fmt.Errorf("failed to unmarshal field %d E (Time). %w", index, err)
			}

			offset += n

		case 6:

			n, err := e.A1.ΦλUnmarshal(data[offset:])
			if err != nil {
				return 0, fmt.Errorf("failed to unmarshal field %d A1 (Bool). %w", index, err)
			}

			offset += n

		case 7:

			n, err := e.B1.ΦλUnmarshal(data[offset:])
			if err != nil {
				return 0, fmt.Errorf("failed to unmarshal field %d B1 (Float). %w", index, err)
			}

			offset += n

		case 8:

			n, err := e.C1.ΦλUnmarshal(data[offset:])
			if err != nil {
				return 0, fmt.Errorf("failed to unmarshal field %d C1 (Int). %w", index, err)
			}

			offset += n

		case 9:

			s, n := plenc.ReadVarUint(data[offset:])
			offset += n
			n, err := e.D1.ΦλUnmarshal(data[offset : offset+int(s)])
			if err != nil {
				return 0, fmt.Errorf("failed to unmarshal field %d D1 (String). %w", index, err)
			}

			offset += n

		case 10:

			s, n := plenc.ReadVarUint(data[offset:])
			offset += n
			n, err := e.E1.ΦλUnmarshal(data[offset : offset+int(s)])
			if err != nil {
				return 0, fmt.Errorf("failed to unmarshal field %d E1 (Time). %w", index, err)
			}

			offset += n

		case 11:

			n, err := e.A2.ΦλUnmarshal(data[offset:])
			if err != nil {
				return 0, fmt.Errorf("failed to unmarshal field %d A2 (Bool). %w", index, err)
			}

			offset += n

		case 12:

			n, err := e.B2.ΦλUnmarshal(data[offset:])
			if err != nil {
				return 0, fmt.Errorf("failed to unmarshal field %d B2 (Float). %w", index, err)
			}

			offset += n

		case 13:

			n, err := e.C2.ΦλUnmarshal(data[offset:])
			if err != nil {
				return 0, fmt.Errorf("failed to unmarshal field %d C2 (Int). %w", index, err)
			}

			offset += n

		case 14:

			s, n := plenc.ReadVarUint(data[offset:])
			offset += n
			n, err := e.D2.ΦλUnmarshal(data[offset : offset+int(s)])
			if err != nil {
				return 0, fmt.Errorf("failed to unmarshal field %d D2 (String). %w", index, err)
			}

			offset += n

		case 15:

			s, n := plenc.ReadVarUint(data[offset:])
			offset += n
			n, err := e.E2.ΦλUnmarshal(data[offset : offset+int(s)])
			if err != nil {
				return 0, fmt.Errorf("failed to unmarshal field %d E2 (Time). %w", index, err)
			}

			offset += n

		default:
			// Field corresponding to index does not exist
			n, err := plenc.Skip(data[offset:], wt)
			if err != nil {
				return 0, fmt.Errorf("failed to skip field %d. %w", index, err)
			}
			offset += n
		}
	}

	return offset, nil
}
