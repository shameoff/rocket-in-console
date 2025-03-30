package objects

import "math/rand"

type Cloud struct {
	X, Y   int
	Sprite []string
}

var CloudSprite = []string{
	"  ~~  ",
	"~~~~~~",
	"  ~~  ",
}

var Clouds []Cloud

// InitClouds генерирует n облаков в заданной зоне по оси Y (например, от 10 до 30)
func InitClouds(n int) {
	Clouds = make([]Cloud, n)
	for i := 0; i < n; i++ {
		Clouds[i] = Cloud{
			X:      rand.Intn(WorldWidth),
			Y:      10 + rand.Intn(20),
			Sprite: CloudSprite,
		}
	}
}
