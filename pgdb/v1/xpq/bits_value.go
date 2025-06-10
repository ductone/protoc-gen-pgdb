package xpq

import (
	"database/sql/driver"

	"github.com/jackc/pgx/v5/pgtype"
)

type BitsValue []byte

func (b BitsValue) Value() (driver.Value, error) {
	bits := &pgtype.Bits{
		Bytes: b,
		Valid: true,
		Len:   int32(len(b)),
	}

	return bits.Value()
}
