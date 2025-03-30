package physics

import (
	"github.com/shameoff/rocket-in-console/pkg/objects"
	"math"
)

// UpdateRocket обновляет состояние ракеты с учетом тяги, гравитации и зоны отключения гравитации.
// Если ракета находится в зоне гравитации (altitude < gravityCutoff), то чистое ускорение равно ThrustY - gravity.
// Если ракета вне зоны (altitude >= gravityCutoff), то чистое ускорение = ThrustY (без гравитации).
// Для отрицательного ускорения (если ThrustY меньше гравитации) применяется дополнительное затухание.
func UpdateRocket(r *objects.Rocket, dt float64, groundLevel int, gravityCutoff float64, hoverThrust float64) {
	// Вычисляем "альтитуду" (расстояние от земли; поскольку y растет вниз, чем меньше r.Y, тем выше ракета)
	altitude := float64(groundLevel - r.Y)

	var netAccY float64
	if altitude < gravityCutoff {
		// В зоне гравитации: чтобы зависать, ThrustY должно быть равно gravity (например, 9.8).
		// Чистое ускорение = ThrustY - gravity.
		netAccY = r.ThrustY - 9.8
		// Если чистое ускорение отрицательное, применим затухание (например, для "смягчения" падения)
		if netAccY < 0 {
			netAccY *= 0.3
		}
	} else {
		// В космосе гравитация не действует, поэтому netAccY = ThrustY.
		netAccY = r.ThrustY
	}

	// Обновляем вертикальную скорость.
	// Заметим, что в нашей системе ось Y растет вниз, поэтому для движения вверх (по желанию) мы уменьшаем Vy.
	r.Vy -= netAccY * dt

	// Обновляем горизонтальную скорость
	r.Vx += r.ThrustX * dt

	// Обновляем положение ракеты
	r.X += int(r.Vx * dt)
	r.Y += int(r.Vy * dt)

	// Предотвращаем проваливание под землю
	if r.Y+len(objects.RocketSprite) > groundLevel {
		r.Y = groundLevel - len(objects.RocketSprite)
		r.Vy = 0
	}

	// Дополнительное затухание горизонтальной скорости (чтобы не накапливалась бесконечно)
	r.Vx *= math.Pow(0.9, dt)
}
