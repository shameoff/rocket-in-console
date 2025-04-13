package objects

// Rocket описывает состояние ракеты
type Rocket struct {
	X, Y         int     // позиция (левый верхний угол спрайта)
	Vx, Vy       float64 // скорости по осям X и Y
	ThrustX      float64 // тяга по горизонтали (положительное значение – вправо)
	ThrustY      float64 // тяга по вертикали (для подъёма; базовая равна HoverThrust)
	Fuel         float64 // оставшееся топливо
	AccumulatedX float64 // аккумулятор дробных перемещений по X
	AccumulatedY float64 // аккумулятор дробных перемещений по Y
}

var RocketSprite = []string{
	"  /\\  ",
	" |==| ",
	" |  | ",
	"  /\\  ",
}

var ExplosionSprite = []string{
	"   ***   ",
	"  *****  ",
	" ******* ",
	"*********",
	" ******* ",
	"  *****  ",
	"   ***   ",
}
