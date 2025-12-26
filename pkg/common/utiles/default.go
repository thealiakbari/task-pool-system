package utiles

import "golang.org/x/exp/constraints"

func ZeroDefault[T constraints.Integer | constraints.Float](val T, defaultValue T) T {
	if val == 0 {
		return defaultValue
	}
	return val
}
