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
var worldWidth, worldHeight = 10000, 20000
var groundLevel = 1

const gravityCutoff = 150.0 // Если "альтитуда" (расстояние от земли) больше 150, гравитация не действует

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

// Структура для дерева
type Tree struct {
	x, y   int      // мировые координаты (левый верхний угол спрайта)
	sprite []string // ASCII-спрайт дерева
}

// Пример ASCII-спрайта дерева
var treeSprite = []string{
	"  ^  ",
	" /|\\ ",
	"  |  ",
}

// Глобальный срез деревьев
var trees []Tree

// Структура для облака
type Cloud struct {
	x, y   int      // позиция в мировых координатах
	sprite []string // ASCII-спрайт облака
}

var explosionSprite = []string{
	"   ***   ",
	"  *****  ",
	" ******* ",
	"*********",
	" ******* ",
	"  *****  ",
	"   ***   ",
}

// Пример спрайта облака (можно изменить по вкусу)
var cloudSprite = []string{
	"  ~~  ",
	"~~~~~~",
	"  ~~  ",
}

var clouds []Cloud

// Функция инициализации облаков
// Здесь облака генерируются только в заданной зоне по оси Y, например, от 10 до 30
func initClouds(n int) {
	clouds = make([]Cloud, n)
	for i := 0; i < n; i++ {
		clouds[i] = Cloud{
			x:      rand.Intn(worldWidth), // по всей ширине мира
			y:      10 + rand.Intn(20),    // y от 10 до 30
			sprite: cloudSprite,
		}
	}
}

// initTrees генерирует n деревьев в зоне около земли.
func initTrees(n int) {
	trees = make([]Tree, n)
	for i := 0; i < n; i++ {
		trees[i] = Tree{
			x:      rand.Intn(worldWidth),         // по всей ширине мира
			y:      groundLevel - len(treeSprite), // чтобы дерево "стоило" на земле
			sprite: treeSprite,
		}
	}
}

type Rocket struct {
	x, y   int     // позиция (левый верхний угол спрайта)
	vx, vy float64 // скорость по осям x и y
	thrust float64 // текущая тяга (значение, влияющее на ускорение)
	fuel   float64 // оставшееся топливо
}

func (r *Rocket) update(dt float64) {
	// Вычисляем "альтитуду" как расстояние от земли:
	altitude := float64(groundLevel - r.y)

	// Определяем, действует ли гравитация:
	var gravityEffect float64
	if altitude < gravityCutoff {
		gravityEffect = 9.8 // стандартное значение гравитации
	} else {
		gravityEffect = 0 // выше cutoff гравитация не действует
	}

	const decelerationFactor = 0.3 // если тяга ниже базовой, применяется ослабление
	// Чистое ускорение – разница между текущей тягой и действующей гравитацией
	netAcc := r.thrust - gravityEffect
	if netAcc < 0 {
		netAcc *= decelerationFactor
	}
	// Обновляем вертикальную скорость:
	// Заметим, что в нашей системе ось Y растёт вниз, поэтому для движения вверх мы уменьшаем vy
	r.vy -= netAcc * dt

	// Обновляем положение ракеты
	r.x += int(r.vx * dt)
	r.y += int(r.vy * dt)

	// Предотвращаем проваливание под землю:
	if r.y+len(rocketSprite) > groundLevel {
		r.y = groundLevel - len(rocketSprite)
		r.vy = 0
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

func drawTrees(cameraX, cameraY, screenWidth, screenHeight int) {
	for _, tree := range trees {
		screenX := tree.x - cameraX
		screenY := tree.y - cameraY
		// Проверяем, что дерево хотя бы частично в видимой области
		if screenX+len(tree.sprite[0]) >= 0 && screenX < screenWidth &&
			screenY+len(tree.sprite) >= 0 && screenY < screenHeight {
			drawSprite(screenX, screenY, tree.sprite, termbox.ColorGreen, termbox.ColorBlack)
		}
	}
}

func drawClouds(cameraX, cameraY, screenWidth, screenHeight int) {
	for _, cloud := range clouds {
		// Вычисляем позицию облака на экране
		screenX := cloud.x - cameraX
		screenY := cloud.y - cameraY
		// Проверяем, находится ли спрайт облака хотя бы частично в видимой области
		if screenX+len(cloud.sprite[0]) >= 0 && screenX < screenWidth &&
			screenY+len(cloud.sprite) >= 0 && screenY < screenHeight {
			drawSprite(screenX, screenY, cloud.sprite, termbox.ColorWhite, termbox.ColorBlack)
		}
	}
}

// Отрисовка звезд в видимой области
func drawStars(cameraX, cameraY, screenWidth, screenHeight int) {
	for sy := 0; sy < screenHeight; sy++ {
		for sx := 0; sx < screenWidth; sx++ {
			worldX := sx + cameraX
			worldY := sy + cameraY
			if isStarAt(worldX, worldY) {
				termbox.SetCell(sx, sy, '*', termbox.ColorYellow, termbox.ColorBlack)
			}
		}
	}
}

// drawNotificationBox рисует уведомление в рамочке в правом верхнем углу.
// screenWidth — ширина экрана, message — текст уведомления.
func drawNotificationBox(screenWidth int, message string) {
	// Определяем размеры рамки: 2 символа отступа слева и справа, 2 символа для рамки
	boxWidth := len(message) + 4 // 1 пробел слева, 1 справа и 2 символа рамки
	// boxHeight := 3               // boxHeight - Верхняя граница, строка с текстом и нижняя граница

	startX := screenWidth - boxWidth // Выравниваем по правому краю
	startY := 0                      // Верхняя строка экрана

	// Верхняя граница: +---+
	for x := startX; x < startX+boxWidth; x++ {
		var ch rune
		if x == startX || x == startX+boxWidth-1 {
			ch = '+'
		} else {
			ch = '-'
		}
		termbox.SetCell(x, startY, ch, termbox.ColorMagenta, termbox.ColorBlack)
	}

	// Средняя строка с текстом: | message |
	midY := startY + 1
	termbox.SetCell(startX, midY, '|', termbox.ColorMagenta, termbox.ColorBlack)
	// Заполняем пробелами и текстом
	for i := 0; i < boxWidth-2; i++ {
		ch := ' '
		if i == 1 && len(message) > 0 {
			// начинаем вставку сообщения
			for j, r := range message {
				// Записываем символы сообщения начиная с позиции startX+1
				termbox.SetCell(startX+1+j, midY, r, termbox.ColorMagenta, termbox.ColorBlack)
			}
			// После вставки сообщения i смещаем до конца сообщения
			i = len(message)
		}
		termbox.SetCell(startX+1+i, midY, ch, termbox.ColorMagenta, termbox.ColorBlack)
	}
	termbox.SetCell(startX+boxWidth-1, midY, '|', termbox.ColorMagenta, termbox.ColorBlack)

	// Нижняя граница: +---+
	bottomY := startY + 2
	for x := startX; x < startX+boxWidth; x++ {
		var ch rune
		if x == startX || x == startX+boxWidth-1 {
			ch = '+'
		} else {
			ch = '-'
		}
		termbox.SetCell(x, bottomY, ch, termbox.ColorMagenta, termbox.ColorBlack)
	}
}

// isStarAt возвращает true, если в мировых координатах (x, y) должна быть звезда.
func isStarAt(x, y int) bool {
	// Простейший хэш: перемножаем координаты на простые числа
	h := int64(x)*73856093 ^ int64(y)*19349663
	if h < 0 {
		h = -h
	}
	// Вероятность появления звезды, например, 3%
	return h%100 < 3
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
	stats := fmt.Sprintf("Speed: %.2f | Thrust: %.2f | Fuel: %.2f", r.vy, r.thrust, r.fuel)
	for i, ch := range stats {
		termbox.SetCell(i, screenHeight-1, ch, termbox.ColorGreen, termbox.ColorBlack)
	}
}

// getSkyColor возвращает цвет фона (неба) в зависимости от высоты ракеты.
// Предположим, что "альтитуда" = groundLevel - rocket.y.
func getSkyColor(r Rocket) termbox.Attribute {
	altitude := groundLevel - r.y
	switch {
	case altitude < 50:
		return termbox.ColorBlue // низкая высота — классический голубой
	case altitude < 100:
		return termbox.ColorCyan // выше — светлее
	case altitude < 150:
		return termbox.ColorMagenta // ещё выше — переход к фиолетовому
	default:
		return termbox.ColorBlack // очень высокая — почти космос
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	initStars(100)

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	// Создаём канал для событий ввода
	eventQueue := make(chan termbox.Event)

	// Запускаем горутину, которая постоянно получает события ввода
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	const hoverThrust = 9.8 // базовая тяга, необходимая для поддержания на месте

	// Инициализация ракеты: стартуем с Земли
	rocket := Rocket{
		x:      worldWidth/2 - len(rocketSprite[0])/2,
		y:      groundLevel - len(rocketSprite),
		vx:     0,
		vy:     0,
		thrust: hoverThrust, // стартуем с "поддерживающей" тягой
		fuel:   100,
	}
	initClouds(200) // например, создаём 10 облаков
	initTrees(200)  // или любое нужное число деревьев

	lastTime := time.Now()

	// Основной игровой цикл
	for {
		now := time.Now()
		dt := now.Sub(lastTime).Seconds()
		lastTime = now

		// Обработка событий ввода (неблокирующая)
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				// Нажатие клавиши Esc – выход из игры
				if ev.Key == termbox.KeyEsc {
					return
				}
				// Обработка клавиш для изменения тяги:
				// Предположим, что KeyCtrlA увеличивает тягу, а KeyCtrlZ – уменьшает
				switch ev.Key {
				case termbox.KeyCtrlA:
					rIncrement := 1.0
					rocket.thrust += rIncrement
					// Опционально: расход топлива можно рассчитывать здесь
				case termbox.KeyCtrlZ:
					rDecrement := 1.0
					rocket.thrust -= rDecrement
					if rocket.thrust < 0 {
						rocket.thrust = 0
					}
				// Можно оставить и обработку стрелок, если она нужна для горизонтального движения
				case termbox.KeyArrowLeft:
					rocket.vx -= 10
				case termbox.KeyArrowRight:
					rocket.vx += 10
				}
			}
		default:
			// если нет событий – ничего не делаем
		}

		// Обновляем состояние ракеты
		rocket.update(dt)

		// Определяем размеры экрана и позицию камеры. Камеру можно смещать так,
		// чтобы ракета была примерно в центре экрана:
		screenWidth, screenHeight := termbox.Size()
		cameraX := rocket.x - screenWidth/2
		cameraY := rocket.y - screenHeight/2

		const safeLandingSpeed = 20.0
		// Если ракета касается земли и скорость превышает безопасное значение,
		// считаем это аварийным приземлением.
		if rocket.y+len(rocketSprite) >= groundLevel && rocket.vy > safeLandingSpeed {
			// Отобразим взрыв на месте ракеты
			drawSprite(rocket.x-cameraX, rocket.y-cameraY, explosionSprite, termbox.ColorRed, termbox.ColorBlack)
			termbox.Flush()
			time.Sleep(2 * time.Second) // задержка для показа взрыва

			// Сброс игры: восстановим состояние ракеты (можно добавить сброс других параметров)
			rocket.x = worldWidth/2 - len(rocketSprite[0])/2
			rocket.y = groundLevel - len(rocketSprite)
			rocket.vx = 0
			rocket.vy = 0
			rocket.fuel = 100

			// При необходимости, можно также сбросить состояние облаков или другие элементы
		}

		skyColor := getSkyColor(rocket)
		// Очистка экрана
		termbox.Clear(termbox.ColorDefault, skyColor)

		// Рисуем звезды, Землю, ракету и статистику с учетом смещения камеры
		drawClouds(cameraX, cameraY, screenWidth, screenHeight)
		drawStars(cameraX, cameraY, screenWidth, screenHeight)
		drawGround(cameraX, cameraY, screenWidth, screenHeight)
		drawTrees(cameraX, cameraY, screenWidth, screenHeight) // отрисовка деревьев
		// Ракету рисуем в мировых координатах с корректировкой на камеру:
		drawSprite(rocket.x-cameraX, rocket.y-cameraY, rocketSprite, termbox.ColorWhite, termbox.ColorBlack)
		drawStats(rocket, screenWidth, screenHeight)

		const cosmicSpeedThreshold = 100.0

		// Если скорость ракеты превышает порог (например, по оси Y)
		if rocket.vy > cosmicSpeedThreshold {
			drawNotificationBox(screenWidth, "COSMIC SPEED!")
		}

		termbox.Flush()
		time.Sleep(30 * time.Millisecond)
	}
}
