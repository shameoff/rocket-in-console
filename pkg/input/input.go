package input

import "github.com/nsf/termbox-go"

// EventQueue создаёт канал, в который постоянно поступают события termbox.
func EventQueue() chan termbox.Event {
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()
	return eventQueue
}
