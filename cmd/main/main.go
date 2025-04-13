// Package main.go
package main

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/shameoff/rocket-in-console/pkg/input"
	"github.com/shameoff/rocket-in-console/pkg/objects"
	"github.com/shameoff/rocket-in-console/pkg/physics"
	"github.com/shameoff/rocket-in-console/pkg/render"
)

const (
	verticalStep     = 0.5
	horizontalStep   = 0.5
	thrustDecayRate  = 0.0 // Скорость снижения тяги при отпускании клавиши
	decayRate        = 0.3
	safeLandingSpeed = 20.0
)

func processInput(rocket *objects.Rocket, eventQueue chan tcell.Event, dt float64) bool {
	// Флаги нажатия клавиш в текущем цикле
	upPressed := false
	downPressed := false
	leftPressed := false
	rightPressed := false
	stageToggled := false

	// Проверка событий клавиатуры
	for {
		select {
		case ev := <-eventQueue:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyCtrlC:
					return true
				case tcell.KeyUp:
					upPressed = true
				case tcell.KeyDown:
					downPressed = true
				case tcell.KeyLeft:
					leftPressed = true
				case tcell.KeyRight:
					rightPressed = true
				}

				switch ev.Rune() {
				case 'w', 'W':
					upPressed = true
				case 's', 'S':
					downPressed = true
				case 'a', 'A':
					leftPressed = true
				case 'd', 'D':
					rightPressed = true
				case ' ': // Использование руны пробела вместо tcell.KeySpace
					stageToggled = true
				case 'q', 'Q':
					return true
				}
			}
		default:
			// Выходим из цикла, если в очереди больше нет событий
			goto processMovement
		}
	}

processMovement:
	// Переключение ступеней
	if stageToggled {
		rocket.ActiveStage = (rocket.ActiveStage + 1) % len(objects.RocketStages)
	}

	// Получаем текущую ступень
	currentStage := objects.RocketStages[rocket.ActiveStage]

	// Обработка вертикального движения с учётом макс. тяги текущей ступени
	if upPressed {
		rocket.ThrustY += verticalStep
		if rocket.ThrustY > currentStage.MaxThrustY {
			rocket.ThrustY = currentStage.MaxThrustY
		}
	} else if downPressed {
		rocket.ThrustY -= verticalStep
	} else {
		// При отпускании клавиш, тяга быстро падает до нуля (не до гравитации!)
		rocket.ThrustY += (0 - rocket.ThrustY) * thrustDecayRate * dt
	}

	// Обработка горизонтального движения с учётом макс. тяги текущей ступени
	if leftPressed {
		rocket.ThrustX -= horizontalStep
		if rocket.ThrustX < -currentStage.MaxThrustX {
			rocket.ThrustX = -currentStage.MaxThrustX
		}
	} else if rightPressed {
		rocket.ThrustX += horizontalStep
		if rocket.ThrustX > currentStage.MaxThrustX {
			rocket.ThrustX = currentStage.MaxThrustX
		}
	} else {
		// Если клавиши не нажаты, плавно снижаем тягу
		rocket.ThrustX += (0 - rocket.ThrustX) * thrustDecayRate * dt
	}

	return false
}

func updateGame(rocket *objects.Rocket, dt float64, hoverThrust float64) {
	// Remove the gravityCutoff parameter since our new physics model handles this
	physics.UpdateRocket(rocket, dt, objects.GroundLevel, hoverThrust)
}

func handleCollisions(screen tcell.Screen, rocket *objects.Rocket) {
	// Получаем текущий спрайт ракеты
	rocketSprite := rocket.GetRocketSprite()
	
	if rocket.Y+len(rocketSprite) >= objects.GroundLevel && rocket.Vy > safeLandingSpeed {
		screenWidth, screenHeight := screen.Size()
		cameraX := rocket.X - screenWidth/2
		cameraY := rocket.Y - screenHeight/2
		render.DrawSprite(screen, rocket.X-cameraX, rocket.Y-cameraY, objects.ExplosionSprite, tcell.ColorRed, tcell.ColorBlack)
		screen.Show()
		time.Sleep(2 * time.Second)
		rocket.X = objects.WorldWidth/2 - len(objects.RocketSprite[0])/2
		rocket.Y = objects.GroundLevel - len(objects.RocketSprite)
		rocket.Vx = 0
		rocket.Vy = 0
		rocket.Fuel = 100
	}
}

func renderFrame(screen tcell.Screen, rocket *objects.Rocket) {
	screenWidth, screenHeight := screen.Size()
	cameraX := rocket.X - screenWidth/2
	cameraY := rocket.Y - screenHeight/2
	skyColor := render.GetSkyColor(rocket, objects.GroundLevel)
	screen.Clear()

	// Устанавливаем фоновый цвет неба
	style := tcell.StyleDefault.Background(skyColor)
	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			screen.SetContent(x, y, ' ', nil, style)
		}
	}

	render.DrawClouds(screen, objects.Clouds, cameraX, cameraY, screenWidth, screenHeight)
	render.DrawStars(screen, cameraX, cameraY, screenWidth, screenHeight, objects.Stars, objects.IsStarAt)
	render.DrawGround(screen, cameraX, cameraY, screenWidth, screenHeight, objects.GroundLevel)
	render.DrawTrees(screen, objects.Trees, cameraX, cameraY, screenWidth, screenHeight)
	
	// Используем динамический спрайт ракеты вместо статичного
	rocketSprite := rocket.GetRocketSprite()
	render.DrawSprite(screen, rocket.X-cameraX, rocket.Y-cameraY, rocketSprite, tcell.ColorWhite, tcell.ColorBlack)
	
	render.DrawExhaust(screen, rocket, cameraX, cameraY)
	render.DrawStats(screen, rocket, objects.GroundLevel)
	
	// Отображаем информацию о текущей ступени ракеты без пробела
	stageName := objects.RocketStages[rocket.ActiveStage].Name
	render.DrawText(screen, 1, 1, "Stage:" + stageName, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	
	const cosmicSpeedThreshold = 100.0
	if rocket.Vy > cosmicSpeedThreshold {
		render.DrawNotificationBox(screen, screenWidth, "COSMIC SPEED!")
	}
	screen.Show()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	objects.InitStars(100)
	objects.InitClouds(200)
	objects.InitTrees(200)

	// Инициализация tcell
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := screen.Init(); err != nil {
		panic(err)
	}
	defer screen.Fini()

	// Настройка экрана
	screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	screen.Clear()

	eventQueue := input.EventQueue(screen)

	hoverThrust := physics.StandardGravity

	rocket := &objects.Rocket{
		X:           objects.WorldWidth/2 - len(objects.RocketSprite[0])/2,
		Y:           objects.GroundLevel - len(objects.RocketSprite),
		Vx:          0,
		Vy:          0,
		ThrustX:     0,
		ThrustY:     hoverThrust,
		Fuel:        10000,
		ActiveStage: 0,
	}

	lastTime := time.Now()

	for {
		now := time.Now()
		dt := now.Sub(lastTime).Seconds()
		lastTime = now

		if processInput(rocket, eventQueue, dt) {
			return
		}

		updateGame(rocket, dt, hoverThrust)
		handleCollisions(screen, rocket)
		renderFrame(screen, rocket)
		time.Sleep(30 * time.Millisecond)
	}
}
