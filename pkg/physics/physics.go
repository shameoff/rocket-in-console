package physics

import (
	"math"

	"github.com/shameoff/rocket-in-console/pkg/objects"
)

// Earth radius in meters
const EarthRadius = 6371000.0

// Standard Earth gravity at sea level (m/s²)
const StandardGravity = 9.80665

// Kármán line - conventional boundary of space (meters)
const KarmanLine = 100000.0

// Масштабный коэффициент для перевода игровых единиц в реальные
const GameToRealScale = 100.0

// Calculates gravity strength at given altitude using inverse square law
func CalculateGravity(altitude float64) float64 {
	// Convert game altitude units to meters
	altitudeInMeters := altitude * GameToRealScale

	// Calculate gravity using inverse square law
	// g = GM/r² = g₀*(R/(R+h))²
	// where g₀ is standard gravity, R is Earth radius, h is altitude
	gravity := StandardGravity * math.Pow(EarthRadius/(EarthRadius+altitudeInMeters), 2)
	return gravity
}

// UpdateRocket обновляет состояние ракеты с учётом реалистичной гравитации и характеристик текущей ступени.
func UpdateRocket(r *objects.Rocket, dt float64, groundLevel int, hoverThrust float64) {
	// Вычисляем "альтитуду" (расстояние от земли)
	altitude := float64(groundLevel - r.Y)

	// Рассчитываем силу гравитации на текущей высоте
	gravity := CalculateGravity(altitude)
	
	// Получаем текущую ступень и её характеристики
	currentStage := objects.RocketStages[r.ActiveStage]

	// Учитываем расход топлива в зависимости от ступени
	if r.ThrustY > gravity {
		fuelUsed := (r.ThrustY - gravity) * dt * currentStage.FuelConsumptionRate
		r.Fuel -= fuelUsed
		if r.Fuel < 0 {
			r.Fuel = 0
			r.ThrustY = 0 // Топливо закончилось, тяги нет
		}
	}

	// Вычисляем модификаторы ускорения на основе характеристик ступени
	// Текущая ступень определяет эффективность тяги
	thrustEfficiencyY := currentStage.MaxThrustY / 15.0 // Нормализуем относительно базовой ступени
	thrustEfficiencyX := currentStage.MaxThrustX / 2.0  // Нормализуем относительно базовой ступени
	
	// Ограничиваем тягу максимальной для текущей ступени
	appliedThrustY := r.ThrustY
	if appliedThrustY > currentStage.MaxThrustY {
		appliedThrustY = currentStage.MaxThrustY
	}
	
	appliedThrustX := r.ThrustX
	if math.Abs(appliedThrustX) > currentStage.MaxThrustX {
		if appliedThrustX > 0 {
			appliedThrustX = currentStage.MaxThrustX
		} else {
			appliedThrustX = -currentStage.MaxThrustX
		}
	}

	// Вычисляем чистое ускорение с учетом эффективности ступени
	netAccY := appliedThrustY*thrustEfficiencyY - gravity
	netAccX := appliedThrustX * thrustEfficiencyX

	// Интегрируем ускорение в скорость
	// Отрицательная Vy означает движение вверх, положительная - вниз
	r.Vy -= netAccY * dt
	r.Vx += netAccX * dt

	// Аккумулируем дробные значения перемещения
	r.AccumulatedX += r.Vx * dt
	r.AccumulatedY += r.Vy * dt

	// Перемещаем ракету только когда накопленное смещение >= 1 пиксель
	deltaX := int(r.AccumulatedX)
	deltaY := int(r.AccumulatedY)

	if deltaX != 0 {
		r.X += deltaX
		r.AccumulatedX -= float64(deltaX)
	}

	if deltaY != 0 {
		r.Y += deltaY
		r.AccumulatedY -= float64(deltaY)
	}

	// Получаем актуальный спрайт ракеты для проверки столкновений
	rocketSprite := r.GetRocketSprite()

	// Проверка на касание земли
	if r.Y+len(rocketSprite) > groundLevel && r.Vy >= 0 {
		r.Y = groundLevel - len(rocketSprite)
		r.Vy = 0
		r.AccumulatedY = 0
	}

	// Затухание горизонтальной скорости
	r.Vx *= math.Pow(0.9, dt)
}
