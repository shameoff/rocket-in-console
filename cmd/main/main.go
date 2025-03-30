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

const (
	verticalStep     = 0.5
	horizontalStep   = 0.5
	decayRate        = 0.3
	safeLandingSpeed = 20.0
	gravityCutoff    = 150.0
)

func processInput(rocket *objects.Rocket, eventQueue chan termbox.Event, dt float64, hoverThrust float64) bool {
	upPressed, downPressed, leftPressed, rightPressed := false, false, false, false

	// Обрабатываем все накопленные события
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				if ev.Key == termbox.KeyEsc {
					return true
				}
				switch ev.Key {
				case termbox.KeyArrowUp:
					upPressed = true
				case termbox.KeyArrowDown:
					downPressed = true
				case termbox.KeyArrowLeft:
					leftPressed = true
				case termbox.KeyArrowRight:
					rightPressed = true
				}
			}
		default:
			goto AfterInput
		}
	}
AfterInput:
	if upPressed {
		rocket.ThrustY += verticalStep
	}
	if downPressed {
		rocket.ThrustY -= verticalStep
	}
	if leftPressed {
		rocket.ThrustX -= horizontalStep
	}
	if rightPressed {
		rocket.ThrustX += horizontalStep
	}

	// Определяем альтитуду (расстояние от земли)
	altitude := float64(objects.GroundLevel - rocket.Y)
	if altitude < gravityCutoff {
		rocket.ThrustY += (hoverThrust - rocket.ThrustY) * decayRate * dt
	} else {
		rocket.ThrustY += (0 - rocket.ThrustY) * decayRate * dt
	}
	rocket.ThrustX += (0 - rocket.ThrustX) * decayRate * dt

	return false
}

func updateGame(rocket *objects.Rocket, dt float64, hoverThrust float64) {
	physics.UpdateRocket(rocket, dt, objects.GroundLevel, gravityCutoff, hoverThrust)
}

func handleCollisions(rocket *objects.Rocket) {
	if rocket.Y+len(objects.RocketSprite) >= objects.GroundLevel && rocket.Vy > safeLandingSpeed {
		screenWidth, screenHeight := termbox.Size()
		cameraX := rocket.X - screenWidth/2
		cameraY := rocket.Y - screenHeight/2
		render.DrawSprite(rocket.X-cameraX, rocket.Y-cameraY, objects.ExplosionSprite, termbox.ColorRed, termbox.ColorBlack)
		termbox.Flush()
		time.Sleep(2 * time.Second)
		rocket.X = objects.WorldWidth/2 - len(objects.RocketSprite[0])/2
		rocket.Y = objects.GroundLevel - len(objects.RocketSprite)
		rocket.Vx = 0
		rocket.Vy = 0
		rocket.Fuel = 100
	}
}

func renderFrame(rocket *objects.Rocket) {
	screenWidth, screenHeight := termbox.Size()
	cameraX := rocket.X - screenWidth/2
	cameraY := rocket.Y - screenHeight/2
	skyColor := render.GetSkyColor(rocket, objects.GroundLevel)
	termbox.Clear(termbox.ColorDefault, skyColor)
	render.DrawClouds(objects.Clouds, cameraX, cameraY, screenWidth, screenHeight)
	render.DrawStars(cameraX, cameraY, screenWidth, screenHeight, objects.Stars, objects.IsStarAt)
	render.DrawGround(cameraX, cameraY, screenWidth, screenHeight, objects.GroundLevel)
	render.DrawTrees(objects.Trees, cameraX, cameraY, screenWidth, screenHeight)
	render.DrawSprite(rocket.X-cameraX, rocket.Y-cameraY, objects.RocketSprite, termbox.ColorWhite, termbox.ColorBlack)
	render.DrawExhaust(rocket, cameraX, cameraY)
	render.DrawStats(rocket, screenWidth, screenHeight)
	const cosmicSpeedThreshold = 100.0
	if rocket.Vy > cosmicSpeedThreshold {
		render.DrawNotificationBox(screenWidth, "COSMIC SPEED!")
	}
	termbox.Flush()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	objects.InitStars(100)
	objects.InitClouds(200)
	objects.InitTrees(200)

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := input.EventQueue()
	hoverThrust := objects.HoverThrust

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

		if processInput(rocket, eventQueue, dt, hoverThrust) {
			return
		}

		updateGame(rocket, dt, hoverThrust)
		handleCollisions(rocket)
		renderFrame(rocket)
		time.Sleep(30 * time.Millisecond)
	}
}
