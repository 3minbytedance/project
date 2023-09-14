package redis

import (
	"strconv"
	"strings"
	"time"
)

func SetMonthlyActiveBit(userId uint) {
	month := time.Now().Month().String()
	baseSlice := []string{MonthlyActive, month, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	day := time.Now().Day()
	Rdb.SetBit(Ctx, key, int64(day), 1)
	return
}

func CountMonthlyActiveBit(month time.Month, userId uint) int64 {
	monthStr := month.String()
	baseSlice := []string{MonthlyActive, monthStr, strconv.Itoa(int(userId))}
	key := strings.Join(baseSlice, Delimiter)
	val := Rdb.BitCount(Ctx, key, nil).Val()
	return val
}
