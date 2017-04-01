package socket

import (
	"bytes"
	"testing"

	"testing/quick"

	"fmt"

	"io"

	"math"

	"github.com/stretchr/testify/assert"
)

func TestSocketWrite(t *testing.T) {
	err := quick.Check(func(messageType uint32, payload string) bool {
		buf := &bytes.Buffer{}
		socket := NewSocket(buf)

		socket.Write(Message{
			Type:    messageType,
			Payload: payload,
		})

		result := buf.Bytes()

		return assert.Equal(t, result[:6], []byte(MagicString)) &&
			assert.EqualValues(t, nativeOrder.Uint32(result[6:]), len(payload)) &&
			assert.EqualValues(t, nativeOrder.Uint32(result[10:]), messageType) &&
			assert.Equal(t, result[14:], []byte(payload))
	}, nil)
	assert.NoError(t, err)
}

type errNWriter struct {
	io.Reader
	n   int
	err error
}

func (e *errNWriter) Write(buf []byte) (int, error) {
	e.n -= len(buf)
	if e.n < 0 {
		// We only want to return the error once, if the write passes the boundary.
		e.n = math.MaxInt32
		return 0, e.err
	}
	return len(buf), nil
}

func TestSocketWriteError(t *testing.T) {
	err := quick.Check(func(messageType uint32, payload string) bool {
		// The write will always fail somewhere before the end of the message.
		totalMessageSize := len(payload) + len(MagicString) + 8
		for maxBytes := 0; maxBytes < totalMessageSize-1; maxBytes++ {
			err := fmt.Errorf("Cannot write past %d bytes", maxBytes)
			socket := NewSocket(&errNWriter{
				n:   maxBytes,
				err: err,
			})
			resultErr := socket.Write(Message{
				Type:    messageType,
				Payload: payload,
			})
			if !assert.Equal(t, err, resultErr) {
				return false
			}
		}

		return true
	}, nil)
	assert.NoError(t, err)
}

func TestSocketRead(t *testing.T) {
	err := quick.Check(func(message Message) bool {
		buf := &bytes.Buffer{}
		sock := NewSocket(buf)
		err := sock.Write(message)
		assert.NoError(t, err)
		result, err := sock.Read()
		assert.NoError(t, err)
		return assert.Equal(t, message, result)
	}, nil)
	assert.NoError(t, err)
}

type errNReader struct {
	io.Writer
	data []byte
	err  error
}

func (e *errNReader) Read(buf []byte) (int, error) {
	if e.data == nil {
		return len(buf), nil
	}
	rest := len(buf)
	if rest > len(e.data) {
		rest = len(e.data)
	}
	copy(buf, e.data[:rest])
	e.data = e.data[rest:]
	if len(e.data) == 0 {
		// We only want to return the error once, if the read passes the boundary.
		e.data = nil
		return 0, e.err
	}
	return len(buf), nil
}

func TestSocketReadError(t *testing.T) {
	err := quick.Check(func(message Message) bool {

		buf := new(bytes.Buffer)
		NewSocket(buf).Write(message)

		// The write will always fail somewhere before the end of the message.
		totalMessageSize := len(message.Payload) + len(MagicString) + 8
		for maxBytes := 0; maxBytes < totalMessageSize-1; maxBytes++ {
			err := fmt.Errorf("Cannot read past %d bytes", maxBytes)
			socket := NewSocket(&errNReader{
				data: buf.Bytes()[:maxBytes],
				err:  err,
			})
			_, resultErr := socket.Read()
			if !assert.Equal(t, err, resultErr) {
				return false
			}
		}

		return true
	}, nil)
	assert.NoError(t, err)
}
