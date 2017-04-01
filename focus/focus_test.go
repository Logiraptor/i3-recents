package focus

import "testing"
import "github.com/Logiraptor/i3-recents/socket"
import "bytes"
import "github.com/stretchr/testify/assert"

type RW struct {
	rd *bytes.Buffer
	wr *bytes.Buffer
}

func NewRW() *RW {
	return &RW{
		rd: new(bytes.Buffer),
		wr: new(bytes.Buffer),
	}
}

func (rw *RW) Write(buf []byte) (int, error) {
	return rw.wr.Write(buf)
}

func (rw *RW) Read(buf []byte) (int, error) {
	return rw.rd.Read(buf)
}

func TestFocusMonitorStart(t *testing.T) {
	buf := NewRW()
	monitor := NewFocusMonitor(socket.NewValueSocket(buf))

	monitor.Start()

	sock := socket.NewValueSocket(buf.wr)
	var sub Subscribe
	sock.Read(&sub)
	assert.EqualValues(t, []string{"window"}, sub)
}

func TestFocusMonitorReceive(t *testing.T) {
	buf := NewRW()
	monitor := NewFocusMonitor(socket.NewValueSocket(buf))

	expectedEvent := Event{
		Change:    "focus",
		Container: Container{12},
	}

	sock := socket.NewValueSocket(buf.rd)
	sock.Write(0, Success{})
	sock.Write(0, expectedEvent)

	ev := <-monitor.Start()
	assert.EqualValues(t, expectedEvent, ev)

}
