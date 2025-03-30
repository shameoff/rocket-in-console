// Package physics.go
package physics

import (
	"github.com/shameoff/rocket-in-console/pkg/objects"
	"math"
)

// UpdateRocket обновляет состояние ракеты.
// Если ракета в зоне гравитации (altitude < gravityCutoff), чистое ускорение = ThrustY - 9.8 (гравитация).
// Если вне зоны, чистое ускорение = ThrustY.
// При отрицательном ускорении применяется небольшое затухание.
func UpdateRocket(r *objects.Rocket, dt float64, groundLevel int, gravityCutoff float64, hoverThrust float64) {
	// Вычисляем "альтитуду" (расстояние от земли; поскольку y растет вниз, чем меньше r.Y, тем выше ракета)
	altitude := float64(groundLevel - r.Y)

	var netAccY float64
	if altitude < gravityCutoff {
		netAccY = r.ThrustY - 9.8
		if netAccY < 0 {
			netAccY *= 0.3
		}
	} else {
		netAccY = r.ThrustY
	}

	// Интегрируем ускорение в скорость
	r.Vy -= netAccY * dt
	r.Vx += r.ThrustX * dt

	// Интегрируем скорость в положение
	r.X += int(r.Vx * dt)
	r.Y += int(r.Vy * dt)

	// Позволяем ракете покидать землю, если Vy < 0 (то есть она движется вверх).
	// Если же Vy >= 0 и ракета касается земли, то фиксируем положение.
	if r.Y+len(objects.RocketSprite) > groundLevel && r.Vy >= 0 {
		r.Y = groundLevel - len(objects.RocketSprite)
		r.Vy = 0
	}

	// Затухание горизонтальной скорости
	r.Vx *= math.Pow(0.9, dt)
}
