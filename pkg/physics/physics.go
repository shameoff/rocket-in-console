package physics

import (
	"github.com/shameoff/rocket-in-console/pkg/objects"
	"math"
)

// UpdateRocket обновляет состояние ракеты с учётом тяги, гравитации и высоты.
func UpdateRocket(r *objects.Rocket, dt float64, groundLevel int, gravityCutoff float64) {
	// Вычисляем "альтитуду" как расстояние от земли:
	altitude := float64(groundLevel - r.Y)

	// Определяем, действует ли гравитация:
	var gravityEffect float64
	if altitude < gravityCutoff {
		gravityEffect = 9.8
	} else {
		gravityEffect = 0
	}

	const decelerationFactor = 0.3
	// Для вертикали: чистое ускорение – разница между вертикальной тягой и гравитацией
	netAccY := r.ThrustY - gravityEffect
	if netAccY < 0 {
		netAccY *= decelerationFactor
	}
	// Заметим: поскольку ось Y растёт вниз, чтобы подняться, надо уменьшать Vy
	r.Vy -= netAccY * dt

	// Для горизонтали: просто обновляем скорость на основе тяги
	r.Vx += r.ThrustX * dt

	// Обновляем положение ракеты
	r.X += int(r.Vx * dt)
	r.Y += int(r.Vy * dt)

	// Предотвращаем проваливание под землю
	if r.Y+len(objects.RocketSprite) > groundLevel {
		r.Y = groundLevel - len(objects.RocketSprite)
		r.Vy = 0
	}

	// (Опционально можно добавить затухание горизонтальной скорости, чтобы ракета не "бесконечно" скользила)
	r.Vx *= math.Pow(0.9, dt)
}
