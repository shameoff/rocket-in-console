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

// UpdateRocket обновляет состояние ракеты с учётом реалистичной гравитации.
func UpdateRocket(r *objects.Rocket, dt float64, groundLevel int, hoverThrust float64) {
	// Вычисляем "альтитуду" (расстояние от земли)
	altitude := float64(groundLevel - r.Y)

	// Рассчитываем силу гравитации на текущей высоте
	gravity := CalculateGravity(altitude)

	// Вычисляем чистое ускорение (гравитация действует вниз, тяга - вверх)
	netAccY := r.ThrustY - gravity

	// Интегрируем ускорение в скорость
	// Отрицательная Vy означает движение вверх, положительная - вниз
	r.Vy -= netAccY * dt
	r.Vx += r.ThrustX * dt

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

	// Проверка на касание земли
	if r.Y+len(objects.RocketSprite) > groundLevel && r.Vy >= 0 {
		r.Y = groundLevel - len(objects.RocketSprite)
		r.Vy = 0
		r.AccumulatedY = 0
	}

	// Затухание горизонтальной скорости
	r.Vx *= math.Pow(0.9, dt)
}
