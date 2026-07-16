package pkg

import "math"

const MoneyScale int64 = 100

// ToMoney 金额转最小单位
// 10.25 -> 1025
func ToMoney(amount float64) int64 {
	return int64(math.Floor(amount * float64(MoneyScale)))
}

// ToAmount 最小单位转金额
// 1025 -> 10.25
func ToAmount(money int64) float64 {
	return float64(money) / float64(MoneyScale)
}

// Percent 百分比
// Percent(1000, 1) = 10
func Percent(value int64, rate int64) int64 {
	return value * rate / 100
}

// Permillage 千分比
// Permillage(1000, 1) = 1
func Permillage(value int64, rate int64) int64 {
	return value * rate / 1000
}