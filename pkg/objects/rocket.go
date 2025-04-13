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
	ActiveStage  int     // индекс текущей активной ступени
}

// RocketBody - основная часть спрайта ракеты (без нижней части)
var RocketBody = []string{
	"  /\\  ",
	" |==| ",
	" |  | ",
}

// GetRocketSprite возвращает полный спрайт ракеты в зависимости от текущей ступени
func (r *Rocket) GetRocketSprite() []string {
	// Проверка валидности индекса ступени
	if r.ActiveStage < 0 || r.ActiveStage >= len(RocketStages) {
		r.ActiveStage = 0
	}
	
	// Создаём полный спрайт, объединяя корпус и нижнюю часть текущей ступени
	fullSprite := make([]string, len(RocketBody)+len(RocketStages[r.ActiveStage].BottomSprite))
	
	// Копируем верхнюю часть (корпус)
	copy(fullSprite, RocketBody)
	
	// Копируем нижнюю часть (текущей ступени)
	copy(fullSprite[len(RocketBody):], RocketStages[r.ActiveStage].BottomSprite)
	
	return fullSprite
}

// Сохраняем традиционный спрайт для обратной совместимости
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
