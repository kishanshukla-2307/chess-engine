package utils

import "math/rand/v2"

func Abs[T int16 | int32 | int64 | int](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

var ZORBIST_TABLE [64][12]uint64
var BLACK_TO__MOVE uint64

func initialize_zobrist() {
	for i := 0; i < 64; i++ {
		for j := 0; j < 12; j++ {
			ZORBIST_TABLE[i][j] = rand.Uint64()
		}
	}
	BLACK_TO__MOVE = rand.Uint64()
}
