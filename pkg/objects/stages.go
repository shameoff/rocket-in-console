package objects

// Stage описывает характеристики ступени ракеты
type Stage struct {
	Name              string   // название ступени
	MaxThrustX        float64  // максимальная горизонтальная тяга
	MaxThrustY        float64  // максимальная вертикальная тяга
	BottomSprite      []string // нижняя часть спрайта (сопла)
	FuelConsumptionRate float64 // скорость потребления топлива
}

// Предустановленные ступени ракеты
var RocketStages = []Stage{
	{
		Name:              "Основная",
		MaxThrustX:        2.0,
		MaxThrustY:        15.0,
		FuelConsumptionRate: 1.0,
		BottomSprite: []string{
			"  /\\  ",
		},
	},
	{
		Name:              "Ускоритель",
		MaxThrustX:        1.0,
		MaxThrustY:        25.0,
		FuelConsumptionRate: 2.0,
		BottomSprite: []string{
			" /||\\ ",
		},
	},
	{
		Name:              "Маневровый",
		MaxThrustX:        3.5,
		MaxThrustY:        10.0,
		FuelConsumptionRate: 0.7,
		BottomSprite: []string{
			" <||> ",
		},
	},
}