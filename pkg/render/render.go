package render

import (
	"fmt"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/shameoff/rocket-in-console/pkg/objects"
	"github.com/shameoff/rocket-in-console/pkg/physics"
)

func DrawSprite(screen tcell.Screen, x, y int, sprite []string, fg, bg tcell.Color) {
	for dy, line := range sprite {
		for dx, ch := range line {
			if ch != ' ' {
				screen.SetContent(x+dx, y+dy, ch, nil, tcell.StyleDefault.Foreground(fg).Background(bg))
			}
		}
	}
}

// DrawText отображает строку текста на экране в указанной позиции с указанным стилем
func DrawText(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	// Проходим по каждой руне (символу) в строке
	for i, r := range text {
		// Устанавливаем содержимое ячейки экрана в указанной позиции
		screen.SetContent(x+i, y, r, nil, style)
	}
}

// DrawTextLines рисует несколько строк текста одну под другой
func DrawTextLines(screen tcell.Screen, x, y int, lines []string, style tcell.Style) {
	currentY := y
	for _, line := range lines {
		DrawText(screen, x, currentY, line, style)
		currentY++
	}
}

// DrawStats отображает статистические данные о полете на экране
func DrawStats(screen tcell.Screen, rocket *objects.Rocket, groundLevel int) {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	// Расчет скорости (сохраняем знак для определения направления)
	speedY := math.Abs(rocket.Vy)
	speedX := rocket.Vx // Сохраняем знак для горизонтальной скорости

	// Расчет расстояния от Земли (высота)
	altitude := float64(groundLevel - rocket.Y)

	// Преобразуем в километры для более наглядного отображения
	altitudeKm := altitude / 10.0

	// Форматируем для отображения: положительное значение Vy - это движение вниз, отрицательное - вверх
	verticalSpeedDirection := ""
	if rocket.Vy < 0 {
		verticalSpeedDirection = "▲" // вверх
	} else if rocket.Vy > 0 {
		verticalSpeedDirection = "▼" // вниз
	}

	horizontalSpeedDirection := ""
	if speedX > 0 {
		horizontalSpeedDirection = "►" // вправо
	} else if speedX < 0 {
		horizontalSpeedDirection = "◄" // влево
	}

	// Рассчитываем текущую гравитацию на этой высоте
	currentGravity := physics.CalculateGravity(altitude)

	// Статистические строки
	stats := []string{
		fmt.Sprintf("Altitude: %.2f km", altitudeKm),
		fmt.Sprintf("Vspeed: %.2f %s", speedY, verticalSpeedDirection),
		fmt.Sprintf("Hspeed: %.2f %s", math.Abs(speedX), horizontalSpeedDirection),
		fmt.Sprintf("Thrust: V=%.2f H=%.2f", rocket.ThrustY, rocket.ThrustX),
		fmt.Sprintf("Gravity: %.2f", currentGravity),
	}

	// Если мы находимся в космосе (выше линии Кармана), добавляем индикатор
	if altitudeKm > physics.KarmanLine/1000.0 {
		stats = append(stats, "*** SPACE ***")
	}

	// Отображаем статистику в углу экрана
	width, _ := screen.Size()
	x := width - 25 // отступ от правого края
	y := 1          // начинаем немного ниже верха экрана

	// Используем новую функцию для отрисовки всех строк статистики
	DrawTextLines(screen, x, y, stats, style)
}

// GetSkyColor returns realistic atmospheric layer colors based on altitude
func GetSkyColor(r *objects.Rocket, groundLevel int) tcell.Color {
	// Calculate altitude in arbitrary units
	altitude := float64(groundLevel - r.Y)

	// Scale to approximate real-world altitudes in kilometers
	// Let's say each unit is about 100 meters
	altitudeKm := altitude / 10.0

	switch {
	case altitudeKm < 12: // Troposphere (0-12 km)
		// Gradual transition from light blue to darker blue
		blueIntensity := uint8(255 - (altitude * 5))
		if blueIntensity < 100 {
			blueIntensity = 100
		}
		return tcell.NewRGBColor(100, 100, int32(blueIntensity))

	case altitudeKm < 50: // Stratosphere (12-50 km)
		// Gradual transition from deep blue to indigo
		progress := (altitudeKm - 12) / 38
		blue := uint8(100 - progress*50)
		return tcell.NewRGBColor(0, 0, int32(blue+50))

	case altitudeKm < 85: // Mesosphere (50-85 km)
		// Dark purple transitioning to near black
		progress := (altitudeKm - 50) / 35
		val := uint8(50 - progress*50)
		return tcell.NewRGBColor(int32(val/2), 0, int32(val))

	default: // Thermosphere (85+ km) / Space
		return tcell.ColorBlack
	}
}

func DrawTrees(screen tcell.Screen, trees []objects.Tree, cameraX, cameraY, screenWidth, screenHeight int) {
	for _, tree := range trees {
		screenX := tree.X - cameraX
		screenY := tree.Y - cameraY
		if screenX+len(tree.Sprite[0]) >= 0 && screenX < screenWidth &&
			screenY+len(tree.Sprite) >= 0 && screenY < screenHeight {
			DrawSprite(screen, screenX, screenY, tree.Sprite, tcell.ColorGreen, tcell.ColorBlack)
		}
	}
}

func DrawClouds(screen tcell.Screen, clouds []objects.Cloud, cameraX, cameraY, screenWidth, screenHeight int) {
	for _, cloud := range clouds {
		screenX := cloud.X - cameraX
		screenY := cloud.Y - cameraY
		if screenX+len(cloud.Sprite[0]) >= 0 && screenX < screenWidth &&
			screenY+len(cloud.Sprite) >= 0 && screenY < screenHeight {
			DrawSprite(screen, screenX, screenY, cloud.Sprite, tcell.ColorWhite, tcell.ColorBlack)
		}
	}
}

func DrawStars(screen tcell.Screen, cameraX, cameraY, screenWidth, screenHeight int, stars []objects.Star, isStarAt func(x, y int) bool) {
	for sy := 0; sy < screenHeight; sy++ {
		for sx := 0; sx < screenWidth; sx++ {
			worldX := sx + cameraX
			worldY := sy + cameraY
			if isStarAt(worldX, worldY) {
				screen.SetContent(sx, sy, '*', nil, tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack))
			}
		}
	}
}

func DrawGround(screen tcell.Screen, cameraX, cameraY, screenWidth, screenHeight, groundLevel int) {
	screenY := groundLevel - cameraY
	if screenY >= 0 && screenY < screenHeight {
		for x := 0; x < screenWidth; x++ {
			screen.SetContent(x, screenY, '=', nil, tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack))
		}
	}
}

func DrawNotificationBox(screen tcell.Screen, screenWidth int, message string) {
	boxWidth := len(message) + 4
	startX := screenWidth - boxWidth
	startY := 0
	style := tcell.StyleDefault.Foreground(tcell.ColorPurple).Background(tcell.ColorBlack)

	// Верхняя граница
	for x := startX; x < startX+boxWidth; x++ {
		var ch rune
		if x == startX || x == startX+boxWidth-1 {
			ch = '+'
		} else {
			ch = '-'
		}
		screen.SetContent(x, startY, ch, nil, style)
	}

	// Средняя строка с текстом
	midY := startY + 1
	screen.SetContent(startX, midY, '|', nil, style)
	for i := 0; i < boxWidth-2; i++ {
		ch := ' '
		if i == 1 && len(message) > 0 {
			for j, r := range message {
				screen.SetContent(startX+1+j, midY, r, nil, style)
			}
			i = len(message)
		}
		screen.SetContent(startX+1+i, midY, ch, nil, style)
	}
	screen.SetContent(startX+boxWidth-1, midY, '|', nil, style)

	// Нижняя граница
	bottomY := startY + 2
	for x := startX; x < startX+boxWidth; x++ {
		var ch rune
		if x == startX || x == startX+boxWidth-1 {
			ch = '+'
		} else {
			ch = '-'
		}
		screen.SetContent(x, bottomY, ch, nil, style)
	}
}

// DrawExhaust рисует след от сопел, если активирована тяга в горизонтальном или вертикальном направлении.
func DrawExhaust(screen tcell.Screen, rocket *objects.Rocket, cameraX, cameraY int) {
	threshold := 0.5
	width := len(objects.RocketSprite[0])
	height := len(objects.RocketSprite)

	// Рассчитаем текущую величину гравитации для определения порога тяги
	altitude := float64(objects.GroundLevel - rocket.Y)
	currentGravity := physics.CalculateGravity(altitude)

	// Горизонтальный след (вспомогательные двигатели):
	blueStyle := tcell.StyleDefault.Foreground(tcell.ColorBlue).Background(tcell.ColorBlack)
	redStyle := tcell.StyleDefault.Foreground(tcell.ColorRed).Background(tcell.ColorBlack)

	if rocket.ThrustX > threshold {
		flame := "=>"
		x := rocket.X - cameraX - len(flame)
		y1 := rocket.Y - cameraY + height/3
		y2 := rocket.Y - cameraY + (2 * height / 3)
		for i, r := range flame {
			screen.SetContent(x+i, y1, r, nil, blueStyle)
			screen.SetContent(x+i, y2, r, nil, blueStyle)
		}
	} else if rocket.ThrustX < -threshold {
		flame := "<="
		x := rocket.X - cameraX + width
		y1 := rocket.Y - cameraY + height/3
		y2 := rocket.Y - cameraY + (2 * height / 3)
		for i, r := range flame {
			screen.SetContent(x+i, y1, r, nil, blueStyle)
			screen.SetContent(x+i, y2, r, nil, blueStyle)
		}
	}

	// Вертикальный след (главный двигатель)
	// Теперь используем динамическую гравитацию вместо константы
	if rocket.ThrustY > currentGravity+threshold {
		flame := "vv"
		x1 := rocket.X - cameraX + width/3
		x2 := rocket.X - cameraX + (2*width)/3 - len(flame)
		y := rocket.Y - cameraY + height
		for i, r := range flame {
			screen.SetContent(x1+i, y, r, nil, redStyle)
			screen.SetContent(x2+i, y, r, nil, redStyle)
		}
	} else if rocket.ThrustY < currentGravity-threshold {
		flame := "^^"
		x1 := rocket.X - cameraX + width/3
		x2 := rocket.X - cameraX + (2*width)/3 - len(flame)
		y := rocket.Y - cameraY - 1
		for i, r := range flame {
			screen.SetContent(x1+i, y, r, nil, blueStyle)
			screen.SetContent(x2+i, y, r, nil, blueStyle)
		}
	}
}
