package socket

import (
	"bytes"
	"testing"
	"testing/quick"

	"errors"

	"github.com/stretchr/testify/assert"
)

type testValue struct {
	Num int
	Str string
	Cat struct {
		Fur  string
		Paws []string
	}
}

func TestValueSocketRoundTrip(t *testing.T) {
	err := quick.Check(func(messageType uint32, v testValue) bool {
		buf := new(bytes.Buffer)
		sock := NewValueSocket(buf)
		err := sock.Write(messageType, v)
		assert.NoError(t, err)
		var result testValue
		resultType, err := sock.Read(&result)
		assert.NoError(t, err)
		assert.Equal(t, messageType, resultType)
		return assert.Equal(t, v, result)
	}, nil)
	assert.NoError(t, err)
}

func TestValueSocketWriteError(t *testing.T) {
	buf := new(bytes.Buffer)
	sock := NewValueSocket(buf)
	err := sock.Write(0, func() {})
	assert.Error(t, err)

	wr := &errNWriter{n: 10, err: errors.New("cannot write")}
	sock = NewValueSocket(wr)
	err = sock.Write(0, []string{"a", "b", "c"})
	assert.Error(t, err)
}

func TestValueSocketReadError(t *testing.T) {
	buf := new(bytes.Buffer)
	rawSocket := NewSocket(buf)
	rawSocket.Write(Message{
		Type:    0,
		Payload: "invalid json",
	})
	sock := NewValueSocket(buf)
	var v interface{}
	_, err := sock.Read(&v)
	assert.Error(t, err)

	expectedError := errors.New("cannot read")
	wr := &errNReader{data: magicStringBytes[:], err: expectedError}
	sock = NewValueSocket(wr)
	_, err = sock.Read(&v)
	assert.EqualValues(t, expectedError, err)
}
