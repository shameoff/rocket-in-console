package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

type Star struct {
	x, y int
}

var stars []Star
var worldWidth, worldHeight = 100, 200
var groundLevel = 180

func initStars(n int) {
	stars = make([]Star, n)
	for i := 0; i < n; i++ {
		stars[i] = Star{
			x: rand.Intn(worldWidth),
			y: rand.Intn(groundLevel),
		}
	}
}

var rocketSprite = []string{
	"  /\\  ",
	" |==| ",
	" |  | ",
	"  ||  ",
}

type Rocket struct {
	x, y   int     // позиция в мире (левый верхний угол спрайта)
	vx, vy float64 // скорость
	fuel   float64
}

func (r *Rocket) update(dt float64) {
	gravity := 9.8
	r.vy += gravity * dt
	r.x += int(r.vx * dt)
	r.y += int(r.vy * dt)

	// Ограничим движение ракеты, чтобы не проваливаться под Землю
	if r.y+len(rocketSprite) > groundLevel {
		r.y = groundLevel - len(rocketSprite)
		r.vy = 0
	}

	if r.fuel > 0 {
		r.fuel -= 0.1 * dt
		if r.fuel < 0 {
			r.fuel = 0
		}
	}
}

func drawSprite(x, y int, sprite []string, fg, bg termbox.Attribute) {
	for dy, line := range sprite {
		for dx, ch := range line {
			if ch != ' ' {
				termbox.SetCell(x+dx, y+dy, ch, fg, bg)
			}
		}
	}
}

func drawStars(cameraX, cameraY, screenWidth, screenHeight int) {
	for _, star := range stars {
		screenX := star.x - cameraX
		screenY := star.y - cameraY
		if screenX >= 0 && screenX < screenWidth && screenY >= 0 && screenY < screenHeight {
			termbox.SetCell(screenX, screenY, '*', termbox.ColorYellow, termbox.ColorBlack)
		}
	}
}

func drawGround(cameraX, cameraY, screenWidth, screenHeight int) {
	screenY := groundLevel - cameraY
	if screenY >= 0 && screenY < screenHeight {
		for x := 0; x < screenWidth; x++ {
			termbox.SetCell(x, screenY, '=', termbox.ColorGreen, termbox.ColorBlack)
		}
	}
}

func drawStats(r Rocket, screenWidth, screenHeight int) {
	stats := fmt.Sprintf("Speed: %.2f | Fuel: %.2f", r.vy, r.fuel)
	for i, ch := range stats {
		termbox.SetCell(i, screenHeight-1, ch, termbox.ColorGreen, termbox.ColorBlack)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	initStars(100)

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	// Инициализация ракеты: стартуем с Земли, по центру мира
	rocket := Rocket{
		x:    worldWidth/2 - len(rocketSprite[0])/2,
		y:    groundLevel - len(rocketSprite),
		vx:   0,
		vy:   0,
		fuel: 100,
	}

	lastTime := time.Now()

	// Основной игровой цикл
	for {
		now := time.Now()
		dt := now.Sub(lastTime).Seconds()
		lastTime = now

		// Обработка событий ввода (например, Esc для выхода)
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				return
			}
			// Можно добавить управление ускорением по клавишам
		}

		// Обновляем состояние ракеты
		rocket.update(dt)

		// Определяем размеры экрана и позицию камеры. Камеру можно смещать так,
		// чтобы ракета была примерно в центре экрана:
		screenWidth, screenHeight := termbox.Size()
		cameraX := rocket.x - screenWidth/2
		cameraY := rocket.y - screenHeight/2

		// Очистка экрана
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		// Рисуем звезды, Землю, ракету и статистику с учетом смещения камеры
		drawStars(cameraX, cameraY, screenWidth, screenHeight)
		drawGround(cameraX, cameraY, screenWidth, screenHeight)
		// Ракету рисуем в мировых координатах с корректировкой на камеру:
		drawSprite(rocket.x-cameraX, rocket.y-cameraY, rocketSprite, termbox.ColorWhite, termbox.ColorBlack)
		drawStats(rocket, screenWidth, screenHeight)

		termbox.Flush()
		time.Sleep(30 * time.Millisecond)
	}
}
