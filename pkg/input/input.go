package input

import "github.com/gdamore/tcell/v2"

// EventQueue создаёт канал, в который постоянно поступают события tcell.
func EventQueue(screen tcell.Screen) chan tcell.Event {
	eventQueue := make(chan tcell.Event)
	go func() {
		for {
			eventQueue <- screen.PollEvent()
		}
	}()
	return eventQueue
}
