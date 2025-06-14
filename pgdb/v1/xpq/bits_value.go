package xpq

import (
	"database/sql/driver"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgtype"
)

type BitsValue []byte

func (b BitsValue) Value() (driver.Value, error) {
	bitsLen, err := safeIntToInt32(len(b) * 8)
	if err != nil {
		return nil, err
	}

	bits := &pgtype.Bits{
		Bytes: b,
		Valid: true,
		Len:   bitsLen,
	}

	return bits.Value()
}

func safeIntToInt32(val int) (int32, error) {
	if val < math.MinInt32 || val > math.MaxInt32 {
		return 0, fmt.Errorf("value %d out of range for int32", val)
	}
	return int32(val), nil
}
