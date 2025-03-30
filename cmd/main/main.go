package main

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/shameoff/rocket-in-console/pkg/input"
	"github.com/shameoff/rocket-in-console/pkg/objects"
	"github.com/shameoff/rocket-in-console/pkg/physics"
	"github.com/shameoff/rocket-in-console/pkg/render"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Инициализация объектов мира
	objects.InitStars(100)
	objects.InitClouds(200)
	objects.InitTrees(200)

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	// Создаём канал для ввода
	eventQueue := input.EventQueue()

	const hoverThrust = 9.8 // базовая вертикальная тяга для зависания

	// Инициализация ракеты: старт с земли
	rocket := &objects.Rocket{
		X:       objects.WorldWidth/2 - len(objects.RocketSprite[0])/2,
		Y:       objects.GroundLevel - len(objects.RocketSprite),
		Vx:      0,
		Vy:      0,
		ThrustX: 0,
		ThrustY: hoverThrust,
		Fuel:    100,
	}

	lastTime := time.Now()

	for {
		now := time.Now()
		dt := now.Sub(lastTime).Seconds()
		lastTime = now

		// Обработка ввода
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				if ev.Key == termbox.KeyEsc {
					return
				}
				switch ev.Key {
				// Ракета двигается с максимальным ускорением вперед, но с ограничением в стороны
				case termbox.KeyArrowUp:
					rocket.ThrustY += 1.0
				case termbox.KeyArrowDown:
					rocket.ThrustY -= 1.0
					if rocket.ThrustY < -10 {
						rocket.ThrustY = -10
					}
				case termbox.KeyArrowLeft:
					rocket.ThrustX -= 1.0
					if rocket.ThrustX < -10 {
						rocket.ThrustX = -10
					}
				case termbox.KeyArrowRight:
					rocket.ThrustX += 1.0
					if rocket.ThrustX > 10 {
						rocket.ThrustX = 10
					}
				}
			}
		default:
			// если нет событий – продолжаем
		}

		// Обновление физики ракеты
		physics.UpdateRocket(rocket, dt, objects.GroundLevel, 150.0)

		screenWidth, screenHeight := termbox.Size()
		cameraX := rocket.X - screenWidth/2
		cameraY := rocket.Y - screenHeight/2

		const safeLandingSpeed = 20.0
		if rocket.Y+len(objects.RocketSprite) >= objects.GroundLevel && rocket.Vy > safeLandingSpeed {
			render.DrawSprite(rocket.X-cameraX, rocket.Y-cameraY, objects.ExplosionSprite, termbox.ColorRed, termbox.ColorBlack)
			termbox.Flush()
			time.Sleep(2 * time.Second)

			rocket.X = objects.WorldWidth/2 - len(objects.RocketSprite[0])/2
			rocket.Y = objects.GroundLevel - len(objects.RocketSprite)
			rocket.Vx = 0
			rocket.Vy = 0
			rocket.Fuel = 100
		}

		skyColor := render.GetSkyColor(rocket, objects.GroundLevel)
		termbox.Clear(termbox.ColorDefault, skyColor)

		render.DrawClouds(objects.Clouds, cameraX, cameraY, screenWidth, screenHeight)
		render.DrawStars(cameraX, cameraY, screenWidth, screenHeight, objects.Stars, objects.IsStarAt)
		render.DrawGround(cameraX, cameraY, screenWidth, screenHeight, objects.GroundLevel)
		render.DrawTrees(objects.Trees, cameraX, cameraY, screenWidth, screenHeight)
		render.DrawSprite(rocket.X-cameraX, rocket.Y-cameraY, objects.RocketSprite, termbox.ColorWhite, termbox.ColorBlack)
		render.DrawStats(rocket, screenWidth, screenHeight)

		const cosmicSpeedThreshold = 100.0
		if rocket.Vy > cosmicSpeedThreshold {
			render.DrawNotificationBox(screenWidth, "COSMIC SPEED!")
		}

		termbox.Flush()
		time.Sleep(30 * time.Millisecond)
	}
}
