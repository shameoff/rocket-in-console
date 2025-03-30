package objects

import "math/rand"

type Star struct {
	X, Y int
}

var Stars []Star

// Параметры мира (при необходимости можно вынести в отдельный конфиг)
var WorldWidth, WorldHeight int = 10000, 20000
var GroundLevel int = 1

// InitStars генерирует массив звёзд
func InitStars(n int) {
	Stars = make([]Star, n)
	for i := 0; i < n; i++ {
		Stars[i] = Star{
			X: rand.Intn(WorldWidth),
			Y: rand.Intn(GroundLevel),
		}
	}
}

// IsStarAt возвращает true, если в мировых координатах (x, y) должна быть звезда.
func IsStarAt(x, y int) bool {
	h := int64(x)*73856093 ^ int64(y)*19349663
	if h < 0 {
		h = -h
	}
	return h%100 < 3
}
