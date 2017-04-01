package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"flag"

	"io"

	"encoding/json"

	"github.com/Logiraptor/i3-recents/focus"
	"github.com/Logiraptor/i3-recents/socket"
)

// To send a message to i3, you have to format in the
// binary message format which i3 expects. This format
// specifies a magic string in the beginning to ensure
// the integrity of messages (to prevent follow-up errors).
// Following the magic string comes the length of the
// payload of the message as 32-bit integer, and the
// type of the message as 32-bit integer (the integers
// are not converted, so they are in native byte order).

// The magic string currently is "i3-ipc" and will only
// be changed when a change in the IPC API is done which
// breaks compatibility (we hope that we don’t need to do that).

// -- MVP
// new states drain the forwards queue
// backwards queue has a limit
// limit is configurable on server
// xxx somehow ignore adding states that are a result of back / forward

const Subscribe = 2

func main() {
	back := flag.Bool("back", false, "Go back to the last focused window")
	flag.Parse()

	if *back {
		r, err := http.Get("http://localhost:7364")
		if err != nil {
			fmt.Println("Failed to connect to i3-recents server: ", err.Error())
		}
		defer r.Body.Close()
		io.Copy(os.Stdout, r.Body)
		fmt.Println()
		return
	}

	socketPath, err := getSocketPath()
	if err != nil {
		fmt.Println("getSocketPath:", err)
		os.Exit(1)
	}
	conn, err := net.Dial("unix", strings.TrimSpace(socketPath))
	if err != nil {
		fmt.Println("DialUnix:", err)
		os.Exit(1)
	}

	sock := socket.NewValueSocket(conn)
	monitor := focus.NewFocusMonitor(sock)

	server := FocusServer{
		goBack: make(chan struct{}),
	}

	go server.dispatch(monitor.Start())

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		server.goBack <- struct{}{}
		json.NewEncoder(rw).Encode(focus.Success{})
	})
	fmt.Println(http.ListenAndServe(":7364", nil))
}

func getSocketPath() (string, error) {
	buf, err := exec.Command("i3", "--get-socketpath").CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

type FocusServer struct {
	goBack chan struct{}
}

func (f *FocusServer) dispatch(changes <-chan focus.Event) {
	var current, last focus.Event
	for {
		select {
		case e, ok := <-changes:
			if !ok {
				panic("channel closed")
			}
			if e != current {
				current, last = e, current
			}
		case <-f.goBack:
			current, last = last, current

			arg := fmt.Sprintf("[con_id=\"%d\"] focus", current.Container.ID)
			exec.Command("i3-msg", arg).CombinedOutput()
		}
	}
}
