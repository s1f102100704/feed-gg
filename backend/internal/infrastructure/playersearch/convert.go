package playersearch

import (
	"database/sql"
	"fmt"
	"math"
	"time"
)

func nowUTC() time.Time {
	return time.Now().UTC()
}

func nullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}

func nullTimeFromMillis(value int64) sql.NullTime {
	if value <= 0 {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: timeFromMillis(value), Valid: true}
}

func timeFromMillis(value int64) time.Time {
	if value <= 0 {
		return time.Unix(0, 0).UTC()
	}
	return time.UnixMilli(value).UTC()
}

func toNullInt32(value int) (sql.NullInt32, error) {
	int32Value, err := toInt32(value)
	if err != nil {
		return sql.NullInt32{}, err
	}
	return sql.NullInt32{Int32: int32Value, Valid: true}, nil
}

func toInt32[T ~int | ~int64](value T) (int32, error) {
	if value > math.MaxInt32 || value < math.MinInt32 {
		return 0, fmt.Errorf("value %d is out of int32 range", value)
	}
	return int32(value), nil
}
