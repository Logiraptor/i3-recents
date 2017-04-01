package focus

import (
	"fmt"

	"github.com/Logiraptor/i3-recents/socket"
)

type Subscribe []string
type Success struct{}

type Event struct {
	Change    string    `json:"change"`
	Container Container `json:"container"`
}

type Container struct {
	ID int `json:"id"`
}

type FocusMonitor struct {
	sock *socket.ValueSocket
	out  chan Event
}

func NewFocusMonitor(sock *socket.ValueSocket) *FocusMonitor {
	return &FocusMonitor{sock: sock}
}

func (f *FocusMonitor) Start() <-chan Event {
	err := f.sock.Write(2, Subscribe{"window"})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	_, err = f.sock.Read(&Success{})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	out := make(chan Event)
	go func() {
		defer close(out)
		for {
			var event Event
			_, err := f.sock.Read(&event)
			if err != nil {
				fmt.Println(err)
				return
			}
			out <- event
		}
	}()
	return out
}
