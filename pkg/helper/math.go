package helper

import (
	"golang.org/x/exp/constraints"
	"math"
)

// RoundUpDivide делит два числа и округляет результат вверх до ближайшего целого числа
func RoundUpDivide[T1, T2 constraints.Integer | constraints.Float](a T1, b T2) int {
	// Преобразуем входные значения в float64
	aFloat := float64(a)
	bFloat := float64(b)

	// Выполняем деление и округление вверх
	result := math.Ceil(aFloat / bFloat)

	// Преобразуем результат в int и возвращаем его
	return int(result)
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
