package render

import (
	"fmt"
	"math"

	"github.com/nsf/termbox-go"
	"github.com/shameoff/rocket-in-console/pkg/objects"
)

func DrawSprite(x, y int, sprite []string, fg, bg termbox.Attribute) {
	for dy, line := range sprite {
		for dx, ch := range line {
			if ch != ' ' {
				termbox.SetCell(x+dx, y+dy, ch, fg, bg)
			}
		}
	}
}

// ... (функции DrawTrees, DrawClouds, DrawStars, DrawGround, DrawNotificationBox остаются без изменений)

func DrawStats(r *objects.Rocket, screenWidth, screenHeight int) {
	// Вычисляем общую скорость как модуль вектора
	speed := math.Sqrt(r.Vx*r.Vx + r.Vy*r.Vy)
	stats := fmt.Sprintf("Speed: %.2f | Thrust: (%.2f, %.2f) | Fuel: %.2f", speed, r.ThrustX, r.ThrustY, r.Fuel)
	for i, ch := range stats {
		termbox.SetCell(i, screenHeight-1, ch, termbox.ColorGreen, termbox.ColorBlack)
	}
}

func GetSkyColor(r *objects.Rocket, groundLevel int) termbox.Attribute {
	altitude := groundLevel - r.Y
	switch {
	case altitude < 50:
		return termbox.ColorBlue
	case altitude < 100:
		return termbox.ColorCyan
	case altitude < 150:
		return termbox.ColorMagenta
	default:
		return termbox.ColorBlack
	}
}

func DrawTrees(trees []objects.Tree, cameraX, cameraY, screenWidth, screenHeight int) {
	for _, tree := range trees {
		screenX := tree.X - cameraX
		screenY := tree.Y - cameraY
		if screenX+len(tree.Sprite[0]) >= 0 && screenX < screenWidth &&
			screenY+len(tree.Sprite) >= 0 && screenY < screenHeight {
			DrawSprite(screenX, screenY, tree.Sprite, termbox.ColorGreen, termbox.ColorBlack)
		}
	}
}

func DrawClouds(clouds []objects.Cloud, cameraX, cameraY, screenWidth, screenHeight int) {
	for _, cloud := range clouds {
		screenX := cloud.X - cameraX
		screenY := cloud.Y - cameraY
		if screenX+len(cloud.Sprite[0]) >= 0 && screenX < screenWidth &&
			screenY+len(cloud.Sprite) >= 0 && screenY < screenHeight {
			DrawSprite(screenX, screenY, cloud.Sprite, termbox.ColorWhite, termbox.ColorBlack)
		}
	}
}

func DrawStars(cameraX, cameraY, screenWidth, screenHeight int, stars []objects.Star, isStarAt func(x, y int) bool) {
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

func DrawGround(cameraX, cameraY, screenWidth, screenHeight, groundLevel int) {
	screenY := groundLevel - cameraY
	if screenY >= 0 && screenY < screenHeight {
		for x := 0; x < screenWidth; x++ {
			termbox.SetCell(x, screenY, '=', termbox.ColorGreen, termbox.ColorBlack)
		}
	}
}

func DrawNotificationBox(screenWidth int, message string) {
	boxWidth := len(message) + 4
	startX := screenWidth - boxWidth
	startY := 0

	// Верхняя граница
	for x := startX; x < startX+boxWidth; x++ {
		var ch rune
		if x == startX || x == startX+boxWidth-1 {
			ch = '+'
		} else {
			ch = '-'
		}
		termbox.SetCell(x, startY, ch, termbox.ColorMagenta, termbox.ColorBlack)
	}

	// Средняя строка с текстом
	midY := startY + 1
	termbox.SetCell(startX, midY, '|', termbox.ColorMagenta, termbox.ColorBlack)
	for i := 0; i < boxWidth-2; i++ {
		ch := ' '
		if i == 1 && len(message) > 0 {
			for j, r := range message {
				termbox.SetCell(startX+1+j, midY, r, termbox.ColorMagenta, termbox.ColorBlack)
			}
			i = len(message)
		}
		termbox.SetCell(startX+1+i, midY, ch, termbox.ColorMagenta, termbox.ColorBlack)
	}
	termbox.SetCell(startX+boxWidth-1, midY, '|', termbox.ColorMagenta, termbox.ColorBlack)

	// Нижняя граница
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

// DrawExhaust рисует след от сопел, если активирована тяга в горизонтальном или вертикальном направлении.
// Для горизонтали: если ThrustX > 0 (движение вправо) – рисуем след с левой стороны, если ThrustX < 0 – с правой.
// Для вертикали: если ThrustY > HoverThrust – след снизу, если ThrustY < HoverThrust – след сверху.
func DrawExhaust(rocket *objects.Rocket, cameraX, cameraY int) {
	threshold := 0.5
	width := len(objects.RocketSprite[0])
	height := len(objects.RocketSprite)

	// Горизонтальный след (вспомогательные двигатели):
	if rocket.ThrustX > threshold {
		flame := "=>"
		x := rocket.X - cameraX - len(flame)
		y1 := rocket.Y - cameraY + height/3
		y2 := rocket.Y - cameraY + (2 * height / 3)
		for i, r := range flame {
			termbox.SetCell(x+i, y1, r, termbox.ColorBlue, termbox.ColorBlack)
			termbox.SetCell(x+i, y2, r, termbox.ColorBlue, termbox.ColorBlack)
		}
	} else if rocket.ThrustX < -threshold {
		flame := "<="
		x := rocket.X - cameraX + width
		y1 := rocket.Y - cameraY + height/3
		y2 := rocket.Y - cameraY + (2 * height / 3)
		for i, r := range flame {
			termbox.SetCell(x+i, y1, r, termbox.ColorBlue, termbox.ColorBlack)
			termbox.SetCell(x+i, y2, r, termbox.ColorBlue, termbox.ColorBlack)
		}
	}

	// Вертикальный след (главный двигатель)
	if rocket.ThrustY > objects.HoverThrust+threshold {
		flame := "vv"
		x1 := rocket.X - cameraX + width/3
		x2 := rocket.X - cameraX + (2*width)/3 - len(flame)
		y := rocket.Y - cameraY + height
		// Если ThrustY превышает базовое значение, главный двигатель работает – пламя красное.
		for i, r := range flame {
			termbox.SetCell(x1+i, y, r, termbox.ColorRed, termbox.ColorBlack)
			termbox.SetCell(x2+i, y, r, termbox.ColorRed, termbox.ColorBlack)
		}
	} else if rocket.ThrustY < objects.HoverThrust-threshold {
		flame := "^^"
		x1 := rocket.X - cameraX + width/3
		x2 := rocket.X - cameraX + (2*width)/3 - len(flame)
		y := rocket.Y - cameraY - 1
		// Если ThrustY ниже hover, пламя вспомогательных (обратного действия) – синее.
		for i, r := range flame {
			termbox.SetCell(x1+i, y, r, termbox.ColorBlue, termbox.ColorBlack)
			termbox.SetCell(x2+i, y, r, termbox.ColorBlue, termbox.ColorBlack)
		}
	}
}
