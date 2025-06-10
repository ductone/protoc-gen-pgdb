package xpq

import (
	"database/sql/driver"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgtype"
)

type BitsValue []byte

func (b BitsValue) Value() (driver.Value, error) {
	if len(b) > math.MaxInt32 {
		return nil, fmt.Errorf("bits value too large: %d bytes", len(b))
	}

	bits := &pgtype.Bits{
		Bytes: b,
		Valid: true,
		Len:   int32(len(b)),
	}

	return bits.Value()
}
