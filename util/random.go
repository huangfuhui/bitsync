package util

import "math/rand"

type Random struct {
}

// 生成64位int随机数
func (r *Random) Rand(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return 0
	}

	return rand.Int63n(max-min) + min
}
