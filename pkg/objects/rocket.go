package objects

// Rocket описывает состояние ракеты
type Rocket struct {
	X, Y    int     // позиция (левый верхний угол спрайта)
	Vx, Vy  float64 // скорости по осям X и Y
	ThrustX float64 // тяга по горизонтали (положительное значение – вправо)
	ThrustY float64 // тяга по вертикали (для подъёма; базовая равна HoverThrust)
	Fuel    float64 // оставшееся топливо
}

var RocketSprite = []string{
	"  /\\  ",
	" |==| ",
	" |  | ",
	"  ||  ",
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
